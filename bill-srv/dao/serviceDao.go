package dao

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"a.a/cu/container"
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/bill-srv/common"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type ServiceDao struct {
}

var ServiceDaoInst ServiceDao

func (ServiceDao) GetServicerPWDFromOpAccNo(sid string) (servierNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var pwdT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select password from servicer where servicer_no =$1 and is_delete='0' limit 1`, []*sql.NullString{&pwdT}, sid)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return pwdT.String
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

// 获取收款,取款权限
func (ServiceDao) GetPermissionFromSrvNo(srvNo string) (string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var incomeAuthorizationT, outgoAuthorizationT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select income_authorization,outgo_authorization from servicer where servicer_no =$1 and is_delete='0' limit 1`, []*sql.NullString{&incomeAuthorizationT, &outgoAuthorizationT}, srvNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", ""
	}
	return incomeAuthorizationT.String, outgoAuthorizationT.String
}

func (ServiceDao) GetServicerNoByCashierNo(cid string) (servicerNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var Sno sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select servicer_no from cashier where uid =$1 and is_delete='0' limit 1`, []*sql.NullString{&Sno}, cid)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return Sno.String
}

func (ServiceDao) GetSerNoBySerAcc(accountNo string) (servicerNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var Sno sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select servicer_no from servicer where account_no = $1 and is_delete='0' limit 1`, []*sql.NullString{&Sno}, accountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return Sno.String
}

//// todo 没有测试，新版本接口改写
//func (ServiceDao) GetTransferToServicerLogs(pageSize, page int32, startTime, endTime, orderStatus, servicerNo string) (returnDatas []*go_micro_srv_bill.TransferToServicerLogData, returnTotal string, returnErr string) {
//	dbHandler := ss_sql2.NewDbInst(constants.DB_CRM)
//	defer dbHandler.Close()
//
//	datas := []*go_micro_srv_bill.TransferToServicerLogData{}
//	dbHandler.InitWhereList([]*ss_sql2.WhereSqlCond{
//		{Key: "lts.create_time", Val: startTime, EqType: ">="},
//		{Key: "lts.create_time", Val: endTime, EqType: "<="},
//		{Key: "lts.order_type", Val: "1", EqType: "="},
//		{Key: "lts.order_status", Val: orderStatus, EqType: "="},
//		{Key: "lts.servicer_no", Val: servicerNo, EqType: "="},
//	})
//
//	total := strext.ToStringNoPoint(dbHandler.GetCnt("log_to_servicer lts"))
//
//	dbHandler.AppendWhereExtra(`order by lts.create_time desc`)
//	dbHandler.AppendWhereLimitI32(pageSize, page)
//	sqlStr := "SELECT lts.currency_type, lts.amount, lts.order_status, lts.card_no, lts.finish_time, ch.channel_name, c.name, c.card_number" +
//		" FROM log_to_servicer lts " +
//		" LEFT JOIN card c ON c.card_no = lts.card_no " +
//		" LEFT JOIN channel ch ON ch.channel_no = c.channel_no "
//	rows, err := dbHandler.QueryWhere(sqlStr)
//	if err != nil {
//		ss_log.Error("err=%v", err)
//		return nil, "0", ss_err.ERR_SYS_DB_GET
//	}
//
//	for rows.Next() {
//		var data go_micro_srv_bill.TransferToServicerLogData
//		err = rows.Scan(
//			&data.CurrencyType,
//			&data.Amount,
//			&data.OrderStatus,
//			&data.CardNo,
//			&data.FinishTime,
//
//			&data.ChannelName,
//			&data.Name,
//			&data.CardNumber,
//		)
//		if err != nil {
//			ss_log.Error("err=%v", err)
//			continue
//		}
//		datas = append(datas, &data)
//	}
//
//	return datas, total, ss_err.ERR_SUCCESS
//}

