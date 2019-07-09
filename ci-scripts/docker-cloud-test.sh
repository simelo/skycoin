#!/bin/bash

export PORT=6420
cd ./..
DIR=$PWD

sleep 60

curl http://integration-test-stable:6420/api/v1/version

curl http://integration-test-disable-wallet-api:6420/api/v1/version
 

# 'integration-test-db-no-unconfirmed'
# 'integration-test-auth'
TESTS=(
    'integration-test-stable'
    'integration-test-stable-disable-csrf'
    'integration-test-disable-wallet-api'
    'integration-test-enable-seed-api'
    'integration-test-disable-gui'
)

for TEST in ${TESTS[@]} ; do
    echo "----- START TEST: $TEST -----"
    if [ -d /wallet ]; then
        rm /wallet
    fi

    if [ -d /data/.skycoin ]; then
        rm /data/.skycoin
    fi
    ln -s /data/.skycoin-$TEST /data/.skycoin
    
    if [ $TEST != 'integration-test-disable-wallet-api' ]; then
        ln -s /wallet-$TEST /wallet
    fi

    
    export SKYCOIN_NODE=$TEST
    make $TEST
    FAIL=$?
    if [ $FAIL -ne 0 ]; then
        echo "----- FAIL TEST: $TEST -----"
        echo "----- MY OUTPUT -----"
        cat /tmp/my_output
        echo "---------------------"
        exit 1
    fi
    echo "----- PASS TEST: $TEST -----"

    #cd /wallet
    #rm -r -f `ls`
    # cd /data/.skycoin
    # rm -r -f `ls`
    #cd $DIR
done
