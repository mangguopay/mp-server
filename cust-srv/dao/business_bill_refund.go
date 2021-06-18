package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type BusinessBillRefundDao struct {
	RefundNo     string
	OutRefundNo  string
	RefundAmount string
	RefundStatus string
	CreateTime   string
	FinishTime   string
	PayOrderNo   string
	TransAmount  string
	CurrencyType string
	Subject      string
	SceneName    string
	PayeeName    string
	PayeeAcc     string
	AppName      string
	BusinessName string
	BusinessId   string
	BusinessAcc  string
	TradeChannel string
}

var BusinessBillRefundDaoInst BusinessBillRefundDao

func (*BusinessBillRefundDao) GetRefundBills(whereStr string, args []interface{}) ([]*BusinessBillRefundDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT br.refund_no, br.out_refund_no, br.pay_order_no, br.amount, br.refund_status, br.finish_time, " +
		"br.create_time, bb.currency_type, bb.amount, bb.subject, app.app_name, acc.account, acc2.account, b.business_id, b.full_name, " +
		"bc.channel_name, bs.scene_name " +
		"FROM business_refund_order br " +
		"LEFT JOIN business_bill bb ON bb.order_no = br.pay_order_no " +
		"LEFT JOIN business_app app ON app.app_id = bb.app_id " +
		"LEFT JOIN account acc ON acc.uid = bb.account_no " +
		"LEFT JOIN account acc2 ON acc2.uid = bb.business_account_no " +
		"LEFT JOIN business b on b.business_no = bb.business_no " +
		"LEFT JOIN business_scene bs ON bs.scene_no = bb.scene_no " +
		"LEFT JOIN business_channel bc ON bc.business_channel_no = bb.business_channel_no "

	sqlStr += whereStr
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		return nil, err2
	}

	var list []*BusinessBillRefundDao
	for rows.Next() {
		var refundNo, outRefundNo, payOrderNo, refundAmount, refundStatus, finishTime, createTime, currencyType,
			transAmount, subject, appName, payeeAcc, businessAcc, businessId, businessName, channlName, sceneName sql.NullString
		err2 = rows.Scan(&refundNo, &outRefundNo, &payOrderNo, &refundAmount, &refundStatus, &finishTime, &createTime,
			&currencyType, &transAmount, &subject, &appName, &payeeAcc, &businessAcc, &businessId, &businessName, &channlName, &sceneName)
		if err2 != nil {
			return nil, err2
		}
		data := BusinessBillRefundDao{
			RefundNo:     refundNo.String,
			OutRefundNo:  outRefundNo.String,
			PayOrderNo:   payOrderNo.String,
			TransAmount:  transAmount.String,
			RefundAmount: refundAmount.String,
			RefundStatus: refundStatus.String,
			FinishTime:   finishTime.String,
			CreateTime:   createTime.String,
			CurrencyType: currencyType.String,
			Subject:      subject.String,
			AppName:      appName.String,
			PayeeAcc:     payeeAcc.String,
			BusinessName: businessName.String,
			BusinessId:   businessId.String,
			BusinessAcc:  businessAcc.String,
			SceneName:    sceneName.String,
			TradeChannel: channlName.String,
		}
		list = append(list, &data)
	}

	return list, nil
}

func (*BusinessBillRefundDao) GetOrderDetail(refundNo string) (*BusinessBillRefundDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `SELECT br.refund_no, br.out_refund_no, br.amount, br.refund_status, br.create_time, br.finish_time,
		bb.order_no, bb.amount, bb.currency_type, bb.subject, bs.scene_name, 
		acc.nickname, acc.account, app.app_name 
		FROM business_refund_order br 
		LEFT JOIN business_bill bb ON bb.order_no=br.pay_order_no
		LEFT JOIN business_scene bs ON bb.scene_no=bs.scene_no
		LEFT JOIN account acc ON acc.uid=bb.account_no
		LEFT JOIN business_app app ON app.app_id = bb.app_id
		WHERE br.refund_no=$1 `

	var refundNoT, outRefundNo, refundAmount, refundStatus, createTime, finishTime, transOrderNo, transAmount, currencyType,
		subject, sceneName, payeeName, payeeAcc, appName sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&refundNoT, &outRefundNo, &refundAmount, &refundStatus,
		&createTime, &finishTime, &transOrderNo, &transAmount, &currencyType, &subject, &sceneName, &payeeName, &payeeAcc,
		&appName}, refundNo)
	if err != nil {
		return nil, err
	}

	data := new(BusinessBillRefundDao)
	data.RefundNo = refundNoT.String
	data.OutRefundNo = outRefundNo.String
	data.RefundAmount = refundAmount.String
	data.RefundStatus = refundStatus.String
	data.CreateTime = createTime.String
	data.FinishTime = finishTime.String
	data.PayOrderNo = transOrderNo.String
	data.TransAmount = transAmount.String
	data.CurrencyType = currencyType.String
	data.Subject = subject.String
	data.SceneName = sceneName.String
	data.PayeeName = payeeName.String
	data.PayeeAcc = payeeAcc.String
	data.AppName = appName.String

	return data, nil
}

func (*BusinessBillRefundDao) GetCnt(whereList []*model.WhereSqlCond) (total int64, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	sqlCnt := "SELECT COUNT(1) " +
		"FROM business_refund_order br " +
		"LEFT JOIN business_bill bb ON bb.order_no = br.pay_order_no " +
		"LEFT JOIN business b on b.business_no = bb.business_no " +
		"LEFT JOIN business_app app ON app.app_id = bb.app_id " +
		"LEFT JOIN account acc ON acc.uid = bb.account_no " +
		"LEFT JOIN business_scene bs ON bs.scene_no = bb.scene_no " +
		"LEFT JOIN business_channel bc ON bc.business_channel_no = bb.business_channel_no "
	var totalT sql.NullString
	sqlCnt += whereModel.WhereStr
	err = ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereModel.Args...)
	if err != nil {
		return -1, err
	}

	return strext.ToInt64(totalT.String), nil
}

func (*BusinessBillRefundDao) GetSum(whereList []*model.WhereSqlCond) (total int64, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	sqlSum := "SELECT SUM(br.amount) " +
		"FROM business_refund_order br " +
		"LEFT JOIN business_bill bb ON bb.order_no = br.pay_order_no " +
		"LEFT JOIN business b on b.business_no = bb.business_no " +
		"LEFT JOIN business_app app ON app.app_id = bb.app_id " +
		"LEFT JOIN account acc ON acc.uid = bb.account_no " +
		"LEFT JOIN business_scene bs ON bs.scene_no = bb.scene_no " +
		"LEFT JOIN business_channel bc ON bc.business_channel_no = bb.business_channel_no "
	var sumT sql.NullString
	sqlSum += whereModel.WhereStr
	if err := ss_sql.QueryRow(dbHandler, sqlSum, []*sql.NullString{&sumT}, whereModel.Args...); err != nil {
		return -1, err
	}

	return strext.ToInt64(sumT.String), nil
}

func (*BusinessBillRefundDao) GetCntAndSum(whereList []*model.WhereSqlCond) (num, amount string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	sqlSum := "SELECT COUNT(1), SUM(br.amount) " +
		"FROM business_refund_order br " +
		"LEFT JOIN business_bill bb ON bb.order_no = br.pay_order_no "

	var numT, sumT sql.NullString
	sqlSum += whereModel.WhereStr
	if err := ss_sql.QueryRow(dbHandler, sqlSum, []*sql.NullString{&numT, &sumT}, whereModel.Args...); err != nil {
		return "", "", err
	}

	return numT.String, sumT.String, nil

}
