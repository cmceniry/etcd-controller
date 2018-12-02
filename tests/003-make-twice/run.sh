#!/bin/sh -x

cd `dirname $0`

TESTNAME=etcd-controller-test-003
TESTNET=172.27.0

export TESTNAME TESTNET

docker-compose up -d
sleep 5

docker-compose exec ctl /etcd-controller-ctl ${TESTNET}.101:4270 init
rc=$?
if [ $rc -ne 0 ]; then
    echo "SETUP FAIL: rc=$rc"
    exit -1
fi
sleep 5

docker-compose exec ctl /etcd-controller-ctl ${TESTNET}.101:4270 init
rc=$?
if [ $rc -eq 0 ]; then
    echo "EXPECTED FAILURE, GOT: $rc"
    exit -1
else
    echo "SUCCESS"
fi

docker-compose down