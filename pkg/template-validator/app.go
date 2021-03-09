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
	"fmt"
	"net/http"

	templatev1 "github.com/openshift/api/template/v1"
	flag "github.com/spf13/pflag"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"kubevirt.io/client-go/log"
	k6tversion "kubevirt.io/client-go/version"

	"github.com/kubevirt/kubevirt-template-validator/internal/pkg/k8sutils"
	"github.com/kubevirt/kubevirt-template-validator/internal/pkg/service"
	"github.com/kubevirt/kubevirt-template-validator/internal/pkg/version"
	"github.com/kubevirt/kubevirt-template-validator/pkg/virtinformers"
	"github.com/kubevirt/kubevirt-template-validator/pkg/webhooks/validating"
)

const (
	defaultPort = 8443
	defaultHost = "0.0.0.0"
)

func init() {
	// The Kubernetes Go client (nested within the OpenShift Go client)
	// automatically registers its types in scheme.Scheme, however the
	// additional OpenShift types must be registered manually.  AddToScheme
	// registers the API group types (e.g. route.openshift.io/v1, Route) only.
	utilruntime.Must(templatev1.Install(scheme.Scheme))
}

type App struct {
	service.ServiceListen
	TLSInfo       k8sutils.TLSInfo
	versionOnly   bool
	skipInformers bool
}

var _ service.Service = &App{}

func (app *App) AddFlags() {
	app.InitFlags()

	app.BindAddress = defaultHost
	app.Port = defaultPort

	app.AddCommonFlags()

	flag.StringVarP(&app.TLSInfo.CertsDirectory, "cert-dir", "c", "", "specify path to the directory containing TLS key and certificate - this enables TLS")
	flag.BoolVarP(&app.versionOnly, "version", "V", false, "show version and exit")
	flag.BoolVarP(&app.skipInformers, "skip-informers", "S", false, "don't initialize informerers - use this only in devel mode")
}

func (app *App) KubevirtVersion() string {
	info := k6tversion.Get()
	return fmt.Sprintf("%s %s %s", info.GitVersion, info.GitCommit, info.BuildDate)
}

func (app *App) Run() {
	log.Log.Infof("%s %s (revision: %s) starting", version.COMPONENT, version.VERSION, version.REVISION)
	log.Log.Infof("%s using kubevirt client-go (%s)", version.COMPONENT, app.KubevirtVersion())
	if app.versionOnly {
		return
	}

	app.TLSInfo.Init()
	defer app.TLSInfo.Clean()

	stopChan := make(chan struct{}, 1)
	defer close(stopChan)

	if app.skipInformers {
		log.Log.Infof("validator app: informers DISALBED")
		virtinformers.SetInformers(nil)
	}

	informers := virtinformers.GetInformers()
	if !informers.Available() {
		log.Log.Infof("validator app: template informer NOT available")
	} else {
		go informers.TemplateInformer.Run(stopChan)
		log.Log.Infof("validator app: started informers")
		cache.WaitForCacheSync(
			stopChan,
			informers.TemplateInformer.HasSynced,
		)
		log.Log.Infof("validator app: synched informers")
	}

	log.Log.Infof("validator app: running with TLSInfo.CertsDirectory%+v", app.TLSInfo.CertsDirectory)

	http.HandleFunc(validating.VMTemplateValidatePath, func(w http.ResponseWriter, r *http.Request) {
		validating.ServeVMTemplateValidate(w, r)
	})

	if app.TLSInfo.IsEnabled() {
		server := &http.Server{Addr: app.Address(), TLSConfig: app.TLSInfo.CrateTlsConfig()}
		log.Log.Infof("validator app: TLS configured, serving over HTTPS on %s", app.Address())
		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Log.Criticalf("Error listening TLS: %s", err)
			panic(err)
		}
	} else {
		log.Log.Infof("validator app: TLS *NOT* configured, serving over HTTP on %s", app.Address())
		if err := http.ListenAndServe(app.Address(), nil); err != nil {
			log.Log.Criticalf("Error listening TLS: %s", err)
			panic(err)
		}
	}
}
