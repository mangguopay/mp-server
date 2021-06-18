module a.a/mp-server/conf-srv

go 1.13

require (
	a.a/cu v0.0.0-incompatible
	a.a/mp-server/common v0.0.0-incompatible
	github.com/lib/pq v1.3.0
	github.com/micro/go-micro v1.16.0
	github.com/micro/go-micro/v2 v2.6.0
)

replace a.a/cu v0.0.0-incompatible => ../../cu

replace a.a/mp-server/common v0.0.0-incompatible => ../common
