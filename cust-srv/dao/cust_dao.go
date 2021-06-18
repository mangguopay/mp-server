package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type CustDao struct {
}

var CustDaoInst CustDao

func (CustDao) AddCustTx(tx *sql.Tx, accountNo, gender string) (err error, custNo string) {
	//创建运营商
	custNo = strext.NewUUID()
	sqlStr := "insert into cust(cust_no, account_no, payment_password, gender) " +
		" values ($1,$2,$3,$4)"
	err = ss_sql.ExecTx(tx, sqlStr, custNo, accountNo, "", gender)
	if err != nil {
		ss_log.Error("CustDao |AddCust err=[%v]", err)
	}
	return err, custNo
}

func (CustDao) DeleteCust(uid string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update set is_delete ='1' from cust where cust_no = $1 "
	err := ss_sql.Exec(dbHandler, sqlStr, uid)

	if err != nil {
		ss_log.Error("err=[%v],删除用户失败", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}

func (CustDao) ModifyCustInfo(tx *sql.Tx, inAuthorization, outAuthorization, inTransferAuthorization, outTransferAuthorization, custNo string) string {

	sqlUpdateCust := "Update cust set in_authorization = $1, out_authorization = $2, in_transfer_authorization = $3, out_transfer_authorization = $4 where cust_no = $5"

	err := ss_sql.ExecTx(tx, sqlUpdateCust, inAuthorization, outAuthorization, inTransferAuthorization, outTransferAuthorization, custNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_UPDATE
	}
	return ss_err.ERR_SUCCESS
}

func (CustDao) GetCustInfo(whereList []*model.WhereSqlCond) (custData *go_micro_srv_cust.CustData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlStr := "SELECT c.cust_no, c.gender, c.in_authorization, c.out_authorization,c.in_transfer_authorization,c.out_transfer_authorization" +
		", acc.nickname, acc.create_time, acc.uid, acc.phone, acc.account, acc.usd_balance, acc.khr_balance, acc.individual_auth_status " +
		", am.auth_name, am.auth_number " +
		" FROM cust c " +
		" LEFT JOIN account acc ON acc.uid = c.account_no " +
		" LEFT JOIN auth_material am ON am.auth_material_no = acc.individual_auth_material_no "

	rows, stmt, errT := ss_sql.QueryRowN(dbHandler, sqlStr+whereModel.WhereStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}

	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, errT
	}

	data := &go_micro_srv_cust.CustData{}
	var nickname, createTime, uid sql.NullString
	var phone, account, usdBalance, khrBalance sql.NullString
	var authName, authNumber, authStatus sql.NullString
	err = rows.Scan(
		&data.CustNo,
		&data.Gender,
		&data.InAuthorization,
		&data.OutAuthorization,
		&data.InTransferAuthorization,

		&data.OutTransferAuthorization,
		&nickname,
		&createTime,
		&uid,
		&phone,

		&account,
		&usdBalance,
		&khrBalance,
		&authStatus,
		&authName,

		&authNumber,
	)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return
	}

	data.Nickname = nickname.String
	data.CreateTime = createTime.String
	data.Uid = uid.String
	data.Phone = phone.String
	data.Account = account.String

	data.UsdBalance = usdBalance.String
	data.KhrBalance = khrBalance.String
	data.AuthStatus = authStatus.String
	data.AuthName = authName.String

	data.AuthNumber = authNumber.String

	return data, nil
}

func (CustDao) UpdateTradingAuthority(custNo string, tradingAuthority int) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	return ss_sql.Exec(dbHandler, `update cust set trading_authority=$1 where cust_no=$2`,
		tradingAuthority, custNo)
}
