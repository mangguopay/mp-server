package handler

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/dao"
	"strings"

	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"context"
)

//商家转账订单列表(商家获取)
func (*CustHandler) GetBusinessTransferOrderList(ctx context.Context, req *custProto.GetBusinessTransferOrderListRequest, reply *custProto.GetBusinessTransferOrderListReply) error {
	ss_log.Info("查询商家转账订单接口请求参数：%v", strext.ToJson(req))
	if req.TransferType != "1" && req.TransferType != "2" {
		req.TransferType = ""
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "bto.from_account_no", Val: req.BusinessAccNo, EqType: "="},
		{Key: "bto.to_account_no", Val: req.ToAccountNo, EqType: "="},
		{Key: "bto.currency_type", Val: strings.ToUpper(req.CurrencyType), EqType: "="},
		{Key: "bto.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "bto.batch_no", Val: req.BatchNo, EqType: "="},
		{Key: "bto.log_no", Val: req.LogNo, EqType: "="},
		{Key: "bto.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "bto.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "bto.transfer_type", Val: req.TransferType, EqType: "="},
		{Key: "bto.out_transfer_no", Val: req.OutTransferNo, EqType: "="},
	})

	orderNum, err := dao.BusinessTransferDaoInst.CountOrderNum(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("查询商家转账订单数量失败，BusinessAccNo=%v, err=%v", req.BusinessAccNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY create_time DESC ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
	orderList, err := dao.BusinessTransferDaoInst.GetOrderList(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("查询商家转账订单列表失败，BusinessAccNo=%v, err=%v", req.BusinessAccNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	var WrongReasonArr []string
	var list []*custProto.BusinessTransferOrder
	for _, v := range orderList {
		order := new(custProto.BusinessTransferOrder)
		order.LogNo = v.LogNo
		order.Amount = v.Amount
		order.Fee = v.Fee
		order.RealAmount = v.RealAmount
		order.CurrencyType = v.CurrencyType
		order.CreateTime = v.CreateTime
		order.OrderStatus = v.OrderStatus
		order.Remarks = v.Remarks
		order.PayeeAccount = v.PayeeAccount
		order.WrongReason = v.WrongReason
		order.AuthName = v.AuthName
		order.ToAccount = v.ToAccount
		order.TransferType = v.TransferType
		order.OutTransferNo = v.OutTransferNo
		list = append(list, order)
		if v.WrongReason != "" {
			WrongReasonArr = append(WrongReasonArr, v.WrongReason)
		}
	}

	//处理返回去的异常原因（记录的是错误码，需换成多语言对应的语言）
	if len(WrongReasonArr) != 0 {
		WrongReasonLangDatas, errLang := dao.LangDaoInst.GetLangTextsByKeys(WrongReasonArr)
		if errLang != nil {
			ss_log.Error("err=[%v]", errLang)
		}

		for _, data := range list {
			for _, v := range WrongReasonLangDatas {
				if v.Key == data.WrongReason {
					switch req.Lang {
					case constants.LangEnUS:
						data.WrongReason = v.LangEn
					case constants.LangKmKH:
						data.WrongReason = v.LangKm
					case constants.LangZhCN:
						data.WrongReason = v.LangCh
					}
				}
			}
		}
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Total = orderNum
	reply.List = list
	return nil
}

//商家转账订单详情
func (*CustHandler) GetBusinessTransferOrderDetail(ctx context.Context, req *custProto.GetBusinessTransferOrderDetailRequest, reply *custProto.GetBusinessTransferOrderDetailReply) error {
	if req.LogNo == "" {
		ss_log.Error("LogNo")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	data, err := dao.BusinessTransferDaoInst.GetOrderDetail(req.LogNo)
	if err != nil {
		ss_log.Error("查询商家转账订单详情失败，LogNo=%v, err=%v", req.LogNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	order := new(custProto.BusinessTransferOrder)
	order.LogNo = data.LogNo
	order.PayerAccount = data.FromBusinessAccount
	order.PayerName = data.FromBusinessName
	order.PayeeAccount = data.ToAccount
	order.PayeeName = data.ToBusinessName
	order.Amount = data.Amount
	order.Fee = data.Fee
	order.RealAmount = data.RealAmount
	order.CurrencyType = data.CurrencyType
	order.CreateTime = data.CreateTime
	order.OrderStatus = data.OrderStatus
	order.Remarks = data.Remarks
	order.BatchNo = data.BatchNo
	order.TransferType = data.TransferType

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = order
	return nil
}

func (*CustHandler) GetBusinessTransferBatchList(ctx context.Context, req *custProto.GetBusinessTransferBatchListRequest, reply *custProto.GetBusinessTransferBatchListReply) error {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "btb.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "btb.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "btb.batch_no", Val: req.BatchNo, EqType: "="},
		{Key: "btb.business_no", Val: req.BusinessNo, EqType: "="},
		{Key: "btb.status", Val: req.Status, EqType: "="},
		{Key: "acc.account", Val: req.Account, EqType: "like"},
		{Key: "btb.currency_type", Val: strings.ToUpper(req.CurrencyType), EqType: "="},
		{Key: "btb.status", Val: constants.BusinessBatchTransferOrderStatusPending, EqType: "!="},
	})

	total, errCnt := dao.BusinessBatchTransferDaoInst.CountOrderNum(whereModel.WhereStr, whereModel.Args)
	if errCnt != nil {
		ss_log.Error("查询商家批量转账订单数量,err[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY btb.create_time DESC")
	ss_sql.SsSqlFactoryInst.AppendWhereLimitI32(whereModel, strext.ToInt32(req.PageSize), strext.ToInt32(req.Page))
	datas, err := dao.BusinessBatchTransferDaoInst.GetOrderList(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("查询商家批量转账订单列表失败,err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Total = total
	reply.Datas = datas
	return nil
}

func (*CustHandler) GetBusinessTransferBatchDetail(ctx context.Context, req *custProto.GetBusinessTransferBatchDetailRequest, reply *custProto.GetBusinessTransferBatchDetailReply) error {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "btb.batch_no", Val: req.BatchNo, EqType: "="},
	})

	batchData, err := dao.BusinessBatchTransferDaoInst.GetBatchOrderDetail(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("查询商家批量转账订单列表失败,err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.BatchData = batchData
	return nil
}

//商家转账订单列表(管理后台获取)
func (*CustHandler) ManagementGetBusinessTransferOrderList(ctx context.Context, req *custProto.ManagementGetBusinessTransferOrderListRequest, reply *custProto.ManagementGetBusinessTransferOrderListReply) error {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "bto.currency_type", Val: strings.ToUpper(req.CurrencyType), EqType: "="},
		{Key: "bto.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "bto.batch_no", Val: req.BatchNo, EqType: "="},
		{Key: "bto.log_no", Val: req.LogNo, EqType: "="},
		{Key: "acc.account", Val: req.ToAccount, EqType: "like"},
		{Key: "acc2.account", Val: req.FromAccount, EqType: "like"},
		{Key: "bto.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "bto.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "bto.transfer_type", Val: req.TransferType, EqType: "="},
	})

	orderNum, err := dao.BusinessTransferDaoInst.GetCnt(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("查询商家转账订单数量失败,err=%v", err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY bto.create_time DESC ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
	datas, err := dao.BusinessTransferDaoInst.GetTransferOrderList(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("查询商家转账订单列表失败，err=%v", err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	var WrongReasonArr []string
	for _, data := range datas {
		if data.WrongReason != "" {
			WrongReasonArr = append(WrongReasonArr, data.WrongReason)
		}
	}

	//处理返回去的异常原因（记录的是错误码，需换成多语言对应的语言）
	if len(WrongReasonArr) != 0 {
		WrongReasonLangDatas, errLang := dao.LangDaoInst.GetLangTextsByKeys(WrongReasonArr)
		if errLang != nil {
			ss_log.Error("err=[%v]", errLang)
		}

		for _, data := range datas {
			for _, v := range WrongReasonLangDatas {
				if v.Key == data.WrongReason {
					data.WrongReason = v.LangCh
				}
			}
		}
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Total = orderNum
	reply.Datas = datas
	return nil
}
