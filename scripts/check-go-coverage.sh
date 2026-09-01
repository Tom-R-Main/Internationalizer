#!/bin/sh
set -eu

profile=${1:-coverage.out}
minimum=${2:-60.0}

if [ ! -f "$profile" ]; then
  echo "coverage profile not found: $profile" >&2
  exit 2
fi

total=$(go tool cover -func="$profile" | awk '$1 == "total:" { gsub(/%/, "", $3); print $3 }')
if [ -z "$total" ]; then
  echo "could not read total coverage from $profile" >&2
  exit 2
fi

awk -v total="$total" -v minimum="$minimum" 'BEGIN {
  if ((total + 0) < (minimum + 0)) {
    printf "coverage %.1f%% is below required %.1f%%\n", total, minimum > "/dev/stderr"
    exit 1
  }
  printf "coverage %.1f%% meets required %.1f%%\n", total, minimum
}'
