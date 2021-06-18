package ss_rsa

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"strings"
)

const (
	kPublicKeyPrefix = "-----BEGIN PUBLIC KEY-----"
	kPublicKeySuffix = "-----END PUBLIC KEY-----"

	kPKCS1Prefix = "-----BEGIN RSA PRIVATE KEY-----"
	KPKCS1Suffix = "-----END RSA PRIVATE KEY-----"

	kPKCS8Prefix = "-----BEGIN PRIVATE KEY-----"
	KPKCS8Suffix = "-----END PRIVATE KEY-----"
)

var (
	ErrPrivateKeyFailedToLoad = errors.New("private key failed to load")
	ErrPublicKeyFailedToLoad  = errors.New("public key failed to load")
)

func formatKey(raw, prefix, suffix string, lineCount int) []byte {
	if raw == "" {
		return nil
	}
	raw = strings.Replace(raw, prefix, "", 1)
	raw = strings.Replace(raw, suffix, "", 1)
	raw = strings.Replace(raw, " ", "", -1)
	raw = strings.Replace(raw, "\n", "", -1)
	raw = strings.Replace(raw, "\r", "", -1)
	raw = strings.Replace(raw, "\t", "", -1)

	var sl = len(raw)
	var c = sl / lineCount
	if sl%lineCount > 0 {
		c = c + 1
	}

	var buf bytes.Buffer
	buf.WriteString(prefix + "\n")
	for i := 0; i < c; i++ {
		var b = i * lineCount
		var e = b + lineCount
		if e > sl {
			buf.WriteString(raw[b:])
		} else {
			buf.WriteString(raw[b:e])
		}
		buf.WriteString("\n")
	}
	buf.WriteString(suffix)
	return buf.Bytes()
}

func FormatPublicKey(raw string) []byte {
	return formatKey(raw, kPublicKeyPrefix, kPublicKeySuffix, 64)
}

func FormatPKCS1PrivateKey(raw string) []byte {
	raw = strings.Replace(raw, kPKCS8Prefix, "", 1)
	raw = strings.Replace(raw, KPKCS8Suffix, "", 1)
	return formatKey(raw, kPKCS1Prefix, KPKCS1Suffix, 64)
}

func FormatPKCS8PrivateKey(raw string) []byte {
	raw = strings.Replace(raw, kPKCS1Prefix, "", 1)
	raw = strings.Replace(raw, KPKCS1Suffix, "", 1)
	return formatKey(raw, kPKCS8Prefix, KPKCS8Suffix, 64)
}

func ParsePKCS1PrivateKey(data []byte) (key *rsa.PrivateKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, ErrPrivateKeyFailedToLoad
	}

	key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, err
}

func ParsePKCS8PrivateKey(data []byte) (key *rsa.PrivateKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, ErrPrivateKeyFailedToLoad
	}

	rawKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := rawKey.(*rsa.PrivateKey)
	if ok == false {
		return nil, ErrPrivateKeyFailedToLoad
	}

	return key, err
}

func ParsePublicKey(data []byte) (key *rsa.PublicKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, ErrPublicKeyFailedToLoad
	}

	var pubInterface interface{}
	pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, ErrPublicKeyFailedToLoad
	}

	return key, err
}

func packageData(originalData []byte, packageSize int) (r [][]byte) {
	var src = make([]byte, len(originalData))
	copy(src, originalData)

	r = make([][]byte, 0)
	if len(src) <= packageSize {
		return append(r, src)
	}
	for len(src) > 0 {
		var p = src[:packageSize]
		r = append(r, p)
		src = src[packageSize:]
		if len(src) <= packageSize {
			r = append(r, src)
			break
		}
	}
	return r
}

// RSA 公钥加密
//
// data 待加密的数据
// pubKey 公钥
func RSAEncrypt(data, pubKey string) ([]byte, error) {
	pub, err := ParsePublicKey([]byte(pubKey))
	if err != nil {
		return nil, err
	}

	return RSAEncryptWithKey([]byte(data), pub)
}

// RSA 私钥解密
//
// cipher 密文 base64编码格式
// priKey 私钥 支持PKCS1和PKCS8格式
func RSADecrypt(cipher, priKey string) ([]byte, error) {
	cBytes, err := base64.StdEncoding.DecodeString(cipher)
	if err != nil {
		return nil, err
	}

	pri, err := ParsePKCS1PrivateKey(FormatPKCS1PrivateKey(priKey))
	if err != nil {
		pri, err = ParsePKCS8PrivateKey(FormatPKCS8PrivateKey(priKey))
		if err != nil {
			return nil, err
		}
	}

	return RSADecryptWithKey(cBytes, pri)
}

// RSAEncryptWithKey 使用公钥 key 对数据 src 进行 RSA 加密
func RSAEncryptWithKey(src []byte, key *rsa.PublicKey) ([]byte, error) {
	var data = packageData(src, key.N.BitLen()/8-11)
	var cipher = make([]byte, 0, 0)

	for _, d := range data {
		var c, e = rsa.EncryptPKCS1v15(rand.Reader, key, d)
		if e != nil {
			return nil, e
		}
		cipher = append(cipher, c...)
	}

	return cipher, nil
}

// RSADecryptWithPKCS1 使用私钥 key 对数据 cipher 进行 RSA 解密，key 的格式为 pkcs1
func RSADecryptWithPKCS1(cipher, key []byte) ([]byte, error) {
	pri, err := ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, err
	}

	return RSADecryptWithKey(cipher, pri)
}

