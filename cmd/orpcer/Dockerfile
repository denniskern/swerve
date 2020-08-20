# BUILD
FROM golang:alpine as builder
LABEL maintainer="dennis.kern@axelspringer.com"

# Install git + SSL ca certificates
RUN apk update && apk add git && apk add ca-certificates

COPY . /orpcer
WORKDIR /orpcer

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-w -s -X github.com/axelspringer/swerve/api.githubHash=`git rev-parse --short HEAD`" -o orpcer main.go

# RUNTIME
FROM scratch
LABEL maintainer="dennis.kern@axelspringer.com"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /orpcer/orpcer /orpcer

ENTRYPOINT [ "/orpcer" ]
