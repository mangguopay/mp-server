package dao

import (
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type EnterpriseTransferOrderDao struct {
	Id              string
	AppId           string
	NotifyUrl       string
	NotifyStatus    string
	NotifyFailTimes string
	NextNotifyTime  string
	CreateTime      string
	TransferLogNo   string
}

var EnterpriseTransferOrderDaoInst EnterpriseTransferOrderDao

func (EnterpriseTransferOrderDao) InsertTx(tx *sql.Tx, d *EnterpriseTransferOrderDao) (id string, err error) {
	idT := strext.GetDailyId()
	sqlStr := "INSERT INTO enterprise_transfer_order(id, app_id, notify_url, notify_status, transfer_log_no, create_time) " +
		"VALUES($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)"

	if err := ss_sql.ExecTx(tx, sqlStr, idT, d.AppId, d.NotifyUrl, constants.NotifyStatusNOT, d.TransferLogNo); err != nil {
		return "", err
	}
	return idT, nil
}
