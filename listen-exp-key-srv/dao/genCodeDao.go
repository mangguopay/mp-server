package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
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
func (*GenCodeDao) GetSrvFromCode(tx *sql.Tx, code string) string {
	var srvAccNoT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select rai.iden_no from gen_code gc  LEFT JOIN rela_acc_iden rai ON gc.account_no = 
										rai.account_no where gc.gen_key=$1 AND rai.account_type = '3' limit 1`,
		[]*sql.NullString{&srvAccNoT}, code)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return srvAccNoT.String
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

// ?????????????????????
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

// ???????????????????????????
func (*GenCodeDao) UpdateGenCodeExp(tx *sql.Tx, status, genKey string) (errCode string) {
	err := ss_sql.ExecTx(tx, `update gen_code set use_status = $1, modify_time =  current_timestamp where gen_key=$2`, status, genKey)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_MODIFY_GEN_KEY_FAILD
	}

	return ss_err.ERR_SUCCESS
}

// ????????????????????? ??????success???????????????
func (*GenCodeDao) CheckCodeTimeExp(tx *sql.Tx, accountNo, genCode, codeType, useStatus string) string {
	var accountNOT, exptT sql.NullString

	if genCode == "" { // ?????????????????????,???????????????????????????????????????
		if err := ss_sql.QueryRowTx(tx, `select account_no,create_time+interval '5 minute'  < current_timestamp  as expt from gen_code where account_no=$1 and code_type=$2  and use_status = $3 order by create_time desc limit 1`,
			//if err := ss_sql.QueryRowTx(tx, `select account_no,create_time+interval '1 minute'  < current_timestamp  as expt from gen_code where account_no=$1 and code_type=$2  and use_status = $3 order by create_time desc limit 1`,
			[]*sql.NullString{&accountNOT, &exptT}, accountNo, codeType, useStatus); err != nil {
			ss_log.Error("err=%v", err)
			//?????????????????????,????????????
			return ss_err.ERR_SUCCESS
		}
		if exptT.String == "false" { //false ?????????,true,??????
			return ss_err.ERR_SUCCESS // ?????????????????????????????????
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

	if exptT.String == "false" { //?????????
		return ss_err.ERR_SUCCESS
	}
	//select account_no,create_time+interval '5 minute'  < current_timestamp  as expt from gen_code WHERE account_no = 'ffd703c1-d78e-43ce-864a-c298c2e0f10a' and code_type='2' and use_status = '1'
	return ""
}

// ???????????????????????????
func (*GenCodeDao) QeurySweepCodeStatusFromAccountNo(tx *sql.Tx, accountNo, genCode string) (string, string, string) {
	var status, orderNoT, sweepAccountNoT sql.NullString

	if err := ss_sql.QueryRowTx(tx, `select use_status,order_no,sweep_account_no  from gen_code where account_no=$1 and code_type=$2 and gen_key=$3 order by create_time desc limit 1`,
		[]*sql.NullString{&status, &orderNoT, &sweepAccountNoT}, accountNo, constants.CODETYPE_SWEEP, genCode); err != nil {
		ss_log.Error("err=%v", err)
		return "", "", ""
	}

	if status.String == "" {
		return "", "", ""
	}
	return status.String, orderNoT.String, sweepAccountNoT.String
}

// ?????????????????????
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
