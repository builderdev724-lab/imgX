BINARY := imgx
PKG := ./cmd/imgx
BIN_DIR := bin

.PHONY: all build test race vet fmt check clean snapshot release install

all: check build

build:
	mkdir -p $(BIN_DIR)
	go build -trimpath -ldflags='-s -w' -o $(BIN_DIR)/$(BINARY) $(PKG)

test:
	go test ./...

race:
	go test -race ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

check: fmt vet test

audit:
	go vet ./...
	go test ./...
	go mod tidy

snapshot:
	goreleaser release --snapshot --clean

release:
	goreleaser release --clean

install: build
	install -Dm755 $(BIN_DIR)/$(BINARY) $(DESTDIR)/usr/local/bin/$(BINARY)

clean:
	rm -rf $(BIN_DIR) dist coverage.out
