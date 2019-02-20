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
	"regexp"
	"strings"

	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

type RuleApplier interface {
	Apply(vm *k6tv1.VirtualMachine) (bool, error)
	String() string
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
	case "regex":
		return NewRegexRule(r, vm)
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

// can't be a Stringer
func (r *Range) ToString(v int64) string {
	cond := ""
	if !r.Includes(v) {
		cond = "not "
	}
	lowerBound := "N/A"
	if r.MinSet {
		lowerBound = fmt.Sprintf("%d", r.Min)
	}
	upperBound := "N/A"
	if r.MaxSet {
		upperBound = fmt.Sprintf("%d", r.Max)
	}
	return fmt.Sprintf("%v %sin [%s, %s]", v, cond, lowerBound, upperBound)
}

// These are the specializedrules
type intRule struct {
	Ref       *Rule
	Value     Range
	Current   int64
	Satisfied bool
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
		p, err := NewPath(minStr)
		if err != nil {
			return 0, err
		}
		err = p.Find(vm)
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
	p, err := NewPath(s)
	if err != nil {
		return "", err
	}
	err = p.Find(vm)
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
	var err error
	ir.Current, err = decodeInt64(ir.Ref.Path, vm)
	if err != nil {
		return false, err
	}
	ir.Satisfied = ir.Value.Includes(ir.Current)
	return ir.Satisfied, nil
}

func (ir *intRule) String() string {
	return ir.Value.ToString(ir.Current)
}

type stringRule struct {
	Ref       *Rule
	Length    Range
	Current   string
	Satisfied bool
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
	var err error
	sr.Current, err = decodeString(sr.Ref.Path, vm)
	if err != nil {
		return false, err
	}
	sr.Satisfied = sr.Length.Includes(int64(len(sr.Current)))
	return sr.Satisfied, nil
}

func (sr *stringRule) String() string {
	return sr.Length.ToString(int64(len(sr.Current)))
}

type enumRule struct {
	Ref       *Rule
	Values    []string
	Current   string
	Satisfied bool
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
	var err error
	er.Current, err = decodeString(er.Ref.Path, vm)
	if err != nil {
		return false, err
	}
	for _, v := range er.Values {
		if er.Current == v {
			er.Satisfied = true
			return er.Satisfied, nil
		}
	}
	er.Satisfied = false // enforce
	return er.Satisfied, nil
}

func (er *enumRule) String() string {
	cond := ""
	if !er.Satisfied {
		cond = "not "
	}
	return fmt.Sprintf("%s %sin [%s]", er.Current, cond, strings.Join(er.Values, ", "))
}

type regexRule struct {
	Ref       *Rule
	Regex     string
	Current   string
	Satisfied bool
}

func NewRegexRule(r *Rule, vm *k6tv1.VirtualMachine) (RuleApplier, error) {
	return &regexRule{
		Ref:   r,
		Regex: r.Regex,
	}, nil
}

func (rr *regexRule) Apply(vm *k6tv1.VirtualMachine) (bool, error) {
	var err error
	rr.Current, err = decodeString(rr.Ref.Path, vm)
	if err != nil {
		return false, err
	}
	rr.Satisfied, err = regexp.MatchString(rr.Regex, rr.Current)
	return rr.Satisfied, err
}

func (rr *regexRule) String() string {
	cond := ""
	if !rr.Satisfied {
		cond = "not "
	}
	return fmt.Sprintf("%s %smatches %s", rr.Current, cond, rr.Regex)
}
