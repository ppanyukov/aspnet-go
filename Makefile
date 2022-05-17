FILES_TO_FMT      ?= $(shell find . -path ./vendor -prune -o -name '*.go' -print)

# Ensure everything works even if GOPATH is not set, which is often the case.
# The `go env GOPATH` will work for all cases for Go 1.8+.
GOPATH            ?= $(shell go env GOPATH)

# TMP_GOPATH is used to install tools with specified version
# in a way that does not mess up with go modules.
TMP_GOPATH        ?= /tmp/github.com/ppanyukov/aspnet-go

# Tooling and target OS/ARCH for the binary.
# Assists builds in docker for Windows users too.
#
# See build.sh
#
# Z_GOOS, Z_GOARCH: the target for the binary.
# GOBIN: location where tools like promu and goimports should be installed.
#
# By default, hosts GOOS and GOARCH will be used.
# Override to target other. E.g.
#   Z_GOOS=window make build
#   Z_GOOS=linux Z_GOARCH=arm make build
#
Z_GOOS            ?= $(shell go env GOOS)
Z_GOARCH          ?= $(shell go env GOARCH)

Z_GOBIN           ?= $(shell pwd)/bin_tools/$(shell go env GOOS)_$(shell go env GOARCH)
GOBIN             = $(Z_GOBIN)


#GO111MODULE       ?= on
#export GO111MODULE
#GOPROXY           ?= https://proxy.golang.org
#export GOPROXY

# Tools.
EMBEDMD           ?= $(GOBIN)/embedmd-$(EMBEDMD_VERSION)
# v2.0.0
EMBEDMD_VERSION   ?= 97c13d6e41602fc6e397eb51c45f38069371a969

GOIMPORTS         ?= $(GOBIN)/goimports-$(GOIMPORTS_VERSION)
GOIMPORTS_VERSION ?= v0.1.10

GOLANGCILINT_VERSION ?= v1.26.0
GOLANGCILINT         ?= $(GOBIN)/golangci-lint-$(GOLANGCILINT_VERSION)

MISSPELL_VERSION     ?= c0b55c8239520f6b5aa15a0207ca8b28027ba49e
MISSPELL             ?= $(GOBIN)/misspell-$(MISSPELL_VERSION)


# The version using modules. Fasterer and betterer.
define fetch_go_bin_version
	@mkdir -p $(GOBIN)
	@mkdir -p $(TMP_GOPATH)

	@echo ">> fetching $(1)@$(2) revision/version"
	cd '$(TMP_GOPATH)' && GOPATH='$(TMP_GOPATH)' GOBIN='$(TMP_GOPATH)/bin' GO111MODULE=on go install '$(1)@$(2)'
	mv -- '$(TMP_GOPATH)/bin/$(shell basename $(1))' '$(GOBIN)/$(shell basename $(1))-$(2)'
endef

.PHONY: all
all: format build test lint

# build builds binaries using `promu`.
.PHONY: build
build: deps
	@echo ">> building..."
	GOOS=$(Z_GOOS) GOARCH=$(Z_GOARCH) go build ...

# same as build but not checking git or modules
.PHONY: build-fast
build-fast:
	@echo ">> building-fast..."
	GOOS=$(Z_GOOS) GOARCH=$(Z_GOARCH) go build ...

# deps ensures fresh go.mod and go.sum.
.PHONY: deps
deps:
	go mod tidy
	go mod verify

# format formats the code (including imports format).
.PHONY: format
format: $(GOIMPORTS)
	@echo ">> formatting code"
	$(GOIMPORTS) -w $(FILES_TO_FMT)

# test runs all standard tests
.PHONY: test
test:
	go test $(shell go list ./... | grep -v /vendor/);

.PHONY: lint
# PROTIP:
# Add
#      --cpu-profile-path string   Path to CPU profile output file
#      --mem-profile-path string   Path to memory profile output file
#
# to debug big allocations during linting.
lint: $(GOLANGCILINT) $(MISSPELL)
	@echo ">> linting all of the Go files GOGC=${GOGC}"
	$(GOLANGCILINT) run --enable goimports --enable goconst --skip-dirs vendor
	@echo ">> detecting misspells"
	find . -type f | grep -v vendor/ | grep -v '.csv' | grep -v go.mod | grep -v go.sum | grep -vE '\./\..*' | xargs $(MISSPELL) -error

# non-phony targets
$(EMBEDMD):
	$(call fetch_go_bin_version,github.com/campoy/embedmd,$(EMBEDMD_VERSION))

$(GOIMPORTS):
	$(call fetch_go_bin_version,golang.org/x/tools/cmd/goimports,$(GOIMPORTS_VERSION))

$(GOLANGCILINT):
	$(call fetch_go_bin_version,github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCILINT_VERSION))

$(MISSPELL):
	$(call fetch_go_bin_version,github.com/client9/misspell/cmd/misspell,$(MISSPELL_VERSION))

