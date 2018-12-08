#!/bin/bash

TESTNAME=etcd-controller-test-007
cd `dirname $0`
. ../common.sh

mkdir -p ./config

echo $TESTNAME

generate_config 101 102 103
docker-compose up -d || exit -1
sleep 30

echo "---------- CHECK CLUSTER STATUS - ENSURE THREE NODES ----------"

cs=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.100:4270 cstatus | awk '$2 == "RUNNING"' | wc -l)
rc=$?
if [ $rc -ne 0 ]; then fail -1 "FAIL: clusterstatus rc=$rc"; fi
if [ -z $cs ]; then fail -1 "FAIL: clusterstatus cs is empty"; fi
if [ $cs -ne 3 ]; then fail -1 "FAIL: clusterstatus not 3. got $cs"; fi

echo SUCCESS
echo "---------- RESTART CONDUCTOR, CONFIRM CLUSTER STATUS ----------"

docker-compose restart conductor || fail -1 "FAIL: conductor restart failed"
sleep 10

cs=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.100:4270 cstatus | awk '$2 == "RUNNING"' | wc -l)
rc=$?
if [ $rc -ne 0 ]; then fail -1 "FAIL: clusterstatus rc=$rc"; fi
if [ -z $cs ]; then fail -1 "FAIL: clusterstatus cs is empty"; fi
if [ $cs -ne 3 ]; then fail -1 "FAIL: clusterstatus not 3. got $cs"; fi

echo SUCCESS
echo "---------- CONFIRM ETCD CLUSTER HEALTHY ----------"

ctl_command /usr/local/bin/etcdctl --endpoints http://${TESTNET}.101:2379,http://${TESTNET}.102:2379,http://${TESTNET}.103:2379 endpoint status

echo SUCCESS
echo "---------- CLEANUP ----------"

docker-compose down
rm -rf ./config
