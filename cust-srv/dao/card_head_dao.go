package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"context"
	"database/sql"
)

type CardHeadDao struct {
}

var CardHeadDaoInst CardHeadDao

func (*CardHeadDao) QueryCardNo(carNum, accountType string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var cardNOT sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT card_no FROM card_head WHERE  card_number= $1 and account_type = $2 and is_delete = '0' ",
		[]*sql.NullString{&cardNOT}, carNum, accountType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}
	return cardNOT.String
}

func (CardHeadDao) ModifyCollectStatus(setStatus, cardNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlUpdate := "update card_head set collect_status = $1 where card_no = $2 and is_delete = '0' "
	if err := ss_sql.Exec(dbHandler, sqlUpdate, setStatus, cardNo); err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_UPDATE
	}
	return ss_err.ERR_SUCCESS
}

func (*CardHeadDao) GetCardHeadAccountTypeByCardNo(cardNo string) (accountType string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select account_type from card_head where card_no = $1 and is_delete = '0' "
	var accountTypeT sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&accountTypeT}, cardNo)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", errT
	}
	return accountTypeT.String, nil
}

func (CardHeadDao) DeleteCard(cardNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "Update card_head set is_delete = '1' where card_no = $1 and is_delete='0' "
	if err := ss_sql.Exec(dbHandler, sqlStr, cardNo); err != nil {
		ss_log.Error("UpdateErr=[%v]", err)
		return ss_err.ERR_SYS_DB_DELETE
	}

	return ss_err.ERR_SUCCESS
}

func (CardHeadDao) InsertHeadSerCard(channelNo, name, cardNumber, note, balanceType, isDefalut, accountType string) (cardNo string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, errTx := dbHandler.BeginTx(context.TODO(), nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		return "", errTx
	}
	defer ss_sql.Rollback(tx)
	//查询平台总账号
	_, accPlat, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")
	if isDefalut == "1" { //只能有一张默认推荐收款卡
		//修改原来默认的卡为不默认推荐
		sqlStr2 := "update card_head set is_defalut = '0' where account_no = $1 and balance_type = $2 and account_type = $3 "
		err2 := ss_sql.ExecTx(tx, sqlStr2, accPlat, balanceType, accountType)
		if err2 != nil {
			ss_log.Error("修改默认卡失败。err2=[%v]", err2)
			return "", err2
		}
	}

	cardNoT := strext.NewUUID()
	sqlStr3 := "insert into card_head(card_no, account_no, channel_no, name, create_time, is_delete, card_number, balance_type, is_defalut, collect_status, audit_status, note, account_type) " +
		" values ($1,$2,$3,$4,current_timestamp,$5,$6,$7,$8,$9,$10,$11,$12)"
	err3 := ss_sql.ExecTx(tx, sqlStr3, cardNoT, accPlat, channelNo, name, "0", cardNumber, balanceType, isDefalut, "1", "0", note, accountType)
	if err3 != nil {
		ss_log.Error("添加卡失败，err3=[%v]", err3)
		return "", err3
	}
	tx.Commit()
	return cardNoT, nil
}

//
func (CardHeadDao) UpdateHeadSerCard(cardNo, channelNo, name, cardNumber, note, balanceType, isDefalut string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(context.TODO(), nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		return ss_err.ERR_SYS_DB_OP
	}
	defer ss_sql.Rollback(tx)

	//查询平台总账号
	_, accPlat, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")
	if isDefalut == "1" {
		//修改原来默认的卡为不默认推荐
		sqlStr2 := "update card_head set is_defalut = '0' where account_no = $1 and balance_type = $2 and is_defalut ='1' "
		err2 := ss_sql.ExecTx(tx, sqlStr2, accPlat, balanceType)
		if err2 != nil {
			ss_log.Error("err2=[%v]", err2)
			return ss_err.ERR_PARAM
		}
	}

	//设置默认卡
	sqlStr3 := "update card_head set channel_no = $2,name = $3,card_number = $4,note = $5,balance_type = $6,is_defalut = $7 where card_no = $1 "
	err3 := ss_sql.ExecTx(tx, sqlStr3, cardNo, channelNo, name, cardNumber, note, balanceType, isDefalut)
	if err3 != nil {
		ss_log.Error("err3=[%v]", err3)
		return ss_err.ERR_PARAM
	}

	tx.Commit()
	return ss_err.ERR_SUCCESS
}

