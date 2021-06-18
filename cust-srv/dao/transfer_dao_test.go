package dao

import (
	"a.a/mp-server/common/ss_sql"
	"fmt"
	"testing"
)

func TestTransferCount(t *testing.T) {
	data, err := TransferDaoInst.GetTransferCountByDate("2020-06-04", "khr")
	if err != nil && err.Error() != ss_sql.DB_NO_ROWS_MSG {
		fmt.Println("统计  转账统计操作失败,err:  ", err.Error())
		return
	}
	// 插入数据
	if data != nil {
		if err := StatisticUserTransferDaoInst.Insert(data); err != nil {
			fmt.Println("统计   转账统计插入数据操作失败,err:  ", err.Error())
		}
	}
}
