#!/usr/bin/make -f

DOCKER_BUILDKIT=1
COSMOS_BUILD_OPTIONS ?= ""
TMVERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
BINDIR ?= $(GOPATH)/bin
# QS_BINARY = quicksilverd
# QS_DIR = quicksilver
BUILDDIR ?= $(CURDIR)/build

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=stride \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=icademod \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=netgo"

BUILD_FLAGS := -tags "netgo" -ldflags "$(ldflags)" -trimpath -mod=readonly -modcacherw

DOCKER := $(shell which docker)
DOCKERCOMPOSE := $(shell which docker-compose)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf
COMMIT_HASH := $(shell git rev-parse --short=7 HEAD)
DOCKER_TAG := $(COMMIT_HASH)

###############################################################################
###                       Install, Build & Clean                            ###
###############################################################################
install: go.sum
		@echo "--> Installing icademod"
		@go install $(BUILD_FLAGS) ./cmd/icademod

build: go.sum
		@echo "--> Building icademod"
		@go build $(BUILD_FLAGS) -o $(BUILDDIR)/icademod ./cmd/icademod

clean:
	rm -rf $(BUILDDIR)


###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=0.15.0
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-gen: 
	@echo "🤖 Generating code from protobuf..."
	@$(protoImage) sh ./scripts/proto-gen.sh
	@echo "✅ Completed code generation!"

proto-lint:
	@echo "🤖 Running protobuf linter..."
	@$(protoImage) buf lint
	@echo "✅ Completed protobuf linting!"

proto-format:
	@echo "🤖 Running protobuf format..."
	@$(protoImage) buf format -w
	@echo "✅ Completed protobuf format!"

proto-breaking-check:
	@echo "🤖 Running protobuf breaking check against main branch..."
	@$(protoImage) buf breaking --against '.git#branch=main'
	@echo "✅ Completed protobuf breaking check!"


###############################################################################
###                                Initialize                               ###
###############################################################################

init-golang-rly: kill-dev install
	@echo "Initializing both blockchains..."
	./network/init.sh
	./network/start.sh
	@echo "Initializing relayer..."
	./network/relayer/rly-init.sh

start: 
	@echo "Starting up test network"
	./network/start.sh

start-golang-rly:
	./network/relayer/rly-start.sh

kill-dev:
	@echo "Killing icademod and removing previous data"
	-@rm -rf ./data
	-@killall icademod 2>/dev/null
