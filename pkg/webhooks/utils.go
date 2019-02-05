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

package webhooks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"

	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
	"kubevirt.io/kubevirt/pkg/kubecli"
	"kubevirt.io/kubevirt/pkg/util/openapi"
	"kubevirt.io/kubevirt/pkg/virt-api/rest"

	"github.com/fromanirh/kubevirt-template-validator/pkg/virtinformers"

	"github.com/fromanirh/kubevirt-template-validator/internal/pkg/k8sutils"
	"github.com/fromanirh/kubevirt-template-validator/internal/pkg/log"
)

var Validator = openapi.CreateOpenAPIValidator(rest.ComposeAPIDefinitions())

var webhookInformers *Informers
var once sync.Once

type Informers struct {
	TemplateInformer       cache.SharedIndexInformer
	VirtualMachineInformer cache.SharedIndexInformer
}

func GetInformers() *Informers {
	once.Do(func() {
		webhookInformers = newInformers()
	})
	return webhookInformers
}

// SetInformers created for unittest usage only
func SetInformers(informers *Informers) {
	once.Do(func() {
		webhookInformers = informers
	})
}

func newInformers() *Informers {
	kubeClient, err := kubecli.GetKubevirtClient()
	if err != nil {
		panic(err)
	}
	namespace, err := k8sutils.GetNamespace()
	if err != nil {
		log.Log.Errorf("Error searching for namespace: %v", err)
		return nil
	}
	kubeInformerFactory := virtinformers.NewKubeInformerFactory(kubeClient.RestClient(), namespace)
	return &Informers{
		TemplateInformer: kubeInformerFactory.Template(),
	}
}

// GetAdmissionReview
func GetAdmissionReview(r *http.Request) (*v1beta1.AdmissionReview, error) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return nil, fmt.Errorf("contentType=%s, expect application/json", contentType)
	}

	ar := &v1beta1.AdmissionReview{}
	err := json.Unmarshal(body, ar)
	return ar, err
}

func ToAdmissionResponseOK() *v1beta1.AdmissionResponse {
	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	return &reviewResponse
}

// ToAdmissionResponseError
func ToAdmissionResponseError(err error) *v1beta1.AdmissionResponse {
	log.Log.Errorf("admission generic error: %v", err)

	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		},
	}
}

func ToAdmissionResponse(causes []metav1.StatusCause) *v1beta1.AdmissionResponse {
	log.Log.Infof("rejected vmi admission")

	globalMessage := ""
	for _, cause := range causes {
		if globalMessage == "" {
			globalMessage = cause.Message
		} else {
			globalMessage = fmt.Sprintf("%s, %s", globalMessage, cause.Message)
		}
	}

	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: globalMessage,
			Reason:  metav1.StatusReasonInvalid,
			Code:    http.StatusUnprocessableEntity,
			Details: &metav1.StatusDetails{
				Causes: causes,
			},
		},
	}
}

func ValidationErrorsToAdmissionResponse(errs []error) *v1beta1.AdmissionResponse {
	var causes []metav1.StatusCause
	for _, e := range errs {
		causes = append(causes,
			metav1.StatusCause{
				Message: e.Error(),
			},
		)
	}
	return ToAdmissionResponse(causes)
}

func ValidateSchema(gvk schema.GroupVersionKind, data []byte) *v1beta1.AdmissionResponse {
	in := map[string]interface{}{}
	err := json.Unmarshal(data, &in)
	if err != nil {
		return ToAdmissionResponseError(err)
	}
	errs := Validator.Validate(gvk, in)
	if len(errs) > 0 {
		return ValidationErrorsToAdmissionResponse(errs)
	}
	return nil
}

func GetAdmissionReviewVM(ar *v1beta1.AdmissionReview) (*k6tv1.VirtualMachine, *k6tv1.VirtualMachine, error) {

	if ar.Request.Resource.Resource != "template" {
		return nil, nil, fmt.Errorf("expect resource %v to be '%s'", ar.Request.Resource, "template")
	}

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
