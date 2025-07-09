DB_URI := $(shell jq -r .db.uri $(ENVFILE))

.PHONY: psql
psql:
	@psql '$(DB_URI)'

.PHONY: psql-query
psql-query:
	@echo "Conductivity:"
	@psql '$(DB_URI)' -c 'SELECT * FROM conductivity ORDER BY ts DESC LIMIT 50'
	@echo "Temperature:"
	@psql '$(DB_URI)' -c 'SELECT * FROM  temperature ORDER BY ts DESC LIMIT 50'
	@echo "Humidity:"
	@psql '$(DB_URI)' -c 'SELECT * FROM     humidity ORDER BY ts DESC LIMIT 50'
