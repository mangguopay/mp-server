module a.a/mp-server/api-pay

go 1.13

require (
	a.a/cu v0.0.0-incompatible
	a.a/mp-server/common v0.0.0-incompatible
	a.a/net v0.0.0-incompatible
	github.com/gin-gonic/gin v1.5.0
	github.com/micro/go-micro v1.18.0 // indirect
	github.com/micro/go-micro/v2 v2.6.0
)

replace a.a/cu v0.0.0-incompatible => ../../cu

replace a.a/net v0.0.0-incompatible => ../../net

replace a.a/mp-server/common v0.0.0-incompatible => ../common
