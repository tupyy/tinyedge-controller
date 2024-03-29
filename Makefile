# VERSION defines the project version.
# Update this value when you upgrade the version of your project.
VERSION ?= 0.1.0

QUAY_REPO = quay.io/ctupangiu
RELEASE_REPO ?= tinyedge-controller
SKIP_TEST_IMAGE_PULL ?= false

# Image URL to use all building/pushing image targets
IMG ?= tinyedge-controller
IMG_TAG ?= latest

# Path to protoc compiler
PROTOC ?= $(PWD)/bin/protoc

# Docker command to use, can be podman
DOCKER ?= podman
DOCKER-COMPOSE ?= podman-compose 

OS = $(shell go env GOOS)
ARCH = $(shell go env GOARCH)
GO = go

# Set quiet mode by default
Q=@

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

# Colors used in this Makefile
escape=$(shell printf '\033')
RESET_COLOR=$(escape)[0m
COLOR_YELLOW=$(escape)[38;5;220m
COLOR_RED=$(escape)[91m
COLOR_BLUE=$(escape)[94m

COLOR_LEVEL_TRACE=$(escape)[38;5;87m
COLOR_LEVEL_DEBUG=$(escape)[38;5;87m
COLOR_LEVEL_INFO=$(escape)[92m
COLOR_LEVEL_WARN=$(escape)[38;5;208m
COLOR_LEVEL_ERROR=$(escape)[91m
COLOR_LEVEL_FATAL=$(escape)[91m

define COLORIZE
sed -u -e "s/\\\\\"/'/g; \
s/debug/$(COLOR_LEVEL_DEBUG)debug$(RESET_COLOR)/;    \
s/info/$(COLOR_LEVEL_INFO)info$(RESET_COLOR)/;       \
s/warning/$(COLOR_LEVEL_WARN)warning$(RESET_COLOR)/; \
s/error/$(COLOR_LEVEL_ERROR)error$(RESET_COLOR)/;    \
s/fatal/level=$(COLOR_LEVEL_FATAL)fatal$(RESET_COLOR)/"
endef

MOQ ?= ${LOCAL_BIN_PATH}/moq
moq:
ifeq (, $(shell which ${LOCAL_BIN_PATH}/moq 2> /dev/null))
	@{ \
	set -e ;\
	MOQ_TMP_DIR=$$(mktemp -d) ;\
	cd $$MOQ_TMP_DIR ;\
	$(GO) mod init tmp ;\
	$(GO) get -d github.com/matryer/moq@v0.2.7 ;\
	mkdir -p ${LOCAL_BIN_PATH} ;\
	$(GO) build -o ${LOCAL_BIN_PATH}/moq github.com/matryer/moq ;\
	rm -rf $$MOQ_TMP_DIR ;\
	}
endif

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
.PHONY: generate generate.models generate.proto generate.mock 
generate: generate.proto

DEST ?= .
generate.proto:
	$(PROTOC) -I=protocol --go_out=$(DEST) --go-grpc_opt=module=github.com/tupyy/tinyedge-controller --go_opt=module=github.com/tupyy/tinyedge-controller --go-grpc_out=$(DEST) protocol/edge.proto protocol/common.proto protocol/admin.proto

BASE_CONNSTR="postgresql://$(PG_USER):$(PG_PWD)@$(DB_HOST):$(DB_PORT)"
GEN_CMD=$(TOOLS_DIR)/gen --sqltype=postgres \
	--module=github.com/tupyy/tinyedge-controller/internal/repo/models/pg --exclude=schema_migrations \
	--gorm --no-json --no-xml --overwrite --out $(CURDIR)/internal/repo/models

#help generate.models: generate models for the database
generate.models:
	sh -c '$(GEN_CMD) --connstr "$(BASE_CONNSTR)/tinyedge?sslmode=disable"  --model=pg --database tinyedge' 						# Generate models for the DB tables

generate.mock:
	$(GO) generate ./...
	$(GO) mod vendor

GO_IMAGE=golang:1.17.8-alpine3.14
GOIMPORTS_IMAGE=golang.org/x/tools/cmd/goimports@latest
FILES_LIST=$(shell ls -d */ | grep -v -E "vendor|tools|test|client|restapi|models|generated")
MODULE_NAME=$(shell head -n 1 go.mod | cut -d '/' -f 3)
imports: ## fix and format go imports
	@# Removes blank lines within import block so that goimports does its magic in a deterministic way
	find $(FILES_LIST) -type f -name "*.go" | xargs -L 1 sed -i '/import (/,/)/{/import (/n;/)/!{/^$$/d}}'
	$(DOCKER) run --rm -v $(CURDIR):$(CURDIR):z -w="$(CURDIR)" $(GO_IMAGE) \
		sh -c 'go install $(GOIMPORTS_IMAGE) && goimports -w -local github.com/tupyy/tinyedge-controller $(FILES_LIST) && goimports -w -local github.com/tupyy/tinyedge-controller/$(MODULE_NAME) $(FILES_LIST)'

LINT_IMAGE=golangci/golangci-lint:v1.45.0
lint: ## Check if the go code is properly written, rules are in .golangci.yml 
	$(DOCKER) run --rm -v $(CURDIR):$(CURDIR) -w="$(CURDIR)" $(LINT_IMAGE) sh -c 'golangci-lint run'


##@ Build
.PHONY: build build.client run run.infra run.infra.stop docker.build docker.stop test ginkgo

build: ## Build binary.
	go build -mod=vendor -o $(PWD)/bin/tinyedge-controller $(PWD)/main.go

build.client: ## Build the client
	go build -mod=vendor -o $(PWD)/bin/tinyedge-cli $(PWD)/client/main.go

FLAGS=--vault_address "localhost:8200"
run: ## Run the controller from your host.
	bin/tinyedge-controller run $(FLAGS) | $(COLORIZE)

run.infra: podman.build.vault
	podman play kube $(CURDIR)/build/kube.yaml

run.infra.stop:
	podman kube down $(CURDIR)/build/kube.yaml
	podman image rm vault:latest

docker.build: ## Build docker image with the manager.
	$(DOCKER) build -f build/Dockerfile -t ${IMG}:${IMG_TAG} .

docker.push: ## Push docker image with the manager.
	$(DOCKER) tag ${IMG}:${IMG_TAG} ${QUAY_REPO}/${IMG}:${IMG_TAG}
	$(DOCKER) push ${IMG}:${IMG_TAG}

podman.build.vault:
	$(DOCKER) build -t vault -f build/vault/Dockerfile

test: ginkgo
	$(GINKGO) -focus=$(FOCUS) -v --cover --coverprofile=cover.out ./...

test-create-coverage:
	go tool cover -func cover.out
	go tool cover --html=cover.out -o coverage.html

test-coverage:
	go tool cover --html=cover.out 

##@ Infra
.PHONY: postgres.setup.clean postgres.setup.init postgres.setup.tables postgres.setup.migrations

DB_HOST ?= localhost
DB_PORT ?= 5433
PG_USER = postgres
PG_PWD = postgres
PGPASSFILE=$(CURDIR)/sql/.pgpass
PSQL_COMMAND=PGPASSFILE=$(PGPASSFILE) psql --quiet --host=$(DB_HOST) --port=$(DB_PORT) -v ON_ERROR_STOP=on --user=$(PG_USER)

#help postgres.setup: Setup postgres from scratch
postgres.setup: postgres.setup.init postgres.setup.tables postgres.setup.fixtures

#help postgres.setup.clean: cleans postgres from all created resources
postgres.setup.clean:
	$(PSQL_COMMAND) -f sql/clean.sql

#help postgres.setup.init: init the database
postgres.setup.init:
	$(PSQL_COMMAND) --dbname=postgres \
		-f sql/init.sql

#help postgres.setup.users: init postgres users
postgres.setup.tables:
	$(PSQL_COMMAND) --dbname=tinyedge \
		-f sql/tables.sql

postgres.setup.fixtures:
	$(PSQL_COMMAND) --dbname=tinyedge \
		-f sql/fixtures.sql

##@ Vault

VAULT_ADDR="http://localhost:8200"
VAULT_FORMAT=json
VAULT_CMD=VAULT_ADDR=$(VAULT_ADDR) VAULT_FORMAT=$(VAULT_FORMAT) vault

.PHONY: vault.login vault.secret.id

vault.login:
	$(VAULT_CMD) login root

vault.secret.id:
	@$(VAULT_CMD) write -f auth/approle/role/dev-role/secret-id | jq '.data.secret_id' | sed 's/"//g'

vault.generate.registration.certificates:
	@$(VAULT_CMD) write pki_int/issue/tinyedge-role common_name="registration.home.net" ttl="3600h" ip_sans="127.0.0.1" > registration
	@cat $(PWD)/registration | jq --raw-output '.data.private_key' > registration_private_key.pem
	@cat $(PWD)/registration | jq --raw-output '.data.certificate' > registration_cert.pem
	@cat $(PWD)/registration | jq --raw-output '.data.issuing_ca' > registration_ca.pem
	@rm $(PWD)/registration

vault.generate.server.certificates:
	@$(VAULT_CMD) write pki_int/issue/tinyedge-role common_name="server.home.net" ttl="3600h" ip_sans="127.0.0.1" > server
	@cat $(PWD)/server | jq --raw-output '.data.private_key' > $(PWD)/resources/certificates/key.pem
	@cat $(PWD)/server | jq --raw-output '.data.certificate' > $(PWD)/resources/certificates/cert.pem
	@cat $(PWD)/server | jq --raw-output '.data.issuing_ca' > $(PWD)/resources/certificates/ca.pem
	@rm $(PWD)/server

##@ Tools
TOOLS_DIR=$(CURDIR)/tools/bin

GINKGO = $(shell pwd)/bin/ginkgo
ginkgo: ## Download ginkgo locally if necessary.
ifeq (, $(shell which ginkgo 2> /dev/null))
	$(call go-install-tool,$(GINKGO),github.com/onsi/ginkgo/v2/ginkgo@v2.1.3)
endif

.PHONY: tools.clean tools.get

#help tools.clean: remove everything in the tools/bin directory
tools.clean:
	rm -fr $(TOOLS_DIR)/*

#help tools.get: retrieve all the tools specified in gex
tools.get:
	cd $(CURDIR)/tools && go generate tools.go

# go-install-tool will 'go install' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-install-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef
