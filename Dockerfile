# Gunakan base image golang untuk ARM64
FROM --platform=linux/arm64 golang:1.24-alpine AS builder

RUN apk update && apk add --no-cache mysql-client curl

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o inventory-app main.go

# Final stage (ARM64)
FROM --platform=linux/arm64 alpine:latest

RUN apk add --no-cache mysql-client curl

WORKDIR /app

COPY --from=builder /app/inventory-app .
COPY entrypoint.sh .

RUN chmod +x entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["./inventory-app"]


