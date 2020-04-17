module github.com/fromanirh/kubevirt-template-validator

go 1.12

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/evanphx/json-patch v4.5.0+incompatible // indirect
	github.com/fromanirh/okdutil v0.0.1
	github.com/go-openapi/jsonreference v0.19.2 // indirect
	github.com/golang/mock v0.0.0-20190713102442-dd8d2a22370e // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.1-0.20190515112211-6a48b4839f85
	github.com/openshift/api v3.9.1-0.20190401220125-3a6077f1f910+incompatible
	github.com/openshift/client-go v0.0.0-20190401163519-84c2b942258a
	github.com/pkg/errors v0.8.1 // indirect
	github.com/spf13/pflag v1.0.1
	k8s.io/api v0.0.0-20190222213804-5cb15d344471
	k8s.io/apimachinery v0.0.0-20190221213512-86fb29eff628
	k8s.io/client-go v0.0.0-20190228174230-b40b2a5939e4
	k8s.io/klog v0.3.0
	kubevirt.io/client-go v0.19.0
	kubevirt.io/containerized-data-importer v1.9.5 // indirect
)

replace github.com/go-kit/kit => github.com/go-kit/kit v0.3.0
