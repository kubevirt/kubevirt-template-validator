#!/bin/bash
{
RET_A=1
RET_B=1
echo "[test_id:5035]: Template with validations, VM with validations"
# Template expects max amount of CPU cores to be 4
$KUBECTL create -f manifests/16-template-with-validations-vm-with-validations/template-with-rules.yaml

# Creating a VM with a rule of maximum 6 cores and actual 5 cores, should pass
echo "[test_id:5036]: should successfully create a VM based on the VM validation rules"
$KUBECTL create -f manifests/16-template-with-validations-vm-with-validations/vm-with-validation-rules-passing.yaml && RET_A=0
$KUBECTL delete vm vm-test-16

# Creating a VM with a rule of maximum 4 cores and actual 5 cores, should fail
echo "[test_id:5035]: Template with validations, VM with validations"
echo "should fail to create a VM based on the VM validation rules"
$KUBECTL create -f manifests/16-template-with-validations-vm-with-validations/vm-with-validation-rules-failing.yaml || RET_B=0

$KUBECTL delete -f manifests/16-template-with-validations-vm-with-validations/template-with-rules.yaml

exit $((RET_A + RET_B))
}
