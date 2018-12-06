#!/bin/bash

TESTNAME=etcd-controller-test-003
cd `dirname $0`
. ../common.sh

docker-compose up -d || exit -1
sleep 5

echo "---------- INITIALIZE SINGLE NODE ----------"

ctl_command /etcd-controller-ctl ${TESTNET}.101:4270 init
sleep 5
ctl_command /etcd-controller-ctl ${TESTNET}.101:4270 status

echo SUCCESS
echo "---------- ATTEMPT TO INIT AGAIN - CHECK THAT IT FAILS ----------"

ctl_command_result /etcd-controller-ctl ${TESTNET}.101:4270 init
rc=$?
if [ $rc -eq 0 ]; then fail -1 "expected failure. got $rc"; fi

echo SUCCESS
echo "---------- CLEANUP ----------"

docker-compose down