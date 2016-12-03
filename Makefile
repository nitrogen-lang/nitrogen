.PHONY: test build

all: build

test:
	go test ./...

build: test
	go install ./cmd/nitrogen/...
