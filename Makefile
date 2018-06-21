VERSION ?= $(shell git describe --tags --always --dirty)
BUILDTIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILDER := $(shell echo "`git config user.name` <`git config user.email`>")
CGO_ENABLED ?= 1

LDFLAGS := -X 'main.version=$(VERSION)' \
			-X 'main.buildTime=$(BUILDTIME)' \
			-X 'main.builder=$(BUILDER)' \
			-s -w

.PHONY: go-test nitrogen-test build modules

all: build modules

build:
	go build -o bin/nitrogen -ldflags="$(LDFLAGS)" ./cmd/nitrogen/...

test: go-test nitrogen-test

go-test:
	go test ./...

nitrogen-test:
	@echo "Run Nitrogen source test suite"
	@p="$$(pwd)"; \
	for test in tests/*.ni; do \
		./bin/nitrogen -M $$p/stdlib -M $$p/built-modules "$$test"; \
	done

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
