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
	"net/http"

	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfield "k8s.io/apimachinery/pkg/util/validation/field"

	templatev1 "github.com/openshift/api/template/v1"

	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"

	"github.com/davecgh/go-spew/spew"

	"github.com/fromanirh/kubevirt-template-validator/pkg/webhooks"

	"github.com/fromanirh/kubevirt-template-validator/internal/pkg/log"
)

const VMTemplateValidatePath string = "/virtualmachine-template-validate"

type admitFunc func(*v1beta1.AdmissionReview) *v1beta1.AdmissionResponse

func validateVirtualMachineFromTemplate(field *k8sfield.Path, newVM *k6tv1.VirtualMachine, oldVM *k6tv1.VirtualMachine, tmpl *templatev1.Template) []metav1.StatusCause {
	var causes []metav1.StatusCause
	return causes
}

func admitVMTemplate(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	newVM, oldVM, err := webhooks.GetAdmissionReviewVM(ar)
	if err != nil {
		return webhooks.ToAdmissionResponseError(err)
	}

	if resp := webhooks.ValidateSchema(k6tv1.VirtualMachineGroupVersionKind, ar.Request.Object.Raw); resp != nil {
		return resp
	}

	//	informers := webhooks.GetInformers()
	//	cacheKey := "" // fmt.Sprintf("%s/%s", migration.Namespace, migration.Spec.VMIName)
	//	obj, exists, err := informers.VirtualMachineInformer.GetStore().GetByKey(cacheKey)
	//	if err != nil {
	//		return webhooks.ToAdmissionResponseError(err)
	//	}

	//	if !exists {
	//		// VM doesn't originate from a template. Totally fine and expected.
	//		return webhooks.ToAdmissionResponseOK()
	//	}
	//	tmplObj := obj.(*templatev1.Template)
	//	tmpl := tmplObj.DeepCopy()

	if IsDumpModeEnabled() {
		log.Log.Infof("admission newVM:\n%s", spew.Sdump(newVM))
		log.Log.Infof("admission oldVM:\n%s", spew.Sdump(oldVM))
		//		log.Log.Infof("admission tmpl:\n%s", spew.Sdump(tmpl))
	}

	causes := validateVirtualMachineFromTemplate(nil, newVM, oldVM, nil) //	tmpl)
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
