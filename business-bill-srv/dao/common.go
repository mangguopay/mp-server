package dao

import (
	"errors"
)

var GetDBConnectFailedErr = errors.New("获取数据库连接失败")

var InsertBusinessCheckingLogErr = errors.New("pq: duplicate key value violates unique constraint \"business_balance\"")
