
GOPATH:=$(shell go env GOPATH)


.PHONY: proto
proto:
	protoc --micro_out=. --go_out=. proto/cust/cust.proto
	protoc --micro_out=. --go_out=. proto/auth/auth.proto
	protoc --micro_out=. --go_out=. proto/push/push.proto
	protoc --micro_out=. --go_out=. proto/bill/bill.proto
	protoc --micro_out=. --go_out=. proto/settle/settle.proto
	protoc --micro_out=. --go_out=. proto/quota/quota.proto
	protoc --micro_out=. --go_out=. proto/riskctrl/riskctrl.proto
	protoc --micro_out=. --go_out=. proto/listen-exp-key/listen_exp_key.proto
	protoc --micro_out=. --go_out=. proto/statlog/statlog.proto
	protoc --micro_out=. --go_out=. proto/gis/gis.proto
	protoc --micro_out=. --go_out=. proto/tm/tm.proto
	protoc --micro_out=. --go_out=. proto/business-bill/business-bill.proto
	protoc --micro_out=. --go_out=. proto/business-settle/business-settle.proto
	protoc --micro_out=. --go_out=. proto/admin-auth/admin_auth.proto
	protoc --micro_out=. --go_out=. proto/notify/notify.proto

.PHONY: build
build: proto
    #go build -o api-srv *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t auth-srv:latest
