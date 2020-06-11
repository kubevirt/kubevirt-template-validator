#!/bin/bash

# This test only works with OKD cluster
{
echo '[test_id:4375] Test refreshing of certificates'

NAMESPACE=${NAMESPACE:-kubevirt}
SECRET_NAME=${SECRET_NAME:-"virt-template-validator-certs"}

PORT_FORWARD_PID=''
TMP_DIR=$(mktemp -dt validator-certs-XXXXXX)

function stopPortForward {
  if [ -n "${PORT_FORWARD_PID}" ] ; then
    kill ${PORT_FORWARD_PID}
    wait ${PORT_FORWARD_PID}
    PORT_FORWARD_PID=''
  fi
}

function cleanup() {
  rm -rf ${TMP_DIR}
  stopPortForward
}
trap cleanup EXIT

function writeCertToFile() {
  local LOCAL_PORT=48443

  # Start a background process to do port forwarding
  oc port-forward -n kubevirt service/virt-template-validator ${LOCAL_PORT}:443 >/dev/null &
  PORT_FORWARD_PID=$!

  # Wait for port to be opened
  while true ; do
    if nc -z localhost ${LOCAL_PORT} ; then
      break
    fi
    if ! ps -o pid= | grep -q ${PORT_FORWARD_PID} ; then
      echo "[ERROR] oc port-forward failed."
      exit 2
    fi
    sleep 1s
  done

  openssl s_client -connect "localhost:${LOCAL_PORT}" </dev/null 2>/dev/null | openssl x509 > $1
  stopPortForward
}

# Get the current cerfiticate used by the service
writeCertToFile "${TMP_DIR}/old.cert"

# Delete the secret so it will be recreated
$KUBECTL delete secret -n ${NAMESPACE} ${SECRET_NAME}

# It takes a while until the new secret is mouned to the pod.
# This loop times out after 5 minutes
for i in {1..30} ; do
  echo "Waiting for certificate update."
  sleep 10s

  # Get the new certificate and compare it to the old one
  writeCertToFile "${TMP_DIR}/new.cert"
  if ! cmp "${TMP_DIR}/old.cert" "${TMP_DIR}/new.cert" >/dev/null 2>&1; then
    exit 0
  fi
done

echo "[ERROR] Timed out waiting for certificate update"
exit 1
}
