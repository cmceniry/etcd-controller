#!/bin/bash

TESTNAME=etcd-controller-test-002
cd `dirname $0`
. ../common.sh

docker-compose up -d || exit -1
sleep 5

echo "---------- INITIALIZE SINGLE NODE AND CHECK STATUS ----------"

ctl_command /etcd-controller-ctl ${TESTNET}.101:4270 init
sleep 5
ctl_command /etcd-controller-ctl ${TESTNET}.101:4270 status
ctl_command /usr/local/bin/etcdctl \
  --endpoints http://${TESTNET}.101:2379 \
  endpoint status

echo SUCCESS
echo "---------- ADD SECOND NODE AND CHECK STATUS ----------"

ctl_command /usr/local/bin/etcdctl \
  --endpoints http://${TESTNET}.101:2379 \
  member add \
    --peer-urls http://${TESTNET}.102:2380 ${TESTNAME}-102
ctl_command /etcd-controller-ctl \
  ${TESTNET}.102:4270 join \
  ${TESTNAME}-101=http://${TESTNET}.101:2380,${TESTNAME}-102=http://${TESTNET}.102:2380
sleep 5
ctl_command /etcd-controller-ctl \
  ${TESTNET}.102:4270 status
ctl_command /usr/local/bin/etcdctl \
  --endpoints http://${TESTNET}.101:2379,http://${TESTNET}.102:2379 \
  endpoint status

echo SUCCESS
echo "---------- ADD THIRD NODE AND CHECK STATUS ----------"

ctl_command /usr/local/bin/etcdctl \
  --endpoints http://${TESTNET}.101:2379 \
  member add \
    --peer-urls http://${TESTNET}.103:2380 ${TESTNAME}-103
ctl_command /etcd-controller-ctl \
  ${TESTNET}.103:4270 join \
  ${TESTNAME}-101=http://${TESTNET}.101:2380,${TESTNAME}-102=http://${TESTNET}.102:2380,${TESTNAME}-103=http://${TESTNET}.103:2380
sleep 5
ctl_command /etcd-controller-ctl \
  ${TESTNET}.103:4270 status
ctl_command /usr/local/bin/etcdctl \
  --endpoints http://${TESTNET}.101:2379,http://${TESTNET}.102:2379,http://${TESTNET}.103:2379 \
  endpoint status

echo SUCCESS
echo "---------- CLEANUP ----------"

docker-compose down