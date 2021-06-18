package handler

import (
	"context"

	"a.a/cu/ss_big"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	authProto "a.a/mp-server/common/proto/auth"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/cust-srv/dao"
	"a.a/mp-server/cust-srv/i"
)

/**
企业商家向总部充值
*/
func (c *CustHandler) AddBusinessToHead(ctx context.Context, req *custProto.AddBusinessToHeadRequest, reply *custProto.AddBusinessToHeadReply) error {

	// 判断转账金额
	if strext.ToFloat64(req.Amount) <= 0 || req.Amount == "" {
		ss_log.Error("err=[商家充值接口, 充值金额为0或者为空,传入的金额为----->%s]", req.Amount)
		reply.ResultCode = ss_err.ERR_WALLET_AMOUNT_NULL
		return nil
	}

	reply1, err := i.AuthHandlerInst.Client.CheckPayPWD(ctx, &authProto.CheckPayPWDRequest{
		AccountUid:  req.AccountUid,
		AccountType: req.AccountType,
		Password:    req.PayPwd,
		IdenNo:      req.IdenNo,
		NonStr:      req.NonStr,
	})
	if err != nil || reply1.ResultCode != ss_err.ERR_SUCCESS {
		ss_log.Error("确认支付密码出错 err=[%v],ResultCode=[%v] ", err, reply1.ResultCode)
		reply.ResultCode = ss_err.ERR_BusinessPayPwd_FAILD
		return nil
	}

	// 校验收款人账户
	whereList := []*model.WhereSqlCond{
		{Key: "ca.is_delete", Val: "0", EqType: "="},
		{Key: "ca.card_no", Val: req.CardNo, EqType: "="},
	}
	cardData, errGet := dao.CardHeadDaoInst.GetHeadCardBusinessDetail(whereList)
	if errGet != nil || cardData.CardNo == "" || cardData.MoneyType != req.MoneyType {
		ss_log.Error("卡id[%v],查询出的平台收款卡信息为[%+v]", req.CardNo, cardData)
		reply.ResultCode = ss_err.ERR_CARD_NOT_EXIST
		return nil
	}

	// 判断限额
	if cardData.SaveMaxAmount != "" {
		if strext.ToFloat64(req.Amount) > strext.ToFloat64(cardData.SaveMaxAmount) {
			ss_log.Error("商家充值,交易金额超出单笔最大金额,交易金额为: %s,单笔提现最大金额为: %s", req.Amount, cardData.SaveMaxAmount)
			reply.ResultCode = ss_err.ERR_LOCAL_RULE_EXCEED_AMOUNT
			return nil
		}
	}

	//获取手续费
	fee := "0"
	switch cardData.SaveChargeType {
	case constants.Charge_Type_Rate: // 手续费比例收取
		feesDeci := ss_count.CountFees(req.Amount, cardData.SaveRate, "0")
		// 取整
		fee = ss_big.SsBigInst.ToRound(feesDeci, 0, ss_big.RoundingMode_HALF_EVEN).String()
	case constants.Charge_Type_Count: // 按单笔手续费收取
		fee = cardData.SaveSingleMinFee
	}

	//上传凭证
	upReg := &custProto.UploadImageRequest{
		ImageStr:     req.ImageBase64,
		AccountUid:   req.AccountUid,
		Type:         constants.UploadImage_Auth,
		AddWatermark: constants.AddWatermark_True,
	}

	upReply := &custProto.UploadImageReply{}
	errU := c.UploadImage(ctx, upReg, upReply)
	if errU != nil || upReply.ResultCode != ss_err.ERR_SUCCESS {
		ss_log.Error("addImg1Err=[%v]", upReply.ResultCode)
		reply.ResultCode = ss_err.ERR_SAVE_IMAGE_FAILD
		return nil
	}

	arriveAmount := ss_count.Sub(req.Amount, fee).String()

	// 存记录
	logNo := dao.LogBusinessToHeadquartersDaoInst.Insert(dao.LogBusinessToHeadquartersDao{
		BusinessNo:     req.IdenNo,
		CurrencyType:   req.MoneyType,                  //币种
		Amount:         req.Amount,                     //金额
		CollectionType: constants.BANK_COLLECTION_TYPE, //收款方式,1-支票;2-现金;3-银行转账;4-其他
		CardNo:         req.CardNo,                     //总部收款卡uid
		ImageId:        upReply.ImageId,                //凭证 图片id
		ArriveAmount:   arriveAmount,                   //实际到账金额
		Fee:            fee,                            //手续费
	})
	if logNo == "" {
		reply.ResultCode = ss_err.ERR_TRANSFER_TO_HEAD_QUARTERS_FAILD
		return nil
	}

	reply.LogNo = logNo
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
app内的个人商家向总部充值申请
*/
func (c *CustHandler) AddIndividualBusinessToHead(ctx context.Context, req *custProto.AddIndividualBusinessToHeadRequest, reply *custProto.AddIndividualBusinessToHeadReply) error {
	// 判断转账金额
	if strext.ToFloat64(req.Amount) <= 0 || req.Amount == "" {
		ss_log.Error("err=[商家充值接口, 充值金额为0或者为空,传入的金额为----->%s]", req.Amount)
		reply.ResultCode = ss_err.ERR_WALLET_AMOUNT_NULL
		return nil
	}

	//验证上传的凭证是在数据库内的
	imgData, imgErr := dao.ImageDaoInstance.GetImageUrlById(req.ImageId)
	if imgErr != nil {
		ss_log.Error("查询上传的凭证图片出错,ImageId=[%v],err=[%v]", req.ImageId, imgErr)
		reply.ResultCode = ss_err.ERR_IMAGE_OP_FAILD
		return nil
	}

	if imgData.ImageId == "" {
		ss_log.Error("查询上传的凭证图片不存在,ImageId=[%v]", req.ImageId)
		reply.ResultCode = ss_err.ERR_IMAGE_OP_FAILD
		return nil
	}

	//验证的是用户支付密码
	replyCheckPayPwd, errCheckPayPwd := i.AuthHandlerInst.Client.CheckPayPWD(context.TODO(), &authProto.CheckPayPWDRequest{
		AccountUid:  req.AccountUid,
		AccountType: req.AccountType,
		Password:    req.PayPwd,
		NonStr:      req.NonStr,
		IdenNo:      req.IdenNo,
	})
	if errCheckPayPwd != nil {
		ss_log.Error("paymentPwdErrLimit 调用验证支付密码接口出错,err[%v]", errCheckPayPwd)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}

	if replyCheckPayPwd.ResultCode != ss_err.ERR_SUCCESS {
		reply.ResultCode = replyCheckPayPwd.ResultCode
		reply.PayPasswordErrTips = replyCheckPayPwd.ErrTips //提示还可以输入几次错误支付密码
		return nil
	}

	// 校验收款人账户
	whereList := []*model.WhereSqlCond{
		{Key: "ca.is_delete", Val: "0", EqType: "="},
		{Key: "ca.card_no", Val: req.CardNo, EqType: "="},
	}
	cardData, errGet := dao.CardHeadDaoInst.GetHeadCardBusinessDetail(whereList)
	if errGet != nil || cardData.CardNo == "" || cardData.MoneyType != req.MoneyType {
		ss_log.Error("卡id[%v],查询出的平台收款卡信息为[%+v]", req.CardNo, cardData)
		reply.ResultCode = ss_err.ERR_CARD_NOT_EXIST
		return nil
	}

	// 判断限额
	if cardData.SaveMaxAmount != "" {
		if strext.ToFloat64(req.Amount) > strext.ToFloat64(cardData.SaveMaxAmount) {
			ss_log.Error("商家充值,交易金额超出单笔最大金额,交易金额为: %s,单笔提现最大金额为: %s", req.Amount, cardData.SaveMaxAmount)
			reply.ResultCode = ss_err.ERR_LOCAL_RULE_EXCEED_AMOUNT
			return nil
		}
	}

	//获取手续费
	fee := "0"
	switch cardData.SaveChargeType {
	case constants.Charge_Type_Rate: // 手续费比例收取
		feesDeci := ss_count.CountFees(req.Amount, cardData.SaveRate, "0")
		// 取整
		fee = ss_big.SsBigInst.ToRound(feesDeci, 0, ss_big.RoundingMode_HALF_EVEN).String()
	case constants.Charge_Type_Count: // 按单笔手续费收取
		fee = cardData.SaveSingleMinFee
	}

	arriveAmount := ss_count.Sub(req.Amount, fee).String() //实际到账金额

	businessNo := dao.RelaAccIdenDaoInst.GetIdenFromAcc(req.AccountUid, constants.AccountType_PersonalBusiness)
	// 存记录
	logNo := dao.LogBusinessToHeadquartersDaoInst.Insert(dao.LogBusinessToHeadquartersDao{
		BusinessNo:     businessNo,
		CurrencyType:   req.MoneyType,                  //币种
		Amount:         req.Amount,                     //金额
		CollectionType: constants.BANK_COLLECTION_TYPE, //收款方式,1-支票;2-现金;3-银行转账;4-其他
		CardNo:         req.CardNo,                     //总部收款卡uid
		ImageId:        req.ImageId,                    //凭证 图片id
		ArriveAmount:   arriveAmount,                   //实际到账金额
		Fee:            fee,                            //手续费
	})
	if logNo == "" {
		reply.ResultCode = ss_err.ERR_TRANSFER_TO_HEAD_QUARTERS_FAILD
		return nil
	}

	reply.LogNo = logNo
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
查询商家充值订单列表
*/
func (c *CustHandler) GetBusinessToHeadList(ctx context.Context, req *custProto.GetBusinessToHeadListRequest, reply *custProto.GetBusinessToHeadListReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "bth.business_no", Val: req.IdenNo, EqType: "="},
		{Key: "bu.account_no", Val: req.AccountNo, EqType: "="},
		{Key: "bth.currency_type", Val: req.MoneyType, EqType: "="},
		{Key: "bth.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "bth.log_no", Val: req.LogNo, EqType: "like"},
		{Key: "bth.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "bth.create_time", Val: req.EndTime, EqType: "<="},
	}

	//获取总数
	total := dao.LogBusinessToHeadquartersDaoInst.GetCnt(whereList)

	if total == "" || total == "0" {
		reply.Total = strext.ToInt32(total)
		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	}

	//获取列表信息
	datas, err := dao.LogBusinessToHeadquartersDaoInst.GetBusinessToHeadList(whereList, strext.ToInt(req.Page), strext.ToInt(req.PageSize))
	if err != nil {
		ss_log.Error("查询商家充值订单失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
查询商家充值订单详情
*/
func (c *CustHandler) GetBusinessToHeadDetail(ctx context.Context, req *custProto.GetBusinessToHeadDetailRequest, reply *custProto.GetBusinessToHeadDetailReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "bth.log_no", Val: req.LogNo, EqType: "="},
	}

	//获取详情
	data, err := dao.LogBusinessToHeadquartersDaoInst.GetBusinessToHeadDetail(whereList)
	if err != nil {
		ss_log.Error("查询商家充值订单失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
