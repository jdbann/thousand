BINARY_NAME=tmp/build/thousand

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | sort -d | column -t -s ':' |  sed -e 's/^/ /'

## build: build the application
.PHONY: build
build:
	go build -o ${BINARY_NAME} -v .

## build/clean: remove generated files
.PHONY: build/clean
build/clean:
	go mod tidy
	go clean
	-rm ${BINARY_NAME}
	-rm -rf tmp/

## ci: setup and run the CI process
.PHONY: ci
ci: lint test/db/migrate test/run

## dev: migrate and run the app (default dev task)
.PHONY: dev
dev: dev/db/migrate dev/run

## dev/db/migrate: migrate the development database
.PHONY: dev/db/migrate
dev/db/migrate: build
	${BINARY_NAME} migrate run

## dev/db/setup: setup the development database
.PHONY: dev/db/setup
dev/db/setup: build
	-${BINARY_NAME} db drop
	${BINARY_NAME} db create

## dev/run: run the app
.PHONY: dev/run
dev/run: build
	${BINARY_NAME} run

## dev/setup: setup the development environment
.PHONY: dev/setup
dev/setup: dev/db/setup dev/db/migrate

## generate: run all code generation steps
.PHONY: generate
generate:
	sqlc generate
	goimports -w .

## lint: run additional linting steps
.PHONY: lint
lint:
	bin/verifymakefile

## local: setup a local developer environment (both dev and test)
.PHONY: local
local: build dev/setup test/setup

## test: setup the test environment and run all tests (default test task)
.PHONY: test
test: test/setup test/run

## test/db/migrate: migrate the test database
.PHONY: test/db/migrate
test/db/migrate: build
	${BINARY_NAME} --environment=test migrate run

## test/db/setup: setup the test database
.PHONY: test/db/setup
test/db/setup: build
	-${BINARY_NAME} --environment=test db drop
	${BINARY_NAME} --environment=test db create

## test/run: run all tests
.PHONY: test/run
test/run:
	go test ./...

## test/setup: setup the test environment
.PHONY: test/setup
test/setup: test/db/setup test/db/migrate
