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

DOCKER := $(shell which docker)
DOCKERCOMPOSE := $(shell which docker-compose)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf
COMMIT_HASH := $(shell git rev-parse --short=7 HEAD)
DOCKER_TAG := $(COMMIT_HASH)

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
