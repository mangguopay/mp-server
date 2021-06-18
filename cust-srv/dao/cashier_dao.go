package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

var (
	CashierDaoInstance CashierDao
)

type CashierDao struct{}

// 获取收银员密码
func (CashierDao) GetCashierPwdFromOpAccNo(opAccNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var opPWD sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select op_password from cashier where uid=$1 and is_delete='0' limit 1`, []*sql.NullString{&opPWD}, opAccNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return opPWD.String
}

// 查询手机号和查询收银员ID
func (*CashierDao) QueryPhoneAndCID(uid, rechierType string) (string, string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// 查询手机号,店员id
	var phone, cashierUID sql.NullString
	if err := ss_sql.QueryRow(dbHandler, "SELECT a.phone,r.iden_no  FROM account a  LEFT JOIN rela_acc_iden r ON a.uid = r.account_no  "+
		" WHERE a.uid = $1 and r.account_type = $2 and a.is_delete = 0 LIMIT 1", []*sql.NullString{&phone, &cashierUID}, uid, rechierType); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return "", "", err
	}
	return phone.String, cashierUID.String, nil

}

func (CashierDao) GetServicerNoFromOpAccNo(opAccNo string) (servierNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var serviceNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select servicer_no from cashier where uid=$1 and is_delete='0' limit 1`, []*sql.NullString{&serviceNoT}, opAccNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}

	return serviceNoT.String
}

func (*CashierDao) DeleteCashier(dbHandler *sql.DB, uid string) (errCode string) {
	err := ss_sql.Exec(dbHandler, `update cashier set is_delete='1' where uid=$1`, uid)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}
	return ss_err.ERR_SUCCESS
}

func (*CashierDao) GetCashierNoByAccNo(accNo string) (cashierNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select iden_no from rela_acc_iden where account_no = $1 and account_type = $2 "
	var cashierNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cashierNoT}, accNo, constants.AccountType_POS)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return cashierNoT.String
}

func (*CashierDao) GetSrvAccNoFromCaNo(caNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT s.account_no FROM rela_acc_iden rai LEFT JOIN cashier c ON rai.iden_no = c.uid LEFT JOIN servicer s ON " +
		"s.servicer_no = c.servicer_no WHERE s.is_delete = '0' AND rai.account_no = $1"
	var srvAccNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&srvAccNo}, caNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return srvAccNo.String
}

func (*CashierDao) GetCashiers(whereStr string, whereArgs []interface{}) (datas []*go_micro_srv_cust.CashierData, errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlStr := "select ca.create_time,ca.uid,acc.account,acc.phone,acc.nickname,acc2.account " +
		" from cashier ca " +
		//" left join rela_acc_iden rai on rai.iden_no = ca.uid and account_type = " + constants.AccountType_POS +
		" left join account acc on acc.uid = ca.account_no " +
		" left join servicer ser on ser.servicer_no = ca.servicer_no " +
		" left join account acc2 on acc2.uid = ser.account_no " + whereStr

	rows, stmt, errT := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, ss_err.ERR_SYS_DB_GET
	}
	var datasT []*go_micro_srv_cust.CashierData
	for rows.Next() {
		var data go_micro_srv_cust.CashierData
		errT = rows.Scan(
			&data.CreateTime,
			&data.CashierNo,
			&data.Account,
			&data.Phone,
			&data.Nickname,
			&data.SerAccount,
		)
		if errT != nil {
			ss_log.Error("errT=[%v]", errT)
			continue
		}
		datasT = append(datasT, &data)
	}

	return datasT, ss_err.ERR_SUCCESS
}

func (*CashierDao) GetCashierDetail(whereStr string, whereArgs []interface{}) (dataT *go_micro_srv_cust.CashierData, errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlStr := "select ca.create_time,ca.uid,acc.account,acc.phone,acc.nickname,acc2.account " +
		" from cashier ca " +
		//" left join rela_acc_iden rai on rai.iden_no = ca.uid and account_type = " + constants.AccountType_POS +
		" left join account acc on acc.uid = ca.account_no " +
		" left join servicer ser on ser.servicer_no = ca.servicer_no " +
		" left join account acc2 on acc2.uid = ser.account_no " + whereStr

	rows, stmt, errT := ss_sql.QueryRowN(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}

	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, ss_err.ERR_SYS_DB_GET
	}
	data := &go_micro_srv_cust.CashierData{}
	errT = rows.Scan(
		&data.CreateTime,
		&data.CashierNo,
		&data.Account,
		&data.Phone,
		&data.Nickname,
		&data.SerAccount,
	)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, ss_err.ERR_PARAM
	}

	return data, ss_err.ERR_SUCCESS
}

func (*CashierDao) GetCnt(whereStr string, whereArgs []interface{}) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlStr := "select count(1) " +
		" from cashier ca " +
		//" left join rela_acc_iden rai on rai.iden_no = ca.uid and account_type = " + constants.AccountType_POS +
		" left join account acc on acc.uid = ca.account_no " +
		" left join servicer ser on ser.servicer_no = ca.servicer_no " +
		" left join account acc2 on acc2.uid = ser.account_no " + whereStr
	var totalT sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, whereArgs...)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return ""
	}
	return totalT.String
}

func (CashierDao) DeleteCashierTx(tx *sql.Tx, cashierNo string) string {
	sqlStr := "update cashier set is_delete = '1' where uid = $1 and is_delete = '0' "
	if err := ss_sql.ExecTx(tx, sqlStr, cashierNo); err != nil {
		ss_log.Error("删除店员失败，err=[%v]", err)
		return ss_err.ERR_SYS_DB_DELETE
	}
	return ss_err.ERR_SUCCESS
}

//查询店员的账号
func (CashierDao) GetCashierAccountByCashierNo(opAccNo string) (account string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select acc.account " +
		" from rela_acc_iden rai " +
		" left join account acc on acc.uid = rai.account_no " +
		" where iden_no = $1 and account_type =$2 "
	var accountT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&accountT}, opAccNo, constants.AccountType_POS)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}
	return accountT.String

}

func (*CashierDao) GetSrvAccountNoFromCashierNo(cashierNo string) (string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT acc.account, acc.uid " +
		" FROM cashier ca " +
		" LEFT JOIN servicer ser ON ser.servicer_no = ca.servicer_no " +
		" left join account acc on acc.uid = ser.account_no " +
		" WHERE ser.is_delete = '0' and ca.is_delete = '0' and ca.uid = $1"
	var serAccount, serAccNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&serAccount, &serAccNo}, cashierNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "", ""
	}

	return serAccount.String, serAccNo.String
}
