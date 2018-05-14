.DEFAULT_GOAL := help
.PHONY: run run-help test test-core test-libc test-lint build-libc check cover
.PHONY: integration-test-stable integration-test-stable-disable-csrf
.PHONY: integration-test-live integration-test-live-wallet
.PHONY: integration-test-disable-wallet-api integration-test-disable-seed-api
.PHONY: install-linters format release clean-release install-deps-ui build-ui help

# Static files directory
GUI_STATIC_DIR = src/gui/static

# Electron files directory
ELECTRON_DIR = electron

# Compilation output
BUILD_DIR = build
BUILDLIB_DIR = $(BUILD_DIR)/libskycoin
LIB_DIR = lib
LIB_FILES = $(shell find ./lib/cgo -type f -name "*.go")
SRC_FILES = $(shell find ./src -type f -name "*.go")
BIN_DIR = bin
INCLUDE_DIR = include

# Compilation flags
CC = gcc
LIBC_LIBS = -lcriterion
LDFLAGS = -I$(INCLUDE_DIR) -I$(BUILD_DIR)/usr/include -L $(BUILDLIB_DIR) -L$(BUILD_DIR)/usr/lib

# Platform specific checks
OSNAME = $(TRAVIS_OS_NAME)

ifeq ($(shell uname -s),Linux)
  LDLIBS=$(LIBC_LIBS) -lpthread
	LDPATH=$(shell printenv LD_LIBRARY_PATH)
	LDPATHVAR=LD_LIBRARY_PATH
ifndef OSNAME
  OSNAME = linux
endif
else ifeq ($(shell uname -s),Darwin)
ifndef OSNAME
  OSNAME = osx
endif
	LDLIBS = $(LIBC_LIBS)
	LDPATH=$(shell printenv DYLD_LIBRARY_PATH)
	LDPATHVAR=DYLD_LIBRARY_PATH
else
	LDLIBS = $(LIBC_LIBS)
	LDPATH=$(shell printenv LD_LIBRARY_PATH)
	LDPATHVAR=LD_LIBRARY_PATH
endif

run:  ## Run the skycoin node. To add arguments, do 'make ARGS="--foo" run'.
	./run.sh ${ARGS}

run-testnet:  ## Run the skycoin testnet node. To add arguments, do 'make ARGS="--foo" run'.
	go run cmd/skycoin/skycoin.go --testchain --gui-dir="./${STATIC_DIR}" ${ARGS}

run-help: ## Show skycoin node help
	@go run cmd/skycoin/skycoin.go --help

test: ## Run tests for Skycoin
	go test ./cmd/... -timeout=5m
	go test ./src/... -timeout=5m

test-386: ## Run tests for Skycoin with GOARCH=386
	GOARCH=386 go test ./cmd/... -timeout=5m
	GOARCH=386 go test ./src/... -timeout=5m

test-amd64: ## Run tests for Skycoin with GOARCH=amd64
	GOARCH=amd64 go test ./cmd/... -timeout=5m
	GOARCH=amd64 go test ./src/... -timeout=5m

configure-build:
	mkdir -p $(BUILD_DIR)/usr/tmp $(BUILD_DIR)/usr/lib $(BUILD_DIR)/usr/include
	mkdir -p $(BUILDLIB_DIR) $(BIN_DIR) $(INCLUDE_DIR)

build-libc: configure-build $(BUILDLIB_DIR)/libskycoin.so $(BUILDLIB_DIR)/libskycoin.a ## Build libskycoin C client library

