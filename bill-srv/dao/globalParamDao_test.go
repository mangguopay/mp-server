package dao

import (
	"a.a/cu/strext"
	_ "a.a/mp-server/bill-srv/test"
	"testing"
)

func TestGlobalParamDao_GetBusinessTransferParamValue(t *testing.T) {
	got, err := GlobalParamDaoInstance.GetBusinessTransferParamValue()
	if err != nil {
		t.Errorf("GetBusinessTransferParamValue() error = %v", err)
		return
	}
	t.Logf(strext.ToJson(got))
}
