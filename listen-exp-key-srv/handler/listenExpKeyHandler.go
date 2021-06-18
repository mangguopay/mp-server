package handler

import (
	"context"
	"strings"
	"time"

	"a.a/cu/strext"

	"a.a/cu/ss_log"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	go_micro_srv_push "a.a/mp-server/common/proto/push"
	"a.a/mp-server/listen-exp-key-srv/common"
	"a.a/mp-server/listen-exp-key-srv/dao"
)

type ListenExpKeyHandler struct{}

// 开始监听过期的key
func DoListenRedisExpKey() {
	go PSubscribeExpireKey()
}

// 监听redis过期key事件
func PSubscribeExpireKey() {
	pubSub := cache.RedisClient.PSubscribe("__keyevent@0__:expired")
	defer pubSub.Close()

	for {
		msg, err := pubSub.ReceiveMessage()
		if err != nil {
			ss_log.Error("获取订阅过期key出错,err:%v", err)
			time.Sleep(time.Millisecond * 500) // 出错时稍微等待
			continue
		}
		HandleExpireKeyEvent(msg.Channel, msg.Pattern, msg.Payload)
	}
}

// 处理过期key事件
func HandleExpireKeyEvent(channel, pattern, payload string) {
	ss_log.Info("PubCallback获的一条过期数据-----> channel:%s, pattern:%s, payload:%s", channel, pattern, payload)

	// key示例: listenExp-RmTQ4Gqf8uOr35OtTCVfGXYX1dirXIlP
	//if strings.HasPrefix(payload, "listenExp-") {
	if strings.HasPrefix(payload, constants.EXP_GEN_CODE_KEY) {

		dataList := strings.Split(payload, constants.EXP_GEN_CODE_KEY)
		if len(dataList) != 2 {
			ss_log.Error("格式不正确")
			return
		}

		lockKey := common.GetListenExpKeyLockKey(dataList[1])
		lockValue := strext.NewUUID()

		// 获取分布式锁
		if !cache.GetDistributedLock(lockKey, lockValue, 30*time.Second) {
			return
		}

		handler(dataList[1])

		// 释放分布式锁
		cache.ReleaseDistributedLock(lockKey, lockValue)
	}
}

func handler(gendCode string) {

	// 得到码,判断码是否已确认,如果没确认,需要退款
	// 推送到mq
	status, orderNo := dao.GenCodeDaoInst.GetStatusFromCode(gendCode)
	ss_log.Info("gendCode------->%s,orderNo----->%s,status------>%v", gendCode, orderNo, status)
	if status == constants.CODE_Pendding_Confirm { // 这个码还在待确认状态.需要退款.
		ev := &go_micro_srv_push.SendListenExpKeyRequest{
			GenCode: gendCode,
			OrderNo: orderNo,
		}
		if err := common.ListenExpKey.Publish(context.TODO(), ev); err != nil {
			ss_log.Error("err=[pos 存款接口,核销码推送到MQ失败,err----->%s]", err.Error())
		}
		ss_log.Info("pos 存款接口,核销码推送到MQ成功,SendListenExpKeyRequest: %+v", ev)
	}
}
