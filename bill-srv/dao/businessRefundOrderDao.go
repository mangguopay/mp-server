package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type BusinessRefundOrderDao struct {
	RefundNo     string
	Amount       string
	RefundStatus string
	PayOrderNo   string
	Remark       string
	CreateTime   string
	OutRefundNo  string
	FinishTime   string
	PayeeAmount  string
	CurrencyType string
	Subject      string
}

var BusinessRefundOrderDaoInst BusinessRefundOrderDao

func (BusinessRefundOrderDao) GetRefundOrderDetail(whereList []*model.WhereSqlCond) (datas *BusinessRefundOrderDao, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlStr := `SELECT la.amount, bro.refund_no, bro.amount, bro.refund_status, bro.pay_order_no, bro.remark,
			bro.finish_time, bb.currency_type, bb.subject 
		 FROM log_vaccount la
		 LEFT JOIN vaccount va ON va.vaccount_no = la.vaccount_no 
		 LEFT JOIN business_refund_order bro ON bro.refund_no = la.biz_log_no 
         LEFT JOIN business_bill bb ON bb.order_no = bro.pay_order_no  ` + whereModel.WhereStr
	row, stmt, errT := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if errT != nil {
		ss_log.Error("err[%v]", errT)
		return nil, errT
	}

	var payeeAmount, refundNo, amount, refundStatus, payOrderNo, remark, finishTime, currencyType, subject sql.NullString
	errT = row.Scan(
		&payeeAmount, &refundNo, &amount, &refundStatus, &payOrderNo,
		&remark, &finishTime, &currencyType, &subject,
	)

	if errT != nil {
		ss_log.Error("err[%v]", errT)
		return nil, errT
	}

	data := &BusinessRefundOrderDao{
		RefundNo:     refundNo.String,
		Amount:       amount.String,
		RefundStatus: refundStatus.String,
		PayOrderNo:   payOrderNo.String,
		Remark:       remark.String,

		FinishTime:   finishTime.String,
		CurrencyType: currencyType.String,
		Subject:      subject.String,
		PayeeAmount:  payeeAmount.String,
	}
	return data, nil
}
