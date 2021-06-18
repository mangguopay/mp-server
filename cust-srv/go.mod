module a.a/mp-server/cust-srv

go 1.13

require (
	a.a/cu v0.0.0-incompatible
	a.a/mp-server/common v0.0.0-incompatible
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
	github.com/json-iterator/go v1.1.9
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/lib/pq v1.3.0
	github.com/micro/go-micro/v2 v2.6.0
	github.com/tealeg/xlsx v1.0.5
)

replace a.a/cu v0.0.0-incompatible => ../../cu

replace a.a/mp-server/common v0.0.0-incompatible => ../common

replace a.a/mp-server/auth-srv v0.0.0-incompatible => ../auth-srv
