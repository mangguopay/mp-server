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
	"a.a/mp-server/common/ss_rsa"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/common"
	"a.a/mp-server/cust-srv/dao"
	"a.a/mp-server/cust-srv/util"
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
)

/**
商家应用列表
*/
func (c *CustHandler) GetBusinessAppList(ctx context.Context, req *custProto.GetBusinessAppListRequest, reply *custProto.GetBusinessAppListReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "ba.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "ba.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "ba.status", Val: constants.BusinessAppStatus_Delete, EqType: "!="}, //不是删除的
		{Key: "ba.business_no", Val: req.IdenNo, EqType: "="},
		{Key: "acc.account", Val: req.Account, EqType: "like"}, //商家账号
		{Key: "ba.app_id", Val: req.AppId, EqType: "like"},
	}

	total := dao.BusinessAppDaoInst.GetBusinessAppCnt(whereList)

	datas, err := dao.BusinessAppDaoInst.GetBusinessAppList(whereList, req.Page, req.PageSize)
	if err != nil {
		ss_log.Error("查询商家应用数据列表失败，req=[%+v],err=[%v]", req, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	for _, data := range datas {
		ids := []string{
			data.SmallImgNo,
			data.BigImgNo,
		}
		imgurls := dao.ImageDaoInstance.GetImgUrlsByImgIds(ids)
		data.SmallImgUrl = imgurls[0]
		data.BigImgUrl = imgurls[1]
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
查询商家应用详情
*/
func (c *CustHandler) GetBusinessAppDetail(ctx context.Context, req *custProto.GetBusinessAppDetailRequest, reply *custProto.GetBusinessAppDetailReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "ba.business_no", Val: req.IdenNo, EqType: "="},
		{Key: "ba.status", Val: constants.BusinessAppStatus_Delete, EqType: "!="}, //不是删除的
		{Key: "ba.app_id", Val: req.AppId, EqType: "="},
	}
	data, err := dao.BusinessAppDaoInst.GetBusinessAppDetail(whereList)
	if err != nil {
		ss_log.Error("查询商家应用数据列表失败，req=[%+v],err=[%v]", req, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	ids := []string{
		data.SmallImgNo,
		data.BigImgNo,
	}

	var base64Str []string
	for _, v := range ids {
		reqImg := &custProto.UnAuthDownloadImageBase64Request{
			ImageId: v,
		}
		replyImg := &custProto.UnAuthDownloadImageBase64Reply{}
		c.UnAuthDownloadImageBase64(ctx, reqImg, replyImg)
		if replyImg.ResultCode != ss_err.ERR_SUCCESS {
			ss_log.Error("获取图片url失败")
		}
		base64Str = append(base64Str, replyImg.ImageBase64)
	}

	data.SmallImgBase64 = base64Str[0]
	data.BigImgBase64 = base64Str[1]

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
添加或修改商家应用
*/
func (c *CustHandler) InsertOrUpdateBusinessApp(ctx context.Context, req *custProto.InsertOrUpdateBusinessAppRequest, reply *custProto.InsertOrUpdateBusinessAppReply) error {

	// 商户公钥需去头去尾
	req.BusinessPubKey = ss_rsa.StripRSAKey(req.BusinessPubKey)

	if req.AppId == "" { //新应用添加

		//判断是否已实名认证
		authInfo, err := dao.AccDaoInstance.GetAuthInfoByAccountNo(req.AccountUid, req.AccountType)
		if err != nil && err != sql.ErrNoRows {
			ss_log.Error("查询商家实名认证状态失败，accountNo=%v, err=%v", req.AccountUid, err)
			reply.ResultCode = ss_err.ERR_SYSTEM
			return nil
		}

		if authInfo == nil {
			ss_log.Error("查询商家未实名认证，accountNo=%v, err=%v", req.AccountUid, err)
			reply.ResultCode = ss_err.ERR_ACCOUNT_NOT_REAL_AUTH
			return nil
		}
		if authInfo.AuthStatus != constants.AuthMaterialStatus_Passed {
			ss_log.Error("账号实名认证未通过，accountNo=%v, AuthStatus=%v", req.AccountUid, authInfo.AuthStatus)
			reply.ResultCode = ss_err.ERR_HaveNotPass_BusinessRealName_FAILD
			return nil
		}

		if _, err := dao.BusinessAppDaoInst.AddBusinessApp(dao.BusinessAppDao{
			BusinessNo: req.IdenNo,
			Status:     constants.BusinessAppStatus_Pending,
			ApplyType:  req.ApplyType, //应用类型 1-移动应用，2-网页应用
			AppName:    req.AppName,   //应用名称
			Describe:   req.Describe,  //应用描述

			SmallImgNo: req.SmallImgNo, //小图标id
			BigImgNo:   req.BigImgNo,   //大图标id
			//BusinessPubKey: req.BusinessPubKey, //商家公钥
			//SignMethod:     req.SignMethod,     //签名方式
			//IpwhiteList: req.IpWhiteList,
		}); err != nil {
			ss_log.Error("添加应用失败，req=[%+v],err=[%v]", req, err)
			reply.ResultCode = ss_err.ERR_SYS_DB_ADD
			return nil
		}
	} else { //审核不通过的应用，商家修改材料，重新进入审核流程
		status := dao.BusinessAppDaoInst.GetBusinessAppStatus(req.AppId)
		if status != constants.BusinessAppStatus_Deny {
			ss_log.Error("应用不是审核不通过状态，不允许修改材料。")
			reply.ResultCode = ss_err.ERR_BUSINESSApp_Status_FAILD
		}

		if err := dao.BusinessAppDaoInst.UpdateBusinessApp(dao.BusinessAppDao{
			AppId:     req.AppId,
			Status:    constants.BusinessAppStatus_Pending, //将应用状态改为未审核
			ApplyType: req.ApplyType,                       //应用类型 1-移动应用，2-网页应用
			AppName:   req.AppName,                         //应用名称
			Describe:  req.Describe,                        //应用描述

			SmallImgNo:     req.SmallImgNo,     //小图标id
			BigImgNo:       req.BigImgNo,       //大图标id
			BusinessPubKey: req.BusinessPubKey, //商家公钥
			SignMethod:     req.SignMethod,     //签名方式
			IpwhiteList:    req.IpWhiteList,
		}); err != nil {
			ss_log.Error("修改应用失败，req=[%+v],err=[%v]", req, err)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}

	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
删除商家应用
*/
func (c *CustHandler) DelBusinessApp(ctx context.Context, req *custProto.DelBusinessAppRequest, reply *custProto.DelBusinessAppReply) error {
	if !dao.BusinessAppDaoInst.CheckBusinessApp(req.AppId, req.IdenNo) {
		ss_log.Error("应用编号[%v]不属于商家[%v],无权限删除他人的应用。", req.AppId, req.IdenNo)
		reply.ResultCode = ss_err.ERR_SYS_NO_API_AUTH
		return nil
	}

	status := dao.BusinessAppDaoInst.GetBusinessAppStatus(req.AppId)
	switch status {
	case constants.BusinessAppStatus_Passed:
	case constants.BusinessAppStatus_Deny:
	default:
		ss_log.Error("应用[%v]状态不是审核不通过和审核通过，不允许删除", req.AppId)
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}

	if err := dao.BusinessAppDaoInst.DelBusinessApp(req.AppId); err != nil {
		ss_log.Error("删除应用失败，req=[%+v],err=[%v]", req, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
web管理平台审核修改应用状态
*/
func (c *CustHandler) UpdateBusinessAppStatus(ctx context.Context, req *custProto.UpdateBusinessAppStatusRequest, reply *custProto.UpdateBusinessAppStatusReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}
	defer ss_sql.Rollback(tx)

	appQrCodeId := ""
	switch req.Status {
	case constants.BusinessAppStatus_Passed: //通过

		////生成固码
		appQrCodeId = constants.GetQrCodeId(req.AppId)

	case constants.BusinessAppStatus_Deny:
	case constants.BusinessAppStatus_Invalid:
	default:
		ss_log.Error("status参数[%v]不合法", req.Status)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if err := dao.BusinessAppDaoInst.UpdateBusinessAppStatusTx(tx, req.AppId, req.Status, req.Notes, appQrCodeId); err != nil {
		ss_log.Error("修改商家应用[%v]审核状态失败，err=[%v]", req.AppId, err)
		reply.ResultCode = ss_err.ERR_Audit_FAILD
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
商家上下架应用（修改状态）
*/
func (c *CustHandler) BusinessUpdateAppStatus(ctx context.Context, req *custProto.BusinessUpdateAppStatusRequest, reply *custProto.BusinessUpdateAppStatusReply) error {
	oldStatus := ""
	switch req.Status {
	case constants.BusinessAppStatus_Passed: //通过
		oldStatus = constants.BusinessAppStatus_Up
	case constants.BusinessAppStatus_Up: //上架
		oldStatus = constants.BusinessAppStatus_Passed
	default:
		ss_log.Error("status参数[%v]不合法", req.Status)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if err := dao.BusinessAppDaoInst.UpdateBusinessAppStatus(req.AppId, oldStatus, req.Status); err != nil {
		ss_log.Error("修改商家应用[%v]状态失败，err=[%v]", req.AppId, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
app管理相关接口
*/
func (*CustHandler) GetAppVersions(ctx context.Context, req *custProto.GetAppVersionsRequest, reply *custProto.GetAppVersionsReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ap.is_delete", Val: "0", EqType: "="},
		{Key: "ap.vs_type", Val: req.VsType, EqType: "="},
		{Key: "ap.system", Val: req.System, EqType: "="},
	})

	total := dao.AppVersionDaoInst.GetCnt(dbHandler, whereModel.WhereStr, whereModel.Args)

	//添加排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by ap.create_time desc, ap.vs_code desc `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.AppVersionDaoInst.GetAppVersions(dbHandler, whereModel.WhereStr, whereModel.Args)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetAppVersionsCount(ctx context.Context, req *custProto.GetAppVersionsCountRequest, reply *custProto.GetAppVersionsCountReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ap.is_delete", Val: "0", EqType: "="},
	})

	//moderpay
	//统计ios  app包的统计
	iosAppCountData := dao.AppVersionDaoInst.GetVersionCount(dbHandler, whereModel, constants.AppVersionSystem_Ios, constants.AppVersionVsType_app)

	//统计Android app包的统计
	androidAppCountData := dao.AppVersionDaoInst.GetVersionCount(dbHandler, whereModel, constants.AppVersionSystem_Android, constants.AppVersionVsType_app)
	//统计Android pos包的统计
	androidPosCountData := dao.AppVersionDaoInst.GetVersionCount(dbHandler, whereModel, constants.AppVersionSystem_Android, constants.AppVersionVsType_pos)

	//mangopay
	//统计ios  app包的统计
	iosMangopayAppCountData := dao.AppVersionDaoInst.GetVersionCount(dbHandler, whereModel, constants.AppVersionSystem_Ios, constants.APPVERSIONVSTYPE_MANGOPAY_APP)

	//统计Android app包的统计
	androidMangopayAppCountData := dao.AppVersionDaoInst.GetVersionCount(dbHandler, whereModel, constants.AppVersionSystem_Android, constants.APPVERSIONVSTYPE_MANGOPAY_APP)
	//统计Android pos包的统计
	androidMangopayPosCountData := dao.AppVersionDaoInst.GetVersionCount(dbHandler, whereModel, constants.AppVersionSystem_Android, constants.APPVERSIONVSTYPE_MANGOPAY_POS)

	var datas []*custProto.GetAppVersionsCountData
	datas = append(datas, iosAppCountData)
	datas = append(datas, androidAppCountData)
	datas = append(datas, androidPosCountData)

	datas = append(datas, iosMangopayAppCountData)
	datas = append(datas, androidMangopayAppCountData)
	datas = append(datas, androidMangopayPosCountData)

	reply.CountData = datas
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetAppVersion(ctx context.Context, req *custProto.GetAppVersionRequest, reply *custProto.GetAppVersionReply) error {

	data, err := dao.AppVersionDaoInst.GetAppVersionDetail(req.VId)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) ModifyAppVersionStatus(ctx context.Context, req *custProto.ModifyAppVersionStatusRequest, reply *custProto.ModifyAppVersionStatusReply) error {
	if req.Status != constants.Status_Disable && req.Status != constants.Status_Enable {
		ss_log.Error("Status状态错误")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.VId == "" {
		ss_log.Error("VId为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.LoginUid == "" {
		ss_log.Error("LoginUid为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if errStr := dao.AppVersionDaoInst.ModifyAppVersionStatus(req.VId, req.Status); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	//添加关键操作日志
	str1, legal1 := util.GetParamZhCn(req.Status, util.UseStatus)
	if !legal1 {
		ss_log.Error("Status %v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	description := fmt.Sprintf("设置版本vid[%v]的状态为[%v]", req.VId, str1)

	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) ModifyAppVersionIsForce(ctx context.Context, req *custProto.ModifyAppVersionIsForceRequest, reply *custProto.ModifyAppVersionIsForceReply) error {

	str1, legal1 := util.GetParamZhCn(req.IsForce, util.IsForce)
	if !legal1 {
		ss_log.Error("IsForce %v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.VId == "" {
		ss_log.Error("VId为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	if req.LoginUid == "" {
		ss_log.Error("LoginUid为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.IsForce == constants.AppVersionIsForce_True {
		//查询该版本信息，如果是最新版本的将不可设置为强制更新
		data, err := dao.AppVersionDaoInst.GetAppVersionDetail(req.VId)
		if err != ss_err.ERR_SUCCESS {
			ss_log.Error("查询版本信息出错 err=[%v]", err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		//最新包的版本version、最后修改时间
		newVersion, _, _ := dao.AppVersionDaoInst.GetNewVersion(data.System, data.VsType)
		if data.Version == newVersion {
			ss_log.Error("不可将最新版本[%v]设置为强制更新", newVersion)
			reply.ResultCode = ss_err.ERR_VERSION_IsForce_Faile
			return nil
		}
	}

	if errStr := dao.AppVersionDaoInst.ModifyIsForce(req.VId, req.IsForce); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	description := fmt.Sprintf("将vid[%v]版本的是否强制更新设置为[%v]", req.VId, str1)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) InsertOrUpdateAppVersion(ctx context.Context, req *custProto.InsertOrUpdateAppVersionRequest, reply *custProto.InsertOrUpdateAppVersionReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	str1, legal1 := util.GetParamZhCn(req.Status, util.UseStatus)
	if !legal1 {
		ss_log.Error("Status %v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.VId != "" { //更新只能更新部分信息
		if errStr := dao.AppVersionDaoInst.ModifyAppVersionInfo(req.VId, req.Description, req.Note, req.Status); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}

		description := fmt.Sprintf("更新版本id[%v]的版本信息描述为[%v],备注为[%v],使用状态为[%v]", req.VId, req.Description, req.Note, str1)
		errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.AccountNo, constants.LogAccountWebType_Config)
		if errAddLog != nil {
			ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}

		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	}

	//校验参数
	systemStr, legal1 := util.GetParamZhCn(req.System, util.System)
	if !legal1 {
		ss_log.Error("System %v", systemStr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	VsTypeStr, legal2 := util.GetParamZhCn(req.VsType, util.VsType)
	if !legal2 {
		ss_log.Error("VsType %v", VsTypeStr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	versionStr := constants.AppVersion_Init
	vsCodeStr := "1"
	consecutiveLitersNumber := req.ConsecutiveLitersNumber
	if consecutiveLitersNumber == "" {
		consecutiveLitersNumber = "1"
	}
	//查看是否有同类的包，有则要查询出最新的版本，再根据升级的版本更新获取插入的版本号
	if dao.AppVersionDaoInst.CheckHaveVersion(req.System, req.VsType) {
		//查询最新的版本号
		versionGet, _, vsCode := dao.AppVersionDaoInst.GetNewVersion(req.System, req.VsType)
		if versionGet == "" {
			ss_log.Error("获取最新版本出错")
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		vsCodeStr = ss_count.Add(vsCode, "1")

		//根据升级的版本更新获取插入的版本号
		numStrArr := strings.Split(versionGet, ".")
		switch req.UpType {
		case constants.AppVersionUpType_Big:
			//大版本更新，版本号第一位数字+1，第二三位数字归零；
			numStrArr[0] = ss_count.Add(numStrArr[0], consecutiveLitersNumber)
			numStrArr[1] = "0"
			numStrArr[2] = "0"
			versionStr = numStrArr[0] + "." + numStrArr[1] + "." + numStrArr[2]
		case constants.AppVersionUpType_Small:
			//小版本更新，版本号第一位数字不变，第二位数字+1，第三位数字归零；
			numStrArr[1] = ss_count.Add(numStrArr[1], consecutiveLitersNumber)
			numStrArr[2] = "0"
			versionStr = numStrArr[0] + "." + numStrArr[1] + "." + numStrArr[2]
		case constants.AppVersionUpType_Bug:
			//BUG版本更新，版本号第一二位数字不变，第三位数+1；
			numStrArr[2] = ss_count.Add(numStrArr[2], consecutiveLitersNumber)
			versionStr = numStrArr[0] + "." + numStrArr[1] + "." + numStrArr[2]
		default:
			ss_log.Error("UpType类型错误")
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	//_, pathStr, errP := cache.ApiDaoInstance.GetGlobalParam("app_version_file_path")
	//if errP != nil {
	//	ss_log.Error("获取保存路径失败，errP=[%v]\n", errP)
	//	reply.ResultCode = ss_err.ERR_PARAM
	//	return nil
	//}

	errCheck, fileName := dao.AppVersionFileLogDaoInstance.CheckAppNo(req.FileId)
	if errCheck != ss_err.ERR_SUCCESS {
		ss_log.Error("获取文件名失败")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//改名
	// khmangopay-pos-1.0.1.apk //pos包的名称
	//app
	//khmangopay-android-1.0.1.apk //安卓包的名称
	//khmangopay-iso-1.0.1.apk  //ios包的名称
	ss_log.Error("fileName:%v", fileName)
	newReName := fmt.Sprintf("%s/%s-%s-%s%s", filepath.Dir(fileName), systemStr, VsTypeStr, versionStr, filepath.Ext(fileName))
	newReName2 := fmt.Sprintf("%s/%s-%s-%s%s", filepath.Dir(fileName), systemStr, VsTypeStr, "latest", filepath.Ext(fileName))

	ss_log.Info("newReName:[%v]", newReName)

	errUpName := dao.AppVersionFileLogDaoInstance.ModifyFileName(tx, req.FileId, newReName)
	if errUpName != ss_err.ERR_SUCCESS {
		ss_log.Error("修改上传日志内的文件名失败.")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if errStr := dao.AppVersionDaoInst.InsertAppVersion(tx, req.Description, versionStr, newReName, vsCodeStr, req.VsType, req.System, req.AccountNo, req.Note, req.Status); errStr != ss_err.ERR_SUCCESS {
		ss_log.Error("插入版本日志失败")
		reply.ResultCode = errStr
		return nil
	}

	//复制一个新的，实现改名（目前未发现有直接改名的方法，所以先复制再删除原来的）
	_, errAwsS3Copy := common.UploadS3.CopyObject(fileName, newReName, true)
	if errAwsS3Copy != nil {
		ss_log.Error("AwsS3复制失败，errAwsS3Copy:[%v]", errAwsS3Copy)
		reply.ResultCode = ss_err.ERR_UPLOAD
		return nil
	}

	//复制一个新的,为最新的
	_, errAwsS3Copy2 := common.UploadS3.CopyObject(fileName, newReName2, true)
	if errAwsS3Copy2 != nil {
		ss_log.Error("AwsS3复制失败，errAwsS3Copy2:[%v]", errAwsS3Copy2)
		reply.ResultCode = ss_err.ERR_UPLOAD
		return nil
	}

	//删除文件
	_, errAwsS3Del := common.UploadS3.DeleteOne(fileName)
	if errAwsS3Del != nil {
		ss_log.Error("AwsS3删除失败，errAwsS3Del:[%v]", errAwsS3Del)
		reply.ResultCode = ss_err.ERR_UPLOAD
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetNewVersion(ctx context.Context, req *custProto.GetNewVersionRequest, reply *custProto.GetNewVersionReply) error {
	version, updateTime, _ := dao.AppVersionDaoInst.GetNewVersion(req.System, req.VsType)
	reply.NewVersion = version
	reply.UpdateTime = updateTime
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
