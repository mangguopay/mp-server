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

type AppVersionDao struct {
}

var AppVersionDaoInst AppVersionDao

func (AppVersionDao) GetCnt(dbHandler *sql.DB, whereStr string, whereArgs []interface{}) (total string) {
	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM app_version ap " + whereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereArgs...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0"
	}
	return totalT.String
}

func (AppVersionDao) CheckHaveVersion(system, vsType string) (boolT bool) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ap.is_delete", Val: "0", EqType: "="},
		{Key: "ap.system", Val: system, EqType: "="},
		{Key: "ap.vs_type", Val: vsType, EqType: "="},
	})

	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM app_version ap " + whereModel.WhereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereModel.Args...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return strext.ToInt(totalT.String) > 0
	}
	return strext.ToInt(totalT.String) > 0
}

//获取最新包和最后修改时间
func (AppVersionDao) GetNewVersion(system, vsType string) (version, updateTime, vsCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ap.is_delete", Val: "0", EqType: "="},
		//{Key: "ap.status", Val: constants.Status_Enable, EqType: "="},
		{Key: "ap.system", Val: system, EqType: "="},
		{Key: "ap.vs_type", Val: vsType, EqType: "="},
	})

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY ap.create_time desc, ap.vs_code desc LIMIT 1 ")
	var versionT, updateTimeT, vsCodeT sql.NullString
	sqlCnt := "SELECT ap.version, ap.update_time, ap.vs_code " +
		" FROM app_version ap " + whereModel.WhereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&versionT, &updateTimeT, &vsCodeT}, whereModel.Args...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "", "", ""
	}
	return versionT.String, updateTimeT.String, vsCodeT.String
}

//获取统计信息与最新版本
func (AppVersionDao) GetVersionCount(dbHandler *sql.DB, whereModel *model.WhereSql, system, vsType string) (versionCountData *go_micro_srv_cust.GetAppVersionsCountData) {
	tempM := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhere(tempM, "ap.system", system, "=")
	ss_sql.SsSqlFactoryInst.AppendWhere(tempM, "ap.vs_type", vsType, "=")

	//合计
	sumTotal := AppVersionDaoInst.GetCnt(dbHandler, tempM.WhereStr, tempM.Args)

	//禁用合计
	disableM := ss_sql.SsSqlFactoryInst.DeepClone(tempM)
	ss_sql.SsSqlFactoryInst.AppendWhere(disableM, "ap.status", constants.Status_Disable, "=")
	disableTotal := AppVersionDaoInst.GetCnt(dbHandler, disableM.WhereStr, disableM.Args)

	//启用合计
	enableM := ss_sql.SsSqlFactoryInst.DeepClone(tempM)
	ss_sql.SsSqlFactoryInst.AppendWhere(enableM, "ap.status", constants.Status_Enable, "=")
	enableMTotal := AppVersionDaoInst.GetCnt(dbHandler, enableM.WhereStr, enableM.Args)

	//最新包的版本version、最后修改时间
	newVersion, updateTime, _ := AppVersionDaoInst.GetNewVersion(system, vsType)

	versionCountDataT := &go_micro_srv_cust.GetAppVersionsCountData{
		System:       system,       //系统（ios,android）
		VsType:       vsType,       //包类型（0-app;1-pos）
		SumTotal:     sumTotal,     //合计
		DisableTotal: disableTotal, //禁用合计
		Enable:       enableMTotal, //禁用合计
		NewVersion:   newVersion,   //最新版本
		UpdateTime:   updateTime,   //更新时间
	}
	return versionCountDataT
}

