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
	"errors"
)

type CardBusinessDao struct {
	AccountNo     string
	ChannelNo     string
	Name          string
	CardNum       string
	BalanceType   string
	IsDefault     string
	CollectStatus string
	AuditStatus   string
	AccountType   string
	ChannelType   string
}

var CardBusinessDaoInst CardBusinessDao

func (CardBusinessDao) GetBusinessCards(accountNo, balanceType, accountType string) (returnDatas []*go_micro_srv_cust.BusinessCardData, returnTotal string, returnErr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ca.account_no", Val: accountNo, EqType: "="},
		{Key: "ca.collect_status", Val: "1", EqType: "="},
		{Key: "ca.is_delete", Val: "0", EqType: "="},
		{Key: "ca.balance_type", Val: balanceType, EqType: "="},
		{Key: "ca.account_type", Val: accountType, EqType: "="},
	})
	where := whereModel.WhereStr
	args := whereModel.Args
	var total sql.NullString
	sqlCnt := "SELECT count(1) " +
		"FROM card_business ca " +
		"LEFT JOIN channel_business_config cbc ON cbc.channel_no = ca.channel_no and cbc.currency_type = ca.balance_type and cbc.is_delete = '0' and cbc.use_status = '1' " +
		"LEFT JOIN channel ch ON ch.channel_no = ca.channel_no " + where
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, "", ss_err.ERR_PARAM
	}

	if total.String == "" || total.String == "0" {
		return nil, "0", ss_err.ERR_SUCCESS
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY ca.is_defalut DESC, case cbc.use_status when "+constants.Status_Enable+" then 1 end, ca.create_time desc ")
	sqlStr := `SELECT 
		ca.card_no, ca.name, ca.card_number, ca.is_defalut, ca.balance_type, ca.account_type,ca.channel_no,ca.create_time,
		ch.channel_name, ch.logo_img_no, ch.logo_img_no_grey, ch.color_begin, ch.color_end,
		cbc.use_status, cbc.withdraw_max_amount, cbc.withdraw_charge_type, cbc.withdraw_rate, cbc.withdraw_single_min_fee,
		cbc.save_max_amount, cbc.save_charge_type, cbc.save_rate, cbc.save_single_min_fee 
	FROM card_business ca
	LEFT JOIN channel_business_config cbc ON cbc.channel_no = ca.channel_no and cbc.currency_type = ca.balance_type and cbc.is_delete = '0'
	LEFT JOIN channel ch ON ch.channel_no = ca.channel_no ` + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, "0", ss_err.ERR_SYS_DB_GET
	}

	datas := []*go_micro_srv_cust.BusinessCardData{}
	for rows.Next() {
		var data go_micro_srv_cust.BusinessCardData
		var accountType, useStatus, logoImgNoGrey, colorBegin, colorEnd, withdrawMaxAmount,
			withdrawChargeType, withdrawRate, withdrawSingleMinFee, saveMaxAmount,
			saveChargeType, saveRate, saveSingleMinFee sql.NullString
		err = rows.Scan(
			&data.CardNo, &data.Name, &data.CardNumber, &data.IsDefalut, &data.BalanceType,
			&accountType, &data.ChannelNo, &data.CreateTime, &data.ChannelName, &data.LogoImgNo,
			&logoImgNoGrey, &colorBegin, &colorEnd, &useStatus, &withdrawMaxAmount,
			&withdrawChargeType, &withdrawRate, &withdrawSingleMinFee, &saveMaxAmount, &saveChargeType,
			&saveRate, &saveSingleMinFee,
		)

		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}

		data.LogoImgNoGrey = logoImgNoGrey.String
		data.ColorBegin = colorBegin.String
		data.ColorEnd = colorEnd.String
		data.WithdrawMaxAmount = withdrawMaxAmount.String
		data.WithdrawChargeType = withdrawChargeType.String
		data.WithdrawRate = withdrawRate.String
		data.WithdrawSingleMinFee = withdrawSingleMinFee.String
		data.SaveMaxAmount = saveMaxAmount.String
		data.SaveChargeType = saveChargeType.String
		data.SaveRate = saveRate.String
		data.SaveSingleMinFee = saveSingleMinFee.String

		if accountType.String == constants.AccountType_USER {
			langInt := len(data.CardNumber)
			if langInt > 4 {
				data.CardNumber = data.CardNumber[len(data.CardNumber)-4:]
			} else {
				continue
			}
		}

		if useStatus.String != "" { //用户卡是否可提现的状态
			data.UseStatus = useStatus.String
		} else {
			data.UseStatus = constants.Status_Disable
		}

		datas = append(datas, &data)
	}

	return datas, total.String, ss_err.ERR_SUCCESS
}

