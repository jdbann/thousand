## help: print this help message
# first target so `make` without args prints help
#
# help docs are generated from comments:
#   single '#' is not output in help message
#   double '#' is used as output in help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | sort -d | column -t -s ':' |  sed -e 's/^/ /'

BUILD_DIRECTORY = tmp/build
BINARY_NAME = thousand
BINARY_PATH = $(BUILD_DIRECTORY)/$(BINARY_NAME)
TEST_DB = $(DATABASE_URL)
ifeq ($(strip $(TEST_DB)),)
TEST_DB = "postgres://localhost:5432/thousand_test?sslmode=disable"
endif

## build: build the application
.PHONY: build
build: build/assets
	go build -o ${BINARY_PATH} -v .

## build/assets: build the css and js assets
.PHONY: build/assets
build/assets:
	./bin/build-assets

## build/clean: remove generated files
.PHONY: build/clean
build/clean:
	go mod tidy
	go clean
	-rm -rf static/*/**/*
	-rm -rf $(BUILD_DIRECTORY)

## build/docker: build the docker image for the app
.PHONY: build/docker
build/docker:
	docker build -t thousand .

# build/path: outputs the configured build path for air.toml
.PHONY: build/path
build/path:
	@echo $(BINARY_PATH)

## check: make sure project is in a tidy state for committing
.PHONY: check
check: generate lint build build/docker local test routes build/clean

## ci: setup and run the CI process
.PHONY: ci
ci: ci/setup lint test/db/migrate test/run

# ci/setup: install dependencies
.PHONY: ci/setup
ci/setup:
	yarn install

## dev: migrate and run the app (default dev task)
.PHONY: dev
dev: dev/db/migrate dev/run

## dev/db/migrate: migrate the development database
.PHONY: dev/db/migrate
dev/db/migrate: build
	${BINARY_PATH} migrate run

## dev/db/setup: setup the development database
.PHONY: dev/db/setup
dev/db/setup: build
	-${BINARY_PATH} db drop
	${BINARY_PATH} db create

## dev/run: run the app
.PHONY: dev/run
dev/run: build
	# ./tmp/bin/air # air currently broken on macOS
	${BINARY_PATH}

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
	bin/verify-makefile
	yarn lint:css
	yarn lint:js
	yarn lint:tmpl

## local: setup a local developer environment (both dev and test)
.PHONY: local
local: local/deps build dev/setup test/setup

## local/deps: install tools required to work on the project
.PHONY: local/deps
local/deps: local/deps/sqlc tmp/bin/air
	yarn install

# local/deps/sqlc: installs sqlc for db querying code generation
.PHONY: local/deps/sqlc
local/deps/sqlc:
	sqlc version || brew install sqlc

# tmp/bin/air: installs air for automatic reloading of application
tmp/bin/air:
	@echo Installing air downloads and executes a script from GitHub:
	@echo   https://raw.githubusercontent.com/cosmtrek/air/master/install.sh
	@echo "Do you want to continue? [y/N] " && read ans && [ $${ans:-N} = y ]
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b ./tmp/bin

## routes: print routes the application serves
.PHONY: routes
routes: build
	${BINARY_PATH} routes

## test: setup the test environment and run all tests (default test task)
.PHONY: test
test: test/setup test/run

## test/db/migrate: migrate the test database
.PHONY: test/db/migrate
test/db/migrate: build
	DATABASE_URL=${TEST_DB} ${BINARY_PATH} migrate run

## test/db/setup: setup the test database
.PHONY: test/db/setup
test/db/setup: build
	-DATABASE_URL=${TEST_DB} ${BINARY_PATH} db drop
	DATABASE_URL=${TEST_DB} ${BINARY_PATH} db create

## test/run: run all tests
.PHONY: test/run
test/run:
	go test ./...

## test/setup: setup the test environment
.PHONY: test/setup
test/setup: test/db/setup test/db/migrate
