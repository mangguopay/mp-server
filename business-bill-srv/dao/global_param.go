package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
	"strings"
)

type GlobalParamDao struct{}

var GlobalParamDaoInstance GlobalParamDao

func (*GlobalParamDao) GetParamValue(paramKey string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return "", GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var paramValueT sql.NullString
	sqlStr := `select param_value from global_param where param_key=$1   limit 1`
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&paramValueT}, paramKey)
	if nil != err {
		return "", err
	}
	return paramValueT.String, nil
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
		//商户USD转账最小/最大金额
		if paramKey.String == constants.GlobalParamKeyBusinessTransferUSDAmount {
			s := strings.Split(paramValue.String, "/")
			transferConf.USDMinAmount = strext.ToInt64(s[0])
			transferConf.USDMaxAmount = strext.ToInt64(s[1])
		}
		//商户KHR转账最小/最大金额
		if paramKey.String == constants.GlobalParamKeyBusinessTransferKHRAmount {
			s := strings.Split(paramValue.String, "/")
			transferConf.KHRMinAmount = strext.ToInt64(s[0])
			transferConf.KHRMaxAmount = strext.ToInt64(s[1])
		}
		//商户USD转账费率
		if paramKey.String == constants.GlobalParamKeyBusinessTransferUSDRate {
			transferConf.USDRate = strext.ToInt64(paramValue.String)
		}
		//商户KHR转账费率
		if paramKey.String == constants.GlobalParamKeyBusinessTransferKHRRate {
			transferConf.KHRRate = strext.ToInt64(paramValue.String)
		}
		//商户批量转账最大总人数
		if paramKey.String == constants.GlobalParamKeyBusinessTransferBatchNum {
			transferConf.BatchMaxNumBER = strext.ToInt64(paramValue.String)
		}

		//商家转账最低收取USD手续费
		if paramKey.String == constants.GlobalParamKeyBusinessTransferUSDMinFee {
			transferConf.USDMinFee = strext.ToInt64(paramValue.String)
		}
		//商家转账最低收取KHR手续费
		if paramKey.String == constants.GlobalParamKeyBusinessTransferKHRMinFee {
			transferConf.KHRMinFee = strext.ToInt64(paramValue.String)
		}

	}

	return transferConf, nil

}
