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

type LogBusinessToHeadquartersDao struct {
	LogNo          string //商家转账至总部流水id
	BusinessNo     string //商家身份id
	CurrencyType   string //币种
	Amount         string //金额
	OrderStatus    string //订单状态(0待审核、1审核通过，2审核不通过)
	CollectionType string //收款方式,1-支票;2-现金;3-银行转账;4-其他
	CardNo         string //总部收款卡uid
	ImageId        string //凭证 图片id
	ArriveAmount   string //实际到账金额
	Fee            string //手续费
}

var LogBusinessToHeadquartersDaoInst LogBusinessToHeadquartersDao

func (LogBusinessToHeadquartersDao) Insert(data LogBusinessToHeadquartersDao) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	logNo := strext.GetDailyId()
	sqlStr := `insert into log_business_to_headquarters(log_no, business_no, currency_type, amount, order_status, collection_type, card_no, image_id, arrive_amount, fee, create_time, modify_time) 
				values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,current_timestamp,current_timestamp)`
	if err := ss_sql.Exec(dbHandler, sqlStr, logNo, data.BusinessNo, data.CurrencyType, data.Amount, constants.AuditOrderStatus_Pending,
		data.CollectionType, data.CardNo, data.ImageId, data.ArriveAmount, data.Fee); err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return logNo
}

func (LogBusinessToHeadquartersDao) GetCnt(whereList []*model.WhereSqlCond) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	var total sql.NullString
	sqlStr := ` select count(1) 
				from log_business_to_headquarters bth
 				left join business bu on bu.business_no = bth.business_no 
				left join account acc on acc.uid = bu.account_no ` + whereModel.WhereStr
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&total}, whereModel.Args...); err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return total.String
}

func (*LogBusinessToHeadquartersDao) GetBusinessToHeadList(whereList []*model.WhereSqlCond, page, pageSize int) (datas []*go_micro_srv_cust.BusinessToHeadData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY bth.create_time desc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(pageSize), strext.ToInt(page))
	sqlStr := ` select 
					bth.log_no, bth.business_no, bth.currency_type, bth.amount, bth.order_status,
					bth.collection_type, bth.card_no, bth.create_time, bth.modify_time, bth.image_id,
					bth.arrive_amount, bth.fee, bth.notes, acc.account, cah.card_number,
					cah.name, ch.channel_name, bu.business_type
				from log_business_to_headquarters bth
 				left join card_head cah on cah.card_no = bth.card_no  
 				left join channel_business_config cbc on cbc.id = cah.channel_business_config_id  
 				left join channel ch on ch.channel_no = cbc.channel_no  
 				left join business bu on bu.business_no = bth.business_no 
				left join account acc on acc.uid = bu.account_no ` + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		return nil, err2
	}

	var datasT []*go_micro_srv_cust.BusinessToHeadData
	for rows.Next() {
		data := go_micro_srv_cust.BusinessToHeadData{}
		var (
			logNo, businessNo, currencyType, amount, orderStatus,
			collectionType, cardNo, createTime, modifyTime, imageId,
			arriveAmount, fee, notes, account, cardNumber,
			name, channelName, businessType sql.NullString
		)

		err2 = rows.Scan(
			&logNo, &businessNo, &currencyType, &amount, &orderStatus,
			&collectionType, &cardNo, &createTime, &modifyTime, &imageId,
			&arriveAmount, &fee, &notes, &account, &cardNumber,
			&name, &channelName, &businessType,
		)

		if err2 != nil {
			ss_log.Error("LogNo[%v],err=[%v]", data.LogNo, err2)
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
		data.CardNumber = cardNumber.String

		data.Name = name.String
		data.ChannelName = channelName.String
		data.BusinessType = businessType.String

		datasT = append(datasT, &data)
	}

	return datasT, nil
}

func (*LogBusinessToHeadquartersDao) GetBusinessToHeadDetail(whereList []*model.WhereSqlCond) (data *go_micro_srv_cust.BusinessToHeadData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	sqlStr := ` select 
					bth.log_no, bth.business_no, bth.currency_type, bth.amount, bth.order_status,
					bth.collection_type, bth.card_no, bth.create_time, bth.modify_time, bth.image_id,
					bth.arrive_amount, bth.fee, bth.notes, acc.account, cah.card_number,
					cah.name, bu.business_type
				from log_business_to_headquarters bth
 				left join card_head cah on cah.card_no = bth.card_no  
 				left join business bu on bu.business_no = bth.business_no 
				left join account acc on acc.uid = bu.account_no ` + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err2 != nil {
		return nil, err2
	}

	dataT := &go_micro_srv_cust.BusinessToHeadData{}
	var (
		logNo, businessNo, currencyType, amount, orderStatus,
		collectionType, cardNo, createTime, modifyTime, imageId,
		arriveAmount, fee, notes, account, cardNumber,
		name, businessType sql.NullString
	)

	err2 = rows.Scan(
		&logNo, &businessNo, &currencyType, &amount, &orderStatus,
		&collectionType, &cardNo, &createTime, &modifyTime, &imageId,
		&arriveAmount, &fee, &notes, &account, &cardNumber,
		&name, &businessType,
	)

	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return nil, err2
	}

	dataT.LogNo = logNo.String
	dataT.BusinessNo = businessNo.String
	dataT.CurrencyType = currencyType.String
	dataT.Amount = amount.String
	dataT.OrderStatus = orderStatus.String

	dataT.CollectionType = collectionType.String
	dataT.CardNo = cardNo.String
	dataT.CreateTime = createTime.String
	dataT.ModifyTime = modifyTime.String
	dataT.ImageId = imageId.String

	dataT.ArriveAmount = arriveAmount.String
	dataT.Fee = fee.String
	dataT.Notes = notes.String
	dataT.Account = account.String
	dataT.CardNumber = cardNumber.String

	dataT.Name = name.String
	dataT.BusinessType = businessType.String

	return dataT, nil
}

