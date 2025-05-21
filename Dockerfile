<<<<<<< HEAD
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
=======
FROM --platform=linux/arm64 golang:1.24.3-alpine3.21 AS builder-arm64

RUN apk add --no-cache git ca-certificates gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -ldflags='-w -s -linkmode external -extldflags "-static"' -a -installsuffix cgo -o main main.go
>>>>>>> cf0077d (dockerfile)

# Build for ARM64 with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build \
    -ldflags='-w -s -linkmode external -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main main.go

<<<<<<< HEAD
# Final stage - use alpine for CGO builds
FROM alpine:3.21
=======
FROM --platform=linux/arm64 alpine:3.21
>>>>>>> cf0077d (dockerfile)

# Install ca-certificates and timezone data
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -s /bin/sh appuser

<<<<<<< HEAD
# Copy binary from builder
COPY --from=builder /app/main /app/main
=======
# Untuk tahap akhir, copy secara manual sesuai platform buildx akan pilih
COPY --from=builder-arm64 /app/main /app/main
>>>>>>> cf0077d (dockerfile)

# Change ownership and switch to non-root user
RUN chown appuser:appuser /app/main
USER appuser

# Expose port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/app/main"]
