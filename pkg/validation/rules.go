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

type ruleKeys struct {
	// mandatory keys
	Rule    string `json:"rule"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	Message string `json:"message"`
}

type Rule struct {
	ruleKeys
	// optional keys
	Valid string `json:"valid",omitempty`
	// arguments (optional keys)
	Values    []string    `json:"values",omitempty"`
	Min       interface{} `json:"min",omitempty`
	Max       interface{} `json:"max",omitempty`
	MinLength interface{} `json:"minLength",omitempty`
	MaxLength interface{} `json:"maxLength",omitempty`
}

func (r Rule) Apply(vm *k6tv1.VirtualMachine) error {
	return nil
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
