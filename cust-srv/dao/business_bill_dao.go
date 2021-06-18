package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/model"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
	"fmt"
)

type BusinessBillDao struct {
	CreateTime      string
	PayTime         string
	OrderNo         string
	OutOrderNo      string
	Subject         string
	Amount          string
	RealAmount      string
	CurrencyType    string
	Rate            string
	Fee             string
	OrderStatus     string
	Remark          string
	SettleId        string
	SettleDate      string
	Cycle           string
	SceneNo         string
	SceneName       string
	AppId           string
	AppName         string
	BusinessNo      string
	BusinessName    string
	BusinessId      string
	BusinessAccount string
	ChannelName     string
	ChannelRate     string
}

var BusinessBillDaoInst BusinessBillDao

func (*BusinessBillDao) GetCnt(whereList []*model.WhereSqlCond) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlCnt := "SELECT count(1) " +
		"FROM business_bill bb " +
		"LEFT JOIN business_scene bs ON bs.scene_no = bb.scene_no " +
		"LEFT JOIN business_app app ON app.app_id = bb.app_id " +
		"LEFT JOIN business b ON b.business_no = bb.business_no " +
		"LEFT JOIN account acc ON acc.uid = b.account_no " +
		"LEFT JOIN business_channel bc ON bc.business_channel_no = bb.business_channel_no " + whereModel.WhereStr
	var totalT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereModel.Args...); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return "0"
	}

	return totalT.String
}

func (*BusinessBillDao) GetSum(whereList []*model.WhereSqlCond) (sum string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlSum := "SELECT sum(amount) " +
		"FROM business_bill bb " +
		"LEFT JOIN business_scene bs ON bs.scene_no = bb.scene_no " +
		"LEFT JOIN business_app app ON app.app_id = bb.app_id " +
		"LEFT JOIN business b ON b.business_no = bb.business_no " +
		"LEFT JOIN account acc ON acc.uid = b.account_no " +
		"LEFT JOIN business_channel bc ON bc.business_channel_no = bb.business_channel_no " + whereModel.WhereStr
	var sumT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlSum, []*sql.NullString{&sumT}, whereModel.Args...); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return "0"
	}
	if sumT.String == "" {
		return "0"
	}
	return sumT.String
}

func (*BusinessBillDao) GetBusinessBills(whereList []*model.WhereSqlCond, page, pageSize int) ([]*BusinessBillDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY bb.create_time desc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, pageSize, page)
	sqlStr := "SELECT bb.create_time, bb.pay_time, bb.order_no, bb.out_order_no, bb.subject, bb.amount, bb.real_amount, " +
		"bb.currency_type, bb.fee, bb.order_status, bb.settle_id, bs.scene_no, bs.scene_name, app.app_id, app.app_name," +
		" b.business_id, b.full_name, acc.account, bb.cycle, bb.rate, bb.settle_date " +
		"FROM business_bill bb " +
		"LEFT JOIN business_scene bs ON bs.scene_no = bb.scene_no " +
		"LEFT JOIN business_app app ON app.app_id = bb.app_id " +
		"LEFT JOIN business b ON b.business_no = bb.business_no " +
		"LEFT JOIN account acc ON acc.uid = b.account_no " +
		"LEFT JOIN business_channel bc ON bc.business_channel_no = bb.business_channel_no " + whereModel.WhereStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var datas []*BusinessBillDao
	for rows.Next() {
		var createTime, payTime, orderNo, outOrderNo, subject, amount, realAmount, currencyType, fee,
			orderStatus, settleId, sceneNo, sceneName, appId, appName, businessId, businessName,
			businessAcc, cycle, rate, settleDate sql.NullString
		err = rows.Scan(&createTime, &payTime, &orderNo, &outOrderNo, &subject, &amount, &realAmount, &currencyType,
			&fee, &orderStatus, &settleId, &sceneNo, &sceneName, &appId, &appName, &businessId, &businessName, &businessAcc,
			&cycle, &rate, &settleDate,
		)
		if err != nil {
			return nil, err
		}
		data := &BusinessBillDao{
			CreateTime:      createTime.String,
			PayTime:         payTime.String,
			OrderNo:         orderNo.String,
			OutOrderNo:      outOrderNo.String,
			Subject:         subject.String,
			Amount:          amount.String,
			RealAmount:      realAmount.String,
			CurrencyType:    currencyType.String,
			Fee:             fee.String,
			OrderStatus:     orderStatus.String,
			SettleId:        settleId.String,
			SceneNo:         sceneNo.String,
			SceneName:       sceneName.String,
			AppId:           appId.String,
			AppName:         appName.String,
			BusinessName:    businessName.String,
			BusinessId:      businessId.String,
			BusinessAccount: businessAcc.String,
			Cycle:           cycle.String,
			Rate:            rate.String,
			SettleDate:      settleDate.String,
		}
		datas = append(datas, data)
	}

	return datas, nil
}

