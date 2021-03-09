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

	"github.com/davecgh/go-spew/spew"

	"kubevirt.io/client-go/log"

	"github.com/kubevirt/kubevirt-template-validator/pkg/webhooks"
)

const (
	VMTemplateValidatePath string = "/virtualmachine-template-validate"
)

func ServeVMTemplateValidate(resp http.ResponseWriter, req *http.Request) {
	serve(resp, req, admitVMTemplate)
}

type admitFunc func(*v1beta1.AdmissionReview) *v1beta1.AdmissionResponse

func admitVMTemplate(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	newVM, oldVM, err := webhooks.GetAdmissionReviewVM(ar)
	if err != nil {
		return webhooks.ToAdmissionResponseError(err)
	}

	if newVM.DeletionTimestamp != nil {
		return webhooks.ToAdmissionResponseOK()
	}

	rules, err := getValidationRulesForVM(newVM)
	if err != nil {
		return webhooks.ToAdmissionResponseError(err)
	}

	log.Log.V(8).Infof("admission newVM:\n%s", spew.Sdump(newVM))
	log.Log.V(8).Infof("admission oldVM:\n%s", spew.Sdump(oldVM))
	log.Log.V(8).Infof("admission rules:\n%s", spew.Sdump(rules))

	causes := ValidateVMTemplate(rules, newVM, oldVM)
	if len(causes) > 0 {
		return webhooks.ToAdmissionResponse(causes)
	}

	return webhooks.ToAdmissionResponseOK()
}

func serve(resp http.ResponseWriter, req *http.Request, admit admitFunc) {
	response := v1beta1.AdmissionReview{}
	review, err := webhooks.GetAdmissionReview(req)

	log.Log.V(8).Infof("evaluating admission")
	defer log.Log.V(8).Infof("evaluated admission")

	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Log.V(8).Infof("admission review:\n%s", spew.Sdump(review))

	reviewResponse := admit(review)

	log.Log.V(8).Infof("admission review response:\n%s", spew.Sdump(reviewResponse))

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
}
