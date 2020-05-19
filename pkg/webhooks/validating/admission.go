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
	"bytes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k6tv1 "kubevirt.io/client-go/api/v1"
	"kubevirt.io/client-go/log"

	"github.com/fromanirh/kubevirt-template-validator/pkg/validation"
)

func ValidateVMTemplate(rules []validation.Rule, newVM, oldVM *k6tv1.VirtualMachine) []metav1.StatusCause {
	var causes []metav1.StatusCause
	if len(rules) == 0 {
		// no rules! everything is permitted, so let's bail out quickly
		log.Log.V(8).Infof("no admission rules for: %s", newVM.Name)
		return causes
	}
	buf := new(bytes.Buffer)
	ev := validation.Evaluator{Sink: buf}
	res := ev.Evaluate(rules, newVM)
	log.Log.V(2).Infof("evalution summary for %s:\n%s\nsucceeded=%v", newVM.Name, buf.String(), res.Succeeded())

	if res.Succeeded() {
		return causes
	}
	return res.ToStatusCauses()
}
