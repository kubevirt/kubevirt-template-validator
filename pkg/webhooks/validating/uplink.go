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
	"fmt"

	templatev1 "github.com/openshift/api/template/v1"

	k6tv1 "kubevirt.io/client-go/api/v1"

	"github.com/fromanirh/kubevirt-template-validator/pkg/validation"
	"github.com/fromanirh/kubevirt-template-validator/pkg/virtinformers"

	"github.com/fromanirh/kubevirt-template-validator/internal/pkg/log"
)

const (
	annotationTemplateNameKey         string = "vm.kubevirt.io/template"
	annotationTemplateNamespaceKey    string = "vm.kubevirt.io/template.namespace"
	annotationTemplateNamespaceOldKey string = "vm.kubevirt.io/template-namespace"
	annotationValidationKey           string = "validations"
)

func getTemplateKey(vm *k6tv1.VirtualMachine) (string, bool) {
	if vm.Annotations == nil {
		log.Log.V(4).Infof("VM %s missing annotations entirely", vm.Name)
		return "", false
	}

	templateNamespace := vm.Annotations[annotationTemplateNamespaceKey]
	if templateNamespace == "" {
		templateNamespace = vm.Annotations[annotationTemplateNamespaceOldKey]
		if templateNamespace != "" {
			log.Log.V(5).Warningf("VM %s has old-style template namespace annotation '%s', should be updated to '%s'", vm.Name, annotationTemplateNamespaceOldKey, annotationTemplateNamespaceKey)
		}
	}

	if templateNamespace == "" {
		log.Log.V(4).Infof("VM %s missing template namespace annotation", vm.Name)
		return "", false
	}

	templateName := vm.Annotations[annotationTemplateNameKey]
	if templateName == "" {
		log.Log.V(4).Infof("VM %s missing template annotation", vm.Name)
		return "", false
	}

	return fmt.Sprintf("%s/%s", templateNamespace, templateName), true
}

func getParentTemplateForVM(vm *k6tv1.VirtualMachine) (*templatev1.Template, error) {
	informers := virtinformers.GetInformers()

	if !informers.Available() {
		log.Log.V(8).Infof("no informer available (deployed on K8S?)")
		return nil, nil
	}

	cacheKey, ok := getTemplateKey(vm)
	if !ok {
		log.Log.V(8).Infof("detected %s as baked (no parent template)", vm.Name)
		return nil, nil
	}

	obj, exists, err := informers.TemplateInformer.GetStore().GetByKey(cacheKey)
	if err != nil {
		log.Log.V(8).Infof("parent template (key=%s) not found for %s: %v", cacheKey, vm.Name, err)
		return nil, err
	}

	if !exists {
		msg := fmt.Sprintf("missing parent template (key=%s) for %s", cacheKey, vm.Name)
		log.Log.V(4).Warning(msg)
		return nil, fmt.Errorf("%s", msg)
	}

	log.Log.V(8).Infof("found parent template for %s", vm.Name)
	tmpl := obj.(*templatev1.Template)
	// TODO explain deepcopy
	return tmpl.DeepCopy(), nil
}

func getValidationRulesFromTemplate(tmpl *templatev1.Template) ([]validation.Rule, error) {
	return validation.ParseRules([]byte(tmpl.Annotations[annotationValidationKey]))
}

func getValidationRulesForVM(vm *k6tv1.VirtualMachine) ([]validation.Rule, error) {
	tmpl, err := getParentTemplateForVM(vm)
	if tmpl == nil || err != nil {
		// no template resources (kubevirt deployed on kubernetes, not OKD/OCP) or
		// no parent template for this VM. In either case, we have nothing to do,
		// and err is automatically correct
		return []validation.Rule{}, err
	}
	return getValidationRulesFromTemplate(tmpl)
}
