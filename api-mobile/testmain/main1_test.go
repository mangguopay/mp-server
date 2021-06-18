package testmain

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"
)

//
//import (
//	"a.a/cu/encrypt"
//	"encoding/base64"
//	"encoding/hex"
//	"fmt"
//)
//
//func main() {
//	//data := `{"account_type":0,"account":"233","password":"ba6ac31b395d3cff156d044ba5ab6931","imei":"1111","nonstr":"RegjXppdI7ahHGO7","lang":"zh_CN"}`
//	//dataByte := base64.StdEncoding.EncodeToString([]byte(data))
//
//	pub := `MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDJ2ZFm5giqgWooFfH/97m6qbq8Iidche9XhI4BpzsO8uLo1ppcGHwqopHOLeg3lMcK55gNfOStBNr8fcbW/mdfJFGoV4k2fSWazGr1wEu/JJO4vWYGkUsrGwUqqos2kUsUBMlgOGryBG5hSwHI0tfSKE45r2dQEkLONV9uM2klywIDAQAB`
//	pri := `MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAMnZkWbmCKqBaigV8f/3ubqpurwiJ1yF71eEjgGnOw7y4ujWmlwYfCqikc4t6DeUxwrnmA185K0E2vx9xtb+Z18kUahXiTZ9JZrMavXAS78kk7i9ZgaRSysbBSqqizaRSxQEyWA4avIEbmFLAcjS19IoTjmvZ1ASQs41X24zaSXLAgMBAAECgYAdlTJ3Nz2L572sNFMsZZ0l21xP2F2WdNa4J+g8G8tPDI/u+gzTdI82HG9sHVRMWRS252xkhapqJS0HFlP/DHgPubJDdykSy+RONnQhUy0U8/s0n1BpdlUtprXZjZSyaVAGll/BB719dIw9v/EFmqLfUZ3P7RXHp9NakKndrfBuuQJBAPl9z/j9IHJRXINrWd+7MMqhAPbdWfv/Etu/qA2UW3uScQ4rHDQgpw2+uBNHT6Cgp1vLxeOBKw0QYS6RKkBDxV8CQQDPHZbRnpxmYPDJB4EHmaQm0ueu8fw6gDOK5hYST1EMwRaSjfhWGbMYVIAkEJLzrOTLrgq1EL4f2XmiBSJmLSsVAkAJQiH1m28Yzuwf9FvhcZDd9BuVDaHJOC36+aHC3z6F12lanT7usEeCuxEZpgvOaifLwEQXrTNryK/SipCG0f0BAkBMuBUpyKr+cWI/1PvCqPLZPr57Wz+nG9370YbTeXX4V33ZA6W9nv9sP8DHmywT/zMxD2L/9xe2DIS1s1kuqASRAkEA1gSaQQdYuZHhhGEA3j/SErogSQVnXoijCsz2eVn58LMr4VqnDpUFXRPBpV37/u4UPH8jn70kH7JP+HKlKxN55w==`
//
//	//dataByte := base64.StdEncoding.EncodeToString([]byte(`{"account_type":0,"account":"233","password":"ba6ac31b395d3cff156d044ba5ab6931","imei":"1111","nonstr":"RegjXppdI7ahHGO7","lang":"zh_CN"}`))
//	s := `{
//    "page":"1",
//    "page_size":"10",
//    "money_type":"khr"
//}`
//
//	dataByte := base64.StdEncoding.EncodeToString([]byte(s))
//
//	// 加密
//	encr, err := encrypt.DoRsa(encrypt.HANDLE_ENCRYPT, encrypt.KEYTYPE_PKIX, encrypt.RETFMT_BASE64, encrypt.HASHLENTYPE_SHA256,
//		[]byte(dataByte), pub, encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_BYTE)
//	if err != nil {
//		fmt.Println("============加密失败", err.Error())
//		return
//	}
//	fmt.Println("encr----->", encr)
//
//	encr = `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`
//
//	// 解密
//	dencr, derr := encrypt.DoRsa(encrypt.HANDLE_DECRYPT, encrypt.KEYTYPE_PKCS8, encrypt.RETFMT_STRING, encrypt.HASHLENTYPE_SHA256,
//		encr, pri, encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_BASE64)
//	if derr != nil {
//		fmt.Println("============解密失败", derr.Error())
//		return
//	}
//	dencrB, _ := base64.StdEncoding.DecodeString(dencr.(string))
//	fmt.Println(string(dencrB))
//
//	// 签名
//	sign, err := encrypt.DoRsa(encrypt.HANDLE_SIGN, encrypt.KEYTYPE_PKCS8, encrypt.RETFMT_BYTE, encrypt.HASHLENTYPE_SHA256,
//		[]byte(dataByte), pri, encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_BYTE)
//	if err != nil {
//		fmt.Println("============签名失败", err.Error())
//		return
//	}
//
//	fmt.Println(hex.EncodeToString(sign.([]byte)))
//
//	sign, _ = hex.DecodeString("1cb0a9c0b195817e3134045e3c1d8442fb3afa0e96461ed4c94b4bebe2ffcff136dcc380f982ce0071486e30e8c296818c3963ebf73402f63cbac7e25e1d0912225c75686a9f300ce1f2c3c4be17fa3fbc38556de2662c1e24875b9d4b2c525ed34e33d2574fa9ba4545211b72abefb71e0c0d37e25ac27eb5fe122bf2acefc5")
//
//	// 验证签名
//	verifySign, err := encrypt.DoRsa(encrypt.HANDLE_VERIFY, encrypt.KEYTYPE_PKCS8, encrypt.RETFMT_NONE, encrypt.HASHLENTYPE_SHA256,
//		map[string]interface{}{
//			"sign": sign,
//			"data": []byte(dencr.(string)),
//		}, pub, encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_SIGN_MAP)
//	if err != nil {
//		fmt.Println("===========验证签名失败", err.Error())
//		return
//	}
//	fmt.Println("============verifySign", verifySign)
//
//}

