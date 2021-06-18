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

type ChannelDao struct {
}

var ChannelDaoInst ChannelDao

type ChannelData struct {
	ChannelName string
	ChannelType string
}

type UseChannelData struct {
	ChannelNo            string
	CurrencyType         string
	SupportType          string
	SaveRate             string
	SaveSingleMinFee     string
	SaveMaxAmount        string
	SaveChargeType       string
	WithdrawRate         string
	WithdrawSingleMinFee string
	WithdrawMaxAmount    string
	WithdrawChargeType   string
	ChannelType          string
}

type BusinessChannelData struct {
	ChannelNo            string
	CurrencyType         string
	SupportType          string
	ChannelType          string
	SaveRate             string
	SaveSingleMinFee     string
	SaveMaxAmount        string
	SaveChargeType       string
	WithdrawRate         string
	WithdrawSingleMinFee string
	WithdrawMaxAmount    string
	WithdrawChargeType   string
}

func (ChannelDao) GetCnt(dbHandler *sql.DB, whereStr string, whereArgs []interface{}) (total string) {
	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		"FROM channel " + whereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereArgs...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0"
	}
	return totalT.String
}

func (ChannelDao) GetChannelList(dbHandler *sql.DB, whereStr string, whereArgs []interface{}) (datas []*go_micro_srv_cust.ChannelDataSimple) {
	sqlStr := "select channel_no, channel_name, is_recom, currency_type, logo_img_no from channel " + whereStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil
	}

	for rows.Next() {
		data := go_micro_srv_cust.ChannelDataSimple{}
		err = rows.Scan(
			&data.ChannelNo,
			&data.ChannelName,
			&data.IsRecom,
			&data.CurrencyType,
			&data.LogoImgNo,
		)
		datas = append(datas, &data)
	}

	return datas
}

func (ChannelDao) ModifyChannelStatus(tx *sql.Tx, channelNo string) error {
	sqlUpdate := "update channel set is_delete='1' where channel_no = $1 and is_delete = '0' "
	return ss_sql.ExecTx(tx, sqlUpdate, channelNo)
}

func (ChannelDao) AddChannel(tx *sql.Tx, channelName, logoImgNo, logoImgNoGrey, colorBegin, colorEnd, channelType string) (channelNoT string, err error) {
	channelNo := strext.NewUUID()
	sqlInsert := "insert into channel(channel_no, channel_name, logo_img_no, logo_img_no_grey, color_begin, color_end, channel_type, create_time) " +
		" values($1,$2,$3,$4,$5,$6,$7,current_timestamp)"
	insertErr := ss_sql.ExecTx(tx, sqlInsert, channelNo, channelName, logoImgNo, logoImgNoGrey, colorBegin, colorEnd, channelType)
	if insertErr != nil {
		ss_log.Error("insertErr=[%v]", insertErr)
		return channelNo, insertErr
	}

	return channelNo, nil
}

func (ChannelDao) UpdateChannel(tx *sql.Tx, channelNo, channelName, logoImgNo, logoImgNoGrey, colorBegin, colorEnd string) (err string) {
	sqlUpdate := "update channel set channel_name = $2, logo_img_no = $3, logo_img_no_grey = $4, color_begin = $5, color_end = $6 where channel_no = $1 "
	updateErr := ss_sql.ExecTx(tx, sqlUpdate, channelNo, channelName, logoImgNo, logoImgNoGrey, colorBegin, colorEnd)
	if updateErr != nil {
		ss_log.Error("updateErr=[%v]", updateErr)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

//由于发生错误，而需要将前面推荐渠道设为不推荐的回滚
func (ChannelDao) ModifyIsRecomBychannelNo(dbHandler *sql.DB, channelNo string) (err string) {
	sqlUpdate := "update channel set is_recom = '1' where channel_no = $1 "
	updateErr := ss_sql.Exec(dbHandler, sqlUpdate, channelNo)
	if updateErr != nil {
		ss_log.Error("updateErr=[%v]", updateErr)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

func (ChannelDao) GetChannelNameByChannelNo(channelNo string) (channelName string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select channel_name from channel where channel_no = $1 and is_delete = '0' "
	var channelNameT sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&channelNameT}, channelNo)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", errT
	}
	return channelNameT.String, nil
}

func (ChannelDao) GetChannelDetail(channelNo string) (*ChannelData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select channel_name, channel_type from channel where channel_no = $1 and is_delete = '0' "
	row, stmt, errT := ss_sql.QueryRowN(dbHandler, sqlStr, channelNo)
	if stmt != nil {
		defer stmt.Close()
	}
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return nil, errT
	}
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, errT
	}

	var channelName, channelType sql.NullString
	errT = row.Scan(
		&channelName,
		&channelType,
	)

	dataT := &ChannelData{
		ChannelName: channelName.String,
		ChannelType: channelType.String,
	}

	return dataT, nil
}

