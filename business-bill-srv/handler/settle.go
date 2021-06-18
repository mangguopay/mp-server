package handler

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/business-bill-srv/dao"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
)

type DisPoseAmountReq struct {
	SettleId        string
	BusinessNo      string
	BusinessAccNo   string
	CurrencyType    string
	TotalAmount     int64
	TotalRealAmount int64
	TotalFees       int64
}

type VAccountType struct {
	AccPlatVaType           int
	BusinessUnSettledVaType int
	BusinessSettledVaType   int
}

//手动结算函数
func (b *BusinessBillHandler) ManualSettle(ctx context.Context, req *businessBillProto.ManualSettleRequest, reply *businessBillProto.ManualSettleReply) error {
	if req.OrderNos == nil {
		ss_log.Error("结算订单列表参数为空")
		reply.ResultCode = ss_err.OrderNoIsEmpty
		return nil
	}

	//结算失败订单数量，订单号-失败原因
	var failNum int64 = 0
	failOrder := make(map[string]string)

	parallelNum := 100 // 并发执行数量
	var wg sync.WaitGroup
	for i := 0; i < len(req.OrderNos); i++ {
		if i > 0 && i%parallelNum == 0 {
			ss_log.Info("并发等待----------------------------%v", i)
			wg.Wait()
		}
		wg.Add(1)
		go func(orderNo string) {
			defer wg.Done()
			_, channelType, err := dao.BusinessBillDaoInst.GetOrderChannelNo(orderNo)
			if err != nil {
				ss_log.Error("查询订单支付渠道失败, orderNo=%v, err=%v", orderNo, err)
				failNum++
				failOrder[orderNo] = ss_err.QueryOrderPaymentChannel
				return
			}

			if channelType == "" {
				ss_log.Error("订单支付渠道为空, orderNo=%v", orderNo)
				failNum++
				failOrder[orderNo] = ss_err.QueryOrderPaymentChannel
				return
			}

			if channelType == constants.ChannelTypeInner {
				ss_log.Error("平台内部交易订单暂不支持手动结算, orderNo=%v", orderNo)
				failNum++
				failOrder[orderNo] = ss_err.OrderNotSupportedManualSettle
				return
			}

			_, errCode := b.SingleOrderSettle(orderNo)
			if errCode != ss_err.Success {
				ss_log.Error("结算失败, orderNo=%v, err=%v", orderNo, errCode)
				failNum++
				failOrder[orderNo] = errCode
				return
			}
		}(req.OrderNos[i])
	}
	wg.Wait()

	reply.ResultCode = ss_err.Success
	reply.FailNum = strext.ToInt64(len(failOrder))
	reply.FailOrder = failOrder
	return nil
}

func (b *BusinessBillHandler) SingleOrderSettle(orderNo string) (id, errCode string) {
	order, err := dao.BusinessBillDaoInst.GetOrderInfoByOrderNo(orderNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("订单不存在, orderNo=%v, err=%v", orderNo, err)
			return "", ss_err.OrderNotExist
		}
		ss_log.Error("查询订单失败, orderNo=%v, err=%v", orderNo, err)
		return "", ss_err.SystemErr
	}

	if !util.InSlice(order.OrderStatus, []string{constants.BusinessOrderStatusPay, constants.BusinessOrderStatusRebatesRefund}) {
		ss_log.Error("订单[%v]当前状态[%v]不能结算", orderNo, order.OrderStatus)
		return "", ss_err.OrderStatusUnknown
	}

	var upFee string
	switch order.TradeType {
	case constants.TradeTypeWeiXinFaceToFace:
		fallthrough
	case constants.TradeTypeWeiXinApp:
		//微信 todo 新版产品签约
		upRate, err := dao.BusinessSceneSignedDaoInst.GetBusinessUpstreamRate(order.BusinessNo, order.SceneNo)
		if err != nil {
			ss_log.Error("查询上游费率失败, err=%v", err)
			return "", ss_err.SystemErr
		}
		upFee = ss_count.CountFees(order.Amount, upRate.UpWxRate, "0").String()

	case constants.TradeTypeAlipayFaceToFace:
		fallthrough
	case constants.TradeTypeAlipayAPP:
		//支付宝 todo 新版产品签约
		upRate, err := dao.BusinessSceneSignedDaoInst.GetBusinessUpstreamRate(order.BusinessNo, order.SceneNo)
		if err != nil {
			ss_log.Error("查询上游费率失败, err=%v", err)
			return "", ss_err.SystemErr
		}
		upFee = ss_count.CountFees(order.Amount, upRate.UpAliPayRate, "0").String()

	case constants.TradeTypeModernpayMWEB:
		fallthrough
	case constants.TradeTypeModernpayAPP:
		fallthrough
	case constants.TradeTypeModernpayFaceToFace:
		fallthrough
	case constants.TradeTypeBusinessPay:
		//ModernPay
		upFee = "0"

	default:
		ss_log.Error("交易类型不匹配")
		return "", ss_err.TradeTypeValueIsIllegality
	}

	platformFee := strext.ToInt64(ss_count.Sub(order.Fee, upFee).String())
	if platformFee < 0 {
		ss_log.Error("平台收益费率配置有误，orderNo=%v, Fee=%v, upFee=%v", orderNo, order.Fee, upFee)
		return "", ss_err.OrderSettleFail
	}

	req := &DisPoseAmountReq{
		BusinessNo:      order.BusinessNo,
		BusinessAccNo:   order.BusinessAccountNo,
		CurrencyType:    order.CurrencyType,
		TotalAmount:     strext.ToInt64(order.Amount),
		TotalRealAmount: strext.ToInt64(order.RealAmount),
		TotalFees:       platformFee,
	}

	//获取数据库连接
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		ss_log.Error("获取数据库连接失败")
		return "", ss_err.SystemErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//开启事务
	tx, err := dbHandler.Begin()
	if err != nil {
		ss_log.Error("开启事务失败, err:%v", err)
		return "", ss_err.SystemErr
	}
	d := new(dao.BusinessSettleOneDao)
	settleId, err := dao.BusinessSettleOneDaoInst.InsertTx(tx, d)
	if err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("插入订单结算记录失败，err=%v", err)
		return "", ss_err.SystemErr
	}

	req.SettleId = settleId
	disAmountErr := DisPoseAmount(tx, "", req)
	if disAmountErr != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("结算失败，req=%v, err=%v", strext.ToJson(req), disAmountErr)
		return "", ss_err.OrderSettleFail
	}

	//修改订单结算批次
	err = dao.BusinessBillDaoInst.UpdateOrderSettleIdTx(tx, settleId, orderNo)
	if err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("添加订单settle_id失败, orderNo=%v, err=%v", orderNo, err)
		return "", ss_err.SystemErr
	}

	ss_sql.Commit(tx)

	return settleId, ss_err.Success
}

