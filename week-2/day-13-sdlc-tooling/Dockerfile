# Stage 1: builder
# Uses the full Go toolchain to compile the binary. This stage is discarded
# after the build — none of the Go compiler, source code, or build cache
# ends up in the final image. Alpine is used over the default Debian base
# to keep even this intermediate stage lean.
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency files first and download before copying source code.
# Docker caches each layer — if go.mod and go.sum haven't changed, this
# layer is reused on rebuild, skipping the download entirely.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source and compile. The binary is written to bin/server.
COPY . .
RUN go build -o bin/server ./cmd/server

# Stage 2: final image
# Starts from a minimal Alpine base with no Go toolchain. Only the compiled
# binary is copied from the builder stage. Result: a ~37MB image instead of
# the ~600MB+ that a full Go image would produce.
FROM alpine:latest

WORKDIR /app

# Copy only the binary from the builder stage — nothing else.
COPY --from=builder /app/bin/server .

# EXPOSE documents the port the server listens on. It does not publish the port
# — that is done at runtime with docker run -p 8081:8081.
EXPOSE 8081

# CMD is the default command when the container starts.
CMD ["./server"]
