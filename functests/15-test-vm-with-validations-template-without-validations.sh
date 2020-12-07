#!/bin/bash
{
RET_A=1
RET_B=1
echo "[test_id:5034]: Template without validations, VM with validations"
$KUBECTL create -f manifests/15-template-without-validations-vm-with-validations/template-without-rules.yaml

# Create a VM that should pass validation
echo "[test_id:5173]: should create a VM that passes validation"
$KUBECTL create -f manifests/15-template-without-validations-vm-with-validations/vm-with-validation-rules-passing.yaml && RET_A=0
$KUBECTL delete vm vm-test-15

# Attempt to create a VM that should fail validation
echo "[test_id:5034]: should fail to create VM that fails validation"
$KUBECTL create -f manifests/15-template-without-validations-vm-with-validations/vm-with-validation-rules-failing.yaml || RET_B=0

$KUBECTL delete -f manifests/15-template-without-validations-vm-with-validations/template-without-rules.yaml

exit $((RET_A + RET_B))
}