/*
func (*LogBusinessToHeadquartersDao) QueryOrderStatusFromLogNo(orderNo string) (string, string, string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var statusT, currencyType, custNo, arriveAmount, fees, amount sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select order_status,currency_type,cust_no,arrive_amount,fees,amount from log_cust_to_headquarters where log_no=$1  limit 1`, []*sql.NullString{&statusT, &currencyType, &custNo, &arriveAmount, &fees, &amount}, orderNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", "", "", "", "", ""
	}
	return statusT.String, currencyType.String, custNo.String, arriveAmount.String, fees.String, amount.String
}

func (*LogBusinessToHeadquartersDao) UpdateStatusFromLogNo(orderNo string, status int32) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, `update log_cust_to_headquarters set order_status= $1 where log_no=$2`, status, orderNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_OPERATE_FAILD
	}
	return ss_err.ERR_SUCCESS
}

func (*LogBusinessToHeadquartersDao) CustToHeadquartersDetail(logNo string) (data *go_micro_srv_bill.CustToHeadquartersDetailData, returnErr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT log_no, amount, order_status, create_time, finish_time, currency_type,  fees,arrive_amount,payment_type  FROM  " +
		"log_cust_to_headquarters  WHERE  log_no = $1"
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, logNo)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	data = &go_micro_srv_bill.CustToHeadquartersDetailData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, ss_err.ERR_SYS_DB_GET
	}
	for rows.Next() {
		var finishTime sql.NullString
		err = rows.Scan(
			&data.LogNo,
			&data.Amount,
			&data.OrderStatus,
			&data.CreateTime,
			&finishTime,
			&data.BalanceType,
			&data.Fees,
			&data.ArriveAmount,
			&data.PaymentType,
		)
		if err != nil {
			ss_log.Error("err=%v", err)
			continue
		}
		if finishTime.String != "" {
			data.FinishTime = finishTime.String
		}
		data.OrderType = constants.VaReason_Cust_Save
	}

	return data, ss_err.ERR_SUCCESS
}

func (*LogBusinessToHeadquartersDao) CustToHeadquartersList(whereStr string, whereArgs []interface{}) (datas []*go_micro_srv_cust.CustToHeadquartersData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT lcth.log_no, lcth.currency_type, lcth.collection_type" +
		", lcth.amount, lcth.create_time, lcth.order_type, lcth.order_status, lcth.finish_time" +
		",lcth.lat, lcth.lng, lcth.fees, lcth.ip, lcth.image_id " +
		", acc.account, ca.name, ca.card_number " +
		" FROM log_cust_to_headquarters lcth " +
		" LEFT JOIN cust cu ON cu.cust_no = lcth.cust_no " +
		" LEFT JOIN account acc ON acc.uid = cu.account_no " +
		" LEFT JOIN card_head ca ON ca.card_no = lcth.card_no " + whereStr

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	var datasT []*go_micro_srv_cust.CustToHeadquartersData
	for rows.Next() {
		var data go_micro_srv_cust.CustToHeadquartersData
		var finishTime, imageId sql.NullString
		if err = rows.Scan(
			&data.LogNo,
			&data.CurrencyType,
			&data.CollectionType,
			&data.Amount,
			&data.CreateTime,

			&data.OrderType,
			&data.OrderStatus,
			&finishTime,
			&data.Lat,
			&data.Lng,

			&data.Fees,
			&data.Ip,
			&imageId,
			&data.Account,
			&data.Name,

			&data.CardNumber,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		if finishTime.String != "" {
			data.FinishTime = finishTime.String
		}
		if imageId.String != "" {
			data.ImageId = imageId.String
		}

		datasT = append(datasT, &data)
	}

	return datasT, nil

}

// GetCustToHeadquartersCountByDate 统计
func (*LogBusinessToHeadquartersDao) GetCustToHeadquartersCountByDate(beforeDay, currencyType string) (*DataCount, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	startTime := beforeDay
	endTime, tErr := ss_time.TimeAfter(beforeDay, ss_time.DateFormat, time.Hour*24) // 下一天的日期
	if tErr != nil {
		return nil, tErr
	}

	sqlStr := `select count(1) as totalCount,sum(amount) as totalAmount,sum(fees) as totalFees from log_cust_to_headquarters WHERE create_time >= $1
		and create_time < $2 AND currency_type = $3 and order_status = 1`

	var num, amount, fee sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&num, &amount, &fee}, startTime, endTime, currencyType)

	if errT != nil {
		return nil, errT
	}
	data := &DataCount{
		Num:    strext.ToInt64(num.String),
		Amount: strext.ToInt64(amount.String),
		Fee:    strext.ToInt64(fee.String),
		Type:   1, // 向总部充值
		Day:    startTime,
		CType:  currencyType,
	}
	return data, nil
}
*/