func (ChannelDao) GetLogoImageUrlByChannelNo(tx *sql.Tx, channelNo string) (string, error) {
	sqlUpdate := "select di.image_url from channel ch LEFT JOIN dict_images di ON ch.logo_img_no = di.image_id  " +
		"where ch.channel_no = $1 and ch.is_delete = 0 limit 1"
	var imageURL sql.NullString
	err := ss_sql.QueryRowTx(tx, sqlUpdate, []*sql.NullString{&imageURL}, channelNo)
	return imageURL.String, err
}

func (ChannelDao) GetLogoImageGreyUrlByChannelNo(tx *sql.Tx, channelNo string) (string, error) {
	sqlUpdate := "select di.image_url from channel ch LEFT JOIN dict_images di ON ch.logo_img_no_grey = di.image_id  " +
		"where ch.channel_no = $1 and ch.is_delete = 0 limit 1"
	var imageURL sql.NullString
	err := ss_sql.QueryRowTx(tx, sqlUpdate, []*sql.NullString{&imageURL}, channelNo)
	return imageURL.String, err
}

//=========================POS=============================

type ChannelServicerStruct struct {
	ChannelNo    string
	CreateTime   string
	Idx          string
	IsRecom      string
	CurrencyType string
	UseStatus    string
	Id           string
}

//将推荐渠道设置不推荐渠道（推荐渠道只能一个币种只有一个）
func (ChannelDao) ModifyPosChannelIsRecom(tx *sql.Tx, currencyType string) (returnErr string) {

	//修改推荐渠道为不推荐
	sqlUpdate := "update channel_servicer set is_recom = '0' where  currency_type = $1 and is_recom= '1' "
	err2 := ss_sql.ExecTx(tx, sqlUpdate, currencyType)
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return ss_err.ERR_PARAM
	}
	return ss_err.ERR_SUCCESS
}

func (ChannelDao) AddPosChannel(tx *sql.Tx, channelNo, isRecom, currencyType string) (id string, err error) {
	idT := strext.GetDailyId()
	sqlInsert := "insert into channel_servicer(id,channel_no, is_recom, currency_type, create_time) " +
		" values($1,$2,$3,$4,current_timestamp)"
	insertErr := ss_sql.ExecTx(tx, sqlInsert, idT, channelNo, isRecom, currencyType)
	if insertErr != nil {
		ss_log.Error("insertErr=[%v]", insertErr)
		return "", insertErr
	}

	return idT, nil
}

func (ChannelDao) GetChannelServicerDetailById(id string) (dataT ChannelServicerStruct, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select channel_no, create_time, idx, is_recom, currency_type, use_status, id " +
		" from channel_servicer where id = $1 and is_delete = '0' "
	row, stmt, errT := ss_sql.QueryRowN(dbHandler, sqlStr, id)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return dataT, errT
	}
	if stmt != nil {
		stmt.Close()
	}
	data := ChannelServicerStruct{}
	if err := row.Scan(
		&data.ChannelNo,
		&data.CreateTime,
		&data.Idx,
		&data.IsRecom,
		&data.CurrencyType,
		&data.UseStatus,
		&data.Id,
	); err != nil {
		ss_log.Error("err=[%v]", err)
		return dataT, err
	}

	return data, nil
}

func (ChannelDao) GetCurrencyByPosChannelId(id string) (err string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlUpdate := "select currency_type from channel_servicer where id = $1 and is_delete = '0' "
	var currencyType sql.NullString
	updateErr := ss_sql.QueryRow(dbHandler, sqlUpdate, []*sql.NullString{&currencyType}, id)
	if updateErr != nil {
		ss_log.Error("updateErr=[%v]", updateErr)
		return ""
	}

	return currencyType.String
}

func (ChannelDao) GetCurrencyByPosChannelIdTx(tx *sql.Tx, id string) (err string) {
	sqlUpdate := "select currency_type from channel_servicer where id = $1 and is_delete = '0' "
	var currencyType sql.NullString
	updateErr := ss_sql.QueryRowTx(tx, sqlUpdate, []*sql.NullString{&currencyType}, id)
	if updateErr != nil {
		ss_log.Error("updateErr=[%v]", updateErr)
		return ""
	}

	return currencyType.String
}

