package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"context"
	"database/sql"
)

type ConsultationConfigDao struct {
}

var ConsultationConfigDaoInst ConsultationConfigDao

func (ConsultationConfigDao) GetCnt(dbHandler *sql.DB, whereStr string, whereArgs []interface{}) (total string) {
	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM consultation_config " + whereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereArgs...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0"
	}
	return totalT.String
}

func (ConsultationConfigDao) GetConsultationConfigs(dbHandler *sql.DB, whereStr string, whereArgs []interface{}) (datas []*go_micro_srv_cust.ConsultationConfigData, err string) {
	sqlStr := " select id, use_status, lang, idx, logo_img_no, name, text " +
		" FROM consultation_config " + whereStr
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
		var data go_micro_srv_cust.ConsultationConfigData
		errT = rows.Scan(
			&data.Id,
			&data.UseStatus,
			&data.Lang,
			&data.Idx,
			&data.LogoImgNo,

			&data.Name,
			&data.Text,
		)
		if errT != nil {
			ss_log.Error("errT=[%v]", errT)
			continue
		}

		datas = append(datas, &data)
	}
	return datas, ss_err.ERR_SUCCESS
}

func (ConsultationConfigDao) GetConsultationConfigDetail(dbHandler *sql.DB, whereStr string, whereArgs []interface{}) (data *go_micro_srv_cust.ConsultationConfigData, err string) {
	sqlStr := " select id, name, text, logo_img_no, lang, idx, use_status " +
		" FROM consultation_config " + whereStr

	rows, stmt, errT := ss_sql.QueryRowN(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}

	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, ss_err.ERR_SYS_DB_GET
	}
	//todo 尚未修改完
	dataT := &go_micro_srv_cust.ConsultationConfigData{}
	errT = rows.Scan(
		&dataT.Id,
		&dataT.Name,
		&dataT.Text,
		&dataT.LogoImgNo,
		&dataT.Lang,

		&dataT.Idx,
		&dataT.UseStatus,
	)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, ss_err.ERR_PARAM
	}

	return dataT, ss_err.ERR_SUCCESS
}

func (ConsultationConfigDao) DeleteConsultationConfig(tx *sql.Tx, id string) string {
	sqlStr := "update consultation_config set is_delete = '1', idx = '-1' where id = $1 and is_delete = '0' "
	err := ss_sql.ExecTx(tx, sqlStr, id)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_UPDATE
	}
	return ss_err.ERR_SUCCESS

}

func (ConsultationConfigDao) ModifyConsultationConfigStatus(id, useStatus string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update consultation_config set use_status = $2 where id = $1 and is_delete = '0' "
	err := ss_sql.Exec(dbHandler, sqlStr, id, useStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_UPDATE
	}
	return ss_err.ERR_SUCCESS

}

func (ConsultationConfigDao) UpdateConsultationConfig(tx *sql.Tx, id, name, text, logoImgNo, useStatus string) string {
	sqlStr := "update consultation_config set name = $2, text = $3, logo_img_no = $4, use_status = $5 where id = $1 and is_delete = '0' "
	err := ss_sql.ExecTx(tx, sqlStr, id, name, text, logoImgNo, useStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_UPDATE
	}
	return ss_err.ERR_SUCCESS
}

func (d ConsultationConfigDao) AddConsultationConfig(tx *sql.Tx, name, text, logoImgNo, useStatus, lang, idx string) string {
	sqlStr := "insert into consultation_config(id, name, text, logo_img_no, use_status, lang, idx, create_time) values($1,$2,$3,$4,$5,$6,$7,current_timestamp)"
	err := ss_sql.ExecTx(tx, sqlStr, strext.GetDailyId(), name, text, logoImgNo, useStatus, lang, idx)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_ADD
	}
	return ss_err.ERR_SUCCESS
}

func (d ConsultationConfigDao) GetMaxIdx(lang string) (idxMax int, errStr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select idx from consultation_config where is_delete = '0' and lang = $1 order by idx desc limit 1 "
	var idxMaxT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&idxMaxT}, lang)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return 1, ss_err.ERR_PARAM
	}

	return strext.ToInt(idxMaxT.String), ss_err.ERR_SUCCESS

}
func (d ConsultationConfigDao) GetLogoURL(id string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select di.image_url  from consultation_config cc LEFT JOIN  dict_images di ON cc.logo_img_no = di.image_id " +
		"WHERE cc.id = $1 and cc.is_delete = 0"
	var imageURL sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&imageURL}, id)
	return imageURL.String, err
}

// 获取idx
func (ConsultationConfigDao) GetNearIdxConsultationNo(idx, swapType, lang string) (id string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var helpNoT sql.NullString
	switch swapType {
	case constants.SwapType_Up: // 上层
		err := ss_sql.QueryRow(dbHandler, `select id from consultation_config where idx=$1-1 and is_delete = '0' and lang = $2 limit 1`, []*sql.NullString{&helpNoT}, idx, lang)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return ""
		}
	case constants.SwapType_Down: // 下层
		err := ss_sql.QueryRow(dbHandler, `select id from consultation_config where idx=$1+1 and is_delete = '0' and lang = $2 limit 1`, []*sql.NullString{&helpNoT}, idx, lang)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return ""
		}
	}

	return helpNoT.String
}

// 获取idx
func (ConsultationConfigDao) ExchangeIdx(consultationNoFrom, consultationNoTo string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, _ := dbHandler.BeginTx(context.TODO(), nil)
	defer ss_sql.Rollback(tx)

	err := ss_sql.ExecTx(tx, `update consultation_config set idx=idx+1 where id=$1`, consultationNoFrom)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}

	err = ss_sql.ExecTx(tx, `update consultation_config set idx=idx-1 where id=$1`, consultationNoTo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}

	tx.Commit()
	return ss_err.ERR_SUCCESS
}

//获取当前idx
func (ConsultationConfigDao) GetIdxById(tx *sql.Tx, id string) (idx int, lang, err string) {
	sqlStr := "select idx, lang from consultation_config where id =$1 and is_delete = '0' "
	var idxT, langT sql.NullString
	errT := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&idxT, &langT}, id)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return 0, "", ss_err.ERR_PARAM
	}

	return strext.ToInt(idxT.String), langT.String, ss_err.ERR_SUCCESS
}

//将当前的idx换成前一个idx(即idx-1)
func (ConsultationConfigDao) ReplaceIdx(tx *sql.Tx, idx int, lang string) (err string) {
	sqlStr := "update consultation_config set idx=$1-1 where idx=$1 and lang = $2 and is_delete = '0' "
	errT := ss_sql.ExecTx(tx, sqlStr, idx, lang)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return ss_err.ERR_SYS_DB_UPDATE
	}

	return ss_err.ERR_SUCCESS
}
