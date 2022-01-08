#VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')

all: install

###############################################################################
# Build / Install
###############################################################################

build: go.sum
ifeq ($(OS),Windows_NT)
	@echo "building habbgo binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/habbgo.exe ./habbgo
else
	@echo "building habbgo binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/habbgo ./habbgo
endif

install:
	@echo "installing habbgo binary..."
	@go install -mod readonly $(BUILD_FLAGS) .

build-habbgo-docker:
	docker build -t jtieri/habbgo:latest -f ./docker/habbgo/Dockerfile .

clean:
	rm -rf build

###############################################################################
# Tests / CI
###############################################################################

test:
	@go test -mod readonly -v ./...

run-docker:
	docker build -t jtieri/habbgo:latest -f ./docker/habbgo/Dockerfile .
	docker run jtieri/habbgo