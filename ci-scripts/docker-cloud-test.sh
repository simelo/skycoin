#!/bin/bash

export PORT=6420
cd ./..
DIR=$PWD

sleep 60

curl http://integration-test-stable:6420/api/v1/version


# 'integration-test-stable'
TESTS=(
    'integration-test-stable-disable-csrf'
    'integration-test-disable-wallet-api'
    'integration-test-enable-seed-api'
    'integration-test-disable-gui'
    'integration-test-auth'
    'integration-test-db-no-unconfirmed'
)

for TEST in ${TESTS[@]} ; do
    echo "----- START TEST: $TEST -----"
    export SKYCOIN_NODE=$TEST
    make $TEST
    FAIL=$?
    if [ $FAIL -ne 0 ]; then
        echo "----- FAIL TEST: $TEST -----"
        cat /tmp/my_output
        exit 1
    fi
    echo "----- PASS TEST: $TEST -----"

    cd /wallet
    rm -r -f `ls`
    # cd /data/.skycoin
    # rm -r -f `ls`
    cd $DIR
done