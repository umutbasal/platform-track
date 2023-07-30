curl --request PUT \
	--url http://localhost:3000/collection/deployment/backend \
	--header 'Content-Type: application/json' \
	--data '{
	"language": "go",
	"openapi": "backend/api/swagger.json",
	"upstream": ["frontend"],
	"downstream": ["db", "cache"],
	"owner": "jack",
	"status": "ok"
}'

curl --request PUT \
	--url http://localhost:3000/collection/deployment/frontend \
	--header 'Content-Type: application/json' \
	--data '{
	"language": "javascript",
	"downstream": ["backend"],
	"status": "ok"
}'


curl --request PUT \
	--url http://localhost:3000/collection/cluster/cluster1 \
	--header 'Content-Type: application/json' \
	--data '{
	"status": "ok"
}'


curl --request PUT \
	--url http://localhost:3000/collection/db/mongodb \
	--header 'Content-Type: application/json' \
	--data '{
	"status": "ok"
}'

curl --request PUT \
	--url http://localhost:3000/collection/cache/redis \
	--header 'Content-Type: application/json' \
	--data '{
	"upstream": ["backend"],
	"status": "ok",
	"type": "cluster"
}'

curl --request PUT \
	--url http://localhost:3000/collection/storage/s3-my-bucket \
	--header 'Content-Type: application/json' \
	--data '{
	"status": "ok",
	"size": 100,
	"account": "1234567890",
	"bucket": "my-bucket"
}'