//
func (CardHeadDao) UpdateHeadUseCard(cardNo, name, cardNumber, channelId string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//设置默认卡
	sqlStr3 := "update card_head set name = $2,card_number = $3,channel_cust_config_id= $4 where card_no = $1 "
	err3 := ss_sql.Exec(dbHandler, sqlStr3, cardNo, name, cardNumber, channelId)
	if err3 != nil {
		ss_log.Error("err3=[%v]", err3)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

//
func (CardHeadDao) UpdateHeadBusinessCard(cardNo, name, cardNumber, channelId string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//设置默认卡
	sqlStr3 := "update card_head set name = $2,card_number = $3,channel_business_config_id= $4 where card_no = $1 "
	err3 := ss_sql.Exec(dbHandler, sqlStr3, cardNo, name, cardNumber, channelId)
	if err3 != nil {
		ss_log.Error("err3=[%v]", err3)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

func (CardHeadDao) InsertHeadUseCard(name, cardNumber, accountType, channelId, balance_type string) (cardNo string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	//查询平台总账号
	_, accPlat, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")

	id := strext.NewUUID()
	sqlStr3 := "insert into card_head(card_no, account_no, name, create_time, is_delete, card_number, collect_status, audit_status, account_type, balance_type, channel_cust_config_id) " +
		" values ($1,$2,$3,current_timestamp,$4,$5,$6,$7,$8,$9,$10)"
	err3 := ss_sql.Exec(dbHandler, sqlStr3, id, accPlat, name, "0", cardNumber, "1", "0", accountType, balance_type, channelId)
	if err3 != nil {
		ss_log.Error("添加卡失败，err3=[%v]", err3)
		return "", err
	}
	return id, nil
}

func (CardHeadDao) GetHeadCardBusinessCnt(whereList []*model.WhereSqlCond) (cnt string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlStr := "select count(1) " +
		" FROM card_head ca " +
		" LEFT JOIN channel_business_config chcu ON chcu.id = ca.channel_business_config_id and chcu.is_delete='0' " +
		" LEFT JOIN channel ch ON ch.channel_no = ca.channel_no and ch.is_delete='0' " + whereModel.WhereStr
	var cntT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cntT}, whereModel.Args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}
	return cntT.String

}

//管理后台获取给商家使用的平台收款账户
func (CardHeadDao) GetHeadCardBusinessList(whereList []*model.WhereSqlCond, page, pageSize int) (datasT []*go_micro_srv_cust.CollectionManagementData, errT error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by ca.create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, pageSize, page)
	//卡的币种是从渠道那里来的
	sqlStr := "SELECT ch.channel_name,ca.name,ca.card_number,ca.collect_status,ca.card_no,ca.note,chcu.currency_type,ca.is_defalut " +
		" FROM card_head ca " +
		" LEFT JOIN channel_business_config chcu ON chcu.id = ca.channel_business_config_id and chcu.is_delete='0' " +
		" LEFT JOIN channel ch ON ch.channel_no = chcu.channel_no and ch.is_delete='0' " + whereModel.WhereStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	var datas []*go_micro_srv_cust.CollectionManagementData
	for rows.Next() {
		data := &go_micro_srv_cust.CollectionManagementData{}
		var channelName, balanceType sql.NullString
		err = rows.Scan(
			&channelName,
			&data.Name,
			&data.CardNumber,
			&data.CollectStatus,
			&data.CardNo,
			&data.Note,
			&balanceType,
			&data.IsDefalut,
		)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return datas, err

		}
		data.BalanceType = balanceType.String
		data.ChannelName = channelName.String
		datas = append(datas, data)
	}
	return datas, nil
}

//管理后台获取给商家使用的平台收款账户（单个）
func (CardHeadDao) GetHeadCardBusinessDatail(whereList []*model.WhereSqlCond) (dataT *go_micro_srv_cust.GetCardInfoData, errT error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlStr := "SELECT ch.channel_name, chcu.channel_no, ca.name, ca.card_number, ca.note," +
		" ca.collect_status, ca.is_defalut, chcu.currency_type, ca.channel_business_config_id, ca.card_no" +
		" FROM card_head ca " +
		" LEFT JOIN channel_business_config chcu ON chcu.id = ca.channel_business_config_id and chcu.is_delete = '0' " +
		" LEFT JOIN channel ch ON ch.channel_no = chcu.channel_no and ch.is_delete = '0' " + whereModel.WhereStr
	rows, stmt, err := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}

	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	data := &go_micro_srv_cust.GetCardInfoData{}
	err = rows.Scan(
		&data.ChannelName,
		&data.ChannelNo,
		&data.Name,
		&data.CardNumber,
		&data.Note,
		&data.CollectStatus,
		&data.IsDefalut,
		&data.BalanceType,
		&data.Id,
		&data.CardNo,
	)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err

	}
	return data, nil
}

