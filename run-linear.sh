#!/bin/sh

rm -rf result
mkdir -p result

trap "exit" INT
trap "kill 0" EXIT

num_types=3
num_traders=500
run_time=$((10*60))

function log {
    echo -e "\033[1;32m`date "+%Y-%m-%d %H:%M:%S"`\t$1\033[0m"
}

function run {
    result_file="result/$1.result"
    json_file="result/$1.json"
    log_file="result/$1.log"

    if [ -f $json_file ]
    then
        log "Loading from $json_file..."
        go run cmd/main.go -load-from=$json_file > $result_file 2> $log_file
        log "Loaded from $json_file."
    else
        alpha=$(echo "$1/100" | bc -l)
        log "Running with alpha=$1%..."
        go run cmd/main.go -type=$num_types -time=$run_time -trader=$num_traders -random=$num_traders -alpha=$alpha -save-to=$json_file > $result_file 2> $log_file
        log "Run with alpha=$1% finished."
    fi
}

for i in $(seq 0 1 100)
do
    run $i &
done

wait

log "All runs finished!"