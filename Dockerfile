# BUILD
FROM golang:alpine as builder
LABEL maintainer="jamie.kolles@axelspringer.com"

# Install git + SSL ca certificates
RUN apk update && apk add git && apk add ca-certificates

COPY . $GOPATH/src/github.com/TetsuyaXD/evade/
WORKDIR $GOPATH/src/github.com/TetsuyaXD/evade/

RUN echo $GOPATH
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-w -s" -o evade main.go

# RUNTIME
FROM scratch
LABEL maintainer="jamie.kolles@axelspringer.com"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/TetsuyaXD/evade/evade /evade

ENTRYPOINT [ "/evade" ]