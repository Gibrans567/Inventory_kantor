# Build stage
FROM golang:1.24.3-alpine3.21 AS builder

# Install build dependencies for CGO (required by GORM)
RUN apk add --no-cache \
    git \
    ca-certificates \
    gcc \
    musl-dev

WORKDIR /app

# Copy go modules files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build for ARM64 with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build \
    -ldflags='-w -s -linkmode external -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main main.go

# Final stage - use alpine for CGO builds
FROM alpine:3.21

# Install ca-certificates and timezone data
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -s /bin/sh appuser

# Copy binary from builder
COPY --from=builder /app/main /app/main

# Change ownership and switch to non-root user
RUN chown appuser:appuser /app/main
USER appuser

# Expose port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/app/main"]
