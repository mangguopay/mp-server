package dao

import (
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"fmt"
	"testing"
)

func TestCount(t *testing.T) {
	date := "2020-06-04"
	// 注册数量
	data, err := AccDaoInstance.GetRegCountByDate(date)
	if err != nil && err.Error() != ss_sql.DB_NO_ROWS_MSG {
		fmt.Printf("统计注册统计操作失败,err: %s\n", err.Error())
	}
	// 新增服务商统计
	srvData, err := ServiceDaoInst.GetRegCountByDate(date)
	if err != nil && err.Error() != ss_sql.DB_NO_ROWS_MSG {
		fmt.Printf("统计新增服务商操作失败,err: %s\n", err.Error())
	}
	data.ServerNum = srvData.ServerNum
	// 插入数据
	if data != nil {
		if err := StatisticDateDaoInst.Insert(data); err != nil {
			fmt.Printf("统计注册统计插入数据操作失败,err: %s\n", err.Error())
		}
	}

}

func TestAccDao_GetAuthInfoByAccountNo(t *testing.T) {
	accountNo := "555d2d86-fef4-42d9-b2f0-a6adb8c3f325"
	accountType := constants.AccountType_EnterpriseBusiness
	auth, err := AccDaoInstance.GetAuthInfoByAccountNo(accountNo, accountType)
	if err != nil {
		t.Errorf("GetAuthStatusByAccountNo() error = %v", err)
		return
	}

	t.Logf("认证信息：%v", strext.ToJson(auth))
}
