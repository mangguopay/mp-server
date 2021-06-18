package middleware

import (
	"fmt"
	"testing"
	"time"

	"a.a/cu/strext"
	"a.a/mp-server/common/global"

	"a.a/cu/jwt"
)

func TestJwtVerifyMw_VerifyToken(t *testing.T) {

	z, zErr := time.LoadLocation(strext.ToStringNoPoint("Asia/Phnom_Penh"))
	if zErr != nil {
		panic(fmt.Sprintf("解析时区出错,err: %v", zErr))
	}
	// 设置time包中的默认时区
	time.Local = z

	global.Tz = z

	value := "4V/MKP/DpvtcET2t/xKO1cKOHjz2rbkhSV7a0S2/bFnGWmb902vbPXDfewE3gRusdCteTKQuZ4h0owIf2vkuKeoPKKuYEt2XwBY9H0G48ohKwY8kPqQNrLNLC0tfsAQWIwJJNUvwvKEevxK98AW/nBenMVGyjdk29Xxo+WA84AwKyMVCypVEaqSwFaXOZd/zT0p3vrT8ogbF08fhJQPffpdfELgctQXRRR9udUuYQS7EfmS7mXtSg2qBLFsgVWi26qJCJU1vK0LnQt7Ql0XV9uUVM4PWNWT3lSXyuBoMJSGPV0Kg175By+3HgR3qw68Ei6p5NPSEM1fBNu8ImmN/gHz6CTgnEqZXwXIzKWN4Zr5Gjq3UmY2GDjeyJL3X6u4pSw/c5G3HoGdcCI09mQNtafs4nk9s4yAvVNURn8VnBG55SMYsYYAno/w5lxj7dLf7ZmlQy1SvBY+pa8eJVjuxU1i7G3brtRbzKC5w3jm0OPJN9z7xx+5GhFD9sFIL3SR2GLOGjEb+h+gH/pJHIchAXfogLDOwh4bXPUsh4JetHFksx1LXYJdA3QfckhnhlNxeSfCHPDEJC1/Ikh6KCIzCnHjK56JpPTNG0PRZAT4smtRyBvTPcnRbK2jXqK3gudQ42zNpLrEly7qEBe7lU4GE5waZ0y9+X+oZjMvGI0Ko5nLXj0LujHrYWJE37aC06B2dAdNpup0BJmEHH2Eylh7+aHawBROARCPWP7oi8x60r4Nc9VgnSP7ix+ijgeDIx+CS8tBNhX+G6+bWVZ2/axyrpe1qDZ+vM8HY8ZoqGPQPOI+zlEXF5E8/ydu+zI2asvnJ"

	loginAesKey := "1234567890123456"
	loginSignKey := ""

	isOk := jwt.ValidateEncryptedJWT(value, loginAesKey, loginSignKey)

	t.Logf("isOk:%v", isOk)
}
