package handler

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_func"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/dao"
	"context"
	"database/sql"
)

/**
 * 获取功能列表
 */
func (*CustHandler) GetFuncList(ctx context.Context, req *go_micro_srv_cust.GetFuncListRequest, reply *go_micro_srv_cust.GetFuncListReply) error {
	datas := dao.FuncConfigDaoInst.GetFuncList(req.ApplicationType)

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	return nil
}

func (c *CustHandler) GetExchangeRate(ctx context.Context, req *go_micro_srv_cust.GetExchangeRateRequest, reply *go_micro_srv_cust.GetExchangeRateReply) error {
	_, usdToKhr, _ := cache.ApiDaoInstance.GetGlobalParam("usd_to_khr")
	_, khrToUsd, _ := cache.ApiDaoInstance.GetGlobalParam("khr_to_usd")
	_, usdToKhrFee, _ := cache.ApiDaoInstance.GetGlobalParam("usd_to_khr_fee")
	_, khrToUsdFee, _ := cache.ApiDaoInstance.GetGlobalParam("khr_to_usd_fee")
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = &go_micro_srv_cust.ExchangeRateData{
		UsdToKhr:    usdToKhr,
		KhrToUsd:    khrToUsd,
		UsdToKhrFee: usdToKhrFee,
		KhrToUsdFee: khrToUsdFee,
	}
	return nil
}

