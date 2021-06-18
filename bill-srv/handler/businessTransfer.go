package handler

import (
	"context"
	"database/sql"
	"io/ioutil"
	"strings"
	"time"

	go_micro_srv_push "a.a/mp-server/common/proto/push"
	"a.a/mp-server/common/ss_func"

	"a.a/cu/db"
	"a.a/cu/ss_big"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/bill-srv/common"
	"a.a/mp-server/bill-srv/cron"
	"a.a/mp-server/bill-srv/dao"
	"a.a/mp-server/bill-srv/i"
	"a.a/mp-server/bill-srv/util"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	authProto "a.a/mp-server/common/proto/auth"
	billProto "a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"github.com/shopspring/decimal"
	"github.com/tealeg/xlsx"
)

//商家转账给企业商家
func transferToBusiness(outBusiness *dao.BusinessStatus, req *billProto.AddBusinessTransferRequest) (codeT, logNoT, orderStatusT string) {
	//查询收款商家的状态和收款权限
	inBusiness, err := dao.BusinessDaoInst.GetBusinessStatusInfo(req.PayeeNo, "")
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("收款商家不存在，PayeeNo=%v, err=%v", req.PayeeNo, err)
			return ss_err.ERR_PayeeNotExist, "", ""
		}
		ss_log.Error("查询付款商家状态失败，PayeeNo=%v, err=%v", req.PayeeNo, err)
		return ss_err.ERR_SYSTEM, "", ""

	}
	if strext.ToInt(inBusiness.AccountStatus) != constants.AccountUseStatusNormal {
		ss_log.Error("收款商家账号已被禁用，PayeeNo=%v", req.PayeeNo)
		return ss_err.ERR_MERC_NO_USE, "", ""

	}
	if strext.ToInt(inBusiness.BusinessStatus) == constants.BusinessUseStatusDisabled {
		ss_log.Error("收款商家已被禁用，PayeeNo=%v", req.PayeeNo)
		return ss_err.ERR_MERC_NO_USE, "", ""

	}
	if strext.ToInt(inBusiness.IncomeAuthorization) == constants.BusinessIncomeAuthDisabled {
		ss_log.Error("收款商家没有收款权限，PayeeNo=%v", req.PayeeNo)
		return ss_err.ERR_PAY_NO_IN_COME_PERMISSION, "", ""
	}

	//计算商家手续费
	fee, rate, getFeeCode := getBusinessFee(req.Amount, req.CurrencyType)
	if getFeeCode != ss_err.ERR_SUCCESS {
		return getFeeCode, "", ""
	}

	//判断付款商余额是否足够
	businessVAccType := global.GetBusinessVAccType(req.CurrencyType, true)
	err, balance := dao.VaccountDaoInst.GetBalanceFromAccNo(req.BusinessAccNo, businessVAccType)
	if err != nil {
		ss_log.Error("查询商家虚账失败, err=%v", err)
		return ss_err.ERR_PARAM, "", ""
	}
	if strext.ToInt(balance) < strext.ToInt(ss_count.Add(req.Amount, fee)) {
		ss_log.Error("商家余额不足，balance=%v, ss_count.Add(req.Amount, fee)=%v,", balance, ss_count.Add(req.Amount, fee))
		ss_log.Error("商家余额不足，BusinessAccNo=%v, CurrencyType=%v,", req.BusinessAccNo, req.CurrencyType)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH, "", ""
	}

	//付款人虚账
	payerVAccNo := dao.VaccountDaoInst.GetVaccountNo(req.BusinessAccNo, strext.ToInt32(businessVAccType))
	if payerVAccNo == "" {
		ss_log.Error("没有查到商户虚拟账号，BusinessAccNo=%v, CurrencyType=%v", req.BusinessAccNo, req.CurrencyType)
		return ss_err.ERR_SYSTEM, "", ""

	}
	//收款商家虚账,防止商家没有初始化此类型的虚账，这里做一下处理
	payeeVAccNo := InternalCallHandlerInst.ConfirmExistVAccount(inBusiness.AccountNo, req.CurrencyType, strext.ToInt32(businessVAccType))

	//插入转账日志
	d := new(dao.BusinessTransferOrderDao)
	d.FromBusinessNo = outBusiness.BusinessNo
	d.FromAccountNo = outBusiness.AccountNo
	d.ToBusinessNo = inBusiness.BusinessNo
	d.ToAccountNo = inBusiness.AccountNo
	d.Amount = req.Amount
	d.RealAmount = req.Amount
	d.Fee = fee
	d.Rate = rate
	d.CurrencyType = req.CurrencyType
	d.OrderStatus = constants.BusinessTransferOrderStatusPending
	d.Remarks = req.Remarks
	d.TransferType = constants.BusinessTransferOrderTypeOrdinary
	logNo, err := dao.BusinessTransferOrderDaoInst.Insert(d)
	if err != nil {
		ss_log.Error("插入商家转账日志失败，data:%v, err=%v", strext.ToJson(d), err)
		return ss_err.ERR_SYSTEM, "", ""
	}

	//开启事务
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, txErr := dbHandler.BeginTx(context.TODO(), nil)
	if txErr != nil {
		ss_log.Error("开启事务失败, err=%v", txErr)
		return ss_err.ERR_SYSTEM, "", ""
	}

	orderStatus := constants.BusinessTransferOrderStatusSuccess

	//如果没有转账成功则修改订单状态为失败
	defer func() error {
		if err := dao.BusinessTransferOrderDaoInst.UpdateOrderStatusByLogNo(logNo, orderStatus, ""); nil != err {
			//ss_sql.Rollback(tx)
			ss_log.Error("修改订单状态失败，logNo=%v, err=%v", logNo, err)
			//reply.ResultCode = ss_err.ERR_SYSTEM
			return err
		}
		//ss_sql.Commit(tx)
		return nil
	}()

	//减少付款商户金额，并记录账户变动日志
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, payerVAccNo, req.Amount, "-", logNo, constants.VaReason_BusinessTransferToBusiness); errStr != ss_err.ERR_SUCCESS {
		ss_sql.Rollback(tx)
		ss_log.Error("扣除转账金额失败, vAccountNo=%v, errStr=%v", payerVAccNo, errStr)
		orderStatus = constants.BusinessTransferOrderStatusFail
		return errStr, logNo, orderStatus
	}

	//扣除手续费，并记录账户变动日志
	if fee != "" && fee != "0" {
		if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, payerVAccNo, fee, "-", logNo, constants.VaReason_FEES); errStr != ss_err.ERR_SUCCESS {
			ss_sql.Rollback(tx)
			ss_log.Error("扣除转账手续费失败, vAccountNo=%v, errStr=%v", payerVAccNo, errStr)
			orderStatus = constants.BusinessTransferOrderStatusFail
			orderStatus = constants.BusinessTransferOrderStatusFail
			return errStr, logNo, orderStatus
		}
	}

	//增加收款商户金额，并记录账户变动日志
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, payeeVAccNo, req.Amount, "+", logNo, constants.VaReason_BusinessTransferToBusiness); errStr != ss_err.ERR_SUCCESS {
		ss_sql.Rollback(tx)
		ss_log.Error("增加收款商户金额失败, vAccountNo=%v, errStr=%v", payeeVAccNo, errStr)
		orderStatus = constants.BusinessTransferOrderStatusFail
		return errStr, logNo, orderStatus
	}

	tx.Commit()
	return ss_err.ERR_SUCCESS, logNo, orderStatus
}

