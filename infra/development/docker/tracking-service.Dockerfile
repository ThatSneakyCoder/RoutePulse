FROM alpine
WORKDIR /app

ADD shared shared
ADD build build

ENTRYPOINT build/tracking-service