func (ChannelDao) CheckPosChannel(channelNo, currencyType string) (boolStr bool) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlUpdate := "select count(1) from channel_servicer where channel_no = $1 and currency_type = $2 and is_delete = '0' "
	var total sql.NullString
	updateErr := ss_sql.QueryRow(dbHandler, sqlUpdate, []*sql.NullString{&total}, channelNo, currencyType)
	if updateErr != nil {
		ss_log.Error("updateErr=[%v]", updateErr)
		return strext.ToInt(total.String) > 0
	}

	return strext.ToInt(total.String) > 0
}

//=========================USE=============================
func (ChannelDao) AddUseChannel(tx *sql.Tx, data UseChannelData) (id string, err error) {
	sqlInsert := "insert into channel_cust_config(id, channel_no, currency_type, support_type, save_rate" +
		", save_single_min_fee, save_max_amount, save_charge_type, withdraw_rate, withdraw_single_min_fee" +
		", withdraw_max_amount, withdraw_charge_type, channel_type, create_time) " +
		" values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,current_timestamp)"
	idT := strext.NewUUID()
	insertErr := ss_sql.ExecTx(tx, sqlInsert, idT, data.ChannelNo, data.CurrencyType, data.SupportType, data.SaveRate,
		data.SaveSingleMinFee, data.SaveMaxAmount, data.SaveChargeType, data.WithdrawRate, data.WithdrawSingleMinFee,
		data.WithdrawMaxAmount, data.WithdrawChargeType, data.ChannelType)
	if insertErr != nil {
		ss_log.Error("insertErr=[%v]", insertErr)
		return "", insertErr
	}

	return idT, nil
}

func (ChannelDao) CheckUseChannel(channelNo, currencyType string) (boolStr bool) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlUpdate := "select count(1) from channel_cust_config where channel_no = $1 and currency_type = $2 and is_delete = '0' "
	var total sql.NullString
	updateErr := ss_sql.QueryRow(dbHandler, sqlUpdate, []*sql.NullString{&total}, channelNo, currencyType)
	if updateErr != nil {
		ss_log.Error("updateErr=[%v]", updateErr)
		return strext.ToInt(total.String) > 0
	}

	return strext.ToInt(total.String) > 0
}

func (ChannelDao) GetUseChannelId(channelNo, currencyType string) (id string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlUpdate := "select id from channel_cust_config where channel_no = $1 and currency_type = $2 and is_delete = '0' "
	var idT sql.NullString
	updateErr := ss_sql.QueryRow(dbHandler, sqlUpdate, []*sql.NullString{&idT}, channelNo, currencyType)
	if updateErr != nil {
		ss_log.Error("updateErr=[%v]", updateErr)
		return ""
	}
	return idT.String
}

func (ChannelDao) ModifyUseChannel(tx *sql.Tx, id, supportType, saveRate, saveSingleMinFee, saveMaxAmount, saveChargeType, withdrawRate, withdrawSingleMinFee, withdrawMaxAmount, withdrawChargeType string) (err string) {
	sqlInsert := "update channel_cust_config set  support_type = $2, save_rate = $3" +
		", save_single_min_fee = $4, save_max_amount = $5, save_charge_type = $6, withdraw_rate = $7, withdraw_single_min_fee = $8" +
		", withdraw_max_amount = $9, withdraw_charge_type = $10 where id = $1 "
	insertErr := ss_sql.ExecTx(tx, sqlInsert, id, supportType, saveRate, saveSingleMinFee, saveMaxAmount, saveChargeType, withdrawRate, withdrawSingleMinFee, withdrawMaxAmount, withdrawChargeType)
	if insertErr != nil {
		ss_log.Error("insertErr=[%v]", insertErr)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

func (ChannelDao) GetCurrencyTypeByChannelId(id string) (currencyType string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select currency_type from channel_cust_config where id = $1 "
	var currencyTypeT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&currencyTypeT}, id)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}
	return currencyTypeT.String
}

func (ChannelDao) GetChannelCustConfigInfoById(id string) (channelName, currencyType string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := " select ch.channel_name, ccc.currency_type " +
		" from channel_cust_config ccc " +
		" left join channel ch on ch.channel_no = ccc.channel_no " +
		" where ccc.id = $1 and ccc.is_delete = '0' "
	var channelNameT, currencyTypeT sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&channelNameT, &currencyTypeT}, id)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", "", errT
	}
	return channelNameT.String, currencyTypeT.String, nil
}

//==========================business==========================
func (ChannelDao) GetBusinessChannelCnt(whereList []*model.WhereSqlCond) (cnt string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlStr := "select count(1) " +
		" from channel_business_config cbc " +
		" left join channel ch on ch.channel_no = cbc.channel_no " + whereModel.WhereStr
	var cntT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cntT}, whereModel.Args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}
	return cntT.String
}