//获取手续费
func getBusinessFee(amount, currencyType string) (feeT, rateT, errCodeT string) {
	//查询商家转账配置
	transferConf, err := dao.GlobalParamDaoInstance.GetBusinessTransferParamValue()
	if err != nil {
		ss_log.Error("查询商家转账配置失败, err=%v", err)
		return "", "", ss_err.ERR_SYSTEM
	}

	//判断交易金额是否超出限制,并计算手续费
	var feesDeci decimal.Decimal
	var rate string
	switch currencyType {
	case constants.CURRENCY_UP_USD:
		if strext.ToInt64(amount) >= transferConf.USDMinAmount && strext.ToInt64(amount) <= transferConf.USDMaxAmount {
			rate = strext.ToString(transferConf.USDRate)
			feesDeci = ss_count.CountFees(amount, rate, strext.ToStringNoPoint(transferConf.USDMinFee))
		} else {
			ss_log.Error("转账金额超出限制，")
			return "", "", ss_err.ERR_PAY_AMOUNT_IS_LIMIT
		}
	case constants.CURRENCY_UP_KHR:
		if strext.ToInt64(amount) >= transferConf.KHRMinAmount && strext.ToInt64(amount) <= transferConf.KHRMaxAmount {
			rate = strext.ToString(transferConf.KHRRate)
			feesDeci = ss_count.CountFees(amount, rate, strext.ToStringNoPoint(transferConf.KHRMinFee))
		} else {
			ss_log.Error("转账金额超出限制，")
			return "", "", ss_err.ERR_PAY_AMOUNT_IS_LIMIT
		}
	}
	// 取整
	fee := ss_big.SsBigInst.ToRound(feesDeci, 0, ss_big.RoundingMode_HALF_EVEN).String()

	return fee, rate, ss_err.ERR_SUCCESS
}

