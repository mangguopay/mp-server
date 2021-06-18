package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type LogToBusinessDao struct {
}

var LogToBusinessDaoInst LogToBusinessDao

func (LogToBusinessDao) GetToBusinessCnt(whereList []*model.WhereSqlCond) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	var totalT sql.NullString
	sqlStr := "SELECT count(1) " +
		" FROM log_to_business tb " +
		" left join card_business cb on cb.card_no = tb.card_no	 " +
		" left join business b on b.business_no = tb.business_no " +
		" left join account acc on acc.uid = b.account_no " + whereModel.WhereStr
	if err2 := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, whereModel.Args...); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return ""
	}

	return totalT.String
}

func (LogToBusinessDao) GetToBusinessList(whereList []*model.WhereSqlCond, page, pageSize string) (datasT []*go_micro_srv_cust.ToBusinessData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY tb.create_time desc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(pageSize), strext.ToInt(page))
	sqlStr := `SELECT 
			tb.log_no, tb.currency_type, tb.business_no, tb.collection_type, tb.card_no,
			tb.amount, tb.create_time, tb.order_type, tb.order_status, tb.finish_time,
			tb.fee, tb.image_id, tb.notes, tb.real_amount, cb.card_number,
			cb.name, acc.account, ch.channel_name
		 FROM log_to_business tb 
		 left join card_business cb on cb.card_no = tb.card_no	 
		 left join channel ch on ch.channel_no = cb.channel_no	 
		 left join business b on b.business_no = tb.business_no 
		 left join account acc on acc.uid = b.account_no ` + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	var datas []*go_micro_srv_cust.ToBusinessData
	for rows.Next() {
		data := go_micro_srv_cust.ToBusinessData{}
		var imageId, cardNumber, name, account, channelName sql.NullString
		err2 = rows.Scan(
			&data.LogNo, &data.CurrencyType, &data.BusinessNo, &data.CollectionType, &data.CardNo,
			&data.Amount, &data.CreateTime, &data.OrderType, &data.OrderStatus, &data.FinishTime,
			&data.Fee, &imageId, &data.Notes, &data.RealAmount, &cardNumber,
			&name, &account, &channelName,
		)

		if err2 != nil {
			ss_log.Error("err=[%v]", err2)
			return nil, err2
		}

		data.ImageId = imageId.String
		data.CardNumber = cardNumber.String
		data.Name = name.String
		data.Account = account.String
		data.ChannelName = channelName.String
		datas = append(datas, &data)
	}

	return datas, nil
}

//查询商家提现详情
func (LogToBusinessDao) GetToBusinessDetail(logNo string) (*go_micro_srv_cust.ToBusinessData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `SELECT 
			cb.name, cb.card_number, acc.account,
			tb.create_time, tb.log_no, tb.fee, tb.currency_type, tb.amount,
			tb.collection_type, tb.real_amount, tb.order_status, tb.image_id
		 FROM log_to_business tb 
		 left join card_business cb on cb.card_no = tb.card_no	 
		 left join business b on b.business_no = tb.business_no 
		 left join account acc on acc.uid = b.account_no
		 left join dict_images di on di.image_id = tb.image_id
         WHERE tb.log_no=$1 `

	var payeeName, cardNumber, businessAccNo, createTime, orderNo, fee, currencyType, amount,
		collectionType, realAmount, orderStatus, imageId sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr,
		[]*sql.NullString{&payeeName, &cardNumber, &businessAccNo, &createTime, &orderNo, &fee, &currencyType, &amount,
			&collectionType, &realAmount, &orderStatus, &imageId}, logNo)
	if err != nil {
		return nil, err
	}

	data := new(go_micro_srv_cust.ToBusinessData)
	data.Name = payeeName.String
	data.CardNumber = cardNumber.String
	data.Account = businessAccNo.String
	data.CreateTime = createTime.String
	data.LogNo = orderNo.String
	data.Fee = fee.String
	data.CurrencyType = currencyType.String
	data.Amount = amount.String
	data.CollectionType = collectionType.String
	data.RealAmount = realAmount.String
	data.OrderStatus = orderStatus.String
	data.ImageId = imageId.String

	return data, nil
}
