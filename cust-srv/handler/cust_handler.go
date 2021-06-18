package handler

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

	"a.a/cu/db"
	"a.a/cu/encrypt"
	"a.a/cu/ss_img"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/aws_s3"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	custProto "a.a/mp-server/common/proto/cust"
	pushProto "a.a/mp-server/common/proto/push"
	quotaProto "a.a/mp-server/common/proto/quota"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"

	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/common"
	"a.a/mp-server/cust-srv/dao"
	"a.a/mp-server/cust-srv/i"
	"a.a/mp-server/cust-srv/util"
)

type CustHandler struct {
	Client custProto.CustService
}

var CustHandlerInst CustHandler

/**
 * 获取会员列表
 */
func (*CustHandler) GetCustList(ctx context.Context, req *custProto.GetCustListRequest, reply *custProto.GetCustListReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "acc.uid", Val: req.Uid, EqType: "="},
		{Key: "acc.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "acc.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "acc.nickname", Val: req.QueryNickname, EqType: "like"},
		{Key: "acc.phone", Val: req.QueryPhone, EqType: "like"},
	})

	//统计
	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	var total sql.NullString
	sqlCnt := "SELECT count(1) " +
		"FROM cust c LEFT JOIN account acc ON acc.uid = c.account_no " + where2
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args2...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	where := whereModel.WhereStr
	args := whereModel.Args
	sqlStr := "SELECT c.cust_no, c.gender, c.in_authorization, c.out_authorization,c.in_transfer_authorization,c.out_transfer_authorization," +
		"acc.nickname, acc.create_time, acc.uid,,acc.phone " +
		" FROM cust c LEFT JOIN account acc ON acc.uid = c.account_no " + where

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("CustHandler|GetCustList|err=%v\nreq=[%v]\nsql=[%v]", err, req, sqlStr)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	var datas []*custProto.CustData
	for rows.Next() {
		var data custProto.CustData
		err = rows.Scan(
			&data.CustNo,
			&data.Gender,
			&data.InAuthorization,
			&data.OutAuthorization,
			&data.InTransferAuthorization,
			&data.OutTransferAuthorization,

			&data.Nickname,
			&data.CreateTime,
			&data.Uid,
			&data.Phone,
		)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		datas = append(datas, &data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	return nil
}

func (*CustHandler) GetCustInfo(ctx context.Context, req *custProto.GetCustRequest, reply *custProto.GetCustReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "acc.uid", Val: req.Uid, EqType: "="},
	}

	data, err := dao.CustDaoInst.GetCustInfo(whereList)
	if err != nil {
		ss_log.Error("err[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	cardWhereList := []*model.WhereSqlCond{
		{Key: "ca.account_no", Val: req.Uid, EqType: "="},
		{Key: "ca.collect_status", Val: "1", EqType: "="},
		{Key: "ca.is_delete", Val: "0", EqType: "="},
		{Key: "ca.account_type", Val: constants.AccountType_USER, EqType: "="},
	}

	cardTotal, errCnt := dao.CardDaoInst.GetCustCardTotal(cardWhereList)
	if errCnt != nil {
		ss_log.Error("查询用户卡数量出错,err[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	cardDatas, errCards := dao.CardDaoInst.WebAdminGetCustCards(cardWhereList)
	if errCards != nil {
		ss_log.Error("查询用户卡信息出错,err[%v]", errCards)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.CustData = data
	reply.CardDatas = cardDatas
	reply.CardTotal = cardTotal
	return nil
}

/**
 * WEB获取指定用户账单明细列表
 */
func (*CustHandler) GetCustBills(ctx context.Context, req *custProto.GetCustBillsRequest, reply *custProto.GetCustBillsReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var datas []*custProto.CustBillDetailData

	opType := ""
	if strings.Contains(req.Reason, "_") { //转账入4_1，转账出4_2
		reason := req.Reason
		req.Reason = reason[0:1]
		opType = reason[2:]
	}

	if req.StartTime != "" {
		if !ss_time.CheckDateTimeIsRight("2006/01/02 15:04:05", req.StartTime) {
			ss_log.Error("参数错误:日期格式错误,StartTime:%s", req.StartTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	if req.EndTime != "" {
		if !ss_time.CheckDateTimeIsRight("2006/01/02 15:04:05", req.EndTime) {
			ss_log.Error("参数错误:日期格式错误,EndTime:%s", req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	// 检查结束时间是否大于等于开始时间
	if req.StartTime != "" && req.EndTime != "" {
		if cmp, _ := ss_time.CompareDate("2006/01/02 15:04:05", req.StartTime, req.EndTime); cmp > 0 {
			ss_log.Error("参数错误:开始时间大于结束时间,StartTime:%s,EndTime:%s", req.StartTime, req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	if req.Page < 1 || req.PageSize < 1 {
		ss_log.Error("参数错误:Page:%d,PageSize:%d", req.Page, req.PageSize)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "lv.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "lv.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "lv.biz_log_no", Val: req.BizLogNo, EqType: "like"},
		{Key: "vacc.balance_type", Val: req.BalanceType, EqType: "="},
		{Key: "lv.reason", Val: req.Reason, EqType: "="},
		{Key: "lv.op_type", Val: opType, EqType: "="},
		{Key: "vacc.account_no", Val: req.Uid, EqType: "="},
	})
	//只要个人的
	strs := " and vacc.va_type in ('" + strext.ToStringNoPoint(constants.VaType_USD_DEBIT) +
		"','" + strext.ToStringNoPoint(constants.VaType_FREEZE_USD_DEBIT) +
		"','" + strext.ToStringNoPoint(constants.VaType_KHR_DEBIT) + "','" +
		strext.ToStringNoPoint(constants.VaType_FREEZE_KHR_DEBIT) + "')"
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, strs)

	extraStr := getShowUseBillExtraStr()
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, extraStr)

	var total sql.NullString
	sqlCnt := "SELECT count(1) " +
		"FROM log_vaccount lv " +
		"LEFT JOIN vaccount vacc ON vacc.vaccount_no = lv.vaccount_no " + whereModel.WhereStr
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by lv.create_time desc,case reason when "+constants.VaReason_FEES+" then 1 end")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	sqlStr := "SELECT lv.create_time, lv.biz_log_no, lv.amount, lv.balance, lv.op_type, lv.reason, vacc.balance_type " +
		"FROM log_vaccount lv " +
		"LEFT JOIN vaccount vacc ON vacc.vaccount_no = lv.vaccount_no " + whereModel.WhereStr

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("err=[%v]\n", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	for rows.Next() {
		var data custProto.CustBillDetailData
		err = rows.Scan(
			&data.CreateTime,
			&data.LogNo,
			&data.Amount,
			&data.Balance,
			&data.OpType,
			&data.Reason,
			&data.BalanceType,
		)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}

		data.OrderType = data.Reason

		switch data.Reason {
		case constants.VaReason_FEES: //如果是手续费，则要返回产生该笔手续费的订单类型,方便可以查询该笔手续费的订单详情
			whereList := []*model.WhereSqlCond{
				{Key: "lv.biz_log_no", Val: data.LogNo, EqType: "="},
				{Key: "vacc.account_no", Val: req.Uid, EqType: "="},
				{Key: "lv.reason", Val: constants.VaReason_FEES, EqType: "!="}, //非服务费的
				{Key: "lv.create_time", Val: data.CreateTime, EqType: "="},     //同一创建时间的
			}
			orderType := dao.LogVaccountDaoInst.GetFeesOrderType(whereList)
			//这是产生该笔手续费的订单类型
			data.OrderType = orderType
			switch orderType {
			case "":
				ss_log.Error("查询产生该笔手续费订单类型失败")
			case constants.VaReason_OUTGO: //取款手续费处理
				//这是产生该笔手续费的订单类型
				data.OrderType = constants.VaReason_OUTGO
			default:

			}
		case constants.VaReason_TRANSFER: //转账
			//如果是转账则要分清是转入还是转出，将其设置为与WEB前端对接的4_1:转账入，4_2:转账出
			switch data.OpType {
			case constants.VaOpType_Add:
				data.Reason = data.Reason + "_" + constants.VaOpType_Add
			case constants.VaOpType_Minus:
				data.Reason = data.Reason + "_" + constants.VaOpType_Minus
			}
		case constants.VaReason_Cancel_withdraw: //取款取消的
			if data.OpType == constants.VaOpType_Defreeze_Add { //op_type是6（解冻），说明是取款取消的，那么op_type应当是1（加钱）
				data.OpType = constants.VaOpType_Add
			}
			data.Amount = "0"
			//因为不显示取消取款的手续费,而取款取消订单的余额与用户的余额对不上，所以该余额应加上取款取消手续费金额，
			//查询该取消取款订单的取消取款手续费。
			fees := dao.OutgoOrderDaoInst.GetOutgoOrderFees(data.LogNo)
			data.Balance = ss_count.Add(data.Balance, fees)
			data.OrderType = constants.VaReason_OUTGO
		default:

		}

		datas = append(datas, &data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * app获取用户账单明细列表
 */
func (*CustHandler) CustBills(ctx context.Context, req *custProto.CustBillsRequest, reply *custProto.CustBillsReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "vacc.balance_type", Val: req.CurrencyType, EqType: "="},
		{Key: "to_char(lv.create_time,'yyyy-MM')", Val: req.QueryTime, EqType: "="},
		{Key: "vacc.account_no", Val: req.AccountUid, EqType: "="},
	})

	//只要个人的
	strs := " and vacc.va_type in ('" + strext.ToStringNoPoint(constants.VaType_USD_DEBIT) +
		"','" + strext.ToStringNoPoint(constants.VaType_FREEZE_USD_DEBIT) +
		"','" + strext.ToStringNoPoint(constants.VaType_KHR_DEBIT) + "','" +
		strext.ToStringNoPoint(constants.VaType_FREEZE_KHR_DEBIT) + "')"
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, strs)

	//获取用户收入、支出统计（此统计必需在白名单之前，因为有些成功账单不显示，具体看白名单）
	incomeSum, spendingSum := getUseBillSum(whereModel)

	//获取白名单
	extraStr := getShowUseBillExtraStr()
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, extraStr)

	//全部数量统计
	total := dao.LogVaccountDaoInst.GetCnt(whereModel.WhereStr, whereModel.Args)

	//添加排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by lv.create_time desc,case reason when "+constants.VaReason_FEES+" then 1 end")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.LogVaccountDaoInst.GetAppCustBills(dbHandler, whereModel.WhereStr, whereModel.Args, req.AccountUid)
	if err != ss_err.ERR_SUCCESS {
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	reply.IncomeSum = incomeSum     //收入统计
	reply.SpendingSum = spendingSum //支出统计
	reply.Total = strext.ToInt32(total)
	return nil
}

//获取用户账单要显示的白名单(todo 有修改的字段记得修改注释，不然后面维护会更乱)
func getShowUseBillExtraStr() string {

	//黑名单
	//银行卡提现成功、手续费(9-4、6-4)
	//pos机扫码取款成功、手续费(3-4、6-4)

	//白名单
	extraStr := " AND ( "

	//手续费(6-2、6-3、6-6)
	extraStr += " ( lv.reason = '" + constants.VaReason_FEES + "' AND lv.op_type in( '" + constants.VaOpType_Minus + "','" + constants.VaOpType_Freeze + "','" + constants.VaOpType_Defreeze_Add + "' )) "

	//兑换发起币种的虚帐和手续费(1-2,6-2)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_Exchange + "' AND lv.op_type = '" + constants.VaOpType_Minus + "' )"

	//兑换成的币种的虚帐(1-1)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_Exchange + "' AND lv.op_type = '" + constants.VaOpType_Add + "' )"

	//pos机存款(2-1)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_INCOME + "' AND lv.op_type = '" + constants.VaOpType_Add + "' )"

	//pos机扫码取款申请、手续费(3-3、6-3)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_OUTGO + "' AND lv.op_type = '" + constants.VaOpType_Freeze + "' )"

	//pos机手机号取款成功、手续费(3-2、6-2)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_OUTGO + "' AND lv.op_type = '" + constants.VaOpType_Minus + "' )"

	//付款账单是和转账账单记账方式一致(实质就是转账)
	//转账出、手续费(4-2、6-2)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_TRANSFER + "' AND lv.op_type = '" + constants.VaOpType_Minus + "' )"

	//转账入(4-1)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_TRANSFER + "' AND lv.op_type = '" + constants.VaOpType_Add + "' )"

	//pos机扫码取款驳回、手续费(7-6、6-6)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_Cancel_withdraw + "' AND lv.op_type = '" + constants.VaOpType_Defreeze_Add + "' )"

	//银行卡提现申请、手续费(9-3、6-3)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_Cust_Withdraw + "' AND lv.op_type = '" + constants.VaOpType_Freeze + "' )"

	//银行卡提现驳回、手续费(10-6、6-6)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_Cust_Cancel_Withdraw + "' AND lv.op_type = '" + constants.VaOpType_Defreeze_Add + "' )"

	//用户银行卡充值成功、手续费(11-1、6-2)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_Cust_Save + "' AND lv.op_type = '" + constants.VaOpType_Add + "' ) "

	//用户支付商家的订单(15-2)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_Cust_Pay_Order + "' AND lv.op_type = '" + constants.VaOpType_Minus + "' ) "

	//平台修改用户余额(22-1,22-2)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_ChangeCustBalance + "' AND lv.op_type = '" + constants.VaOpType_Add + "' )"
	extraStr += " OR ( lv.reason = '" + constants.VaReason_ChangeCustBalance + "' AND lv.op_type = '" + constants.VaOpType_Minus + "' )"

	//商家转账给个人(24-1)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_BusinessTransferToBusiness + "' AND lv.op_type = '" + constants.VaOpType_Add + "' )"

	//商家退款(25-1)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_BusinessRefund + "' AND lv.op_type = '" + constants.VaOpType_Add + "' )"

	//商家退款(27-1,3,6)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_PlatformFreeze + "' AND lv.op_type in(" +
		"'" + constants.VaOpType_Add + "', " + "'" + constants.VaOpType_Defreeze_Add + "', " + "'" + constants.VaOpType_Freeze + "'" + "))"

	extraStr += " )"

	return extraStr
}

//获取用户收入、支出统计
func getUseBillSum(whereModel *model.WhereSql) (incomeSumT string, spendingSumT string) {
	//收入统计
	incomeSumM := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhereList(incomeSumM, []*model.WhereSqlCond{
		{Key: "lv.reason", Val: "('" + constants.VaReason_INCOME + "','" +
			constants.VaReason_COLLECTION + "','" +
			constants.VaReason_TRANSFER + "','" +
			constants.VaReason_Exchange + "','" +
			constants.VaReason_Cust_Save + "','" +
			constants.VaReason_BusinessTransferToBusiness +
			"') ", EqType: "in"},
		{Key: "lv.op_type", Val: "('" + constants.VaOpType_Add + "','" + constants.VaOpType_Defreeze_Add + "') ", EqType: "in"},
	}, 0, 0)
	incomeSum := dao.LogVaccountDaoInst.GetSumAmt(incomeSumM.WhereStr, incomeSumM.Args)

	//核销码注销总金额
	platformFreeze := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhereList(platformFreeze, []*model.WhereSqlCond{
		{Key: "lv.reason", Val: constants.VaReason_PlatformFreeze, EqType: "="},
		{Key: "lv.op_type", Val: constants.VaOpType_Defreeze, EqType: "="},
	}, 0, 0)
	platformFreezeSum := dao.LogVaccountDaoInst.GetSumAmt(platformFreeze.WhereStr, platformFreeze.Args)

	incomeSum = ss_count.Sub(incomeSum, platformFreezeSum).String()

	//支出统计
	spendingSumM := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhereList(spendingSumM, []*model.WhereSqlCond{
		{Key: "lv.reason", Val: "('" + constants.VaReason_OUTGO + "','" + constants.VaReason_TRANSFER + "','" + constants.VaReason_Exchange + "','" + constants.VaReason_Cust_Pay_Order + "')", EqType: "in"},
		{Key: "lv.op_type", Val: "('" + constants.VaOpType_Minus + "','" + constants.VaOpType_Defreeze_Minus + "') ", EqType: "in"},
	}, 0, 0)
	spendingSum := dao.LogVaccountDaoInst.GetSumAmt(spendingSumM.WhereStr, spendingSumM.Args)

	//支出再加上 后面才开发的银行卡提现成功金额
	toCustSumM := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhereList(toCustSumM, []*model.WhereSqlCond{
		{Key: "lv.reason", Val: constants.VaReason_Cust_Withdraw, EqType: "="},
		{Key: "lv.op_type", Val: constants.VaOpType_Defreeze, EqType: "="},
	}, 0, 0)
	toCustSum := dao.LogVaccountDaoInst.GetSumAmt(toCustSumM.WhereStr, toCustSumM.Args)

	//此处要加上手续费的付出
	feesM := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhereList(feesM, []*model.WhereSqlCond{
		{Key: "lv.reason", Val: constants.VaReason_FEES, EqType: "="},
		{Key: "lv.op_type", Val: "('" + constants.VaOpType_Minus + "','" + constants.VaOpType_Defreeze + "') ", EqType: "in"},
	}, 0, 0)
	feeSum := dao.LogVaccountDaoInst.GetSumAmt(feesM.WhereStr, feesM.Args)

	spendingSum = ss_count.Add(spendingSum, toCustSum) //后来添加的银行卡提现成功的订单金额
	spendingSum = ss_count.Add(spendingSum, feeSum)    //加手续费

	return incomeSum, spendingSum
}

/**
 * pos获取总部收款账号
 */
func (c *CustHandler) GetHeadquartersCards(ctx context.Context, req *custProto.GetHeadquartersCardsRequest, reply *custProto.GetHeadquartersCardsReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var datas []*custProto.HeadquartersCardsData
	var total sql.NullString
	//获取总部账号
	_, accPlat, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")
	switch req.AccountType {
	case constants.AccountType_SERVICER:
		whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
			{Key: "ca.account_no", Val: accPlat, EqType: "="},
			{Key: "ca.collect_status", Val: "1", EqType: "="},
			{Key: "ca.is_delete", Val: "0", EqType: "="},
			{Key: "ca.balance_type", Val: req.MoneyType, EqType: "="},
			{Key: "ca.account_type", Val: req.AccountType, EqType: "="},
		})
		where := whereModel.WhereStr
		args := whereModel.Args
		sqlCnt := "SELECT count(1) " +
			"FROM card_head ca " +
			"LEFT JOIN channel ch ON ch.channel_no = ca.channel_no " + where
		err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args...)

		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}

		ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY ca.is_defalut DESC ")
		ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

		sqlStr := "SELECT ca.card_no, ch.channel_name, ca.name, ca.card_number, ca.balance_type, ch.logo_img_no, ch.logo_img_no_grey, ch.color_begin, ch.color_end " +
			" FROM card_head ca " +
			" LEFT JOIN channel ch ON ch.channel_no = ca.channel_no " + whereModel.WhereStr
		rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
		if stmt != nil {
			defer stmt.Close()
		}
		defer rows.Close()
		if err2 != nil {
			ss_log.Error("err2=[%v]", err2)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}

		for rows.Next() {
			var data custProto.HeadquartersCardsData
			var channelName, logoImgNo, logoImgNoGrey, colorBegin, colorEnd sql.NullString
			err = rows.Scan(
				&data.CardNo,
				&channelName,
				&data.Name,
				&data.CardNumber,
				&data.MoneyType,
				&logoImgNo,
				&logoImgNoGrey,
				&colorBegin,
				&colorEnd,
			)
			if err != nil {
				ss_log.Error("err=[%v]", err)
				continue
			}
			data.ChannelName = channelName.String
			data.ColorBegin = colorBegin.String
			data.ColorEnd = colorEnd.String

			if logoImgNo.String != "" {
				reqImg := &custProto.UnAuthDownloadImageRequest{
					ImageId: logoImgNo.String,
				}
				replyImg := &custProto.UnAuthDownloadImageReply{}
				err2 := c.UnAuthDownloadImage(ctx, reqImg, replyImg)
				if err2 != nil || replyImg.ResultCode != ss_err.ERR_SUCCESS {
					ss_log.Error("获取图片url失败")
				} else {
					data.LogoImgUrl = replyImg.ImageUrl
				}
			}
			if logoImgNoGrey.String != "" {
				reqImg := &custProto.UnAuthDownloadImageRequest{
					ImageId: logoImgNoGrey.String,
				}
				replyImg := &custProto.UnAuthDownloadImageReply{}
				err2 := c.UnAuthDownloadImage(ctx, reqImg, replyImg)
				if err2 != nil || replyImg.ResultCode != ss_err.ERR_SUCCESS {
					ss_log.Error("获取图片url失败")
				} else {
					data.LogoImgUrlGrey = replyImg.ImageUrl
				}
			}
			datas = append(datas, &data)
		}
	case constants.AccountType_USER:
		whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
			{Key: "ca.account_no", Val: accPlat, EqType: "="},
			{Key: "ca.collect_status", Val: "1", EqType: "="},
			{Key: "ca.is_delete", Val: "0", EqType: "="},
			{Key: "chcu.currency_type", Val: req.MoneyType, EqType: "="},
			{Key: "ca.account_type", Val: req.AccountType, EqType: "="},
			{Key: "chcu.use_status", Val: "1", EqType: "="},
		})
		where := whereModel.WhereStr
		args := whereModel.Args
		sqlCnt := "SELECT count(1) " +
			"FROM card_head ca " +
			"LEFT JOIN channel_cust_config chcu ON chcu.id = ca.channel_cust_config_id " +
			"LEFT JOIN channel ch ON ch.channel_no = ca.channel_no " + where
		err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args...)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}

		ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY ca.is_defalut DESC ")
		ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

		sqlStr := "SELECT ca.card_no, ch.channel_name, ca.name, ca.card_number, chcu.currency_type, ch.logo_img_no, ch.logo_img_no_grey, ch.color_begin, ch.color_end " +
			", chcu.save_rate, chcu.withdraw_rate, chcu.withdraw_max_amount, chcu.save_single_min_fee, chcu.withdraw_single_min_fee " +
			", chcu.save_charge_type, chcu.withdraw_charge_type, chcu.support_type, chcu.save_max_amount, chcu.channel_no " +
			"FROM card_head ca " +
			"LEFT JOIN channel_cust_config chcu ON chcu.id = ca.channel_cust_config_id " +
			"LEFT JOIN channel ch ON ch.channel_no = chcu.channel_no " + whereModel.WhereStr
		rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
		if stmt != nil {
			defer stmt.Close()
		}
		defer rows.Close()
		if err2 != nil {
			ss_log.Error("err2=[%v]", err2)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}

		for rows.Next() {
			var data custProto.HeadquartersCardsData
			var channelName, logoImgNo, logoImgNoGrey, colorBegin, colorEnd sql.NullString
			err = rows.Scan(
				&data.CardNo,
				&channelName,
				&data.Name,
				&data.CardNumber,
				&data.MoneyType,

				&logoImgNo,
				&logoImgNoGrey,
				&colorBegin,
				&colorEnd,
				&data.SaveRate,
				&data.WithdrawRate,
				&data.WithdrawMaxAmount,
				&data.SaveSingleMinFee,

				&data.WithdrawSingleMinFee,
				&data.SaveChargeType,
				&data.WithdrawChargeType,
				&data.SupportType,
				&data.SaveMaxAmount,
				&data.ChannelNo,
			)
			if err != nil {
				ss_log.Error("err=[%v]", err)
				continue
			}
			data.ChannelName = channelName.String
			data.ColorBegin = colorBegin.String
			data.ColorEnd = colorEnd.String

			if logoImgNo.String != "" {
				reqImg := &custProto.UnAuthDownloadImageRequest{
					ImageId: logoImgNo.String,
				}
				replyImg := &custProto.UnAuthDownloadImageReply{}
				err2 := c.UnAuthDownloadImage(ctx, reqImg, replyImg)
				if err2 != nil || replyImg.ResultCode != ss_err.ERR_SUCCESS {
					ss_log.Error("获取图片url失败")
				} else {
					data.LogoImgUrl = replyImg.ImageUrl
				}
			}
			if logoImgNoGrey.String != "" {
				reqImg := &custProto.UnAuthDownloadImageRequest{
					ImageId: logoImgNoGrey.String,
				}
				replyImg := &custProto.UnAuthDownloadImageReply{}
				err2 := c.UnAuthDownloadImage(ctx, reqImg, replyImg)
				if err2 != nil || replyImg.ResultCode != ss_err.ERR_SUCCESS {
					ss_log.Error("获取图片url失败")
				} else {
					data.LogoImgUrlGrey = replyImg.ImageUrl
				}
			}

			data.Temp = "0"

			datas = append(datas, &data)
		}
	case constants.AccountType_PersonalBusiness: //收个人商户的平台账户和收企业商户的平台账户是一样的。。
		fallthrough
	case constants.AccountType_EnterpriseBusiness:
		whereList := []*model.WhereSqlCond{
			{Key: "ca.account_no", Val: accPlat, EqType: "="},
			{Key: "ca.collect_status", Val: "1", EqType: "="},
			{Key: "ca.is_delete", Val: "0", EqType: "="},
			{Key: "chcu.currency_type", Val: req.MoneyType, EqType: "="},
			{Key: "ca.account_type", Val: constants.AccountType_EnterpriseBusiness, EqType: "="},
			{Key: "chcu.use_status", Val: "1", EqType: "="},
		}
		total.String = dao.CardHeadDaoInst.GetHeadCardBusinessCnt(whereList)
		datasT, err := dao.CardHeadDaoInst.GetHeadCardBusinessList2(whereList, req.Page, req.PageSize)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}
		datas = datasT
	default:
		ss_log.Error("AccountType参数错误:[%v]", req.AccountType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	reply.DataList = datas
	reply.Total = strext.ToInt32(total.String)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
查询绑定的卡列表信息（用户、服务商）
*/
func (c *CustHandler) GetUserCards(ctx context.Context, req *custProto.GetUserCardsRequest, reply *custProto.GetUserCardsReply) error {

	var total string
	var datas []*custProto.UserCardsData
	switch req.AccountType {
	case constants.AccountType_SERVICER:
		err := ""
		whereList := []*model.WhereSqlCond{
			{Key: "ca.account_no", Val: req.AccountNo, EqType: "="},
			{Key: "ca.collect_status", Val: "1", EqType: "="},
			{Key: "ca.is_delete", Val: "0", EqType: "="},
			{Key: "ca.balance_type", Val: req.BalanceType, EqType: "="},
			{Key: "ca.account_type", Val: constants.AccountType_SERVICER, EqType: "="},
		}
		datas, total, err = dao.CardDaoInst.GetServicerCards(whereList)
		if err != ss_err.ERR_SUCCESS {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}

	case constants.AccountType_POS:
		err, servicerNo := dao.ServiceDaoInst.GetServiceByCashierNo(req.IdenNo)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_DB_OP_SER
			return nil
		}
		accNo, err2 := dao.ServiceDaoInst.GetAccountNoByServicerNo(servicerNo)
		if err2 != ss_err.ERR_SUCCESS {
			ss_log.Error("查询店员的服务商账号失败,err=[%v]", err2)
			reply.ResultCode = err2
			return nil
		}

		errCards := ""
		whereList := []*model.WhereSqlCond{
			{Key: "ca.account_no", Val: accNo, EqType: "="},
			{Key: "ca.collect_status", Val: "1", EqType: "="},
			{Key: "ca.is_delete", Val: "0", EqType: "="},
			{Key: "ca.balance_type", Val: req.BalanceType, EqType: "="},
			{Key: "ca.account_type", Val: constants.AccountType_SERVICER, EqType: "="},
		}
		datas, total, errCards = dao.CardDaoInst.GetServicerCards(whereList)
		if errCards != ss_err.ERR_SUCCESS {
			ss_log.Error("err=[%v]", errCards)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}

	case constants.AccountType_USER:
		err := ""
		whereList := []*model.WhereSqlCond{
			{Key: "ca.account_no", Val: req.AccountNo, EqType: "="},
			{Key: "ca.collect_status", Val: "1", EqType: "="},
			{Key: "ca.is_delete", Val: "0", EqType: "="},
			{Key: "ca.balance_type", Val: req.BalanceType, EqType: "="},
			{Key: "ca.account_type", Val: constants.AccountType_USER, EqType: "="},
		}

		datas, total, err = dao.CardDaoInst.GetCustCards(whereList)
		if err != ss_err.ERR_SUCCESS {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}

	default:
		ss_log.Error("AccountType类型错误[%v]", req.AccountType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	for _, v := range datas {
		if v.LogoImgNo != "" { //查询图片对应的url
			reqImg := &custProto.UnAuthDownloadImageRequest{
				ImageId: v.LogoImgNo,
			}
			replyImg := &custProto.UnAuthDownloadImageReply{}
			err2 := c.UnAuthDownloadImage(ctx, reqImg, replyImg)
			if err2 != nil || replyImg.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取图片url失败")
			} else {
				v.LogoImgUrl = replyImg.ImageUrl
			}
		}
		if v.LogoImgNoGrey != "" { //查询图片对应的url
			reqImg := &custProto.UnAuthDownloadImageRequest{
				ImageId: v.LogoImgNoGrey,
			}
			replyImg := &custProto.UnAuthDownloadImageReply{}
			err2 := c.UnAuthDownloadImage(ctx, reqImg, replyImg)
			if err2 != nil || replyImg.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取图片url失败")
			} else {
				v.LogoImgUrlGrey = replyImg.ImageUrl
			}
		}
	}

	reply.DataList = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
查询绑定的银行卡详情信息（用户、服务商）
*/
func (c *CustHandler) GetUserCardDetail(ctx context.Context, req *custProto.GetUserCardDetailRequest, reply *custProto.GetUserCardDetailReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	data := &custProto.UserCardsData{}
	var err error
	switch req.AccountType {
	case constants.AccountType_SERVICER:
		fallthrough
	case constants.AccountType_POS:
		data, err = dao.CardDaoInst.GetServicerCardDetail(req.CardNo)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}
	case constants.AccountType_USER:
		data, err = dao.CardDaoInst.GetCustCardDetail(req.CardNo)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}
	default:
		ss_log.Error("AccountType类型错误[%v]", req.AccountType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if data.LogoImgNo != "" { //查询图片对应的url
		reqImg := &custProto.UnAuthDownloadImageRequest{
			ImageId: data.LogoImgNo,
		}
		replyImg := &custProto.UnAuthDownloadImageReply{}
		err2 := c.UnAuthDownloadImage(ctx, reqImg, replyImg)
		if err2 != nil || replyImg.ResultCode != ss_err.ERR_SUCCESS {
			ss_log.Error("获取图片url失败")
		} else {
			data.LogoImgUrl = replyImg.ImageUrl
		}
	}

	if data.LogoImgNoGrey != "" { //查询图片对应的url
		reqImg := &custProto.UnAuthDownloadImageRequest{
			ImageId: data.LogoImgNoGrey,
		}
		replyImg := &custProto.UnAuthDownloadImageReply{}
		err2 := c.UnAuthDownloadImage(ctx, reqImg, replyImg)
		if err2 != nil || replyImg.ResultCode != ss_err.ERR_SUCCESS {
			ss_log.Error("获取图片url失败")
		} else {
			data.LogoImgUrlGrey = replyImg.ImageUrl
		}
	}

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
 *	修改默认卡
 */
func (*CustHandler) UpdateDefCard(ctx context.Context, req *custProto.UpdateDefCardRequest, reply *custProto.UpdateDefCardReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlUpdate := "update cust set def_pay_no = $2 where cust_no=$1 "
	err := ss_sql.Exec(dbHandler, sqlUpdate, req.CustNo, req.DefPayNo)
	if err != nil {
		ss_log.Error("err=[%v]\nreq=[%v]\nsql=[%v]", err, req, sqlUpdate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//删除原来币种的默认卡
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 上传图片,前端提交图片的base64字符串
func (c *CustHandler) UploadImage(ctx context.Context, req *custProto.UploadImageRequest, reply *custProto.UploadImageReply) error {
	// 生成文件名文件名
	idName, _ := i.IDW.NextId()
	point := strext.ToStringNoPoint(idName)

	if strings.HasPrefix(req.ImageStr, "data:image/jpeg;base64,") {
		s := strings.Split(req.ImageStr, "data:image/jpeg;base64,")
		req.ImageStr = s[1]
	} else if strings.HasPrefix(req.ImageStr, "data:image/png;base64,") {
		s := strings.Split(req.ImageStr, "data:image/png;base64,")
		req.ImageStr = s[1]
	}

	// 对文件名进行hash
	var baseName string
	img, err := encrypt.DoBase64(encrypt.HANDLE_DECRYPT, req.ImageStr)
	if err != nil {
		ss_log.Error("base64解码参数失败, err:%v, ImageStr:%s", err.Error(), req.ImageStr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 获取文件后缀
	ext := ss_img.SsImgInst.GetFileTypeFromMagic(img.([]byte))
	switch ext {
	case "png":
		baseName = encrypt.DoMd5(point) + ".png"
	case "jpeg":
		baseName = encrypt.DoMd5(point) + ".jpeg"
	default:
		ss_log.Error("上传的图片格式不被支持:%s", ext)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 重新生成文件，防止文件内容中有攻击代码
	newImg, _ := ss_img.SsImgInst.RecreateImage(img.([]byte))
	if len(newImg) == 0 {
		ss_log.Error("重新生成文件失败, 原始数据: %s", string(img.([]byte)))
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//todo 添加水印
	if req.AddWatermark == constants.AddWatermark_True {
		ss_log.Info("添加水印开始")

		//获取要添加的水印图片
		_, _, imgId, errGetImg := dao.LangDaoInst.GetLangByKey("水印图片", constants.LANG_TYPE_IMG)
		if errGetImg != nil {
			ss_log.Error("获取水印图片id出错,err[%v]", errGetImg)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		imgDao, errGetImgUrl := dao.ImageDaoInstance.GetImageUrlById(imgId)
		if errGetImgUrl != nil {
			ss_log.Error("获取水印图片url出错，imgId[%v],err[%v]", imgId, errGetImgUrl)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		// 从s3获取图片
		result, s3Err := common.UploadS3.GetObject(imgDao.ImageUrl)
		if s3Err != nil {
			ss_log.Error("从s3获取图片失败,ImageUrl:%s, err:%v", imgDao.ImageUrl, s3Err)
			reply.ResultCode = ss_err.ERR_IMAGE_OP_FAILD
			return nil
		}

		// 读取body内容
		bytes, err := ioutil.ReadAll(result.Body)
		if err != nil {
			ss_log.Error("读取图片内容失败, result:%+v, err:%v", result, err)
			reply.ResultCode = ss_err.ERR_IMAGE_OP_FAILD
			return nil
		}

		//为图片添加水印
		newImgT, errAddImg := ss_img.SsImgInst.AddImgWatermark(newImg, bytes)
		if errAddImg != nil {
			ss_log.Error("添加水印失败,err[%v]", errAddImg)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		newImg = newImgT
		ss_log.Info("添加水印结束")
	}

	var upErr error
	fileName := ""
	isEncrypt := 0
	contentType := http.DetectContentType(newImg)

	ss_log.Info("开始上传文件到s3")
	switch req.Type {
	case constants.UploadImage_Auth: // 需要授权
		isEncrypt = 1
		fileName = path.Join(aws_s3.Private_Img_Dir, strings.TrimRight(baseName, "."+ext)) // 去掉文件名后缀
		_, upErr = common.UploadS3.UploadByContentEncrypt(newImg, fileName)
	case constants.UploadImage_UnAuth: // 不需要授权
		fileName = path.Join(aws_s3.Public_Img_Dir, baseName)
		_, upErr = common.UploadS3.UploadByContent(newImg, fileName, true)
	default:
		reply.ResultCode = ss_err.ERR_PARAM
		ss_log.Error("是否需要授权参数错误:%d", req.Type)
		return nil
	}

	ss_log.Info("上传文件到s3结束, fileName: %s", fileName)

	// 上传失败
	if upErr != nil {
		ss_log.Error("上传文件s3失败, err:%v, newImg:%s, fileName: %s", upErr, string(newImg), fileName)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 入库存记录
	imageId, err := dao.ImageDaoInstance.InsertImage(req.AccountUid, fileName, contentType, ext, isEncrypt)
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err.Error(), "图片保存进数据库失败")
		reply.ResultCode = ss_err.ERR_SAVE_IMAGE_FAILD
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.ImageName = fileName
	reply.ImageId = imageId
	return nil
}

// 下载图片,前端提交图片的url字符串
func (c *CustHandler) AuthDownloadImage(ctx context.Context, req *custProto.AuthDownloadImageRequest, reply *custProto.AuthDownloadImageReply) error {
	// 获取需要授权的图片路径
	//_, imageBaseUrl, _ := cache.ApiDaoInstance.GetGlobalParam("image_base_url")

	//pathStr := dao.GlobalParamDaoInstance.QeuryParamValue(constants.KEY_STORE_AUTH_IMAGE_PATH)
	//根据id查询路径
	imgDao, err2 := dao.ImageDaoInstance.GetImageUrlById(req.ImageId)
	if err2 != nil {
		ss_log.Error("查询图片记录失败,ImageId:%s, err:%v", req.ImageId, err2)
		reply.ResultCode = ss_err.ERR_IMAGE_OP_FAILD
		return nil
	}

	if imgDao.ImageId == "" {
		ss_log.Error("查询图片记录不存在,ImageId:%s", req.ImageId)
		reply.ResultCode = ss_err.ERR_IMAGE_OP_FAILD
		return nil
	}

	if imgDao.ImageUrl == "" {
		ss_log.Error("查询图片url为空,ImageId:%s", req.ImageId)
		reply.ResultCode = ss_err.ERR_IMAGE_OP_FAILD
		return nil
	}

	// 从s3获取图片
	bytes, s3Err := common.UploadS3.GetObjectEncrypt(imgDao.ImageUrl)
	if s3Err != nil {
		ss_log.Error("从s3获取图形失败,ImageUrl:%s, err:%v", imgDao.ImageUrl, s3Err)
		reply.ResultCode = ss_err.ERR_IMAGE_OP_FAILD
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.ImageStr = base64.StdEncoding.EncodeToString(bytes)
	return nil
}

func (c *CustHandler) UnAuthDownloadImage(ctx context.Context, req *custProto.UnAuthDownloadImageRequest, reply *custProto.UnAuthDownloadImageReply) error {
	// 获取不需要授权的图片路径
	_, imageBaseUrl, _ := cache.ApiDaoInstance.GetGlobalParam("image_base_url")

	//pathStr := dao.GlobalParamDaoInstance.QeuryParamValue(constants.KEY_STORE_UNAUTH_IMAGE_PATH)
	//根据id查询路径
	imgDao, err2 := dao.ImageDaoInstance.GetImageUrlById(req.ImageId)
	if err2 != nil {
		ss_log.Error("查询图片记录失败,ImageId:%s, err:%v", req.ImageId, err2)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//boolean, _ := file.Exists(pathStr + "/" + name)
	//if !boolean {
	//	ss_log.Error("图片不存在")
	//	reply.ResultCode = ss_err.ERR_IMAGE_OP_FAILD
	//	return nil
	//}

	reply.ResultCode = ss_err.ERR_SUCCESS
	//reply.ImageUrl = path.Join(imageBaseUrl, imgDao.ImageUrl)
	reply.ImageUrl = imageBaseUrl + "/" + imgDao.ImageUrl
	return nil
}

//由图片id获取图片的base64字符串
func (c *CustHandler) UnAuthDownloadImageBase64(ctx context.Context, req *custProto.UnAuthDownloadImageBase64Request, reply *custProto.UnAuthDownloadImageBase64Reply) error {
	// 获取不需要授权的图片路径
	//_, imageBaseUrl, _ := cache.ApiDaoInstance.GetGlobalParam("image_base_url")

	//pathStr := dao.GlobalParamDaoInstance.QeuryParamValue(constants.KEY_STORE_UNAUTH_IMAGE_PATH)
	//根据id查询路径
	imgDao, err2 := dao.ImageDaoInstance.GetImageUrlById(req.ImageId)
	if err2 != nil {
		ss_log.Error("查询图片记录失败,ImageId:%s, err:%v", req.ImageId, err2)
		reply.ResultCode = ss_err.ERR_IMAGE_OP_FAILD
		return nil
	}

	if imgDao.ImageUrl == "" {
		ss_log.Error("查询图片url为空,ImageId:%s", req.ImageId)
		reply.ResultCode = ss_err.ERR_IMAGE_OP_FAILD
		return nil
	}

	// 从s3获取图片
	result, s3Err := common.UploadS3.GetObject(imgDao.ImageUrl)
	if s3Err != nil {
		ss_log.Error("从s3获取图片失败,ImageUrl:%s, err:%v", imgDao.ImageUrl, s3Err)
		reply.ResultCode = ss_err.ERR_IMAGE_OP_FAILD
		return nil
	}

	// 读取body内容
	bytes, err := ioutil.ReadAll(result.Body)
	if err != nil {
		ss_log.Error("读取图片内容失败, result:%+v, err:%v", result, err)
		reply.ResultCode = ss_err.ERR_IMAGE_OP_FAILD
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.ImageBase64 = base64.StdEncoding.EncodeToString(bytes)
	return nil
}

// 后台编辑用户向总部存款的订单状态
func (c *CustHandler) UpdateCustSave(ctx context.Context, req *custProto.UpdateCustSaveRequest, reply *custProto.UpdateCustSaveReply) error {
	if req.OrderNo == "" {
		ss_log.Error("订单号为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// 判断订单是否是初始化状态
	orderStatus, currencyType, custNo, _, fees, amount := dao.LogCustToHeadquartersDaoInst.QueryOrderStatusFromLogNo(req.OrderNo)
	if orderStatus != constants.AuditOrderStatus_Pending {
		ss_log.Error(" 订单不是在待审核状态,不能充值,OrderNo: %s,订单状态为: %s", req.OrderNo, orderStatus)
		reply.ResultCode = ss_err.ERR_ORDER_STATUS_NO_INIT
		return nil
	}

	// 根据币种获取虚账类型
	var vaType, plantVaType int32
	//var plantVaType int32
	switch currencyType {
	case constants.CURRENCY_USD:
		vaType = constants.VaType_USD_DEBIT
		plantVaType = constants.VaType_USD_FEES
	case constants.CURRENCY_KHR:
		vaType = constants.VaType_KHR_DEBIT
		plantVaType = constants.VaType_KHR_FEES
	default:
		ss_log.Error("用户向总部充值审核,币种错误,MoneyType: %s", currencyType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	// 获取用户账号
	accNo := dao.RelaAccIdenDaoInst.GetAccNo(custNo, constants.AccountType_USER)
	// 确保虚拟账号存在
	recvVaccNo := dao.VaccountDaoInst.ConfirmExistVaccount(accNo, currencyType, strext.ToInt32(vaType))
	appLang, _ := dao.AccDaoInstance.QueryAccountLang(accNo)
	if appLang == "" {
		appLang = constants.LangEnUS
	}

	appMessType := ""
	orderType := ""
	ss_log.Info("用户 %s 当前的语言为--->%s", accNo, appLang)
	msgType := constants.Template_AddSuccess

	description := ""

	switch strext.ToStringNoPoint(req.OrderStatus) {
	case constants.AuditOrderStatus_Passed:
		//余额加
		if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, recvVaccNo, amount, "+", req.OrderNo, constants.VaReason_Cust_Save); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}

		if fees != "" && fees != "0" {
			// 余额减掉手续费
			if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, recvVaccNo, fees, "-", req.OrderNo, constants.VaReason_FEES); errStr != ss_err.ERR_SUCCESS {
				reply.ResultCode = errStr
				return nil
			}

			// ============平台收益====================
			// 查询总部的账号
			_, headAcc, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")
			headVacc := dao.VaccountDaoInst.ConfirmExistVaccount(headAcc, currencyType, plantVaType)
			// 修改总部的临时虚账余额
			if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, headVacc, fees, "+", req.OrderNo, constants.VaReason_Cust_Save); errStr != ss_err.ERR_SUCCESS {
				reply.ResultCode = errStr
				return nil
			}

			//插入手续费盈利
			d := &dao.HeadquartersProfit{
				OrderNo:      req.OrderNo,
				Amount:       fees,
				OrderStatus:  constants.OrderStatus_Paid,
				BalanceType:  strings.ToLower(currencyType),
				ProfitSource: constants.ProfitSource_INCOME,
				OpType:       constants.PlatformProfitAdd,
			}
			_, err := dao.HeadquartersProfitDao.InsertHeadquartersProfit(tx, d)
			if err != nil {
				ss_log.Error("插入手续费盈利失败，data=%v, err=%v", strext.ToJson(d), err)
				reply.ResultCode = ss_err.ERR_SYSTEM
				return nil
			}

			// 修改收益 总部虚账的余额是等于收益表中的可提现余额
			if errStr := dao.HeadquartersProfitCashableDaoInstance.SyncAccProfit(tx, headVacc, fees, currencyType); errStr != ss_err.ERR_SUCCESS {
				reply.ResultCode = errStr
				return nil
			}
		}

		appMessType = constants.LOG_APP_MESSAGES_ORDER_TYPE_Cust_Save_Apply
		orderType = constants.VaReason_Cust_Save

		description = fmt.Sprintf("审核用户银行卡提现订单[%v],操作[%v]", req.OrderNo, "通过")
	case constants.AuditOrderStatus_Deny: // 驳回
		appMessType = constants.LOG_APP_MESSAGES_ORDER_TYPE_Cust_Save_Fail
		orderType = constants.VaReason_Cust_Cancel_Save
		msgType = constants.Template_AddFail
		description = fmt.Sprintf("审核用户银行卡提现订单[%v],操作[%v]", req.OrderNo, "驳回")
	default:
		ss_log.Error("需要修改的订单状态有误,req.OrderStatus: %v", req.OrderStatus)
		reply.ResultCode = ss_err.ERR_OPERATE_FAILD
		return nil
	}
	// 修改状态
	if errStr := dao.LogCustToHeadquartersDaoInst.UpdateStatusFromLogNo(req.OrderNo, req.OrderStatus); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	//todo 添加关键操作日志
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Trading_Order)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	//==================添加消息，推送消息===============
	orderStatus2 := ""
	switch strext.ToStringNoPoint(req.OrderStatus) {
	case constants.AuditOrderStatus_Deny:
		orderStatus2 = constants.OrderStatus_Err
	case constants.AuditOrderStatus_Passed:
		orderStatus2 = constants.OrderStatus_Paid
	}

	errAddMessages := dao.LogAppMessagesDaoInst.AddLogAppMessages(tx, req.OrderNo, appMessType, orderType, accNo, orderStatus2)
	if errAddMessages != ss_err.ERR_SUCCESS {
		ss_log.Error("errAddMessages=[%v]", errAddMessages)
	}

	toAccountType := dao.AccDaoInstance.GetAccountTypeFromAccNo(accNo)

	moneyType := dao.LangDaoInst.GetLangTextByKey(dbHandler, currencyType, appLang)

	// 修正各币种的金额
	amountB := common.NormalAmountByMoneyType(currencyType, amount)

	timeString := time.Now().Format("2006-01-02 15:04:05")
	args := []string{
		timeString, amountB, moneyType,
	}
	lang, _ := dao.AccDaoInstance.QueryAccountLang(accNo)
	if lang == "" || lang == constants.LangEnUS {
		args = []string{
			amount, moneyType, timeString,
		}
	}

	// 消息推送
	ev := &pushProto.PushReqest{
		Accounts: []*pushProto.PushAccout{
			{
				AccountNo:   accNo,
				AccountType: toAccountType,
			},
		},
		TempNo: msgType,
		Args:   args,
	}
	ss_log.Info("publishing %+v\n", ev)
	// publish an event
	if err := common.PushEvent.Publish(context.TODO(), ev); err != nil {
		ss_log.Error("error publishing: %v", err)
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 用户向总部提现订单审核
func (c *CustHandler) UpdateCustWithdraw(ctx context.Context, req *custProto.UpdateCustWithdrawRequest, reply *custProto.UpdateCustWithdrawReply) error {
	if req.OrderNo == "" {
		ss_log.Error("订单号为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// 判断订单是否是初始化状态
	orderStatus, currencyType, custNo, amount, fees := dao.LogToCustDaoInst.QueryOrderStatusFromLogNo(req.OrderNo)
	if orderStatus != constants.AuditOrderStatus_Pending {
		ss_log.Error(" 订单不是在待审核状态,不能申请提现,OrderNo: %s,订单状态为: %s", req.OrderNo, orderStatus)
		reply.ResultCode = ss_err.ERR_ORDER_STATUS_NO_INIT
		return nil
	}
	// 根据币种获取虚账类型
	var vaType, plantVaType int32
	switch currencyType {
	case constants.CURRENCY_USD:
		vaType = constants.VaType_USD_DEBIT
		plantVaType = constants.VaType_USD_FEES
	case constants.CURRENCY_KHR:
		vaType = constants.VaType_KHR_DEBIT
		plantVaType = constants.VaType_KHR_FEES
	default:
		ss_log.Error("用户向总部提现,币种错误,MoneyType: %s", currencyType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil

	}
	// 获取用户账号
	accNo := dao.RelaAccIdenDaoInst.GetAccNo(custNo, constants.AccountType_USER)
	// 确保虚拟账号存在
	recvVaccNo := dao.VaccountDaoInst.ConfirmExistVaccount(accNo, currencyType, strext.ToInt32(vaType))

	appLang, _ := dao.AccDaoInstance.QueryAccountLang(accNo)
	if appLang == "" {
		appLang = constants.LangEnUS
	}

	appMessType := ""
	orderType := ""
	//langText := ""
	//title := ""
	ss_log.Info("用户 %s 当前的语言为--->%s", accNo, appLang)

	msgType := constants.Template_WithdrawSuccess
	imageId := ""

	description := ""

	switch strext.ToStringNoPoint(req.OrderStatus) {
	case constants.AuditOrderStatus_Passed:
		//上传凭证
		if req.ImageBase64 == "" {
			ss_log.Error("未上传凭证图片req.ImageBase64为空")
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		upReg := &custProto.UploadImageRequest{
			ImageStr:     req.ImageBase64,
			AccountUid:   req.LoginUid,
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

		imageId = upReply.ImageId

		// 修改冻结余额
		//amountB := ss_count.Add(amount,fees)
		if errStr := dao.VaccountDaoInst.ModifyVaccFrozenUpperZero(tx, recvVaccNo, amount, req.OrderNo, constants.VaReason_Cust_Withdraw, fees); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}
		if fees != "" && fees != "0" {
			if errStr := dao.VaccountDaoInst.ModifyVaccFrozenUpperZero(tx, recvVaccNo, fees, req.OrderNo, constants.VaReason_FEES, ""); errStr != ss_err.ERR_SUCCESS {
				reply.ResultCode = errStr
				return nil
			}
		}

		// ============平台收益====================
		// 查询总部的账号
		_, headAcc, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")
		headVacc := dao.VaccountDaoInst.ConfirmExistVaccount(headAcc, currencyType, plantVaType)
		// 修改总部的临时虚账余额
		if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, headVacc, fees, "+", req.OrderNo, constants.VaReason_Cust_Withdraw); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}

		// 插入利润表
		d := &dao.HeadquartersProfit{
			OrderNo:      req.OrderNo,
			Amount:       fees,
			OrderStatus:  constants.OrderStatus_Paid,
			BalanceType:  strings.ToLower(currencyType),
			ProfitSource: constants.ProfitSource_WithdrawFee,
			OpType:       constants.PlatformProfitAdd,
		}
		_, err := dao.HeadquartersProfitDao.InsertHeadquartersProfit(tx, d)
		if err != nil {
			ss_log.Error("插入手续费盈利失败，data=%v, err=%v", strext.ToJson(d), err)
			reply.ResultCode = ss_err.ERR_SYSTEM
			return nil
		}

		// 修改收益 总部虚账的余额是等于收益表中的可提现余额
		if errStr := dao.HeadquartersProfitCashableDaoInstance.SyncAccProfit(tx, headVacc, fees, currencyType); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}

		appMessType = constants.LOG_APP_MESSAGES_ORDER_TYPE_Cust_Withdraw_Apply
		orderType = constants.VaReason_Cust_Withdraw

		description = fmt.Sprintf("审核用户银行转账充值订单[%v],操作[%v] ", req.OrderNo, "通过")

	case constants.AuditOrderStatus_Deny:
		// 恢复余额
		if errStr := dao.VaccountDaoInst.ModifyVaccRemainAndFrozenUpperZero(tx, "+", recvVaccNo, amount, req.OrderNo, constants.VaReason_Cust_Cancel_Withdraw, constants.VaOpType_Defreeze_Add); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}
		if fees != "" && fees != "0" {
			if errStr := dao.VaccountDaoInst.ModifyVaccRemainAndFrozenUpperZero(tx, "+", recvVaccNo, fees, req.OrderNo, constants.VaReason_FEES, constants.VaOpType_Defreeze_Add); errStr != ss_err.ERR_SUCCESS {
				reply.ResultCode = errStr
				return nil
			}
		}

		appMessType = constants.LOG_APP_MESSAGES_ORDER_TYPE_Cust_Withdraw_Fail
		orderType = constants.VaReason_Cust_Cancel_Withdraw
		msgType = constants.Template_WithdrawFail
		description = fmt.Sprintf("审核用户银行转账充值订单[%v],操作[%v]", req.OrderNo, "驳回")
	default:
		ss_log.Error("需要修改的订单状态有误,req.OrderStatus: %v", req.OrderStatus)
		reply.ResultCode = ss_err.ERR_OPERATE_FAILD
		return nil
	}
	// 修改状态
	if errStr := dao.LogToCustDaoInst.UpdateStatusFromLogNo(tx, req.OrderNo, req.Notes, imageId, req.OrderStatus); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	//todo 添加关键操作日志
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Trading_Order)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	//todo 推送消息
	orderStatus2 := ""
	switch strext.ToStringNoPoint(req.OrderStatus) {
	case constants.AuditOrderStatus_Passed:
		orderStatus2 = constants.OrderStatus_Paid
	case constants.AuditOrderStatus_Deny:
		orderStatus2 = constants.OrderStatus_Err
	}
	errAddMessages := dao.LogAppMessagesDaoInst.AddLogAppMessages(tx, req.OrderNo, appMessType, orderType, accNo, orderStatus2)
	if errAddMessages != ss_err.ERR_SUCCESS {
		ss_log.Error("errAddMessages=[%v]", errAddMessages)
	}

	toAccountType := dao.AccDaoInstance.GetAccountTypeFromAccNo(accNo)

	moneyType := dao.LangDaoInst.GetLangTextByKey(dbHandler, currencyType, appLang)

	// 修正各币种的金额
	amountB := common.NormalAmountByMoneyType(currencyType, amount)

	timeString := time.Now().Format("2006-01-02 15:04:05")
	args := []string{
		timeString, amountB, moneyType,
	}
	lang, _ := dao.AccDaoInstance.QueryAccountLang(accNo)
	if lang == "" || lang == constants.LangEnUS {
		args = []string{
			amountB, moneyType, timeString,
		}
	}

	// 消息推送
	ev := &pushProto.PushReqest{
		Accounts: []*pushProto.PushAccout{
			{
				AccountNo:   accNo,
				AccountType: toAccountType,
			},
		},
		TempNo: msgType,
		Args:   args,
	}
	ss_log.Info("publishing %+v\n", ev)
	// publish an event
	if err := common.PushEvent.Publish(context.TODO(), ev); err != nil {
		ss_log.Error("error publishing: %v", err)
	}
	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 后台管理系统,对log_to_service订单的操作 1-通过,2-关闭
func (c *CustHandler) UpdateServiceWithdraw(ctx context.Context, req *custProto.UpdateServiceWithdrawRequest, reply *custProto.UpdateServiceWithdrawReply) error {
	// 获取订单状态
	status := dao.LogToServiceDaoInstance.QueryOrderStatusFromlogNo(req.OrderNo)
	if status == "" {
		ss_log.Error("查询[%v]订单状态失败", req.OrderNo)
		reply.ResultCode = ss_err.ERR_OPERATE_FAILD
		return nil
	}
	//只有待审核的订单才可以申请
	if status != constants.AuditOrderStatus_Pending {
		ss_log.Error("err=[订单不是在待审核状态,不能申请提现,status----->%s]", status)
		reply.ResultCode = ss_err.ERR_ORDER_STATUS_NO_INIT
		return nil
	}

	currencyType, servicerNo, amount := dao.LogToServiceDaoInstance.QueryLogFromlogNo(req.OrderNo)
	srvAccNo := dao.RelaAccIdenDaoInst.GetAccNo(servicerNo, constants.AccountType_SERVICER)
	if srvAccNo == "" {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	description := ""
	switch req.Status {
	case 1: // 通过操作（冻结-，实时不变）
		quotaReq := &quotaProto.ModifyQuotaRequest{
			CurrencyType: currencyType,
			Amount:       amount,
			AccountNo:    srvAccNo,
			OpType:       constants.QuotaOp_SvrWithdraw,
			LogNo:        req.OrderNo,
		}
		quotaRepl, err2 := i.QuotaHandleInstance.Client.ModifyQuota(ctx, quotaReq)
		if err2 != nil || quotaRepl.ResultCode != ss_err.ERR_SUCCESS {
			ss_log.Error("err=[--------------->%s]", "服务商取款,调用八神的服务失败,操作为服务商取款")
			reply.ResultCode = quotaRepl.ResultCode
			return nil
		}

		// 根据 serviceNo 查询accNo
		serviceAccNo := dao.ServiceDaoInst.GetAccNoFromSrvNo(servicerNo)
		// 插入服务商交易明细
		if errStr := dao.BillingDetailsResultsDaoInst.InsertResult(amount, currencyType, serviceAccNo, constants.AccountType_SERVICER, req.OrderNo, "0", constants.OrderStatus_Paid, constants.BillDetailTypewithdraw, "0", amount); errStr == "" {
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		description = fmt.Sprintf("审核服务商提现订单[%v],操作[%v]", req.OrderNo, "通过")
	case 2: // 关闭操作（冻结-，实时-）
		quotaReq := &quotaProto.ModifyQuotaRequest{
			CurrencyType: currencyType,
			Amount:       amount,
			AccountNo:    srvAccNo,
			OpType:       constants.QuotaOp_SvrWithdraw_Cancel,
			LogNo:        req.OrderNo,
		}
		quotaRepl, err2 := i.QuotaHandleInstance.Client.ModifyQuota(ctx, quotaReq)
		if err2 != nil || quotaRepl.ResultCode != ss_err.ERR_SUCCESS {
			ss_log.Error("err=[--------------->%s]", "服务商取款不通过,调用八神的服务失败,操作为服务商取款不通过")
			reply.ResultCode = quotaRepl.ResultCode
			return nil
		}
		description = fmt.Sprintf("审核服务商提现订单[%v],操作[%v]", req.OrderNo, "驳回")
	default:
		ss_log.Error("Status参数错误")
		reply.ResultCode = ss_err.ERR_OPERATE_FAILD
		return nil
	}

	// 修改状态
	if errStr := dao.LogToServiceDaoInstance.UpdateStatusFromLogNo(req.OrderNo, req.Status); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Trading_Order)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
 * 修改用户配置（充值权限、提现权限、可转账入权限、可转账出权限）
 */
func (*CustHandler) ModifyCustInfo(ctx context.Context, req *custProto.ModifyCustInfoRequest, reply *custProto.ModifyCustInfoReply) error {

	//参数校验
	inAuthorizationsStr, legal1 := util.GetParamZhCn(req.InAuthorization, util.InAuthorization)
	if !legal1 {
		ss_log.Error("InAuthorization %v", inAuthorizationsStr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	outAuthorizationStr, legal2 := util.GetParamZhCn(req.OutAuthorization, util.OutAuthorization)
	if !legal2 {
		ss_log.Error("OutAuthorization %v", outAuthorizationStr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	inTransferAuthorizationStr, legal3 := util.GetParamZhCn(req.InTransferAuthorization, util.InTransferAuthorization)
	if !legal3 {
		ss_log.Error("InTransferAuthorization %v", inTransferAuthorizationStr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	outTransferAuthorizationStr, legal4 := util.GetParamZhCn(req.OutTransferAuthorization, util.OutTransferAuthorization)
	if !legal4 {
		ss_log.Error("OutTransferAuthorization %v", outTransferAuthorizationStr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//查询出老数据，为插入关键操作日志做准备
	whereList := []*model.WhereSqlCond{
		{Key: "c.cust_no", Val: req.CustNo, EqType: "="},
	}
	custDataOld, err := dao.CustDaoInst.GetCustInfo(whereList)
	if err != nil {
		ss_log.Error("获取老数据失败")
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	//修改权限
	if errStr := dao.CustDaoInst.ModifyCustInfo(tx, req.InAuthorization, req.OutAuthorization, req.InTransferAuthorization, req.OutTransferAuthorization, req.CustNo); errStr != ss_err.ERR_SUCCESS {
		ss_log.Error("修改用户信息出错，err=[%v]", errStr)
		reply.ResultCode = errStr
		return nil
	}

	description := "" //描述
	if custDataOld.InAuthorization != req.InAuthorization {
		oldDataStr, _ := util.GetParamZhCn(custDataOld.InAuthorization, util.InAuthorization)
		description = fmt.Sprintf("用户账号[%v] 充值权限由[%v]更改为[%v] ", custDataOld.Account, oldDataStr, inAuthorizationsStr)
	}
	if custDataOld.OutAuthorization != req.OutAuthorization {
		oldDataStr, _ := util.GetParamZhCn(custDataOld.OutAuthorization, util.OutAuthorization)
		if description == "" {
			description = fmt.Sprintf("用户账号[%v] 提现权限由[%v]更改为[%v]", custDataOld.Account, oldDataStr, outAuthorizationStr)
		} else {
			description = fmt.Sprintf("%v,提现权限由[%v]更改为[%v] ", description, oldDataStr, outAuthorizationStr)
		}
	}
	if custDataOld.InTransferAuthorization != req.InTransferAuthorization {
		oldDataStr, _ := util.GetParamZhCn(custDataOld.InTransferAuthorization, util.InTransferAuthorization)
		if description == "" {
			description = fmt.Sprintf("用户账号[%v] 可转账入权限由[%v]更改为[%v]", custDataOld.Account, oldDataStr, inTransferAuthorizationStr)
		} else {
			description = fmt.Sprintf("%v,可转账入权限由[%v]更改为[%v] ", description, oldDataStr, inTransferAuthorizationStr)
		}
	}
	if custDataOld.OutTransferAuthorization != req.OutTransferAuthorization {
		oldDataStr, _ := util.GetParamZhCn(custDataOld.OutTransferAuthorization, util.OutTransferAuthorization)
		if description == "" {
			description = fmt.Sprintf("用户账号[%v] 可转账出权限由[%v]更改为[%v]", custDataOld.Account, oldDataStr, outTransferAuthorizationStr)
		} else {
			description = fmt.Sprintf("%v,可转账出权限由[%v]更改为[%v] ", description, oldDataStr, outTransferAuthorizationStr)
		}
	}

	if description != "" {
		errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Account)
		if errAddLog != nil {
			ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetCommonHelpCount(ctx context.Context, req *custProto.GetCommonHelpCountRequest, reply *custProto.GetCommonHelpCountReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var datas []*custProto.CommonHelpCountData

	//app
	//cnAppData := dao.CommonHelpDaoInst.GetCommonHelpCount(dbHandler, constants.LangZhCN, constants.AppVersionVsType_app)
	//enAppData := dao.CommonHelpDaoInst.GetCommonHelpCount(dbHandler, constants.LangEnUS, constants.AppVersionVsType_app)
	//kmAppData := dao.CommonHelpDaoInst.GetCommonHelpCount(dbHandler, constants.LangKmKH, constants.AppVersionVsType_app)
	//
	////pos
	//cnPosData := dao.CommonHelpDaoInst.GetCommonHelpCount(dbHandler, constants.LangZhCN, constants.AppVersionVsType_pos)
	//enPosData := dao.CommonHelpDaoInst.GetCommonHelpCount(dbHandler, constants.LangEnUS, constants.AppVersionVsType_pos)
	//kmPosData := dao.CommonHelpDaoInst.GetCommonHelpCount(dbHandler, constants.LangKmKH, constants.AppVersionVsType_pos)
	//
	//datas = append(datas, cnAppData)
	//datas = append(datas, enAppData)
	//datas = append(datas, kmAppData)
	//
	//datas = append(datas, cnPosData)
	//datas = append(datas, enPosData)
	//datas = append(datas, kmPosData)

	//mangopay app
	cnMangopayAppData := dao.CommonHelpDaoInst.GetCommonHelpCount(dbHandler, constants.LangZhCN, constants.APPVERSIONVSTYPE_MANGOPAY_APP)
	enMangopayAppData := dao.CommonHelpDaoInst.GetCommonHelpCount(dbHandler, constants.LangEnUS, constants.APPVERSIONVSTYPE_MANGOPAY_APP)
	kmMangopayAppData := dao.CommonHelpDaoInst.GetCommonHelpCount(dbHandler, constants.LangKmKH, constants.APPVERSIONVSTYPE_MANGOPAY_APP)

	//mangopay pos
	cnMangopayPosData := dao.CommonHelpDaoInst.GetCommonHelpCount(dbHandler, constants.LangZhCN, constants.APPVERSIONVSTYPE_MANGOPAY_POS)
	enMangopayPosData := dao.CommonHelpDaoInst.GetCommonHelpCount(dbHandler, constants.LangEnUS, constants.APPVERSIONVSTYPE_MANGOPAY_POS)
	kmMangopayPosData := dao.CommonHelpDaoInst.GetCommonHelpCount(dbHandler, constants.LangKmKH, constants.APPVERSIONVSTYPE_MANGOPAY_POS)

	datas = append(datas, cnMangopayAppData)
	datas = append(datas, enMangopayAppData)
	datas = append(datas, kmMangopayAppData)

	datas = append(datas, cnMangopayPosData)
	datas = append(datas, enMangopayPosData)
	datas = append(datas, kmMangopayPosData)

	reply.Datas = datas
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetCommonHelps(ctx context.Context, req *custProto.GetCommonHelpsRequest, reply *custProto.GetCommonHelpsReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ch.is_delete", Val: "0", EqType: "="},
		{Key: "ch.problem", Val: req.Problem, EqType: "like"},
		{Key: "ch.lang", Val: req.Lang, EqType: "="},
		{Key: "ch.vs_type", Val: req.VsType, EqType: "="},
		{Key: "ch.use_status", Val: req.UseStatus, EqType: "="},
	})

	total := dao.CommonHelpDaoInst.GetCnt(dbHandler, whereModel.WhereStr, whereModel.Args)

	//添加排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by ch.idx asc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.CommonHelpDaoInst.GetCommonHelps(dbHandler, whereModel.WhereStr, whereModel.Args)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetAppCommonHelps(ctx context.Context, req *custProto.GetAppCommonHelpsRequest, reply *custProto.GetAppCommonHelpsReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
		{Key: "problem", Val: req.Problem, EqType: "like"},
		{Key: "lang", Val: req.Lang, EqType: "="},
		{Key: "vs_type", Val: req.VsType, EqType: "="},
		{Key: "use_status", Val: "1", EqType: "="}, //后台管理页面获取时可看到禁用和启用的，app那边则只看到启用的
	})

	total := dao.CommonHelpDaoInst.GetCnt(dbHandler, whereModel.WhereStr, whereModel.Args)

	//添加排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by idx asc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.CommonHelpDaoInst.GetAppCommonHelps(dbHandler, whereModel.WhereStr, whereModel.Args)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetCommonHelpDetail(ctx context.Context, req *custProto.GetCommonHelpDetailRequest, reply *custProto.GetCommonHelpDetailReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ch.is_delete", Val: "0", EqType: "="},
		//{Key: "use_status", Val: "0", EqType: "="},
		{Key: "ch.help_no", Val: req.HelpNo, EqType: "="},
	})

	data, err := dao.CommonHelpDaoInst.GetCommonHelpDetail(dbHandler, whereModel.WhereStr, whereModel.Args)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) ModifyHelpStatus(ctx context.Context, req *custProto.ModifyHelpStatusRequest, reply *custProto.ModifyHelpStatusReply) error {
	switch req.UseStatus {
	case "0":
	case "1":
	default:
		ss_log.Error("UseStatus is no in(0,1)")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	err := dao.CommonHelpDaoInst.ModifyHelpStatus(req.HelpNo, req.UseStatus)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = err
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeleteHelp(ctx context.Context, req *custProto.DeleteHelpRequest, reply *custProto.DeleteHelpReply) error {
	if req.HelpNo == "" {
		ss_log.Error("HelpNo is NULL")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	id := req.HelpNo

	//获取当前要删除记录的idx,语言
	startIdx, lang, errGet1 := dao.CommonHelpDaoInst.GetIdxById(tx, id)
	if errGet1 != ss_err.ERR_SUCCESS {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	endIdx, errGet1 := dao.CommonHelpDaoInst.GetMaxidx(tx, lang)
	if errGet1 != ss_err.ERR_SUCCESS {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//将其要删除的排序序号后面的元素往前移
	for j := startIdx + 1; j <= endIdx; j++ {
		errUp := dao.CommonHelpDaoInst.ReplaceIdx(tx, j, lang)
		if errUp != ss_err.ERR_SUCCESS {
			ss_log.Error("errUp=[%v],i=[%v]", errUp, j)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	err := dao.CommonHelpDaoInst.DeleteHelp(tx, id)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = err
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) InsertOrUpdateHelp(ctx context.Context, req *custProto.InsertOrUpdateHelpRequest, reply *custProto.InsertOrUpdateHelpReply) error {
	if req.HelpNo == "" {
		err := dao.CommonHelpDaoInst.AddHelp(req.Problem, req.Answer, req.Lang, req.VsType)
		if err != nil {
			ss_log.Error("插入帮助信息失败，err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_ADD
			return nil
		}
	} else {
		err := dao.CommonHelpDaoInst.UpdateHelp(req.HelpNo, req.Problem, req.Answer)
		if err != nil {
			ss_log.Error("修改帮助信息失败，err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 交换位置
func (c *CustHandler) SwapHelpIdx(ctx context.Context, req *custProto.SwapHelpIdxRequest, reply *custProto.SwapHelpIdxReply) error {
	switch req.SwapType {
	case constants.SwapType_Up: // up
	case constants.SwapType_Down: // down
	default:
		ss_log.Error("swapType错误")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.Idx == "1" && req.SwapType == constants.SwapType_Up {
		ss_log.Error("swapType错误")
		reply.ResultCode = ss_err.ERR_MERC_IS_UPTOP
		return nil
	}

	//helpNoFrom := dao.CommonHelpDaoInst.GetHelpNo(req.Idx)
	//if helpNoFrom == "" {
	//	ss_log.Error("获取helpNo失败")
	//	reply.ResultCode = ss_err.ERR_PARAM
	//	return nil
	//}

	helpNoFrom := req.HelpNo

	helpNoTo := dao.CommonHelpDaoInst.GetNearIdxHelpNo(req.Idx, req.SwapType, req.Lang, req.VsType)
	if helpNoTo == "" {
		ss_log.Error("获取funcNo失败")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	switch req.SwapType {
	case constants.SwapType_Up:
		// 向上需要交换一下
		x := helpNoTo
		helpNoTo = helpNoFrom
		helpNoFrom = x
	}

	errCode := dao.CommonHelpDaoInst.ExchangeIdx(helpNoFrom, helpNoTo)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", errCode)
		reply.ResultCode = errCode
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (c *CustHandler) ModifyCardsDefalut(context.Context, *custProto.ModifyCardsDefalutRequest, *custProto.ModifyCardsDefalutReply) error {
	panic("implement me")
}

//批量修改消息的是否已读状态
func (c *CustHandler) ModifyAppMessagesIsRead(ctx context.Context, req *custProto.ModifyAppMessagesIsReadRequest, reply *custProto.ModifyAppMessagesIsReadReply) error {
	if err := dao.LogAppMessagesDaoInst.ModiftAllRead(req.AccountNo); err != nil {
		ss_log.Error("将账号的消息全部设置为已读出错误,err=[%v]", err)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil

}

//用户查询自己的消息
func (c *CustHandler) GetLogAppMessages(ctx context.Context, req *custProto.GetLogAppMessagesRequest, reply *custProto.GetLogAppMessagesReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "account_no", Val: req.AccountNo, EqType: "="},
	})

	total, errGet := dao.LogAppMessagesDaoInst.GetLogAppMessagesCnt(dbHandler, whereModel.WhereStr, whereModel.Args)
	if errGet != ss_err.ERR_SUCCESS {
		ss_log.Error("获取数量total,err=[%v]", errGet)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by create_time desc,is_read asc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
	datas, err := dao.LogAppMessagesDaoInst.GetLogAppMessages(dbHandler, whereModel.WhereStr, whereModel.Args)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("GetLogAppMessages err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	return nil
}

func (c *CustHandler) GetConsultationConfigs(ctx context.Context, req *custProto.GetConsultationConfigsRequest, reply *custProto.GetConsultationConfigsReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
		{Key: "use_status", Val: req.UseStatus, EqType: "="},
		{Key: "lang", Val: req.Lang, EqType: "="},
		{Key: "name", Val: req.Name, EqType: "like"},
	})

	total := dao.ConsultationConfigDaoInst.GetCnt(dbHandler, whereModel.WhereStr, whereModel.Args)

	//添加排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by lang desc, idx asc `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.ConsultationConfigDaoInst.GetConsultationConfigs(dbHandler, whereModel.WhereStr, whereModel.Args)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}
	for _, data := range datas {
		if data.LogoImgNo != "" { //返回图片url
			reqImg := &custProto.UnAuthDownloadImageRequest{
				ImageId: data.LogoImgNo,
			}
			replyImg := &custProto.UnAuthDownloadImageReply{}
			err2 := c.UnAuthDownloadImage(ctx, reqImg, replyImg)
			if err2 != nil || replyImg.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取图片url失败")
			} else {
				data.LogoImgUrl = replyImg.ImageUrl
			}
		}
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (c *CustHandler) GetConsultationConfigDetail(ctx context.Context, req *custProto.GetConsultationConfigDetailRequest, reply *custProto.GetConsultationConfigDetailReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
		{Key: "id", Val: req.Id, EqType: "="},
	})

	data, err := dao.ConsultationConfigDaoInst.GetConsultationConfigDetail(dbHandler, whereModel.WhereStr, whereModel.Args)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//LogoImgBase64        string   `protobuf:"bytes,7,opt,name=logo_img_base64,json=logoImgBase64,proto3" json:"logo_img_base64,omitempty"`
	if data.LogoImgNo != "" {
		//由图片id获取base64字符串(用于修改处回显图片)
		reqImg := &custProto.UnAuthDownloadImageBase64Request{
			ImageId: data.LogoImgNo,
		}
		replyImg := &custProto.UnAuthDownloadImageBase64Reply{}
		c.UnAuthDownloadImageBase64(ctx, reqImg, replyImg)
		if replyImg.ResultCode != ss_err.ERR_SUCCESS {
			ss_log.Error("获取图片base64字符串失败")
		} else {
			data.LogoImgBase64 = replyImg.ImageBase64
		}
	}

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//app、pos的返回,现在 联系我们的一个页面的每行是一条记录
func (c *CustHandler) GetAppConsultationConfigDetail(ctx context.Context, req *custProto.GetAppConsultationConfigDetailRequest, reply *custProto.GetAppConsultationConfigDetailReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
		{Key: "use_status", Val: "1", EqType: "="},
		{Key: "lang", Val: req.Lang, EqType: "="},
	})

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " order by idx asc ")
	datas, err := dao.ConsultationConfigDaoInst.GetConsultationConfigs(dbHandler, whereModel.WhereStr, whereModel.Args)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}
	for _, data := range datas {
		if data.LogoImgNo != "" { //返回图片url
			reqImg := &custProto.UnAuthDownloadImageRequest{
				ImageId: data.LogoImgNo,
			}
			replyImg := &custProto.UnAuthDownloadImageReply{}
			err2 := c.UnAuthDownloadImage(ctx, reqImg, replyImg)
			if err2 != nil || replyImg.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取图片url失败")
			} else {
				data.LogoImgUrl = replyImg.ImageUrl
			}
		}
	}

	reply.Datas = datas
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) InsertOrUpdateConsultationConfig(ctx context.Context, req *custProto.InsertOrUpdateConsultationConfigRequest, reply *custProto.InsertOrUpdateConsultationConfigReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	//验证语言是否合法
	str1, legal1 := util.GetParamZhCn(req.Lang, util.Lang)
	if !legal1 {
		ss_log.Error("Lang %v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.Id == "" {
		idxMax, errGet := dao.ConsultationConfigDaoInst.GetMaxIdx(req.Lang)

		if errGet == ss_err.ERR_SUCCESS {
			idxMax = idxMax + 1
		}
		err := dao.ConsultationConfigDaoInst.AddConsultationConfig(tx, req.Name, req.Text, req.LogoImgNo, req.UseStatus, req.Lang, strext.ToStringNoPoint(idxMax))
		if err != ss_err.ERR_SUCCESS {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = err
			return nil
		}
		ss_sql.Commit(tx)
		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	}

	// 删除s3上面的图片
	// 找到对应的url
	logoURL, err2 := dao.ConsultationConfigDaoInst.GetLogoURL(req.Id)
	if err2 != nil {
		ss_log.Error("InsertOrUpdateConsultationConfig 失败,err: %s", err2.Error())
	}
	if logoURL != "" {
		// 删除s3上面的图片
		if _, err := common.UploadS3.DeleteOne(logoURL); err != nil {
			notes := fmt.Sprintf("InsertOrUpdateConsultationConfig s3上删除图片失败,图片路劲为: %s,err: %s", logoURL, err.Error())
			ss_log.Error(notes)
			dao.DictimagesDaoInst.AddDelFaildLog(notes)
		}
		if err := dao.DictimagesDaoInst.Delete(logoURL); err != nil {
			notes := fmt.Sprintf("InsertOrUpdateConsultationConfig 删除图片记录失败,图片路劲为: %s,err: %s", logoURL, err.Error())
			ss_log.Error(notes)
			dao.DictimagesDaoInst.AddDelFaildLog(notes)
		}
	}

	err := dao.ConsultationConfigDaoInst.UpdateConsultationConfig(tx, req.Id, req.Name, req.Text, req.LogoImgNo, req.UseStatus)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = err
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeleteConsultationConfig(ctx context.Context, req *custProto.DeleteConsultationConfigRequest, reply *custProto.DeleteConsultationConfigReply) error {
	if req.Id == "" {
		ss_log.Error("id is NULL")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)
	id := req.Id

	//获取当前要删除记录的idx,语言
	startIdx, lang, errGet1 := dao.ConsultationConfigDaoInst.GetIdxById(tx, id)
	if errGet1 != ss_err.ERR_SUCCESS {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	endIdx, errGet1 := dao.ConsultationConfigDaoInst.GetMaxIdx(lang)
	if errGet1 != ss_err.ERR_SUCCESS {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//将其要删除的排序序号后面的元素往前移
	for i := startIdx + 1; i <= endIdx; i++ {
		errUp := dao.ConsultationConfigDaoInst.ReplaceIdx(tx, i, lang)
		if errUp != ss_err.ERR_SUCCESS {
			ss_log.Error("errUp=[%v],i=[%v]", errUp, i)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	logoURL, err2 := dao.ConsultationConfigDaoInst.GetLogoURL(req.Id)
	if err2 != nil {
		ss_log.Error("InsertOrUpdateConsultationConfig 失败,err: %s", err2.Error())
	}

	if errStr := dao.ConsultationConfigDaoInst.DeleteConsultationConfig(tx, req.Id); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	description := fmt.Sprintf("删除服务商渠道id:[%v]", req.Id)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Servicer)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}
	if logoURL != "" {
		if _, err := common.UploadS3.DeleteOne(logoURL); err != nil {
			notes := fmt.Sprintf("DeleteConsultationConfig s3上删除图片失败,图片路劲为: %s,err: %s", logoURL, err.Error())
			ss_log.Error(notes)
			dao.DictimagesDaoInst.AddDelFaildLog(notes)
		}

		if err := dao.DictimagesDaoInst.Delete(logoURL); err != nil {
			notes := fmt.Sprintf("DeleteConsultationConfig 删除图片记录失败,图片路劲为: %s,err: %s", logoURL, err.Error())
			ss_log.Error(notes)
			dao.DictimagesDaoInst.AddDelFaildLog(notes)
		}
	}
	ss_sql.Commit(tx)

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) ModifyConsultationConfigStatus(ctx context.Context, req *custProto.ModifyConsultationConfigStatusRequest, reply *custProto.ModifyConsultationConfigStatusReply) error {
	if req.Id == "" {
		ss_log.Error("参数Id为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	str1, legal1 := util.GetParamZhCn(req.UseStatus, util.UseStatus)
	if !legal1 {
		ss_log.Error("UseStatus %v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if errStr := dao.ConsultationConfigDaoInst.ModifyConsultationConfigStatus(req.Id, req.UseStatus); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	//todo 添加关键操作日志

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetAgreements(ctx context.Context, req *custProto.GetAgreementsRequest, reply *custProto.GetAgreementsReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
		{Key: "type", Val: req.Type, EqType: "="},
		{Key: "lang", Val: req.Lang, EqType: "="},
	})

	total := dao.AgreementDaoInst.GetCnt(dbHandler, whereModel.WhereStr, whereModel.Args)

	//添加排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by use_status desc `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.AgreementDaoInst.GetAgreements(dbHandler, whereModel.WhereStr, whereModel.Args)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//这是管理系统后台获取协议（单个）
func (*CustHandler) GetAgreementDetail(ctx context.Context, req *custProto.GetAgreementDetailRequest, reply *custProto.GetAgreementDetailReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
		{Key: "id", Val: req.Id, EqType: "="},
	})

	data, err := dao.AgreementDaoInst.GetAgreementDetail(dbHandler, whereModel.WhereStr, whereModel.Args)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//这是app获取协议（单个）
func (*CustHandler) GetAgreementAppDetail(ctx context.Context, req *custProto.GetAgreementAppDetailRequest, reply *custProto.GetAgreementAppDetailReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
		{Key: "use_status", Val: "1", EqType: "="},
		{Key: "type", Val: req.Type, EqType: "="}, //0用户协议 1隐私协议 2实名认证协议
		{Key: "lang", Val: req.Lang, EqType: "="},
	})

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " limit 1 ")

	data, err := dao.AgreementDaoInst.GetAgreementDetail(dbHandler, whereModel.WhereStr, whereModel.Args)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	data.Id = ""
	data.CreateTime = ""
	data.ModifyTime = ""
	data.UseStatus = ""

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) InsertOrUpdateAgreement(ctx context.Context, req *custProto.InsertOrUpdateAgreementRequest, reply *custProto.InsertOrUpdateAgreementReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	//参数校验
	langStr, legal1 := util.GetParamZhCn(req.Lang, util.Lang)
	if !legal1 {
		ss_log.Error("Lang %v", langStr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	typeStr, legal2 := util.GetParamZhCn(req.Type, util.AgreementType)
	if !legal2 {
		ss_log.Error("Type %v", typeStr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	statusStr, legal3 := util.GetParamZhCn(req.UseStatus, util.UseStatus)
	if !legal3 {
		ss_log.Error("UseStatus %v", statusStr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//处理描述
	description := ""
	if req.Id == "" {
		id, err := dao.AgreementDaoInst.AddAgreement(tx, req.Text, req.Lang, req.Type, req.UseStatus)
		if err != ss_err.ERR_SUCCESS {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = err
			return nil
		}

		description = fmt.Sprintf("添加id[%v]的[%v][%v],内容:[%v],使用状态:[%v]", id, langStr, typeStr, req.Text, statusStr)
	} else {
		err := dao.AgreementDaoInst.UpdateAgreement(tx, req.Id, req.Text, req.Lang, req.Type, req.UseStatus)
		if err != ss_err.ERR_SUCCESS {
			ss_log.Error("更新失败，err=[%v]", err)
			reply.ResultCode = err
			return nil
		}

		//查询旧数据,为插入Web关键操作日志 做修改内容描述做准备.
		whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
			{Key: "is_delete", Val: "0", EqType: "="},
			{Key: "id", Val: req.Id, EqType: "="},
		})
		oldData, getErr := dao.AgreementDaoInst.GetAgreementDetail(dbHandler, whereModel.WhereStr, whereModel.Args)
		if getErr != ss_err.ERR_SUCCESS {
			ss_log.Error("获取旧数据失败,err=[%v]", getErr)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}

		if oldData.Type != req.Type || oldData.Text != req.Text || oldData.Lang != req.Lang || oldData.UseStatus != req.UseStatus {
			oldData.Type, _ = util.GetParamZhCn(oldData.Type, util.AgreementType)
			oldData.Lang, _ = util.GetParamZhCn(oldData.Lang, util.Lang)
			oldData.UseStatus, _ = util.GetParamZhCn(oldData.UseStatus, util.UseStatus)

			description = fmt.Sprintf("原有id[%v]的[%v][%v],内容:[%v],状态:[%v] ", oldData.Id, oldData.Lang, oldData.Type, oldData.Text, oldData.UseStatus)
			description = fmt.Sprintf("%v 修改为 [%v][%v],内容:[%v],状态:[%v]", description, langStr, typeStr, req.Text, statusStr)
		}

	}

	if description != "" {
		errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Account)
		if errAddLog != nil {
			ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeleteAgreement(ctx context.Context, req *custProto.DeleteAgreementRequest, reply *custProto.DeleteAgreementReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	//查询旧数据,为插入Web关键操作日志,内容描述做准备.
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
		{Key: "id", Val: req.Id, EqType: "="},
	})
	oldData, getErr := dao.AgreementDaoInst.GetAgreementDetail(dbHandler, whereModel.WhereStr, whereModel.Args)
	if getErr != ss_err.ERR_SUCCESS {
		ss_log.Error("获取旧数据失败,err=[%v]", getErr)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	oldData.Type, _ = util.GetParamZhCn(oldData.Type, util.AgreementType)
	oldData.Lang, _ = util.GetParamZhCn(oldData.Lang, util.Lang)

	description := fmt.Sprintf("删除id[%v]的[%v][%v],内容:[%v]", oldData.Id, oldData.Lang, oldData.Type, oldData.Text)

	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Account)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	err := dao.AgreementDaoInst.DeleteAgreement(tx, req.Id)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = err
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) ModifyAgreementStatus(ctx context.Context, req *custProto.ModifyAgreementStatusRequest, reply *custProto.ModifyAgreementStatusReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	//
	//查询要修改状态的语言和协议类型
	var lang, typeT sql.NullString
	sqlStr1 := "select lang, type from agreement where id = $1 and is_delete = '0' "
	err1 := ss_sql.QueryRowTx(tx, sqlStr1, []*sql.NullString{&lang, &typeT}, req.Id)
	if err1 != nil {
		ss_log.Error("查询旧数据失败,err=[%v]", err1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if errStr := dao.AgreementDaoInst.ModifyAgreementStatus(tx, req.Id, lang.String, typeT.String, req.UseStatus); errStr != ss_err.ERR_SUCCESS {
		ss_log.Error("修改数据时失败,err=[%v]", errStr)
		reply.ResultCode = errStr
		return nil
	}
	langStr, _ := util.GetParamZhCn(lang.String, util.Lang)
	agreementTypeStr, _ := util.GetParamZhCn(typeT.String, util.AgreementType)

	description := fmt.Sprintf("将id[%v]的[%v][%v]使用状态设置成启用", req.Id, langStr, agreementTypeStr)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Account)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//根据前端的筛选条件拼成字符串
func getExchangeQueryStr(req *custProto.GetExchangeOrderListRequest) string {

	queryStr := "" //筛选条件
	if req.LogNo != "" {
		queryStr = fmt.Sprintf("%v 兑换订单流水号：[%v]", queryStr, req.LogNo)
	}

	if req.StartTime != "" {
		queryStr = fmt.Sprintf("%v 开始时间：[%v]", queryStr, req.StartTime)
	}

	if req.EndTime != "" {
		queryStr = fmt.Sprintf("%v 结束时间：[%v]", queryStr, req.EndTime)
	}

	if req.OrderStatus != "" {
		str1, legal1 := util.GetParamZhCn(req.OrderStatus, util.OrderStatus)
		if !legal1 { //不合法
			ss_log.Error("%v", str1)
			queryStr = fmt.Sprintf("%v 订单状态：[%v]", queryStr, req.OrderStatus)
		} else {
			queryStr = fmt.Sprintf("%v 订单状态：[%v]", queryStr, str1)
		}
	}

	if req.Phone != "" {
		queryStr = fmt.Sprintf("%v 发起兑换的手机号：[%v]", queryStr, req.Phone)
	}

	if req.ExchangeType != "" {
		str1, legal1 := util.GetParamZhCn(req.ExchangeType, util.ExchangeType)
		if !legal1 { //不合法
			ss_log.Error("%v", str1)
			queryStr = fmt.Sprintf("%v 兑换类型：[%v]", queryStr, req.ExchangeType)
		} else {
			queryStr = fmt.Sprintf("%v 兑换类型：[%v]", queryStr, str1)
		}
	}

	return queryStr
}
func getIncomeQueryStr(req *custProto.GetIncomeOrderListRequest) string {

	queryStr := "" //筛选条件
	if req.LogNo != "" {
		queryStr = fmt.Sprintf("%v 存款订单流水号：[%v]", queryStr, req.LogNo)
	}
	if req.StartTime != "" {
		queryStr = fmt.Sprintf("%v 开始时间：[%v]", queryStr, req.StartTime)
	}
	if req.EndTime != "" {
		queryStr = fmt.Sprintf("%v 结束时间：[%v]", queryStr, req.EndTime)
	}
	if req.OrderStatus != "" {
		str1, legal1 := util.GetParamZhCn(req.OrderStatus, util.OrderStatus)
		if !legal1 { //不合法
			ss_log.Error("%v", str1)
			queryStr = fmt.Sprintf("%v 订单状态：[%v]", queryStr, req.OrderStatus)
		} else {
			queryStr = fmt.Sprintf("%v 订单状态：[%v]", queryStr, str1)
		}
	}
	if req.BalanceType != "" {
		queryStr = fmt.Sprintf("%v 币种：[%v]", queryStr, req.BalanceType)
	}
	if req.IncomePhone != "" {
		queryStr = fmt.Sprintf("%v 存款人的手机号：[%v]", queryStr, req.IncomePhone)
	}
	if req.Account != "" {
		queryStr = fmt.Sprintf("%v 收款的服务商账号：[%v]", queryStr, req.Account)
	}
	if req.RecvPhone != "" {
		queryStr = fmt.Sprintf("%v 收款人手机号：[%v]", queryStr, req.RecvPhone)
	}

	return queryStr
}
func getOutgoQueryStr(req *custProto.GetOutgoOrderListRequest) string {

	queryStr := "" //筛选条件
	if req.LogNo != "" {
		queryStr = fmt.Sprintf("%v 取款订单流水号：[%v]", queryStr, req.LogNo)
	}
	if req.StartTime != "" {
		queryStr = fmt.Sprintf("%v 开始时间：[%v]", queryStr, req.StartTime)
	}
	if req.EndTime != "" {
		queryStr = fmt.Sprintf("%v 结束时间：[%v]", queryStr, req.EndTime)
	}
	if req.OrderStatus != "" {
		str1, legal1 := util.GetParamZhCn(req.OrderStatus, util.OrderStatus)
		if !legal1 {
			ss_log.Error("%v", str1)
			queryStr = fmt.Sprintf("%v 订单状态：[%v]", queryStr, req.OrderStatus)
		} else {
			queryStr = fmt.Sprintf("%v 订单状态：[%v]", queryStr, str1)
		}

	}
	if req.BalanceType != "" {
		queryStr = fmt.Sprintf("%v 币种：[%v]", queryStr, req.BalanceType)
	}
	if req.Phone != "" {
		queryStr = fmt.Sprintf("%v 取款的手机号：[%v]", queryStr, req.Phone)
	}
	if req.Account != "" {
		queryStr = fmt.Sprintf("%v 服务商账号：[%v]", queryStr, req.Account)
	}

	return queryStr
}
func getTransferQueryStr(req *custProto.GetTransferOrderListRequest) string {

	queryStr := "" //筛选条件
	if req.LogNo != "" {
		queryStr = fmt.Sprintf("%v 转账订单流水号：[%v]", queryStr, req.LogNo)
	}
	if req.StartTime != "" {
		queryStr = fmt.Sprintf("%v 开始时间：[%v]", queryStr, req.StartTime)
	}
	if req.EndTime != "" {
		queryStr = fmt.Sprintf("%v 结束时间：[%v]", queryStr, req.EndTime)
	}
	if req.OrderStatus != "" {

		str1, legal1 := util.GetParamZhCn(req.OrderStatus, util.OrderStatus)
		if !legal1 {
			ss_log.Error("%v", str1)
			queryStr = fmt.Sprintf("%v 订单状态：[%v]", queryStr, req.OrderStatus)
		} else {
			queryStr = fmt.Sprintf("%v 订单状态：[%v]", queryStr, str1)
		}

	}
	if req.BalanceType != "" {
		queryStr = fmt.Sprintf("%v 币种：[%v]", queryStr, req.BalanceType)
	}
	if req.FromAccount != "" {
		queryStr = fmt.Sprintf("%v 转账发起人账号：[%v]", queryStr, req.FromAccount)
	}
	if req.ToAccount != "" {
		queryStr = fmt.Sprintf("%v 转账收款人账号：[%v]", queryStr, req.ToAccount)
	}

	return queryStr
}
func getCollectionQueryStr(req *custProto.GetCollectionOrdersRequest) string {
	queryStr := "" //筛选条件
	if req.LogNo != "" {
		queryStr = fmt.Sprintf("%v 收款订单流水号：[%v]", queryStr, req.LogNo)
	}
	if req.StartTime != "" {
		queryStr = fmt.Sprintf("%v 开始时间：[%v]", queryStr, req.StartTime)
	}
	if req.EndTime != "" {
		queryStr = fmt.Sprintf("%v 结束时间：[%v]", queryStr, req.EndTime)
	}
	if req.OrderStatus != "" {
		str1, legal1 := util.GetParamZhCn(req.OrderStatus, util.OrderStatus)
		if !legal1 {
			ss_log.Error("%v", str1)
			queryStr = fmt.Sprintf("%v 订单状态：[%v]", queryStr, req.OrderStatus)
		} else {
			queryStr = fmt.Sprintf("%v 订单状态：[%v]", queryStr, str1)
		}
	}
	if req.BalanceType != "" {
		queryStr = fmt.Sprintf("%v 币种：[%v]", queryStr, req.BalanceType)
	}
	if req.FromAccount != "" {
		queryStr = fmt.Sprintf("%v 转账发起人账号：[%v]", queryStr, req.FromAccount)
	}
	if req.ToAccount != "" {
		queryStr = fmt.Sprintf("%v 转账收款人账号：[%v]", queryStr, req.ToAccount)
	}
	return queryStr
}
func getSerToHeadQueryStr(req *custProto.GetToHeadquartersListRequest) string {
	queryStr := "" //筛选条件
	if req.LogNo != "" {
		queryStr = fmt.Sprintf("%v 服务商充值订单流水号：[%v]", queryStr, req.LogNo)
	}
	if req.StartTime != "" {
		queryStr = fmt.Sprintf("%v 开始时间：[%v]", queryStr, req.StartTime)
	}
	if req.EndTime != "" {
		queryStr = fmt.Sprintf("%v 结束时间：[%v]", queryStr, req.EndTime)
	}
	if req.OrderStatus != "" {
		str1, legal1 := util.GetParamZhCn(req.OrderStatus, util.AuditOrderStatus)
		if !legal1 {
			ss_log.Error("%v", str1)
			queryStr = fmt.Sprintf("%v 订单状态：[%v]", queryStr, req.OrderStatus)
		} else {
			queryStr = fmt.Sprintf("%v 订单状态：[%v]", queryStr, str1)
		}
	}
	if req.Account != "" {
		queryStr = fmt.Sprintf("%v 服务商账号：[%v]", queryStr, req.Account)
	}
	if req.OrderType != "" {
		switch req.OrderType {
		case "1":
			queryStr = fmt.Sprintf("%v 订单类型：[%v]", queryStr, "交易转账")
		case "2":
			queryStr = fmt.Sprintf("%v 订单类型：[%v]", queryStr, "结算转账")
		default:
			queryStr = fmt.Sprintf("%v 订单类型：[%v]", queryStr, "未知订单类型")
		}
	}
	if req.CurrencyType != "" {
		queryStr = fmt.Sprintf("%v 币种：[%v]", queryStr, req.CurrencyType)
	}
	return queryStr
}
func getToSerQueryStr(req *custProto.GetToServicerListRequest) string {
	queryStr := "" //筛选条件
	if req.LogNo != "" {
		queryStr = fmt.Sprintf("%v 服务商请款订单流水号：[%v]", queryStr, req.LogNo)
	}
	if req.StartTime != "" {
		queryStr = fmt.Sprintf("%v 开始时间：[%v]", queryStr, req.StartTime)
	}
	if req.EndTime != "" {
		queryStr = fmt.Sprintf("%v 结束时间：[%v]", queryStr, req.EndTime)
	}
	if req.OrderStatus != "" {
		str1, legal1 := util.GetParamZhCn(req.OrderStatus, util.AuditOrderStatus)
		if !legal1 {
			ss_log.Error("%v", str1)
			queryStr = fmt.Sprintf("%v 订单状态：[%v]", queryStr, req.OrderStatus)
		} else {
			queryStr = fmt.Sprintf("%v 订单状态：[%v]", queryStr, str1)
		}
	}
	if req.Account != "" {
		queryStr = fmt.Sprintf("%v 服务商账号：[%v]", queryStr, req.Account)
	}
	if req.OrderType != "" {
		switch req.OrderType {
		case "1":
			queryStr = fmt.Sprintf("%v 订单类型：[%v]", queryStr, "交易转账")
		case "2":
			queryStr = fmt.Sprintf("%v 订单类型：[%v]", queryStr, "结算转账")
		default:
			queryStr = fmt.Sprintf("%v 订单类型：[%v]", queryStr, "未知订单类型")
		}
	}
	if req.CurrencyType != "" {
		queryStr = fmt.Sprintf("%v 币种：[%v]", queryStr, req.CurrencyType)
	}
	if req.Nickname != "" {
		queryStr = fmt.Sprintf("%v 服务商昵称：[%v]", queryStr, req.Nickname)
	}
	return queryStr
}
func getLogVaccQueryStr(req *custProto.GetLogVaccountsRequest) string {
	queryStr := "" //筛选条件
	if req.LogNo != "" {
		queryStr = fmt.Sprintf("%v 虚拟账户日志流水号：[%v]", queryStr, req.LogNo)
	}
	if req.StartTime != "" {
		queryStr = fmt.Sprintf("%v 开始时间：[%v]", queryStr, req.StartTime)
	}
	if req.EndTime != "" {
		queryStr = fmt.Sprintf("%v 结束时间：[%v]", queryStr, req.EndTime)
	}
	if req.BalanceType != "" {
		queryStr = fmt.Sprintf("%v 币种：[%v]", queryStr, req.BalanceType)
	}
	if req.BizLogNo != "" {
		queryStr = fmt.Sprintf("%v 业务流水号：[%v]", queryStr, req.BizLogNo)
	}
	if req.OpType != "" {

		switch req.OpType {
		case constants.VaOpType_Add:
			queryStr = fmt.Sprintf("%v 操作类型：[%v]", queryStr, "+")
		case constants.VaOpType_Minus:
			queryStr = fmt.Sprintf("%v 操作类型：[%v]", queryStr, "-")
		case constants.VaOpType_Freeze:
			queryStr = fmt.Sprintf("%v 操作类型：[%v]", queryStr, "冻结")
		case constants.VaOpType_Defreeze:
			queryStr = fmt.Sprintf("%v 操作类型：[%v]", queryStr, "解冻")
		case constants.VaOpType_Defreeze_Minus:
			queryStr = fmt.Sprintf("%v 操作类型：[%v]", queryStr, "解冻并扣减")
		case constants.VaOpType_Defreeze_Add:
			queryStr = fmt.Sprintf("%v 操作类型：[%v]", queryStr, "解冻并增加")
		case constants.VaOpType_Defreeze_But_Minus:
			queryStr = fmt.Sprintf("%v 操作类型：[%v]", queryStr, "解冻不减")
		default:
			ss_log.Error("未知操作类型：[%v]", req.OpType)
		}

	}
	if req.Reason != "" {
		switch req.Reason {
		case constants.VaReason_Exchange:
			queryStr = fmt.Sprintf("%v 原因：[%v]", queryStr, "兑换")
		case constants.VaReason_INCOME:
			queryStr = fmt.Sprintf("%v 原因：[%v]", queryStr, "存款")
		case constants.VaReason_OUTGO:
			queryStr = fmt.Sprintf("%v 原因：[%v]", queryStr, "取款")
		case constants.VaReason_TRANSFER:
			queryStr = fmt.Sprintf("%v 原因：[%v]", queryStr, "转账")
		case constants.VaReason_COLLECTION:
			queryStr = fmt.Sprintf("%v 原因：[%v]", queryStr, "收款")
		case constants.VaReason_FEES:
			queryStr = fmt.Sprintf("%v 原因：[%v]", queryStr, "手续费")
		case constants.VaReason_Cancel_withdraw:
			queryStr = fmt.Sprintf("%v 原因：[%v]", queryStr, "pos端取消提现")
		case constants.VaReason_PROFIT_OUTGO:
			queryStr = fmt.Sprintf("%v 原因：[%v]", queryStr, "平台盈利提现")
		default:
			ss_log.Error("未知原因：[%v]", req.Reason)
		}
	}

	return queryStr
}

//处理数据
func processExchangeData(data *custProto.ExchangeOrderData) {
	str1, legal1 := util.GetParamZhCn(data.OrderStatus, util.OrderStatus)
	if !legal1 { //不合法
		ss_log.Error("%v", str1)
	} else {
		data.OrderStatus = str1
	}

	if data.InType == "usd" {
		data.Amount = "$" + strext.ToStringNoPoint(strext.ToFloat64(data.Amount)/100.0)
		data.Fees = "$" + strext.ToStringNoPoint(strext.ToFloat64(data.Fees)/100.0)
		data.Rate = "1USD = " + data.Rate + "KHR"
	} else if data.InType == "khr" {
		data.Amount = "៛" + data.Amount
		data.Fees = "៛" + data.Fees
		data.Rate = data.Rate + "KHR = 1USD"
	}
	if data.OutType == "usd" {
		data.TransAmount = "$" + strext.ToStringNoPoint(strext.ToFloat64(data.TransAmount)/100.0)
	} else if data.OutType == "khr" {
		data.TransAmount = "៛" + data.TransAmount
	}
	if data.CreateTime != "" {
		data.CreateTime = data.CreateTime[:10] + " " + data.CreateTime[11:19]
	}
	if data.FinishTime != "" {
		data.FinishTime = data.FinishTime[:10] + " " + data.FinishTime[11:19]
	}
}
func processIncomeData(data *custProto.IncomeOrderData) {
	if data.OrderStatus != "" {
		str1, legal1 := util.GetParamZhCn(data.OrderStatus, util.OrderStatus)
		if !legal1 { //不合法
			ss_log.Error("%v", str1)
		} else {
			data.OrderStatus = str1
		}
	}
	if data.CreateTime != "" {
		data.CreateTime = data.CreateTime[:10] + " " + data.CreateTime[11:19]
	}
	if data.FinishTime != "" {
		data.FinishTime = data.FinishTime[:10] + " " + data.FinishTime[11:19]
	}
	if data.BalanceType == "usd" {
		data.Amount = strext.ToStringNoPoint(strext.ToFloat64(data.Amount) / 100.0)
		data.Fees = strext.ToStringNoPoint(strext.ToFloat64(data.Fees) / 100.0)
	}
}
func processOutgoData(data *custProto.OutgoOrderData) {
	if data.OrderStatus != "" {
		str1, legal1 := util.GetParamZhCn(data.OrderStatus, util.OrderStatus)
		if !legal1 {
			ss_log.Error("%v", str1)
		} else {
			data.OrderStatus = str1
		}
	}
	if data.CreateTime != "" {
		data.CreateTime = data.CreateTime[:10] + " " + data.CreateTime[11:19]
	}
	if data.FinishTime != "" {
		data.FinishTime = data.FinishTime[:10] + " " + data.FinishTime[11:19]
	}
	if data.BalanceType == "usd" {
		data.Amount = strext.ToStringNoPoint(strext.ToFloat64(data.Amount) / 100.0)
		data.Fees = strext.ToStringNoPoint(strext.ToFloat64(data.Fees) / 100.0)
	}
}
func processTransferData(data *custProto.TransferOrderData) {
	if data.OrderStatus != "" {
		str1, legal1 := util.GetParamZhCn(data.OrderStatus, util.OrderStatus)
		if !legal1 {
			ss_log.Error("%v", str1)
		} else {
			data.OrderStatus = str1
		}
	}
	if data.CreateTime != "" {
		data.CreateTime = data.CreateTime[:10] + " " + data.CreateTime[11:19]
	}
	if data.FinishTime != "" {
		data.FinishTime = data.FinishTime[:10] + " " + data.FinishTime[11:19]
	}
	if data.BalanceType == "usd" {
		data.Amount = strext.ToStringNoPoint(strext.ToFloat64(data.Amount) / 100.0)
		data.Fees = strext.ToStringNoPoint(strext.ToFloat64(data.Fees) / 100.0)
	}
}
func processCollectionData(data *custProto.CollectionOrderData) {
	if data.OrderStatus != "" {
		str1, legal1 := util.GetParamZhCn(data.OrderStatus, util.OrderStatus)
		if !legal1 {
			ss_log.Error("%v", str1)
		} else {
			data.OrderStatus = str1
		}
	}
	if data.CreateTime != "" {
		data.CreateTime = data.CreateTime[:10] + " " + data.CreateTime[11:19]
	}
	if data.FinishTime != "" {
		data.FinishTime = data.FinishTime[:10] + " " + data.FinishTime[11:19]
	}
	if data.BalanceType == "usd" {
		data.Amount = strext.ToStringNoPoint(strext.ToFloat64(data.Amount) / 100.0)
		data.Fees = strext.ToStringNoPoint(strext.ToFloat64(data.Fees) / 100.0)
	}
}
func processSerToHeadData(data *custProto.ToHeadquartersData) {
	str1, legal1 := util.GetParamZhCn(data.OrderStatus, util.AuditOrderStatus)
	if !legal1 {
		ss_log.Error("%v", str1)
	} else {
		data.OrderStatus = str1
	}

	switch data.OrderType {
	case "1":
		data.OrderType = "交易转账"
	case "2":
		data.OrderType = "结算转账"
	default:
		data.OrderType = "未知订单类型"
	}
	str2, legal2 := util.GetParamZhCn(data.CollectionType, util.CollectionType)
	if !legal2 {
		ss_log.Error("%v", str2)
	} else {
		data.CollectionType = str2
	}
	if data.CreateTime != "" {
		data.CreateTime = data.CreateTime[:10] + " " + data.CreateTime[11:19]
	}
	if data.CurrencyType == "usd" {
		data.Amount = strext.ToStringNoPoint(strext.ToFloat64(data.Amount) / 100.0)
	}
}
func processToServicerData(data *custProto.ToServicerData) {
	str1, legal1 := util.GetParamZhCn(data.OrderStatus, util.AuditOrderStatus)
	if !legal1 {
		ss_log.Error("%v", str1)
	} else {
		data.OrderStatus = str1
	}
	switch data.OrderType {
	case "1":
		data.OrderType = "交易转账"
	case "2":
		data.OrderType = "结算转账"
	default:
		data.OrderType = "未知订单类型"
	}
	str2, legal2 := util.GetParamZhCn(data.CollectionType, util.CollectionType)
	if !legal2 {
		ss_log.Error("%v", str2)
	} else {
		data.CollectionType = str2
	}
	if data.CreateTime != "" {
		data.CreateTime = data.CreateTime[:10] + " " + data.CreateTime[11:19]
	}
	if data.CurrencyType == "usd" {
		data.Amount = strext.ToStringNoPoint(strext.ToFloat64(data.Amount) / 100.0)
	}
}
func processLogVaccountData(data *custProto.LogVaccountData) {
	switch data.OpType {
	case constants.VaOpType_Add:
		data.OpType = "+"
	case constants.VaOpType_Minus:
		data.OpType = "-"
	case constants.VaOpType_Freeze:
		data.OpType = "冻结"
	case constants.VaOpType_Defreeze:
		data.OpType = "解冻"
	case constants.VaOpType_Defreeze_Minus:
		data.OpType = "解冻并扣减"
	case constants.VaOpType_Defreeze_Add:
		data.OpType = "解冻并增加"
	case constants.VaOpType_Defreeze_But_Minus:
		data.OpType = "解冻不减"
	default:
		ss_log.Error("未知操作类型：[%v]", data.OpType)
		data.OpType = "未知操作类型"
	}
	switch data.Reason {
	case constants.VaReason_Exchange:
		data.Reason = "兑换"
	case constants.VaReason_INCOME:
		data.Reason = "存款"
	case constants.VaReason_OUTGO:
		data.Reason = "取款"
	case constants.VaReason_TRANSFER:
		data.Reason = "转账"
	case constants.VaReason_COLLECTION:
		data.Reason = "收款"
	case constants.VaReason_FEES:
		data.Reason = "手续费"
	case constants.VaReason_Cancel_withdraw:
		data.Reason = "pos端取消提现"
	case constants.VaReason_PROFIT_OUTGO:
		data.Reason = "平台盈利提现"
	default:
		data.Reason = "未知原因"
		ss_log.Error("未知原因：[%v]", data.Reason)
	}
	if data.CreateTime != "" {
		data.CreateTime = data.CreateTime[:10] + " " + data.CreateTime[11:19]
	}
	if data.BalanceType == "usd" {
		data.Amount = strext.ToStringNoPoint(strext.ToFloat64(data.Amount) / 100.0)
		data.FrozenBalance = strext.ToStringNoPoint(strext.ToFloat64(data.FrozenBalance) / 100.0)
		data.Balance = strext.ToStringNoPoint(strext.ToFloat64(data.Balance) / 100.0)
	}
}

func (custHandler *CustHandler) GetCreateXlsxFileContent(ctx context.Context, req *custProto.GetCreateXlsxFileContentRequest, reply *custProto.GetCreateXlsxFileContentReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	switch req.BillFileType {
	case constants.XLSX_FILE_TYPE_EXCHANGE:
		queryStr := getExchangeQueryStr(req.ExchangeReq)

		//查询结果
		var datas []*custProto.ExchangeOrderData
		for {
			replyGet := &custProto.GetExchangeOrderListReply{}
			custHandler.GetExchangeOrderList(ctx, req.ExchangeReq, replyGet)
			if replyGet.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取兑换订单流水失败------------>err=[%v]", replyGet.ResultCode)
				reply.ResultCode = replyGet.ResultCode
				return nil
			}
			if replyGet.DataList == nil {
				break
			} else {
				req.ExchangeReq.Page = req.ExchangeReq.Page + 1
				for _, data := range replyGet.DataList {
					processExchangeData(data)

					datas = append(datas, data)
				}
			}

		}

		if datas == nil { //当其仍然是空时，将不允许创建。直接返回
			ss_log.Error("即将生成的xlsx文件无内容,已阻止生成文件。")
			reply.ResultCode = ss_err.ERR_CreateFileDataNull
			return nil
		}

		//返回处理过的数据
		reply.QueryStr = queryStr
		reply.BillType = constants.XLSX_FILE_TYPE_EXCHANGE
		reply.Datas = &custProto.XlsxFileContentData{
			ExchangeDatas: datas,
		}
		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	case constants.XLSX_FILE_TYPE_INCOME:

		//筛选条件
		queryStr := getIncomeQueryStr(req.IncomeReq)

		var datas []*custProto.IncomeOrderData
		for {
			replyGet := &custProto.GetIncomeOrderListReply{}
			custHandler.GetIncomeOrderList(ctx, req.IncomeReq, replyGet)

			if replyGet.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("查询存款数据失败")
				reply.ResultCode = replyGet.ResultCode
				return nil
			}

			if replyGet.DataList == nil {
				break
			} else {
				req.IncomeReq.Page = req.IncomeReq.Page + 1
				for _, data := range replyGet.DataList {
					//处理数据
					processIncomeData(data)
					datas = append(datas, data)
				}
			}
		}

		if datas == nil { //当其仍然是空时，将不允许创建。直接返回
			ss_log.Error("即将生成的xlsx文件无内容,已阻止生成文件。")
			reply.ResultCode = ss_err.ERR_CreateFileDataNull
			return nil
		}

		reply.QueryStr = queryStr
		reply.BillType = constants.XLSX_FILE_TYPE_INCOME
		reply.Datas = &custProto.XlsxFileContentData{
			IncomeDatas: datas,
		}
		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	case constants.XLSX_FILE_TYPE_OUTGO:

		queryStr := getOutgoQueryStr(req.OutgoReq) //筛选条件

		var datas []*custProto.OutgoOrderData
		for {
			//查询结果
			replyGet := &custProto.GetOutgoOrderListReply{}
			custHandler.GetOutgoOrderList(ctx, req.OutgoReq, replyGet)

			if replyGet.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取取款数据失败")
				reply.ResultCode = replyGet.ResultCode
				return nil
			}

			if replyGet.DataList == nil {
				break
			} else {
				req.OutgoReq.Page = req.OutgoReq.Page + 1
				for _, data := range replyGet.DataList {
					//处理数据
					processOutgoData(data)

					datas = append(datas, data)
				}
			}

		}

		if datas == nil { //当其仍然是空时，将不允许创建。直接返回
			ss_log.Error("即将生成的xlsx文件无内容,已阻止生成文件。")
			reply.ResultCode = ss_err.ERR_CreateFileDataNull
			return nil
		}

		reply.QueryStr = queryStr
		reply.BillType = constants.XLSX_FILE_TYPE_OUTGO
		reply.Datas = &custProto.XlsxFileContentData{
			OutgoDatas: datas,
		}
		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	case constants.XLSX_FILE_TYPE_TRANSFER: //转
		queryStr := getTransferQueryStr(req.TransferReq) //筛选条件

		var datas []*custProto.TransferOrderData
		for {
			replyGet := &custProto.GetTransferOrderListReply{}
			custHandler.GetTransferOrderList(ctx, req.TransferReq, replyGet)

			if replyGet.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取转账数据失败")
				reply.ResultCode = replyGet.ResultCode
				return nil
			}

			if replyGet.DataList == nil {
				break
			} else {
				req.TransferReq.Page = req.TransferReq.Page + 1
				for _, data := range replyGet.DataList {
					//处理数据
					processTransferData(data)

					datas = append(datas, data)
				}
			}
		}

		if datas == nil { //当其仍然是空时，将不允许创建。直接返回
			ss_log.Error("即将生成的xlsx文件无内容,已阻止生成文件。")
			reply.ResultCode = ss_err.ERR_CreateFileDataNull
			return nil
		}

		//返回处理过的数据
		reply.QueryStr = queryStr
		reply.BillType = constants.XLSX_FILE_TYPE_TRANSFER
		reply.Datas = &custProto.XlsxFileContentData{
			TransferDatas: datas,
		}
		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	case constants.XLSX_FILE_TYPE_COLLECTION: //收
		queryStr := getCollectionQueryStr(req.CollectionReq) //筛选条件

		var datas []*custProto.CollectionOrderData
		for {
			//查询结果
			replyGet := &custProto.GetCollectionOrdersReply{}
			custHandler.GetCollectionOrders(ctx, req.CollectionReq, replyGet)

			if replyGet.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取收款订单数据失败")
				reply.ResultCode = replyGet.ResultCode
				return nil
			}

			if replyGet.DataList == nil {
				break
			} else {
				req.CollectionReq.Page = req.CollectionReq.Page + 1
				for _, data := range replyGet.DataList {
					//处理数据
					processCollectionData(data)
					datas = append(datas, data)
				}
			}
		}

		if datas == nil { //当其仍然是空时，将不允许创建。直接返回
			ss_log.Error("即将生成的xlsx文件无内容,已阻止生成文件。")
			reply.ResultCode = ss_err.ERR_CreateFileDataNull
			return nil
		}

		//返回处理过的数据
		reply.QueryStr = queryStr
		reply.BillType = constants.XLSX_FILE_TYPE_COLLECTION
		reply.Datas = &custProto.XlsxFileContentData{
			CollectionDatas: datas,
		}
		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	case constants.XLSX_FILE_TYPE_TO_HEADQUARTERS: //服务商充值
		queryStr := getSerToHeadQueryStr(req.ToHeadReq) //筛选条件

		var datas []*custProto.ToHeadquartersData
		for {
			//查询结果
			replyGet := &custProto.GetToHeadquartersListReply{}
			custHandler.GetToHeadquartersList(ctx, req.ToHeadReq, replyGet)

			if replyGet.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取服务商充值数据失败")
				reply.ResultCode = replyGet.ResultCode
				return nil
			}
			if replyGet.DataList == nil {
				break
			} else {
				req.ToHeadReq.Page = req.ToHeadReq.Page + 1
				for _, data := range replyGet.DataList {
					//处理数据
					processSerToHeadData(data)
					datas = append(datas, data)
				}
			}

		}

		if datas == nil { //当其仍然是空时，将不允许创建。直接返回
			ss_log.Error("即将生成的xlsx文件无内容,已阻止生成文件。")
			reply.ResultCode = ss_err.ERR_CreateFileDataNull
			return nil
		}

		//返回处理过的数据
		reply.QueryStr = queryStr
		reply.BillType = constants.XLSX_FILE_TYPE_TO_HEADQUARTERS
		reply.Datas = &custProto.XlsxFileContentData{
			ToHeadquartersDatas: datas,
		}
		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	case constants.XLSX_FILE_TYPE_TO_SERVICER: //服务商请款
		queryStr := getToSerQueryStr(req.ToSerReq) //筛选条件

		//查询结果
		var datas []*custProto.ToServicerData
		for {
			replyGet := &custProto.GetToServicerListReply{}
			custHandler.GetToServicerList(ctx, req.ToSerReq, replyGet)

			if replyGet.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取服务商请款数据失败")
				reply.ResultCode = replyGet.ResultCode
				return nil
			}
			if replyGet.DataList == nil {
				break
			} else {
				req.ToSerReq.Page = req.ToSerReq.Page + 1
				for _, data := range replyGet.DataList {
					processToServicerData(data)

					datas = append(datas, data)
				}
			}
		}

		if datas == nil { //当其仍然是空时，将不允许创建。直接返回
			ss_log.Error("即将生成的xlsx文件无内容,已阻止生成文件。")
			reply.ResultCode = ss_err.ERR_CreateFileDataNull
			return nil
		}

		//返回处理过的数据
		reply.QueryStr = queryStr
		reply.BillType = constants.XLSX_FILE_TYPE_TO_SERVICER
		reply.Datas = &custProto.XlsxFileContentData{
			ToServicerDatas: datas,
		}
		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	case constants.XLSX_FILE_TYPE_VACCOUNT_LOG: //虚拟账户日志流水
		queryStr := getLogVaccQueryStr(req.LogVaccReq) //筛选条件

		//查询结果
		var datas []*custProto.LogVaccountData
		for {
			replyGet := &custProto.GetLogVaccountsReply{}
			custHandler.GetLogVaccounts(ctx, req.LogVaccReq, replyGet)

			if replyGet.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取虚拟账户日志流水失败")
				reply.ResultCode = replyGet.ResultCode
				return nil
			}
			if replyGet.Datas == nil {
				break
			} else {
				req.LogVaccReq.Page = req.LogVaccReq.Page + 1
				for _, data := range replyGet.Datas {
					//处理数据
					processLogVaccountData(data)
					datas = append(datas, data)
				}
			}
		}

		if datas == nil { //当其仍然是空时，将不允许创建。直接返回
			ss_log.Error("即将生成的xlsx文件无内容,已阻止生成文件。")
			reply.ResultCode = ss_err.ERR_CreateFileDataNull
			return nil
		}

		//返回处理过的数据
		reply.BillType = constants.XLSX_FILE_TYPE_VACCOUNT_LOG
		reply.QueryStr = queryStr
		reply.Datas = &custProto.XlsxFileContentData{
			LogVaccountDatas: datas,
		}
		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	default:
		ss_log.Error("订单类型参数错误")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (c *CustHandler) PosFence(ctx context.Context, req *custProto.PosFenceRequest, reply *custProto.PosFenceReply) error {
	// 先去redis里面查找看有没有此数据
	redisKey := "pos_sn_" + req.PosSn
	var lat, lng, scope string
	var err error
	//if value, _ := ss_data.GetPosSnFromCache(redisKey, cache.RedisCli, constants.DefPoolName); value == "" { // 查询数据库并设置进redis
	if value, _ := cache.RedisClient.Get(redisKey).Result(); value == "" { // 查询数据库并设置进redis
		lat, lng, scope, err = GetPosLatLng(req.PosSn)
		if err != nil {
			ss_log.Error("%s", err.Error())
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		redisValue := fmt.Sprintf("%s,%s,%s", lat, lng, scope)
		// 设进redis
		//if err := ss_data.SetPosSnToCache(redisKey, redisValue, cache.RedisCli, constants.DefPoolName); err != nil {
		if err := cache.RedisClient.Set(redisKey, redisValue, constants.PosNoKeySecV2).Err(); err != nil {
			ss_log.Error("经纬度存进redis失败,posNo--->%s,lat--->%s,lng--->%s,scope--->%s", req.PosSn, lat, lng, scope)
		}
	} else {
		split := strings.Split(value, ",")
		lat = split[0]
		lng = split[1]
		scope = split[2]
	}
	// 计算距离
	distance := ss_count.CountCircleDistance(strext.ToFloat64(lat), strext.ToFloat64(lng), strext.ToFloat64(req.Lat), strext.ToFloat64(req.Lng))
	if distance > strext.ToFloat64(scope) {
		ss_log.Error("pos机超出使用范围,计算的范围是--->%v,限定的范围是--->%v", distance, scope)
		reply.ResultCode = ss_err.ERR_POS_OUT_OF_RANGE
		return nil
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func GetPosLatLng(posSn string) (lat, lng, scope string, err error) {
	// 根据pos_sn找服务商no
	serverNo := dao.ServicerTerminalDao.GetSerPosServicerNoByPosNo(posSn)
	if serverNo == "" {
		//ss_log.Error("%s", "根据posNo查找服务商 id 失败")
		return "", "", "", errors.New("根据posNo查找服务商 id 失败")
	}
	lat, lng, scope = dao.ServiceDaoInst.GetLatLngInfoFromNo(serverNo)
	if lat == "" || lng == "" || scope == "" {
		//ss_log.Error("根据 serverNo 查找服务商范围失败 lat--->%s,lng--->%s,scope--->%s", lat, lng, scope)
		return "", "", "", errors.New("根据 serverNo 查找服务商范围失败")
	}
	return lat, lng, scope, nil
}

func (*CustHandler) GetCashiers(ctx context.Context, req *custProto.GetCashiersRequest, reply *custProto.GetCashiersReply) error {

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ca.is_delete", Val: "0", EqType: "="},
		{Key: "ca.servicer_no", Val: req.ServicerNo, EqType: "="},
	})

	total := dao.CashierDaoInstance.GetCnt(whereModel.WhereStr, whereModel.Args)

	//添加排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by ca.create_time desc `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.CashierDaoInstance.GetCashiers(whereModel.WhereStr, whereModel.Args)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetCashierDetail(ctx context.Context, req *custProto.GetCashierDetailRequest, reply *custProto.GetCashierDetailReply) error {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ca.is_delete", Val: "0", EqType: "="},
		{Key: "ca.uid", Val: req.CashierNo, EqType: "="},
	})

	data, err := dao.CashierDaoInstance.GetCashierDetail(whereModel.WhereStr, whereModel.Args)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeleteCashier(ctx context.Context, req *custProto.DeleteCashierRequest, reply *custProto.DeleteCashierReply) error {

	//店员的账号
	account := dao.CashierDaoInstance.GetCashierAccountByCashierNo(req.CashierNo)
	//店员服务商的账号
	serAccount, serAccNo := dao.CashierDaoInstance.GetSrvAccountNoFromCashierNo(req.CashierNo)

	if req.ServicerAccNo != "" && req.ServicerAccNo != serAccNo {
		ss_log.Error("服务商不能删除他人的店员。")
		reply.ResultCode = ss_err.ERR_SYS_NO_API_AUTH
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	if errStr := dao.CashierDaoInstance.DeleteCashierTx(tx, req.CashierNo); errStr != ss_err.ERR_SUCCESS {
		ss_log.Error("删除店员[%v]失败", req.CashierNo)
		reply.ResultCode = errStr
		return nil
	}
	//删除店员与其账号的关系
	if err := dao.RelaAccIdenDaoInst.DeleteRelaAccIden(tx, req.CashierNo, constants.AccountType_POS); err != nil {
		ss_log.Error("删除店员[%v]与账号的关系失败，err=[%v]", req.CashierNo, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}

	if req.LoginUid != "" { //此参数只有管理平台为服务商添加店员时才会传。

		description := fmt.Sprintf("为服务商[%v]删除店员[%v]", serAccount, account)

		errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Servicer)
		if errAddLog != nil {
			ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) ModifyCashier(ctx context.Context, req *custProto.ModifyCashierRequest, reply *custProto.ModifyCashierReply) error {
	accNo := dao.RelaAccIdenDaoInst.GetAccNo(req.CashierNo, constants.AccountType_POS)
	if accNo == "" {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if errStr := dao.AccDaoInstance.UpdatePhone(accNo, req.Phone); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 交换位置
func (c *CustHandler) SwapConsultationIdx(ctx context.Context, req *custProto.SwapConsultationIdxRequest, reply *custProto.SwapConsultationIdxReply) error {
	switch req.SwapType {
	case constants.SwapType_Up: // up
	case constants.SwapType_Down: // down
	default:
		ss_log.Error("swapType错误")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.Idx == "1" && req.SwapType == constants.SwapType_Up {
		ss_log.Error("swapType错误")
		reply.ResultCode = ss_err.ERR_MERC_IS_UPTOP
		return nil
	}

	consultationNoFrom := req.Id

	consultationNoTo := dao.ConsultationConfigDaoInst.GetNearIdxConsultationNo(req.Idx, req.SwapType, req.Lang)
	if consultationNoTo == "" {
		ss_log.Error("获取funcNo失败")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	switch req.SwapType {
	case constants.SwapType_Up:
		// 向上需要交换一下
		x := consultationNoTo
		consultationNoTo = consultationNoFrom
		consultationNoFrom = x
	}

	errCode := dao.ConsultationConfigDaoInst.ExchangeIdx(consultationNoFrom, consultationNoTo)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", errCode)
		reply.ResultCode = errCode
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetClientInfos(ctx context.Context, req *custProto.GetClientInfosRequest, reply *custProto.GetClientInfosReply) error {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "cli.id", Val: req.Id, EqType: "="},
		{Key: "cli.platform", Val: req.Platform, EqType: "="},
		{Key: "cli.account", Val: req.Account, EqType: "like"},
		{Key: "cli.uuid", Val: req.Uuid, EqType: "like"},
	})
	datasT := []*custProto.ClientInfoData{}
	totalT := ""
	switch req.ClientType {
	case "cust":
		total, errCnt := dao.ClientInfoDaoInst.GetCustCnt(whereModel.WhereStr, whereModel.Args)
		if errCnt != nil {
			ss_log.Error("err=[%v]", errCnt)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}

		//添加排序
		ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by cli.create_time desc `)
		ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

		datas, err := dao.ClientInfoDaoInst.GetCustClientInfos(whereModel.WhereStr, whereModel.Args)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}

		totalT = total
		datasT = datas
	case "servicer":
		total, errCnt := dao.ClientInfoDaoInst.GetSerCnt(whereModel.WhereStr, whereModel.Args)
		if errCnt != nil {
			ss_log.Error("err=[%v]", errCnt)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}

		//添加排序
		ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by cli.create_time desc `)
		ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

		datas, err := dao.ClientInfoDaoInst.GetSerClientInfos(whereModel.WhereStr, whereModel.Args)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}

		totalT = total
		datasT = datas
	default:
		ss_log.Error("ClientType[%v]参数错误", req.ClientType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	reply.Datas = datasT
	reply.Total = strext.ToInt32(totalT)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetLogAccounts(ctx context.Context, req *custProto.GetLogAccountsRequest, reply *custProto.GetLogAccountsReply) error {
	if req.StartTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.StartTime) {
			ss_log.Error("GetLogAccounts StartTime格式不正确,StartTime: %s", req.StartTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}
	if req.EndTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.EndTime) {
			ss_log.Error("GetLogAccounts EndTime格式不正确,StartTime: %s", req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	var datasT []*custProto.LogAccountData
	totalT := ""
	switch req.LogType {
	case "web":
		whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
			{Key: "web.create_time", Val: req.StartTime, EqType: ">="},
			{Key: "web.create_time", Val: req.EndTime, EqType: "<="},
			{Key: "web.log_no", Val: req.LogNo, EqType: "like"},
			{Key: "web.type", Val: req.Type, EqType: "="},
			{Key: "acc.account", Val: req.Account, EqType: "like"},
		})
		total, errCnt := dao.LogDaoInstance.GetLogAccountWebCnt(whereModel.WhereStr, whereModel.Args)
		if errCnt != nil {
			ss_log.Error("err=[%v]", errCnt)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}

		//添加排序
		ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by web.create_time desc `)
		ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

		datas, err := dao.LogDaoInstance.GetLogAccountWebInfos(whereModel.WhereStr, whereModel.Args)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}

		totalT = total
		datasT = datas
	case "cli_app": //后续这里查app的操作日志
		whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
			{Key: "la.log_time", Val: req.StartTime, EqType: ">="},
			{Key: "la.log_time", Val: req.EndTime, EqType: "<="},
			{Key: "la.log_no", Val: req.LogNo, EqType: "like"},
			{Key: "la.type", Val: req.Type, EqType: "="},
			{Key: "acc.account", Val: req.Account, EqType: "like"},
		})
		total, errCnt := dao.LogDaoInstance.GetLogAccountCnt(whereModel.WhereStr, whereModel.Args)
		if errCnt != nil {
			ss_log.Error("err=[%v]", errCnt)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}

		//添加排序
		ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by la.log_time desc `)
		ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

		datas, err := dao.LogDaoInstance.GetLogAccountInfos(whereModel.WhereStr, whereModel.Args)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}

		totalT = total
		datasT = datas
	case "cli_pos": //这是预留的pos操作日志(如果后面需要的话)
	default:
		ss_log.Error("LogType[%v]参数错误", req.LogType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	reply.Datas = datasT
	reply.Total = strext.ToInt32(totalT)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetEvents(ctx context.Context, req *custProto.GetEventsRequest, reply *custProto.GetEventsReply) error {

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "create_time", Val: req.StartTime, EqType: ">="},
		{Key: "create_time", Val: req.EndTime, EqType: "<="},
		{Key: "is_delete", Val: "0", EqType: "="},
	})
	total, errCnt := dao.RiskDaoInstance.GetEventCnt(whereModel.WhereStr, whereModel.Args)
	if errCnt != nil {
		ss_log.Error("err=[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//添加排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by create_time desc `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.RiskDaoInstance.GetEventInfos(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetEvent(ctx context.Context, req *custProto.GetEventRequest, reply *custProto.GetEventReply) error {
	if req.EventNo == "" {
		ss_log.Error("EventNo为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "event_no", Val: req.EventNo, EqType: "="},
		{Key: "is_delete", Val: "0", EqType: "="},
	})

	data, err := dao.RiskDaoInstance.GetEventDetail(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) InsertOrUpdateEvent(ctx context.Context, req *custProto.InsertOrUpdateEventRequest, reply *custProto.InsertOrUpdateEventReply) error {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	if req.EventNo == "" {
		eventNo, err := dao.RiskDaoInstance.AddEvent(tx, req.EventName)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_ADD
			return nil
		}

		ss_log.Error("插入一条新数据 eventNo:[%s]", eventNo)
	} else {
		err := dao.RiskDaoInstance.ModifyEvent(tx, req.EventNo, req.EventName)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}

	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeleteEvent(ctx context.Context, req *custProto.DeleteEventRequest, reply *custProto.DeleteEventReply) error {
	if req.EventNo == "" {
		ss_log.Error("EventNo为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	err := dao.RiskDaoInstance.DelEvent(tx, req.EventNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetEvaParams(ctx context.Context, req *custProto.GetEvaParamsRequest, reply *custProto.GetEvaParamsReply) error {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		//{Key: "is_delete", Val: "0", EqType: "<="},
	})
	total, errCnt := dao.RiskDaoInstance.GetEvaParamsCnt(whereModel.WhereStr, whereModel.Args)
	if errCnt != nil {
		ss_log.Error("err=[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//添加排序
	//ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by create_time desc `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.RiskDaoInstance.GetEvaParamsInfos(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) InsertOrUpdateEvaParam(ctx context.Context, req *custProto.InsertOrUpdateEvaParamRequest, reply *custProto.InsertOrUpdateEvaParamReply) error {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)
	err := dao.RiskDaoInstance.AddOrUpdateEvaParam(tx, req.Key, req.Val)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeleteEvaParam(ctx context.Context, req *custProto.DeleteEvaParamRequest, reply *custProto.DeleteEvaParamReply) error {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	err := dao.RiskDaoInstance.DelEvaParam(tx, req.Key)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetGlobalParam(ctx context.Context, req *custProto.GetGlobalParamRequest, reply *custProto.GetGlobalParamReply) error {

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		//{Key: "is_delete", Val: "0", EqType: "<="},
	})
	total, errCnt := dao.RiskDaoInstance.GetGlobalParamCnt(whereModel.WhereStr, whereModel.Args)
	if errCnt != nil {
		ss_log.Error("err=[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//添加排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by param_key  `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.RiskDaoInstance.GetGlobalParamInfos(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) InsertOrUpdateGlobalParam(ctx context.Context, req *custProto.InsertOrUpdateGlobalParamRequest, reply *custProto.InsertOrUpdateGlobalParamReply) error {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)
	err := dao.RiskDaoInstance.AddOrUpdateGlobalParam(tx, req.ParamKey, req.ParamValue, req.Remark)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeleteGlobalParam(ctx context.Context, req *custProto.DeleteGlobalParamRequest, reply *custProto.DeleteGlobalParamReply) error {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	err := dao.RiskDaoInstance.DelGlobalParam(tx, req.ParamKey)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetLogResults(ctx context.Context, req *custProto.GetLogResultsRequest, reply *custProto.GetLogResultsReply) error {

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		//{Key: "is_delete", Val: "0", EqType: "<="},
	})
	total, errCnt := dao.RiskDaoInstance.GetLogResultsCnt(whereModel.WhereStr, whereModel.Args)
	if errCnt != nil {
		ss_log.Error("err=[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//添加排序
	//ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by create_time desc `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.RiskDaoInstance.GetLogResultsInfos(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetOps(ctx context.Context, req *custProto.GetOpsRequest, reply *custProto.GetOpsReply) error {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
	})
	total, errCnt := dao.RiskDaoInstance.GetOpsCnt(whereModel.WhereStr, whereModel.Args)
	if errCnt != nil {
		ss_log.Error("err=[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//添加排序
	//ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by create_time desc `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.RiskDaoInstance.GetOpsInfos(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetOp(ctx context.Context, req *custProto.GetOpRequest, reply *custProto.GetOpReply) error {
	if req.OpNo == "" {
		ss_log.Error("OpNo为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "op_no", Val: req.OpNo, EqType: "="},
		{Key: "is_delete", Val: "0", EqType: "="},
	})

	data, err := dao.RiskDaoInstance.GetOpDetail(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) InsertOrUpdateOp(ctx context.Context, req *custProto.InsertOrUpdateOpRequest, reply *custProto.InsertOrUpdateOpReply) error {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	if req.OpNo == "" {
		eventNo, err := dao.RiskDaoInstance.AddOp(tx, req.OpName, req.ScriptName, req.Param, req.Score)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_ADD
			return nil
		}

		ss_log.Error("插入一条新数据 eventNo:[%s]", eventNo)
	} else {
		err := dao.RiskDaoInstance.ModifyOp(tx, req.OpNo, req.OpName, req.ScriptName, req.Param, req.Score)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}

	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeleteOp(ctx context.Context, req *custProto.DeleteOpRequest, reply *custProto.DeleteOpReply) error {
	if req.OpNo == "" {
		ss_log.Error("OpNo为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	err := dao.RiskDaoInstance.DelOp(tx, req.OpNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetRelaApiEvents(ctx context.Context, req *custProto.GetRelaApiEventsRequest, reply *custProto.GetRelaApiEventsReply) error {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		//{Key: "is_delete", Val: "0", EqType: "<="},
	})
	total, errCnt := dao.RiskDaoInstance.GetGetRelaApiEventsCnt(whereModel.WhereStr, whereModel.Args)
	if errCnt != nil {
		ss_log.Error("err=[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//添加排序
	//ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by create_time desc `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.RiskDaoInstance.GetGetRelaApiEventsInfos(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) InsertOrUpdateRelaApiEvent(ctx context.Context, req *custProto.InsertOrUpdateRelaApiEventRequest, reply *custProto.InsertOrUpdateRelaApiEventReply) error {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)
	err := dao.RiskDaoInstance.AddOrUpdateRelaApiEvent(tx, req.ApiType, req.EventNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeleteRelaApiEvent(ctx context.Context, req *custProto.DeleteRelaApiEventRequest, reply *custProto.DeleteRelaApiEventReply) error {
	if req.ApiType == "" {
		ss_log.Error("ApiType为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	err := dao.RiskDaoInstance.DelRelaApiEvent(tx, req.ApiType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetRelaEventRules(ctx context.Context, req *custProto.GetRelaEventRulesRequest, reply *custProto.GetRelaEventRulesReply) error {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		//{Key: "is_delete", Val: "0", EqType: "<="},
	})
	total, errCnt := dao.RiskDaoInstance.GetRelaEventRulesCnt(whereModel.WhereStr, whereModel.Args)
	if errCnt != nil {
		ss_log.Error("err=[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//添加排序
	//ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by create_time desc `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.RiskDaoInstance.GetRelaEventRulesInfos(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) InsertOrUpdateRelaEventRule(ctx context.Context, req *custProto.InsertOrUpdateRelaEventRuleRequest, reply *custProto.InsertOrUpdateRelaEventRuleReply) error {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	if req.RelaNo == "" {
		relaNo, err := dao.RiskDaoInstance.AddRelaEventRule(tx, req.EventNo, req.RuleNo)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_ADD
			return nil
		}

		ss_log.Error("插入一条新数据 relaNo:[%s]", relaNo)
	} else {
		err := dao.RiskDaoInstance.ModifyRelaEventRule(tx, req.RelaNo, req.EventNo, req.RuleNo)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}

	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeleteRelaEventRule(ctx context.Context, req *custProto.DeleteRelaEventRuleRequest, reply *custProto.DeleteRelaEventRuleReply) error {
	if req.RelaNo == "" {
		ss_log.Error("RelaNo为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	err := dao.RiskDaoInstance.DelRelaEventRule(tx, req.RelaNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetRiskThresholds(ctx context.Context, req *custProto.GetRiskThresholdsRequest, reply *custProto.GetRiskThresholdsReply) error {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
	})
	total, errCnt := dao.RiskDaoInstance.GetRiskThresholdsCnt(whereModel.WhereStr, whereModel.Args)
	if errCnt != nil {
		ss_log.Error("err=[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//添加排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by create_time desc `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.RiskDaoInstance.GetRiskThresholdsInfos(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) InsertOrUpdateRiskThreshold(ctx context.Context, req *custProto.InsertOrUpdateRiskThresholdRequest, reply *custProto.InsertOrUpdateRiskThresholdReply) error {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	err := dao.RiskDaoInstance.AddOrUpdateRiskThreshold(tx, req.RuleNo, req.RiskThreshold)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_ADD
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeleteRiskThreshold(ctx context.Context, req *custProto.DeleteRiskThresholdRequest, reply *custProto.DeleteRiskThresholdReply) error {
	if req.RuleNo == "" {
		ss_log.Error("RuleNo为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	err := dao.RiskDaoInstance.DelRiskThreshold(tx, req.RuleNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetRules(ctx context.Context, req *custProto.GetRulesRequest, reply *custProto.GetRulesReply) error {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
	})
	total, errCnt := dao.RiskDaoInstance.GetRulesCnt(whereModel.WhereStr, whereModel.Args)
	if errCnt != nil {
		ss_log.Error("err=[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//添加排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by create_time desc `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.RiskDaoInstance.GetRulesInfos(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetRule(ctx context.Context, req *custProto.GetRuleRequest, reply *custProto.GetRuleReply) error {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
		{Key: "rule_no", Val: req.RuleNo, EqType: "="},
	})

	data, err := dao.RiskDaoInstance.GetRuleDetail(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) InsertOrUpdateRule(ctx context.Context, req *custProto.InsertOrUpdateRuleRequest, reply *custProto.InsertOrUpdateRuleReply) error {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	if req.RuleNo == "" {
		ruleNo, err := dao.RiskDaoInstance.AddRule(tx, req.RuleName, req.Rule)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_ADD
			return nil
		}
		ss_log.Info("插入一条新规则[%s]", ruleNo)

	} else {
		err := dao.RiskDaoInstance.ModifyRule(tx, req.RuleNo, req.RuleName, req.Rule)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeleteRule(ctx context.Context, req *custProto.DeleteRuleRequest, reply *custProto.DeleteRuleReply) error {
	if req.RuleNo == "" {
		ss_log.Error("RuleNo为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	err := dao.RiskDaoInstance.DelRule(tx, req.RuleNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
查询现金充值订单
*/
func (c *CustHandler) GetCashRechargeOrderList(ctx context.Context, req *custProto.GetCashRechargeOrderListRequest, reply *custProto.GetCashRechargeOrderListReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "scro.log_no", Val: req.LogNo, EqType: "like"},
		{Key: "scro.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "scro.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "scro.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "scro.currency_type", Val: req.CurrencyType, EqType: "="},
		{Key: "acc.account", Val: req.AccAccount, EqType: "like"},
		{Key: "acc2.account", Val: req.OpAccAccount, EqType: "like"},
	}

	//获取总数
	total := dao.CashRechargeOrderDaoInst.GetCnt(whereList)

	//获取列表信息
	datas, err := dao.CashRechargeOrderDaoInst.GetCashRechargeOrders(whereList, strext.ToInt(req.Page), strext.ToInt(req.PageSize))
	if err != nil {
		ss_log.Error("查询现金存款订单列表失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
查询改变余额订单列表
*/
func (c *CustHandler) GetChangeBalanceOrders(ctx context.Context, req *custProto.GetChangeBalanceOrdersRequest, reply *custProto.GetChangeBalanceOrdersReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "cbo.log_no", Val: req.LogNo, EqType: "like"},
		{Key: "cbo.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "cbo.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "cbo.account_no", Val: req.AccountNo, EqType: "="},
		{Key: "cbo.account_type", Val: req.AccountType, EqType: "="},
		{Key: "acc.account", Val: req.Account, EqType: "like"},
	}

	//获取总数
	total := dao.ChangeBalanceOrderDaoInst.GetCnt(whereList)

	//获取列表信息
	datas, err := dao.ChangeBalanceOrderDaoInst.GetChangeBalanceOrders(whereList, strext.ToInt(req.Page), strext.ToInt(req.PageSize))
	if err != nil {
		ss_log.Error("查询修改余额订单列表失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	for _, data := range datas {
		if data.AccountType == constants.AccountType_SERVICER { //服务商的需要做处理。负数代表它可用多少钱，而web前端显示的是正数
			data.AfterBalance = ss_count.Sub("0", data.AfterBalance).String()
			data.BeforeBalance = ss_count.Sub("0", data.BeforeBalance).String()
		}
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
查询改变余额订单详情
*/
func (c *CustHandler) ChangeBalanceDetail(ctx context.Context, req *custProto.ChangeBalanceDetailRequest, reply *custProto.ChangeBalanceDetailReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "cbo.log_no", Val: req.OrderNo, EqType: "="},
	}

	data, err := dao.ChangeBalanceOrderDaoInst.GetChangeBalanceOrderDetail(whereList)
	if err != nil {
		ss_log.Error("查询修改余额订单详情失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	if data.AccountType == constants.AccountType_SERVICER { //服务商的需要做处理。负数代表它可用多少钱，而前端显示的是正数
		data.AfterBalance = ss_count.Sub("0", data.AfterBalance).String()
		data.BeforeBalance = ss_count.Sub("0", data.BeforeBalance).String()
	}

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
  查询公告列表
*/
func (c *CustHandler) GetBulletins(ctx context.Context, req *custProto.GetBulletinsRequest, reply *custProto.GetBulletinsReply) error {

	useStatus := ""
	//排序条件
	orderByStr := " ORDER BY case use_status when " + constants.BulletinUseStatus_UnBulletin + " then 1 end," +
		" case top_status when " + constants.BulletinTopStatus_True + " then 2 end "
	switch req.AccountType {
	case constants.AccountType_ADMIN:
		fallthrough
	case constants.AccountType_OPERATOR:
		orderByStr += ", create_time desc "
	case constants.AccountType_PersonalBusiness:
		//个人商家只看到发布所有和发布给个人商家的公告
		useStatus = "('" + constants.BulletinUseStatus_All + "','" + constants.BulletinUseStatus_personal + "')"
		orderByStr += ", bulletin_time desc "
	case constants.AccountType_EnterpriseBusiness:
		//企业商家只看到发布所有和发布给企业商家的公告
		useStatus = "('" + constants.BulletinUseStatus_All + "','" + constants.BulletinUseStatus_Enterprise + "')"
		orderByStr += ", bulletin_time desc "
	}

	whereList := []*model.WhereSqlCond{
		{Key: "create_time", Val: req.StartTime, EqType: ">="},
		{Key: "create_time", Val: req.EndTime, EqType: "<="},
		{Key: "use_status", Val: req.UseStatus, EqType: "="},
		{Key: "use_status", Val: useStatus, EqType: "in"},
		{Key: "title", Val: req.Title, EqType: "like"},
	}

	//获取总数
	total := dao.BulletinDaoInst.GetCnt(whereList)

	//获取列表信息
	datas, err := dao.BulletinDaoInst.GetBulletins(whereList, strext.ToInt(req.Page), strext.ToInt(req.PageSize), orderByStr)
	if err != nil {
		ss_log.Error("查询公告列表失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
  查询公告详情
*/
func (c *CustHandler) GetBulletinDetail(ctx context.Context, req *custProto.GetBulletinDetailRequest, reply *custProto.GetBulletinDetailReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "bulletin_id", Val: req.BulletinId, EqType: "="},
	}

	//获取列表信息
	data, err := dao.BulletinDaoInst.GetBulletinDetail(whereList)
	if err != nil {
		ss_log.Error("查询公告列表失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
  插入或更新公告
*/
func (c *CustHandler) InsertOrUpdateBulletin(ctx context.Context, req *custProto.InsertOrUpdateBulletinRequest, reply *custProto.InsertOrUpdateBulletinReply) error {
	if req.BulletinId == "" {
		if _, err := dao.BulletinDaoInst.AddBulletin(req.Title, req.Content); err != nil {
			ss_log.Error("添加公告失败，err[%v]", err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	} else {
		//确认公告是未发布的才可以修改
		if useStatus := dao.BulletinDaoInst.GetBulletinUseStatus(req.BulletinId); useStatus != constants.BulletinUseStatus_UnBulletin {
			ss_log.Error("无法编辑不是未发布状态的公告[%v]", req.BulletinId)
			reply.ResultCode = ss_err.ERR_BulletinUseStatus_FAILD
			return nil
		}

		if err := dao.BulletinDaoInst.UpdateBulletin(req.BulletinId, req.Title, req.Content); err != nil {
			ss_log.Error("修改公告失败，err[%v]", err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
  删除公告
*/
func (c *CustHandler) DelBulletin(ctx context.Context, req *custProto.DelBulletinRequest, reply *custProto.DelBulletinReply) error {
	if err := dao.BulletinDaoInst.DeleteBulletin(req.BulletinId); err != nil {
		ss_log.Error("删除公告失败，err[%v]", err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
  修改公告的状态
*/
func (c *CustHandler) UpdateBulletinStatus(ctx context.Context, req *custProto.UpdateBulletinStatusRequest, reply *custProto.UpdateBulletinStatusReply) error {
	switch req.StatusType {
	case "use_status":
		switch req.Status {
		case constants.BulletinUseStatus_All:
		case constants.BulletinUseStatus_personal:
		case constants.BulletinUseStatus_Enterprise:
		default:
			ss_log.Error("Status参数不合法")
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		//确认公告是未发布的才可以修改
		if useStatus := dao.BulletinDaoInst.GetBulletinUseStatus(req.BulletinId); useStatus != constants.BulletinUseStatus_UnBulletin {
			ss_log.Error("无法编辑不是未发布状态的公告[%v]", req.BulletinId)
			reply.ResultCode = ss_err.ERR_BulletinUseStatus_FAILD
			return nil
		}

		if err := dao.BulletinDaoInst.UpdateUseStatus(req.BulletinId, req.Status); err != nil {
			ss_log.Error("修改公告发布状态失败，err[%v]", err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	case "top_status":
		switch req.Status {
		case constants.BulletinTopStatus_False:
		case constants.BulletinTopStatus_True:
		default:
			ss_log.Error("Status参数不合法")
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		if err := dao.BulletinDaoInst.UpdateTopStatus(req.BulletinId, req.Status); err != nil {
			ss_log.Error("修改公告置顶状态失败，err[%v]", err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	default:
		ss_log.Error("StatusType参数不合法")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//修改商家认证信息
func (c *CustHandler) UpdateBusinessAuthMaterialInfo(ctx context.Context, req *custProto.UpdateBusinessAuthMaterialInfoRequest, reply *custProto.UpdateBusinessAuthMaterialInfoReply) error {
	//确认没有正审核的修改认证信息（目前只修改简称）
	if !dao.AuthMaterialDaoInst.CheckAuthMaterialPendingStatusUnique(req.AccountUid, req.AccountType) {
		ss_log.Error("已提交过修改认证信息的申请，不允许重复提交")
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}

	authMaterialNo := ""
	oldSimplifyName := ""
	switch req.AccountType {
	case constants.AccountType_PersonalBusiness: //
		//查询该账号的通过的商家认证信息
		whereList := []*model.WhereSqlCond{
			{Key: "amb.account_uid", Val: req.AccountUid, EqType: "="},
			{Key: "amb.status", Val: constants.AuthMaterialStatus_Passed, EqType: "="},
		}

		data, err := dao.AuthMaterialDaoInst.GetAuthMaterialBusinessDetail(whereList)
		if err != nil {
			ss_log.Error("查询账号[%v]的个人商家认证信息失败，err=[%v]", req.AccountUid, err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}
		authMaterialNo = data.AuthMaterialNo
		oldSimplifyName = data.SimplifyName
	case constants.AccountType_EnterpriseBusiness:
		//查询该账号的通过的商家认证信息
		whereList := []*model.WhereSqlCond{
			{Key: "account_uid", Val: req.AccountUid, EqType: "="},
			{Key: "status", Val: constants.AuthMaterialStatus_Passed, EqType: "="},
		}

		data, err := dao.AuthMaterialDaoInst.GetAuthMaterialEnterpriseDetail(whereList)
		if err != nil {
			ss_log.Error("查询账号[%v]的企业商家认证信息失败，err=[%v]", req.AccountUid, err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}

		authMaterialNo = data.AuthMaterialNo
		oldSimplifyName = data.SimplifyName
	default:
		ss_log.Error("账号类型不是个人商家、企业商家, AccountType[%v]", req.AccountType)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}

	if authMaterialNo == "" {
		ss_log.Error("未查询到通过的商家认证信息，uid[%v], accType[%v] ", req.AccountUid, req.AccountType)
		reply.ResultCode = ss_err.ERR_HaveNotPass_BusinessRealName_FAILD
		return nil
	}

	if oldSimplifyName == req.SimplifyName {
		ss_log.Error("商家认证信息的简称和原来的相同，uid[%v], accType[%v] ", req.AccountUid, req.AccountType)
		reply.ResultCode = ss_err.ERR_BusinessSimplifyName_Unique_FAILD
		return nil
	}

	if err := dao.AuthMaterialDaoInst.AddAuthMaterialUpdateInfo(authMaterialNo, req.AccountUid, req.AccountType, req.SimplifyName, oldSimplifyName); err != nil {
		ss_log.Error("添加修改商家简称失败，err[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE //对商家来说是修改失败
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
