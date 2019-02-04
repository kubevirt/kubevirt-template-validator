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
 * Copyright 2017, 2018 Red Hat, Inc.
 *
 */

package virtinformers

import (
	"math/rand"
	"sync"
	"time"

	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	templatev1 "github.com/openshift/api/template/v1"

	"github.com/fromanirh/kubevirt-template-validator/internal/pkg/log"
)

type newSharedInformer func() cache.SharedIndexInformer

type KubeInformerFactory interface {
	// Starts any informers that have not been started yet
	// This function is thread safe and idempotent
	Start(stopCh <-chan struct{})

	Template() cache.SharedIndexInformer
}

type kubeInformerFactory struct {
	restClient    *rest.RESTClient
	lock          sync.Mutex
	defaultResync time.Duration

	informers         map[string]cache.SharedIndexInformer
	startedInformers  map[string]bool
	kubevirtNamespace string
}

func NewKubeInformerFactory(restClient *rest.RESTClient, kubevirtNamespace string) KubeInformerFactory {
	return &kubeInformerFactory{
		restClient: restClient,
		// Resulting resync period will be between 12 and 24 hours, like the default for k8s
		defaultResync:     resyncPeriod(12 * time.Hour),
		informers:         make(map[string]cache.SharedIndexInformer),
		startedInformers:  make(map[string]bool),
		kubevirtNamespace: kubevirtNamespace,
	}
}

// Start can be called from multiple controllers in different go routines safely.
// Only informers that have not started are triggered by this function.
// Multiple calls to this function are idempotent.
func (f *kubeInformerFactory) Start(stopCh <-chan struct{}) {
	f.lock.Lock()
	defer f.lock.Unlock()

	for name, informer := range f.informers {
		if f.startedInformers[name] {
			// skip informers that have already started.
			log.Log.Infof("SKIPPING informer %s", name)
			continue
		}
		log.Log.Infof("STARTING informer %s", name)
		go informer.Run(stopCh)
		f.startedInformers[name] = true
	}
}

// internal function used to retrieve an already created informer
// or create a new informer if one does not already exist.
// Thread safe
func (f *kubeInformerFactory) getInformer(key string, newFunc newSharedInformer) cache.SharedIndexInformer {
	f.lock.Lock()
	defer f.lock.Unlock()

	informer, exists := f.informers[key]
	if exists {
		return informer
	}
	informer = newFunc()
	f.informers[key] = informer

	return informer
}

func (f *kubeInformerFactory) Template() cache.SharedIndexInformer {
	return f.getInformer("templateInformer", func() cache.SharedIndexInformer {
		lw := cache.NewListWatchFromClient(f.restClient, "templates", k8sv1.NamespaceAll, fields.Everything())
		return cache.NewSharedIndexInformer(lw, &templatev1.Template{}, f.defaultResync, cache.Indexers{})
	})
}

// resyncPeriod computes the time interval a shared informer waits before resyncing with the api server
func resyncPeriod(minResyncPeriod time.Duration) time.Duration {
	factor := rand.Float64() + 1
	return time.Duration(float64(minResyncPeriod.Nanoseconds()) * factor)
}