func (ChannelDao) GetBusinessChannelList(whereList []*model.WhereSqlCond, page, pageSize int32) (datasT []*go_micro_srv_cust.BusinessChannelData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY cbc.create_time desc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimitI32(whereModel, pageSize, page)
	sqlStr := "SELECT  cbc.channel_no, cbc.create_time, cbc.use_status, cbc.currency_type, ch.logo_img_no, ch.logo_img_no_grey, ch.channel_name, cbc.id" +
		", cbc.save_rate, cbc.withdraw_rate, cbc.withdraw_max_amount, cbc.save_single_min_fee, cbc.withdraw_single_min_fee " +
		", cbc.save_charge_type, cbc.withdraw_charge_type, cbc.support_type, cbc.save_max_amount, cbc.channel_type " +
		" FROM channel_business_config cbc " +
		" left join channel ch on ch.channel_no = cbc.channel_no  " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	var datas []*go_micro_srv_cust.BusinessChannelData
	for rows.Next() {
		data := go_micro_srv_cust.BusinessChannelData{}
		var logoImgNo, logoImgNoGrey, channelType sql.NullString
		if err = rows.Scan(
			&data.ChannelNo,
			&data.CreateTime,
			&data.UseStatus,
			&data.CurrencyType,
			&logoImgNo,
			&logoImgNoGrey,
			&data.ChannelName,
			&data.Id,

			&data.SaveRate,
			&data.WithdrawRate,
			&data.WithdrawMaxAmount,
			&data.SaveSingleMinFee,
			&data.WithdrawSingleMinFee,

			&data.SaveChargeType,
			&data.WithdrawChargeType,
			&data.SupportType,
			&data.SaveMaxAmount,
			&channelType,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}

		data.ChannelType = channelType.String
		data.LogoImgNo = logoImgNo.String
		data.LogoImgNoGrey = logoImgNoGrey.String
		if logoImgNo.String != "" && logoImgNoGrey.String != "" {
			imgIds := []string{
				logoImgNo.String,
				logoImgNoGrey.String,
			}

			//获取图片url使前端可以展示图片
			imgUrls := ImageDaoInstance.GetImgUrlsByImgIds(imgIds)

			data.LogoImgUrl = imgUrls[0]
			data.LogoImgUrlGrey = imgUrls[1]
		}

		datas = append(datas, &data)
	}

	return datas, nil
}

func (ChannelDao) GetBusinessChannelDetail(id string) (dataT *go_micro_srv_cust.BusinessChannelData, errT error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "cbc.is_delete", Val: "0", EqType: "="},
		{Key: "cbc.id", Val: id, EqType: "="},
	})

	sqlStr := "SELECT  cbc.channel_no, cbc.create_time, cbc.use_status, cbc.currency_type, ch.logo_img_no, ch.logo_img_no_grey, ch.channel_name, cbc.id" +
		", cbc.save_rate, cbc.withdraw_rate, cbc.withdraw_max_amount, cbc.save_single_min_fee, cbc.withdraw_single_min_fee " +
		", cbc.save_charge_type, cbc.withdraw_charge_type, cbc.support_type, cbc.save_max_amount, cbc.channel_type " +
		" FROM channel_business_config cbc " +
		" left join channel ch on ch.channel_no = cbc.channel_no  " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	data := go_micro_srv_cust.BusinessChannelData{}
	var logoImgNo, logoImgNoGrey, channelType sql.NullString
	if err2 = rows.Scan(
		&data.ChannelNo,
		&data.CreateTime,
		&data.UseStatus,
		&data.CurrencyType,
		&logoImgNo,
		&logoImgNoGrey,
		&data.ChannelName,
		&data.Id,

		&data.SaveRate,
		&data.WithdrawRate,
		&data.WithdrawMaxAmount,
		&data.SaveSingleMinFee,
		&data.WithdrawSingleMinFee,

		&data.SaveChargeType,
		&data.WithdrawChargeType,
		&data.SupportType,
		&data.SaveMaxAmount,
		&channelType,
	); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return nil, err2
	}

	data.ChannelType = channelType.String
	data.LogoImgNo = logoImgNo.String
	data.LogoImgNoGrey = logoImgNoGrey.String
	if logoImgNo.String != "" && logoImgNoGrey.String != "" {
		imgIds := []string{
			logoImgNo.String,
			logoImgNoGrey.String,
		}

		//获取图片url使前端可以展示图片
		imgUrls := ImageDaoInstance.GetImgUrlsByImgIds(imgIds)

		data.LogoImgUrl = imgUrls[0]
		data.LogoImgUrlGrey = imgUrls[1]
	}

	return &data, nil
}

