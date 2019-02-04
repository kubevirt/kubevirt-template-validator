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
 * Copyright 2018 Red Hat, Inc.
 */

package main

import (
	"fmt"
	"net/http"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/fromanirh/kubevirt-template-validator/internal/pkg/k8sutils"
	"github.com/fromanirh/kubevirt-template-validator/internal/pkg/log"
	"github.com/fromanirh/kubevirt-template-validator/pkg/webhooks/validating"
	_ "github.com/fromanirh/okdutil/okd"
)

func Main() int {
	log.Log = log.Logger("kubevirt-template-validator")

	tlsInfo := &k8sutils.TLSInfo{}
	addr := flag.StringP("addr", "L", "", "address on which the server is listening to")
	port := flag.StringP("port", "P", "19999", "port on which the server is listening to")
	flag.StringVarP(&tlsInfo.CertFilePath, "cert-file", "c", "", "override path to TLS certificate - you need also the key to enable TLS")
	flag.StringVarP(&tlsInfo.KeyFilePath, "key-file", "k", "", "override path to TLS key - you need also the cert to enable TLS")
	flag.Parse()

	listenAddress := fmt.Sprintf("%s:%s", *addr, *port)

	log.Log.Infof("kubevirt-template-validator started on %v", listenAddress)
	defer log.Log.Infof("kubevirt-template-validator stopped")

	tlsInfo.UpdateFromK8S()
	defer tlsInfo.Clean()

	http.HandleFunc(validating.VMTemplateCreateValidatePath, func(w http.ResponseWriter, r *http.Request) {
		validating.ServeVMTemplateCreate(w, r)
	})
	http.HandleFunc(validating.VMTemplateUpdateValidatePath, func(w http.ResponseWriter, r *http.Request) {
		validating.ServeVMTemplateUpdate(w, r)
	})
	if tlsInfo.IsEnabled() {
		log.Log.Infof("TLS configured, serving over HTTPS")
		log.Log.Infof("%s", http.ListenAndServeTLS(listenAddress, tlsInfo.CertFilePath, tlsInfo.KeyFilePath, nil))
	} else {
		log.Log.Infof("TLS *NOT* configured, serving over HTTP")
		log.Log.Infof("%s", http.ListenAndServe(listenAddress, nil))
	}
	return 0
}

func main() {
	os.Exit(Main())
}
