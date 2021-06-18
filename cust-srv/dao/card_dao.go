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
)

type CardDao struct {
}

var CardDaoInst CardDao

func (CardDao) InsertCard(accountNo, channelNo, name, cardNum, balanceType, isDefault, collectStatus, auditStatus, accountType, channelType string) (errCode string) {
	if isDefault == "" {
		isDefault = "0"
	}
	if collectStatus == "" {
		collectStatus = "0"
	}
	if auditStatus == "" {
		auditStatus = constants.AuditOrderStatus_Pending
	}
	if accountType == "" {
		ss_log.Error("添加的银行卡账号类型accountType为空")
		return ss_err.ERR_PARAM
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "insert into card(card_no,account_no,channel_no,name,create_time,is_delete,card_number,balance_type,is_defalut,collect_status,audit_status,account_type,channel_type) " +
		" values ($1,$2,$3,$4,current_timestamp,$5,$6,$7,$8,$9,$10,$11,$12)"
	err := ss_sql.Exec(dbHandler, sqlStr, strext.NewUUID(), accountNo, channelNo, name, "0", cardNum, balanceType, isDefault, collectStatus, auditStatus, accountType, channelType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}
	return ss_err.ERR_SUCCESS
}

//删除账号下的所有卡（用户卡，服务商卡）
func (CardDao) DeleteCard(tx *sql.Tx, accountNo string) (errCode error) {
	sqlStr := "update card set is_delete = '1' where account_no = $1 and is_delete = '0' "
	err := ss_sql.ExecTx(tx, sqlStr, accountNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

func (CardDao) GetServicerCards(whereList []*model.WhereSqlCond) (returnDatas []*go_micro_srv_cust.UserCardsData, returnTotal string, returnErr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	var total sql.NullString
	sqlCnt := "SELECT count(1) " +
		"FROM card ca " +
		"LEFT JOIN channel ch ON ch.channel_no = ca.channel_no " + whereModel.WhereStr
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, "", ss_err.ERR_PARAM
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY ca.is_defalut DESC ")

	sqlStr := "SELECT ca.card_no, ch.channel_name, ca.name, ca.card_number, ca.is_defalut, ca.balance_type, ch.logo_img_no, ch.logo_img_no_grey, ch.color_begin, ch.color_end, ca.account_type,ca.channel_no " +
		"FROM card ca " +
		"LEFT JOIN channel ch ON ch.channel_no = ca.channel_no " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, "0", ss_err.ERR_SYS_DB_GET
	}
	datas := []*go_micro_srv_cust.UserCardsData{}
	for rows.Next() {
		var data go_micro_srv_cust.UserCardsData
		var accountType, logoImgNoGrey, colorBegin, colorEnd sql.NullString
		err = rows.Scan(
			&data.CardNo,
			&data.ChannelName,
			&data.Name,
			&data.CardNumber,
			&data.IsDefalut,
			&data.BalanceType,
			&data.LogoImgNo,
			&logoImgNoGrey,
			&colorBegin,
			&colorEnd,
			&accountType,
			&data.ChannelNo,
		)

		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}

		data.LogoImgNoGrey = logoImgNoGrey.String
		data.ColorBegin = colorBegin.String
		data.ColorEnd = colorEnd.String

		if accountType.String == constants.AccountType_USER {
			langInt := len(data.CardNumber)
			if langInt > 4 {
				data.CardNumber = data.CardNumber[len(data.CardNumber)-4:]
			} else {
				continue
			}
		}

		datas = append(datas, &data)
	}

	return datas, total.String, ss_err.ERR_SUCCESS
}

func (CardDao) GetServicerCardDetail(cardNo string) (data *go_micro_srv_cust.UserCardsData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ca.card_no", Val: cardNo, EqType: "="},
		{Key: "ca.collect_status", Val: "1", EqType: "="},
		{Key: "ca.is_delete", Val: "0", EqType: "="},
	})

	sqlStr := "SELECT ca.card_no, ch.channel_name, ca.name, ca.card_number, ca.is_defalut, ca.balance_type, ch.logo_img_no, ch.logo_img_no_grey, ch.color_begin, ch.color_end, ca.account_type,ca.channel_no " +
		"FROM card ca " +
		"LEFT JOIN channel ch ON ch.channel_no = ca.channel_no " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	dataT := &go_micro_srv_cust.UserCardsData{}
	var accountType, logoImgNoGrey, colorBegin, colorEnd sql.NullString
	err2 = rows.Scan(
		&dataT.CardNo,
		&dataT.ChannelName,
		&dataT.Name,
		&dataT.CardNumber,
		&dataT.IsDefalut,
		&dataT.BalanceType,
		&dataT.LogoImgNo,
		&logoImgNoGrey,
		&colorBegin,
		&colorEnd,
		&accountType,
		&dataT.ChannelNo,
	)

	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return nil, err2
	}

	dataT.LogoImgNoGrey = logoImgNoGrey.String
	dataT.ColorBegin = colorBegin.String
	dataT.ColorEnd = colorEnd.String

	return dataT, nil

}

