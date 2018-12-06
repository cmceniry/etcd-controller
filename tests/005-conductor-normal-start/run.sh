#!/bin/bash

TESTNAME=etcd-controller-test-005
cd `dirname $0`
. ../common.sh

mkdir -p ./config
generate_config 101 102 103

docker-compose up -d || exit -1

echo "---------- LONG WAIT FOR EVERYTHING TO COME UP ----------"

sleep 30

echo SUCCESS
echo "---------- CHECK NODE ALL THREE NODE STATUSES ----------"

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

echo SUCCESS
echo "---------- CHECK NODE ALL THREE NODE STATUSES ----------"

ctl_command /etcd-controller-ctl ${TESTNET}.100:4270 cstatus
cs=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.100:4270 cstatus | awk '$2 == "RUNNING"' | wc -l)
rc=$?
if [ $rc -ne 0 ]; then fail -1 "FAIL: clusterstatus rc=$rc"; fi
if [ -z $cs ]; then fail -1 "FAIL: clusterstatus cs is empty"; fi
if [ $cs -ne 3 ]; then fail -1 "FAIL: clusterstatus not 3. got $cs"; fi

echo SUCCESS
echo "---------- CLEANUP ----------"

docker-compose down
rm -rf ./config
