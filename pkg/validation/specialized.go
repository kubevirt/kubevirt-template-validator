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
	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

type RuleApplier interface {
	GetRef() *Rule
	Apply(vm *k6tv1.VirtualMachine) bool
}

// we need a vm reference to specialize a rule because few key fields may
// be JSONPath, and we need to walk them to get e.g. the value to check,
// or the limits to enforce.
func (r *Rule) Specialize(vm *k6tv1.VirtualMachine) (RuleApplier, error) {
	return nil, nil
}

// These are the specializedrules
type intRule struct {
	Ref    *Rule
	MinSet bool
	Min    int
	MaxSet bool
	Max    int
}

func (ir *intRule) GetRef() *Rule {
	return ir.Ref
}

func (ir *intRule) Apply(vm *k6tv1.VirtualMachine) bool {
	return false
}

type stringRule struct {
	Ref          *Rule
	MinLengthSet bool
	MinLength    int
	MaxLengthSet bool
	MaxLength    int
}

func (sr *stringRule) GetRef() *Rule {
	return sr.Ref
}

func (sr *stringRule) Apply(vm *k6tv1.VirtualMachine) bool {
	return false
}

type enumRule struct {
	Ref    *Rule
	Values []string
}

func (er *enumRule) GetRef() *Rule {
	return er.Ref
}

func (er *enumRule) Apply(vm *k6tv1.VirtualMachine) bool {
	return false
}