func (*BusinessBillDao) GetChannelBillsCnt(whereList []*model.WhereSqlCond) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlCnt := "SELECT count(1) " +
		"FROM business_bill bb " +
		"LEFT JOIN business_scene bs ON bs.scene_no = bb.scene_no " +
		"LEFT JOIN business_scene_signed bss ON bss.scene_no = bs.scene_no AND bss.status = 3 " +
		"LEFT JOIN business_app app ON app.app_id = bb.app_id " +
		"LEFT JOIN business_industry_rate_cycle bi ON bi.code = bss.industry_no AND bi.business_channel_no =  bb.business_channel_no AND bi.is_delete = 0" +
		"LEFT JOIN business b ON b.business_no = bb.business_no " +
		"LEFT JOIN account acc ON acc.uid = b.account_no " +
		"LEFT JOIN business_channel bc ON bc.business_channel_no = bb.business_channel_no "

	sqlCnt += whereModel.WhereStr
	var totalT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereModel.Args...); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return "0"
	}

	return totalT.String
}

func (*BusinessBillDao) GetBusinessChannelBills(whereList []*model.WhereSqlCond, page, pageSize int) ([]*BusinessBillDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY bb.create_time desc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, pageSize, page)
	sqlStr := "SELECT bb.create_time, bb.pay_time, bb.order_no, bb.out_order_no, bb.subject, bb.amount, bb.real_amount, " +
		"bb.currency_type, bb.fee, bb.order_status, bb.settle_id, bs.scene_no, bs.scene_name, app.app_id, app.app_name, " +
		"b.business_id, b.full_name, acc.account, bb.cycle, bb.rate, bb.settle_date, " +
		"bc.channel_name, " +
		"bi.rate channel_rate " +
		"FROM business_bill bb " +
		"LEFT JOIN business_scene bs ON bs.scene_no = bb.scene_no " +
		"LEFT JOIN business_scene_signed bss ON bss.scene_no = bs.scene_no AND bss.status = 3 " +
		"LEFT JOIN business_app app ON app.app_id = bb.app_id " +
		"LEFT JOIN business_industry_rate_cycle bi ON bi.code = bss.industry_no AND bi.business_channel_no =  bb.business_channel_no AND bi.is_delete = 0" +
		"LEFT JOIN business b ON b.business_no = bb.business_no " +
		"LEFT JOIN account acc ON acc.uid = b.account_no " +
		"LEFT JOIN business_channel bc ON bc.business_channel_no = bb.business_channel_no "

	sqlStr += whereModel.WhereStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var datas []*BusinessBillDao
	for rows.Next() {
		var createTime, payTime, orderNo, outOrderNo, subject, amount, realAmount, currencyType, fee,
			orderStatus, settleId, sceneNo, sceneName, appId, appName, businessId, businessName,
			businessAcc, cycle, rate, settleDate, channelName, channelRate sql.NullString
		err = rows.Scan(&createTime, &payTime, &orderNo, &outOrderNo, &subject, &amount, &realAmount, &currencyType,
			&fee, &orderStatus, &settleId, &sceneNo, &sceneName, &appId, &appName, &businessId, &businessName, &businessAcc,
			&cycle, &rate, &settleDate, &channelName, &channelRate,
		)
		if err != nil {
			return nil, err
		}
		data := &BusinessBillDao{
			CreateTime:      createTime.String,
			PayTime:         payTime.String,
			OrderNo:         orderNo.String,
			OutOrderNo:      outOrderNo.String,
			Subject:         subject.String,
			Amount:          amount.String,
			RealAmount:      realAmount.String,
			CurrencyType:    currencyType.String,
			Fee:             fee.String,
			OrderStatus:     orderStatus.String,
			SettleId:        settleId.String,
			SceneNo:         sceneNo.String,
			SceneName:       sceneName.String,
			AppId:           appId.String,
			AppName:         appName.String,
			BusinessName:    businessName.String,
			BusinessId:      businessId.String,
			BusinessAccount: businessAcc.String,
			Cycle:           cycle.String,
			Rate:            rate.String,
			SettleDate:      settleDate.String,
			ChannelName:     channelName.String,
			ChannelRate:     channelRate.String,
		}
		datas = append(datas, data)
	}

	return datas, nil
}

