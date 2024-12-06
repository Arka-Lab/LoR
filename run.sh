#!/bin/sh

if [ "$1" == "cleanup" ]
then
    rm -rf result
    mkdir -p result
fi

trap "exit" INT
trap "kill 0" EXIT

num_types=3
num_traders=500
run_time=$((10*60))

function log {
    echo -e "\033[1;32m`date "+%Y-%m-%d %H:%M:%S"`\t$1\033[0m"
}

function run {
    num_random=$(($num_traders*$1/100))
    num_bad=$(($num_traders*$2/100))
    
    result_file="result/$1-$2.result"
    json_file="result/$1-$2.json"
    log_file="result/$1-$2.log"

    if [ -f $json_file ]
    then
        log "Loading from $json_file..."
        go run cmd/main.go -load-from=$json_file > $result_file 2> $log_file
        log "Loaded from $json_file."
    else
        log "Running with $1% random traders and $2% bad traders..."
        go run cmd/main.go -type=$num_types -time=$run_time -trader=$num_traders -random=$num_random -bad=$num_bad -save-to=$json_file > $result_file 2> $log_file
        log "Run with $1% random traders and $2% bad traders finished."
    fi
}

for i in $(seq 0 5 100)
do
    for j in $(seq 0 5 $((100-i)))
    do
        run $i $j
    done
done

wait

rm -rf output output.zip && mkdir -p output
cp result/*.result output/
zip -r output.zip output && rm -rf output
log "Output saved to output.zip."

rm -rf backup backup.zip && mkdir -p backup
cp result/*.json backup/
zip -r backup.zip backup && rm -rf backup
log "Backup saved to backup.zip."

log "All runs finished!"