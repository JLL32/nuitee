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
	go run ./cmd/api -db-dsn=${NUITEE_DB_DSN} -openai-key=${OPEN_AI_KEY}

## run/sync: run the sync command
.PHONY: run/sync
run/sync:
	go run ./cmd/sync -db-dsn=${NUITEE_DB_DSN} -api-key=${CUPID_API_KEY} -api-url=${CUPID_API_URL} -input='input.txt'

# ===============================================================================
# DATABASE
# ===============================================================================

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

## db/migrations/down: rollback last migration
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Rolling back last migration...'
	migrate -path ./migrations -database ${NUITEE_DB_DSN} down 1

## db/migrations/version: show current migration version
.PHONY: db/migrations/version
db/migrations/version:
	@echo 'Current migration version:'
	migrate -path ./migrations -database ${NUITEE_DB_DSN} version

## db/reset: reset database (drop and recreate)
.PHONY: db/reset
db/reset: confirm
	@echo 'Resetting database...'
	migrate -path ./migrations -database ${NUITEE_DB_DSN} drop
	migrate -path ./migrations -database ${NUITEE_DB_DSN} up

# ===============================================================================
# TESTING
# ===============================================================================

## test: run all unit tests
.PHONY: test
test:
	@echo 'Running unit tests...'
	go test ./internal/data ./cmd/api

## test/verbose: run tests with verbose output
.PHONY: test/verbose
test/verbose:
	@echo 'Running tests with verbose output...'
	go test -v ./internal/data ./cmd/api

## test/coverage: run tests with coverage report
.PHONY: test/coverage
test/coverage:
	@echo 'Running tests with coverage report...'
	go test -cover ./internal/data ./cmd/api

## test/coverage/html: run tests with HTML coverage report
.PHONY: test/coverage/html
test/coverage/html:
	@echo 'Running tests with HTML coverage report...'
	go test -coverprofile=coverage.out ./internal/data ./cmd/api
	go tool cover -html=coverage.out -o coverage.html
	@echo 'HTML coverage report generated: coverage.html'

## test/race: run tests with race condition detection
.PHONY: test/race
test/race:
	@echo 'Running tests with race condition detection...'
	go test -race ./internal/data ./cmd/api

## test/bench: run benchmark tests
.PHONY: test/bench
test/bench:
	@echo 'Running benchmark tests...'
	go test -bench=. -benchmem ./internal/data ./cmd/api


# ===============================================================================
# QUALITY CONTROL
# ===============================================================================

## audit: run comprehensive quality checks
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
