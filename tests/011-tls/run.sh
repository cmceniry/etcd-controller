#!/bin/bash

TESTNAME=etcd-controller-test-011
cd `dirname $0`
. ../common.sh

mkdir -p ./config

echo $TESTNAME

go run generate_certs.go ${TESTNET} || exit -1
generate_config 101 102 103
docker-compose up -d || exit -1

echo "---------- CHECK CLUSTER STATUS - ENSURE THREE NODES WITH 101 AS CONDUCTOR ----------"
sleep 30

cs=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.101:4270 cstatus | awk '$2 == "RUNNING"' | wc -l)
rc=$?
if [ $rc -ne 0 ]; then fail -1 "FAIL: clusterstatus rc=$rc"; fi
if [ -z $cs ]; then fail -1 "FAIL: clusterstatus cs is empty"; fi
if [ $cs -ne 3 ]; then fail -1 "FAIL: clusterstatus not 3. got $cs"; fi

ctl_command /usr/local/bin/etcdctl --endpoints https://${TESTNET}.101:2379,https://${TESTNET}.102:2379,https://${TESTNET}.103:2379 endpoint status

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

echo SUCCESS
echo "---------- CLEANUP ----------"

docker-compose down
rm -rf ./config
