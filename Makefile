# Makefile for baron.
#
SHA := $(shell git rev-parse --short=10 HEAD)

MAKEFILE_PATH := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
VERSION_DATE := $(shell $(MAKEFILE_PATH)/commit_date.sh)

# Build native baron by default.
default: baron

dependencies:
	go get -u github.com/golang/protobuf/protoc-gen-go@21df5aa0e680850681b8643f0024f92d3b09930c
	go get -u github.com/golang/protobuf/protoc-gen-gofaster@21df5aa0e680850681b8643f0024f92d3b09930c
	go get -u github.com/golang/protobuf/proto@21df5aa0e680850681b8643f0024f92d3b09930c
	go get -u github.com/kevinburke/go-bindata/go-bindata

# Generate go files containing the all template files in []byte form
gobindata:
	go generate github.com/teamlint/baron/gengokit/template

# Install baron
baron: gobindata
	go install -ldflags '-X "main.version=$(SHA)" -X "main.date=$(VERSION_DATE)"' github.com/teamlint/baron/cmd/baron

# Run the go tests and the baron integration tests
test: test-go test-integration

test-go:
	#GO111MODULE=on go test -v ./... -covermode=atomic -coverprofile=./coverage.out --coverpkg=./...
	GO111MODULE=on go test -v ./...

test-integration:
	GO111MODULE=on $(MAKE) -C cmd/_integration-tests

# Removes generated code from tests
testclean:
	$(MAKE) -C cmd/_integration-tests clean

.PHONY: testclean test-integration test-go test baron gobindata dependencies
