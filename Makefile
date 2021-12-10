GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=dcache
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_Win=$(BINARY_NAME).exe

all: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o $(BINARY_UNIX) cmd/server/main.go
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v  cmd/server/main.go
build-win:
	#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_Win) -v  cmd/server/main.go