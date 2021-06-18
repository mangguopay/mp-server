package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type LogBusinessToHeadquartersDao struct{}

var LogBusinessToHeadquartersDaoInst LogBusinessToHeadquartersDao

func (*LogBusinessToHeadquartersDao) GetBusinessToHeadDetail(whereList []*model.WhereSqlCond) (dataT *go_micro_srv_bill.BusinessToHeadData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	sqlStr := ` select 
					bth.log_no, bth.business_no, bth.currency_type, bth.amount, bth.order_status,
					bth.collection_type, bth.card_no, bth.create_time, bth.modify_time, bth.image_id,
					bth.arrive_amount, bth.fee, bth.notes, acc.account, bu.account_no
				from log_business_to_headquarters bth
 				left join business bu on bu.business_no = bth.business_no 
				left join account acc on acc.uid = bu.account_no ` + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err2 != nil {
		return nil, err2
	}

	data := go_micro_srv_bill.BusinessToHeadData{}
	var (
		logNo, businessNo, currencyType, amount, orderStatus,
		collectionType, cardNo, createTime, modifyTime, imageId,
		arriveAmount, fee, notes, account, businessAccNo sql.NullString
	)
	err2 = rows.Scan(
		&logNo,
		&businessNo,
		&currencyType,
		&amount,
		&orderStatus,

		&collectionType,
		&cardNo,
		&createTime,
		&modifyTime,
		&imageId,

		&arriveAmount,
		&fee,
		&notes,
		&account,
		&businessAccNo,
	)

	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return nil, err2
	}

	data.LogNo = logNo.String
	data.BusinessNo = businessNo.String
	data.CurrencyType = currencyType.String
	data.Amount = amount.String
	data.OrderStatus = orderStatus.String

	data.CollectionType = collectionType.String
	data.CardNo = cardNo.String
	data.CreateTime = createTime.String
	data.ModifyTime = modifyTime.String
	data.ImageId = imageId.String

	data.ArriveAmount = arriveAmount.String
	data.Fee = fee.String
	data.Notes = notes.String
	data.Account = account.String
	data.BusinessAccNo = businessAccNo.String

	return &data, nil
}

func (*LogBusinessToHeadquartersDao) ModifyStatus(tx *sql.Tx, logNo, status string) error {
	sqlStr := "update log_business_to_headquarters set order_status = $3 where log_no = $1 and order_status = $2 "
	if err := ss_sql.ExecTx(tx, sqlStr, logNo, constants.AuditOrderStatus_Pending, status); err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}
