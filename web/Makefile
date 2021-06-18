
GOPATH:=$(shell go env GOPATH)

.PHONY: proto
proto:
	protoc --micro_out=. --go_out=. proto/web/web.proto

.PHONY: build
build:

	go build -o web-web *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t web-web:latest
