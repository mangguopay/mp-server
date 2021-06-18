package handler

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/dao"
	"a.a/mp-server/cust-srv/util"
	"context"
	"database/sql"
	"fmt"
)

//添加个人商家认证材料
func (c *CustHandler) AddAuthMaterialBusiness(ctx context.Context, req *custProto.AddAuthMaterialBusinessRequest, reply *custProto.AddAuthMaterialBusinessReply) error {

	//确认当前账号是已经个人实名认证的才可以申请个人商家认证
	individualAuthStatus, errGetStatus := dao.AccDaoInstance.GetAuthStatusFromUid(req.AccountUid)
	if errGetStatus != nil {
		ss_log.Error("AddAuthMaterialBusiness 查询用户实名认证的姓名失败,uid为: %s,err: %s", req.AccountUid, errGetStatus.Error())
		reply.ResultCode = ss_err.ERR_ACCOUNT_NOT_REAL_AUTH
		return nil
	}

	if individualAuthStatus != constants.AuthMaterialStatus_Passed {
		ss_log.Error("AddAuthMaterialBusiness 查询用户实名认证不是通过状态，不允许添加个人商家认证材料,uid为: %s", req.AccountUid)
		reply.ResultCode = ss_err.ERR_ACCOUNT_NOT_REAL_AUTH
		return nil
	}

	//确认当前账号是否可上传商家认证材料(只有不通过和未认证的才可以上传)
	total, errCheck := dao.AuthMaterialDaoInst.CheckAccountIndividualBusinessAuthStatusByUid(req.AccountUid)
	if errCheck != nil {
		ss_log.Error("确认账号[%v]是否可上传商家认证材料失败，err=[%v]", req.AccountUid, errCheck)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}

	if total != "0" {
		ss_log.Error("当前账户[%v]的个人商家认证状态不是 审核不通过和未认证状态，不被允许上传商家认证材料", req.AccountUid)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}

	//如果上传有营业执照，查看图片是否是库中的保存有的图片
	if req.ImgId != "" {
		if imgData, err := dao.ImageDaoInstance.GetImageUrlById(req.ImgId); err != nil || imgData.ImageId == "" {
			ss_log.Error("查询不到图片,id[%v]", req.ImgId)
			reply.ResultCode = ss_err.ERR_IMAGE_OP_FAILD
			return nil
		}
	}

	//确认经营类目存在
	industryData, errGetData := dao.BusinessIndustryDaoInst.GetBusinessIndustryDetail(req.IndustryNo)
	if errGetData != nil {
		ss_log.Error("主要行业应用[%v]查询失败,err[%v]", req.IndustryNo, errGetData)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if industryData == nil {
		ss_log.Error("主要行业应用[%v]不存在", req.IndustryNo)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if industryData.Level == constants.Businesslevel_One {
		ss_log.Error("主要行业应用[%v]的等级是一级，属性分类，不具备费率和结算周期。不允许选择。", req.IndustryNo)
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
	defer ss_sql.Rollback(tx)

	//插入实名认证资料
	_, err := dao.AuthMaterialDaoInst.AddBusinessAuthMaterial(tx, dao.BusinessAuthMaterial{
		LicenseImgNo: req.ImgId,
		AuthName:     req.AuthName,
		AuthNumber:   req.AuthNumber,
		AccountUid:   req.AccountUid,
		StartDate:    req.StartDate,

		EndDate:      req.EndDate,
		TermType:     req.TermType,
		IndustryNo:   req.IndustryNo,
		SimplifyName: req.SimplifyName,
	})
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_ADD
		return nil
	}

	//修改账号的商家认证状态为申请中
	err2 := dao.AuthMaterialDaoInst.ModifyBusinessAccountAuthStatus(tx, req.AccountUid, constants.AuthMaterialStatus_Pending)
	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		reply.ResultCode = ss_err.ERR_SYS_DB_ADD
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//添加企业商家认证材料
func (c *CustHandler) AddAuthMaterialEnterprise(ctx context.Context, req *custProto.AddAuthMaterialEnterpriseRequest, reply *custProto.AddAuthMaterialEnterpriseReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	//确认当前账号是否可上传商家认证材料(只有不通过和未认证的才可以上传)
	total, errCheck := dao.AuthMaterialDaoInst.CheckAccountBusinessAuthStatusByUid(req.AccountUid)
	if errCheck != nil {
		ss_log.Error("确认账号[%v]是否可上传商家认证材料失败，err=[%v]", req.AccountUid, errCheck)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}

	if total != "0" {
		ss_log.Error("当前账户[%v]的企业商家认证状态不是 审核不通过和未认证状态，不被允许上传商家认证材料", req.AccountUid)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}

	//插入图片
	upReg := &custProto.UploadImageRequest{
		ImageStr:     req.ImgBase64,
		AccountUid:   req.AccountUid,
		Type:         constants.UploadImage_Auth,
		AddWatermark: constants.AddWatermark_True,
	}
	upReply := &custProto.UploadImageReply{}

	if errImg := c.UploadImage(ctx, upReg, upReply); errImg != nil {
		ss_log.Error("调用上传图片接口失败,err=[%v]", errImg)
		reply.ResultCode = ss_err.ERR_SAVE_IMAGE_FAILD
		return nil
	}
	if upReply.ResultCode != ss_err.ERR_SUCCESS {
		ss_log.Error("addImg1Err=[%v]", upReply.ResultCode)
		reply.ResultCode = ss_err.ERR_SAVE_IMAGE_FAILD
		return nil
	}

	//插入实名认证资料
	_, err := dao.AuthMaterialDaoInst.AddEnterpriseBusinessAuthMaterial(tx, dao.BusinessAuthMaterial{
		LicenseImgNo: upReply.ImageId,
		AuthName:     req.AuthName,
		AuthNumber:   req.AuthNumber,
		AccountUid:   req.AccountUid,
		StartDate:    req.StartDate,

		EndDate:      req.EndDate,
		TermType:     req.TermType,
		Addr:         req.Addr,
		SimplifyName: req.SimplifyName,
	})
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_ADD
		return nil
	}

	//修改账号的商家认证状态为申请中
	err2 := dao.AuthMaterialDaoInst.ModifyBusinessAccountAuthStatus(tx, req.AccountUid, constants.AuthMaterialStatus_Pending)
	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		reply.ResultCode = ss_err.ERR_SYS_DB_ADD
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//个人商家查询 个人企业认证信息
func (*CustHandler) GetAuthMaterialBusinessDetail(ctx context.Context, req *custProto.GetAuthMaterialBusinessDetailRequest, reply *custProto.GetAuthMaterialBusinessDetailReply) error {

	whereList := []*model.WhereSqlCond{
		{Key: "amb.account_uid", Val: req.AccountUid, EqType: "="},
	}
	data, err := dao.AuthMaterialDaoInst.GetAuthMaterialBusinessDetail(whereList)
	if err != nil {
		ss_log.Error("查询账号[%v]的个人商家认证信息失败，err=[%v]", req.AccountUid, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//企业商家查询 企业商家认证信息
func (*CustHandler) GetAuthMaterialEnterpriseDetail(ctx context.Context, req *custProto.GetAuthMaterialEnterpriseDetailRequest, reply *custProto.GetAuthMaterialEnterpriseDetailReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "account_uid", Val: req.AccountUid, EqType: "="},
	}

	data, err := dao.AuthMaterialDaoInst.GetAuthMaterialEnterpriseDetail(whereList)
	if err != nil {
		ss_log.Error("查询账号[%v]的企业商家认证信息失败，err=[%v]", req.AccountUid, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetAuthMaterials(ctx context.Context, req *custProto.GetAuthMaterialsRequest, reply *custProto.GetAuthMaterialsReply) error {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "am.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "am.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "am.auth_material_no", Val: req.AuthMaterialNo, EqType: "="},
		{Key: "am.status", Val: req.Status, EqType: "="},
		{Key: "acc.account", Val: req.Account, EqType: "like"},
	})
	total, errCnt := dao.AuthMaterialDaoInst.GetCnt(whereModel.WhereStr, whereModel.Args)
	if errCnt != nil {
		ss_log.Error("err=[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//添加排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by case am.status when "+constants.AuthMaterialStatus_Pending+" then 1 end, am.create_time desc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.AuthMaterialDaoInst.GetAuthMaterials(whereModel.WhereStr, whereModel.Args)
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

func (*CustHandler) ModifyAuthMaterialStatus(ctx context.Context, req *custProto.ModifyAuthMaterialStatusRequest, reply *custProto.ModifyAuthMaterialStatusReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	status := ""
	oldStatus := "" //只有未审核状态的时候才可以修改未通过和不通过，审核通过的才可以作废
	switch req.Status {
	case constants.AuthMaterialStatus_Passed: //通过
		fallthrough
	case constants.AuthMaterialStatus_Deny: //不通过
		status = req.Status
		oldStatus = constants.AuthMaterialStatus_Pending
	case constants.AuthMaterialStatus_Appeal_Passed: //作废通过的认证材料,账号内的认证状态将修改为未认证，原来的认证材料则是作废
		status = constants.AuthMaterialStatus_UnAuth
		oldStatus = constants.AuthMaterialStatus_Passed
		errDel := dao.CardDaoInst.DeleteCard(tx, req.AccountUid)
		if errDel != nil {
			ss_log.Error("删除使用原有作废认证材料申请的卡失败，AccountUid[%v]", req.AccountUid)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}
	default:
		ss_log.Error("Status参数错误[%v]", req.Status)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//更改账号内的认证状态
	err1 := dao.AuthMaterialDaoInst.ModifyAccountIndividualAuthStatus(tx, req.AccountUid, status, oldStatus)
	if err1 != nil {
		ss_log.Error("err=[%v]", err1)
		reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
		return nil
	}

	//修改认证材料的认证状态
	err2 := dao.AuthMaterialDaoInst.ModifyAuthMaterialStatus(tx, req.AuthMaterialNo, req.Status, oldStatus)
	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
		return nil
	}

	// 添加关键操作日志
	str1, legal1 := util.GetParamZhCn(req.Status, util.AuthStatus)
	if !legal1 {
		ss_log.Error("Status %v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	account := dao.AccDaoInstance.GetAccountByUid(dbHandler, req.AccountUid)

	description := fmt.Sprintf("处理账号[%v]的认证材料[%v] 认证审核操作为[%v] ", account, req.AuthMaterialNo, str1)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Account)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//获取个人商家认证材料列表(管理后台获取)
func (*CustHandler) GetAuthMaterialBusinessList(ctx context.Context, req *custProto.GetAuthMaterialBusinessListRequest, reply *custProto.GetAuthMaterialBusinessListReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "amb.auth_material_no", Val: req.AuthMaterialNo, EqType: "="},
		{Key: "amb.status", Val: req.Status, EqType: "="},
		{Key: "amb.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "amb.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "acc.account", Val: req.Account, EqType: "like"},
	}

	total, errCnt := dao.AuthMaterialDaoInst.GetBusinessMaterialCnt(whereList)
	if errCnt != nil {
		ss_log.Error("err=[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	if total == "" || total == "0" {
		reply.Total = strext.ToInt32(total)
		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	}

	datas, err := dao.AuthMaterialDaoInst.GetBusinessMaterials(whereList, strext.ToInt(req.Page), strext.ToInt(req.PageSize))
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

//审核个人商家认证材料
func (*CustHandler) ModifyAuthMaterialBusinessStatus(ctx context.Context, req *custProto.ModifyAuthMaterialBusinessStatusRequest, reply *custProto.ModifyAuthMaterialBusinessStatusReply) error {
	//status := ""
	oldStatus := "" //只有未审核状态的时候才可以修改未通过和不通过，审核通过的才可以作废
	switch req.Status {
	case constants.AuthMaterialStatus_Passed: //通过
		fallthrough
	case constants.AuthMaterialStatus_Deny: //不通过
		//status = req.Status
		oldStatus = constants.AuthMaterialStatus_Pending
	//case constants.AuthMaterialStatus_Appeal_Passed: //作废通过的个人商家认证材料,账号内的认证状态将修改为未认证，原来的认证材料则是作废
	//	status = constants.AuthMaterialStatus_UnAuth
	//	oldStatus = constants.AuthMaterialStatus_Passed
	default:
		ss_log.Error("Status参数错误[%v]", req.Status)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//查询要审核的认证材料
	authData, err := dao.AuthMaterialDaoInst.GetAuthMaterialBusinessDetail([]*model.WhereSqlCond{
		{Key: "amb.auth_material_no", Val: req.AuthMaterialNo, EqType: "="},
	})
	if err != nil {
		ss_log.Error("查询材料信息出错，err[%v]", err)
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
	defer ss_sql.Rollback(tx)

	//更改账号内的认证状态
	err1 := dao.AuthMaterialDaoInst.ModifyAccountBusinessAuthStatus(tx, authData.AccountUid, req.Status, oldStatus)
	if err1 != nil {
		ss_log.Error("err=[%v]", err1)
		reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
		return nil
	}

	//修改认证材料的认证状态
	err2 := dao.AuthMaterialDaoInst.ModifyAuthMaterialBusinessStatus(tx, req.AuthMaterialNo, req.Status, oldStatus)
	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
		return nil
	}

	//如果是通过的
	if req.Status == constants.AuthMaterialStatus_Passed {
		//查询要添加的指定签约产品
		sceneData, errGetData := dao.BusinessSceneDaoInst.GetBusinessAutoSceneDetail()
		if errGetData != nil || sceneData.SceneNo == "" {
			ss_log.Error("查询不到要添加的指定签约产品,err=[%v]", errGetData)
			reply.ResultCode = ss_err.ERR_SYSTEM
			return nil
		}

		//查询是否有通过的同样认证名称
		if !dao.AuthMaterialDaoInst.CheckAuthNameUnique(authData.AuthName) {
			ss_log.Error("该公司名称[%v]已认证", authData.AuthName)
			reply.ResultCode = ss_err.ERR_BusinessAuthName_Unique_FAILD
			return nil
		}

		//查询是否有通过的同样简称
		if !dao.AuthMaterialDaoInst.CheckSimplifyNameUnique(authData.SimplifyName) {
			ss_log.Error("该商家简称[%v]已认证", authData.SimplifyName)
			reply.ResultCode = ss_err.ERR_BusinessSimplifyName_Unique_FAILD
			return nil
		}

		//查询是否有通过的同样注册号/机构组织代码
		if authData.AuthNumber != "" && !dao.AuthMaterialDaoInst.CheckAuthNumberUnique(authData.AuthNumber) {
			ss_log.Error("该注册号/机构组织代码[%v]已认证", authData.AuthNumber)
			reply.ResultCode = ss_err.ERR_BusinessAuthNumber_Unique_FAILD
			return nil
		}

		//做商家身份的各种初始化
		accountType := constants.AccountType_PersonalBusiness
		//添加商家business
		errAdd2, businessNo := dao.BusinessDaoInst.AddBusinessTx(tx, authData.AccountUid, authData.AuthName, authData.SimplifyName, authData.AuthNumber, accountType)
		if errAdd2 != nil {
			ss_log.Error("新增商家身份出错,err=%v", errAdd2.Error())
			reply.ResultCode = ss_err.ERR_ACCOUNT_INIT_ACCOUNT_ERR
			return nil
		}

		//添加账号与商家关联关系
		if errCode := dao.RelaAccIdenDaoInst.InsertRelaAccIden(tx, authData.AccountUid, businessNo, accountType); errCode != ss_err.ERR_SUCCESS {
			ss_log.Info("添加账号与商家关联关系出错,errCode==[%v]", errCode)
			reply.ResultCode = errCode
			return nil
		}

		//添加账号与角色关联关系
		if errCode := dao.AccDaoInstance.AuthAccountRetCode(tx, accountType, authData.AccountUid); errCode != ss_err.ERR_SUCCESS {
			ss_log.Info("添加账号与角色关联关系出错,errCode==[%v]", errCode)
			reply.ResultCode = errCode
			return nil
		}

		//初始化钱包（商家要用到的）
		dao.VaccountDaoInst.InitVaccountNoTx(tx, authData.AccountUid, constants.CURRENCY_USD, constants.VaType_USD_BUSINESS_SETTLED)
		dao.VaccountDaoInst.InitVaccountNoTx(tx, authData.AccountUid, constants.CURRENCY_USD, constants.VaType_USD_BUSINESS_UNSETTLED)
		dao.VaccountDaoInst.InitVaccountNoTx(tx, authData.AccountUid, constants.CURRENCY_KHR, constants.VaType_KHR_BUSINESS_SETTLED)
		dao.VaccountDaoInst.InitVaccountNoTx(tx, authData.AccountUid, constants.CURRENCY_KHR, constants.VaType_KHR_BUSINESS_UNSETTLED)

		//1.为添加一条指定产品签约记录做准备
		//根据code和business_channel_no查询经营类目的费率和结算周期
		rateData, errGetRate := dao.BusinessIndustryRateCycleDaoInst.GetDetail([]*model.WhereSqlCond{
			{Key: "birc.business_channel_no", Val: sceneData.BusinessChannelNo, EqType: "="},
			{Key: "birc.code", Val: authData.IndustryNo, EqType: "="},
			{Key: "birc.is_delete", Val: "0", EqType: "="},
		})

		if errGetRate != nil {
			ss_log.Error("获取行业费率、结算周期失败, IndustryNo[%v], BusinessChannelNo[%v], err[%v]", authData.IndustryNo, sceneData.BusinessChannelNo, errGetRate)
			reply.ResultCode = ss_err.ERR_NoBusinessIndustryRateCycle_FAILD
			return nil
		}

		if rateData.Id == "" {
			reply.ResultCode = ss_err.ERR_NoBusinessIndustryRateCycle_FAILD
			return nil
		}

		//签约费率 = 基础费率（行业费率） + 产品浮动费率
		rate := ss_count.Add(rateData.Rate, sceneData.FloatRate)
		cycle := rateData.Cycle

		if strext.ToInt(rate) < 0 { //小于0一律按0处理
			rate = "0"
		}

		//添加一条指定产品签约记录
		_, addSignedErr := dao.BusinessSceneSignedDaoInst.AddBusinessSignedTx(tx, authData.AccountUid, businessNo, sceneData.SceneNo, authData.IndustryNo, rate, cycle, constants.SignedStatusPassed)
		if addSignedErr != nil {
			ss_log.Error("添加指定产品签约记录出错,err=[%v]", addSignedErr)
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}

		//2.生成固定收款码（新建表存储, 那个商家， 那个签约， 码内容）
		_, addCodeErr := dao.BusinessFixedCodeDaoInst.AddBusinessFixedCodeTx(tx, authData.AccountUid, businessNo)
		if addCodeErr != nil {
			ss_log.Error("添加固定二维码出错,err=[%v]", addCodeErr)
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}

	}

	//添加WEb关键操作日志
	str1, legal1 := util.GetParamZhCn(req.Status, util.AuthStatus)
	if !legal1 {
		ss_log.Error("Status %v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	account := dao.AccDaoInstance.GetAccountByUid(dbHandler, authData.AccountUid)
	description := fmt.Sprintf("处理账号[%v]的个人商家认证材料[%v] 认证审核操作为[%v] ", account, req.AuthMaterialNo, str1)
	if errAddLog := dao.LogDaoInstance.InsertWebAccountLogTx(tx, description, req.LoginUid, constants.LogAccountWebType_Business); errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v]", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//获取企业商家认证材料列表（管理后台获取）
func (*CustHandler) GetAuthMaterialEnterpriseList(ctx context.Context, req *custProto.GetAuthMaterialEnterpriseListRequest, reply *custProto.GetAuthMaterialEnterpriseListReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "am.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "am.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "am.auth_material_no", Val: req.AuthMaterialNo, EqType: "="},
		{Key: "am.status", Val: req.Status, EqType: "="},
		{Key: "acc.account", Val: req.Account, EqType: "like"},
	}

	total, errCnt := dao.AuthMaterialDaoInst.GetEnterpriseMaterialCnt(whereList)
	if errCnt != nil {
		ss_log.Error("err=[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	datas, err := dao.AuthMaterialDaoInst.GetEnterpriseMaterials(whereList, strext.ToInt(req.Page), strext.ToInt(req.PageSize))
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

//审核企业商家认证材料
func (*CustHandler) ModifyAuthMaterialEnterpriseStatus(ctx context.Context, req *custProto.ModifyAuthMaterialEnterpriseStatusRequest, reply *custProto.ModifyAuthMaterialEnterpriseStatusReply) error {
	authData, err := dao.AuthMaterialDaoInst.GetAuthMaterialEnterpriseDetail([]*model.WhereSqlCond{
		{Key: "auth_material_no", Val: req.AuthMaterialNo, EqType: "="},
	})
	if err != nil {
		ss_log.Error("查询材料信息出错，err[%v]", err)
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
	defer ss_sql.Rollback(tx)

	//status := ""
	oldStatus := "" //只有未审核状态的时候才可以修改未通过和不通过，审核通过的才可以作废
	switch req.Status {
	case constants.AuthMaterialStatus_Passed: //通过
		//status = req.Status
		oldStatus = constants.AuthMaterialStatus_Pending

		//查询是否有通过的同样认证名称
		if !dao.AuthMaterialDaoInst.CheckAuthNameUnique(authData.AuthName) {
			ss_log.Error("该公司名称[%v]已认证", authData.AuthName)
			reply.ResultCode = ss_err.ERR_BusinessAuthName_Unique_FAILD
			return nil
		}

		//查询是否有通过的同样注册号/机构组织代码
		if !dao.AuthMaterialDaoInst.CheckAuthNumberUnique(authData.AuthNumber) {
			ss_log.Error("该注册号/机构组织代码[%v]已认证", authData.AuthNumber)
			reply.ResultCode = ss_err.ERR_BusinessAuthNumber_Unique_FAILD
			return nil
		}

		//查询是否有通过的同样简称
		if !dao.AuthMaterialDaoInst.CheckSimplifyNameUnique(authData.SimplifyName) {
			ss_log.Error("该商家简称[%v]已认证", authData.SimplifyName)
			reply.ResultCode = ss_err.ERR_BusinessSimplifyName_Unique_FAILD
			return nil
		}

		idenNo := dao.RelaAccIdenDaoInst.GetIdenFromAcc(authData.AccountUid, constants.AccountType_EnterpriseBusiness)
		if err := dao.BusinessDaoInst.UpdateBusinessInfoTx(tx, idenNo, authData.AuthName, authData.SimplifyName, authData.AuthNumber); err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}
	case constants.AuthMaterialStatus_Deny: //不通过
		//status = req.Status
		oldStatus = constants.AuthMaterialStatus_Pending
	//case constants.AuthMaterialStatus_Appeal_Passed: //作废通过的企业商家认证材料,账号内的认证状态将修改为未认证，原来的认证材料则是作废
	//	status = constants.AuthMaterialStatus_UnAuth
	//	oldStatus = constants.AuthMaterialStatus_Passed
	default:
		ss_log.Error("Status参数错误[%v]", req.Status)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//更改账号内的认证状态
	err1 := dao.AuthMaterialDaoInst.ModifyAccountBusinessAuthStatus(tx, authData.AccountUid, req.Status, oldStatus)
	if err1 != nil {
		ss_log.Error("err=[%v]", err1)
		reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
		return nil
	}

	//修改认证材料的认证状态
	err2 := dao.AuthMaterialDaoInst.ModifyAuthMaterialEnterpriseStatus(tx, req.AuthMaterialNo, req.Status, oldStatus)
	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
		return nil
	}

	str1, legal1 := util.GetParamZhCn(req.Status, util.AuthStatus)
	if !legal1 {
		ss_log.Error("Status %v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	account := dao.AccDaoInstance.GetAccountByUid(dbHandler, authData.AccountUid)

	description := fmt.Sprintf("处理账号[%v]的企业商家认证材料[%v] 认证审核操作为[%v] ", account, req.AuthMaterialNo, str1)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLogTx(tx, description, req.LoginUid, constants.LogAccountWebType_Business)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//插入用户实名认证信息
func (*CustHandler) AddAuthMaterialInfo(ctx context.Context, req *custProto.AddAuthMaterialInfoRequest, reply *custProto.AddAuthMaterialInfoReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	//确认当前账号是否可上传实名认证材料(只有不通过和未认证的才可以上传)
	total, errCheck := dao.AuthMaterialDaoInst.CheckAccountAuthStatusByUid(req.AccountUid)
	if errCheck != nil {
		ss_log.Error("确认账号[%v]是否可上传实名认证材料失败，err=[%v]", req.AccountUid, errCheck)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}

	if total == "0" {
		ss_log.Error("当前账户[%v]的实名认证状态不是 审核不通过和未认证状态，不被允许上传实名认证材料", req.AccountUid)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}

	// 确认上传的图片是存在数据库的。
	imgTotal, errImg := dao.ImageDaoInstance.CheckImageById1Id2(req.FrontImgNo, req.BackImgNo)
	if errImg != nil {
		ss_log.Error("确认图片[%v,%v]是否存在数据库时出错，err=[%v]", req.FrontImgNo, req.BackImgNo, errImg)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}

	if imgTotal != "2" { //不等于两种图片数量
		ss_log.Error("图片不存在数据库[%v,%v]，err=[%v]", req.FrontImgNo, req.BackImgNo, errImg)
		reply.ResultCode = ss_err.ERR_IMAGE_OP_FAILD
		return nil
	}

	//插入实名认证资料
	authMaterialNo, err := dao.AuthMaterialDaoInst.AddAuthMaterialInfo(tx, req.FrontImgNo, req.BackImgNo, req.AuthName, req.AuthNumber, req.AccountUid)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_ADD
		return nil
	}

	//修改账号的实名认证状态为申请中，绑定刚添加的实名认证资料id
	err2 := dao.AuthMaterialDaoInst.ModifyAccountAuthStatusAndAuthMaterialNo(tx, req.AccountUid, constants.AuthMaterialStatus_Pending, authMaterialNo)
	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		reply.ResultCode = ss_err.ERR_SYS_DB_ADD
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//查询用户实名认证信息
func (*CustHandler) GetAuthMaterialInfo(ctx context.Context, req *custProto.GetAuthMaterialInfoRequest, reply *custProto.GetAuthMaterialInfoReply) error {

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "acc.uid", Val: req.AccountUid, EqType: "="},
	})

	data, err := dao.AuthMaterialDaoInst.GetAuthMaterialDetailByAccountUid(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("查询账号[%v]的认证信息失败，err=[%v]", req.AccountUid, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetAuthMaterialBusinessUpdateList(ctx context.Context, req *custProto.GetAuthMaterialBusinessUpdateListRequest, reply *custProto.GetAuthMaterialBusinessUpdateListReply) error {

	whereList := []*model.WhereSqlCond{
		{Key: "am.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "am.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "am.auth_material_no", Val: req.AuthMaterialNo, EqType: "="},
		{Key: "am.status", Val: req.Status, EqType: "="},
		{Key: "acc.account", Val: req.Account, EqType: "="},
	}

	total, errCnt := dao.AuthMaterialDaoInst.GetAuthMaterialBusinessUpdateCnt(whereList)
	if errCnt != nil && errCnt != sql.ErrNoRows {
		ss_log.Error("查询数量出错,err=%v", errCnt)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	datas, err := dao.AuthMaterialDaoInst.GetAuthMaterialBusinessUpdateList(whereList, strext.ToInt32(req.Page), strext.ToInt32(req.PageSize))
	if err != nil {
		ss_log.Error("查询修改商家认证信息列表失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) ModifyAuthMaterialBusinessUpdateStatus(ctx context.Context, req *custProto.ModifyAuthMaterialBusinessUpdateStatusRequest, reply *custProto.ModifyAuthMaterialBusinessUpdateStatusReply) error {

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	switch req.Status {
	case constants.AuthMaterialStatus_Passed: //通过

		whereList := []*model.WhereSqlCond{
			{Key: "am.id", Val: req.Id, EqType: "="},
		}
		authData, err := dao.AuthMaterialDaoInst.GetAuthMaterialBusinessUpdateDetail(whereList)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_Audit_FAILD
			return nil
		}

		if authData.Id == "" {
			ss_log.Error("获取信息失败,Id=%v", req.Id)
			reply.ResultCode = ss_err.ERR_Audit_FAILD
			return nil
		}

		//查询是否有通过的同样简称
		if !dao.AuthMaterialDaoInst.CheckSimplifyNameUnique(authData.SimplifyName) {
			ss_log.Error("该商家简称[%v]已认证", authData.SimplifyName)
			reply.ResultCode = ss_err.ERR_BusinessSimplifyName_Unique_FAILD
			return nil
		}

		//修改认证材料的简称
		if err := dao.AuthMaterialDaoInst.ModifyAuthMaterialEnterpriseSimplifyName(tx, authData.AuthMaterialNo, authData.SimplifyName, constants.AuthMaterialStatus_Passed); err != nil {
			ss_sql.Rollback(tx)
			ss_log.Error("修改商家认证信息商家简称失败,err=[%v]", err)
			reply.ResultCode = ss_err.ERR_Audit_FAILD
			return nil
		}

		//修改商家简称
		idenNo := dao.RelaAccIdenDaoInst.GetIdenFromAcc(authData.AccountUid, authData.AccountType)
		if err := dao.BusinessDaoInst.UpdateBusinessSimplifyNameTx(tx, idenNo, authData.SimplifyName); err != nil {
			ss_sql.Rollback(tx)
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}
	case constants.AuthMaterialStatus_Deny:
	default:
		ss_log.Error("Status不合法")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//修改审核状态
	if err := dao.AuthMaterialDaoInst.ModifyAuthMaterialBusinessUpdateStatus(tx, req.Id, req.Status, req.Notes); err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("修改商家认证信息审核状态失败,err=[%v]", err)
		reply.ResultCode = ss_err.ERR_Audit_FAILD
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
