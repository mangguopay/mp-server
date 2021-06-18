module a.a/mp-server/bill-srv

go 1.13

require (
	a.a/cu v0.0.0-incompatible
	a.a/mp-server/common v0.0.0-incompatible
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
	github.com/json-iterator/go v1.1.9
	github.com/kellydunn/golang-geo v0.7.0
	github.com/micro/go-micro/v2 v2.6.0
	github.com/micro/go-plugins v1.5.1
	github.com/nats-io/stan.go v0.6.0
	github.com/shopspring/decimal v0.0.0-20191129051706-bc70c3beb98b
	github.com/tealeg/xlsx v1.0.5
)

replace a.a/cu v0.0.0-incompatible => ../../cu

replace a.a/mp-server/common v0.0.0-incompatible => ../common
