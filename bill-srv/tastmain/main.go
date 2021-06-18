package main

import (
	"a.a/cu/encrypt"
	"a.a/cu/strext"
	"encoding/base64"
	"fmt"
)

func main1() {
	dataByte := `{
    "page":1,
    "page_size":10
}`

	//// 加密
	//encr,err := encrypt.DoRsa(encrypt.HANDLE_ENCRYPT,encrypt.KEYTYPE_PKIX,encrypt.RETFMT_HEX,encrypt.HASHLENTYPE_SHA1,
	//	[]byte(dataByte),pub,encrypt.KEYFMT_PEM,encrypt.PRE_HANDLE_BYTE)
	//if err != nil {
	//	fmt.Println("============加密失败")
	//	return
	//}
	//fmt.Println(encr.(string))
	//
	//
	//// 解密
	//dencr,derr := encrypt.DoRsa(encrypt.HANDLE_DECRYPT,encrypt.KEYTYPE_PKCS8,encrypt.RETFMT_BYTE,encrypt.HASHLENTYPE_SHA1,
	//	encr,pri,encrypt.KEYFMT_PEM,encrypt.PRE_HANDLE_HEX)
	//if derr != nil {
	//	fmt.Println("============解密失败")
	//	return
	//}
	//fmt.Println(string(dencr.([]byte)))
	//=======================================================================================

	dataByteT := base64.StdEncoding.EncodeToString([]byte(dataByte))

	// 加密
	encr, err := encrypt.DoRsa(encrypt.HANDLE_ENCRYPT, encrypt.KEYTYPE_PKIX, encrypt.RETFMT_BASE64, encrypt.HASHLENTYPE_NONE,
		[]byte(dataByteT), pub, encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_BYTE)
	if err != nil {
		fmt.Println("============加密失败")
		return
	}

	fmt.Println("密文--->", strext.ToStringNoPoint(encr))

	dataByteMd5 := encrypt.DoMd5(dataByteT)
	sign, err := encrypt.DoRsa(encrypt.HANDLE_SIGN, encrypt.KEYTYPE_PKCS8, encrypt.RETFMT_BASE64, encrypt.HASHLENTYPE_SHA256,
		[]byte(dataByteMd5), pri, encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_BYTE)
	if err != nil {
		fmt.Println("============签名失败", err.Error())
		return
	}

	fmt.Println("签名--->", strext.ToStringNoPoint(sign))

	//=============================================
	// 发过来
	//encr,sign
	// 解密
	dencr, derr := encrypt.DoRsa(encrypt.HANDLE_DECRYPT, encrypt.KEYTYPE_PKCS8, encrypt.RETFMT_STRING, encrypt.HASHLENTYPE_NONE,
		encr, pri, encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_BASE64)
	if derr != nil {
		fmt.Println("============解密失败")
		return
	}
	dencrMd5 := encrypt.DoMd5(strext.ToStringNoPoint(dencr))
	verifySign, err := encrypt.DoRsa(encrypt.HANDLE_VERIFY, encrypt.KEYTYPE_PKCS8, encrypt.RETFMT_NONE, encrypt.HASHLENTYPE_SHA256,
		map[string]interface{}{
			"sign": sign,
			"data": []byte(dencrMd5),
		}, pub, encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_SIGN_MAP_BASE64)
	if err != nil {
		fmt.Println("===========验证签名失败", err.Error())
		return
	}
	fmt.Println("============verifySign", verifySign)

	result, _ := base64.StdEncoding.DecodeString(dencr.(string))

	fmt.Println(strext.ToStringNoPoint(result))
}

