FROM golang:1.24.3-alpine3.21 AS builder-amd64

RUN apk add --no-cache git ca-certificates gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s -linkmode external -extldflags "-static"' -a -installsuffix cgo -o main main.go


FROM golang:1.24.3-alpine3.21 AS builder-arm64

RUN apk add --no-cache git ca-certificates gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -ldflags='-w -s -linkmode external -extldflags "-static"' -a -installsuffix cgo -o main main.go


FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata
RUN adduser -D -s /bin/sh appuser
WORKDIR /app

# Untuk tahap akhir, copy secara manual sesuai platform buildx akan pilih
COPY --from=builder-amd64 /app/main /app/main

RUN chown appuser:appuser /app/main
USER appuser

EXPOSE 8080
ENTRYPOINT ["/app/main"]
