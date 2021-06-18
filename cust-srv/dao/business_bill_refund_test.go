package dao

import (
	"a.a/cu/strext"
	"a.a/mp-server/common/model"
	"a.a/mp-server/common/ss_sql"
	_ "a.a/mp-server/cust-srv/test"
	"testing"
)

func TestBusinessBillRefundDao_GetOrderDetail(t *testing.T) {
	orderNo := "2020082719372727728603"
	got, err := BusinessBillRefundDaoInst.GetOrderDetail(orderNo)
	if err != nil {
		t.Errorf("GetOrderDetail() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(got))
}

func TestBusinessBillRefundDao_GetRefundBills(t *testing.T) {
	whereList := []*model.WhereSqlCond{
		{Key: "bb.business_no", Val: "53be2111-b1cb-4041-8143-d4a2ccf7d995", EqType: "="},
	}
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY br.finish_time desc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt("10"), strext.ToInt("1"))

	list, err := BusinessBillRefundDaoInst.GetRefundBills(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		t.Errorf("GetRefundBills() error = %v", err)
		return
	}
	t.Logf("列表：%v", strext.ToJson(list))
}
