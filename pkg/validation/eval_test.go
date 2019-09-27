package validation_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	k6tv1 "kubevirt.io/client-go/api/v1"

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
					Path:    "jsonpath::.spec.domain.cpu.cores",
					Message: "testing",
				},
				validation.Rule{
					Name: "rule-1",
					Rule: "string",
					// any legal path is fine
					Path:    "jsonpath::.spec.domain.cpu.cores",
					Message: "testing",
				},
			}
			vm := k6tv1.VirtualMachine{}

			res := validation.NewEvaluator().Evaluate(rules, &vm)
			Expect(res.Succeeded()).To(BeFalse())
			Expect(len(res.Status)).To(Equal(2))
			Expect(res.Status[0].Error).To(BeNil())
			Expect(res.Status[1].Error).To(Equal(validation.ErrDuplicateRuleName))

		})

		It("Should detect missing keys", func() {
			rules := []validation.Rule{
				validation.Rule{
					Name: "rule-1",
					Rule: "integer",
					// any legal path is fine
					Path: "jsonpath::.spec.domain.cpu.cores",
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
			Expect(res.Status[0].Error).To(Equal(validation.ErrMissingRequiredKey))
			Expect(res.Status[1].Error).To(Equal(validation.ErrMissingRequiredKey))
		})

		It("Should detect invalid rules", func() {
			rules := []validation.Rule{
				validation.Rule{
					Name: "rule-1",
					Rule: "foobar",
					// any legal path is fine
					Path:    "jsonpath::.spec.domain.cpu.cores",
					Message: "testing",
				},
			}
			vm := k6tv1.VirtualMachine{}

			res := validation.NewEvaluator().Evaluate(rules, &vm)
			Expect(res.Succeeded()).To(BeFalse())
			Expect(len(res.Status)).To(Equal(1))
			Expect(res.Status[0].Error).To(Equal(validation.ErrUnrecognizedRuleType))
		})
		It("Should detect unappliable rules", func() {
			rules := []validation.Rule{
				validation.Rule{
					Name: "rule-1",
					Rule: "integer",
					// any legal path is fine
					Path:    "jsonpath::.spec.domain.cpu.cores",
					Message: "testing",
					Valid:   "jsonpath::.spec.domain.some.inexistent.path",
				},
			}
			vm := k6tv1.VirtualMachine{}

			ev := validation.Evaluator{Sink: GinkgoWriter}
			res := ev.Evaluate(rules, &vm)

			Expect(res.Succeeded()).To(BeTrue())
			Expect(len(res.Status)).To(Equal(1))
			Expect(res.Status[0].Skipped).To(BeTrue())
			Expect(res.Status[0].Satisfied).To(BeFalse())
			Expect(res.Status[0].Error).To(BeNil())
		})

		It("Should not fail, when justWarning is set", func() {
			rules := []validation.Rule{
				{
					Name: "rule-1",
					Rule: "integer",
					Min:  8,
					// any legal path is fine
					Path:        "jsonpath::.spec.domain.cpu.cores",
					Message:     "testing",
					JustWarning: true,
				},
			}
			vm := k6tv1.VirtualMachine{}

			ev := validation.Evaluator{Sink: GinkgoWriter}
			res := ev.Evaluate(rules, &vm)

			Expect(res.Succeeded()).To(BeTrue(), "succeeded")
			Expect(len(res.Status)).To(Equal(1), "status length")
			Expect(res.Status[0].Skipped).To(BeFalse(), "skipped")
			Expect(res.Status[0].Satisfied).To(BeFalse(), "satisfied")
			Expect(res.Status[0].Error).To(BeNil(), "error")
		})
	})

	Context("With an initialized VM object", func() {
		var (
			vmCirros *k6tv1.VirtualMachine
		)

		BeforeEach(func() {
			vmCirros = NewVMCirros()
		})

		It("should skip uninitialized paths if requested", func() {
			rules := []validation.Rule{
				validation.Rule{
					Name:    "LimitCores",
					Rule:    "integer",
					Path:    "jsonpath::.spec.domain.cpu.cores",
					Valid:   "jsonpath::.spec.domain.cpu.cores",
					Message: "testing",
					Min:     1,
					Max:     8,
				},
			}

			ev := validation.Evaluator{Sink: GinkgoWriter}
			res := ev.Evaluate(rules, vmCirros)

			Expect(res.Succeeded()).To(BeTrue())
			Expect(len(res.Status)).To(Equal(1))
			Expect(res.Status[0].Skipped).To(BeTrue())
			Expect(res.Status[0].Satisfied).To(BeFalse())
			Expect(res.Status[0].Error).To(BeNil())

		})

		It("should handle uninitialized paths", func() {
			rules := []validation.Rule{
				validation.Rule{
					Name:    "LimitCores",
					Rule:    "integer",
					Path:    "jsonpath::.spec.domain.cpu.cores",
					Message: "testing",
					Min:     1,
					Max:     8,
				},
			}

			ev := validation.Evaluator{Sink: GinkgoWriter}
			res := ev.Evaluate(rules, vmCirros)

			Expect(res.Succeeded()).To(BeFalse())
		})

		It("should handle uninitialized paths intermixed with valid paths", func() {
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
					Rule:    "integer",
					Name:    "LimitCores",
					Path:    "jsonpath::.spec.domain.cpu.cores",
					Message: "Core amount not within range",
					Min:     1,
					Max:     4,
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
			Expect(res.Succeeded()).To(BeFalse())

			causes := res.ToStatusCauses()
			Expect(len(causes)).To(Equal(1))
		})

		It("should not fail, when justWarning is set", func() {
			rules := []validation.Rule{
				{
					Name:        "disk bus",
					Rule:        "enum",
					Path:        "jsonpath::.spec.domain.devices.disks[*].disk.bus",
					Message:     "testing",
					Values:      []string{"sata"},
					JustWarning: true,
				},
			}

			ev := validation.Evaluator{Sink: GinkgoWriter}
			res := ev.Evaluate(rules, vmCirros)

			Expect(res.Succeeded()).To(BeTrue(), "succeeded")
			Expect(len(res.Status)).To(Equal(1), "status length")
			Expect(res.Status[0].Skipped).To(BeFalse(), "skipped")
			Expect(res.Status[0].Satisfied).To(BeFalse(), "satisfied")
			Expect(res.Status[0].Error).To(BeNil(), "error")
		})

		It("should fail, when one rule does not have justWarning set", func() {
			rules := []validation.Rule{
				{
					Name:        "disk bus",
					Rule:        "enum",
					Path:        "jsonpath::.spec.domain.devices.disks[*].disk.bus",
					Message:     "testing",
					Values:      []string{"sata"},
					JustWarning: true,
				}, {
					Name: "rule-2",
					Rule: "integer",
					Min:  6,
					Max:  8,
					// any legal path is fine
					Path:    "jsonpath::.spec.domain.cpu.cores",
					Message: "enough cores",
				},
			}

			ev := validation.Evaluator{Sink: GinkgoWriter}
			res := ev.Evaluate(rules, vmCirros)

			Expect(res.Succeeded()).To(BeFalse(), "succeeded")
			Expect(len(res.Status)).To(Equal(2), "status length")
			Expect(res.Status[0].Skipped).To(BeFalse(), "skipped")
			Expect(res.Status[0].Satisfied).To(BeFalse(), "satisfied")
			Expect(res.Status[0].Error).To(BeNil(), "error")
			Expect(res.Status[1].Skipped).To(BeFalse(), "skipped")
			Expect(res.Status[1].Satisfied).To(BeFalse(), "satisfied")
			Expect(res.Status[1].Error).To(BeNil(), "error")
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
				Expect(res.Status[ix].Error).To(BeNil())
			}
		})

		It("Should fail applying a ruleset with at least one malformed rule", func() {
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
					Rule:    "value-set",
					Name:    "SupportedChipset",
					Path:    "jsonpath::.spec.domain.machine.type",
					Message: "machine type must be a supported value",
					Values:  []string{"q35"},
				},
			}

			ev := validation.Evaluator{Sink: GinkgoWriter}
			res := ev.Evaluate(rules, vmCirros)
			Expect(res.Succeeded()).To(BeFalse())
		})

		It("Should fail, when rule with justWarning has incorrect path and another rule is correct", func() {
			rules := []validation.Rule{
				{
					Name:        "disk bus",
					Rule:        "enum",
					Path:        "jsonpath::.spec.domain.devices.some.non.existing.path",
					Message:     "testing",
					Values:      []string{"sata"},
					JustWarning: true,
				}, {
					Name: "rule-2",
					Rule: "integer",
					Min:  0,
					Max:  8,
					// any legal path is fine
					Path:    "jsonpath::.spec.domain.cpu.cores",
					Message: "enough cores",
				},
			}

			ev := validation.Evaluator{Sink: GinkgoWriter}
			res := ev.Evaluate(rules, vmCirros)

			for ix := range res.Status {
				fmt.Fprintf(GinkgoWriter, "%+#v", res.Status[ix])
			}

			Expect(res.Succeeded()).To(BeFalse(), "succeeded")
			Expect(len(res.Status)).To(Equal(2), "status length")
			//status for second rule which should pass
			Expect(res.Status[0].Skipped).To(BeFalse(), "skipped")
			Expect(res.Status[0].Satisfied).To(BeFalse(), "satisfied")
			Expect(res.Status[0].Error).NotTo(BeNil(), "error") // expects invalid JSONPath

			Expect(res.Status[1].Skipped).To(BeFalse(), "skipped")
			Expect(res.Status[1].Satisfied).To(BeTrue(), "satisfied")
			Expect(res.Status[1].Error).To(BeNil(), "error")
		})
	})
})

// TODO:
// test with 2+ rules failed
// test to exercise the translation logic
