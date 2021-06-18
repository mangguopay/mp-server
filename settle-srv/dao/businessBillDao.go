package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
	"errors"
)

type BusinessBillDao struct {
	OrderNo           string
	OrderStatus       string
	Amount            string
	RealAmount        string
	CurrencyType      string
	Fee               string
	BusinessNo        string
	BusinessAccountNo string
}

var BusinessBillDaoInst BusinessBillDao

func (BusinessBillDao) GetOrderInfoByOrderNo(orderNo string) (*BusinessBillDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, errors.New("获取数据库连接失败")
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var (
		orderNoT, orderStatus, amount, realAmount, currencyType, businessNo, fee, businessAccNo sql.NullString
	)
	err := ss_sql.QueryRow(dbHandler, "SELECT order_no, order_status, amount, real_amount, currency_type, "+
		"business_no, fee, business_account_no "+
		"FROM business_bill WHERE order_no= $1 limit 1 ",
		[]*sql.NullString{&orderNoT, &orderStatus, &amount, &realAmount, &currencyType, &businessNo, &fee, &businessAccNo}, orderNo)
	if err != nil {
		return nil, err
	}

	obj := new(BusinessBillDao)
	obj.OrderNo = orderNo
	obj.OrderStatus = orderStatus.String
	obj.Amount = amount.String
	obj.RealAmount = realAmount.String
	obj.CurrencyType = currencyType.String
	obj.BusinessNo = businessNo.String
	obj.Fee = fee.String
	obj.BusinessAccountNo = businessAccNo.String
	return obj, nil
}

type BusinessBillSettleData struct {
	BusinessNo      string
	BusinessAccNo   string
	AppId           string
	CurrencyType    string
	TotalAmount     int64
	TotalRealAmount int64
	TotalFees       int64
	TotalOrder      int64
}

// 获取商户App的结算金额
func (b *BusinessBillDao) GetSettleData(startTime, endTime int64) ([]*BusinessBillSettleData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `SELECT business_no, business_account_no, app_id, currency_type,SUM(amount) AS t_amount,
		SUM(real_amount) AS t_real_amount, SUM(fee) AS t_fee,COUNT(order_no) AS t_order 
		FROM business_bill 
		WHERE settle_date >= $1 AND settle_date <= $2 AND order_status = $3 AND settle_id = '' 
		GROUP BY business_no, business_account_no, app_id, currency_type `
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, startTime, endTime, constants.BusinessOrderStatusPay)
	if nil != err {
		return nil, err
	}

	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var dataList []*BusinessBillSettleData
	for rows.Next() {
		var businessNo, businessAccNo, appId, currencyType, tAmount, tRealAmount, tFee, tOrder sql.NullString
		err := rows.Scan(&businessNo, &businessAccNo, &appId, &currencyType, &tAmount, &tRealAmount, &tFee, &tOrder)
		if err != nil {
			return nil, err
		}
		ret := new(BusinessBillSettleData)
		ret.BusinessNo = businessNo.String
		ret.BusinessAccNo = businessAccNo.String
		ret.AppId = appId.String
		ret.CurrencyType = currencyType.String
		ret.TotalAmount = strext.ToInt64(tAmount.String)
		ret.TotalRealAmount = strext.ToInt64(tRealAmount.String)
		ret.TotalFees = strext.ToInt64(tFee.String)
		ret.TotalOrder = strext.ToInt64(tOrder.String)

		dataList = append(dataList, ret)
	}

	return dataList, nil
}

func (b *BusinessBillDao) UpdateOrderSettleIdTx(tx *sql.Tx, settleId, orderNo string) error {
	sqlStr := "UPDATE business_bill SET settle_id=$1 WHERE order_no=$2 AND settle_id='' "
	return ss_sql.ExecTx(tx, sqlStr, settleId, orderNo)

}
