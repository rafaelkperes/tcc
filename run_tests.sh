#!/bin/bash

set -euo pipefail

RESULTS_DIR="./results"
MAX_INT64=9223372036854775807

p="${PRODUCER_BIN:-"./bin/producer"} -r 1000 -i 0 -l 1000000"
[[ ! -z ${CONSUMER_ENDPOINT:-} ]] && p="./bin/producer -c ${CONSUMER_ENDPOINT} -r 100 -i 0 -l 1000000"

mkdir -p ${RESULTS_DIR:-"./results"}

formats=("application/json" "application/x-protobuf" "application/x-avro" "application/x-msgpack")
types=("int" "float" "string" "object")
for format in ${formats[*]}; do
    for type in ${types[*]}; do
        printf "Running %s.%s tests\n" $type ${format//'application/'}

        pjson="$p -f $format"
        logfile=${RESULTS_DIR:-"./results"}/$type.${format//'application/'}.log
        $pjson -t $type -intmin 0 -intmax 9223372036854775807 -strlen 100 2>&1 | tee $logfile | jq -M -c 'select(.event=="progress") | .msg'
    done
done
