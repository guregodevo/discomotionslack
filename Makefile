BINARY=currency-fetcher

VERSION=1.0.0
GIT_COMMIT=$(shell git rev-parse HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S_UTC')

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME}"

GIT_VERSION=$(shell git describe --tags --always --long --dirty --abbrev=7)
GO_PKG=github.com/dailymotion-leo/discomotionslack

GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
GOVERSION=1.9.0
OUTPUT_DIR=.
OUTPUT_NAME=${BINARY}
DOCKER_IMAGE_VERSION=${GIT_VERSION}
DOCKER_REGISTRY=374211583830.dkr.ecr.us-west-2.amazonaws.com
DOCKER_REPOSITORY=dmx_discomotionslack_api
DOCKER_REPOSITORY_VERSION=latest
UID=$(shell id -u)
GID=$(shell id -g)
DOCKER_IMAGE_NAME=dmx/${BINARY}



build:
	go build ${LDFLAGS} -o ${BINARY}

fmt:
	gofmt -w ./$*

tests:

build-in-docker:
	docker pull golang:${GOVERSION}
	docker run --rm -v ${PWD}:/go/src/${GO_PKG} -w /go/src/${GO_PKG} golang:${GOVERSION} make build GOOS=${GOOS} GOARCH=${GOARCH}

build-docker-image: 
	docker build --tag ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_VERSION} .


run-in-docker: build-docker-image
	docker run ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_VERSION}


install:
	go install ${LDFLAGS}

dist: clean tests
	GOOS=linux go build ${LDFLAGS} -o ${BINARY}

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install
