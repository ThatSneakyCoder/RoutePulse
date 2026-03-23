FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates curl \
 && curl -L https://github.com/pressly/goose/releases/download/v3.21.1/goose_linux_x86_64 \
    -o /usr/local/bin/goose \
 && chmod +x /usr/local/bin/goose

COPY services/tracking-service/internal/migrate/migrations /migrations

ENTRYPOINT ["/usr/local/bin/goose"]