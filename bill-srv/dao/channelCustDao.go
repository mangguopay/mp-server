package dao

import (
	"database/sql"

	"a.a/cu/ss_log"

	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type ChannelCustDao struct {
}

var ChannelCustDaoInst ChannelCustDao

func (*ChannelCustDao) QueryChannelCustWithdrawInfoFromNo(channelNo, currencyType string) (string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	// 提现手续费率      单笔提现最大金额      提现单笔手续费          提现计算手续费类型(1-按比例收取手续费，2按单笔手续费收取)
	var withdrawRate, withdrawMaxAmount, withdrawSingleMinFee, withdrawChargeType sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select withdraw_rate,withdraw_max_amount,withdraw_single_min_fee,withdraw_charge_type from channel_cust_config 
		where channel_no = $1 and currency_type = $2 and is_delete='0'  and use_status = '1' and support_type in (2,3) limit 1`,
		[]*sql.NullString{&withdrawRate, &withdrawMaxAmount, &withdrawSingleMinFee, &withdrawChargeType}, channelNo, currencyType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return "", "", "", ""
	}
	return withdrawRate.String, withdrawMaxAmount.String, withdrawSingleMinFee.String, withdrawChargeType.String
}
func (*ChannelCustDao) QueryChannelCustSaveInfoFromNo(channelNo, currencyType string) (string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	// 存款手续费率      单笔存款最大金额      存款单笔手续费          存款计算手续费类型(1-按比例收取手续费，2按单笔手续费收取)
	var saveRate, saveMaxAmount, saveSingleMinFee, saveChargeType sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select save_rate,save_max_amount,save_single_min_fee,save_charge_type from channel_cust_config 
		where channel_no = $1 and currency_type = $2 and is_delete='0'  and use_status = '1' and support_type in (1,3) limit 1`,
		[]*sql.NullString{&saveRate, &saveMaxAmount, &saveSingleMinFee, &saveChargeType}, channelNo, currencyType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return "", "", "", ""
	}
	return saveRate.String, saveMaxAmount.String, saveSingleMinFee.String, saveChargeType.String
}
func (*ChannelCustDao) QueryCountFeeInfoFromNo(channelNo, currencyType string) (string, string, string, string, string, string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	// 存款手续费率      单笔存款最大金额      存款单笔手续费          存款计算手续费类型(1-按比例收取手续费，2按单笔手续费收取)
	var saveRate, withdrawRate, withdrawMaxAmount, saveSingleMinFee, withdrawSingleMinFee, saveChargeType, withdrawChargeType, supportType, saveMaxAmount sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select save_rate,withdraw_rate,withdraw_max_amount,save_single_min_fee,
withdraw_single_min_fee,save_charge_type,withdraw_charge_type,support_type,save_max_amount from channel_cust_config 
		where channel_no = $1 and currency_type = $2 and is_delete='0'  and use_status = '1' and support_type in (1,3) limit 1`,
		[]*sql.NullString{&saveRate, &withdrawRate, &withdrawMaxAmount, &saveSingleMinFee, &withdrawSingleMinFee, &saveChargeType, &withdrawChargeType, &supportType, &saveMaxAmount}, channelNo, currencyType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return "", "", "", "", "", "", "", "", ""
	}
	return saveRate.String, withdrawRate.String, withdrawMaxAmount.String, saveSingleMinFee.String, withdrawSingleMinFee.String,
		saveChargeType.String, withdrawChargeType.String, supportType.String, saveMaxAmount.String
}
