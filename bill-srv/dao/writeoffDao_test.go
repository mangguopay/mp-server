package dao

import (
	"a.a/cu/strext"
	_ "a.a/mp-server/bill-srv/test"
	"testing"
)

func TestWriteoff_GetCodeByIncomeOrderNo(t *testing.T) {
	incomeOrderNo := "2020012011555415638414"
	code, err := WriteoffInst.GetCodeByIncomeOrderNo(incomeOrderNo)
	if err != nil {
		t.Errorf("GetCodeByIncomeOrderNo() error = %v", err)
		return
	}

	t.Logf("核销码:%v", code)

}

func TestWriteoff_QueryOrderNo(t *testing.T) {
	code := "8871470928"
	recvPhone := "10092"
	writeoff, err := WriteoffInst.QueryOrderNo(code, recvPhone)
	if err != nil {
		t.Errorf("QueryOrderNo() error = %v", err)
		return
	}

	t.Logf("核销码详情：%v", strext.ToJson(writeoff))

}
