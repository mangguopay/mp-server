package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type EnterpriseTransferOrderDao struct {
	TransferNo        string
	OutTransferNo     string
	Amount            int64
	RealAmount        int64
	Rate              int64
	Fee               int64
	CurrencyType      string
	OrderStatus       string
	PaymentType       string
	BusinessAccountNo string
	AppId             string
	PayeeAccountNo    string
	PayeeAccount      string
	NotifyUrl         string
	NotifyStatus      string
	NotifyFailTimes   string
	NextNotifyTime    string
	Remark            string
	CreateTime        string
}

var EnterpriseTransferOrderDaoInst EnterpriseTransferOrderDao

func (EnterpriseTransferOrderDao) UpdateOrderStatusById(transferNo, orderStatus, wrongReason string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "UPDATE enterprise_transfer_order SET order_status=$3, wrong_reason = $4, finish_time=CURRENT_TIMESTAMP " +
		"WHERE transfer_no = $1 AND order_status = $2 "
	return ss_sql.Exec(dbHandler, sqlStr, transferNo, constants.BusinessTransferOrderStatusPending, orderStatus, wrongReason)
}
