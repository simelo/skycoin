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
$(expr $1 + 6)
$(expr $1 + 7)
$(expr $1 + 8)
$(expr $1 + 9)")

PORTS_PEERS_SEED=$(echo "$PORTS" | head -n 3)
PORTS_PEERS_LATEST=$(echo "$PORTS" | head -n 6)
PORTS_PEERS_STABLE=$(echo "$PORTS" | tail -n 4 | head -n 2)
PORTS_PEERS_PREV=$(echo "$PORTS" | tail -n 2)

TEMP_DIR="/tmp/skytestnet.$1"
PEERS_FILE="${TEMP_DIR}/localhost-peers.txt"
PID_FILE="${TEMP_DIR}/pid.txt"

VERSION_STABLE=0.24.1
VERSION_LEGACY=0.23.0

echo "Creating temp dirs at ${TEMP_DIR}"
echo "${PORTS}" | xargs -I PORT mkdir -p "${TEMP_DIR}/PORT" && \
  echo "........................................ [OK]"

echo "Creating local peers file"
echo "${PORTS_PEERS_SEED}" | head -n 3 | sed 's/^/127.0.0.1:/g' > ${PEERS_FILE} && \
  cat ${PEERS_FILE} && \
  echo "........................................ [OK]"

echo "Launching Skycoin nodes from working copy"
echo "" > ${PID_FILE} && \
  echo "${PORTS_PEERS_LATEST}" | while read PORT; do \
    screen -dmS "skytest.$(echo ${PORT})" /bin/bash -c "./run-client.sh -localhost-only -custom-peers-file=$(echo ${TEMP_DIR})/localhost-peers.txt -download-peerlist=false -launch-browser=false -data-dir=$(echo ${TEMP_DIR})/$(echo ${PORT}) -web-interface-port=\$(expr $(echo ${PORT}) + 420) -port=$(echo ${PORT}) | sed 's/^/skytest.$(echo ${PORT}) I /g'" ; \
  done && \
  echo "${PORTS_PEERS_LATEST}" | sed 's/^/skytest./g' > ${PID_FILE} && \
  echo "........................................ [OK]"

# TODO : Run legacy Skycoin versions inside Docker containers
#echo "Launching Skycoin ${VERSION_STABLE} nodes in Docker containers"
#  echo "${PORTS_PEERS_STABLE}" | while read PORT; do \
#  docker run --name "skydocker.$(echo ${PORT})" skycoin/skycoin:release-v${VERSION_STABLE} -localhost-only -custom-peers-file=/data/localhost-peers.txt -download-peerlist=false -launch-browser=false ; \
## FIXME: Volumes and log to shared file
#  done && \
#  echo "${PORTS_PEERS_STABLE}" | sed 's/^/skydocker./g' >> ${PID_FILE} && \
#  echo "........................................ [OK]"
#
#echo "Launching Skycoin ${VERSION_LEGACY} nodes in Docker containers"
#  echo "${PORTS_PEERS_LEGACY}" | while read PORT; do \
#  docker run --name "skydocker.$(echo ${PORT})" -localhost-only -custom-peers-file=/data/localhost-peers.txt -download-peerlist=false -launch-browser=false ; \
## FIXME: Volumes and log to shared file
#  done && \
#  echo "${PORTS_PEERS_LEGACY}" | sed 's/^/skydocker./g' >> ${PID_FILE} && \
#  echo "........................................ [OK]"

