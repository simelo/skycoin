#!/usr/bin/env bash

# Runs a local Skycoin testnet

# Disclaimer : For the moment this script is a hack used for testing purposes
#              Proper implementation consists in launching/connecting to (sim)testnet fiber coin network

set -x

if [ -z "$1" ]; then
  exit "Missing argument"
fi

PORTS=$(echo "$1
$(expr $1 + 1)
$(expr $1 + 2)
$(expr $1 + 3)
$(expr $1 + 4)
$(expr $1 + 5)
$(expr $1 + 6)")

TEMP_DIR="/tmp/skytestnet.$1"
PEERS_FILE="${TEMP_DIR}/localhost-peers.txt"
PID_FILE="${TEMP_DIR}/pid.txt"

echo "Creating temp dirs starting at $1"
echo "${PORTS}" | tr , '\n' | xargs -I PORT mkdir -p "${TEMP_DIR}/PORT" && \
  echo "..................................... [OK]"

echo "Creating local peers file"
echo "${PORTS}" | sed 's/^/127.0.0.1:/g' > ${PEERS_FILE} && \
  cat ${PEERS_FILE} && \
  echo "........................... [OK]"

echo "Launching Skycoin nodes"
echo "" > ${PID_FILE} && \
  echo "${PORTS}" | while read PORT; do \
    screen -dmS "skytest.$(echo ${PORT})" /bin/bash -c "./run-client.sh -localhost-only -custom-peers-file=$(echo ${TEMP_DIR})/localhost-peers.txt -download-peerlist=false -launch-browser=false -data-dir=$(echo ${TEMP_DIR})/$(echo ${PORT}) -web-interface-port=\$(expr $(echo ${PORT}) + 420) -port=$(echo ${PORT}) | sed 's/^/skytest.$(echo ${PORT}) I /g'" ; \
  done && \
  echo "${PORTS}" | sed 's/^/skytest./g' > ${PID_FILE} && \
  echo "........................... [OK]"


