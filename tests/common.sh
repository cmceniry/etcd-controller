
if [ "$1" == "-v" ]; then
    set -x
fi

TESTNET=172.27.0
export TESTNAME TESTNET

function ctl_command() {
    docker-compose exec ctl \
        "$@"
    rc=$?
    if [ $rc -ne 0 ]; then
        echo "FAIL: rc=$rc"
        exit -1
    fi
}

function ctl_command_result() {
    docker-compose exec -T ctl "$@"
}

function fail() {
    echo $2
    exit $1
}

function generate_config() {
    tmp=$(mktemp config/node-list.yaml.XXXXXX)
    echo "---" >>${tmp}
    for nodenumber in $@; do
        cat >>${tmp} <<EONS
- name: ${TESTNAME}-${nodenumber}
  IP: ${TESTNET}.${nodenumber}
  CommandPort: 4270
  Insecure: true
  PeerPort: 2380
  ClientPort: 2379
EONS
    done
    mv ${tmp} config/node-list.yaml
}


