# BUILD
#FROM golang:latest as build
FROM golang:latest

LABEL maintainer="jan.michalowsky@axelspringer.com"

WORKDIR /go/src/github.com/axelspringer/swerve
COPY . .

RUN echo $GOPATH
RUN go get -d -v ./...

#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o swerve -v -ldflags "-extldflags '-static'" -a -installsuffix cgo main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o swerve main.go

# RUNTIME
#FROM scratch

#MAINTAINER Jan Michalowsky <sejamich@googlemail.com>

#COPY --from=build /go/src/github.com/axelspringer/swerve/swerve /swerve
RUN cp /go/src/github.com/axelspringer/swerve/swerve /swerve
CMD [ "/swerve" ]