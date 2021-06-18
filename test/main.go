package main

import (
	"a.a/cu/strext"
	"a.a/mp-server/common/ss_rsa"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"
)

func main1() {
	//Pay()
	//Query()
	//RSA2Sgin()
	//RSA2Verify()
}

func Pay() {
	// 商家的私钥
	var privKey = "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCg7xVS54c4O0I0pakkCH2hA3qTei+J7Xxc29YEqlMxutcx2fJVOCwbscYhUbUNYWWuAeaetItjUFhGsLRRR5G4cEql8EyKEbSFWk+/Awc0b2rwvFpF958C+BJ96I0R4dVhN2Huuy4N9wxZ+W0dIaNaoGc0RGy4uKTB2YNXpwEI1xrQ1LTek3XmrzRf8+vI+Mr/6ogN+uPTgrRByuOIPnfor5fYJAb3/kD06Gf6vR+Vy7am/cD5JOzSlcQML4D3Pes9oEXKkCt1nloWRFJpXkoU1cxJnEzYUa8jlKoWPYhLiMwEw0AFTtBf0MqIATj8cbkn5xA1R1jnpaXgn2r0fWdXAgMBAAECggEBAIq5rSr606/gPRC+4I9kFk8ufYIVKFd/9Nziz5jT7cUHZyrc0a0uL69rzfu4wBYZqBmYH+Echq8EeiPtfHI3/F/9xAtImeHGo1L0Z5ujE5nalVeRNUvsyRA5IU+Rn9ETV+lmYS/2ABwfonEItksPTQ35CR2gAgw1dih3xGVRW23vCHkZY+akHjGwYD4X0oqutmLvHN2568gvm6ycd8iUeTaMKDwlhP0r7kcQEUV3xA8VKoF/HMc8miCsy0cNfFnmesj5/CppUNh7ujmtoeFjMCqTiSnuVEmUoj86LwQcVmTcdxKiOR2bPgH2bRxuIGIKvdha6EbFbT+/WYmvrjZG8lkCgYEAzvPoqe+dsSNQjyzShU8cNNX1NZOm0n28HsGhJUD2YJXypYlVtBkiln8yv6oTwnRwCo0Y9FsNTUGUuSP7YJ0EjFeJmwylAKllDb2uKNYL08S5waTkHjaS1USDDnSMhzUPgf38vKxy3sKZb126yziXXBDWVSygABzYdBMSKQ/XNt0CgYEAxxMmy0dhFvYoFNdNm2k/duAQmAdQgO2biPhOl4iaPoaBMK2yiSVo7zZX9Rl589YUR6El8Fjcv3nqj9ce7RASIjE48V5a/Z6Hnjl5RE1oxweaD2m8GddvfI3ZHwpOCatr48R0gKNVU8Z6ehR1j+Z8qBKuETF+RLUn/dpvF2fawcMCgYA6bcSnjd5Ar87DzYzWVGKLTEkBymEUFqmxKUvc371vwYYTVHXc9ie8w8bJNDSF9yfW4sVD4B0eTcC2kMEdItew49oW63f+etTsDzyHjP8j1+v2Dx7UpOXJzqENyLwQRFvPgK0Fe86ms9xsA9OEIsMhHCPXQlUeEwbNpsC+1RkXBQKBgQC1hHvyDJK1qhuf7TVhSJVKokHfLYQ1GvKf8LFQsIjcDC7OIQNS1B6bR7Tp0qIFOKVjLsf2IECgAt1i7KbRR78RGEqwovVanetQ1V0Cb4bjO8Y42ZNfCLYqHvjjubSwUnLcyuvjw4pxCd/xYqhTXrk5U1cObE+S/I+Lg1maQOMRmwKBgDiVAJfLy1fOgGh41ei1jiDquGFtIpsee7BhToZ/ueoXwXDPNebKsOXEmuBXbGiF1dnNxVyecIoPEGYg5p4mns9glFno6i0/gTKibazVptjuBHVE27DMvh0Q51c4qkGf5kQcvxyvGOm/M4mn5kPDpO46QU5R0Vj2+ETLLka99PAl"

	data := make(map[string]interface{})
	data["app_id"] = "2020063014584944251674"
	data["sign_type"] = "RSA2"
	data["timestamp"] = strext.ToString(time.Now().Unix())
	data["amount"] = 600
	data["order_no"] = "50012020070810330656635622"
	data["account_no"] = "e9586425-bfb7-4054-88b2-f1dfa47bdfa3"
	//data["payment_password"] = ""
	//data["non_str"] = ""

	dataStr := ParamsMapToString(data)
	fmt.Println("dataStr:", dataStr)

	sign, err := ss_rsa.RSA2Sign(dataStr, privKey)
	if err != nil {
		fmt.Println("RSA2Sign-err:", err)
		return
	}
	data["sign"] = sign

	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json-err:", err)
		return
	}

	fmt.Println("jsonStr:", string(bytes))
}
func Query() {
	// 商家的私钥
	var privKey = "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCg7xVS54c4O0I0pakkCH2hA3qTei+J7Xxc29YEqlMxutcx2fJVOCwbscYhUbUNYWWuAeaetItjUFhGsLRRR5G4cEql8EyKEbSFWk+/Awc0b2rwvFpF958C+BJ96I0R4dVhN2Huuy4N9wxZ+W0dIaNaoGc0RGy4uKTB2YNXpwEI1xrQ1LTek3XmrzRf8+vI+Mr/6ogN+uPTgrRByuOIPnfor5fYJAb3/kD06Gf6vR+Vy7am/cD5JOzSlcQML4D3Pes9oEXKkCt1nloWRFJpXkoU1cxJnEzYUa8jlKoWPYhLiMwEw0AFTtBf0MqIATj8cbkn5xA1R1jnpaXgn2r0fWdXAgMBAAECggEBAIq5rSr606/gPRC+4I9kFk8ufYIVKFd/9Nziz5jT7cUHZyrc0a0uL69rzfu4wBYZqBmYH+Echq8EeiPtfHI3/F/9xAtImeHGo1L0Z5ujE5nalVeRNUvsyRA5IU+Rn9ETV+lmYS/2ABwfonEItksPTQ35CR2gAgw1dih3xGVRW23vCHkZY+akHjGwYD4X0oqutmLvHN2568gvm6ycd8iUeTaMKDwlhP0r7kcQEUV3xA8VKoF/HMc8miCsy0cNfFnmesj5/CppUNh7ujmtoeFjMCqTiSnuVEmUoj86LwQcVmTcdxKiOR2bPgH2bRxuIGIKvdha6EbFbT+/WYmvrjZG8lkCgYEAzvPoqe+dsSNQjyzShU8cNNX1NZOm0n28HsGhJUD2YJXypYlVtBkiln8yv6oTwnRwCo0Y9FsNTUGUuSP7YJ0EjFeJmwylAKllDb2uKNYL08S5waTkHjaS1USDDnSMhzUPgf38vKxy3sKZb126yziXXBDWVSygABzYdBMSKQ/XNt0CgYEAxxMmy0dhFvYoFNdNm2k/duAQmAdQgO2biPhOl4iaPoaBMK2yiSVo7zZX9Rl589YUR6El8Fjcv3nqj9ce7RASIjE48V5a/Z6Hnjl5RE1oxweaD2m8GddvfI3ZHwpOCatr48R0gKNVU8Z6ehR1j+Z8qBKuETF+RLUn/dpvF2fawcMCgYA6bcSnjd5Ar87DzYzWVGKLTEkBymEUFqmxKUvc371vwYYTVHXc9ie8w8bJNDSF9yfW4sVD4B0eTcC2kMEdItew49oW63f+etTsDzyHjP8j1+v2Dx7UpOXJzqENyLwQRFvPgK0Fe86ms9xsA9OEIsMhHCPXQlUeEwbNpsC+1RkXBQKBgQC1hHvyDJK1qhuf7TVhSJVKokHfLYQ1GvKf8LFQsIjcDC7OIQNS1B6bR7Tp0qIFOKVjLsf2IECgAt1i7KbRR78RGEqwovVanetQ1V0Cb4bjO8Y42ZNfCLYqHvjjubSwUnLcyuvjw4pxCd/xYqhTXrk5U1cObE+S/I+Lg1maQOMRmwKBgDiVAJfLy1fOgGh41ei1jiDquGFtIpsee7BhToZ/ueoXwXDPNebKsOXEmuBXbGiF1dnNxVyecIoPEGYg5p4mns9glFno6i0/gTKibazVptjuBHVE27DMvh0Q51c4qkGf5kQcvxyvGOm/M4mn5kPDpO46QU5R0Vj2+ETLLka99PAl"

	data := make(map[string]interface{})
	data["app_id"] = "2020063014584944251674"
	data["sign_type"] = "RSA2"
	data["timestamp"] = "2020-07-07 13:37:16"
	data["order_no"] = "50002020070713563280380470"
	data["out_order_no"] = "outOrderNo4444444444444444"

	dataStr := ParamsMapToString(data)
	fmt.Println("dataStr:", dataStr)

	sign, err := ss_rsa.RSA2Sign(dataStr, privKey)
	if err != nil {
		fmt.Println("RSA2Sign-err:", err)
		return
	}
	data["sign"] = sign

	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json-err:", err)
		return
	}

	fmt.Println("jsonStr:", string(bytes))
}