// RSADecryptWithPKCS1 使用私钥 key 对数据 cipher 进行 RSA 解密，key 的格式为 pkcs8
func RSADecryptWithPKCS8(cipher, key []byte) ([]byte, error) {
	pri, err := ParsePKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}

	return RSADecryptWithKey(cipher, pri)
}

// RSADecryptWithKey 使用私钥 key 对数据 cipher 进行 RSA 解密
func RSADecryptWithKey(cipher []byte, key *rsa.PrivateKey) ([]byte, error) {
	var data = packageData(cipher, key.PublicKey.N.BitLen()/8)
	var plainData = make([]byte, 0, 0)

	for _, d := range data {
		var p, e = rsa.DecryptPKCS1v15(rand.Reader, key, d)
		if e != nil {
			return nil, e
		}
		plainData = append(plainData, p...)
	}
	return plainData, nil
}

func RSASignWithPKCS1(src, key []byte, hash crypto.Hash) ([]byte, error) {
	pri, err := ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, err
	}
	return RSASignWithKey(src, pri, hash)
}

func RSASignWithPKCS8(src, key []byte, hash crypto.Hash) ([]byte, error) {
	pri, err := ParsePKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	return RSASignWithKey(src, pri, hash)
}

func RSASignWithKey(src []byte, key *rsa.PrivateKey, hash crypto.Hash) ([]byte, error) {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, key, hash, hashed)
}

func RSAVerify(src, sig, key []byte, hash crypto.Hash) error {
	pub, err := ParsePublicKey(key)
	if err != nil {
		return err
	}
	return RSAVerifyWithKey(src, sig, pub, hash)
}

func RSAVerifyWithKey(src, sig []byte, key *rsa.PublicKey, hash crypto.Hash) error {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)
	return rsa.VerifyPKCS1v15(key, hash, hashed, sig)
}

// data 待签名的数据
// privKey 私钥 支持PKCS1和PKCS8
//
// 返回base64编码的sign
func RSA2Sign(data string, privKey string) (string, error) {
	priKey, err := ParsePKCS1PrivateKey(FormatPKCS1PrivateKey(privKey))
	if err != nil {
		priKey, err = ParsePKCS8PrivateKey(FormatPKCS8PrivateKey(privKey))
		if err != nil {
			return "", err
		}
	}

	sign, err := RSASignWithKey([]byte(data), priKey, crypto.SHA256)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sign), nil
}

// data 待验签的数据
// sign base64格式的字符串
// pubKey 公钥
func RSA2Verify(data, sign, pubKey string) error {
	pub, err := ParsePublicKey(FormatPublicKey(pubKey))
	if err != nil {
		return err
	}

	byteSgin, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}

	return RSAVerifyWithKey([]byte(data), byteSgin, pub, crypto.SHA256)
}

// 生成RSA密钥对-私钥为PKCS1格式
// privKey, pubKey, err := GenRsaKeyPairPKCS1(2048)
func GenRsaKeyPairPKCS1(bits int) (string, string, error) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}

	privateBuf := new(bytes.Buffer)
	if err := pem.Encode(privateBuf, block); err != nil {
		return "", "", err
	}

	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}

	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}

	publicBuf := new(bytes.Buffer)
	if err := pem.Encode(publicBuf, block); err != nil {
		return "", "", err
	}

	return privateBuf.String(), publicBuf.String(), nil
}

// 生成RSA密钥对-私钥为PKCS8格式
// privKey, pubKey, err := GenRsaKeyPairPKCS8(2048)
func GenRsaKeyPairPKCS8(bits int) (string, string, error) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	derStream, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}

	privateBuf := new(bytes.Buffer)
	if err := pem.Encode(privateBuf, block); err != nil {
		return "", "", err
	}

	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}

	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}

	publicBuf := new(bytes.Buffer)
	if err := pem.Encode(publicBuf, block); err != nil {
		return "", "", err
	}

	return privateBuf.String(), publicBuf.String(), nil
}

// 判断私钥是否是PKCS1格式
func IsPKCS1PrivateKey(privKey string) bool {
	// 无论什么格式，都先转换为PKCS1格式
	_, err := ParsePKCS1PrivateKey(FormatPKCS1PrivateKey(privKey))
	if err != nil {
		return false
	}

	return true
}

// 判断私钥是否是PKCS8格式
func IsPKCS8PrivateKey(privKey string) bool {
	// 无论什么格式，都先转换为PKCS8格式
	_, err := ParsePKCS8PrivateKey(FormatPKCS8PrivateKey(privKey))
	if err != nil {
		return false
	}

	return true
}

// 去掉RSA密钥的头尾、换行和空格，转成一行字符串
func StripRSAKey(key string) string {
	key = strings.Replace(key, kPKCS1Prefix, "", 1)
	key = strings.Replace(key, KPKCS1Suffix, "", 1)

	key = strings.Replace(key, kPKCS8Prefix, "", 1)
	key = strings.Replace(key, KPKCS8Suffix, "", 1)

	key = strings.Replace(key, kPublicKeyPrefix, "", 1)
	key = strings.Replace(key, kPublicKeySuffix, "", 1)

	key = strings.Replace(key, " ", "", -1)
	key = strings.Replace(key, "\n", "", -1)
	key = strings.Replace(key, "\r", "", -1)
	key = strings.Replace(key, "\t", "", -1)

	return key
}
