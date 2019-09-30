package kubevirtobjs_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKubevirtobjs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kubevirtobjs Suite")
}
