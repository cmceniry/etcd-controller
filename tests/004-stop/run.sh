#!/bin/sh -x

cd `dirname $0`

TESTNAME=etcd-controller-test-004
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
sleep 5

ctl_command /etcd-controller-ctl ${TESTNET}.101:4270 init
sleep 5

ctl_command /etcd-controller-ctl ${TESTNET}.101:4270 status
sleep 5

ctl_command /etcd-controller-ctl ${TESTNET}.101:4270 stop
sleep 5

docker-compose exec ctl /etcd-controller-ctl ${TESTNET}.101:4270 stop
rc=$?
if [ $rc -eq 0 ]; then
    echo "EXPECTED FAILURE, GOT: $rc"
    exit -1
else
    echo "SUCCESS"
fi

docker-compose down
