#!/bin/sh -x

cd `dirname $0`

TESTNAME=etcd-controller-test-002
TESTNET=172.27.0

export TESTNAME TESTNET

function ctl_command() {
    docker-compose exec ctl \
        "$@"
    rc=$?
    if [ $rc -eq 0 ]; then
        echo "SUCCESS"
    else
        echo "FAIL: rc=$rc"
        exit -1
    fi
}


docker-compose up -d
sleep 5

ctl_command /etcd-controller-ctl ${TESTNET}.101:4270 init
sleep 5

ctl_command /etcd-controller-ctl ${TESTNET}.101:4270 status
sleep 5

ctl_command /usr/local/bin/etcdctl \
  --endpoints ${TESTNET}.101:2379 \
  member add \
    --peer-urls http://${TESTNET}.102:2380 ${TESTNAME}-102
ctl_command /etcd-controller-ctl \
  ${TESTNET}.102:4270 join \
  ${TESTNAME}-101=http://${TESTNET}.101:2380,${TESTNAME}-102=http://${TESTNET}.102:2380
sleep 5

ctl_command /etcd-controller-ctl \
  ${TESTNET}.102:4270 status
sleep 5

ctl_command /usr/local/bin/etcdctl \
  --endpoints ${TESTNET}.101:2379 \
  member add \
    --peer-urls http://${TESTNET}.103:2380 ${TESTNAME}-103
ctl_command /etcd-controller-ctl \
  ${TESTNET}.103:4270 join \
  ${TESTNAME}-101=http://${TESTNET}.101:2380,${TESTNAME}-102=http://${TESTNET}.102:2380,${TESTNAME}-103=http://${TESTNET}.103:2380
sleep 5

ctl_command /etcd-controller-ctl \
  ${TESTNET}.103:4270 status
sleep 5

ctl_command /usr/local/bin/etcdctl \
  --endpoints http://${TESTNET}.101:2379,http://${TESTNET}.102:2379,http://${TESTNET}.103:2379 \
  endpoint status

docker-compose down