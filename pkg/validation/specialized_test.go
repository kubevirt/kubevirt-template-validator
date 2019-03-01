package validation_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"

	k6tobjs "github.com/fromanirh/kubevirt-template-validator/pkg/kubevirtobjs"
	"github.com/fromanirh/kubevirt-template-validator/pkg/validation"
)

var _ = Describe("Specialized", func() {
	Context("With invvalid data", func() {

		var (
			vmCirros *k6tv1.VirtualMachine
			vmRef    *k6tv1.VirtualMachine
		)

		BeforeEach(func() {
			vmCirros = NewVMCirros()
			vmRef = k6tobjs.NewDefaultVirtualMachine()
		})

		It("Should detect bogus rules", func() {
			r := validation.Rule{
				Rule:    "integer-value",
				Name:    "EnoughMemory",
				Path:    "jsonpath::.spec.domain.resources.requests.memory",
				Message: "Memory size not specified",
				Valid:   "jsonpath::.spec.domain.this.path.does.not.exist",
				Min:     64 * 1024 * 1024,
				Max:     512 * 1024 * 1024,
			}

			ra, err := r.Specialize(vmCirros, vmRef)
			Expect(err).To(Not(BeNil()))
			Expect(ra).To(BeNil())
		})
	})

	Context("With valid data", func() {

		var (
			vmCirros *k6tv1.VirtualMachine
			vmRef    *k6tv1.VirtualMachine
		)

		BeforeEach(func() {
			vmCirros = NewVMCirros()
			vmRef = k6tobjs.NewDefaultVirtualMachine()
		})

		It("Should apply simple integer rules", func() {
			r := validation.Rule{
				Rule:    "integer",
				Name:    "EnoughMemory",
				Path:    "jsonpath::.spec.domain.resources.requests.memory",
				Message: "Memory size not specified",
				Min:     64 * 1024 * 1024,
				Max:     512 * 1024 * 1024,
			}

			checkRuleApplication(&r, vmCirros, vmRef, true)
		})

		It("Should apply simple string rules", func() {
			r := validation.Rule{
				Rule:      "string",
				Name:      "HasChipset",
				Path:      "jsonpath::.spec.domain.machine.type",
				Message:   "machine type must be specified",
				MinLength: 1,
				MaxLength: 32,
			}
			checkRuleApplication(&r, vmCirros, vmRef, true)
		})

		It("Should apply simple enum rules", func() {
			r := validation.Rule{
				Rule:    "enum",
				Name:    "SupportedChipset",
				Path:    "jsonpath::.spec.domain.machine.type",
				Message: "machine type must be a supported value",
				Values:  []string{"q35"},
			}
			checkRuleApplication(&r, vmCirros, vmRef, true)
		})

		It("Should apply simple regex rules", func() {
			r := validation.Rule{
				Rule:    "regex",
				Name:    "SupportedChipset",
				Path:    "jsonpath::.spec.domain.machine.type",
				Message: "machine type must be a supported value",
				Regex:   "q35|440fx",
			}
			checkRuleApplication(&r, vmCirros, vmRef, true)

		})
	})

	Context("With INvalid data", func() {

		var (
			vmCirros *k6tv1.VirtualMachine
			vmRef    *k6tv1.VirtualMachine
		)

		BeforeEach(func() {
			vmCirros = NewVMCirros()
			vmRef = k6tobjs.NewDefaultVirtualMachine()
		})

		It("Should fail simple integer rules", func() {
			r1 := validation.Rule{
				Rule:    "integer",
				Name:    "EnoughMemory",
				Path:    "jsonpath::.spec.domain.resources.requests.memory",
				Message: "Memory size not specified",
				Valid:   "jsonpath::.spec.domain.this.path.does.not.exist",
				Min:     512 * 1024 * 1024,
			}
			checkRuleApplication(&r1, vmCirros, vmRef, false)

			r2 := validation.Rule{
				Rule:    "integer",
				Name:    "EnoughMemory",
				Path:    "jsonpath::.spec.domain.resources.requests.memory",
				Message: "Memory size not specified",
				Valid:   "jsonpath::.spec.domain.this.path.does.not.exist",
				Max:     64 * 1024 * 1024,
			}
			checkRuleApplication(&r2, vmCirros, vmRef, false)
		})

		It("Should apply simple string rules", func() {
			r1 := validation.Rule{
				Rule:      "string",
				Name:      "HasChipset",
				Path:      "jsonpath::.spec.domain.machine.type",
				Message:   "machine type must be specified",
				MinLength: 64,
			}
			checkRuleApplication(&r1, vmCirros, vmRef, false)

			r2 := validation.Rule{
				Rule:      "string",
				Name:      "HasChipset",
				Path:      "jsonpath::.spec.domain.machine.type",
				Message:   "machine type must be specified",
				MaxLength: 1,
			}
			checkRuleApplication(&r2, vmCirros, vmRef, false)

		})

		It("Should apply simple enum rules", func() {
			r := validation.Rule{
				Rule:    "enum",
				Name:    "SupportedChipset",
				Path:    "jsonpath::.spec.domain.machine.type",
				Message: "machine type must be a supported value",
				Values:  []string{"foo", "bar"},
			}
			checkRuleApplication(&r, vmCirros, vmRef, false)
		})

		It("Should apply simple regex rules", func() {
			r := validation.Rule{
				Rule:    "regex",
				Name:    "SupportedChipset",
				Path:    "jsonpath::.spec.domain.machine.type",
				Message: "machine type must be a supported value",
				Regex:   "\\d[a-z]+\\d\\d",
			}
			checkRuleApplication(&r, vmCirros, vmRef, false)

		})
	})

})

func checkRuleApplication(r *validation.Rule, vm, ref *k6tv1.VirtualMachine, expected bool) {

	ra, err := r.Specialize(vm, ref)
	Expect(err).To(BeNil())
	Expect(ra).To(Not(BeNil()))

	ok, err := ra.Apply(vm, ref)
	Expect(err).To(BeNil())
	Expect(ok).To(Equal(expected))
}
