package handler

import (
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/common/constants"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/cust-srv/cron"
	"a.a/mp-server/cust-srv/dao"
	"context"
)

const (
	// 重新统计数据
	ReStatisticTypeServicerCheckList = "servicer_check_list" // 生成服务商对账列表
	ReStatisticTypeServicerCount     = "servicer_count"      // 生成服务商对账总计
	ReStatisticTypeUserWithdraw      = "user_withdraw"       // 用户提现
	ReStatisticTypeUserRecharge      = "user_recharge"       // 用户充值
	ReStatisticTypeUserTransfer      = "user_transfer"       // 用户转账
	ReStatisticTypeUserExchange      = "user_exchange"       // 用户兑换
	ReStatisticTypeDate              = "date"                // 按天统计
)

// 获取统计数据-用户提现-用来展示统计图表
func (*CustHandler) GetStatisticUserWithdraw(ctx context.Context, req *go_micro_srv_cust.GetStatisticUserWithdrawRequest, reply *go_micro_srv_cust.GetStatisticUserWithdrawReply) error {
	// 检查日期格式是否正确
	if !ss_time.CheckDateIsRight(req.StartDate) || !ss_time.CheckDateIsRight(req.EndDate) {
		ss_log.Error("参数错误:日期格式错误,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 检查结束时间是否大于等于开始时间
	if cmp, _ := ss_time.CompareDate(ss_time.DateFormat, req.StartDate, req.EndDate); cmp > 0 {
		ss_log.Error("参数错误:开始时间大于结束时间,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 生成日期范围
	dateList, err := ss_time.GetDateRange(ss_time.DateFormat, req.StartDate, req.EndDate)
	if err != nil {
		ss_log.Error("获取日期范围slice失败,err:%v,format:%s,StartDate:%s,EndDate:%s", err, ss_time.DateFormat, req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 初始化各项数据
	dateLen := len(dateList)
	reply.DateList = dateList
	reply.UsdNumList = make([]int64, dateLen)
	reply.UsdAmountList = make([]int64, dateLen)
	reply.UsdFeeList = make([]int64, dateLen)
	reply.KhrNumList = make([]int64, dateLen)
	reply.KhrAmountList = make([]int64, dateLen)
	reply.KhrFeeList = make([]int64, dateLen)

	for _, currencyType := range []string{constants.CURRENCY_USD, constants.CURRENCY_KHR} {
		// 查询统计数据
		list, err := dao.StatisticUserWithdrawDaoInst.GetStatisticData(req.StartDate, req.EndDate, currencyType)
		if err != nil {
			ss_log.Error("查询用户提现统计数据失败,err:%v,StartDate:%s,EndDate:%s, currencyType:%s", err, req.StartDate, req.EndDate, currencyType)
			continue
		}

		// 以生成的日期为准去组织数据
		for i, day := range reply.DateList {
			for _, v := range list {
				if day == v.Day {
					if currencyType == constants.CURRENCY_USD {
						reply.UsdNumList[i] = v.TotalNum
						reply.UsdAmountList[i] = v.TotalAmount
						reply.UsdFeeList[i] = v.TotalFee
					} else if currencyType == constants.CURRENCY_KHR {
						reply.KhrNumList[i] = v.TotalNum
						reply.KhrAmountList[i] = v.TotalAmount
						reply.KhrFeeList[i] = v.TotalFee
					}
					break
				}
			}
		}
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 获取统计数据-用户提现-列表
func (*CustHandler) GetStatisticUserWithdrawList(ctx context.Context, req *go_micro_srv_cust.GetStatisticUserWithdrawListRequest, reply *go_micro_srv_cust.GetStatisticUserWithdrawListReply) error {
	// 检查日期格式是否正确
	if !ss_time.CheckDateIsRight(req.StartDate) || !ss_time.CheckDateIsRight(req.EndDate) {
		ss_log.Error("参数错误:日期格式错误,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 检查结束时间是否大于等于开始时间
	if cmp, _ := ss_time.CompareDate(ss_time.DateFormat, req.StartDate, req.EndDate); cmp > 0 {
		ss_log.Error("参数错误:开始时间大于结束时间,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 检查提现类型
	if req.WithdrawType != "" && !util.InSlice(req.WithdrawType, []string{dao.StatisticWithdrawTypeCard, dao.StatisticRechargeTypeWriteoff, dao.StatisticRechargeTypeScan}) {
		ss_log.Error("参数错误:提现类型错误,WithdrawType:%v", req.WithdrawType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 币种类型
	if req.CurrencyType != "" && !util.InSlice(req.CurrencyType, []string{constants.CURRENCY_USD, constants.CURRENCY_KHR}) {
		ss_log.Error("参数错误:币种类型错误,CurrencyType:%v", req.CurrencyType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = 10
	}

	// 查询列表数据
	dataList, total, err := dao.StatisticUserWithdrawDaoInst.GetStatisticDataList(req)
	if err != nil {
		ss_log.Error("GetStatisticDataList失败,err:%v, req:%+v", err, req)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = dataList
	reply.Total = total
	return nil
}

// 获取统计数据-用户充值-用来展示统计图表
func (*CustHandler) GetStatisticUserRecharge(ctx context.Context, req *go_micro_srv_cust.GetStatisticUserRechargeRequest, reply *go_micro_srv_cust.GetStatisticUserRechargeReply) error {
	// 检查日期格式是否正确
	if !ss_time.CheckDateIsRight(req.StartDate) || !ss_time.CheckDateIsRight(req.EndDate) {
		ss_log.Error("参数错误:日期格式错误,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 检查结束时间是否大于等于开始时间
	if cmp, _ := ss_time.CompareDate(ss_time.DateFormat, req.StartDate, req.EndDate); cmp > 0 {
		ss_log.Error("参数错误:开始时间大于结束时间,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 生成日期范围
	dateList, err := ss_time.GetDateRange(ss_time.DateFormat, req.StartDate, req.EndDate)
	if err != nil {
		ss_log.Error("获取日期范围slice失败,err:%v,format:%s,StartDate:%s,EndDate:%s", err, ss_time.DateFormat, req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 初始化各项数据
	dateLen := len(dateList)
	reply.DateList = dateList
	reply.UsdNumList = make([]int64, dateLen)
	reply.UsdAmountList = make([]int64, dateLen)
	reply.UsdFeeList = make([]int64, dateLen)
	reply.KhrNumList = make([]int64, dateLen)
	reply.KhrAmountList = make([]int64, dateLen)
	reply.KhrFeeList = make([]int64, dateLen)

	for _, currencyType := range []string{constants.CURRENCY_USD, constants.CURRENCY_KHR} {
		// 查询统计数据
		list, err := dao.StatisticUserRechargeDaoInst.GetStatisticData(req.StartDate, req.EndDate, currencyType)
		if err != nil {
			ss_log.Error("查询用户充值统计数据失败,err:%v,StartDate:%s,EndDate:%s, currencyType:%s", err, req.StartDate, req.EndDate, currencyType)
			continue
		}

		// 以生成的日期为准去组织数据
		for i, day := range reply.DateList {
			for _, v := range list {
				if day == v.Day {
					if currencyType == constants.CURRENCY_USD {
						reply.UsdNumList[i] = v.TotalNum
						reply.UsdAmountList[i] = v.TotalAmount
						reply.UsdFeeList[i] = v.TotalFee
					} else if currencyType == constants.CURRENCY_KHR {
						reply.KhrNumList[i] = v.TotalNum
						reply.KhrAmountList[i] = v.TotalAmount
						reply.KhrFeeList[i] = v.TotalFee
					}
					break
				}
			}
		}
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 获取统计数据-用户充值-列表
func (*CustHandler) GetStatisticUserRechargeList(ctx context.Context, req *go_micro_srv_cust.GetStatisticUserRechargeListRequest, reply *go_micro_srv_cust.GetStatisticUserRechargeListReply) error {
	// 检查日期格式是否正确
	if !ss_time.CheckDateIsRight(req.StartDate) || !ss_time.CheckDateIsRight(req.EndDate) {
		ss_log.Error("参数错误:日期格式错误,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 检查结束时间是否大于等于开始时间
	if cmp, _ := ss_time.CompareDate(ss_time.DateFormat, req.StartDate, req.EndDate); cmp > 0 {
		ss_log.Error("参数错误:开始时间大于结束时间,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 检查充值类型
	if req.RechargeType != "" && !util.InSlice(req.RechargeType, []string{dao.StatisticRechargeTypeToHeadquarters, dao.StatisticRechargeTypeToservicer}) {
		ss_log.Error("参数错误:充值类型错误,RechargeType:%v", req.RechargeType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 币种类型
	if req.CurrencyType != "" && !util.InSlice(req.CurrencyType, []string{constants.CURRENCY_USD, constants.CURRENCY_KHR}) {
		ss_log.Error("参数错误:币种类型错误,CurrencyType:%v", req.CurrencyType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = 10
	}

	// 查询列表数据
	dataList, total, err := dao.StatisticUserRechargeDaoInst.GetStatisticDataList(req)
	if err != nil {
		ss_log.Error("GetStatisticDataList失败,err:%v, req:%+v", err, req)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = dataList
	reply.Total = total
	return nil
}

// 获取统计数据-用户转账-用来展示统计图表
func (*CustHandler) GetStatisticUserTransfer(ctx context.Context, req *go_micro_srv_cust.GetStatisticUserTransferRequest, reply *go_micro_srv_cust.GetStatisticUserTransferReply) error {
	// 检查日期格式是否正确
	if !ss_time.CheckDateIsRight(req.StartDate) || !ss_time.CheckDateIsRight(req.EndDate) {
		ss_log.Error("参数错误:日期格式错误,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 检查结束时间是否大于等于开始时间
	if cmp, _ := ss_time.CompareDate(ss_time.DateFormat, req.StartDate, req.EndDate); cmp > 0 {
		ss_log.Error("参数错误:开始时间大于结束时间,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 生成日期范围
	dateList, err := ss_time.GetDateRange(ss_time.DateFormat, req.StartDate, req.EndDate)
	if err != nil {
		ss_log.Error("获取日期范围slice失败,err:%v,format:%s,StartDate:%s,EndDate:%s", err, ss_time.DateFormat, req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 初始化各项数据
	dateLen := len(dateList)
	reply.DateList = dateList
	reply.UsdNumList = make([]int64, dateLen)
	reply.UsdAmountList = make([]int64, dateLen)
	reply.UsdFeeList = make([]int64, dateLen)
	reply.KhrNumList = make([]int64, dateLen)
	reply.KhrAmountList = make([]int64, dateLen)
	reply.KhrFeeList = make([]int64, dateLen)

	for _, currencyType := range []string{constants.CURRENCY_USD, constants.CURRENCY_KHR} {
		// 查询统计数据
		list, err := dao.StatisticUserTransferDaoInst.GetStatisticData(req.StartDate, req.EndDate, currencyType)
		if err != nil {
			ss_log.Error("查询用户充值统计数据失败,err:%v,StartDate:%s,EndDate:%s, currencyType:%s", err, req.StartDate, req.EndDate, currencyType)
			continue
		}

		// 以生成的日期为准去组织数据
		for i, day := range reply.DateList {
			for _, v := range list {
				if day == v.Day {
					if currencyType == constants.CURRENCY_USD {
						reply.UsdNumList[i] = v.TotalNum
						reply.UsdAmountList[i] = v.TotalAmount
						reply.UsdFeeList[i] = v.TotalFee
					} else if currencyType == constants.CURRENCY_KHR {
						reply.KhrNumList[i] = v.TotalNum
						reply.KhrAmountList[i] = v.TotalAmount
						reply.KhrFeeList[i] = v.TotalFee
					}
					break
				}
			}
		}
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 获取统计数据-用户转账-列表
func (*CustHandler) GetStatisticUserTransferList(ctx context.Context, req *go_micro_srv_cust.GetStatisticUserTransferListRequest, reply *go_micro_srv_cust.GetStatisticUserTransferListReply) error {
	// 检查日期格式是否正确
	if !ss_time.CheckDateIsRight(req.StartDate) || !ss_time.CheckDateIsRight(req.EndDate) {
		ss_log.Error("参数错误:日期格式错误,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 检查结束时间是否大于等于开始时间
	if cmp, _ := ss_time.CompareDate(ss_time.DateFormat, req.StartDate, req.EndDate); cmp > 0 {
		ss_log.Error("参数错误:开始时间大于结束时间,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 币种类型
	if req.CurrencyType != "" && !util.InSlice(req.CurrencyType, []string{constants.CURRENCY_USD, constants.CURRENCY_KHR}) {
		ss_log.Error("参数错误:币种类型错误,CurrencyType:%v", req.CurrencyType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = 10
	}

	// 查询列表数据
	dataList, total, err := dao.StatisticUserTransferDaoInst.GetStatisticDataList(req)
	if err != nil {
		ss_log.Error("GetStatisticDataList失败,err:%v, req:%+v", err, req)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = dataList
	reply.Total = total
	return nil
}

// 获取统计数据-用户兑换-用来展示统计图表
func (*CustHandler) GetStatisticUserExchange(ctx context.Context, req *go_micro_srv_cust.GetStatisticUserExchangeRequest, reply *go_micro_srv_cust.GetStatisticUserExchangeReply) error {
	// 检查日期格式是否正确
	if !ss_time.CheckDateIsRight(req.StartDate) || !ss_time.CheckDateIsRight(req.EndDate) {
		ss_log.Error("参数错误:日期格式错误,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 检查结束时间是否大于等于开始时间
	if cmp, _ := ss_time.CompareDate(ss_time.DateFormat, req.StartDate, req.EndDate); cmp > 0 {
		ss_log.Error("参数错误:开始时间大于结束时间,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 生成日期范围
	dateList, err := ss_time.GetDateRange(ss_time.DateFormat, req.StartDate, req.EndDate)
	if err != nil {
		ss_log.Error("获取日期范围slice失败,err:%v,format:%s,StartDate:%s,EndDate:%s", err, ss_time.DateFormat, req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 初始化各项数据
	dateLen := len(dateList)
	reply.DateList = dateList
	reply.Usd2KhrNumList = make([]int64, dateLen)
	reply.Usd2KhrAmountList = make([]int64, dateLen)
	reply.Usd2KhrFeeList = make([]int64, dateLen)
	reply.Khr2UsdNumList = make([]int64, dateLen)
	reply.Khr2UsdAmountList = make([]int64, dateLen)
	reply.Khr2UsdFeeList = make([]int64, dateLen)

	// 查询统计数据
	list, err := dao.StatisticUserExchangeDaoInst.GetStatisticData(req.StartDate, req.EndDate)
	if err != nil {
		ss_log.Error("查询用户兑换统计数据失败,err:%v,StartDate:%s,EndDate:%s", err, req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	// 以生成的日期为准去组织数据
	for i, day := range reply.DateList {
		for _, v := range list {
			if day == v.Day {
				reply.Usd2KhrNumList[i] = v.Usd2khrNum
				reply.Usd2KhrAmountList[i] = v.Usd2khrAmount
				reply.Usd2KhrFeeList[i] = v.Usd2khrFee
				reply.Khr2UsdNumList[i] = v.Khr2usdNum
				reply.Khr2UsdAmountList[i] = v.Khr2usdAmount
				reply.Khr2UsdFeeList[i] = v.Khr2usdFee
				break
			}
		}
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 获取统计数据-用户兑换-列表
func (*CustHandler) GetStatisticUserExchangeList(ctx context.Context, req *go_micro_srv_cust.GetStatisticUserExchangeListRequest, reply *go_micro_srv_cust.GetStatisticUserExchangeListReply) error {
	// 检查日期格式是否正确
	if !ss_time.CheckDateIsRight(req.StartDate) || !ss_time.CheckDateIsRight(req.EndDate) {
		ss_log.Error("参数错误:日期格式错误,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 检查结束时间是否大于等于开始时间
	if cmp, _ := ss_time.CompareDate(ss_time.DateFormat, req.StartDate, req.EndDate); cmp > 0 {
		ss_log.Error("参数错误:开始时间大于结束时间,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = 10
	}

	// 查询列表数据
	dataList, total, err := dao.StatisticUserExchangeDaoInst.GetStatisticDataList(req)
	if err != nil {
		ss_log.Error("GetStatisticDataList失败,err:%v, req:%+v", err, req)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = dataList
	reply.Total = total
	return nil
}

// 获取按天统计数据-用来展示统计图表
func (*CustHandler) GetStatisticDate(ctx context.Context, req *go_micro_srv_cust.GetStatisticDateRequest, reply *go_micro_srv_cust.GetStatisticDateReply) error {
	// 检查日期格式是否正确
	if !ss_time.CheckDateIsRight(req.StartDate) || !ss_time.CheckDateIsRight(req.EndDate) {
		ss_log.Error("参数错误:日期格式错误,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 检查结束时间是否大于等于开始时间
	if cmp, _ := ss_time.CompareDate(ss_time.DateFormat, req.StartDate, req.EndDate); cmp > 0 {
		ss_log.Error("参数错误:开始时间大于结束时间,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 生成日期范围
	dateList, err := ss_time.GetDateRange(ss_time.DateFormat, req.StartDate, req.EndDate)
	if err != nil {
		ss_log.Error("获取日期范围slice失败,err:%v,format:%s,StartDate:%s,EndDate:%s", err, ss_time.DateFormat, req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 初始化各项数据
	dateLen := len(dateList)
	reply.DateList = dateList
	reply.RegUserNumList = make([]int64, dateLen)
	reply.RegServicerNumList = make([]int64, dateLen)

	// 查询统计数据
	list, err := dao.StatisticDateDaoInst.GetStatisticData(req.StartDate, req.EndDate)
	if err != nil {
		ss_log.Error("查询按天统计数据失败,err:%v,StartDate:%s,EndDate:%s", err, req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	// 以生成的日期为准去组织数据
	for i, day := range reply.DateList {
		for _, v := range list {
			if day == v.Day {
				reply.RegUserNumList[i] = v.RegUserNum
				reply.RegServicerNumList[i] = v.RegServicerNum
				break
			}
		}
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 获取按天统计数据-列表
func (*CustHandler) GetStatisticDateList(ctx context.Context, req *go_micro_srv_cust.GetStatisticDateListRequest, reply *go_micro_srv_cust.GetStatisticDateListReply) error {
	// 检查日期格式是否正确
	if !ss_time.CheckDateIsRight(req.StartDate) || !ss_time.CheckDateIsRight(req.EndDate) {
		ss_log.Error("参数错误:日期格式错误,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 检查结束时间是否大于等于开始时间
	if cmp, _ := ss_time.CompareDate(ss_time.DateFormat, req.StartDate, req.EndDate); cmp > 0 {
		ss_log.Error("参数错误:开始时间大于结束时间,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = 10
	}

	// 查询列表数据
	dataList, total, err := dao.StatisticDateDaoInst.GetStatisticDataList(req)
	if err != nil {
		ss_log.Error("GetStatisticDataList失败,err:%v, req:%+v", err, req)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = dataList
	reply.Total = total
	return nil
}

// 对统计数据进行重新统计
func (*CustHandler) ReStatistic(ctx context.Context, req *go_micro_srv_cust.ReStatisticRequest, reply *go_micro_srv_cust.ReStatisticReply) error {
	// 检查日期格式是否正确
	if !ss_time.CheckDateIsRight(req.StartDate) || !ss_time.CheckDateIsRight(req.EndDate) {
		ss_log.Error("参数错误:日期格式错误,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 检查结束时间是否大于等于开始时间
	if cmp, _ := ss_time.CompareDate(ss_time.DateFormat, req.StartDate, req.EndDate); cmp > 0 {
		ss_log.Error("参数错误:开始时间大于结束时间,StartDate:%s,EndDate:%s", req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 生成日期范围
	dateList, err := ss_time.GetDateRange(ss_time.DateFormat, req.StartDate, req.EndDate)
	if err != nil {
		ss_log.Error("获取日期范围slice失败,err:%v,format:%s,StartDate:%s,EndDate:%s", err, ss_time.DateFormat, req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if len(dateList) > 31 {
		ss_log.Error("参数错误,日期范围太大,len:%d,StartDate:%s,EndDate:%s", len(dateList), req.StartDate, req.EndDate)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	switch req.Type {
	case ReStatisticTypeServicerCheckList: // 生成服务商对账列表
		servicerCheckList := &cron.ServicerCheckList{cron.CronBase{LogCat: "重新统计服务商对账列表:"}}
		go func() {
			for _, date := range dateList {
				servicerCheckList.HandleByDate(date)
			}
		}()
	//case ReStatisticTypeServicerCount: // 生成服务商对账总计
	//servicerCount := &cron.ServicerCount{cron.CronBase{LogCat: "重新统计服务商对账总计:"}}
	//go func() {
	//	for _, date := range dateList {
	//		servicerCount.HandleByDate(date)
	//	}
	//}()
	case ReStatisticTypeUserWithdraw: // 用户提现
		withdrawCount := &cron.WithdrawCount{cron.CronBase{LogCat: "重新统计用户提现数据:"}}
		go func() {
			for _, date := range dateList {
				withdrawCount.HandleByDate(date)
			}
		}()
	case ReStatisticTypeUserRecharge: // 用户充值
		saveCount := &cron.SaveCount{cron.CronBase{LogCat: "重新统计用户充值数据:"}}
		go func() {
			for _, date := range dateList {
				saveCount.HandleByDate(date)
			}
		}()
	case ReStatisticTypeUserTransfer: // 用户转账
		transferCount := &cron.TransferCount{cron.CronBase{LogCat: "重新统计用户转账数据:"}}
		go func() {
			for _, date := range dateList {
				transferCount.HandleByDate(date)
			}
		}()
	case ReStatisticTypeUserExchange: // 用户兑换
		exchangeCount := &cron.ExchangeCount{cron.CronBase{LogCat: "重新统计用户兑换数据:"}}
		go func() {
			for _, date := range dateList {
				exchangeCount.HandleByDate(date)
			}
		}()
	case ReStatisticTypeDate: // 按天统计
		regCount := &cron.RegCount{cron.CronBase{LogCat: "重新统计用户注册和服务商注册的数据:"}}
		go func() {
			for _, date := range dateList {
				regCount.HandleByDate(date)
			}
		}()
	default:
		ss_log.Error("参数错误:类型不支持,Type:%s", req.Type)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 获取统计数据-用户资金总留存-用来展示统计图表
func (*CustHandler) GetStatisticUserMoney(ctx context.Context, req *go_micro_srv_cust.GetStatisticUserMoneyRequest, reply *go_micro_srv_cust.GetStatisticUserMoneyReply) error {
	// 查询统计数据
	list, err := dao.StatisticUserMoneyDaoInst.GetStatisticUserMoneyTime(req.StartTime, req.EndTime)
	if err != nil {
		ss_log.Error("查询用户兑换统计数据失败,err:%v,StartDate:%s,EndDate:%s", err, req.StartTime, req.EndTime)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	// 生成时间范围
	timeList, err := ss_time.GetTimeRange("2006/01/02 15:04:05", req.StartTime, req.EndTime)
	if err != nil {
		ss_log.Error("获取时间范围slice失败,err:%v,format:%s,StartDate:%s,EndDate:%s", err, "2006/01/02 15:04:05", req.StartTime, req.EndTime)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	dateLen := len(timeList)
	reply.TimeList = timeList
	reply.UsdBalanceList = make([]int64, dateLen)
	reply.KhrBalanceList = make([]int64, dateLen)
	reply.UsdFrozenBalanceList = make([]int64, dateLen)
	reply.KhrFrozenBalanceList = make([]int64, dateLen)

	// 以生成的日期为准去组织数据
	for i, time := range reply.TimeList {
		for _, v := range list {
			if time == v.CreateTime {
				reply.UsdBalanceList[i] = v.UserUseBalance
				reply.KhrBalanceList[i] = v.UserKhrBalance
				reply.UsdFrozenBalanceList[i] = v.UserUseFrozenBalance
				reply.KhrFrozenBalanceList[i] = v.UserKhrFrozenBalance
				break
			}
		}
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 获取统计数据-用户资金总留存-列表
func (*CustHandler) GetStatisticUserMoneyList(ctx context.Context, req *go_micro_srv_cust.GetStatisticUserMoneyListRequest, reply *go_micro_srv_cust.GetStatisticUserMoneyListReply) error {
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = 10
	}

	// 查询列表数据
	dataList, total, err := dao.StatisticUserMoneyDaoInst.GetStatisticUserMoneyTimeList(req)
	if err != nil {
		ss_log.Error("GetStatisticDataList失败,err:%v, req:%+v", err, req)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = dataList
	reply.Total = strext.ToInt32(total)
	return nil
}
