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
	"encoding/json"

	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

// the flow:
// 1. first you do ParseRule and get []Rule. This is little more than raw data rearranged in Go structs.
//    You can work with that programmatically, but the first thing you may want to do is
// 2. ...

type Rule struct {
	// mandatory keys
	Rule    string `json:"rule"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	Message string `json:"message"`
	// optional keys
	Valid       string `json:"valid",omitempty`
	JustWarning bool   `json:"justWarning",omitempty"`
	// arguments (optional keys)
	Values    []string    `json:"values",omitempty"`
	Min       interface{} `json:"min",omitempty`
	Max       interface{} `json:"max",omitempty`
	MinLength interface{} `json:"minLength",omitempty`
	MaxLength interface{} `json:"maxLength",omitempty`
	Regex     string      `json:"regex",omitempty"`
}

func (r *Rule) findPathOn(vm *k6tv1.VirtualMachine) (bool, error) {
	var err error
	p, err := NewPath(r.Valid)
	if err != nil {
		return false, err
	}
	err = p.Find(vm)
	if err != nil {
		return false, err
	}
	return p.Len() > 0, nil
}

func (r *Rule) IsAppliableOn(vm *k6tv1.VirtualMachine) (bool, error) {
	if r.Valid == "" {
		// nothing to check against, so it is OK
		return true, nil
	}
	ok, err := r.findPathOn(vm)
	if err == ErrInvalidJSONPath {
		return false, nil
	}
	return ok, err
}

func ParseRules(data []byte) ([]Rule, error) {
	var rules []Rule
	if len(data) == 0 {
		// nothing to do
		return rules, nil
	}
	err := json.Unmarshal(data, &rules)
	return rules, err
}
