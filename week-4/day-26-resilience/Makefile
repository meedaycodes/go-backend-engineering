# .PHONY tells make these are command names, not files. Without this, if a file
# named "test" or "build" exists in the directory, make would do nothing.
.PHONY: build run test lint clean tidy

# build compiles the server binary into bin/server. Using -o keeps the binary
# out of the project root and makes clean trivial.
build:
	go build -o bin/server ./cmd/server

# run starts the server directly without producing a binary. Use for development.
run:
	go run ./cmd/server

# test runs all tests with the race detector enabled. -race catches concurrent
# memory access bugs at runtime — always run tests with it, never without.
test:
	go test -v -race ./...

# lint runs golangci-lint using the rules defined in .golangci.yml.
# Catches unchecked errors, style violations, and logic bugs beyond go vet.
lint:
	golangci-lint run ./...

# clean removes compiled binaries. Run before a fresh build or before committing.
clean:
	rm -rf bin/

# tidy removes unused dependencies and adds any missing ones to go.mod and go.sum.
# Run after adding or removing imports.
tidy:
	go mod tidy
