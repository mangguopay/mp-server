package handler

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/util"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/dao"
	"context"
	"fmt"
)

/**
查询指纹支付功能是否开启
*/
func (CustHandler) GetAppFingerprintOn(ctx context.Context, req *custProto.GetAppFingerprintOnRequest, reply *custProto.GetAppFingerprintOnReply) error {
	paramValue := dao.GlobalParamDaoInstance.QeuryParamValue(constants.AppFingerprintKey)
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.IsOpen = paramValue
	return nil
}

/**
用户指纹列表
*/
func (CustHandler) GetFingerprintList(ctx context.Context, req *custProto.GetFingerprintListRequest, reply *custProto.GetFingerprintListReply) error {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "acc.account", Val: req.Account, EqType: "="},
		{Key: "afs.device_uuid", Val: req.DeviceNo, EqType: "="},
		{Key: "afs.use_status", Val: req.UseStatus, EqType: "="},
		{Key: "afs.open_time", Val: req.StartTime, EqType: ">="},
		{Key: "afs.open_time", Val: req.EndTime, EqType: "<="},
	})

	total, err := dao.AppFingerprintSignDaoInst.Count(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("查询列表失败, err=%v", err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if total == 0 {
		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	}

	ss_sql.SsSqlFactoryInst.AppendWhereLimitI32(whereModel, req.PageSize, req.Page)
	list, err := dao.AppFingerprintSignDaoInst.GetList(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("查询列表失败, err=%v", err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	var dataList []*custProto.DeviceFingerprint
	for _, v := range list {
		data := &custProto.DeviceFingerprint{
			Id:        v.Id,
			Account:   v.Account,
			DeviceNo:  v.DeviceUuid,
			UseStatus: v.UseStatus,
			OpenTime:  v.OpenTime,
		}
		dataList = append(dataList, data)
	}

	reply.ResultCode = ss_err.ACErrSuccess
	reply.Total = total
	reply.List = dataList
	return nil
}

/**
关闭指纹支付功能
*/
func (CustHandler) CloseFingerprintFunction(ctx context.Context, req *custProto.CloseFingerprintFunctionRequest, reply *custProto.CloseFingerprintFunctionReply) error {
	if req.IsOpen == "" || !util.InSlice(req.IsOpen, []string{"true", "false"}) {
		ss_log.Error("IsOpen参数[%v]有误", req.IsOpen)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, txErr := dbHandler.BeginTx(context.TODO(), nil)
	if txErr != nil {
		ss_log.Error("开启事务失败，err=%v", txErr)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	err := dao.GlobalParamDaoInstance.UpdateParamValueTx(tx, req.IsOpen, constants.AppFingerprintKey)
	if err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("修改全局参数[%v]失败, err=%v", constants.AppFingerprintKey, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	description := fmt.Sprintf("对'是否开启用户指纹支付的录入功能'进行操作[%v]", req.IsOpen)
	addLogErr := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config)
	if addLogErr != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, addLogErr)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
注销单个指纹或者清楚所有数据
*/
func (CustHandler) CleanFingerprintData(ctx context.Context, req *custProto.CleanFingerprintDataRequest, reply *custProto.CleanFingerprintDataReply) error {
	if req.OpType == "" {
		ss_log.Error("OpType为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, txErr := dbHandler.BeginTx(context.TODO(), nil)
	if txErr != nil {
		ss_log.Error("开启事务失败，err=%v", txErr)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	switch req.OpType {
	case "all":
		err := dao.AppFingerprintSignDaoInst.BanAll(tx)
		if err != nil {
			ss_sql.Rollback(tx)
			ss_log.Error("清除所有数据失败")
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	case "single":
		if req.Id == "" {
			ss_log.Error("Id参数为空")
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		err := dao.AppFingerprintSignDaoInst.UpdateUseStatusSingle(tx, constants.AppFingerprintUseStatus_Disable, req.Id)
		if err != nil {
			ss_sql.Rollback(tx)
			ss_log.Error("修改记录失败, id=%v, err=%v", req.Id, err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

	default:
		ss_log.Error("OpType[%v]有误", req.OpType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	description := fmt.Sprintf("清除用户指纹数据，optype=%v", req.OpType)
	addLogErr := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Account)
	if addLogErr != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, addLogErr)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
