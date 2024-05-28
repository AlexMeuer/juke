# Define binary names and their source directories
BINARIES := juke juke-auth juke-serve
SRC_DIRS := ./cmd/juke ./cmd/juke-auth ./cmd/juke-serve

# Default target: build all binaries
all: $(BINARIES)

# Define how to build each binary
juke:
	go build -o bin/juke ./cmd/juke

juke-auth:
	go build -o bin/juke-auth ./cmd/juke-auth

juke-serve:
	go build -o bin/juke-serve ./cmd/juke-serve

# Test all modules
test:
	go test ./...

# Clean up built binaries
clean:
	rm -rf bin/*

# PHONY targets to avoid conflicts with files named 'all', 'clean', etc.
.PHONY: all clean $(BINARIES)