func (ChannelDao) CheckBusinessChannel(channelNo, currencyType string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlUpdate := "select count(1) from channel_business_config where channel_no = $1 and currency_type = $2 and is_delete = '0' "
	var total sql.NullString
	updateErr := ss_sql.QueryRow(dbHandler, sqlUpdate, []*sql.NullString{&total}, channelNo, currencyType)
	if updateErr != nil {
		ss_log.Error("updateErr=[%v]", updateErr)
		return strext.ToInt(total.String) > 0
	}

	return strext.ToInt(total.String) > 0
}

func (ChannelDao) AddBusinessChannel(tx *sql.Tx, data BusinessChannelData) (id string, err error) {
	ss_log.Info("data: --- %+v", data)
	sqlInsert := "insert into channel_business_config(id, channel_no, currency_type, support_type, save_rate" +
		", save_single_min_fee, save_max_amount, save_charge_type, withdraw_rate, withdraw_single_min_fee" +
		", withdraw_max_amount, withdraw_charge_type, channel_type, create_time) " +
		" values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,current_timestamp)"
	idT := strext.NewUUID()
	insertErr := ss_sql.ExecTx(tx, sqlInsert, idT, data.ChannelNo, data.CurrencyType, data.SupportType, data.SaveRate,
		data.SaveSingleMinFee, data.SaveMaxAmount, data.SaveChargeType, data.WithdrawRate, data.WithdrawSingleMinFee,
		data.WithdrawMaxAmount, data.WithdrawChargeType, data.ChannelType)
	if insertErr != nil {
		ss_log.Error("insertErr=[%v]", insertErr)
		return "", insertErr
	}

	return idT, nil
}

func (ChannelDao) GetBusinessChannelId(channelNo, currencyType string) (id string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlUpdate := "select id from channel_business_config where channel_no = $1 and currency_type = $2 and is_delete = '0' "
	var idT sql.NullString
	updateErr := ss_sql.QueryRow(dbHandler, sqlUpdate, []*sql.NullString{&idT}, channelNo, currencyType)
	if updateErr != nil {
		ss_log.Error("updateErr=[%v]", updateErr)
		return ""
	}
	return idT.String
}

func (ChannelDao) ModifyBusinessChannel(tx *sql.Tx, id, supportType, saveRate, saveSingleMinFee, saveMaxAmount, saveChargeType, withdrawRate, withdrawSingleMinFee, withdrawMaxAmount, withdrawChargeType string) (err string) {
	sqlInsert := "update channel_business_config set  support_type = $2, save_rate = $3" +
		", save_single_min_fee = $4, save_max_amount = $5, save_charge_type = $6, withdraw_rate = $7, withdraw_single_min_fee = $8" +
		", withdraw_max_amount = $9, withdraw_charge_type = $10 where id = $1 "
	insertErr := ss_sql.ExecTx(tx, sqlInsert, id, supportType, saveRate, saveSingleMinFee, saveMaxAmount, saveChargeType, withdrawRate, withdrawSingleMinFee, withdrawMaxAmount, withdrawChargeType)
	if insertErr != nil {
		ss_log.Error("insertErr=[%v]", insertErr)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

func (ChannelDao) DeleteBusinessChannelById(id string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlUpdate := "update channel_business_config set is_delete = '1' where id = $1 and is_delete = '0' "
	if updateErr := ss_sql.Exec(dbHandler, sqlUpdate, id); updateErr != nil {
		ss_log.Error("updateErr=[%v]", updateErr)
		return updateErr
	}
	return nil
}

func (ChannelDao) ModifyBusinessChannelStatusById(id, status string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlUpdate := "update channel_business_config set use_status = $2 where id = $1 "
	if updateErr := ss_sql.Exec(dbHandler, sqlUpdate, id, status); updateErr != nil {
		ss_log.Error("updateErr=[%v]", updateErr)
		return updateErr
	}
	return nil
}

func (ChannelDao) GetBusinessChannelInfoById(id string) (channelName, currencyType string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := " select ch.channel_name, ccc.currency_type " +
		" from channel_business_config ccc " +
		" left join channel ch on ch.channel_no = ccc.channel_no " +
		" where ccc.id = $1 and ccc.is_delete = '0' "
	var channelNameT, currencyTypeT sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&channelNameT, &currencyTypeT}, id)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", "", errT
	}
	return channelNameT.String, currencyTypeT.String, nil
}
