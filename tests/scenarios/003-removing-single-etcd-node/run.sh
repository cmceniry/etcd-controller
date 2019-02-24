#!/bin/bash

TESTNAME=etcd-controller-test-003
cd `dirname $0`
. ../../common.sh

mkdir -p ./config

echo $TESTNAME

generate_config 101 102 103
docker-compose up -d || exit -1
sleep 30

echo "---------- UPDATE NODELIST AND GIVE TIME TO STABALIZE ----------"
generate_config 101 102
sleep 30

echo DONE
echo "---------- CHECK ETCD CLUSTER STATUS IS RUNNING ON 2 NODES ----------"
c=$(ctl_command_result /usr/local/bin/etcdctl --endpoints http://${TESTNET}.101:2379,http://${TESTNET}.102:2379 endpoint status)
rc=$?
if [ $rc -ne 0 ]; then fail -1 "FAIL: etcdctl rc=$rc"; fi

cs=$(ctl_command_result /usr/local/bin/etcdctl --endpoints http://${TESTNET}.101:2379,http://${TESTNET}.102:2379 endpoint status | grep 3.1.19 | wc -l)
rc=$?
if [ $rc -ne 0 ]; then fail -1 "FAIL: etcdctl count rc=$rc"; fi
if [ -z $cs ]; then fail -1 "FAIL: etcdctl count is empty"; fi
if [ $cs -ne 2 ]; then fail -1 "FAIL: etcdctl count is not 2. got $cs"; fi

cs=$(ctl_command_result /usr/local/bin/etcdctl --endpoints http://${TESTNET}.103:2379 endpoint status 2>&1)
rc=$?
if [ $rc -eq 0 ]; then fail -1 "FAIL: etcdctl 103 rc=$rc"; fi

echo PASS
echo "---------- CHECK GROUP STATUS IS RUNNING AND 3 NODES ----------"
cs=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.101:4270 cstatus | wc -l)
rc=$?
if [ $rc -ne 0 ]; then fail -1 "FAIL: clusterstatus rc=$rc"; fi
if [ -z $cs ]; then fail -1 "FAIL: clusterstatus cs is empty"; fi
if [ $cs -ne 3 ]; then fail -1 "FAIL: clusterstatus not 3. got $cs"; fi

cs=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.101:4270 cstatus | awk '$2 == "WATCHING"' | wc -l)
rc=$?
if [ $rc -ne 0 ]; then fail -1 "FAIL: clusterstatus watching rc=$rc"; fi
if [ -z $cs ]; then fail -1 "FAIL: clusterstatus watching cs is empty"; fi
if [ $cs -ne 1 ]; then fail -1 "FAIL: clusterstatus watching not 1. got $cs"; fi

echo PASS
echo "---------- CLEANUP ----------"

docker-compose down
rm -rf ./config

echo TEST SUCCESS