//商家转账给个人企业
func transferToUser(outBusiness *dao.BusinessStatus, req *billProto.AddBusinessTransferRequest) (codeT, logNoT, orderStatusT string) {
	//计算付款商家手续费
	fee, rate, getFeeCode := getBusinessFee(req.Amount, req.CurrencyType)
	if getFeeCode != ss_err.ERR_SUCCESS {
		return getFeeCode, "", ""
	}

	//判断付款商余额是否足够
	businessVAccType := global.GetBusinessVAccType(req.CurrencyType, true)
	err, balance := dao.VaccountDaoInst.GetBalanceFromAccNo(req.BusinessAccNo, businessVAccType)
	if err != nil {
		ss_log.Error("查询商家虚账失败, err=%v", err)
		return ss_err.ERR_PARAM, "", ""
	}
	if strext.ToInt(balance) < strext.ToInt(ss_count.Add(req.Amount, fee)) {
		ss_log.Error("商家余额不足，BusinessAccNo=%v, CurrencyType=%v,", req.BusinessAccNo, req.CurrencyType)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH, "", ""
	}

	//付款人虚账
	payerVAccNo := dao.VaccountDaoInst.GetVaccountNo(req.BusinessAccNo, strext.ToInt32(businessVAccType))
	if payerVAccNo == "" {
		ss_log.Error("没有查到商户虚拟账号，BusinessAccNo=%v, CurrencyType=%v", req.BusinessAccNo, req.CurrencyType)
		return ss_err.ERR_SYSTEM, "", ""
	}

	//查询收款用户是否存在
	accountNo, err := dao.AccDaoInstance.GetAccNoFromAccount(req.PayeeNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("收款人不存在，PayeeNo=%v, err=%v", req.PayeeNo, err)
			return ss_err.ERR_PayeeNotExist, "", ""
		}
		ss_log.Error("查询收款人失败，PayeeNo=%v, err=%v", req.PayeeNo, err)
		return ss_err.ERR_SYSTEM, "", ""
	}

	idenNo := dao.RelaAccIdenDaoInst.GetIdenNo(accountNo, constants.AccountType_USER)
	if idenNo == "" {
		ss_log.Error("查询账号的用户身份id为空，uid[%v]", accountNo)
		return ss_err.ERR_PayeeNotExist, "", ""
	}

	//确认用户打开收款权限
	_, _, _, inTransferAuthorization, _ := dao.CustDaoInstance.QueryRateRoleFrom(idenNo)
	if strext.ToInt(inTransferAuthorization) == constants.CustInTransferAuthorizationDisabled {
		ss_log.Error("收款人没有转入权限,accountNo: %s,custNo: %s", accountNo, idenNo)
		return ss_err.ERR_PAY_NO_IN_COME_PERMISSION, "", ""
	}

	//收款商家虚账,防止商家没有初始化此类型的虚账，这里做一下处理
	userVaccType, vaErr1 := common.VirtualAccountTypeByMoneyType(req.CurrencyType, "1")
	if vaErr1 != nil {
		ss_log.Error("Transfer 转账的人必须是有账号的 err: %s", vaErr1.Error())
		return ss_err.ERR_PARAM, "", ""
	}
	payeeVAccNo := InternalCallHandlerInst.ConfirmExistVAccount(accountNo, req.CurrencyType, strext.ToInt32(userVaccType))

	//插入转账日志
	d := new(dao.BusinessTransferOrderDao)
	d.FromBusinessNo = outBusiness.BusinessNo
	d.FromAccountNo = outBusiness.AccountNo
	d.ToAccountNo = accountNo
	d.Amount = req.Amount
	d.RealAmount = req.Amount
	d.Fee = fee
	d.Rate = rate
	d.CurrencyType = req.CurrencyType
	d.OrderStatus = constants.BusinessTransferOrderStatusPending
	d.Remarks = req.Remarks
	d.TransferType = req.TransferType
	d.OutTransferNo = req.OutOrderNo
	logNo, err := dao.BusinessTransferOrderDaoInst.Add(d)
	if err != nil {
		ss_log.Error("插入商家转账日志失败，data:%v, err=%v", strext.ToJson(d), err)
		return ss_err.ERR_SYSTEM, "", ""
	}

	//开启事务
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, txErr := dbHandler.BeginTx(context.TODO(), nil)
	if txErr != nil {
		ss_log.Error("开启事务失败, err=%v", txErr)
		return ss_err.ERR_SYSTEM, "", ""
	}

	orderStatus := constants.BusinessTransferOrderStatusSuccess

	//如果没有转账成功则修改订单状态为失败
	defer func() error {
		if err := dao.BusinessTransferOrderDaoInst.UpdateOrderStatusByLogNo(logNo, orderStatus, ""); nil != err {
			//ss_sql.Rollback(tx)
			ss_log.Error("修改订单状态失败，logNo=%v, err=%v", logNo, err)
			//reply.ResultCode = ss_err.ERR_SYSTEM
			return err
		}
		//ss_sql.Commit(tx)
		return nil
	}()

	//减少付款商户金额，并记录账户变动日志
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, payerVAccNo, req.Amount, "-", logNo, constants.VaReason_BusinessTransferToBusiness); errStr != ss_err.ERR_SUCCESS {
		ss_sql.Rollback(tx)
		ss_log.Error("扣除转账金额失败, vAccountNo=%v, errStr=%v", payerVAccNo, errStr)
		orderStatus = constants.BusinessTransferOrderStatusFail
		return errStr, logNo, orderStatus
	}

	//扣除手续费，并记录账户变动日志
	if fee != "" && fee != "0" {
		if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, payerVAccNo, fee, "-", logNo, constants.VaReason_FEES); errStr != ss_err.ERR_SUCCESS {
			ss_sql.Rollback(tx)
			ss_log.Error("扣除转账手续费失败, vAccountNo=%v, errStr=%v", payerVAccNo, errStr)
			orderStatus = constants.BusinessTransferOrderStatusFail
			return errStr, logNo, orderStatus
		}
	}

	//增加收款用户金额，并记录账户变动日志
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, payeeVAccNo, req.Amount, "+", logNo, constants.VaReason_BusinessTransferToBusiness); errStr != ss_err.ERR_SUCCESS {
		ss_sql.Rollback(tx)
		ss_log.Error("增加收款商户金额失败, vAccountNo=%v, errStr=%v", payeeVAccNo, errStr)
		orderStatus = constants.BusinessTransferOrderStatusFail
		return errStr, logNo, orderStatus
	}

	tx.Commit()
	return ss_err.ERR_SUCCESS, logNo, orderStatus
}

