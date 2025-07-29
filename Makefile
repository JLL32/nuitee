include .env

# ===============================================================================
# HELPERS
# ===============================================================================

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

# ===============================================================================
# DEVELOPMENT
# ===============================================================================

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api -db-dsn=${NUITEE_DB_DSN}

run/sync:
	go run ./cmd/sync -db-dsn=${NUITEE_DB_DSN} -api-key=${CUPID_API_KEY} -api-url=${CUPID_API_URL} -input='input.txt'

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${NUITEE_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -ext=.sql -dir=./migrations -seq ${name}

## db/migrations/up: apply all up migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${NUITEE_DB_DSN} up

# ===============================================================================
# QUALITY CONTROL
# ===============================================================================

.PHONY: audit
audit: vendor
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ===============================================================================
# BUILD
# ===============================================================================

.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	GOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api
