# Tahap build
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Setup untuk cross-compilation
ARG TARGETPLATFORM
ARG BUILDPLATFORM
RUN echo "Building on $BUILDPLATFORM for $TARGETPLATFORM"

# Set up proper Go architecture flags khusus untuk amd64 dan arm64
RUN case "$TARGETPLATFORM" in \
    "linux/amd64") \
        echo "Building for AMD64" && \
        GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o main . ;; \
    "linux/arm64") \
        echo "Building for ARM64" && \
        GOOS=linux CGO_ENABLED=0 GOARCH=arm64 go build -o main . ;; \
    *) \
        echo "Unsupported platform: $TARGETPLATFORM" && exit 1 ;; \
    esac

# Tahap runtime (minimal)
FROM --platform=$TARGETPLATFORM alpine:latest

WORKDIR /app

# Salin binary dari tahap builder
COPY --from=builder /app/main .

# Jalankan aplikasi
CMD ["./main"]

