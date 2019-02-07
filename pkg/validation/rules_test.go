package validation_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fromanirh/kubevirt-template-validator/pkg/validation"
)

var _ = Describe("Rules", func() {
	Context("Without validation text", func() {
		It("Should return no rules", func() {
			rules, err := validation.ParseRules([]byte(""))

			Expect(err).To(Not(HaveOccurred()))
			Expect(len(rules)).To(Equal(0))
		})
	})

	Context("With validation text", func() {
		It("Should parse a single rule", func() {
			text := `[{
            "name": "core-limits",
            "valid": "spec.domain.cpu.cores",
            "path": "spec.domain.cpu.cores",
            "rule": "integer",
            "message": "cpu cores must be limited",
            "min": 1,
            "max": 8
          }]`
			rules, err := validation.ParseRules([]byte(text))

			Expect(err).To(Not(HaveOccurred()))
			Expect(len(rules)).To(Equal(1))
		})
	})
})