//处理相关金额 logCat: 日志前缀
func DisPoseAmount(tx *sql.Tx, logCat string, req *DisPoseAmountReq) error {
	if req.BusinessAccNo == "" {
		businessAccNo, err := dao.BusinessDaoInst.QueryAccNoByBusinessNo(req.BusinessNo)
		if err != nil {
			ss_log.Error(logCat+"查询商户账号失败,businessNo:%v, err:%v", req.BusinessNo, err)
			return err
		}
		req.BusinessAccNo = businessAccNo
	}

	//获取商户账号类型(未结算,已结算)
	vAccType := GetVAccType(req.CurrencyType)

	//=================================================================================
	//减少商户未结算虚账金额
	busUnSettleVAccNo, err := dao.VaccountDaoInst.GetVaccountNo(req.BusinessAccNo, vAccType.BusinessUnSettledVaType)
	if err != nil {
		ss_log.Error(logCat+"查询商户未结算虚账失败,businessNo:%v, currencyType:%v, err:%v", req.BusinessNo, req.CurrencyType, err)
		return err
	}
	balance1, frozenBalance1, err1 := dao.VaccountDaoInst.MinusBalance(tx, busUnSettleVAccNo, strext.ToString(req.TotalAmount))
	if err1 != nil {
		ss_log.Error(logCat+"减少商户未结算金额失败(美元), businessNo:%v, amount:%v, err:%v", req.BusinessNo, req.TotalAmount, err1)
		return err1
	}
	//正常情况应该是大于0或者等于0
	//如果未结算账户余额减去结算金额后小于0,也许是哪里有问题,需要检查
	if strext.ToInt(balance1) < 0 {
		return errors.New(fmt.Sprintf("未结算账户金额(%v)有误,商户:%v,应扣金额:%v,扣除后余额:%v", req.CurrencyType, req.BusinessNo, req.TotalAmount, balance1))
	}
	//记录商户账户变动日志
	log1 := dao.LogVaccountDao{
		VaccountNo:    busUnSettleVAccNo,
		OpType:        constants.VaOpType_Minus,
		Amount:        strext.ToString(req.TotalAmount),
		Balance:       balance1,
		FrozenBalance: frozenBalance1,
		Reason:        constants.VaReason_Business_Settle,
		BizLogNo:      req.SettleId,
	}
	if err := dao.LogVaccountDaoInst.InsertLogTx(tx, log1); err != nil {
		ss_log.Error(logCat+"插入商户账户变动日志失败,err:%v,data:%+v", err, log1)
		return err
	}

	//=================================================================================
	//增加商户已结算账户余额
	busSettledVAccNo, err := dao.VaccountDaoInst.GetVaccountNo(req.BusinessAccNo, vAccType.BusinessSettledVaType)
	if err != nil {
		ss_log.Error(logCat+"查询商户已结算虚账失败,businessNo:%v, vaType:%v,err:%v", req.BusinessNo, vAccType.BusinessSettledVaType, err)
		return err
	}
	balance2, frozenBalance2, err2 := dao.VaccountDaoInst.PlusBalance(tx, busSettledVAccNo, strext.ToString(req.TotalRealAmount))
	if err2 != nil {
		ss_log.Error(logCat+"增加商户已结算金额失败(%v), businessNo:%v, amount:%v, err:%v", req.CurrencyType, req.BusinessNo, req.TotalRealAmount, err2)
		return err
	}
	//记录商户账户变动日志
	log2 := dao.LogVaccountDao{
		VaccountNo:    busSettledVAccNo,
		OpType:        constants.VaOpType_Add,
		Amount:        strext.ToString(req.TotalRealAmount),
		Balance:       balance2,
		FrozenBalance: frozenBalance2,
		Reason:        constants.VaReason_Business_Settle,
		BizLogNo:      req.SettleId,
	}
	if err := dao.LogVaccountDaoInst.InsertLogTx(tx, log2); err != nil {
		ss_log.Error(logCat+"插入商户账户变动日志失败,err:%v,data:%+v", err, log2)
		return err
	}

	//=================================================================================
	//平台盈利，增加平台虚账金额
	if req.TotalFees > 0 {
		_, accPlat, err := cache.ApiDaoInstance.GetGlobalParam(constants.GlobalParamKeyAccPlat)
		if err != nil {
			ss_log.Error(logCat+"查询平台账号失败,paramKey:%v, err:%v", constants.GlobalParamKeyAccPlat, err)
			return err
		}
		vAccPlat, err := dao.VaccountDaoInst.GetVaccountNo(accPlat, vAccType.AccPlatVaType)
		if err != nil {
			ss_log.Error(logCat+"查询平台虚账失败,accountNo:%v, err:%v", accPlat, err)
			return err
		}
		balance3, frozenBalance3, err3 := dao.VaccountDaoInst.PlusBalance(tx, vAccPlat, strext.ToString(req.TotalFees))
		if err3 != nil {
			ss_log.Error(logCat+"增加平台手续费失败(%v), account:%v, fee:%v, err:%v", req.CurrencyType, accPlat, req.TotalFees, err3)
			return err
		}
		//记录商户账户变动日志
		log3 := dao.LogVaccountDao{
			VaccountNo:    vAccPlat,
			OpType:        constants.VaOpType_Add,
			Amount:        strext.ToString(req.TotalFees),
			Balance:       balance3,
			FrozenBalance: frozenBalance3,
			Reason:        constants.VaReason_Business_Settle,
			BizLogNo:      req.SettleId,
		}
		if err := dao.LogVaccountDaoInst.InsertLogTx(tx, log3); err != nil {
			ss_log.Error(logCat+"插入平台手续费变动日志失败,err:%v,data:%+v", err, log3)
			return err
		}

		//插入手续费盈利记录
		d := &dao.HeadquartersProfit{
			GeneralLedgerNo: req.SettleId,
			Amount:          strext.ToString(req.TotalFees),
			OrderStatus:     constants.OrderStatus_Paid,
			BalanceType:     strings.ToLower(req.CurrencyType),
			ProfitSource:    constants.ProfitSource_ModernPayOrderFee,
			OpType:          constants.PlatformProfitAdd,
		}
		_, err = dao.HeadquartersProfitDao.InsertHeadquartersProfit(tx, d)
		if err != nil {
			ss_log.Error("插入平台盈利失败, data=%v, err=%v", strext.ToJson(d), err)
			return err
		}

		//同步总部虚账的余额(等于收益表中的可提现余额)
		err = dao.HeadquartersProfitDao.SyncHeadquartersProfit(tx, vAccPlat, strext.ToString(req.TotalFees), strings.ToLower(req.CurrencyType))
		if err != nil {
			ss_log.Error("同步收益余额失败, headVacc=%v, err=%v", vAccPlat, err)
			return err
		}
	}
	return nil
}

// 通过币种类型获商户的虚拟账户类型
func GetVAccType(currencyType string) *VAccountType {
	vAccountType := new(VAccountType)
	switch currencyType {
	case constants.CURRENCY_UP_USD:
		vAccountType.AccPlatVaType = constants.VaType_USD_FEES
		vAccountType.BusinessUnSettledVaType = constants.VaType_USD_BUSINESS_UNSETTLED
		vAccountType.BusinessSettledVaType = constants.VaType_USD_BUSINESS_SETTLED
	case constants.CURRENCY_UP_KHR:
		vAccountType.AccPlatVaType = constants.VaType_KHR_FEES
		vAccountType.BusinessUnSettledVaType = constants.VaType_KHR_BUSINESS_UNSETTLED
		vAccountType.BusinessSettledVaType = constants.VaType_KHR_BUSINESS_SETTLED
	}
	return vAccountType
}