// 插入新卡信息
func (*CustHandler) InsertOrUpdateCard(ctx context.Context, req *go_micro_srv_cust.InsertOrUpdateCardRequest, reply *go_micro_srv_cust.InsertOrUpdateCardReply) error {

	switch req.AccountType { //只有登陆的是用户和服务商才有权限调用该接口
	case constants.AccountType_USER:
	case constants.AccountType_SERVICER:
	default:
		ss_log.Error("AccountType类型错误:[%v]", req.AccountType)
		reply.ResultCode = ss_err.ERR_SYS_NO_API_AUTH
		return nil
	}

	// 获取实名制姓名,报错或者为空说明未实名制
	authName, err := dao.AccDaoInstance.GetAuthNameFromUid(req.AccountNo)
	if err != nil {
		ss_log.Error("InsertOrUpdateCard 查询用户实名认证的姓名失败,uid为: %s,err: %s", req.AccountNo, err.Error())
		reply.ResultCode = ss_err.ERR_ACCOUNT_NOT_REAL_AUTH
		return nil
	}

	// 实名认证的姓名要和持卡人的姓名一致
	if authName != req.CardAccName {
		ss_log.Error("用户实名制的姓名和银行卡姓名不一致,实名制姓名为: %s,需要添加银行卡的姓名为: %s", authName, req.CardAccName)
		reply.ResultCode = ss_err.ERR_ACCOUNT_REAL_NAME_NOT_SAME
		return nil
	}

	//获取渠道类型
	channelData, getChannelDataErr := dao.ChannelDaoInst.GetChannelDetail(req.ChannelNo)
	if getChannelDataErr != nil {
		ss_log.Error("err=[%v]", getChannelDataErr)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	switch channelData.ChannelType {
	case constants.CHANNELTYPE_ORDINARY: //普通渠道
		//普通渠道卡号长度参数校验
		cardLength := len(req.CardNumber)
		if cardLength > constants.CardNumberLengthMax || cardLength < constants.CardNumberLengthMin {
			ss_log.Error("普通渠道银行卡长度必需在%v到%v之间, CardNumber长度[%v]", constants.CardNumberLengthMin, constants.CardNumberLengthMax, cardLength)
			reply.ResultCode = ss_err.ERR_CardNumberLength_FAILD
			return nil
		}

		//  判断卡是否存在
		if cnt := dao.CardDaoInst.QueryCardCnt(req.CardNumber, req.AccountType); cnt != "0" {
			ss_log.Error("err=[普通渠道银行卡号已存在,卡号为--->%s]", req.CardNumber)
			reply.ResultCode = ss_err.ERR_CARD_IS_EXIST
			return nil
		}
	case constants.CHANNELTYPE_THIRDPARTY: //第三方渠道
		//第三方的渠道卡号都是柬埔寨的手机号，所以要删除用户可能输入的0前缀
		req.CardNumber = ss_func.PreThirdpartyCardNumber(req.CardNumber)

		//第三方渠道卡号长度参数校验
		cardLength := len(req.CardNumber)
		if cardLength > constants.ThirdPartyCardNumberLengthMax || cardLength < constants.ThirdPartyCardNumberLengthMin {
			ss_log.Error("第三方渠道银行卡长度必需在%v到%v之间, CardNumber长度[%v]", constants.ThirdPartyCardNumberLengthMin, constants.ThirdPartyCardNumberLengthMax, cardLength)
			reply.ResultCode = ss_err.ERR_CardNumberLength_FAILD
			return nil
		}

		//  判断卡是否存在
		if cnt := dao.CardDaoInst.QueryThirdPartyCardCnt(req.CardNumber, req.AccountType, req.ChannelNo, channelData.ChannelType, req.BalanceType); cnt != "0" {
			ss_log.Error("err=[第三方渠道银行卡号已存在,卡号为--->%s]", req.CardNumber)
			reply.ResultCode = ss_err.ERR_CARD_IS_EXIST
			return nil
		}
	default:
		ss_log.Error("查询出的渠道类型出错 ChannelNo[%v]", channelData.ChannelType)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	errCode := dao.CardDaoInst.InsertCard(req.AccountNo, req.ChannelNo, req.CardAccName, req.CardNumber, req.BalanceType, req.IsDefault, req.CollectStatus, req.AuditStatus, req.AccountType, channelData.ChannelType)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("插入银行卡出错，req[%+v],err=[%v]", req, errCode)
		reply.ResultCode = ss_err.ERR_SYS_DB_ADD
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 插入商家新卡信息
func (*CustHandler) InsertOrUpdateBusinessCard(ctx context.Context, req *go_micro_srv_cust.InsertOrUpdateBusinessCardRequest, reply *go_micro_srv_cust.InsertOrUpdateBusinessCardReply) error {
	//参数校验
	switch req.AccountType { //只有个人商家和企业商家调用该接口
	case constants.AccountType_PersonalBusiness:
		//现在是再app端登录的个人商家添加卡时需要卡的名称和账号的实名认证名称一致
		accountAuthName, err := dao.AccDaoInstance.GetAuthNameFromUid(req.AccountNo)
		if err != nil {
			ss_log.Error("账号的实名认证名称查询出错,err[%v],accountNo[%v]", err, req.AccountNo)
			reply.ResultCode = ss_err.ERR_ACCOUNT_REAL_NAME_NOT_SAME
			return nil
		}

		if accountAuthName != req.Name {
			ss_log.Error("账号的实名认证名称[%v]和要添加的持卡人姓名[%v]不一致,不允许添加。。", accountAuthName, accountAuthName)
			reply.ResultCode = ss_err.ERR_ACCOUNT_REAL_NAME_NOT_SAME
			return nil
		}
	case constants.AccountType_EnterpriseBusiness:
	default:
		ss_log.Error("AccountType类型错误:[%v]", req.AccountType)
		reply.ResultCode = ss_err.ERR_SYS_NO_API_AUTH
		return nil
	}

	authInfo, err := dao.AccDaoInstance.GetAuthInfoByAccountNo(req.AccountNo, req.AccountType)
	if err != nil && err != sql.ErrNoRows {
		ss_log.Error("查询商家实名认证状态失败，accountNo=%v, err=%v", req.AccountNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	if authInfo == nil {
		ss_log.Error("查询商家未实名认证，accountNo=%v, err=%v", req.AccountNo, err)
		reply.ResultCode = ss_err.ERR_AddCard_No_BusinessRealName_FAILD
		return nil
	}
	if authInfo.AuthStatus != constants.AuthMaterialStatus_Passed {
		ss_log.Error("账号实名认证状态未通过，无法添加银行卡，accountNo=%v, AuthStatus=%v", req.AccountNo, authInfo.AuthStatus)
		reply.ResultCode = ss_err.ERR_HaveNotPass_BusinessRealName_FAILD
		return nil
	}

	//根据渠道id查询出渠道信息
	channelData, err := dao.ChannelDaoInst.GetBusinessChannelDetail(req.ChannelId)
	if err != nil {
		ss_log.Error("查询出错err[%v],渠道id[%v]", err, req.ChannelId)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	switch channelData.ChannelType {
	case constants.CHANNELTYPE_ORDINARY:
		//普通渠道卡号长度参数校验
		cardLength := len(req.CardNumber)
		if cardLength > constants.CardNumberLengthMax || cardLength < constants.CardNumberLengthMin {
			ss_log.Error("普通渠道银行卡长度必需在%v到%v之间, CardNumber长度[%v]", constants.CardNumberLengthMin, constants.CardNumberLengthMax, cardLength)
			reply.ResultCode = ss_err.ERR_CardNumberLength_FAILD
			return nil
		}

		//  判断卡是否存在
		if cnt := dao.CardBusinessDaoInst.QueryCardCnt(req.CardNumber, req.AccountType); cnt != "0" {
			ss_log.Error("err=[普通渠道银行卡号已存在,卡号为--->%s]", req.CardNumber)
			reply.ResultCode = ss_err.ERR_CARD_IS_EXIST
			return nil
		}
	case constants.CHANNELTYPE_THIRDPARTY:
		//第三方的渠道卡号都是柬埔寨的手机号，所以要删除用户可能输入的0前缀
		req.CardNumber = ss_func.PreThirdpartyCardNumber(req.CardNumber)

		//第三方渠道卡号长度参数校验
		cardLength := len(req.CardNumber)
		if cardLength > constants.ThirdPartyCardNumberLengthMax || cardLength < constants.ThirdPartyCardNumberLengthMin {
			ss_log.Error("第三方渠道银行卡长度必需在%v到%v之间, CardNumber长度[%v]", constants.ThirdPartyCardNumberLengthMin, constants.ThirdPartyCardNumberLengthMax, cardLength)
			reply.ResultCode = ss_err.ERR_CardNumberLength_FAILD
			return nil
		}

		//  判断卡是否存在
		if cnt := dao.CardBusinessDaoInst.QueryThirdPartyCardCnt(req.CardNumber, req.AccountType, channelData.ChannelNo, channelData.ChannelType, channelData.CurrencyType); cnt != "0" {
			ss_log.Error("err=[第三方渠道银行卡号已存在,卡号为--->%s]", req.CardNumber)
			reply.ResultCode = ss_err.ERR_CARD_IS_EXIST
			return nil
		}
	default:
		ss_log.Error("查询出的渠道类型出错 ChannelNo[%v]", channelData.ChannelType)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//添加卡
	if _, errCode := dao.CardBusinessDaoInst.InsertCard(dao.CardBusinessDao{
		AccountNo:     req.AccountNo,
		ChannelNo:     channelData.ChannelNo,
		Name:          req.Name,
		CardNum:       req.CardNumber,
		BalanceType:   channelData.CurrencyType,
		IsDefault:     req.IsDefault,
		CollectStatus: "1",
		AuditStatus:   constants.AuditOrderStatus_Passed,
		AccountType:   req.AccountType,
		ChannelType:   channelData.ChannelType,
	}); errCode != nil {
		ss_log.Error("err=[%v]", errCode.Error())
		reply.ResultCode = ss_err.ERR_SYS_DB_ADD
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (c *CustHandler) GetPosChannelList(ctx context.Context, req *go_micro_srv_cust.GetPosChannelListRequest, reply *go_micro_srv_cust.GetPosChannelListReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "cs.is_delete", Val: "0", EqType: "="},
		{Key: "ch.use_status", Val: "1", EqType: "="},                 //仓库渠道使用状态为使用
		{Key: "cs.use_status", Val: "1", EqType: "="},                 //服务商渠道使用状态为使用
		{Key: "cs.currency_type", Val: req.CurrencyType, EqType: "="}, //币种
	})

	sqlCnt := "select count(1) " +
		" from channel_servicer cs " +
		" left join channel ch on ch.channel_no = cs.channel_no  " + whereModel.WhereStr
	var total sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by cs.is_recom desc, cs.create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT  cs.channel_no, cs.is_recom, cs.currency_type, ch.logo_img_no, ch.logo_img_no_grey, ch.color_begin, ch.color_end, ch.channel_name, cs.id " +
		" FROM channel_servicer cs " +
		" left join channel ch on ch.channel_no = cs.channel_no  " + where2
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	var datas []*go_micro_srv_cust.ChannelDataSimple
	for rows.Next() {
		data := go_micro_srv_cust.ChannelDataSimple{}
		var logoImgNo, logoImgNoGrey, colorBegin, colorEnd sql.NullString
		if err = rows.Scan(
			&data.ChannelNo,
			&data.IsRecom,

			&data.CurrencyType,
			&logoImgNo,
			&logoImgNoGrey,
			&colorBegin,
			&colorEnd,
			&data.ChannelName,
			&data.Id,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		data.ColorBegin = colorBegin.String
		data.ColorEnd = colorEnd.String

		if logoImgNo.String != "" {
			data.LogoImgNo = logoImgNo.String
			//查询处图片的url，使前端可显示出来
			reqImg := &go_micro_srv_cust.UnAuthDownloadImageRequest{
				ImageId: logoImgNo.String,
			}
			replyImg := &go_micro_srv_cust.UnAuthDownloadImageReply{}
			c.UnAuthDownloadImage(ctx, reqImg, replyImg)
			if replyImg.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取图片url失败")
			} else {
				data.LogoImgUrl = replyImg.ImageUrl
			}
		}
		if logoImgNoGrey.String != "" {
			data.LogoImgNoGrey = logoImgNoGrey.String
			//查询处图片的url，使前端可显示出来
			reqImg := &go_micro_srv_cust.UnAuthDownloadImageRequest{
				ImageId: logoImgNoGrey.String,
			}
			replyImg := &go_micro_srv_cust.UnAuthDownloadImageReply{}
			c.UnAuthDownloadImage(ctx, reqImg, replyImg)
			if replyImg.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取图片url失败")
			} else {
				data.LogoImgUrlGrey = replyImg.ImageUrl
			}
		}
		datas = append(datas, &data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

func (c *CustHandler) GetUseChannelList(ctx context.Context, req *go_micro_srv_cust.GetUseChannelListRequest, reply *go_micro_srv_cust.GetUseChannelListReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ch.is_delete", Val: "0", EqType: "="},
		{Key: "ch.use_status", Val: "1", EqType: "="}, //仓库渠道使用状态为使用
		//{Key: "ch.currency_type", Val: req.CurrencyType, EqType: "="}, //币种
	})

	sqlCnt := "select count(1) from channel ch " + whereModel.WhereStr
	var total sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...); err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if total.String == "" || total.String == "0" {
		reply.ResultCode = ss_err.ERR_SUCCESS
		reply.Total = strext.ToInt32(total.String)
		return nil
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by ch.create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
	sqlStr := "SELECT  ch.channel_no, ch.logo_img_no, ch.logo_img_no_grey, ch.color_begin, ch.color_end," +
		" ch.channel_name, ch.channel_type " +
		" FROM channel ch " + whereModel.WhereStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	var datas []*go_micro_srv_cust.ChannelDataSimple
	for rows.Next() {
		data := go_micro_srv_cust.ChannelDataSimple{}
		var logoImgNo, logoImgNoGrey, colorBegin, colorEnd, channelType sql.NullString
		if err = rows.Scan(
			&data.ChannelNo,
			&logoImgNo,
			&logoImgNoGrey,
			&colorBegin,
			&colorEnd,
			&data.ChannelName,
			&channelType,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}

		data.ColorBegin = colorBegin.String
		data.ColorEnd = colorEnd.String
		data.ChannelType = channelType.String

		if logoImgNo.String != "" {
			data.LogoImgNo = logoImgNo.String
			//查询处图片的url，使前端可显示出来
			reqImg := &go_micro_srv_cust.UnAuthDownloadImageRequest{
				ImageId: logoImgNo.String,
			}
			replyImg := &go_micro_srv_cust.UnAuthDownloadImageReply{}
			c.UnAuthDownloadImage(ctx, reqImg, replyImg)
			if replyImg.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取图片url失败")
			} else {
				data.LogoImgUrl = replyImg.ImageUrl
			}
		}
		if logoImgNoGrey.String != "" {
			data.LogoImgNoGrey = logoImgNoGrey.String
			//查询处图片的url，使前端可显示出来
			reqImg := &go_micro_srv_cust.UnAuthDownloadImageRequest{
				ImageId: logoImgNoGrey.String,
			}
			replyImg := &go_micro_srv_cust.UnAuthDownloadImageReply{}
			c.UnAuthDownloadImage(ctx, reqImg, replyImg)
			if replyImg.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取图片url失败")
			} else {
				data.LogoImgUrlGrey = replyImg.ImageUrl
			}
		}

		data.Temp = "0"
		datas = append(datas, &data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}
