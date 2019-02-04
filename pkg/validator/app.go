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

package validator

import (
	"net/http"

	"k8s.io/client-go/tools/cache"

	"github.com/fromanirh/kubevirt-template-validator/internal/pkg/k8sutils"
	"github.com/fromanirh/kubevirt-template-validator/internal/pkg/log"

	"github.com/fromanirh/kubevirt-template-validator/pkg/webhooks"
	"github.com/fromanirh/kubevirt-template-validator/pkg/webhooks/validating"
)

type App struct {
	ListenAddress string
	TLSInfo       *k8sutils.TLSInfo
}

func (app *App) Run() error {
	if app.TLSInfo == nil {
		app.TLSInfo = &k8sutils.TLSInfo{}
	}

	log.Log.Infof("webhook App: running with TLSInfo %#v", app.TLSInfo)
	// Run informers for webhooks usage
	webhookInformers := webhooks.GetInformers()

	stopChan := make(chan struct{}, 1)
	defer close(stopChan)
	go webhookInformers.TemplateInformer.Run(stopChan)

	log.Log.Infof("webhook App: started informers")

	cache.WaitForCacheSync(
		stopChan,
		webhookInformers.TemplateInformer.HasSynced,
	)

	log.Log.Infof("webhook App: synched informers")

	http.HandleFunc(validating.VMTemplateCreateValidatePath, func(w http.ResponseWriter, r *http.Request) {
		validating.ServeVMTemplateCreate(w, r)
	})
	http.HandleFunc(validating.VMTemplateUpdateValidatePath, func(w http.ResponseWriter, r *http.Request) {
		validating.ServeVMTemplateUpdate(w, r)
	})
	if !app.TLSInfo.IsEnabled() {
		log.Log.Infof("webhook App: TLS *NOT* configured, serving over HTTP")
		return http.ListenAndServe(app.ListenAddress, nil)
	}
	log.Log.Infof("webhook App: TLS configured, serving over HTTPS")
	return http.ListenAndServeTLS(app.ListenAddress, app.TLSInfo.CertFilePath, app.TLSInfo.KeyFilePath, nil)
}
