package test

import (
	"a.a/mp-server/common/cache"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestPwd(t *testing.T) {
	//pwd := "fADf24ddaaff2323FG3Sd545fgdd"
	pwd := "2$G!f"
	fmt.Println(len(pwd))
	// 判断密码长度
	if len(pwd) < 3 || len(pwd) > 5 {
		fmt.Println("长度应该为3-5位")
	}
	// 是否包含小写
	reg1 := regexp.MustCompile(`[a-z]+`)
	fmt.Printf("是否包含小写%v\n", reg1.MatchString(pwd))
	reg2 := regexp.MustCompile(`[A-Z]+`)
	fmt.Printf("是否包含大写%v\n", reg2.MatchString(pwd))

	// 是否包含数字
	reg3 := regexp.MustCompile(`[0-9]+`)
	fmt.Printf("是否包含数字%v\n", reg3.MatchString(pwd))

	// 是否包含特殊字符
	f := "!@#$%^&*"
	fmt.Println("是否包含特殊字符", strings.ContainsAny(pwd, f))
}

func CheckPWDRule() {
	pwd := "2$G!f"
	fmt.Println(len(pwd))
	// 判断密码长度
	if len(pwd) < 3 || len(pwd) > 5 {
		fmt.Println("长度应该为3-5位")
	}
	// 是否包含小写
	reg1 := regexp.MustCompile(`[a-z]+`)
	fmt.Printf("是否包含小写%v\n", reg1.MatchString(pwd))
	reg2 := regexp.MustCompile(`[A-Z]+`)
	fmt.Printf("是否包含大写%v\n", reg2.MatchString(pwd))

	// 是否包含数字
	reg3 := regexp.MustCompile(`[0-9]+`)
	fmt.Printf("是否包含数字%v\n", reg3.MatchString(pwd))

	// 是否包含特殊字符
	f := "!@#$%^&*"
	fmt.Println("是否包含特殊字符", strings.ContainsAny(pwd, f))
}

func TestJoin(t *testing.T) {
	a := "https://modernpay-test.s3-ap-southeast-1.amazonaws.com"
	b := "img/452a9125af452d477efb606b58589012.jpeg"
	fmt.Println(fmt.Sprintf("%s/%s", a, b))

}

func TestCache(t *testing.T) {
	cache.RedisClient.Set("abc", "123", 0)
}
