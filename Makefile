# easyGZH Makefile
# Single-binary Go CLI build, test, and cross-platform release.

BINARY   := easygzh
VERSION  := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT   := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE     := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS  := -s -w -X main.Version=$(VERSION)
GO       := go

# Cross-compile targets: GOOS_GOARCH
TARGETS := darwin_arm64 darwin_amd64 linux_arm64 linux_amd64 windows_amd64

.PHONY: all generate build test clean install release help

all: build

## generate: refresh generated Go assets from canonical source files
generate:
	node scripts/generate-memory-scaffold-go.mjs

## build: compile the binary into ./easygzh
build: generate
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/easygzh

## test: run all tests
test: generate
	$(GO) test ./...

## vet: go vet
vet:
	$(GO) vet ./...

## tidy: go mod tidy
tidy:
	$(GO) mod tidy

## install: install to $$GOBIN (or $$GOPATH/bin)
install: generate
	$(GO) install -ldflags "$(LDFLAGS)" ./cmd/easygzh

## clean: remove built binary and dist/
clean:
	rm -f $(BINARY)
	rm -rf dist/

## release: cross-compile into dist/<target>/easygzh[.exe]
release: clean generate
	@mkdir -p dist
	@set -e; for target in $(TARGETS); do \
		os=$${target%_*}; arch=$${target#*_}; \
		ext=""; [ $$os = windows ] && ext=".exe"; \
		echo "→ $$os/$$arch"; \
		GOOS=$$os GOARCH=$$arch $(GO) build -ldflags "$(LDFLAGS)" \
			-o dist/$(BINARY)-$$target$$ext ./cmd/easygzh; \
	done
	@ls -lh dist/

## golden: regenerate golden fixtures from Node juice (requires node)
golden:
	@if command -v node >/dev/null 2>&1; then \
		node scripts/_gen-golden.mjs; \
	else echo "node not found; skipping golden regen"; fi

help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //'
