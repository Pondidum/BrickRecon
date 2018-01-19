aws dynamodb create-table \
    --table-name boidlookup_test \
    --attribute-definitions AttributeName=boid,AttributeType=S \
    --key-schema AttributeName=boid,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --endpoint-url http://localhost:8000