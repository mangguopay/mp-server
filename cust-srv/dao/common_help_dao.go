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
	"context"
	"database/sql"
)

type CommonHelpDao struct {
}

var CommonHelpDaoInst CommonHelpDao

func (CommonHelpDao) GetCnt(dbHandler *sql.DB, whereStr string, whereArgs []interface{}) (total string) {
	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM common_help ch " + whereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereArgs...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0"
	}
	return totalT.String
}

func (CommonHelpDao) GetCommonHelps(dbHandler *sql.DB, whereStr string, whereArgs []interface{}) (datas []*go_micro_srv_cust.CommonHelpData, err string) {
	sqlStr := " select ch.help_no, ch.problem, ch.answer, ch.use_status, ch.idx, ch.lang, ch.vs_type  " +
		" FROM common_help ch " + whereStr
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
		var data go_micro_srv_cust.CommonHelpData
		errT = rows.Scan(
			&data.HelpNo,
			&data.Problem,
			&data.Answer,
			&data.UseStatus,
			&data.Idx,
			&data.Lang,
			&data.VsType,
		)
		if errT != nil {
			ss_log.Error("errT=[%v]", errT)
			continue
		}

		datas = append(datas, &data)
	}
	return datas, ss_err.ERR_SUCCESS
}

func (CommonHelpDao) GetCommonHelpCount(dbHandler *sql.DB, lang, vsType string) (data *go_micro_srv_cust.CommonHelpCountData) {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
		{Key: "lang", Val: lang, EqType: "="},
		{Key: "vs_type", Val: vsType, EqType: "="},
	})

	//合计
	sumTotal := CommonHelpDaoInst.GetCnt(dbHandler, whereModel.WhereStr, whereModel.Args)

	//禁用合计
	disableM := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhere(disableM, "use_status", constants.Status_Disable, "=")
	disableTotal := CommonHelpDaoInst.GetCnt(dbHandler, disableM.WhereStr, disableM.Args)

	//启用合计
	enableM := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhere(enableM, "use_status", constants.Status_Enable, "=")
	enableMTotal := CommonHelpDaoInst.GetCnt(dbHandler, enableM.WhereStr, enableM.Args)

	dataT := &go_micro_srv_cust.CommonHelpCountData{
		SumTotal:     sumTotal,
		DisableTotal: disableTotal,
		Enable:       enableMTotal,
		Lang:         lang,
		VsType:       vsType,
	}

	return dataT
}

//app只看到问题
func (CommonHelpDao) GetAppCommonHelps(dbHandler *sql.DB, whereStr string, whereArgs []interface{}) (datas []*go_micro_srv_cust.CommonHelpData, err string) {
	sqlStr := " select help_no, problem, use_status, idx, lang, vs_type  " +
		" FROM common_help " + whereStr
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
		var data go_micro_srv_cust.CommonHelpData
		errT = rows.Scan(
			&data.HelpNo,
			&data.Problem,
			&data.UseStatus,
			&data.Idx,
			&data.Lang,
			&data.VsType,
		)
		if errT != nil {
			ss_log.Error("errT=[%v]", errT)
			continue
		}
		datas = append(datas, &data)
	}
	return datas, ss_err.ERR_SUCCESS
}

func (CommonHelpDao) GetCommonHelpDetail(dbHandler *sql.DB, whereStr string, whereArgs []interface{}) (data *go_micro_srv_cust.CommonHelpData, err string) {

	sqlStr := " select ch.help_no, ch.problem, ch.answer, ch.use_status, ch.lang, ch.vs_type  " +
		" FROM common_help ch " + whereStr

	rows, stmt, errT := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, ss_err.ERR_SYS_DB_GET
	}

	dataT := &go_micro_srv_cust.CommonHelpData{}
	for rows.Next() {
		errT = rows.Scan(
			&dataT.HelpNo,
			&dataT.Problem,
			&dataT.Answer,
			&dataT.UseStatus,
			&dataT.Lang,
			&dataT.VsType,
		)
		if errT != nil {
			ss_log.Error("errT=[%v]", errT)
			continue
		}
	}
	return dataT, ss_err.ERR_SUCCESS
}

//获取当前idx
func (CommonHelpDao) GetIdxById(tx *sql.Tx, id string) (idx int, lang, err string) {
	sqlStr := "select idx, lang from common_help where help_no =$1 and is_delete = '0' "
	var idxT, langT sql.NullString
	errT := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&idxT, &langT}, id)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return 0, "", ss_err.ERR_PARAM
	}

	return strext.ToInt(idxT.String), langT.String, ss_err.ERR_SUCCESS
}

