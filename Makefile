# Copyright 2017 Openprovider Authors. All rights reserved.
# Use of this source code is governed by a MIT-style
# license that can be found in the LICENSE file.

GO_BIN ?= go
ENV_BIN ?= env

export PATH := $(PATH):/usr/local/go/bin

APP=whoisd
PROJECT=github.com/openprovider/whoisd

# Use the 0.0.0 tag for testing, it shouldn't clobber any release builds
RELEASE?=0.5.0
GOOS?=linux
GOARCH?=amd64

REPO_INFO=$(shell git config --get remote.origin.url)
RELEASE_DATE=$(shell date +%FT%T%Z:00)

ifndef COMMIT
	COMMIT := git-$(shell git rev-parse --short HEAD)
endif

build:
	$(GO_BIN) mod tidy
	$(ENV_BIN) CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -a -installsuffix cgo \
		-ldflags "-s -w -X ${PROJECT}/pkg/version.RELEASE=${RELEASE} -X ${PROJECT}/pkg/version.DATE=${RELEASE_DATE} -X ${PROJECT}/pkg/version.COMMIT=${COMMIT} -X ${PROJECT}/pkg/version.REPO=${REPO_INFO}" \
		-o bin/${GOOS}-${GOARCH}/${APP} ${PROJECT}/cmd