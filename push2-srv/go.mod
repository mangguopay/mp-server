module a.a/mp-server/push2-srv

go 1.13

require (
	a.a/cu v0.0.0-incompatible
	a.a/mp-server/common v0.0.0-incompatible
	a.a/net v0.0.0-incompatible
	github.com/NaySoftware/go-fcm v0.0.0-20190516140123-808e978ddcd2
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/micro/go-micro/v2 v2.6.0
)

replace a.a/cu v0.0.0-incompatible => ../../cu

replace a.a/mp-server/common v0.0.0-incompatible => ../common

replace a.a/net v0.0.0-incompatible => ../../net
