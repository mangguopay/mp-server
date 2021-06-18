package myrsa

import (
	"encoding/base64"
	"testing"
	"time"
)

// 可通过openssl产生
//openssl genrsa -out rsa_private_key.pem 1024
var pkcs1PrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAtJmAhSZd/t0Kpu8+Bggzu9D87MCUOU6v/jkwhwZrNhL/s5N1
uv2y/tYIAhYDMY6ffPvycxa1BubRRA7X4T7YiN5NKdlZ63FBN/hVQQIpiWv/nM1E
PqSX0+H5TJNpyE5gPxUsEYAMIxVgSOOKCriOhlA5kMQ43EsIzAC0izemZTtrfaMR
Zfiiya4ASozUc2ULgDzKEx9Mk3EyEU1kZzx7jFeaaag6sp14pWt4u6yK5KtZ3aZN
b4KRCDGMR0GW1QobzAbTF1S6tSApKKE8QKwoCSahyiilyuX3o/bIu3l923h7e39T
+ICrer4Z81gFE+qw2YQJHaHRB/y0qIrqfTtquQIDAQABAoIBABQ6DPbcsTJWN1cy
/FNYn9KtOWaauv8gHP3wEctNoLxRrGnEXi8vMtjvx73UEU9/lcU9wo25QVUgjNd+
ghcsHrxViXbibPu9l3xZR9McFuLZQauiIb6SWJ7WdEFrUTGp9YWbLXBYUwOq5TnE
ojVJLk1Etu3LzEJ/9OBEQ/RDq1MzABGO1sjnohN75Vl4QxlpMuHKg65gguoIa+DK
bWClxOiV676dJWCfjBHnaDgYgS565Qct0u9DStdor2W3Nhmx04qgW5kHn/Qq/HTi
t+4aVPavoEjVwxTCammsLxwLtvSzzg47x/Vj1A/eJ483zy/PNzRPuAx0NMPQb8Dd
AiY9L/ECgYEA7fOTxvg+HzVwgOIDWfvGEAnLrUjWL0WrwlK0E9m/1Q2T+dskgwnC
ry4iR+DriLfDEDLbn/2go6B0jzReW2guxo8VxZEo6KsY3f9rUQaAy+bJoQW4LtCv
KOWbUMM4hhj3RYe5ACXrt7lzAQKzOYC1xxbJ60pGU35fIoQ0pyefEo8CgYEAwkxL
jAHNqY5zmSuD4M8VRrGLkWOZkQIKbVSheFhdESNJZ8H2ZQIbxjTlyTdc74Q82QK5
Bo207pEMsK2m+G3V/Q1fU8ncRAaT5b304QOuiAtS2s3+fQci4vgpxqUlqKSxYoYy
F4N8xDFht8sjPxorCtdFwMRIDlkGN9V8zVXpsjcCgYEAuL7sFoh4mvx/u+E+3udn
EN66H3E0soEyaO6TV/IxSbaAFHa7s22plR+JiCsuU/jw3yvNbzuZNFGJDgKH3ApY
fttq+PjKPVNSPFJqPP+Ckk0+cOGi7d4ikOssGpln0l2h5n8I+P94My4uBzPUeSng
eJHN9fu1/G9aZ88jnkBZ9isCgYAy+CC5UZ/J4vygKbImvywtp1Wdhis6xvZFR/Yz
w7pmTINtHIyuYqc2j5nX9xYCHwZ3RyeSeIoGKzbRAjzS3r1L7L4dFM8baT5S/knG
3Vhjh9TsYS1pTv3v3HnZCmmem9WMqvdpA60vKmUf+cH9Q7gW1/IMZ3Efkmr3KqHa
m7b6cQKBgD+J3TPRVieONc0p7Q0dtkeVdrGmWEOoOjY6I7DkBLcbtzA0XMwGRCTY
j0ghNWf9wfsdtrm6tvEHkacqms6tfamgOq7B9SHM+o7JWo0aGG6epfhPH3ng8wAf
D61eH41IeIw21BFG503EuGqEsHogcdBEbJD6ljWepv6QwBmdSSHz
-----END RSA PRIVATE KEY-----
`

//openssl
//openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem
var pkcs1PublicKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtJmAhSZd/t0Kpu8+Bggz
u9D87MCUOU6v/jkwhwZrNhL/s5N1uv2y/tYIAhYDMY6ffPvycxa1BubRRA7X4T7Y
iN5NKdlZ63FBN/hVQQIpiWv/nM1EPqSX0+H5TJNpyE5gPxUsEYAMIxVgSOOKCriO
hlA5kMQ43EsIzAC0izemZTtrfaMRZfiiya4ASozUc2ULgDzKEx9Mk3EyEU1kZzx7
jFeaaag6sp14pWt4u6yK5KtZ3aZNb4KRCDGMR0GW1QobzAbTF1S6tSApKKE8QKwo
CSahyiilyuX3o/bIu3l923h7e39T+ICrer4Z81gFE+qw2YQJHaHRB/y0qIrqfTtq
uQIDAQAB
-----END PUBLIC KEY-----
`

