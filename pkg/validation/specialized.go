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
	"reflect"
	"regexp"
	"strings"

	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

type RuleApplier interface {
	Apply(vm, ref *k6tv1.VirtualMachine) (bool, error)
	String() string
}

// we need a vm reference to specialize a rule because few key fields may
// be JSONPath, and we need to walk them to get e.g. the value to check,
// or the limits to enforce.
func (r *Rule) Specialize(vm, ref *k6tv1.VirtualMachine) (RuleApplier, error) {
	switch r.Rule {
	case "integer":
		return NewIntRule(r, vm, ref)
	case "string":
		return NewStringRule(r, vm, ref)
	case "enum":
		return NewEnumRule(r, vm, ref)
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

func (r *Range) Decode(Min, Max interface{}, vm, ref *k6tv1.VirtualMachine) error {
	if Min != nil {
		v, err := decodeInt64(Min, vm, ref)
		if err != nil {
			return err
		}
		r.Min = v
		r.MinSet = true
	}
	if Max != nil {
		v, err := decodeInt64(Max, vm, ref)

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

// JSONPATH lookup logic, aka what this "ref" object and why we need it
//
// When we need to fetch the value of a Rule.Path which happens to be a JSONPath,
// first we just try if the given VM object has the path we need.
// If the lookup succeeds, everyone's happy and we stop here.
// Else, the vm obj has not the path we were looking for.
// It could be either:
// - the path is bogus. We check lazily, so this is the first time we see this
//   and we need to make a decision. But mayne
// - the path is legal, but it refers to an optional subpath which is missing.
//   so we try again with the zero-initialized "reference" object.
//   if even this lookup fails, we mark the path as bogus.
//   Otherwise we use the zero, default, value for our logic.

func decodeInt64(obj interface{}, vm, ref *k6tv1.VirtualMachine) (int64, error) {
	if val, ok := toInt64(obj); ok {
		return val, nil
	}
	v, err := decodeInt64FromJSONPath(obj, vm)
	if err != nil {
		v, err = decodeInt64FromJSONPath(obj, ref)
	}
	return v, err
}

func decodeString(s string, vm, ref *k6tv1.VirtualMachine) (string, error) {
	if !isJSONPath(s) {
		return s, nil
	}
	v, err := decodeJSONPathString(s, vm)
	if err != nil {
		v, err = decodeJSONPathString(s, ref)
	}
	return v, err
}

func decodeInt64FromJSONPath(obj interface{}, vm *k6tv1.VirtualMachine) (int64, error) {
	if strVal, ok := obj.(string); ok && isJSONPath(strVal) {
		p, err := NewPath(strVal)
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
	return 0, fmt.Errorf("Unsupported type %v (%v)", obj, reflect.TypeOf(obj).Name())
}

func decodeJSONPathString(s string, vm *k6tv1.VirtualMachine) (string, error) {
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

func NewIntRule(r *Rule, vm, ref *k6tv1.VirtualMachine) (RuleApplier, error) {
	ir := intRule{Ref: r}
	err := ir.Value.Decode(r.Min, r.Max, vm, ref)
	if err != nil {
		return nil, err
	}
	return &ir, nil
}

func (ir *intRule) Apply(vm, ref *k6tv1.VirtualMachine) (bool, error) {
	v, err := decodeInt64(ir.Ref.Path, vm, ref)
	if err != nil {
		return false, err
	}
	ir.Current = v
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

func NewStringRule(r *Rule, vm, ref *k6tv1.VirtualMachine) (RuleApplier, error) {
	sr := stringRule{Ref: r}
	err := sr.Length.Decode(r.MinLength, r.MaxLength, vm, ref)
	if err != nil {
		return nil, err
	}
	return &sr, nil
}

func (sr *stringRule) Apply(vm, ref *k6tv1.VirtualMachine) (bool, error) {
	v, err := decodeString(sr.Ref.Path, vm, ref)
	if err != nil {
		return false, err
	}
	sr.Current = v
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

func NewEnumRule(r *Rule, vm, ref *k6tv1.VirtualMachine) (RuleApplier, error) {
	er := enumRule{Ref: r}
	for _, v := range r.Values {
		s, err := decodeString(v, vm, ref)
		if err != nil {
			return nil, err
		}
		er.Values = append(er.Values, s)
	}
	return &er, nil
}

func (er *enumRule) Apply(vm, ref *k6tv1.VirtualMachine) (bool, error) {
	v, err := decodeString(er.Ref.Path, vm, ref)
	if err != nil {
		return false, err
	}
	er.Current = v
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

func (rr *regexRule) Apply(vm, ref *k6tv1.VirtualMachine) (bool, error) {
	v, err := decodeString(rr.Ref.Path, vm, ref)
	if err != nil {
		return false, err
	}
	rr.Current = v
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
