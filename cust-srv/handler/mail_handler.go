package handler

import (
	"a.a/cu/ss_log"
	"a.a/cu/util"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	go_micro_srv_push "a.a/mp-server/common/proto/push"
	"a.a/mp-server/common/ss_data"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/cust-srv/common"
	"a.a/mp-server/cust-srv/dao"
	"context"
	"errors"
	"time"
)

func (CustHandler) SendMail(ctx context.Context, req *go_micro_srv_cust.SendMailRequest, reply *go_micro_srv_cust.SendMailReply) error {
	num := util.RandomDigitStrOnlyNum(6)
	ss_log.Info("mailCode=[%v],lang=[%v]", num, req.Lang)

	// 防刷机制,限制40秒发一次
	key := ss_data.GetMailKey(req.Function, req.Email+req.Function)
	value, err := cache.RedisClient.Get(key).Result() // 这个key作为防刷特殊处理
	if err != nil && !redisErrorIsNil(err) {
		ss_log.Error("不能频繁发送邮件  key=[%v], err=[%v]", key, err)
		return errors.New("不能频繁发送邮件")
	}
	if value != "" { // redis里面有值,证明上次发送还没有过期
		ss_log.Error("不能频繁发送邮件  key=[%v], err=[%v]", key, err)
		return errors.New("不能频繁发送邮件")
	}

	ss_log.Info("---------redisKey: %s", key)
	// 设置做防刷标记使用
	if err := cache.RedisClient.Set(key, num, constants.MailKeySecV2).Err(); err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err.Error(), "存放短信验证码进redis失败,key为"+key)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 保存验证码
	if err := cache.RedisClient.Set(ss_data.GetMailKey(req.Function, req.Email), num, time.Minute*15).Err(); err != nil {
		ss_log.Error("设置redis失败,key:%s,value:%s, err:%v", ss_data.GetMailKey(req.Function, req.Email), num, err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 消息推送
	ev := go_micro_srv_push.PushReqest{
		Accounts: []*go_micro_srv_push.PushAccout{
			{
				Email: req.Email,
				Lang:  req.Lang,
			},
		},
		TempNo: constants.Template_VerifyEmail,
		Args: []string{
			req.Email,
			num,
		},
		TitleArgs: []string{
			num,
		},
	}
	ss_log.Info("publishing %+v\n", ev)
	// publish an event
	if err := common.SmsEvent.Publish(context.TODO(), &ev); err != nil {
		ss_log.Error("error publishing: %v", err)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 验证邮箱收到的验证码是否正确接口
func (CustHandler) CheckMailCode(ctx context.Context, req *go_micro_srv_cust.CheckMailCodeRequest, reply *go_micro_srv_cust.CheckMailCodeReply) error {
	if !cache.CheckMailCode(req.Function, req.Mail, req.MailCode) {
		reply.ResultCode = ss_err.ERR_ACCOUNT_MailCode_FAILD
		return nil
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 验证邮箱是否唯一
func (CustHandler) CheckMail(ctx context.Context, req *go_micro_srv_cust.CheckMailRequest, reply *go_micro_srv_cust.CheckMailReply) error {
	cnt, err := dao.AccDaoInstance.CheckAccount(req.Email)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_ACCOUNT_Mail_FAILD
		return nil
	}

	if cnt != 0 {
		ss_log.Error("邮箱[%v]已存在，err=[%v]", req.Email, err)
		reply.ResultCode = ss_err.ERR_ACCOUNT_Mail_FAILD
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
