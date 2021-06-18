package dao

import (
	"database/sql"
	"fmt"
	"time"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type ServiceDao struct {
}

var ServiceDaoInst ServiceDao

func (ServiceDao) AddService(tx *sql.Tx, accountNo, addr string) (error, string) {
	//创建运营商
	servicer_no := strext.NewUUID()
	sqlStr := "insert into servicer(servicer_no, account_no, addr, create_time) " +
		" values ($1,$2,$3, current_timestamp)"
	err := ss_sql.ExecTx(tx, sqlStr, servicer_no, accountNo, addr)
	if err != nil {
		ss_log.Error("ServiceDao |AddService err=[%v]", err)
	}
	return err, servicer_no
}

func (ServiceDao) GetServiceByCashierNo(dbHandler *sql.DB, cashier_no string) (error, string) {
	sqlStr := "select servicer_no from cashier where uid = $1"
	var servicerNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&servicerNo}, cashier_no)
	if err != nil {
		ss_log.Error("ServiceDao |GetServiceByCashierNo err=[%v]", err)
	}
	return err, servicerNo.String
}

func (ServiceDao) GetServicerBillingDetails(accountNo, accoutType, currencyType string, page, pageSize int32) (returnDatas []*go_micro_srv_cust.ServicerBillingDetailsData, returnTotal, returnErr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var datas []*go_micro_srv_cust.ServicerBillingDetailsData

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "account_no", Val: accountNo, EqType: "="},
		{Key: "account_type", Val: accoutType, EqType: "="},
		{Key: "currency_type", Val: currencyType, EqType: "="},
	})
	//统计
	where := whereModel.WhereStr
	args := whereModel.Args
	var total sql.NullString
	sqlCnt := "SELECT count(1) " +
		"FROM billing_details_results " + where
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}
	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(pageSize), strext.ToInt(page))

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT create_time, bill_no, amount, currency_type, bill_type " +
		" FROM billing_details_results " + where2

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err == nil {
		for rows.Next() {
			var data go_micro_srv_cust.ServicerBillingDetailsData
			err = rows.Scan(
				&data.CreateTime,
				&data.BillNo,
				&data.Amount,
				&data.CurrencyType,
				&data.BillType,
			)
			if err != nil {
				ss_log.Error("err=[%v]", err)
				continue
			}
			datas = append(datas, &data)
		}
	} else {
		ss_log.Error("err=[%v]", err)
		return nil, "0", ss_err.ERR_SYS_DB_GET
	}

	return datas, total.String, ss_err.ERR_SUCCESS
}

func (ServiceDao) GetAccountNoByServicerNo(servicerNo string) (string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountNo sql.NullString
	sqlStr := "select account_no from servicer where servicer_no =$1 "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&accountNo}, servicerNo)
	if err != nil {
		ss_log.Error("ServiceDao | GetAccountNoByServicerNo |  err= [%v]", err)
		return "", ss_err.ERR_SYS_DB_GET
	}
	return accountNo.String, ss_err.ERR_SUCCESS

}

func (ServiceDao) GetLatLngInfoFromNo(servicerNo string) (string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var latT, lngT, scopeT sql.NullString
	sqlStr := "select  lat, lng, scope from servicer where servicer_no =$1 and is_delete = $2 and use_status = $3"
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&latT, &lngT, &scopeT}, servicerNo, 0, 1)
	if err != nil {
		ss_log.Error("ServiceDao | GetLatLngInfoFromNo |  err= [%v]", err)
		return "", "", ""
	}
	return latT.String, lngT.String, scopeT.String

}

