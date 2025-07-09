# What can we remove?
TRASH := bootstrap .aws-sam *.zip

# We don't need the patch anymore!
BUILD_TAGS := lambda.norpc

# In case we end up moving the environment file around, parametrise it
# from the get-go!
ENVFILE = .env

# AWS stuff
AWS_ROLE_ARN := $(shell jq -r .aws.roleArn $(ENVFILE))

targets:
	@echo "Take a look at the Makefile for now..."

help: targets

# Just print all the variables we're deriving from the environment file
# to make sure things look okay. Bear in mind variables are evaluated
# BEFORE any targets are run, that's why we can define variables in
# their corresponding *.mk files and centralise the 'debugging' here.
.PHONY: debug
debug:
	@echo "base64-encoded credentials: $(BASE64_CREDS)"
	@echo "            integration ID: $(INTEGRATION_ID)"
	@echo "                lambda URL: $(LAMBDA_URL)"
	@echo "              AWS role ARN: $(AWS_ROLE_ARN)"
	@echo "                    DB URI: $(DB_URI)"

# This target is leveraged by 'sam build'
build-local: $(wildcard *.go) go.mod
	GOOS=linux go build -tags $(BUILD_TAGS) -o bootstrap

# Invoking 'sam build' will itself invoke this specific target...
build-XToDexma: build-local
	@cp ./bootstrap $(ARTIFACTS_DIR)/.

# Be sure to check https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html
build-deployment: $(wildcard *.go) go.mod
	GOOS=linux GOARCH=amd64 go build -tags $(BUILD_TAGS) -o bootstrap

bundle: build-deployment
	@zip somon.zip bootstrap

# Note that the IAM Role ARN was acquired from the AWS Console
publish: bundle
	@aws lambda create-function --function-name somon      \
		--runtime provided.al2023 --handler bootstrap      \
		--architectures x86_64                             \
		--role $(AWS_ARN)  \
		--zip-file fileb://somon.zip

update: bundle
	@aws lambda update-function-code --function-name somon \
		--zip-file fileb://somon.zip

sam-build:
	@sam build

.PHONY: sam-local-api
sam-local-api: sam-build
	@sam local start-api

.PHONY: clean
clean:
	@rm -rf $(TRASH)

# Include the other targets too!
include mk/*.mk
