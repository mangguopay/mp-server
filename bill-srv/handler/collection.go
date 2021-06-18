package handler

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/bill-srv/common"
	"a.a/mp-server/bill-srv/dao"
	"a.a/mp-server/bill-srv/i"
	"a.a/mp-server/common/constants"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	go_micro_srv_push "a.a/mp-server/common/proto/push"
	go_micro_srv_riskctrl "a.a/mp-server/common/proto/riskctrl"
	go_micro_srv_settle "a.a/mp-server/common/proto/settle"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"context"
	"time"
)

// 收款
func (b BillHandler) Collection(ctx context.Context, req *go_micro_srv_bill.CollectionRequest, reply *go_micro_srv_bill.CollectionReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// 判断金额
	if strext.ToFloat64(req.Amount) <= 0 || req.Amount == "" {
		ss_log.Error("err=[收款, 金额为0或者为空,传入的金额为----->%s]", req.Amount)
		reply.ResultCode = ss_err.ERR_WALLET_AMOUNT_NULL
		return nil
	}

	// 校验扫码人的支付密码
	if errStr := checkoutPWD(req.AccountType, req.OpAccNo, req.NonStr, req.Password); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	//  校验参数
	accountNo, amount, moneyType, _, useStatus := dao.GenCodeDaoInst.GetRecvCode(req.GenCode, constants.CodeType_Recv)
	if req.Amount != amount {
		reply.ResultCode = ss_err.ERR_WRONG_AMOUNT
		return nil
	}
	if req.MoneyType != moneyType || accountNo != req.RecAccountUid {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	//码的状态不是被扫状态
	if useStatus != constants.CODE_USE_STATUS_IS_SWEEP {
		ss_log.Error("err=[------->%s]", "码的状态不是被扫状态")
		reply.ResultCode = ss_err.ERR_QUERY_SWEEP_CODE_STATUS_FAILD
		return nil
	}
	// 根据 币种 和 account_uid 获取虚拟账户id
	fromVAccount := dao.VaccountDaoInst.GetVaccountNoFromMoneyType(req.SweepAccountUid, req.MoneyType)
	if fromVAccount == "" {
		ss_log.Error("err=[付款人虚拟账户不存在,账户id为----->%s,moneyType为---->%s]", req.SweepAccountUid, req.MoneyType)
		reply.ResultCode = ss_err.ERR_PAY_VACC_ACCOUNT_NO_EXIST
		return nil
	}
	collectionVAccount := dao.VaccountDaoInst.GetVaccountNoFromMoneyType(req.RecAccountUid, req.MoneyType)
	if collectionVAccount == "" {
		ss_log.Error("err=[收款人虚拟账户不存在,账户id为----->%s,moneyType为---->%s]", req.RecAccountUid, req.MoneyType)
		reply.ResultCode = ss_err.ERR_COLLECTION_VACC_ACCOUNT_NO_EXIST
		return nil
	}

	// 手续费,属于转账的手续费,谁扫码是算谁的
	var feesType int32
	switch req.MoneyType {
	case "usd":
		feesType = 2
	case "khr":
		feesType = 7
	}

	//  获取手续费
	_, fees, feeErr := doFees(feesType, req.Amount)
	if feeErr != nil {
		ss_log.Error("收款 计算手续费失败,err:--->%s", feeErr.Error())
		reply.ResultCode = ss_err.ERR_RATE_FAILD
		return nil
	}
	ss_log.Info("币种为------->%s,金额为----->%s,手续费为-------->%s", req.MoneyType, req.Amount, fees)

	// 判断本金加手续费是否超额
	balance, _ := dao.VaccountDaoInst.GetBalance(fromVAccount)
	allAmount := ss_count.Add(req.Amount, fees) // 加上手续费后的amount
	if strext.ToInt64(balance) < strext.ToInt64(allAmount) {
		reply.ResultCode = ss_err.ERR_PAY_AMT_NOT_ENOUGH
		return nil
	}

	// 记录日志
	logNo := dao.CollectionDaoInst.InsertCollectionOrder(tx, fromVAccount, collectionVAccount, req.Amount, req.MoneyType, fees, req.Lat, req.Lng, req.Ip)
	if logNo == "" {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//风控
	riskReply, _ := i.RiskCtrlHandleInstance.Client.GetRiskCtrlReuslt(context.TODO(), &go_micro_srv_riskctrl.GetRiskCtrlResultRequest{
		ApiType: "collection",
		// 发起支付的账号
		PayerAccNo: req.SweepAccountUid,
		ActionTime: time.Now().String(),
		Amount:     req.Amount,
		Ip:         req.Ip,
		PayType:    constants.Risk_Pay_Type_Collection, //
		// 收款人账号
		PayeeAccNo:  req.RecAccountUid,
		ProductType: "collection",
		// 币种
		MoneyType: req.MoneyType,
		// 订单号
		OrderNo: logNo,
	})

	ss_log.Info("收款 风控返回结果,操作结果是---->%s,RiskNo为----->%s", riskReply.OpResult, riskReply.RiskNo)

	if riskReply.OpResult == constants.Risk_Result_No_Pass_Str {
		reply.ResultCode = ss_err.ERR_RISK_IS_RISK
		reply.RiskNo = riskReply.RiskNo
		return nil
	}

	// 更改金额
	if errCode := dao.VaccountDaoInst.AccFromAToBUpperZero(tx, fromVAccount, collectionVAccount, req.Amount, logNo, constants.VaReason_COLLECTION); errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", errCode)
		if errStr := dao.CollectionDaoInst.UpdateCollectionOrderStatus(tx, logNo, constants.OrderStatus_Err); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}
		reply.ResultCode = errCode
		return nil
	}

	if fees != "0" {
		if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, fromVAccount, fees, "-", logNo, constants.VaReason_TRANSFER); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}

		// 发送手续费进MQ
		feeEv := &go_micro_srv_settle.SettleTransferRequest{
			BillNo:    logNo,
			FeesType:  constants.FEES_TYPE_COLLECTION,
			Fees:      fees,
			MoneyType: req.MoneyType,
		}
		ss_log.Info("publishing %+v\n", feeEv)
		// publish an event
		if err := common.SettleEvent.Publish(context.TODO(), feeEv); err != nil {
			ss_log.Error("err=[收款接口,手续费推送到MQ失败,err----->%s]", err.Error())
		}
	}

	//if errStr := dao.VaccountDaoInst.SyncAccRemain(tx, req.SweepAccountUid); errStr != ss_err.ERR_SUCCESS {
	//	reply.ResultCode = errStr
	//	return nil
	//}
	//if errStr := dao.VaccountDaoInst.SyncAccRemain(tx, req.RecAccountUid); errStr != ss_err.ERR_SUCCESS {
	//	reply.ResultCode = errStr
	//	return nil
	//}
	// 下单成功
	if errStr := dao.CollectionDaoInst.UpdateCollectionOrderStatus(tx, logNo, constants.OrderStatus_Paid); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	//添加扫码人的付款推送消息
	errAddMessages := dao.LogAppMessagesDaoInst.AddLogAppMessagesTx(tx, logNo, constants.LOG_APP_MESSAGES_ORDER_TYPE_COLLECTION_Apply, constants.VaReason_COLLECTION, req.SweepAccountUid, constants.OrderStatus_Paid)
	if errAddMessages != ss_err.ERR_SUCCESS {
		ss_log.Error("errAddMessages=[%v]", errAddMessages)
		reply.ResultCode = ss_err.ERR_SYS_DB_ADD
		return nil
	}
	//添加收款人的收款推送消息
	errAddMessages2 := dao.LogAppMessagesDaoInst.AddLogAppMessagesTx(tx, logNo, constants.LOG_APP_MESSAGES_ORDER_TYPE_COLLECTION_Apply, constants.VaReason_COLLECTION, req.RecAccountUid, constants.OrderStatus_Paid)
	if errAddMessages2 != ss_err.ERR_SUCCESS {
		ss_log.Error("errAddMessages2=[%v]", errAddMessages2)
		reply.ResultCode = ss_err.ERR_SYS_DB_ADD
		return nil
	}

	recAccountType := dao.AccDaoInstance.GetAccountTypeFromAccNo(req.RecAccountUid)
	moneyTypeB := dao.LangDaoInstance.GetLangTextByKey(req.MoneyType, req.Lang)
	amountB := ""
	switch req.MoneyType {
	case "usd":
		amountB = strext.ToStringNoPoint(strext.ToFloat64(req.Amount) / 100)
	case "khr":
		amountB = strext.ToStringNoPoint(strext.ToFloat64(req.Amount))
	}

	// 消息推送
	ev := &go_micro_srv_push.PushReqest{
		Accounts: []*go_micro_srv_push.PushAccout{
			{
				AccountNo:   req.RecAccountUid,
				AccountType: recAccountType,
			},
		},
		TempNo: constants.Template_TransferSuccess,
		Args: []string{
			time.Now().Format("2006-01-02 15:04:05"), amountB, moneyTypeB,
		},
	}
	ss_log.Info("publishing %+v\n", ev)
	// publish an event
	if err := common.PushEvent.Publish(context.TODO(), ev); err != nil {
		ss_log.Error("error publishing: %v", err)
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.OrderNo = logNo
	reply.RiskNo = riskReply.RiskNo
	return nil
}
