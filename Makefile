current_dir = $(shell pwd)
GOPATH=$(current_dir)
install:
	
	go get github.com/gin-gonic/gin
	go get github.com/fvbock/endless
	go get github.com/lib/pq
	
	createuser -s -d -r -l -i postgres
	createdb -U postgres kiss
	
build:
	gofmt -w  src/kiss/*
	go install -v kiss/cmd/kiss
test:
	go test -v kiss/...