module github.com/kubevirt/kubevirt-template-validator

go 1.15

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/fsnotify/fsnotify v1.4.9
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/openshift/api v0.0.0
	github.com/openshift/client-go v0.0.0
	github.com/spf13/pflag v1.0.5
	k8s.io/api v0.20.6
	k8s.io/apimachinery v0.20.6
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	kubevirt.io/client-go v0.38.1
)

replace (
	github.com/go-kit/kit => github.com/go-kit/kit v0.3.0
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
	github.com/openshift/api => github.com/openshift/api v0.0.0-20200930075302-db52bc4ef99f // release-4.6
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20200929181438-91d71ef2122c // release-4.6
	github.com/operator-framework/operator-lifecycle-manager => github.com/operator-framework/operator-lifecycle-manager v0.0.0-20190128024246-5eb7ae5bdb7a
	k8s.io/client-go => k8s.io/client-go v0.20.6
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.20.6
	sigs.k8s.io/structured-merge-diff => sigs.k8s.io/structured-merge-diff v1.0.0
)
