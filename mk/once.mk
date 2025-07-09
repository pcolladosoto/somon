# Simply escape the nested command substitutions. Bear in mind the (necessary) -n flag is not
# working with echo, that's why we've settled on printf(1)...
BASE64_CREDS = $(shell printf "%s" "$$(jq -r .once.user $(ENVFILE)):$$(jq -r .once.pass $(ENVFILE))" | base64)

# Just pull the integration identifier: it's not supposed to change too much...
INTEGRATION_ID = $(shell jq -r .once.integrationId $(ENVFILE))

# Get a hold of a token based on the regular user and password. These tokens
# usually last for 1 hour.
.PHONY: onceGetToken
onceGetToken:
	@curl --request POST --silent \
		-H 'accept: application/json' \
		-H 'authorization: Basic $(BASE64_CREDS)' \
		-H 'content-type: application/json' \
		--data '{"grant_type": "client_credentials"}' \
	https://api.1nce.com/management-api/oauth/token | jq .access_token

# Simply list the available integrations: bear in mind the TOKEN should've
# been exported based on the output of the onceGetToken target. We need
# to escape it twice to prevent make(1) from trying to substitute $(T)...
.PHONY: onceListIntegrations
onceListIntegrations:
	@curl -X GET -L --silent \
		-H "Authorization: Bearer $$TOKEN" \
	'https://api.1nce.com/management-api/v1/integrate/clouds' | jq .

# Test the integration as specified by its integration ID. This ID can be observed
# in the payload obtained with onceListIntegrations.
.PHONY: onceTestIntegration
onceTestIntegration:
	@curl -X POST -L --silent \
		-H "Authorization: Bearer $$TOKEN" \
		-H "Content-Type: application/json" \
		--data '{"msg": "miau"}' \
	https://api.1nce.com/management-api/v1/integrate/clouds/webhooks/$(INTEGRATION_ID)/test | jq .
