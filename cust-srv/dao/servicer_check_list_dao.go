package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type ServicerCheckListDao struct {
}

var ServicerCheckListDaoInst ServicerCheckListDao

const (
	// 记录数据是否已经统计过了
	ServicerCheckListDao_IsCounted_Ok = 1 // 是
	ServicerCheckListDao_IsCounted_No = 0 // 否
)

type ServicerCheckListStatis struct {
	ServicerNo     string
	CurrencyType   string
	Dates          string
	InNum          int64
	InAmount       int64
	OutNum         int64
	OutAmount      int64
	ProfitNum      int64
	ProfitAmount   int64
	RechargeNum    int64
	RechargeAmount int64
	WithdrawNum    int64
	WithdrawAmount int64
}

// 通过货币类型插入数据
func (ServicerCheckListDao) InsertByCurrency(data ServicerCheckListStatis) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "INSERT INTO servicer_count_list (servicer_no, currency_type, "
	sqlStr += " in_num, in_amount, out_num, out_amount, profit_num, profit_amount, "
	sqlStr += " recharge_num, recharge_amount, withdraw_num, withdraw_amount, id, dates, create_time) "
	sqlStr += " values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12, $13, $14, current_timestamp)"

	execErr := ss_sql.Exec(dbHandler,
		sqlStr, data.ServicerNo, data.CurrencyType,
		data.InNum, data.InAmount,
		data.OutNum, data.OutAmount,
		data.ProfitNum, data.ProfitAmount,
		data.RechargeNum, data.RechargeAmount,
		data.WithdrawNum, data.WithdrawAmount,
		strext.GetDailyId(), data.Dates,
	)

	if execErr != nil {
		ss_log.Error("ServicerCheckListDao |InsertByCurrency execErr=[%v]", execErr)
		return execErr
	}

	return nil
}

// 获取对应某天的pageSize条is_counted=0的记录
func (ServicerCheckListDao) GetCheckListStatis(date string, pageSize int) ([]ServicerCheckListStatis, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	list := []ServicerCheckListStatis{}

	sqlStr := "SELECT servicer_no, dates, currency_type, in_num, in_amount, "
	sqlStr += " out_num, out_amount, profit_num, profit_amount, "
	sqlStr += " recharge_num, recharge_amount, withdraw_num, withdraw_amount "
	sqlStr += " FROM servicer_count_list "
	sqlStr += " WHERE dates=$1 AND is_counted=$2 ORDER BY servicer_no ASC LIMIT $3 OFFSET 0 "

	rows, stmt, qErr := ss_sql.Query(dbHandler, sqlStr, date, ServicerCheckListDao_IsCounted_No, pageSize)
	if stmt != nil {
		defer stmt.Close()
		defer rows.Close()
	}

	if qErr != nil {
		if qErr.Error() != ss_sql.DB_NO_ROWS_MSG {
			return list, qErr
		}
		return list, nil
	}

	for rows.Next() {
		var sn, dates, currencyType, inNum, inAmount, outNum, outAmount sql.NullString
		var profitNum, profitAmount, rechargeNum, rechargeAmount, withdrawNum, withdrawAmount sql.NullString

		err := rows.Scan(&sn, &dates, &currencyType, &inNum, &inAmount, &outNum, &outAmount,
			&profitNum, &profitAmount, &rechargeNum, &rechargeAmount, &withdrawNum, &withdrawAmount,
		)
		if err != nil {
			ss_log.Error("err=%v", err)
			continue
		}
		list = append(list, ServicerCheckListStatis{
			ServicerNo:     sn.String,
			Dates:          dates.String,
			CurrencyType:   currencyType.String,
			InNum:          strext.ToInt64(inNum.String),
			InAmount:       strext.ToInt64(inAmount.String),
			OutNum:         strext.ToInt64(outNum.String),
			OutAmount:      strext.ToInt64(outAmount.String),
			ProfitNum:      strext.ToInt64(profitNum.String),
			ProfitAmount:   strext.ToInt64(profitAmount.String),
			RechargeNum:    strext.ToInt64(rechargeNum.String),
			RechargeAmount: strext.ToInt64(rechargeAmount.String),
			WithdrawNum:    strext.ToInt64(withdrawNum.String),
			WithdrawAmount: strext.ToInt64(withdrawAmount.String),
		})
	}

	return list, nil
}

// 更新记录已经统计完成
func (ServicerCheckListDao) UpdateIsCountedTx(tx *sql.Tx, data ServicerCheckListStatis) error {
	var isCounted sql.NullString

	sqlStr := "UPDATE servicer_count_list SET is_counted=$1 WHERE servicer_no=$2 AND currency_type=$3 AND dates=$4 AND is_counted=$5 RETURNING is_counted "

	// 此处不要修改为
	eErr := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&isCounted}, ServicerCheckListDao_IsCounted_Ok,
		data.ServicerNo, data.CurrencyType,
		data.Dates, ServicerCheckListDao_IsCounted_No,
	)

	if eErr != nil { // 如果记录已经更新过了， 重复执行的时候会返回 no rows in result set 错误，前面的操作会自动回滚
		return eErr
	}

	return nil
}
