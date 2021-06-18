package dao

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type LogVaccountDao struct {
}

var LogVaccountDaoInst LogVaccountDao

func (LogVaccountDao) GetCnt(whereStr string, whereArgs []interface{}) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM log_vaccount lv " +
		" LEFT JOIN vaccount vacc ON vacc.vaccount_no = lv.vaccount_no "
	sqlCnt += fmt.Sprintf(" AND vacc.va_type in (%v, %v)", constants.VaType_USD_BUSINESS_SETTLED, constants.VaType_KHR_BUSINESS_SETTLED)
	sqlCnt += whereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereArgs...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0"
	}
	return totalT.String
}

func (LogVaccountDao) GetSumAmt(whereStr string, whereArgs []interface{}) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT sum(lv.amount) " +
		" FROM log_vaccount lv " +
		" LEFT JOIN vaccount vacc ON vacc.vaccount_no = lv.vaccount_no " + whereStr
	if cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereArgs...); cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0"
	}
	return totalT.String
}

func (LogVaccountDao) GetAppCustBills(dbHandler *sql.DB, whereStr string, whereArgs []interface{}, accountNo string) (datas []*go_micro_srv_cust.CustBillsData, err string) {
	sqlStr := "SELECT lv.log_no, lv.create_time, lv.amount, vacc.balance_type, lv.reason, lv.balance, lv.biz_log_no, lv.op_type " +
		" FROM log_vaccount lv " +
		" LEFT JOIN vaccount vacc ON vacc.vaccount_no = lv.vaccount_no " + whereStr
	rows, stmt, errT := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, ss_err.ERR_SYS_DB_GET
	}

	for rows.Next() {
		var data go_micro_srv_cust.CustBillsData
		errT = rows.Scan(
			&data.LogNo,
			&data.CreateTime,
			&data.Amount,
			&data.CurrencyType,
			&data.Reason,
			&data.Balance,
			&data.OrderNo,
			&data.OpType,
		)
		if errT != nil {
			ss_log.Error("errT=[%v]", errT)
			continue
		}

		data.OrderType = data.Reason

		switch data.OpType { //将opType转成前端认识的1+2-
		case constants.VaOpType_Add:
		case constants.VaOpType_Minus:
		case constants.VaOpType_Freeze:
			data.OpType = constants.VaOpType_Minus
		case constants.VaOpType_Defreeze_Minus:
			data.OpType = constants.VaOpType_Minus
		case constants.VaOpType_Defreeze_Add:
			data.OpType = constants.VaOpType_Add
		default:
			ss_log.Error("发现未处理过的OpType[%v],biz_log_no[%v]", data.OpType, data.OrderNo)
		}

		switch data.Reason {
		case constants.VaReason_FEES:
			whereList := []*model.WhereSqlCond{
				{Key: "lv.biz_log_no", Val: data.OrderNo, EqType: "="},
				{Key: "vacc.account_no", Val: accountNo, EqType: "="},
				{Key: "lv.reason", Val: constants.VaReason_FEES, EqType: "!="}, //非服务费的
				{Key: "lv.create_time", Val: data.CreateTime, EqType: "="},     //同一创建时间的
			}

			//如果是手续费，则要返回产生该笔手续费的订单类型,方便可以查询该笔手续费的订单详情
			orderType := LogVaccountDaoInst.GetFeesOrderType(whereList)
			//这是产生该笔手续费的订单类型
			data.OrderType = orderType
			switch orderType {
			case "":
				ss_log.Error("查询产生该笔手续费订单类型失败")
			case constants.VaReason_Cancel_withdraw: //取款取消的手续费处理
				//这是产生该笔手续费的订单类型
				data.OrderType = constants.VaReason_OUTGO
			case constants.VaReason_Cust_Save: //银行卡存款成功的手续费处理
				if data.OpType == constants.VaOpType_Freeze { //如果是手续费并且op_type是3（冻结），说明是存款申请手续费，那么应当是丢弃掉
					continue
				}
			default:

			}
		case constants.VaReason_Cancel_withdraw: //取款失败的(pos取消取款)
			//用来查看详情时用的，不管它是成功还是失败订单，此处都应该去查的是取款订单
			data.OrderType = constants.VaReason_OUTGO
		case constants.VaReason_TRANSFER: //转账的现需要返回核销码
			code, errGetCode := WriteoffDaoInst.GetCode(data.OrderNo, data.OrderType)
			if errGetCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取核销码Code失败,err=[%v]", errGetCode)
			}
			data.Code = code
		default:

		}

		datas = append(datas, &data)
	}
	return datas, ss_err.ERR_SUCCESS
}

