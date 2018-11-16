GO=go
RACE := $(shell test $$(go env GOARCH) != "amd64" || (echo "-race"))
GOFLAGS= 
BIN=bin/swerve
VERSION := $(shell git rev-parse HEAD)
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
IMAGE="axelspringer/swerve"

all: build/local

help:
	@echo 'Available commands:'
	@echo
	@echo 'Usage:'
	@echo '    make restore     	Install go deps.'
	@echo '    make build/local    	Compile the project.'
	@echo '    make test/local    	Run ginkgo test suites.'
	@echo '    make test/local    	Run ginkgo test suites.'
	@echo '    make build/docker    Create docker container'
	@echo '    make clean    		Clean the directory tree.'
	@echo

test/local:
	ginkgo --race --cover --coverprofile "$(ROOT_DIR)/swerve.coverprofile" ./...
	go tool cover -html=swerve.coverprofile -o swerve_test_coverage.html

build/local:
	$(GO) get ./...
	$(GO) build -ldflags "-X main.Version=$(VERSION)" -o $(BIN) $(GOFLAGS) $(RACE) main.go

build/docker:
	docker build -t ${IMAGE}.

push/docker:
	docker login -u="${DOCKER_USERNAME}" -p="${DOCKER_PASSWORD}"
	docker push ${IMAGE}

run/dynamo:
	docker-compose -f example/stack/stack.yml up

restore:
	dep ensure