func (CardBusinessDao) GetBusinessCardDetail(cardNo string) (*go_micro_srv_cust.BusinessCardData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ca.card_no", Val: cardNo, EqType: "="},
		{Key: "ca.collect_status", Val: "1", EqType: "="},
		{Key: "ca.is_delete", Val: "0", EqType: "="},
	})

	sqlStr := "SELECT ca.card_no, ch.channel_name, ca.name, ca.card_number, ca.is_defalut, ca.balance_type," +
		"ch.logo_img_no, ch.logo_img_no_grey, ch.color_begin, ch.color_end, ca.account_type,ca.channel_no," +
		"cbc.use_status, cbc.withdraw_max_amount, cbc.withdraw_charge_type, cbc.withdraw_rate, cbc.withdraw_single_min_fee," +
		"cbc.save_max_amount, cbc.save_charge_type, cbc.save_rate, cbc.save_single_min_fee " +
		"FROM card_business ca " +
		"LEFT JOIN channel_business_config cbc ON cbc.channel_no = ca.channel_no and cbc.currency_type = ca.balance_type and cbc.is_delete = '0'  " +
		"LEFT JOIN channel ch ON ch.channel_no = ca.channel_no " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	dataT := &go_micro_srv_cust.BusinessCardData{}
	var accountType, useStatus, logoImgNoGrey, colorBegin, colorEnd, withdrawMaxAmount, withdrawChargeType,
		withdrawRate, withdrawSingleMinFee, saveMaxAmount, saveChargeType, saveRate, saveSingleMinFee sql.NullString
	err2 = rows.Scan(
		&dataT.CardNo, &dataT.ChannelName, &dataT.Name, &dataT.CardNumber, &dataT.IsDefalut,
		&dataT.BalanceType, &dataT.LogoImgNo, &logoImgNoGrey, &colorBegin, &colorEnd,
		&accountType, &dataT.ChannelNo, &useStatus, &withdrawMaxAmount, &withdrawChargeType,
		&withdrawRate, &withdrawSingleMinFee, &saveMaxAmount, &saveChargeType, &saveRate, &saveSingleMinFee,
	)
	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return nil, err2
	}

	dataT.WithdrawMaxAmount = withdrawMaxAmount.String
	dataT.WithdrawChargeType = withdrawChargeType.String
	dataT.WithdrawRate = withdrawRate.String
	dataT.WithdrawSingleMinFee = withdrawSingleMinFee.String
	dataT.SaveMaxAmount = saveMaxAmount.String
	dataT.SaveChargeType = saveChargeType.String
	dataT.SaveRate = saveRate.String
	dataT.SaveSingleMinFee = saveSingleMinFee.String
	dataT.LogoImgNoGrey = logoImgNoGrey.String
	dataT.ColorBegin = colorBegin.String
	dataT.ColorEnd = colorEnd.String

	if useStatus.String != "" { //用户卡是否可提现的状态
		dataT.UseStatus = useStatus.String
	} else {
		dataT.UseStatus = constants.Status_Disable
	}

	return dataT, nil
}

//func (CardBusinessDao) InsertCard(accountNo, channelNo, name, cardNum, balanceType, isDefault, collectStatus, auditStatus, accountType string) (errCode string) {
func (CardBusinessDao) InsertCard(data CardBusinessDao) (cardNumber string, err error) {
	if data.IsDefault == "" {
		data.IsDefault = "0"
	}
	if data.CollectStatus == "" {
		data.CollectStatus = "0"
	}
	if data.AuditStatus == "" {
		data.AuditStatus = constants.AuditOrderStatus_Pending
	}
	if data.AccountType == "" {
		ss_log.Error("添加的银行卡账号类型accountType为空")
		return "", errors.New("accountType参数为空")
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	cardNumberT := strext.NewUUID()
	sqlStr := "insert into card_business(card_no,account_no,channel_no,name,create_time,is_delete,card_number,balance_type,is_defalut,collect_status,audit_status,account_type,channel_type) " +
		" values ($1,$2,$3,$4,current_timestamp,$5,$6,$7,$8,$9,$10,$11,$12)"
	if err := ss_sql.Exec(dbHandler, sqlStr, cardNumberT, data.AccountNo, data.ChannelNo, data.Name, "0", data.CardNum, data.BalanceType, data.IsDefault,
		data.CollectStatus, data.AuditStatus, data.AccountType, data.ChannelType); err != nil {
		ss_log.Error("err=[%v]", err)
		return "", err
	}
	return cardNumberT, nil
}

//删除银行卡
func (CardBusinessDao) DeleteCard(cardNo string) (errCode error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update card_business set is_delete = '1' where card_no = $1 and is_delete = '0' "
	err := ss_sql.Exec(dbHandler, sqlStr, cardNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

//查询卡号是否存在
func (*CardBusinessDao) QueryCardNo(carNum, accountType string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var cardNOT sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT card_no FROM card_business WHERE  card_number= $1 and account_type = $2 and is_delete = '0' ",
		[]*sql.NullString{&cardNOT}, carNum, accountType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}
	return cardNOT.String
}

func (*CardBusinessDao) QueryCardCnt(carNum, accountType string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var cnt sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT count(1) FROM card_business WHERE  card_number= $1 and account_type = $2 and is_delete = '0' ",
		[]*sql.NullString{&cnt}, carNum, accountType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "-1"
	}
	return cnt.String
}

//查询第三方渠道卡是否存在
func (*CardBusinessDao) QueryThirdPartyCardCnt(carNum, accountType, channelNo, channelType, balanceType string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT count(1) " +
		" FROM card_business " +
		" WHERE card_number= $1 " +
		" AND account_type = $2 " +
		" AND channel_no = $3 " +
		" AND channel_type = $4 " +
		" AND balance_type = $5 " +
		" AND is_delete = '0' "
	var cnt sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cnt}, carNum, accountType, channelNo, channelType, balanceType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "-1"
	}
	return cnt.String
}