func (*BusinessBillDao) GetBusinessBillDetail(orderNo string) (*custProto.BusinessBillData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT bb.create_time, bb.order_no, bb.out_order_no, bb.subject, bb.amount, bb.real_amount, bb.currency_type, bb.fee," +
		"bb.order_status, bb.settle_id, bs.scene_name, bb.rate, bb.cycle, bb.settle_date, bb.pay_time, app.app_name, acc.account " +
		"FROM business_bill bb " +
		"LEFT JOIN business_scene bs ON bs.scene_no = bb.scene_no " +
		"LEFT JOIN business_app app ON app.app_id = bb.app_id " +
		"LEFT JOIN account acc ON acc.uid = bb.account_no " +
		"WHERE bb.order_no=$1 "

	var createTime, orderNoT, outOrderNo, subject, amount, realAmount, currencyType, fee, orderStatus,
		settleId, sceneName, rate, cycle, settleDate, payTime, appName, account sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{
		&createTime, &orderNoT, &outOrderNo, &subject, &amount,
		&realAmount, &currencyType, &fee, &orderStatus, &settleId,
		&sceneName, &rate, &cycle, &settleDate, &payTime,
		&appName, &account,
	}, orderNo)
	if err != nil {
		return nil, err
	}

	data := &custProto.BusinessBillData{
		CreateTime:   createTime.String,
		OrderNo:      orderNoT.String,
		OutOrderNo:   outOrderNo.String,
		Subject:      subject.String,
		Amount:       amount.String,
		RealAmount:   realAmount.String,
		CurrencyType: currencyType.String,
		Fee:          fee.String,
		OrderStatus:  orderStatus.String,
		SettleId:     settleId.String,
		SceneName:    sceneName.String,
		Rate:         rate.String,
		Cycle:        cycle.String,
		SettleDate:   settleDate.String,
		PayTime:      payTime.String,
		AppName:      appName.String,
		Account:      account.String,
	}

	return data, nil
}

func (*BusinessBillDao) GetTodayPaySum(businessNo, currencyType string) (totalNum, totalAmount string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	startTime := fmt.Sprintf("%s %s", ss_time.Now(global.Tz).Format(ss_time.DateFormat), "00:00:00")
	endTime := fmt.Sprintf("%s %s", ss_time.Now(global.Tz).Format(ss_time.DateFormat), "23:59:59")

	sqlSum := "SELECT count(1), sum(amount) " +
		"FROM business_bill " +
		"WHERE business_no=$1 AND currency_type = $2 AND pay_time >= $3 AND pay_time <= $4 AND order_status IN ($5, $6, $7) "
	var num, amountSum sql.NullString
	err = ss_sql.QueryRow(dbHandler, sqlSum, []*sql.NullString{&num, &amountSum}, businessNo, currencyType, startTime, endTime,
		constants.BusinessOrderStatusPay, constants.BusinessOrderStatusRefund, constants.BusinessOrderStatusRebatesRefund)
	if err != nil {
		return "", "", err
	}

	if amountSum.String == "" || num.String == "" {
		return "0", "0", nil
	}
	return num.String, amountSum.String, nil
}

