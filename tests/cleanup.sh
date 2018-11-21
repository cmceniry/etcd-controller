#!/bin/sh

docker kill \
    etcd-controller-test-001-0 \
    etcd-controller-test-001-1 \
    etcd-controller-test-002-0 \
    etcd-controller-test-002-1 \
    etcd-controller-test-002-2 \
    etcd-controller-test-002-3 \
    etcd-controller-test-003-0 \
    etcd-controller-test-003-1 \
    etcd-controller-test-004-0 \
    etcd-controller-test-004-1 \
    etcd-controller-test-005-0 \
    etcd-controller-test-005-1 \
    etcd-controller-test-005-2 \
    etcd-controller-test-005-3
    
docker rm \
    etcd-controller-test-001-0 \
    etcd-controller-test-001-1 \
    etcd-controller-test-002-0 \
    etcd-controller-test-002-1 \
    etcd-controller-test-002-2 \
    etcd-controller-test-002-3 \
    etcd-controller-test-003-0 \
    etcd-controller-test-003-1 \
    etcd-controller-test-004-0 \
    etcd-controller-test-004-1 \
    etcd-controller-test-005-0 \
    etcd-controller-test-005-1 \
    etcd-controller-test-005-2 \
    etcd-controller-test-005-3
