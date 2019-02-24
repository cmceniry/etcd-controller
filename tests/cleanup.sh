#!/bin/sh

cd $(dirname $0)

for d in `find . -type d -maxdepth 1 -not -path . -not -path ./scenarios`; do
    pushd $d
    docker-compose down
    popd
done

rm -f results-*.out