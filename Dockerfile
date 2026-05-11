# Build stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /ddl_guard ./cmd/ddl_guard/

# Final stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata netcat-openbsd

WORKDIR /app

COPY --from=builder /ddl_guard /app/ddl_guard
COPY configs/config.yaml /app/configs/config.yaml

RUN mkdir -p /app/data

EXPOSE 8080

ENV TZ=Asia/Shanghai

RUN echo '#!/bin/sh' > /app/entrypoint.sh && \
    echo 'set -e' >> /app/entrypoint.sh && \
    echo 'echo "Waiting for PostgreSQL..."' >> /app/entrypoint.sh && \
    echo 'while ! nc -z postgres 5432; do sleep 1; done' >> /app/entrypoint.sh && \
    echo 'echo "PostgreSQL is ready"' >> /app/entrypoint.sh && \
    echo 'echo "Waiting for Redis..."' >> /app/entrypoint.sh && \
    echo 'while ! nc -z redis 6379; do sleep 1; done' >> /app/entrypoint.sh && \
    echo 'echo "Redis is ready"' >> /app/entrypoint.sh && \
    echo 'if [ ! -f /app/data/conf/config.yaml ]; then' >> /app/entrypoint.sh && \
    echo '    mkdir -p /app/data/conf' >> /app/entrypoint.sh && \
    echo '    cp /app/configs/config.yaml /app/data/conf/config.yaml' >> /app/entrypoint.sh && \
    echo 'fi' >> /app/entrypoint.sh && \
    echo '/app/ddl_guard init -d /app/data' >> /app/entrypoint.sh && \
    echo 'exec /app/ddl_guard run -d /app/data' >> /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
