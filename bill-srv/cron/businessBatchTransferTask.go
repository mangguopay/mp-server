package cron

import (
	"a.a/mp-server/bill-srv/common"
	go_micro_srv_push "a.a/mp-server/common/proto/push"
	"a.a/mp-server/common/ss_func"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"a.a/mp-server/bill-srv/util"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/bill-srv/dao"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/model"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	jsoniter "github.com/json-iterator/go"
)

// 定时处理已支付的批量转账批次订单
var BusinessBatchTransferTask = &BusinessBatchTransfer{CronBase{LogCat: "定时任务处理已支付的批量转账批次订单:", LockExpire: time.Hour * 2}}

type BusinessBatchTransfer struct {
	CronBase
}

// 运行定时任务
func (s *BusinessBatchTransfer) Run() {
	if s.Runing { // 正在运行中
		return
	}

	lockKey := GetLockKey(s)
	lockValue := strext.NewUUID()

	// 获取分布式锁
	if !cache.GetDistributedLock(lockKey, lockValue, s.LockExpire) {
		return
	}
	s.Runing = true

	s.Doing()

	// 释放分布式锁
	cache.ReleaseDistributedLock(lockKey, lockValue)
	s.Runing = false
}

func (s *BusinessBatchTransfer) Doing() {
	ss_log.Info(s.LogCat+"开始time:%s", time.Now())

	//查询出要处理的批量转账批次（已支付的）
	whereList := []*model.WhereSqlCond{
		{Key: "status", Val: constants.BusinessBatchTransferOrderStatusPaySuccess, EqType: "="},
	}
	batchDatas, err := dao.BusinessBatchTransferOrderDaoInst.GetBatchOrderList(whereList)
	if err != nil {
		ss_log.Error("查询已支付的批量转账批次订单出错，err=[%v]", err)
		return
	}

	if batchDatas == nil {
		ss_log.Info(s.LogCat+"未发现需要处理的数据，结束time：%s", time.Now())
		return
	}

	//查询商家转账配置
	transferConf, err := dao.GlobalParamDaoInstance.GetBusinessTransferParamValue()
	if err != nil {
		ss_log.Error("查询商家转账配置失败, err=%v", err)
		return
	}

	for _, batchData := range batchDatas {
		s.DoSingleBatchTransfer(batchData, transferConf)
	}

	ss_log.Info(s.LogCat+"结束time：%s", time.Now())
}

//batchData 批次信息
//transferConf 商家转账配置信息
func (s *BusinessBatchTransfer) DoSingleBatchTransfer(batchData *dao.BusinessBatchTransferOrderDao, transferConf *dao.BusinessTransferParamValue) error {
	ss_log.Info(s.LogCat+"开始处理转账批次[%v],time：%s", batchData.BatchNo, time.Now())

	businessAccNo := dao.BusinessDaoInst.GetAccNoByBusinessNo(batchData.BusinessNo)

	//该转账批次未全部生成转账订单的话
	if batchData.GenerateAll == "false" {
		//生成转账订单
		if err := createBusinessTransfer(batchData, transferConf, businessAccNo); err != nil {
			ss_log.Error(s.LogCat+"处理转账批次创建转账订单发生错误err=[%v],BatchNo[%v]", batchData.BatchNo)
			return err
		}

	}

	//处理该转账批次的全部转账订单
	if err := doTransfer(batchData.BatchNo); err != nil && err != noDataErr {
		ss_log.Error(s.LogCat+"处理转账批次,发生错误，err=[%v],BatchNo[%v]", batchData.BatchNo)
		return err
	}

	//检查转账批量，如果转账全部完成则修改转账批次的状态
	if err := checkBatchTransfer(batchData, businessAccNo); err != nil {
		ss_log.Error(s.LogCat+"处理转账批次发生错误，err=[%v],BatchNo[%v]", batchData.BatchNo)
		return err
	}

	ss_log.Info(s.LogCat+" 处理转账批次[%v]结束,time：%s", batchData.BatchNo, time.Now())
	return nil
}

