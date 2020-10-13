#!/bin/bash
{
RET_A=1
RET_B=1
echo "[test_id:5174]: VM with validations and deleted template"
echo "should successfully create a VM based on the VM validation rules"
$KUBECTL create -f manifests/17-vm-with-validations-no-template/vm-with-validation-rules-deleted-template-passing.yaml && RET_A=0
$KUBECTL delete vm vm-test-17

echo "[test_id:5046]: should fail to create a VM based on the VM validation rules"
$KUBECTL create -f manifests/17-vm-with-validations-no-template/vm-with-validation-rules-deleted-template-failing.yaml && RET_A=1

echo "[test_id:5175]: VM with validations without template"
echo "should successfully create a VM based on the VM validation rules"
$KUBECTL create -f manifests/17-vm-with-validations-no-template/vm-with-validation-rules-without-template-passing.yaml && RET_B=0
$KUBECTL delete vm vm-test-17

echo "[test_id:5047]: should fail to create a VM based on the VM validation rules"
$KUBECTL create -f manifests/17-vm-with-validations-no-template/vm-with-validation-rules-without-template-failing.yaml && RET_B=1

exit $((RET_A + RET_B))
}
