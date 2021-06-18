package ss_count

import (
	"a.a/cu/ss_big"
	geo "github.com/kellydunn/golang-geo"
	"github.com/shopspring/decimal"
)

// 计算手续费 结果已经除了100倍,
func CountFees(amount, rate, minDefaultFees string) decimal.Decimal {
	amountDeci, _ := decimal.NewFromString(amount)
	rateDeci, _ := decimal.NewFromString(rate)
	dd, _ := decimal.NewFromString("10000")
	defFeeDeci, _ := decimal.NewFromString(minDefaultFees)

	//实际手续费
	feesDeci := amountDeci.Mul(rateDeci).Div(dd)

	res, _ := defFeeDeci.Sub(feesDeci).Float64()
	if res > 0 {
		return defFeeDeci
	}
	//return (amount * rate) / 10000
	//四舍五入
	return feesDeci.Round(0)
}

func CountSharing(fees, rate string) (string, string) {
	feesDeci, _ := decimal.NewFromString(fees)
	rateDeci, _ := decimal.NewFromString(rate)
	dd, _ := decimal.NewFromString("10000")
	srvFeesDeci := rateDeci.Div(dd).Mul(feesDeci)
	// 取整
	srvFeesDeci = ss_big.SsBigInst.ToRound(srvFeesDeci, 0, ss_big.RoundingMode_HALF_EVEN)

	headFeesDeci := feesDeci.Sub(srvFeesDeci)
	return headFeesDeci.String(), srvFeesDeci.String()
}

// 加法
func Add(a, b string) string {
	switch "" {
	case a:
		a = "0"
	case b:
		b = "0"
	}
	aDeci, _ := decimal.NewFromString(a)
	bDeci, _ := decimal.NewFromString(b)
	return aDeci.Add(bDeci).String()
}

// 除法
func Div(a, b string) decimal.Decimal {
	switch "" {
	case a:
		a = "0"
	}
	aDeci, _ := decimal.NewFromString(a)
	bDeci, _ := decimal.NewFromString(b)
	return aDeci.Div(bDeci)
}

// 乘法
func Multiply(a, b string) decimal.Decimal {
	switch "" {
	case a:
		a = "0"
	case b:
		b = "0"
	}

	aDeci, _ := decimal.NewFromString(a)
	bDeci, _ := decimal.NewFromString(b)
	return aDeci.Mul(bDeci)
}

// 减法
func Sub(a, b string) decimal.Decimal {
	switch "" {
	case a:
		a = "0"
	case b:
		b = "0"
	}
	aDeci, _ := decimal.NewFromString(a)
	bDeci, _ := decimal.NewFromString(b)
	return aDeci.Sub(bDeci)
}

func CountCircleDistance(lat1, lng1, lat2, lng2 float64) float64 {
	// Make a few points
	p := geo.NewPoint(lat1, lng1)
	p2 := geo.NewPoint(lat2, lng2)

	// find the great circle distance between them
	dist := p.GreatCircleDistance(p2)
	return dist
}
