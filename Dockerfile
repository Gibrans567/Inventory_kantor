# Tahap build
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

# Tahap runtime dengan nginx
FROM alpine:latest

# Install nginx dan supervisor
RUN apk add --no-cache nginx supervisor

WORKDIR /app

# Salin binary dari tahap builder
COPY --from=builder /app/main .

# Salin seluruh folder Inventaris untuk serve frontend
COPY --from=builder /app/Inventaris ./Inventaris

# Salin file konfigurasi nginx
COPY --from=builder /app/nginx.conf /etc/nginx/nginx.conf

# Buat direktori yang diperlukan nginx
RUN mkdir -p /var/log/nginx /var/lib/nginx/tmp /run/nginx /var/log/supervisor

# Buat file konfigurasi supervisor
RUN echo '[supervisord]' > /etc/supervisord.conf && \
    echo 'nodaemon=true' >> /etc/supervisord.conf && \
    echo 'logfile=/var/log/supervisor/supervisord.log' >> /etc/supervisord.conf && \
    echo 'pidfile=/var/run/supervisord.pid' >> /etc/supervisord.conf && \
    echo '' >> /etc/supervisord.conf && \
    echo '[program:nginx]' >> /etc/supervisord.conf && \
    echo 'command=nginx -g "daemon off;"' >> /etc/supervisord.conf && \
    echo 'autostart=true' >> /etc/supervisord.conf && \
    echo 'autorestart=true' >> /etc/supervisord.conf && \
    echo 'stderr_logfile=/var/log/nginx/error.log' >> /etc/supervisord.conf && \
    echo 'stdout_logfile=/var/log/nginx/access.log' >> /etc/supervisord.conf && \
    echo '' >> /etc/supervisord.conf && \
    echo '[program:goapp]' >> /etc/supervisord.conf && \
    echo 'command=/app/main' >> /etc/supervisord.conf && \
    echo 'directory=/app' >> /etc/supervisord.conf && \
    echo 'autostart=true' >> /etc/supervisord.conf && \
    echo 'autorestart=true' >> /etc/supervisord.conf && \
    echo 'stderr_logfile=/var/log/supervisor/goapp.err.log' >> /etc/supervisord.conf && \
    echo 'stdout_logfile=/var/log/supervisor/goapp.out.log' >> /etc/supervisord.conf

# Expose port
EXPOSE 80

# Jalankan supervisor untuk mengelola nginx dan aplikasi Go
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]