package dao

import (
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	_ "a.a/mp-server/cust-srv/test"
	"testing"
	"time"
)

func TestBusinessSignedDao_GetExpireSigned(t *testing.T) {
	expireTime := time.Now().AddDate(0, 0, 3).Format(ss_time.DateTimeDashFormat)
	signedList, err := BusinessSignedDaoInst.GetExpireSignedList(expireTime)
	if err != nil {
		t.Errorf("GetExpireSigned() error = %v", err)
		return
	}

	t.Logf("过期签约：%v", strext.ToJson(signedList))
}

func TestBusinessSignedDao_AutoSignedTx(t *testing.T) {
	serviceTerm := GlobalParamDaoInstance.QeuryParamValue(constants.GlobalParamKeyAppSignedTerm)
	startTime := time.Now()
	d := new(BusinessSignedDao)
	d.AppId = "2020072914053990298762"
	d.StartTime = startTime.Format(ss_time.DateTimeDashFormat)
	d.EndTime = startTime.AddDate(strext.ToInt(serviceTerm), 0, 0).Format(ss_time.DateTimeDashFormat)
	d.BusinessNo = "e371896b-b660-4093-9447-9152ccf72303"
	d.BusinessAccNo = "8deff2bf-0cf0-4a4c-95c4-0074fcfc2b55"
	d.SceneNo = ""
	d.Rate = ""
	d.Cycle = ""
	d.IndustryNo = ""

	signedNoT, err := BusinessSignedDaoInst.AutoSigned(d)
	if err != nil {
		t.Errorf("AutoSignedTx() error = %v", err)
		return
	}

	t.Logf("新的签约号：%v", signedNoT)
}

func TestBusinessSignedDao_UpdateStatusBySignedNo(t *testing.T) {
	expireTime := ss_time.Now(global.Tz).Format(ss_time.DateTimeDashFormat)
	idList, err := BusinessSignedDaoInst.UpdateStatusBySignedNo(expireTime, constants.SignedStatusInvalid)
	if err != nil {
		t.Errorf("UpdateStatusBySignedNo() error = %v", err)
		return
	}
	t.Logf("过期的签约：%v", strext.ToJson(idList))

}

func TestBusinessAppDao_GetSignedExpireApp(t *testing.T) {
	expireTime := time.Now().AddDate(0, 0, 3).Format(ss_time.DateTimeDashFormat)
	app, err := BusinessSignedDaoInst.GetExpireSignedList(expireTime)
	if err != nil {
		t.Errorf("GetSignedExpireApp() error = %v", err)
		return
	}

	t.Logf("签约即将过期的应用:%v", strext.ToJson(app))
}
