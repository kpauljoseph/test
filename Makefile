.PHONY: test build clean proto run-server run-client lint fmt

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt

# Binary names
SERVER_BINARY=grpc-server
CLIENT_BINARY=grpc-client

# Build binaries
build: build-server build-client

build-server:
	$(GOBUILD) -o bin/$(SERVER_BINARY) cmd/server/main.go

build-client:
	$(GOBUILD) -o bin/$(CLIENT_BINARY) cmd/client/main.go

test:
	$(GOTEST) -v ./...

test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

proto:
	cd proto && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative blog.proto

run-server: build-server
	./bin/$(SERVER_BINARY)

run-client: build-client
	./bin/$(CLIENT_BINARY)

lint:
	golangci-lint run

fmt:
	$(GOFMT) -w .
	$(GOCMD) mod tidy

clean:
	$(GOCLEAN)
	rm -rf bin/
	rm -f coverage.out coverage.html