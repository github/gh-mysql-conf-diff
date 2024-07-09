# Define variables
BINARY_NAME=gh-mysql-conf-diff
BINARY_PATH=bin/$(BINARY_NAME)
CMD_PATH=./cmd/$(BINARY_NAME)

# Default target
all: build test

# Build the binary
build:
	go build -o $(BINARY_PATH)/$(BINARY_NAME) $(CMD_PATH)

# Run tests
test:
	go test -v ./...

# Clean up binary
clean:
	rm -rf $(BINARY_PATH)

.PHONY: all build test clean