func (LogVaccountDao) GetServicerBills(dbHandler *sql.DB, whereStr string, whereArgs []interface{}, accountNo string) (datas []*go_micro_srv_cust.ServicerBillDetailData, err error) {
	sqlStr := "SELECT lv.create_time, lv.amount, vacc.balance_type, lv.reason, lv.balance, lv.biz_log_no, lv.op_type " +
		" FROM log_vaccount lv " +
		" LEFT JOIN vaccount vacc ON vacc.vaccount_no = lv.vaccount_no " + whereStr
	rows, stmt, errT := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, errT
	}

	for rows.Next() {
		var data go_micro_srv_cust.ServicerBillDetailData
		errT = rows.Scan(
			&data.CreateTime,
			&data.Amount,
			&data.BalanceType,
			&data.Reason,
			&data.Balance,
			&data.LogNo,
			&data.OpType,
		)
		if errT != nil {
			ss_log.Error("errT=[%v]", errT)
			continue
		}

		data.OrderType = data.Reason

		//if data.Amount
		if strings.Contains(data.Amount, "-") {
			data.Amount = ss_count.Sub("0", data.Amount).String()
		}
		switch data.OrderType {
		case constants.VaReason_FEES: //服务商的手续费并不是同一时间给到
			whereList := []*model.WhereSqlCond{
				{Key: "lv.biz_log_no", Val: data.LogNo, EqType: "="},
				{Key: "vacc.account_no", Val: accountNo, EqType: "="},
				{Key: "lv.reason", Val: constants.VaReason_FEES, EqType: "!="}, //非服务费的
			}

			//如果是手续费，则要返回产生该笔手续费的订单类型,方便可以查询该笔手续费的订单详情
			orderType := LogVaccountDaoInst.GetFeesOrderType(whereList)
			//这是产生该笔手续费的订单类型
			data.OrderType = orderType

			switch orderType {
			case constants.VaReason_INCOME:
			case constants.VaReason_OUTGO: //用户取款产生的手续费对服务商来说，手续费是+
				data.OpType = constants.VaOpType_Add
			default:

			}
		case constants.VaReason_OUTGO: //用户取款产生的订单金额对服务商来说，订单金额是+
			data.OpType = constants.VaOpType_Add
		case constants.VaReason_Srv_Withdraw:
			if data.OpType == constants.VaOpType_Balance_Frozen_Add {
				//服务商提现申请成功将冻结部分余额到冻结金额
				data.OpType = constants.VaOpType_Minus
			} else if data.OpType == constants.VaOpType_Balance_Defreeze_Add { //驳回的
				data.OpType = constants.VaOpType_Add
			}

		default:

		}

		data.Balance = ss_count.Sub("0", data.Balance).String()

		datas = append(datas, &data)
	}
	return datas, nil
}

func (LogVaccountDao) GetBusinessBills(whereStr string, whereArgs []interface{}, accountNo string) (datas []*go_micro_srv_cust.GetBusinessBillData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT lv.create_time, lv.amount, vacc.balance_type, lv.reason, lv.balance, lv.biz_log_no, lv.op_type, vacc.va_type " +
		" FROM log_vaccount lv " +
		" LEFT JOIN vaccount vacc ON vacc.vaccount_no = lv.vaccount_no " + whereStr
	rows, stmt, errT := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, errT
	}

	for rows.Next() {
		var data go_micro_srv_cust.GetBusinessBillData
		errT = rows.Scan(
			&data.CreateTime,
			&data.Amount,
			&data.CurrencyType,
			&data.Reason,
			&data.Balance,
			&data.LogNo,
			&data.OpType,
			&data.VaType,
		)
		if errT != nil {
			ss_log.Error("errT=[%v]", errT)
			continue
		}

		data.OrderType = data.Reason

		switch data.OrderType {
		case constants.VaReason_FEES: //服务商的手续费并不是同一时间给到
			whereList := []*model.WhereSqlCond{
				{Key: "lv.biz_log_no", Val: data.LogNo, EqType: "="},
				{Key: "vacc.account_no", Val: accountNo, EqType: "="},
				{Key: "lv.reason", Val: constants.VaReason_FEES, EqType: "!="}, //非服务费的
			}

			//如果是手续费，则要返回产生该笔手续费的订单类型,方便可以查询该笔手续费的订单详情
			orderType := LogVaccountDaoInst.GetFeesOrderType(whereList)
			//这是产生该笔手续费的订单类型
			data.OrderType = orderType

		default:

		}

		datas = append(datas, &data)
	}
	return datas, nil
}

//获取服务费的订单类型
func (LogVaccountDao) GetFeesOrderType(whereList []*model.WhereSqlCond) (orderType string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " limit 1")
	sqlStr := "select lv.reason from log_vaccount lv" +
		" LEFT JOIN vaccount vacc ON vacc.vaccount_no = lv.vaccount_no " + whereModel.WhereStr
	var orderTypeT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&orderTypeT}, whereModel.Args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return orderTypeT.String
}

