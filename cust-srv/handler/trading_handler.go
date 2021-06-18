package handler

import (
	"context"
	"database/sql"
	"fmt"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	go_micro_srv_quota "a.a/mp-server/common/proto/quota"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/dao"
	"a.a/mp-server/cust-srv/i"
)

/**
 * 转账至总部
 */
func (*CustHandler) GetToHeadquartersList(ctx context.Context, req *go_micro_srv_cust.GetToHeadquartersListRequest, reply *go_micro_srv_cust.GetToHeadquartersListReply) error {
	if req.StartTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.StartTime) {
			ss_log.Error("GetToHeadquartersList StartTime格式不正确,StartTime: %s", req.StartTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}
	if req.EndTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.EndTime) {
			ss_log.Error("GetToHeadquartersList EndTime格式不正确,StartTime: %s", req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "lth.log_no", Val: req.LogNo, EqType: "like"},
		{Key: "lth.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "lth.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "lth.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "lth.order_type", Val: req.OrderType, EqType: "="},
		{Key: "lth.currency_type", Val: req.CurrencyType, EqType: "="},
		{Key: "acc.account", Val: req.Account, EqType: "like"}, //服务商账号
	})
	//统计
	var total sql.NullString
	sqlCnt := "SELECT count(1) FROM log_to_headquarters lth " +
		" LEFT JOIN servicer ser ON ser.servicer_no = lth.servicer_no " +
		" LEFT JOIN account acc ON acc.uid = ser.account_no " +
		" LEFT JOIN card c ON c.card_no= lth.card_no " + whereModel.WhereStr
	if err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...); err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	if total.String == "" || total.String == "0" {
		reply.ResultCode = ss_err.ERR_SUCCESS
		reply.Total = strext.ToInt32(total.String)
		return nil
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by case lth.order_status when "+constants.AuditOrderStatus_Pending+" then 1 end, lth.create_time desc")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	sqlStr := "SELECT lth.log_no, lth.servicer_no, lth.currency_type, lth.amount, lth.order_status, lth.collection_type, lth.card_no, lth.create_time, lth.finish_time, lth.order_type, lth.image_id " +
		", acc.account,c.name,c.card_number, ch.channel_name " +
		" FROM log_to_headquarters lth " +
		" LEFT JOIN servicer ser ON ser.servicer_no = lth.servicer_no " +
		" LEFT JOIN account acc ON acc.uid = ser.account_no " +
		" LEFT JOIN card_head c ON c.card_no= lth.card_no " +
		" LEFT JOIN channel ch ON ch.channel_no = c.channel_no " + whereModel.WhereStr

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
	var datas []*go_micro_srv_cust.ToHeadquartersData
	for rows.Next() {
		data := &go_micro_srv_cust.ToHeadquartersData{}
		var name, cardNumber, channelName sql.NullString
		if err = rows.Scan(
			&data.LogNo,
			&data.ServicerNo,
			&data.CurrencyType,
			&data.Amount,
			&data.OrderStatus,

			&data.CollectionType,
			&data.CardNo,
			&data.CreateTime,
			&data.FinishTime,
			&data.OrderType,

			&data.ImageId,
			&data.Account,
			&name,
			&cardNumber,
			&channelName,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		data.Name = name.String
		data.CardNumber = cardNumber.String
		data.ChannelName = channelName.String

		datas = append(datas, data)
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * 修改转账至总部订单的状态
 */
func (*CustHandler) UpdateToHeadquarters(ctx context.Context, req *go_micro_srv_cust.UpdateToHeadquartersRequest, reply *go_micro_srv_cust.UpdateToHeadquartersReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	if req.LogNo == "" {
		ss_log.Error("logNo is nil")
		reply.ResultCode = ss_err.ERR_PAY_NO_THIS_ORDER
		return nil
	}

	if req.LoginUid == "" {
		ss_log.Error("LoginUid参数为空")
		reply.ResultCode = ss_err.ERR_ACCOUNT_JWT_OUTDATED
		return nil
	}

	// 检查状态
	if req.OrderStatus != constants.AuditOrderStatus_Passed && req.OrderStatus != constants.AuditOrderStatus_Deny {
		ss_log.Error("OrderStatus no in (1,2)")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 获取旧状态
	err, orderStatus := dao.LogToHeadquartersDaoInst.GetLogToHeadquarterOrderStatus(tx, req.LogNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	// 检查状态
	if orderStatus != constants.AuditOrderStatus_Pending {
		ss_log.Error("logNo is nil")
		reply.ResultCode = ss_err.ERR_PAY_ORDER_STATUS_MISTAKE
		return nil
	}
	// 审核
	err = dao.LogToHeadquartersDaoInst.UpdateLogToHeadquarterOrderStatus(tx, req.LogNo, req.OrderStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	// 根据logNo查询
	servicerNo, amount, currentType := dao.LogToHeadquartersDaoInst.GetLogToHeadquarterFromNo(tx, req.LogNo)
	accNo := dao.RelaAccIdenDaoInst.GetAccNo(servicerNo, constants.AccountType_SERVICER)
	if amount == "" || servicerNo == "" {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	if accNo == "" {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	quotaReq := &go_micro_srv_quota.ModifyQuotaRequest{
		CurrencyType: currentType,
		//Amount:       amount,
		AccountNo: accNo,
		//OpType:       constants.QuotaOp_SvrSave,
		LogNo: req.LogNo,
	}

	description := ""
	switch req.OrderStatus { //审核通过.
	case constants.AuditOrderStatus_Passed:
		quotaReq.OpType = constants.QuotaOp_SvrSave
		quotaReq.Amount = amount
		// 调用ps服务商预存
		//quotaRepl := &go_micro_srv_quota.ModifyQuotaReply{}
		quotaRepl, err2 := i.QuotaHandleInstance.Client.ModifyQuota(ctx, quotaReq)
		if err2 != nil || quotaRepl.ResultCode != ss_err.ERR_SUCCESS {
			ss_log.Error("err=[--------------->%s]", "服务商存款,调用八神的服务失败,操作为服务商存款")
			reply.ResultCode = quotaRepl.ResultCode
			return nil
		}
		// 根据 serviceNo 查询accNo
		serviceAccNo := dao.ServiceDaoInst.GetAccNoFromSrvNo(servicerNo)
		// 插入服务商交易明细
		if errStr := dao.BillingDetailsResultsDaoInst.InsertResult(amount, currentType, serviceAccNo, constants.AccountType_SERVICER, req.LogNo, "0", constants.OrderStatus_Paid, constants.BillDetailTypeRecharge, "0", amount); errStr == "" {
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		description = fmt.Sprintf("审核服务商充值订单[%v],操作[%v]", req.LogNo, "通过")
	case constants.AuditOrderStatus_Deny: //驳回,现发起充值时，冻结金额不改变，所以不需改变

		description = fmt.Sprintf("审核服务商充值订单[%v],操作[%v]", req.LogNo, "驳回")
	}

	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Trading_Order)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
 * 获取转账至服务商流水
 */
func (*CustHandler) GetToServicerList(ctx context.Context, req *go_micro_srv_cust.GetToServicerListRequest, reply *go_micro_srv_cust.GetToServicerListReply) error {
	if req.StartTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.StartTime) {
			ss_log.Error("GetToServicerList StartTime格式不正确,StartTime: %s", req.StartTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}
	if req.EndTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.EndTime) {
			ss_log.Error("GetToServicerList EndTime格式不正确,StartTime: %s", req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "lts.log_no", Val: req.LogNo, EqType: "="},
		{Key: "lts.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "lts.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "lts.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "lts.order_type", Val: req.OrderType, EqType: "="},
		{Key: "lts.currency_type", Val: req.CurrencyType, EqType: "="},
		{Key: "acc.account", Val: req.Account, EqType: "like"},
		{Key: "acc.nickname", Val: req.Nickname, EqType: "like"},
	})
	//统计
	var total sql.NullString
	sqlCnt := "SELECT count(1) FROM log_to_servicer lts " +
		" LEFT JOIN servicer ser ON ser.servicer_no = lts.servicer_no " +
		" LEFT JOIN account acc ON acc.uid = ser.account_no " +
		" LEFT JOIN card ca ON ca.card_no= lts.card_no " + whereModel.WhereStr
	if err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...); err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	if total.String == "" || total.String == "0" {
		reply.ResultCode = ss_err.ERR_SUCCESS
		reply.Total = strext.ToInt32(total.String)
		return nil
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by case lts.order_status when "+constants.AuditOrderStatus_Pending+" then 1 end, lts.create_time desc")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	sqlStr := "SELECT lts.log_no, lts.servicer_no, lts.currency_type, lts.amount, lts.order_status, lts.collection_type, lts.card_no, lts.create_time, lts.finish_time, lts.order_type " +
		", acc.account, acc.phone, acc.nickname, ca.name, ca.card_number, ch.channel_name " +
		" FROM log_to_servicer lts " +
		" LEFT JOIN servicer ser ON ser.servicer_no = lts.servicer_no " +
		" LEFT JOIN account acc ON acc.uid = ser.account_no " +
		" LEFT JOIN card ca ON ca.card_no= lts.card_no " +
		" LEFT JOIN channel ch ON ch.channel_no = ca.channel_no " + whereModel.WhereStr
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
	var datas []*go_micro_srv_cust.ToServicerData
	for rows.Next() {
		data := &go_micro_srv_cust.ToServicerData{}
		var channelName sql.NullString
		if err = rows.Scan(
			&data.LogNo,
			&data.ServicerNo,
			&data.CurrencyType,
			&data.Amount,
			&data.OrderStatus,

			&data.CollectionType,
			&data.CardNo,
			&data.CreateTime,
			&data.FinishTime,
			&data.OrderType,

			&data.Account,
			&data.Phone,
			&data.Nickname,
			&data.Name,
			&data.CardNumber,
			&channelName,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		data.ChannelName = channelName.String
		datas = append(datas, data)
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * 总部打钱给服务商
 */
func (*CustHandler) AddToServicer(ctx context.Context, req *go_micro_srv_cust.AddToServicerRequest, reply *go_micro_srv_cust.AddToServicerReply) error {
	if errStr := dao.LogToServiceDaoInstance.InsertLogToService1(req.CurrencyType, req.ServicerNo, req.CollectionType,
		req.CardNo, req.Amount, req.OrderType, req.OrderStatus); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
 * 获取转账流水
 */
func (*CustHandler) GetTransferOrderList(ctx context.Context, req *go_micro_srv_cust.GetTransferOrderListRequest, reply *go_micro_srv_cust.GetTransferOrderListReply) error {

	if req.StartTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.StartTime) {
			ss_log.Error("GetTransferOrderList StartTime格式不正确,StartTime: %s", req.StartTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}
	if req.EndTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.EndTime) {
			ss_log.Error("GetTransferOrderList EndTime格式不正确,StartTime: %s", req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var datas []*go_micro_srv_cust.TransferOrderData

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "tr.log_no", Val: req.LogNo, EqType: "like"},
		{Key: "tr.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "tr.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "tr.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "w.code", Val: req.WriteOff, EqType: "like"},
		{Key: "tr.balance_type", Val: req.BalanceType, EqType: "="},
		{Key: "acc.account", Val: req.FromAccount, EqType: "like"},
		{Key: "acc2.account", Val: req.ToAccount, EqType: "like"},
	})
	//统计
	where := whereModel.WhereStr
	args := whereModel.Args
	var total sql.NullString
	sqlCnt := "SELECT count(1) FROM transfer_order tr " +
		" LEFT JOIN writeoff w ON w.transfer_order_no = tr.log_no " +
		" LEFT JOIN vaccount vacc ON vacc.vaccount_no = tr.from_vaccount_no " +
		" LEFT JOIN account acc ON acc.uid = vacc.account_no " +
		" LEFT JOIN vaccount vacc2 ON vacc2.vaccount_no = tr.to_vaccount_no " +
		" LEFT JOIN account acc2 ON acc2.uid = vacc2.account_no " + where
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by tr.create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT tr.log_no, tr.from_vaccount_no, tr.to_vaccount_no, tr.amount, tr.create_time, tr.finish_time, tr.order_status, tr.balance_type, tr.exchange_type, tr.fees, w.code " +
		", acc.account, acc2.account, tr.ree_rate, tr.real_amount " +
		" FROM transfer_order tr " +
		" LEFT JOIN writeoff w ON w.transfer_order_no = tr.log_no " +
		" LEFT JOIN vaccount vacc ON vacc.vaccount_no = tr.from_vaccount_no " +
		" LEFT JOIN account acc ON acc.uid = vacc.account_no " +
		" LEFT JOIN vaccount vacc2 ON vacc2.vaccount_no = tr.to_vaccount_no " +
		" LEFT JOIN account acc2 ON acc2.uid = vacc2.account_no " + where2
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("err[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	for rows.Next() {
		data := go_micro_srv_cust.TransferOrderData{}
		var writeoff sql.NullString
		var fromVaccountNo, toVaccountNo, feeRate, realAmount, fromAccount, toAccount sql.NullString
		if err = rows.Scan(
			//  w.code
			&data.LogNo,
			&fromVaccountNo,
			&toVaccountNo,
			&data.Amount,
			&data.CreateTime,

			&data.FinishTime,
			&data.OrderStatus,
			&data.BalanceType,
			&data.ExchangeType,
			&data.Fees,

			&writeoff,
			&fromAccount,
			&toAccount,
			&feeRate,
			&realAmount,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		data.FromAccount = fromAccount.String
		data.ToAccount = toAccount.String

		if writeoff.String != "" {
			data.WriteOff = writeoff.String
		}
		if feeRate.String != "" {
			data.FeeRate = feeRate.String
		}
		if realAmount.String != "" {
			data.RealAmount = realAmount.String
		}
		datas = append(datas, &data)
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * 获取收款流水
 */

func (*CustHandler) GetCollectionOrders(ctx context.Context, req *go_micro_srv_cust.GetCollectionOrdersRequest, reply *go_micro_srv_cust.GetCollectionOrdersReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var datas []*go_micro_srv_cust.CollectionOrderData

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "coo.log_no", Val: req.LogNo, EqType: "like"},
		{Key: "coo.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "coo.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "coo.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "coo.balance_type", Val: req.BalanceType, EqType: "="},
		{Key: "acc.account", Val: req.FromAccount, EqType: "like"},
		{Key: "acc2.account", Val: req.ToAccount, EqType: "like"},
	})
	//统计
	where := whereModel.WhereStr
	args := whereModel.Args
	var total sql.NullString
	sqlCnt := "SELECT count(1) FROM collection_order coo " +
		" LEFT JOIN vaccount vacc ON vacc.vaccount_no = coo.from_vaccount_no " +
		" LEFT JOIN account acc ON acc.uid = vacc.account_no " +
		" LEFT JOIN vaccount vacc2 ON vacc2.vaccount_no = coo.to_vaccount_no " +
		" LEFT JOIN account acc2 ON acc2.uid = vacc2.account_no " + where
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by coo.create_time desc`)
	if req.NoPaging == "" {
		ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
	}

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT coo.log_no, coo.amount, coo.create_time, coo.finish_time, coo.order_status,coo.payment_type, coo.balance_type, coo.fees " +
		", acc.account, acc2.account " +
		" FROM collection_order coo " +
		" LEFT JOIN vaccount vacc ON vacc.vaccount_no = coo.from_vaccount_no " +
		" LEFT JOIN account acc ON acc.uid = vacc.account_no " +
		" LEFT JOIN vaccount vacc2 ON vacc2.vaccount_no = coo.to_vaccount_no " +
		" LEFT JOIN account acc2 ON acc2.uid = vacc2.account_no " + where2
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
	for rows.Next() {
		data := go_micro_srv_cust.CollectionOrderData{}
		var paymentType, fees sql.NullString
		if err = rows.Scan(
			&data.LogNo,
			&data.Amount,
			&data.CreateTime,
			&data.FinishTime,
			&data.OrderStatus,

			&paymentType,
			&data.BalanceType,
			&fees,
			//&data.Fees,
			&data.FromAccount,
			&data.ToAccount,
		); err != nil {
			ss_log.Error("logNo=[%v],err=[%v]", data.LogNo, err)
			continue
		}

		if paymentType.String != "" {
			data.PaymentType = paymentType.String
		}
		if fees.String != "" {
			data.Fees = fees.String
		}

		datas = append(datas, &data)
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * 获取取款流水
 */
func (*CustHandler) GetOutgoOrderList(ctx context.Context, req *go_micro_srv_cust.GetOutgoOrderListRequest, reply *go_micro_srv_cust.GetOutgoOrderListReply) error {
	if req.StartTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.StartTime) {
			ss_log.Error("GetOutgoOrderList StartTime格式不正确,StartTime: %s", req.StartTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}
	if req.EndTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.EndTime) {
			ss_log.Error("GetOutgoOrderList EndTime格式不正确,StartTime: %s", req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var datas []*go_micro_srv_cust.OutgoOrderData

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "oo.log_no", Val: req.LogNo, EqType: "="},
		{Key: "oo.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "oo.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "oo.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "oo.balance_type", Val: req.BalanceType, EqType: "="},
		{Key: "oo.withdraw_type", Val: req.WithdrawType, EqType: "="}, //提现类型
		{Key: "w.code", Val: req.WriteOff, EqType: "="},
		{Key: "acc.phone", Val: req.Phone, EqType: "like"},        //取款的手机号
		{Key: "acc.account", Val: req.OutAccount, EqType: "like"}, //取款的账号
		{Key: "acc2.account", Val: req.Account, EqType: "like"},   //服务商账号
	})
	//统计
	where := whereModel.WhereStr
	args := whereModel.Args
	var total sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM outgo_order oo " +
		" LEFT JOIN vaccount vacc ON oo.vaccount_no = vacc.vaccount_no " +
		" LEFT JOIN account acc ON acc.uid = vacc.account_no " +
		" LEFT JOIN servicer ser ON oo.servicer_no = ser.servicer_no " +
		" LEFT JOIN account acc2 ON acc2.uid = ser.account_no " +
		" LEFT JOIN writeoff w ON w.outgo_order_no = oo.log_no " + where
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by oo.create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT oo.log_no, oo.vaccount_no, oo.amount, oo.create_time, oo.order_status, oo.modify_time, oo.balance_type, oo.fees, w.code, oo.servicer_no" +
		", acc.phone, acc.account, acc2.account, oo.rate, oo.real_amount, oo.cancel_reason, oo.withdraw_type " +
		" FROM outgo_order oo " +
		" LEFT JOIN vaccount vacc ON oo.vaccount_no = vacc.vaccount_no " +
		" LEFT JOIN account acc ON acc.uid = vacc.account_no " +
		" LEFT JOIN servicer ser ON oo.servicer_no = ser.servicer_no " +
		" LEFT JOIN account acc2 ON acc2.uid = ser.account_no " +
		" LEFT JOIN writeoff w ON w.outgo_order_no = oo.log_no " + where2
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
	for rows.Next() {
		data := go_micro_srv_cust.OutgoOrderData{}
		var writeoff, finishTime, account, rate, realAmount, phone, outAccount, cancelReason sql.NullString
		if err = rows.Scan(
			&data.LogNo,
			&data.VaccountNo,
			&data.Amount,
			&data.CreateTime,
			&data.OrderStatus,

			&finishTime,
			&data.BalanceType,
			&data.Fees,
			&writeoff,
			&data.ServicerNo,

			&phone,
			&outAccount,
			&account,
			&rate,
			&realAmount,
			&cancelReason,
			&data.WithdrawType,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		data.WriteOff = writeoff.String
		data.FinishTime = finishTime.String
		data.Account = account.String
		data.Phone = phone.String
		data.Rate = rate.String
		data.RealAmount = realAmount.String
		data.CancelReason = cancelReason.String
		data.OutAccount = outAccount.String
		datas = append(datas, &data)
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * 获取存款流水
 */
func (*CustHandler) GetIncomeOrderList(ctx context.Context, req *go_micro_srv_cust.GetIncomeOrderListRequest, reply *go_micro_srv_cust.GetIncomeOrderListReply) error {
	if req.StartTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.StartTime) {
			ss_log.Error("GetIncomeOrderList StartTime格式不正确,StartTime: %s", req.StartTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}
	if req.EndTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.EndTime) {
			ss_log.Error("GetIncomeOrderList EndTime格式不正确,StartTime: %s", req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var datas []*go_micro_srv_cust.IncomeOrderData

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "inc.log_no", Val: req.LogNo, EqType: "like"},
		{Key: "inc.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "inc.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "inc.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "inc.balance_type", Val: req.BalanceType, EqType: "="},
		{Key: "w.code", Val: req.WriteOff, EqType: "="},
		{Key: "acc.phone", Val: req.IncomePhone, EqType: "like"},     //存款人的手机号
		{Key: "acc.account", Val: req.IncomeAccount, EqType: "like"}, //存款人的账号
		{Key: "acc2.account", Val: req.Account, EqType: "like"},      //收款的商户账号
		{Key: "acc3.phone", Val: req.RecvPhone, EqType: "like"},      //收款人手机号
		{Key: "acc3.account", Val: req.RecvAccount, EqType: "like"},  //收款人手机号
	})
	//统计
	where := whereModel.WhereStr
	args := whereModel.Args
	var total sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM income_order inc " +
		" LEFT JOIN account acc ON acc.uid = inc.act_acc_no " +
		" LEFT JOIN servicer ser ON ser.servicer_no = inc.servicer_no " +
		" LEFT JOIN account acc2 ON acc2.uid = ser.account_no " +
		" LEFT JOIN account acc3 ON acc3.uid = inc.recv_acc_no " +
		" LEFT JOIN writeoff w ON w.income_order_no = inc.log_no " + where
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by inc.create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT inc.log_no, inc.act_acc_no, inc.amount, inc.servicer_no, inc.create_time, inc.order_status, inc.finish_time, inc.query_time, inc.balance_type, inc.fees, w.code," +
		" acc.phone, acc.account, acc2.account, acc3.phone, acc3.account, inc.ree_rate, inc.real_amount " +
		" FROM income_order inc " +
		" LEFT JOIN account acc ON acc.uid = inc.act_acc_no " +
		" LEFT JOIN servicer ser ON ser.servicer_no = inc.servicer_no " +
		" LEFT JOIN account acc2 ON acc2.uid = ser.account_no " +
		" LEFT JOIN account acc3 ON acc3.uid = inc.recv_acc_no " +
		" LEFT JOIN writeoff w ON w.income_order_no = inc.log_no " + where2
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
	for rows.Next() {
		var data go_micro_srv_cust.IncomeOrderData
		var writeoff sql.NullString
		var queryTime sql.NullString
		var feeRate sql.NullString
		var realAmount sql.NullString
		if err = rows.Scan(
			&data.LogNo,
			&data.ActAccNo,
			&data.Amount,
			&data.ServicerNo,
			&data.CreateTime,

			&data.OrderStatus,
			&data.FinishTime,
			&queryTime,
			&data.BalanceType,
			&data.Fees,

			&writeoff,
			&data.IncomePhone,
			&data.IncomeAccount,
			&data.Account,
			&data.RecvPhone,
			&data.RecvAccount,

			&feeRate,
			&realAmount,
		); err != nil {
			ss_log.Error("logNo=[%v] /n err=[%v]", data.LogNo, err)
			continue
		}
		data.WriteOff = writeoff.String
		data.QueryTime = queryTime.String
		data.FeeRate = feeRate.String
		data.RealAmount = realAmount.String
		datas = append(datas, &data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * 获取兑换流水
 */
func (*CustHandler) GetExchangeOrderList(ctx context.Context, req *go_micro_srv_cust.GetExchangeOrderListRequest, reply *go_micro_srv_cust.GetExchangeOrderListReply) error {
	if req.StartTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.StartTime) {
			ss_log.Error("GetExchangeOrderList StartTime格式不正确,StartTime: %s", req.StartTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}
	if req.EndTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.EndTime) {
			ss_log.Error("GetExchangeOrderList EndTime格式不正确,StartTime: %s", req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var datas []*go_micro_srv_cust.ExchangeOrderData

	in_type := ""
	out_type := ""
	switch req.ExchangeType {
	case "":
	case constants.Exchange_Usd_To_Khr:
		in_type = "usd"
		out_type = "khr"
	case constants.Exchange_Khr_To_Usd:
		in_type = "khr"
		out_type = "usd"
	default:
		ss_log.Error("ExchangeType参数异常")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "eo.log_no", Val: req.LogNo, EqType: "="},
		{Key: "acc.phone", Val: req.Phone, EqType: "like"},
		{Key: "acc.account", Val: req.Account, EqType: "like"},
		{Key: "eo.create_time", Val: req.StartTime, EqType: ">"},
		{Key: "eo.create_time", Val: req.EndTime, EqType: "<"},
		{Key: "eo.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "eo.in_type", Val: in_type, EqType: "="},
		{Key: "eo.out_type", Val: out_type, EqType: "="},
	})

	//统计
	where := whereModel.WhereStr
	args := whereModel.Args
	var total sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM exchange_order eo " +
		" LEFT JOIN account acc ON acc.uid = eo.account_no " + where
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by eo.create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT eo.log_no, eo.in_type, eo.out_type, eo.amount, eo.create_time, eo.rate, eo.order_status, eo.finish_time, eo.account_no, eo.trans_from, eo.trans_amount, eo.err_reason, eo.fees" +
		", acc.phone, acc.account " +
		" FROM exchange_order eo " +
		" LEFT JOIN account acc ON acc.uid = eo.account_no " + where2
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
	for rows.Next() {
		data := go_micro_srv_cust.ExchangeOrderData{}
		var finishTime sql.NullString
		if err = rows.Scan(
			&data.LogNo,
			&data.InType,
			&data.OutType,
			&data.Amount,
			&data.CreateTime,
			&data.Rate,
			&data.OrderStatus,
			&finishTime,
			&data.AccountNo,
			&data.TransFrom,
			&data.TransAmount,
			&data.ErrReason,

			&data.Fees,
			&data.Phone,
			&data.Account,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		inTypeB := data.InType
		outTypeB := data.OutType
		data.ExchangeType = fmt.Sprintf("%s_to_%s", inTypeB, outTypeB)

		//if data.Rate != "" {
		//	tempD, _ := decimal.NewFromString(data.Rate)
		//	tempD = ss_big.SsBigInst.ToRound(tempD.Div(decimal.NewFromInt(10000)), 2, ss_big.RoundingMode_CEILING)
		//	data.Rate = tempD.String()
		//}

		if finishTime.String != "" {
			data.FinishTime = finishTime.String
		}

		datas = append(datas, &data)
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * 获取虚拟账户日志
 */
func (*CustHandler) GetLogVaccounts(ctx context.Context, req *go_micro_srv_cust.GetLogVaccountsRequest, reply *go_micro_srv_cust.GetLogVaccountsReply) error {
	if req.StartTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.StartTime) {
			ss_log.Error("GetLogVaccounts StartTime格式不正确,StartTime: %s", req.StartTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}
	if req.EndTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.EndTime) {
			ss_log.Error("GetLogVaccounts EndTime格式不正确,StartTime: %s", req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "lv.log_no", Val: req.LogNo, EqType: "like"},
		{Key: "lv.biz_log_no", Val: req.BizLogNo, EqType: "like"},
		{Key: "acc.account", Val: req.Account, EqType: "like"},
		{Key: "lv.reason", Val: req.Reason, EqType: "="},
		{Key: "lv.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "lv.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "lv.op_type", Val: req.OpType, EqType: "="},
		{Key: "vacc.va_type", Val: req.VaType, EqType: "="},
		{Key: "vacc.balance_type", Val: req.BalanceType, EqType: "="},
	})

	//统计
	where := whereModel.WhereStr
	args := whereModel.Args
	var total sql.NullString
	sqlCnt := "SELECT count(1) FROM log_vaccount lv " +
		" LEFT JOIN vaccount vacc ON vacc.vaccount_no = lv.vaccount_no " +
		" LEFT JOIN account acc ON acc.uid = vacc.account_no " + where
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by lv.create_time desc,case reason when "+constants.VaReason_FEES+" then 1 end")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT lv.log_no, lv.vaccount_no, lv.create_time, lv.amount, lv.op_type, lv.frozen_balance, lv.balance, lv.reason, lv.settle_hourly_log_no, lv.settle_daily_log_no, lv.biz_log_no," +
		" vacc.balance_type, vacc.account_no, vacc.va_type, " +
		" acc.nickname, acc.account, acc2.account " +
		" FROM log_vaccount lv " +
		" LEFT JOIN vaccount vacc ON vacc.vaccount_no = lv.vaccount_no " +
		" LEFT JOIN account acc ON acc.uid = vacc.account_no " +
		" LEFT JOIN admin_account acc2 ON acc2.uid = vacc.account_no " + where2

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

	var datas []*go_micro_srv_cust.LogVaccountData
	for rows.Next() {
		data := &go_micro_srv_cust.LogVaccountData{}
		var settleHourlyLogNo, settleDailyLogNo sql.NullString
		var nickname, account, account2 sql.NullString
		if err = rows.Scan(
			&data.LogNo,
			&data.VaccountNo,
			&data.CreateTime,
			&data.Amount,
			&data.OpType,

			&data.FrozenBalance,
			&data.Balance,
			&data.Reason,
			&settleHourlyLogNo,
			&settleDailyLogNo,

			&data.BizLogNo,
			&data.BalanceType,
			&data.AccountNo,
			&data.VaType,
			&nickname,
			&account,
			&account2,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}

		data.OrderType = data.Reason

		switch data.Reason {
		case constants.VaReason_FEES:
			whereList := []*model.WhereSqlCond{
				{Key: "lv.biz_log_no", Val: data.BizLogNo, EqType: "="},
				{Key: "vacc.account_no", Val: data.AccountNo, EqType: "="},
				{Key: "lv.reason", Val: constants.VaReason_FEES, EqType: "!="}, //非服务费的
				{Key: "lv.create_time", Val: data.CreateTime, EqType: "="},     //同一创建时间的
			}

			//如果是手续费，则要返回产生该笔手续费的订单类型,方便可以查询该笔手续费的订单详情
			orderType := dao.LogVaccountDaoInst.GetFeesOrderType(whereList)
			if orderType != "" {
				data.OrderType = orderType
			}

		default:
		}

		//我需要冻结变化金额，余额变化金额，两者的符号
		data.AlterBalnace, data.AlterBalnaceSymbol, data.AlterFrozenBalance, data.AlterFrozenBalanceSymbol = computeAmount(data.Amount, data.OpType, data.OrderType, data.VaType)

		data.Nickname = nickname.String
		data.Account = account.String
		if account2.String != "" {
			data.Account = account2.String
		}

		data.SettleHourlyLogNo = settleHourlyLogNo.String
		data.SettleDailyLogNo = settleDailyLogNo.String

		datas = append(datas, data)
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

//为更好的直观的看清除金额的变化，这里是获取冻结变化金额，余额变化金额，两者的符号(0不变，1+,2-)
func computeAmount(amount, opType string, reason string, vaType string) (alterBalnace, alterBalnaceSymbol, alterFrozenBalance, alterFrozenBalanceSymbol string) {
	alterBalnaceT := ""
	alterBalnaceSymbolT := ""

	alterFrozenBalanceT := ""
	alterFrozenBalanceSymbolT := ""

	switch reason {
	case constants.VaReason_Exchange:
		switch opType {
		case constants.VaOpType_Add:
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		case constants.VaOpType_Minus:
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Minus

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_INCOME: //充值
		switch opType {
		case constants.VaOpType_Add:
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		case constants.VaOpType_Defreeze_Add:
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Minus

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Minus
		case constants.VaOpType_Defreeze_But_Minus:
			alterBalnaceT = "0"
			alterBalnaceSymbolT = "0"

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Add
		default:
			ss_log.Error("未知opType:[%v]，reason:[%v]", opType, reason)
		}
	case constants.VaReason_OUTGO: //提现
		switch opType {
		case constants.VaOpType_Minus:
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Minus

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"

			if strext.ToInt(vaType) == constants.VaType_QUOTA_USD_REAL { //美金实时额度
				alterBalnaceT = amount
				alterBalnaceSymbolT = constants.VaOpType_Add

				alterFrozenBalanceT = "0"
				alterFrozenBalanceSymbolT = "0"
			}

		case constants.VaOpType_Freeze:
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Minus

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Add
		case constants.VaOpType_Defreeze:
			alterBalnaceT = "0"
			alterBalnaceSymbolT = "0"

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Minus
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_TRANSFER: //转账
		switch opType {
		case constants.VaOpType_Minus: //申请就冻结,
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Minus

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		case constants.VaOpType_Add: //申请就冻结,
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_COLLECTION: //收款
	case constants.VaReason_FEES: //手续费
		switch opType {
		case constants.VaOpType_Add:
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		case constants.VaOpType_Minus:
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_Cancel_withdraw: //pos 端取消提现
		switch opType {
		case constants.VaOpType_Defreeze_Add:
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Minus
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_PROFIT_OUTGO: //平台盈利提现
		switch opType {
		case constants.VaOpType_Add:
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		case constants.VaOpType_Minus:
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Minus

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_Cust_Withdraw: //客户向总部提现
		switch opType {
		case constants.VaOpType_Add:
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		case constants.VaOpType_Freeze: //申请就冻结,
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Minus

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Add
		case constants.VaOpType_Defreeze: //成功就解冻
			alterBalnaceT = "0"
			alterBalnaceSymbolT = "0"

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Minus

		case constants.VaOpType_Defreeze_Add: //驳回就解冻并增加
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Minus

		default:
			ss_log.Error("未知reason[%v],opType[%v],amount[%v]", reason, opType, amount)
		}
	case constants.VaReason_Cust_Cancel_Withdraw: // 驳回客户向总部提现
		switch opType {
		case constants.VaOpType_Defreeze_Add: //驳回就解冻并增加
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Minus

		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_Cust_Save: //客户向总部充值
		switch opType {
		case constants.VaOpType_Freeze: //申请就冻结,
			alterBalnaceT = "0"
			alterBalnaceSymbolT = "0"

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Add
		case constants.VaOpType_Defreeze_Minus:
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Minus

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		case constants.VaOpType_Defreeze_Add:
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Minus
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_Srv_Save: //服务商向总部充值
		switch opType {
		case constants.VaOpType_Add:
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		case constants.VaOpType_Defreeze_But_Minus:
			alterBalnaceT = "0"
			alterBalnaceSymbolT = "0"

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Minus
		case constants.VaOpType_Defreeze_Add:
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Add
		case constants.VaOpType_Defreeze_Minus:
			alterBalnaceT = "0"
			alterBalnaceSymbolT = "0"

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Add
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_Srv_Withdraw: //服务商向总部提现
		switch opType {
		case constants.VaOpType_Balance_Frozen_Add: //申请
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Minus

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Add
		case constants.VaOpType_Balance_Defreeze_Add: //驳回
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Minus
		case constants.VaOpType_Defreeze: //通过
			alterBalnaceT = "0"
			alterBalnaceSymbolT = "0"

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Minus

		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_Cust_Cancel_Save: //驳回客户向总部存款
		switch opType {
		case constants.VaOpType_Defreeze_But_Minus: //
			alterBalnaceT = "0"
			alterBalnaceSymbolT = "0"

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Minus

		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_Cust_Pay_Order: //客户支付商家的订单
		switch opType {
		case constants.VaOpType_Minus: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Minus

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_Business_Payee: //商家收到客户的订单
		switch opType {
		case constants.VaOpType_Add: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_Srv_CashRecharge: //服务商现金充值
		switch opType {
		case constants.VaOpType_Add: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_Business_Settle: //商家结算
		switch opType {
		case constants.VaOpType_Add: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		case constants.VaOpType_Minus: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Minus

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_Business_Save: //商家充值
		switch opType {
		case constants.VaOpType_Add: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_Business_Withdraw: //商家提现
		switch opType {
		case constants.VaOpType_Freeze: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Minus

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Add
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_Business_Cancel_Withdraw: //商家提现
		switch opType {
		case constants.VaOpType_Defreeze_Add: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Minus
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_ChangeCustBalance: //改变服务商余额
		switch opType {
		case constants.VaOpType_Add: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		case constants.VaOpType_Minus: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Minus

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_ChangeSrvBalance: //改变服务商余额
		switch opType {
		case constants.VaOpType_Add: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		case constants.VaOpType_Minus: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Minus

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_BusinessTransferToBusiness: //商家转账(注意现在商家转个人 是转用户账户而不是商家账户)
		switch opType {
		case constants.VaOpType_Add: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		case constants.VaOpType_Minus: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Minus

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		case constants.VaOpType_Defreeze: //批量转账的转账成功，对转账方来说是解冻
			alterBalnaceT = "0"
			alterBalnaceSymbolT = "0"

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Minus
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_BusinessRefund: //商家退款
		switch opType {
		case constants.VaOpType_Add: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		case constants.VaOpType_Minus: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Minus

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_BusinessBatchTransferToBusiness: //商家批量转账
		switch opType {
		case constants.VaOpType_Freeze: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Minus

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Add
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	case constants.VaReason_PlatformFreeze: //商家批量转账
		switch opType {
		case constants.VaOpType_Add: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = "0"
			alterFrozenBalanceSymbolT = "0"
		case constants.VaOpType_Defreeze: //
			alterBalnaceT = "0"
			alterBalnaceSymbolT = "0"

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Minus
		case constants.VaOpType_Freeze: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Minus

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Add
		case constants.VaOpType_Defreeze_Add: //
			alterBalnaceT = amount
			alterBalnaceSymbolT = constants.VaOpType_Add

			alterFrozenBalanceT = amount
			alterFrozenBalanceSymbolT = constants.VaOpType_Minus
		default:
			ss_log.Error("reason[%v],未知opType:[%v]", reason, opType)
		}
	default:
		ss_log.Error("reason[%v],opType[%v],amount[%v]未有变化金额和符号的方法", reason, opType, amount)
		return "", "", "", ""
	}

	return alterBalnaceT, alterBalnaceSymbolT, alterFrozenBalanceT, alterFrozenBalanceSymbolT
}

/**
 * 转账至用户
 */
func (*CustHandler) GetToCustList(ctx context.Context, req *go_micro_srv_cust.GetToCustListRequest, reply *go_micro_srv_cust.GetToCustListReply) error {
	if req.StartTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.StartTime) {
			ss_log.Error("GetToCustList StartTime格式不正确,StartTime: %s", req.StartTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}
	if req.EndTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.EndTime) {
			ss_log.Error("GetToCustList EndTime格式不正确,StartTime: %s", req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var datas []*go_micro_srv_cust.ToCustData

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ltc.log_no", Val: req.LogNo, EqType: "like"},
		{Key: "ltc.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "ltc.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "ltc.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "acc.account", Val: req.Account, EqType: "like"}, //用户账号
	})
	//统计
	var total sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM log_to_cust ltc " +
		" LEFT JOIN cust cu ON cu.cust_no = ltc.cust_no " +
		" LEFT JOIN account acc ON acc.uid = cu.account_no " +
		" LEFT JOIN card c ON c.card_no = ltc.card_no " + whereModel.WhereStr
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...)
	if err != nil {
		ss_log.Error("查询数量失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//添加limit
	//ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by ltc.create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by case ltc.order_status when "+constants.AuditOrderStatus_Pending+" then 1 end, ltc.create_time desc")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	sqlStr := "SELECT  ltc.log_no,ltc.currency_type,ltc.collection_type" +
		",ltc.amount,ltc.create_time,ltc.order_type,ltc.order_status,ltc.finish_time" +
		",ltc.lat,ltc.lng,ltc.fees,ltc.ip,ltc.image_id " +
		", acc.account,ca.name,ca.card_number,ch.channel_name " +
		" FROM log_to_cust ltc " +
		" LEFT JOIN cust cu ON cu.cust_no = ltc.cust_no " +
		" LEFT JOIN account acc ON acc.uid = cu.account_no " +
		" LEFT JOIN card ca ON ca.card_no = ltc.card_no " +
		" LEFT JOIN channel ch ON ch.channel_no = ca.channel_no " + whereModel.WhereStr

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
	for rows.Next() {
		var data go_micro_srv_cust.ToCustData
		var finishTime, imageId, channelName sql.NullString
		if err = rows.Scan(
			&data.LogNo,
			&data.CurrencyType,
			&data.CollectionType,
			&data.Amount,
			&data.CreateTime,

			&data.OrderType,
			&data.OrderStatus,
			&finishTime,
			&data.Lat,
			&data.Lng,

			&data.Fees,
			&data.Ip,
			&imageId,
			&data.Account,
			&data.Name,

			&data.CardNumber,
			&channelName,
		); err != nil {
			ss_log.Error("err=[%v]，LogNo[%v]", err, data.LogNo)
			continue
		}
		data.FinishTime = finishTime.String
		data.ImageId = imageId.String
		data.ChannelName = channelName.String

		datas = append(datas, &data)
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * 用户转账至总部（充值）
 */
func (*CustHandler) GetCustToHeadquartersList(ctx context.Context, req *go_micro_srv_cust.GetCustToHeadquartersListRequest, reply *go_micro_srv_cust.GetCustToHeadquartersListReply) error {

	if req.StartTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.StartTime) {
			ss_log.Error("GetCustToHeadquartersList StartTime格式不正确,StartTime: %s", req.StartTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}
	if req.EndTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.EndTime) {
			ss_log.Error("GetCustToHeadquartersList EndTime格式不正确,StartTime: %s", req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "lcth.log_no", Val: req.LogNo, EqType: "like"},
		{Key: "lcth.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "lcth.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "lcth.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "acc.account", Val: req.Account, EqType: "like"}, //用户账号
	})
	//统计
	var total sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM log_cust_to_headquarters lcth " +
		" LEFT JOIN cust cu ON cu.cust_no = lcth.cust_no " +
		" LEFT JOIN account acc ON acc.uid = cu.account_no " +
		" LEFT JOIN card_head c ON c.card_no = lcth.card_no " + whereModel.WhereStr
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...)
	if err != nil {
		ss_log.Error("查询数量失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//添加limit
	//ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by lcth.create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by case lcth.order_status when "+constants.AuditOrderStatus_Pending+" then 1 end, lcth.create_time desc")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, errGet := dao.LogCustToHeadquartersDaoInst.CustToHeadquartersList(whereModel.WhereStr, whereModel.Args)
	if errGet != nil {
		ss_log.Error("获取数据失败")
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * app用户查询提现订单
 */
func (*CustHandler) GetLogToCusts(ctx context.Context, req *go_micro_srv_cust.GetLogToCustsRequest, reply *go_micro_srv_cust.GetLogToCustsReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var datas []*go_micro_srv_cust.ToCustData

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "acc.uid", Val: req.AccountUid, EqType: "="}, //用户账号
		{Key: "ltc.currency_type", Val: req.CurrencyType, EqType: "="},
	})
	//统计
	var total sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM log_to_cust ltc " +
		" LEFT JOIN cust cu ON cu.cust_no = ltc.cust_no " +
		" LEFT JOIN account acc ON acc.uid = cu.account_no " + whereModel.WhereStr
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...)
	if err != nil {
		ss_log.Error("查询数量失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by ltc.create_time desc")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	sqlStr := "SELECT  ltc.log_no,ltc.currency_type ,ltc.amount,ltc.create_time" +
		",ltc.order_status,ltc.finish_time,ltc.fees " +
		" FROM log_to_cust ltc " +
		" LEFT JOIN cust cu ON cu.cust_no = ltc.cust_no " +
		" LEFT JOIN account acc ON acc.uid = cu.account_no " + whereModel.WhereStr

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
	for rows.Next() {
		data := go_micro_srv_cust.ToCustData{}
		var finishTime sql.NullString
		if err = rows.Scan(
			&data.LogNo,
			&data.CurrencyType,
			&data.Amount,
			&data.CreateTime,

			&data.OrderStatus,
			&finishTime,
			&data.Fees,
		); err != nil {
			ss_log.Error("err=[%v]，LogNo[%v]", err, data.LogNo)
			continue
		}
		if finishTime.String != "" {
			data.FinishTime = finishTime.String
		}

		//统一转换成前端对接的订单状态。
		switch data.OrderStatus {
		case constants.AuditOrderStatus_Pending:
			data.OrderStatus = constants.OrderStatus_Pending
			data.OpType = constants.VaOpType_Minus
		case constants.AuditOrderStatus_Passed:
			data.OrderStatus = constants.OrderStatus_Paid
			data.OpType = constants.VaOpType_Minus
		case constants.AuditOrderStatus_Deny:
			data.OrderStatus = constants.OrderStatus_Err
			data.OpType = constants.VaOpType_Add
		default:
			ss_log.Error("OrderStatus[%v]错误, LogNo[%v]", data.OrderStatus, data.LogNo)
			continue
		}

		data.OrderType = constants.VaReason_Cust_Withdraw

		datas = append(datas, &data)
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * app用户查询银行卡充值订单
 */
func (*CustHandler) GetLogCustToHeadquarters(ctx context.Context, req *go_micro_srv_cust.GetLogCustToHeadquartersRequest, reply *go_micro_srv_cust.GetLogCustToHeadquartersReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var datas []*go_micro_srv_cust.CustToHeadquartersData

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "lcth.currency_type", Val: req.CurrencyType, EqType: "="},
		{Key: "acc.uid", Val: req.AccountUid, EqType: "="}, //用户
	})
	//统计
	var total sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM log_cust_to_headquarters lcth " +
		" LEFT JOIN cust cu ON cu.cust_no = lcth.cust_no " +
		" LEFT JOIN account acc ON acc.uid = cu.account_no " + whereModel.WhereStr
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...)
	if err != nil {
		ss_log.Error("查询数量失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//添加limit
	//ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by lcth.create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by lcth.create_time desc")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	sqlStr := "SELECT lcth.log_no,lcth.currency_type,lcth.collection_type" +
		",lcth.amount,lcth.create_time,lcth.order_status,lcth.finish_time,lcth.fees " +
		" FROM log_cust_to_headquarters lcth " +
		" LEFT JOIN cust cu ON cu.cust_no = lcth.cust_no " +
		" LEFT JOIN account acc ON acc.uid = cu.account_no " + whereModel.WhereStr

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
	for rows.Next() {
		data := go_micro_srv_cust.CustToHeadquartersData{}
		var finishTime sql.NullString
		if err = rows.Scan(
			&data.LogNo,
			&data.CurrencyType,
			&data.CollectionType,
			&data.Amount,
			&data.CreateTime,

			&data.OrderStatus,
			&finishTime,
			&data.Fees,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		if finishTime.String != "" {
			data.FinishTime = finishTime.String
		}
		//统一转换成前端对接的订单状态。
		switch data.OrderStatus {
		case constants.AuditOrderStatus_Pending:
			data.OrderStatus = constants.OrderStatus_Pending
		case constants.AuditOrderStatus_Passed:
			data.OrderStatus = constants.OrderStatus_Paid
			data.OpType = constants.VaOpType_Add
		case constants.AuditOrderStatus_Deny:
			data.OrderStatus = constants.OrderStatus_Err
		default:
			ss_log.Error("OrderStatus[%v]错误, LogNo[%v]", data.OrderStatus, data.LogNo)
			continue
		}

		data.OrderType = constants.VaReason_Cust_Save

		datas = append(datas, &data)
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * app用户查询网点提现订单
 */
func (*CustHandler) CustOutgoBills(ctx context.Context, req *go_micro_srv_cust.CustOutgoBillsRequest, reply *go_micro_srv_cust.CustOutgoBillsReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ou.balance_type", Val: req.CurrencyType, EqType: "="},
		{Key: "vacc.account_no", Val: req.AccountUid, EqType: "="}, //用户
	})
	//统计
	var total sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM outgo_order ou " +
		" LEFT JOIN vaccount vacc ON vacc.vaccount_no = ou.vaccount_no " + whereModel.WhereStr
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...)
	if err != nil {
		ss_log.Error("查询数量失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//添加limit
	//ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by ou.create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by ou.create_time desc")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	sqlStr := "SELECT ou.log_no, ou.amount, ou.create_time, ou.order_status, ou.balance_type, ou.fees, ou.modify_time " +
		" FROM outgo_order ou " +
		" LEFT JOIN vaccount vacc ON vacc.vaccount_no = ou.vaccount_no " + whereModel.WhereStr

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
	var datas []*go_micro_srv_cust.CustOutgoBillData
	for rows.Next() {
		data := go_micro_srv_cust.CustOutgoBillData{}
		if err = rows.Scan(
			&data.LogNo,
			&data.Amount,
			&data.CreateTime,

			&data.OrderStatus,
			&data.CurrencyType,
			&data.Fees,
			&data.FinishTime,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}

		if data.OrderStatus == constants.OrderStatus_Cancel {
			data.OrderStatus = constants.OrderStatus_Err
		}
		if data.OrderStatus == constants.OrderStatus_Pending_Confirm {
			data.OrderStatus = constants.OrderStatus_Pending
		}

		data.OrderType = constants.VaReason_OUTGO

		datas = append(datas, &data)
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * app用户查询网点充值订单
 */
func (*CustHandler) CustIncomeBills(ctx context.Context, req *go_micro_srv_cust.CustIncomeBillsRequest, reply *go_micro_srv_cust.CustIncomeBillsReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "inc.balance_type", Val: req.CurrencyType, EqType: "="},
		{Key: "inc.recv_acc_no", Val: req.AccountUid, EqType: "="}, //用户
	})
	//统计
	var total sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM income_order inc " + whereModel.WhereStr
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...)
	if err != nil {
		ss_log.Error("查询数量失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by inc.create_time desc")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	sqlStr := "SELECT inc.log_no, inc.amount, inc.create_time, inc.order_status, inc.balance_type, inc.fees, inc.finish_time  " +
		" FROM income_order inc " + whereModel.WhereStr

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
	var datas []*go_micro_srv_cust.CustIncomeBillData
	for rows.Next() {
		data := go_micro_srv_cust.CustIncomeBillData{}
		if err = rows.Scan(
			&data.LogNo,
			&data.Amount,
			&data.CreateTime,

			&data.OrderStatus,
			&data.CurrencyType,
			&data.Fees,
			&data.FinishTime,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}

		data.OrderType = constants.VaReason_INCOME

		datas = append(datas, &data)
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}
