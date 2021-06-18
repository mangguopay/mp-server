package dao

import (
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/global"
	_ "a.a/mp-server/cust-srv/test"
	"testing"
)

func TestSceneSignedDao_GetExpireSignedList(t *testing.T) {
	expireTime := ss_time.Now(global.Tz).AddDate(0, 0, 3).Format(ss_time.DateTimeDashFormat)
	got, err := SceneSignedDaoInst.GetExpireSignedList(expireTime)
	if err != nil {
		t.Errorf("GetExpireSignedList() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToString(got))
}

func TestSceneSignedDao_AutoSigned(t *testing.T) {
	d := &SceneSignedDao{
		SignedNo:      strext.GetDailyId(),
		Status:        "",
		CreateTime:    "",
		StartTime:     "",
		EndTime:       "",
		Rate:          "",
		Cycle:         "",
		SceneNo:       "",
		IndustryNo:    "",
		BusinessNo:    "",
		BusinessAccNo: "",
		LastSignedNo:  "",
	}
	gotSignedNoT, err := SceneSignedDaoInst.AutoSigned(d)
	if err != nil {
		t.Errorf("AutoSigned() error = %v", err)
		return
	}
	t.Logf("查询成功，%v", gotSignedNoT)
}
