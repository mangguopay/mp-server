package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_sql"
)

type LogToBusinessDao struct{}

var LogToBusinessDaoInst LogToBusinessDao

func (LogToBusinessDao) Insert(tx *sql.Tx, currencyType, businessNo, collectionType, cardNo, amount, realAmount, fees string) string {
	logNoT := strext.GetDailyId()
	sqlStr := `insert into log_to_business(log_no,currency_type,business_no,collection_type,card_no,amount,real_amount,fee,order_type,order_status,create_time) 
				values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,current_timestamp)`
	err := ss_sql.ExecTx(tx, sqlStr,
		logNoT, currencyType, businessNo, collectionType, cardNo,
		amount, realAmount, fees, constants.TRANSFER_TYPE_BILL, constants.WithdrawalOrderStatusPending)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return logNoT
}

func (LogToBusinessDao) GetToBusinessDetail(logNo string) (dataT *go_micro_srv_cust.ToBusinessData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereList := []*model.WhereSqlCond{
		{Key: "tb.log_no", Val: logNo, EqType: "="},
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlStr := `SELECT 
			tb.log_no, tb.currency_type, tb.business_no, tb.collection_type, tb.card_no,
			tb.amount, tb.create_time, tb.order_type, tb.order_status, tb.finish_time,
			tb.fee, tb.image_id, tb.notes, tb.real_amount, cb.card_number,
			acc.account 
		 FROM log_to_business tb 
		 left join card_business cb on cb.card_no = tb.card_no	 
		 left join business b on b.business_no = tb.business_no 
		 left join account acc on acc.uid = b.account_no ` + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	data := &go_micro_srv_cust.ToBusinessData{}
	err2 = rows.Scan(
		&data.LogNo,
		&data.CurrencyType,
		&data.BusinessNo,
		&data.CollectionType,
		&data.CardNo,

		&data.Amount,
		&data.CreateTime,
		&data.OrderType,
		&data.OrderStatus,
		&data.FinishTime,

		&data.Fee,
		&data.ImageId,
		&data.Notes,
		&data.RealAmount,
		&data.CardNumber,

		&data.Account,
	)

	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return nil, err2
	}

	return data, nil
}

func (LogToBusinessDao) UpdateStatusFromLogNo(tx *sql.Tx, logNo, orderStatus, notes, imgNo string) error {
	sqlStr := `update log_to_business set order_status = $2, notes = $3, image_id = $4 where log_no = $1 `
	if err := ss_sql.ExecTx(tx, sqlStr, logNo, orderStatus, notes, imgNo); err != nil {
		ss_log.Error("err=%v", err)
		return err
	}
	return nil
}
