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

	"github.com/fromanirh/kubevirt-template-validator/pkg/webhooks"

	"github.com/fromanirh/kubevirt-template-validator/internal/pkg/log"
)

const VMTemplateCreateValidatePath string = "/virtualmachine-template-validate-create"
const VMTemplateUpdateValidatePath string = "/virtualmachine-template-validate-update"

type admitFunc func(*v1beta1.AdmissionReview) *v1beta1.AdmissionResponse

func validateVirtualMachineFromTemplate(field *k8sfield.Path, newVM *k6tv1.VirtualMachine, oldVM *k6tv1.VirtualMachine, tmpl *templatev1.Template) []metav1.StatusCause {
	var causes []metav1.StatusCause
	return causes
}

func admitVMTemplate(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	newVM, oldVM, err := getAdmissionReviewVM(ar)
	if err != nil {
		return webhooks.ToAdmissionResponseError(err)
	}

	informers := webhooks.GetInformers()
	cacheKey := "" // fmt.Sprintf("%s/%s", migration.Namespace, migration.Spec.VMIName)
	obj, exists, err := informers.VirtualMachineInformer.GetStore().GetByKey(cacheKey)
	if err != nil {
		return webhooks.ToAdmissionResponseError(err)
	}

	if !exists {
		return webhooks.ToAdmissionResponseError(fmt.Errorf("the VMI %s does not exist under the cache", ""))
	}
	tmplObj := obj.(*templatev1.Template)
	tmpl := tmplObj.DeepCopy()

	causes := validateVirtualMachineFromTemplate(nil, newVM, oldVM, tmpl)
	if len(causes) > 0 {
		return webhooks.ToAdmissionResponse(causes)
	}

	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	return &reviewResponse
}

func ServeVMTemplateCreate(resp http.ResponseWriter, req *http.Request) {
	serve(resp, req, admitVMTemplate)
}

func ServeVMTemplateUpdate(resp http.ResponseWriter, req *http.Request) {
	serve(resp, req, admitVMTemplate)
}

func serve(resp http.ResponseWriter, req *http.Request, admit admitFunc) {
	response := v1beta1.AdmissionReview{}
	review, err := webhooks.GetAdmissionReview(req)

	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	reviewResponse := admit(review)
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

func getAdmissionReviewVM(ar *v1beta1.AdmissionReview) (*k6tv1.VirtualMachine, *k6tv1.VirtualMachine, error) {
	var err error
	raw := ar.Request.Object.Raw
	newVM := k6tv1.VirtualMachine{}

	err = json.Unmarshal(raw, &newVM)
	if err != nil {
		return nil, nil, err
	}

	if ar.Request.Operation == v1beta1.Update {
		raw := ar.Request.OldObject.Raw
		oldVM := k6tv1.VirtualMachine{}
		err = json.Unmarshal(raw, &oldVM)
		if err != nil {
			return nil, nil, err
		}
		return &newVM, &oldVM, nil
	}

	return &newVM, nil, nil
}
