# Stage 1: Build
ARG GO_VERSION=1.25.3
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the rest of the code
COPY . .

# Build the binary
# Adjust ./cmd if your main.go is elsewhere
RUN go build -v -o /run-app ./cmd

# Stage 2: Minimal runtime image
FROM debian:bookworm-slim

WORKDIR /app

# Copy binary from builder
COPY --from=builder /run-app /app/run-app

# Expose port
ENV PORT=8080

# Run the binary
CMD ["/app/run-app"]