//用户查询自己的银行卡,银行卡卡号有做处理
func (CardDao) GetCustCards(whereList []*model.WhereSqlCond) (returnDatas []*go_micro_srv_cust.UserCardsData, returnTotal string, returnErr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	var total sql.NullString
	sqlCnt := "SELECT count(1) " +
		"FROM card ca " +
		"LEFT JOIN channel_cust_config chcc ON chcc.channel_no = ca.channel_no and chcc.currency_type = ca.balance_type and chcc.is_delete = '0' and chcc.use_status = '1' " +
		"LEFT JOIN channel ch ON ch.channel_no = ca.channel_no " + whereModel.WhereStr
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, "", ss_err.ERR_PARAM
	}

	if total.String == "" || total.String == "0" {
		return nil, "0", ss_err.ERR_SUCCESS
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY ca.is_defalut DESC, case chcc.use_status when "+constants.Status_Enable+" then 1 end ")

	sqlStr := "SELECT ca.card_no, ch.channel_name, ca.name, ca.card_number, ca.is_defalut," +
		" ca.balance_type, ch.logo_img_no, ch.logo_img_no_grey, ch.color_begin, ch.color_end," +
		" ca.account_type, ca.channel_no, chcc.use_status, ca.channel_type " +
		"FROM card ca " +
		"LEFT JOIN channel_cust_config chcc ON chcc.channel_no = ca.channel_no and chcc.currency_type = ca.balance_type and chcc.is_delete = '0'  " +
		"LEFT JOIN channel ch ON ch.channel_no = ca.channel_no " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, "0", ss_err.ERR_SYS_DB_GET
	}

	datas := []*go_micro_srv_cust.UserCardsData{}
	for rows.Next() {
		var data go_micro_srv_cust.UserCardsData
		var accountType, useStatus, logoImgNoGrey, colorBegin, colorEnd, channelType sql.NullString
		err = rows.Scan(
			&data.CardNo,
			&data.ChannelName,
			&data.Name,
			&data.CardNumber,
			&data.IsDefalut,
			&data.BalanceType,
			&data.LogoImgNo,
			&logoImgNoGrey,
			&colorBegin,
			&colorEnd,
			&accountType,
			&data.ChannelNo,
			&useStatus,
			&channelType,
		)

		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}

		data.LogoImgNoGrey = logoImgNoGrey.String
		data.ColorBegin = colorBegin.String
		data.ColorEnd = colorEnd.String
		data.ChannelType = channelType.String

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

func (CardDao) GetCustCardTotal(whereList []*model.WhereSqlCond) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		"FROM card ca " +
		"LEFT JOIN channel_cust_config chcc ON chcc.channel_no = ca.channel_no and chcc.currency_type = ca.balance_type and chcc.is_delete = '0' and chcc.use_status = '1' " +
		"LEFT JOIN channel ch ON ch.channel_no = ca.channel_no " + whereModel.WhereStr
	if err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereModel.Args...); err != nil {
		ss_log.Error("err=[%v]", err)
		return "0", err
	}

	return totalT.String, nil
}

