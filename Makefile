default: help

.PHONY: help
help: ## Print this help message
	@echo "Available make commands:"; grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

###############################################################################
# Build / Install
###############################################################################

.PHONY: build
build: go.sum ## Build the binary in ./build
ifeq ($(OS),Windows_NT)
	@echo "building habbgo binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/habbgo.exe ./habbgo
else
	@echo "building habbgo binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/habbgo ./habbgo
endif

.PHONY: install
install: ## Install the binary in go/bin
	@echo "installing habbgo binary..."
	@go generate
	@go install -race -mod readonly $(BUILD_FLAGS) .

.PHONY: build-docker
build-docker: ## Build a Docker image for the game server
	docker build -t jtieri/habbgo:latest -f ./docker/habbgo/Dockerfile .

.PHONY: clean
clean: ## Delete the ./build directory
	rm -rf build

###############################################################################
# Tests / CI
###############################################################################

.PHONY: test
test: ## Run Go tests
	@go test -mod readonly --race -v ./...

.PHONY: run-docker
run-docker: build-habbgo-docker ## Build the Docker image if needed and start container
	docker run jtieri/habbgo