var pkcs8PrivateKey = `
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC0mYCFJl3+3Qqm
7z4GCDO70PzswJQ5Tq/+OTCHBms2Ev+zk3W6/bL+1ggCFgMxjp98+/JzFrUG5tFE
DtfhPtiI3k0p2VnrcUE3+FVBAimJa/+czUQ+pJfT4flMk2nITmA/FSwRgAwjFWBI
44oKuI6GUDmQxDjcSwjMALSLN6ZlO2t9oxFl+KLJrgBKjNRzZQuAPMoTH0yTcTIR
TWRnPHuMV5ppqDqynXila3i7rIrkq1ndpk1vgpEIMYxHQZbVChvMBtMXVLq1ICko
oTxArCgJJqHKKKXK5fej9si7eX3beHt7f1P4gKt6vhnzWAUT6rDZhAkdodEH/LSo
iup9O2q5AgMBAAECggEAFDoM9tyxMlY3VzL8U1if0q05Zpq6/yAc/fARy02gvFGs
acReLy8y2O/HvdQRT3+VxT3CjblBVSCM136CFywevFWJduJs+72XfFlH0xwW4tlB
q6IhvpJYntZ0QWtRMan1hZstcFhTA6rlOcSiNUkuTUS27cvMQn/04ERD9EOrUzMA
EY7WyOeiE3vlWXhDGWky4cqDrmCC6ghr4MptYKXE6JXrvp0lYJ+MEedoOBiBLnrl
By3S70NK12ivZbc2GbHTiqBbmQef9Cr8dOK37hpU9q+gSNXDFMJqaawvHAu29LPO
DjvH9WPUD94njzfPL883NE+4DHQ0w9BvwN0CJj0v8QKBgQDt85PG+D4fNXCA4gNZ
+8YQCcutSNYvRavCUrQT2b/VDZP52ySDCcKvLiJH4OuIt8MQMtuf/aCjoHSPNF5b
aC7GjxXFkSjoqxjd/2tRBoDL5smhBbgu0K8o5ZtQwziGGPdFh7kAJeu3uXMBArM5
gLXHFsnrSkZTfl8ihDSnJ58SjwKBgQDCTEuMAc2pjnOZK4PgzxVGsYuRY5mRAgpt
VKF4WF0RI0lnwfZlAhvGNOXJN1zvhDzZArkGjbTukQywrab4bdX9DV9TydxEBpPl
vfThA66IC1Lazf59ByLi+CnGpSWopLFihjIXg3zEMWG3yyM/GisK10XAxEgOWQY3
1XzNVemyNwKBgQC4vuwWiHia/H+74T7e52cQ3rofcTSygTJo7pNX8jFJtoAUdruz
bamVH4mIKy5T+PDfK81vO5k0UYkOAofcClh+22r4+Mo9U1I8Umo8/4KSTT5w4aLt
3iKQ6ywamWfSXaHmfwj4/3gzLi4HM9R5KeB4kc31+7X8b1pnzyOeQFn2KwKBgDL4
ILlRn8ni/KApsia/LC2nVZ2GKzrG9kVH9jPDumZMg20cjK5ipzaPmdf3FgIfBndH
J5J4igYrNtECPNLevUvsvh0UzxtpPlL+ScbdWGOH1OxhLWlO/e/cedkKaZ6b1Yyq
92kDrS8qZR/5wf1DuBbX8gxncR+SavcqodqbtvpxAoGAP4ndM9FWJ441zSntDR22
R5V2saZYQ6g6NjojsOQEtxu3MDRczAZEJNiPSCE1Z/3B+x22ubq28QeRpyqazq19
qaA6rsH1Icz6jslajRoYbp6l+E8feeDzAB8PrV4fjUh4jDbUEUbnTcS4aoSweiBx
0ERskPqWNZ6m/pDAGZ1JIfM=
-----END PRIVATE KEY-----
`

