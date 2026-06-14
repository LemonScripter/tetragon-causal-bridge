# DCC Tetragon Bridge Makefile

GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=tetragon-dcc-bridge

all: test build

build:
	@echo "Building Tetragon DCC Bridge..."
	cd src && $(GOBUILD) -v -o ../bin/$(BINARY_NAME) .

test:
	@echo "Running Go Unit Tests..."
	$(GOTEST) -v ./src/...

test-integration:
	@echo "Running Logic Verification (Python)..."
	python3 tests/verify_bridge.py

clean:
	rm -f bin/$(BINARY_NAME)

.PHONY: all build test test-integration clean
