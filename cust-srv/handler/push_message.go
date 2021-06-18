package handler

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/dao"
	"a.a/mp-server/cust-srv/util"
	"context"
	"fmt"
	"strings"
)

/**
添加或修改消息推送模板
*/
func (c *CustHandler) InsertOrUpdatePushTemp(ctx context.Context, req *custProto.InsertOrUpdatePushTempRequest, reply *custProto.InsertOrUpdatePushTempReply) error {
	// 判断是否存在pushNo
	for _, v := range strings.Split(req.PushNos, ",") {
		if !dao.PushDaoInstance.ConfirmPusherIsExist(v) {
			ss_log.Error("添加模板配置,检测pushNO不存在,pushNo: %s", v)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	description := ""

	// 检查参数的个数
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	txt := dao.LangDaoInst.GetLangTextByKey(dbHandler, req.ContentKey, constants.LangZhCN)
	n := strext.HasN(txt, "%s")
	ss_log.Info("内容key=[%v], n=[%v]", txt, n)
	req.LenArgs = strext.ToStringNoPoint(n)

	switch req.OpType {
	case constants.PushTemp_Op_Add: // 新增
		if !dao.PushTempDaoInst.CheckPushTempNo(req.TempNo) {
			ss_log.Error("TempNo[%v]重复，不可重复插入", req.TempNo)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		if errStr := dao.PushTempDaoInst.Insert(req.TempNo, req.PushNos, req.TitleKey, req.ContentKey); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}
		description = fmt.Sprintf("插入消息推送模板配置 模板id[%v],PushNos[%v],标题key[%v],内容key[%v]", req.TempNo, req.PushNos, req.TitleKey, req.ContentKey)
	case constants.PushTemp_Op_Modify:
		if errStr := dao.PushTempDaoInst.ModifyPushTemp(req.TempNo, req.PushNos, req.TitleKey, req.ContentKey); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}
		description = fmt.Sprintf("修改消息推送模板配置 模板id[%v],PushNos[%v],标题key[%v],内容key[%v]", req.TempNo, req.PushNos, req.TitleKey, req.ContentKey)

	default:
		ss_log.Error("操作参数有误,optype: %d", req.OpType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeletePushTemp(ctx context.Context, req *custProto.DeletePushTempRequest, reply *custProto.DeletePushTempReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	err := dao.PushDaoInstance.DelPushTemp(tx, req.TempNo)
	if err != nil {
		ss_log.Error("删除失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}
	description := fmt.Sprintf("删除 消息推送模板 id[%v]", req.TempNo)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetPushTemps(ctx context.Context, req *custProto.GetPushTempsRequest, reply *custProto.GetPushTempsReply) error {

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
		//{Key: "create_time", Val: req.StartTime, EqType: ">="},
		//{Key: "create_time", Val: req.EndTime, EqType: "<="},
	})
	total, errCnt := dao.PushDaoInstance.GetPushTempsCnt(whereModel.WhereStr, whereModel.Args)
	if errCnt != nil {
		ss_log.Error("err=[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//添加排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by temp_no asc `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.PushDaoInstance.GetPushTempsInfos(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	whereModel2 := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		//{Key: "create_time", Val: req.StartTime, EqType: ">="},
		//{Key: "create_time", Val: req.EndTime, EqType: "<="},
	})
	datas2, err2 := dao.PushDaoInstance.GetPushConfInfos(whereModel2.WhereStr, whereModel2.Args)
	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	for _, v := range datas {
		if v.PushNos != "" {
			pushNos := "" //存处理完的pushNos
			pushNoArr := strings.Split(v.PushNos, ",")
			for _, pushNo := range pushNoArr {
				tempStr := pushNo
				for _, v := range datas2 {
					if v.PusherNo == pushNo {
						//tempStr = "[" + v.Pusher + "]" + pushNo
						tempStr = "[" + v.Pusher + "]"
						break
					}
				}
				if pushNos != "" {
					pushNos = fmt.Sprintf("%v,%v", pushNos, tempStr)
				} else {
					pushNos = fmt.Sprintf("%v", tempStr)
				}
			}

			v.PushNos2 = pushNos
		}
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (c *CustHandler) GetPushTemp(ctx context.Context, req *custProto.GetPushTempRequest, reply *custProto.GetPushTempReply) error {
	data, err := dao.PushTempDaoInst.GetPushTempInfoFromNo(req.TempNo)
	if err != nil {
		ss_log.Error("查询推送模板失败,templeNo为: %s,err: %s", req.TempNo, err.Error())
		reply.ResultCode = ss_err.ERR_PARAM
		return nil

	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

func (*CustHandler) GetPushConfs(ctx context.Context, req *custProto.GetPushConfsRequest, reply *custProto.GetPushConfsReply) error {

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "create_time", Val: req.StartTime, EqType: ">="},
		{Key: "create_time", Val: req.EndTime, EqType: "<="},
		{Key: "use_status", Val: req.UseStatus, EqType: "="},
		{Key: "pusher", Val: req.Pusher, EqType: "like"},
	})
	total, errCnt := dao.PushDaoInstance.GetPushConfCnt(whereModel.WhereStr, whereModel.Args)
	if errCnt != nil {
		ss_log.Error("err=[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//添加排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by create_time desc `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.PushDaoInstance.GetPushConfInfos(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetPushConf(ctx context.Context, req *custProto.GetPushConfRequest, reply *custProto.GetPushConfReply) error {

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "pusher_no", Val: req.PusherNo, EqType: "="},
	})

	data, err := dao.PushDaoInstance.GetPushConfDetail(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) InsertOrUpdatePushConfs(ctx context.Context, req *custProto.InsertOrUpdatePushConfsRequest, reply *custProto.InsertOrUpdatePushConfsReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	str1, legal1 := util.GetParamZhCn(req.UseStatus, util.UseStatus)
	if !legal1 {
		ss_log.Error("UseStatus %v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	switch req.ConditionType {
	case "0":
	case "1": //国家码
		if req.ConditionValue == "" {
			ss_log.Error("ConditionValue条件值为空")
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	default:
		ss_log.Error("ConditionType类型错误")
	}

	description := ""

	if req.PusherNo != "" {
		err := dao.PushDaoInstance.ModifyPushConf(tx, req.PusherNo, req.Pusher, req.Config, req.UseStatus, req.ConditionType, req.ConditionValue)
		if err != nil {
			ss_log.Error("修改失败，err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}
		description = fmt.Sprintf("修改消息推送配置 id:[%v],消息推送商[%v],配置信息[%v],使用状态[%v]", req.PusherNo, req.Pusher, req.Config, str1)
	} else {
		pusherNo, err := dao.PushDaoInstance.AddPushConf(tx, req.Pusher, req.Config, req.UseStatus, req.ConditionType, req.ConditionValue)
		ss_log.Info("添加成功[%v]", pusherNo)
		if err != nil {
			ss_log.Error("添加失败，err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_ADD
			return nil
		}
		description = fmt.Sprintf("插入消息推送配置 id:[%v],消息推送商[%v],配置信息[%v],使用状态[%v]", pusherNo, req.Pusher, req.Config, str1)
	}

	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeletePushConf(ctx context.Context, req *custProto.DeletePushConfRequest, reply *custProto.DeletePushConfReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	err := dao.PushDaoInstance.DelPushConf(tx, req.PusherNo)
	if err != nil {
		ss_log.Error("删除失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}
	description := fmt.Sprintf("删除 消息推送配置 id[%v]", req.PusherNo)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) ModifyPushConfStatus(ctx context.Context, req *custProto.ModifyPushConfStatusRequest, reply *custProto.ModifyPushConfStatusReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	str1, legal1 := util.GetParamZhCn(req.UseStatus, util.UseStatus)
	if !legal1 {
		ss_log.Error("UseStatus %v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	sqlStr := "update push_conf set use_status = $2 where pusher_no = $1 "
	err := ss_sql.ExecTx(tx, sqlStr, req.PusherNo, req.UseStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
		return nil
	}

	description := fmt.Sprintf("修改消息推送配置 id[%v]的使用状态为[%v]", req.PusherNo, str1)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetPushRecords(ctx context.Context, req *custProto.GetPushRecordsRequest, reply *custProto.GetPushRecordsReply) error {

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "create_time", Val: req.StartTime, EqType: ">="},
		{Key: "create_time", Val: req.EndTime, EqType: "<="},
		{Key: "status", Val: req.Status, EqType: "="},
	})
	total, errCnt := dao.PushDaoInstance.GetPushRecordCnt(whereModel.WhereStr, whereModel.Args)
	if errCnt != nil {
		ss_log.Error("err=[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//添加排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by create_time desc `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.PushDaoInstance.GetPushRecordInfos(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