var pkcs8PublicKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtJmAhSZd/t0Kpu8+Bggz
u9D87MCUOU6v/jkwhwZrNhL/s5N1uv2y/tYIAhYDMY6ffPvycxa1BubRRA7X4T7Y
iN5NKdlZ63FBN/hVQQIpiWv/nM1EPqSX0+H5TJNpyE5gPxUsEYAMIxVgSOOKCriO
hlA5kMQ43EsIzAC0izemZTtrfaMRZfiiya4ASozUc2ULgDzKEx9Mk3EyEU1kZzx7
jFeaaag6sp14pWt4u6yK5KtZ3aZNb4KRCDGMR0GW1QobzAbTF1S6tSApKKE8QKwo
CSahyiilyuX3o/bIu3l923h7e39T+ICrer4Z81gFE+qw2YQJHaHRB/y0qIrqfTtq
uQIDAQAB
-----END PUBLIC KEY-----
`

func TestRSA2Sign(t *testing.T) {
	data := "app_id=20201234234234&charset=utf-8&data=dddddddddddddddddddddddd&sign_type=RSA2&timestamp=" + time.Now().Format("2006-01-02 15:04:05")

	t.Logf("data:%s", data)

	// rsa2签名
	sign, err := RSA2Sign(data, pkcs8PrivateKey)
	if err != nil {
		t.Errorf("RSA2Sign-err:%v", err)
		return
	}
	t.Logf("sign:%s", sign)
}

func TestRSA2Verify(t *testing.T) {
	data := "app_id=20201234234234&charset=utf-8&data=dddddddddddddddddddddddd&sign_type=RSA2&timestamp=2020-07-02 17:37:16"
	sign := "eAiXN8h/HfkZOqpZ/z5vvflkVI3TmueOHJBK1BBodmfzKiUZrM2VHChV6t6fsU5tiGREWXO9qzH3FrDWYrsEDkOT0bu4AIuuaj7PADT4fRNONJ7wy1Jmk78X7KkttseDNZ4SpW4itLVu+Yl7rgxfgJTvvkhph17TdMAOwCZDjtwgZkZOmyJMPlJpDk1Z1RBS5iloLaoHCKHsS7ZPg/D7rouEut2sDWOiBpVGHXAybclwoZIDcUj8t5EzEmBNOgrLuVBbGOOQ2+EUoH9o+g0tcpqkRrTDocUSzcmZQ8dSQTjTmOcW7tmjEaJif0k4iEOhCxXNOgBPa8OqRGWG7YDXSQ=="

	err := RSA2Verify(data, sign, pkcs1PublicKey)
	if err != nil {
		t.Errorf("RSA2Verify-err:%v", err)
		return
	}

	t.Logf("RSA2Verify-ok")
}

func TestRSAEncrypt(t *testing.T) {
	data := "abcdefg"

	cipher, err := RSAEncrypt(data, pkcs1PublicKey)
	if err != nil {
		t.Errorf("RSAEncrypt-err:%v", err)
		return
	}

	t.Logf("RSAEncrypt-cipher:%s", base64.StdEncoding.EncodeToString(cipher))
}

func TestRSADecrypt(t *testing.T) {
	cipher := "VqPQuUJPcr0BuUX7i+79Wt1+WuVCeSc/HHWDrTGAs94mPW38eGveFNycbtzWFNVuJgbCKj8D6uHlzTBX1usyx+638UtUtHUrIX+95ZzMHTrOQCfSoTc6LE1D/io5jvJvTFX0SJ4DvC1q827gBESLKgmiMCNMBij1diBj5b9kehUrVO1OlzdsWPMPISOq2RXTz1LlKGvQaL/7Q26C0ZzGEBARt/hjLUVh/VG4vPN5AFzwFjFlEmcgc8tlefMhiRjRy6h3HNWVgu/qjqIcu2P20DTL35Mxl+vo5GfuxW5lJY+D3oIj0h/dLdXDkJZgzJMQRT5xGp8Y3WG8QU9dZR1mIg=="

	//textplain, err := RSADecrypt(cipher, pkcs1PrivateKey)
	textplain, err := RSADecrypt(cipher, pkcs8PrivateKey)
	if err != nil {
		t.Errorf("RSADecrypt-err:%v", err)
		return
	}

	t.Logf("RSADecrypt-textplain:%s", string(textplain))
}

func TestRSADecryptWithPKCS1(t *testing.T) {
	cipher := "L1jjcy9ky6mPg/Fgnf+f0WYEamUk4jnlNT4osVvOuN4Kk1y/2Snre6PAc2UlaOmkXgY4Hw1GWm/ClFTj+dE6MckGi8jUb1V8IMSIgFhE2XHFgMT42f8mZ5ig/BCDEoVch0Tz48lO2qVzTt76ALbnMc67OP1F5b/7ralKHL42gmMABm+DTMCHkCT4hFz1VLlVFEcvTHgWQz+lYN3pz8UBLD08pvE8aJ0FybGfWyLV9VdvPUYAjHKnIes6O1p+LwGBfuzW9fkqU756HfkEqYMA0fdj+v3rJssQD3RqT2TQzZJz3s1sQvIhlrWeO91DHiDQAg4T6f79Ti1w+UOwJ0uy/Q=="

	cipherBytes, err := base64.StdEncoding.DecodeString(cipher)
	if err != nil {
		t.Errorf("RSADecryptWithPKCS1-base64-Decode-err:%v", err)
		return
	}

	textplain, err := RSADecryptWithPKCS1(cipherBytes, []byte(pkcs1PrivateKey))
	if err != nil {
		t.Errorf("RSADecryptWithPKCS1-err:%v", err)
		return
	}

	t.Logf("RSADecryptWithPKCS1-textplain:%s", string(textplain))
}

func TestRSADecryptWithPKCS8(t *testing.T) {
	cipher := "L1jjcy9ky6mPg/Fgnf+f0WYEamUk4jnlNT4osVvOuN4Kk1y/2Snre6PAc2UlaOmkXgY4Hw1GWm/ClFTj+dE6MckGi8jUb1V8IMSIgFhE2XHFgMT42f8mZ5ig/BCDEoVch0Tz48lO2qVzTt76ALbnMc67OP1F5b/7ralKHL42gmMABm+DTMCHkCT4hFz1VLlVFEcvTHgWQz+lYN3pz8UBLD08pvE8aJ0FybGfWyLV9VdvPUYAjHKnIes6O1p+LwGBfuzW9fkqU756HfkEqYMA0fdj+v3rJssQD3RqT2TQzZJz3s1sQvIhlrWeO91DHiDQAg4T6f79Ti1w+UOwJ0uy/Q=="

	cipherBytes, err := base64.StdEncoding.DecodeString(cipher)
	if err != nil {
		t.Errorf("RSADecryptWithPKCS8-base64-Decode-err:%v", err)
		return
	}

	textplain, err := RSADecryptWithPKCS8(cipherBytes, []byte(pkcs8PrivateKey))
	if err != nil {
		t.Errorf("RSADecryptWithPKCS8-err:%v", err)
		return
	}

	t.Logf("RSADecryptWithPKCS8-textplain:%s", string(textplain))
}
