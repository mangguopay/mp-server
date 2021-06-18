package dao

import (
	"a.a/cu/strext"
	"testing"
)

func TestBusinessNotifyLogDao_InsertLog(t *testing.T) {
	data := &BusinessNotifyLog{
		OrderNo:    "1111111",
		OutOrderNo: "2222222",
		Status:     0,
	}
	id, err := BusinessNotifyLogDaoInst.InsertLog(data)
	if err != nil {
		t.Errorf("InsertLog() error = %v", err)
	}
	t.Logf("插入成功, id=%v", id)
}

func TestBusinessNotifyLogDao_UpdateNotifyResultById(t *testing.T) {
	resultMap := make(map[string]interface{})
	data := &BusinessNotifyLog{
		LogId:  "",
		Result: strext.ToJson(resultMap),
		Status: 2,
	}
	if err := BusinessNotifyLogDaoInst.UpdateNotifyResultById(data); err != nil {
		t.Errorf("UpdateNotifyResultById() error = %v", err)
	}

	t.Logf("修改成功")
}
