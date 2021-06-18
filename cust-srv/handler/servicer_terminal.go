package handler

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/model"
	"a.a/mp-server/common/ss_func"
	"a.a/mp-server/common/ss_sql"
	"context"
	"database/sql"
	"fmt"

	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/cust-srv/dao"
)

//终端列表
func (*CustHandler) GetTerminalList(ctx context.Context, req *custProto.GetTerminalListRequest, reply *custProto.GetTerminalListReply) error {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "acc.account", Val: req.ServicerAccount, EqType: "="},
		{Key: "st.terminal_number", Val: req.TerminalNumber, EqType: "="},
		{Key: "st.pos_sn", Val: req.PosSn, EqType: "="},
		{Key: "st.use_status", Val: req.UseStatus, EqType: "="},
		{Key: "st.is_delete", Val: "0", EqType: "="},
	})

	total, err := dao.ServicerTerminalDao.Cnt(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("查询终端数量失败，req=%v, err=%v", strext.ToJson(req), err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	ss_sql.SsSqlFactoryInst.AppendWhereOrderBy(whereModel, "st.create_time", false)
	ss_sql.SsSqlFactoryInst.AppendWhereLimitI32(whereModel, req.PageSize, req.Page)
	dataList, err := dao.ServicerTerminalDao.GetTerminalList(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("查询终端列表失败，err=%v", err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	var list []*custProto.Terminal
	for _, v := range dataList {
		terminal := new(custProto.Terminal)
		terminal.TerminalNo = v.TerminalNo
		terminal.TerminalNumber = v.TerminalNumber
		terminal.PosSn = v.PosSn
		terminal.UseStatus = v.UseStatus
		terminal.ServicerAccount = v.ServiceAccount
		terminal.CreateTime = v.CreateTime
		terminal.UpdateTime = v.UpdateTime
		list = append(list, terminal)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Total = strext.ToInt32(total)
	reply.List = list
	return nil
}

//增加服务商终端
func (*CustHandler) AddTerminal(ctx context.Context, req *custProto.AddTerminalRequest, reply *custProto.AddTerminalReply) error {
	errCode := ss_func.CheckCountryCode(req.CountryCode)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("CountryCode参数值不合法， CountryCode=%v", req.CountryCode)
		reply.ResultCode = ss_err.ERR_CountryCode_FAILD
		return nil
	}

	if req.ServicerAccount == "" {
		ss_log.Error("ServicerAccount参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	if req.TerminalNumber == "" {
		ss_log.Error("TerminalNumber")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	if req.PosSn == "" {
		ss_log.Error("PosSn")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//账号拼接
	req.ServicerAccount = ss_func.PreCountryCode(req.CountryCode) + req.ServicerAccount

	if req.UseStatus != constants.Status_Disable && req.UseStatus != constants.Status_Enable {
		ss_log.Error("useStatus参数值错误, useStatus=[%v]", req.UseStatus)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	serviceNo, err := dao.ServiceDaoInst.GetServiceNoByAccount(req.ServicerAccount)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("服务商[%v]不存在", req.ServicerAccount)
			reply.ResultCode = ss_err.ERR_ACC_NO_SER_FAILD
			return nil
		}
		ss_log.Error("查询服务商[%v]失败，err=%v", req.ServicerAccount, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//确保同一个(terminal_number,pos_sn)未删除，使用中的只有一个
	count, errCheck := dao.ServiceDaoInst.CheckUniqueServicerTerminal("terminal_number", req.TerminalNumber)
	if errCheck != nil {
		ss_log.Error("确认唯一terminal_number查询出错，errCheck=[%v]", errCheck)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	if count != 0 {
		if req.UseStatus == constants.Status_Enable {
			ss_log.Error("terminal_number被使用中")
			reply.ResultCode = ss_err.ERR_TERMINALNUM_IN_USE
			return nil
		}
	}
	count2, errCheck2 := dao.ServiceDaoInst.CheckUniqueServicerTerminal("pos_sn", req.PosSn)
	if errCheck2 != nil {
		ss_log.Error("errCheck2=[%v]", errCheck2)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	if count2 != 0 {
		if req.UseStatus == constants.Status_Enable {
			ss_log.Error("pos_sn被使用中")
			reply.ResultCode = ss_err.ERR_TERMINALPOSSN_IN_USE
			return nil
		}
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}

	terminalNo, errStr := dao.ServicerTerminalDao.InsertTx(tx, serviceNo, req.TerminalNumber, req.PosSn, req.UseStatus)
	if errStr != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("添加终端失败, err=%v", err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	description := fmt.Sprintf("添加服务商[%v]终端[%v]", req.ServicerAccount, terminalNo)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginAccount, constants.LogAccountWebType_Servicer)
	if errAddLog != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	ss_sql.Commit(tx)

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) UpdateTerminal(ctx context.Context, req *custProto.UpdateTerminalRequest, reply *custProto.UpdateTerminalReply) error {
	if req.TerminalNo == "" {
		ss_log.Error("TerminalNo参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.UseStatus != constants.Status_Disable && req.UseStatus != constants.Status_Enable {
		ss_log.Error("useStatus参数值错误, useStatus=[%v]", req.UseStatus)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	terminal, err := dao.ServicerTerminalDao.GetTerminalById(req.TerminalNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("终端不存在，TerminalNo=%v", req.TerminalNo)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		ss_log.Error("查询终端失败，TerminalNo=%v，err=%v", req.TerminalNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}

	if req.UseStatus == constants.Status_Disable {
		if terminal.UseStatus == constants.Status_Enable {
			if err := dao.ServicerTerminalDao.UpdateUseStatusTx(tx, req.TerminalNo, req.UseStatus); err != nil {
				ss_sql.Rollback(tx)
				ss_log.Error("修改终端使用状态失败, TerminalNo=%v, useStatus=%v, err=%v", req.TerminalNo, req.UseStatus, err)
				reply.ResultCode = ss_err.ERR_PARAM
				return nil
			}
		}
	} else {
		if terminal.UseStatus == constants.Status_Disable {
			//根据终端号查询是否有已启用的终端，有则不用再启用此终端
			terminals, err := dao.ServicerTerminalDao.GetTerminalByNumber(terminal.TerminalNumber, constants.Status_Enable)
			if err != nil {
				ss_log.Error("查询终端失败, TerminalNumber=%v, useStatus=%v, err=%v", terminal.TerminalNumber, constants.Status_Enable, err)
				reply.ResultCode = ss_err.ERR_SYSTEM
				return nil
			}
			if len(terminals) == 0 {
				if err := dao.ServicerTerminalDao.UpdateUseStatusTx(tx, req.TerminalNo, req.UseStatus); err != nil {
					ss_sql.Rollback(tx)
					ss_log.Error("修改终端使用状态失败, TerminalNo=%v, useStatus=%v, err=%v", req.TerminalNo, req.UseStatus, err)
					reply.ResultCode = ss_err.ERR_PARAM
					return nil
				}
			} else {
				ss_log.Error("终端[%v]已被使用", terminal.TerminalNumber)
				reply.ResultCode = ss_err.ERR_TERMINALNUM_IN_USE
				return nil

			}
		}
	}

	description := fmt.Sprintf("修改终端[%v]使用状态[%v]", req.TerminalNo, req.UseStatus)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginAccount, constants.LogAccountWebType_Servicer)
	if errAddLog != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeleteTerminal(ctx context.Context, req *custProto.DeleteTerminalRequest, reply *custProto.DeleteTerminalReply) error {
	if req.TerminalNo == "" {
		ss_log.Error("TerminalNo参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}

	if err := dao.ServicerTerminalDao.DeleteTerminalByIdTx(tx, req.TerminalNo); err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("删除终端失败, TerminalNo=%v, err=%v", req.TerminalNo, err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	description := fmt.Sprintf("删除终端[%v]", req.TerminalNo)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginAccount, constants.LogAccountWebType_Servicer)
	if errAddLog != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	ss_sql.Commit(tx)

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
