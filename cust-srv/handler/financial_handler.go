package handler

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/dao"
	"a.a/mp-server/cust-srv/util"
)

/**
 * 获取每天服务商对账统计列表
 */
func (*CustHandler) GetFinancialServicerCheckList(ctx context.Context, req *go_micro_srv_cust.GetFinancialServicerCheckListRequest, reply *go_micro_srv_cust.GetFinancialServicerCheckListReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	if req.StartTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.StartTime) {
			ss_log.Error("GetFinancialServicerCheckList StartTime格式不正确,StartTime: %s", req.StartTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}
	if req.EndTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.EndTime) {
			ss_log.Error("GetFinancialServicerCheckList EndTime格式不正确,StartTime: %s", req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}
	var datas []*go_micro_srv_cust.FinancialServicerCheckData
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "scl.dates", Val: req.StartTime, EqType: ">="},
		{Key: "scl.dates", Val: req.EndTime, EqType: "<="},
		{Key: "scl.servicer_no", Val: req.ServicerNo, EqType: "="},
		{Key: "acc.phone", Val: req.Phone, EqType: "like"},
		{Key: "acc.account", Val: req.Account, EqType: "like"},
		{Key: "scl.id", Val: req.Id, EqType: "like"},
		{Key: "scl.currency_type", Val: req.CurrencyType, EqType: "="},
	})

	//统计
	where := whereModel.WhereStr
	args := whereModel.Args
	var total sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM servicer_count_list scl " +
		" LEFT JOIN servicer ser ON ser.servicer_no = scl.servicer_no " +
		" LEFT JOIN account acc ON acc.uid = ser.account_no " + where
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by scl.create_time desc,scl.id desc `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT scl.servicer_no, scl.currency_type, scl.create_time, scl.in_num, scl.in_amount, scl.out_num" +
		", scl.out_amount, scl.profit_num, scl.profit_amount, scl.recharge_num, scl.recharge_amount" +
		", scl.withdraw_num, scl.withdraw_amount, scl.id, scl.dates " +
		", acc.phone, acc.account " +
		" FROM servicer_count_list scl " +
		" LEFT JOIN servicer ser ON ser.servicer_no = scl.servicer_no " +
		" LEFT JOIN account acc ON acc.uid = ser.account_no  " + where2

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	for rows.Next() {
		var data go_micro_srv_cust.FinancialServicerCheckData
		var phone, account sql.NullString
		err = rows.Scan(
			&data.ServicerNo,
			&data.CurrencyType,
			&data.CreateTime,
			&data.InNum,
			&data.InAmount,

			&data.OutNum,
			&data.OutAmount,
			&data.ProfitNum,
			&data.ProfitAmount,
			&data.RechargeNum,

			&data.RechargeAmount,
			&data.WithdrawNum,
			&data.WithdrawAmount,
			&data.Id,
			&data.Dates,
			&phone,
			&account,
		)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		if phone.String != "" {
			data.Phone = phone.String
		}
		if account.String != "" {
			data.Account = account.String
		}
		datas = append(datas, &data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * 查看指定服务商的某天账单明细
 */
func (*CustHandler) GetBillingDetailsResultsList(ctx context.Context, req *go_micro_srv_cust.GetBillingDetailsResultsListRequest, reply *go_micro_srv_cust.GetBillingDetailsResultsListReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var datas []*go_micro_srv_cust.BillingDetailsResultsData

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "bdr.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "bdr.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "bdr.order_no", Val: req.OrderNo, EqType: "like"},
		{Key: "bdr.bill_type", Val: req.BillType, EqType: "="},
		{Key: "ser.servicer_no", Val: req.ServicerNo, EqType: "="},
		{Key: "bdr.currency_type", Val: req.CurrencyType, EqType: "="},
		{Key: "bdr.order_status", Val: constants.OrderStatus_Paid, EqType: "="},
	})
	//统计
	where := whereModel.WhereStr
	args := whereModel.Args
	var total sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM billing_details_results bdr " +
		" LEFT JOIN servicer ser ON ser.account_no = bdr.account_no " + where
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by bdr.create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT bdr.create_time, bdr.bill_no, bdr.amount, bdr.currency_type, bdr.bill_type, ser.servicer_no, bdr.order_no " +
		" FROM billing_details_results bdr " +
		" LEFT JOIN servicer ser ON ser.account_no = bdr.account_no " + where2

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}
	for rows.Next() {
		var data go_micro_srv_cust.BillingDetailsResultsData
		err = rows.Scan(
			&data.CreateTime,
			&data.BillNo,
			&data.Amount,
			&data.CurrencyType,
			&data.BillType,
			&data.ServicerNo,
			&data.OrderNo,
		)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		datas = append(datas, &data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * //收款账户收款状态修改
 */
func (*CustHandler) ModifyCollectStatus(ctx context.Context, req *go_micro_srv_cust.ModifyCollectStatusRequest, reply *go_micro_srv_cust.ModifyCollectStatusReply) error {
	str1, legal1 := util.GetParamZhCn(req.SetStatus, util.UseStatus)
	if !legal1 {
		ss_log.Error("SetStatus %v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if errStr := dao.CardHeadDaoInst.ModifyCollectStatus(req.SetStatus, req.CardNo); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	accountType, errGet1 := dao.CardHeadDaoInst.GetCardHeadAccountTypeByCardNo(req.CardNo)
	if errGet1 != nil {
		ss_log.Error("根据cardNo获取accountType出错,err=[%v]", errGet1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	description := ""
	switch accountType {
	case constants.AccountType_SERVICER:
		description = fmt.Sprintf("修改平台收款账户 id[%v] 的收款状态为[%v]", req.CardNo, str1)
	case constants.AccountType_USER:
		description = fmt.Sprintf("修改银行卡充值账户 id[%v] 的收款状态为[%v]", req.CardNo, str1)
	}

	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Financial)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
 *	获取平台收款账户
 */
func (c *CustHandler) CollectionManagementList(ctx context.Context, req *go_micro_srv_cust.CollectionManagementListRequest, reply *go_micro_srv_cust.CollectionManagementListReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var datas []*go_micro_srv_cust.CollectionManagementData
	var total sql.NullString

	//查询平台总账号
	_, accPlat, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")

	switch req.AccountType { //要查看的是用户的收款账户还是服务商的收款账户
	case constants.AccountType_SERVICER:

		whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
			{Key: "ca.is_delete", Val: "0", EqType: "="},
			{Key: "ca.account_no", Val: accPlat, EqType: "="},
			{Key: "ca.account_type", Val: req.AccountType, EqType: "="},
		})

		//统计
		where := whereModel.WhereStr
		args := whereModel.Args
		sqlCnt := "SELECT count(1) " +
			" FROM card_head ca " +
			" LEFT JOIN channel ch ON ch.channel_no = ca.channel_no and ch.is_delete='0' " + where
		err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args...)
		if err != nil {
			ss_log.Error("获取统计数目出错，err=[%v]", err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		//添加limit
		ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by ca.create_time desc`)
		ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
		where2 := whereModel.WhereStr
		args2 := whereModel.Args
		sqlStr := "SELECT ch.channel_name,ca.name,ca.card_number,ca.collect_status,ca.card_no,ca.note,ca.balance_type,ca.is_defalut " +
			" FROM card_head ca " +
			" LEFT JOIN channel ch ON ch.channel_no = ca.channel_no and ch.is_delete='0' " + where2
		rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
		if stmt != nil {
			defer stmt.Close()
		}
		defer rows.Close()

		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}

		for rows.Next() {
			var data go_micro_srv_cust.CollectionManagementData
			var channelName sql.NullString
			err = rows.Scan(
				&channelName,
				&data.Name,
				&data.CardNumber,
				&data.CollectStatus,
				&data.CardNo,
				&data.Note,
				&data.BalanceType,
				&data.IsDefalut,
			)
			if err != nil {
				ss_log.Error("err=[%v]", err)
				continue
			}
			if channelName.String != "" {
				data.ChannelName = channelName.String
			}
			datas = append(datas, &data)
		}
	case constants.AccountType_USER:
		whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
			{Key: "ca.is_delete", Val: "0", EqType: "="},
			{Key: "ca.account_no", Val: accPlat, EqType: "="},
			{Key: "ca.account_type", Val: req.AccountType, EqType: "="},
		})

		//统计
		where := whereModel.WhereStr
		args := whereModel.Args
		sqlCnt := "SELECT count(1) " +
			" FROM card_head ca " +
			" LEFT JOIN channel ch ON ch.channel_no = ca.channel_no and ch.is_delete='0' " + where
		err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args...)
		if err != nil {
			ss_log.Error("获取统计数目出错，err=[%v]", err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		//添加limit
		ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by ca.create_time desc`)
		ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
		where2 := whereModel.WhereStr
		args2 := whereModel.Args
		//卡的币种是从渠道那里来的
		sqlStr := "SELECT ch.channel_name,ca.name,ca.card_number,ca.collect_status,ca.card_no,ca.note,chcu.currency_type,ca.is_defalut " +
			" FROM card_head ca " +
			" LEFT JOIN channel_cust_config chcu ON chcu.id = ca.channel_cust_config_id and chcu.is_delete='0' " +
			" LEFT JOIN channel ch ON ch.channel_no = chcu.channel_no and ch.is_delete='0' " + where2
		rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
		if stmt != nil {
			defer stmt.Close()
		}
		defer rows.Close()

		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}

		for rows.Next() {
			var data go_micro_srv_cust.CollectionManagementData
			var channelName, balanceType sql.NullString
			err = rows.Scan(
				&channelName,
				&data.Name,
				&data.CardNumber,
				&data.CollectStatus,
				&data.CardNo,
				&data.Note,
				&balanceType,
				&data.IsDefalut,
			)
			if err != nil {
				ss_log.Error("err=[%v]", err)
				continue
			}
			if balanceType.String != "" {
				data.BalanceType = balanceType.String
			}

			if channelName.String != "" {
				data.ChannelName = channelName.String
			}
			datas = append(datas, &data)
		}
	case constants.AccountType_PersonalBusiness:
		fallthrough
	case constants.AccountType_EnterpriseBusiness: //存的时候，个人、企业统一存成用于企业的收款账户
		whereList := []*model.WhereSqlCond{
			{Key: "ca.is_delete", Val: "0", EqType: "="},
			{Key: "ca.account_no", Val: accPlat, EqType: "="},
			{Key: "ca.account_type", Val: constants.AccountType_EnterpriseBusiness, EqType: "="},
		}

		total.String = dao.CardHeadDaoInst.GetHeadCardBusinessCnt(whereList)
		datasT, err := dao.CardHeadDaoInst.GetHeadCardBusinessList(whereList, strext.ToInt(req.Page), strext.ToInt(req.PageSize))
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		datas = datasT
	default:
		ss_log.Error("AccountType参数错误[%v]", req.AccountType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * 删除平台收款账户
 */
func (*CustHandler) DelectCard(ctx context.Context, req *go_micro_srv_cust.DelectCardRequest, reply *go_micro_srv_cust.DelectCardReply) error {
	accountType, errGet1 := dao.CardHeadDaoInst.GetCardHeadAccountTypeByCardNo(req.CardNo)
	if errGet1 != nil {
		ss_log.Error("根据cardNo获取accountType出错,err=[%v]", errGet1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	description := ""
	switch accountType {
	case constants.AccountType_SERVICER:
		description = fmt.Sprintf("删除平台收服务商款账户 卡id[%v]", req.CardNo)
	case constants.AccountType_USER:
		description = fmt.Sprintf("删除银行卡充值账户 id[%v]", req.CardNo)
	}

	if errStr := dao.CardHeadDaoInst.DeleteCard(req.CardNo); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Financial)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
 *	获取平台收款账户(单个)
 */
func (c *CustHandler) GetCardInfo(ctx context.Context, req *go_micro_srv_cust.GetCardInfoRequest, reply *go_micro_srv_cust.GetCardInfoReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	data := &go_micro_srv_cust.GetCardInfoData{}
	if req.CardNo == "" {
		ss_log.Error("CardNo参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	accountType, errGet1 := dao.CardHeadDaoInst.GetCardHeadAccountTypeByCardNo(req.CardNo)
	if errGet1 != nil {
		ss_log.Error("根据cardNo获取accountType出错,err=[%v]", errGet1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	switch accountType {
	case constants.AccountType_SERVICER:
		whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
			{Key: "ca.is_delete", Val: "0", EqType: "="},
			{Key: "ca.card_no", Val: req.CardNo, EqType: "="},
		})

		where2 := whereModel.WhereStr
		args2 := whereModel.Args
		sqlStr := "SELECT ch.channel_name, ch.channel_no, ca.name, ca.card_number, ca.note, ca.collect_status, ca.is_defalut, ca.balance_type " +
			" FROM card_head ca " +
			" LEFT JOIN channel ch ON ch.channel_no = ca.channel_no " + where2
		rows, stmt, err := ss_sql.QueryRowN(dbHandler, sqlStr, args2...)
		if stmt != nil {
			defer stmt.Close()
		}
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}
		err = rows.Scan(
			&data.ChannelName,
			&data.ChannelNo,
			&data.Name,
			&data.CardNumber,
			&data.Note,
			&data.CollectStatus,
			&data.IsDefalut,
			&data.BalanceType,
		)

		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}
	case constants.AccountType_USER:
		//todo 此处由于币种存储的位置不同。所以和服务商的分开查询
		whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
			{Key: "ca.is_delete", Val: "0", EqType: "="},
			{Key: "ca.card_no", Val: req.CardNo, EqType: "="},
		})

		where2 := whereModel.WhereStr
		args2 := whereModel.Args
		sqlStr := "SELECT ch.channel_name, chcu.channel_no, ca.name, ca.card_number, ca.note, ca.collect_status, ca.is_defalut, chcu.currency_type, ca.channel_cust_config_id, ca.card_no " +
			" FROM card_head ca " +
			" LEFT JOIN channel_cust_config chcu ON chcu.id = ca.channel_cust_config_id and chcu.is_delete = '0' " +
			" LEFT JOIN channel ch ON ch.channel_no = chcu.channel_no and ch.is_delete = '0' " + where2
		rows, stmt, err := ss_sql.QueryRowN(dbHandler, sqlStr, args2...)
		if stmt != nil {
			defer stmt.Close()
		}
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}
		err = rows.Scan(
			&data.ChannelName,
			&data.ChannelNo,
			&data.Name,
			&data.CardNumber,
			&data.Note,
			&data.CollectStatus,
			&data.IsDefalut,
			&data.BalanceType,
			&data.Id,
			&data.CardNo,
		)

		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}
	case constants.AccountType_PersonalBusiness:
		fallthrough
	case constants.AccountType_EnterpriseBusiness:
		whereList := []*model.WhereSqlCond{
			{Key: "ca.is_delete", Val: "0", EqType: "="},
			{Key: "ca.card_no", Val: req.CardNo, EqType: "="},
		}

		dataT, err := dao.CardHeadDaoInst.GetHeadCardBusinessDatail(whereList)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}
		data = dataT
	default:
		ss_log.Error("accountType参数错误:[%v]", accountType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

/**
 *插入或修改平台收款账户
 */
func (*CustHandler) UpdateOrInsertCard(ctx context.Context, req *go_micro_srv_cust.UpdateOrInsertCardRequest, reply *go_micro_srv_cust.UpdateOrInsertCardReply) error {
	//插入
	switch req.AccountType {
	case constants.AccountType_SERVICER:
		description := "" //插入关键操作日志的描述
		str1, legal1 := util.GetParamZhCn(req.IsDefalut, util.IsRecom)
		if !legal1 {
			ss_log.Error("IsDefalut %v", str1)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		str2, legal2 := util.GetParamZhCn(req.BalanceType, util.CurrencyType)
		if !legal2 {
			ss_log.Error("BalanceType %v", str2)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		channelName, getErr := dao.ChannelDaoInst.GetChannelNameByChannelNo(req.ChannelNo)
		if getErr != nil {
			ss_log.Error("err=[%v]", getErr)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		if req.CardNo != "" {
			//  判断卡是否存在
			if cardNo := dao.CardHeadDaoInst.QueryCardNo(req.CardNumber, req.AccountType); cardNo != "" && cardNo != req.CardNo { //当其修改成的另一个卡号和其他人的卡号相同时，不允许修改
				ss_log.Error("err=[银行卡号已存在,卡号为--->%s]", req.CardNumber)
				reply.ResultCode = ss_err.ERR_CARD_IS_EXIST
				return nil
			}
			if updateErr := dao.CardHeadDaoInst.UpdateHeadSerCard(req.CardNo, req.ChannelNo, req.Name, req.CardNumber, req.Note, req.BalanceType, req.IsDefalut); updateErr != ss_err.ERR_SUCCESS {
				ss_log.Error("UpdateErr=[%v]", updateErr)
				reply.ResultCode = updateErr
				return nil
			}

			description = fmt.Sprintf("修改平台收款账户 卡id[%v]的信息为 渠道名称[%v],收款账户名[%v],卡号[%v],币种[%v],备注[%v],是否推荐使用[%v]", req.CardNo, channelName, req.Name, req.CardNumber, req.BalanceType, req.Note, str1)

		} else {
			//  判断卡是否存在
			if cardNo := dao.CardHeadDaoInst.QueryCardNo(req.CardNumber, req.AccountType); cardNo != "" {
				ss_log.Error("err=[银行卡号已存在,卡号为--->%s]", req.CardNumber)
				reply.ResultCode = ss_err.ERR_CARD_IS_EXIST
				return nil
			}
			CardNo, insertErr2 := dao.CardHeadDaoInst.InsertHeadSerCard(req.ChannelNo, req.Name, req.CardNumber, req.Note, req.BalanceType, req.IsDefalut, req.AccountType)
			if insertErr2 != nil {
				ss_log.Error("err2=[%v]", insertErr2)
				reply.ResultCode = ss_err.ERR_SYS_DB_ADD
				return nil
			}

			description = fmt.Sprintf("插入平台收款账户 卡id[%v] 插入的信息为 渠道名称[%v],收款账户名[%v],卡号[%v],币种[%v],备注[%v],是否推荐使用[%v]", CardNo, channelName, req.Name, req.CardNumber, req.BalanceType, req.Note, str1)
		}
		errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Financial)
		if errAddLog != nil {
			ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}
	case constants.AccountType_USER: //
		description := "" //插入关键操作日志的描述
		if req.CardNo != "" {
			//  判断卡是否存在
			if cardNo := dao.CardHeadDaoInst.QueryCardNo(req.CardNumber, req.AccountType); cardNo != "" && cardNo != req.CardNo { //当其修改成的另一个卡号和其他人的卡号相同时，不允许修改
				ss_log.Error("err=[银行卡号已存在,卡号为--->%s]", req.CardNumber)
				reply.ResultCode = ss_err.ERR_CARD_IS_EXIST
				return nil
			}
			if updateErr := dao.CardHeadDaoInst.UpdateHeadUseCard(req.CardNo, req.Name, req.CardNumber, req.ChannelId); updateErr != ss_err.ERR_SUCCESS {
				ss_log.Error("UpdateErr=[%v]", updateErr)
				reply.ResultCode = updateErr
				return nil
			}

			channelName, currencyType, getErr := dao.ChannelDaoInst.GetChannelCustConfigInfoById(req.ChannelId)
			if getErr != nil {
				ss_log.Error("err=[%v]", getErr)
				reply.ResultCode = ss_err.ERR_PARAM
				return nil
			}

			description = fmt.Sprintf("修改银行卡充值账户 卡id[%v] 修改信息为 渠道[%v],收款账户名[%v],卡号[%v],币种[%v],备注[%v]", req.CardNo, channelName, req.Name, req.CardNumber, currencyType, req.Note)
		} else {
			//  判断卡是否存在
			if cardNo := dao.CardHeadDaoInst.QueryCardNo(req.CardNumber, req.AccountType); cardNo != "" {
				ss_log.Error("err=[银行卡号已存在,卡号为--->%s]", req.CardNumber)
				reply.ResultCode = ss_err.ERR_CARD_IS_EXIST
				return nil
			}

			//查询渠道的币种
			channelName, currencyType, getErr := dao.ChannelDaoInst.GetChannelCustConfigInfoById(req.ChannelId)
			if getErr != nil {
				ss_log.Error("err=[%v]", getErr)
				reply.ResultCode = ss_err.ERR_PARAM
				return nil
			}

			cardNo, insertErr2 := dao.CardHeadDaoInst.InsertHeadUseCard(req.Name, req.CardNumber, req.AccountType, req.ChannelId, currencyType)
			if insertErr2 != nil {
				ss_log.Error("err2=[%v]", insertErr2)
				reply.ResultCode = ss_err.ERR_SYS_DB_ADD
				return nil
			}

			description = fmt.Sprintf("插入新银行卡充值账户 卡id[%v] 插入的信息为 渠道名称[%v],收款账户名[%v],卡号[%v],币种[%v],备注[%v]", cardNo, channelName, req.Name, req.CardNumber, currencyType, req.Note)
		}

		errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Financial)
		if errAddLog != nil {
			ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}
	case constants.AccountType_PersonalBusiness: //
		fallthrough
	case constants.AccountType_EnterpriseBusiness: //
		description := "" //插入关键操作日志的描述
		if req.CardNo != "" {
			//  判断卡是否存在
			if cardNo := dao.CardHeadDaoInst.QueryCardNo(req.CardNumber, constants.AccountType_EnterpriseBusiness); cardNo != "" && cardNo != req.CardNo { //当其修改成的另一个卡号和其他人的卡号相同时，不允许修改
				ss_log.Error("err=[银行卡号已存在,卡号为--->%s]", req.CardNumber)
				reply.ResultCode = ss_err.ERR_CARD_IS_EXIST
				return nil
			}

			if updateErr := dao.CardHeadDaoInst.UpdateHeadBusinessCard(req.CardNo, req.Name, req.CardNumber, req.ChannelId); updateErr != ss_err.ERR_SUCCESS {
				ss_log.Error("UpdateErr=[%v]", updateErr)
				reply.ResultCode = updateErr
				return nil
			}

			channelName, currencyType, getErr := dao.ChannelDaoInst.GetBusinessChannelInfoById(req.ChannelId)
			if getErr != nil {
				ss_log.Error("err=[%v]", getErr)
				reply.ResultCode = ss_err.ERR_PARAM
				return nil
			}

			description = fmt.Sprintf("修改平台收商家款账户 卡id[%v] 修改信息为 渠道[%v],收款账户名[%v],卡号[%v],币种[%v],备注[%v]", req.CardNo, channelName, req.Name, req.CardNumber, currencyType, req.Note)
		} else {
			//  判断卡是否存在
			if cardNo := dao.CardHeadDaoInst.QueryCardNo(req.CardNumber, constants.AccountType_EnterpriseBusiness); cardNo != "" {
				ss_log.Error("err=[银行卡号已存在,卡号为--->%s]", req.CardNumber)
				reply.ResultCode = ss_err.ERR_CARD_IS_EXIST
				return nil
			}

			//查询渠道的币种
			channelName, currencyType, getErr := dao.ChannelDaoInst.GetBusinessChannelInfoById(req.ChannelId)
			if getErr != nil {
				ss_log.Error("err=[%v]", getErr)
				reply.ResultCode = ss_err.ERR_PARAM
				return nil
			}

			cardNo, insertErr2 := dao.CardHeadDaoInst.InsertHeadBusinessCard(req.Name, req.CardNumber, req.AccountType, req.ChannelId, currencyType)
			if insertErr2 != nil {
				ss_log.Error("err2=[%v]", insertErr2)
				reply.ResultCode = ss_err.ERR_SYS_DB_ADD
				return nil
			}

			description = fmt.Sprintf("插入新平台收商家款账户 卡id[%v] 插入的信息为 渠道名称[%v],收款账户名[%v],卡号[%v],币种[%v],备注[%v]", cardNo, channelName, req.Name, req.CardNumber, currencyType, req.Note)
		}

		errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Financial)
		if errAddLog != nil {
			ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}

	default:
		ss_log.Error("AccountType参数错误:[%v]", req.AccountType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
 *平台盈利统计
 */
func (*CustHandler) GetHeadquartersProfitList(ctx context.Context, req *go_micro_srv_cust.GetHeadquartersProfitListRequest, reply *go_micro_srv_cust.GetHeadquartersProfitListReply) error {
	if req.StartTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.StartTime) {
			ss_log.Error("GetHeadquartersProfitList StartTime格式不正确,StartTime: %s", req.StartTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}
	if req.EndTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.EndTime) {
			ss_log.Error("GetHeadquartersProfitList EndTime格式不正确,StartTime: %s", req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "general_ledger_no", Val: req.OrderNo, EqType: "like"},
		{Key: "create_time", Val: req.StartTime, EqType: ">="},
		{Key: "create_time", Val: req.EndTime, EqType: "<="},
		{Key: "balance_type", Val: req.BalanceType, EqType: "="},
		{Key: "profit_source", Val: req.ProfitSource, EqType: "="},
	})

	total, err := dao.HeadquartersProfitDao.CountLog(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("统计失败, err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//usd(进)统计
	usdProfitAdd := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhereList(usdProfitAdd, []*model.WhereSqlCond{
		{Key: "balance_type", Val: strings.ToLower(constants.CURRENCY_UP_USD), EqType: "="},
		{Key: "op_type", Val: constants.PlatformProfitAdd, EqType: "="},
	}, 0, 0)

	usdAddCount, usdAddSum, err := dao.HeadquartersProfitDao.CountProfit(usdProfitAdd.WhereStr, usdProfitAdd.Args)
	if err != nil {
		ss_log.Error("统计收益失败，op_type=%v, balance_type=%v, err=%v", constants.PlatformProfitAdd, constants.CURRENCY_UP_USD, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//usd(出)统计
	usdProfitMinus := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhereList(usdProfitMinus, []*model.WhereSqlCond{
		{Key: "balance_type", Val: strings.ToLower(constants.CURRENCY_UP_USD), EqType: "="},
		{Key: "op_type", Val: constants.PlatformProfitMinus, EqType: "="},
	}, 0, 0)

	usdMinusCount, usdMinusSum, err := dao.HeadquartersProfitDao.CountProfit(usdProfitMinus.WhereStr, usdProfitMinus.Args)
	if err != nil {
		ss_log.Error("统计收益失败，op_type=%v, balance_type=%v, err=%v", constants.PlatformProfitMinus, constants.CURRENCY_UP_USD, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//khr(进)统计
	khrProfitAdd := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhereList(khrProfitAdd, []*model.WhereSqlCond{
		{Key: "balance_type", Val: strings.ToLower(constants.CURRENCY_UP_KHR), EqType: "="},
		{Key: "op_type", Val: constants.PlatformProfitAdd, EqType: "="},
	}, 0, 0)

	khrAddCount, khrAddSum, err := dao.HeadquartersProfitDao.CountProfit(khrProfitAdd.WhereStr, khrProfitAdd.Args)
	if err != nil {
		ss_log.Error("统计收益失败，op_type=%v, balance_type=%v, err=%v", constants.PlatformProfitAdd, constants.CURRENCY_UP_KHR, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}
	//khr(出)统计
	khrProfitMinus := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhereList(khrProfitMinus, []*model.WhereSqlCond{
		{Key: "balance_type", Val: strings.ToLower(constants.CURRENCY_UP_KHR), EqType: "="},
		{Key: "op_type", Val: constants.PlatformProfitMinus, EqType: "="},
	}, 0, 0)

	khrMinusCount, khrMinusSum, err := dao.HeadquartersProfitDao.CountProfit(khrProfitMinus.WhereStr, khrProfitMinus.Args)
	if err != nil {
		ss_log.Error("统计收益失败，op_type=%v, balance_type=%v, err=%v", constants.PlatformProfitMinus, constants.CURRENCY_UP_KHR, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//可提现余额
	usdCashableBalance, err := dao.HeadquartersProfitCashableDaoInstance.GetCashableBalance(strings.ToLower(constants.CURRENCY_UP_USD))
	if err != nil {
		ss_log.Error("查询可提现余额失败,, balance_type=%v, err=%v", constants.CURRENCY_UP_USD, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	khrCashableBalance, err := dao.HeadquartersProfitCashableDaoInstance.GetCashableBalance(strings.ToLower(constants.CURRENCY_UP_KHR))
	if err != nil {
		ss_log.Error("查询可提现余额失败,, balance_type=%v, err=%v", constants.CURRENCY_UP_KHR, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	list, err := dao.HeadquartersProfitDao.GetList(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("查询数据列表失败， err=%v", err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	var datas []*go_micro_srv_cust.HeadquartersProfitData
	for _, v := range list {
		data := &go_micro_srv_cust.HeadquartersProfitData{
			LogNo:        v.LogNo,
			OrderNo:      v.OrderNo,
			Amount:       v.Amount,
			CreateTime:   v.CreateTime,
			OrderStatus:  v.OrderStatus,
			FinishTime:   v.FinishTime,
			BalanceType:  v.BalanceType,
			ProfitSource: v.ProfitSource,
			OpType:       v.OpType,
		}
		datas = append(datas, data)
	}

	//收益金额统计
	usdProfitCnt := ss_count.Sub(usdAddCount, usdMinusCount).String()
	usdProfitAmount := ss_count.Sub(usdAddSum, usdMinusSum).String()
	khrProfitCnt := ss_count.Sub(khrAddCount, khrMinusCount).String()
	khrProfitAmount := ss_count.Sub(khrAddSum, khrMinusSum).String()

	reply.CountData = &go_micro_srv_cust.HeadquartersProfitCountData{}
	reply.CountData.UsdProfitCount = usdProfitCnt
	reply.CountData.UsdProfitSum = usdProfitAmount
	reply.CountData.KhrProfitCount = khrProfitCnt
	reply.CountData.KhrProfitSum = khrProfitAmount
	reply.CountData.UsdCashableBalance = usdCashableBalance
	reply.CountData.KhrCashableBalance = khrCashableBalance

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	reply.Total = strext.ToInt32(total)
	return nil
}

/**
 *
 */
func (c *CustHandler) GetChannels(ctx context.Context, req *go_micro_srv_cust.GetChannelsRequest, reply *go_micro_srv_cust.GetChannelsReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var datas []*go_micro_srv_cust.ChannelDetailData

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
		{Key: "channel_type", Val: req.ChannelType, EqType: "="},
		{Key: "use_status", Val: req.UseStatus, EqType: "="},
		{Key: "channel_name", Val: req.ChannelName, EqType: "like"},
	})

	sqlCnt := "select count(1) from channel " + whereModel.WhereStr
	var total sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...); err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	if total.String == "" || total.String == "0" {
		reply.Total = strext.ToInt32(total.String)
		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
	sqlStr := "SELECT  channel_no, channel_name, create_time, use_status, logo_img_no, logo_img_no_grey, color_begin, color_end, channel_type  " +
		" FROM channel " + whereModel.WhereStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	for rows.Next() {
		data := go_micro_srv_cust.ChannelDetailData{}
		var logoImgNo, logoImgNoGrey sql.NullString
		var colorBegin, colorEnd, channelType sql.NullString
		if err = rows.Scan(
			&data.ChannelNo,
			&data.ChannelName,
			&data.CreateTime,
			&data.UseStatus,
			&logoImgNo,
			&logoImgNoGrey,
			&colorBegin,
			&colorEnd,
			&channelType,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}

		data.ColorBegin = colorBegin.String
		data.ColorEnd = colorEnd.String
		data.ChannelType = channelType.String

		if logoImgNo.String != "" {
			data.LogoImgNo = logoImgNo.String
			//查询处图片的url，使前端可显示出来
			reqImg := &go_micro_srv_cust.UnAuthDownloadImageRequest{
				ImageId: logoImgNo.String,
			}
			replyImg := &go_micro_srv_cust.UnAuthDownloadImageReply{}
			c.UnAuthDownloadImage(ctx, reqImg, replyImg)
			if replyImg.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取图片url失败")
			} else {
				data.LogoImgUrl = replyImg.ImageUrl
			}
		}

		if logoImgNoGrey.String != "" {
			data.LogoImgNoGrey = logoImgNoGrey.String
			//查询处图片的url，使前端可显示出来
			reqImg := &go_micro_srv_cust.UnAuthDownloadImageRequest{
				ImageId: logoImgNoGrey.String,
			}
			replyImg := &go_micro_srv_cust.UnAuthDownloadImageReply{}
			c.UnAuthDownloadImage(ctx, reqImg, replyImg)
			if replyImg.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取图片url失败")
			} else {
				data.LogoImgUrlGrey = replyImg.ImageUrl
			}
		}

		datas = append(datas, &data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 *
 */
func (c *CustHandler) GetPosChannels(ctx context.Context, req *go_micro_srv_cust.GetPosChannelsRequest, reply *go_micro_srv_cust.GetPosChannelsReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "cs.is_delete", Val: "0", EqType: "="},
		{Key: "ch.channel_name", Val: req.ChannelName, EqType: "like"},
	})

	sqlCnt := "select count(1) from channel_servicer cs " +
		" left join channel ch on ch.channel_no = cs.channel_no  " + whereModel.WhereStr
	var total sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by cs.create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT  cs.channel_no, cs.create_time, cs.use_status, cs.idx, cs.is_recom, cs.currency_type, ch.logo_img_no, ch.logo_img_no_grey, ch.channel_name, cs.id " +
		" FROM channel_servicer cs " +
		" left join channel ch on ch.channel_no = cs.channel_no  " + where2
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	var datas []*go_micro_srv_cust.PosChannelData
	for rows.Next() {
		data := go_micro_srv_cust.PosChannelData{}
		var logoImgNo, logoImgNoGrey sql.NullString
		if err = rows.Scan(
			&data.ChannelNo,
			&data.CreateTime,
			&data.UseStatus,
			&data.Idx,
			&data.IsRecom,

			&data.CurrencyType,
			&logoImgNo,
			&logoImgNoGrey,
			&data.ChannelName,
			&data.Id,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}

		if logoImgNo.String != "" {
			data.LogoImgNo = logoImgNo.String
			//查询处图片的url，使前端可显示出来
			reqImg := &go_micro_srv_cust.UnAuthDownloadImageRequest{
				ImageId: logoImgNo.String,
			}
			replyImg := &go_micro_srv_cust.UnAuthDownloadImageReply{}
			c.UnAuthDownloadImage(ctx, reqImg, replyImg)
			if replyImg.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取图片url失败")
			} else {
				data.LogoImgUrl = replyImg.ImageUrl
			}
		}
		if logoImgNoGrey.String != "" {
			data.LogoImgNoGrey = logoImgNoGrey.String
			//查询处图片的url，使前端可显示出来
			reqImg := &go_micro_srv_cust.UnAuthDownloadImageRequest{
				ImageId: logoImgNoGrey.String,
			}
			replyImg := &go_micro_srv_cust.UnAuthDownloadImageReply{}
			c.UnAuthDownloadImage(ctx, reqImg, replyImg)
			if replyImg.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取图片url失败")
			} else {
				data.LogoImgUrlGrey = replyImg.ImageUrl
			}
		}
		datas = append(datas, &data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

//
///**
// *平台盈利提现
// */
//func (*CustHandler) InsertHeadquartersProfitWithdraw(ctx context.Context, req *go_micro_srv_cust.InsertHeadquartersProfitWithdrawRequest, reply *go_micro_srv_cust.InsertHeadquartersProfitWithdrawReply) error {
//	dbHandler := db.GetDB(constants.DB_CRM)
//	defer db.PutDB(constants.DB_CRM, dbHandler)
//
//	tx, _ := dbHandler.BeginTx(ctx, nil)
//	defer ss_sql.Rollback(tx)
//
//	//查询平台盈利可提现余额是否满足此次提现
//	amount ,errGet1:=dao.HeadquartersProfitWithdrawDaoInst.GetProfitCashable(tx, req.CurrencyType)
//	if errGet1 !=ss_err.ERR_SUCCESS{
//		ss_log.Error("查询平台盈利余额失败",errGet1)
//	}
//
//	// 判断金额是否包含小数点
//	if req.CurrencyType == constants.CURRENCY_USD {
//		if strings.Contains(req.Amount, ".") {
//			reply.ResultCode = ss_err.ERR_PAY_AMOUNT_IS_NO_INTEGER
//			return nil
//		}
//	}
//
//	if strext.ToFloat64(amount) < strext.ToFloat64(req.Amount){
//		ss_log.Error("平台盈利可提现余额不足此次提现")
//		reply.ResultCode = ss_err.ERR_PROFITC_ASHABLE_FAILD
//		return nil
//	}
//
//	vaType := 0
//	//更新提现账号的余额
//	switch req.CurrencyType {
//	case "usd":
//		vaType = constants.VaType_USD_DEBIT
//	case "khr":
//		vaType = constants.VaType_KHR_DEBIT
//	default:
//		reply.ResultCode = ss_err.ERR_PARAM
//		return nil
//	}
//
//	recvVaccNo := b.confirmExistVaccount(req.AccountNo, req.CurrencyType, strext.ToInt32(vaType))
//
//	//添加平台盈利提现日志
//	sqlInsert := "insert into headquarters_profit_withdraw(order_no, currency_type, amount, note, create_time, account_no) " +
//		" values($1,$2,$3,$4,current_timestamp,$5)"
//	err := ss_sql.ExecTx(tx, sqlInsert, strext.GetDailyId(), req.CurrencyType, req.Amount, req.Note, req.AccountNo)
//	if err != nil {
//		ss_log.Error("err=[%v]", err)
//		reply.ResultCode = ss_err.ERR_PARAM
//	}
//
//	tx.Commit()
//	reply.ResultCode = ss_err.ERR_SUCCESS
//	return nil
//}

/**
 *
 */
func (*CustHandler) HeadquartersProfitWithdraws(ctx context.Context, req *go_micro_srv_cust.HeadquartersProfitWithdrawsRequest, reply *go_micro_srv_cust.HeadquartersProfitWithdrawsReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var datas []*go_micro_srv_cust.HeadquartersProfitWithdrawsData

	ss_log.Info("FinanciaHandler | HeadquartersProfitWithdraws req==[%v]", req)
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		//"page", "page_size", "order_no", "start_time", "end_time", "currency_type"
		{Key: "hpw.order_no", Val: req.OrderNo, EqType: "="},
		{Key: "hpw.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "hpw.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "hpw.currency_type", Val: req.CurrencyType, EqType: "="},
	})
	whereCnt := whereModel.WhereStr
	argsCnt := whereModel.Args

	var total sql.NullString
	sqlCnt := "select count(1) from headquarters_profit_withdraw hpw " +
		" LEFT JOIN account acc ON hpw.account_no = acc.uid	and acc.is_delete = '0' " + whereCnt
	cntErr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, argsCnt...)
	if cntErr != nil {
		ss_log.Error("cntErr=[%v]", cntErr)
	}

	var usdCnt, usdSum, khrCnt, khrSum sql.NullString
	sqlUsdCnt := "select count(1),sum(hpw.amount) from headquarters_profit_withdraw hpw " +
		" LEFT JOIN account acc ON hpw.account_no = acc.uid	and acc.is_delete = '0' " + whereCnt + " and hpw.currency_type = 'usd' "
	usdCntErr := ss_sql.QueryRow(dbHandler, sqlUsdCnt, []*sql.NullString{&usdCnt, &usdSum}, argsCnt...)
	if usdCntErr != nil {
		ss_log.Error("usdCntErr=[%v]", usdCntErr)
	}
	sqlKhrCnt := "select count(1),sum(hpw.amount) from headquarters_profit_withdraw hpw " +
		" LEFT JOIN account acc ON hpw.account_no = acc.uid	and acc.is_delete = '0' " + whereCnt + " and hpw.currency_type = 'khr'"
	usdKhrErr := ss_sql.QueryRow(dbHandler, sqlKhrCnt, []*sql.NullString{&khrCnt, &khrSum}, argsCnt...)
	if usdKhrErr != nil {
		ss_log.Error("usdKhrErr=[%v]", usdKhrErr)
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` ORDER BY hpw.create_time DESC `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT hpw.order_no, hpw.currency_type, hpw.amount, hpw.note, hpw.create_time, acc.nickname " +
		" FROM headquarters_profit_withdraw hpw " +
		" LEFT JOIN account acc ON hpw.account_no = acc.uid	and acc.is_delete = '0' " + where2
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("FinanciaHandler | HeadquartersProfitWithdraws | err=%v\nreq=[%v]\nsql=[%v]", err, req, sqlStr)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	for rows.Next() {
		var data go_micro_srv_cust.HeadquartersProfitWithdrawsData
		if err = rows.Scan(
			&data.OrderNo,
			&data.CurrencyType,
			&data.Amount,
			&data.Note,
			&data.CreateTime,
			&data.Nickname,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		datas = append(datas, &data)
	}

	reply.CountData = &go_micro_srv_cust.CountHeadquartersProfitWithdrawsData{}

	if usdCnt.String != "" {
		reply.CountData.UsdCnt = strext.ToInt32(usdCnt.String)
	}
	if usdSum.String != "" {
		reply.CountData.UsdSum = strext.ToInt32(usdSum.String)
	}
	if khrCnt.String != "" {
		reply.CountData.KhrCnt = strext.ToInt32(khrCnt.String)
	}
	if khrSum.String != "" {
		reply.CountData.KhrCnt = strext.ToInt32(khrSum.String)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * POS查看服务商的某天账单明细
 */
func (*CustHandler) GetServicerBillingDetails(ctx context.Context, req *go_micro_srv_cust.GetServicerBillingDetailsRequest, reply *go_micro_srv_cust.GetServicerBillingDetailsReply) error {
	var datas []*go_micro_srv_cust.ServicerBillingDetailsData
	//如果是收银员来查询
	if req.AccountType == constants.AccountType_POS { //现pos机店员看不到账单明细
		reply.ResultCode = ss_err.ERR_SYS_NO_API_AUTH
		return nil
		//cashierNo := dao.RelaAccIdenDaoInst.GetIdenFromAcc(req.AccountNo, constants.AccountType_POS)
		//if cashierNo == "" {
		//	ss_log.Error("店员账号查询店员id出错,店员账号uid---------------->%s", req.AccountNo)
		//	reply.ResultCode = ss_err.ERR_PARAM
		//	return nil
		//}
		//
		//err, servicerNo := dao.ServiceDaoInst.GetServiceByCashierNo(cashierNo)
		//if err != nil {
		//	ss_log.Error("Cashier is no servicer /nreq=[%v]", req)
		//	reply.ResultCode = ss_err.ERR_DB_OP_SER
		//	return nil
		//}
		//
		//accountNo, err2 := dao.ServiceDaoInst.GetAccountNoByServicerNo(servicerNo)
		//if err2 != ss_err.ERR_SUCCESS {
		//	ss_log.Error("查询店员的服务商账号失败,err=[%v]", err2)
		//	reply.ResultCode = err2
		//	return nil
		//}
		//req.AccountNo = accountNo
	}

	datas, total, err := dao.ServiceDaoInst.GetServicerBillingDetails(req.AccountNo, constants.AccountType_SERVICER, req.CurrencyType, req.Page, req.PageSize)
	if err != ss_err.ERR_SUCCESS {
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	return nil
}

/**
 *
 */
func (c *CustHandler) GetUseChannels(ctx context.Context, req *go_micro_srv_cust.GetUseChannelsRequest, reply *go_micro_srv_cust.GetUseChannelsReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "cs.is_delete", Val: "0", EqType: "="},
		{Key: "ch.channel_name", Val: req.ChannelName, EqType: "like"},
		{Key: "cs.currency_type", Val: req.CurrencyType, EqType: "="},
		{Key: "cs.channel_type", Val: req.ChannelType, EqType: "="},
	})

	sqlCnt := "select count(1) from channel_cust_config cs " +
		" left join channel ch on ch.channel_no = cs.channel_no " + whereModel.WhereStr
	var total sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by cs.create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
	sqlStr := "SELECT  cs.channel_no, cs.create_time, cs.use_status, cs.currency_type, ch.logo_img_no, ch.logo_img_no_grey, ch.channel_name, cs.id" +
		", cs.save_rate, cs.withdraw_rate, cs.withdraw_max_amount, cs.save_single_min_fee, cs.withdraw_single_min_fee " +
		", cs.save_charge_type, cs.withdraw_charge_type, cs.support_type, cs.save_max_amount, cs.channel_type " +
		" FROM channel_cust_config cs " +
		" left join channel ch on ch.channel_no = cs.channel_no  " + whereModel.WhereStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	var datas []*go_micro_srv_cust.UseChannelData
	for rows.Next() {
		data := go_micro_srv_cust.UseChannelData{}
		var logoImgNo, logoImgNoGrey, channelType sql.NullString
		if err = rows.Scan(
			&data.ChannelNo,
			&data.CreateTime,
			&data.UseStatus,
			&data.CurrencyType,
			&logoImgNo,
			&logoImgNoGrey,
			&data.ChannelName,
			&data.Id,

			&data.SaveRate,
			&data.WithdrawRate,
			&data.WithdrawMaxAmount,
			&data.SaveSingleMinFee,
			&data.WithdrawSingleMinFee,

			&data.SaveChargeType,
			&data.WithdrawChargeType,
			&data.SupportType,
			&data.SaveMaxAmount,
			&channelType,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}

		data.ChannelType = channelType.String

		if logoImgNo.String != "" {
			data.LogoImgNo = logoImgNo.String
			//查询处图片的url，使前端可显示出来
			reqImg := &go_micro_srv_cust.UnAuthDownloadImageRequest{
				ImageId: logoImgNo.String,
			}
			replyImg := &go_micro_srv_cust.UnAuthDownloadImageReply{}
			c.UnAuthDownloadImage(ctx, reqImg, replyImg)
			if replyImg.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取图片url失败")
			} else {
				data.LogoImgUrl = replyImg.ImageUrl
			}
		}
		if logoImgNoGrey.String != "" {
			data.LogoImgNoGrey = logoImgNoGrey.String
			//查询处图片的url，使前端可显示出来
			reqImg := &go_micro_srv_cust.UnAuthDownloadImageRequest{
				ImageId: logoImgNoGrey.String,
			}
			replyImg := &go_micro_srv_cust.UnAuthDownloadImageReply{}
			c.UnAuthDownloadImage(ctx, reqImg, replyImg)
			if replyImg.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取图片url失败")
			} else {
				data.LogoImgUrlGrey = replyImg.ImageUrl
			}
		}
		datas = append(datas, &data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

func (c *CustHandler) GetUseChannelDetail(ctx context.Context, req *go_micro_srv_cust.GetUseChannelDetailRequest, reply *go_micro_srv_cust.GetUseChannelDetailReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "cs.is_delete", Val: "0", EqType: "="},
		{Key: "cs.id", Val: req.Id, EqType: "="},
	})

	sqlStr := "SELECT  cs.channel_no, cs.create_time, cs.use_status, cs.currency_type, ch.logo_img_no, ch.channel_name, cs.id" +
		", cs.save_rate, cs.withdraw_rate, cs.withdraw_max_amount, cs.save_single_min_fee, cs.withdraw_single_min_fee " +
		", cs.save_charge_type, cs.withdraw_charge_type, cs.support_type, cs.save_max_amount, cs.channel_type " +
		" FROM channel_cust_config cs " +
		" left join channel ch on ch.channel_no = cs.channel_no  " + whereModel.WhereStr
	rows, stmt, err := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}

	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	data := &go_micro_srv_cust.UseChannelData{}
	var logoImgNo, channelType sql.NullString
	if err = rows.Scan(
		&data.ChannelNo,
		&data.CreateTime,
		&data.UseStatus,
		&data.CurrencyType,
		&logoImgNo,
		&data.ChannelName,
		&data.Id,

		&data.SaveRate,
		&data.WithdrawRate,
		&data.WithdrawMaxAmount,
		&data.SaveSingleMinFee,
		&data.WithdrawSingleMinFee,

		&data.SaveChargeType,
		&data.WithdrawChargeType,
		&data.SupportType,
		&data.SaveMaxAmount,
		&channelType,
	); err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	data.LogoImgNo = logoImgNo.String
	data.ChannelType = channelType.String

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

/**
 *
 */
func (c *CustHandler) GetBusinessChannels(ctx context.Context, req *go_micro_srv_cust.GetBusinessChannelsRequest, reply *go_micro_srv_cust.GetBusinessChannelsReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "cbc.is_delete", Val: "0", EqType: "="},
		{Key: "cbc.currency_type", Val: req.CurrencyType, EqType: "="},
		{Key: "cbc.use_status", Val: req.UseStatus, EqType: "="},
		{Key: "cbc.channel_type", Val: req.ChannelType, EqType: "="},
		{Key: "ch.channel_name", Val: req.ChannelName, EqType: "like"},
	}

	total := dao.ChannelDaoInst.GetBusinessChannelCnt(whereList)

	datas, err := dao.ChannelDaoInst.GetBusinessChannelList(whereList, req.Page, req.PageSize)
	if err != nil {
		ss_log.Error("err=%v", err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	return nil
}

func (c *CustHandler) GetBusinessChannelDetail(ctx context.Context, req *go_micro_srv_cust.GetBusinessChannelDetailRequest, reply *go_micro_srv_cust.GetBusinessChannelDetailReply) error {
	data, err := dao.ChannelDaoInst.GetBusinessChannelDetail(req.Id)
	if err != nil {
		ss_log.Error("err=%v", err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

func (*CustHandler) InsertBusinessChannel(ctx context.Context, req *go_micro_srv_cust.InsertBusinessChannelRequest, reply *go_micro_srv_cust.InsertBusinessChannelReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	str1, legal1 := util.GetParamZhCn(req.CurrencyType, util.CurrencyType)
	if !legal1 {
		ss_log.Error("CurrencyType %v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	str2, legal2 := util.GetParamZhCn(req.SupportType, util.SupportType)
	if !legal2 {
		ss_log.Error("SupportType %v", str2)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	str3, legal3 := util.GetParamZhCn(req.SaveChargeType, util.ChargeType)
	if !legal3 {
		ss_log.Error("SaveChargeType %v", str3)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	str4, legal4 := util.GetParamZhCn(req.WithdrawChargeType, util.ChargeType)
	if !legal4 {
		ss_log.Error("WithdrawChargeType %v", str4)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	channelData, errGet := dao.ChannelDaoInst.GetChannelDetail(req.ChannelNo)
	if errGet != nil {
		ss_log.Error("[%v]渠道异常，查询渠道名称失败  err=[%v]", req.ChannelNo, errGet)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	str5, legal5 := util.GetParamZhCn(channelData.ChannelType, util.ChannelType)
	if !legal5 {
		ss_log.Error("ChannelType %v", str5)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	saveRate := ss_count.Div(req.SaveRate, "100").String() + "%"
	saveSingleMinFee := req.SaveSingleMinFee
	saveMaxAmount := req.SaveMaxAmount
	withdrawRate := ss_count.Div(req.WithdrawRate, "100").String() + "%"
	withdrawSingleMinFee := req.WithdrawSingleMinFee
	withdrawMaxAmount := req.WithdrawMaxAmount

	if req.CurrencyType == "usd" {
		saveSingleMinFee = ss_count.Div(req.SaveSingleMinFee, "100").String()
		saveMaxAmount = ss_count.Div(req.SaveMaxAmount, "100").String()
		withdrawSingleMinFee = ss_count.Div(req.WithdrawSingleMinFee, "100").String()
		withdrawMaxAmount = ss_count.Div(req.WithdrawMaxAmount, "100").String()
	}

	description := fmt.Sprintf(" 渠道名称[%v],币种[%v],渠道业务类型[%v],存款计算手续费类型[%v],取款计算手续费类型[%v],渠道类型[%v]", channelData.ChannelName, str1, str2, str3, str4, str5)
	description = fmt.Sprintf("%v,存款手续费率[%v],存款单笔手续费[%v],存款单笔最大金额[%v]", description, saveRate, saveSingleMinFee, saveMaxAmount)
	description = fmt.Sprintf("%v,取款手续费率[%v],取款单笔手续费[%v],取款单笔最大金额[%v]", description, withdrawRate, withdrawSingleMinFee, withdrawMaxAmount)

	//InsertBusinessChannel
	if req.Id == "" {
		//确认是否有同样channelNo与CurrencyType的记录
		if dao.ChannelDaoInst.CheckBusinessChannel(req.ChannelNo, req.CurrencyType) {
			ss_log.Error("银行卡存取款渠道存在相同币种的渠道")
			reply.ResultCode = ss_err.ERR_UseChannel_FAILD
			return nil
		}

		id, err2 := dao.ChannelDaoInst.AddBusinessChannel(tx, dao.BusinessChannelData{
			ChannelNo:            req.ChannelNo,
			CurrencyType:         req.CurrencyType,
			SupportType:          req.SupportType,
			ChannelType:          channelData.ChannelType,
			SaveRate:             req.SaveRate,
			SaveSingleMinFee:     req.SaveSingleMinFee,
			SaveMaxAmount:        req.SaveMaxAmount,
			SaveChargeType:       req.SaveChargeType,
			WithdrawRate:         req.WithdrawRate,
			WithdrawSingleMinFee: req.WithdrawSingleMinFee,
			WithdrawMaxAmount:    req.WithdrawMaxAmount,
			WithdrawChargeType:   req.WithdrawChargeType,
		})
		if err2 != nil {
			ss_log.Error("err2=[%v]", err2)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		description = fmt.Sprintf("插入新商家存取款渠道 id[%v], %v ", id, description)
	} else {
		//确认是否有同样channelNo与CurrencyType的记录
		if dao.ChannelDaoInst.CheckBusinessChannel(req.ChannelNo, req.CurrencyType) {
			if id := dao.ChannelDaoInst.GetBusinessChannelId(req.ChannelNo, req.CurrencyType); id != "" && id != req.Id {
				ss_log.Error("商家存取款渠道存在相同币种的渠道")
				reply.ResultCode = ss_err.ERR_UseChannel_FAILD
				return nil
			}
		}

		err2 := dao.ChannelDaoInst.ModifyBusinessChannel(tx, req.Id, req.SupportType, req.SaveRate, req.SaveSingleMinFee,
			req.SaveMaxAmount, req.SaveChargeType, req.WithdrawRate, req.WithdrawSingleMinFee, req.WithdrawMaxAmount, req.WithdrawChargeType)
		if err2 != ss_err.ERR_SUCCESS {
			ss_log.Error("err2=[%v]", err2)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		description = fmt.Sprintf("修改商家存取款渠道 id[%v], %v", req.Id, description)
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

func (*CustHandler) DeleteBusinessChannel(ctx context.Context, req *go_micro_srv_cust.DeleteBusinessChannelRequest, reply *go_micro_srv_cust.DeleteBusinessChannelReply) error {
	if req.Id == "" {
		ss_log.Error("参数Id为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if err := dao.ChannelDaoInst.DeleteBusinessChannelById(req.Id); err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}

	description := fmt.Sprintf("删除旧商家存取款渠道 id[%v]", req.Id)

	if errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config); errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) ModifyBusinessChannelStatus(ctx context.Context, req *go_micro_srv_cust.ModifyBusinessChannelStatusRequest, reply *go_micro_srv_cust.ModifyBusinessChannelStatusReply) error {
	if req.Id == "" {
		ss_log.Error("Id参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	str1, legal1 := util.GetParamZhCn(req.UseStatus, util.UseStatus)
	if !legal1 {
		ss_log.Error("UseStatus %v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if err := dao.ChannelDaoInst.ModifyBusinessChannelStatusById(req.Id, req.UseStatus); err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	description := fmt.Sprintf("修改旧银行卡存取款渠道 id[%v] 的渠道状态为[%v]", req.Id, str1)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