func (ServiceDao) GetTransferToServicerLogs(dbHandler *sql.DB, whereModelStr string, whereModelArgs []interface{}) (returnDatas []*go_micro_srv_bill.TransferToServicerLogData, returnErr string) {

	datas := []*go_micro_srv_bill.TransferToServicerLogData{}
	sqlStr := "SELECT lts.currency_type, lts.amount, lts.order_status, lts.card_no, lts.create_time" +
		", ch.channel_name, c.name, c.card_number, ch.logo_img_no " +
		" FROM log_to_servicer lts " +
		" LEFT JOIN card c ON c.card_no = lts.card_no " +
		" LEFT JOIN channel ch ON ch.channel_no = c.channel_no " + whereModelStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModelArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	times := []string{}
	if err == nil {
		for rows.Next() {
			var data go_micro_srv_bill.TransferToServicerLogData
			var logoImgNo, channelName, nameT, cardNumber sql.NullString
			err = rows.Scan(
				&data.CurrencyType,
				&data.Amount,
				&data.OrderStatus,
				&data.CardNo,
				&data.CreateTime,

				&channelName,
				&nameT,
				&cardNumber,
				&logoImgNo,
			)
			if err != nil {
				ss_log.Error("err=%v", err)
				continue
			}

			if channelName.String != "" {
				data.ChannelName = channelName.String
			}
			if nameT.String != "" {
				data.Name = nameT.String
			}
			if cardNumber.String != "" {
				data.CardNumber = cardNumber.String
			}

			isNewTime := container.GetKey(data.CreateTime[:10], times) //-1找不到
			if isNewTime == -1 {
				times = append(times, data.CreateTime[:10])
				data.Time = data.CreateTime[:10]
			}

			// 获取不需要授权的图片路径
			//_, imageBaseUrl, _ := cache.ApiDaoInstance.GetGlobalParam("image_base_url")
			//pathStr := GlobalParamDaoInstance.QeuryParamValue(constants.KEY_STORE_UNAUTH_IMAGE_PATH)
			////根据id查询路径
			//name, err2 := ImageDaoInstance.GetImageUrlById(logoImgNo.String)
			//
			//if err2 != ss_err.ERR_SUCCESS {
			//	ss_log.Error("err2=[%v]", err2)
			//}
			//
			//boolean, _ := file.Exists(pathStr + "/" + name)
			//if !boolean {
			//	ss_log.Error("图片[%v]不存在", logoImgNo.String)
			//} else {
			//	data.LogoImgUrl = imageBaseUrl + "/" + name
			//}

			datas = append(datas, &data)
		}
	} else {
		ss_log.Error("err=[%v]", err)
		return nil, ss_err.ERR_SYS_DB_GET
	}

	return datas, ss_err.ERR_SUCCESS
}

// pos端获取对账信息
func (ServiceDao) GetServicerCheckListTotal(startTime, endTime, servicerNo string) (int32, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// 日期往后加1天，where条件中用小于
	endTime, retErr := ss_time.TimeAfter(endTime, ss_time.DateFormat, time.Hour*24) // 下一天的日期
	if retErr != nil {
		return 0, retErr
	}

	sqlStr := "SELECT count(1) FROM servicer_count_list "
	sqlStr += " WHERE dates >= $1 AND dates < $2 AND servicer_no=$3  "

	var total sql.NullString

	qErr := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&total}, startTime, endTime, servicerNo)
	if qErr != nil {
		return 0, qErr
	}

	return strext.ToInt32(total.String), nil
}

type ServicerCountListStatis struct {
	ServicerNo     string
	CurrencyType   string
	Dates          string
	InNum          int64
	InAmount       int64
	OutNum         int64
	OutAmount      int64
	ProfitNum      int64
	ProfitAmount   int64
	RechargeNum    int64
	RechargeAmount int64
	WithdrawNum    int64
	WithdrawAmount int64
}

