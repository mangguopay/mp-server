package dao

import (
	"a.a/mp-server/common/ss_sql"
	"fmt"
	"testing"
)

func TestLogCustToHeadquartersDao(t *testing.T) {
	//data, err := LogCustToHeadquartersDaoInst.GetCustToHeadquartersCountByDate("2020-05-29", "usd")
	//if err != nil {
	//	fmt.Println("--------", err.Error())
	//	return
	//}
	//if data != nil {
	//	if err := StatisticUserRechargeDaoInst.Insert(data); err != nil {
	//		fmt.Println("统计  线上充值统计插入数据操作失败,err: s", err.Error())
	//		return
	//	}
	//}

	incomeCount, err := IncomeOrderDaoInst.GetIncomeOrderCountByDate("2020-05-29", "khr")
	if err != nil && err.Error() != ss_sql.DB_NO_ROWS_MSG {
		fmt.Println("统计 pos端充值统计操作失败,err: ", err.Error())
		return
	}
	// 插入数据
	if incomeCount != nil {
		if err := StatisticUserRechargeDaoInst.Insert(incomeCount); err != nil {
			fmt.Println("统计  pos端充值统计插入数据操作失败,err:", err.Error())
		}
	}

	fmt.Println("------------成功")
}
