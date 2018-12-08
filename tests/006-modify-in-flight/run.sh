#!/bin/bash

TESTNAME=etcd-controller-test-006
cd `dirname $0`
. ../common.sh

mkdir -p ./config

generate_config #none

docker-compose up -d || exit -1
sleep 5

echo "---------- ENSURING NO NODES IN CSTATUS ----------"

cs=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.100:4270 cstatus | wc -l)
rc=$?
if [ $rc -ne 0 ]; then fail -1 "FAIL: clusterstatus rc=$rc"; fi
if [ -z $cs ]; then fail -1 "FAIL: clusterstatus is empty"; fi
if [ $cs -ne 0 ]; then fail -1 "FAIL: clusterstatus not 0. got $cs"; fi

echo "SUCCESS"
echo "---------- UPDATING NODE LIST TO INIT/START A NEW CLUSTER ----------"

generate_config 101 102 103
sleep 20

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

ctl_command /usr/local/bin/etcdctl --endpoints http://${TESTNET}.101:2379,http://${TESTNET}.102:2379,http://${TESTNET}.103:2379 endpoint status

echo SUCCESS
echo "---------- STOPPING A NODE ###AND DETECTING IT AS DOWN ----------"

ctl_command /etcd-controller-ctl ${TESTNET}.102:4270 stop

# cs=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.100:4270 cstatus | awk '$2 == "RUNNING"' | wc -l)
# rc=$?
# if [ $rc -ne 0 ]; then fail -1 "FAIL: clusterstatus rc=$rc"; fi
# if [ -z $cs ]; then fail -1 "FAIL: clusterstatus cs is empty"; fi
# if [ $cs -ne 2 ]; then fail -1 "FAIL: clusterstatus not 2. got $cs"; fi

echo SUCCESS
sleep 10
echo "---------- ENSURING THAT IT IS ALREADY RESTARTED BY CONDUCTOR ---------"

cs=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.100:4270 cstatus | awk '$2 == "RUNNING"' | wc -l)
rc=$?
if [ $rc -ne 0 ]; then fail -1 "FAIL: clusterstatus rc=$rc"; fi
if [ -z $cs ]; then fail -1 "FAIL: clusterstatus cs is empty"; fi
if [ $cs -ne 3 ]; then fail -1 "FAIL: clusterstatus not 3. got $cs"; fi

ctl_command /usr/local/bin/etcdctl --endpoints http://${TESTNET}.101:2379,http://${TESTNET}.102:2379,http://${TESTNET}.103:2379 endpoint status

echo SUCCESS
echo "---------- CLEANUP ----------"

docker-compose down
rm -rf ./config
