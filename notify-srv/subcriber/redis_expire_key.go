package subcriber

import (
	"fmt"
	"strings"
	"time"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/notify-srv/common"
	"a.a/mp-server/notify-srv/handler"
)

type ListenExpKeyHandler struct{}

// 开始监听过期的key
func DoListenRedisExpKey() {
	go PSubscribeExpireKey()
}

// 监听redis过期key事件
func PSubscribeExpireKey() {
	pubSub := cache.RedisClient.PSubscribe(fmt.Sprintf("__keyevent@%v__:expired", constants.NotifySrvRedisDb))
	defer pubSub.Close()

	for {
		msg, err := pubSub.ReceiveMessage()
		if err != nil {
			ss_log.Error("获取订阅过期key出错,err:%v", err)
			time.Sleep(time.Millisecond * 500) // 出错时稍微等待
			continue
		}
		HandleExpireKeyEvent(msg.Payload)
	}
}

// 处理过期key事件
func HandleExpireKeyEvent(payload string) {
	expireKeys := strings.Split(payload, "_")
	if len(expireKeys) != 2 {
		return
	}

	key := expireKeys[0]
	if key == "" {
		return
	}
	orderNo := expireKeys[1]
	if orderNo == "" {
		return
	}

	ss_log.Info("PubCallback获的一条过期数据----->payload:%s", payload)

	lockKey := common.GetLockKey(orderNo)
	lockValue := strext.NewUUID()

	// 获取分布式锁
	if !cache.GetDistributedLock(lockKey, lockValue, 30*time.Second) {
		return
	}

	switch key {
	case common.PayNotifyExpireKey:
		handler.PayNotifyH.PaySuccessNotify(orderNo)
	case common.TransferNotifyExpireKey:
		handler.TransferNotifyH.TransferSuccessNotify(orderNo)
	case common.RefundNotifyExpireKey:
		handler.RefundNotifyH.RefundNotify(orderNo)
	}

	// 释放分布式锁
	cache.ReleaseDistributedLock(lockKey, lockValue)
}
