current_dir = $(shell pwd)
GOPATH=$(current_dir)
all: install build
install:
	go get github.com/gin-gonic/gin
	go get github.com/fvbock/endless
	go get github.com/lib/pq
	go get github.com/gin-gonic/contrib/renders/multitemplate
	go get github.com/mediocregopher/radix.v2
	
build:
	gofmt -w  src/kiss/*
	go install -v kiss/cmd/kiss
test:
	go test -v kiss/...
bundle:
	git archive -v -o myapp.zip --format=zip HEAD