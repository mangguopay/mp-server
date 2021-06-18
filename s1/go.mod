module a.a/mp-server/s1

go 1.13

require (
	a.a/cu v0.0.0-incompatible
	a.a/mp-server/common v0.0.0-incompatible
	a.a/net v0.0.0-incompatible
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
	github.com/golang/protobuf v1.3.5
	github.com/hashicorp/consul v1.6.2 // indirect
	github.com/json-iterator/go v1.1.9
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v2 v2.6.0
	github.com/micro/go-plugins v1.5.1
	github.com/qiniu/qlang v0.0.0-20190907152943-bae5c3dacbdf
	github.com/shopspring/decimal v0.0.0-20191129051706-bc70c3beb98b
	github.com/tealeg/xlsx v1.0.5
	github.com/yuin/gopher-lua v0.0.0-20191213034115-f46add6fdb5c
	google.golang.org/genproto v0.0.0-20191216164720-4f79533eabd1
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)

replace a.a/cu v0.0.0-incompatible => ../../cu

replace a.a/mp-server/common v0.0.0-incompatible => ../common

replace a.a/net v0.0.0-incompatible => ../../net