func checkBatchTransfer(batchData *dao.BusinessBatchTransferOrderDao, businessAccNo string) error {
	//如果总数和订单总数对得上,则标记转账订单已经全部生成
	cnt := dao.BusinessTransferOrderDaoInst.GetBatchCnt(batchData.BatchNo)
	if batchData.TotalNumber != "" && batchData.TotalNumber == cnt { //转账订单数量上等于该转账批次内的订单数
		//修改 转账订单是否全部生成订单状态为true
		errUp := dao.BusinessBatchTransferOrderDaoInst.UpdateGenerateAllByBatchNo(batchData.BatchNo, "true")
		if errUp != nil {
			ss_log.Error("修改批量转账订单全部生成状态为true失败,err=[%v]", errUp)
			return errors.New("修改批量转账订单全部生成状态为true失败")
		}

		err := updateBatchTransferStatus(batchData, businessAccNo)
		if err != nil {
			ss_log.Error("修改转账批量状态失败,err=[%v]", err)
			return errors.New("修改转账批量状态失败")
		}

	} else {
		ss_log.Error("转账订单批次BatchNo[%v]仍然有未生成的订单。", batchData.BatchNo)
	}

	return nil
}

//确认转账金额和数量与转账批次的订单金额和数量一致,如果一致则更新转账批次订单的数据(成功和失败的金额、数量)
func updateBatchTransferStatus(batchData *dao.BusinessBatchTransferOrderDao, businessAccNo string) error {
	//确认是否全部已处理完，金额是否对得上
	successRealAmountSum, successFeeSum, successCnt, errSum := dao.BusinessTransferOrderDaoInst.GetTransferOrderSum(batchData.BatchNo, constants.BusinessTransferOrderStatusSuccess)
	if errSum != nil && errSum != sql.ErrNoRows {
		ss_log.Error("查询转账批次[%v]的成功转账金额、手续费统计失败", batchData.BatchNo)
		return errors.New("查询转账批次的成功转账金额、手续费统计失败")
	}
	failRealAmountSum, _, failCnt, errSum2 := dao.BusinessTransferOrderDaoInst.GetTransferOrderSum(batchData.BatchNo, constants.BusinessTransferOrderStatusFail)
	if errSum2 != nil && errSum2 != sql.ErrNoRows {
		ss_log.Error("查询转账批次[%v]的失败转账金额、手续费统计失败", batchData.BatchNo)
		return errors.New("查询转账批次的失败转账金额、手续费统计失败")
	}

	if batchData.TotalNumber != ss_count.Add(successCnt, failCnt) {
		ss_log.Error("仍有转账订单未处理,订单总数[%v],成功笔数[%v],失败笔数[%v]", batchData.TotalNumber, successCnt, failCnt)
		return errors.New("仍有转账订单未处理")
	}

	//批量转账付款金额-成功金额-成功手续费 = 未用到的批量转账金额
	successSum := ss_count.Add(successRealAmountSum, successFeeSum)
	refundAmount := ss_count.Sub(batchData.RealAmount, successSum).String()

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//开启事务
	tx, txErr := dbHandler.BeginTx(context.TODO(), nil)
	if txErr != nil {
		errText := fmt.Sprintf("开启事务失败, err=[%v]", txErr)
		ss_log.Error(errText)
		return errors.New(errText)
	}

	//以下是付款商家批量转账结束了，但有些金额没用到的情况，开始退还商家没用到的付款
	if strext.ToInt(refundAmount) > 0 {
		businessVAccType := global.GetBusinessVAccType(batchData.CurrencyType, true)
		if businessVAccType == 0 { //查询后仍然是0的话，说明是发生了错误，查询不到
			errText := fmt.Sprintf("获取虚帐账号类型失败,CurrencyType[%v]", batchData.CurrencyType)
			ss_log.Error(errText)
			return errors.New(errText)
		}

		fromVAccNo := dao.VaccountDaoInst.GetVaccountNo(businessAccNo, strext.ToInt32(businessVAccType))
		//付款人虚账
		if fromVAccNo == "" { //查询后仍然是空的话，说明是发生了错误，查询不到
			errText := fmt.Sprintf("没有查到商户虚拟账号，BusinessAccNo=%v, CurrencyType=%v", businessAccNo, batchData.CurrencyType)
			ss_log.Error(errText)
			return errors.New(errText)
		}

		if errStr := dao.VaccountDaoInst.ModifyVaccFrozenToBalance(tx, fromVAccNo, refundAmount, batchData.BatchNo, constants.VaReason_BusinessBatchTransferToBusiness, constants.VaOpType_Defreeze_Add); errStr != ss_err.ERR_SUCCESS {
			ss_sql.Rollback(tx)
			errText := fmt.Sprintf("返还付款商户多付批量转账订单金额失败,refundAmount[%v], toVAccNo=[%v], errMsg=[%v]", refundAmount, fromVAccNo, errStr)
			ss_log.Error(errText)
			return errors.New(errText)
		}

	}

	//将批量转账订单的状态由处理中改为完成,并更新成功\失败的金额和数量
	if err := dao.BusinessBatchTransferOrderDaoInst.UpdateOrderStatusSuccess(tx, batchData.BatchNo, successCnt, successRealAmountSum, failCnt, failRealAmountSum); nil != err {
		tx.Rollback()
		errText := fmt.Sprintf("修改订单状态失败，logNo=%v, err=%v", batchData.BatchNo, err)
		ss_log.Error(errText)
		return errors.New(errText)
	}

	tx.Commit()
	return nil
}

