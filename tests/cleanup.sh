#!/bin/sh

docker kill \
    etcd-controller-test-001-0 \
    etcd-controller-test-001-1

docker rm \
    etcd-controller-test-001-0 \
    etcd-controller-test-001-1