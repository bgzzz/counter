NAME = counter
GOENV = CGO_ENABLED=0
GO = $(GOENV) go

all: mod tools test-all build
.PHONY: all

mod:
	$(GO) mod download
mod-tidy:
	$(GO) mod tidy
.PHONY: mod mod-tidy

tools: 
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint
.PHONY: tools

test:
	$(GO) test -v ./...
lint:
	$(GOENV) golangci-lint run
test-all: test lint
.PHONY: test lint test-all

build: build-server build-ctl
.PHONY: build

build-server:
	$(GO) build -o ./bin/$(NAME) ./cmd/counter
.PHONY: build-server

build-ctl:
	$(GO) build -o ./bin/$(NAME)ctl ./cmd/counterctl
.PHONY: build-ctl

clean:
	$(GO) clean
	rm -rf bin
.PHONY: clean
