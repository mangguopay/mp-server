package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"context"
	"database/sql"
	"time"
)

type ServicerCountDao struct {
}

var ServicerCountDaoInst ServicerCountDao

// 更新账户数据
func (s *ServicerCountDao) UpdateCountData(data ServicerCheckListStatis) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//更新服务商的统计	servicer_count
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*120)
	defer cancel()

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		return errTx
	}

	// 1.更新数据
	sqlStr := "UPDATE servicer_count SET in_num=in_num+$1, in_amount=in_amount+$2, out_num=out_num+$3, out_amount=out_amount+$4, "
	sqlStr += " profit_num=profit_num+$5, profit_amount=profit_amount+$6, recharge_num=recharge_num+$7, "
	sqlStr += " recharge_amount=recharge_amount+$8, withdraw_num=withdraw_num+$9, withdraw_amount=withdraw_amount+$10, modify_time=current_timestamp "
	sqlStr += " WHERE servicer_no=$11 AND currency_type=$12 RETURNING servicer_no  "

	var servicerNo sql.NullString

	qErr := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&servicerNo},
		data.InNum, data.InAmount, data.OutNum, data.OutAmount, data.ProfitNum, data.ProfitAmount,
		data.RechargeNum, data.RechargeAmount, data.WithdrawNum, data.WithdrawAmount,
		data.ServicerNo, data.CurrencyType,
	)
	if qErr != nil {
		if qErr.Error() != ss_sql.DB_NO_ROWS_MSG {
			tx.Rollback()
			return qErr
		}
	}

	// 2.更新后没有返回数据，证明记录不存在， 进行插入一条记录
	if servicerNo.String == "" {
		if err := s.InsertRecordTx(tx, data); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 3.更新数据，标记已经处理完成
	if err := ServicerCheckListDaoInst.UpdateIsCountedTx(tx, data); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// 插入一条记录
func (s *ServicerCountDao) InsertRecordTx(tx *sql.Tx, data ServicerCheckListStatis) error {
	sqlStr := " INSERT INTO servicer_count (servicer_no, currency_type, in_num, in_amount, out_num, out_amount, "
	sqlStr += " profit_num, profit_amount, recharge_num, recharge_amount, withdraw_num, withdraw_amount, modify_time) "
	sqlStr += " VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, current_timestamp)  "

	eErr := ss_sql.ExecTx(tx, sqlStr,
		data.ServicerNo, data.CurrencyType, data.InNum, data.InAmount, data.OutNum, data.OutAmount,
		data.ProfitNum, data.ProfitAmount, data.RechargeNum, data.RechargeAmount,
		data.WithdrawNum, data.WithdrawAmount,
	)

	if eErr != nil {
		return eErr
	}

	return nil
}
