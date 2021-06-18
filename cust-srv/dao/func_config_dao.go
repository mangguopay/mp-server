package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type FuncConfigDao struct {
}

var FuncConfigDaoInst FuncConfigDao

func (FuncConfigDao) GetFuncList(applicationType int32) []*go_micro_srv_cust.FuncData {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var datas []*go_micro_srv_cust.FuncData

	sqlStr := "select func_no,func_name,img,jump_url from func_config where use_status='1' and is_delete='0' and application_type = $1 order by idx asc"
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, applicationType)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil
	}

	for rows.Next() {
		data := go_micro_srv_cust.FuncData{}
		err = rows.Scan(
			&data.Id,
			&data.Name,
			&data.Img,
			&data.Target,
		)

		imageBaseUrl := GlobalParamDaoInstance.QeuryParamValue("image_base_url")
		data.Img = imageBaseUrl + "/" + data.Img

		datas = append(datas, &data)
	}

	return datas
}

// 获取idx
func (FuncConfigDao) GetNearIdxFuncNo(appType, idx, swapType string) (funcNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var funcNoT sql.NullString
	switch swapType {
	case constants.SwapType_Up: // 上层
		err := ss_sql.QueryRow(dbHandler, `select func_no from func_config where idx=$1-1 and application_type=$2 limit 1`, []*sql.NullString{&funcNoT}, idx, appType)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return ""
		}
	case constants.SwapType_Down: // 下层
		err := ss_sql.QueryRow(dbHandler, `select func_no from func_config where idx=$1+1 and application_type=$2 limit 1`, []*sql.NullString{&funcNoT}, idx, appType)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return ""
		}
	}

	return funcNoT.String
}

// 获取FuncNo
func (FuncConfigDao) GetFuncNo(appType, idx string) (funcNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var funcNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select func_no from func_config where idx=$1 and application_type=$2 limit 1`, []*sql.NullString{&funcNoT}, idx, appType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return funcNoT.String
}

func (FuncConfigDao) GetImgURLFromNo(funcNO string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var imageURL sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select di.image_url from func_config fc LEFT JOIN dict_images di 
		ON fc.img_id = di.image_id  where fc.func_no = $1 and fc.is_delete = 0 limit 1`,
		[]*sql.NullString{&imageURL}, funcNO)
	if err != nil {
		return "", err
	}
	return imageURL.String, nil
}

// 获取idx
func (FuncConfigDao) ExchangeIdx(funcNoFrom, funcNoTo string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, `update func_config set idx=idx+1 where func_no=$1`, funcNoFrom)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}

	err = ss_sql.Exec(dbHandler, `update func_config set idx=idx-1 where func_no=$1`, funcNoTo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

//获取当前idx与类型（0-手机端;1-pos）
func (FuncConfigDao) GetIdxAndApplicationTypeById(tx *sql.Tx, id string) (idx int, applicationType, err string) {
	sqlStr := "select idx,application_type from func_config  where func_no =$1 and is_delete = '0' "
	var idxT, applicationTypeT sql.NullString
	errT := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&idxT, &applicationTypeT}, id)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return 0, "", ss_err.ERR_PARAM
	}

	return strext.ToInt(idxT.String), applicationTypeT.String, ss_err.ERR_SUCCESS
}

//将当前的idx换成前一个idx(即idx-1)
func (FuncConfigDao) ReplaceIdx(tx *sql.Tx, idx int, applicationType string) (err string) {
	sqlStr := "update func_config set idx=$1-1 where idx=$1 and application_type = $2 and is_delete = '0' "
	errT := ss_sql.ExecTx(tx, sqlStr, idx, applicationType)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return ss_err.ERR_SYS_DB_UPDATE
	}

	return ss_err.ERR_SUCCESS
}

//获取最大idx
func (FuncConfigDao) GetMaxidx(tx *sql.Tx, applicationType string) (maxIdx int, err string) {
	sqlStr := "select max(idx) from func_config where application_type = $1 and is_delete='0'"
	var maxIdxT sql.NullString
	errT := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&maxIdxT}, applicationType)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return 0, ss_err.ERR_PARAM
	}
	return strext.ToInt(maxIdxT.String), ss_err.ERR_SUCCESS
}

func (FuncConfigDao) DeleteFuncConfig(tx *sql.Tx, funcNo string) (errCode string) {

	err := ss_sql.ExecTx(tx, `update func_config set is_delete='1',idx='-1' where func_no=$1 and is_delete='0' `, funcNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

func (FuncConfigDao) ModifyUseStatusFuncConfig(funcNo, useStatus string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, `update func_config set use_status=$2 where func_no=$1`, funcNo, useStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

func (FuncConfigDao) AddFuncConfig(funcName, jumpUrl, img, imgId, applicationType string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var idx sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select max(idx)+1 from func_config where application_type=$1 and is_delete='0' `, []*sql.NullString{&idx}, applicationType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}

	funcNo := strext.NewUUID()
	sqlInsert := "insert into func_config(func_no, func_name, jump_url, img, img_id, application_type, use_status, is_delete, idx) " +
		" values($1,$2,$3,$4,$5,$6,'0',$7, $8)"
	err = ss_sql.Exec(dbHandler, sqlInsert, funcNo, funcName, jumpUrl, img, imgId, applicationType, "0", idx.String)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

func (FuncConfigDao) UpdateFuncConfig(tx *sql.Tx, funcNo, funcName, jumpUrl, img, imgId string) (errCode string) {
	sqlStr := " update func_config set func_name=$2, jump_url=$3, img=$4, img_id=$5 where func_no=$1"
	err := ss_sql.ExecTx(tx, sqlStr, funcNo, funcName, jumpUrl, img, imgId)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

/*
func (FuncConfigDao) GetFuncData(funcNo string) *go_micro_srv_cust.FuncConfigData {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select func_no, func_name, img, img_id, jump_url, idx, use_status from func_config where func_no=$1 and is_delete='0' limit 1"
	row, stmt, err := ss_sql.QueryRowN(dbHandler, sqlStr, funcNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil
	}
	defer stmt.Close()

	var data go_micro_srv_cust.FuncConfigData
	err = row.Scan(
		&data.FuncNo,
		&data.FuncName,
		&data.Img,
		&data.ImgId,
		&data.JumpUrl,
		&data.Idx,
		&data.UseStatus,
	)

	return &data
}*/
