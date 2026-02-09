include .envrc

# ================================================================================================== #
# DEVELOPMENT
# ================================================================================================== #

.PHONY: run/api
run/api:
	@. ./.envrc && go run ./cmd/api -db-dsn=$${APP_DROP_DSN}

.PHONY: db/mig/new
db/mig/new:
	@echo 'creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

.PHONY: db/mig/up
db/mig/up:
	@echo 'running up migrations'
	migrate -path ./migrations -database ${APP_DROP_DSN} up

.PHONY: db/mig/down
db/mig/down:
	@echo 'running down migrations'
	migrate -path ./migrations -database ${APP_DROP_DSN} down
