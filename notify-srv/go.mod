module a.a/mp-server/notify-srv

go 1.13

require (
	a.a/cu v0.0.0-incompatible
	a.a/mp-server/common v0.0.0-incompatible
	a.a/net v0.0.0-incompatible
	github.com/gin-gonic/gin v1.5.0
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
	github.com/golang/protobuf v1.3.5
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/json-iterator/go v1.1.9
	github.com/lib/pq v1.3.0
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v2 v2.6.0
	github.com/micro/go-plugins v1.5.1
	github.com/nats-io/stan.go v0.6.0
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	github.com/shopspring/decimal v0.0.0-20191129051706-bc70c3beb98b
	github.com/wiwii/pool v0.0.0-20171030022714-e6b389d645e3
	github.com/wiwii/redis-go-cluster v0.0.0-20171117095555-0cc33a51bfec
	go.uber.org/atomic v1.5.0
	golang.org/x/crypto v0.0.0-20200323165209-0ec3e9974c59
)

replace a.a/cu v0.0.0-incompatible => ../../cu

replace a.a/net v0.0.0-incompatible => ../../net

replace a.a/mp-server/common v0.0.0-incompatible => ../common
