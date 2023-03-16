GOCMD = go
GOBUILD = $(GOCMD) build
GOMOD = $(GOCMD) mod
GOTEST = $(GOCMD) test

build:
	go build -o dapp-backend.exe ./cmd/

lint:
	golangci-lint run --fix

download:
	echo Download go.mod dependencies
	go mod download

install-tools: download
	echo Installing tools from tools.go
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %@latest

test:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...