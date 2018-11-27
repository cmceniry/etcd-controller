#!/bin/sh -x

TESTNAME=etcd-controller-test-006
TESTDIR=`pwd`/tests

echo "expect to fail"

function start_etcd_controller() {
    NAME=${1}
    IP=${2}
    if [ "$3" != "" ]; then
        MNTOPTION="-v ${3}:/config"
    fi
    docker run --rm -d \
        --name ${NAME} \
        --network atlas2_default \
        --ip ${IP} \
        -e ETCDCONTROLLER_NAME=${NAME} \
        -e ETCDCONTROLLER_IP=${IP} \
        ${MNTOPTION} \
        etcd-controller:snapshot /etcd-controller
    return $?
}

function run_command() {
    NAME=${1}
    shift
    IP=${1}
    shift
    if [ -d "$1" ]; then
        MNTOPTION="-v ${1}:/config"
        shift
    fi
    docker run --rm -d \
        --name ${NAME} \
        --network atlas2_default \
        --ip ${IP} \
        -e ETCDCTL_API=3 \
        ${MNTOPTION} \
        etcd-controller:snapshot \
        "$@"
    rc=$?
    if [ $rc -eq 0 ]; then
        echo "Start Successful"
    else
        echo "Start Failed: rc=$rc"
        exit -1
    fi
}

function ctl_command() {
    if [ -d "$1" ]; then
        MNTOPTION="-v ${1}:/config"
        shift
    fi
    docker run -it \
        --name ${TESTNAME}-0 \
        --network atlas2_default \
        --ip 172.18.1.1 \
        -e ETCDCTL_API=3 \
        ${MNTOPTION} \
        etcd-controller:snapshot \
        "$@"
}

function check_ctl_command() {
    ctl_command $@
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
start_etcd_controller ${TESTNAME}-2 172.18.100.2
start_etcd_controller ${TESTNAME}-3 172.18.100.3
run_command ${TESTNAME}-10 172.18.100.100 ${TESTDIR}/data/${TESTNAME} /etcd-controller-conductor

sleep 15
n1=$(ctl_command ${TESTDIR}/data/${TESTNAME} /etcd-controller-ctl 172.18.100.100:4270 nodestatus ${TESTNAME}-1)
echo ${n1}
docker rm ${TESTNAME}-0
n2=$(ctl_command ${TESTDIR}/data/${TESTNAME} /etcd-controller-ctl 172.18.100.100:4270 nodestatus ${TESTNAME}-2)
docker rm ${TESTNAME}-0
n3=$(ctl_command ${TESTDIR}/data/${TESTNAME} /etcd-controller-ctl 172.18.100.100:4270 nodestatus ${TESTNAME}-3)
docker rm ${TESTNAME}-0

docker kill \
    ${TESTNAME}-1 \
    ${TESTNAME}-2 \
    ${TESTNAME}-3 \
    ${TESTNAME}-10
