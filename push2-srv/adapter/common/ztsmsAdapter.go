package common

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	go_micro_srv_push "a.a/mp-server/common/proto/push"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/push2-srv/dao"
	"a.a/mp-server/push2-srv/m"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	ZtSmsAdapterInst ZtSmsAdapter
)

type ZtSmsAdapter struct {
}

func (f ZtSmsAdapter) GetAccToken(accReq *go_micro_srv_push.PushAccout) (string, string) {
	ss_log.Info("GetAccToken -------------->%v", accReq)
	var phone string
	if accReq.Phone != "" {
		phone = accReq.Phone
	} else {
		_, phone = dao.AccDaoInst.GetPhone(accReq.AccountNo)
	}
	return phone, ""
}

func (f ZtSmsAdapter) Send(req m.SendReq) (string, string) {
	action := "send"
	account := strext.ToStringNoPoint(req.Config["account"])
	password := strext.ToStringNoPoint(req.Config["password"])
	userid := strext.ToStringNoPoint(req.Config["userid"])
	path := strext.ToStringNoPoint(req.Config["path"])
	return ztDoSms(action, account, password, userid, path, req.AccToken, req.Content)
}

func ztDoSms(action, account, password, userid, path, accountToken, msg string) (string, string) {
	// 1.获取发送短信的配置信息
	switch "" {
	case action:
		ss_log.Error("err=[%v]", "发送短信,获取 action 失败")
		return "发送短信,获取 action 失败", ss_err.ERR_PARAM
	case account:
		ss_log.Error("err=[%v]", "发送短信,获取 account 失败")
		return "发送短信,获取 account 失败", ss_err.ERR_PARAM
	case password:
		ss_log.Error("err=[%v]", "发送短信,获取 password 失败")
		return "发送短信,获取 password 失败", ss_err.ERR_PARAM
	case userid:
		ss_log.Error("err=[%v]", "发送短信,获取 userid 失败")
		return "发送短信,获取 userid 失败", ss_err.ERR_PARAM
	case path:
		ss_log.Error("err=[%v]", "发送短信,获取 path 失败")
		return "发送短信,获取 path 失败", ss_err.ERR_PARAM
	}

	// 执行发送
	message, ztErr := ztSendSms(action, account, password, userid, path, accountToken, msg)
	if ztErr != nil {
		ss_log.Error("err=[%v],missing key=[%v]", ztErr.Error(), "发送消息的日志存放进数据库失败")
		return message, ss_err.ERR_PARAM
	}
	return message, ss_err.ERR_SUCCESS
}

// 创蓝发送短信
func ztSendSms(action, account, password, userid, path, accountToken, msg string) (string, error) {
	body := url.Values{
		"action":   {action},
		"account":  {account},
		"password": {password},
		"content":  {msg},
		"mobile":   {accountToken},
		"userid":   {userid},
		"sendTime": {""},
		"extno":    {""},
	}

	resp, err := http.PostForm(path, body)
	defer resp.Body.Close()
	if err != nil {
		ss_log.Error("err=%v\n", err)
		return "", err
	}

	body2, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		ss_log.Error("err1=%v\n", err1)
		return "", err1
	}

	zt := new(ztResp)
	body3 := strings.Replace(strext.ToStringNoPoint(body2), "\n", "", -1)
	if err := xml.Unmarshal([]byte(body3), &zt); err != nil {
		ss_log.Error("---------------%s", err.Error())
		return "", err

	}
	ss_log.Info("------->zt_sms_resp : %+v", zt)
	if zt.Message != "ok" {
		return zt.Message, errors.New("发送失败" + zt.Message)
	}
	return zt.Message, nil
}

type ztResp struct {
	Returnstatus  string `xml:"returnstatus"`
	Message       string `xml:"message"`
	Remainpoint   string `xml:"remainpoint"`
	TaskID        string `xml:"taskID"`
	SuccessCounts string `xml:"successCounts"`
}
