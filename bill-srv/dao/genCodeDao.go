package dao

import (
	"database/sql"

	"a.a/cu/ss_log"

	"a.a/cu/db"
	"a.a/cu/util"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type GenCodeDao struct {
}

var GenCodeDaoInst GenCodeDao

func (*GenCodeDao) GenCode(accountNo, amount, moneyType, codeType string) (code string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	genKey := util.RandomDigitStr(32)
	err := ss_sql.Exec(dbHandler, `insert into gen_code(gen_key,account_no,amount,money_type,create_time,code_type,use_status) values($1,$2,$3,$4,current_timestamp,$5,'1')`,
		genKey, accountNo, amount, moneyType, codeType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return genKey
}

func (*GenCodeDao) PosGenCode(serverNo, amount, moneyType, codeType, opAccNo string, opAccType int) (code string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	genKey := util.RandomDigitStr(32)
	err := ss_sql.Exec(dbHandler, `insert into gen_code(gen_key,account_no,amount,money_type,create_time,code_type,use_status,op_acc_no,op_acc_type) values($1,$2,$3,$4,current_timestamp,$5,'1',$6,$7)`,
		genKey, serverNo, amount, moneyType, codeType, opAccNo, opAccType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return genKey
}

func (*GenCodeDao) GetRecvCode(code, codeType string) (accountNo, amount, moneyType, createTime, useStatus string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountNoT, amountT, moneyTypeT, createTimeT, useStatusT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select account_no,amount,money_type,create_time,use_status from gen_code where gen_key=$1 and code_type=$2 limit 1`,
		[]*sql.NullString{&accountNoT, &amountT, &moneyTypeT, &createTimeT, &useStatusT}, code, codeType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return "", "", "", "", ""
	}

	return accountNoT.String, amountT.String, moneyTypeT.String, createTimeT.String, useStatusT.String
}
func (*GenCodeDao) GetCodeStatus(tx *sql.Tx, code, codeType string) string {
	var useStatusT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select use_status from gen_code where gen_key=$1 and code_type=$2 limit 1`,
		[]*sql.NullString{&useStatusT}, code, codeType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return useStatusT.String
}

// 此方法获取的account_no就是服务商no,在扫码取款的时候生成码的时候存进去的时候已处理成服务商id
func (*GenCodeDao) GetSrvFromCode(tx *sql.Tx, code string) (string, string, string) {
	var srvAccNoT, opAccNo, opAccType sql.NullString
	err := ss_sql.QueryRowTx(tx, `select account_no,op_acc_no,op_acc_type from gen_code where gen_key=$1 limit 1`,
		[]*sql.NullString{&srvAccNoT, &opAccNo, &opAccType}, code)
	if err != nil {
		ss_log.Error("err=%v", err)
		return "", "", ""
	}
	return srvAccNoT.String, opAccNo.String, opAccType.String
}
func (*GenCodeDao) GetCodeAccNoFromLogNo(logNo string) (string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var genCodeT, sweepAccNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select gen_key,sweep_account_no from gen_code where order_no=$1  limit 1`,
		[]*sql.NullString{&genCodeT, &sweepAccNoT}, logNo)
	if err != nil {
		ss_log.Error("err=%v", err)
		return "", ""
	}
	return genCodeT.String, sweepAccNoT.String
}

// 修改扫码的状态
func (*GenCodeDao) UpdateGenCodeUseStatus(tx *sql.Tx, status int32, genKey, sweepUID, accountType, moneyType string) (errCode string) {
	if accountType == constants.AccountType_SERVICER || accountType == constants.AccountType_POS {
		err := ss_sql.ExecTx(tx, `update gen_code set use_status = $1, modify_time =  current_timestamp where gen_key=$2`, status, genKey)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_MODIFY_GEN_KEY_FAILD
		}
	}

	err := ss_sql.ExecTx(tx, `update gen_code set use_status = $1,sweep_account_no = $2,money_type = $4, modify_time =  current_timestamp where gen_key=$3`, status, sweepUID, genKey, moneyType)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_MODIFY_GEN_KEY_FAILD
	}

	return ss_err.ERR_SUCCESS
}
func (*GenCodeDao) UpdateGenCodeStatus(tx *sql.Tx, status int32, genKey, accountUID string) (errCode string) {

	err := ss_sql.ExecTx(tx, `update gen_code set use_status = $1, sweep_account_no= $3, modify_time =  current_timestamp where gen_key=$2`, status, genKey, accountUID)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_MODIFY_GEN_KEY_FAILD
	}

	return ss_err.ERR_SUCCESS
}

// 修改码为已过期状态
func (*GenCodeDao) UpdateGenCodeExp(tx *sql.Tx, status, genKey string) (errCode string) {
	err := ss_sql.ExecTx(tx, `update gen_code set use_status = $1, modify_time =  current_timestamp where gen_key=$2`, status, genKey)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_MODIFY_GEN_KEY_FAILD
	}

	return ss_err.ERR_SUCCESS
}

// 修改扫码的状态 返回success说明没过期
func (*GenCodeDao) CheckCodeTimeExp(tx *sql.Tx, accountNo, genCode, codeType, useStatus string) string {
	var accountNOT, exptT sql.NullString

	if genCode == "" { // 获取扫一扫的码,判断是否有码还没有处理完的
		if err := ss_sql.QueryRowTx(tx, `select account_no,create_time+interval '5 minute'  < current_timestamp  as expt from gen_code where account_no=$1 and code_type=$2  and use_status = $3 order by create_time desc limit 1`,
			//if err := ss_sql.QueryRowTx(tx, `select account_no,create_time+interval '1 minute'  < current_timestamp  as expt from gen_code where account_no=$1 and code_type=$2  and use_status = $3 order by create_time desc limit 1`,
			[]*sql.NullString{&accountNOT, &exptT}, accountNo, codeType, useStatus); err != nil {
			ss_log.Error("err=%v", err)
			//报错说明找不到,可以生成
			return ss_err.ERR_SUCCESS
		}
		if exptT.String == "false" { //false 没过期,true,过期
			return ss_err.ERR_SUCCESS // 二维码五分钟后才能生成
		}
		return ""
	}

	err := ss_sql.QueryRowTx(tx, `select account_no,create_time+interval '5 minute'  < current_timestamp  as expt from gen_code where gen_key=$1 and code_type=$2 limit 1`,
		//err := ss_sql.QueryRowTx(tx, `select account_no,create_time+interval '1 minute'  < current_timestamp  as expt from gen_code where gen_key=$1 and code_type=$2 limit 1`,
		[]*sql.NullString{&accountNOT, &exptT}, genCode, codeType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}

	if exptT.String == "false" { //没过期
		return ss_err.ERR_SUCCESS
	}
	//select account_no,create_time+interval '5 minute'  < current_timestamp  as expt from gen_code WHERE account_no = 'ffd703c1-d78e-43ce-864a-c298c2e0f10a' and code_type='2' and use_status = '1'
	return ""
}

// 获取扫一扫码的状态
func (*GenCodeDao) QeurySweepCodeStatusFromAccountNo(tx *sql.Tx, accountNo, genCode string) (string, string, string) {
	var status, orderNoT, sweepAccountNoT sql.NullString

	if err := ss_sql.QueryRowTx(tx, `select use_status,order_no,sweep_account_no  from gen_code where op_acc_no=$1 and code_type=$2 and gen_key=$3 order by create_time desc limit 1`,
		[]*sql.NullString{&status, &orderNoT, &sweepAccountNoT}, accountNo, constants.CODETYPE_SWEEP, genCode); err != nil {
		ss_log.Error("err=%v", err)
		return "", "", ""
	}

	if status.String == "" {
		return "", "", ""
	}
	return status.String, orderNoT.String, sweepAccountNoT.String
}

// 修改扫码的状态
func (*GenCodeDao) UpdateOrderNoFormGenCode(tx *sql.Tx, orderNo, genCode string) string {
	err := ss_sql.ExecTx(tx, `update gen_code set order_no = $1, modify_time=current_timestamp where gen_key=$2`, orderNo, genCode)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_MODIFY_GEN_KEY_FAILD
	}
	return ss_err.ERR_SUCCESS
}

func (*GenCodeDao) GetStatusFromCode(code string) (string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var useStatusT, orderNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select use_status,order_no from gen_code where gen_key=$1  limit 1`,
		[]*sql.NullString{&useStatusT, &orderNoT}, code)
	if err != nil {
		ss_log.Error("err=%v", err)
		return "", ""
	}
	return useStatusT.String, orderNoT.String
}

func (*GenCodeDao) GetGenRecvCodeInfo(code, codeType string) (string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var opAccTypeT, opAccNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select  op_acc_type,op_acc_no from gen_code where gen_key=$1 and code_type=$2 limit 1`,
		[]*sql.NullString{&opAccTypeT, &opAccNoT}, code, codeType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return "", ""
	}

	return opAccTypeT.String, opAccNoT.String
}
