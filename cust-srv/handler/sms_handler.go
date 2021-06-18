package handler

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	go_micro_srv_push "a.a/mp-server/common/proto/push"
	"a.a/mp-server/common/ss_data"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/cust-srv/common"
	util2 "a.a/mp-server/cust-srv/util"
	"context"
	"errors"
	"time"
)

func redisErrorIsNil(redisErr error) bool {
	return redisErr.Error() == "redis: nil"
}

func (CustHandler) RegSms(ctx context.Context, req *go_micro_srv_cust.RegSmsRequest, reply *go_micro_srv_cust.RegSmsReply) error {
	num := util.RandomDigitStrOnlyNum(6)
	ss_log.Info("sms=[%v],lang=[%v]", num, req.Lang)

	// 防刷机制,限制40秒发一次
	key := ss_data.GetSMSKey(req.Function, req.Phone+req.Function)
	value, err := cache.RedisClient.Get(key).Result() // 这个key作为防刷特殊处理
	if err != nil && !redisErrorIsNil(err) {
		ss_log.Error("不能频繁发送短信  key=[%v], err=[%v]", key, err)
		return errors.New("不能频繁发送短信")
	}
	if value != "" { // redis里面有值,证明上次发送还没有过期
		ss_log.Error("不能频繁发送短信  key=[%v], err=[%v]", key, err)
		return errors.New("不能频繁发送短信")
	}

	//// 获取短信商户
	//k1, business, err := cache.ApiDaoInstance.GetGlobalParam("sms_business") // cl
	//if err != nil {
	//	ss_log.Error("err=[%v],missing key=[%v]", err, k1)
	//}

	// 消息推送
	ev := go_micro_srv_push.PushReqest{
		Accounts: []*go_micro_srv_push.PushAccout{
			{
				Phone:       req.Phone,
				Lang:        req.Lang,
				CountryCode: req.CountryCode,
			},
		},
		TempNo: constants.Template_Reg,
		Args: []string{
			num,
		},
	}
	ss_log.Info("publishing %+v\n", ev)
	// publish an event
	if err := common.SmsEvent.Publish(context.TODO(), &ev); err != nil {
		ss_log.Error("error publishing: %v", err)
	}

	ss_log.Info("---------redisKey: %s", ss_data.GetSMSKey(req.Function, req.Phone))
	// 设置做防刷标记使用
	if err := cache.RedisClient.Set(ss_data.GetSMSKey(req.Function, req.Phone+req.Function), num, constants.SmsKeySecV2).Err(); err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err.Error(), "存放短信验证码进redis失败,手机号为"+ss_data.GetSMSKey(req.Function, req.Phone))
	}

	// 保存验证码
	if err := cache.RedisClient.Set(ss_data.GetSMSKey(req.Function, req.Phone), num, time.Minute*20).Err(); err != nil {
		ss_log.Error("设置redis失败,key:%s,value:%s, err:%v", util2.MkSmsCodeName(num), strext.ToStringNoPoint(req.Phone), err)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 验证短信验证码是否正确接口
func (CustHandler) CheckSms(ctx context.Context, req *go_micro_srv_cust.CheckSmsRequest, reply *go_micro_srv_cust.CheckSmsReply) error {
	// 校验短信验证码开关
	k1, isCheck, err := cache.ApiDaoInstance.GetGlobalParam("is_check_sms") // is_check_sms 0-需要校验,1-不需要校验
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k1)
	}

	// 校验短信验证码是否正确
	if isCheck == "0" {
		if !cache.CheckSMS(req.Function, req.Phone, req.Sms) {
			reply.ResultCode = ss_err.ERR_Business_Verification_Code_FAILD
			return nil
		}
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
