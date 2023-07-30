json_data=$(cat logs.jsonl)

#echo "Number of logs: $(echo "$json_data" | jq -s 'length')"

# Extract timestamps from JSON data and find the minimum and maximum values
start_ts=$(echo "$json_data" | jq -s 'map(select(.collection == "deployment" and .identifier == "frontend")) | min_by(.timestamp) | .timestamp')
end_ts=$(echo "$json_data" | jq -s 'map(select(.collection == "deployment" and .identifier == "frontend")) | max_by(.timestamp) | .timestamp')

# Calculate the timestamp range for the 24-hour period
end_ts=$((start_ts + (24 * 60 * 60)))

# Use jq with the determined timestamps to count frontend logs in the 24-hour period
count_frontend_logs=$(echo "$json_data" | jq -s --argjson start_ts $start_ts --argjson end_ts $end_ts 'map(select(.collection == "deployment" and .identifier == "frontend" and (.timestamp >= $start_ts and .timestamp <= $end_ts))) | length')

# Print the result
#echo "Number of frontend deployment in the 24-hour period: $count_frontend_logs"

start_ts=$(echo "$json_data" | jq -s 'map(select(.collection == "deployment" and .identifier == "backend")) | min_by(.timestamp) | .timestamp')
end_ts=$(echo "$json_data" | jq -s 'map(select(.collection == "deployment" and .identifier == "backend")) | max_by(.timestamp) | .timestamp')

# Calculate the timestamp range for the 24-hour period
end_ts=$((start_ts + (24 * 60 * 60)))

# Use jq with the determined timestamps to count frontend logs in the 24-hour period
count_backend_logs=$(echo "$json_data" | jq -s --argjson start_ts $start_ts --argjson end_ts $end_ts 'map(select(.collection == "deployment" and .identifier == "backend" and (.timestamp >= $start_ts and .timestamp <= $end_ts))) | length')



echo "{\"frontend_deployments_day\": $count_frontend_logs, \"backend_deployments_day\": $count_backend_logs, "total_logs": $(echo "$json_data" | jq -s 'length')}"
