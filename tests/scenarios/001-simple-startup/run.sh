#!/bin/bash

TESTNAME=etcd-controller-test-001
cd `dirname $0`
. ../../common.sh

mkdir -p ./config

echo $TESTNAME

generate_config
docker-compose up -d || exit -1
sleep 15

echo "---------- UPDATE NODELIST AND GIVE TIME TO STABALIZE ----------"
generate_config 101 102 103
sleep 30

echo DONE
echo "---------- CHECK ETCD CLUSTER STATUS IS RUNNING ON 3 NODES ----------"
c=$(ctl_command_result /usr/local/bin/etcdctl --endpoints http://${TESTNET}.101:2379,http://${TESTNET}.102:2379,http://${TESTNET}.103:2379 endpoint status)
rc=$?
if [ $rc -ne 0 ]; then fail -1 "FAIL: etcdctl rc=$rc"; fi

ces=$(ctl_command_result /usr/local/bin/etcdctl --endpoints http://${TESTNET}.101:2379,http://${TESTNET}.102:2379,http://${TESTNET}.103:2379 endpoint status | grep 3.1.19 | wc -l)
rc=$?
if [ $rc -ne 0 ]; then fail -1 "FAIL: etcdctl count rc=$rc"; fi
if [ -z $ces ]; then fail -1 "FAIL: etcdctl count is empty"; fi
if [ $ces -ne 3 ]; then fail -1 "FAIL: etcdctl count is not 3. got $cs"; fi

echo PASS
echo "---------- CHECK GROUP STATUS IS RUNNING AND 3 NODES ----------"
cs=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.101:4270 cstatus | awk '$2 == "RUNNING"' | wc -l)
rc=$?
if [ $rc -ne 0 ]; then fail -1 "FAIL: clusterstatus rc=$rc"; fi
if [ -z $cs ]; then fail -1 "FAIL: clusterstatus cs is empty"; fi
if [ $cs -ne 3 ]; then fail -1 "FAIL: clusterstatus not 3. got $cs"; fi

echo PASS
echo "---------- CHECK ALL NODES SAY ORCHESTRATOR IS 101 NODE ----------"
ci=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.101:4270 conductor)
rc=$?
if [ $rc -ne 0 ]; then fail -1 "FAIL: {TESTNET}.101 conductor command rc=$rc"; fi
if [ -z $ci ]; then fail -1 "FAIL: {TESTNET}.101 conductor response is empty"; fi
if [ $ci != "yes" ]; then fail -1 "FAIL: ${TESTNET}.101 is not conductor"; fi

ci=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.102:4270 conductor)
rc=$?
if [ $rc -ne 1 ]; then fail -1 "FAIL: {TESTNET}.102 conductor command rc=$rc"; fi
if [ -z $ci ]; then fail -1 "FAIL: {TESTNET}.102 conductor response is empty"; fi
if [ $ci == "yes" ]; then fail -1 "FAIL: ${TESTNET}.102 is conductor"; fi

ci=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.103:4270 conductor)
rc=$?
if [ $rc -ne 1 ]; then fail -1 "FAIL: {TESTNET}.103 conductor command rc=$rc"; fi
if [ -z $ci ]; then fail -1 "FAIL: {TESTNET}.103 conductor response is empty"; fi
if [ $ci == "yes" ]; then fail -1 "FAIL: ${TESTNET}.103 is conductor"; fi

echo PASS
echo "---------- CLEANUP ----------"

docker-compose down
rm -rf ./config

echo TEST SUCCESS
