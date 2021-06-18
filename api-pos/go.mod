module a.a/mp-server/api-pos

go 1.13

require (
	a.a/cu v0.0.0-incompatible
	a.a/mp-server/common v0.0.0-incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.5.0
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
	github.com/json-iterator/go v1.1.9
	github.com/lib/pq v1.3.0
	github.com/micro/go-micro/v2 v2.6.0
	golang.org/x/tools v0.0.0-20191216173652-a0e659d51361
)

replace a.a/cu v0.0.0-incompatible => ../../cu

replace a.a/mp-server/common v0.0.0-incompatible => ../common
