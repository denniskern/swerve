# BUILD
FROM golang:alpine as builder
LABEL maintainer="dennis.kern@axelspringer.com"

# Install git + SSL ca certificates
RUN apk update && apk add git && apk add ca-certificates

# RUNTIME
FROM debian:latest
LABEL maintainer="dennis.kern@axelspringer.com"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY orpcer /orpcer

ENTRYPOINT [ "/orpcer" ]
