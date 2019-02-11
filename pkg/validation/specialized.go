/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2019 Red Hat, Inc.
 */

package validation

import (
	"fmt"
	"strings"

	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

type RuleApplier interface {
	Apply(vm *k6tv1.VirtualMachine) (bool, error)
}

func isJSONPath(s string) bool {
	return strings.HasPrefix(s, "$")
}

// we need a vm reference to specialize a rule because few key fields may
// be JSONPath, and we need to walk them to get e.g. the value to check,
// or the limits to enforce.
func (r *Rule) Specialize(vm *k6tv1.VirtualMachine) (RuleApplier, error) {
	switch r.Rule {
	case "integer":
		return NewIntRule(r, vm)
	case "string":
		return NewStringRule(r, vm)
	case "enum":
		return NewEnumRule(r, vm)
	}
	return nil, fmt.Errorf("Usupported rule: %s", r.Rule)
}

type Range struct {
	MinSet bool
	Min    int64
	MaxSet bool
	Max    int64
}

func (r *Range) Decode(Min, Max interface{}, vm *k6tv1.VirtualMachine) error {
	if Min != nil {
		v, err := decodeInt64(Min, vm)
		if err != nil {
			return err
		}
		r.Min = v
		r.MinSet = true
	}
	if Max != nil {
		v, err := decodeInt64(Max, vm)

		if err != nil {
			return err
		}
		r.Max = v
		r.MaxSet = true
	}
	return nil
}

func (r *Range) Includes(v int64) bool {
	if r.MinSet && v < r.Min {
		return false
	}
	if r.MaxSet && v > r.Max {
		return false
	}
	return true
}

// These are the specializedrules
type intRule struct {
	Ref   *Rule
	Value Range
}

func decodeInt64(obj interface{}, vm *k6tv1.VirtualMachine) (int64, error) {
	if minVal, ok := obj.(int); ok {
		return int64(minVal), nil
	}
	if minVal, ok := obj.(int32); ok {
		return int64(minVal), nil
	}
	if minVal, ok := obj.(int64); ok {
		return int64(minVal), nil
	}
	if minStr, ok := obj.(string); ok && isJSONPath(minStr) {
		p, err := Find(vm, minStr)
		if err != nil {
			return 0, err
		}
		if p.Len() != 1 {
			return 0, fmt.Errorf("expected one value, found %v", p.Len())
		}
		vals, err := p.AsInt64()
		if err != nil {
			return 0, err
		}
		return vals[0], nil
	}
	return 0, fmt.Errorf("Unrecognized type")
}

func decodeString(s string, vm *k6tv1.VirtualMachine) (string, error) {
	if !isJSONPath(s) {
		return s, nil
	}
	p, err := Find(vm, s)
	if err != nil {
		return "", err
	}
	if p.Len() != 1 {
		return "", fmt.Errorf("expected one value, found %v", p.Len())
	}
	vals, err := p.AsString()
	if err != nil {
		return "", err
	}
	return vals[0], nil
}

func NewIntRule(r *Rule, vm *k6tv1.VirtualMachine) (RuleApplier, error) {
	ir := intRule{Ref: r}
	err := ir.Value.Decode(r.Min, r.Max, vm)
	if err != nil {
		return nil, err
	}
	return &ir, nil
}

func (ir *intRule) Apply(vm *k6tv1.VirtualMachine) (bool, error) {
	v, err := decodeInt64(ir.Ref.Path, vm)
	if err != nil {
		return false, err
	}
	return ir.Value.Includes(v), nil
}

type stringRule struct {
	Ref    *Rule
	Length Range
}

func NewStringRule(r *Rule, vm *k6tv1.VirtualMachine) (RuleApplier, error) {
	sr := stringRule{Ref: r}
	err := sr.Length.Decode(r.MinLength, r.MaxLength, vm)
	if err != nil {
		return nil, err
	}
	return &sr, nil
}

func (sr *stringRule) Apply(vm *k6tv1.VirtualMachine) (bool, error) {
	s, err := decodeString(sr.Ref.Path, vm)
	if err != nil {
		return false, err
	}
	return sr.Length.Includes(int64(len(s))), nil
}

type enumRule struct {
	Ref    *Rule
	Values []string
}

func NewEnumRule(r *Rule, vm *k6tv1.VirtualMachine) (RuleApplier, error) {
	er := enumRule{Ref: r}
	for _, v := range r.Values {
		s, err := decodeString(v, vm)
		if err != nil {
			return nil, err
		}
		er.Values = append(er.Values, s)
	}
	return &er, nil
}

func (er *enumRule) Apply(vm *k6tv1.VirtualMachine) (bool, error) {
	s, err := decodeString(er.Ref.Path, vm)
	if err != nil {
		return false, err
	}
	for _, v := range er.Values {
		if s == v {
			return true, nil
		}
	}
	return false, nil
}

// TODO: regex