func createBusinessTransfer(batchData *dao.BusinessBatchTransferOrderDao, transferConf *dao.BusinessTransferParamValue, businessAccNo string) error {
	var fileDatas []util.FileContentJsonStruct
	err := jsoniter.Unmarshal([]byte(batchData.FileContent), &fileDatas)
	if err != nil {
		ss_log.Error("解析文件json内容出错,err=[%v], BatchNo[%v]", batchData.BatchNo, err)
		return err
	}
	for _, fileData := range fileDatas {
		//插入到转账订单表里，如果有该条记录则不插入
		total, err := dao.BusinessTransferOrderDaoInst.CheckBusinessTransferOrder(batchData.BatchNo, fileData.Row)
		if err != nil {
			ss_log.Error("确认是否有该批次的一个订单失败,BatchNo[%v],Row[%v},err[%v]", batchData.BatchNo, fileData.Row, err)
			continue //直接结束该转账批次的该订单
		}

		if total == "0" {
			d := new(dao.BusinessTransferOrderDao)
			d.FromBusinessNo = batchData.BusinessNo
			d.FromAccountNo = businessAccNo
			d.BatchNo = batchData.BatchNo
			d.BatchRowNum = fileData.Row
			d.Remarks = fileData.Remarks
			d.Rate = "0"
			d.Fee = "0"

			errAdd := insertBusinessTransfer(d, fileData, transferConf)
			if errAdd != nil {
				ss_log.Error("插入转账订单失败，BatchNo[%v],Row[%v},err[%v]", batchData.BatchNo, fileData.Row, errAdd)
				continue //直接结束该转账批次的该订单
			}

		}

	}

	return nil
}

var noDataErr = errors.New("无需要处理的数据")

