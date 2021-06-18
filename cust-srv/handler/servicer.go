package handler

import (
	"a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/cust-srv/dao"
	"context"
)

/**
 */
func (*CustHandler) ServicerBills(ctx context.Context, req *custProto.ServicerBillsRequest, reply *custProto.ServicerBillsReply) error {

	accountUid := ""  //从jwt获取的
	accountType := "" //从jwt获取的
	idenNo := ""      //从jwt获取的

	opNo := "" //前端传来的

	switch req.AccountType {
	case constants.AccountType_SERVICER:
		accountUid = req.AccountUid
		opNo = req.OpNo
	case constants.AccountType_POS:
		accountType = req.AccountType
		idenNo = req.IdenNo
	default:
		ss_log.Error("账号类型不是服务商或店员，参数accountType:[%v]", req.AccountType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	whereList := []*model.WhereSqlCond{
		{Key: "to_char(bdr.create_time,'yyyy-MM')", Val: req.StartTime, EqType: "="},
		{Key: "bdr.currency_type", Val: req.CurrencyType, EqType: "="},
		{Key: "bdr.account_no", Val: accountUid, EqType: "="},
		{Key: "bdr.account_type", Val: accountType, EqType: "="},
		{Key: "bdr.op_acc_no", Val: idenNo, EqType: "="},
		{Key: "bdr.op_acc_no", Val: opNo, EqType: "="},
	}

	// 账单类型条件
	billTypeWhere := &model.WhereSqlCond{
		Key: "bdr.bill_type", Val: "('" +
			constants.BILL_TYPE_INCOME + "','" +
			constants.BILL_TYPE_OUTGO + "','" +
			constants.BILL_TYPE_ChangeBalance +
			"') ", EqType: "in"}
	if req.BillType != "" {
		billTypeWhere = &model.WhereSqlCond{Key: "bdr.bill_type", Val: req.BillType, EqType: "="}
	}
	whereList = append(whereList, billTypeWhere)

	// 订单状态条件
	orderStatusWhere := &model.WhereSqlCond{Key: "bdr.order_status", Val: "('" + constants.OrderStatus_Paid + "','" + constants.OrderStatus_Cancel + "','" + constants.OrderStatus_Pending_Confirm + "') ", EqType: "in"}
	if req.OrderStatus != "" {
		orderStatusWhere = &model.WhereSqlCond{Key: "bdr.order_status", Val: req.OrderStatus, EqType: "="}
	}
	whereList = append(whereList, orderStatusWhere)

	//全部数量统计
	total := dao.BillingDetailsResultsDaoInst.GetCnt(whereList)

	//统计用的
	datas, err := dao.BillingDetailsResultsDaoInst.GetBillingDetailsResults(whereList, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//
	whereList2 := []*model.WhereSqlCond{
		{Key: "account_no", Val: accountUid, EqType: "="},
		{Key: "account_type", Val: accountType, EqType: "="},
		{Key: "op_acc_no", Val: idenNo, EqType: "="},
		{Key: "op_acc_no", Val: opNo, EqType: "="},
	}

	orderStatusWhere2 := &model.WhereSqlCond{
		Key: "order_status", Val: "('" +
			constants.OrderStatus_Paid + "','" +
			constants.OrderStatus_Cancel + "','" +
			constants.OrderStatus_Pending_Confirm +
			"') ", EqType: "in"}
	if req.OrderStatus != "" {
		orderStatusWhere2 = &model.WhereSqlCond{Key: "order_status", Val: req.OrderStatus, EqType: "="}
	}
	whereList2 = append(whereList2, orderStatusWhere2)

	times := []string{}
	for _, data := range datas {
		isNewTime := container.GetKey(data.CreateTime[:10], times) //-1找不到
		if isNewTime == -1 {
			data.Time = data.CreateTime[:10]
			times = append(times, data.Time)
			//存取款统计
			da := dao.BillingDetailsResultsDaoInst.GetServicerBillsDaySum(data.Time, whereList2, req.BillType)
			data.UsdIncomeSum = da.UsdIncomeSum
			data.KhrIncomeSum = da.KhrIncomeSum
			data.UsdOutgoSum = da.UsdOutgoSum
			data.KhrOutgoSum = da.KhrOutgoSum
		}
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	return nil
}

/**
 * pos获取服务商授权收款额度
 */
func (*CustHandler) ServicerCollectLimit(ctx context.Context, req *custProto.ServicerCollectLimitRequest, reply *custProto.ServicerCollectLimitReply) error {
	//dbHandler := db.GetDB(constants.DB_CRM)
	//defer db.PutDB(constants.DB_CRM, dbHandler)
	//reply.Datas = &go_micro_srv_cust.ServicerCollectLimitData{}

	////查询服务商的授权收款额度
	//usdAuthCollectLimit, _ := dao.ServiceDaoInst.GetSerAuthCollectLimit(dbHandler, req.AccountNo, "usd")
	//reply.Datas.UsdAuthCollectLimit = usdAuthCollectLimit
	//
	//khrAuthCollectLimit, _ := dao.ServiceDaoInst.GetSerAuthCollectLimit(dbHandler, req.AccountNo, "khr")
	//reply.Datas.KhrAuthCollectLimit = khrAuthCollectLimit
	//
	//usdSpentAuthCollectLimit, _ := dao.ServiceDaoInst.GetSerAuthCollectLimit(dbHandler, req.AccountNo, "usd_spent")
	//reply.Datas.UsdNoSpentCollectLimit = strext.ToStringNoPoint(strext.ToInt64(usdAuthCollectLimit) - strext.ToInt64(usdSpentAuthCollectLimit))
	//
	//khrSpentAuthCollectLimit, _ := dao.ServiceDaoInst.GetSerAuthCollectLimit(dbHandler, req.AccountNo, "khr_spent")
	//reply.Datas.KhrNoSpentCollectLimit = strext.ToStringNoPoint(strext.ToInt64(khrAuthCollectLimit) - strext.ToInt64(khrSpentAuthCollectLimit))

	// 获取 serAccNo
	var serAccNo string
	switch req.AccountType {
	case constants.AccountType_POS: // 销售员
		serAccNo = dao.CashierDaoInstance.GetSrvAccNoFromCaNo(req.AccountUid)
		if serAccNo == "" {
			ss_log.Error("err=[服务商账号id错误,销售员账号id为----->%s]", req.AccountUid)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

	case constants.AccountType_SERVICER: // 服务商
		serAccNo = req.AccountUid
	}

	// 获取虚账额度
	// usd 实时额度
	_, usdRealBalance := dao.VaccountDaoInst.GetBalanceFromAccNo(serAccNo, constants.VaType_QUOTA_USD_REAL)
	// usd 授权额度
	_, usdAuthBalance := dao.VaccountDaoInst.GetBalanceFromAccNo(serAccNo, constants.VaType_QUOTA_USD)
	// khr 实时额度
	_, khrRealBalance := dao.VaccountDaoInst.GetBalanceFromAccNo(serAccNo, constants.VaType_QUOTA_KHR_REAL)
	// khr 授权额度
	_, khrAuthBalance := dao.VaccountDaoInst.GetBalanceFromAccNo(serAccNo, constants.VaType_QUOTA_KHR)
	ss_log.Info("美金授权额度为--->%s ,使用额度为--->%s, 瑞尔授权额度为--->%s, 使用额度为--->%s", usdAuthBalance, usdRealBalance, khrAuthBalance, khrRealBalance)
	var useKhr, UseUsd string

	UseUsd = ss_count.Sub(usdAuthBalance, usdRealBalance).String()

	useKhr = ss_count.Sub(khrAuthBalance, khrRealBalance).String()

	data := &custProto.ServicerCollectLimitData{
		UsdNoSpentCollectLimit: UseUsd,
		KhrNoSpentCollectLimit: useKhr,
	}

	reply.Datas = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