// 商家中心-单笔转账
func (b *BillHandler) AddBusinessTransfer(ctx context.Context, req *billProto.AddBusinessTransferRequest, reply *billProto.AddBusinessTransferReply) error {
	if strext.ToInt64(req.Amount) <= 0 {
		ss_log.Error("转账金额有误，req.Amount=%v", req.Amount)
		reply.ResultCode = ss_err.ERR_WALLET_AMOUNT_NULL
		return nil
	}
	if req.CurrencyType == "" {
		ss_log.Error("CurrencyType参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	req.CurrencyType = strings.ToUpper(req.CurrencyType)

	//查询付款商家的状态和出款权限
	outBusiness, err := dao.BusinessDaoInst.GetBusinessStatusInfo("", req.BusinessAccNo)
	if err != nil {
		ss_log.Error("查询付款商家状态失败，BusinessAccNo=%v, err=%v", req.BusinessAccNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}
	if strext.ToInt(outBusiness.BusinessStatus) == constants.BusinessUseStatusDisabled {
		ss_log.Error("付款商家已被禁用，BusinessAccNo=%v", req.BusinessAccNo)
		reply.ResultCode = ss_err.ERR_MERC_NO_USE
		return nil
	}
	if strext.ToInt(outBusiness.OutgoAuthorization) == constants.BusinessOutGoAuthDisabled {
		ss_log.Error("付款商家没有出款权限，BusinessAccNo=%v", req.BusinessAccNo)
		reply.ResultCode = ss_err.ERR_PAY_NO_OUT_GO_PERMISSION
		return nil
	}

	if req.TransferType == constants.BusinessTransferOrderTypeOrdinary {
		//验证支付密码
		replyCheckPayPwd, errCheckPayPwd := i.AuthHandlerInst.Client.CheckPayPWD(ctx, &authProto.CheckPayPWDRequest{
			AccountUid:  req.BusinessAccNo,
			AccountType: req.AccountType,
			Password:    req.PaymentPwd,
			NonStr:      req.NonStr,
			IdenNo:      req.BusinessNo,
		})
		if errCheckPayPwd != nil {
			ss_log.Error("调用验证支付密码接口出错,err[%v]", errCheckPayPwd)
			reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
			return nil
		}
		if replyCheckPayPwd.ResultCode != ss_err.ERR_SUCCESS {
			ss_log.Error("支付密码错误。req[%+v]", req)
			reply.ResultCode = ss_err.ERR_BusinessPayPwd_FAILD
			return nil
		}

		if req.CountryCode != "" { //转给的是个人 才带这个参数
			//处理手机号，拼成用户的账号
			req.PayeeNo = ss_func.ComposeAccountByPhoneCountryCode(ss_func.PrePhone(req.CountryCode, req.PayeeNo), req.CountryCode)

			code, logNo, orderStatus := transferToUser(outBusiness, req)
			reply.ResultCode = code
			reply.LogNo = logNo
			reply.OrderStatus = orderStatus

			if reply.ResultCode == ss_err.ERR_SUCCESS {
				//查询收款用户是否存在
				toAccountNo, err := dao.AccDaoInstance.GetAccNoFromAccount(req.PayeeNo)
				if err != nil {
					if err == sql.ErrNoRows {
						ss_log.Error("收款人不存在,无法推送消息，PayeeNo=%v, err=%v", req.PayeeNo, err)
						return nil
					}
					ss_log.Error("查询收款人失败,无法推送消息，PayeeNo=%v, err=%v", req.PayeeNo, err)
					return nil
				}
				//添加转账到的账号推送消息
				errAddMessages2 := dao.LogAppMessagesDaoInst.AddLogAppMessages(logNo, constants.LOG_APP_MESSAGES_ORDER_TYPE_TRANSFER_Apply, constants.VaReason_BusinessTransferToBusiness, toAccountNo, constants.OrderStatus_Paid)
				if errAddMessages2 != ss_err.ERR_SUCCESS {
					ss_log.Error("errAddMessages2=[%v]", errAddMessages2)
				}

				appLang, _ := dao.AccDaoInstance.QueryAccountLang(toAccountNo)
				if appLang == "" {
					req.Lang = constants.LangEnUS
				} else {
					req.Lang = appLang
				}

				ss_log.Info("用户 %s 当前的语言为--->%s", toAccountNo, req.Lang)
				moneyType := dao.LangDaoInstance.GetLangTextByKey(strings.ToLower(req.CurrencyType), req.Lang)
				timeString := time.Now().Format("2006-01-02 15:04:05")
				// 修正各币种的金额
				amount := common.NormalAmountByMoneyType(strings.ToLower(req.CurrencyType), req.Amount)

				args := []string{
					timeString, amount, moneyType,
				}
				if req.Lang == constants.LangEnUS {
					args = []string{
						amount, moneyType, timeString,
					}
				}

				// 消息推送
				ev := &go_micro_srv_push.PushReqest{
					Accounts: []*go_micro_srv_push.PushAccout{
						{
							AccountNo:   toAccountNo,
							AccountType: constants.AccountType_USER,
						},
					},
					TempNo: constants.Template_TransferSuccess,
					Args:   args,
				}

				ss_log.Info("publishing %+v\n", ev)
				// publish an event
				if err := common.PushEvent.Publish(context.TODO(), ev); err != nil {
					ss_log.Error("消息推送到用户toAccountNo[%v]出错。error : %v", toAccountNo, err)
				}

			}
		} else { //其他情况则是转给企业商家
			req.PayeeNo = strings.ToLower(req.PayeeNo)

			code, logNo, orderStatus := transferToBusiness(outBusiness, req)
			reply.ResultCode = code
			reply.LogNo = logNo
			reply.OrderStatus = orderStatus
		}
	} else {
		ss_log.Error("TransferType参数错误, req.TransferType=%v", req.TransferType)
		reply.ResultCode = ss_err.ERR_PARAM
	}

	return nil
}

// 商家中心-批量转账-分析批量转账文件内容
func (*BillHandler) GetBatchAnalysisResult(ctx context.Context, req *billProto.GetBatchAnalysisResultRequest, reply *billProto.GetBatchAnalysisResultReply) error {

	//查询上传的文件信息
	data, err := dao.UploadFileLogDaoInstance.GetUploadFileInfo(req.FileId)
	if err != nil {
		ss_log.Error("查询上传的文件失败，fileId[%v]", req.FileId)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if data.AccountNo != "" && data.AccountNo != req.AccountUid {
		ss_log.Error("文件不是该商家上传的,FileId[%v]", req.FileId)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if data.FileType != constants.UploadFileType_XLSX {
		ss_log.Error("文件不是xlsx文件,FileId[%v]", req.FileId)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if data.FileName == "" {
		ss_log.Error("上传的文件异常,FileId[%v]", req.FileId)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//获取数据，分析、处理数据
	batchDatas, errStr := getBatchAnalysisData(data.FileName, data.AccountNo)
	if errStr != ss_err.ERR_SUCCESS {
		ss_log.Error("分析批量付款文件获取结果失败,err[%v]", errStr)
		reply.ResultCode = errStr
		return nil
	}

	//获取批量转账限制的数量
	batchNum := dao.GlobalParamDaoInstance.QeuryParamValue(constants.GlobalParamKeyBusinessTransferBatchNum)
	if strext.ToInt(batchNum) < len(batchDatas) {
		ss_log.Error("批量转账超出限额笔数")
		reply.ResultCode = ss_err.ERR_BusinessBatchTransNumber_FAILD
		return nil
	}

	//统计数据
	dataResults := getBatchAnalysisResult(batchDatas)

	//生成文件的内容json字符串
	var temps []*util.FileContentJsonStruct
	var wrongDatas []*billProto.BatchAnalysisWrongResultData

	var wrongReasonArr []string
	for _, data := range batchDatas {
		temp := &util.FileContentJsonStruct{
			Row:          data.Row,
			ToAccount:    data.ToAccount,
			Amount:       data.Amount,
			CurrencyType: data.CurrencyType,
			Name:         data.AuthName,
			Remarks:      data.Remarks,
		}
		temps = append(temps, temp)

		if data.OrderStatus == constants.BusinessTransferOrderStatusFail {
			wrongDatas = append(wrongDatas, &billProto.BatchAnalysisWrongResultData{
				Row:          data.Row,          //序号
				Account:      data.ToAccount,    //账号
				Name:         data.AuthName,     //认证名称
				CurrencyType: data.CurrencyType, //币种
				Amount:       data.Amount,       //金额
				Remarks:      data.Remarks,      //备注
				WrongReason:  data.WrongReason,  //异常原因
			})
			wrongReasonArr = append(wrongReasonArr, data.WrongReason)

		}
	}

	//文件内容json字符串
	fileContent := strext.ToJson(temps)

	//插入转账批次表
	batchNo, errInsert := dao.BusinessBatchTransferOrderDaoInst.InsertOrder(dao.BusinessBatchTransferOrderDao{
		TotalNumber:      dataResults.TotalNumber,
		TotalAmount:      dataResults.TotalAmount,
		ProcessingNumber: "0",
		CurrencyType:     dataResults.CurrencyType,
		Remarks:          "",
		BusinessNo:       req.BusinessNo,
		FileContent:      fileContent,
		GenerateAll:      "false", //转账订单是否全部生成
		RealAmount:       dataResults.RealAmount,
		Fee:              dataResults.Fee,
	})
	if errInsert != nil {
		ss_log.Error("插入批次记录失败,err=[%v]", errInsert)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if len(wrongReasonArr) != 0 {
		//处理返回去的异常原因（记录的是错误码，需换成多语言对应的语言）
		var wrongReasonDatas []*dao.LangDao
		wrongReasonDatas, err = dao.LangDaoInstance.GetLangTextsByKeys(wrongReasonArr)
		if err != nil {
			ss_log.Error("查询多语言多个key对应文字出错，wrongReasonArr[%v],err[%v]", wrongReasonArr, err)
		}

		for _, data := range wrongDatas {
			for _, v := range wrongReasonDatas {
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
	reply.Data = dataResults
	reply.WrongDatas = wrongDatas
	reply.BatchNo = batchNo
	return nil
}

type BusinessTransferOrderData struct {
	AuthName     string
	ToAccountNo  string
	ToAccount    string
	ToBusinessNo string
	Amount       string
	CurrencyType string
	Rate         string
	Fee          string
	RealAmount   string
	OrderStatus  string
	Remarks      string
	Row          string
	WrongReason  string
}

//从文件中获取数据，分清哪些是可以成功的，哪些是会失败的。。返回两者的数据集合
func getBatchAnalysisData(fileName, fromAccNo string) ([]*BusinessTransferOrderData, string) {
	ss_log.Info("数据分析开始 start ")

	//从s3获取文件
	result, s3Err := common.UploadS3.GetObject(fileName)
	if s3Err != nil {
		ss_log.Error("从s3获取文件失败,FileName:%s, err:%v", fileName, s3Err)
		return nil, ss_err.ERR_FILE_OP_FAILD
	}

	// 读取body内容
	bytes, err := ioutil.ReadAll(result.Body)
	if err != nil {
		ss_log.Error("读取内容失败, result:%+v, err:%v", result, err)
		return nil, ss_err.ERR_FILE_OP_FAILD
	}

	//分析文件信息，得到结果
	xlFile, err := xlsx.OpenBinary(bytes)
	if err != nil {
		ss_log.Error("open failed: %s\n", err)
		return nil, ss_err.ERR_FILE_OP_FAILD
	}

	//查询商家转账配置
	transferConf, err := dao.GlobalParamDaoInstance.GetBusinessTransferParamValue()
	if err != nil {
		ss_log.Error("查询商家转账配置失败, err=%v", err)
		return nil, ss_err.ERR_SYS_DB_GET
	}

	var currencyType string
	var datas []*BusinessTransferOrderData
	for _, sheet := range xlFile.Sheets {
		for j, row := range sheet.Rows {
			if j == 0 { //去掉标题
				continue
			}
			if len(row.Cells) == 0 { //去掉空行
				continue
			}
			data := &BusinessTransferOrderData{
				ToAccountNo:  "",
				ToAccount:    "",
				ToBusinessNo: "",
				Amount:       "",
				CurrencyType: "",
				Rate:         "",
				Fee:          "",
				RealAmount:   "",
				OrderStatus:  constants.BusinessTransferOrderStatusSuccess,
				Remarks:      "",
				AuthName:     "",
				Row:          strext.ToStringNoPoint(j),
			}
			for k, cell := range row.Cells {
				text := cell.String()
				switch k {
				case 0: //序号（这是商家批量转账的序号，后台不做处理）
				case 1: //商家账号
					data.ToAccount = text
				case 2: //认证名称
					data.AuthName = text

					_, _, wrongReason := util.CheckToAccountAndAuthName(data.ToAccount, data.AuthName, fromAccNo)
					if wrongReason != "" && data.WrongReason == "" {
						data.WrongReason = wrongReason
					}
				case 3: //币种
					switch text {
					case constants.CURRENCY_UP_USD:
					case constants.CURRENCY_UP_KHR:
					default: //币种错误，直接打回
						return nil, ss_err.ERR_MONEY_TYPE_FAILD
					}

					if currencyType == "" { //文档内的币种必需是只有一种币种
						currencyType = text
					} else if currencyType != text {
						return nil, ss_err.ERR_MONEY_TYPE_FAILD
					}

					data.CurrencyType = text
				case 4: //金额
					switch data.CurrencyType {
					case constants.CURRENCY_UP_USD:
						data.Amount = ss_count.Multiply(text, "100").String()
					case constants.CURRENCY_UP_KHR:
						data.Amount = text
					default: //币种错误，直接打回
						return nil, ss_err.ERR_MONEY_TYPE_FAILD
					}

					//校验转账金额
					if wrongReason := util.CheckTransferAmount(data.Amount, data.CurrencyType, transferConf); wrongReason != "" && data.WrongReason == "" {
						data.WrongReason = wrongReason
					}

					//计算手续费
					fee, _, wrongReason := util.QueryTransferFeeAndRate(data.Amount, data.CurrencyType, transferConf)
					data.Fee = fee
					if wrongReason != "" && data.WrongReason == "" {
						data.WrongReason = wrongReason
					}
				case 5: //备注
					data.Remarks = text
				}

				ss_log.Info("行号%v 第%v列, 内容: %s\n", j, k, text)
			}

			if data.ToAccount == "" {
				ss_log.Error("行号:%v,未填写账号", j)
				if data.WrongReason == "" {
					data.WrongReason = ss_err.ERR_MSG_ACCOUNT_NOT_EXISTS //"未填写账号"
				}
			}

			if data.AuthName == "" {
				ss_log.Error("行号:%v,未填写认证名称", j)
				if data.WrongReason == "" {
					data.WrongReason = ss_err.ERR_UnFilledAuthName_FAILD //"未填写认证名称"
				}
			}

			if data.Amount == "" {
				ss_log.Error("行号:%v,未填写金额", j)
				data.Amount = "0"
				if data.WrongReason == "" {
					data.WrongReason = ss_err.ERR_WALLET_AMOUNT_NULL // "未填写金额"
				}
			}

			if data.CurrencyType == "" {
				ss_log.Error("行号:%v,未填写币种", j)
				if data.WrongReason == "" {
					data.WrongReason = ss_err.ERR_UnFilledCurrencyType_FAILD // "未填写币种"
				}
			}

			if data.WrongReason != "" {
				data.OrderStatus = constants.BusinessTransferOrderStatusFail
			}

			datas = append(datas, data)
		}
	}

	//
	ss_log.Info("数据分析完 end ")

	return datas, ss_err.ERR_SUCCESS
}

//获取分析结果
func getBatchAnalysisResult(datas []*BusinessTransferOrderData) *billProto.BatchAnalysisResultData {
	totalNumber := "0"      //   总笔数
	totalAmount := "0"      //   总金额
	failNumber := "0"       //  异常笔数
	failAmount := "0"       //  异常金额
	successfulNumber := "0" //  成功笔数
	successfulAmount := "0" //  成功金额
	realAmount := "0"       //实际支付金额(订单金额加手续费，生成时还未付款，要付款是付的这个金额)
	fee := "0"              //手续费
	currencyType := ""      //币种
	wrongStatus := "0"      //0-全部正常，1-全部异常，2-部分异常
	for _, data := range datas {
		totalNumber = ss_count.Add(totalNumber, "1")
		totalAmount = ss_count.Add(data.Amount, totalAmount)

		if currencyType == "" {
			currencyType = data.CurrencyType
		}

		switch data.OrderStatus {
		case constants.BusinessTransferOrderStatusFail:
			failNumber = ss_count.Add(failNumber, "1")
			failAmount = ss_count.Add(data.Amount, failAmount)
		case constants.BusinessTransferOrderStatusSuccess:
			successfulNumber = ss_count.Add(successfulNumber, "1")
			successfulAmount = ss_count.Add(data.Amount, successfulAmount)

			//实际支付金额
			realAmount = ss_count.Add(data.Amount, realAmount)
			realAmount = ss_count.Add(data.Fee, realAmount)

			fee = ss_count.Add(data.Fee, fee)
		}

	}

	if failNumber == totalNumber { //全部异常
		wrongStatus = "1"
	} else if failNumber != "0" && failNumber < totalNumber { //部分异常
		wrongStatus = "2"
	}

	return &billProto.BatchAnalysisResultData{
		TotalNumber:      totalNumber,      //总笔数
		TotalAmount:      totalAmount,      //总金额
		SuccessfulAmount: successfulAmount, //成功金额
		SuccessfulNumber: successfulNumber, //成功笔数
		RealAmount:       realAmount,       //实际支付金额
		RealNumber:       successfulNumber, //实际支付笔数
		FailNumber:       failNumber,       //异常笔数
		FailAmount:       failAmount,       //异常金额
		Fee:              fee,              //手续费
		CurrencyType:     currencyType,     //币种
		WrongStatus:      wrongStatus,      //0-全部正常，1-全部异常，2-部分异常
	}
}

// 商家中心-批量转账-输入支付密码确认
func (b *BillHandler) BusinessBatchTransferConfirm(ctx context.Context, req *billProto.BusinessBatchTransferConfirmRequest, reply *billProto.BusinessBatchTransferConfirmReply) error {

	//验证支付密码
	replyCheckPayPwd, errCheckPayPwd := i.AuthHandlerInst.Client.CheckPayPWD(ctx, &authProto.CheckPayPWDRequest{
		AccountUid:  req.BusinessAccNo,
		AccountType: req.AccountType,
		Password:    req.PayPwd,
		NonStr:      req.NonStr,
		IdenNo:      req.BusinessNo,
	})
	if errCheckPayPwd != nil {
		ss_log.Error("调用验证支付密码接口出错,err[%v]", errCheckPayPwd)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	if replyCheckPayPwd.ResultCode != ss_err.ERR_SUCCESS {
		ss_log.Error("支付密码错误。req[%+v]", req)
		reply.ResultCode = ss_err.ERR_WALLET_PAY_PWD_ERR
		return nil
	}

	//查询转账批次信息
	batchData, errGet1 := dao.BusinessBatchTransferOrderDaoInst.GetBatchOrderDetail(req.BatchNo)
	if errGet1 != nil {
		ss_log.Error("获取批次记录失败,BatchNo[%v],err=[%v]", req.BatchNo, errGet1)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	if batchData.RealAmount == "0" {
		ss_log.Error("批次转账要付款的金额为0,不需要支付,BatchNo[%v]", req.BatchNo)
		reply.ResultCode = ss_err.ERR_BusinessBatchTransAmount_FAILD
		return nil
	}

	if batchData.Status != constants.BusinessBatchTransferOrderStatusPending {
		ss_log.Error("当前批量转账批次不是待支付状态,无法重复支付,BatchNo[%v]", req.BatchNo)
		reply.ResultCode = ss_err.ERR_BusinessBatchTransferOrderStatus_FAILD
		return nil
	}

	if batchData.BusinessNo != req.BusinessNo {
		ss_log.Error("当前支付的批量转账批次不是该商家的，无法进行支付")
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//查询付款商家的状态和出款权限
	outBusiness, err := dao.BusinessDaoInst.GetBusinessStatusInfo("", req.BusinessAccNo)
	if err != nil {
		ss_log.Error("查询付款商家状态失败，BusinessAccNo=%v, err=%v", req.BusinessAccNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}
	if strext.ToInt(outBusiness.BusinessStatus) == constants.BusinessUseStatusDisabled {
		ss_log.Error("付款商家已被禁用，BusinessAccNo=%v", req.BusinessAccNo)
		reply.ResultCode = ss_err.ERR_MERC_NO_USE
		return nil
	}
	if strext.ToInt(outBusiness.OutgoAuthorization) == constants.BusinessOutGoAuthDisabled {
		ss_log.Error("付款商家没有出款权限，BusinessAccNo=%v", req.BusinessAccNo)
		reply.ResultCode = ss_err.ERR_PAY_NO_OUT_GO_PERMISSION
		return nil
	}

	//判断付款商余额是否足够
	businessVAccType := global.GetBusinessVAccType(batchData.CurrencyType, true)
	err, balance := dao.VaccountDaoInst.GetBalanceFromAccNo(req.BusinessAccNo, businessVAccType)
	if err != nil {
		ss_log.Error("查询商家虚账失败, err=%v", err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if strext.ToInt(balance) < strext.ToInt(batchData.RealAmount) {
		ss_log.Error("商家余额不足以支付该转账批次的钱，balance=%v, RealAmount=%v,", balance, batchData.RealAmount)
		reply.ResultCode = ss_err.ERR_PAY_AMT_NOT_ENOUGH
		return nil
	}

	//付款人虚账
	payerVAccNo := dao.VaccountDaoInst.GetVaccountNo(req.BusinessAccNo, strext.ToInt32(businessVAccType))
	if payerVAccNo == "" {
		ss_log.Error("没有查到商户虚拟账号，BusinessAccNo=%v, CurrencyType=%v", req.BusinessAccNo, batchData.CurrencyType)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//开启事务
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, txErr := dbHandler.BeginTx(ctx, nil)
	if txErr != nil {
		ss_log.Error("开启事务失败, err=%v", txErr)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}
	defer tx.Rollback()

	//减少付款商户余额到冻结金额，并记录账户变动日志
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainAndFrozenUpperZero(tx, "-", payerVAccNo, ss_count.Sub(batchData.RealAmount, batchData.Fee).String(), batchData.BatchNo, constants.VaReason_BusinessBatchTransferToBusiness, constants.VaOpType_Freeze); errStr != ss_err.ERR_SUCCESS {
		ss_log.Error("扣除批量转账金额失败, vAccountNo=%v, errMsg=%v", payerVAccNo, errStr)
		reply.ResultCode = errStr
		return nil
	}

	if batchData.Fee != "" && batchData.Fee != "0" {
		if errStr := dao.VaccountDaoInst.ModifyVaccRemainAndFrozenUpperZero(tx, "-", payerVAccNo, batchData.Fee, batchData.BatchNo, constants.VaReason_FEES, constants.VaOpType_Freeze); errStr != ss_err.ERR_SUCCESS {
			ss_log.Error("扣除批量转账金额失败, vAccountNo=%v, errMsg=%v", payerVAccNo, errStr)
			reply.ResultCode = errStr
			return nil
		}
	}

	//这里是修改批量转账批次的状态，如果是支付成功了，到后面会有定时任务一个一个的开始转账。
	if err := dao.BusinessBatchTransferOrderDaoInst.UpdateOrderStatusPaySuccess(tx, batchData.BatchNo); nil != err {
		ss_log.Error("修改订单状态失败，logNo=%v, err=%v", batchData.BatchNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	ss_sql.Commit(tx)

	go func() error { //异步处理该转账批次的任务
		doTask := &cron.BusinessBatchTransfer{
			CronBase: cron.CronBase{LogCat: "异步处理已支付的转账批量订单:", LockExpire: time.Hour * 2},
		}

		//查询商家转账配置
		transferConf, err := dao.GlobalParamDaoInstance.GetBusinessTransferParamValue()
		if err != nil {
			ss_log.Error("查询商家转账配置失败, err=[%v]", err)
			return err
		}

		if err := doTask.DoSingleBatchTransfer(batchData, transferConf); err != nil {
			ss_log.Error("err=%v", err)
			return err
		}

		return nil
	}()

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//企业转账给普通用户(只给business-bill调用，如果别的服务要用请另写)
func (b *BillHandler) EnterpriseTransferToUser(ctx context.Context, req *billProto.EnterpriseTransferToUserRequest, reply *billProto.EnterpriseTransferToUserReply) error {
	errCode, err := CheckEnterpriseTransferToUserParam(req)
	if err != nil {
		ss_log.Error("参数有误：%v", err)
		reply.ResultCode = errCode
		return nil
	}

	order, err := dao.BusinessTransferOrderDaoInst.GetTransferOrderByLogNo(req.TransferNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("转账订单[%v]不存在", req.TransferNo)
			reply.ResultCode = ss_err.ERR_PayOrderNoNotExist
			return nil
		}
		ss_log.Error("查询转账订单[%v]失败，err=%v", req.TransferNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	if order.OrderStatus != constants.BusinessTransferOrderStatusPending {
		ss_log.Error("订单[%v]目前状态[%v]不支持转账", req.TransferNo, order.OrderStatus)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if order.FromAccountNo == order.ToAccountNo {
		ss_log.Error("不能转账给自己, payer=%v, payee=%v", order.FromAccountNo, order.ToAccountNo)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//查询付款商家的状态和出款权限
	outBusiness, err := dao.BusinessDaoInst.GetBusinessStatusInfo("", order.FromAccountNo)
	if err != nil {
		ss_log.Error("查询付款商家状态失败，BusinessAccNo=%v, err=%v", order.FromAccountNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}
	if strext.ToInt(outBusiness.BusinessStatus) == constants.BusinessUseStatusDisabled {
		ss_log.Error("付款商家已被禁用，BusinessAccNo=%v", order.FromAccountNo)
		reply.ResultCode = ss_err.ERR_MERC_NO_USE
		return nil
	}
	if strext.ToInt(outBusiness.OutgoAuthorization) == constants.BusinessOutGoAuthDisabled {
		ss_log.Error("付款商家没有出款权限，BusinessAccNo=%v", order.FromAccountNo)
		reply.ResultCode = ss_err.ERR_PAY_NO_OUT_GO_PERMISSION
		return nil
	}

	//判断付款商余额是否足够
	businessVAccType := global.GetBusinessVAccType(order.CurrencyType, true)
	err, balance := dao.VaccountDaoInst.GetBalanceFromAccNo(order.FromAccountNo, businessVAccType)
	if err != nil {
		ss_log.Error("查询商家虚账失败, err=%v", err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	if strext.ToInt(balance) < strext.ToInt(ss_count.Add(strext.ToString(order.Amount), strext.ToString(order.Fee))) {
		ss_log.Error("商家余额不足，BusinessAccNo=%v, CurrencyType=%v,", order.FromAccountNo, order.CurrencyType)
		reply.ResultCode = ss_err.ERR_PAY_AMT_NOT_ENOUGH
		return nil
	}

	//查询收款用户是否存在
	payee, err := dao.AccDaoInstance.GetAccountByAccNo(order.ToAccountNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("收款人[%v]不存在", order.ToAccountNo)
			reply.ResultCode = ss_err.ERR_PayeeNotExist
			return nil
		}
		ss_log.Error("查询收款人[%v]账号失败, err=%v", order.ToAccountNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	var payeeVAccType int
	if payee.AccountType == constants.AccountType_USER {
		payeeVAccType = global.GetUserVAccType(order.CurrencyType, true)
		//确认用户打开收款权限
		_, _, _, inTransferAuthorization, _ := dao.CustDaoInstance.QueryRateRoleFrom(payee.IdentityNo)
		if strext.ToInt(inTransferAuthorization) == constants.CustInTransferAuthorizationDisabled {
			ss_log.Error("收款人没有转入权限,accountNo: %s,custNo: %s", order.ToAccountNo, payee.IdentityNo)
			reply.ResultCode = ss_err.ERR_PAY_NO_IN_COME_PERMISSION
			return nil
		}
	} else {
		payeeVAccType = global.GetBusinessVAccType(order.CurrencyType, true)
		payeeBusiness, err := dao.BusinessDaoInst.GetBusinessStatusInfo("", order.ToAccountNo)
		if err != nil {
			ss_log.Error("查询账号的用户身份id为空，uid[%v]", order.ToAccountNo)
			reply.ResultCode = ss_err.ERR_PayeeNotExist
			return nil
		}
		if strext.ToInt(payeeBusiness.AccountStatus) != constants.AccountUseStatusNormal {
			ss_log.Error("收款商家账号已被禁用，PayeeNo=%v", order.ToAccountNo)
			reply.ResultCode = ss_err.ERR_MERC_NO_USE
			return nil

		}
		if strext.ToInt(payeeBusiness.BusinessStatus) == constants.BusinessUseStatusDisabled {
			ss_log.Error("收款商家已被禁用，PayeeNo=%v", order.ToAccountNo)
			reply.ResultCode = ss_err.ERR_MERC_NO_USE
			return nil

		}
		if strext.ToInt(payeeBusiness.IncomeAuthorization) == constants.BusinessIncomeAuthDisabled {
			ss_log.Error("收款商家没有收款权限，PayeeNo=%v", order.ToAccountNo)
			reply.ResultCode = ss_err.ERR_PAY_NO_IN_COME_PERMISSION
			return nil
		}
	}

	//付款商户虚账
	payerVAccNo := dao.VaccountDaoInst.GetVaccountNo(order.FromAccountNo, strext.ToInt32(businessVAccType))
	if payerVAccNo == "" {
		ss_log.Error("没有查到商户虚拟账号，BusinessAccNo=%v, CurrencyType=%v", order.FromAccountNo, order.CurrencyType)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//收款商家虚账,防止商家没有初始化此类型的虚账，这里做一下处理
	payeeVAccNo := InternalCallHandlerInst.ConfirmExistVAccount(order.ToAccountNo, order.CurrencyType, strext.ToInt32(payeeVAccType))

	result, wrongReason := turnAmount(&TurnAmount{
		LogNo:       order.LogNo,
		Amount:      strext.ToString(order.Amount),
		Fee:         strext.ToString(order.Fee),
		PayerVAccNo: payerVAccNo,
		PayeeVAccNo: payeeVAccNo,
	})
	if result != ss_err.ERR_SUCCESS {
		//修改平台转账订单的状态
		err = dao.BusinessTransferOrderDaoInst.UpdateOrderStatusByLogNo(order.LogNo, constants.BusinessTransferOrderStatusFail, wrongReason)
		if err != nil {
			ss_log.Error("修改转账订单状态失败, logNo=%v, errS=%v", payeeVAccNo, err)
			reply.ResultCode = ss_err.ERR_SYSTEM
			return nil
		}
	}

	enterpriseOrder := &billProto.EnterpriseTransfer{
		OrderNo:      order.LogNo,
		Amount:       strext.ToString(order.Amount),
		CurrencyType: order.CurrencyType,
		PayeeAccNo:   payee.AccountNo,
		PayeeAccType: payee.AccountType,
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Order = enterpriseOrder
	return nil
}

type TurnAmount struct {
	LogNo       string
	Amount      string
	Fee         string
	PayerVAccNo string
	PayeeVAccNo string
}

func turnAmount(req *TurnAmount) (errCode string, wrongReason string) {
	//开启事务
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, txErr := dbHandler.BeginTx(context.TODO(), nil)
	if txErr != nil {
		ss_log.Error("开启事务失败, err=%v", txErr)
		return ss_err.ERR_SYSTEM, "系统错误"
	}

	//减少付款商户金额，并记录账户变动日志
	errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, req.PayerVAccNo, strext.ToString(req.Amount), "-", req.LogNo, constants.VaReason_BusinessTransferToBusiness)
	if errStr != ss_err.ERR_SUCCESS {
		ss_sql.Rollback(tx)
		ss_log.Error("扣除转账金额失败, vAccountNo=%v, errStr=%v", req.PayerVAccNo, errStr)
		return errStr, "付款方转账金额扣除失败"
	}

	//扣除手续费，并记录账户变动日志
	if req.Fee != "" && req.Fee != "0" {
		if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, req.PayerVAccNo, strext.ToString(req.Fee), "-", req.LogNo, constants.VaReason_FEES); errStr != ss_err.ERR_SUCCESS {
			ss_sql.Rollback(tx)
			ss_log.Error("扣除转账手续费失败, vAccountNo=%v, errStr=%v", req.PayerVAccNo, errStr)
			return errStr, "扣除付款方手续费失败"
		}
	}

	//增加收款用户金额，并记录账户变动日志
	errStr = dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, req.PayeeVAccNo, strext.ToString(req.Amount), "+", req.LogNo, constants.VaReason_BusinessTransferToBusiness)
	if errStr != ss_err.ERR_SUCCESS {
		ss_sql.Rollback(tx)
		ss_log.Error("增加收款商户金额失败, vAccountNo=%v, errStr=%v", req.PayeeVAccNo, errStr)
		return errStr, "收款方接收失败"
	}

	//订单的状态
	err := dao.BusinessTransferOrderDaoInst.UpdateOrderStatusByLogNoTx(tx, req.LogNo, constants.BusinessTransferOrderStatusSuccess, "")
	if err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("修改转账订单状态失败, logNo=%v, err=%v", req.LogNo, err)
		return ss_err.ERR_SYSTEM, "系统错误"
	}

	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS, ""
}
