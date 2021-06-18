package dao

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type LogToHeadquartersDao struct {
}

var LogToHeadquartersDaoInst LogToHeadquartersDao

func (LogToHeadquartersDao) GetLogToHeadquarterOrderStatus(tx *sql.Tx, logNo string) (err error, orderStatus string) {

	var orderStatusT sql.NullString
	sqlStr := "select order_status from log_to_headquarters where log_no=$1 limit 1"
	err = ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&orderStatusT}, logNo)
	if err != nil {
		ss_log.Error("CustDao |AddCust err=[%v]", err)
	}
	return err, orderStatusT.String
}
func (LogToHeadquartersDao) GetLogToHeadquarterFromNo(tx *sql.Tx, logNo string) (string, string, string) {
	var servicerNoT, amountT, currentTypeT sql.NullString
	sqlStr := "select servicer_no,amount,currency_type from log_to_headquarters where log_no=$1 limit 1"
	err := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&servicerNoT, &amountT, &currentTypeT}, logNo)
	if err != nil {
		ss_log.Error("CustDao |AddCust err=[%v]", err)
		return "", "", ""
	}
	return servicerNoT.String, amountT.String, currentTypeT.String
}

func (LogToHeadquartersDao) UpdateLogToHeadquarterOrderStatus(tx *sql.Tx, logNo, orderStatus string) (err error) {
	sqlStr := "update log_to_headquarters set order_status=$2 where log_no=$1 and order_status=$3"
	err = ss_sql.ExecTx(tx, sqlStr, logNo, orderStatus, constants.AuditOrderStatus_Pending)
	if err != nil {
		ss_log.Error("CustDao |AddCust err=[%v]", err)
	}
	return err
}
