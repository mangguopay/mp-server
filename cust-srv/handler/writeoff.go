package handler

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/model"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/dao"
	"context"
	"database/sql"
	"fmt"
)

/**
 * 获取核销码列表
 */
func (*CustHandler) GetWriteOffList(ctx context.Context, req *custProto.GetWriteOffListRequest, reply *custProto.GetWriteOffListReply) error {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "wo.code", Val: req.Code, EqType: "="},
		{Key: "wo.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "wo.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "wo.use_status", Val: req.UseStatus, EqType: "="},
		{Key: "acc1.account", Val: req.PayerAccount, EqType: "="},
		{Key: "acc2.account", Val: req.PayeeAccount, EqType: "="},
	})

	total, err := dao.WriteoffDaoInst.CntCode(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("统计核销码数量失败，request=%v, err=%v", strext.ToJson(req), err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "ORDER BY wo.create_time DESC ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
	list, err := dao.WriteoffDaoInst.GetCodeList(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("查询核销码列表失败，request=%v, err=%v", strext.ToJson(req), err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	var dataList []*custProto.WriteOff
	for _, v := range list {
		data := new(custProto.WriteOff)
		data.Code = v.Code
		data.CodeType = v.OrderSource
		data.UseStatus = v.UseStatus
		data.CreateTime = v.CreateTime
		data.FinishTime = v.FinishTime
		data.DurationTime = v.DurationTime
		data.OrderNo = v.OrderNo
		data.OrderAmount = v.OrderAmount
		data.RealAmount = v.RealAmount
		data.CurrencyType = v.CurrencyType
		data.PayerAccount = v.PayerAccount
		data.PayeeAccount = v.PayeeAccount
		data.PayeeAccountType = v.PayeeAccountType
		dataList = append(dataList, data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Total = total
	reply.List = dataList
	return nil
}

/**
处理核销码（freeze冻结、unfreeze解冻、cancel注销）
*/
func (*CustHandler) DisposeWriteOffCode(ctx context.Context, req *custProto.DisposeWriteOffCodeRequest, reply *custProto.DisposeWriteOffCodeReply) error {
	if req.Code == "" {
		ss_log.Error("Code参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	writeOff, err := dao.WriteoffDaoInst.GetCodeDetailByCode(req.Code)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("核销码[%v]不存在", req.Code)
			reply.ResultCode = ss_err.ERR_WRITE_OFF_CODE_FAILD
			return nil
		}
		ss_log.Error("查询核销码[%v]失败, err=%v", req.Code, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	if writeOff.UseStatus == constants.WriteOffCodeExpired {
		ss_log.Error("核销码[%v]已过期，不能再操作", req.Code)
		reply.ResultCode = ss_err.ERR_WRITE_OFF_CODE_Expired
		return nil
	}

	if writeOff.UseStatus == constants.WriteOffCodeCancelled {
		ss_log.Error("核销码[%v]已注销，不能再操作", req.Code)
		reply.ResultCode = ss_err.ERR_WRITE_OFF_CODE_Cancelled
		return nil
	}

	//收款人账号已激活，则不能冻结已激活账号的核销码
	if writeOff.PayeeAccountType == constants.AccountActived && req.OpType == constants.WrittenOffCodeOpFreeze {
		ss_log.Error("不能对已激活账号[%v]进行此操作[%v]", writeOff.PayeeAccount, req.OpType)
		reply.ResultCode = ss_err.ERR_ACCOUNT_Actived_NOT_PERFORMED_OP
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, err := dbHandler.BeginTx(ctx, nil)
	if err != nil {
		ss_log.Error("开启事务失败")
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	var useStatus string
	switch req.OpType {
	//解冻, 核销码为已冻结状态才能进行此操作
	case constants.WrittenOffCodeOpUnFreeze:
		if writeOff.UseStatus != constants.WriteOffCodeFrozen {
			ss_log.Error("核销码当前状态为[%v],不能直接[%v]操作", writeOff.UseStatus, req.OpType)
			reply.ResultCode = ss_err.ERR_WRITE_OFF_NOT_Operate
			return nil
		}

		//如果收款人账号已激活则将钱转入已激活虚账，否则转入未激活虚账
		if writeOff.PayeeAccountType == constants.AccountActived {
			//查询虚账
			payeeAccNo, err := dao.AccDaoInstance.GetUidByAccount(writeOff.PayeeAccount)
			if err != nil {
				ss_log.Error("解冻-查询收款人账号ID失败, account=%v, err=%v", writeOff.PayeeAccount, err)
				reply.ResultCode = ss_err.ERR_SYSTEM
				return nil
			}
			vAccNoType := global.GetUserVAccType(writeOff.CurrencyType, true)
			vAccNo, err := dao.VaccountDaoInst.GetVaccountNo(payeeAccNo, strext.ToInt32(vAccNoType))
			if err != nil {
				ss_log.Error("解冻-查询收款人虚账失败, accountNo=%v, err=%v", payeeAccNo, err)
				reply.ResultCode = ss_err.ERR_SYSTEM
				return nil
			}

			//减少未激活虚账的冻结金额
			errCode := dao.VaccountDaoInst.ModifyVaccFrozenUpperZero1(tx, "-", writeOff.PayeeVAccNo, writeOff.RealAmount, writeOff.Code,
				constants.VaReason_PlatformFreeze, constants.VaOpType_Defreeze)
			if errCode != ss_err.ERR_SUCCESS {
				ss_sql.Rollback(tx)
				ss_log.Error("解冻-减少未激活虚账[%v]冻结金额失败, errCode=%v", writeOff.PayeeVAccNo, errCode)
				reply.ResultCode = errCode
				return nil
			}

			//增加已激活虚账的余额
			errCode = dao.VaccountDaoInst.ModifyVAccBalance(tx, "+", vAccNo, writeOff.RealAmount, writeOff.Code,
				constants.VaReason_PlatformFreeze, constants.VaOpType_Add)
			if errCode != ss_err.ERR_SUCCESS {
				ss_sql.Rollback(tx)
				ss_log.Error("解冻-增加已激活虚账[%v]冻结金额失败, errCode=%v", writeOff.RealAmount, errCode)
				reply.ResultCode = errCode
				return nil
			}
			useStatus = constants.WriteOffCodeIsUse
		} else {
			errCode := dao.VaccountDaoInst.ModifyVaccRemainAndFrozenUpperZero1(tx, "+", writeOff.PayeeVAccNo, writeOff.RealAmount, writeOff.Code,
				constants.VaReason_PlatformFreeze, constants.VaOpType_Defreeze_Add)
			if errCode != ss_err.ERR_SUCCESS {
				ss_sql.Rollback(tx)
				ss_log.Error("解冻-减少虚账[%v]冻结金额失败, errCode=%v", writeOff.PayeeVAccNo, errCode)
				reply.ResultCode = errCode
				return nil
			}
			useStatus = constants.WriteOffCodeWaitUse
		}

	//冻结, 核销码为已冻结状态才能进行此操作
	case constants.WrittenOffCodeOpFreeze:
		if writeOff.UseStatus != constants.WriteOffCodeWaitUse {
			ss_log.Error("核销码当前状态为[%v],不能直接[%v]操作", writeOff.UseStatus, req.OpType)
			reply.ResultCode = ss_err.ERR_WRITE_OFF_NOT_Operate
			return nil
		}
		errCode := dao.VaccountDaoInst.ModifyVaccRemainAndFrozenUpperZero1(tx, "-", writeOff.PayeeVAccNo, writeOff.RealAmount, writeOff.Code,
			constants.VaReason_PlatformFreeze, constants.VaOpType_Freeze)
		if errCode != ss_err.ERR_SUCCESS {
			ss_sql.Rollback(tx)
			ss_log.Error("冻结-增加虚账[%v]冻结金额失败, errCode=%v", writeOff.PayeeVAccNo, errCode)
			reply.ResultCode = errCode
			return nil
		}
		useStatus = constants.WriteOffCodeFrozen

	//注销, 核销码为初始化状态才能进行此操作
	case constants.WrittenOffCodeOpCancel:
		if writeOff.UseStatus != constants.WriteOffCodeFrozen {
			ss_log.Error("核销码当前状态为[%v],不能直接[%v]操作", writeOff.UseStatus, req.OpType)
			reply.ResultCode = ss_err.ERR_WRITE_OFF_NOT_Operate
			return nil
		}
		errCode := dao.VaccountDaoInst.ModifyVaccFrozenUpperZero1(tx, "-", writeOff.PayeeVAccNo, writeOff.RealAmount, writeOff.Code,
			constants.VaReason_PlatformFreeze, constants.VaOpType_Defreeze)
		if errCode != ss_err.ERR_SUCCESS {
			ss_sql.Rollback(tx)
			ss_log.Error("注销-减少虚账[%v]冻结金额失败, errCode=%v", writeOff.PayeeVAccNo, errCode)
			reply.ResultCode = errCode
			return nil
		}
		useStatus = constants.WriteOffCodeCancelled

	default:
		ss_log.Error("操作[%v]不正确", req.OpType)
		reply.ResultCode = ss_err.ERR_WRITE_OFF_CODE_Cancelled
		return nil
	}

	err = dao.WriteoffDaoInst.UpdateExpiredCodeStatus(writeOff.Code, useStatus)
	if err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("修改核销码[%v]状态[%v]失败，err=%v", writeOff.Code, useStatus, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//如果没有成功记录此次操作，视为失败
	description := fmt.Sprintf("对用户核销码进行操作[%v]", req.OpType)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Account)
	if errAddLog != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
