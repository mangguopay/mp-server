package test

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"unsafe"
)

func TestA(t *testing.T) {
	//s := `【Modern Pay】 您的验证码为：666666，请勿泄露！`
	s := `【Modern Pay】您的验证码是：666666。请不要把验证码泄露给其他人`
	SendSms("3590", "khmodern", "modern12abcZ", s, "13570213647")
}
func SendSms(userid, account, passwd, msg, phone string) {
	body := url.Values{
		"action":   {"send"},
		"account":  {account},
		"password": {passwd},
		"content":  {msg},
		"mobile":   {phone},
		"userid":   {userid},
		"sendTime": {""},
		"extno":    {""},
	}

	resp, err := http.PostForm("http://122.114.79.52:6688/sms.aspx", body)
	defer resp.Body.Close()
	if err != nil {
		ss_log.Error("err=%v\n", err)
		return
	}

	body2, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		ss_log.Error("err1=%v\n", err1)
		return
	}
	zt := new(ztResp)
	body3 := strings.Replace(strext.ToStringNoPoint(body2), "\n", "", -1)
	if err := xml.Unmarshal([]byte(body3), &zt); err != nil {
		ss_log.Error("---------------%s", err.Error())
		return

	}
	fmt.Printf("%+v\n", zt)

	//retI := make(xml2.XmlStringMapDecoder)
	/*	retI := make(xml2.XmlStringMapDecoder)
		body3 := strings.Replace(strext.ToStringNoPoint(body2), "\n", "", -1)
		errUnmarshal := xml.Unmarshal([]byte(body3), &retI)
		ret := map[string]interface{}(retI)
		log.Printf("----1\n%v\n---2\n%v\n", errUnmarshal, ret)
		toString, _ := jsoniter.MarshalToString(ret)
		fmt.Println("==", toString)
		fmt.Println("message---------->", ret["returnsms"].(map[string]interface{})["message"])
		fmt.Println("remainpoint---------->", ret["returnsms"].(map[string]interface{})["remainpoint"])
		fmt.Println("returnstatus---------->", ret["returnsms"].(map[string]interface{})["returnstatus"])
		fmt.Println("successCounts---------->", ret["returnsms"].(map[string]interface{})["successCounts"])
		fmt.Println("taskID---------->", ret["returnsms"].(map[string]interface{})["taskID"])
		//map[returnsms:map[message:ok remainpoint:2 returnstatus:Success successCounts:1 taskID:8892290]]*/
}

type ztResp struct {
	Returnstatus  string `xml:"returnstatus"`
	Message       string `xml:"message"`
	Remainpoint   string `xml:"remainpoint"`
	TaskID        string `xml:"taskID"`
	SuccessCounts string `xml:"successCounts"`
}

func TestLat(t *testing.T) {
	lat := `1354`
	fmt.Println(lat)
	fmt.Println(strext.ToFloat64(lat))
	fmt.Println(fmt.Sprintf("%.8f", strext.ToFloat64(lat)))
}

func TestLat111(t *testing.T) {
	fun1()
	fun2()
	fun3()
}

func fun1() {
	a := 2
	c := (*string)(unsafe.Pointer(&a))
	//fmt.Println(*c)
	*c = "44"
	fmt.Println(*c)
	fmt.Println(&c)
	fmt.Println(a)
}
func fun2() {
	a := "654"
	c := (*string)(unsafe.Pointer(&a))
	*c = "44"
	fmt.Println(*c)
}
func fun3() {
	a := 3
	c := *(*string)(unsafe.Pointer(&a))
	c = "445"
	fmt.Println(c)
}
