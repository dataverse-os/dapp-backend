GOCMD = go
GOBUILD = $(GOCMD) build
GOMOD = $(GOCMD) mod
GOTEST = $(GOCMD) test

build-all: generate-js build-rs build-go

build-go:
	LD_LIBRARY_PATH=$LD_LIBRARY_PATH:./lib/ go build -o dapp-backend ./cmd/

lint:
	golangci-lint run --fix

build-rs:
	cd rs-binding && cargo build --release
ifeq ($(shell uname),Darwin)
	cp rs-binding/target/release/librs_binding.dylib ./lib
else
	cp rs-binding/target/release/librs_binding.so ./lib
endif
	cp rs-binding/target/rs-binding.h ./lib

generate-js:
	cd js-scripts && pnpm install && pnpm run build

download:
	echo Download go.mod dependencies
	go mod download

install-tools: download
	echo Installing tools from tools.go
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %@latest

test: test-go test-rs

test-go:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...

test-rs:
	cd rs-binding && cargo test

.PHONY: check-clippy
check-clippy:
	# Check with default features
	cd rs-binding && cargo clippy --all-targets -- -D warnings
	# Check with all features
	cd rs-binding && cargo clippy --all-targets --all-features -- -D warnings