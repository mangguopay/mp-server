module rsaTest

go 1.13

require (
	a.a/cu v0.0.0-incompatible
	a.a/mp-server/common v0.0.0-incompatible
	github.com/360EntSecGroup-Skylar/excelize v1.4.1
	github.com/shopspring/decimal v0.0.0-20191129051706-bc70c3beb98b

)

replace a.a/cu v0.0.0-incompatible => ../../cu

replace a.a/mp-server/common v0.0.0-incompatible => ../common
