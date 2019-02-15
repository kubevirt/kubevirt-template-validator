package validation_test

import (
	"bytes"
	"k8s.io/apimachinery/pkg/util/yaml"

	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

func NewVMCirros() *k6tv1.VirtualMachine {
	vm := k6tv1.VirtualMachine{}
	b := bytes.NewBufferString(`
apiVersion: kubevirt.io/v1alpha3
kind: VirtualMachine
metadata:
  labels:
    kubevirt.io/vm: vm-cirros
  name: vm-cirros
spec:
  running: false
  template:
    metadata:
      labels:
        kubevirt.io/vm: vm-cirros
    spec:
      domain:
        devices:
          disks:
          - disk:
              bus: virtio
            name: containerdisk
          - disk:
              bus: virtio
            name: cloudinitdisk
        machine:
          type: ""
        resources:
          requests:
            memory: 64M
      terminationGracePeriodSeconds: 0
      volumes:
      - containerDisk:
          image: registry:5000/kubevirt/cirros-container-disk-demo:devel
        name: containerdisk
      - cloudInitNoCloud:
          userData: |
            #!/bin/sh

            echo 'printed from cloud-init userdata'
        name: cloudinitdisk`)
	decoder := yaml.NewYAMLOrJSONDecoder(b, 1024) // FIXME explain magic number
	err := decoder.Decode(&vm)
	if err != nil {
		panic(err)
	}
	return &vm
}
