package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/bill-srv/m"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type CollectionDao struct {
}

var CollectionDaoInst CollectionDao

func (CollectionDao) InsertCollectionOrder(tx *sql.Tx, fromVacc, toVacc, amount, balanceType, rees, lat, lng, ip string) (logNo string) {
	logNo = strext.GetDailyId()
	err := ss_sql.ExecTx(tx, `insert into collection_order(log_no,from_vaccount_no,to_vaccount_no,amount,create_time,order_status,balance_type,fees,is_count,lat,lng,ip) values($1,$2,$3,$4,current_timestamp,$5,$6,$7,0,$8,$9,$10)`,
		logNo, fromVacc, toVacc, amount, constants.OrderStatus_Pending, balanceType, rees, lat, lng, ip)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return logNo
}

func (CollectionDao) UpdateCollectionOrderStatus(tx *sql.Tx, logNo, orderStatus string) (errCode string) {
	err := ss_sql.ExecTx(tx, `update collection_order set order_status=$1,finish_time = current_timestamp,modify_time = current_timestamp where log_no=$2`,
		orderStatus, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

func (*CollectionDao) CustCollectionBillsDetail(logNo string) (data *go_micro_srv_bill.CustCollectionBillsDetailData, returnErr string) {
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
	data = &go_micro_srv_bill.CustCollectionBillsDetailData{}
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

// 获取未统计手续费的订单信息
func (*CollectionDao) CollectionFeesTaskResult() (collResults []*m.CommonFeesResult) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	rows, stmt, err := ss_sql.Query(dbHandler, "SELECT  log_no,  balance_type,  fees  FROM collection_order where is_count = '0' and  create_time >  current_timestamp+interval  '-1 hour'")
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err != nil {
		ss_log.Error("err=[查询 ---------> %s 结果失败,err -----> %s]", "CollectionFeesTask", err.Error())
		return nil
	}
	for rows.Next() {
		data := &m.CommonFeesResult{}
		err := rows.Scan(
			&data.LogNo,
			&data.MoneyType,
			&data.Fees,
		)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		data.FeesType = constants.FEES_TYPE_COLLECTION
		collResults = append(collResults, data)
	}
	return collResults
}
