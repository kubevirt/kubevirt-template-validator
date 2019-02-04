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
	"k8s.io/apimachinery/pkg/runtime"
	//	templatev1 "github.com/openshift/api/template/v1"

	"github.com/fromanirh/kubevirt-template-validator/pkg/webhooks"

	"github.com/fromanirh/kubevirt-template-validator/internal/pkg/log"
)

const VMTemplateCreateValidatePath string = "/virtualmachinetemplate-validate-create"
const VMTemplateUpdateValidatePath string = "/virtualmachinetemplate-validate-update"

type admitFunc func(*v1beta1.AdmissionReview) *v1beta1.AdmissionResponse

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

func admitVMTemplateCreate(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	return &reviewResponse
}

func ServeVMTemplateCreate(resp http.ResponseWriter, req *http.Request) {
	serve(resp, req, admitVMTemplateCreate)
}

func admitVMTemplateUpdate(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	return &reviewResponse
}

func ServeVMTemplateUpdate(resp http.ResponseWriter, req *http.Request) {
	serve(resp, req, admitVMTemplateUpdate)
}
