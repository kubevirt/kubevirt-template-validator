package validating_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestValidating(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Validating Suite")
}
