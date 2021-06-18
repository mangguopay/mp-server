package handler

import (
	"a.a/cu/encrypt"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	custProto "a.a/mp-server/common/proto/cust"
	_ "a.a/mp-server/cust-srv/test"
	"context"
	"encoding/xml"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestCustHandler_GetBusinessVAccLogList(t *testing.T) {
	req := &custProto.GetBusinessVAccLogListRequest{
		BusinessAccNo: "49307c6f-9e03-4535-b3ef-8aaa581f91bd",
		MoneyType:     "usd",
		StartTime:     "2020-10-01 00:00:00",
		EndTime:       "2020-10-09 23:00:00",
	}
	reply := &custProto.GetBusinessVAccLogListReply{}
	if err := CustHandlerInst.GetBusinessVAccLogList(context.TODO(), req, reply); err != nil {
		t.Errorf("GetBusinessVAccLogList() error = %v", err)
		return
	}
	t.Logf(strext.ToJson(reply))
}

func TestCustHandler_GetBusinessVAccLogDetail(t *testing.T) {
	req := &custProto.GetBusinessVAccLogDetailRequest{
		LogNo:  "2020090716171900137385",
		Reason: "18",
	}
	reply := &custProto.GetBusinessVAccLogDetailReply{}
	if err := CustHandlerInst.GetBusinessVAccLogDetail(context.TODO(), req, reply); err != nil {
		t.Errorf("GetBusinessVAccLogDetail() error = %v", err)
		return
	}
	t.Logf("流水详情：%v", strext.ToJson(reply))
}

func TestCustHandler_GenerateKeys(t *testing.T) {
	req := &custProto.GenerateKeysRequest{
		KeyType: constants.SecretKeyPKCS1,
	}
	reply := &custProto.GenerateKeysReply{}
	if err := CustHandlerInst.GenerateKeys(context.TODO(), req, reply); err != nil {
		t.Errorf("GenerateKeys() error = %v", err)
		return
	}

	t.Logf("结果：%v", strext.ToJson(reply))
}

type MyMap map[string]interface{}

type xmlMapEntry struct {
	XMLName xml.Name
	Value   interface{} `xml:",chardata"`
}

// map转xml
func (m MyMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m {
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
	}

	return e.EncodeToken(start.End())
}
func Test_Xml(t *testing.T) {
	m := make(map[string]interface{})
	m["name"] = "wang"
	m["age"] = 25

	buf, err := xml.Marshal(MyMap(m))
	if err != nil {
		t.Errorf("转化失败：%v", err)
		return
	}
	t.Logf("结果：\n%v", strings.Replace(string(buf), "MyMap", "xml", -1))
}

func Test_Sign(t *testing.T) {
	t.Logf("时间戳：%v", time.Now().Unix())
	ti, err := time.Parse(ss_time.DateTimeDashFormat, "1.598499932e+09")
	if err != nil {
		t.Logf("解析失败: %v", err)
		return
	}

	t.Logf("时间: %v", ti.Format(ss_time.DateTimeDashFormat))

	_, passwordSalt, _ := cache.ApiDaoInstance.GetGlobalParam("password_salt")
	signBefore := fmt.Sprintf("x-login-token=%s&x-lang=%s&x-path=%s&x-ran=%s&key=%s", "xLoginToken", "xLang", "xPath", strext.ToString(time.Now().Unix()), passwordSalt)
	t.Logf("signBefore：\n%v", signBefore)
	sign := encrypt.DoMd5(signBefore)
	t.Logf("sign=[%v]", sign)
}
