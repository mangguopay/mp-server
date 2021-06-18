package handler

import (
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/model"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/dao"
	"a.a/mp-server/cust-srv/util"
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
)

//查询商家交易订单
func (c *CustHandler) GetBusinessBills(ctx context.Context, req *custProto.GetBusinessBillsRequest, reply *custProto.GetBusinessBillsReply) error {
	/**
	待支付订单: orderStatus=0, isSettled=0或为空
	待结算订单：orderStatus为空, isSettled=0
	已结算订单：orderStatus为空, isSettled=1
	*/
	whereList := []*model.WhereSqlCond{
		{Key: "bc.channel_type", Val: req.ChannelType, EqType: "="},
		{Key: "bb.business_no", Val: req.BusinessNo, EqType: "="},
		{Key: "bb.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "bb.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "bb.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "bb.order_no", Val: req.OrderNo, EqType: "="},
		{Key: "bb.settle_id", Val: req.SettleId, EqType: "="},
		{Key: "b.business_id", Val: req.BusinessId, EqType: "="},
		{Key: "bs.scene_no", Val: req.SceneNo, EqType: "="},
		{Key: "b.full_name", Val: req.BusinessName, EqType: "like"},
		{Key: "acc.account", Val: req.BusinessAccount, EqType: "like"},
		{Key: "bb.subject", Val: req.Subject, EqType: "like"},
	}

	countWhereList := whereList

	isSettled := ""
	eqType := "="
	if req.IsSettled != "" {
		if req.OrderStatus == "" {
			//查询待结算或者已结算订单都不允许未支付和支付失败的订单出现
			whereList = append(whereList, &model.WhereSqlCond{Key: "bb.order_status", Val: constants.BusinessOrderStatusPending, EqType: "!="})
			whereList = append(whereList, &model.WhereSqlCond{Key: "bb.order_status", Val: constants.BusinessOrderStatusPayTimeOut, EqType: "!="})
		}

		switch req.IsSettled {
		case "0": //待结算,  bb.settle_id == ''
			isSettled = "__empty_string"
			eqType = "="
			if req.OrderStatus == "" {
				whereList = append(whereList, &model.WhereSqlCond{Key: "bb.order_status", Val: constants.BusinessOrderStatusPay, EqType: "="})
			}
			countWhereList = whereList
		case "1": //已结算, bb.settle_id != ''
			isSettled = "__empty_string"
			eqType = "!="

			inStatus := fmt.Sprintf("(%v, %v)", constants.BusinessOrderStatusPay, constants.BusinessOrderStatusRebatesRefund)
			countWhereList = append(countWhereList, &model.WhereSqlCond{Key: "bb.order_status", Val: inStatus, EqType: "in"})
		}

	}
	whereList = append(whereList, &model.WhereSqlCond{Key: "bb.settle_id", Val: isSettled, EqType: eqType})
	countWhereList = append(countWhereList, &model.WhereSqlCond{Key: "bb.settle_id", Val: isSettled, EqType: eqType})

	//获取usd统计信息
	reply.UsdCnt = dao.BusinessBillDaoInst.GetCnt(append(countWhereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_USD, EqType: "="}))
	reply.UsdSum = dao.BusinessBillDaoInst.GetSum(append(countWhereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_USD, EqType: "="}))

	//获取khr统计信息
	reply.KhrCnt = dao.BusinessBillDaoInst.GetCnt(append(countWhereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_KHR, EqType: "="}))
	reply.KhrSum = dao.BusinessBillDaoInst.GetSum(append(countWhereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_KHR, EqType: "="}))

	//追加条件
	whereList = append(whereList, &model.WhereSqlCond{
		Key: "bb.currency_type", Val: req.CurrencyType, EqType: "=",
	})

	//获取总数
	total := dao.BusinessBillDaoInst.GetCnt(whereList)

	//获取列表信息
	orderList, err := dao.BusinessBillDaoInst.GetBusinessBills(whereList, strext.ToInt(req.Page), strext.ToInt(req.PageSize))
	if err != nil {
		ss_log.Error("查询订单失败，req=%v, err=[%v]", strext.ToJson(req), err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	var keys []string
	for _, data := range orderList {
		keys = append(keys, data.SceneName)
	}

	//读取多语言
	langDatas, errLang := dao.LangDaoInst.GetLangTextsByKeys(keys)
	if errLang != nil {
		ss_log.Error("查询多语言出错,keys[%v]", keys)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	for _, data := range orderList {
		for _, langData := range langDatas {
			switch req.Lang {
			case constants.LangZhCN:
				if data.SceneName == langData.Key {
					data.SceneName = langData.LangCh
				}
			case constants.LangEnUS:
				if data.SceneName == langData.Key {
					data.SceneName = langData.LangEn
				}
			case constants.LangKmKH:
				if data.SceneName == langData.Key {
					data.SceneName = langData.LangKm
				}
			default:

			}
		}
	}

	var list []*custProto.BusinessBillData
	for _, v := range orderList {
		settleDate := time.Unix(strext.ToInt64(v.SettleDate), 0).Format(ss_time.DateTimeDashFormat)
		order := &custProto.BusinessBillData{
			CreateTime:      v.CreateTime,
			OrderNo:         v.OrderNo,
			OutOrderNo:      v.OutOrderNo,
			Subject:         v.Subject,
			Amount:          v.Amount,
			RealAmount:      v.RealAmount,
			CurrencyType:    v.CurrencyType,
			Fee:             v.Fee,
			OrderStatus:     v.OrderStatus,
			SettleId:        v.SettleId,
			SceneName:       v.SceneName,
			AppName:         v.AppName,
			BusinessName:    v.BusinessName,
			BusinessId:      v.BusinessId,
			BusinessAccount: v.BusinessAccount,
			Cycle:           v.Cycle,
			Rate:            v.Rate,
			SettleDate:      settleDate,
		}
		if order.OrderStatus == constants.BusinessOrderStatusRebatesRefund {
			whereList := []*model.WhereSqlCond{
				{Key: "br.pay_order_no", Val: order.OrderNo, EqType: "="},
				{Key: "br.refund_status", Val: constants.BusinessRefundStatusSuccess, EqType: "="},
			}
			//查询订单已退款金额
			refundedAmount, err := dao.BusinessBillRefundDaoInst.GetSum(whereList)
			if err != nil {
				ss_log.Error("查询订单[%v]已退款金额失败，err=[%v]", order.OrderNo, err)
				reply.ResultCode = ss_err.ERR_SYSTEM
				return nil
			}
			switch order.CurrencyType {
			case constants.CURRENCY_UP_USD:
				reply.UsdSum = ss_count.Sub(reply.UsdSum, strext.ToString(refundedAmount)).String()
			case constants.CURRENCY_UP_KHR:
				reply.KhrSum = ss_count.Sub(reply.KhrSum, strext.ToString(refundedAmount)).String()
			}
		}
		list = append(list, order)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = list
	reply.Total = strext.ToInt32(total)

	return nil
}

//查询商家交易订单详情
func (c *CustHandler) GetBusinessBillDetail(ctx context.Context, req *custProto.GetBusinessBillDetailRequest, reply *custProto.GetBusinessBillDetailReply) error {
	if req.OrderNo == "" {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//获取列表信息
	data, err := dao.BusinessBillDaoInst.GetBusinessBillDetail(req.OrderNo)
	if err != nil {
		ss_log.Error("查询商家订单[%v]失败，err=[%v]", req.OrderNo, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	var keys []string
	keys = append(keys, data.SceneName)

	//读取多语言
	langDatas, errLang := dao.LangDaoInst.GetLangTextsByKeys(keys)
	if errLang != nil {
		ss_log.Error("查询多语言出错,keys[%v]", keys)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	for _, langData := range langDatas {
		switch req.Lang {
		case constants.LangZhCN:
			if data.SceneName == langData.Key {
				data.SceneName = langData.LangCh
			}
		case constants.LangEnUS:
			if data.SceneName == langData.Key {
				data.SceneName = langData.LangEn
			}

		case constants.LangKmKH:
			if data.SceneName == langData.Key {
				data.SceneName = langData.LangKm
			}
		}
	}

	if data.SettleDate != "" {
		data.SettleDate = ss_time.Unixtime2Time(data.SettleDate, global.Tz).Format(ss_time.DateFormat)
	}

	whereList := []*model.WhereSqlCond{
		{Key: "br.pay_order_no", Val: data.OrderNo, EqType: "="},
		{Key: "br.refund_status", Val: constants.BusinessRefundStatusSuccess, EqType: "="},
	}
	//查询订单已退款金额
	refundedAmount, err := dao.BusinessBillRefundDaoInst.GetSum(whereList)
	if err != nil {
		ss_log.Error("查询订单[%v]已退款金额失败，err=[%v]", data.OrderNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}
	data.RefundedAmount = strext.ToString(refundedAmount)

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

//生成商家交易订单文件
func (c *CustHandler) CreateBillFile(ctx context.Context, req *custProto.CreateBillFileRequest, reply *custProto.CreateBillFileReply) error {
	/**
	待支付订单: orderStatus=1, isSettled=为空
	待结算订单：orderStatus为空, isSettled=0
	已结算订单：orderStatus为空, isSettled=1
	*/
	whereList := []*model.WhereSqlCond{
		{Key: "bc.channel_type", Val: req.ChannelType, EqType: "="}, //ModernPay渠道交易订单
		{Key: "bb.business_no", Val: req.IdenNo, EqType: "="},
		{Key: "bb.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "bb.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "bb.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "bb.order_no", Val: req.OrderNo, EqType: "="},
		{Key: "bs.scene_name", Val: req.SceneName, EqType: "like"},
	}

	isSettled := ""
	eqType := "="
	fileHead := "#Mango Pay"
	if req.IsSettled != "" {
		switch req.IsSettled {
		case "0": //待结算,  bb.settle_id == ''
			isSettled = "__empty_string"
			eqType = "="
			fileHead += "待结算订单明细"
		case "1": //已结算, bb.settle_id != ''
			isSettled = "__empty_string"
			eqType = "!="
			fileHead += "已结算订单明细"
		}
		if req.OrderStatus == "" {
			//查询待结算或者已结算订单都不允许未支付和支付失败的订单出现
			whereList = append(whereList, []*model.WhereSqlCond{
				{Key: "bb.order_status", Val: constants.BusinessOrderStatusPending, EqType: "!="},
				{Key: "bb.order_status", Val: constants.BusinessOrderStatusPayTimeOut, EqType: "!="},
			}...)
		}
	}
	whereList = append(whereList, &model.WhereSqlCond{Key: "bb.settle_id", Val: isSettled, EqType: eqType})

	sucWhereList := append(whereList, []*model.WhereSqlCond{
		{Key: "bb.order_status", Val: constants.BusinessOrderStatusPay, EqType: "="},
	}...)

	refWhereList := append(whereList, []*model.WhereSqlCond{
		{Key: "bb.order_status", Val: constants.BusinessOrderStatusRefund, EqType: "="},
		{Key: "bb.order_status", Val: constants.BusinessOrderStatusRebatesRefund, EqType: "="},
	}...)

	//usd
	usdCnt := dao.BusinessBillDaoInst.GetCnt(append(sucWhereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_USD, EqType: "="}))
	usdSum := dao.BusinessBillDaoInst.GetSum(append(sucWhereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_USD, EqType: "="}))
	usdRefundCnt := dao.BusinessBillDaoInst.GetCnt(append(refWhereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_USD, EqType: "="}))
	usdRefundSum := dao.BusinessBillDaoInst.GetSum(append(refWhereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_USD, EqType: "="}))

	//获取khr统计信息
	khrCnt := dao.BusinessBillDaoInst.GetCnt(append(sucWhereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_KHR, EqType: "="}))
	khrSum := dao.BusinessBillDaoInst.GetSum(append(sucWhereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_KHR, EqType: "="}))
	khrRefundCnt := dao.BusinessBillDaoInst.GetCnt(append(refWhereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_KHR, EqType: "="}))
	khrRefundSum := dao.BusinessBillDaoInst.GetSum(append(refWhereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_KHR, EqType: "="}))

	//追加条件
	whereList = append(whereList, &model.WhereSqlCond{
		Key: "bb.currency_type", Val: req.CurrencyType, EqType: "=",
	})

	//获取总数
	total := dao.BusinessBillDaoInst.GetCnt(whereList)

	//获取列表信息
	orderList, err := dao.BusinessBillDaoInst.GetBusinessBills(whereList, strext.ToInt(req.Page), strext.ToInt(req.PageSize))
	if err != nil {
		ss_log.Error("查询商家[%v]订单失败，err=[%v]", req.IdenNo, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	var keys []string
	var langDatas []*custProto.LangData
	keyMap := make(map[string]string) //用于去重，不用重复查询一些key
	for k, data := range orderList {
		if data.SceneName != "" { //产品名称记录的是多语言的key
			if _, ok := keyMap[data.SceneName]; !ok { //只有没添加过的才去查询
				keyMap[data.SceneName] = data.SceneName
				keys = append(keys, data.SceneName)
			}
		}

		//一次最多查30个key对应的语言
		if len(keys) == 30 || k == len(orderList)-1 {
			//读取多语言
			langDatas2, errLang := dao.LangDaoInst.GetLangTextsByKeys(keys)
			if errLang != nil {
				ss_log.Error("查询多语言出错,keys[%v]", keys)
				reply.ResultCode = ss_err.ERR_SYS_DB_GET
				return nil
			}
			langDatas = append(langDatas, langDatas2...)
			keys = keys[0:0]
		}

	}

	var list []*custProto.BusinessBillData
	for _, v := range orderList {
		for _, langData := range langDatas {
			switch req.Lang {
			case constants.LangZhCN:
				switch langData.Key {
				case v.SceneName:
					v.SceneName = langData.LangCh
				}

			case constants.LangEnUS:
				switch langData.Key {
				case v.SceneName:
					v.SceneName = langData.LangEn
				}
			case constants.LangKmKH:
				switch langData.Key {
				case v.SceneName:
					v.SceneName = langData.LangKm
				}
			}
		}

		order := &custProto.BusinessBillData{
			CreateTime:   v.CreateTime,
			OrderNo:      v.OrderNo,
			OutOrderNo:   v.OutOrderNo,
			Subject:      v.Subject,
			Amount:       v.Amount,
			RealAmount:   v.RealAmount,
			CurrencyType: v.CurrencyType,
			Fee:          v.Fee,
			OrderStatus:  v.OrderStatus,
			SettleId:     v.SettleId,
			SceneName:    v.SceneName,
			AppName:      v.AppName,
			BusinessName: v.BusinessName,
			BusinessId:   v.BusinessId,
		}
		list = append(list, order)
	}

	account, err := dao.AccDaoInstance.GetAccByUid(req.Uid)
	if err != nil {
		ss_log.Error("查询商户账号失败，accountNo=%v, err=%v", req.Uid, err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	fileName := strext.GetDailyId()
	pathStr := os.TempDir()
	filePath := fmt.Sprintf("%v/%v.xlsx", pathStr, fileName)
	ss_log.Info("fileName:[%v]", fileName)
	titles := []string{"支付订单号", "订单号", "创建时间", "币种", "金额", "手续费", "应用名称", "商品名称", "订单状态"}
	var query []string
	if req.StartTime != "" {
		query = append(query, fmt.Sprintf("#起始日期：[%v]", req.StartTime))
	}
	if req.EndTime != "" {
		query = append(query, fmt.Sprintf("终止日期：[%v]", req.EndTime))
	}
	f := &util.FileCommonInfo{
		FileName:     fileName,
		FilePath:     filePath,
		Head:         fileHead,
		Account:      account,
		Total:        total,
		USDCnt:       usdCnt,
		USDSum:       usdSum,
		USDRefundCnt: usdRefundCnt,
		USDRefundSum: usdRefundSum,
		KHRCnt:       khrCnt,
		KHRSum:       khrSum,
		KHRRefundCnt: khrRefundCnt,
		KHRRefundSum: khrRefundSum,
		QueryStr:     query,
		Titles:       titles,
		OrderList:    list,
	}
	errC := util.CreateBillFile(f)
	if errC != nil {
		ss_log.Error("创建临时文件出错,err[%v]", errC)
		reply.ResultCode = ss_err.ERR_SYS_IO_ERR
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.FileName = fileName
	reply.FilePath = filePath
	return nil
}

//查询商家退款订单
func (c *CustHandler) GetRefundBills(ctx context.Context, req *custProto.GetRefundBillsRequest, reply *custProto.GetRefundBillsReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "bc.channel_type", Val: req.ChannelType, EqType: "="},
		{Key: "bb.business_no", Val: req.BusinessNo, EqType: "="},
		{Key: "br.refund_status", Val: req.OrderStatus, EqType: "="},
		{Key: "br.finish_time", Val: req.StartTime, EqType: ">="},
		{Key: "br.finish_time", Val: req.EndTime, EqType: "<="},
		{Key: "br.refund_no", Val: req.RefundNo, EqType: "="},
		{Key: "br.pay_order_no", Val: req.TransOrderNo, EqType: "="},
		{Key: "b.business_id", Val: req.BusinessId, EqType: "="},
		{Key: "b.full_name", Val: req.BusinessName, EqType: "like"},
		{Key: "bs.scene_no", Val: req.SceneNo, EqType: "="},
		{Key: "acc.account", Val: req.PayeeAccount, EqType: "like"},
	}

	//获取usd统计信息
	usdWhereList := append(whereList, []*model.WhereSqlCond{
		{Key: "bb.currency_type", Val: constants.CURRENCY_UP_USD, EqType: "="},
		{Key: "br.refund_status", Val: constants.BusinessRefundStatusSuccess, EqType: "="},
	}...)
	usdCnt, err := dao.BusinessBillRefundDaoInst.GetCnt(usdWhereList)
	if err != nil {
		ss_log.Error("统计usd数量失败，err=%v", err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}
	usdSum, err := dao.BusinessBillRefundDaoInst.GetSum(usdWhereList)
	if err != nil {
		ss_log.Error("统计usd金额失败，err=%v", err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//获取khr统计信息
	khrWhereList := append(whereList, []*model.WhereSqlCond{
		{Key: "bb.currency_type", Val: constants.CURRENCY_UP_KHR, EqType: "="},
		{Key: "br.refund_status", Val: constants.BusinessRefundStatusSuccess, EqType: "="},
	}...)
	khrCnt, err := dao.BusinessBillRefundDaoInst.GetCnt(khrWhereList)
	if err != nil {
		ss_log.Error("统计khr数量失败，err=%v", err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}
	khrSum, err := dao.BusinessBillRefundDaoInst.GetSum(khrWhereList)
	if err != nil {
		ss_log.Error("统计khr金额失败，err=%v", err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//追加条件
	whereList = append(whereList, &model.WhereSqlCond{
		Key: "bb.currency_type", Val: req.CurrencyType, EqType: "=",
	})

	//获取总数
	total, err := dao.BusinessBillRefundDaoInst.GetCnt(whereList)

	//获取列表信息
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY br.finish_time desc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.BusinessBillRefundDaoInst.GetRefundBills(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		if err != sql.ErrNoRows {
			ss_log.Error("查询商家[%v]退款订单失败，err=[%v]", req.BusinessNo, err)
			reply.ResultCode = ss_err.ERR_SYSTEM
			return nil
		}
	}

	var keys []string
	var langDatas []*custProto.LangData
	keyMap := make(map[string]string) //用于去重，不用重复查询一些key
	for k, data := range datas {
		if data.SceneName != "" { //产品名称记录的是多语言的key
			if _, ok := keyMap[data.SceneName]; !ok { //只有没添加过的才去查询
				keyMap[data.SceneName] = data.SceneName
				keys = append(keys, data.SceneName)
			}
		}

		//一次最多查30个key对应的语言
		if len(keys) == 30 || k == len(datas)-1 {
			//读取多语言
			langDatas2, errLang := dao.LangDaoInst.GetLangTextsByKeys(keys)
			if errLang != nil {
				ss_log.Error("查询多语言出错,keys[%v]", keys)
				reply.ResultCode = ss_err.ERR_SYS_DB_GET
				return nil
			}
			langDatas = append(langDatas, langDatas2...)
			keys = keys[0:0]
		}

	}

	for _, data := range datas {
		for _, langData := range langDatas {
			switch req.Lang {
			case constants.LangZhCN:
				switch langData.Key {
				case data.SceneName:
					data.SceneName = langData.LangCh
				}
			case constants.LangEnUS:
				switch langData.Key {
				case data.SceneName:
					data.SceneName = langData.LangEn
				}
			case constants.LangKmKH:
				switch langData.Key {
				case data.SceneName:
					data.SceneName = langData.LangKm
				}

			default:

			}
		}

	}

	var list []*custProto.RefundBill
	for _, v := range datas {
		data := &custProto.RefundBill{
			RefundNo:     v.RefundNo,
			OutRefundNo:  v.OutRefundNo,
			OrderStatus:  v.RefundStatus,
			TransOrderNo: v.PayOrderNo,
			Amount:       v.RefundAmount,
			CurrencyType: v.CurrencyType,
			CreateTime:   v.CreateTime,
			FinishTime:   v.FinishTime,
			Subject:      v.Subject,
			AppName:      v.AppName,
			PayeeAccount: v.PayeeAcc,
			BusinessName: v.BusinessName,
			BusinessId:   v.BusinessId,
			SceneName:    v.SceneName,
		}
		list = append(list, data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.List = list
	reply.Total = strext.ToInt32(total)
	reply.UsdCnt = usdCnt
	reply.UsdSum = usdSum
	reply.KhrCnt = khrCnt
	reply.KhrSum = khrSum
	return nil
}

//查询商家退款订单详情
func (c *CustHandler) GetRefundDetail(ctx context.Context, req *custProto.GetRefundDetailRequest, reply *custProto.GetRefundDetailReply) error {
	if req.RefundNo == "" {
		ss_log.Error("RefundNo参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	detail, err := dao.BusinessBillRefundDaoInst.GetOrderDetail(req.RefundNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("查询不到退款记录，RefundNo=%v", req.RefundNo)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		ss_log.Error("查询退款记录失败，RefundNo=%v, err=%v", req.RefundNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//读取多语言
	langDatas, errLang := dao.LangDaoInst.GetLangTextsByKeys([]string{detail.SceneName})
	if errLang != nil {
		ss_log.Error("查询多语言出错,keys[%v]", detail.SceneName)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	for _, langData := range langDatas {
		switch req.Lang {
		case constants.LangZhCN:
			switch langData.Key {
			case detail.SceneName:
				detail.SceneName = langData.LangCh
			}
		case constants.LangEnUS:
			switch langData.Key {
			case detail.SceneName:
				detail.SceneName = langData.LangEn
			}
		case constants.LangKmKH:
			switch langData.Key {
			case detail.SceneName:
				detail.SceneName = langData.LangKm
			}
		default:

		}
	}

	data := new(custProto.RefundBill)
	data.RefundNo = detail.RefundNo
	data.OutRefundNo = detail.OutRefundNo
	data.OrderStatus = detail.RefundStatus
	data.Amount = detail.RefundAmount
	data.CurrencyType = detail.CurrencyType
	data.CreateTime = detail.CreateTime
	data.FinishTime = detail.FinishTime
	data.TransOrderNo = detail.PayOrderNo
	data.TransAmount = detail.TransAmount
	data.Subject = detail.Subject
	data.SceneName = detail.SceneName
	data.PayeeAccount = detail.PayeeAcc
	data.AppName = detail.AppName

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

//生成商家退款订单文件
func (c *CustHandler) CreateRefundFile(ctx context.Context, req *custProto.CreateRefundFileRequest, reply *custProto.CreateRefundFileReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "bc.channel_type", Val: req.ChannelType, EqType: "="},
		{Key: "bb.business_no", Val: req.BusinessNo, EqType: "="},
		{Key: "br.refund_status", Val: req.OrderStatus, EqType: "="},
		{Key: "br.finish_time", Val: req.StartTime, EqType: ">="},
		{Key: "br.finish_time", Val: req.EndTime, EqType: "<="},
		{Key: "br.refund_no", Val: req.RefundNo, EqType: "="},
		{Key: "br.pay_order_no", Val: req.TransOrderNo, EqType: "="},
	}

	var f = new(util.FileCommonInfo)
	var err error
	//获取usd统计信息
	usdWhereList := append(whereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_USD, EqType: "="})
	f.USDRefundBeingCnt, f.USDRefundBeingSum, err = dao.BusinessBillRefundDaoInst.GetCntAndSum(append(usdWhereList, &model.WhereSqlCond{Key: "br.refund_status", Val: constants.BusinessRefundStatusPending, EqType: "="}))
	if err != nil {
		ss_log.Error("统计USD退款处理中订单失败，err=%v", err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}
	f.USDRefundCnt, f.USDRefundSum, err = dao.BusinessBillRefundDaoInst.GetCntAndSum(append(usdWhereList, &model.WhereSqlCond{Key: "br.refund_status", Val: constants.BusinessRefundStatusSuccess, EqType: "="}))
	if err != nil {
		ss_log.Error("统计USD退款成功订单失败，err=%v", err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}
	f.USDRefundFailCnt, f.USDRefundFailSum, err = dao.BusinessBillRefundDaoInst.GetCntAndSum(append(usdWhereList, &model.WhereSqlCond{Key: "br.refund_status", Val: constants.BusinessRefundStatusFail, EqType: "="}))
	if err != nil {
		ss_log.Error("统计USD退款失败订单失败，err=%v", err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//获取khr统计信息
	khrWhereList := append(whereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_KHR, EqType: "="})
	f.KHRRefundBeingCnt, f.KHRRefundBeingSum, err = dao.BusinessBillRefundDaoInst.GetCntAndSum(append(khrWhereList, &model.WhereSqlCond{Key: "br.refund_status", Val: constants.BusinessRefundStatusPending, EqType: "="}))
	if err != nil {
		ss_log.Error("统计KHR退款处理中订单失败，err=%v", err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}
	f.KHRRefundCnt, f.KHRRefundSum, err = dao.BusinessBillRefundDaoInst.GetCntAndSum(append(khrWhereList, &model.WhereSqlCond{Key: "br.refund_status", Val: constants.BusinessRefundStatusSuccess, EqType: "="}))
	if err != nil {
		ss_log.Error("统计KHR退款成功订单失败，err=%v", err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}
	f.KHRRefundFailCnt, f.KHRRefundFailSum, err = dao.BusinessBillRefundDaoInst.GetCntAndSum(append(khrWhereList, &model.WhereSqlCond{Key: "br.refund_status", Val: constants.BusinessRefundStatusFail, EqType: "="}))
	if err != nil {
		ss_log.Error("统计KHR退款失败订单失败，err=%v", err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//追加条件
	whereList = append(whereList, &model.WhereSqlCond{
		Key: "bb.currency_type", Val: req.CurrencyType, EqType: "=",
	})

	//获取总数
	total, err := dao.BusinessBillRefundDaoInst.GetCnt(whereList)
	if err != nil {
		ss_log.Error("查询商家退款订单总数失败，err=[%v]", err)
	}

	//获取列表信息
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY br.finish_time desc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.BusinessBillRefundDaoInst.GetRefundBills(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		if err != sql.ErrNoRows {
			ss_log.Error("查询商家[%v]退款订单失败，err=[%v]", req.BusinessNo, err)
			reply.ResultCode = ss_err.ERR_SYSTEM
			return nil
		}
	}

	var list []*custProto.RefundBill
	for _, v := range datas {
		data := &custProto.RefundBill{
			RefundNo:     v.RefundNo,
			OutRefundNo:  v.OutRefundNo,
			OrderStatus:  v.RefundStatus,
			TransOrderNo: v.PayOrderNo,
			Amount:       v.RefundAmount,
			TransAmount:  v.TransAmount,
			CurrencyType: v.CurrencyType,
			CreateTime:   v.CreateTime,
			FinishTime:   v.FinishTime,
			Subject:      v.Subject,
			AppName:      v.AppName,
			PayeeAccount: v.PayeeAcc,
		}
		list = append(list, data)
	}

	account, err := dao.AccDaoInstance.GetAccByUid(req.AccountNo)
	if err != nil {
		ss_log.Error("查询商户账号失败，accountNo=%v, err=%v", req.AccountNo, err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	fileName := strext.GetDailyId()
	pathStr := os.TempDir()
	filePath := fmt.Sprintf("%v/%v.xlsx", pathStr, fileName)
	ss_log.Info("filePath:[%v]", filePath)
	titles := []string{"支付订单号", "退款订单号", "创建时间", "币种", "支付订单金额", "实际退款金额", "应用名称", "商品名称", "订单状态"}
	var query []string
	if req.StartTime != "" {
		query = append(query, fmt.Sprintf("#起始日期：[%v]", req.StartTime))
	}
	if req.EndTime != "" {
		query = append(query, fmt.Sprintf("终止日期：[%v]", req.EndTime))
	}

	f.FileName = fileName
	f.FilePath = filePath
	f.Head = "#Mango Pay退款订单明细"
	f.Account = account
	f.QueryStr = query
	f.Titles = titles
	f.Total = strext.ToString(total)
	f.RefundOrder = list

	errC := util.CreateRefundFile(f)
	if errC != nil {
		ss_log.Error("创建临时文件出错,err[%v]", errC)
		reply.ResultCode = ss_err.ERR_SYS_IO_ERR
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.FileName = fileName
	reply.FilePath = filePath
	return nil
}

//查询渠道交易订单
func (c *CustHandler) GetChannelBills(ctx context.Context, req *custProto.GetChannelBillsRequest, reply *custProto.GetChannelBillsReply) error {
	/**
	待结算订单：orderStatus为空, isSettled=0
	已结算订单：orderStatus为空, isSettled=1
	*/
	whereList := []*model.WhereSqlCond{
		{Key: "bc.channel_type", Val: constants.ChannelTypeOut, EqType: "="}, //外部渠道交易订单
		{Key: "bb.business_channel_no", Val: req.ChannelNo, EqType: "="},
		{Key: "acc.account", Val: req.Account, EqType: "="},
		{Key: "bb.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "bb.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "bb.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "bs.scene_no", Val: req.SceneNo, EqType: "="},
		{Key: "bb.currency_type", Val: req.CurrencyType, EqType: "="},
	}

	eqType := "="
	if req.IsSettled != "" {
		switch req.IsSettled {
		case "0": //待结算
			eqType = "="
			if req.OrderStatus == "" {
				whereList = append(whereList, &model.WhereSqlCond{Key: "bb.order_status", Val: constants.BusinessOrderStatusPay, EqType: "="})
			}
		case "1": //已结算
			eqType = "!="
		}
		if req.OrderStatus == "" {
			//查询待结算或者已结算订单都不允许未支付和支付失败的订单出现
			whereList = append(whereList, &model.WhereSqlCond{Key: "bb.order_status", Val: constants.BusinessOrderStatusPending, EqType: "!="})
			whereList = append(whereList, &model.WhereSqlCond{Key: "bb.order_status", Val: constants.BusinessOrderStatusPayTimeOut, EqType: "!="})
		}
	}
	whereList = append(whereList, &model.WhereSqlCond{Key: "bb.settle_id", Val: "__empty_string", EqType: eqType})

	//获取总数
	total := dao.BusinessBillDaoInst.GetCnt(whereList)

	//获取列表信息
	orderList, err := dao.BusinessBillDaoInst.GetBusinessChannelBills(whereList, strext.ToInt(req.Page), strext.ToInt(req.PageSize))
	if err != nil {
		ss_log.Error("查询订单失败，req=%v, err=[%v]", strext.ToJson(req), err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}
	var list []*custProto.ChannelBillData
	for _, v := range orderList {
		settleDate := time.Unix(strext.ToInt64(v.SettleDate), 0).Format(ss_time.DateTimeDashFormat)
		order := &custProto.ChannelBillData{
			OrderNo:         v.OrderNo,
			OutOrderNo:      v.OutOrderNo,
			Subject:         v.Subject,
			Amount:          v.Amount,
			RealAmount:      v.RealAmount,
			CurrencyType:    v.CurrencyType,
			OrderStatus:     v.OrderStatus,
			Rate:            v.Rate,
			Fee:             v.Fee,
			Remark:          v.Remark,
			CreateTime:      v.CreateTime,
			PayTime:         v.PayTime,
			SceneName:       v.SceneName,
			AppName:         v.AppName,
			Cycle:           v.Cycle,
			SettleDate:      settleDate,
			BusinessName:    v.BusinessName,
			BusinessId:      v.BusinessId,
			BusinessAccount: v.BusinessAccount,
			ChannelName:     v.ChannelName,
			ChannelRate:     v.ChannelRate,
		}
		list = append(list, order)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.List = list
	reply.Total = strext.ToInt32(total)
	return nil
}

//查询渠道退款订单
func (c *CustHandler) GetChannelRefundBills(ctx context.Context, req *custProto.GetChannelRefundBillsRequest, reply *custProto.GetChannelRefundBillsReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "bc.channel_type", Val: constants.ChannelTypeOut, EqType: "="}, //外部渠道退款订单
		{Key: "bc.business_channel_no", Val: req.ChannelNo, EqType: "="},
		{Key: "br.refund_no", Val: req.RefundNo, EqType: "="},
		{Key: "br.out_refund_no", Val: req.OutRefundNo, EqType: "="},
		{Key: "br.refund_status", Val: req.RefundStatus, EqType: "="},
		{Key: "br.finish_time", Val: req.StartTime, EqType: ">="},
		{Key: "br.finish_time", Val: req.EndTime, EqType: "<="},
		{Key: "bb.currency_type", Val: req.CurrencyType, EqType: "="},
		{Key: "bs.scene_name", Val: req.CurrencyType, EqType: "like"},

		{Key: "acc2.account", Val: req.BusinessAccount, EqType: "="},
		{Key: "b.business_id", Val: req.BusinessId, EqType: "="},
	}

	total, err := dao.BusinessBillRefundDaoInst.GetCnt(whereList)

	//获取列表信息
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY br.finish_time desc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.BusinessBillRefundDaoInst.GetRefundBills(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		if err != sql.ErrNoRows {
			ss_log.Error("退款订单失败，err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYSTEM
			return nil
		}
	}

	var list []*custProto.ChannelRefundBill
	for _, v := range datas {
		data := &custProto.ChannelRefundBill{
			RefundNo:        v.RefundNo,
			OutRefundNo:     v.OutRefundNo,
			RefundStatus:    v.RefundStatus,
			RefundAmount:    v.RefundAmount,
			CurrencyType:    v.CurrencyType,
			FinishTime:      v.FinishTime,
			BusinessAccount: v.BusinessAcc,
			BusinessId:      v.BusinessId,
			BusinessName:    v.BusinessName,
			Subject:         v.Subject,
			TransOrderNo:    v.PayOrderNo,
			TransAmount:     v.TransAmount,
			ChannelName:     v.TradeChannel,
			SceneName:       v.SceneName,
		}
		list = append(list, data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.List = list
	reply.Total = strext.ToInt32(total)
	return nil
}
