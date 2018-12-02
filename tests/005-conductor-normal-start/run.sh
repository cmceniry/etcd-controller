#!/bin/sh -x

cd `dirname $0`

TESTNAME=etcd-controller-test-005
TESTNET=172.27.0

export TESTNAME TESTNET

function ctl_command() {
    docker-compose exec ctl \
        "$@"
    rc=$?
    if [ $rc -eq 0 ]; then
        echo "SUCCESS"
    else
        echo "FAIL: rc=$rc"
        exit -1
    fi
}


docker-compose up -d
sleep 30

ctl_command /etcd-controller-ctl ${TESTNET}.100:4270 nodestatus ${TESTNAME}-101
ctl_command /etcd-controller-ctl ${TESTNET}.100:4270 nodestatus ${TESTNAME}-102
ctl_command /etcd-controller-ctl ${TESTNET}.100:4270 nodestatus ${TESTNAME}-103

docker-compose down