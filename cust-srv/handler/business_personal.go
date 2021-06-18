package handler

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/dao"
	"context"
	"database/sql"
	"fmt"
	"strings"

	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
)

func (*CustHandler) GetPersonalBusinessInfo(ctx context.Context, req *custProto.GetPersonalBusinessInfoRequest, reply *custProto.GetPersonalBusinessInfoReply) error {
	if req.AccountNo == "" {
		ss_log.Error("AccountNo参数缺失")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	accInfo, err := dao.AccDaoInstance.GetUserAccountInfo(req.AccountNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("账号[%v]未实名认证", req.AccountNo)
			reply.ResultCode = ss_err.ERR_ACCOUNT_NOT_REAL_AUTH
			return nil
		}
		ss_log.Error("查询账号[%v]实名认证信息失败，err=%v", req.AccountNo, err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if accInfo.AuthStatus != constants.AuthMaterialStatus_Passed {
		ss_log.Error("账号[%v]未实名认证", req.AccountNo)
		reply.ResultCode = ss_err.ERR_ACCOUNT_NOT_REAL_AUTH
		return nil
	}

	data := &custProto.PersonalBusinessInfo{
		Account:          accInfo.Account,
		RealName:         accInfo.AuthName,
		BusinessName:     "",
		SimplifyName:     "",
		OperatingPeriod:  "",
		OrganizationCode: "",
		LicenseImg:       "",
	}

	business, err := dao.AccDaoInstance.GetBusinessAccountInfo(req.AccountNo)
	if err != nil && err != sql.ErrNoRows {
		ss_log.Error("查询账号[%v]，商家认证资料失败，err=%v", req.AccountNo, err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	data.BusinessName = business.BusinessName
	data.SimplifyName = business.SimplifyName
	data.BusinessAuthStatus = business.AuthStatus
	data.OperatingPeriod = fmt.Sprintf("%s——%s", business.StartDate, business.EndDate)
	data.OrganizationCode = business.AuthNumber
	data.LicenseImg = strext.ToStringNoPoint(business.LicenseImgNo != "")

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

func (*CustHandler) GetPersonalBusinessBalance(ctx context.Context, req *custProto.GetPersonalBusinessBalanceRequest, reply *custProto.GetPersonalBusinessBalanceReply) error {
	if req.AccountNo == "" {
		ss_log.Error("AccountNo参数缺失")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//查询用户是否有个人商家这重身份
	businessNo, _, err := dao.BusinessDaoInst.GetBusinessNo(req.AccountNo, constants.AccountType_PersonalBusiness)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("账号[%v]没有对应的商户号", req.AccountNo)
			reply.ResultCode = ss_err.ERR_NOT_PERSONAL_BUSINESS
			return nil
		}
		ss_log.Error("查询商户号失败，accountNo=%v, err=%v", req.AccountNo, err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//查询个人商家余额，包含美金，瑞尔，已结算，没结算
	businessBalance, err := dao.VaccountDaoInst.GetAllVAccBalanceByAccNo(req.AccountNo)
	if err != nil {
		ss_log.Error("查询个人商家[%v]余额失败, err=%v", req.AccountNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//查询usd的今日收款数和今日收款金额
	usdNum, usdAmount, err := dao.BusinessBillDaoInst.GetTodayPaySum(businessNo, constants.CURRENCY_UP_USD)
	if err != nil {
		ss_log.Error("个人商家[%v]USD余额失败, err=%v", businessNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//查询khr的今日收款数和今日收款金额
	khrNum, khrAmount, err := dao.BusinessBillDaoInst.GetTodayPaySum(businessNo, constants.CURRENCY_UP_KHR)
	if err != nil {
		ss_log.Error("查询个人商家[%v]余额失败, err=%v", req.AccountNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	usdData := new(custProto.PersonalBusiness)
	usdData.CurrencyType = constants.CURRENCY_UP_USD
	usdData.CollectionAmount = usdAmount
	usdData.CollectionNum = usdNum

	khrData := new(custProto.PersonalBusiness)
	khrData.CurrencyType = constants.CURRENCY_UP_KHR
	khrData.CollectionAmount = khrAmount
	khrData.CollectionNum = khrNum

	for _, v := range businessBalance {
		if strings.ToUpper(v.CurrencyType) == constants.CURRENCY_UP_USD {
			if v.VAccType == constants.VaType_USD_BUSINESS_SETTLED {
				usdData.AccountBalance = v.Balance
			} else if v.VAccType == constants.VaType_USD_BUSINESS_UNSETTLED {
				usdData.NoSettleAmount = v.Balance
			}
		} else if strings.ToUpper(v.CurrencyType) == constants.CURRENCY_UP_KHR {
			if v.VAccType == constants.VaType_KHR_BUSINESS_SETTLED {
				khrData.AccountBalance = v.Balance
			} else if v.VAccType == constants.VaType_KHR_BUSINESS_UNSETTLED {
				khrData.NoSettleAmount = v.Balance
			}
		}
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = []*custProto.PersonalBusiness{khrData, usdData}
	return nil
}

func (*CustHandler) GetPersonalBusinessBills(ctx context.Context, req *custProto.GetPersonalBusinessBillsRequest, reply *custProto.GetPersonalBusinessBillsReply) error {
	if req.AccountNo == "" {
		ss_log.Error("AccountNo参数缺失")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//查询用户是否有个人商家这重身份
	businessNo, _, err := dao.BusinessDaoInst.GetBusinessNo(req.AccountNo, constants.AccountType_PersonalBusiness)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("账号[%v]没有对应的商户号", req.AccountNo)
			reply.ResultCode = ss_err.ERR_NOT_PERSONAL_BUSINESS
			return nil
		}
		ss_log.Error("查询商户号失败，accountNo=%v, err=%v", req.AccountNo, err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "bb.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "bb.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "bb.currency_type", Val: req.CurrencyType, EqType: "="},
		{Key: "bb.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "bb.business_no", Val: businessNo, EqType: "="},
	})

	total, err := dao.BusinessBillDaoInst.GetPersonalBusinessBillsCnt(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("统计个人商户订单数量失败，req=%v, err=%v", strext.ToJson(req), err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY bb.create_time desc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimitI32(whereModel, req.PageSize, req.Page)
	datas, err := dao.BusinessBillDaoInst.GetPersonalBusinessBills(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("查询个人商户订单列表失败，req=%v, err=%v", strext.ToJson(req), err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	var list []*custProto.BusinessBillData
	for _, v := range datas {
		data := &custProto.BusinessBillData{
			CreateTime:   v.CreateTime,
			PayTime:      v.PayTime,
			OrderNo:      v.OrderNo,
			Subject:      v.Subject,
			Amount:       v.Amount,
			CurrencyType: v.CurrencyType,
			OrderStatus:  v.OrderStatus,
		}
		if v.SettleId != "" {
			data.IsSettled = "1"
		} else {
			data.IsSettled = "0"
		}

		list = append(list, data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Total = strext.ToInt32(total)
	reply.List = list
	return nil
}

func (*CustHandler) GetPersonalBusinessBillDetail(ctx context.Context, req *custProto.GetPersonalBusinessBillDetailRequest, reply *custProto.GetPersonalBusinessBillDetailReply) error {
	if req.AccountNo == "" {
		ss_log.Error("AccountNo参数缺失")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.OrderNo == "" {
		ss_log.Error("OrderNo参数缺失")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	order, err := dao.BusinessBillDaoInst.GetPersonalBusinessBillDetail(req.OrderNo, req.AccountNo)
	if err != nil {
		ss_log.Error("查询个人商户订单失败，req=%v, err=%v", strext.ToJson(req), err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	data := &custProto.BusinessBillData{
		OrderNo:      order.OrderNo,
		Subject:      order.Subject,
		Amount:       order.Amount,
		RealAmount:   order.RealAmount,
		CurrencyType: order.CurrencyType,
		Fee:          order.Fee,
		OrderStatus:  order.OrderStatus,
		SettleId:     order.SettleId,
		SettleDate:   order.SettleDate,
		PayTime:      order.PayTime,
		SceneName:    order.SceneName,
	}

	if util.InSlice(order.OrderStatus, []string{constants.BusinessOrderStatusRefund, constants.BusinessOrderStatusRebatesRefund}) {
		refundOrder, err := dao.BusinessRefundOrderDaoInst.GetPersonalRefundOrderDetail(order.OrderNo, req.AccountNo)
		if err != nil {
			ss_log.Error("查询个人商户订单退款记录失败，req=%v, err=%v", strext.ToJson(req), err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		data.RefundedAmount = refundOrder.Amount
		data.RefundTime = refundOrder.FinishTime
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

func (*CustHandler) GetPersonalBusinessFixedCode(ctx context.Context, req *custProto.GetPersonalBusinessFixedCodeRequest, reply *custProto.GetPersonalBusinessFixedCodeReply) error {
	if req.AccountNo == "" {
		ss_log.Error("AccountNo参数缺失")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.AccountType != constants.AccountType_USER && req.AccountType != constants.AccountType_PersonalBusiness {
		ss_log.Error("账号类型错误")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//查询用户是否有个人商家这重身份
	businessNo, simplifyName, err := dao.BusinessDaoInst.GetBusinessNo(req.AccountNo, constants.AccountType_PersonalBusiness)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("账号[%v]没有对应的商户号", req.AccountNo)
			reply.ResultCode = ss_err.ERR_NOT_PERSONAL_BUSINESS
			return nil
		}
		ss_log.Error("查询商户号失败，accountNo=%v, err=%v", req.AccountNo, err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	fixedCode, err := dao.BusinessFixedCodeDaoInst.GetFixedCodeByBusinessNo(businessNo)
	if err != nil {
		ss_log.Error("查询商户固定二维码失败，businessNo=%v, err=%v", businessNo, err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.FixedCode = constants.GetPersonalBusinessFixedQrCodeUrl(fixedCode)
	reply.SimplifyName = simplifyName
	return nil
}
