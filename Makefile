.PHONY: clean run build install dep test lint format docker

SHELL := /bin/bash
PATHINSTBIN = $(abspath ./bin)
export PATH := $(PATHINSTBIN):$(PATH)

BIN_NAME					?= rewards-api
DEFAULT_INSTALL_DIR			:= $(go env GOPATH)/$(PATHINSTBIN)
DEFAULT_ARCH				:= $(shell go env GOARCH)
DEFAULT_GOOS				:= $(shell go env GOOS)
ARCH						?= $(DEFAULT_ARCH)
GOOS						?= $(DEFAULT_GOOS)
INSTALL_DIR					?= $(DEFAULT_INSTALL_DIR)
.DEFAULT_GOAL 				:= run


VERSION   := $(shell git describe --tags || echo "v0.0.0")
VER_CUT   := $(shell echo $(VERSION) | cut -c2-)

# Dependency versions
GOLANGCI_VERSION = latest

help:
	@echo "\nSpecify a subcommand:\n"
	@grep -hE '^[0-9a-zA-Z_-]+:.*?## .*$$' ${MAKEFILE_LIST} | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[0;36m%-20s\033[m %s\n", $$1, $$2}'
	@echo ""

build:
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(ARCH) \
		go build -o $(PATHINSTBIN)/$(BIN_NAME) ./cmd/$(BIN_NAME)

run: build
	@./$(PATHINSTBIN)/$(BIN_NAME)
all: clean target

clean:
	@rm -rf $(PATHINSTBIN)
	
install: build
	@install -d $(INSTALL_DIR)
	@rm -f $(INSTALL_DIR)/$(BIN_NAME)
	@cp $(PATHINSTBIN)/$(BIN_NAME) $(INSTALL_DIR)/$(BIN_NAME)

tidy: 
	@go mod tidy

test: ## run tests
	@go test ./...

lint: ## run linter
	@PATH=$$PATH golangci-lint --version
	@PATH=$$PATH golangci-lint run --timeout 10m

docker: dep ## build docker image
	@docker build -f ./docker/dockerfile . -t dimozone/$(BIN_NAME):$(VER_CUT)
	@docker tag dimozone/$(BIN_NAME):$(VER_CUT) dimozone/$(BIN_NAME):latest

tools-golangci-lint: ## install golangci-lint
	@mkdir -p $(PATHINSTBIN)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | BINARY=golangci-lint bash -s -- ${GOLANGCI_VERSION}


make tools: tools-golangci-lint 
	GOBIN=$(PATHINSTBIN) go install tool

generate: generate-swagger generate-go generate-sqlboiler ## run all file generation for the project

generate-swagger: ## generate swagger documentation
	@go tool swag -version
	go tool swag init -g cmd/$(BIN_NAME)/main.go --parseDependency --parseInternal 

generate-go:## run go generate
	@go generate ./...

generate-sqlboiler: build
	docker compose up -d postgresql
	sleep 5
	DB_USER=dimo DB_PASSWORD=dimo DB_HOST=localhost DB_PORT=5432 DB_NAME=rewards_api $(PATHINSTBIN)/$(BIN_NAME) migrate
	@PATH=$$PATH sqlboiler psql --no-tests --wipe
	docker compose down

