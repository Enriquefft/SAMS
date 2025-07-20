BINARY := app

.PHONY: all build test lint fmt clean vet govulncheck

all: build

build:
	go build -o bin/$(BINARY) ./cmd/app

test:
	go test ./...

lint:
	golangci-lint run

fmt:
	gofmt -w $(shell find . -name '*.go' -not -path './vendor/*')

vet:
	go vet ./...

govulncheck:
	go run golang.org/x/vuln/cmd/govulncheck ./...

clean:
	rm -rf bin
