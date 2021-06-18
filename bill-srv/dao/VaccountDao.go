package dao

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	tmProto "a.a/mp-server/common/proto/tm"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/common/ss_struct"
)

type VaccountDao struct {
}

var VaccountDaoInst VaccountDao

func (VaccountDao) GetVaccountNo(accountNo string, vaType int32) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var vaccountNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select vaccount_no from vaccount where account_no=$1 and va_type=$2 and is_delete='0' limit 1`,
		[]*sql.NullString{&vaccountNo}, accountNo, vaType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return vaccountNo.String
}
func (VaccountDao) GetVaccountNoFromMoneyType(accountNo, moneyType string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var vaccountNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select vaccount_no from vaccount where account_no=$1 and balance_type=$2 and is_delete='0' limit 1`,
		[]*sql.NullString{&vaccountNo}, accountNo, moneyType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return vaccountNo.String
}

func (VaccountDao) InitVaccountNo(accountNo, balanceType string, vaType int32) (vaccountNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	balanceType = strings.ToLower(balanceType)

	vaccountNo = strext.NewUUID()
	err := ss_sql.Exec(dbHandler, `insert into vaccount(vaccount_no,account_no,va_type,balance,create_time,is_delete,use_status,frozen_balance,balance_type) values ($1,$2,$3,$4,current_timestamp,$5,$6,$7,$8)`,
		vaccountNo, accountNo, vaType, "0", "0", "1", "0", balanceType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return vaccountNo
}

func (VaccountDao) GetBalanceV2(tmProxy *ss_struct.TmServerProxy, vaccountNo string) (string, string, error) {
	sql := `select balance,frozen_balance from vaccount where vaccount_no=$1 limit 1`
	args := []string{vaccountNo}

	rsp, err := tmProxy.GetTmServer().TxQueryRow(context.TODO(), &tmProto.TxQueryRowRequest{FromServerId: tmProxy.GetFromServerId(), TxNo: tmProxy.GetTxNo(), Sql: sql, Args: args})
	if err != nil {
		return "0", "0", err
	}
	ss_log.Info("GetBalanceV2 map-------------> %v", rsp.Datas)
	return rsp.Datas["balance"], rsp.Datas["frozen_balance"], nil
}
func (VaccountDao) GetBalance(vaccountNo string) (balance, frozenBalance string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var balanceT, frozenBalanceT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select balance,frozen_balance from vaccount where vaccount_no=$1 limit 1`,
		[]*sql.NullString{&balanceT, &frozenBalanceT}, vaccountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "0", "0"
	}
	return balanceT.String, frozenBalanceT.String
}
func (VaccountDao) GetBalance2Tx(tx *sql.Tx, vaccountNo string) (balance, frozenBalance string) {
	var balanceT, frozenBalanceT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select balance,frozen_balance from vaccount where vaccount_no=$1 limit 1`,
		[]*sql.NullString{&balanceT, &frozenBalanceT}, vaccountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "0", "0"
	}
	return balanceT.String, frozenBalanceT.String
}
func (VaccountDao) GetAccNoFromVaccNo(tx *sql.Tx, vaccountNo string) string {
	var accNoT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select account_no from vaccount where vaccount_no=$1 limit 1`,
		[]*sql.NullString{&accNoT}, vaccountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return accNoT.String
}

// 同名进出，余额必须正
func (VaccountDao) SameAccFromAToBUpperZero(tx *sql.Tx, toAmount, accountNo, vaccountNoFrom, vaccountNoTo, amount, logNo, reason string) (errCode string) {
	// 出
	err := ss_sql.ExecTx(tx, `update vaccount set balance=balance-$1 where vaccount_no=$2 and account_no=$3`, toAmount, vaccountNoFrom, accountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoFrom, constants.VaOpType_Minus, toAmount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return errCode
	}

	// 进
	err = ss_sql.ExecTx(tx, `update vaccount set balance=balance+$1 where vaccount_no=$2 and account_no=$3`, amount, vaccountNoTo, accountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoTo, constants.VaOpType_Add, amount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return errCode
	}

	var tmp sql.NullString
	err = ss_sql.QueryRowTx(tx, `select balance from vaccount where vaccount_no=$1 and account_no=$2 and is_delete='0' limit 1`,
		[]*sql.NullString{&tmp}, vaccountNoFrom, accountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	if strext.ToInt64(tmp.String) < 0 {
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	tmp.String = "-1"
	tmp.Valid = false
	err = ss_sql.QueryRowTx(tx, `select balance from vaccount where vaccount_no=$1 and account_no=$2 and is_delete='0' limit 1`,
		[]*sql.NullString{&tmp}, vaccountNoTo, accountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	if strext.ToInt64(tmp.String) < 0 {
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	return ss_err.ERR_SUCCESS
}

// 虚拟账号进出，余额必须正
func (r VaccountDao) AccFromAToBUpperZero(tx *sql.Tx, vaccountNoFrom, vaccountNoTo, amount, logNo, reason string) (errCode string) {

	// 出
	err := ss_sql.ExecTx(tx, `update vaccount set balance=balance-$1 where vaccount_no=$2 and is_delete='0'`, amount, vaccountNoFrom)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoFrom, constants.VaOpType_Minus, amount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return errCode
	}

	// 进
	err = ss_sql.ExecTx(tx, `update vaccount set balance=balance+$1 where vaccount_no=$2 and is_delete='0'`, amount, vaccountNoTo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoTo, constants.VaOpType_Add, amount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return errCode
	}

	errCode, balance := r.GetBalanceTx(tx, vaccountNoFrom)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	if strext.ToInt64(balance) < 0 {
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode, balance = r.GetBalanceTx(tx, vaccountNoTo)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	if strext.ToInt64(balance) < 0 {
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	return ss_err.ERR_SUCCESS
}

// 存款,虚拟账号进出，余额必须正
func (r VaccountDao) SaveMoneyAccFromAToBUpperZero(tx *sql.Tx, vaccountNoFrom, vaccountNoTo, fromAmount, toAmount, logNo, reason string) (errCode string) {

	// 出
	err := ss_sql.ExecTx(tx, `update vaccount set balance=balance-$1,modify_time=current_timestamp where vaccount_no=$2 and is_delete='0'`, fromAmount, vaccountNoFrom)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoFrom, constants.VaOpType_Minus, fromAmount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return errCode
	}

	// 进
	err = ss_sql.ExecTx(tx, `update vaccount set balance=balance+$1,modify_time=current_timestamp where vaccount_no=$2 and is_delete='0'`, toAmount, vaccountNoTo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoTo, constants.VaOpType_Add, toAmount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		return errCode
	}

	errCode, balance := r.GetBalanceTx(tx, vaccountNoFrom)
	if errCode != ss_err.ERR_SUCCESS {
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	if strext.ToInt64(balance) < 0 {
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode, balance = r.GetBalanceTx(tx, vaccountNoTo)
	if errCode != ss_err.ERR_SUCCESS {
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	if strext.ToInt64(balance) < 0 {
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	return ss_err.ERR_SUCCESS
}

//已经不需要同步余额到account表中的usd_balance、khr_balance了，直接查询虚帐账户1、2、3、4就可以了
func (VaccountDao) SyncAccRemainV2(tmProxy *ss_struct.TmServerProxy, accNo string) error {
	var sql string
	var args []string

	sql = `update account set usd_balance=(select sum(balance) from vaccount where account_no=$1 and va_type in($2,$3) and is_delete='0') where uid=$1 and is_delete='0'`
	args = []string{accNo, strext.ToStringNoPoint(constants.VaType_USD_DEBIT), strext.ToStringNoPoint(constants.VaType_FREEZE_USD_DEBIT)}
	rsp, err := tmProxy.GetTmServer().TxExec(context.TODO(), &tmProto.TxExecRequest{FromServerId: tmProxy.GetFromServerId(), TxNo: tmProxy.GetTxNo(), Sql: sql, Args: args})
	if err != nil {
		return err
	}
	if rsp.Err != "" {
		return errors.New(rsp.Err)
	}
	sql = `update account set khr_balance=(select sum(balance) from vaccount where account_no=$1 and va_type in($2,$3) and is_delete='0') where uid=$1 and is_delete='0'`
	args = []string{accNo, strext.ToStringNoPoint(constants.VaType_KHR_DEBIT), strext.ToStringNoPoint(constants.VaType_FREEZE_KHR_DEBIT)}
	rsp1, err1 := tmProxy.GetTmServer().TxExec(context.TODO(), &tmProto.TxExecRequest{FromServerId: tmProxy.GetFromServerId(), TxNo: tmProxy.GetTxNo(), Sql: sql, Args: args})
	if err1 != nil {
		return err
	}
	if rsp1.Err != "" {
		return errors.New(rsp1.Err)
	}

	return nil
}
func (VaccountDao) SyncAccRemain(tx *sql.Tx, accNo string) (errCode string) {
	err := ss_sql.ExecTx(tx, `update account set usd_balance=(select sum(balance) from vaccount where account_no=$1 and va_type in($2,$3) and is_delete='0') where uid=$1 and is_delete='0'`, accNo, constants.VaType_USD_DEBIT, constants.VaType_FREEZE_USD_DEBIT)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	err = ss_sql.ExecTx(tx, `update account set khr_balance=(select sum(balance) from vaccount where account_no=$1 and va_type in($2,$3) and is_delete='0') where uid=$1 and is_delete='0'`, accNo, constants.VaType_KHR_DEBIT, constants.VaType_FREEZE_KHR_DEBIT)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	return ss_err.ERR_SUCCESS
}

// 修改虚拟账户余额，余额必须正
func (r VaccountDao) ModifyVaccRemainUpperZeroV2(tmProxy *ss_struct.TmServerProxy, vaccountNo, amount, op, logNo, reason string) error {
	var opType, sql string
	var args []string
	switch op {
	case "+":
		opType = constants.VaOpType_Add
		sql = `update vaccount set balance=balance+$1 where vaccount_no=$2 and is_delete='0'`
		args = []string{amount, vaccountNo}
	case "-":
		opType = constants.VaOpType_Minus
		sql = `update vaccount set balance=balance-$1 where vaccount_no=$2 and is_delete='0'`
		args = []string{amount, vaccountNo}
	default:
		return errors.New("操作类型错误")
	}

	rsp, err := tmProxy.GetTmServer().TxExec(context.TODO(), &tmProto.TxExecRequest{FromServerId: tmProxy.GetFromServerId(), TxNo: tmProxy.GetTxNo(), Sql: sql, Args: args})
	if err != nil {
		return err
	}
	if rsp.Err != "" {
		return errors.New(rsp.Err)
	}
	ss_log.Info("  vaccountNo : %v, opType: %v, amount: %v, logNo: %v, reason: %v", vaccountNo, opType, amount, logNo, reason)
	insertLogErr := LogVaccountDaoInst.InsertLogTxV2(tmProxy, vaccountNo, opType, amount, logNo, reason)
	if insertLogErr != nil {
		return insertLogErr
	}

	balance, balanceErr := r.GetBalanceTxV2(tmProxy, vaccountNo)
	if balanceErr != nil {
		return balanceErr
	}
	if strext.ToInt64(balance) < 0 {
		return errors.New("余额不足")
	}

	return nil
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
		ss_log.Error("errCode=%v", errCode)
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

// 提现修改虚拟账户余额,冻结资金，余额必须正
func (r VaccountDao) ModifyVaccRemainAndFrozenUpperZero(tx *sql.Tx, op, vaccountNoFrom, amount, logNo, reason, logVaccOpType string) (errCode string) {
	switch op {
	case "+":
		err := ss_sql.ExecTx(tx, `update vaccount set balance=balance+$1,frozen_balance=frozen_balance-$1,modify_time = current_timestamp where vaccount_no=$2 and is_delete='0' `, amount, vaccountNoFrom)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	case "-":
		err := ss_sql.ExecTx(tx, `update vaccount set balance=balance-$1,frozen_balance=frozen_balance+$1,modify_time = current_timestamp where vaccount_no=$2 and is_delete='0'`, amount, vaccountNoFrom)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoFrom, logVaccOpType, amount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		return errCode
	}

	errCode, balance := r.GetBalanceTx(tx, vaccountNoFrom)
	if ss_err.ERR_SUCCESS != errCode {
		return errCode
	}
	if strext.ToInt64(balance) < 0 {
		return ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH
	}
	return ss_err.ERR_SUCCESS
}

// 返还转账订单冻结的钱到余额里去 修改虚拟账户余额,冻结资金，余额必须正
func (r VaccountDao) ModifyVaccFrozenToBalance(tx *sql.Tx, vaccountNoFrom, amount, logNo, reason, logVaccOpType string) (errCode string) {
	err := ss_sql.ExecTx(tx, `update vaccount set balance=balance+$1,frozen_balance=frozen_balance-$1,modify_time = current_timestamp where vaccount_no=$2 and is_delete='0' `, amount, vaccountNoFrom)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoFrom, logVaccOpType, amount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		return errCode
	}

	var tmp sql.NullString
	errT := ss_sql.QueryRowTx(tx, `select frozen_balance from vaccount where vaccount_no = $1 and is_delete='0' limit 1`,
		[]*sql.NullString{&tmp}, vaccountNoFrom)
	if nil != errT {
		ss_log.Error("err=%v", errT)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	if strext.ToInt64(tmp.String) < 0 {
		ss_log.Error("----->%s", "超出解冻金额")
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	return ss_err.ERR_SUCCESS
}

//  冻结资金，余额必须正
func (r VaccountDao) ModifyVaccFrozenUpperZero1(tx *sql.Tx, op, vaccountNoFrom, amount, logNo, reason, logVaccOpType string) (errCode string) {
	switch op {
	case "+":
		err := ss_sql.ExecTx(tx, `update vaccount set  frozen_balance=frozen_balance+$1,modify_time = current_timestamp where vaccount_no=$2 and is_delete='0' `, amount, vaccountNoFrom)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	case "-":
		err := ss_sql.ExecTx(tx, `update vaccount set  frozen_balance=frozen_balance-$1,modify_time = current_timestamp where vaccount_no=$2 and is_delete='0'`, amount, vaccountNoFrom)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoFrom, logVaccOpType, amount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		return errCode
	}
	var tmp sql.NullString
	err := ss_sql.QueryRowTx(tx, `select frozen_balance from vaccount where vaccount_no=$1 and is_delete='0' limit 1`,
		[]*sql.NullString{&tmp}, vaccountNoFrom)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH
	}
	if strext.ToInt64(tmp.String) < 0 {
		ss_log.Error("----->%s", "超出解冻金额")
		return ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH
	}

	return ss_err.ERR_SUCCESS
}

// 提现修改虚拟账,解冻，余额必须正
func (VaccountDao) ModifyVaccFrozenUpperZero(tx *sql.Tx, vaccountNoFrom, amount, logNo, reason, fees string) (errCode string) {
	err := ss_sql.ExecTx(tx, `update vaccount set frozen_balance=frozen_balance-$1 where vaccount_no=$2 and is_delete='0'`, amount, vaccountNoFrom)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode = LogVaccountDaoInst.InsertPosConfirmWithdrawLogTx(tx, vaccountNoFrom, constants.VaOpType_Defreeze, amount, logNo, reason, fees)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return errCode
	}

	var tmp sql.NullString
	err = ss_sql.QueryRowTx(tx, `select frozen_balance from vaccount where vaccount_no=$1 and is_delete='0' limit 1`,
		[]*sql.NullString{&tmp}, vaccountNoFrom)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH
	}
	if strext.ToInt64(tmp.String) < 0 {
		ss_log.Error("----->%s", "超出解冻金额")
		return ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH
	}
	return ss_err.ERR_SUCCESS
}

func (VaccountDao) GetBalanceTxV2(tmProxy *ss_struct.TmServerProxy, vaccountNo string) (string, error) {
	sql := `select balance from vaccount where vaccount_no=$1 and is_delete='0' limit 1`
	args := []string{vaccountNo}

	rsp, err := tmProxy.GetTmServer().TxQueryRow(context.TODO(), &tmProto.TxQueryRowRequest{FromServerId: tmProxy.GetFromServerId(), TxNo: tmProxy.GetTxNo(), Sql: sql, Args: args})
	if err != nil {
		return "", err
	}
	return rsp.Datas["balance"], nil
}
func (VaccountDao) GetBalanceTx(tx *sql.Tx, vaccountNo string) (errCode, balance string) {
	var tmp sql.NullString
	err := ss_sql.QueryRowTx(tx, `select balance from vaccount where vaccount_no=$1 and is_delete='0' limit 1`,
		[]*sql.NullString{&tmp}, vaccountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH, ""
	}
	return ss_err.ERR_SUCCESS, tmp.String
}

func (VaccountDao) GetBalanceFromAccNo(accNo string, vaType int) (error, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var balanceT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select balance from vaccount where account_no=$1 and va_type = $2 and is_delete='0' limit 1`,
		[]*sql.NullString{&balanceT}, accNo, vaType)
	if nil != err {
		return err, "0"
	}

	if balanceT.String == "" {
		balanceT.String = "0"
	}

	return nil, balanceT.String
}

func (VaccountDao) ConfirmExistVAccount(accountNo, balanceType string, vaType int32) (vAccountNo string) {
	vAccountNo = VaccountDaoInst.GetVaccountNo(accountNo, vaType)
	if vAccountNo == "" {
		vAccountNo = VaccountDaoInst.InitVaccountNo(accountNo, balanceType, vaType)
	}
	return vAccountNo
}