// pos端获取对账信息列表
func (s ServiceDao) GetServicerCheckList(startTime, endTime, servicerNo string, page, pageSize int32) (datas []*go_micro_srv_bill.GetServicerCheckListReplyData, total int32, retErr error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// 获取总记录数
	total, retErr = s.GetServicerCheckListTotal(startTime, endTime, servicerNo)
	if retErr != nil {
		return
	}

	// 日期往后加1天，where条件中用小于
	endTime, retErr = ss_time.TimeAfter(endTime, ss_time.DateFormat, time.Hour*24) // 下一天的日期
	if retErr != nil {
		return
	}

	// 按日期查询，按日期desc排序
	startNum := (page - 1) * pageSize

	sqlStr := "SELECT dates FROM servicer_count_list "
	sqlStr += " WHERE dates >= $1 AND dates < $2 AND servicer_no=$3 "
	sqlStr += " GROUP BY dates ORDER BY dates DESC LIMIT $4 OFFSET $5"

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, startTime, endTime, servicerNo, pageSize, startNum)
	if stmt != nil {
		defer stmt.Close()
		defer rows.Close()
	}

	if err != nil {
		retErr = err
		return
	}

	dateList := []string{}

	for rows.Next() {
		var d string
		err := rows.Scan(&d)
		if err != nil {
			ss_log.Error("err=%v", err)
			continue
		}
		dateList = append(dateList, d)
	}

	for _, date := range dateList {
		// 1.获取某一天的数据
		currencyList, err := s.GetServicerCheckListByDate(date, servicerNo)
		if err != nil {
			ss_log.Error("err=%v", err)
			continue
		}

		// 2.添加到返回
		datas = append(datas, &go_micro_srv_bill.GetServicerCheckListReplyData{
			Date:         common.GetPostgresDate(date),
			CurrencyList: currencyList,
		})
	}

	return
}

// pos端获取对账信息列表-按日期获取
func (s ServiceDao) GetServicerCheckListByDate(date, servicerNo string) (currencyList []*go_micro_srv_bill.GetServicerCheckListReplyData_CurrencyList, retErr error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// 按日期查询，按日期desc排序
	sqlStr := "SELECT dates, currency_type, in_num, in_amount, out_num, out_amount, profit_num, profit_amount, "
	sqlStr += " recharge_num, recharge_amount, withdraw_num, withdraw_amount  "
	sqlStr += " FROM servicer_count_list WHERE dates = $1 AND servicer_no=$2 "
	sqlStr += " ORDER BY dates DESC,currency_type DESC "

	rows, stmt, qErr := ss_sql.Query(dbHandler, sqlStr, date, servicerNo)
	if stmt != nil {
		defer stmt.Close()
		defer rows.Close()
	}

	if qErr != nil {
		retErr = qErr
		return
	}

	list := []ServicerCountListStatis{}

	for rows.Next() {
		var dates, currencyType, inNum, inAmount, outNum, outAmount sql.NullString
		var profitNum, profitAmount, rechargeNum, rechargeAmount, withdrawNum, withdrawAmount sql.NullString

		err := rows.Scan(&dates, &currencyType, &inNum, &inAmount, &outNum, &outAmount,
			&profitNum, &profitAmount, &rechargeNum, &rechargeAmount, &withdrawNum, &withdrawAmount,
		)
		if err != nil {
			ss_log.Error("err=%v", err)
			continue
		}
		list = append(list, ServicerCountListStatis{
			Dates:          dates.String,
			CurrencyType:   currencyType.String,
			InNum:          strext.ToInt64(inNum.String),
			InAmount:       strext.ToInt64(inAmount.String),
			OutNum:         strext.ToInt64(outNum.String),
			OutAmount:      strext.ToInt64(outAmount.String),
			ProfitNum:      strext.ToInt64(profitNum.String),
			ProfitAmount:   strext.ToInt64(profitAmount.String),
			RechargeNum:    strext.ToInt64(rechargeNum.String),
			RechargeAmount: strext.ToInt64(rechargeAmount.String),
			WithdrawNum:    strext.ToInt64(withdrawNum.String),
			WithdrawAmount: strext.ToInt64(withdrawAmount.String),
		})
	}

	// 美元数据
	usdRsult := []*go_micro_srv_bill.GetServicerCheckListReplyData_CurrencyList_Result{}
	// 瑞尔数据
	khrRsult := []*go_micro_srv_bill.GetServicerCheckListReplyData_CurrencyList_Result{}

	for _, row := range list {
		switch row.CurrencyType {
		case constants.CURRENCY_USD: // 美元
			usdRsult = GetCurrencyListResult(usdRsult, row)
		case constants.CURRENCY_KHR: // 瑞尔
			khrRsult = GetCurrencyListResult(khrRsult, row)
		}
	}

	// 美元
	currencyList = append(currencyList, &go_micro_srv_bill.GetServicerCheckListReplyData_CurrencyList{
		CurrencyType: strings.ToUpper(constants.CURRENCY_USD),
		Results:      usdRsult,
	})

	// 瑞尔
	currencyList = append(currencyList, &go_micro_srv_bill.GetServicerCheckListReplyData_CurrencyList{
		CurrencyType: strings.ToUpper(constants.CURRENCY_KHR),
		Results:      khrRsult,
	})

	return
}

