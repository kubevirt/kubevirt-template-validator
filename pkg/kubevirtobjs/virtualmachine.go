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
	"reflect"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

const (
	MaxDisks         uint = 64
	MaxIfaces        uint = 64
	MaxPortsPerIface uint = 16
	MaxNTPServers    uint = 8
	MaxItems         int  = 64
)

var (
	ErrTooManyDisks         error = fmt.Errorf("Too many disks requested, max = %d", MaxDisks)
	ErrTooManyIfaces        error = fmt.Errorf("Too many network interface requested, max = %d", MaxIfaces)
	ErrTooManyPortsPerIface error = fmt.Errorf("Too many ports per network interface requested, max = %d", MaxPortsPerIface)
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

func NewFeatureHyperv() *k6tv1.FeatureHyperv {
	enabled := false
	featState := k6tv1.FeatureState{
		Enabled: &enabled,
	}
	featSpinLocks := k6tv1.FeatureSpinlocks{
		Enabled: &enabled,
		Retries: new(uint32),
	}
	featVendorID := k6tv1.FeatureVendorID{
		Enabled: &enabled,
	}

	return &k6tv1.FeatureHyperv{
		Relaxed:    &featState,
		VAPIC:      &featState,
		Spinlocks:  &featSpinLocks,
		VPIndex:    &featState,
		Runtime:    &featState,
		SyNIC:      &featState,
		SyNICTimer: &featState,
		Reset:      &featState,
		VendorID:   &featVendorID,
	}
}

// NewDomainSpec returns a fully zero-value DomainSpec with all optional fields, or error if requested parameters exceeds limits
// TODO: build using instrospection (aka the reflect package)
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

	enabled := false

	dom := k6tv1.DomainSpec{
		Resources: k6tv1.ResourceRequirements{
			Requests: make(v1.ResourceList),
			Limits:   make(v1.ResourceList),
		},
		CPU: &k6tv1.CPU{},
		Memory: &k6tv1.Memory{
			Hugepages: &k6tv1.Hugepages{},
			Guest:     &resource.Quantity{},
		},
		Firmware: &k6tv1.Firmware{},
		Clock: &k6tv1.Clock{
			ClockOffset: k6tv1.ClockOffset{
				UTC: &k6tv1.ClockOffsetUTC{
					OffsetSeconds: new(int),
				},
				Timezone: new(k6tv1.ClockOffsetTimezone),
			},
			Timer: &k6tv1.Timer{
				HPET:   &k6tv1.HPETTimer{Enabled: &enabled},
				KVM:    &k6tv1.KVMTimer{Enabled: &enabled},
				PIT:    &k6tv1.PITTimer{Enabled: &enabled},
				RTC:    &k6tv1.RTCTimer{Enabled: &enabled},
				Hyperv: &k6tv1.HypervTimer{Enabled: &enabled},
			},
		},
		Features: &k6tv1.Features{
			ACPI:   k6tv1.FeatureState{Enabled: &enabled},
			APIC:   &k6tv1.FeatureAPIC{Enabled: &enabled},
			Hyperv: NewFeatureHyperv(),
		},
		Devices: k6tv1.Devices{
			Disks: disks,
			Watchdog: &k6tv1.Watchdog{
				WatchdogDevice: k6tv1.WatchdogDevice{
					I6300ESB: &k6tv1.I6300ESBWatchdog{},
				},
			},
			Interfaces:                 ifaces,
			AutoattachPodInterface:     &enabled,
			AutoattachGraphicsDevice:   &enabled,
			Rng:                        &k6tv1.Rng{},
			BlockMultiQueue:            &enabled,
			NetworkInterfaceMultiQueue: &enabled,
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

	vm := k6tv1.VirtualMachine{}
	vm.Spec.Template = &tmpl
	k6tv1.SetObjectDefaults_VirtualMachine(&vm)
	return &vm, err
}

// NewVirtualMachine returns a fully zero-value VirtualMachine with all optional fields
func NewDefaultVirtualMachine() *k6tv1.VirtualMachine {
	vm, _ := NewVirtualMachine(MaxDisks, MaxIfaces, MaxPortsPerIface)
	return vm
}

func NewDefaultVirtualMachine2() *k6tv1.VirtualMachine {
	domSpec := k6tv1.DomainSpec{}
	// this is important. The reflect.Value must be addressable. You may want
	// to read carefully https://blog.golang.org/laws-of-reflection
	numItems := NumItems(map[string]int{
		"Disks":      int(MaxDisks),
		"Interfaces": int(MaxIfaces),
		"Ports":      int(MaxPortsPerIface),
		"NTPServers": int(MaxNTPServers),
	})
	makeStruct(reflect.TypeOf(domSpec), reflect.ValueOf(&domSpec).Elem(), numItems)

	tmpl := k6tv1.VirtualMachineInstanceTemplateSpec{}
	tmpl.Spec.Domain = domSpec

	vm := k6tv1.VirtualMachine{}
	vm.Spec.Template = &tmpl
	k6tv1.SetObjectDefaults_VirtualMachine(&vm)
	// workaround for k6t
	setObjectDefaults_VirtualMachine(&vm)
	return &vm
}

func setObjectDefaults_VirtualMachine(in *k6tv1.VirtualMachine) {
	if in.Spec.Template != nil {
		for i := range in.Spec.Template.Spec.Domain.Devices.Disks {
			a := &in.Spec.Template.Spec.Domain.Devices.Disks[i]
			if a.DiskDevice.CDRom != nil {
				setDefaults_CDRomTarget(a.DiskDevice.CDRom)
			}
		}
	}
}

func setDefaults_CDRomTarget(obj *k6tv1.CDRomTarget) {
	_true := true
	obj.ReadOnly = &_true
	if obj.Tray == "" {
		obj.Tray = k6tv1.TrayStateClosed
	}
}
