package common

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	go_micro_srv_push "a.a/mp-server/common/proto/push"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/push2-srv/dao"
	"a.a/mp-server/push2-srv/m"
	"fmt"
	"github.com/NaySoftware/go-fcm"
	"net/http"
)

const (
	topic = "/topics/"
)

var (
	FcmAdapterInst FcmAdapter
)

type FcmAdapter struct {
}

func (f FcmAdapter) GetAccToken(accReq *go_micro_srv_push.PushAccout) (string, string) {
	//原账号没有国家码,现在的账号有国家码,所以得换手机号
	_, phone := dao.AccDaoInst.GetPhone(accReq.AccountNo)
	if accReq.AccountType == "" {
		return "", "账号类型缺失"
	}
	if phone == "" {
		return "", "账号不存在"
	}
	return fmt.Sprintf("%s_%s", accReq.AccountType, phone), ""
}

func (f FcmAdapter) Send(req m.SendReq) (string, string) {
	if req.AccToken == "" {
		ss_log.Error("no account")
		return "no account", ss_err.ERR_PUSH_ACCOUNT_IS_NIL
	}
	servKey := strext.ToStringNoPoint(req.Config["server_key"])
	if servKey == "" {
		ss_log.Error("谷歌fcm推送,%s", "serverKey 为空")
		return fmt.Sprintf("谷歌fcm推送,%s", "serverKey 为空"), ss_err.ERR_PUSH_SERVER_KEY_IS_NIL
	}
	ss_log.Info("servKey=[%v]", servKey)
	c := fcm.NewFcmClient(servKey)
	payload := &fcm.NotificationPayload{
		Title: req.Title,
		Body:  req.Content,
	}
	c.SetNotificationPayload(payload)

	if len(req.AccToken) == 0 {
		ss_log.Error("%s", "推送消息目标手机号为空")
		return fmt.Sprintf("%s", "推送消息目标手机号为空"), ss_err.ERR_PUSH_PHONE_IS_NIL
	}
	ss_log.Info("accToken==================== %s", req.AccToken)
	//p := strings.Split(req.AccToken, "_")
	myTopic := fmt.Sprintf("%s%s", topic, req.AccToken)
	ss_log.Info("topic=[%v],payload=[%v]", myTopic, payload)
	c.NewFcmMsgTo(myTopic, payload)
	status, err := c.Send()
	ss_log.Info("fcm resp: %+v", status)
	pushErr := status.Err
	if err != nil {
		ss_log.Error("err=[%s]", err.Error())
		return pushErr, ss_err.ERR_PUSH_FAIL
	}

	if status.StatusCode != http.StatusOK || status.Fail != 0 || status.Err != "" {
		status.PrintResults()
		return pushErr, ss_err.ERR_PUSH_FAIL
	}
	return pushErr, ss_err.ERR_SUCCESS
}
