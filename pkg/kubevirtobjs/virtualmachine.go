/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2019 Red Hat, Inc.
 */

package kubevirtobjs

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"

	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

const (
	MaxDisks         uint = 64
	MaxIfaces        uint = 64
	MaxPortsPerIface uint = 16
	MaxNTPServers    uint = 8
)

var (
	ErrTooManyDisks         error = fmt.Errorf("Too many disks requested, max = %u", MaxDisks)
	ErrTooManyIfaces        error = fmt.Errorf("Too many network interface requested, max = %u", MaxIfaces)
	ErrTooManyPortsPerIface error = fmt.Errorf("Too many ports per network interface requested, max = %u", MaxPortsPerIface)
)

func newDisks(num uint) ([]k6tv1.Disk, error) {
	var disks []k6tv1.Disk
	if num > MaxDisks {
		return disks, ErrTooManyDisks
	}
	for i := uint(0); i < num; i++ {
		disk := k6tv1.Disk{
			DiskDevice: k6tv1.DiskDevice{
				Disk:   &k6tv1.DiskTarget{},
				LUN:    &k6tv1.LunTarget{},
				Floppy: &k6tv1.FloppyTarget{},
				CDRom:  &k6tv1.CDRomTarget{},
			},
			BootOrder:         new(uint),
			DedicatedIOThread: new(bool),
		}
		disks = append(disks, disk)
	}
	return disks, nil
}

func newInterfaces(num, ports uint) ([]k6tv1.Interface, error) {
	var ifaces []k6tv1.Interface
	if num > MaxIfaces {
		return ifaces, ErrTooManyIfaces
	}
	if ports > MaxPortsPerIface {
		return ifaces, ErrTooManyPortsPerIface
	}
	for i := uint(0); i < num; i++ {
		iface := k6tv1.Interface{
			InterfaceBindingMethod: k6tv1.InterfaceBindingMethod{
				Bridge:     &k6tv1.InterfaceBridge{},
				Slirp:      &k6tv1.InterfaceSlirp{},
				Masquerade: &k6tv1.InterfaceMasquerade{},
				SRIOV:      &k6tv1.InterfaceSRIOV{},
			},
			Ports:     make([]k6tv1.Port, ports),
			BootOrder: new(uint),
			DHCPOptions: &k6tv1.DHCPOptions{
				NTPServers: make([]string, MaxNTPServers), // intentionally not user-configurable. It looks like a too obscure detail to be exposed
			},
		}
		ifaces = append(ifaces, iface)
	}
	return ifaces, nil
}

// NewDomainSpec returns a fully zero-value DomainSpec with all optional fields, or error if requested parameters exceeds limits
func NewDomainSpec(numDisks, numIfaces, numPortsPerIface uint) (*k6tv1.DomainSpec, error) {
	var err error

	disks, err := newDisks(numDisks)
	if err != nil {
		return nil, err
	}
	ifaces, err := newInterfaces(numIfaces, numPortsPerIface)
	if err != nil {
		return nil, err
	}

	dom := k6tv1.DomainSpec{
		// TODO: resources
		CPU: &k6tv1.CPU{},
		Memory: &k6tv1.Memory{
			Hugepages: &k6tv1.Hugepages{},
			Guest:     &resource.Quantity{},
		},
		Firmware: &k6tv1.Firmware{},
		Clock: &k6tv1.Clock{
			Timer: &k6tv1.Timer{
				HPET:   &k6tv1.HPETTimer{Enabled: new(bool)},
				KVM:    &k6tv1.KVMTimer{Enabled: new(bool)},
				PIT:    &k6tv1.PITTimer{Enabled: new(bool)},
				RTC:    &k6tv1.RTCTimer{Enabled: new(bool)},
				Hyperv: &k6tv1.HypervTimer{Enabled: new(bool)},
			},
		},
		Features: &k6tv1.Features{
			ACPI:   k6tv1.FeatureState{Enabled: new(bool)},
			APIC:   &k6tv1.FeatureAPIC{Enabled: new(bool)},
			Hyperv: &k6tv1.FeatureHyperv{}, // TODO
		},
		Devices: k6tv1.Devices{
			Disks: disks,
			Watchdog: &k6tv1.Watchdog{
				WatchdogDevice: k6tv1.WatchdogDevice{
					I6300ESB: &k6tv1.I6300ESBWatchdog{},
				},
			},
			Interfaces:                 ifaces,
			AutoattachPodInterface:     new(bool),
			AutoattachGraphicsDevice:   new(bool),
			Rng:                        &k6tv1.Rng{},
			BlockMultiQueue:            new(bool),
			NetworkInterfaceMultiQueue: new(bool),
		},
		IOThreadsPolicy: new(k6tv1.IOThreadsPolicy),
	}
	return &dom, nil
}

// NewVirtualMachine returns a fully zero-value VirtualMachine with all optional fields, or error if requested parameters exceeds limits
func NewVirtualMachine(numDisks, numIfaces, numPortsPerIface uint) (*k6tv1.VirtualMachine, error) {
	dom, err := NewDomainSpec(numDisks, numIfaces, numPortsPerIface)

	tmpl := k6tv1.VirtualMachineInstanceTemplateSpec{}
	tmpl.Spec.Domain = *dom

	var vm k6tv1.VirtualMachine
	vm.Spec.Template = &tmpl
	k6tv1.SetObjectDefaults_VirtualMachine(&vm)
	return &vm, err
}

// NewVirtualMachine returns a fully zero-value VirtualMachine with all optional fields
func NewDefaultVirtualMachine() *k6tv1.VirtualMachine {
	vm, _ := NewVirtualMachine(MaxDisks, MaxIfaces, MaxPortsPerIface)
	return vm
}