func (VaccountDao) GetBalance(tx *sql.Tx, vaccountNo string) (balance, frozenBalance string) {
	var balanceT, frozenBalanceT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select balance,frozen_balance from vaccount where vaccount_no=$1 limit 1`,
		[]*sql.NullString{&balanceT, &frozenBalanceT}, vaccountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "0", "0"
	}
	return balanceT.String, frozenBalanceT.String
}

func (LogVaccountDao) InsertPosConfirmWithdrawLogTx(tx *sql.Tx, vaccountNo, opType, amount, bizLogNo, reason, fees string) string {
	balance, fbalance := VaccountDaoInst.GetBalance(tx, vaccountNo)
	var resultBalance, resultFbalance string
	if fees != "" && fees != "0" { // 还没扣除手续费的时候,需要把手续费加进balance
		// 当前余额需要相加
		resultBalance = ss_count.Add(balance, fees)
		resultFbalance = fbalance
	} else {
		resultBalance = balance
		resultFbalance = fbalance
	}
	err := ss_sql.ExecTx(tx, `insert into log_vaccount(log_no,create_time,vaccount_no,amount,op_type,frozen_balance,balance,reason,biz_log_no) values ($1,current_timestamp,$2,$3,$4,$5,$6,$7,$8)`,
		strext.GetDailyId(), vaccountNo, amount, opType, resultFbalance, resultBalance, reason, bizLogNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}

// 修改虚拟账户余额，余额必须正
func (r VaccountDao) ModifyVaccRemainUpperZero(tx *sql.Tx, vaccountNo, amount, op, logNo, reason string) (errCode string) {
	var opType string
	switch op {
	case "+":
		opType = constants.VaOpType_Add
		err := ss_sql.ExecTx(tx, `update vaccount set balance=balance+$1 where vaccount_no=$2 and is_delete='0'`, amount, vaccountNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	case "-":
		opType = constants.VaOpType_Minus
		err := ss_sql.ExecTx(tx, `update vaccount set balance=balance-$1 where vaccount_no=$2 and is_delete='0'`, amount, vaccountNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	default:
		return ss_err.ERR_PAY_VACCOUNT_OP_MISSING
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNo, opType, amount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", errCode)
		return errCode
	}

	errCode, balance := r.GetBalanceTx(tx, vaccountNo)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("errCode=%v", errCode)
		return errCode
	}
	if strext.ToInt64(balance) < 0 {
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	return ss_err.ERR_SUCCESS
}

func (LogVaccountDao) InsertLogTx(tx *sql.Tx, vaccountNo, opType, amount, bizLogNo, reason string) string {
	balance, fbalance := VaccountDaoInst.GetBalance(tx, vaccountNo)

	err := ss_sql.ExecTx(tx, `insert into log_vaccount(log_no,create_time,vaccount_no,amount,op_type,frozen_balance,balance,reason,biz_log_no) values ($1,current_timestamp,$2,$3,$4,$5,$6,$7,$8)`,
		strext.GetDailyId(), vaccountNo, amount, opType, fbalance, balance, reason, bizLogNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}

// 获取用户一小时内的虚帐日志
func (*LogVaccountDao) GetUserFrozenBalance(nowTime string) (datasList []*go_micro_srv_cust.LogVaccountData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	endTime := nowTime
	startTime, tErr := ss_time.TimeAfter(nowTime, ss_time.DateFormat, -time.Hour*1) //获取一小时前的时间
	if tErr != nil {
		return nil, tErr
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "lv.create_time", Val: startTime, EqType: ">="},
		{Key: "lv.create_time", Val: endTime, EqType: "<="},
		{Key: "vacc.va_type", Val: "('" + strext.ToStringNoPoint(constants.VaType_USD_DEBIT) + "','" + strext.ToStringNoPoint(constants.VaType_KHR_DEBIT) + "')", EqType: "in"},
	})

	sqlStr := `select lv.frozen_balance, lv.reason, lv.op_type, lv.amount, vacc.balance_type, lv.create_time, lv.biz_log_no, vacc.account_no
	from log_vaccount lv
	left join vaccount vacc on vacc.vaccount_no = lv.vaccount_no `
	rows, stmt, errT := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, errT
	}
	var datas []*go_micro_srv_cust.LogVaccountData
	for rows.Next() {
		var data go_micro_srv_cust.LogVaccountData
		errT = rows.Scan(
			&data.FrozenBalance,
			&data.Reason,
			&data.OpType,
			&data.Amount,
			&data.BalanceType,
			&data.CreateTime,
			&data.BizLogNo,
			&data.AccountNo,
		)
		if errT != nil {
			ss_log.Error("errT=[%v]", errT)
			continue
		}
		datas = append(datas, &data)
	}
	return datas, nil
}

type LogVAccount struct {
	LogNo         string
	VAccountNo    string
	Amount        string
	OpType        string
	FrozenBalance string
	Balance       string
	Reason        string
	BizLogNo      string
	CreateTime    string
}

//查询虚账日志
func (*LogVaccountDao) CountLogVAccount(whereStr string, whereArgs []interface{}) (int32, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlCount := "SELECT COUNT(1) FROM log_vaccount "

	var totalNum sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlCount+whereStr, []*sql.NullString{&totalNum}, whereArgs...); nil != err {
		return -1, err
	}

	return strext.ToInt32(totalNum.String), nil
}

func (*LogVaccountDao) GetLogVAccountList(whereStr string, whereArgs []interface{}) ([]*LogVAccount, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlSelect := "SELECT log_no, op_type, amount, balance, reason, create_time FROM log_vaccount "
	rows, stmt, err := ss_sql.Query(dbHandler, sqlSelect+whereStr, whereArgs...)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var list []*LogVAccount
	for rows.Next() {
		var logNo, opType, amount, balance, reason, createTime sql.NullString
		err := rows.Scan(&logNo, &opType, &amount, &balance, &reason, &createTime)
		if err != nil {
			return nil, err
		}
		obj := new(LogVAccount)
		obj.LogNo = logNo.String
		obj.OpType = opType.String
		obj.Amount = amount.String
		obj.Balance = balance.String
		obj.Reason = reason.String
		obj.CreateTime = createTime.String
		list = append(list, obj)
	}
	return list, nil
}

//查询虚账日志的业务流水号
func (*LogVaccountDao) GetBizLogNoAndOpType(logNo, reason string) (bizLogNo string, opType string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var bizLogNoT, opTypeT sql.NullString
	sqlStr := "SELECT biz_log_no,op_type FROM log_vaccount WHERE log_no=$1 AND reason=$2 "
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&bizLogNoT, &opTypeT}, logNo, reason); nil != err {
		return "", "", err
	}
	return bizLogNoT.String, opTypeT.String, nil
}

//获取商户入账金额
func (r LogVaccountDao) GetBusinessRecordedAmount(vAccountNo string, startTime, endTime string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT SUM(amount) FROM log_vaccount WHERE vaccount_no=$1 AND op_type in($2,$3,$4) AND create_time>=$5 AND create_time<=$6 "
	var totalAmount sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalAmount},
		vAccountNo, constants.VaOpType_Minus, constants.VaOpType_Defreeze, constants.VaOpType_Defreeze_Minus, startTime, endTime)
	if err != nil {
		return "", err
	}

	return totalAmount.String, nil
}

//获取商户出账金额
func (r LogVaccountDao) GetBusinessExpenditureAmount(vAccountNo string, startTime, endTime string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT sum(amount) FROM log_vaccount WHERE vaccount_no=$1 AND op_type=$2 AND create_time>=$3 AND create_time<=$4 "
	var totalAmount sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalAmount}, vAccountNo, constants.VaOpType_Add, startTime, endTime)
	if err != nil {
		return "", err
	}

	return totalAmount.String, nil
}

type BusinessProfitLog struct {
	CreateTime   string
	Amount       string
	CurrencyType string
	Reason       string
	Balance      string
	BizLogNo     string
	LogNo        string
	OpType       string
	VAccType     string
}

//查询商家收益明细列表
func (LogVaccountDao) GetBusinessProfitList(whereStr string, whereArgs []interface{}) ([]*BusinessProfitLog, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT lv.log_no, lv.create_time, lv.amount, lv.reason, lv.balance, lv.biz_log_no, lv.op_type, vacc.balance_type, vacc.va_type " +
		" FROM log_vaccount lv " +
		" LEFT JOIN vaccount vacc ON vacc.vaccount_no = lv.vaccount_no "
	sqlStr += fmt.Sprintf(" AND vacc.va_type in (%v, %v)", constants.VaType_USD_BUSINESS_SETTLED, constants.VaType_KHR_BUSINESS_SETTLED)
	sqlStr += whereStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var list []*BusinessProfitLog
	for rows.Next() {
		var logNo, createTime, amount, reason, balance, bizLogNo, opType, currencyType, vAccType sql.NullString
		if err = rows.Scan(&logNo, &createTime, &amount, &reason, &balance, &bizLogNo, &opType, &currencyType, &vAccType); err != nil {
			return nil, err
		}
		data := &BusinessProfitLog{
			LogNo:        logNo.String,
			CreateTime:   createTime.String,
			Amount:       amount.String,
			Reason:       reason.String,
			Balance:      balance.String,
			BizLogNo:     bizLogNo.String,
			OpType:       opType.String,
			CurrencyType: currencyType.String,
			VAccType:     vAccType.String,
		}
		list = append(list, data)
	}
	return list, nil
}
