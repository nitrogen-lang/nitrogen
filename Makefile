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
	@echo
	@p="$$(pwd)"; \
	for test in tests/**/*.ni; do \
		/bin/echo -n -e "$$test - \e[31m"; \
		./bin/nitrogen -M $$p/nitrogen -M $$p/built-modules "$$test"; \
		if [ $$? -ne 0 ]; then /bin/echo -e "\e[0m"; exit 1; fi; \
		/bin/echo -e "\e[32mpassed\e[0m"; \
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
