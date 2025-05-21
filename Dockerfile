# Tahap build - set platform ke arm64
FROM --platform=linux/arm64 golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Set GOARCH ke arm64 untuk memastikan kompilasi spesifik ARM64
ENV GOOS=linux GOARCH=arm64 CGO_ENABLED=0
RUN go build -o main .

# Tahap runtime (minimal) - juga set platform ke arm64
FROM --platform=linux/arm64 alpine:latest

WORKDIR /app

# Salin binary dari tahap builder
COPY --from=builder /app/main .

# Jalankan aplikasi
CMD ["./main"]