module a.a/mp-server/business-bill-srv

go 1.13

require (
	a.a/cu v0.0.0-incompatible
	a.a/mp-server/common v0.0.0-incompatible
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/micro/go-micro/v2 v2.6.0
	github.com/shopspring/decimal v0.0.0-20191129051706-bc70c3beb98b
)

replace a.a/cu v0.0.0-incompatible => ../../cu

replace a.a/net v0.0.0-incompatible => ../../net

replace a.a/mp-server/common v0.0.0-incompatible => ../common