// 获取每种获取类型的结果数据
func GetCurrencyListResult(list []*go_micro_srv_bill.GetServicerCheckListReplyData_CurrencyList_Result, row ServicerCountListStatis) []*go_micro_srv_bill.GetServicerCheckListReplyData_CurrencyList_Result {
	//存款
	list = append(list, &go_micro_srv_bill.GetServicerCheckListReplyData_CurrencyList_Result{
		Type:   strext.ToInt64(constants.BILL_TYPE_INCOME),
		Num:    row.InNum,
		Amount: row.InAmount,
	})

	//取款
	list = append(list, &go_micro_srv_bill.GetServicerCheckListReplyData_CurrencyList_Result{
		Type:   strext.ToInt64(constants.BILL_TYPE_OUTGO),
		Num:    row.OutNum,
		Amount: row.OutAmount,
	})

	// 收益(佣金)
	list = append(list, &go_micro_srv_bill.GetServicerCheckListReplyData_CurrencyList_Result{
		Type:   strext.ToInt64(constants.BILL_TYPE_PROFIT),
		Num:    row.ProfitNum,
		Amount: row.ProfitAmount,
	})

	// 充值
	list = append(list, &go_micro_srv_bill.GetServicerCheckListReplyData_CurrencyList_Result{
		Type:   strext.ToInt64(constants.BILL_TYPE_RECHARGE),
		Num:    row.RechargeNum,
		Amount: row.RechargeAmount,
	})

	// 提现
	list = append(list, &go_micro_srv_bill.GetServicerCheckListReplyData_CurrencyList_Result{
		Type:   strext.ToInt64(constants.BILL_TYPE_WITHDRAWALS),
		Num:    row.WithdrawNum,
		Amount: row.WithdrawAmount,
	})

	return list
}

func (ServiceDao) GetCntServicerProfitLedgers(dbHandler *sql.DB, whereModelStr string, whereModelArgs []interface{}) (totalR string) {
	//统计
	var total sql.NullString
	sqlCnt := "SELECT count(1) FROM servicer_profit_ledger spl " + whereModelStr
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModelArgs...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "0"
	}
	return total.String
}

func (ServiceDao) GetSumAmount(dbHandler *sql.DB, whereModelStr string, whereModelArgs []interface{}) (sumAmountT string) {
	//统计
	var sumAmount sql.NullString
	sqlCnt := "SELECT count(spl.actual_income) FROM servicer_profit_ledger spl " + whereModelStr
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&sumAmount}, whereModelArgs...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "0"
	}

	if sumAmount.String == "" {
		sumAmount.String = "0"
	}

	return sumAmount.String
}

func (ServiceDao) GetCnt(dbHandler *sql.DB, tabname, whereModelStr string, whereModelArgs []interface{}) (totalR string) {
	//统计
	var total sql.NullString
	sqlCnt := fmt.Sprintf(`SELECT count(1) FROM %s %s`, tabname, whereModelStr)
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModelArgs...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "0"
	}
	return total.String
}

