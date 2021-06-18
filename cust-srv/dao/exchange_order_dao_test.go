package dao

import (
	"a.a/mp-server/common/ss_sql"
	"fmt"
	"testing"
)

func TestExchngeCount(t *testing.T) {
	data, err := ExchangeOrderDaoInst.GetExchangeCountByDate("2020-05-22")
	if err != nil && err.Error() != ss_sql.DB_NO_ROWS_MSG {
		fmt.Println("兑换统计操作失败,err: ", err.Error())
		return
	}
	// 插入数据
	if data != nil {
		if err := StatisticUserExchangeDaoInst.Insert(data); err != nil {
			fmt.Println("兑换统计插入数据操作失败,err: ", err.Error())
		}
	}
}