func (ServiceDao) GetServicerInfo(dbHandler *sql.DB, servicerNo string) (*go_micro_srv_cust.ServicerData, string) {

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ser.servicer_no", Val: servicerNo, EqType: "="},
		{Key: "ser.is_delete", Val: "0", EqType: "="},
	})

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT ser.servicer_no,ser.account_no,ser.addr,ser.create_time" +
		",ser.use_status,ser.commission_sharing,ser.income_sharing,ser.income_authorization,ser.outgo_authorization,ser.lat,ser.lng,ser.scope" +
		",acc.nickname, acc.phone " +
		" FROM servicer ser " +
		" LEFT JOIN account acc ON acc.uid = ser.account_no and acc.is_delete = '0' " +
		" LEFT JOIN servicer_terminal mter ON mter.servicer_no = ser.servicer_no " + where2

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	data := &go_micro_srv_cust.ServicerData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return data, ss_err.ERR_SYS_DB_GET
	} else {
		for rows.Next() {
			err = rows.Scan(
				&data.ServicerNo,
				&data.AccountNo,
				&data.Addr,
				&data.CreateTime,
				&data.UseStatus,

				&data.CommissionSharing,
				&data.IncomeSharing,
				&data.IncomeAuthorization,
				&data.OutgoAuthorization,
				&data.Lat,

				&data.Lng,
				&data.Scope,
				&data.Nickname,
				&data.Phone,
			)
			//查询服务商的授权收款额度
			usdAuthCollectLimit, _ := ServiceDaoInst.GetSerAuthCollectLimit(dbHandler, data.AccountNo, "usd")
			data.UsdAuthCollectLimit = usdAuthCollectLimit

			khrAuthCollectLimit, _ := ServiceDaoInst.GetSerAuthCollectLimit(dbHandler, data.AccountNo, "khr")
			data.KhrAuthCollectLimit = khrAuthCollectLimit
		}
	}
	return data, ss_err.ERR_SUCCESS
}

func (ServiceDao) GetServicerCards(accountNo string) ([]*go_micro_srv_cust.ServicerCardPackData, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var datas []*go_micro_srv_cust.ServicerCardPackData
	//查询其卡包
	cardPacksql := "select ca.card_no ,ca.account_no ,ca.channel_no ,ca.name ,ca.create_time ,ca.card_number ,ca.balance_type ,ca.is_defalut," +
		" c.channel_name " +
		" from card ca " +
		" left join channel c on ca.channel_no = c.channel_no " +
		" where ca.account_no = $1 and ca.is_delete = $2 "
	rows2, stmt2, err2 := ss_sql.Query(dbHandler, cardPacksql, accountNo, 0)

	if stmt2 != nil {
		defer stmt2.Close()
	}
	defer rows2.Close()

	if err2 != nil {
		ss_log.Error("ServiceDao | GetServicerCards |  err=%v\nsql=[%v]", err2, cardPacksql)
		return nil, ss_err.ERR_SYS_DB_GET
	} else {
		for rows2.Next() {
			var data go_micro_srv_cust.ServicerCardPackData
			err := rows2.Scan(
				&data.CardNo,
				&data.AccountNo,
				&data.ChannelNo,
				&data.Name,
				&data.CreateTime,
				&data.CardNumber,
				&data.BalanceType,
				&data.IsDefalut,
				&data.ChannelName,
			)
			if err != nil {
				ss_log.Error("err=[%v]", err)
				continue
			}
			datas = append(datas, &data)
		}
	}

	return datas, ss_err.ERR_SUCCESS
}

/*
func (ServiceDao) GetServicerImg(servicerNo string) ([]*go_micro_srv_cust.ServicerImgData, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	servicerImgSQl := "select servicer_img_no, img_id, img_type, create_time, servicer_no " +
		" from servicer_img " +
		" where servicer_no = $1 "
	rows3, stmt3, err3 := ss_sql.Query(dbHandler, servicerImgSQl, servicerNo)

	if stmt3 != nil {
		defer stmt3.Close()
	}
	defer rows3.Close()

	var servicerImgDatalist []*go_micro_srv_cust.ServicerImgData
	if err3 != nil {
		ss_log.Error("err=[%v]", err3)
		return nil, ss_err.ERR_SYS_DB_GET
	} else {
		for rows3.Next() {
			var servicerImgData go_micro_srv_cust.ServicerImgData
			err := rows3.Scan(
				&servicerImgData.ServicerImgNo,
				&servicerImgData.ImgId,
				&servicerImgData.ImgType,
				&servicerImgData.CreateTime,
				&servicerImgData.ServicerNo,
			)
			if err != nil {
				ss_log.Error("err=[%v]", err)
			}

			servicerImgDatalist = append(servicerImgDatalist, &servicerImgData)
		}
	}

	return servicerImgDatalist, ss_err.ERR_SUCCESS
}
*/
func (ServiceDao) GetServicerTerminal(whereModelStr string, whereModelAgrs []interface{}) ([]*go_micro_srv_cust.ServicerTerminalData, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	servicerTerminalSQl := "select terminal_no, terminal_number, pos_sn, use_status " +
		" from servicer_terminal " + whereModelStr
	rows3, stmt3, err3 := ss_sql.Query(dbHandler, servicerTerminalSQl, whereModelAgrs...)

	if stmt3 != nil {
		defer stmt3.Close()
	}
	defer rows3.Close()

	var servicerTerminalDatas []*go_micro_srv_cust.ServicerTerminalData
	if err3 != nil {
		ss_log.Error("err=[%v]", err3)
		return nil, ss_err.ERR_SYS_DB_GET
	} else {
		for rows3.Next() {
			var servicerTerminalData go_micro_srv_cust.ServicerTerminalData
			var posSn, useStatus sql.NullString
			err := rows3.Scan(
				&servicerTerminalData.TerminalNo,
				&servicerTerminalData.TerminalNumber,
				&posSn,
				&useStatus,
			)
			if err != nil {
				ss_log.Error("err=[%v]", err)
			}
			if posSn.String != "" {
				servicerTerminalData.PosSn = posSn.String
			}
			if useStatus.String != "" {
				servicerTerminalData.UseStatus = useStatus.String
			}
			servicerTerminalDatas = append(servicerTerminalDatas, &servicerTerminalData)
		}
	}

	return servicerTerminalDatas, ss_err.ERR_SUCCESS
}

