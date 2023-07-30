json_data=$(cat logs.jsonl)

#count_backend_logs=$(echo "$json_data" | jq -s --argjson start_ts $start_ts --argjson end_ts $end_ts 'map(select(.collection == "deployment" and .identifier == "backend" and (.timestamp >= $start_ts and .timestamp <= $end_ts))) | length')
# identifiers
identifiers=$(echo "$json_data" | jq -s 'map(.identifier) | unique')


echo "{\"identifiers\": $identifiers}"
