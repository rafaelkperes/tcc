#!/bin/bash

set -euo pipefail

RESULTS_DIR="./results"
MAX_INT64=9223372036854775807

p="${PRODUCER_BIN:-"./bin/producer"} -r 100 -i 0 -l 1000000"
[[ ! -z ${CONSUMER_ENDPOINT:-} ]] && p="./bin/producer -c ${CONSUMER_ENDPOINT} -r 100 -i 0 -l 1000000"

mkdir -p ${RESULTS_DIR:-"./results"}

formats=("application/json" "application/x-protobuf" "application/x-avro" "application/x-msgpack")
for format in ${formats[*]}; do
    printf "Running %s tests\n" $format
    pjson="$p -f $format"

    type="int"
    printf "Running %s.%s tests\n" $format $type
    # MAX_INT64=9223372036854775807
    logfile=${RESULTS_DIR:-"./results"}/$type.${format///}.log
    $pjson -t $type -intmin 0 -intmax 9223372036854775807 2>&1 | tee $logfile | jq -M -c 'select(.event=="progress") | .msg'

    type="float"
    printf "Running %s.%s tests\n" $format $type
    logfile=${RESULTS_DIR:-"./results"}/$type.${format///}.log
    $pjson -t $type 2>&1 | tee $logfile | jq -M -c 'select(.event=="progress") | .msg'

    printf "Running %s.%s tests\n" $format $type
    $pjson -t $type 2>&1 | tee $logfile | jq -M -c 'select(.event=="progress") | .msg'

    printf "Running %s.%s tests\n" $format $type
    $pjson -t $type -strlen 100 2>&1 | tee $logfile | jq -M -c 'select(.event=="progress") | .msg'
done


{"event":"progress","level":"debug","msg":"sending request 1/1","time":"2019-10-06T13:11:04-03:00"}