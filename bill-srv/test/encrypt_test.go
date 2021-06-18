package test

import (
	"testing"

	"a.a/cu/encrypt"
)

func TestGenPriPub(t *testing.T) {
	// 生成公钥私钥
	encrypt.GenRsaKeyFile(1024, "private.pem", "public.pem")
}
