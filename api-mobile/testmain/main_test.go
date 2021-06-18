package testmain

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
//	dataByte := base64.StdEncoding.EncodeToString([]byte(`{"account_type":0,"account":"233","password":"ba6ac31b395d3cff156d044ba5ab6931","imei":"1111","nonstr":"RegjXppdI7ahHGO7","lang":"zh_CN"}`))
//	//dataByte:= base64.StdEncoding.EncodeToString([]byte(`{"account_type":0,"account":"233","password":"ba6ac31b395d3cff156d044ba5ab6931"}`))
//	// 加密
//	encr, err := encrypt.DoRsa(encrypt.HANDLE_ENCRYPT, encrypt.KEYTYPE_PKIX, encrypt.RETFMT_BASE64, encrypt.HASHLENTYPE_SHA256,
//		[]byte(dataByte), pub, encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_BYTE)
//	if err != nil {
//		fmt.Println("============加密失败", err.Error())
//		return
//	}
//	fmt.Println("encr----->", encr)
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
//	// 验证签名
//	verifySign, err := encrypt.DoRsa(encrypt.HANDLE_VERIFY, encrypt.KEYTYPE_PKCS8, encrypt.RETFMT_NONE, encrypt.HASHLENTYPE_SHA256,
//		map[string]interface{}{
//			"sign": sign,
//			"data": []byte(dataByte),
//		}, pub, encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_SIGN_MAP)
//	if err != nil {
//		fmt.Println("===========验证签名失败", err.Error())
//		return
//	}
//	fmt.Println("============verifySign", verifySign)
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
//}