func Test1(t *testing.T) {
	//a := `q9Z468qOKSkkm8NopAjLs+p6AxJDTy9JnrZrFRwpunXs3rfiaN+jUBkiUx3k2TBwpqJgTDKjH6nsYLCbjEmdQQj/xq7A9BLw32295b5ZiPvqhe55iEVQYMvLUX2ywMoDX6bq3du2J8c5gnXlz+No4tC0ZQ1psfG7RH3stKOc/Mw=`
	//e := `afBjfyaaHqpCC8brmJXeM9WenxPnC1QcL9dwaOHTFJkUD9wlruNSrkOPqM5OnOl93fKRJSMcKx7Lyx1ZZVWof0Yv8lj9sodbI26ICo2yNGVrkc/I322/utkNkj19iTLbwnvKMxak0dqPtSzvnbPlkMBK7nYExidA3X/44oMzBg0=`
	//m := encrypt.DoMd5(e)
	//fmt.Println(m)
	//pubKey := `MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDGnqhG6CVltWvOOxhnZfq18Qs6xaPy31+oZ5ltCgqNIO1KxslpHQEIMUj1YVg82bTXSWX/RUr7VrdzJJeNOYFKNoAoHXpuDf5RH8zDMEue7z6YLr/OUOEhKvKYQWGx9RdKkrCfrqUCy0RZlFvuriCLhuaQAZVsEhJYxwWj+4gjLwIDAQAB`
	//_, err := encrypt.DoRsa(encrypt.HANDLE_VERIFY, encrypt.KEYTYPE_PKCS8, encrypt.RETFMT_NONE, encrypt.HASHLENTYPE_SHA256,
	//	map[string]interface{}{
	//		"sign": a,
	//		"data": []byte(m),
	//	}, pubKey, encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_SIGN_MAP_BASE64)
	//fmt.Println(err)
	s := `MIGJAoGBAJ0yE72N51iL8+e4dYXvIJoaJwfbD/oHOAJuSPKOxwnaqSUAZ2d30H55BR49fIMD7+lyQiGWiZwXfbk/0BKlKgJAUOvwrxFxfPT9mAHmeXOR4umu5JicMzRO6hptrZsCurTi5wexWWv55P8nOofRp0PJogJ+ud0Yx+6l0PCdIQanAgMBAAE=`
	block, _ := pem.Decode([]byte(s))
	if block == nil {
		fmt.Printf("err=%v\n", "public key error!")
		return
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	fmt.Println(pubInterface, err)
}
