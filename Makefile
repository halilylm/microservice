.PHONY: cover start test test-integration

cover:
	go tool cover -html=cover.out

start:
	go run -ldflags="-X 'main.release=`git rev-parse --short=8 HEAD`'" app/*.go

test:
	go test -coverprofile=cover.out -short ./...

test-integration:
	go test -coverprofile=cover.out -p 1 ./...

build:
	go build -ldflags="-X 'main.release=`git rev-parse --short=8 HEAD`'" -o bin/server app/*.go