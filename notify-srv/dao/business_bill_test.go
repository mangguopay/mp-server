package dao

import (
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	_ "a.a/mp-server/notify-srv/test"
	"testing"
	"time"
)

func TestBusinessBillDao_QueryNotifyBreak(t *testing.T) {
	nextNotifyTime := ss_time.Now(global.Tz).Format(ss_time.DateTimeDashFormat)
	ret, err := BusinessBillDaoInst.QueryNotifyBreak(constants.BusinessOrderStatusPay, constants.NotifyStatusDoing, nextNotifyTime)
	if err != nil {
		t.Errorf("QueryNotifyBreak() error = %v", err)
		return
	}

	t.Logf("ret: %v", ret)

}

func TestBusinessBillDao_QueryNotifyOmission(t *testing.T) {
	payTime := ss_time.Now(global.Tz).Add(-15 * time.Second).Format(ss_time.DateTimeDashFormat)
	ret, err := BusinessBillDaoInst.QueryNotifyOmission(constants.BusinessOrderStatusPay, constants.NotifyStatusNOT, payTime)
	if err != nil {
		t.Errorf("QueryNotifyOmission() error = %v", err)
		return
	}
	t.Logf("ret: %v", ret)
}

func TestBusinessBillDao_UpdateNotifyStatusByOrderNo(t *testing.T) {
	updateDate := new(UpdateNotifyStatus)
	updateDate.OrderNo = ""
	updateDate.NextTime = ""
	updateDate.NotifyStatus = constants.NotifyStatusDoing
	updateDate.NotifyFailTimes = 1
	err := BusinessBillDaoInst.UpdateNotifyStatusByOrderNo(updateDate)
	if err != nil {
		t.Errorf("UpdateNotifyStatusByOrderNo() error = %v", err)
	}

	t.Log("修改成功")
}

func TestBusinessBillDao_QueryOrderInfoByOrderNo(t *testing.T) {
	orderNo := "2020090811284053782038"
	orderInfo, err := BusinessBillDaoInst.QueryOrderInfoByOrderNo(orderNo)
	if err != nil {
		t.Errorf("QueryOrderInfoByOrderNo() error = %v", err)
		return
	}

	t.Logf("订单:%v", strext.ToJson(orderInfo))
}
