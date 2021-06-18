package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type AgreementDao struct {
}

var AgreementDaoInst AgreementDao

func (AgreementDao) GetCnt(dbHandler *sql.DB, whereStr string, whereArgs []interface{}) (total string) {
	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM agreement " + whereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereArgs...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0"
	}
	return totalT.String
}

func (AgreementDao) GetAgreements(dbHandler *sql.DB, whereStr string, whereArgs []interface{}) (datas []*go_micro_srv_cust.AgreementData, err string) {

	sqlStr := " select id, text, lang, type, create_time, modify_time, use_status " +
		" FROM agreement " + whereStr
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
		var data go_micro_srv_cust.AgreementData
		errT = rows.Scan(
			//select id, text, lang, type, create_time, modify_time, use_status
			&data.Id,
			&data.Text,
			&data.Lang,
			&data.Type,
			&data.CreateTime,
			&data.ModifyTime,
			&data.UseStatus,
		)
		if errT != nil {
			ss_log.Error("errT=[%v]", errT)
			continue
		}
		datas = append(datas, &data)
	}
	return datas, ss_err.ERR_SUCCESS
}

func (AgreementDao) GetAgreementDetail(dbHandler *sql.DB, whereStr string, whereArgs []interface{}) (data *go_micro_srv_cust.AgreementData, err string) {
	sqlStr := " select id, text, lang, type, create_time, modify_time, use_status " +
		" FROM agreement " + whereStr

	rows, stmt, errT := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, ss_err.ERR_SYS_DB_GET
	}

	dataT := &go_micro_srv_cust.AgreementData{}
	for rows.Next() {
		errT = rows.Scan(
			&dataT.Id,
			&dataT.Text,
			&dataT.Lang,
			&dataT.Type,
			&dataT.CreateTime,
			&dataT.ModifyTime,
			&dataT.UseStatus,
		)
		if errT != nil {
			ss_log.Error("errT=[%v]", errT)
			continue
		}
	}
	return dataT, ss_err.ERR_SUCCESS
}

func (AgreementDao) GetAgreementUseStatus(tx *sql.Tx, id string) (useStatus, err string) {
	sqlStr := "select use_status from agreement where id = $1 and is_delete = '0' "
	var useStatusT sql.NullString
	errT := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&useStatusT}, id)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", ss_err.ERR_SYS_DB_GET
	}

	return useStatusT.String, ss_err.ERR_SUCCESS

}

func (AgreementDao) DeleteAgreement(tx *sql.Tx, id string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//确认删除的不是正在使用状态中的协议
	useStatusOld, errGet := AgreementDaoInst.GetAgreementUseStatus(tx, id)
	if errGet != ss_err.ERR_SUCCESS {
		ss_log.Error("errGet=[%v]", errGet)
		return ss_err.ERR_PARAM
	}
	if useStatusOld == "1" { //"协议使用中，不可删除或修改为不使用"
		return ss_err.ERR_AGREEMENT_BEING_USE
	}

	sqlStr := "update agreement set is_delete = '1' where id = $1 and is_delete = '0' and use_status = '0' "
	err := ss_sql.ExecTx(tx, sqlStr, id)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_UPDATE
	}
	return ss_err.ERR_SUCCESS

}

func (AgreementDao) ModifyAgreementStatus(tx *sql.Tx, id, lang, typeT, useStatus string) string {
	//查询同种语言和同种协议下使用状态是启用的数目。
	total, errCnt := AgreementDaoInst.CheckAgreementUseStatusTotal(tx, lang, typeT)
	if errCnt != ss_err.ERR_SUCCESS {
		ss_log.Error("errCnt=[%v]", errCnt)
		return ss_err.ERR_PARAM
	}

	if total > 0 { //将其原来使用协议的使用状态改变为不使用
		errUpdate := AgreementDaoInst.ModifyAgreementStatusTx(tx, lang, typeT)
		if errUpdate != ss_err.ERR_SUCCESS {
			ss_log.Error("errUpdate=[%v]", errUpdate)
			return ss_err.ERR_PARAM
		}
	}

	sqlStr := "update agreement set use_status = $2 where id = $1 and is_delete = '0' "
	err := ss_sql.ExecTx(tx, sqlStr, id, useStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_UPDATE
	}
	return ss_err.ERR_SUCCESS

}

