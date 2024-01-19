VERSION ?= $(shell git describe --tags --always --dirty)
BUILDTIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILDER := $(shell echo "`git config user.name` <`git config user.email`>")
CGO_ENABLED ?= 1
MODULE_PATHS ?= "nitrogen"

LDFLAGS := -X 'main.version=$(VERSION)' \
			-X 'main.buildTime=$(BUILDTIME)' \
			-X 'main.builder=$(BUILDER)' \
			-X 'main.builtinModPaths=$(MODULE_PATHS)' \
			-s -w

.PHONY: go-test nitrogen-test build modules build-tools buildc

all: build-tools

build-tools: buildc build

buildc:
	go build -o bin/nitrogenc -ldflags="$(LDFLAGS)" ./cmd/nitrogenc/...

build:
	go build -o bin/nitrogen -ldflags="$(LDFLAGS)" ./cmd/nitrogen/...

test: go-test nitrogen-test

go-test:
	go test ./...

nitrogen-test:
	./tests/run_tests.sh

modules:
ifeq ($(CGO_ENABLED),1)
	rm -f ./built-modules/*
	@p="$$(pwd)"; \
	for m in ./modules/*; do \
		cd "$$m"; \
		echo "Building module $$(basename $$m).so"; \
		go build -buildmode=plugin -o "../../built-modules/$$(basename $$m).so" .; \
		cd "$$p"; \
	done
else
	@echo "CGO disabled, not building modules"
endif
