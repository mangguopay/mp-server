package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_func"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
	"fmt"
)

type LogAppMessagesDao struct {
}

var LogAppMessagesDaoInst LogAppMessagesDao

//统计未读的消息数量
func (*LogAppMessagesDao) GetNoReadCnt(accountNo string) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select count(1) from log_app_messages where is_read = '0' and account_no = $1 "
	var totalT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, accountNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}
	return totalT.String
}

func (*LogAppMessagesDao) CustOrderBillDetail(accountNo, orderNo, orderType string) (data *go_micro_srv_bill.CustOrderBillDetailData, returnErr string) {

	data = &go_micro_srv_bill.CustOrderBillDetailData{}
	switch orderType {
	case constants.VaReason_Exchange: //兑换
		exchangeData, err := ExchangeOrderDaoInst.CustExchangeBillsDetail(orderNo)
		if err != ss_err.ERR_SUCCESS {
			ss_log.Error("getIncomeDataErr= [%v]", err)
			return nil, ss_err.ERR_SYS_DB_GET
		}
		data.LogNo = exchangeData.LogNo
		data.InType = exchangeData.InType
		data.OutType = exchangeData.OutType
		data.Amount = exchangeData.Amount
		data.CreateTime = exchangeData.CreateTime
		data.OrderStatus = exchangeData.OrderStatus
		data.FinishTime = exchangeData.FinishTime
		data.TransFrom = exchangeData.TransFrom     //app,trade
		data.TransAmount = exchangeData.TransAmount //转换后金额
		data.Fees = exchangeData.Fees
		data.BalanceType = exchangeData.InType
		data.PaymentType = constants.ORDER_PAYMENT_BALANCE

		data.OrderType = orderType
	case constants.VaReason_INCOME: //存款
		incomeData, err := IncomeOrderDaoInst.CustIncomeBillsDetail(orderNo)

		if err != ss_err.ERR_SUCCESS {
			ss_log.Error("getIncomeDataErr= [%v]", err)
			return nil, ss_err.ERR_SYS_DB_GET
		}
		data.OrderStatus = incomeData.OrderStatus
		data.CreateTime = incomeData.CreateTime
		data.Amount = incomeData.Amount
		data.BalanceType = incomeData.BalanceType
		data.PaymentType = incomeData.PaymentType
		data.FinishTime = incomeData.FinishTime
		data.LogNo = orderNo
		data.Fees = incomeData.Fees

		//data.OpType = incomeData.OpType
		switch incomeData.OrderStatus {
		case constants.OrderStatus_Paid:
			data.OpType = constants.VaOpType_Add
		case constants.OrderStatus_Err:
			data.OpType = constants.VaOpType_Minus
		default:
			ss_log.Error("orderStatus[%v],未处理,opType可能错误", incomeData.OrderStatus)
		}

		data.OrderType = orderType
	case constants.VaReason_OUTGO: //取款
		outgoData, err := OutgoOrderDaoInst.CustOutgoBillsDetail(orderNo)
		if err != ss_err.ERR_SUCCESS {
			ss_log.Error("err=%v", err)
			return nil, ss_err.ERR_SYS_DB_GET
		}
		data.OrderStatus = outgoData.OrderStatus
		data.CreateTime = outgoData.CreateTime
		data.Amount = outgoData.Amount
		data.BalanceType = outgoData.BalanceType
		data.PaymentType = outgoData.PaymentType
		data.FinishTime = outgoData.FinishTime
		data.LogNo = orderNo
		//data.OpType = constants.VaOpType_Minus
		data.Fees = outgoData.Fees

		switch outgoData.OrderStatus {
		case constants.OrderStatus_Paid: //已支付
			data.OpType = constants.VaOpType_Minus
		case constants.OrderStatus_Err: //失败
			fallthrough
		case constants.OrderStatus_Cancel: //虽然说取消提现的不显示订单，但这里还是给出了OpType的转换
			data.OpType = constants.VaOpType_Add
		case constants.OrderStatus_Pending_Confirm: //待确认统一为2，等待
			data.OrderStatus = constants.OrderStatus_Pending
			data.OpType = constants.VaOpType_Minus
		default:
			ss_log.Error("orderStatus[%v]未处理，opType可能错误。", outgoData.OrderStatus)
		}

		data.OrderType = orderType
	case constants.VaReason_TRANSFER: //转账
		transferData, errStr := TransferDaoInst.CustTransferBillsDetail(orderNo)
		if errStr != ss_err.ERR_SUCCESS {
			ss_log.Error("err=%v", errStr)
			return nil, ss_err.ERR_SYS_DB_GET
		}

		log, err := LogVaccountDaoInst.GetLogVAccountByBizLogNo(accountNo, orderNo, constants.VaReason_TRANSFER)
		if err != nil {
			ss_log.Error("查询虚账日志失败，err=%v", err)
			return nil, ss_err.ERR_SYSTEM
		}
		//收款方， 付款方
		var toPhone, fromPhone string
		if log.OpType == constants.VaOpType_Add {
			//手机号脱敏
			toPhone = fmt.Sprintf("%v-%v", transferData.ToPhoneCountryCode, transferData.ToPhone)
			fromPhone, err = ss_func.GetDesensitizationPhoneByCountryCode(transferData.FromPhoneCountryCode, transferData.FromPhone)
			if err != nil {
				ss_log.Error("手机号脱敏失败，CountryCode=%v, Phone=%v", transferData.ToPhoneCountryCode, transferData.ToPhone)
				return nil, ss_err.ERR_SYSTEM
			}

		} else if log.OpType == constants.VaOpType_Minus {
			//手机号脱敏
			toPhone, err = ss_func.GetDesensitizationPhoneByCountryCode(transferData.ToPhoneCountryCode, transferData.ToPhone)
			if err != nil {
				ss_log.Error("手机号脱敏失败，CountryCode=%v, Phone=%v", transferData.ToPhoneCountryCode, transferData.ToPhone)
				return nil, ss_err.ERR_SYSTEM
			}
			fromPhone = fmt.Sprintf("%v-%v", transferData.FromPhoneCountryCode, transferData.FromPhone)
		}

		data.OrderStatus = transferData.OrderStatus
		data.CreateTime = transferData.CreateTime
		data.Amount = transferData.Amount
		data.BalanceType = transferData.BalanceType
		data.PaymentType = transferData.PaymentType
		data.FinishTime = transferData.FinishTime
		data.LogNo = orderNo
		data.Fees = transferData.Fees
		data.ToPhone = toPhone
		data.FromPhone = fromPhone
		data.OrderType = orderType
	case constants.VaReason_COLLECTION: //收款
		collectionData, err := CollectionDaoInst.CustCollectionBillsDetail(orderNo)
		if err != ss_err.ERR_SUCCESS {
			ss_log.Error("err=%v", err)
			return nil, ss_err.ERR_SYS_DB_GET
		}

		data.OrderStatus = collectionData.OrderStatus
		data.CreateTime = collectionData.CreateTime
		data.Amount = collectionData.Amount
		data.BalanceType = collectionData.BalanceType
		data.PaymentType = collectionData.PaymentType
		data.FinishTime = collectionData.FinishTime
		data.LogNo = orderNo
		data.FromAccount = collectionData.FromAccount
		data.OpType = constants.VaOpType_Add

		data.OrderType = orderType
	case constants.VaReason_Cust_Cancel_Save: //银行卡存款
		fallthrough
	case constants.VaReason_Cust_Save:
		toHeadData, errStr := LogCustToHeadquartersDaoInst.CustToHeadquartersDetail(orderNo)
		if errStr != ss_err.ERR_SUCCESS {
			ss_log.Error("CustToHeadquartersDetailErr= [%v]", errStr)
			return nil, ss_err.ERR_SYS_DB_GET
		}
		//data.OrderStatus = toHeadData.OrderStatus

		//要审核的订单状态和不需审核的订单状态不同，而现在用户消息这里又审核和不需审核的混淆在一起，现进行转换，统一转换成不需审核的订单状态。
		switch toHeadData.OrderStatus {
		case constants.AuditOrderStatus_Pending:
			data.OrderStatus = constants.OrderStatus_Pending
		case constants.AuditOrderStatus_Passed:
			data.OpType = constants.VaOpType_Add
			data.OrderStatus = constants.OrderStatus_Paid
		case constants.AuditOrderStatus_Deny:
			data.OrderStatus = constants.OrderStatus_Err
		default:

		}
		data.CreateTime = toHeadData.CreateTime
		data.Amount = toHeadData.Amount
		data.BalanceType = toHeadData.BalanceType
		data.FinishTime = toHeadData.FinishTime
		data.LogNo = orderNo
		//data.OpType = toHeadData.OpType
		data.Fees = toHeadData.Fees
		data.PaymentType = toHeadData.PaymentType
		data.OrderType = orderType
		data.CardNumber = toHeadData.CardNumber
		data.Name = toHeadData.Name
		data.ChannelName = toHeadData.ChannelName
	case constants.VaReason_Cust_Cancel_Withdraw:
		fallthrough
	case constants.VaReason_Cust_Withdraw: //银行卡提现
		toCustData, errStr := LogToCustDaoInst.LogToCustDetail(orderNo)
		if errStr != ss_err.ERR_SUCCESS {
			ss_log.Error("LogToCustDetail= [%v]", errStr)
			return nil, ss_err.ERR_SYS_DB_GET
		}
		//data.OrderStatus = toCustData.OrderStatus
		//要审核的订单状态和不需审核的订单状态不同，而现在用户消息这里又审核和不需审核的混淆在一起，现进行转换，统一转换成不需审核的订单状态。
		switch toCustData.OrderStatus {
		case constants.AuditOrderStatus_Pending:
			data.OrderStatus = constants.OrderStatus_Pending
			data.OpType = constants.VaOpType_Minus
		case constants.AuditOrderStatus_Passed:
			data.OpType = constants.VaOpType_Minus
			data.OrderStatus = constants.OrderStatus_Paid
		case constants.AuditOrderStatus_Deny:
			data.OpType = constants.VaOpType_Add
			data.OrderStatus = constants.OrderStatus_Err
		default:

		}
		data.CreateTime = toCustData.CreateTime
		data.Amount = toCustData.Amount
		data.BalanceType = toCustData.BalanceType
		data.FinishTime = toCustData.FinishTime
		data.LogNo = orderNo
		data.Fees = toCustData.Fees
		data.PaymentType = toCustData.PaymentType
		data.OrderType = orderType
		data.CardNumber = toCustData.CardNumber
		data.Name = toCustData.Name
		data.ChannelName = toCustData.ChannelName
	case constants.VaReason_Cust_Pay_Order:
		businessBillData, err := BusinessBillDaoInst.GetBusinessBillDetail(orderNo)
		if err != nil {
			ss_log.Error("orderNo[%v],orderType[%v],err=[%v]", orderNo, orderType, err)
			return nil, ss_err.ERR_SYS_DB_GET
		}
		switch businessBillData.OrderStatus {
		case constants.BusinessOrderStatusPending:
			data.OrderStatus = constants.OrderStatus_Pending
			//data.OpType = constants.VaOpType_Minus
		case constants.BusinessOrderStatusPay:
			data.OrderStatus = constants.OrderStatus_Paid
			data.OpType = constants.VaOpType_Minus
		case constants.BusinessOrderStatusPayTimeOut:
			data.OrderStatus = constants.OrderStatus_Err
			//data.OpType = constants.VaOpType_Add
		case constants.BusinessOrderStatusRefund:
			//对用户来说原来的订单还是交易成功的
			data.OrderStatus = constants.OrderStatus_Paid
			data.OpType = constants.VaOpType_Add
		default:

		}
		data.CreateTime = businessBillData.CreateTime
		data.Amount = businessBillData.Amount
		data.BalanceType = businessBillData.CurrencyType
		data.FinishTime = businessBillData.PayTime
		data.LogNo = orderNo
		//data.Fees = businessBillData.Fee
		data.PaymentType = constants.ORDER_PAYMENT_BALANCE
		data.OrderType = orderType
		//data.ToAccount = businessBillData.ReceiveAccount  //收款人账号
		//data.BusinessName = businessBillData.BusinessName //商家名称
		//data.BusinessAppName = businessBillData.AppName //商家app名称
		data.SimplifyName = businessBillData.SimplifyName //商家简称
		data.Subject = businessBillData.Subject           //商品名称
		data.Notes = businessBillData.Remark              //备注
	case constants.VaReason_ChangeCustBalance: // 改变用户余额
		changeBalanceData, err := ChangeBalanceOrderDaoInst.GetChangeBalanceOrderDetail(orderNo)
		if err != nil {
			ss_log.Error("orderNo[%v],orderType[%v],err=[%v]", orderNo, orderType, err)
			return nil, ss_err.ERR_SYS_DB_GET
		}
		data.OrderStatus = changeBalanceData.OrderStatus
		data.OpType = changeBalanceData.OpType
		data.CreateTime = changeBalanceData.CreateTime
		data.BeforeBalance = changeBalanceData.BeforeBalance
		data.Amount = changeBalanceData.ChangeAmount
		data.AfterBalance = changeBalanceData.AfterBalance
		data.BalanceType = changeBalanceData.CurrencyType
		data.FinishTime = changeBalanceData.CreateTime
		data.LogNo = orderNo
		data.OrderType = orderType
		//产品说原因给平台自己人看，就不返回前端了
		//data.Notes = changeBalanceData.ChangeReason
	case constants.VaReason_BusinessTransferToBusiness: // 商家转账给个人
		whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
			{Key: "bto.log_no", Val: orderNo, EqType: "="},
		})

		businessTransferData, err := BusinessTransferOrderDaoInst.BusinessTransferBillsDetail(whereModel.WhereStr, whereModel.Args)
		if err != nil {
			ss_log.Error("err=%v", err)
			return nil, ss_err.ERR_SYS_DB_GET
		}
		businessTransferData.FromAccount = ss_func.GetDesensitizationEmail(businessTransferData.FromAccount)
		data.Amount = businessTransferData.Amount
		data.BalanceType = businessTransferData.CurrencyType
		data.FromAccount = businessTransferData.FromAccount
		data.ToAccount = businessTransferData.ToAccount
		// 注意这里是给用户看的，不应该知道付款方付了多少手续费，对用户来说就是收到一笔转账，手续费为0
		//data.Fees = businessTransferData.Fee
		data.Fees = "0"
		data.CreateTime = businessTransferData.CreateTime

		//转换成前端统一认识的订单状态
		switch businessTransferData.OrderStatus {
		case constants.BusinessTransferOrderStatusPending:
			data.OrderStatus = constants.OrderStatus_Pending
		case constants.BusinessTransferOrderStatusSuccess:
			data.OrderStatus = constants.OrderStatus_Paid
		case constants.BusinessTransferOrderStatusFail:
			data.OrderStatus = constants.OrderStatus_Err
		default:
			ss_log.Error("未知订单状态，OrderStatus[%v]", businessTransferData.OrderStatus)
			data.OrderStatus = businessTransferData.OrderStatus
		}

		data.LogNo = orderNo
		data.OrderType = orderType
	case constants.VaReason_BusinessRefund: // 商家退款
		if accountNo == "" {
			return nil, ss_err.ERR_PARAM
		}
		whereList := []*model.WhereSqlCond{
			{Key: "va.account_no", Val: accountNo, EqType: "="},
			{Key: "bro.refund_no", Val: orderNo, EqType: "="},
		}

		refundData, err := BusinessRefundOrderDaoInst.GetRefundOrderDetail(whereList)
		if err != nil {
			ss_log.Error("orderNo[%v],orderType[%v],err=[%v]", orderNo, orderType, err)
			return nil, ss_err.ERR_SYS_DB_GET
		}

		switch refundData.RefundStatus {
		case constants.BusinessRefundStatusPending:
			data.OrderStatus = constants.OrderStatus_Pending
			//data.OpType = constants.VaOpType_Minus
		case constants.BusinessRefundStatusSuccess:
			data.OrderStatus = constants.OrderStatus_Paid
			data.OpType = constants.VaOpType_Add
		case constants.BusinessRefundStatusFail:
			data.OrderStatus = constants.OrderStatus_Err
			//data.OpType = constants.VaOpType_Add
		default:
		}
		data.Amount = refundData.PayeeAmount //用户收款金额
		data.TransAmount = refundData.Amount //商家实际退款金额
		data.BalanceType = refundData.CurrencyType
		data.FinishTime = refundData.FinishTime

		data.LogNo = orderNo //这是退款订单号,不是付款订单号
		data.OrderType = orderType
		data.Subject = refundData.Subject
		data.PayOrderNo = refundData.PayOrderNo //付款订单号
	default:
		ss_log.Error("参数异常 LogNo：[%v]  orderType:[%v] ", orderNo, orderType)
		return nil, ss_err.ERR_PARAM
	}

	return data, ss_err.ERR_SUCCESS
}

func (*LogAppMessagesDao) AddLogAppMessages(orderNo, appMessType, orderType, accountNo, orderStatus string) (errR string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "insert into log_app_messages(log_no, order_no, order_type, is_read, is_push, account_no, app_mess_type, order_status, create_time) " +
		"values($1,$2,$3,$4,$5,$6,$7,$8,current_timestamp)"
	err := ss_sql.Exec(dbHandler, sqlStr, strext.GetDailyId(), orderNo, orderType, "0", "0", accountNo, appMessType, orderStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_ADD
	}

	return ss_err.ERR_SUCCESS

}
func (*LogAppMessagesDao) AddLogAppMessagesTx(tx *sql.Tx, orderNo, appMessType, orderType, accountNo, orderStatus string) (errR string) {
	sqlStr := "insert into log_app_messages(log_no, order_no, order_type, is_read, is_push, account_no, app_mess_type, order_status, create_time) " +
		"values($1,$2,$3,$4,$5,$6,$7,$8,current_timestamp)"
	err := ss_sql.ExecTx(tx, sqlStr, strext.GetDailyId(), orderNo, orderType, "0", "0", accountNo, appMessType, orderStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_ADD
	}

	return ss_err.ERR_SUCCESS

}
