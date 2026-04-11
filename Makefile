.PHONY: help deps fmt lint check test-unit test-all test-integration

GO ?= go
MODULE := github.com/pattonkan/sui-go

# Packages with tests that require local Sui/faucet or public network.
INTEGRATION_PKGS := \
	$(MODULE)/sui/movebcs \
	$(MODULE)/sui/suiptb \
	$(MODULE)/suiclient \
	$(MODULE)/suiclient/conn \
	$(MODULE)/suisigner

ALL_PKGS := $(shell $(GO) list ./...)
UNIT_PKGS := $(filter-out $(INTEGRATION_PKGS),$(ALL_PKGS))

help:
	@echo "Targets:"
	@echo "  make deps             - download Go dependencies"
	@echo "  make fmt              - go fmt all packages"
	@echo "  make lint             - run go vet"
	@echo "  make test-unit        - run offline/unit tests only"
	@echo "  make test-all         - run full test suite"
	@echo "  make test-integration - run integration/network tests"
	@echo "  make check            - fmt + lint + test-unit"

deps:
	$(GO) mod download

fmt:
	$(GO) fmt ./...

lint:
	$(GO) vet ./...

test-unit:
	$(GO) test -v $(UNIT_PKGS)

test-all:
	$(GO) test -v ./...

test-integration:
	$(GO) test -v $(INTEGRATION_PKGS)

check: fmt lint test-unit
