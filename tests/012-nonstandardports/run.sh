#!/bin/bash

TESTNAME=etcd-controller-test-012
cd `dirname $0`
. ../common.sh

mkdir -p ./config

echo $TESTNAME

COMMAND_PORT=5270
SERF_PORT=5271
CLIENT_PORT=4379
PEER_PORT=4380

generate_config 101 102 103
sed -ie "s/2379/${CLIENT_PORT}/g; s/2380/${PEER_PORT}/g; s/4270/${COMMAND_PORT}/g" config/node-list.yaml
COMMAND_PORT=${COMMAND_PORT} SERF_PORT=${SERF_PORT} CLIENT_PORT=${CLIENT_PORT} PEER_PORT=${PEER_PORT} \
  docker-compose up -d || exit -1

echo "---------- CHECK CLUSTER STATUS - ENSURE THREE NODES WITH 101 AS CONDUCTOR ----------"
sleep 30

cs=$(ctl_command_result /etcd-controller-ctl ${TESTNET}.101:${COMMAND_PORT} cstatus | awk '$2 == "RUNNING"' | wc -l)
rc=$?
if [ $rc -ne 0 ]; then fail -1 "FAIL: clusterstatus rc=$rc"; fi
if [ -z $cs ]; then fail -1 "FAIL: clusterstatus cs is empty"; fi
if [ $cs -ne 3 ]; then fail -1 "FAIL: clusterstatus not 3. got $cs"; fi

ctl_command /usr/local/bin/etcdctl --endpoints http://${TESTNET}.101:${CLIENT_PORT},http://${TESTNET}.102:${CLIENT_PORT},http://${TESTNET}.103:${CLIENT_PORT} endpoint status

echo SUCCESS
echo "---------- CLEANUP ----------"

docker-compose down
rm -rf ./config
