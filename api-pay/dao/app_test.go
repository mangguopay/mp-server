package dao

import (
	"database/sql"
	"testing"

	_ "a.a/mp-server/api-pay/test"
)

func TestAppDao_GetSignInfo(t *testing.T) {
	appId := "2020082716544263950282"
	obj, err := AppDaoInstance.GetSignInfo(appId)
	if err != nil {
		if err == sql.ErrNoRows {
			t.Errorf("GetSignInfo-empty:%v", err)
			return
		}

		t.Errorf("GetSignInfo-err:%v", err)
		return
	}

	t.Logf("GetSignInfo-result:%+v", obj)
}
