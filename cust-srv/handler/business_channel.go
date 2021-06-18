package handler

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
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

//添加支付渠道
func (c *CustHandler) AddPaymentChannel(ctx context.Context, req *custProto.AddPaymentChannelRequest, reply *custProto.AddPaymentChannelReply) error {
	if req.ChannelName == "" {
		ss_log.Error("ChannelName参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if !util.InSlice(req.ChannelType, []string{constants.ChannelTypeInner, constants.ChannelTypeOut}) {
		ss_log.Error("ChannelType参数错误, req.ChannelType=%v", req.ChannelType)
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

	d := new(dao.BusinessChannel)
	d.ChannelNo = strext.GetDailyId()
	d.ChannelName = req.ChannelName
	d.ChannelType = req.ChannelType
	d.UpstreamNo = req.UpstreamNo
	err := dao.BusinessChannelDao.InsertTx(tx, d)
	if err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("添加渠道失败，data:=%v, err=%v", strext.ToJson(d), err)
		reply.ResultCode = ss_err.ERR_SYSTEM
	}

	description := fmt.Sprintf("添加支付渠道[%v], ChannelNo=%v", req.ChannelName, d.ChannelNo)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginAccount, constants.LogAccountWebType_Business)
	if errAddLog != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	ss_sql.Commit(tx)

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.ChannelNo = d.ChannelNo
	return nil
}

//修改支付渠道
func (c *CustHandler) UpdatePaymentChannel(ctx context.Context, req *custProto.UpdatePaymentChannelRequest, reply *custProto.UpdatePaymentChannelReply) error {
	if req.ChannelNo == "" {
		ss_log.Error("ChannelNo参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.ChannelType != "" {
		if !util.InSlice(req.ChannelType, []string{constants.ChannelTypeInner, constants.ChannelTypeOut}) {
			ss_log.Error("ChannelType参数错误, req.ChannelType=%v", req.ChannelType)
			reply.ResultCode = ss_err.ERR_PARAM
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

	d := new(dao.BusinessChannel)
	d.ChannelNo = req.ChannelNo
	d.ChannelName = req.ChannelName
	d.ChannelType = req.ChannelType
	d.UpstreamNo = req.UpstreamNo
	err := dao.BusinessChannelDao.UpdateTx(tx, d)
	if err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("修改渠道失败，data:=%v, err=%v", strext.ToJson(d), err)
		reply.ResultCode = ss_err.ERR_SYSTEM
	}

	description := fmt.Sprintf("修改支付渠道[%v], req=%v", req.ChannelNo, strext.ToJson(d))
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginAccount, constants.LogAccountWebType_Business)
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

//获取所有的支付渠道
func (c *CustHandler) GetAllPaymentChannel(ctx context.Context, req *custProto.GetAllPaymentChannelRequest, reply *custProto.GetAllPaymentChannelReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "channel_name", Val: req.ChannelName, EqType: "like"},
		{Key: "create_time", Val: req.StartTime, EqType: ">="},
		{Key: "create_time", Val: req.EndTime, EqType: "<="},
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	totalNum, err := dao.BusinessChannelDao.CntChannel(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("统计渠道数量失败， err=%v", err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	if totalNum == 0 {
		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	}

	ss_sql.SsSqlFactoryInst.AppendWhereOrderBy(whereModel, "create_time", false)
	ss_sql.SsSqlFactoryInst.AppendWhereLimitI32(whereModel, req.PageSize, req.Page)
	channels, err := dao.BusinessChannelDao.GetAllChannel(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("查询渠道列表失败， err=%v", err)
		reply.ResultCode = ss_err.ERR_SYSTEM
	}

	var list []*custProto.PaymentChannel
	for _, v := range channels {
		channel := new(custProto.PaymentChannel)
		channel.ChannelNo = v.ChannelNo
		channel.ChannelName = v.ChannelName
		channel.ChannelType = v.ChannelType
		channel.UpstreamNo = v.UpstreamNo
		channel.CreateTime = v.CreateTime
		list = append(list, channel)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Channel = list
	reply.Total = strext.ToInt32(totalNum)
	return nil
}
