package handler

/*
import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/business-bill-srv/common"
	"a.a/mp-server/business-bill-srv/dao"
	"a.a/mp-server/business-bill-srv/i"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	authProto "a.a/mp-server/common/proto/auth"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	pushProto "a.a/mp-server/common/proto/push"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_func"
	"a.a/mp-server/common/ss_sql"
)

//商家中心-支付完成订单退款
func (b *BusinessBillHandler) BusinessBillRefund(ctx context.Context, req *businessBillProto.BusinessBillRefundRequest, reply *businessBillProto.BusinessBillRefundReply) error {
	//支付密码校验
	authReq := authProto.CheckPayPWDRequest{
		AccountUid:  req.BusinessAccNo,
		AccountType: req.AccountType,
		Password:    req.PaymentPwd,
		IdenNo:      req.BusinessNo,
		NonStr:      req.NonStr,
	}
	authRet, authErr := i.AuthHandlerInst.Client.CheckPayPWD(context.TODO(), &authReq)
	if authErr != nil {
		ss_log.Error("调用auth-srv服务失败,err:%v", authErr)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	authRet.ResultCode = ss_err.AuthSrvRetCode(authRet.ResultCode)
	if authRet.ResultCode != ss_err.Success {
		ss_log.Error("支付密码校验失败, resultCode=%v, CheckPayPWDRequest=%v", authRet.ResultCode, strext.ToJson(authReq))
		reply.ResultCode = authRet.ResultCode
		reply.Msg = ss_err.GetMsg(authRet.ResultCode, req.Lang)
		return nil
	}

	//查询订单详情
	order, err := dao.BusinessBillDaoInst.GetBillByLogNoAndBusinessNo(req.OrderNo, req.BusinessNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("订单不存在，orderNo=%v", req.OrderNo)
			reply.ResultCode = ss_err.OrderNotExist
			reply.Msg = ss_err.GetMsg(ss_err.OrderNotExist, req.Lang)
			return nil
		}
		ss_log.Error("查询订单失败，orderNo=%v， err=%v", req.OrderNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	if !util.InSlice(order.OrderStatus, []string{constants.BusinessOrderStatusPay, constants.BusinessOrderStatusRebatesRefund}) {
		ss_log.Error("订单状态不对，不能退款")
		reply.ResultCode = ss_err.OrderNotRefundable
		reply.Msg = ss_err.GetMsg(ss_err.OrderNotRefundable, req.Lang)
		return nil
	}

	//当前只支持全额退款
	if req.RefundAmount != order.Amount {
		ss_log.Error("退款金额与交易订单金额不一致，refundAmount=%v, transAmount=%v", req.RefundAmount, order.Amount)
		reply.ResultCode = ss_err.RefundAmountDisagree
		reply.Msg = ss_err.GetMsg(ss_err.RefundAmountDisagree, req.Lang)
		return nil
	}

	payeeAcc, err := dao.AccountDaoInst.GetAccountById(order.AccountNo)
	if err != nil {
		ss_log.Error("查询退款接收账号失败，AccountNo=%v, err=%v", order.AccountNo, err)
		reply.ResultCode = ss_err.PayeeAccountNoNotExist
		reply.Msg = ss_err.GetMsg(ss_err.PayeeAccountNoNotExist, req.Lang)
		return nil
	}

	p := &RefundParam{
		RefundAmount:     req.RefundAmount,
		RefundReason:     req.RefundReason,
		PayOrderNo:       order.OrderNo,
		PayOrderStatus:   order.OrderStatus,
		PayAmount:        order.Amount,
		PayFee:           order.Fee,
		PayRealAmount:    order.RealAmount,
		CurrencyType:     order.CurrencyType,
		SettleId:         order.SettleId,
		BusinessNo:       order.BusinessNo,
		BusinessUnVAccNo: order.BusinessVaccountNo,
		CustAccNo:        order.AccountNo,
	}
	resp, err := b.Refund(ctx, p)
	if err != nil {
		ss_log.Error("%v", err)
	}

	//添加成功退款消息提示（推给用户）
	if resp.ResultCode == ss_err.Success {
		toAccountNo := order.AccountNo

		//添加退款到的账号推送消息
		msg := new(dao.LogAppMessagesDao)
		msg.OrderNo = resp.RefundNo
		msg.AppMessType = constants.LOG_APP_MESSAGES_ORDER_TYPE_Business_Refund_Cust
		msg.OrderType = constants.VaReason_BusinessRefund
		msg.AccountNo = toAccountNo
		msg.OrderStatus = constants.OrderStatus_Paid
		err = dao.LogAppMessagesDaoInst.AddLogAppMessages(msg)
		if err != nil {
			ss_log.Error("aAddMessagesErr=[%v]", err)
		}

		ss_log.Info("用户 %s 当前的语言为--->%s", toAccountNo, req.Lang)
		moneyType := dao.LangDaoInstance.GetLangTextByKey(strings.ToLower(order.CurrencyType), req.Lang)
		timeString := time.Now().Format("2006-01-02 15:04:05")
		// 修正各币种的金额
		amount := common.NormalAmountByMoneyType(order.CurrencyType, resp.RefundAmount)

		args := []string{
			timeString, amount, moneyType,
		}

		if req.Lang == constants.LangEnUS {
			args = []string{
				amount, moneyType, timeString,
			}
		}

		// 消息推送
		ev := &pushProto.PushReqest{
			Accounts: []*pushProto.PushAccout{
				{
					AccountNo:   toAccountNo,
					AccountType: constants.AccountType_USER,
				},
			},
			TempNo: constants.Template_RefundSuccess,
			Args:   args,
		}

		ss_log.Info("publishing %+v\n", ev)
		// publish an event
		if err := common.PushEvent.Publish(context.TODO(), ev); err != nil {
			ss_log.Error("消息推送到用户toAccountNo[%v]出错。error : %v", toAccountNo, err)
		}
	}

	reply.ResultCode = resp.ResultCode
	reply.RefundNo = resp.RefundNo
	reply.RefundAmount = resp.RefundAmount
	reply.RefundStatus = resp.RefundStatus
	reply.RefundPayeeAcc = payeeAcc
	return nil
}

//接口-支付完成订单退款
func (b *BusinessBillHandler) ApiPayRefund(ctx context.Context, req *businessBillProto.ApiPayRefundRequest, reply *businessBillProto.ApiPayRefundReply) error {
	resultCode, err := CheckApiPayRefundReq(req)
	if err != nil {
		ss_log.Error("参数有误, err=%v", err)
		reply.ResultCode = resultCode
		return nil
	}

	//查询订单详情
	order, err := dao.BusinessBillDaoInst.AppQueryOrder(req.AppId, req.OrderNo, req.OutOrderNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("订单不存在，orderNo=%v", req.OrderNo)
			reply.ResultCode = ss_err.OrderNotExist
			reply.Msg = ss_err.GetMsg(ss_err.OrderNotExist, req.Lang)
			return nil
		}
		ss_log.Error("查询订单失败，orderNo=%v， err=%v", req.OrderNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	if !util.InSlice(order.OrderStatus, []string{constants.BusinessOrderStatusPay, constants.BusinessOrderStatusRebatesRefund}) {
		ss_log.Error("订单状态不满足退款条件")
		switch order.OrderStatus {
		case constants.BusinessOrderStatusPending:
			fallthrough
		case constants.BusinessOrderStatusPayTimeOut:
			reply.ResultCode = ss_err.OrderUnpaid
		case constants.BusinessOrderStatusRefund:
			reply.ResultCode = ss_err.OrderFullRefund
		}
		reply.Msg = ss_err.GetMsg(reply.ResultCode, req.Lang)
		return nil
	}
	//
	//payeeAcc, err := dao.AccountDaoInst.GetAccountById(order.AccountNo)
	//if err != nil {
	//	ss_log.Error("查询退款接收账号失败，AccountNo=%v, err=%v", order.AccountNo, err)
	//	reply.ResultCode = ss_err.PayeeAccountNoNotExist
	//	reply.Msg = ss_err.GetMsg(ss_err.PayeeAccountNoNotExist, req.Lang)
	//	return nil
	//}

	p := &RefundParam{
		OutRefundNo:      req.OutRefundNo,
		RefundAmount:     req.RefundAmount,
		RefundReason:     req.RefundReason,
		PayOrderNo:       order.OrderNo,
		PayOrderStatus:   order.OrderStatus,
		PayAmount:        order.Amount,
		PayFee:           order.Fee,
		PayRealAmount:    order.RealAmount,
		CurrencyType:     order.CurrencyType,
		SettleId:         order.SettleId,
		BusinessNo:       order.BusinessNo,
		BusinessUnVAccNo: order.BusinessVaccountNo,
		CustAccNo:        order.AccountNo,
	}
	resp, err := b.Refund(ctx, p)
	if err != nil {
		ss_log.Error("%v", err)
	}
	if resp.ResultCode == ss_err.Success {
		toAccountNo := order.AccountNo
		//添加退款到的账号推送消息
		msg := new(dao.LogAppMessagesDao)
		msg.OrderNo = resp.RefundNo
		msg.AppMessType = constants.LOG_APP_MESSAGES_ORDER_TYPE_Business_Refund_Cust
		msg.OrderType = constants.VaReason_BusinessRefund
		msg.AccountNo = toAccountNo
		msg.OrderStatus = constants.OrderStatus_Paid
		aAddMessagesErr := dao.LogAppMessagesDaoInst.AddLogAppMessages(msg)
		if aAddMessagesErr != nil {
			ss_log.Error("errAddMessages2=[%v]", aAddMessagesErr)
		}
		//推给用户成功退款消息
		pushMessageToUser(order.AccountNo, req.RefundAmount, order.CurrencyType, constants.Template_RefundSuccess, req.Lang)

	}

	reply.ResultCode = resp.ResultCode
	reply.OrderNo = order.OrderNo
	reply.OutOrderNo = order.OutOrderNo
	reply.RefundNo = resp.RefundNo
	reply.OutRefundNo = req.OutRefundNo
	reply.RefundAmount = resp.RefundAmount
	reply.RefundStatus = resp.RefundStatus
	reply.RealRefundAmount = resp.RefundAmount
	return nil
}

type RefundParam struct {
	OutRefundNo      string
	RefundAmount     string
	RefundReason     string
	PayOrderNo       string
	PayOrderStatus   string
	PayAmount        string
	PayFee           string
	PayRealAmount    string
	CurrencyType     string
	SettleId         string
	BusinessNo       string
	BusinessUnVAccNo string
	CustAccNo        string
}

type RefundResp struct {
	ResultCode   string
	RefundNo     string
	RefundAmount string
	RefundStatus string
}

func (b *BusinessBillHandler) Refund(ctx context.Context, p *RefundParam) (*RefundResp, error) {
	var resp = new(RefundResp)

	//获取退款金额
	ret, err := getRefundAmount(p.RefundAmount, p.PayAmount, p.PayOrderNo)
	if err != nil {
		ss_log.Error("退款金额有误，err=%v", err)
		resp.ResultCode = ret.ResultCode
		return resp, err
	}
	//退款金额
	ss_log.Info("退款金额：%v", strext.ToJson(ret))

	//查询支付订单的付款人虚账
	userVaType := global.GetUserVAccType(p.CurrencyType, true)
	userVAccNo, err := dao.VaccountDaoInst.GetVaccountNo(p.CustAccNo, userVaType)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("用户虚账不存在，AccountNo=%v", p.CustAccNo)
			resp.ResultCode = ss_err.VirtualAccountNotExist
			return resp, err
		}
		ss_log.Error("查询用户虚账失败，AccountNo=%v, err=%v", p.CustAccNo, err)
		resp.ResultCode = ss_err.SystemErr
		return resp, err
	}

	//查商户虚账
	var businessVAccNo, deductAmount string
	if p.SettleId == "" {
		//未结算
		businessVAccNo = p.BusinessUnVAccNo
		deductAmount = p.PayAmount
	} else {
		//已结算
		deductAmount = p.PayRealAmount

		businessVaType := global.GetBusinessVAccType(p.CurrencyType, true)
		businessAccNo, err := dao.BusinessDaoInst.QueryAccNoByBusinessNo(p.BusinessNo)
		businessVAccNo, err = dao.VaccountDaoInst.GetVaccountNo(businessAccNo, businessVaType)
		if err != nil {
			if err == sql.ErrNoRows {
				ss_log.Error("商户虚账不存在，BusinessAccNo=%v", businessAccNo)
				resp.ResultCode = ss_err.VirtualAccountNotExist
				return resp, err
			}
			ss_log.Error("查询商户虚账失败，BusinessAccNo=%v, err=%v", businessAccNo, err)
			resp.ResultCode = ss_err.SystemErr
			return resp, err
		}
	}

	balance, err := dao.VaccountDaoInst.GetBalanceByVAccNo(businessVAccNo)
	if err != nil {
		ss_log.Error("查询商户未结算虚账余额失败, BusinessVAccNo=%v, err=%v", businessVAccNo, err)
		resp.ResultCode = ss_err.SystemErr
		return resp, err
	}

	//判断余额是否足够退款， 计算结果小于0，Sign()返回-1；等于0，返回0；大于0，返回1
	if ss_count.Sub(strext.ToString(balance), deductAmount).Sign() < 0 {
		ss_log.Error("商户虚账余额不足, BusinessVAccNo=%v, CurrencyType=%v, balance=%v", businessVAccNo, p.CurrencyType, balance)
		resp.ResultCode = ss_err.BalanceNotEnough
		return resp, err
	}

	refundParam := new(AmountRollback)
	refundParam.BusinessVAccNo = businessVAccNo
	refundParam.CustVAccNo = userVAccNo
	refundParam.Amount = deductAmount
	refundParam.CurrencyType = p.CurrencyType
	if p.SettleId != "" {
		//已结算订单退款金额
		//查询平台虚账
		//todo 目前已结算订单考虑全额退款
		if ret.OrderStatus != constants.BusinessOrderStatusRefund {
			ss_log.Error("目前已结算订单不支持部分退款")
			resp.ResultCode = ss_err.OrderNotRefundable
			return resp, errors.New("目前已结算订单不支持部分退款")

		} else {
			platAccNo, err := dao.GlobalParamDaoInstance.GetParamValue(constants.GlobalParamKeyAccPlat)
			if err != nil {
				ss_log.Error("查询平台账号失败，paramKey=%v, err=%v", constants.GlobalParamKeyAccPlat, platAccNo)
				resp.ResultCode = ss_err.SystemErr
				return resp, err
			}

			platVAccType := global.GetPlatFormVAccType(p.CurrencyType)
			platVAccNo, err := dao.VaccountDaoInst.GetVaccountNo(platAccNo, platVAccType)
			if err != nil {
				ss_log.Error("查询平台虚账失败，paramKey=%v, err=%v", constants.GlobalParamKeyAccPlat, platAccNo)
				resp.ResultCode = ss_err.SystemErr
				return resp, err
			}

			refundParam.PlatformVAccNo = platVAccNo
			refundParam.Fee = p.PayFee
		}
	}

	//插入退款订单日志
	d := new(dao.BusinessRefundOrderDao)
	d.Amount = refundParam.Amount
	d.PayOrderNo = p.PayOrderNo
	d.RefundStatus = constants.BusinessRefundStatusPending
	d.Remarks = p.RefundReason
	d.OutRefundNo = p.OutRefundNo
	refundOrderNo, err := dao.BusinessRefundOrderDaoInst.Insert(d)
	if err != nil {
		ss_log.Error("插入退款日志失败，data:=%v, err=%v", strext.ToJson(d), err)
		resp.ResultCode = ss_err.SystemErr
		return resp, err
	}

	//开启事务
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, txErr := dbHandler.BeginTx(ctx, nil)
	if txErr != nil {
		ss_log.Error("开启事务失败，err=%v", txErr)
		resp.ResultCode = ss_err.SystemErr
		return resp, err
	}

	defer func() {
		if resp.ResultCode == ss_err.Success {
			//退款成功，修改原支付订单状态
			err := dao.BusinessBillDaoInst.UpdateStatusByOrderNoTx(tx, ret.OrderStatus, p.PayOrderNo)
			if err != nil {
				ss_sql.Rollback(tx)
				ss_log.Error("修改支付订单状态失败，orderNo=%v, err=%v", p.PayOrderNo, err)
				resp.ResultCode = ss_err.SystemErr
				return
			}
		} else {
			//退款失败，修改退款订单状态
			resp.RefundStatus = constants.BusinessRefundStatusFail
			if err := dao.BusinessRefundOrderDaoInst.UpdatePendingOrderByOrderNoTx(tx, resp.RefundStatus, refundOrderNo); nil != err {
				ss_sql.Rollback(tx)
				ss_log.Error("修改退款订单状态失败，orderNo=%v, err=%v", refundOrderNo, err)
				resp.ResultCode = ss_err.SystemErr
				return
			}
		}
		ss_sql.Commit(tx)
	}()

	//退款-减扣虚账金额
	refundParam.RefundOrderNo = refundOrderNo
	resp.ResultCode, err = amountRollback(tx, refundParam)
	if err != nil {
		ss_sql.Rollback(tx)
		return resp, err
	}

	//退款成功，修改退款订单状态
	resp.RefundStatus = constants.BusinessRefundStatusSuccess
	if err := dao.BusinessRefundOrderDaoInst.UpdatePendingOrderByOrderNoTx(tx, resp.RefundStatus, refundOrderNo); nil != err {
		ss_sql.Rollback(tx)
		ss_log.Error("修改退款订单状态失败，orderNo=%v, err=%v", refundOrderNo, err)
		resp.ResultCode = ss_err.SystemErr
		return resp, err
	}

	resp.ResultCode = ss_err.Success
	resp.RefundNo = refundOrderNo
	resp.RefundAmount = refundParam.Amount
	return resp, err
}

//部分退款检查金额：已退金额 + 退款金额 <= 交易金额
type RefundAmount struct {
	Amount      string
	OrderStatus string
	ResultCode  string
}

//计算退款金额和退款状态
func getRefundAmount(refundAmount, transAmount, orderNo string) (*RefundAmount, error) {
	resp := new(RefundAmount)
	//订单已退金额
	totalAmount, err := dao.BusinessRefundOrderDaoInst.GetTotalAmountByPayOrderNo(orderNo)
	if err != nil {
		if err != sql.ErrNoRows {
			ss_log.Error("查询订单已退金额失败，orderNo=%v, err=%v", orderNo, err)
			resp.ResultCode = ss_err.SystemErr
			return resp, err
		}
	}

	//退款金额大
	subRet := ss_count.Sub(transAmount, refundAmount).Sign()
	if subRet == -1 {
		resp.ResultCode = ss_err.RefundAmountExcessBalance
		return resp, errors.New(fmt.Sprintf("退款金额大于订单金额|交易金额:%v|退款金额:%v", transAmount, refundAmount))
	}

	//金额参数相等
	if subRet == 0 {
		//退还剩余全部金额
		if ss_count.Sub(transAmount, ss_count.Add(refundAmount, totalAmount)).Sign() < 0 {
			resp.ResultCode = ss_err.RefundAmountExcessBalance
			return resp, errors.New(fmt.Sprintf("退款金额大于订单剩余金额|交易金额:%v|退款金额:%v|已退金额:%v", transAmount, refundAmount, totalAmount))
		}
		resp.ResultCode = ss_err.Success
		resp.Amount = refundAmount
		resp.OrderStatus = constants.BusinessOrderStatusRefund
		return resp, nil
	}

	//退款金额小
	if subRet == 1 {
		//退剩下的所有金额
		if refundAmount == "" {
			resp.ResultCode = ss_err.Success
			resp.Amount = ss_count.Sub(transAmount, totalAmount).String()
			resp.OrderStatus = constants.BusinessOrderStatusRefund
			return resp, nil
		}

		//退还部分金额
		subRet = ss_count.Sub(transAmount, ss_count.Add(refundAmount, totalAmount)).Sign()
		if subRet < 0 {
			//金额有误
			resp.ResultCode = ss_err.RefundAmountExcessBalance
			return resp, errors.New(fmt.Sprintf("退款金额超出可退金额|交易金额:%v|退款金额:%v|已退金额:%v", transAmount, refundAmount, totalAmount))
		} else if subRet == 0 {
			//金额退完
			resp.ResultCode = ss_err.Success
			resp.Amount = refundAmount
			resp.OrderStatus = constants.BusinessOrderStatusRefund
		} else {
			//金额有剩余
			resp.ResultCode = ss_err.Success
			resp.Amount = refundAmount
			resp.OrderStatus = constants.BusinessOrderStatusRebatesRefund
		}
	}

	return resp, nil
}

type AmountRollback struct {
	PlatformVAccNo string
	BusinessVAccNo string
	CustVAccNo     string
	Amount         string
	CurrencyType   string
	Fee            string
	RefundOrderNo  string
}

//金额回退
func amountRollback(tx *sql.Tx, req *AmountRollback) (string, error) {
	//平台退还手续费
	if req.PlatformVAccNo != "" && strext.ToInt64(req.Fee) > 0 {
		//减少商平台虚账的余额
		platformBalance, platformFrozenBalance, err := dao.VaccountDaoInst.MinusBalance(tx, req.PlatformVAccNo, req.Fee)
		if err != nil {
			ss_sql.Rollback(tx)
			ss_log.Error("减少平台账户余额, PlatformVAccNo=%v, err=%v", req.PlatformVAccNo, err)
			return ss_err.SystemErr, err
		}

		// 判断余额是否为负数
		if r, err := ss_func.JudgeAmountPositiveOrNegative(platformBalance); err != nil || r < 0 {
			ss_sql.Rollback(tx)
			ss_log.Error("平台账户余额不足,err:%v,platformBalance:%v, result:%v", err, platformBalance, r)
			return ss_err.BalanceNotEnough, err
		}

		//记录账户变动日志
		log1 := dao.LogVaccountDao{
			VaccountNo:    req.PlatformVAccNo,
			OpType:        constants.VaOpType_Minus,
			Amount:        req.Fee,
			Balance:       platformBalance,
			FrozenBalance: platformFrozenBalance,
			Reason:        constants.VaReason_BusinessRefund,
			BizLogNo:      req.RefundOrderNo,
		}
		if err := dao.LogVaccountDaoInst.InsertLogTx(tx, log1); err != nil {
			ss_sql.Rollback(tx)
			ss_log.Error("插入平台账户变动日志失败,err:%v,data:%+v", err, log1)
			return ss_err.SystemErr, err
		}

		//插入手续费盈利扣减记录
		d := &dao.HeadquartersProfit{
			GeneralLedgerNo: req.RefundOrderNo,
			Amount:          req.Fee,
			OrderStatus:     constants.OrderStatus_Paid,
			BalanceType:     strings.ToLower(req.CurrencyType),
			ProfitSource:    constants.ProfitSource_ModernPayOrderRefund,
			OpType:          constants.PlatformProfitMinus,
		}
		_, err = dao.HeadquartersProfitDao.InsertHeadquartersProfit(tx, d)
		if err != nil {
			ss_sql.Rollback(tx)
			ss_log.Error("插入平台盈利失败, data=%v, err=%v", strext.ToJson(d), err)
			return ss_err.SystemErr, err
		}

		//同步总部虚账的余额(等于收益表中的可提现余额)
		err = dao.HeadquartersProfitDao.SyncHeadquartersProfit(tx, req.PlatformVAccNo, req.Fee, strings.ToLower(req.CurrencyType))
		if err != nil {
			ss_sql.Rollback(tx)
			ss_log.Error("同步收益余额失败, headVacc=%v, err=%v", req.PlatformVAccNo, err)
			return ss_err.SystemErr, err
		}
	}

	//减少商家未结算虚账的余额
	businessBalance, businessFrozenBalance, err := dao.VaccountDaoInst.MinusBalance(tx, req.BusinessVAccNo, req.Amount)
	if err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("减少商户账户余额, businessVAccNo=%v, err=%v", req.BusinessVAccNo, err)
		return ss_err.SystemErr, err
	}

	// 判断余额是否为负数
	if r, err := ss_func.JudgeAmountPositiveOrNegative(businessBalance); err != nil || r < 0 {
		ss_sql.Rollback(tx)
		ss_log.Error("商户账户余额不足,err:%v,businessBalance:%v, result:%v", err, businessBalance, r)
		return ss_err.BalanceNotEnough, err
	}

	//记录账户变动日志
	log1 := dao.LogVaccountDao{
		VaccountNo:    req.BusinessVAccNo,
		OpType:        constants.VaOpType_Minus,
		Amount:        req.Amount,
		Balance:       businessBalance,
		FrozenBalance: businessFrozenBalance,
		Reason:        constants.VaReason_BusinessRefund,
		BizLogNo:      req.RefundOrderNo,
	}
	if err := dao.LogVaccountDaoInst.InsertLogTx(tx, log1); err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("插入商家账户变动日志失败,err:%v,data:%+v", err, log1)
		return ss_err.SystemErr, err
	}

	//增加用户虚账余额，并记录账户变动日志
	userAmount := ss_count.Add(req.Amount, req.Fee)
	userBalance, userFrozenBalance, err := dao.VaccountDaoInst.PlusBalance(tx, req.CustVAccNo, userAmount)
	if err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("增加用户虚账余额失败，CustVAccNo=%v, err=%v", req.CustVAccNo, err)
		return ss_err.SystemErr, err
	}

	//记录账户变动日志
	log2 := dao.LogVaccountDao{
		VaccountNo:    req.CustVAccNo,
		OpType:        constants.VaOpType_Add,
		Amount:        userAmount,
		Balance:       userBalance,
		FrozenBalance: userFrozenBalance,
		Reason:        constants.VaReason_BusinessRefund,
		BizLogNo:      req.RefundOrderNo,
	}
	if err := dao.LogVaccountDaoInst.InsertLogTx(tx, log2); err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("插入用户账户变动日志失败,err:%v,data:%+v", err, log2)
		return ss_err.SystemErr, err
	}

	return ss_err.Success, nil
}
*/
