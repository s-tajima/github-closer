setup:
	go get -u github.com/golang/dep/...
	dep ensure

build:
	go build

run:
	go run *.go

test:
	go test

travis:
	go fmt
	go test
