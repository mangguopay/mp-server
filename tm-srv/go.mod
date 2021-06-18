module a.a/mp-server/tm-srv

go 1.13

require (
	a.a/cu v0.0.0-incompatible
	a.a/mp-server/common v0.0.0-incompatible
	github.com/micro/go-micro/v2 v2.6.0
)

replace (
	a.a/cu v0.0.0-incompatible => ../../cu
	a.a/mp-server/common v0.0.0-incompatible => ../common
)
