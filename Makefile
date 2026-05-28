.PHONY: build test test-verbose test-e2e test-e2e-verbose clean install

BINARY := liepin-cli
BUILD_DIR := bin

build:
	go build -o $(BUILD_DIR)/$(BINARY) .

test:
	go test ./... -count=1

test-verbose:
	go test ./... -v -count=1

test-e2e:
	go test -tags=e2e ./test/e2e/ -count=1

test-e2e-verbose:
	go test -tags=e2e ./test/e2e/ -v -count=1

lint:
	golangci-lint run ./...

clean:
	rm -rf $(BUILD_DIR)

PREFIX ?= $(HOME)/.local
install: build
	install -d $(PREFIX)/bin
	install $(BUILD_DIR)/$(BINARY) $(PREFIX)/bin/$(BINARY)
