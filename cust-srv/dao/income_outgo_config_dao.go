package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type IncomeOutgoConfigDao struct {
}

var IncomeOutgoConfigDaoInst IncomeOutgoConfigDao

// 根据idx获取id
func (IncomeOutgoConfigDao) GetIncomeOutgoConfigId(tx *sql.Tx, configType, idx string) (incomeOugoNo string) {

	var incomeOugoConfigNo sql.NullString
	sqlStr := "select income_ougo_config_no from income_ougo_config where idx=$1 and config_type=$2 limit 1 "
	err := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&incomeOugoConfigNo}, idx, configType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}
	return incomeOugoConfigNo.String
}

// 根据上下移获取上与下的idx
func (IncomeOutgoConfigDao) GetNearIdxIncomeOugoNo(tx *sql.Tx, configType, idx, swapType string) (funcNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var IncomeOugoNoT sql.NullString
	switch swapType {
	case constants.SwapType_Up: // 上层
		err := ss_sql.QueryRow(dbHandler, `select income_ougo_config_no from income_ougo_config where idx=$1-1 and config_type=$2 limit 1`, []*sql.NullString{&IncomeOugoNoT}, idx, configType)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return ""
		}
	case constants.SwapType_Down: // 下层
		err := ss_sql.QueryRow(dbHandler, `select income_ougo_config_no from income_ougo_config where idx=$1+1 and config_type=$2 limit 1`, []*sql.NullString{&IncomeOugoNoT}, idx, configType)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return ""
		}
	}

	return IncomeOugoNoT.String
}

//获取当前idx与类型（1.充值方式。2.提现方式）
func (IncomeOutgoConfigDao) GetIdxAndConfigTypeById(tx *sql.Tx, id string) (idx int, configType, err string) {
	sqlStr := "select idx,config_type from income_ougo_config  where income_ougo_config_no =$1 and is_delete = '0' "
	var idxT, configTypeT sql.NullString
	errT := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&idxT, &configTypeT}, id)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return 0, "", ss_err.ERR_PARAM
	}

	return strext.ToInt(idxT.String), configTypeT.String, ss_err.ERR_SUCCESS
}

//获取最大idx
func (IncomeOutgoConfigDao) GetMaxidx(tx *sql.Tx, configType string) (maxIdx int, err string) {
	sqlStr := "select max(idx) from income_ougo_config where config_type = $1 and is_delete='0'"
	var maxIdxT sql.NullString
	errT := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&maxIdxT}, configType)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return 0, ss_err.ERR_PARAM
	}
	return strext.ToInt(maxIdxT.String), ss_err.ERR_SUCCESS
}

//将当前的idx换成前一个idx(即idx-1)
func (IncomeOutgoConfigDao) ReplaceIdx(tx *sql.Tx, idx int, configType string) (err string) {
	sqlStr := "update income_ougo_config set idx=$1-1 where idx=$1 and config_type = $2 and is_delete = '0' "
	errT := ss_sql.ExecTx(tx, sqlStr, idx, configType)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return ss_err.ERR_SYS_DB_UPDATE
	}

	return ss_err.ERR_SUCCESS
}

// 交换
func (IncomeOutgoConfigDao) ExchangeIdx(tx *sql.Tx, incomeOutgoNoFrom, incomeOutgoNoTo string) (errCode string) {
	err := ss_sql.ExecTx(tx, `update income_ougo_config set idx=idx+1 where income_ougo_config_no=$1`, incomeOutgoNoFrom)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}

	err = ss_sql.ExecTx(tx, `update income_ougo_config set idx=idx-1 where income_ougo_config_no=$1`, incomeOutgoNoTo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}
