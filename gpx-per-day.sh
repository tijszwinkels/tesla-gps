#!/bin/bash
executable_path="./tesla-gps"
params="--token ./t.token --verbose"

while true; do
  today=$(date +%Y-%m-%d)
  output_file="${today}.gpx"

  # Calculate seconds until midnight
  current_time=$(date +%s)
  midnight=$(date -d "$(date -d tomorrow +%Y-%m-%d) 00:00:00" +%s)
  sleep_duration=$((midnight - current_time))

  trap 'kill ${!}; break' SIGINT SIGTERM

  "${executable_path}" ${params} > "${output_file}" &
  pid=$!

  # Sleep until midnight, then kill the executable
  sleep $sleep_duration
  kill $pid
  wait $pid
done