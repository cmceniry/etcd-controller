#!/bin/sh -x

docker run --rm -d \
    --name etcd-controller-test-001-1 \
    --network atlas2_default \
    --ip 172.18.100.1 \
    -e ETCDCONTROLLER_IP=172.18.100.1 \
    etcd-controller:snapshot /etcd-controller
sleep 5

docker run --name etcd-controller-test-001-0 --network atlas2_default --ip 172.18.1.1 -it \
    etcd-controller:snapshot \
    /etcd-controller-ctl 172.18.100.1:4270 init
rc=$(docker inspect etcd-controller-test-001-0 --format='{{.State.ExitCode}}')
if [ $rc -eq 0 ]; then
    echo "SUCCESS"
else
    echo "FAIL: rc=$rc"
    exit -1
fi
docker rm etcd-controller-test-001-0

sleep 5
docker run --name etcd-controller-test-001-0 --network atlas2_default --ip 172.18.1.1 -it \
    etcd-controller:snapshot \
    /etcd-controller-ctl 172.18.100.1:4270 status
rc=$(docker inspect etcd-controller-test-001-0 --format='{{.State.ExitCode}}')
if [ $rc -eq 0 ]; then
    echo "SUCCESS"
else
    echo "FAIL: rc=$rc"
    exit -1
fi
docker rm etcd-controller-test-001-0

docker kill \
    etcd-controller-test-001-1
