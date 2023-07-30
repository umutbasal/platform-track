clusterName=$(kubectl config current-context)
nodeNames=$(kubectl get nodes -o jsonpath='{.items[*].metadata.name}')

JSON_STRING=$(jq -n \
	--arg clusterName "$clusterName" \
	'{nodes: [], cluster: $clusterName}')

for nodeName in $nodeNames; do
	JSON_STRING=$(jq \
		--arg nodeName "$nodeName" \
		'.nodes += [$nodeName]' \
		<<<"$JSON_STRING")
done

curl --request PUT \
	--url http://localhost:3000/collection/cluster/$clusterName --header 'Content-Type: application/json' \
	--data "$JSON_STRING"
