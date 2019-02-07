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
	"encoding/json"
	"fmt"
	"net/http"

	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfield "k8s.io/apimachinery/pkg/util/validation/field"

	templatev1 "github.com/openshift/api/template/v1"

	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"

	"github.com/davecgh/go-spew/spew"

	"github.com/fromanirh/kubevirt-template-validator/pkg/virtinformers"
	"github.com/fromanirh/kubevirt-template-validator/pkg/webhooks"

	"github.com/fromanirh/kubevirt-template-validator/internal/pkg/log"
)

const (
	VMTemplateValidatePath string = "/virtualmachine-template-validate"

	annotationTemplateNameKey      string = "vm.cnv.io/template"
	annotationTemplateNamespaceKey string = "vm.cnv.io/template-namespace"
)

type admitFunc func(*v1beta1.AdmissionReview) *v1beta1.AdmissionResponse

func validateVirtualMachineFromTemplate(field *k8sfield.Path, newVM *k6tv1.VirtualMachine, oldVM *k6tv1.VirtualMachine, tmpl *templatev1.Template) []metav1.StatusCause {
	var causes []metav1.StatusCause
	return causes
}

func getTemplateKey(vm *k6tv1.VirtualMachine) (string, bool) {
	if vm.Annotations == nil {
		log.Log.Warningf("VM %s missing annotations entirely", vm.Name)
		return "", false
	}

	templateNamespace := vm.Annotations[annotationTemplateNamespaceKey]
	if templateNamespace == "" {
		log.Log.Warningf("VM %s missing template namespace annotation", vm.Name)
		return "", false
	}

	templateName := vm.Annotations[annotationTemplateNameKey]
	if templateNamespace == "" {
		log.Log.Warningf("VM %s missing template annotation", vm.Name)
		return "", false
	}

	return fmt.Sprintf("%s/%s", templateNamespace, templateName), true
}

func getParentTemplateForVM(vm *k6tv1.VirtualMachine) (*templatev1.Template, error) {
	informers := virtinformers.GetInformers()

	if informers == nil || informers.TemplateInformer == nil {
		// no error, it can happen: we're been deployed ok K8S, not OKD/OCD.
		return nil, nil
	}

	cacheKey, ok := getTemplateKey(vm)
	if !ok {
		// baked VM (aka no parent template). May happen, it's OK.
		return nil, nil
	}

	obj, exists, err := informers.TemplateInformer.GetStore().GetByKey(cacheKey)
	if err != nil {
		return nil, err
	}

	if !exists {
		// ok, this is weird
		return nil, fmt.Errorf("unable to find template object %s for VM %s", cacheKey, vm.Name)
	}

	tmpl := obj.(*templatev1.Template)
	// TODO explain deepcopy
	return tmpl.DeepCopy(), nil
}

func admitVMTemplate(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	newVM, oldVM, err := webhooks.GetAdmissionReviewVM(ar)
	if err != nil {
		return webhooks.ToAdmissionResponseError(err)
	}

	if resp := webhooks.ValidateSchema(k6tv1.VirtualMachineGroupVersionKind, ar.Request.Object.Raw); resp != nil {
		return resp
	}

	templ, err := getParentTemplateForVM(newVM)
	if err != nil {
		return webhooks.ToAdmissionResponseError(err)
	}
	if templ == nil {
		// no template resources (kubevirt deployed on kubernetes, not OKD/OCP) or
		// no parent template for this VM. In either case, we have nothing to do.
		return webhooks.ToAdmissionResponseOK()
	}

	if IsDumpModeEnabled() {
		log.Log.Infof("admission newVM:\n%s", spew.Sdump(newVM))
		log.Log.Infof("admission oldVM:\n%s", spew.Sdump(oldVM))
		log.Log.Infof("admission Templ:\n%s", spew.Sdump(templ))
	}

	causes := validateVirtualMachineFromTemplate(nil, newVM, oldVM, templ)
	if len(causes) > 0 {
		return webhooks.ToAdmissionResponse(causes)
	}

	return webhooks.ToAdmissionResponseOK()
}

func ServeVMTemplateValidate(resp http.ResponseWriter, req *http.Request) {
	serve(resp, req, admitVMTemplate)
}

func serve(resp http.ResponseWriter, req *http.Request, admit admitFunc) {
	response := v1beta1.AdmissionReview{}
	review, err := webhooks.GetAdmissionReview(req)

	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	if IsDumpModeEnabled() {
		log.Log.Infof("admission review:\n%s", spew.Sdump(review))
	}

	reviewResponse := admit(review)

	if IsDumpModeEnabled() {
		log.Log.Infof("admission review response:\n%s", spew.Sdump(reviewResponse))
	}

	if reviewResponse != nil {
		response.Response = reviewResponse
		response.Response.UID = review.Request.UID
	}
	// reset the Object and OldObject, they are not needed in a response.
	review.Request.Object = runtime.RawExtension{}
	review.Request.OldObject = runtime.RawExtension{}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Log.Errorf("failed json encode webhook response: %v", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := resp.Write(responseBytes); err != nil {
		log.Log.Errorf("failed to write webhook response: %v", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	resp.WriteHeader(http.StatusOK)
}