var pri = `MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAMaeqEboJWW1a847GGdl+rXxCzrFo/LfX6hnmW0KCo0g7UrGyWkdAQgxSPVhWDzZtNdJZf9FSvtWt3Mkl405gUo2gCgdem4N/lEfzMMwS57vPpguv85Q4SEq8phBYbH1F0qSsJ+upQLLRFmUW+6uIIuG5pABlWwSEljHBaP7iCMvAgMBAAECgYEAmwn7xyWtE4CTP29HtGaZVm5q2pyuRoZwsA72Z3QgGlvqfqU/UQqv5Io8LjanXxB9plOIB8Q6LIlbN8kHd9i8fyQ8a8Y4RQIxNFGqfJ3lf4XjkOs9d/+kLLEevcYr0KlI3V2a56SprYW7LLzLAJXwK/nq24X/M0QwImtJWr42j5kCQQDI7F6/a4R15IX/LctgmfjZ5OGWzTPNShQfRroNuApJTNEICyMrNBjCnTS48a+VFeQCEwhjMPCAs7jRnp2IEeULAkEA/RCi6rsuVu4RtlzZ386Q4jeQxxr3Xc/8s0zLLR/QpmTNFuNKw4nN3jvZ/Uee2TlncIj63cHR276+FjicbyRI7QJAUBHKTFRHjEfOknuocc3KWuMYd2U9QJFF5ZTk7jSqfL2NC7yMflobh+roKM+/3hTEMYNuM0E8hr2YaIjiVGh1MwJAFzzp2OgrTyw5UCeikhyjzUIQ91eQk3q/168bkR80x7LF6m4gtWf4EYopcEqdWZEd4IWTk71yid0wE1ZLdyE72QJBAJ2TUwsgS/Xr93AMLjj2rN4vBo2WFehBOHBOsKbslc/dATZLZnCi8AfZCTEt82cZJIHsMWBEwubBeV+0BWViVsg=`

//
var pub = `MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCvf5bTBdJMhCJkZ4FQLQwXfY4RuMXFbmQQqPNr5P6f2CMoGNLiuzC0VbyOqz9u7/rt8RH3VuF3iXcxvwCLxaVruANxc64V9n6RGYUGpyDcqu0jF0SUHVSHevUNJN6kvEIuvyGwSAtowCRTsaqD2k+hDJoE/oANBPA/rpG/p9g7HQIDAQAB`

