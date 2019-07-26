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
 * Copyright 2017 Red Hat, Inc.
 *
 */

package service

import (
	goflag "flag"
	"fmt"
	"strconv"

	flag "github.com/spf13/pflag"

	"k8s.io/klog"
)

type Service interface {
	Run()
	AddFlags()
}

type ServiceListen struct {
	Name        string
	BindAddress string
	Port        int
}

type ServiceLibvirt struct {
	LibvirtUri string
}

func (service *ServiceListen) Address() string {
	return fmt.Sprintf("%s:%s", service.BindAddress, strconv.Itoa(service.Port))
}

func (service *ServiceListen) InitFlags() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

func (service *ServiceListen) AddCommonFlags() {
	flag.StringVar(&service.BindAddress, "listen", service.BindAddress, "Address where to listen on")
	flag.IntVar(&service.Port, "port", service.Port, "Port to listen on")
}

func (service *ServiceLibvirt) AddLibvirtFlags() {
	flag.StringVar(&service.LibvirtUri, "libvirt-uri", service.LibvirtUri, "Libvirt connection string")

}

func Setup(service Service) {
	service.AddFlags()

	defVerbose := "2"

	// set new default verbosity, was set to 0 by glog
	flag.Set("v", defVerbose)
	flag.Set("logtostderr", "true")

	flag.Parse()

	// borrowed from cdi/apiserver 1.9.5
	klogFlags := goflag.NewFlagSet("klog", goflag.ExitOnError)
	klog.InitFlags(klogFlags)
	flag.CommandLine.VisitAll(func(f1 *flag.Flag) {
		f2 := klogFlags.Lookup(f1.Name)
		if f2 != nil {
			value := f1.Value.String()
			f2.Value.Set(value)
		}
	})
}
