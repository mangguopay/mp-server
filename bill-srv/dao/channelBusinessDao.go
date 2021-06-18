package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type ChannelBusinessDao struct{}

type ChannelBusinessConfigObj struct {
	ConfigId             string
	ChannelNo            string
	IsDelete             string
	CurrencyType         string
	SupportType          string
	UseStatus            string
	SaveRate             string
	SaveSingleMinFee     string
	SaveChargeType       string
	SaveMaxAmount        string
	WithdrawRate         string
	WithdrawMaxAmount    string
	WithdrawSingleMinFee string
	WithdrawChargeType   string
	CreateTime           string
}

var ChannelBusinessDaoInst ChannelBusinessDao

func (*ChannelBusinessDao) QueryChannelBusinessWithdrawInfoFromNo(channelNo, currencyType string) (string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	// 提现手续费率      单笔提现最大金额      提现单笔手续费          提现计算手续费类型(1-按比例收取手续费，2按单笔手续费收取)
	var withdrawRate, withdrawMaxAmount, withdrawSingleMinFee, withdrawChargeType sql.NullString
	sqlStr := `select withdraw_rate,withdraw_max_amount,withdraw_single_min_fee,withdraw_charge_type
				from channel_business_config 
				where channel_no = $1 
					and currency_type = $2 
					and is_delete='0' 
					and use_status = '1'
					and support_type in (2,3)
					limit 1`
	err := ss_sql.QueryRow(dbHandler, sqlStr,
		[]*sql.NullString{&withdrawRate, &withdrawMaxAmount, &withdrawSingleMinFee, &withdrawChargeType}, channelNo, currencyType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return "", "", "", ""
	}
	return withdrawRate.String, withdrawMaxAmount.String, withdrawSingleMinFee.String, withdrawChargeType.String
}
func (*ChannelBusinessDao) QueryChannelBusinessSaveInfoFromNo(channelNo, currencyType string) (string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	// 存款手续费率      单笔存款最大金额      存款单笔手续费          存款计算手续费类型(1-按比例收取手续费，2按单笔手续费收取)
	var saveRate, saveMaxAmount, saveSingleMinFee, saveChargeType sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select save_rate,save_max_amount,save_single_min_fee,save_charge_type from channel_business_config 
		where channel_no = $1 and currency_type = $2 and is_delete='0'  and use_status = '1' and support_type in (1,3) limit 1`,
		[]*sql.NullString{&saveRate, &saveMaxAmount, &saveSingleMinFee, &saveChargeType}, channelNo, currencyType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return "", "", "", ""
	}
	return saveRate.String, saveMaxAmount.String, saveSingleMinFee.String, saveChargeType.String
}
func (*ChannelBusinessDao) QueryCountFeeInfoFromNo(channelNo, currencyType string) (*ChannelBusinessConfigObj, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	// 存款手续费率      单笔存款最大金额      存款单笔手续费          存款计算手续费类型(1-按比例收取手续费，2按单笔手续费收取)
	var saveRate, withdrawRate, withdrawMaxAmount, saveSingleMinFee, withdrawSingleMinFee, saveChargeType, withdrawChargeType, supportType, saveMaxAmount sql.NullString
	sqlStr := `select save_rate,withdraw_rate,withdraw_max_amount,save_single_min_fee,
		withdraw_single_min_fee,save_charge_type,withdraw_charge_type,support_type,save_max_amount 
		from channel_business_config 
		where channel_no = $1 and currency_type = $2 and is_delete='0'  and use_status = '1' and support_type in (1,3) limit 1`
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&saveRate, &withdrawRate, &withdrawMaxAmount, &saveSingleMinFee,
		&withdrawSingleMinFee, &saveChargeType, &withdrawChargeType, &supportType, &saveMaxAmount}, channelNo, currencyType)
	if err != nil {
		return nil, err
	}

	obj := new(ChannelBusinessConfigObj)
	obj.SupportType = supportType.String
	obj.SaveRate = saveRate.String
	obj.SaveSingleMinFee = saveSingleMinFee.String
	obj.SaveMaxAmount = saveMaxAmount.String
	obj.SaveChargeType = saveChargeType.String
	obj.WithdrawRate = withdrawRate.String
	obj.WithdrawSingleMinFee = withdrawSingleMinFee.String
	obj.WithdrawMaxAmount = withdrawMaxAmount.String
	obj.WithdrawChargeType = withdrawChargeType.String

	return obj, nil
}