//商家获取给商家使用的平台收款账户列表
func (CardHeadDao) GetHeadCardBusinessList2(whereList []*model.WhereSqlCond, page, pageSize int32) (datasT []*go_micro_srv_cust.HeadquartersCardsData, errT error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY ca.is_defalut DESC ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimitI32(whereModel, pageSize, page)

	sqlStr := "SELECT ca.card_no, ch.channel_name, ca.name, ca.card_number, chcu.currency_type, ch.logo_img_no, ch.logo_img_no_grey, ch.color_begin, ch.color_end " +
		", chcu.save_rate, chcu.withdraw_rate, chcu.withdraw_max_amount, chcu.save_single_min_fee, chcu.withdraw_single_min_fee " +
		", chcu.save_charge_type, chcu.withdraw_charge_type, chcu.support_type, chcu.save_max_amount, chcu.channel_no " +
		"FROM card_head ca " +
		"LEFT JOIN channel_business_config chcu ON chcu.id = ca.channel_business_config_id " +
		"LEFT JOIN channel ch ON ch.channel_no = chcu.channel_no " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	var datas []*go_micro_srv_cust.HeadquartersCardsData
	for rows.Next() {
		var channelName, logoImgNo, logoImgNoGrey, colorBegin, colorEnd sql.NullString
		data := &go_micro_srv_cust.HeadquartersCardsData{}
		err2 = rows.Scan(
			&data.CardNo,
			&channelName,
			&data.Name,
			&data.CardNumber,
			&data.MoneyType,

			&logoImgNo,
			&logoImgNoGrey,
			&colorBegin,
			&colorEnd,
			&data.SaveRate,
			&data.WithdrawRate,
			&data.WithdrawMaxAmount,
			&data.SaveSingleMinFee,

			&data.WithdrawSingleMinFee,
			&data.SaveChargeType,
			&data.WithdrawChargeType,
			&data.SupportType,
			&data.SaveMaxAmount,
			&data.ChannelNo,
		)
		if err2 != nil {
			ss_log.Error("err=[%v]", err2)
			return nil, err2
		}
		data.ChannelName = channelName.String
		data.ColorBegin = colorBegin.String
		data.ColorEnd = colorEnd.String

		if logoImgNo.String != "" || logoImgNoGrey.String != "" {
			imgIds := []string{
				logoImgNo.String,
				logoImgNoGrey.String,
			}

			imgUrls := ImageDaoInstance.GetImgUrlsByImgIds(imgIds)
			data.LogoImgUrl = imgUrls[0]
			data.LogoImgUrlGrey = imgUrls[1]
		}

		data.Temp = "0"

		datas = append(datas, data)
	}
	return datas, nil
}

//商家获取给商家使用的平台收款账户（单个）
func (CardHeadDao) GetHeadCardBusinessDetail(whereList []*model.WhereSqlCond) (dataT *go_micro_srv_cust.HeadquartersCardsData, errT error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlStr := "SELECT ca.card_no, ch.channel_name, ca.name, ca.card_number, chcu.currency_type, ch.logo_img_no, ch.logo_img_no_grey, ch.color_begin, ch.color_end " +
		", chcu.save_rate, chcu.withdraw_rate, chcu.withdraw_max_amount, chcu.save_single_min_fee, chcu.withdraw_single_min_fee " +
		", chcu.save_charge_type, chcu.withdraw_charge_type, chcu.support_type, chcu.save_max_amount, chcu.channel_no " +
		"FROM card_head ca " +
		"LEFT JOIN channel_business_config chcu ON chcu.id = ca.channel_business_config_id " +
		"LEFT JOIN channel ch ON ch.channel_no = chcu.channel_no " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	var channelName, logoImgNo, logoImgNoGrey, colorBegin, colorEnd sql.NullString
	data := &go_micro_srv_cust.HeadquartersCardsData{}
	err2 = rows.Scan(
		&data.CardNo,
		&channelName,
		&data.Name,
		&data.CardNumber,
		&data.MoneyType,

		&logoImgNo,
		&logoImgNoGrey,
		&colorBegin,
		&colorEnd,
		&data.SaveRate,
		&data.WithdrawRate,
		&data.WithdrawMaxAmount,
		&data.SaveSingleMinFee,

		&data.WithdrawSingleMinFee,
		&data.SaveChargeType,
		&data.WithdrawChargeType,
		&data.SupportType,
		&data.SaveMaxAmount,
		&data.ChannelNo,
	)
	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return nil, err2
	}
	data.ChannelName = channelName.String
	data.ColorBegin = colorBegin.String
	data.ColorEnd = colorEnd.String

	if logoImgNo.String != "" || logoImgNoGrey.String != "" {
		imgIds := []string{
			logoImgNo.String,
			logoImgNoGrey.String,
		}

		imgUrls := ImageDaoInstance.GetImgUrlsByImgIds(imgIds)
		data.LogoImgUrl = imgUrls[0]
		data.LogoImgUrlGrey = imgUrls[1]
	}

	data.Temp = "0"

	return data, nil
}

func (CardHeadDao) InsertHeadBusinessCard(name, cardNumber, accountType, channelId, balanceType string) (cardNo string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	//查询平台总账号
	_, accPlat, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")

	id := strext.NewUUID()
	sqlStr3 := "insert into card_head(card_no, account_no, name, create_time, is_delete, card_number, collect_status, audit_status, account_type, balance_type, channel_business_config_id) " +
		" values ($1,$2,$3,current_timestamp,$4,$5,$6,$7,$8,$9,$10)"
	err3 := ss_sql.Exec(dbHandler, sqlStr3, id, accPlat, name, "0", cardNumber, "1", "0", accountType, balanceType, channelId)
	if err3 != nil {
		ss_log.Error("添加卡失败，err3=[%v]", err3)
		return "", err
	}
	return id, nil
}