//获取拥有数据的时间(去重后的年月日)
func (ServiceDao) GetHaveDataTime(dbHandler *sql.DB, tabname, timeColumn, whereModelStr string, whereModelArgs []interface{}) (timeR []string, errR string) {
	sqlGetTime := fmt.Sprintf("select distinct %s from %s %s ", timeColumn, tabname, whereModelStr)

	var timeS []string
	rowsGetTime, stmtGetTime, errGetTime := ss_sql.Query(dbHandler, sqlGetTime, whereModelArgs...)
	if stmtGetTime != nil {
		defer stmtGetTime.Close()
	}
	defer rowsGetTime.Close()
	if errGetTime == nil {
		for rowsGetTime.Next() {
			var time sql.NullString
			errGetTime = rowsGetTime.Scan(
				&time,
			)
			if errGetTime != nil {
				ss_log.Error("errGetTime=[%v]", errGetTime)
				continue
			}
			timeS = append(timeS, time.String)
		}
	} else {
		ss_log.Error("err=[%v]", errGetTime)
		return nil, ss_err.ERR_SYS_DB_GET
	}

	return timeS, ss_err.ERR_SUCCESS

}

func (ServiceDao) GetServicerProfitLedgers(dbHandler *sql.DB, whereModelStr string, whereModelArgs []interface{}, mModel *model.WhereSql) (datas []*go_micro_srv_bill.ServicerProfitLedgersData, returnErr string) {
	//以下是查询数据
	sqlStr := "SELECT spl.log_no, spl.amount_order, spl.actual_income, spl.payment_time, spl.currency_type, spl.order_type " +
		" FROM servicer_profit_ledger spl " + whereModelStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModelArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	//用于记录已统计的天（年月日）
	time := []string{}
	if err == nil {
		for rows.Next() {
			data := &go_micro_srv_bill.ServicerProfitLedgersData{}
			err = rows.Scan(
				&data.OrderNo,
				&data.AmountOrder,
				&data.ActualIncome, //实际所得
				&data.CreateTime,
				&data.CurrencyType,

				&data.OrderType,
			)
			if err != nil {
				ss_log.Error("err=%v", err)
				continue
			}

			isNewTime := container.GetKey(data.CreateTime[:10], time) //-1找不到
			if isNewTime == -1 {
				time = append(time, data.CreateTime[:10])
				data.Time = data.CreateTime[:10]

				whereModelUsdSum := ss_sql.SsSqlFactoryInst.DeepClone(mModel)
				ss_sql.SsSqlFactoryInst.AppendWhere(whereModelUsdSum, "to_char(spl.payment_time,'yyyy-MM-dd')", data.CreateTime[:10], "=")
				ss_sql.SsSqlFactoryInst.AppendWhere(whereModelUsdSum, "spl.currency_type", "usd", "=")
				usdSum, _ := ServiceDaoInst.GetServicerProfitByTime(dbHandler, whereModelUsdSum.WhereStr, whereModelUsdSum.Args)
				data.UsdSum = usdSum

				whereModelKhrSum := ss_sql.SsSqlFactoryInst.DeepClone(mModel)
				ss_sql.SsSqlFactoryInst.AppendWhere(whereModelKhrSum, "to_char(spl.payment_time,'yyyy-MM-dd')", data.CreateTime[:10], "=")
				ss_sql.SsSqlFactoryInst.AppendWhere(whereModelKhrSum, "spl.currency_type", "khr", "=")
				khrSum, _ := ServiceDaoInst.GetServicerProfitByTime(dbHandler, whereModelKhrSum.WhereStr, whereModelKhrSum.Args)
				data.KhrSum = khrSum
			}

			datas = append(datas, data)
		}
	} else {
		ss_log.Error("err=[%v]", err)
		return nil, ss_err.ERR_SYS_DB_GET
	}

	return datas, ss_err.ERR_SUCCESS
}

func (ServiceDao) GetServicerProfitLedgerDetail(logNo string) (data *go_micro_srv_bill.ServicerProfitLedgersData, returnErr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "spl.log_no", Val: logNo, EqType: "="},
	})

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT spl.amount_order, spl.actual_income, spl.servicefee_amount_sum, spl.split_proportion, spl.payment_time, spl.log_no, spl.currency_type" +
		", bdr.order_no, bdr.bill_type " +
		" FROM servicer_profit_ledger spl " +
		//" LEFT JOIN billing_details_results bdr ON spl.general_ledger_no = bdr.bill_no " + where2
		" LEFT JOIN billing_details_results bdr ON spl.log_no = bdr.order_no " + where2
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	data = &go_micro_srv_bill.ServicerProfitLedgersData{}
	if err == nil {
		for rows.Next() {
			err = rows.Scan(
				&data.AmountOrder,
				&data.ActualIncome,
				&data.ServicefeeAmountSum,
				&data.SplitProportion,
				&data.CreateTime,
				&data.LogNo,
				&data.CurrencyType,
				&data.OrderNo,
				&data.OrderType,
			)
			if err != nil {
				ss_log.Error("err=%v", err)
				continue
			}
		}
	} else {
		ss_log.Error("err=[%v]", err)
		return nil, ss_err.ERR_SYS_DB_GET
	}
	return data, ss_err.ERR_SUCCESS
}

