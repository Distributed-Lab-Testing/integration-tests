FROM golang:1.21.1-alpine

WORKDIR /github.com/apm-aoc/integration-tests
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go test -v -c -o /usr/local/bin/integration_tests.bin ./tests

FROM alpine:3.9
COPY --from=0 /usr/local/bin/integration_tests.bin /usr/local/bin/integration_tests.bin
RUN apk add --no-cache ca-certificates