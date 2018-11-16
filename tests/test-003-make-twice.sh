#!/bin/sh -x

TESTNAME=etcd-controller-test-003

echo "expect to fail"

function start_etcd_controller() {
    NAME=${1}
    IP=${2}
    docker run --rm -d \
        --name ${NAME} \
        --network atlas2_default \
        --ip ${IP} \
        -e ETCDCONTROLLER_NAME=${NAME} \
        -e ETCDCONTROLLER_IP=${IP} \
        etcd-controller:snapshot /etcd-controller
    return $?
}

function ctl_command() {
    docker run \
        --name ${TESTNAME}-0 \
        --network atlas2_default \
        --ip 172.18.1.1 -it \
        -e ETCDCTL_API=3 \
        etcd-controller:snapshot \
        "$@"
    rc=$(docker inspect ${TESTNAME}-0 --format='{{.State.ExitCode}}')
    if [ $rc -eq 0 ]; then
        echo "SUCCESS"
    else
        echo "FAIL: rc=$rc"
        exit -1
    fi
    docker rm ${TESTNAME}-0
}

start_etcd_controller ${TESTNAME}-1 172.18.100.1

ctl_command /etcd-controller-ctl 172.18.100.1:4270 init
sleep 5

ctl_command /etcd-controller-ctl 172.18.100.1:4270 init
sleep 5

docker kill \
    ${TESTNAME}-1