func doTransfer(batchNo string) (err error) {
	//查询该转账批次的转账订单信息
	whereList2 := []*model.WhereSqlCond{
		{Key: "order_status", Val: constants.BusinessTransferOrderStatusPending, EqType: "="},
		{Key: "batch_no", Val: batchNo, EqType: "="},
	}

	transferDatas, errGet := dao.BusinessTransferOrderDaoInst.GetTransferOrderList(whereList2)
	if errGet != nil {
		ss_log.Error("查询转账批次的转账订单出错,BatchNo", batchNo)
		return errors.New("查询转账批次的转账订单出错") //		结束该批次转账的任务。
	}

	if len(transferDatas) == 0 {
		ss_log.Error("该转账批次[%v]无需要处理的数据", batchNo)
		return noDataErr
	}

	//同一批次的转账发起方、币种是相同的，所以查一遍就可以了
	businessVAccType := global.GetBusinessVAccType(transferDatas[0].CurrencyType, true)
	if businessVAccType == 0 { //查询后仍然是0的话，说明是发生了错误，查询不到
		return errors.New("获取虚帐账号类型失败")
	}

	fromVAccNo := dao.VaccountDaoInst.GetVaccountNo(transferDatas[0].FromAccountNo, strext.ToInt32(businessVAccType))
	//付款人虚账
	if fromVAccNo == "" { //查询后仍然是空的话，说明是发生了错误，查询不到
		ss_log.Error("没有查到商户虚拟账号，BusinessAccNo=%v, CurrencyType=%v", transferDatas[0].FromAccountNo, transferDatas[0].CurrencyType)
		return errors.New("没有查到发起批量转账的商户虚拟账号")
	}

	// 查询总部的账号
	_, headAcc, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")

	plantVaType := global.GetPlatFormVAccType(transferDatas[0].CurrencyType)
	if plantVaType == 0 { //查询后仍然是0的话，说明是发生了错误，查询不到
		return errors.New("获取平台虚帐账号类型失败")
	}

	// 确保平台虚拟账号存在
	headVacc := dao.VaccountDaoInst.ConfirmExistVAccount(headAcc, transferDatas[0].CurrencyType, strext.ToInt32(plantVaType))

	for _, data := range transferDatas {
		toVAccType := 0
		pushMessages := false //是否推送消息(用户要推送)

		//查询该账号是企业账号还是个人账号
		toAccount := dao.AccDaoInstance.GetAccountFromAccNo(data.ToAccountNo)
		if strings.Contains(toAccount, "@") {
			toVAccType = businessVAccType
		} else {
			userVaccType, vaErr1 := common.VirtualAccountTypeByMoneyType(transferDatas[0].CurrencyType, "1")
			if vaErr1 != nil {
				return errors.New("获取用户虚帐账号类型失败")
			}
			pushMessages = true
			toVAccType = userVaccType
		}

		//收款账号虚账
		toVAccNo := dao.VaccountDaoInst.ConfirmExistVAccount(data.ToAccountNo, data.CurrencyType, strext.ToInt32(toVAccType))

		//转账并修改订单状态
		errTransfer := fromVAccNoTransferToVaccNo(fromVAccNo, toVAccNo, headVacc, data.RealAmount, data.Fee, data.CurrencyType, data.LogNo)
		if errTransfer != nil {
			ss_log.Error("errTransfer=%v", errTransfer)
			//return errTransfer
		}

		if pushMessages && errTransfer == nil { //是用户并且是转账成功的就要推送消息
			//添加转账到的账号推送消息
			if errAddMessages2 := dao.LogAppMessagesDaoInst.AddLogAppMessages(data.LogNo, constants.LOG_APP_MESSAGES_ORDER_TYPE_TRANSFER_Apply, constants.VaReason_BusinessTransferToBusiness, data.ToAccountNo, constants.OrderStatus_Paid); errAddMessages2 != ss_err.ERR_SUCCESS {
				ss_log.Error("errAddMessages2=[%v]", errAddMessages2)
			}

			appLang, _ := dao.AccDaoInstance.QueryAccountLang(data.ToAccountNo)
			if appLang == "" {
				appLang = constants.LangEnUS
			}

			ss_log.Info("用户 AccountNo %s 当前的语言为--->%s", data.ToAccountNo, appLang)
			moneyType := dao.LangDaoInstance.GetLangTextByKey(strings.ToLower(data.CurrencyType), appLang)
			timeString := time.Now().Format("2006-01-02 15:04:05")
			// 修正各币种的金额
			amount := common.NormalAmountByMoneyType(strings.ToLower(data.CurrencyType), data.RealAmount)

			args := []string{
				timeString, amount, moneyType,
			}
			if appLang == constants.LangEnUS {
				args = []string{
					amount, moneyType, timeString,
				}
			}

			// 消息推送
			ev := &go_micro_srv_push.PushReqest{
				Accounts: []*go_micro_srv_push.PushAccout{
					{
						AccountNo:   data.ToAccountNo,
						AccountType: constants.AccountType_USER,
					},
				},
				TempNo: constants.Template_TransferSuccess,
				Args:   args,
			}

			ss_log.Info("publishing %+v\n", ev)
			// publish an event
			if err := common.PushEvent.Publish(context.TODO(), ev); err != nil {
				ss_log.Error("消息推送到用户toAccountNo[%v]出错。error : %v", data.ToAccountNo, err)
			}
		}

	}
	return nil
}

