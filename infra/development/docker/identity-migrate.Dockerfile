FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache curl \
 && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz \
 | tar xz \
 && mv migrate /usr/local/bin/migrate \
 && chmod +x /usr/local/bin/migrate

COPY services/identity-service/internal/migrate/migrations /migrations

ENTRYPOINT ["/usr/local/bin/migrate"]
