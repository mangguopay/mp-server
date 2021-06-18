module a.a/mp-server/api-webadmin

go 1.13

require (
	a.a/cu v0.0.0-incompatible
	a.a/mp-server/common v0.0.0-incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.5.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/json-iterator/go v1.1.9
	github.com/lib/pq v1.3.0
	github.com/micro/go-micro/v2 v2.6.0
	github.com/tealeg/xlsx v1.0.5
	github.com/wiwii/base64Captcha v0.0.0-20190617051550-75921257df0f
)

replace a.a/cu v0.0.0-incompatible => ../../cu

replace a.a/mp-server/common v0.0.0-incompatible => ../common