//转账并修改订单状态
func fromVAccNoTransferToVaccNo(fromVAccNo, toVAccNo, headVacc, realAmount, fee, currencyType, logNo string) (err error) {
	ss_log.Info("fromVAccNoTransferToVaccNo 开始 fromVAccNo[%s], toVAccNo[%s], headVacc[%s], realAmount[%s], fee[%s], currencyType[%v], logNo[%s]", fromVAccNo, toVAccNo, headVacc, realAmount, fee, currencyType, logNo)

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//开启事务
	tx, txErr := dbHandler.BeginTx(context.TODO(), nil)
	if txErr != nil {
		ss_log.Error("开启事务失败, err=%v", txErr)
		return errors.New("开启事务失败")
	}

	var errR error
	defer func() error {
		if errR != nil { //如果返还时发现发生错误，则改变订单以失败做处理
			ss_log.Error("转账时发生错误,errR=%v", errR)
			// 不必按转账订单金额返还钱商家了。到后面该批次转账完成后会有确定批量转账所有完成的转账订单金额、与数量是否和批量转账订单的数据一致，多付的钱到那时候再返还
			if err := dao.BusinessTransferOrderDaoInst.UpdateOrderStatusByLogNo(logNo, constants.BusinessTransferOrderStatusFail, "转账失败"); nil != err {
				ss_log.Error("修改订单状态为失败出错，logNo=%v, err=%v", logNo, err)
				return errors.New("修改订单状态失败")
			}
		}
		return errR
	}()

	//	付款商家冻结金额-
	if errStr := dao.VaccountDaoInst.ModifyVaccFrozenUpperZero1(tx, "-", fromVAccNo, realAmount, logNo, constants.VaReason_BusinessTransferToBusiness, constants.VaOpType_Defreeze); errStr != ss_err.ERR_SUCCESS {
		ss_sql.Rollback(tx)
		ss_log.Error("付款商户冻结余额减少订单金额失败, fromVAccNo=%v, errMsg=%v", fromVAccNo, errStr)
		errR = errors.New("付款商户冻结余额减少订单金额失败")
		return errR
	}

	if fee != "" && fee != "0" {
		if errStr := dao.VaccountDaoInst.ModifyVaccFrozenUpperZero1(tx, "-", fromVAccNo, fee, logNo, constants.VaReason_FEES, constants.VaOpType_Defreeze); errStr != ss_err.ERR_SUCCESS {
			ss_sql.Rollback(tx)
			ss_log.Error("付款商户冻结余额减少订单手续费失败, fromVAccNo=%v, errMsg=%v", fromVAccNo, errStr)
			errR = errors.New("付款商户冻结余额减少订单手续费失败")
			return errR
		}

		// 修改总部的临时虚账余额
		if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, headVacc, fee, "+", logNo, constants.VaReason_BusinessTransferToBusiness); errStr != ss_err.ERR_SUCCESS {
			ss_sql.Rollback(tx)
			ss_log.Error("修改总部的临时虚账余额失败, headVacc=%v, errMsg=%v", headVacc, errStr)
			errR = errors.New("修改总部的临时虚账余额失败")
			return errR
		}

		// 插入利润表
		d := &dao.HeadquartersProfit{
			GeneralLedgerNo: logNo,
			Amount:          fee,
			OrderStatus:     constants.OrderStatus_Paid,
			BalanceType:     strings.ToLower(currencyType),
			ProfitSource:    constants.ProfitSource_BusinessTransferFee,
			OpType:          constants.PlatformProfitAdd,
		}
		if errStr := dao.HeadquartersProfitDaoInstance.InsertHeadquartersProfit(tx, d); errStr != ss_err.ERR_SUCCESS {
			ss_sql.Rollback(tx)
			ss_log.Error("修改总部的临时虚账余额失败, headVacc=%v, errMsg=%v", headVacc, errStr)
			errR = errors.New("修改总部的临时虚账余额失败")
			return errR
		}

		// 修改收益 总部虚账的余额是等于收益表中的可提现余额
		if errStr := dao.HeadquartersProfitCashableDaoInstance.SyncAccProfit(tx, headVacc, fee, currencyType); errStr != ss_err.ERR_SUCCESS {
			ss_sql.Rollback(tx)
			ss_log.Error("修改收益余额失败, headVacc=%v, errMsg=%v", headVacc, errStr)
			errR = errors.New("修改收益余额失败")
			return errR
		}
	}

	//收款虚帐账户余额+
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, toVAccNo, realAmount, "+", logNo, constants.VaReason_BusinessTransferToBusiness); errStr != ss_err.ERR_SUCCESS {
		ss_sql.Rollback(tx)
		ss_log.Error("增加收款商户余额失败, toVAccNo=%v, errMsg=%v", toVAccNo, errStr)
		errR = errors.New("增加收款商户余额失败")
		return errR
	}

	//走到这里来，说明没出任何问题，开始修改订单的状态为成功
	if err := dao.BusinessTransferOrderDaoInst.UpdateOrderStatusByLogNoTx(tx, logNo, constants.BusinessTransferOrderStatusSuccess, ""); nil != err {
		ss_sql.Rollback(tx)
		ss_log.Error("修改订单状态失败，logNo=%v, err=%v", logNo, err)
		errR = errors.New("修改订单状态失败")
		return errR
	}

	ss_sql.Commit(tx)
	ss_log.Info("fromVAccNoTransferToVaccNo结束 ")
	return nil
}

