#!/bin/bash

TESTNAME=etcd-controller-test-001
cd `dirname $0`
. ../common.sh

docker-compose up -d || exit -1
sleep 5

echo "---------- INITIALIZE SINGLE NODE ----------"

ctl_command /etcd-controller-ctl ${TESTNET}.101:4270 init

echo SUCCESS
sleep 5
echo "---------- CHECK STATUS ----------"

ctl_command /etcd-controller-ctl ${TESTNET}.101:4270 status

ctl_command /usr/local/bin/etcdctl \
  --endpoints http://${TESTNET}.101:2379 \
  endpoint status

echo SUCCESS
echo "---------- CLEANUP ----------"

docker-compose down