func RSA2Sgin() {
	// 商家的私钥
	var privKey = "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCg7xVS54c4O0I0pakkCH2hA3qTei+J7Xxc29YEqlMxutcx2fJVOCwbscYhUbUNYWWuAeaetItjUFhGsLRRR5G4cEql8EyKEbSFWk+/Awc0b2rwvFpF958C+BJ96I0R4dVhN2Huuy4N9wxZ+W0dIaNaoGc0RGy4uKTB2YNXpwEI1xrQ1LTek3XmrzRf8+vI+Mr/6ogN+uPTgrRByuOIPnfor5fYJAb3/kD06Gf6vR+Vy7am/cD5JOzSlcQML4D3Pes9oEXKkCt1nloWRFJpXkoU1cxJnEzYUa8jlKoWPYhLiMwEw0AFTtBf0MqIATj8cbkn5xA1R1jnpaXgn2r0fWdXAgMBAAECggEBAIq5rSr606/gPRC+4I9kFk8ufYIVKFd/9Nziz5jT7cUHZyrc0a0uL69rzfu4wBYZqBmYH+Echq8EeiPtfHI3/F/9xAtImeHGo1L0Z5ujE5nalVeRNUvsyRA5IU+Rn9ETV+lmYS/2ABwfonEItksPTQ35CR2gAgw1dih3xGVRW23vCHkZY+akHjGwYD4X0oqutmLvHN2568gvm6ycd8iUeTaMKDwlhP0r7kcQEUV3xA8VKoF/HMc8miCsy0cNfFnmesj5/CppUNh7ujmtoeFjMCqTiSnuVEmUoj86LwQcVmTcdxKiOR2bPgH2bRxuIGIKvdha6EbFbT+/WYmvrjZG8lkCgYEAzvPoqe+dsSNQjyzShU8cNNX1NZOm0n28HsGhJUD2YJXypYlVtBkiln8yv6oTwnRwCo0Y9FsNTUGUuSP7YJ0EjFeJmwylAKllDb2uKNYL08S5waTkHjaS1USDDnSMhzUPgf38vKxy3sKZb126yziXXBDWVSygABzYdBMSKQ/XNt0CgYEAxxMmy0dhFvYoFNdNm2k/duAQmAdQgO2biPhOl4iaPoaBMK2yiSVo7zZX9Rl589YUR6El8Fjcv3nqj9ce7RASIjE48V5a/Z6Hnjl5RE1oxweaD2m8GddvfI3ZHwpOCatr48R0gKNVU8Z6ehR1j+Z8qBKuETF+RLUn/dpvF2fawcMCgYA6bcSnjd5Ar87DzYzWVGKLTEkBymEUFqmxKUvc371vwYYTVHXc9ie8w8bJNDSF9yfW4sVD4B0eTcC2kMEdItew49oW63f+etTsDzyHjP8j1+v2Dx7UpOXJzqENyLwQRFvPgK0Fe86ms9xsA9OEIsMhHCPXQlUeEwbNpsC+1RkXBQKBgQC1hHvyDJK1qhuf7TVhSJVKokHfLYQ1GvKf8LFQsIjcDC7OIQNS1B6bR7Tp0qIFOKVjLsf2IECgAt1i7KbRR78RGEqwovVanetQ1V0Cb4bjO8Y42ZNfCLYqHvjjubSwUnLcyuvjw4pxCd/xYqhTXrk5U1cObE+S/I+Lg1maQOMRmwKBgDiVAJfLy1fOgGh41ei1jiDquGFtIpsee7BhToZ/ueoXwXDPNebKsOXEmuBXbGiF1dnNxVyecIoPEGYg5p4mns9glFno6i0/gTKibazVptjuBHVE27DMvh0Q51c4qkGf5kQcvxyvGOm/M4mn5kPDpO46QU5R0Vj2+ETLLka99PAl"

	data := make(map[string]interface{})
	//data["charset"] = "utf-8"
	data["app_id"] = "2020063014584944251674"
	data["sign_type"] = "RSA2"
	data["timestamp"] = strext.ToString(time.Now().Unix())
	data["data"] = "dddddddddddddddddddddddd"
	data["out_order_no"] = "outOrderNo555555555555"
	data["amount"] = 600
	data["notify_url"] = "www.baidu.com"
	data["return_url"] = "www.google.com"
	data["currency_type"] = "USD"

	dataStr := ParamsMapToString(data)
	fmt.Println("dataStr:", dataStr)

	sign, err := ss_rsa.RSA2Sign(dataStr, privKey)
	if err != nil {
		fmt.Println("RSA2Sign-err:", err)
		return
	}
	data["sign"] = sign

	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json-err:", err)
		return
	}

	fmt.Println("jsonStr:", string(bytes))
}
func RSA2Verify() {
	// 平台的公钥
	var pubKey = "IIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAj9vjUaaxrZOd2thOcHdRaBL0IVR6OOWDyiCSsuorQ3BMH+/U54B+Q7uKrQXMb8zWnP0PpyrWerUBlXTSdzFPFoM4qyniLQlN3/2QKdRVwevv2O29f3pMpM1vX01mFkXMKjnCqRpmV+WUAO5gbgBKPNBFDXJxxfuJ5cXWRG96/4qSXyfm4PJ+tDwtqOwEJOI6bfZZ/gXWzxk+9gtl3V1/x+Mr29ZxTYaw9JXZhETVzrw5gkpGLioELu4nmVR6VomBTMyYihhdahDnKxN53lCbePZ6YmowxzG7DXt3wuKfg/CN3NMBNqOke1Yjhfa0d3sDRI4RTm0gPK5mUkO3MqydUQIDAQAB"

	// 接口返回的数据
	respData := `{"code":"0","msg":"正常","order_no":"orderNo111111111111111","out_order_no":"outOrderNo222222222222","sign":"X8emIk+bC67ui3ypO0fL3pcmWucgGnvO2mV7RwJ64zYK6rdDOEC9NagjxmRvgRjgESRsnRq+vvI4k/V4pfI6pj5TmpMNsVU13iw5ZYmJChofp1nR9rErOoLJYsqXw7h/Za3o0JHZqa1uCvgjZKO4fZK4mPO7FrEOaBxs7Za+7DOhzW9rul8OM0rUDM1GVWMAcRe2yw5YmRPjI0YMkr5vaAqPFGJXQbK1tOq1jQfOIo27qEIN7cMODQw/qE962hSq9dAI6VyQHHDHCmCUJIf3xqXELc2/xPxpUwiFP2DZD/+rc7uWMkr9J5kc30AfhqGCGqZLBwawI7RsJvjomAOPYw=="}`

	paramsMap := make(map[string]interface{})

	jerr := json.Unmarshal([]byte(respData), &paramsMap)
	if jerr != nil {
		fmt.Println("接口返回数据,json解码失败, err:", jerr)
		return
	}

	sign, exists := paramsMap["sign"]
	if !exists {
		fmt.Println("返回数据中没有sign字段")
		return
	}

	// 将参数排序并拼接成字符串
	data := ParamsMapToString(paramsMap)

	// rsa2验签
	err := ss_rsa.RSA2Verify(data, fmt.Sprintf("%s", sign), pubKey)
	fmt.Println("验证结果:", err)
}
func ParamsMapToString(params map[string]interface{}) string {
	var pList = make([]string, 0)

	for key, value := range params {
		if strings.TrimSpace(key) == "sign" {
			continue
		}
		// 将interface转换为字符串
		val := strings.TrimSpace(fmt.Sprintf("%v", value))
		if len(val) > 0 {
			pList = append(pList, key+"="+val)
		}
	}

	// 按键排序
	sort.Strings(pList)

	// 使用&符号拼接
	return strings.Join(pList, "&")
}

func main() {
	list := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

	var wg sync.WaitGroup
	parallelNum := 5 // 并发执行数量
	n := 0

	for i, v := range list {
		if n > 0 && n%parallelNum == 0 {
			wg.Wait()
			fmt.Printf("--------------------------------n:%d \n", n)
		}

		n++

		wg.Add(1)

		go func(index int, value string) {
			defer wg.Done()
			show(index, value)
		}(i, v)
	}
	wg.Wait()

	fmt.Println("end...")
}

func show(index int, value string) {
	tt := time.Millisecond * time.Duration(RandInt64(1, 11)*100)

	time.Sleep(tt)

	fmt.Printf("index:%d, value:%s, tt:%v \n", index, value, tt)
}

// 根据数字范围获取随机数(注意：取值范围小于max)
func RandInt64(min, max int64) int64 {
	if min < 0 || max < 1 || min >= max {
		return 0
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63n(max-min) + min
}
