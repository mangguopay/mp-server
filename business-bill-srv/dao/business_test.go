package dao

import (
	"a.a/cu/strext"
	_ "a.a/mp-server/business-bill-srv/test"
	"testing"
)

func TestBusinessDao_QueryAccNoByBusinessNo(t *testing.T) {
	businessNo := "d950f705-688a-4b35-b773-de4497fa7602"
	businessAccNo, err := BusinessDaoInst.QueryAccNoByBusinessNo(businessNo)
	if err != nil {
		t.Errorf("QueryAccNoByBusinessNo() error = %v", err)
		return
	}

	t.Logf("商户账号=%v", businessAccNo)
}

func TestBusinessDao_QueryNameByBusinessNo(t *testing.T) {
	businessNo := "d950f705-688a-4b35-b773-de4497fa7602"

	businessName, simplifyName, err := BusinessDaoInst.QueryNameByBusinessNo(businessNo)
	if err != nil {
		t.Errorf("QueryNameByBusinessNo() error = %v", err)
		return
	}

	t.Logf("商户名=%v,商户简称=%v", businessName, simplifyName)

}

func TestBusinessDao_GetTransConfig(t *testing.T) {
	businessNo := "a79123bd-0d06-418f-9c12-976a8643f82a"
	businessTransConf, err := BusinessDaoInst.GetTransConfig(businessNo)
	if err != nil {
		t.Errorf("GetTransConfig() error = %v", err)
		return
	}

	t.Logf("商户交易配置：%v", strext.ToJson(businessTransConf))
}
