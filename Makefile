.PHONY: go-test nitrogen-test nitrogen modules

all: go-test modules nitrogen nitrogen-test nitrogen-test-vm

go-test:
	go test ./...

nitrogen-test:
	@echo "Run Nitrogen source test suite"
	@for test in tests/*.ni; do \
		./bin/nitrogen -nocompile "$$test"; \
	done

nitrogen-test-vm:
	@echo "Run Nitrogen source test suite using VM"
	@for test in tests/*.ni; do \
		./bin/nitrogen "$$test"; \
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
