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

	flag "github.com/spf13/pflag"

	"k8s.io/client-go/tools/cache"

	_ "github.com/fromanirh/okdutil/okd"

	"github.com/fromanirh/kubevirt-template-validator/internal/pkg/k8sutils"
	"github.com/fromanirh/kubevirt-template-validator/internal/pkg/log"
	"github.com/fromanirh/kubevirt-template-validator/internal/pkg/service"

	"github.com/fromanirh/kubevirt-template-validator/pkg/virtinformers"
	"github.com/fromanirh/kubevirt-template-validator/pkg/webhooks/validating"
)

const (
	defaultPort = 8443
	defaultHost = "0.0.0.0"
)

type App struct {
	service.ServiceListen
	TLSInfo  k8sutils.TLSInfo
	DumpMode bool
}

var _ service.Service = &App{}

func (app *App) AddFlags() {
	app.InitFlags()

	app.BindAddress = defaultHost
	app.Port = defaultPort

	app.AddCommonFlags()

	flag.StringVarP(&app.TLSInfo.CertFilePath, "cert-file", "c", "", "override path to TLS certificate - you need also the key to enable TLS")
	flag.StringVarP(&app.TLSInfo.KeyFilePath, "key-file", "k", "", "override path to TLS key - you need also the cert to enable TLS")
	flag.BoolVarP(&app.DumpMode, "dump", "D", false, "dump data involved in the admission control")
}

func (app *App) Run() {
	app.TLSInfo.UpdateFromK8S()
	defer app.TLSInfo.Clean()

	informers := virtinformers.GetInformers()
	stopChan := make(chan struct{}, 1)
	defer close(stopChan)

	if informers.TemplateInformer == nil {
		log.Log.Infof("validator app: no template informer available")
	} else {
		go informers.TemplateInformer.Run(stopChan)
		log.Log.Infof("validator app: started informers")
		cache.WaitForCacheSync(
			stopChan,
			informers.TemplateInformer.HasSynced,
		)
		log.Log.Infof("validator app: synched informers")
	}

	validating.SetDumpMode(app.DumpMode)

	log.Log.Infof("validator app: running with TLSInfo%+v", app.TLSInfo)

	http.HandleFunc(validating.VMTemplateValidatePath, func(w http.ResponseWriter, r *http.Request) {
		validating.ServeVMTemplateValidate(w, r)
	})

	if app.TLSInfo.IsEnabled() {
		log.Log.Infof("validator app: TLS configured, serving over HTTPS on %s", app.Address())
		http.ListenAndServeTLS(app.Address(), app.TLSInfo.CertFilePath, app.TLSInfo.KeyFilePath, nil)
	} else {
		log.Log.Infof("validator app: TLS *NOT* configured, serving over HTTP on %s", app.Address())
		http.ListenAndServe(app.Address(), nil)
	}
}
