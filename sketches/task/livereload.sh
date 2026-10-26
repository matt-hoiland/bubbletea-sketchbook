#!/usr/bin/env bash

# watch code changes, trigger re-build, and kill process
while true; do
  go build -o _build/debug && pkill -f '_build/debug'
  inotifywait -e modify $(find . -name '*.go') || exit
done
