#!/bin/sh

rm -rf result
mkdir -p result

trap "exit" INT
trap "kill 0" EXIT

num_types=3
num_traders=500
run_time=$((10*60))

function run {
    num_random=$(($num_traders*$1/100))
    num_bad=$(($num_traders*$2/100))
    file_name="result/$1-$2.json"

    echo "Running with $1% random traders and $2% bad traders..."
    go run cmd/main.go -type=$num_types -time=$run_time -trader=$num_traders -random=$num_random -bad=$num_bad -file=$file_name > result/$1-$2.result 2> result/$1-$2.log
    echo "Run with $1% random traders and $2% bad traders finished."
}

for i in $(seq 0 10 100)
do
    for j in $(seq 0 10 100)
    do
        if [ $(($i+$j)) -gt 70 ]
        then
            continue
        fi
        run $i $j
    done
done

wait

rm -rf output output.zip
mkdir -p output
cp result/*.result output/
zip -r output.zip output
rm -rf output

echo "All runs finished."