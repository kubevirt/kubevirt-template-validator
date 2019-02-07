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

package validating

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfield "k8s.io/apimachinery/pkg/util/validation/field"

	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"

	"github.com/fromanirh/kubevirt-template-validator/pkg/validation"
	//	"github.com/fromanirh/kubevirt-template-validator/internal/pkg/log"
)

func ValidateVirtualMachine(field *k8sfield.Path, newVM *k6tv1.VirtualMachine, oldVM *k6tv1.VirtualMachine, rules []validation.Rule) []metav1.StatusCause {
	var causes []metav1.StatusCause
	if len(rules) == 0 {
		// no rules! everything is permitted, so let's bail out quickly
		return causes
	}
	return causes
}
