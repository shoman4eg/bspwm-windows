VERSION ?= 0.0.4
BIN_NAME ?= "bspcw"

GREEN='\033[0;32m'
NC='\033[0m'

all: build strip ## Run all

GO_MOD_NAME = github.com/shoman4eg/bspwm-windows

GIT_IMPORT = ${GO_MOD_NAME}/cmd
GIT_COMMIT = $(shell git rev-parse --short HEAD)

init: ## Install required tools
	@echo -e ${GREEN}[Init]${NC}

	@cd tools && go mod tidy && go generate -x -tags=tools

build: ## Build
	@echo -e ${GREEN}[Build]${NC}

	@go build -ldflags "-X $(GIT_IMPORT).version=$(VERSION) -X $(GIT_IMPORT).commit=$(GIT_COMMIT)" -o $(BIN_NAME)

strip: ## Strip
	@echo -e ${GREEN}[Srtip]${NC}

	strip $(BIN_NAME)

clean: ## Clean
	@echo -e ${GREEN}[Clean]${NC}

	rm -rf $(BIN_NAME)

help: ## Show this help screen
	@printf 'Usage: make \033[36m<TARGETS>\033[0m ... \033[36m<OPTIONS>\033[0m\n\nAvailable targets are:'
	@awk 'BEGIN {FS = ":.*##"; printf "\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "    \033[36m%-17s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

FILES = $(shell find . -type f -name '*.go')

format: ## Format source code
	@echo -e ${GREEN}[Format]${NC}

	@bin/gofumpt -l -w $(FILES)
	@bin/goimports -local ${GO_MOD_NAME} -l -w $(FILES)
	@bin/gci write --section Standard --section Default --section "Prefix(${GO_MOD_NAME})" $(FILES)

lint: ## Run required checkers and linters
	@echo -e ${GREEN}[Lint]${NC}

	@LOG_LEVEL=error bin/golangci-lint run
	@bin/go-consistent -pedantic ./...