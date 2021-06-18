module a.a/mp-server/common

go 1.13

require (
	a.a/cu v0.0.0-incompatible
	github.com/aws/aws-sdk-go v1.25.31
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5 // indirect
	github.com/gin-gonic/gin v1.5.0
	github.com/go-redis/redis/v7 v7.2.0
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
	github.com/golang/protobuf v1.3.5
	github.com/kellydunn/golang-geo v0.7.0
	github.com/kylelemons/go-gypsy v0.0.0-20160905020020-08cad365cd28 // indirect
	github.com/lib/pq v1.3.0
	github.com/micro/go-micro/v2 v2.6.0
	github.com/micro/go-plugins v1.5.1
	github.com/nats-io/nats-streaming-server v0.16.2 // indirect
	github.com/nats-io/stan.go v0.6.0
	github.com/pkg/errors v0.9.1
	github.com/shopspring/decimal v0.0.0-20191129051706-bc70c3beb98b
	github.com/ziutek/mymysql v1.5.4 // indirect
)

replace a.a/cu v0.0.0-incompatible => ../../cu
