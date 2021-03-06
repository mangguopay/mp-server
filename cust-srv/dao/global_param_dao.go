package dao

import (
	"database/sql"
	"strings"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type GlobalParamDao struct{}

var GlobalParamDaoInstance GlobalParamDao

func (*GlobalParamDao) QeuryParamValue(paramKey string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var paramValue sql.NullString
	sqlStr := `select param_value from global_param where param_key=$1   limit 1`
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&paramValue}, paramKey)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return paramValue.String
}

func (*GlobalParamDao) UpdateParamValueTx(tx *sql.Tx, paramValue, paramKey string) error {
	sqlStr := "update global_param set param_value = $1 where param_key = $2 "
	return ss_sql.ExecTx(tx, sqlStr, paramValue, paramKey)
}

func (*GlobalParamDao) InsertOrUpdateGlobalParam(paramKey, paramValue, remark string) (err string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlInsert := "insert into global_param(param_key,param_value,remark) values($1,$2,$3) " +
		" on conflict(param_key) do update set param_value=$2,remark=$3"

	errT := ss_sql.Exec(dbHandler, sqlInsert, paramKey, paramValue, remark)
	if errT != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

type BusinessTransferParamValue struct {
	USDMaxAmount   int64
	USDMinAmount   int64
	KHRMaxAmount   int64
	KHRMinAmount   int64
	BatchMaxNumBER int64
	USDRate        int64
	KHRRate        int64
	USDMinFee      int64
	KHRMinFee      int64
}

func (*GlobalParamDao) GetBusinessTransferParamValue() (*BusinessTransferParamValue, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `select param_key, param_value from global_param where param_key in($1, $2, $3, $4, $5, $6, $7) `
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr,
		constants.GlobalParamKeyBusinessTransferUSDAmount,
		constants.GlobalParamKeyBusinessTransferKHRAmount,
		constants.GlobalParamKeyBusinessTransferUSDRate,
		constants.GlobalParamKeyBusinessTransferKHRRate,
		constants.GlobalParamKeyBusinessTransferUSDMinFee,
		constants.GlobalParamKeyBusinessTransferKHRMinFee,
		constants.GlobalParamKeyBusinessTransferBatchNum,
	)
	if nil != err {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	transferConf := new(BusinessTransferParamValue)
	for rows.Next() {
		var paramKey, paramValue sql.NullString
		err := rows.Scan(&paramKey, &paramValue)
		if err != nil {
			return nil, err
		}
		//??????USD????????????/????????????
		if paramKey.String == constants.GlobalParamKeyBusinessTransferUSDAmount {
			s := strings.Split(paramValue.String, "/")
			transferConf.USDMinAmount = strext.ToInt64(s[0])
			transferConf.USDMaxAmount = strext.ToInt64(s[1])
		}
		//??????KHR????????????/????????????
		if paramKey.String == constants.GlobalParamKeyBusinessTransferKHRAmount {
			s := strings.Split(paramValue.String, "/")
			transferConf.KHRMinAmount = strext.ToInt64(s[0])
			transferConf.KHRMaxAmount = strext.ToInt64(s[1])
		}
		//??????USD????????????
		if paramKey.String == constants.GlobalParamKeyBusinessTransferUSDRate {
			transferConf.USDRate = strext.ToInt64(paramValue.String)
		}
		//??????KHR????????????
		if paramKey.String == constants.GlobalParamKeyBusinessTransferKHRRate {
			transferConf.KHRRate = strext.ToInt64(paramValue.String)
		}
		//?????????????????????????????????
		if paramKey.String == constants.GlobalParamKeyBusinessTransferBatchNum {
			transferConf.BatchMaxNumBER = strext.ToInt64(paramValue.String)
		}

		//????????????????????????USD?????????
		if paramKey.String == constants.GlobalParamKeyBusinessTransferUSDMinFee {
			transferConf.USDMinFee = strext.ToInt64(paramValue.String)
		}
		//????????????????????????KHR?????????
		if paramKey.String == constants.GlobalParamKeyBusinessTransferKHRMinFee {
			transferConf.KHRMinFee = strext.ToInt64(paramValue.String)
		}

	}

	return transferConf, nil

}
