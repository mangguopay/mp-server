package common

import (
	"a.a/cu/ss_log"
	"a.a/cu/ss_mail"
	"a.a/cu/strext"
	go_micro_srv_push "a.a/mp-server/common/proto/push"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/push2-srv/dao"
	"a.a/mp-server/push2-srv/m"
	"fmt"
)

var (
	GoogleEmailAdapterInst GoogleEmailAdapter
)

type GoogleEmailAdapter struct {
}

func (f GoogleEmailAdapter) GetAccToken(accReq *go_micro_srv_push.PushAccout) (string, string) {
	ss_log.Info("GetAccToken -------------->%v", accReq)
	var email string
	if accReq.Email != "" {
		email = accReq.Email
	} else {
		email = dao.AccDaoInst.GetEmail(accReq.AccountNo)
	}
	return email, ""
}

func (f GoogleEmailAdapter) Send(req m.SendReq) (string, string) {
	ss_log.Info("req=[%v]", req)
	fromUser := strext.ToStringNoPoint(req.Config["from_user"])
	authCode := strext.ToStringNoPoint(req.Config["auth_code"])
	smtpHost := strext.ToStringNoPoint(req.Config["smtp_host"])
	bodyTemplete := strext.ToStringNoPoint(req.Config["body_templete"])
	userName := strext.ToStringNoPoint(req.Config["user_name"])
	mimeType := strext.ToStringNoPoint(req.Config["mime_type"])

	body := fmt.Sprintf(bodyTemplete, req.Content)
	ss_log.Info("bodyTemplete=[%v],Content[%v]", bodyTemplete, req.Content)
	ss_log.Info("Title=[%v]", req.Title)

	//title := fmt.Sprintf(req.Title, req.Content)
	// 邮件推送
	if err := ss_mail.SsMailInst.SendMail(userName, fromUser, authCode, smtpHost, req.AccToken, req.Title, body, mimeType); err != nil {
		ss_log.Error("发送邮件失败!err=[%v]", err)
		return err.Error(), ss_err.ERR_PARAM
	}

	return "", ss_err.ERR_SUCCESS
}
