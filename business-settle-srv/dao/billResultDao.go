package dao

import (
	"database/sql"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type BillResultDao struct{}

var BillResultDaoInst BillResultDao

func (BillResultDao) ModifyBillEnterResult(tx *sql.Tx, submitAmountSum, realAmountSum, enterFees, enterAmountSum, accNo, roleType string) string {
	// 确认是否存在,不存在则新增
	var countT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select count(1) from business_bill_results where acc_no=$1 and role_type = $2  limit 1`,
		[]*sql.NullString{&countT}, accNo, roleType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PARAM
	}
	if strext.ToInt(countT.String) > 0 { // 修改
		err := ss_sql.ExecTx(tx, `update business_bill_results set submit_amount_sum = submit_amount_sum + $1, 
					real_amount_sum = real_amount_sum + $2, enter_fees = enter_fees +$3,enter_amount_sum = enter_amount_sum + $4,modify_time = current_timestamp where acc_no=$5 and role_type = $6 `,
			submitAmountSum, realAmountSum, enterFees, enterAmountSum, accNo, roleType)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return ss_err.ERR_PAY_UPDATE_ORDER
		}
		return ss_err.ERR_SUCCESS
	}

	// 新增
	err = ss_sql.ExecTx(tx, `insert into business_bill_results(id,submit_amount_sum,real_amount_sum,enter_fees,enter_amount_sum,acc_no,role_type,modify_time) `+
		`values($1,$2,$3,$4,$5,$6,$7,current_timestamp)`,
		strext.GetDailyId(), submitAmountSum, realAmountSum, enterFees, enterAmountSum, accNo, roleType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_ADD
	}

	return ss_err.ERR_SUCCESS
}
func (BillResultDao) ModifyBillOutGoResult(tx *sql.Tx, outGoAmount, outGoFee, outGoRealAmount, accNo, roleType string) string {
	// 确认是否存在,不存在则新增
	var countT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select count(1) from business_bill_results where acc_no=$1 and role_type = $2  limit 1`,
		[]*sql.NullString{&countT}, accNo, roleType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PARAM
	}
	if strext.ToInt(countT.String) > 0 { // 修改
		err := ss_sql.ExecTx(tx, `update business_bill_results set out_go_amount = out_go_amount + $1, 
					out_go_fee = out_go_fee + $2, out_go_real_amount = out_go_real_amount +$3,modify_time = current_timestamp  where acc_no=$4 and role_type = $5 `,
			outGoAmount, outGoFee, outGoRealAmount, accNo, roleType)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return ss_err.ERR_PAY_UPDATE_ORDER
		}
		return ss_err.ERR_SUCCESS
	}

	// 新增
	err = ss_sql.ExecTx(tx, `insert into business_bill_results(id, out_go_amount, out_go_fee, out_go_real_amount,acc_no,role_type,modify_time) `+
		`values($1,$2,$3,$4,$5,$6,current_timestamp)`,
		strext.GetDailyId(), outGoAmount, outGoFee, outGoRealAmount, accNo, roleType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_ADD
	}

	return ss_err.ERR_SUCCESS
}
