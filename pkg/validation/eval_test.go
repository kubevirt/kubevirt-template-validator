package validation_test

import (
	"fmt"

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

			res := validation.NewEvaluator().Evaluate(rules, &vm)
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

			res := validation.NewEvaluator().Evaluate(rules, &vm)
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

			res := validation.NewEvaluator().Evaluate(rules, &vm)
			Expect(res.Succeeded()).To(BeFalse())
			Expect(len(res.Status)).To(Equal(1))
			Expect(res.Status[0].IllegalReason).To(Equal(validation.ErrUnrecognizedRuleType))
		})
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

			res := validation.NewEvaluator().Evaluate(rules, &vm)
			Expect(res.Succeeded()).To(BeFalse())
			Expect(len(res.Status)).To(Equal(1))
			Expect(res.Status[0].Satisfied).To(BeFalse())
			Expect(res.Status[0].IllegalReason).To(Equal(validation.ErrInvalidJSONPath))
		})
	})

	Context("With valid rule set", func() {

		var (
			vmCirros *k6tv1.VirtualMachine
		)

		BeforeEach(func() {
			vmCirros = NewVMCirros()
		})

		It("Should succeed applying a ruleset", func() {
			rules := []validation.Rule{
				validation.Rule{
					Rule:    "integer",
					Name:    "EnoughMemory",
					Path:    "jsonpath::.spec.domain.resources.requests.memory",
					Message: "Memory size not specified",
					Min:     64 * 1024 * 1024,
					Max:     512 * 1024 * 1024,
				},
				validation.Rule{
					Rule:    "enum",
					Name:    "SupportedChipset",
					Path:    "jsonpath::.spec.domain.machine.type",
					Message: "machine type must be a supported value",
					Values:  []string{"q35"},
				},
			}

			ev := validation.Evaluator{Sink: GinkgoWriter}
			res := ev.Evaluate(rules, vmCirros)
			Expect(res.Succeeded()).To(BeTrue())

			for ix := range res.Status {
				fmt.Fprintf(GinkgoWriter, "%+#v", res.Status[ix])
			}

			Expect(len(res.Status)).To(Equal(len(rules)))
			for ix := range res.Status {
				Expect(res.Status[ix].Satisfied).To(BeTrue())
				Expect(res.Status[ix].IllegalReason).To(BeNil())
			}
		})
	})
})

// TODO:
// test with 2+ rules failed
// test to exercise the translation logic
