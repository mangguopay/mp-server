package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type businessBillDao struct{}

var BusinessBillDaoInst businessBillDao

func (*businessBillDao) GetBusinessBillDetail(logNo string) (dataT *go_micro_srv_bill.BusinessBillData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "bb.order_no", Val: logNo, EqType: "="},
	})
	sqlStr := ` select 
					bb.order_no, bb.fee, bb.create_time, bb.amount,	bb.real_amount, 
					bb.order_status, bb.rate, bb.business_no, bb.app_id, bb.account_no, 
					bb.currency_type, bb.pay_time, bb.subject, bb.remark, acc.account, app.app_name, bu.simplify_name
				from business_bill bb
 				left join account acc on acc.uid = bb.business_account_no  
 				left join business bu on bu.business_no = bb.business_no  
 				left join business_app app on app.app_id = bb.app_id  
				` + whereModel.WhereStr

	var orderNo, fee, createTime, amount, realAmount, orderStatus, rate, businessNo, appId,
		accountNo, currencyType, payTime, subject, remark, receiveAccount, appName, simplifyName sql.NullString
	err2 := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{
		&orderNo, &fee, &createTime, &amount, &realAmount, &orderStatus, &rate, &businessNo, &appId,
		&accountNo, &currencyType, &payTime, &subject, &remark, &receiveAccount, &appName, &simplifyName,
	}, whereModel.Args...)

	if err2 != nil {
		return nil, err2
	}

	data := go_micro_srv_bill.BusinessBillData{}
	data.OrderNo = orderNo.String
	data.Fee = fee.String
	data.CreateTime = createTime.String
	data.Amount = amount.String
	data.RealAmount = realAmount.String

	data.OrderStatus = orderStatus.String
	data.Rate = rate.String
	data.BusinessNo = businessNo.String
	data.AppId = appId.String               //
	data.AccountNo = accountNo.String       //收款人账号uid
	data.CurrencyType = currencyType.String //币种
	data.PayTime = payTime.String           //支付时间
	data.Subject = subject.String           //商品名称
	data.Remark = remark.String
	data.ReceiveAccount = receiveAccount.String //收款人账号
	data.AppName = appName.String               //商家app名称
	data.SimplifyName = simplifyName.String     //商家简称

	return &data, nil
}
