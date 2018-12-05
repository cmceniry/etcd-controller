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

function ctl_command_result() {
    docker-compose exec -T ctl "$@"
}

function fail() {
    echo $2
    exit $1
}

docker-compose up -d
sleep 30

n1=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.100:4270 nodestatus ${TESTNAME}-101)
rc=$?
if [ $? -ne 0 ]; then fail -1 "FAIL: nodestatus101. Expected 0. Got $rc: $n1"; fi
if [ "$n1" != "RUNNING" ]; then fail -1 "FAIL: nodestatus101. Not RUNNING: $n1"; fi

n2=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.100:4270 nodestatus ${TESTNAME}-102)
rc=$?
if [ $? -ne 0 ]; then fail -1 "FAIL: nodestatus102. Expected 0. Got $rc: $n1"; fi
if [ "$n2" != "RUNNING" ]; then fail -1 "FAIL: nodestatus102. Not RUNNING: $n2"; fi

n3=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.100:4270 nodestatus ${TESTNAME}-103)
rc=$?
if [ $? -ne 0 ]; then fail -1 "FAIL: nodestatus103. Expected 0. Got $rc: $n1"; fi
if [ "$n3" != "RUNNING" ]; then fail -1 "FAIL: nodestatus103. Not RUNNING: $n3"; fi

ctl_command /etcd-controller-ctl ${TESTNET}.100:4270 cstatus

cs=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.100:4270 cstatus | awk '$2 == "RUNNING"' | wc -l)
rc=$?
if [ $rc -ne 0 ]; then fail -1 "FAIL: clusterstatus rc=$rc"; fi
if [ -z $cs ]; then fail -1 "FAIL: clusterstatus cs is empty"; fi
if [ $cs -ne 3 ]; then fail -1 "FAIL: clusterstatus not 3. got $cs"; fi

docker-compose down
