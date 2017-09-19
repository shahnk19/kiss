current_dir = $(shell pwd)
GOPATH=$(current_dir)
install:
	go get github.com/gin-gonic/gin
	go get github.com/fvbock/endless
build:
	gofmt -w  src/kiss/*
	go install -v kiss/cmd/kiss
test:
	go test -v kiss/...