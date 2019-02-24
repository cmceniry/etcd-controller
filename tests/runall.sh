#!/bin/bash

cd $(dirname $0)

failed=0
for d in `find . -type d -maxdepth 1 -not -path . -not -path ./scenarios | sed -e 's/^.\///'`; do
    printf "TEST:%-35s " $d
    ./$d/run.sh > results-${d}.out 2>&1
    if [ $? -eq 0 ]; then
        echo SUCCESS
    else
        failed=1
        echo FAILED
        echo "======================================================"
        cat results-$d.out
        echo "======================================================"
    fi
done

if [ $failed -ne 0 ]; then
    echo "AT LEAST ONE TEST FAILED"
    exit -1
fi
