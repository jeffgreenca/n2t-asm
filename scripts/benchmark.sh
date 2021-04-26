#!/bin/bash
./scripts/build.sh || exit 1

for i in $(seq 1 100); do
	./scripts/time.sh 2>&1 | grep real | awk '{print $2}' | cut -dm -f 2 | tr -d s
done | jq -s add/length

#| sort | uniq -c
