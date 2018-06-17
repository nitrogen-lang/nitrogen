.PHONY: go-test nitrogen-test nitrogen nitrogen-no-modules modules

all: go-test modules nitrogen nitrogen-test

test: go-test nitrogen-test

go-test:
	go test ./...

nitrogen-test:
	@echo "Run Nitrogen source test suite"
	@p="$$(pwd)"; \
	for test in tests/*.ni; do \
		./bin/nitrogen -M $$p/stdlib -M $$p/built-modules "$$test"; \
	done

nitrogen:
	go build -o bin/nitrogen ./cmd/nitrogen/...

nitrogen-no-modules:
	CGO_ENABLED=0 go build -o bin/nitrogen ./cmd/nitrogen/...

modules:
	rm -f ./built-modules/*
	@p="$$(pwd)"; \
	for m in ./modules/*; do \
		cd "$$m"; \
		echo "Building module $$(basename $$m).so"; \
		go build -buildmode=plugin -o "../../built-modules/$$(basename $$m).so" .; \
		cd "$$p"; \
	done