func (ServiceDao) DeleteAccount(tx *sql.Tx, uid string) string {
	sqlStr := "update account set is_delete ='1' where uid =$1 "
	err := ss_sql.ExecTx(tx, sqlStr, uid)

	if err != nil {
		ss_log.Error("err=[%v],删除账号失败", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}

func (ServiceDao) DeleteServicer(tx *sql.Tx, uid string) string {

	sqlStr := "update servicer set is_delete ='1' where servicer_no = $1 "
	err := ss_sql.ExecTx(tx, sqlStr, uid)

	if err != nil {
		ss_log.Error("err=[%v],删除服务商失败", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}

//删除服务商营业执照和营业场所
func (ServiceDao) DeleteServicerImg(tx *sql.Tx, servicerNo string) string {
	sqlStr := "update servicer_img set is_delete ='1' where servicer_no = $1 and is_delete = '0' "
	err := ss_sql.ExecTx(tx, sqlStr, servicerNo)

	if err != nil {
		ss_log.Error("err=[%v],删除服务商营业执照和营业场所图片失败", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}

//删除服务商营业执照或营业场所
func (ServiceDao) DeleteServicerImgByImgType(tx *sql.Tx, servicerNo, imgType string) string {
	sqlStr := "update servicer_img set is_delete ='1' where servicer_no = $1 and img_type = $2 and is_delete = '0' "
	err := ss_sql.ExecTx(tx, sqlStr, servicerNo, imgType)
	if err != nil {
		ss_log.Error("err=[%v],删除图片失败", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}

//获取服务商授权收款额度
func (ServiceDao) GetSerAuthCollectLimit(dbHandler *sql.DB, accountNo, balanceType string) (authCollectLimitR, errR string) {
	var vaType int
	switch balanceType {
	case "usd":
		vaType = constants.VaType_QUOTA_USD
	case "khr":
		vaType = constants.VaType_QUOTA_KHR
	case "usd_spent":
		balanceType = "usd"
		vaType = constants.VaType_QUOTA_USD_REAL
	case "khr_spent":
		balanceType = "khr"
		vaType = constants.VaType_QUOTA_KHR_REAL
	}

	var authCollectLimit sql.NullString
	sqlStr := "select balance from vaccount where account_no =$1 and balance_type = $2 and va_type = $3 and is_delete = '0' "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&authCollectLimit}, accountNo, balanceType, vaType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "0", ss_err.ERR_PARAM
	}

	if authCollectLimit.String == "" {
		return "0", ss_err.ERR_SUCCESS
	} else {
		return authCollectLimit.String, ss_err.ERR_SUCCESS
	}

}

//更新服务商授权收款额度
func (ServiceDao) ModifySerAuthCollectLimit(tx *sql.Tx, accountNo, balanceType, authCollectLimit string) (errR string) {
	var vaType int
	switch balanceType {
	case "usd":
		vaType = constants.VaType_QUOTA_USD
	case "khr":
		vaType = constants.VaType_QUOTA_KHR
	}

	//先确认目标是否创建了授权收款额度的虚拟账户
	count, errCheck := ServiceDaoInst.CheckSerAuthCollectLimit(tx, accountNo, balanceType, vaType)
	if errCheck != nil {
		ss_log.Error("errCheck=[%v]", errCheck)
		return ss_err.ERR_PARAM
	}
	if count == 0 {
		//开始创建该账号的授权虚拟账户
		_, err := ServiceDaoInst.AddSerAuthCollectLimit(tx, accountNo, balanceType, vaType)
		if err != ss_err.ERR_SUCCESS {
			ss_log.Error("err=[%v]", err)
			return ss_err.ERR_PARAM
		}
	}

	sqlStr := "update vaccount set balance = $1 where account_no = $2 and balance_type = $3 and va_type = $4 and is_delete = '0' "
	err := ss_sql.ExecTx(tx, sqlStr, authCollectLimit, accountNo, balanceType, vaType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

//确认是否创建了授权收款额度的虚拟账户
func (ServiceDao) CheckSerAuthCollectLimit(tx *sql.Tx, accountNo, balanceType string, vaType int) (int, error) {
	sqlCnt := "select count(1) from vaccount " +
		" where account_no = $1 and balance_type = $2 and va_type = $3 and is_delete = '0' "
	var count sql.NullString
	err := ss_sql.QueryRowTx(tx, sqlCnt, []*sql.NullString{&count}, accountNo, balanceType, vaType)

	return strext.ToInt(count.String), err

}

func (ServiceDao) AddSerAuthCollectLimit(tx *sql.Tx, accountNo, balanceType string, vaType int) (vaccountNo, errR string) {
	vaccountNoT := strext.NewUUID()
	sqlAdd := "insert into vaccount(vaccount_no, account_no, va_type, balance, create_time, balance_type) " +
		"values($1,$2,$3,$4,current_timestamp,$5)"
	err := ss_sql.ExecTx(tx, sqlAdd, vaccountNoT, accountNo, vaType, "0", balanceType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "", ss_err.ERR_SYS_DB_ADD
	}
	return vaccountNoT, ss_err.ERR_SUCCESS
}

func (ServiceDao) CheckUniqueServicerTerminal(tx *sql.Tx, column, columnData string) (countR int, errR error) {
	sqlCnt := fmt.Sprintf("select count(1) from servicer_terminal where %s = $1 and is_delete ='0' ", column)
	var count sql.NullString
	err := ss_sql.QueryRowTx(tx, sqlCnt, []*sql.NullString{&count}, columnData)

	return strext.ToInt(count.String), err

}

func (ServiceDao) GetServicerNoByAccNo(dbHandler *sql.DB, accNo string) (servicerNo string) {
	sqlStr := "select servicer_no from servicer where account_no = $1 and is_delete = '0' "
	var servicerNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&servicerNoT}, accNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return servicerNoT.String
}

func (ServiceDao) ModifySerPosStatus(terminalNo, useStatus string) (errR string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update servicer_terminal set use_status = $2 where terminal_no = $1 "
	err := ss_sql.Exec(dbHandler, sqlStr, terminalNo, useStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_UPDATE
	}
	return ss_err.ERR_SUCCESS
}

//查看pos属于哪个服务商
func (ServiceDao) GetSerPosServicerNoByTerminalNo(dbHandler *sql.DB, terminalNo string) (servicerNo string) {
	sqlStr := "select servicer_no from servicer_terminal where terminal_no = $1 "
	var servicerNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&servicerNoT}, terminalNo)
	if err != nil {
		return ""
	}
	return servicerNoT.String

}

func (ServiceDao) GetServicerOrderCountDetail(whereModelStr string, whereModelAgrs []interface{}) (dataT *go_micro_srv_cust.GetServicerOrderCountData, err string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	data := &go_micro_srv_cust.GetServicerOrderCountData{}
	sqlStr := "SELECT servicer_no, currency_type, in_num, in_amount, out_num, out_amount, profit_num, profit_amount, recharge_num, recharge_amount, withdraw_num, withdraw_amount, modify_time " +
		" from servicer_count " + whereModelStr

	rows, stmt, errT := ss_sql.Query(dbHandler, sqlStr, whereModelAgrs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return nil, ss_err.ERR_SYS_DB_GET
	} else {
		for rows.Next() {
			//string recharge_total_sum = 9;
			//    string recharge_amount_sum = 10;
			//    string withdraw_total_sum = 11;
			//    string withdraw_amount_sum = 12;
			var servicerNo, balanceType, incomeTotalSum, incomeAmountSum sql.NullString
			var outgoTotalSum, outgoAmountSum, profitAmountSum, profitTotalSum, modifyTime sql.NullString
			var rechargeTotalSum, rechargeAmountSum, withdrawTotalSum, withdrawAmountSum sql.NullString
			errT = rows.Scan(
				&servicerNo,
				&balanceType,
				&incomeTotalSum,
				&incomeAmountSum,
				&outgoTotalSum,
				&outgoAmountSum,
				&profitTotalSum,
				&profitAmountSum,
				&rechargeTotalSum,
				&rechargeAmountSum,
				&withdrawTotalSum,
				&withdrawAmountSum,
				&modifyTime,
			)
			if errT != nil {
				ss_log.Error("err=[%v]", errT)
				return data, ss_err.ERR_SYS_DB_GET
			}

			if incomeTotalSum.String != "" {
				data.IncomeTotalSum = incomeTotalSum.String
			} else {
				data.IncomeTotalSum = "0"
			}

			if incomeAmountSum.String != "" {
				data.IncomeAmountSum = incomeAmountSum.String
			} else {
				data.IncomeAmountSum = "0"
			}

			if outgoTotalSum.String != "" {
				data.OutgoTotalSum = outgoTotalSum.String
			} else {
				data.OutgoTotalSum = "0"
			}

			if outgoAmountSum.String != "" {
				data.OutgoAmountSum = outgoAmountSum.String
			} else {
				data.OutgoAmountSum = "0"
			}

			if profitTotalSum.String != "" {
				data.ProfitTotalSum = profitTotalSum.String
			} else {
				data.ProfitTotalSum = "0"
			}

			if profitAmountSum.String != "" {
				data.ProfitAmountSum = profitAmountSum.String
			} else {
				data.ProfitAmountSum = "0"
			}

			if rechargeTotalSum.String != "" {
				data.RechargeTotalSum = rechargeTotalSum.String
			} else {
				data.RechargeTotalSum = "0"
			}

			if rechargeAmountSum.String != "" {
				data.RechargeAmountSum = rechargeAmountSum.String
			} else {
				data.RechargeAmountSum = "0"
			}

			if withdrawTotalSum.String != "" {
				data.WithdrawTotalSum = withdrawTotalSum.String
			} else {
				data.WithdrawTotalSum = "0"
			}

			if withdrawAmountSum.String != "" {
				data.WithdrawAmountSum = withdrawAmountSum.String
			} else {
				data.WithdrawAmountSum = "0"
			}

			if modifyTime.String != "" {
				data.ModifyTime = modifyTime.String
			} else {
				data.ModifyTime = time.Now().Format("2006-01-02") + " 00:00:00"
			}

			if servicerNo.String != "" {
				data.ServicerNo = servicerNo.String
			} else {
				data.ServicerNo = ""
			}

			if balanceType.String != "" {
				data.BalanceType = balanceType.String
			} else {
				data.BalanceType = ""
			}

		}
	}

	ss_log.Info("data=[%v]", data)
	return data, ss_err.ERR_SUCCESS

}

func (ServiceDao) AddServicerImg(tx *sql.Tx, imgId, imgType, servicerNo string) (err string) { //添加服务商营业执照和营业场所

	sqlStr := "insert into servicer_img(servicer_img_no, img_id, img_type, create_time, servicer_no) values($1,$2,$3,current_timestamp,$4)"
	errT := ss_sql.ExecTx(tx, sqlStr, strext.GetDailyId(), imgId, imgType, servicerNo)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return ss_err.ERR_SAVE_IMAGE_FAILD
	}
	return ss_err.ERR_SUCCESS
}

func (ServiceDao) GetAllSerNoAndAccNo(dbHandler *sql.DB) (datas []*go_micro_srv_cust.GetServicerAccountsData, err string) {
	sqlStr := "select servicer_no, account_no from servicer where is_delete = '0'  "
	rows, stmt, errT := ss_sql.Query(dbHandler, sqlStr)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	var dataT []*go_micro_srv_cust.GetServicerAccountsData
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return nil, ss_err.ERR_SYS_DB_GET
	} else {
		for rows.Next() {
			var data go_micro_srv_cust.GetServicerAccountsData
			errT = rows.Scan(
				&data.ServicerNo,
				&data.AccountNo,
			)
			if errT != nil {
				ss_log.Error("err=[%v]", errT)
				return dataT, ss_err.ERR_SYS_DB_GET
			}
			dataT = append(dataT, &data)
		}
	}
	return dataT, ss_err.ERR_SUCCESS
}

func (ServiceDao) GetAccNoFromSrvNo(sid string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select account_no from servicer where servicer_no =$1 and is_delete='0' limit 1`, []*sql.NullString{&accountNoT}, sid)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return accountNoT.String
}

func (ServiceDao) GetScopeOffNoByPosSn(posSn string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select scope_off from servicer s left join servicer_terminal st on s.servicer_no = st.servicer_no" +
		"  where st.pos_sn = $1 and st.is_delete = '0' "
	var scopeOff sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&scopeOff}, posSn)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return scopeOff.String
}
