package validation_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"

	"github.com/fromanirh/kubevirt-template-validator/pkg/validation"
)

var _ = Describe("Eval", func() {
	Context("With invalid rule set", func() {
		It("Should detect duplicate names", func() {
			rules := []validation.Rule{
				validation.Rule{
					Name: "rule-1",
					Rule: "integer",
					// any legal path is fine
					Path:    ".spec.domain.cpu.cores",
					Message: "testing",
				},
				validation.Rule{
					Name: "rule-1",
					Rule: "string",
					// any legal path is fine
					Path:    ".spec.domain.cpu.cores",
					Message: "testing",
				},
			}
			vm := k6tv1.VirtualMachine{}

			res := validation.Evaluate(&vm, rules)
			Expect(res.Succeeded()).To(BeFalse())
			Expect(len(res.Status)).To(Equal(2))
			Expect(res.Status[0].IllegalReason).To(Not(BeNil()))
			Expect(res.Status[1].IllegalReason).To(Equal(validation.ErrDuplicateRuleName))

		})

		It("Should detect missing keys", func() {
			rules := []validation.Rule{
				validation.Rule{
					Name: "rule-1",
					Rule: "integer",
					// any legal path is fine
					Path: ".spec.domain.cpu.cores",
				},
				validation.Rule{
					Name:    "rule-2",
					Rule:    "string",
					Message: "testing",
				},
			}
			vm := k6tv1.VirtualMachine{}

			res := validation.Evaluate(&vm, rules)
			Expect(res.Succeeded()).To(BeFalse())
			Expect(len(res.Status)).To(Equal(2))
			Expect(res.Status[0].IllegalReason).To(Equal(validation.ErrMissingRequiredKey))
			Expect(res.Status[1].IllegalReason).To(Equal(validation.ErrMissingRequiredKey))
		})

		It("Should detect invalid rules", func() {
			rules := []validation.Rule{
				validation.Rule{
					Name: "rule-1",
					Rule: "foobar",
					// any legal path is fine
					Path:    ".spec.domain.cpu.cores",
					Message: "testing",
				},
			}
			vm := k6tv1.VirtualMachine{}

			res := validation.Evaluate(&vm, rules)
			Expect(res.Succeeded()).To(BeFalse())
			Expect(len(res.Status)).To(Equal(1))
			Expect(res.Status[0].IllegalReason).To(Equal(validation.ErrUnrecognizedRuleType))
		})
		/*
			It("Should detect unappliable rules", func() {
				rules := []validation.Rule{
					validation.Rule{
						Name: "rule-1",
						Rule: "integer",
						// any legal path is fine
						Path:    ".spec.domain.cpu.cores",
						Message: "testing",
						Valid:   ".spec.domain.some.inexistent.path",
					},
				}
				vm := k6tv1.VirtualMachine{}

				res := validation.Evaluate(&vm, rules)
				Expect(res.Succeeded()).To(BeFalse())
				Expect(len(res.Status)).To(Equal(1))
				Expect(res.Status[0].Satisfied).To(BeFalse())
				Expect(res.Status[0].IllegalReason).To(BeNil())
			})
		*/
	})
})
