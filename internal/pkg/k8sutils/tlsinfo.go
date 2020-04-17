package k8sutils

import (
	"io/ioutil"
	"os"

	"k8s.io/client-go/rest"
	"kubevirt.io/client-go/log"
)

type TLSInfo struct {
	CertFilePath   string
	KeyFilePath    string
	certsDirectory string
}

func (ti *TLSInfo) Clean() {
	if ti.certsDirectory != "" {
		os.RemoveAll(ti.certsDirectory)
		ti.certsDirectory = ""
		log.Log.V(4).Infof("TLSInfo cleaned up!")
	}
}

func (ti *TLSInfo) IsEnabled() bool {
	return ti.CertFilePath != "" && ti.KeyFilePath != ""
}

func (ti *TLSInfo) UpdateFromK8S() error {
	var err error
	if _, err = rest.InClusterConfig(); err != nil {
		// is not a real error, rather a supported case. So, let's swallow the error
		log.Log.V(3).Infof("running outside a K8S cluster")
		return nil
	}
	if ti.IsEnabled() {
		log.Log.V(3).Infof("TLSInfo already fully set")
		return nil
	}

	// at least one between cert and key need to be set
	ti.certsDirectory, err = ioutil.TempDir("", "certsdir")
	if err != nil {
		return err
	}
	namespace, err := GetNamespace()
	if err != nil {
		log.Log.Warningf("Error searching for namespace: %v", err)
		return err
	}
	certStore, err := GenerateSelfSignedCert(ti.certsDirectory, "kubevirt-template-validator", namespace)
	if err != nil {
		log.Log.Warningf("unable to generate certificates: %v", err)
		return err
	}

	if ti.CertFilePath == "" {
		ti.CertFilePath = certStore.CurrentPath()
	} else {
		log.Log.V(2).Infof("NOT overriding cert file %s with %s", ti.CertFilePath, certStore.CurrentPath())
	}
	if ti.KeyFilePath == "" {
		ti.KeyFilePath = certStore.CurrentPath()
	} else {
		log.Log.V(2).Infof("NOT overriding key file %s with %s", ti.KeyFilePath, certStore.CurrentPath())
	}
	log.Log.V(3).Infof("running in a K8S cluster: with configuration %#v", *ti)
	return nil
}
