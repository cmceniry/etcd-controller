#!/bin/sh

cd $(dirname $0)

for d in `find . -type d -maxdepth 1 -not -path .`; do
    pushd $d
    docker-compose down
    popd
done