/*
var pri = `-----BEGIN RSA PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAMaeqEboJWW1a847
GGdl+rXxCzrFo/LfX6hnmW0KCo0g7UrGyWkdAQgxSPVhWDzZtNdJZf9FSvtWt3Mk
l405gUo2gCgdem4N/lEfzMMwS57vPpguv85Q4SEq8phBYbH1F0qSsJ+upQLLRFmU
W+6uIIuG5pABlWwSEljHBaP7iCMvAgMBAAECgYEAmwn7xyWtE4CTP29HtGaZVm5q
2pyuRoZwsA72Z3QgGlvqfqU/UQqv5Io8LjanXxB9plOIB8Q6LIlbN8kHd9i8fyQ8
a8Y4RQIxNFGqfJ3lf4XjkOs9d/+kLLEevcYr0KlI3V2a56SprYW7LLzLAJXwK/nq
24X/M0QwImtJWr42j5kCQQDI7F6/a4R15IX/LctgmfjZ5OGWzTPNShQfRroNuApJ
TNEICyMrNBjCnTS48a+VFeQCEwhjMPCAs7jRnp2IEeULAkEA/RCi6rsuVu4RtlzZ
386Q4jeQxxr3Xc/8s0zLLR/QpmTNFuNKw4nN3jvZ/Uee2TlncIj63cHR276+Fjic
byRI7QJAUBHKTFRHjEfOknuocc3KWuMYd2U9QJFF5ZTk7jSqfL2NC7yMflobh+ro
KM+/3hTEMYNuM0E8hr2YaIjiVGh1MwJAFzzp2OgrTyw5UCeikhyjzUIQ91eQk3q/
168bkR80x7LF6m4gtWf4EYopcEqdWZEd4IWTk71yid0wE1ZLdyE72QJBAJ2TUwsg
S/Xr93AMLjj2rN4vBo2WFehBOHBOsKbslc/dATZLZnCi8AfZCTEt82cZJIHsMWBE
wubBeV+0BWViVsg=
-----END RSA PRIVATE KEY-----`

var pub = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDGnqhG6CVltWvOOxhnZfq18Qs6
xaPy31+oZ5ltCgqNIO1KxslpHQEIMUj1YVg82bTXSWX/RUr7VrdzJJeNOYFKNoAo
HXpuDf5RH8zDMEue7z6YLr/OUOEhKvKYQWGx9RdKkrCfrqUCy0RZlFvuriCLhuaQ
AZVsEhJYxwWj+4gjLwIDAQAB
-----END PUBLIC KEY-----`
*/
func main() {
	//s := `eyJhY2NvdW50X3R5cGUiOjAsImFjY291bnQiOiIyMzMiLCJwYXNzd29yZCI6ImZhNWVlMmVjNzNkNjY0MWVjZjMxMjg2MDI5ZWY2YjE4IiwiaW1laSI6IjExMTEiLCJub25zdHIiOiJzSHd5TFlvUGl2SDlubVBvIiwibGFuZyI6InpoX0NOIn0=`
	//decodeString, _ := base64.StdEncoding.DecodeString(s)
	//fmt.Println(string(decodeString))

	m := `bgDKtLe3wYa5s9F90t7yk6ZBzfjtRustqyOXhoIFm5QUbUVY9cyTR4GFKLWnql/68GFqbMDvk/SIsAxFi7Lwth5VhIba/N+Dmtu+fbTv0mPYY9YFaivkZEG0GtcXfHtQ2CbAXufvzg5rAuV3jyJ0OmKI2HkmVngZfKRo/glsDbB0lUUJb8LuVi6WMhawuVxxjxhGIaNw/aQ0r9+O7JaDz1ER9MMPg3vHgu2O4EY5l6sy78COYGmADni//BgtCjFRNLK9A50y9ASvL09vnZ8LcCYUuVigSAdA7HO8yp2t0pbzDzVKKBtUgli0E5IBJzTVG2xX9ypU83k+Evv4kpeEBYuZlPHZ3fmU5deco2bg265OYiIW4LI4ap1aVYjA6xtX+gSDwqK+/JYTSgKdexy0VwOASz84zYBbF7bkeXxO1RCIiVPZqDgWkFtLc6T3/hcOTLym5f+M/8nVWq1eCWVmGYGznH1aSc74J+RZEN5QUVRVqcUbpkXUYQxStDhz49yJE54+D5NnPMmidrcYLSaROL8GYMlkivRAG+jPh17P3Ih8Uk7fvGyDAEUVvdGDkKm24K3ERMfWJPdxzKogmcyPx055mQ7fLc67OvARFr/6SmQgFe/TLcEd4pnjYIDgOA4fXGQqoHt5FXlGJ35IIkYIZeKV4oarUv1uHctP29hxKUk=`
	//c7219b9afbda03a00b2bb767f2e7ab29
	// 解密
	dencr, derr := encrypt.DoRsa(encrypt.HANDLE_DECRYPT, encrypt.KEYTYPE_PKCS8, encrypt.RETFMT_STRING, encrypt.HASHLENTYPE_NONE,
		m, pri, encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_BASE64)
	if derr != nil {
		fmt.Println("============解密失败")
		return
	}

	//fmt.Println(string(dencr.(string)))
	decodeString1, _ := base64.StdEncoding.DecodeString(dencr.(string))
	fmt.Println("明文", string(decodeString1))

	// 对解密的明文MD5
	m2 := encrypt.DoMd5(strext.ToStringNoPoint(dencr))
	sign := `Ko2QGEMgVsizUM8hGGeeYh4zF6EiNpbfO0rLos+Vm1RkKNP92qCACcM5u4EhZGyOFGK4+oqkDlAltRuAoKiwBlJvdZQUVxRoten8vgDFxw3MgGNLe3WjWK4KxZ61SeLztSZBzq/z57ZKi6HQzGwe3PC+c5B4rjEo84rPdXrfpKo=`
	//fmt.Println(m2)

	// 验签
	verifySign, err := encrypt.DoRsa(encrypt.HANDLE_VERIFY, encrypt.KEYTYPE_PKCS8, encrypt.RETFMT_NONE, encrypt.HASHLENTYPE_SHA256,
		map[string]interface{}{
			"sign": sign,
			"data": []byte(m2),
		}, pub, encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_SIGN_MAP_BASE64)
	if err != nil {
		fmt.Println("===========验证签名失败", err.Error())
		return
	}
	fmt.Println("============verifySign", verifySign)

}
