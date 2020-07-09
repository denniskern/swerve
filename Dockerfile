# BUILD
FROM golang:alpine as builder
LABEL maintainer="jamie.kolles@axelspringer.com"

# Install git + SSL ca certificates
RUN apk update && apk add git && apk add ca-certificates

COPY . /swerve
WORKDIR /swerve

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-w -s" -o swerve main.go

# RUNTIME
FROM scratch
LABEL maintainer="jamie.kolles@axelspringer.com"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /swerve/swerve /swerve

ENTRYPOINT [ "/swerve" ]