func (ServiceDao) GetServicerProfitByTime(dbHandler *sql.DB, whereModelStr string, whereModelArgs []interface{}) (sumR string, errR string) {
	var sum sql.NullString
	sqlStr := "select sum(spl.actual_income) from servicer_profit_ledger spl " + whereModelStr
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&sum}, whereModelArgs...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "0", ss_err.ERR_PARAM
	}
	if sum.String == "" {
		return "0", ss_err.ERR_SUCCESS
	} else {
		return sum.String, ss_err.ERR_SUCCESS
	}

}

func (ServiceDao) GetTransferToHeadquartersLog(dbHandler *sql.DB, whereModelStr string, whereModelArgs []interface{}) (datas []*go_micro_srv_bill.TransferToHeadquartersLogData, returnErr string) {
	//以下是查询数据
	sqlStr := "SELECT lth.log_no,lth.currency_type, lth.amount, lth.order_status, lth.card_no, lth.create_time " +
		", ch.channel_name, c.name, c.card_number, ch.logo_img_no " +
		" FROM log_to_headquarters lth " +
		" LEFT JOIN card_head c ON c.card_no = lth.card_no " +
		" LEFT JOIN channel ch ON ch.channel_no = c.channel_no " + whereModelStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModelArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	times := []string{}
	if err == nil {
		for rows.Next() {
			var data go_micro_srv_bill.TransferToHeadquartersLogData
			var logoImgNo, logNo, channelName, name, cardNumber sql.NullString
			err = rows.Scan(
				&logNo,
				&data.CurrencyType,
				&data.Amount,
				&data.OrderStatus,
				&data.CardNo,
				&data.CreateTime,

				&channelName,
				&name,
				&cardNumber,
				&logoImgNo,
			)
			if err != nil {
				ss_log.Error("err=[%v],logNo=[%v]", err, logNo.String)
				continue
			}

			if channelName.String != "" {
				data.ChannelName = channelName.String
			}
			if name.String != "" {
				data.Name = name.String
			}
			if cardNumber.String != "" {
				data.CardNumber = cardNumber.String
			}

			isNewTime := container.GetKey(data.CreateTime[:10], times) //-1找不到
			if isNewTime == -1 {
				times = append(times, data.CreateTime[:10])
				data.Time = data.CreateTime[:10]
			}

			//if logoImgNo.String != "" {
			//	// 获取不需要授权的图片路径
			//	_, imageBaseUrl, _ := cache.ApiDaoInstance.GetGlobalParam("image_base_url")
			//	pathStr := GlobalParamDaoInstance.QeuryParamValue(constants.KEY_STORE_UNAUTH_IMAGE_PATH)
			//	//根据id查询路径
			//	name, err2 := ImageDaoInstance.GetImageUrlById(logoImgNo.String)
			//
			//	if err2 != ss_err.ERR_SUCCESS {
			//		ss_log.Error("err2=[%v]", err2)
			//	}
			//
			//	boolean, _ := file.Exists(pathStr + "/" + name)
			//	if !boolean {
			//		ss_log.Error("图片[%v]不存在", logoImgNo.String)
			//	} else {
			//		data.LogoImgUrl = imageBaseUrl + "/" + name
			//	}
			//}

			datas = append(datas, &data)
		}
	} else {
		ss_log.Error("err=[%v]", err)
		return nil, ss_err.ERR_SYS_DB_GET
	}
	return datas, ss_err.ERR_SUCCESS
}
