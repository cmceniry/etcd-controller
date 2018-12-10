#!/bin/bash

TESTNAME=etcd-controller-test-009
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

ctl_command /usr/local/bin/etcdctl --endpoints http://${TESTNET}.101:2379,http://${TESTNET}.102:2379,http://${TESTNET}.103:2379 endpoint status

echo SUCCESS
echo "---------- REMOVE NODES FROM LIST - CHECK THAT THEY LEFT AND CLUSTER STILL UP ----------"

generate_config 101
sleep 30

cs=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.100:4270 cstatus | awk '$2 == "RUNNING"' | wc -l)
rc=$?
if [ $rc -ne 0 ]; then fail -1 "FAIL: clusterstatus rc=$rc"; fi
if [ -z $cs ]; then fail -1 "FAIL: clusterstatus cs is empty"; fi
if [ $cs -ne 1 ]; then fail -1 "FAIL: clusterstatus not 1. got $cs"; fi

cs=$(ctl_command_result /usr/local/bin/etcdctl --endpoints http://${TESTNET}.101:2379,http://${TESTNET}.102:2379,http://${TESTNET}.103:2379 endpoint status 2>&1)
rc=$?
if [ $rc -eq 0 ]; then fail -1 "FAIL: clusterstatus succeed when should not. rc=$rc. $cs"; fi

ctl_command /usr/local/bin/etcdctl --endpoints http://${TESTNET}.101:2379 endpoint status

echo SUCCESS
echo "---------- ADDING NODES BACK ----------"

docker-compose up --force-recreate -d controller002 controller003 || fail -1 "unable to recreate 002 container"
sleep 10
generate_config 101 102 103
sleep 30

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