func (AgreementDao) UpdateAgreement(tx *sql.Tx, id, text, lang, typeT, useStatus string) string {

	//确认同种语言和同种协议下只有一个的使用状态是启用的。
	if useStatus == "1" {
		total, errCnt := AgreementDaoInst.CheckAgreementUseStatusTotal(tx, lang, typeT)
		if errCnt != ss_err.ERR_SUCCESS {
			ss_log.Error("errCnt=[%v]", errCnt)
			return ss_err.ERR_PARAM
		}

		if total > 0 { //将其原来使用协议的使用状态改变为不使用
			errUpdate := AgreementDaoInst.ModifyAgreementStatusTx(tx, lang, typeT)
			if errUpdate != ss_err.ERR_SUCCESS {
				ss_log.Error("errUpdate=[%v]", errUpdate)
				return ss_err.ERR_PARAM
			}
		}
	} else {
		//确认当前的使用状态不是启用
		useStatusOld, errGet := AgreementDaoInst.GetAgreementUseStatus(tx, id)
		if errGet != ss_err.ERR_SUCCESS {
			ss_log.Error("errGet=[%v]", errGet)
			return ss_err.ERR_PARAM
		}
		if useStatusOld == "1" { //"协议使用中，不可删除或修改为不使用"
			return ss_err.ERR_AGREEMENT_BEING_USE
		}
	}

	sqlStr := "update agreement set text = $2, lang = $3, type = $4, use_status = $5 ,modify_time = current_timestamp where id = $1 "
	err := ss_sql.ExecTx(tx, sqlStr, id, text, lang, typeT, useStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_UPDATE
	}
	return ss_err.ERR_SUCCESS
}

//查询同种语言和同种协议下使用状态是启用的数目。
func (AgreementDao) CheckAgreementUseStatusTotal(tx *sql.Tx, lang, typeT string) (total int, err string) {
	sqlCnt := "select count(1) from agreement where lang= $1 and type= $2 and use_status='1' "
	var totalT sql.NullString
	errCnt := ss_sql.QueryRowTx(tx, sqlCnt, []*sql.NullString{&totalT}, lang, typeT)
	if errCnt != nil {
		ss_log.Error("errCnt=[%v]", errCnt)
		return strext.ToInt(totalT.String), ss_err.ERR_PARAM
	}
	return strext.ToInt(totalT.String), ss_err.ERR_SUCCESS
}

func (AgreementDao) ModifyAgreementStatusTx(tx *sql.Tx, lang, typeT string) string {
	sqlStr := "update agreement set use_status = '0' where lang = $1 and type = $2 and use_status = '1' "
	err := ss_sql.ExecTx(tx, sqlStr, lang, typeT)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_UPDATE
	}
	return ss_err.ERR_SUCCESS

}

func (d AgreementDao) AddAgreement(tx *sql.Tx, text, lang, typeT, useStatus string) (idT, errT string) {
	//确认同种语言和同种协议下只有一个的使用状态是启用的。
	if useStatus == "1" {
		total, errCnt := AgreementDaoInst.CheckAgreementUseStatusTotal(tx, lang, typeT)
		if errCnt != ss_err.ERR_SUCCESS {
			ss_log.Error("errCnt=[%v]", errCnt)
			return "", ss_err.ERR_PARAM
		}

		if total > 0 { //将其原来使用协议的使用状态改变为不使用
			errUpdate := AgreementDaoInst.ModifyAgreementStatusTx(tx, lang, typeT)
			if errUpdate != ss_err.ERR_SUCCESS {
				ss_log.Error("errUpdate=[%v]", errUpdate)
				return "", ss_err.ERR_PARAM
			}
		}
	}

	id := strext.GetDailyId()
	sqlStr := "insert into agreement(id, text, lang, type, use_status, create_time, modify_time) " +
		"values($1,$2,$3,$4,$5,current_timestamp,current_timestamp)"
	err := ss_sql.ExecTx(tx, sqlStr, id, text, lang, typeT, useStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "", ss_err.ERR_SYS_DB_ADD
	}
	return id, ss_err.ERR_SUCCESS
}
