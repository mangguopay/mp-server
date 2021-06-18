package dao

import (
	"a.a/cu/strext"
	"a.a/mp-server/common/ss_err"
	_ "a.a/mp-server/cust-srv/test"
	"testing"
)

func TestCardBusinessDao_GetBusinessCards(t *testing.T) {
	accountNo := "0e8d24af-bec7-4f95-b038-c48045f51abf"
	balanceType := "usd"
	accountType := "7"
	datas, total, err := CardBusinessDaoInst.GetBusinessCards(accountNo, balanceType, accountType)
	if err != ss_err.ERR_SUCCESS {
		t.Errorf("查询失败")
		return
	}
	t.Logf("数量：%v", total)
	t.Logf("数据；%v", strext.ToJson(datas))

}

func TestCardBusinessDao_GetBusinessCardDetail(t *testing.T) {
	cardNo := "d6ca1a43-56d0-4119-841c-f726127db88e"
	data, err := CardBusinessDaoInst.GetBusinessCardDetail(cardNo)
	if err != nil {
		t.Errorf("GetBusinessCardDetail() error = %v", err)
		return
	}

	t.Logf("数据；%v", strext.ToJson(data))
}
