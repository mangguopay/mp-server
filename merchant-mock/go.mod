module a.a/mp-server/merchant-mock

go 1.13

require (
	a.a/cu v0.0.0-incompatible
	a.a/mp-server/common v0.0.0-incompatible
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc
	github.com/gin-gonic/gin v1.5.0
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/lib/pq v1.3.0
	github.com/tuotoo/qrcode v0.0.0-20190222102259-ac9c44189bf2
)

replace a.a/cu v0.0.0-incompatible => ../../cu

replace a.a/net v0.0.0-incompatible => ../../net

replace a.a/mp-server/common v0.0.0-incompatible => ../common