//管理后台查询用户的银行卡
func (CardDao) WebAdminGetCustCards(whereList []*model.WhereSqlCond) ([]*go_micro_srv_cust.UserCardsData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY ca.is_defalut DESC,case chcc.use_status when "+constants.Status_Enable+" then 1 end ")
	//ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(pageSize), strext.ToInt(page))

	sqlStr := "SELECT ca.card_no, ch.channel_name, ca.name, ca.card_number, ca.is_defalut," +
		" ca.balance_type, ch.logo_img_no, ch.logo_img_no_grey, ch.color_begin, ch.color_end," +
		" ca.account_type, ca.channel_no, ca.create_time, chcc.use_status " +
		"FROM card ca " +
		"LEFT JOIN channel_cust_config chcc ON chcc.channel_no = ca.channel_no and chcc.currency_type = ca.balance_type and chcc.is_delete = '0'  " +
		"LEFT JOIN channel ch ON ch.channel_no = ca.channel_no " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	var datas []*go_micro_srv_cust.UserCardsData
	for rows.Next() {
		var data go_micro_srv_cust.UserCardsData
		var accountType, useStatus, logoImgNoGrey, colorBegin, colorEnd sql.NullString
		err2 = rows.Scan(
			&data.CardNo,
			&data.ChannelName,
			&data.Name,
			&data.CardNumber,
			&data.IsDefalut,
			&data.BalanceType,
			&data.LogoImgNo,
			&logoImgNoGrey,
			&colorBegin,
			&colorEnd,
			&accountType,
			&data.ChannelNo,
			&data.CreateTime,
			&useStatus,
		)

		if err2 != nil {
			ss_log.Error("err=[%v]", err2)
			continue
		}

		data.LogoImgNoGrey = logoImgNoGrey.String
		data.ColorBegin = colorBegin.String
		data.ColorEnd = colorEnd.String

		if useStatus.String != "" { //用户卡是否可提现的状态
			data.UseStatus = useStatus.String
		} else {
			data.UseStatus = constants.Status_Disable
		}

		datas = append(datas, &data)
	}

	return datas, nil
}

func (CardDao) GetCustCardDetail(cardNo string) (data *go_micro_srv_cust.UserCardsData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ca.card_no", Val: cardNo, EqType: "="},
		{Key: "ca.collect_status", Val: "1", EqType: "="},
		{Key: "ca.is_delete", Val: "0", EqType: "="},
	})

	sqlStr := "SELECT ca.card_no, ch.channel_name, ca.name, ca.card_number, ca.is_defalut, ca.balance_type" +
		", ch.logo_img_no, ch.logo_img_no_grey, ch.color_begin, ch.color_end, ca.account_type, ca.channel_no" +
		", chcc.use_status, chcc.channel_type " +
		" FROM card ca " +
		" LEFT JOIN channel_cust_config chcc ON chcc.channel_no = ca.channel_no and chcc.currency_type = ca.balance_type and chcc.is_delete = '0'  " +
		" LEFT JOIN channel ch ON ch.channel_no = ca.channel_no " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	dataT := &go_micro_srv_cust.UserCardsData{}
	var accountType, useStatus, logoImgNoGrey, colorBegin, colorEnd, channelType sql.NullString
	err2 = rows.Scan(
		&dataT.CardNo,
		&dataT.ChannelName,
		&dataT.Name,
		&dataT.CardNumber,
		&dataT.IsDefalut,
		&dataT.BalanceType,
		&dataT.LogoImgNo,
		&logoImgNoGrey,
		&colorBegin,
		&colorEnd,
		&accountType,
		&dataT.ChannelNo,
		&useStatus,
		&channelType,
	)

	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return nil, err2
	}

	dataT.LogoImgNoGrey = logoImgNoGrey.String
	dataT.ColorBegin = colorBegin.String
	dataT.ColorEnd = colorEnd.String
	dataT.ChannelType = channelType.String

	if useStatus.String != "" { //用户卡是否可提现的状态
		dataT.UseStatus = useStatus.String
	} else {
		dataT.UseStatus = constants.Status_Disable
	}

	return dataT, nil
}

//todo 修改用户或服务商查询卡号是否存在的逻辑 2020/5/7
func (*CardDao) QueryCardCnt(carNum, accountType string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var cnt sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT count(1) FROM card WHERE  card_number= $1 and account_type = $2 and is_delete = '0' ",
		[]*sql.NullString{&cnt}, carNum, accountType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "-1"
	}
	return cnt.String
}

//查询第三方渠道卡是否存在（现在可同一号码添加usd、khr两张渠道卡）
func (*CardDao) QueryThirdPartyCardCnt(carNum, accountType, channelNo, channelType, balanceType string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT count(1) " +
		" FROM card " +
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