func insertBusinessTransfer(data *dao.BusinessTransferOrderDao, fileData util.FileContentJsonStruct, transferConf *dao.BusinessTransferParamValue) (err error) {
	wrongReason := "" //发生错误的原因
	orderStatus := constants.BusinessTransferOrderStatusPending

	//查询账号的uid和商家uid,并校验
	toAccUid, businessNo, wrongReasonT := util.CheckToAccountAndAuthName(fileData.ToAccount, fileData.Name, data.FromAccountNo)
	if wrongReasonT != "" {
		wrongReason = wrongReasonT
		orderStatus = constants.BusinessTransferOrderStatusFail
	}

	//校验金额
	amountWrongReason := util.CheckTransferAmount(fileData.Amount, fileData.CurrencyType, transferConf)
	if wrongReason == "" && amountWrongReason != "" {
		wrongReason = amountWrongReason
		orderStatus = constants.BusinessTransferOrderStatusFail
	}

	//计算手续费
	fee, rate, wrongReasonT := util.QueryTransferFeeAndRate(fileData.Amount, fileData.CurrencyType, transferConf)
	if wrongReason == "" && wrongReasonT != "" {
		wrongReason = wrongReasonT
		orderStatus = constants.BusinessTransferOrderStatusFail
	}

	toAccount := ""
	if strings.Contains(fileData.ToAccount, "@") { //只有企业商家的账号有@
		toAccount = fileData.ToAccount
	} else { //todo 其他情况视为转账给个人
		accountArr := strings.Split(fileData.ToAccount, "-")

		if len(accountArr) == 2 {
			//处理国家码将其变成前缀无0的格式
			countryCode := strext.ToString(strext.ToInt(accountArr[0]))
			account := ss_func.PrePhone(countryCode, accountArr[1])

			//将国家码变成0086、0855的格式，组成账号
			toAccount = ss_func.PreCountryCode(countryCode) + account
		} else {
			toAccount = fileData.ToAccount
		}

	}

	//插入转账日志
	data.ToBusinessNo = businessNo
	data.ToAccountNo = toAccUid
	data.Amount = fileData.Amount
	data.RealAmount = fileData.Amount //外扣
	data.Fee = fee
	data.Rate = rate
	data.CurrencyType = fileData.CurrencyType
	data.AuthName = fileData.Name
	data.OrderStatus = orderStatus
	data.WrongReason = wrongReason
	data.PaymentType = constants.ORDER_PAYMENT_BALANCE
	data.ToAccount = toAccount
	data.TransferType = constants.BusinessTransferOrderTypeOrdinary
	_, errAdd := dao.BusinessTransferOrderDaoInst.Insert2(data)
	if errAdd != nil {
		ss_log.Error("插入转账订单表发生错误,err[%v]", errAdd)
		return errAdd
	}

	return nil
}
