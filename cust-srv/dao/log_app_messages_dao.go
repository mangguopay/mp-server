package dao

import (
	"a.a/mp-server/common/model"
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type LogAppMessagesDao struct {
}

var LogAppMessagesDaoInst LogAppMessagesDao

//批量修改推送的消息状态为已读，logNos为订单流水号的集合
func (*LogAppMessagesDao) ModiftLogAppMessagesIsRead(dbHandler *sql.DB, logNos []string) (errR string) {
	for _, logNo := range logNos {
		sqlStr := "update log_app_messages set is_read = '1' where log_no = $1 and is_read = '0' "
		err := ss_sql.Exec(dbHandler, sqlStr, logNo)
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
	}

	return ss_err.ERR_SUCCESS
}

//设置账号的所有消息为已读
func (*LogAppMessagesDao) ModiftAllRead(accNo string) (errR error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update log_app_messages set is_read = '1' where account_no = $1 and is_read = '0' "
	err := ss_sql.Exec(dbHandler, sqlStr, accNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

func (*LogAppMessagesDao) GetLogAppMessages(dbHandler *sql.DB, whereStr string, whereArgs []interface{}) (datas []*go_micro_srv_cust.LogAppMessagesData, returnErr string) {
	//111

	sqlStr := "SELECT log_no, order_no, order_type, is_read, is_push, account_no, app_mess_type, order_status  " +
		" FROM log_app_messages " + whereStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err == nil {
		for rows.Next() {
			data := &go_micro_srv_cust.LogAppMessagesData{}
			var orderNo, orderType, isRead, isPush, appMessType sql.NullString
			err = rows.Scan(
				&data.LogNo, //消息id
				&orderNo,    //订单id
				&orderType,
				&isRead,
				&isPush,
				&data.AccountNo,
				&appMessType,
				&data.OrderStatus,
			)

			if err != nil {
				ss_log.Error("err=%v", err)
				continue
			}

			data.AppMessType = appMessType.String
			data.IsRead = isRead.String
			data.OrderNo = orderNo.String

			switch orderType.String {
			case constants.VaReason_Exchange: //兑换
				exchangeData, err := ExchangeOrderDaoInst.CustExchangeBillsDetail(orderNo.String)
				if err != ss_err.ERR_SUCCESS {
					ss_log.Error("获取兑换订单数据失败orderNo[%v],Err= [%v]", orderNo.String, err)
					continue
				}
				data.InType = exchangeData.InType
				data.OutType = exchangeData.OutType
				data.Amount = exchangeData.Amount
				data.CreateTime = exchangeData.CreateTime
				//data.Rate = exchangeData.Rate
				//data.OrderStatus = exchangeData.OrderStatus
				data.FinishTime = exchangeData.FinishTime
				//data.AccountNo = exchangeData.AccountNo
				data.TransFrom = exchangeData.TransFrom     //app,trade
				data.TransAmount = exchangeData.TransAmount //转换后金额
				//data.ErrReason = exchangeData.ErrReason
				data.Fees = exchangeData.Fees
				data.OrderType = orderType.String
				data.BalanceType = exchangeData.InType
				data.PaymentType = constants.ORDER_PAYMENT_BALANCE

				datas = append(datas, data)
			case constants.VaReason_INCOME: //存款
				incomeData, err := IncomeOrderDaoInst.CustIncomeBillsDetail(orderNo.String)

				if err != ss_err.ERR_SUCCESS {
					ss_log.Error("获取存款订单数据失败orderNo[%v],Err= [%v]", orderNo.String, err)
					continue
				}
				//data.OrderStatus = incomeData.OrderStatus
				data.CreateTime = incomeData.CreateTime
				data.Amount = incomeData.Amount
				data.BalanceType = incomeData.BalanceType
				data.PaymentType = incomeData.PaymentType
				data.OrderType = orderType.String
				data.FinishTime = incomeData.FinishTime
				data.IsRead = isRead.String

				switch incomeData.OrderStatus {
				case constants.OrderStatus_Paid:
					data.OpType = constants.VaOpType_Add
				case constants.OrderStatus_Err:
					data.OpType = constants.VaOpType_Minus
				default:
					ss_log.Error("orderStatus[%v],未处理,opType可能错误", incomeData.OrderStatus)
				}

				datas = append(datas, data)
			case constants.VaReason_OUTGO: //取款
				outgoData, err := OutgoOrderDaoInst.CustOutgoBillsDetail(orderNo.String)
				if err != ss_err.ERR_SUCCESS {
					ss_log.Error("获取取款订单数据失败orderNo[%v],Err= [%v]", orderNo.String, err)
					continue
				}
				//data.OrderStatus = outgoData.OrderStatus
				data.CreateTime = outgoData.CreateTime
				data.Amount = outgoData.Amount
				data.BalanceType = outgoData.BalanceType
				data.PaymentType = outgoData.PaymentType
				data.OrderType = orderType.String
				data.FinishTime = outgoData.FinishTime
				data.IsRead = isRead.String
				data.OpType = constants.VaOpType_Minus

				switch outgoData.OrderStatus {
				case constants.OrderStatus_Paid:
					data.OpType = constants.VaOpType_Minus
				case constants.OrderStatus_Err:
					data.OpType = constants.VaOpType_Add
				case constants.OrderStatus_Cancel: //虽然说取消提现的不显示订单，但这里还是给出了OpType的转换
					data.OpType = constants.VaOpType_Add
				default:
					ss_log.Error("orderStatus[%v]未处理，opType可能错误。", outgoData.OrderStatus)
				}

				datas = append(datas, data)
			case constants.VaReason_TRANSFER: //转账
				transferData, err := TransferDaoInst.CustTransferBillsDetailByAccount(data.AccountNo, orderNo.String)
				if err != ss_err.ERR_SUCCESS {
					ss_log.Error("获取转账订单数据失败orderNo[%v],Err= [%v]", orderNo.String, err)
					continue
				}
				//data.OrderStatus = transferData.OrderStatus
				data.CreateTime = transferData.CreateTime
				data.Amount = transferData.Amount
				data.BalanceType = transferData.BalanceType
				data.PaymentType = transferData.PaymentType
				data.OrderType = orderType.String
				data.FinishTime = transferData.FinishTime
				data.IsRead = isRead.String
				data.OpType = transferData.OpType
				datas = append(datas, data)
			case constants.VaReason_COLLECTION: //收款
				collectionData, err := CollectionDaoInst.CustCollectionBillsDetail(orderNo.String)
				if err != ss_err.ERR_SUCCESS {
					ss_log.Error("获取收款订单数据失败orderNo[%v],Err= [%v]", orderNo.String, err)
					continue
				}

				//data.OrderStatus = collectionData.OrderStatus
				data.CreateTime = collectionData.CreateTime
				data.Amount = collectionData.Amount
				data.BalanceType = collectionData.BalanceType
				data.PaymentType = collectionData.PaymentType
				data.OrderType = orderType.String
				data.FinishTime = collectionData.FinishTime
				data.IsRead = isRead.String
				data.FromAccount = collectionData.FromAccount
				data.OpType = constants.VaOpType_Add
				datas = append(datas, data)
			case constants.VaReason_Cust_Cancel_Save: //驳回客户向总部存款
				fallthrough
			case constants.VaReason_Cust_Save: //用户线上充值
				toHeadData, errStr := LogCustToHeadquartersDaoInst.CustToHeadquartersDetail(orderNo.String)
				if errStr != ss_err.ERR_SUCCESS {
					ss_log.Error("获取用户线上充值订单数据失败orderNo[%v],Err= [%v]", orderNo.String, err)
					return nil, ss_err.ERR_SYS_DB_GET
				}

				//要审核的订单状态和不需审核的订单状态不同，而现在用户消息这里又审核和不需审核的混淆在一起，现进行转换，统一转换成不需审核的订单状态。
				switch toHeadData.OrderStatus {
				case constants.AuditOrderStatus_Passed:
					data.OrderStatus = constants.OrderStatus_Paid
				case constants.AuditOrderStatus_Deny:
					data.OrderStatus = constants.OrderStatus_Err
				default:

				}

				if orderType.String == constants.VaReason_Cust_Save {
					data.OpType = constants.VaOpType_Add
				}

				data.CreateTime = toHeadData.CreateTime
				data.Amount = toHeadData.Amount
				data.BalanceType = toHeadData.BalanceType
				data.FinishTime = toHeadData.FinishTime
				//data.OpType = toHeadData.OpType
				data.Fees = toHeadData.Fees
				data.PaymentType = toHeadData.PaymentType
				data.OrderType = constants.VaReason_Cust_Save
				datas = append(datas, data)
			case constants.VaReason_Cust_Cancel_Withdraw: //驳回客户向总部提现
				data.OpType = constants.VaOpType_Add
				fallthrough
			case constants.VaReason_Cust_Withdraw: //用户线上提现
				toCustData, errStr := LogToCustDaoInst.LogToCustDetail(orderNo.String)
				if errStr != ss_err.ERR_SUCCESS {
					ss_log.Error("获取用户线上提现订单数据失败orderNo[%v],Err= [%v]", orderNo.String, err)
					return nil, ss_err.ERR_SYS_DB_GET
				}

				switch toCustData.OrderStatus { //要审核的订单状态和不需审核的订单状态不同，现暂时进行转换
				case constants.AuditOrderStatus_Passed:
					data.OrderStatus = constants.OrderStatus_Paid
				case constants.AuditOrderStatus_Deny:
					data.OrderStatus = constants.OrderStatus_Err
				default:
				}

				if data.OpType == "" {
					data.OpType = constants.VaOpType_Minus
				}

				data.CreateTime = toCustData.CreateTime
				data.Amount = toCustData.Amount
				data.BalanceType = toCustData.BalanceType
				data.FinishTime = toCustData.FinishTime
				//data.OpType = toCustData.OpType
				data.Fees = toCustData.Fees
				data.PaymentType = toCustData.PaymentType
				data.OrderType = constants.VaReason_Cust_Withdraw
				datas = append(datas, data)
			case constants.VaReason_Cust_Pay_Order: //用户付款
				businessBillData, err := BusinessBillDaoInst.GetBusinessBillDetail(orderNo.String)
				if err != nil {
					ss_log.Error("orderNo[%v],orderType[%v],err=[%v]", orderNo, orderType, err)
					return nil, ss_err.ERR_SYS_DB_GET
				}
				switch businessBillData.OrderStatus {
				case constants.BusinessOrderStatusPending:
					data.OrderStatus = constants.OrderStatus_Pending
				case constants.BusinessOrderStatusPay:
					data.OrderStatus = constants.OrderStatus_Paid
					data.OpType = constants.VaOpType_Minus
				case constants.BusinessOrderStatusPayTimeOut:
					data.OrderStatus = constants.OrderStatus_Err
				case constants.BusinessOrderStatusRefund:
					data.OrderStatus = constants.OrderStatus_Cancel
					data.OpType = constants.VaOpType_Add
				default:

				}
				data.CreateTime = businessBillData.CreateTime
				data.Amount = businessBillData.Amount
				data.BalanceType = businessBillData.CurrencyType
				data.FinishTime = businessBillData.PayTime
				data.OrderType = constants.VaReason_Cust_Pay_Order
				datas = append(datas, data)
			case constants.VaReason_BusinessTransferToBusiness: //商家转账
				businessToUserData, err := BusinessTransferDaoInst.GetOrderDetail(orderNo.String)
				if err != nil {
					ss_log.Error("获商家转账至用户订单数据失败orderNo[%v],Err= [%v]", orderNo.String, err)
					return nil, ss_err.ERR_SYS_DB_GET
				}

				switch businessToUserData.OrderStatus {
				case constants.BusinessTransferOrderStatusPending: //订单状态：0处理中，1成功，2失败
					data.OrderStatus = constants.OrderStatus_Pending
				case constants.BusinessTransferOrderStatusSuccess:
					data.OrderStatus = constants.OrderStatus_Paid
				case constants.BusinessTransferOrderStatusFail:
					data.OrderStatus = constants.OrderStatus_Err
				default:
				}

				if data.OpType == "" {
					data.OpType = constants.VaOpType_Add
				}

				data.CreateTime = businessToUserData.CreateTime
				data.Amount = businessToUserData.Amount
				data.BalanceType = businessToUserData.CurrencyType
				data.PaymentType = businessToUserData.PaymentType
				data.OrderType = orderType.String
				data.FinishTime = businessToUserData.CreateTime
				data.IsRead = isRead.String
				datas = append(datas, data)
			case constants.VaReason_BusinessRefund: // 商家退款
				whereList := []*model.WhereSqlCond{
					{Key: "vacc.account_no", Val: data.AccountNo, EqType: "="},
					{Key: "bro.refund_no", Val: orderNo.String, EqType: "="},
				}

				refundData, err := BusinessRefundOrderDaoInst.GetRefundOrderDetail(whereList)
				if err != nil {
					ss_log.Error("orderNo[%v],orderType[%v],err=[%v]", orderNo.String, orderType, err)
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
				data.CreateTime = refundData.FinishTime
				data.FinishTime = refundData.FinishTime

				data.OrderNo = orderNo.String //这是退款订单号,不是付款订单号
				data.OrderType = orderType.String
				datas = append(datas, data)
			default:
				ss_log.Error("LogNo：[%v]  orderType[%v] 类型没有处理方法，可能部分数据不全。。 ", data.LogNo, orderType.String)
			}
		}
	} else {
		ss_log.Error("err=[%v]", err)
		return nil, ss_err.ERR_SYS_DB_GET
	}
	return datas, ss_err.ERR_SUCCESS
}

//获取消息数量
func (*LogAppMessagesDao) GetLogAppMessagesCnt(dbHandler *sql.DB, whereStr string, whereArgs []interface{}) (total, err string) {
	sqlStr := "SELECT count(1)  " +
		" FROM log_app_messages " + whereStr
	var totalT sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, whereArgs...)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return "", ss_err.ERR_PARAM
	}

	return totalT.String, ss_err.ERR_SUCCESS
}

func (*LogAppMessagesDao) AddLogAppMessages(tx *sql.Tx, orderNo, appMessType, orderType, accountNo, orderStatus string) (errR string) {
	sqlStr := "insert into log_app_messages(log_no, order_no, order_type, is_read, is_push, account_no, app_mess_type, order_status, create_time) " +
		"values($1,$2,$3,$4,$5,$6,$7,$8,current_timestamp)"
	err := ss_sql.ExecTx(tx, sqlStr, strext.GetDailyId(), orderNo, orderType, "0", "0", accountNo, appMessType, orderStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_ADD
	}

	return ss_err.ERR_SUCCESS
}
