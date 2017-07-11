.PHONY: go-test nitrogen-test build-exec build

all: build

go-test:
	go test ./...

nitrogen-test:
	@echo "Run Nitrogen source test suite"
	@for test in tests/*.ni; do \
		nitrogen -f "$$test"; \
	done

build-exec:
	go install ./cmd/nitrogen/...

build: go-test build-exec nitrogen-test
