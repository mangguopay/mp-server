
GOPATH:=$(shell go env GOPATH)


.PHONY: proto
proto:
	protoc --micro_out=. --go_out=. proto/tm/tm.proto

.PHONY: build
build: proto

	go build -o tm-srv *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t tm-srv:latest
