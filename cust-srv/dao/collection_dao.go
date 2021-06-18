package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type CollectionDao struct {
}

var CollectionDaoInst CollectionDao

func (*CollectionDao) CustCollectionBillsDetail(logNo string) (data *go_micro_srv_cust.CustCollectionBillsDetailData, returnErr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "log_no", Val: logNo, EqType: "="},
	})

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT coo.log_no, coo.amount, coo.order_status, coo.payment_type, coo.create_time, coo.finish_time, coo.balance_type " +
		", acc.account " +
		" FROM collection_order coo " +
		" LEFT JOIN vaccount vacc ON vacc.vaccount_no = coo.from_vaccount_no " +
		" LEFT JOIN account acc ON acc.uid= vacc.account_no " + where2
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	data = &go_micro_srv_cust.CustCollectionBillsDetailData{}
	if err == nil {
		for rows.Next() {
			var finishTime, paymentType sql.NullString
			err = rows.Scan(
				&data.LogNo,
				&data.Amount,
				&data.OrderStatus,
				&paymentType,
				&data.CreateTime,
				&finishTime,
				&data.BalanceType,
				&data.FromAccount,
			)
			if err != nil {
				ss_log.Error("err=%v", err)
				continue
			}
			if finishTime.String != "" {
				data.FinishTime = finishTime.String
			}
			if paymentType.String != "" {
				data.PaymentType = paymentType.String
			}

			data.OrderType = constants.VaReason_COLLECTION
		}
	} else {
		ss_log.Error("OutgoOrderDao | CustTransferBillsDetail | err=%v\n\nsql=[%v]", err, sqlStr)
		return nil, ss_err.ERR_SYS_DB_GET
	}
	return data, ss_err.ERR_SUCCESS
}