//获取最大idx
func (CommonHelpDao) GetMaxidx(tx *sql.Tx, lang string) (maxIdx int, err string) {
	sqlStr := "select max(idx) from common_help where is_delete='0' and lang = $1 "
	var maxIdxT sql.NullString
	errT := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&maxIdxT}, lang)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return 0, ss_err.ERR_PARAM
	}
	return strext.ToInt(maxIdxT.String), ss_err.ERR_SUCCESS
}

//将当前的idx换成前一个idx(即idx-1)
func (CommonHelpDao) ReplaceIdx(tx *sql.Tx, idx int, lang string) (err string) {
	sqlStr := "update common_help set idx=$1-1 where idx=$1 and lang = $2 and is_delete = '0' "
	errT := ss_sql.ExecTx(tx, sqlStr, idx, lang)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return ss_err.ERR_SYS_DB_UPDATE
	}

	return ss_err.ERR_SUCCESS
}

func (CommonHelpDao) DeleteHelp(tx *sql.Tx, helpNo string) string {

	sqlStr := "update common_help set is_delete = '1',idx = '-1' where help_no = $1 and is_delete = '0' "
	err := ss_sql.ExecTx(tx, sqlStr, helpNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_UPDATE
	}
	return ss_err.ERR_SUCCESS

}

func (CommonHelpDao) ModifyHelpStatus(helpNo, useStatus string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update common_help set use_status = $2 where help_no = $1 "
	err := ss_sql.Exec(dbHandler, sqlStr, helpNo, useStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_UPDATE
	}
	return ss_err.ERR_SUCCESS

}

func (CommonHelpDao) UpdateHelp(helpNo, problem, answer string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update common_help set problem = $2, answer = $3, modify_time = current_timestamp  where help_no = $1 and is_delete = '0' "
	err := ss_sql.Exec(dbHandler, sqlStr, helpNo, problem, answer)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

func (d CommonHelpDao) AddHelp(problem, answer, lang, vsType string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//查询最大idx
	idx, errGet := d.GetHelpIdxMax(dbHandler, lang, vsType)

	if errGet == ss_err.ERR_SUCCESS {
		idx = strext.ToStringNoPoint(strext.ToInt32(idx) + 1)
	}

	sqlStr := "insert into common_help(help_no, problem, answer, idx, lang, vs_type, create_time, modify_time) values($1,$2,$3,$4,$5,$6,current_timestamp,current_timestamp)"
	err := ss_sql.Exec(dbHandler, sqlStr, strext.GetDailyId(), problem, answer, idx, lang, vsType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

func (CommonHelpDao) GetHelpIdxMax(dbHandler *sql.DB, lang, vsType string) (idx, rErr string) {
	//查询最大idx
	var idxMax sql.NullString
	sqlStr := "select idx from common_help where lang = $1 and vs_type = $2 and is_delete = '0' order by idx desc limit 1 "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&idxMax}, lang, vsType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "1", ss_err.ERR_SYS_DB_GET
	}
	return idxMax.String, ss_err.ERR_SUCCESS
}

// 获取helpNo
func (CommonHelpDao) GetHelpNo(idx string) (helpNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var helpNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select help_no from common_help where idx=$1 limit 1`, []*sql.NullString{&helpNoT}, idx)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return helpNoT.String
}

// 获取idx
func (CommonHelpDao) GetNearIdxHelpNo(idx, swapType, lang, vsType string) (helpNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var helpNoT sql.NullString
	switch swapType {
	case constants.SwapType_Up: // 上层
		err := ss_sql.QueryRow(dbHandler, `select help_no from common_help where idx=$1-1 and is_delete = '0' and lang = $2 and vs_type = $3 limit 1`, []*sql.NullString{&helpNoT}, idx, lang, vsType)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return ""
		}
	case constants.SwapType_Down: // 下层
		err := ss_sql.QueryRow(dbHandler, `select help_no from common_help where idx=$1+1 and is_delete = '0' and lang = $2 and vs_type = $3 limit 1`, []*sql.NullString{&helpNoT}, idx, lang, vsType)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return ""
		}
	}

	return helpNoT.String
}

// 获取idx
func (CommonHelpDao) ExchangeIdx(helpNoFrom, helpNoTo string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, _ := dbHandler.BeginTx(context.TODO(), nil)
	defer ss_sql.Rollback(tx)

	err := ss_sql.ExecTx(tx, `update common_help set idx=idx+1 where help_no=$1`, helpNoFrom)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}

	err = ss_sql.ExecTx(tx, `update common_help set idx=idx-1 where help_no=$1`, helpNoTo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}

	tx.Commit()
	return ss_err.ERR_SUCCESS
}
