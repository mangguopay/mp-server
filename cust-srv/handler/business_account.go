package handler

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/model"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_func"
	"a.a/mp-server/common/ss_rsa"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/dao"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

/**
忘记密码里的身份验证，返回脱敏的手机号和邮箱
*/
func (c *CustHandler) IdentityVerify(ctx context.Context, req *custProto.IdentityVerifyRequest, resp *custProto.IdentityVerifyReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//校验验证码
	ret, errVerifyCode := cache.RedisClient.Get("verify_" + req.Verifyid).Result()
	if ret == "" || errVerifyCode != nil {
		ss_log.Error("----------->crypt%s", errVerifyCode.Error())
		resp.ResultCode = ss_err.ERR_ACCOUNT_SMS_CODE
		return nil
	}

	errVerifyCode = cache.RedisClient.Del("verify_" + req.Verifyid).Err()
	if errVerifyCode != nil {
		ss_log.Error("err=[%v]", errVerifyCode)
	}

	if strings.ToLower(strext.ToStringNoPoint(ret)) != strings.ToLower(req.Verifynum) { //验证码错误
		resp.ResultCode = ss_err.ERR_ACCOUNT_LOGIN_CODE
		return nil
	}

	accNo, err := dao.AccDaoInstance.GetUidByAccount(req.Account)
	if err != nil || accNo == "" {
		ss_log.Error("err=[%v]", err)
		resp.ResultCode = ss_err.ERR_MSG_ACCOUNT_NOT_EXISTS
		return nil
	}

	//企业商家的手机号用的是account表的business_phone字段
	phone := dao.AccDaoInstance.GetBusinessPhoneFromAccNo(accNo)
	email := dao.AccDaoInstance.GetEmailFromAccNo(accNo)

	//手机号、邮箱脱敏
	if phone != "" { //9位以下保留前后两位，9位和9位以上保留前后三位
		if len(phone) >= 5 {
			//获取脱敏后的手机号，
			resp.Phone = ss_func.GetDesensitizationPhone(phone)
		} else {
			ss_log.Error("账号uid[%v]手机号码[%v]位数异常(长度小于5),", accNo, phone)
		}
	}

	if email != "" { //9位以下保留前后两位，9位和9位以上保留前后三位
		resp.Email = ss_func.GetDesensitizationEmail(email)
	}

	resp.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
 * 获取商户账户流水
 */
func (c *CustHandler) GetBusinessVAccLogList(ctx context.Context, req *custProto.GetBusinessVAccLogListRequest, reply *custProto.GetBusinessVAccLogListReply) error {
	if req.BusinessAccNo == "" {
		ss_log.Error("BusinessAccNo参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.MoneyType == "" {
		ss_log.Error("MoneyType参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	vAccType := global.GetBusinessVAccType(req.MoneyType, true)
	//查询商家虚账
	vAccountNo, err := dao.VaccountDaoInst.GetVaccountNo(req.BusinessAccNo, strext.ToInt32(vAccType))
	if err != nil {
		ss_log.Error("查询商家虚账失败，BusinessAccNo=%v, MoneyType=%v", req.BusinessAccNo, req.MoneyType)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//暂时不显示冻结金额变动日志
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "vaccount_no", Val: vAccountNo, EqType: "="},
		{Key: "create_time", Val: req.StartTime, EqType: ">="},
		{Key: "create_time", Val: req.EndTime, EqType: "<="},
		{Key: "op_type", Val: constants.VaOpType_Defreeze, EqType: "!="},
		{Key: "op_type", Val: constants.VaOpType_Defreeze_Minus, EqType: "!="},
		{Key: "op_type", Val: constants.VaOpType_Defreeze_But_Minus, EqType: "!="},
	})

	totalNum, err := dao.LogVaccountDaoInst.CountLogVAccount(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("统计虚账日志失败，vAccountNo=%v, err=%v", vAccountNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "ORDER BY create_time DESC, balance ASC ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
	list, err := dao.LogVaccountDaoInst.GetLogVAccountList(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("查询虚账日志列表失败，vAccount=%v, err=%v", vAccountNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	var dataList []*custProto.BusinessVAccLog
	for _, v := range list {
		data := new(custProto.BusinessVAccLog)
		data.LogNo = v.LogNo
		data.ChangeAmount = v.Amount
		data.ChangeAfterBalance = v.Balance
		data.Reason = v.Reason
		data.CreateTime = v.CreateTime
		data.OpType = v.OpType
		switch v.OpType {
		case constants.VaOpType_Add:
			fallthrough
		case constants.VaOpType_Defreeze_Add:
			data.ChangeBeforeBalance = ss_count.Sub(v.Balance, v.Amount).String()
			data.OpType = constants.VaOpType_Add //前端是粗暴的1+,2-,这里是换成前端认识的OpType
		case constants.VaOpType_Minus:
			fallthrough
		case constants.VaOpType_Freeze:
			data.ChangeBeforeBalance = ss_count.Add(v.Balance, v.Amount)
		}

		dataList = append(dataList, data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Total = totalNum
	reply.LogList = dataList
	return nil
}

/**
 * 获取商户账户流水详情（结算、提现、充值、退款、转账）
 */
func (c *CustHandler) GetBusinessVAccLogDetail(ctx context.Context, req *custProto.GetBusinessVAccLogDetailRequest, reply *custProto.GetBusinessVAccLogDetailReply) error {
	if req.LogNo == "" {
		ss_log.Error("LogNo参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.Reason == "" {
		ss_log.Error("Reason参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//查询虚账日志的业务流水号
	bizLogNo, opType, err := dao.LogVaccountDaoInst.GetBizLogNoAndOpType(req.LogNo, req.Reason)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("虚账日志不存在，logNo=%v", req.LogNo)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		ss_log.Error("查询虚账日志的业务流水号失败，logNo=%v, err=%v", req.LogNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	data := new(custProto.BusinessVAccLogDetail)
	data.OpType = opType
	switch req.Reason {
	case constants.VaReason_Business_Settle: //商家结算
		settleLog, err := dao.BusinessBillSettleDaoInst.GetSettleLogById(bizLogNo)
		if err != nil {
			if err != sql.ErrNoRows {
				ss_log.Error("查询商家结算详情失败，settleId=%v, err=%v", bizLogNo, err)
				reply.ResultCode = ss_err.ERR_SYSTEM
				return nil
			}
			settleLog, err = dao.BusinessBillSettleDaoInst.GetSingleSettleDetail(bizLogNo)
			if err != nil {
				ss_log.Error("查询商家结算详情失败，settleId=%v, err=%v", bizLogNo, err)
				reply.ResultCode = ss_err.ERR_SYSTEM
				return nil
			}
		}
		data.LogNo = settleLog.SettleId             //结算日志号
		data.Amount = settleLog.TotalAmount         //结算总金额
		data.RealAmount = settleLog.TotalRealAmount //总实际到账金额
		data.Fee = settleLog.TotalFees              //总手续费
		data.CurrencyType = settleLog.CurrencyType  //币种
		data.CreateTime = settleLog.CreateTime      //创建时间
		data.PayeeName = settleLog.BusinessName     //收款商家名称
	case constants.VaReason_Business_Save: //商家充值
		whereList := []*model.WhereSqlCond{
			{Key: "bth.log_no", Val: bizLogNo, EqType: "="},
		}
		//获取详情
		toHead, err := dao.LogBusinessToHeadquartersDaoInst.GetBusinessToHeadDetail(whereList)
		if err != nil {
			ss_log.Error("查询商家充值详情失败，bizLogNo=%v, err=%v", bizLogNo, err)
			reply.ResultCode = ss_err.ERR_SYSTEM
			return nil
		}
		data.LogNo = toHead.LogNo               //充值日志号
		data.Amount = toHead.Amount             //充值金额
		data.RealAmount = toHead.ArriveAmount   //实际到账金额
		data.Fee = toHead.Fee                   //手续费
		data.CurrencyType = toHead.CurrencyType //币种
		data.CreateTime = toHead.CreateTime     //创建时间
	case constants.VaReason_Business_Withdraw: //商家提现
		withdrow, err := dao.LogToBusinessDaoInst.GetToBusinessDetail(bizLogNo)
		if err != nil {
			ss_log.Error("查询提现详情失败, bizLogNo=%v, err=%v", bizLogNo, err)
			reply.ResultCode = ss_err.ERR_SYSTEM
			return nil
		}
		data.LogNo = withdrow.LogNo               //提现日志号
		data.Amount = withdrow.Amount             //提现金额
		data.RealAmount = withdrow.RealAmount     //提现实际到账金额
		data.Fee = withdrow.Fee                   //提现手续费
		data.CurrencyType = withdrow.CurrencyType //币种
		data.CreateTime = withdrow.CreateTime     //创建时间
	case constants.VaReason_BusinessRefund: //商家退款
		refundDetail, err := dao.BusinessBillRefundDaoInst.GetOrderDetail(bizLogNo)
		if err != nil {
			ss_log.Error("查询退款订单详情失败，bizLogNo=%v, err=%v", bizLogNo, err)
			reply.ResultCode = ss_err.ERR_SYSTEM
			return nil
		}
		data.LogNo = refundDetail.RefundNo            //退款订单号
		data.Amount = refundDetail.RefundAmount       //退款发起金额
		data.RealAmount = refundDetail.RefundAmount   //退款到账金额
		data.Fee = "0"                                //退款手续费，目前为0
		data.CurrencyType = refundDetail.CurrencyType //币种
		data.CreateTime = refundDetail.CreateTime     //创建时间
		data.PayeeName = refundDetail.PayeeName       //收款人(昵称)
		//data.PayeeAccount = ss_func.GetDesensitizationAccount(refundDetail.PayeeAcc) //收款人账号(脱敏)
	case constants.VaReason_BusinessTransferToBusiness: //商家转账
		transferDetail, err := dao.BusinessTransferDaoInst.GetOrderDetail(bizLogNo)
		if err != nil {
			ss_log.Error("查询商家转账详情失败，logNo=%v, err=%v", bizLogNo, err)
			reply.ResultCode = ss_err.ERR_SYSTEM
			return nil
		}
		data.LogNo = transferDetail.LogNo                //转账日志号
		data.Amount = transferDetail.Amount              //转账金额
		data.RealAmount = transferDetail.RealAmount      //实际到账金额
		data.Fee = transferDetail.Fee                    //手续费
		data.CurrencyType = transferDetail.CurrencyType  //币种
		data.CreateTime = transferDetail.CreateTime      //创建时间
		data.PayeeName = transferDetail.ToBusinessName   //收款商家名称
		data.PayerName = transferDetail.FromBusinessName //付款商家名称
	case constants.VaReason_BusinessBatchTransferToBusiness:
		whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
			{Key: "btb.batch_no", Val: bizLogNo, EqType: "="},
		})
		batchData, err := dao.BusinessBatchTransferDaoInst.GetBatchOrderDetail(whereModel.WhereStr, whereModel.Args)
		if err != nil {
			ss_log.Error("查询商家转账详情失败，logNo=%v, err=%v", bizLogNo, err)
			reply.ResultCode = ss_err.ERR_SYSTEM
			return nil
		}
		data.LogNo = batchData.BatchNo               //批量转账批次号
		data.Amount = batchData.TotalAmount          //发起金额(订单总金额)
		data.RealAmount = batchData.SuccessfulAmount //实际到账（成功金额）
		data.CurrencyType = batchData.CurrencyType   //币种
		data.CreateTime = batchData.CreateTime       //创建时间
	case constants.VaReason_FEES:
		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	default:
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

/**
生成商户APP秘钥对文件，支持PKCS1、PKCS8格式
*/
func (c *CustHandler) GenerateKeys(ctx context.Context, req *custProto.GenerateKeysRequest, reply *custProto.GenerateKeysReply) error {
	if req.KeyType == "" {
		req.KeyType = constants.SecretKeyPKCS1
	}
	if !util.InSlice(req.KeyType, []string{constants.SecretKeyPKCS1, constants.SecretKeyPKCS8}) {
		ss_log.Error("秘钥格式错误,keyType=[%v]", req.KeyType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//平台私钥、平台公钥
	var privateKey, publicKey string
	var err error
	switch req.KeyType {
	case constants.SecretKeyPKCS1:
		privateKey, publicKey, err = ss_rsa.GenRsaKeyPairPKCS1(2048)
	case constants.SecretKeyPKCS8:
		privateKey, publicKey, err = ss_rsa.GenRsaKeyPairPKCS8(2048)
	}
	if err != nil {
		ss_log.Error("生成秘钥对失败, KeyType=%v, err=%v", req.KeyType, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	privateKey = ss_rsa.StripRSAKey(privateKey)
	publicKey = ss_rsa.StripRSAKey(publicKey)
	fileName := fmt.Sprintf("public_private_%v.key", req.KeyType)
	ss_log.Info("fileName:[%v]", fileName)

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.FileName = fileName
	reply.FileContent = fmt.Sprintf("私钥:\r\n%s \r\n\r\n公钥:\r\n%s", privateKey, publicKey)
	return nil
}
