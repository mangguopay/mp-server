package handler

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"a.a/cu/db"
	"a.a/cu/ss_time"
	"a.a/mp-server/common/global"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/business-bill-srv/dao"
	"a.a/mp-server/common/constants"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type BusinessBillHandler struct{}

var BusinessBillHandlerInst BusinessBillHandler

/*
下单接口
1.用户扫码商家带金额的二维码
2.商家扫码用户的付款二维码
3.APP支付下单
*/
func (b *BusinessBillHandler) Prepay(ctx context.Context, req *businessBillProto.PrepayRequest, reply *businessBillProto.PrepayReply) error {
	req.CurrencyType = strings.ToUpper(req.CurrencyType)
	resultCode, err := CheckPrepayReq(req)
	if err != nil {
		ss_log.Error("Prepay() 参数错误，err=%v", err)
		reply.ResultCode = resultCode
		reply.Msg = ss_err.GetMsg(resultCode, req.Lang)
		return nil
	}

	//查询app
	app, err := dao.BusinessAppDaoInst.GetAppInfoByAppId(req.AppId)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("查询应用信息失败,应用不存在,AppId:%v", req.AppId)
			reply.ResultCode = ss_err.AppNotExist
			reply.Msg = ss_err.GetMsg(ss_err.AppNotExist, req.Lang)
			return nil
		}
		ss_log.Error("查询应用信息失败,err:%v,AppId:%v", err, req.AppId)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	if app.Status != constants.BusinessAppStatus_Up {
		ss_log.Error("APP(%v)未上架", app.AppId)
		reply.ResultCode = ss_err.AppNotPutOn
		reply.Msg = ss_err.GetMsg(ss_err.AppNotPutOn, req.Lang)
		return nil
	}

	//查询app签约 todo 新版产品签约
	signed, err := dao.BusinessSceneSignedDaoInst.GetSignedByTradeType(app.BusinessNo, req.TradeType)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("商家(%v)未签约,TradeType=%v", app.BusinessNo, req.TradeType)
			reply.ResultCode = ss_err.ProductUnsigned
			reply.Msg = ss_err.GetMsg(ss_err.ProductUnsigned, req.Lang)
			return nil
		}
		ss_log.Error("查询商家(%v)签约失败，TradeType=%v, err=%v", app.BusinessNo, req.TradeType, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	ss_log.Info("商户应用(%v)签约信息:%+v", app.AppId, strext.ToJson(signed))

	//校验APP的签约
	retCode, err := CheckBusinessSigned(signed)
	if err != nil {
		ss_log.Error("app校验结果：%v", err)
		reply.ResultCode = retCode
		reply.Msg = ss_err.GetMsg(retCode, req.Lang)
		return nil
	}

	//检查产品是否可用
	isEnabled, err := dao.BusinessSceneDao.GetSceneIsEnabled(signed.SceneNo)
	if err != nil {
		ss_log.Error("查询产品是否可用失败，SceneNo=%v, err=%v", signed.SceneNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	if !isEnabled {
		ss_log.Error("产品[%v]已被禁用，暂时不可交易", signed.SceneNo)
		reply.ResultCode = ss_err.SceneDisabled
		reply.Msg = ss_err.GetMsg(ss_err.SceneDisabled, req.Lang)
		return nil
	}

	// 查询商户交易配置（商户账号，商户是否启用，收款权限）
	businessConf, err := dao.BusinessDaoInst.GetTransConfig(app.BusinessNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("商户不存在,BusinessNo:%v", app.BusinessNo)
			reply.ResultCode = ss_err.BusinessNotExist
			reply.Msg = ss_err.GetMsg(ss_err.BusinessNotExist, req.Lang)
			return nil
		}
		ss_log.Error("查询商户交易配置失败,err:%v,BusinessNo:%v", err, app.BusinessNo)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	//校验商家是否能正常收款
	resultCode, err = CheckBusinessIncomeAuth(businessConf)
	if err != nil {
		ss_log.Error("商户交易配置校验结果: %v", err)
		reply.ResultCode = resultCode
		reply.Msg = ss_err.GetMsg(resultCode, req.Lang)
		return nil
	}

	// 查询商户的订单号是否已经存在
	exit, err := dao.BusinessBillDaoInst.OutOrderNoExist(app.BusinessNo, req.OutOrderNo)
	if err != nil {
		ss_log.Error("查询商户订单是否存在失败,err:%v,BusinessNo:%v,OutOrderNo:%v", err, app.BusinessNo, req.OutOrderNo)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	if exit {
		ss_log.Error("商户订单已经存在,BusinessNo:%v,OutOrderNo:%v", app.BusinessNo, req.OutOrderNo)
		reply.ResultCode = ss_err.OrderAlreadyExist
		reply.Msg = ss_err.GetMsg(ss_err.OrderAlreadyExist, req.Lang)
		return nil
	}

	// 应用没有设置费率，需要设置为字符串0，方便后续统一处理
	if signed.Rate == "" {
		ss_log.Info("应用没有设置费率,AppId:%v", req.AppId)
		signed.Rate = "0"
	}

	// todo 计算手续费, 四舍五入后结果
	fee := ss_count.CountFees(req.Amount, signed.Rate, "0").String()
	realAmount := ss_count.Sub(req.Amount, fee).String()

	ss_log.Info("计算手续费结果,OutOrderNo:%v,金额为:%s,费率为:%s,手续费为:%s,实际金额为:%s", req.OutOrderNo, req.Amount, signed.Rate, fee, realAmount)

	//根据付款码查询付款人
	var payerAccNo string
	var orderNo = strext.GetDailyId()
	if req.PaymentCode != "" {
		custAccNo, err := GetAccountNoByCode(req.PaymentCode)
		if err != nil {
			ss_log.Error("根据付款码查询付款人失败，PaymentCode=%v, err=%v", req.PaymentCode, err)
			reply.ResultCode = ss_err.PaymentCodeExpire
			reply.Msg = ss_err.GetMsg(ss_err.PaymentCodeExpire, req.Lang)
			return nil
		}
		payerAccNo = custAccNo
		orderNo = req.PaymentCode
	}

	// 获取平台订单号
	order := dao.BusinessBillDao{
		OrderNo:           orderNo,
		Fee:               fee,
		Amount:            req.Amount,
		RealAmount:        realAmount,
		OrderStatus:       constants.BusinessOrderStatusPending,
		Rate:              signed.Rate,
		AccountNo:         payerAccNo,
		BusinessNo:        app.BusinessNo,
		BusinessAccountNo: businessConf.BusinessAccNo,
		AppId:             req.AppId,
		AppName:           app.AppName,
		SimplifyName:      app.SimplifyName, //商家简称
		CurrencyType:      req.CurrencyType,
		Remark:            req.Remark,
		NotifyUrl:         req.NotifyUrl,
		RreturnUrl:        req.ReturnUrl,
		OutOrderNo:        req.OutOrderNo,
		Subject:           req.Subject,
		SceneNo:           signed.SceneNo,
		ExpireTime:        strext.ToInt64(req.TimeExpire),
		TradeType:         req.TradeType,
		BusinessChannelNo: signed.BusinessChannelNo,
	}

	// 获取数据库连接
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		ss_log.Error("获取数据连接失败")
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// 开启事物
	tx, err := dbHandler.BeginTx(ctx, nil)
	if err != nil {
		ss_log.Error("开启事物失败,err:%v", err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	// 插入bill表记录
	if err := dao.BusinessBillDaoInst.InsertOrderTx(tx, order); err != nil {
		if strings.Contains(err.Error(), ss_sql.DbDuplicateKey) {
			ss_log.Error("订单号重复,err:%v,orderInfo:%v", err, strext.ToJson(order))
		} else {
			ss_log.Error("插入订单失败,err:%v,orderInfo:%v", err, strext.ToJson(order))
		}
		ss_sql.Rollback(tx)
		reply.ResultCode = ss_err.PlaceAnOrderFail
		reply.Msg = ss_err.GetMsg(ss_err.PlaceAnOrderFail, req.Lang)
		return nil
	}

	// 面对面支付， 用户扫码商家
	if req.TradeType == constants.TradeTypeModernpayFaceToFace && req.PaymentCode == "" {
		//插入订单支付二维码ID
		qrCodeId := constants.GetQrCodeId(order.OrderNo)
		qrCodErr := dao.BusinessBillQrCodeInst.InsertOrderQrCode(tx, order.OrderNo, qrCodeId)
		if qrCodErr != nil {
			if strings.Contains(qrCodErr.Error(), ss_sql.DbDuplicateKey) {
				ss_log.Error("订单号重复,err:%v,orderInfo:%v", qrCodErr, strext.ToJson(order))
			} else {
				ss_log.Error("插入订单二维码失败,err:%v,orderInfo:%v", qrCodErr, strext.ToJson(order))
			}
			ss_sql.Rollback(tx)
			reply.ResultCode = ss_err.PlaceAnOrderFail
			reply.Msg = ss_err.GetMsg(ss_err.PlaceAnOrderFail, req.Lang)
			return nil
		}
		reply.QrCodeId = constants.GetQrCodeUrl(qrCodeId)
	}
	ss_sql.Commit(tx)

	reply.ResultCode = ss_err.Success
	reply.Msg = ss_err.GetMsg(ss_err.Success, req.Lang)
	reply.OrderNo = order.OrderNo
	reply.OutOrderNo = order.OutOrderNo

	// APP支付独有返回参数
	if req.TradeType == constants.TradeTypeModernpayAPP {
		reply.AppPayContent = GetAppPayContent(order)
	}

	return nil
}

/*
用户扫码下单(用户扫商家的固定二维码下单)
*/
func (b *BusinessBillHandler) QrCodeFixedPrePay(ctx context.Context, req *businessBillProto.QrCodeFixedPrePayRequest, reply *businessBillProto.QrCodeFixedPrePayReply) error {
	req.CurrencyType = strings.ToUpper(req.CurrencyType)
	resultCode, err := CheckQrCodeFixedPrePay(req)
	if err != nil {
		ss_log.Error("Prepay() 参数错误，err=%v", err)
		reply.ResultCode = resultCode
		reply.Msg = ss_err.GetMsg(resultCode, req.Lang)
		return nil
	}

	//用户基本信息
	custInfo, err := dao.CustDaoInst.QueryCustInfo(req.AccountNo, "")
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("查询用户账号信息失败,用户不存在,AccountNo:%v", req.AccountNo)
			reply.ResultCode = ss_err.AccountNoNotExist
			return nil
		}
		ss_log.Error("查询用户账号信息失败,err:%v,AccountNo:%v", err, req.AccountNo)
		reply.ResultCode = ss_err.SystemErr
		return nil
	}

	//检查用户是否已实名
	if custInfo.IndividualAuthStatus != constants.AuthMaterialStatus_Passed {
		ss_log.Error("用户账号未实名认证, AccountNo:%v, IndividualAuthStatus=%v", req.AccountNo, custInfo.IndividualAuthStatus)
		reply.ResultCode = ss_err.UserHasNoRealName
		return nil
	}

	//判断用户是否有交易权限
	if custInfo.TradingAuthority == constants.TradingAuthorityForbid {
		ss_log.Error("用户[%v]被禁止交易", req.AccountNo)
		reply.ResultCode = ss_err.AccountNoNotTradeForbid
		return nil
	}

	//判断用户是否有出款权限
	if custInfo.OutgoAuthorization == constants.CustOutTransferAuthorizationDisabled {
		ss_log.Error("用户[%s]没有出款权限", req.AccountNo)
		reply.ResultCode = ss_err.AccountNoNotTradeForbid
		return nil
	}

	//校验APP是否能交易
	app, err := dao.BusinessAppDaoInst.GetAppInfoByFixedQrCode(req.QrCodeId)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("应用固码不存在，QrCodeId=%v", req.QrCodeId)
			reply.ResultCode = ss_err.QrCodeNotInvalid
			reply.Msg = ss_err.GetMsg(ss_err.QrCodeNotInvalid, req.Lang)
			return nil
		}
		ss_log.Error("查询商家APP失败，QrCodeId=%v, err=%v", req.QrCodeId, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	if app.Status != constants.BusinessAppStatus_Up {
		ss_log.Error("APP(%v)未上架", app.AppId)
		reply.ResultCode = ss_err.AppNotPutOn
		reply.Msg = ss_err.GetMsg(ss_err.AppNotPutOn, req.Lang)
		return nil
	}

	//查询app是否签约了"当面付"
	signed, err := dao.BusinessSceneSignedDaoInst.GetSignedByTradeType(app.BusinessNo, constants.TradeTypeModernpayFaceToFace)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("App(%v)未签约,TradeType=%v", app.AppId, constants.TradeTypeModernpayFaceToFace)
			reply.ResultCode = ss_err.QrCodeNotAvailable
			reply.Msg = ss_err.GetMsg(ss_err.QrCodeNotAvailable, req.Lang)
			return nil
		}
		ss_log.Error("查询app(%v)签约失败，TradeType=%v, err=%v", app.AppId, constants.TradeTypeModernpayFaceToFace, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	ss_log.Info("商户应用(%v)签约信息:%+v", app.AppId, strext.ToJson(signed))

	//校验APP的签约
	retCode, err := CheckBusinessSigned(signed)
	if err != nil {
		ss_log.Error("app校验结果：%v", err)
		reply.ResultCode = retCode
		reply.Msg = ss_err.GetMsg(retCode, req.Lang)
		return nil
	}

	//检查产品是否可用
	isEnabled, err := dao.BusinessSceneDao.GetSceneIsEnabled(signed.SceneNo)
	if err != nil {
		ss_log.Error("查询产品是否可用失败，SceneNo=%v, err=%v", signed.SceneNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	if !isEnabled {
		ss_log.Error("产品[%v]已被系统禁用，暂时不可交易", signed.SceneNo)
		reply.ResultCode = ss_err.SceneDisabled
		reply.Msg = ss_err.GetMsg(ss_err.SceneDisabled, req.Lang)
		return nil
	}

	// 查询商户交易配置（商户账号，商户是否启用，收款权限）
	businessConf, err := dao.BusinessDaoInst.GetTransConfig(app.BusinessNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("商户不存在,BusinessNo:%v", app.BusinessNo)
			reply.ResultCode = ss_err.QrCodeNotInvalid
			reply.Msg = ss_err.GetMsg(ss_err.QrCodeNotInvalid, req.Lang)
			return nil
		}
		ss_log.Error("查询商户交易配置失败,err:%v,BusinessNo:%v", err, app.BusinessNo)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	resultCode, err = CheckBusinessIncomeAuth(businessConf)
	if err != nil {
		ss_log.Error("商户交易配置校验结果: %v", err)
		reply.ResultCode = resultCode
		reply.Msg = ss_err.GetMsg(resultCode, req.Lang)
		return nil
	}

	//用户余额是否足够支付
	custVaType := global.GetUserVAccType(req.CurrencyType, true)
	custVAccNo, err := dao.VaccountDaoInst.GetVaccountNo(req.AccountNo, custVaType)
	if err != nil {
		ss_log.Error("查询用户虚账失败，AccountNo=%v, err=%v", req.AccountNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	balance, err := dao.VaccountDaoInst.GetBalanceByVAccNo(custVAccNo)
	if err != nil {
		ss_log.Error("查询用户余额失败，AccountNo=%v, err=%v", req.AccountNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	if strext.ToInt64(req.Amount) > balance {
		ss_log.Error("用户[%v]余额不足，req.Amount=%v, balance=%v", req.AccountNo, req.Amount, balance)
		reply.ResultCode = ss_err.BalanceNotEnough
		reply.Msg = ss_err.GetMsg(ss_err.BalanceNotEnough, req.Lang)
		return nil
	}

	// 应用没有设置费率，需要设置为字符串0，方便后续统一处理
	if signed.Rate == "" {
		ss_log.Info("应用没有设置费率,AppId:%v", app.AppId)
		signed.Rate = "0"
	}

	// todo 计算手续费, 四舍五入后结果
	fee := ss_count.CountFees(req.Amount, signed.Rate, "0").String()
	realAmount := ss_count.Sub(req.Amount, fee).String()

	ss_log.Info("计算手续费结果,金额为:%s,费率为:%s,手续费为:%s,实际金额为:%s", req.Amount, signed.Rate, fee, realAmount)

	// 获取平台订单号
	order := dao.BusinessBillDao{
		OrderNo:           strext.GetDailyId(),
		Fee:               fee,
		Amount:            req.Amount,
		RealAmount:        realAmount,
		OrderStatus:       constants.BusinessOrderStatusPending,
		Rate:              signed.Rate,
		BusinessNo:        app.BusinessNo,
		BusinessAccountNo: businessConf.BusinessAccNo,
		AppId:             app.AppId,
		CurrencyType:      req.CurrencyType,
		Subject:           req.Subject,
		Remark:            req.Remark,
		SceneNo:           signed.SceneNo,
		ExpireTime:        ss_time.Now(global.Tz).Add(constants.BusinessOrderExpireTime * time.Minute).Unix(),
		TradeType:         constants.TradeTypeModernpayFaceToFace,
		BusinessChannelNo: signed.BusinessChannelNo,
	}

	// 插入bill表记录
	if err := dao.BusinessBillDaoInst.InsertOrder(order); err != nil {
		if strings.Contains(err.Error(), ss_sql.DbDuplicateKey) {
			ss_log.Error("订单号重复,err:%v,orderInfo:%v", err, strext.ToJson(order))
		} else {
			ss_log.Error("插入订单失败,err:%v,orderInfo:%v", err, strext.ToJson(order))
		}
		reply.ResultCode = ss_err.PlaceAnOrderFail
		reply.Msg = ss_err.GetMsg(ss_err.PlaceAnOrderFail, req.Lang)
		return nil
	}

	reply.ResultCode = ss_err.Success
	reply.Msg = ss_err.GetMsg(ss_err.Success, req.Lang)
	reply.OrderNo = order.OrderNo
	reply.Amount = order.Amount
	reply.Subject = order.Subject
	return nil
}

/*
个人商家设置金额下单
*/
func (b *BusinessBillHandler) PersonalBusinessPrepay(ctx context.Context, req *businessBillProto.PersonalBusinessPrepayRequest, reply *businessBillProto.PersonalBusinessPrepayReply) error {
	req.CurrencyType = strings.ToUpper(req.CurrencyType)
	resultCode, err := CheckPersonalBusinessPrepayReq(req)
	if err != nil {
		ss_log.Error("PersonalBusinessPrepay() 参数错误，err=%v", err)
		reply.ResultCode = resultCode
		reply.Msg = ss_err.GetMsg(resultCode, req.Lang)
		return nil
	}

	business, err := dao.BusinessDaoInst.GetTransConfigByAccNo(req.AccountNo, constants.AccountType_PersonalBusiness)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("用户[%v]没有个人商家身份", req.AccountNo)
			reply.ResultCode = ss_err.BusinessNotExist
			reply.Msg = ss_err.GetMsg(ss_err.BusinessNotExist, req.Lang)
			return nil
		}
		ss_log.Error("查询用户[%v]个人商家信息失败， err=%v", req.AccountNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	//校验商家是否能正常收款
	resultCode, err = CheckBusinessIncomeAuth(business)
	if err != nil {
		ss_log.Error("商户交易配置校验结果: %v", err)
		reply.ResultCode = resultCode
		reply.Msg = ss_err.GetMsg(resultCode, req.Lang)
		return nil
	}

	//商家产品签约
	signed, err := dao.BusinessSceneSignedDaoInst.GetSignedByTradeType(business.BusinessNo, constants.TradeTypeBusinessPay)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("个人商家(%v)未签约[%v]", business.BusinessNo, constants.TradeTypeBusinessPay)
			reply.ResultCode = ss_err.ProductUnsigned
			reply.Msg = ss_err.GetMsg(ss_err.ProductUnsigned, req.Lang)
			return nil
		}
		ss_log.Error("查询个人商家(%v)签约[%v]失败，err=%v", business.BusinessNo, constants.TradeTypeBusinessPay, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	//校验APP的签约
	retCode, err := CheckBusinessSigned(signed)
	if err != nil {
		ss_log.Error("app校验结果：%v", err)
		reply.ResultCode = retCode
		reply.Msg = ss_err.GetMsg(retCode, req.Lang)
		return nil
	}

	//检查产品是否可用
	isEnabled, err := dao.BusinessSceneDao.GetSceneIsEnabled(signed.SceneNo)
	if err != nil {
		ss_log.Error("查询产品是否可用失败，SceneNo=%v, err=%v", signed.SceneNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	if !isEnabled {
		ss_log.Error("产品[%v]已被禁用，暂时不可交易", signed.SceneNo)
		reply.ResultCode = ss_err.SceneDisabled
		reply.Msg = ss_err.GetMsg(ss_err.SceneDisabled, req.Lang)
		return nil
	}

	// 没有设置费率，需要设置为字符串0，方便后续统一处理
	if signed.Rate == "" {
		ss_log.Info("签约没有设置费率,signedNo:%v", signed.SceneNo)
		signed.Rate = "0"
	}

	// todo 计算手续费, 四舍五入后结果
	fee := ss_count.CountFees(req.Amount, signed.Rate, "0").String()
	realAmount := ss_count.Sub(req.Amount, fee).String()

	var orderNo = strext.GetDailyId()

	// 获取平台订单号
	order := dao.BusinessBillDao{
		OrderNo:           orderNo,
		Fee:               fee,
		Amount:            req.Amount,
		RealAmount:        realAmount,
		OrderStatus:       constants.BusinessOrderStatusPending,
		Rate:              signed.Rate,
		BusinessNo:        business.BusinessNo,
		BusinessAccountNo: req.AccountNo,
		SimplifyName:      business.SimplifyName, //商家简称
		CurrencyType:      req.CurrencyType,
		Remark:            req.Remark,
		Subject:           req.Subject,
		SceneNo:           signed.SceneNo,
		ExpireTime:        ss_time.Now(global.Tz).Add(constants.BusinessOrderExpireTime * time.Minute).Unix(),
		TradeType:         constants.TradeTypeBusinessPay,
		BusinessChannelNo: signed.BusinessChannelNo,
	}

	// 获取数据库连接
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		ss_log.Error("获取数据连接失败")
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// 开启事物
	tx, err := dbHandler.BeginTx(ctx, nil)
	if err != nil {
		ss_log.Error("开启事物失败,err:%v", err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	// 插入bill表记录
	if err := dao.BusinessBillDaoInst.InsertOrderTx(tx, order); err != nil {
		if strings.Contains(err.Error(), ss_sql.DbDuplicateKey) {
			ss_log.Error("订单号重复,err:%v,orderInfo:%v", err, strext.ToJson(order))
		} else {
			ss_log.Error("插入订单失败,err:%v,orderInfo:%v", err, strext.ToJson(order))
		}
		ss_sql.Rollback(tx)
		reply.ResultCode = ss_err.PlaceAnOrderFail
		reply.Msg = ss_err.GetMsg(ss_err.PlaceAnOrderFail, req.Lang)
		return nil
	}

	//插入订单支付二维码ID
	qrCodeId := constants.GetQrCodeId(order.OrderNo)
	qrCodErr := dao.BusinessBillQrCodeInst.InsertOrderQrCode(tx, order.OrderNo, qrCodeId)
	if qrCodErr != nil {
		if strings.Contains(qrCodErr.Error(), ss_sql.DbDuplicateKey) {
			ss_log.Error("订单号重复,err:%v,orderInfo:%v", qrCodErr, strext.ToJson(order))
		} else {
			ss_log.Error("插入订单二维码失败,err:%v,orderInfo:%v", qrCodErr, strext.ToJson(order))
		}
		ss_sql.Rollback(tx)
		reply.ResultCode = ss_err.PlaceAnOrderFail
		reply.Msg = ss_err.GetMsg(ss_err.PlaceAnOrderFail, req.Lang)
		return nil
	}

	ss_sql.Commit(tx)

	reply.ResultCode = ss_err.Success
	reply.Msg = ss_err.GetMsg(ss_err.Success, req.Lang)
	reply.OrderNo = order.OrderNo
	reply.QrCodeId = constants.GetQrCodeUrl(qrCodeId)
	return nil
}

/*
用户扫个人商家固定二维码下单
*/
func (b *BusinessBillHandler) PersonalBusinessCodeFixedPrePay(ctx context.Context, req *businessBillProto.PersonalBusinessCodeFixedPrePayRequest, reply *businessBillProto.PersonalBusinessCodeFixedPrePayReply) error {
	req.CurrencyType = strings.ToUpper(req.CurrencyType)
	resultCode, err := CheckPersonalBusinessCodeFixedPrePay(req)
	if err != nil {
		ss_log.Error("Prepay() 参数错误，err=%v", err)
		reply.ResultCode = resultCode
		reply.Msg = ss_err.GetMsg(resultCode, req.Lang)
		return nil
	}

	//用户基本信息
	custInfo, err := dao.CustDaoInst.QueryCustInfo(req.AccountNo, "")
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("查询用户账号信息失败,用户不存在,AccountNo:%v", req.AccountNo)
			reply.ResultCode = ss_err.AccountNoNotExist
			return nil
		}
		ss_log.Error("查询用户账号信息失败,err:%v,AccountNo:%v", err, req.AccountNo)
		reply.ResultCode = ss_err.SystemErr
		return nil
	}

	//检查用户是否已实名
	if custInfo.IndividualAuthStatus != constants.AuthMaterialStatus_Passed {
		ss_log.Error("用户账号未实名认证, AccountNo:%v, IndividualAuthStatus=%v", req.AccountNo, custInfo.IndividualAuthStatus)
		reply.ResultCode = ss_err.UserHasNoRealName
		return nil
	}

	//判断用户是否有交易权限
	if custInfo.TradingAuthority == constants.TradingAuthorityForbid {
		ss_log.Error("用户[%v]被禁止交易", req.AccountNo)
		reply.ResultCode = ss_err.AccountNoNotTradeForbid
		return nil
	}

	//判断用户是否有出款权限
	if custInfo.OutgoAuthorization == constants.CustOutTransferAuthorizationDisabled {
		ss_log.Error("用户[%s]没有出款权限", req.AccountNo)
		reply.ResultCode = ss_err.AccountNoNotTradeForbid
		return nil
	}

	businessNo, err := dao.BusinessFixedCodeDaoInst.GetBusinessByCode(req.QrCodeId)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("二维码不存在，QrCodeId=%v", req.QrCodeId)
			reply.ResultCode = ss_err.QrCodeNotExist
			reply.Msg = ss_err.GetMsg(ss_err.QrCodeNotExist, req.Lang)
			return nil
		}
		ss_log.Error("查询二维码失败，QrCodeId=%v, err=%v", req.QrCodeId, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	// 查询商户交易配置（商户账号，商户是否启用，收款权限）
	businessConf, err := dao.BusinessDaoInst.GetTransConfig(businessNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("商户不存在,BusinessNo:%v", businessNo)
			reply.ResultCode = ss_err.QrCodeNotInvalid
			reply.Msg = ss_err.GetMsg(ss_err.QrCodeNotInvalid, req.Lang)
			return nil
		}
		ss_log.Error("查询商户交易配置失败,err:%v,BusinessNo:%v", err, businessNo)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	resultCode, err = CheckBusinessIncomeAuth(businessConf)
	if err != nil {
		ss_log.Error("商户交易配置校验结果: %v", err)
		reply.ResultCode = resultCode
		reply.Msg = ss_err.GetMsg(resultCode, req.Lang)
		return nil
	}

	//查询商家是否签约了"商家收款"
	signed, err := dao.BusinessSceneSignedDaoInst.GetSignedByTradeType(businessNo, constants.TradeTypeBusinessPay)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("商家(%v)未签约,TradeType=%v", businessNo, constants.TradeTypeBusinessPay)
			reply.ResultCode = ss_err.QrCodeNotAvailable
			reply.Msg = ss_err.GetMsg(ss_err.QrCodeNotAvailable, req.Lang)
			return nil
		}
		ss_log.Error("查询商家(%v)签约失败，TradeType=%v, err=%v", businessNo, constants.TradeTypeBusinessPay, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	//校验签约
	retCode, err := CheckBusinessSigned(signed)
	if err != nil {
		ss_log.Error("校验结果：%v", err)
		reply.ResultCode = retCode
		reply.Msg = ss_err.GetMsg(retCode, req.Lang)
		return nil
	}

	//检查产品是否可用
	isEnabled, err := dao.BusinessSceneDao.GetSceneIsEnabled(signed.SceneNo)
	if err != nil {
		ss_log.Error("查询产品是否可用失败，SceneNo=%v, err=%v", signed.SceneNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	if !isEnabled {
		ss_log.Error("产品[%v]已被系统禁用，暂时不可交易", signed.SceneNo)
		reply.ResultCode = ss_err.SceneDisabled
		reply.Msg = ss_err.GetMsg(ss_err.SceneDisabled, req.Lang)
		return nil
	}

	if req.AccountNo == businessConf.BusinessAccNo {
		ss_log.Error("不能自己付款给自己，存在刷单行为，account=%v", req.AccountNo)
		reply.ResultCode = ss_err.PayeeAccountErr
		reply.Msg = ss_err.GetMsg(ss_err.PayeeAccountErr, req.Lang)
		return nil
	}

	//用户余额是否足够支付
	custVaType := global.GetUserVAccType(req.CurrencyType, true)
	custVAccNo, err := dao.VaccountDaoInst.GetVaccountNo(req.AccountNo, custVaType)
	if err != nil {
		ss_log.Error("查询用户虚账失败，AccountNo=%v, err=%v", req.AccountNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	balance, err := dao.VaccountDaoInst.GetBalanceByVAccNo(custVAccNo)
	if err != nil {
		ss_log.Error("查询用户余额失败，AccountNo=%v, err=%v", req.AccountNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	if strext.ToInt64(req.Amount) > balance {
		ss_log.Error("用户[%v]余额不足，req.Amount=%v, balance=%v", req.AccountNo, req.Amount, balance)
		reply.ResultCode = ss_err.BalanceNotEnough
		reply.Msg = ss_err.GetMsg(ss_err.BalanceNotEnough, req.Lang)
		return nil
	}

	//没有设置费率，需要设置为字符串0，方便后续统一处理
	if signed.Rate == "" {
		signed.Rate = "0"
	}

	// todo 计算手续费, 四舍五入后结果
	fee := ss_count.CountFees(req.Amount, signed.Rate, "0").String()
	realAmount := ss_count.Sub(req.Amount, fee).String()

	ss_log.Info("计算手续费结果,金额为:%s,费率为:%s,手续费为:%s,实际金额为:%s", req.Amount, signed.Rate, fee, realAmount)

	// 获取平台订单号
	order := dao.BusinessBillDao{
		OrderNo:           strext.GetDailyId(),
		Fee:               fee,
		Amount:            req.Amount,
		RealAmount:        realAmount,
		OrderStatus:       constants.BusinessOrderStatusPending,
		Rate:              signed.Rate,
		BusinessNo:        businessNo,
		BusinessAccountNo: businessConf.BusinessAccNo,
		CurrencyType:      req.CurrencyType,
		Subject:           req.Subject,
		Remark:            req.Remark,
		SceneNo:           signed.SceneNo,
		ExpireTime:        ss_time.Now(global.Tz).Add(constants.BusinessOrderExpireTime * time.Minute).Unix(),
		TradeType:         constants.TradeTypeBusinessPay,
		BusinessChannelNo: signed.BusinessChannelNo,
	}

	// 插入bill表记录
	if err := dao.BusinessBillDaoInst.InsertOrder(order); err != nil {
		if strings.Contains(err.Error(), ss_sql.DbDuplicateKey) {
			ss_log.Error("订单号重复,err:%v,orderInfo:%v", err, strext.ToJson(order))
		} else {
			ss_log.Error("插入订单失败,err:%v,orderInfo:%v", err, strext.ToJson(order))
		}
		reply.ResultCode = ss_err.PlaceAnOrderFail
		reply.Msg = ss_err.GetMsg(ss_err.PlaceAnOrderFail, req.Lang)
		return nil
	}

	reply.ResultCode = ss_err.Success
	reply.Msg = ss_err.GetMsg(ss_err.Success, req.Lang)
	reply.OrderNo = order.OrderNo
	reply.Amount = order.Amount
	reply.Subject = order.Subject
	return nil
}
