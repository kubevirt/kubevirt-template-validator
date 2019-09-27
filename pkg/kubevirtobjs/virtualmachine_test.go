package kubevirtobjs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fromanirh/kubevirt-template-validator/pkg/kubevirtobjs"

	k6tv1 "kubevirt.io/client-go/api/v1"
)

var _ = Describe("Virtualmachine", func() {
	var (
		vm *k6tv1.VirtualMachine
	)

	BeforeEach(func() {
		vm = kubevirtobjs.NewDefaultVirtualMachine()
	})

	Context("Sanity checks", func() {
		It("Should create objects", func() {
			Expect(vm).NotTo(BeNil())
		})

		It("Should create objects with domain", func() {
			Expect(vm.Spec.Template.Spec.Domain).NotTo(BeNil())
		})

	})

	Context("Check field access", func() {
		It("Should NOT have defaults set outside DomainSpec", func() {
			Expect(vm.Spec.Template.Spec.EvictionStrategy).To(BeNil())
		})

		It("Should have CPU", func() {
			Expect(vm.Spec.Template.Spec.Domain.CPU.Cores).To(Equal(uint32(0)))
			Expect(vm.Spec.Template.Spec.Domain.CPU.Sockets).To(Equal(uint32(0)))
			Expect(vm.Spec.Template.Spec.Domain.CPU.Threads).To(Equal(uint32(0)))
		})

		It("Should have guest memory", func() {
			Expect(vm.Spec.Template.Spec.Domain.Memory.Guest.IsZero()).To(BeTrue())
		})
		It("Should have hugepages", func() {
			Expect(vm.Spec.Template.Spec.Domain.Memory.Hugepages.PageSize).To(Equal(""))
		})
	})
})
