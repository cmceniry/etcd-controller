#!/bin/bash

TESTNAME=etcd-controller-test-008
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
echo "---------- REMOVE NODE FROM LIST - ENSURE STILL RUNNING ----------"

generate_config 101 103
sleep 10

cs=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.100:4270 cstatus | awk '$2 == "RUNNING"' | wc -l)
rc=$?
if [ $rc -ne 0 ]; then fail -1 "FAIL: clusterstatus rc=$rc"; fi
if [ -z $cs ]; then fail -1 "FAIL: clusterstatus cs is empty"; fi
if [ $cs -ne 3 ]; then fail -1 "FAIL: clusterstatus not 3. got $cs"; fi

ctl_command /usr/local/bin/etcdctl --endpoints http://${TESTNET}.101:2379,http://${TESTNET}.102:2379,http://${TESTNET}.103:2379 endpoint status

echo SUCCESS
echo "---------- REMOVE STOPPED NODE ----------"

docker-compose stop controller002
sleep 30

cs=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.100:4270 cstatus | wc -l)
rc=$?
if [ $rc -ne 0 ]; then fail -1 "FAIL: clusterstatus rc=$rc"; fi
if [ -z $cs ]; then fail -1 "FAIL: clusterstatus cs is empty"; fi
if [ $cs -ne 2 ]; then fail -1 "FAIL: clusterstatus not 2. got $cs"; fi

cs=$(ctl_command_result /usr/local/bin/etcdctl --endpoints http://${TESTNET}.101:2379,http://${TESTNET}.102:2379,http://${TESTNET}.103:2379 endpoint status 2>&1)
rc=$?
if [ $rc -eq 0 ]; then fail -1 "FAIL: clusterstatus succeed when should not. rc=$rc. $cs"; fi

ctl_command /usr/local/bin/etcdctl --endpoints http://${TESTNET}.101:2379,http://${TESTNET}.103:2379 endpoint status

echo SUCCESS
echo "---------- ADDING NODE BACK ----------"

docker-compose up --force-recreate -d controller002 || fail -1 "unable to recreate 002 container"
sleep 30
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
