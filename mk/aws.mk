# Endpoint made available after running 'sam local start-api'
LOCAL_API := http://localhost:3000/somon

# Obtained after deploying the lambda
LAMBDA_URL := $(shell jq .aws.lambdaUrl $(ENVFILE))

.PHONY: localPostSimple
localPostSimple:
	curl -X POST --header "Content-Type: application/json" --data '{"msg": "miau"}' $(LOCAL_API)

.PHONY: localPostOnce
localPostOnce:
	curl -X POST --header "Content-Type: application/json" --data @testdata/raw.json $(LOCAL_API)

.PHONY: awsPostOnce
awsPostOnce:
	curl -X POST --header "Content-Type: application/json" --data @testdata/raw.json $(LAMBDA_URL)
