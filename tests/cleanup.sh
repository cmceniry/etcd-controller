#!/bin/sh

docker kill \
    etcd-controller-test-001-0 \
    etcd-controller-test-001-1 \
    etcd-controller-test-002-0 \
    etcd-controller-test-002-1 \
    etcd-controller-test-002-2 \
    etcd-controller-test-002-3

docker rm \
    etcd-controller-test-001-0 \
    etcd-controller-test-001-1 \
    etcd-controller-test-002-0 \
    etcd-controller-test-002-1 \
    etcd-controller-test-002-2 \
    etcd-controller-test-002-3