func (AppVersionDao) GetAppVersions(dbHandler *sql.DB, whereStr string, whereArgs []interface{}) (datas []*go_micro_srv_cust.AppVersionData, err string) {

	sqlStr := " select ap.v_id, ap.description, ap.version, ap.create_time, ap.update_time, ap.app_url" +
		", ap.vs_code, ap.vs_type, ap.is_force, ap.system, ap.note, ap.status " +
		", acc.account " +
		" FROM app_version ap " +
		" left join admin_account acc on acc.uid = ap.account_uid " + whereStr
	rows, stmt, errT := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, ss_err.ERR_SYS_DB_GET
	}
	appBaseUrl := GlobalParamDaoInstance.QeuryParamValue("app_base_url")

	for rows.Next() {
		var data go_micro_srv_cust.AppVersionData
		errT = rows.Scan(
			&data.VId,
			&data.Description,
			&data.Version,
			&data.CreateTime,
			&data.UpdateTime,

			&data.AppUrl,
			&data.VsCode,
			&data.VsType,
			&data.IsForce,
			&data.System,

			&data.Note,
			&data.Status,
			&data.Account,
		)
		if errT != nil {
			ss_log.Error("errT=[%v]", errT)
			continue
		}
		data.AppUrl = appBaseUrl + "/" + data.AppUrl

		datas = append(datas, &data)
	}
	return datas, ss_err.ERR_SUCCESS
}

func (AppVersionDao) GetAppVersionDetail(vid string) (data *go_micro_srv_cust.AppVersionData, err string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ap.is_delete", Val: "0", EqType: "="},
		{Key: "ap.v_id", Val: vid, EqType: "="},
	})

	sqlStr := " select ap.v_id, ap.description, ap.version, ap.create_time, ap.update_time, ap.app_url" +
		", ap.vs_code, ap.vs_type, ap.is_force, ap.system, ap.note, ap.status " +
		", acc.account " +
		" FROM app_version ap " +
		" left join admin_account acc on acc.uid = ap.account_uid " + whereModel.WhereStr
	rows, stmt, errT := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, ss_err.ERR_SYS_DB_GET
	}

	dataT := &go_micro_srv_cust.AppVersionData{}
	errT = rows.Scan(
		&dataT.VId,
		&dataT.Description,
		&dataT.Version,
		&dataT.CreateTime,
		&dataT.UpdateTime,

		&dataT.AppUrl,
		&dataT.VsCode,
		&dataT.VsType,
		&dataT.IsForce,
		&dataT.System,

		&dataT.Note,
		&dataT.Status,
		&dataT.Account,
	)
	appBaseUrl := GlobalParamDaoInstance.QeuryParamValue("app_base_url")
	dataT.AppUrl = appBaseUrl + "/" + dataT.AppUrl
	return dataT, ss_err.ERR_SUCCESS
}

func (AppVersionDao) ModifyAppVersionStatus(vid, status string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlStr := " update app_version set status = $2 where v_id = $1 "
	err := ss_sql.Exec(dbHandler, sqlStr, vid, status)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_UPDATE
	}
	return ss_err.ERR_SUCCESS
}

func (AppVersionDao) ModifyIsForce(vid, isForce string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := " update app_version set is_force = $2 where v_id = $1 "
	if err := ss_sql.Exec(dbHandler, sqlStr, vid, isForce); err != nil {
		ss_log.Error("err = [%v]", err)
		return ss_err.ERR_SYS_DB_UPDATE
	}
	return ss_err.ERR_SUCCESS
}

func (AppVersionDao) InsertAppVersion(tx *sql.Tx, description, versionStr, appUrl, vsCodeStr, vsType, system, accountNo, note, status string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "insert into app_version(v_id, description, version, create_time, update_time, app_url " +
		", vs_code, vs_type, system, account_uid, note, status) " +
		" values($1,$2,$3,current_timestamp,current_timestamp,$4,$5,$6,$7,$8,$9,$10) "

	if err := ss_sql.ExecTx(tx, sqlStr, strext.GetDailyId(), description, versionStr, appUrl, vsCodeStr, vsType, system, accountNo, note, status); err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_ADD
	}
	return ss_err.ERR_SUCCESS
}

func (AppVersionDao) ModifyAppVersionInfo(vId, description, note, status string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := " update app_version set description = $2, note = $3, status = $4 where v_id = $1 "
	if err := ss_sql.Exec(dbHandler, sqlStr, vId, description, note, status); err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_UPDATE
	}
	return ss_err.ERR_SUCCESS
}