func (*BusinessBillDao) GetPersonalBusinessBills(whereStr string, args []interface{}) ([]*BusinessBillDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT bb.create_time, bb.pay_time, bb.order_no, bb.out_order_no, bb.subject, bb.amount, bb.real_amount, " +
		"bb.currency_type, bb.fee, bb.order_status, bb.settle_id, bb.cycle, bb.rate, bb.settle_date, bb.remark " +
		"FROM business_bill bb " +
		"LEFT JOIN business b ON b.business_no = bb.business_no "
	sqlStr += whereStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var list []*BusinessBillDao
	for rows.Next() {
		var createTime, payTime, orderNo, outOrderNo, subject, amount, realAmount, currencyType, fee,
			orderStatus, settleId, cycle, rate, settleDate, remark sql.NullString
		err = rows.Scan(&createTime, &payTime, &orderNo, &outOrderNo, &subject, &amount, &realAmount, &currencyType,
			&fee, &orderStatus, &settleId, &cycle, &rate, &settleDate, &remark)
		if err != nil {
			return nil, err
		}
		data := &BusinessBillDao{
			CreateTime:   createTime.String,
			PayTime:      payTime.String,
			OrderNo:      orderNo.String,
			OutOrderNo:   outOrderNo.String,
			Subject:      subject.String,
			Amount:       amount.String,
			RealAmount:   realAmount.String,
			CurrencyType: currencyType.String,
			Fee:          fee.String,
			OrderStatus:  orderStatus.String,
			SettleId:     settleId.String,
			Cycle:        cycle.String,
			Rate:         rate.String,
			SettleDate:   settleDate.String,
			Remark:       remark.String,
		}
		list = append(list, data)
	}

	return list, nil
}

func (*BusinessBillDao) GetPersonalBusinessBillDetail(orderNo, businessAccNo string) (*BusinessBillDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT bb.subject, bb.amount, bb.real_amount, bb.currency_type, bb.fee, bb.order_status, bb.settle_id, " +
		"bb.settle_date, bb.pay_time, bs.scene_name " +
		"FROM business_bill bb " +
		"LEFT JOIN business_scene bs ON bs.scene_no = bb.scene_no " +
		"WHERE bb.order_no = $1 AND bb.business_account_no = $2 "

	var subject, amount, realAmount, currencyType, fee, orderStatus, settleId, settleDate, payTime, sceneName sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&subject, &amount, &realAmount, &currencyType, &fee,
		&orderStatus, &settleId, &settleDate, &payTime, &sceneName}, orderNo, businessAccNo)
	if err != nil {
		return nil, err
	}

	data := &BusinessBillDao{
		OrderNo:      orderNo,
		Subject:      subject.String,
		Amount:       amount.String,
		RealAmount:   realAmount.String,
		CurrencyType: currencyType.String,
		Fee:          fee.String,
		OrderStatus:  orderStatus.String,
		SettleId:     settleId.String,
		SettleDate:   settleDate.String,
		PayTime:      payTime.String,
		SceneName:    sceneName.String,
	}

	return data, nil
}

func (*BusinessBillDao) GetPersonalBusinessBillsCnt(whereStr string, args []interface{}) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlCnt := "SELECT COUNT(1) FROM business_bill bb " +
		"LEFT JOIN business b ON b.business_no = bb.business_no "
	sqlCnt += whereStr
	var total sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args...); err != nil {
		return "", err
	}

	return total.String, nil
}