$(BUILDLIB_DIR)/libskycoin.so $(BUILDLIB_DIR)/libskycoin.a: $(LIB_FILES) $(SRC_FILES)
	rm -Rf $(BUILDLIB_DIR)/*
	go build -buildmode=c-shared  -o $(BUILDLIB_DIR)/libskycoin.so $(LIB_FILES)
	go build -buildmode=c-archive -o $(BUILDLIB_DIR)/libskycoin.a  $(LIB_FILES)
	mv $(BUILDLIB_DIR)/libskycoin.h $(INCLUDE_DIR)/

## Build libskycoin C client library and executable C test suites
## with debug symbols. Use this target to debug the source code
## with the help of an IDE
build-libc-dbg: configure-build $(BUILDLIB_DIR)/libskycoin.so $(BUILDLIB_DIR)/libskycoin.a
	$(CC) -g -o $(BIN_DIR)/test_libskycoin_shared $(LIB_DIR)/cgo/tests/*.c -lskycoin                    $(LDLIBS) $(LDFLAGS)
	$(CC) -g -o $(BIN_DIR)/test_libskycoin_static $(LIB_DIR)/cgo/tests/*.c $(BUILDLIB_DIR)/libskycoin.a $(LDLIBS) $(LDFLAGS)

test-libc: build-libc ## Run tests for libskycoin C client library
	cp $(LIB_DIR)/cgo/tests/*.c $(BUILDLIB_DIR)/
	$(CC) -o $(BIN_DIR)/test_libskycoin_shared $(BUILDLIB_DIR)/*.c -lskycoin                    $(LDLIBS) $(LDFLAGS)
	$(CC) -o $(BIN_DIR)/test_libskycoin_static $(BUILDLIB_DIR)/*.c $(BUILDLIB_DIR)/libskycoin.a $(LDLIBS) $(LDFLAGS)
	$(LDPATHVAR)="$(LDPATH):$(BUILD_DIR)/usr/lib"                 $(BIN_DIR)/test_libskycoin_static
	$(LDPATHVAR)="$(LDPATH):$(BUILD_DIR)/usr/lib:$(BUILDLIB_DIR)" $(BIN_DIR)/test_libskycoin_shared

lint: ## Run linters. Use make install-linters first.
	vendorcheck ./...
	gometalinter --deadline=3m --concurrency=2 --disable-all --tests --vendor --skip=lib/cgo --warn-unmatched-nolint \
		-E goimports \
		-E golint \
		-E varcheck \
		-E unparam \
		./...
	# lib cgo can't use golint because it needs export directives in function docstrings that do not obey golint rules
	gometalinter --deadline=3m --concurrency=2 --disable-all --tests --vendor --warn-unmatched-nolint \
		-E goimports \
		-E varcheck \
		-E unparam \
		./lib/cgo/...

check: lint test integration-test-stable integration-test-stable-disable-csrf integration-test-disable-wallet-api integration-test-disable-seed-api ## Run tests and linters

integration-test-stable: ## Run stable integration tests
	./ci-scripts/integration-test-stable.sh -c

integration-test-stable-disable-csrf: ## Run stable integration tests with CSRF disabled
	./ci-scripts/integration-test-stable.sh

integration-test-live: ## Run live integration tests
	./ci-scripts/integration-test-live.sh -c

integration-test-live-wallet: ## Run live integration tests with wallet
	./ci-scripts/integration-test-live.sh -w

integration-test-live-disable-csrf: ## Run live integration tests against a node with CSRF disabled
	./ci-scripts/integration-test-live.sh

integration-test-disable-wallet-api: ## Run disable wallet api integration tests
	./ci-scripts/integration-test-disable-wallet-api.sh

integration-test-disable-seed-api: ## Run enable seed api integration test
	./ci-scripts/integration-test-disable-seed-api.sh

cover: ## Runs tests on ./src/ with HTML code coverage
	go test -cover -coverprofile=cover.out -coverpkg=github.com/skycoin/skycoin/... ./src/...
	go tool cover -html=cover.out

install-linters: ## Install linters
	go get -u github.com/FiloSottile/vendorcheck
	go get -u github.com/alecthomas/gometalinter
	gometalinter --vendored-linters --install

install-deps-libc: configure-build ## Install locally dependencies for testing libskycoin
	wget -O $(BUILD_DIR)/usr/tmp/criterion-v2.3.2-$(OSNAME)-x86_64.tar.bz2 https://github.com/Snaipe/Criterion/releases/download/v2.3.2/criterion-v2.3.2-$(OSNAME)-x86_64.tar.bz2
	tar -x -C $(BUILD_DIR)/usr/tmp/ -j -f $(BUILD_DIR)/usr/tmp/criterion-v2.3.2-$(OSNAME)-x86_64.tar.bz2
	ls $(BUILD_DIR)/usr/tmp/criterion-v2.3.2/include
	ls -1 $(BUILD_DIR)/usr/tmp/criterion-v2.3.2/lib     | xargs -I NAME mv $(BUILD_DIR)/usr/tmp/criterion-v2.3.2/lib/NAME     $(BUILD_DIR)/usr/lib/NAME
	ls -1 $(BUILD_DIR)/usr/tmp/criterion-v2.3.2/include | xargs -I NAME mv $(BUILD_DIR)/usr/tmp/criterion-v2.3.2/include/NAME $(BUILD_DIR)/usr/include/NAME

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w -local github.com/skycoin/skycoin ./cmd
	goimports -w -local github.com/skycoin/skycoin ./src
	goimports -w -local github.com/skycoin/skycoin ./docker
	goimports -w -local github.com/skycoin/skycoin ./lib

install-deps-ui:  ## Install the UI dependencies
	cd $(GUI_STATIC_DIR) && npm install

lint-ui:  ## Lint the UI code
	cd $(GUI_STATIC_DIR) && npm run lint

test-ui:  ## Run UI tests
	cd $(GUI_STATIC_DIR) && npm run test

test-ui-e2e:  ## Run UI e2e tests
	./ci-scripts/ui-e2e.sh

build-ui:  ## Builds the UI
	cd $(GUI_STATIC_DIR) && npm run build

build-ui-travis:  ## Builds the UI for travis
	cd $(GUI_STATIC_DIR) && npm run build-travis

release: ## Build electron apps, the builds are located in electron/release folder.
	cd $(ELECTRON_DIR) && ./build.sh
	@echo release files are in the folder of electron/release

clean-release: ## Clean dist files and delete all builds in electron/release
	rm $(ELECTRON_DIR)/release/*

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
