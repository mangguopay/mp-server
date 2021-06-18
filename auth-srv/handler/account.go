package handler

import (
	"fmt"
	"strconv"

	"a.a/cu/db"
	"a.a/cu/encrypt"
	"a.a/cu/jwt"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	util2 "a.a/cu/util"
	"a.a/mp-server/auth-srv/common"
	"a.a/mp-server/auth-srv/dao"
	"a.a/mp-server/auth-srv/i"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_data"
	"a.a/mp-server/common/ss_func"
	"a.a/mp-server/common/ss_struct"

	//"a.a/mp-server/auth-srv/i"
	"a.a/mp-server/auth-srv/util"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_auth "a.a/mp-server/common/proto/auth"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"

	//"a.a/mp-server/common/proto/auth"
	//"a.a/mp-server/common/proto/cust"
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

/**
 * 获取账户
 */
func (*Auth) GetAccount(ctx context.Context, req *go_micro_srv_auth.GetAccountRequest, resp *go_micro_srv_auth.GetAccountReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	row, stmt, err := ss_sql.QueryRowN(dbHandler, "SELECT ac.country_code,ac.gen_key, ac.uid, ac.nickname, ac.account, ac.use_status, ac.create_time, "+
		"ac.modify_time, ac.drop_time, ac.master_acc,acc2.account,ac.phone,ac.email,ac.head_portrait_img_no, ac.business_phone "+
		"FROM account ac left join account acc2 on acc2.uid=ac.master_acc "+
		"WHERE ac.uid = $1 LIMIT 1", req.AccountUid)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
		return nil
	}

	var masterAcc, masterAccount, headPortraitImgNo, email, businessPhone sql.NullString
	err = row.Scan(
		&resp.CountryCode,
		&resp.GenKey,
		&resp.Uid,
		&resp.Nickname,
		&resp.Account,
		&resp.UseStatus,
		&resp.CreateTime,
		&resp.ModifyTime,
		&resp.DropTime,
		&masterAcc,
		&masterAccount,
		&resp.Phone,
		&email,
		&headPortraitImgNo,
		&businessPhone,
	)
	resp.HeadPortraitImgNo = headPortraitImgNo.String
	resp.BusinessPhone = businessPhone.String
	resp.Email = email.String
	resp.MasterAccount = masterAccount.String
	resp.MasterAcc = masterAcc.String
	if nil != err {
		ss_log.Error("err=[%v]", err)
		resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
		return nil
	}

	rowsRole, stmt, err2 := ss_sql.Query(dbHandler, "SELECT role_uid FROM rela_account_role WHERE account_uid = $1", req.AccountUid)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rowsRole.Close()
	roleUids := db.ToStringList(rowsRole)
	ss_log.Info("账号[%v]关联的角色ids[%v]", req.AccountUid, roleUids)
	if nil != err2 {
		ss_log.Error("err|2=%v", err2)
		resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
		return nil
	}

	// 读取菜单列表
	if len(roleUids) > 0 {
		rowsUrl, stmt, errSel := ss_sql.Query(dbHandler, "SELECT url_uid,url_name,url,parent_uid,title,icon,component_name,component_path,redirect,idx,is_hidden FROM url "+
			"WHERE url_uid in (SELECT url_uid FROM rela_role_url WHERE role_uid in ('"+strings.Join(roleUids, "','")+"'))")
		if stmt != nil {
			defer stmt.Close()
		}
		defer rowsUrl.Close()
		if errSel == nil {
			for rowsUrl.Next() {
				data := go_micro_srv_auth.RouteData{}
				err := rowsUrl.Scan(&data.UrlUid, &data.UrlName, &data.Url, &data.ParentUid, &data.Title, &data.Icon, &data.ComponentName,
					&data.ComponentPath, &data.Redirect, &data.Idx, &data.IsHidden)
				if err != nil {
					ss_log.Error("err=[%v]", err)
					continue
				}
				resp.DataList = append(resp.DataList, &data)
			}
		}

		if nil != errSel {
			ss_log.Error("err|3=%v", errSel)
			resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
			return nil
		}
	}

	idenNo := dao.AccDaoInstance.GetIdenNoFromAcc(req.AccountUid, req.AccountType)

	//if idenNo == "" {
	//	ss_log.Info("获取账户角色信息失败,uid为-------->%s", req.AccountUid)
	//	resp.ResultCode = ss_err.ERR_PARAM
	//	return nil
	//}

	// 修改account中的app_lang
	if errStr := dao.AccDaoInstance.UpdateLang(req.AccountUid, req.Lang); errStr != ss_err.ERR_SUCCESS {
		ss_log.Error("%s", "修改用户语言失败")
		resp.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	ss_log.Info("获取账户[%v]角色信息成功", req.AccountUid)
	resp.ResultCode = ss_err.ERR_SUCCESS
	resp.MerchantUid = idenNo
	resp.AccountType = req.AccountType
	return nil
}

func (r *Auth) GetVersinInfo(ctx context.Context, req *go_micro_srv_auth.GetVersinInfoRequest, reply *go_micro_srv_auth.GetVersinInfoReply) error {
	var appVersion, versionDescription, appUrl, isForce string

	switch req.VsType {
	case constants.AppVersionVsType_app:
		fallthrough
	case constants.APPVERSIONVSTYPE_MANGOPAY_APP:
		if req.System == "" {
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		appVersion, versionDescription, appUrl, isForce = getAppVersion(req.VsType, req.System, req.AppVersion)
	case constants.AppVersionVsType_pos:
		fallthrough
	case constants.APPVERSIONVSTYPE_MANGOPAY_POS:
		appVersion, versionDescription, appUrl, isForce = getVersion(req.AppVersion, req.VsType)
	}

	appBaseUrl := dao.GlobalParamDaoInstance.QeuryParamValue("app_base_url")

	reply.AppVersion = appVersion
	reply.VersionDescription = versionDescription
	reply.AppUrl = appBaseUrl + "/" + appUrl
	reply.IsForce = isForce
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 上传客户端信息
func (r *Auth) UploadClientInfo(ctx context.Context, req *go_micro_srv_auth.UploadClientInfoRequest, reply *go_micro_srv_auth.UploadClientInfoReply) error {

	switch req.ClientType {
	case 1: // app端
		instance := dao.ClientInfoAppDao{}
		instance.DeviceBrand = req.DeviceBrand
		instance.DeviceModel = req.DeviceModel
		instance.Resolution = req.Resolution
		instance.ScreenSize = req.ScreenSize
		instance.Imei1 = req.Imei1
		instance.Imei2 = req.Imei2
		instance.SystemVer = req.SystemVer
		instance.UploadPoint = req.UploadPoint
		instance.UserAgent = req.UserAgent
		instance.Platform = req.Platform
		instance.AppVer = req.AppVer
		instance.Account = req.Account
		instance.Uuid = req.Uuid

		// 插入数据
		if err := instance.Insert(); err != nil {
			ss_log.Error("UploadClientInfo-Cust-Insert,err=[%v],data:%+v", err, instance)
		}
	case 2: // pos端
		instance := dao.ClientInfoPosDao{}
		instance.DeviceBrand = req.DeviceBrand
		instance.DeviceModel = req.DeviceModel
		instance.Resolution = req.Resolution
		instance.ScreenSize = req.ScreenSize
		instance.Imei1 = req.Imei1
		instance.Imei2 = req.Imei2
		instance.SystemVer = req.SystemVer
		instance.UploadPoint = req.UploadPoint
		instance.UserAgent = req.UserAgent
		instance.Platform = req.Platform
		instance.AppVer = req.AppVer
		instance.Account = req.Account
		instance.Uuid = req.Uuid

		// 插入数据
		if err := instance.Insert(); err != nil {
			ss_log.Error("UploadClientInfo-Servicer-Insert,err=[%v],data:%+v", err, instance)
		}
	default:
		ss_log.Error("UploadClientInfo-客户端类型错误,req:%+v", req)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func getVersion(version, vsType string) (appVersion, versionDescription, appUrl, isForce string) {

	// ==============暂时不删除=================
	// 查找最新的版本号
	//latDescription, latVersion, latAppURL, latVsCode, isFor := dao.AppVersionDaoInstance.QueryAppVersion("", vsType)
	// 历史版本信息
	//_, _, _, hisVsCode, _ := dao.AppVersionDaoInstance.QueryAppVersion(version, vsType)

	//if strext.ToFloat64(hisVsCode) < strext.ToFloat64(latVsCode) {
	//	appVersion = latVersion
	//	versionDescription = latDescription
	//	appUrl = latAppURL
	//	isForce = isFor
	//	return
	//}
	//=================================

	latDescription, latVersion, latAppURL, _, isFor := dao.AppVersionDaoInstance.QueryAppVersion("", vsType)
	if latVersion == "" {
		appVersion = version
	} else {
		appVersion = latVersion
	}
	versionDescription = latDescription
	appUrl = latAppURL
	isForce = isFor

	return
}
func getAppVersion(vsType, system, oldAppVersion string) (appVersion, versionDescription, appUrl, isForce string) {
	latDescription, latVersion, latAppURL, _, isFor := dao.AppVersionDaoInstance.QueryAppVersion1(vsType, system)
	if latVersion == "" {
		appVersion = oldAppVersion
	} else {
		appVersion = latVersion
	}
	versionDescription = latDescription
	appUrl = latAppURL
	isForce = isFor
	return
}

//修改用户使用状态（冻结）
func (*Auth) ModifyUserStatus(ctx context.Context, req *go_micro_srv_auth.ModifyUserStatusRequest, resp *go_micro_srv_auth.ModifyUserStatusReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	account := dao.AccDaoInstance.GetAccountByUid(req.Uid)
	description := ""
	switch req.SetStatus { //0.禁用，1.正常
	case constants.Status_Enable: //如果是解冻则还要将连续密码错误的次数变为0
		key := cache.GetPwdErrCountKey(cache.PrePwdErrCountKey, account)
		cache.RedisClient.Set(key, 0, 0)

		description = fmt.Sprintf("解冻账号[%v]", account)
	case constants.Status_Disable:
		description = fmt.Sprintf("冻结账号[%v]", account)
	default:
		ss_log.Error("参数SetStatus[%v]错误", req.SetStatus)
		resp.ResultCode = ss_err.ERR_ACCOUNT_STATUS
		return nil
	}

	sqlUpdate := "Update account set use_status = $1 where uid = $2 and is_delete ='0' "
	err := ss_sql.Exec(dbHandler, sqlUpdate, req.SetStatus, req.Uid)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		resp.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Account)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败,err=[%v],", description, errAddLog)
	}
	resp.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (r *Auth) BusinessLogin(ctx context.Context, req *go_micro_srv_auth.BusinessLoginRequest, resp *go_micro_srv_auth.BusinessLoginReply) error {
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

	var accUid, pwdMD5 sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT uid,password FROM account WHERE account=$1 LIMIT 1", []*sql.NullString{&accUid, &pwdMD5}, req.Account)
	if err != nil {
		ss_log.Error("login failed|acc=[%v]|password=[%v]\n", req.Account, pwdMD5.String)
		resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
		return nil
	}

	//md5(md5(sha1(原)+salt)+nonstr)
	pwdMD5Fixed := req.Password
	pwdMD5FixedDB := encrypt.DoMd5Salted(pwdMD5.String, req.Nonstr)
	ss_log.Info("pwdMD5=[%v], pwdMD5Fixed=[%v],pwdMD5FixedDB=[%v]", pwdMD5.String, pwdMD5Fixed, pwdMD5FixedDB)
	if req.Account == "" || pwdMD5Fixed != pwdMD5FixedDB {
		ss_log.Error("login failed|acc=[%v]|password=[%v]\n", req.Account, pwdMD5.String)
		resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
		return nil
	}

	if accUid.Valid == false {
		ss_log.Error("login failed|acc=[%v]|password=[%v]\n", req.Account, pwdMD5)
		err = ss_sql.Exec(dbHandler, "insert into log_login(log_time,acc_no,ip,result,client,log_no) values (current_timestamp,$1,$2,$3,$4,$5)",
			accUid.String, "", ss_err.ERR_ACCOUNT_NOT_EXISTS, req.Client, strext.GetDailyId())
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
		resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
		return nil
	}

	getAcc := go_micro_srv_auth.GetAccountReply{}
	_ = r.GetAccount(ctx, &go_micro_srv_auth.GetAccountRequest{
		AccountUid: accUid.String,
	}, &getAcc)
	if getAcc.UseStatus != "1" {
		err = ss_sql.Exec(dbHandler, "insert into log_login(log_time,acc_no,ip,result,client,log_no) values (current_timestamp,$1,$2,$3,$4,$5)",
			accUid.String, "", ss_err.ERR_ACCOUNT_NO_PERMISSION, req.Client, strext.GetDailyId())
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
		resp.ResultCode = ss_err.ERR_ACCOUNT_NO_PERMISSION
		return nil
	}

	accountType := dao.AccDaoInstance.GetAccountTypeFromAccNoBusiness(accUid.String)
	if accountType == "" {
		err = ss_sql.Exec(dbHandler, "insert into log_login(log_time,acc_no,ip,result,client,log_no) values (current_timestamp,$1,$2,$3,$4,$5)",
			accUid.String, "", ss_err.ERR_ACCOUNT_NO_PERMISSION, req.Client, strext.GetDailyId())
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
		resp.ResultCode = ss_err.ERR_ACCOUNT_NO_PERMISSION
		return nil
	}

	idenNo := dao.AccDaoInstance.GetIdenNoFromAcc(accUid.String, accountType)

	//根据身份id查询是否已初始化支付密码(0否，1是)
	initPayPwdStatus, err := dao.BusinessDaoInstance.GetBusinessInitPayPwdStatus(idenNo)
	if err != nil {
		ss_log.Error("查询账号uid[%v]的商家[%v]支付密码是否初始化状态失败,err=[%v]", accUid.String, idenNo, err)
		resp.ResultCode = ss_err.ERR_ACCOUNT_NO_PERMISSION
		return nil
	}
	ss_log.Info("initPayPwdStatus:[%v]", initPayPwdStatus)
	resp.InitPayPwdStatus = initPayPwdStatus

	phone := ""
	switch accountType {
	case constants.AccountType_PersonalBusiness: //个人商家的手机号是和app用户一样用的phone
		phone = getAcc.Phone
	case constants.AccountType_EnterpriseBusiness: //企业商家的手机号是新创建的字段business_phone
		phone = getAcc.BusinessPhone
	}
	jwt2 := common.CreateWebBusinessJWT(ss_struct.JwtDataWebBusiness{
		Account:        req.Account,
		AccountUid:     accUid.String,
		IdenNo:         idenNo, //BusinessNo
		AccountType:    accountType,
		LoginAccountNo: accUid.String,
		Email:          getAcc.Email,
		Phone:          phone,
		CountryCode:    getAcc.CountryCode,
		//JumpIdenNo:     "",
		//JumpIdenType:   "",
		//MasterAccNo:    "",
		IsMasterAcc: "1",
	})

	err = ss_sql.Exec(dbHandler, "insert into log_login(log_time,acc_no,ip,result,client,log_no) values (current_timestamp,$1,$2,$3,$4,$5)",
		accUid.String, req.Ip, ss_err.ERR_SUCCESS, req.Client, strext.GetDailyId())
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	resp.ResultCode = ss_err.ERR_SUCCESS
	resp.AccountUid = accUid.String
	resp.Jwt = jwt2
	resp.IsEncrypted = true
	resp.AccountType = accountType
	return nil
}

func (r *Auth) InitBusinessPayPwd(ctx context.Context, req *go_micro_srv_auth.InitBusinessPayPwdRequest, reply *go_micro_srv_auth.InitBusinessPayPwdReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	//确认该商家是第一次登录来设置支付密码
	initPayPwdStatus, err := dao.BusinessDaoInstance.GetBusinessInitPayPwdStatus(req.IdenNo)
	if err != nil {
		ss_log.Error("查询账号uid[%v]的商家[%v]支付密码是否初始化状态失败,err=[%v]", req.Uid, req.IdenNo, err)
		reply.ResultCode = ss_err.ERR_ACCOUNT_NO_PERMISSION
		return nil
	}
	if initPayPwdStatus != constants.InitPayPwdStatus_false {
		ss_log.Error("商家[%v]支付密码已初始化过，不允许再次初始化。", req.IdenNo)
		reply.ResultCode = ss_err.ERR_HasPayPwd_FAILD
		return nil
	}

	//修改支付密码
	if err := dao.BusinessDaoInstance.ModifyBusinessPayPwdByUidTx(tx, req.IdenNo, req.PayPwd); err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_MODIFY_PAY_PWD_FAILD
		return nil
	}

	//将未初始化支付密码状态改变
	if err := dao.BusinessDaoInstance.ModifyBusinessInitPayPwdStatus(tx, req.IdenNo); err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_MODIFY_PAY_PWD_FAILD
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//注册并登陆
func (r *Auth) RegLogin(ctx context.Context, req *go_micro_srv_auth.RegLoginRequest, reply *go_micro_srv_auth.RegLoginReply) error {
	// 校验短信验证码开关
	k1, isCheck, err := cache.ApiDaoInstance.GetGlobalParam("is_check_sms") // is_check_sms 0-需要校验,1-不需要校验
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k1)
	}

	// 校验短信验证码是否正确
	if isCheck == "0" {
		if !cache.CheckSMS(constants.FUNCTIONREG, req.Phone, req.Sms) {
			//todo 插入注册失败的日志
			err := dao.LogDaoInstance.InsertLogAppRegister("注册失败:短信验证码错误", req.Uuid, "2")
			if err != nil {
				ss_log.Error("插入注册失败的日志出错，err[%v],uuid=[%v]", err, req.Uuid)
			}
			reply.ResultCode = ss_err.ERR_ACCOUNT_SMS_MSG_FAILD
			return nil
		}
	}
	//是否需要添加用户账号
	needAddCust := true
	accountUid := ""

	accUid, isActived := dao.AccDaoInstance.GetAccountIsActived(ss_func.PreCountryCode(req.CountryCode), req.CountryCode, req.Phone)
	if accUid != "" { //账号存在
		needAddCust = false
		accountUid = accUid
		if isActived == "1" { //并且是激活的，直接返回错误账号已存在
			//todo 插入注册失败的日志
			err := dao.LogDaoInstance.InsertLogAppRegister("注册失败:账户已存在", req.Uuid, "2")
			if err != nil {
				ss_log.Error("插入注册失败的日志出错，err[%v],uuid=[%v]", err, req.Uuid)
			}
			reply.ResultCode = ss_err.ERR_ACCOUNT_ALREADY_EXISTS
			return nil
		} else {
			//则创建cust并关联账号
			dbHandler := db.GetDB(constants.DB_CRM)
			defer db.PutDB(constants.DB_CRM, dbHandler)

			tx, errTx := dbHandler.BeginTx(ctx, nil)
			if errTx != nil {
				ss_log.Error("开启事务失败,errTx=[%v]", errTx)
				reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
				return nil
			}
			defer ss_sql.Rollback(tx)

			//激活账号
			if err := dao.AccDaoInstance.UpdateAccountIsActived(tx, accUid); err != nil {
				ss_log.Error("修改激活状态失败accUid[%v],err=[%v]", accUid, err)
				reply.ResultCode = ss_err.ERR_ACCOUNT_ALREADY_EXISTS
				return nil
			}

			//修改密码和昵称
			nickName := util2.RandomStringLower(6)
			if err := dao.AccDaoInstance.UpdateAccountPwdTx(tx, req.Phone, req.Password, req.CountryCode, nickName); err != nil {
				ss_log.Error("修改密码失败Phone[%v],CountryCode[%v],err=[%v]", req.Phone, req.CountryCode, err)
				reply.ResultCode = ss_err.ERR_ACCOUNT_ALREADY_EXISTS
				return nil
			}

			//添加用户
			err, custNo := dao.CustDaoInstance.AddCustTx(tx, accUid, "1")
			if err != nil {
				ss_log.Error("创建用户失败,err=[%v]", err)
				reply.ResultCode = ss_err.ERR_ACCOUNT_ALREADY_EXISTS
				return nil
			}

			//添加关系
			errCode := dao.RelaAccIdenDaoInst.InsertRelaAccIdenTx(tx, accUid, custNo, constants.AccountType_USER)
			if errCode != ss_err.ERR_SUCCESS {
				ss_log.Error("添加关联关系失败，accUid[%v],custNo[%v],errCode=[%v]", accUid, custNo, errCode)
				reply.ResultCode = ss_err.ERR_ACCOUNT_ALREADY_EXISTS
				return nil
			}

			//同步冻结虚账金额
			err = VAccountHandlerInst.SyncFreezeVAccountBalance(tx, accUid)
			if err != nil {
				ss_log.Error("同步冻结虚账金额失败，accountNo=%v, err=%v", accUid, err)
				reply.ResultCode = ss_err.ERR_ACCOUNT_ALREADY_EXISTS
				return nil
			}

			tx.Commit()
		}
	}

	if needAddCust {
		//以下是不存在账号的处理
		addCustReq := &go_micro_srv_auth.AddCustRequest{
			Gender:      "1",
			Nickname:    util2.RandomStringLower(6),
			Phone:       req.Phone,
			Password:    req.Password,
			CountryCode: req.CountryCode,
		}

		addCustReply := &go_micro_srv_auth.AddCustReply{}
		errAdd := AuthHandlerInst.AddCust(context.TODO(), addCustReq, addCustReply)
		if errAdd != nil {
			ss_log.Error("phone[%v]添加用户调用api失败，err=[%v]", req.Phone, errAdd)
			//todo 插入注册失败的日志
			if err := dao.LogDaoInstance.InsertLogAppRegister("注册失败:调用api失败", req.Uuid, "2"); err != nil {
				ss_log.Error("插入注册失败的日志出错，err[%v],uuid=[%v]", err, req.Uuid)
			}
			reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
			return nil
		}

		if addCustReply.ResultCode != ss_err.ERR_SUCCESS {
			//todo 插入注册失败的日志
			if err := dao.LogDaoInstance.InsertLogAppRegister("注册失败:调用api失败", req.Uuid, "2"); err != nil {
				ss_log.Error("插入注册失败的日志出错，err[%v],uuid=[%v]", err, req.Uuid)
			}
			ss_log.Error("phone[%v]添加用户失败,err=[%v]", req.Phone, addCustReply.ResultCode)
			reply.ResultCode = addCustReply.ResultCode
			return nil
		}

		accountUid = addCustReply.AccUid

		//初始化钱包1、2
		if _, err := dao.VaccountDaoInst.InitVAccountNo(accountUid, constants.VaType_USD_DEBIT, constants.CURRENCY_USD); err != nil {
			ss_log.Error("初始化个人虚账失败，accountNo=%v, vaccType=%v, err=%v", accountUid, constants.VaType_USD_DEBIT, err)
			//return err
		}
		if _, err := dao.VaccountDaoInst.InitVAccountNo(accountUid, constants.VaType_KHR_DEBIT, constants.CURRENCY_KHR); err != nil {
			ss_log.Error("初始化个人虚账失败，accountNo=%v, vaccType=%v, err=%v", accountUid, constants.VaType_KHR_DEBIT, err)
			//return err
		}

	}

	//todo 插入注册成功的日志
	if err := dao.LogDaoInstance.InsertLogAppRegister("注册成功", req.Uuid, "1"); err != nil {
		ss_log.Error("插入注册成功的日志出错，err[%v],uuid=[%v]", err, req.Uuid)
	}

	idenType := constants.AccountType_USER
	idenNo := dao.RelaAccIdenDaoInst.GetIdenFromAcc(accountUid, idenType)

	getAcc := go_micro_srv_auth.GetAccountReply{}
	_ = r.GetAccount(ctx, &go_micro_srv_auth.GetAccountRequest{
		AccountUid: accountUid,
	}, &getAcc)
	if getAcc.UseStatus != "1" {
		dao.LogLoginDaoInstance.InsertLogLogin(accountUid, "", ss_err.ERR_ACCOUNT_NO_PERMISSION, req.Lat, req.Lng)
		reply.ResultCode = ss_err.ERR_ACCOUNT_NO_PERMISSION
		return nil
	}

	// todo 返回登录信息
	retMap := common.JwtStructToMapApp(ss_struct.JwtDataApp{
		Account:        req.Phone,
		AccountUid:     accountUid,
		IdenNo:         idenNo,
		AccountType:    idenType,
		LoginAccountNo: accountUid,
		PubKey:         req.PubKey,
		JumpIdenNo:     "",
		JumpIdenType:   "",
		MasterAccNo:    "",
		IsMasterAcc:    "1",
	})

	jwt2 := jwt.GetNewEncryptedJWTToken(constants.JwtLifeTime, retMap, constants.MobileKeyJwtAes, constants.MobileKeyJwtSign)
	genKey := dao.AccDaoInstance.QueryGenKeyFromAccNo(accountUid)

	// 查找最新的版本号
	latDescription, latVersion, latAppURL, latVsCode, _ := dao.AppVersionDaoInstance.QueryAppVersion("", constants.APPVERSIONVSTYPE_MANGOPAY_APP)
	// 历史版本信息
	_, _, _, hisVsCode, _ := dao.AppVersionDaoInstance.QueryAppVersion(req.AppVersion, constants.APPVERSIONVSTYPE_MANGOPAY_APP)
	var appVersion, versionDescription, appUrl string
	if strext.ToFloat64(hisVsCode) < strext.ToFloat64(latVsCode) {
		appVersion = latVersion
		versionDescription = latDescription
		appUrl = latAppURL
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.AccountUid = accountUid
	reply.Jwt = jwt2
	reply.IsEncrypted = true
	reply.GenKey = genKey
	reply.Phone = req.Phone
	reply.AccountType = idenType
	reply.AppVersion = appVersion
	reply.VersionDescription = versionDescription
	reply.AppUrl = appUrl
	reply.RefreshTokenInterval = constants.RefreshTokenInterval
	reply.IsSetPaymentPwd = false
	return nil
}

// 刷新token
func (r *Auth) RefreshToken(ctx context.Context, req *go_micro_srv_auth.RefreshTokenRequest, reply *go_micro_srv_auth.RefreshTokenReply) error {
	if req.Account == "" || req.JwtIat == "" {
		reply.ResultCode = ss_err.ERR_PARAM
		ss_log.Error("缺少参数,Account:%s, JwtIat:%s,", req.Account, req.JwtIat)
		return nil
	}

	// 将jwt的签发时间转成int64
	jwtIat, err := strconv.ParseInt(req.JwtIat, 10, 64)
	if err != nil {
		reply.ResultCode = ss_err.ERR_PARAM
		ss_log.Error("将jwt的签发时间转成int64失败,err:%v, JwtIat:%s", err, req.JwtIat)
		return nil
	}

	nowTimestamp := ss_time.NowTimestamp(global.Tz)
	if nowTimestamp-jwtIat < int64(constants.RefreshTokenInterval)/2 { // 除以2是防止刚好和前端同步上，导致刷新失败
		ss_log.Error("刷新token太过频繁,nowTimestamp:%d, jwtIat:%d, RefreshTokenInterval:%d", nowTimestamp, jwtIat, constants.RefreshTokenInterval)
		reply.ResultCode = ss_err.ERR_REFRESH_TOKEN_FREQUENTLY
		return nil
	}

	jwtMap := common.JwtStructToMapApp(ss_struct.JwtDataApp{
		Account:        req.Account,
		AccountUid:     req.AccountUid,
		IdenNo:         req.IdenNo,
		AccountType:    req.AccountType,
		LoginAccountNo: req.LoginAccountNo,
		PubKey:         req.PubKey,
		JumpIdenNo:     req.JumpIdenNo,
		JumpIdenType:   req.JumpIdenType,
		MasterAccNo:    req.MasterAccNo,
		IsMasterAcc:    req.IsMasterAcc,
		PosSn:          req.PosSn,
	})

	// 重新生成jwt
	jwt := jwt.GetNewEncryptedJWTToken(constants.JwtLifeTime, jwtMap, constants.MobileKeyJwtAes, constants.MobileKeyJwtSign)

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Jwt = jwt
	reply.RefreshTokenInterval = constants.RefreshTokenInterval
	return nil
}

func (r *Auth) MobileLogin(ctx context.Context, req *go_micro_srv_auth.MobileLoginRequest, resp *go_micro_srv_auth.MobileLoginReply) error {
	if req.Account == "" {
		ss_log.Error("缺少登录参数,req.Account为空")
		resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
		return nil
	}

	switch req.AccountType {
	case constants.AccountType_USER:
	case constants.AccountType_SERVICER:
	case constants.AccountType_POS:
	default:
		ss_log.Error("登录的账号类型不正确,登录的账号类型 account_type 为---> %s", req.AccountType)
		resp.ResultCode = ss_err.ERR_ACCOUNT_NO_PERMISSION
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	account := ss_func.ComposeAccountByPhoneCountryCode(req.Account, req.CountryCode)
	var err error
	var accUid, pwdMD5, genKey, phone, headPortraitImgNoT, isFirstLoginT sql.NullString
	err = ss_sql.QueryRow(dbHandler, "SELECT uid,password,gen_key,phone,head_portrait_img_no,is_first_login FROM account WHERE account=$1 and is_delete = '0' and is_actived='1' LIMIT 1",
		[]*sql.NullString{&accUid, &pwdMD5, &genKey, &phone, &headPortraitImgNoT, &isFirstLoginT}, account)
	if err != nil {
		resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
		return nil
	}

	// 验证密码
	//md5(md5(sha1(原)+salt)+nonstr)
	pwdMD5Fixed := req.Password
	pwdMD5FixedDB := encrypt.DoMd5Salted(pwdMD5.String, req.Nonstr)
	ss_log.Info("pwdMD5=[%v], pwdMD5Fixed=[%v],pwdMD5FixedDB=[%v],nonstr=[%v]", pwdMD5.String, pwdMD5Fixed, pwdMD5FixedDB, req.Nonstr)
	if req.Account == "" || pwdMD5Fixed != pwdMD5FixedDB {
		ss_log.Error("login failed|acc=[%v]|password=[%v]\n", req.Account, pwdMD5.String)

		if accUid.String != "" {
			//todo 增加错误密码日志（关键操作日志）
			description := "验证登录密码错误"
			if err := dao.LogDaoInstance.InsertAccountLog(description, accUid.String, req.AccountType, constants.MODIFYACCOUNTLOGTYPE); err != nil {
				ss_log.Error("err=[%v],missing key=[%v]", err, "插入验证登录密码错误的日志失败")
			}
		}

		// 处理密码试错次数
		if !pwdErrLimit(account) {
			// 修改改账号的使用状态
			if errStr := dao.AccDaoInstance.UpdateAccountStatusByAccount(account, constants.AccountUseStatusTemporaryDisabled); errStr != ss_err.ERR_SUCCESS {
				ss_log.Error("禁用状态失败")
			}
			resp.ResultCode = ss_err.ERR_ACCOUNT_ERR_PWD_LIMIT
			return nil
		}
		resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
		return nil
	}
	// 密码正确,判断是否被限制
	limitKey := cache.GetPwdErrCountKey(cache.PrePwdErrCountKey, req.Account)
	result, err := cache.RedisClient.Get(limitKey).Result()
	if err != nil && err.Error() != constants.RedisNilValue {
		ss_log.Error("判断密码错误是否超过规定次数,查询redis失败,err: %s", err.Error())
		resp.ResultCode = ss_err.ERR_ACCOUNT_ERR_PWD_LIMIT
		return nil
	}
	count := strext.ToInt(result)
	if b, limitCount := isLimit(count); b {
		ss_log.Error("登录受限制,当前错误次数为 %d,限制的次数为: %d", count, limitCount)
		resp.ResultCode = ss_err.ERR_ACCOUNT_ERR_PWD_LIMIT
		return nil
	}

	// 密码正确清除密码错误次数限制
	delErr := cache.RedisClient.Del(limitKey).Err()
	if delErr != nil {
		ss_log.Error("Del err: %s", delErr.Error())
	}

	// 进行风控检测
	/* 暂时未使用到 modify by xiaoyanchun
	riskReply, rErr := i.RiskCtrlHandleInst.Client.Login(context.TODO(), &go_micro_srv_riskctrl.LoginRequest{
		Uid:      accUid.String,
		Ip:       req.Ip,
		DeviceId: req.Uuid,
	})

	if rErr == nil { // 请求风控系统成功
		if riskReply.ResultCode == ss_err.ERR_SUCCESS { // 请求成功
			if riskReply.OpResult == constants.Risk_Result_No_Pass_Str { // 风控不通过
				ss_log.Error("Login请求风控检测不通过,Uid:%s", accUid.String)
				resp.ResultCode = ss_err.ERR_RISK_IS_RISK
				return nil
			}
		} else {
			ss_log.Error("Login请求风控检测失败,Uid:%s,ResultCode:%v,Msg:%s", accUid.String, riskReply.ResultCode, riskReply.Msg)
		}
	} else {
		ss_log.Error("Login请求风控检测出错,Uid:%s,Err:%v", accUid.String, rErr)
	}
	*/

	// 判断关联关系
	idenNo := dao.RelaAccIdenDaoInst.GetIdenFromAcc(accUid.String, req.AccountType)
	if idenNo == "" {
		resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
		return nil
	}

	getAcc := go_micro_srv_auth.GetAccountReply{}
	_ = r.GetAccount(ctx, &go_micro_srv_auth.GetAccountRequest{
		AccountUid: accUid.String,
	}, &getAcc)

	if getAcc.ResultCode != ss_err.ERR_SUCCESS {
		resp.ResultCode = getAcc.ResultCode
		return nil
	}

	if getAcc.UseStatus != "1" {
		dao.LogLoginDaoInstance.InsertLogLogin(accUid.String, "", ss_err.ERR_Account_UseStatus_Frozen_FAILD, req.Lat, req.Lng)
		resp.ResultCode = ss_err.ERR_Account_UseStatus_Frozen_FAILD
		return nil
	}

	retMap := common.JwtStructToMapApp(ss_struct.JwtDataApp{
		Account:        account,
		AccountUid:     accUid.String,
		IdenNo:         idenNo,
		AccountType:    req.AccountType,
		LoginAccountNo: accUid.String,
		PubKey:         req.PubKey,
		JumpIdenNo:     "",
		JumpIdenType:   "",
		MasterAccNo:    "",
		IsMasterAcc:    "1",
		PosSn:          req.PosSn,
	})

	jwt2 := jwt.GetNewEncryptedJWTToken(constants.JwtLifeTime, retMap, constants.MobileKeyJwtAes, constants.MobileKeyJwtSign)

	dao.LogLoginDaoInstance.InsertLogLogin(accUid.String, req.Ip, ss_err.ERR_SUCCESS, req.Lat, req.Lng)

	// 记录设备和pubkey关联
	if errStr := dao.RelaImeiPubkeyDaoInst.InsertRelaImeiPubKey(req.Imei, req.PubKey); errStr != ss_err.ERR_SUCCESS {
		ss_log.Error("%s", "设备和 pub_key 关联失败")
		resp.ResultCode = errStr
		return nil
	}

	if req.AccountType == constants.AccountType_SERVICER || req.AccountType == constants.AccountType_POS { // 收银员或者服务商登录
		var cashierName, terminalNumber, servicerNoT sql.NullString
		if req.AccountType == constants.AccountType_POS {
			// 根据收银员ID去cashier中找到收银员姓名和服务商
			err = ss_sql.QueryRow(dbHandler, "SELECT name,servicer_no FROM cashier WHERE uid=$1 LIMIT 1", []*sql.NullString{&cashierName, &servicerNoT}, idenNo)
			if err != nil {
				ss_log.Error("err=[%v]", err)
			}
		} else {
			cashierName.String = ""
			servicerNoT.String = idenNo
		}

		// todo 刷新商户的 pub_key
		//if errStr := dao.ServicerDaoInstance.ModifyServicerPubKey(servicerNoT.String, req.PubKey); errStr != ss_err.ERR_SUCCESS {
		//	ss_log.Error("%s", "刷新商户的 pub_key 失败")
		//	resp.ResultCode = errStr
		//	return nil
		//}

		// todo 测试阶段注释下面代码
		// ------------------------------------
		// 根据前端传的pos机器码sn找到number和服务商
		var servicerNoPos sql.NullString
		err = ss_sql.QueryRow(dbHandler, "SELECT terminal_number,servicer_no  FROM servicer_terminal WHERE pos_sn=$1 and is_delete = 0  and use_status = 1 LIMIT 1", []*sql.NullString{&terminalNumber, &servicerNoPos}, req.PosSn)
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
		// 校验服务商是否正确
		if servicerNoT.String != servicerNoPos.String {
			ss_log.Error("服务商号对应不上,当前登录的服务商为------> %s,根据pos_sn获取的服务商号为---> %s", servicerNoT.String, servicerNoPos.String)
			resp.ResultCode = ss_err.ERR_TERMINALNUM_Account_NoRelation_FAILD
			return nil
		}
		// ------------------------------------

		resp.Phone = phone.String
		resp.ResultCode = ss_err.ERR_SUCCESS
		resp.CashierName = cashierName.String
		resp.TerminalNo = terminalNumber.String
		resp.AccountUid = accUid.String
		resp.Phone = phone.String
		resp.Jwt = jwt2
		resp.AccountType = req.AccountType
		resp.IsFirstLogin = isFirstLoginT.String
		resp.RefreshTokenInterval = constants.RefreshTokenInterval
		return nil
	}

	// todo 刷新用户的pub_key
	//if errStr := dao.CustDaoInstance.UpdateAccountPub(accUid.String, req.PubKey); errStr != ss_err.ERR_SUCCESS {
	//	ss_log.Error("%s", "刷新用户的 pub_key 失败")
	//	resp.ResultCode = ss_err.ERR_PARAM
	//	return nil
	//}
	// 查找最新的版本号
	latDescription, latVersion, latAppURL, latVsCode, isForce := dao.AppVersionDaoInstance.QueryAppVersion("", constants.APPVERSIONVSTYPE_MANGOPAY_APP)
	// 历史版本信息
	_, _, _, hisVsCode, isForce := dao.AppVersionDaoInstance.QueryAppVersion(req.AppVersion, constants.APPVERSIONVSTYPE_MANGOPAY_APP)
	var appVersion, versionDescription, appUrl string
	if strext.ToFloat64(hisVsCode) < strext.ToFloat64(latVsCode) {
		appVersion = latVersion
		versionDescription = latDescription
		appUrl = latAppURL
	}

	//查询支付密码
	isSetPaymentPwd, err := dao.CustDaoInstance.GetPaymentPwdIsNull(idenNo)
	if err != nil {
		if err != nil {
			ss_log.Error("查询用户是否设置支付密码失败, custNo=%v, err=%v", idenNo, err)
			resp.ResultCode = ss_err.ERR_SYSTEM
			return nil
		}
	}

	resp.ResultCode = ss_err.ERR_SUCCESS
	resp.AccountUid = accUid.String
	resp.Jwt = jwt2
	resp.IsEncrypted = true
	resp.GenKey = genKey.String
	resp.Phone = phone.String
	resp.AccountType = req.AccountType
	resp.HeadPortraitImgNo = headPortraitImgNoT.String
	resp.AppVersion = appVersion
	resp.VersionDescription = versionDescription
	resp.AppUrl = appUrl
	resp.PubKey = strext.ToStringNoPoint(common.EncryptMap["back_pub_key"])
	resp.IsForce = isForce
	resp.RefreshTokenInterval = constants.RefreshTokenInterval
	resp.IsSetPaymentPwd = isSetPaymentPwd
	return nil
}

// pwdErrCount 密码试错
func pwdErrLimit(account string) bool {
	key := cache.GetPwdErrCountKey(cache.PrePwdErrCountKey, account)
	result, err := cache.RedisClient.Get(key).Result()
	if err != nil && err.Error() != constants.RedisNilValue {
		return false
	}
	count := strext.ToInt(result)
	count = count + 1
	// 设置错误次数
	cache.RedisClient.Set(key, count, 0)
	if b, limitCount := isLimit(count); b {
		ss_log.Error("登录受限制,当前错误次数为 %d,限制的次数为: %d", count, limitCount)
		return false
	}
	return true
}
func isLimit(count int) (bool, int) {
	var limitCount int
	// 判断次数
	value := dao.GlobalParamDaoInstance.QeuryParamValue(constants.GlobalParamKeyLoginPwdErrCount)
	if value == "" {
		limitCount = constants.ErrPwdLimitDefaultCount // 默认5条
	} else {
		limitCount = strext.ToInt(value)
	}
	return count >= limitCount, limitCount
}

func (r *Auth) Login2Nd(ctx context.Context, req *go_micro_srv_auth.Login2NdRequest, resp *go_micro_srv_auth.Login2NdReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	ss_log.Info("req=[%v]", req)
	if req.LoginAccountNo == "" || ((req.JumpIdenType == "" || req.JumpIdenNo == "") && req.JumpAccNo == "") {
		resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
		resp.IsEncrypted = false
		return nil
	}

	var masterAccNo string
	isMasterAcc := "1"
	if req.JumpAccNo != "" {
		getAccJump := go_micro_srv_auth.GetAccountReply{}
		_ = r.GetAccount(ctx, &go_micro_srv_auth.GetAccountRequest{
			AccountUid: req.JumpAccNo,
		}, &getAccJump)
		if getAccJump.Uid == "" {
			resp.ResultCode = ss_err.ERR_ACCOUNT_NO_PERMISSION
			resp.IsEncrypted = false
			return nil
		}
		req.JumpIdenNo = getAccJump.MerchantUid
		req.JumpIdenType = getAccJump.AccountType
		if getAccJump.MasterAcc != "" && getAccJump.MasterAcc != ss_sql.UUID {
			masterAccNo = getAccJump.MasterAcc
			isMasterAcc = "0"
		}
	}

	ss_log.Info("req.JumpIdenNo=[%v], req.JumpIdenType=[%v]", req.JumpIdenNo, req.JumpIdenType)
	isExists := dao.AccDaoInstance.HasJumpRela(req.LoginAccountNo, req.JumpIdenNo, req.JumpIdenType)
	if !isExists {
		resp.ResultCode = ss_err.ERR_ACCOUNT_NF_IDEN_NO
		resp.IsEncrypted = false
		return nil
	}

	accountNo := req.JumpAccNo
	if req.JumpAccNo == "" {
		accountNoT, err := dao.AccDaoInstance.GetAccNoFromIden(req.JumpIdenNo, req.JumpIdenType)
		if err != nil {
			ss_log.Error("login failed|err=[%v]", err)
			resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
			resp.IsEncrypted = false
			return nil
		}
		accountNo = accountNoT
	}

	getAcc := go_micro_srv_auth.GetAccountReply{}
	_ = r.GetAccount(ctx, &go_micro_srv_auth.GetAccountRequest{
		AccountUid: accountNo,
	}, &getAcc)
	if getAcc.UseStatus != "1" {
		resp.ResultCode = ss_err.ERR_ACCOUNT_NO_PERMISSION
		resp.IsEncrypted = false
		return nil
	}

	getAccLogin := go_micro_srv_auth.GetAccountReply{}
	_ = r.GetAccount(ctx, &go_micro_srv_auth.GetAccountRequest{
		AccountUid: req.LoginAccountNo,
	}, &getAccLogin)
	if getAccLogin.UseStatus != "1" {
		resp.ResultCode = ss_err.ERR_ACCOUNT_NO_PERMISSION
		resp.IsEncrypted = false
		return nil
	}

	retMap := common.JwtStructToMapWebAdmin(ss_struct.JwtDataWebAdmin{
		Account:        getAcc.Account,
		AccountUid:     accountNo,
		IdenNo:         "",
		AccountType:    getAcc.AccountType,
		LoginAccountNo: req.LoginAccountNo,
		JumpIdenNo:     req.JumpIdenNo,
		JumpIdenType:   req.JumpIdenType,
		MasterAccNo:    masterAccNo,
		IsMasterAcc:    isMasterAcc,
	})

	k1, loginSignKey, err := cache.ApiDaoInstance.GetGlobalParam("login_sign_key")
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k1)
	}
	k1, loginAesKey, err := cache.ApiDaoInstance.GetGlobalParam("login_aes_key")
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k1)
	}

	jwt2 := jwt.GetNewEncryptedJWTToken(-1, retMap, loginAesKey, loginSignKey)
	resp.ResultCode = ss_err.ERR_SUCCESS
	resp.IsEncrypted = true
	resp.LoginNickName = getAccLogin.Nickname
	resp.JumpAccNo = accountNo
	resp.Jwt = jwt2
	return nil
}

func (r *Auth) LoginReturn(ctx context.Context, req *go_micro_srv_auth.LoginReturnRequest, resp *go_micro_srv_auth.LoginReturnReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	if req.LoginAccountNo == "" || req.JumpIdenType == "" || req.JumpIdenNo == "" {
		resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
		resp.IsEncrypted = false
		return nil
	}

	isExists := dao.AccDaoInstance.HasJumpRela(req.LoginAccountNo, req.JumpIdenNo, req.JumpIdenType)
	if !isExists {
		resp.ResultCode = ss_err.ERR_ACCOUNT_NF_IDEN_NO
		resp.IsEncrypted = false
		return nil
	}

	getAccLogin := go_micro_srv_auth.GetAccountReply{}
	_ = r.GetAccount(ctx, &go_micro_srv_auth.GetAccountRequest{
		AccountUid: req.LoginAccountNo,
	}, &getAccLogin)
	if getAccLogin.UseStatus != "1" {
		resp.ResultCode = ss_err.ERR_ACCOUNT_NO_PERMISSION
		resp.IsEncrypted = false
		return nil
	}

	masterAccNo := ""
	isMasterAcc := "1"
	if getAccLogin.MasterAcc != "" && getAccLogin.MasterAcc != ss_sql.UUID {
		masterAccNo = getAccLogin.MasterAcc
		isMasterAcc = "0"
	}

	retMap := common.JwtStructToMapWebAdmin(ss_struct.JwtDataWebAdmin{
		Account:        getAccLogin.Account,
		AccountUid:     req.LoginAccountNo,
		IdenNo:         "",
		AccountType:    getAccLogin.AccountType,
		LoginAccountNo: req.LoginAccountNo,
		JumpIdenNo:     "",
		JumpIdenType:   "",
		MasterAccNo:    masterAccNo,
		IsMasterAcc:    isMasterAcc,
	})

	k1, loginSignKey, err := cache.ApiDaoInstance.GetGlobalParam("login_sign_key")
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k1)
	}
	k1, loginAesKey, err := cache.ApiDaoInstance.GetGlobalParam("login_aes_key")
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k1)
	}
	jwt2 := jwt.GetNewEncryptedJWTToken(-1, retMap, loginAesKey, loginSignKey)
	resp.ResultCode = ss_err.ERR_SUCCESS
	resp.IsEncrypted = true
	resp.LoginNickName = getAccLogin.Nickname
	resp.Jwt = jwt2
	return nil
}

/**
重置密码
*/
func (r *Auth) ResetPw(ctx context.Context, req *go_micro_srv_auth.ResetPwRequest, reply *go_micro_srv_auth.ResetPwReply) error {
	// 检查是否存在
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	if req.Cli == constants.LOGIN_CLI_WEB {
		redisKey, defPass, err := cache.ApiDaoInstance.GetGlobalParam("def_password")
		if err != nil {
			ss_log.Error("redisKey=[%v],err=[%v]", redisKey, err)
			reply.ResultCode = ss_err.ERR_ACCOUNT_SMS_CODE
			return nil
		}
		req.NewPw = defPass

		k1, passwordSalt, err := cache.ApiDaoInstance.GetGlobalParam("password_salt")
		if err != nil {
			ss_log.Error("err=[%v],missing key=[%v]", err, k1)
		}
		var tmp sql.NullString
		pwdMD5 := encrypt.DoMd5Salted(req.LoginPw, passwordSalt)
		err = ss_sql.QueryRow(dbHandler, "select 1 from account where uid = $1 and password=$2 limit 1", []*sql.NullString{&tmp}, req.LoginAccNo, pwdMD5)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}
		if tmp.String != "1" {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
			return nil
		}
	} else {
		redisKey, isNoNeedSmsRegChk, err := cache.ApiDaoInstance.GetGlobalParam(constants.GlparamNoSmsRegChk)
		if err != nil {
			ss_log.Error("err=[%v]\nredisKey=[%v]", err, redisKey)
		}
		if isNoNeedSmsRegChk != "1" {
			// 验证码检查
			isFitCode := util.DoChkSmsCode(req.SmsCode, req.Account)
			if !isFitCode {
				reply.ResultCode = ss_err.ERR_ACCOUNT_SMS_CODE
				return nil
			}
		}
	}

	var tmp sql.NullString
	err := ss_sql.QueryRow(dbHandler, "select 1 from account where uid = $1 limit 1", []*sql.NullString{&tmp}, req.Account)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	// 已存在
	if tmp.String != "" {
		redisKey, defPass, err := cache.ApiDaoInstance.GetGlobalParam("def_password")
		if err != nil {
			ss_log.Error("redisKey=[%v],err=[%v]", redisKey, err)
			reply.ResultCode = ss_err.ERR_ACCOUNT_SMS_CODE
			return nil
		}
		req.NewPw = defPass

		k1, passwordSalt, err := cache.ApiDaoInstance.GetGlobalParam("password_salt")
		if err != nil {
			ss_log.Error("err=[%v],missing key=[%v]", err, k1)
		}
		pwdMD5 := encrypt.DoMd5Salted(req.NewPw, passwordSalt)
		err = ss_sql.Exec(dbHandler, "update account set password=$1 where uid=$2", pwdMD5, req.Account)
		if err != nil {
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}

		//cache.RedisCli.Del(constants.DefPoolName, util.MkSmsCodeName(req.SmsCode))

		cache.RedisClient.Del(util.MkSmsCodeName(req.SmsCode))
		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	}

	reply.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
	return nil
}

/**
修改密码
*/
func (*Auth) ModifyPwMerc(ctx context.Context, req *go_micro_srv_auth.ModifyPwMercRequest, reply *go_micro_srv_auth.ModifyPwMercReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var accountType, mercNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select account_type from account where uid=$1`, []*sql.NullString{&accountType}, req.AccountLogin)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	var tmp sql.NullString

	switch accountType.String {
	case constants.AccountType_SERVICER:
		err := ss_sql.QueryRow(dbHandler, `select merchant_no from merchant where account_no=$1 limit 1`, []*sql.NullString{&mercNo}, req.AccountLogin)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}

		switch req.IdenType {
		default:
		}
	}

	if tmp.String == "" {
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	k1, passwordSalt, err := cache.ApiDaoInstance.GetGlobalParam("password_salt")
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k1)
	}
	pwdMD5 := encrypt.DoMd5Salted(req.NewPw, passwordSalt)
	err = ss_sql.Exec(dbHandler, `update account set "password"=$1 where affiliation_uid=$2 and account_type=$3 `, pwdMD5, req.IdenNo, req.IdenType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*Auth) ModifyPw(ctx context.Context, req *go_micro_srv_auth.ModifyPwRequest, reply *go_micro_srv_auth.ModifyPwReply) error {
	switch req.Cli {
	case constants.LOGIN_CLI_MICROPRO:
		redisKey, isNoNeedSmsRegChk, err := cache.ApiDaoInstance.GetGlobalParam(constants.GlparamNoSmsRegChk)
		if err != nil {
			ss_log.Error("err=[%v]\nredisKey=[%v]", err, redisKey)
		}
		if isNoNeedSmsRegChk != "1" {
			// 验证码检查
			isFitCode := util.DoChkSmsCode(req.SmsCode, req.Account)
			if !isFitCode {
				reply.ResultCode = ss_err.ERR_ACCOUNT_SMS_CODE
				return nil
			}
		}

		dbHandler := db.GetDB(constants.DB_CRM)
		defer db.PutDB(constants.DB_CRM, dbHandler)
		// 检查是否有权限
		if req.TopAgency != "" {
			tmp, err := dao.AccDaoInstance.GetAccountCnt(req.Account)
			if err != nil || tmp <= 0 {
				ss_log.Error("err=[%v]", err)
				reply.ResultCode = ss_err.ERR_ACCOUNT_NO_PERMISSION
				return nil
			}
		}

		// 检查是否存在
		var tmp sql.NullString
		k1, passwordSalt, err := cache.ApiDaoInstance.GetGlobalParam("password_salt")
		if err != nil {
			ss_log.Error("err=[%v],missing key=[%v]", err, k1)
		}
		oldPwdMD5 := encrypt.DoMd5Salted(req.OldPw, passwordSalt)
		err = ss_sql.QueryRow(dbHandler, `select 1 from account where account = $1 and "password" = $2 and is_delete='0' limit 1`, []*sql.NullString{&tmp}, req.Account, oldPwdMD5)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
			return nil
		}

		ss_log.Info("ModifyPw|tmp=[%v]\n", tmp)
		// 已存在
		if tmp.String != "" {
			k1, passwordSalt, err := cache.ApiDaoInstance.GetGlobalParam("password_salt")
			if err != nil {
				ss_log.Error("err=[%v],missing key=[%v]", err, k1)
			}
			pwdMD5 := encrypt.DoMd5Salted(req.NewPw, passwordSalt)
			err = ss_sql.Exec(dbHandler, `update account set "password"=$1 where account=$2 and is_delete='0'`, pwdMD5, req.Account)

			if err != nil {
				ss_log.Error("err=[%v]", err)
				reply.ResultCode = ss_err.ERR_SYS_DB_OP
				return nil
			}

			//cache.RedisCli.Del(constants.DefPoolName, util.MkSmsCodeName(req.SmsCode))
			cache.RedisClient.Del(util.MkSmsCodeName(req.SmsCode))
			reply.ResultCode = ss_err.ERR_SUCCESS
			return nil
		}
	case constants.LOGIN_CLI_WEB:
		// 检查是否存在
		dbHandler := db.GetDB(constants.DB_CRM)
		defer db.PutDB(constants.DB_CRM, dbHandler)
		var tmp sql.NullString

		k1, passwordSalt, err := cache.ApiDaoInstance.GetGlobalParam("password_salt")
		if err != nil {
			ss_log.Error("err=[%v],missing key=[%v]", err, k1)
		}
		oldPwdMD5 := encrypt.DoMd5Salted(req.OldPw, passwordSalt)
		err = ss_sql.QueryRow(dbHandler, `select 1 from account where uid = $1 and "password" = $2 limit 1`, []*sql.NullString{&tmp}, req.Account, oldPwdMD5)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_ACCOUNT_OLD_PWD_FAILD
			return nil
		}

		ss_log.Error("ModifyPw|tmp=[%v]\n", tmp.String)
		// 已存在
		if tmp.String != "" {
			k1, passwordSalt, err := cache.ApiDaoInstance.GetGlobalParam("password_salt")
			if err != nil {
				ss_log.Error("err=[%v],missing key=[%v]", err, k1)
			}
			pwdMD5 := encrypt.DoMd5Salted(req.NewPw, passwordSalt)
			err = ss_sql.Exec(dbHandler, `update account set "password"=$1 where uid=$2`, pwdMD5, req.Account)

			if err != nil {
				ss_log.Error("err=[%v]", err)
				reply.ResultCode = ss_err.ERR_SYS_DB_OP
				return nil
			}

			//cache.RedisCli.Del(constants.DefPoolName, util.MkSmsCodeName(req.SmsCode))

			cache.RedisClient.Del(util.MkSmsCodeName(req.SmsCode))
			reply.ResultCode = ss_err.ERR_SUCCESS
			return nil
		}
	}
	reply.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
	return nil
}

func (*Auth) UpdateOrInsertAccountAuth(ctx context.Context, req *go_micro_srv_auth.UpdateOrInsertAccountAuthRequest, resp *go_micro_srv_auth.UpdateOrInsertAccountAuthReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, err := dbHandler.Begin()
	if nil != err {
		tx.Rollback()
		resp.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}
	defer ss_sql.Rollback(tx)
	switch req.AccountType {
	case constants.AccountType_ADMIN:
	case constants.AccountType_OPERATOR:
	default:
		ss_log.Error("登陆的账号角色不是管理员或运营,无权更改权限")
		resp.ResultCode = ss_err.ERR_SYS_NO_API_AUTH
		return nil
	}

	// 实际执行
	errDelete := dao.RelaAccountRoleDaoInst.DeleteRelaAccountRoleTx(tx, req.Uid)
	if ss_err.ERR_SUCCESS != errDelete {
		ss_log.Error("err=%v", errDelete)
		resp.ResultCode = errDelete
		return nil
	}

	account := dao.AccDaoInstance.GetAccountByUid(req.Uid)
	if account == "" {
		ss_log.Error("获取账号失败，uid:[%v]", req.Uid)
		account = req.Uid
	}
	description := fmt.Sprintf("授权账号[%v]的角色为[", account)

	if req.Roles != "" {
		roles := strings.Split(req.Roles, ",")
		for _, role := range roles {
			accType, err := dao.RelaAccountRoleDaoInst.GetRoleAccTypeTx(tx, role)
			if ss_err.ERR_SUCCESS != err {
				ss_log.Error("err=%v", err)
				resp.ResultCode = err
				return nil
			}

			switch accType {
			case constants.AccountType_ADMIN:
				description = fmt.Sprintf("%v %v", description, "管理员")

			case constants.AccountType_OPERATOR:
				description = fmt.Sprintf("%v %v", description, "运营")
			case constants.AccountType_SERVICER: //如果要授权账号角色是服务商
				description = fmt.Sprintf("%v %v", description, "服务商")
				//查询该账号有没有创建过服务商
				serNo := dao.RelaAccIdenDaoInst.GetIdenFromAcc(req.Uid, constants.AccountType_SERVICER)
				if serNo == "" { //如果没有关联关系,说明该账号未创建过服务商
					//创建服务商(servicer)
					servicerNo, err := dao.ServicerDaoInstance.InsertInitService(tx, req.Uid)
					if err != nil {
						ss_log.Error("添加账号角色为服务商,插入服务商表失败,err: %s", err.Error())
						resp.ResultCode = ss_err.ERR_SYS_DB_OP
						return nil
					}
					// 建立关联关系
					errCode := dao.RelaAccIdenDaoInst.InsertRelaAccIdenTx(tx, req.Uid, servicerNo, constants.AccountType_SERVICER)
					if errCode != ss_err.ERR_SUCCESS {
						ss_log.Info("errCode==[%v]", errCode)
						resp.ResultCode = errCode
						return nil
					}
				} else {
					ss_log.Error("账号已添加过服务商id[%v],不允许再创建服务商。", serNo)
				}
			case constants.AccountType_USER:
				description = fmt.Sprintf("%v %v", description, "用户")

			case constants.AccountType_POS:
				description = fmt.Sprintf("%v %v", description, "店员")

			case constants.AccountType_Headquarters:
				description = fmt.Sprintf("%v %v", description, "总部")

			default:
				ss_log.Error("未知账号类型[%v]", accType)
			}

			//为账号添加角色
			errCode := dao.RelaAccountRoleDaoInst.InsertRelaAccountRoleTx(tx, req.Uid, role)
			if ss_err.ERR_SUCCESS != errCode {
				ss_log.Error("err=%v", errCode)
				resp.ResultCode = errCode
				return nil
			}
		}

	}

	description = fmt.Sprintf("%v]", description)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Account_Menu)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		resp.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}
	ss_sql.Commit(tx)
	resp.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
 * 获取账户列表
 */
func (*Auth) GetAccountList(ctx context.Context, req *go_micro_srv_auth.GetAccountListRequest, resp *go_micro_srv_auth.GetAccountListReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	condStr, _ := json.Marshal([]*model.WhereSqlCond{
		{Key: "acc.account", Val: req.Search, EqType: "like"},
		{Key: "acc.nickname", Val: req.Search, EqType: "like"},
	})

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "acc.use_status", Val: "1", EqType: "="},
		{Key: "", Val: string(condStr), EqType: "or_group"},
		{Key: "acc.master_acc", Val: req.MasterAcc, EqType: "="},
		{Key: "acc.is_delete", Val: "0", EqType: "="},
		{Key: "acc.is_actived", Val: req.IsActived, EqType: "="},
	})
	where := whereModel.WhereStr
	args := whereModel.Args

	var total sql.NullString
	sqlCnt := "SELECT count(1) FROM account acc " + where
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by acc.create_time desc")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
	where = whereModel.WhereStr
	args = whereModel.Args
	rows, stmt, err := ss_sql.Query(dbHandler, "SELECT acc.uid,acc.nickname,acc.account,acc.use_status,acc.create_time,"+
		"acc.modify_time,acc.master_acc,acc2.account,acc.is_actived "+
		"FROM account acc left join account acc2 on acc2.uid=acc.master_acc "+where, args...)
	if stmt != nil {
		defer stmt.Close()
		defer rows.Close()
	}
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}
	datas := []*go_micro_srv_auth.Account{}
	if err == nil {
		for rows.Next() {
			data := go_micro_srv_auth.Account{}
			var masterAccount, masterAcc, isActived sql.NullString
			err = rows.Scan(&data.Uid, &data.Nickname, &data.Account, &data.UseStatus, &data.CreateTime, &data.ModifyTime,
				&masterAcc, &masterAccount, &isActived)
			if err != nil {
				ss_log.Error("err=[%v]", err)
				continue
			}
			data.MasterAcc = masterAcc.String
			data.MasterAccount = masterAccount.String
			data.IsActived = isActived.String

			//accountTypes := dao.AccDaoInstance.GetAccountTypeByAccountNo(data.Uid)
			//
			////查询账号类型
			//data.AccountType = accountTypes

			datas = append(datas, &data)
		}
	}
	resp.ResultCode = ss_err.ERR_SUCCESS
	resp.AccountData = datas
	resp.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * 添加账户
 */
func (*Auth) SaveAccount(ctx context.Context, req *go_micro_srv_auth.SaveAccountRequest, resp *go_micro_srv_auth.SaveAccountReply) error {
	ss_log.Info("req == [%v]", req)
	if req.UseStatus == "" {
		req.UseStatus = "1"
	}

	if req.Phone == "" {
		resp.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//账号长度小于6位不给与创建
	if len(req.Account) < 6 {
		resp.ResultCode = ss_err.ERR_SAVE_ACCOUNT_LENGTH_FAILD
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx := ss_sql.BeginTx(dbHandler)
	if nil == tx {
		resp.ResultCode = ss_err.ERR_SYS_DB_INIT
		return nil
	}
	defer ss_sql.Rollback(tx)
	// 验证国家码和手机号是否唯一
	if err := dao.CountryCodePhoneDaoInst.Insert(tx, req.CountryCode, req.Phone); err != nil {
		ss_log.Error("新增或更改账号, 新增手机号和国家码进唯一表失败,accountUid: %s,err: %s", req.AccountUid, err.Error())
		resp.ResultCode = ss_err.ERR_ACCOUNT_ALREADY_EXISTS
		return nil
	}
	account := ss_func.ComposeAccountByPhoneCountryCode(req.Account, req.CountryCode)
	if req.AccountUid == "" {
		if req.Password == "" {
			resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_PASSWORD
			return nil
		}
		accountCount, err := dao.AccDaoInstance.CheckAccountTx(tx, req.Account)
		if nil != err {
			ss_log.Error("SaveAccount|CheckAccount|err=[%v]", err)
			resp.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}
		if accountCount != 0 {
			resp.ResultCode = ss_err.ERR_ACCOUNT_ALREADY_EXISTS
			return nil
		}
		ss_log.Info("SaveAccount |AddAccount req=%v", req)
		// 新建
		req.AccountUid, err = dao.AccDaoInstance.AddAccount(tx, req.Nickname, account, req.Password, req.UseStatus,
			req.MasterAcc, req.Phone, req.CountryCode, req.UtmSource, "")
		if nil != err {
			ss_log.Error("err=%v", err.Error())
			resp.ResultCode = ss_err.ERR_ACCOUNT_INIT_ACCOUNT_ERR
			return nil
		}
		//管理员和运营才添加关联关系
		if req.AccountType == constants.AccountType_ADMIN || req.AccountType == constants.AccountType_OPERATOR {
			errCode := dao.RelaAccIdenDaoInst.InsertRelaAccIdenTx(tx, req.AccountUid, "00000000-0000-0000-0000-000000000000", req.AccountType)
			if errCode != ss_err.ERR_SUCCESS {
				ss_log.Info("errCode==[%v]", errCode)
				resp.ResultCode = errCode
				return nil
			}
			ss_log.Info("AccountUid == [%v],AccountType == [%v]", req.AccountUid, req.AccountType)

			errCode = dao.AccDaoInstance.AuthAccountRetCode(tx, req.AccountType, req.AccountUid)
			if errCode != ss_err.ERR_SUCCESS {
				ss_log.Info("errCode==[%v]", errCode)
				resp.ResultCode = errCode
				return nil
			}
		}
		//以下是说后台可以创建服务商角色的账号，但现在l3说不可以创建服务商角色的账号，只能授权
		//else if req.AccountType == constants.AccountType_SERVICER { // 添加服务商需要建立关联关系,后台添加服务端是就相当于是修改操作了.
		//	servNo, err := dao.ServicerDaoInstance.InsertInitService(tx, req.AccountUid)
		//	if err != nil {
		//		ss_log.Error("添加账号操作,角色为服务商,插入服务商表失败,err: %s", err.Error())
		//		resp.ResultCode = ss_err.ERR_SYS_DB_OP
		//		return nil
		//	}
		//	// 建立关联关系
		//	errCode := dao.RelaAccIdenDaoInst.InsertRelaAccIden(tx, req.AccountUid, servNo, req.AccountType)
		//	if errCode != ss_err.ERR_SUCCESS {
		//		ss_log.Info("errCode==[%v]", errCode)
		//		resp.ResultCode = errCode
		//		return nil
		//	}
		//}
	} else {
		accountCount, err := dao.AccDaoInstance.CheckAccountUpdate(tx, req.AccountUid)
		if nil != err {
			ss_log.Error("SaveAccount|CheckAccount|err=[%v]", err.Error())
			resp.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}
		if accountCount == 0 {
			resp.ResultCode = ss_err.ERR_MSG_ACCOUNT_NOT_EXISTS
			return nil
		}
		// 删除以前的唯一
		oldCountryCode, oldPhone, err := dao.AccDaoInstance.GetCountryCodePhoneByUidTx(tx, req.AccountUid)
		if err != nil {
			ss_log.Error("err=[addAccount 查询原账号的手机号和国家码失败,uid: %s,err: %s]", req.AccountUid, err.Error())
			resp.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		if err := dao.CountryCodePhoneDaoInst.Delete(tx, oldCountryCode, oldPhone); err != nil {
			ss_log.Error("err=[addAccount 删除原账号的手机号和国家码失败, err: %s]", err.Error())
			resp.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		// 更新
		err = dao.AccDaoInstance.UpdateAccount(tx, req.Nickname, account, req.Password, req.UseStatus,
			req.AccountUid, req.Phone, req.Email)
		if nil != err {
			ss_log.Error("SaveAccount | UpdateAccount | err=[%v]", err)
			resp.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}
		ss_log.Info("修改用户[%v]信息成功", req.AccountUid)
	}

	ss_sql.Commit(tx)
	resp.ResultCode = ss_err.ERR_SUCCESS
	resp.Uid = req.AccountUid
	return nil
}

/**
 * 添加企业商家账户
 */
func (*Auth) SaveBusinessAccount(ctx context.Context, req *go_micro_srv_auth.SaveBusinessAccountRequest, resp *go_micro_srv_auth.SaveBusinessAccountReply) error {
	if req.Phone == "" {
		ss_log.Error("参数Phone为空")
		resp.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.Email == "" {
		ss_log.Error("参数Email为空")
		resp.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.Password == "" {
		ss_log.Error("参数Password为空")
		resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_PASSWORD
		return nil
	}

	if errStr := ss_func.CheckCountryCode(req.CountryCode); errStr != ss_err.ERR_SUCCESS {
		ss_log.Error("国家码[%v]不合法", req.CountryCode)
		resp.ResultCode = errStr
		return nil
	}

	//确认企业商家的手机号与国家码组成是唯一的,现在企业商家的手机号是存的business_phone
	if !dao.AccDaoInstance.CheckBusinessPhoneAndCountryCodeUnique(req.Phone, req.CountryCode) {
		ss_log.Error("account存在相同的商家手机号[%v]和国家码[%v]组合", req.Phone, req.CountryCode)
		resp.ResultCode = ss_err.ERR_BusinessPhone_Unique_FAILD
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx := ss_sql.BeginTx(dbHandler)
	if nil == tx {
		resp.ResultCode = ss_err.ERR_SYS_DB_INIT
		return nil
	}
	defer ss_sql.Rollback(tx)

	accountCount, err := dao.AccDaoInstance.CheckAccountTx(tx, req.Email)
	if nil != err {
		ss_log.Error("确认账号[%v]唯一性出错,err=[%v]", req.Email, err)
		resp.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}
	if accountCount != 0 {
		resp.ResultCode = ss_err.ERR_Account_Registered_FAILD
		return nil
	}

	// 新建账号
	accountUid, errAdd := dao.AccDaoInstance.AddBusinessAccount(tx, "", req.Email, req.Password, strext.ToStringNoPoint(constants.AccountUseStatusNormal), "", req.Phone, req.CountryCode, "", req.Email)
	if nil != errAdd {
		ss_log.Error("新增账号[%v]出错,err=%v", req.Email, errAdd.Error())
		resp.ResultCode = ss_err.ERR_ACCOUNT_INIT_ACCOUNT_ERR
		return nil
	}

	accountType := constants.AccountType_EnterpriseBusiness

	//添加商家business
	errAdd2, businessNo := dao.BusinessDaoInstance.AddBusinessTx(tx, accountUid, "", accountType)
	if errAdd2 != nil {
		ss_log.Error("新增商家出错,err=%v", errAdd.Error())
		resp.ResultCode = ss_err.ERR_ACCOUNT_INIT_ACCOUNT_ERR
		return nil
	}

	errCode := dao.RelaAccIdenDaoInst.InsertRelaAccIdenTx(tx, accountUid, businessNo, accountType)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Info("errCode==[%v]", errCode)
		resp.ResultCode = errCode
		return nil
	}

	ss_log.Info("AccountUid == [%v],AccountType == [%v]", accountUid, accountType)

	errCode = dao.AccDaoInstance.AuthAccountRetCode(tx, accountType, accountUid)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Info("errCode==[%v]", errCode)
		resp.ResultCode = errCode
		return nil
	}

	//初始化虚拟账户
	err = VAccountHandlerInst.InitVAccount(tx, accountUid)
	if err != nil {
		ss_log.Error("初始化企业商家虚账失败；accountNo=%v, err=%v", accountUid, err)
		resp.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	ss_sql.Commit(tx)

	//登录
	jwt2 := common.CreateWebBusinessJWT(ss_struct.JwtDataWebBusiness{
		Account:        req.Email,
		AccountUid:     accountUid,
		IdenNo:         businessNo, //BusinessNo
		AccountType:    accountType,
		LoginAccountNo: accountUid,
		Email:          req.Email,
		Phone:          req.Phone,
		CountryCode:    req.CountryCode,
		JumpIdenNo:     "",
		JumpIdenType:   "",
		MasterAccNo:    "",
		IsMasterAcc:    "1",
	})

	err = ss_sql.Exec(dbHandler, "insert into log_login(log_time,acc_no,ip,result,client,log_no) values (current_timestamp,$1,$2,$3,$4,$5)",
		accountUid, req.Ip, ss_err.ERR_SUCCESS, req.Client, strext.GetDailyId())
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	resp.ResultCode = ss_err.ERR_SUCCESS
	resp.Uid = accountUid
	resp.AccountType = accountType
	resp.Jwt = jwt2
	resp.InitPayPwdStatus = "0" //是否已初始化支付密码(0否，1是)
	return nil
}

/**
 * 添加操作员，注意!!此接口是后台与pos添加店员时共用一个接口。。
 */
func (*Auth) AddCashier(ctx context.Context, req *go_micro_srv_auth.AddCashierRequest, resp *go_micro_srv_auth.AddCashierReply) error {
	// 判断店员账号是否存在
	accNo := dao.AccDaoInstance.GetAccNoFromPhone(req.Phone, req.CountryCode)
	if accNo == "" {
		ss_log.Error("添加店员,根据手机号查询店员账号id失败")
		resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXIST
		return nil
	}

	//确认该服务商没有添加过该店员
	if dao.CashierDaoInstance.CheckServicerNoCashier(accNo, req.ServicerNo) {
		ss_log.Error("不允许重复添加账号[%v]为服务商[%v]的店员", accNo, req.ServicerNo)
		resp.ResultCode = ss_err.ERR_Servicer_Have_Cashier_FAILD
		return nil
	}

	//根据查询出来的店员账号去服务商表查询，看此账号是否是服务商（现不允许一个服务商添加另一个服务商为店员）
	servicerNo := dao.ServicerDaoInstance.GetServicerNoByAccNo(accNo)
	if servicerNo != "" { //不为空即表示账号有服务商身份
		ss_log.Error("不允许添加服务商[%v]为店员", servicerNo)
		resp.ResultCode = ss_err.ERR_ACCOUNT_IS_SERVICER
		return nil
	}

	//添加店员和关联关系
	cashierNo, relaErr := relaAccCachier(req.ServicerNo, accNo)
	if relaErr != ss_err.ERR_SUCCESS {
		ss_log.Error("账号关联店员失败,relaErr=[%v]", relaErr)
		resp.ResultCode = relaErr
		return nil
	}

	if req.LoginUid != "" { //后台管理系统调用的添加店员接口
		serAcc, err1 := dao.ServicerDaoInstance.GetAccountBySerNo(req.ServicerNo)
		if err1 != nil {
			ss_log.Error("获取服务商的账号失败")
			resp.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		cashierAccount := dao.AccDaoInstance.GetAccountByUid(accNo)
		description := fmt.Sprintf("为服务商[%v]添加店员[%v]", serAcc, cashierAccount)
		errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Servicer)
		if errAddLog != nil {
			ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
			resp.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}
	}

	resp.ResultCode = ss_err.ERR_SUCCESS
	resp.Uid = cashierNo
	return nil
}

// 关联用户和店员的关系
func relaAccCachier(servicerNo, accountUid string) (string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(context.TODO(), nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		return "", ss_err.ERR_SYS_DB_OP
	}
	defer ss_sql.Rollback(tx)

	//新增操作员信息
	cashierNo, err2 := dao.CashierDaoInstance.AddCashier(tx, servicerNo, accountUid)
	if err2 != nil {
		return "", ss_err.ERR_PARAM
	}

	//添加操作员与账号的关联关系
	errCode := dao.RelaAccIdenDaoInst.InsertRelaAccIdenTx(tx, accountUid, cashierNo, constants.AccountType_POS)
	if errCode != ss_err.ERR_SUCCESS {
		return "", errCode
	}
	ss_sql.Commit(tx)
	return cashierNo, ss_err.ERR_SUCCESS
}

func (*Auth) CheckAccount(ctx context.Context, req *go_micro_srv_auth.CheckAccountRequest, resp *go_micro_srv_auth.CheckAccountReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	data := &go_micro_srv_auth.CheckAccountData{}

	cnt, err := dao.AccDaoInstance.GetAccountCnt(req.Account)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}
	data.Count = strext.ToStringNoPoint(cnt)
	resp.ResultCode = ss_err.ERR_SUCCESS
	resp.Data = data
	return nil
}

/**
 * 删除账户
 */
func (*Auth) DeleteAccount(ctx context.Context, req *go_micro_srv_auth.DeleteAccountRequest, resp *go_micro_srv_auth.DeleteAccountReply) error {
	dao.AccDaoInstance.DeleteAccountList(req.AccountUids)
	// 不管是否成功
	resp.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
 * 获取账户
 */
func (*Auth) GetAccountByNickname(ctx context.Context, req *go_micro_srv_auth.GetAccountByNicknameRequest, resp *go_micro_srv_auth.GetAccountByNicknameReply) error {
	account := dao.AccDaoInstance.GetAccByNickname(req.Nickname)
	if account == nil {
		resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
		return nil
	}

	resp.Uid = account.Uid
	resp.Nickname = account.Nickname
	resp.UseStatus = account.UseStatus
	resp.Account = account.Account
	resp.CreateTime = account.CreateTime
	resp.ModifyTime = account.ModifyTime
	resp.DropTime = ""

	// 菜单
	_, resp.DataList = dao.AccDaoInstance.GetRoleAuthedUrlList(account.Uid)
	_, resp.DataList_2 = dao.AccDaoInstance.GetRoleAuthedUrlList(account.Uid)

	resp.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
 * 获取账户列表
 */
func (*Auth) GetRoleFromAcc(ctx context.Context, req *go_micro_srv_auth.GetRoleFromAccRequest, reply *go_micro_srv_auth.GetRoleFromAccReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	rowsRole, stmt, err2 := ss_sql.Query(dbHandler, "SELECT r.role_no,r.role_name FROM rela_account_role rela "+
		" left join role r on r.role_no=rela.role_uid WHERE account_uid = $1", req.AccNo)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rowsRole.Close()
	if nil != err2 {
		ss_log.Error("err|2=%v", err2)
		reply.ResultCode = ss_err.ERR_SYS_IO_ERR
		return nil
	}

	ds := []*go_micro_srv_auth.RoleSimpleData{}
	for rowsRole.Next() {
		d := go_micro_srv_auth.RoleSimpleData{}
		err := rowsRole.Scan(&d.RoleNo, &d.RoleName)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		ds = append(ds, &d)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = ds
	return nil
}

func (*Auth) GetLogLoginList(ctx context.Context, req *go_micro_srv_auth.GetLogLoginListRequest, reply *go_micro_srv_auth.GetLogLoginListReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlStr := "select lg.log_time,lg.acc_no,lg.ip,lg.result,lg.client,lg.log_no,a.account from log_login lg left join account a on lg.acc_no=a.uid"
	sqlCount := "select count(1) from log_login lg left join account a on lg.acc_no=a.uid"
	where := " where 1=1 "
	if req.AccountType == constants.AccountType_OPERATOR || req.AccountType == constants.AccountType_ADMIN {
		where += " and a.account like '%" + req.Search + "%' and a.nickname like '%" + req.Search + "%' "
	} else {
		where += " and lg.acc_no='" + req.AccountNo + "' "
	}
	if req.StartTime != "" && req.EndTime != "" {
		where += " AND log_time > TO_TIMESTAMP('" + req.StartTime + "', 'YYYY-MM-DD HH24:MI:SS') "
		where += " AND log_time <= TO_TIMESTAMP('" + req.EndTime + "', 'YYYY-MM-DD HH24:MI:SS') "
	}
	row, stmt, err := ss_sql.QueryRowN(dbHandler, sqlCount+where)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}
	var total int32
	err = row.Scan(&total)
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr+where+" LIMIT $1 OFFSET $2 ORDER BY log_time DESC ", req.PageSize, (req.Page-1)*req.PageSize)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	var datas []*go_micro_srv_auth.LogLogin
	if err == nil {
		for rows.Next() {
			var data go_micro_srv_auth.LogLogin
			var logtime string
			err = rows.Scan(&logtime, &data.AccountNo, &data.Ip, &data.Result, &data.Client, &data.LoginNo, &data.AccountName)
			data.LoginTime = ss_time.PostgresTimeToTime(logtime, global.Tz)
			datas = append(datas, &data)
		}
	} else {
		ss_log.Error("ApiHandler|GetApiList|err=%v\n{page,pageSize }=[%v]\n", err.Error(), []interface{}{
			req.Page, req.PageSize,
		})
		reply.ResultCode = ss_err.ERR_SUCCESS
		reply.Datas = datas
		reply.Total = total
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = total
	return nil
}

func (*Auth) GetRemain(ctx context.Context, req *go_micro_srv_auth.GetRemainRequest, reply *go_micro_srv_auth.GetRemainReply) error {
	khr, usd := dao.AccDaoInstance.GetRemain(req.AccountNo)
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = &go_micro_srv_auth.RemainData{
		Khr: khr,
		Usd: usd,
	}
	return nil
}

// 找回密码
func (*Auth) MobileBackPwd(ctx context.Context, req *go_micro_srv_auth.MobileBackPwdRequest, reply *go_micro_srv_auth.MobileBackPwdReply) error {
	// 校验短信验证码开关
	k1, isCheck, err := cache.ApiDaoInstance.GetGlobalParam("is_check_sms") // is_check_sms 0-需要校验,1-不需要校验
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k1)
	}

	// 校验短信验证码是否正确
	if isCheck == "0" {
		if !cache.CheckSMS(constants.BACKPWD, req.Phone, req.Sms) {
			reply.ResultCode = ss_err.ERR_ACCOUNT_SMS_MSG_FAILD
			return nil
		}
	}

	if err := dao.AccDaoInstance.UpdateAccountPwd(req.Phone, req.Password, req.CountryCode); err != nil {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	uid, _ := dao.AccDaoInstance.QueryUIDByPhone(req.Phone, req.CountryCode)

	// 插入日志
	description := "修改密码"
	if err := dao.LogDaoInstance.InsertAccountLog(description, uid, constants.AccountType_USER, constants.MODIFYACCOUNTLOGTYPE); err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, "插入修改密码的日志失败")
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (r *Auth) MobileModifyPwd(ctx context.Context, req *go_micro_srv_auth.MobileModifyPwdRequest, reply *go_micro_srv_auth.MobileModifyPwdReply) error {
	// 修改密码(包含密码校验)
	if err := dao.AccDaoInstance.UpdateAccountPwdByUID(req.Uid, req.OldPassword, req.NewPassword, req.NonStr); err != ss_err.ERR_SUCCESS {
		reply.ResultCode = err
		return nil
	}

	// 插入日志
	description := "修改密码"
	if err := dao.LogDaoInstance.InsertAccountLog(description, req.Uid, req.AccountType, constants.MODIFYACCOUNTLOGTYPE); err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, "插入修改密码的日志失败")
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (r *Auth) BusinessModifyPWDBySms(ctx context.Context, req *go_micro_srv_auth.BusinessModifyPWDBySmsRequest, reply *go_micro_srv_auth.BusinessModifyPWDBySmsReply) error {
	//校验短信验证码
	key := ss_data.GetSMSKey(constants.BACKPWD_Business, req.Phone)
	if value, err := cache.RedisClient.Get(key).Result(); err != nil || value != req.Sms {
		ss_log.Error("校验验证码错误")
		reply.ResultCode = ss_err.ERR_ACCOUNT_LOGIN_CODE
		return nil
	}

	uid := req.Uid
	if uid == "" { //如果是没有token情况下uid是为空的(忘记密码的找回密码，使用手机号设置新密码)，所以应该从前端传来的账户account查询出uid
		uid = dao.AccDaoInstance.GetAccNoFromAccount(req.Account)
		if uid == "" { //查询后仍等于空,说明账号查询不到
			ss_log.Error("查询账户uid出错,account=[%v]", req.Account)
			reply.ResultCode = ss_err.ERR_MSG_ACCOUNT_NOT_EXISTS
			return nil
		}
	}

	//修改登录密码
	if err := dao.AccDaoInstance.UpdateLoginPwdByUID(uid, req.Password); err != nil {
		ss_log.Error("修改登录密码错误，uid[%v],password[%v],err=[%v]", uid, req.Password, err)
		reply.ResultCode = ss_err.ERR_MODIFY_ACCOUNT_PWD_FAILD
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 修改支付密码
func (r *Auth) MobileModifyPayPwd(ctx context.Context, req *go_micro_srv_auth.MobileModifyPayPwdRequest, reply *go_micro_srv_auth.MobileModifyPayPwdReply) error {
	// 校验短信验证码开关
	k1, isCheck, err := cache.ApiDaoInstance.GetGlobalParam("is_check_sms") // is_check_sms 0-需要校验,1-不需要校验
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k1)
	}

	var phone, idenNo, function string
	var errGetPhone error
	switch req.AccountType {
	case constants.AccountType_SERVICER:
		fallthrough
	case constants.AccountType_POS:
		fallthrough
	case constants.AccountType_USER:
		function = constants.PAYPWD
		phone, idenNo, errGetPhone = dao.CashierDaoInstance.QueryPhoneAndCID(req.Uid, req.AccountType)
	case constants.AccountType_PersonalBusiness:
		function = constants.PAYPWD_Business
		phone, idenNo, errGetPhone = dao.CashierDaoInstance.QueryPhoneAndCID(req.Uid, req.AccountType)
	case constants.AccountType_EnterpriseBusiness:
		function = constants.PAYPWD_Business
		phone, idenNo, errGetPhone = dao.CashierDaoInstance.QueryBusinessPhoneAndCID(req.Uid, req.AccountType)
	}

	if errGetPhone != nil {
		ss_log.Error("查询账号的手机号和身份id出错,uid[%v],AccountType[%v],err=[%v]", req.Uid, req.AccountType, errGetPhone)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	// 校验短信验证码是否正确
	if isCheck == "0" {
		if !cache.CheckSMS(function, phone, req.Sms) {
			if req.AccountType == constants.AccountType_PersonalBusiness || req.AccountType == constants.AccountType_EnterpriseBusiness {
				//商家的提示语不一样
				reply.ResultCode = ss_err.ERR_Business_Verification_Code_FAILD
				return nil
			}

			reply.ResultCode = ss_err.ERR_ACCOUNT_SMS_MSG_FAILD
			return nil
		}
	}

	// 判断是服务商还是收银员
	switch req.AccountType {
	case constants.AccountType_SERVICER:
		if err := dao.ServicerDaoInstance.ModifyServicerPwdByUID(idenNo, req.Password); err != nil {
			ss_log.Error("err=[%v],missing key=[%v]", err.Error(), "修改服务商支付密码失败")
			reply.ResultCode = ss_err.ERR_MODIFY_PAY_PWD_FAILD
			return nil
		}
	case constants.AccountType_POS:
		// 修改密码
		if err := dao.CashierDaoInstance.ModifyCashierPwdByUID(idenNo, req.Password); err != nil {
			ss_log.Error("err=[%v],missing key=[%v]", err.Error(), "修改收银员支付密码失败")
			reply.ResultCode = ss_err.ERR_MODIFY_PAY_PWD_FAILD
			return nil
		}
	case constants.AccountType_USER: // 修改用户支付密码
		if err := dao.CustDaoInstance.ModifyCustPwdByUID(idenNo, req.Password); err != nil {
			ss_log.Error("err=[%v],missing key=[%v]", err.Error(), "修改用户支付密码失败")
			reply.ResultCode = ss_err.ERR_MODIFY_PAY_PWD_FAILD
			return nil
		}
	case constants.AccountType_PersonalBusiness: // 修改个人商家支付密码
		fallthrough
	case constants.AccountType_EnterpriseBusiness: // 修改企业商家支付密码
		if err := dao.BusinessDaoInstance.ModifyBusinessPayPwdByUid(idenNo, req.Password); err != nil {
			ss_log.Error("err=[%v],missing key=[%v]", err.Error(), "修改商家支付密码失败")
			reply.ResultCode = ss_err.ERR_MODIFY_PAY_PWD_FAILD
			return nil
		}
	}

	// 插入日志
	description := "修改支付密码"
	if err := dao.LogDaoInstance.InsertAccountLog(description, req.Uid, req.AccountType, constants.MODIFYPAYLOGTYPE); err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, "插入修改密码的日志失败")
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (r *Auth) MobileModifyPhone(ctx context.Context, req *go_micro_srv_auth.MobileModifyPhonedRequest, reply *go_micro_srv_auth.MobileModifyPhonedReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// 判断新的手机号和国家码是否唯一
	if err := dao.CountryCodePhoneDaoInst.Insert(tx, req.CountryCode, req.Phone); err != nil {
		ss_log.Error("MobileModifyPhone 新增手机号和国家码进唯一表失败,err: %s", err.Error())
		reply.ResultCode = ss_err.ERR_ACCOUNT_ALREADY_EXISTS
		return nil
	}

	// 校验短信验证码开关
	k1, isCheck, err := cache.ApiDaoInstance.GetGlobalParam("is_check_sms") // is_check_sms 0-需要校验,1-不需要校验
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k1)
	}

	// 校验短信验证码是否正确
	if isCheck == "0" {
		if !cache.CheckSMS(constants.MODIFYPHONE, req.Phone, req.Sms) {
			reply.ResultCode = ss_err.ERR_ACCOUNT_SMS_MSG_FAILD
			return nil
		}
	}

	// 判断是否存在此账号
	isExists := dao.AccDaoInstance.CheckAccountExists(ss_func.PreCountryCode(req.CountryCode), req.CountryCode, req.Phone)
	if isExists {
		ss_log.Error("err=[该手机号码已存在,CountryCode为: %s,reqPhone为: %s]", req.CountryCode, req.Phone)
		reply.ResultCode = ss_err.ERR_ACCOUNT_ALREADY_EXISTS
		return nil
	}

	// 删除以前的唯一
	oldCountryCode, oldPhone, err := dao.AccDaoInstance.GetCountryCodePhoneByUidTx(tx, req.Uid)
	if err != nil {
		ss_log.Error("err=[查询原账号的手机号和国家码失败,uid: %s,err: %s]", req.Uid, err.Error())
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	if err := dao.CountryCodePhoneDaoInst.Delete(tx, oldCountryCode, oldPhone); err != nil {
		ss_log.Error("err=[删除原账号的手机号和国家码失败, err: %s]", err.Error())
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 修改手机号
	account := ss_func.ComposeAccountByPhoneCountryCode(req.Phone, req.CountryCode)
	if err := dao.AccDaoInstance.ModifyPhone(tx, req.Uid, account, req.CountryCode, req.Phone); err != nil {
		ss_log.Error("修改手机号失败,err: %s", err.Error())
		reply.ResultCode = ss_err.ERR_MODIFY_PHONE_FAILD
		return nil
	}
	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (r *Auth) BusinessModifyPhone(ctx context.Context, req *go_micro_srv_auth.BusinessModifyPhoneRequest, reply *go_micro_srv_auth.BusinessModifyPhoneReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	switch req.AccountType {
	case constants.AccountType_PersonalBusiness: //个人商家的账号是国家码加手机号
		// 判断新的手机号和国家码是否唯一
		if err := dao.CountryCodePhoneDaoInst.Insert(tx, req.CountryCode, req.Phone); err != nil {
			ss_log.Error("MobileModifyPhone 新增手机号和国家码进唯一表失败,err: %s", err.Error())
			reply.ResultCode = ss_err.ERR_ACCOUNT_ALREADY_EXISTS
			return nil
		}

		// 校验短信验证码开关
		k1, isCheck, err := cache.ApiDaoInstance.GetGlobalParam("is_check_sms") // is_check_sms 0-需要校验,1-不需要校验
		if err != nil {
			ss_log.Error("err=[%v],missing key=[%v]", err, k1)
		}

		// 校验短信验证码是否正确
		if isCheck == "0" {
			if !cache.CheckSMS(constants.MODIFYPHONE_Business, req.Phone, req.Sms) {
				reply.ResultCode = ss_err.ERR_Business_Verification_Code_FAILD
				return nil
			}
		}

		// 判断是否存在此账号
		isExists := dao.AccDaoInstance.CheckAccountExists(ss_func.PreCountryCode(req.CountryCode), req.CountryCode, req.Phone)
		if isExists {
			ss_log.Error("err=[该手机号码已存在,CountryCode为: %s,reqPhone为: %s]", req.CountryCode, req.Phone)
			reply.ResultCode = ss_err.ERR_ACCOUNT_ALREADY_EXISTS
			return nil
		}

		// 删除以前的唯一
		oldCountryCode, oldPhone, err := dao.AccDaoInstance.GetCountryCodePhoneByUidTx(tx, req.Uid)
		if err != nil {
			ss_log.Error("err=[查询原账号的手机号和国家码失败,uid: %s,err: %s]", req.Uid, err.Error())
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		if err := dao.CountryCodePhoneDaoInst.Delete(tx, oldCountryCode, oldPhone); err != nil {
			ss_log.Error("err=[删除原账号的手机号和国家码失败, err: %s]", err.Error())
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		// 修改手机号
		account := ss_func.ComposeAccountByPhoneCountryCode(req.Phone, req.CountryCode)
		if err := dao.AccDaoInstance.ModifyPhone(tx, req.Uid, account, req.CountryCode, req.Phone); err != nil {
			ss_log.Error("修改手机号失败,err: %s", err.Error())
			reply.ResultCode = ss_err.ERR_MODIFY_PHONE_FAILD
			return nil
		}

		reply.Token = ""
	case constants.AccountType_EnterpriseBusiness:
		//校验新手机收到的验证码
		if !cache.CheckSMS(constants.MODIFYPHONE_Business, req.Phone, req.Sms) {
			reply.ResultCode = ss_err.ERR_ACCOUNT_SMS_MSG_FAILD
			return nil
		}

		if !dao.AccDaoInstance.CheckBusinessPhoneAndCountryCodeUnique(req.Phone, req.CountryCode) {
			ss_log.Error("account存在相同的商家手机号[%v]和国家码[%v]组合", req.Phone, req.CountryCode)
			reply.ResultCode = ss_err.ERR_BusinessPhone_Unique_FAILD
			return nil
		}

		if err := dao.AccDaoInstance.BusinessModifyBusinessPhoneTx(tx, req.Uid, req.Phone, req.CountryCode); err != nil {
			ss_log.Error("修改手机号失败,err: %s", err.Error())
			reply.ResultCode = ss_err.ERR_MODIFY_PHONE_FAILD
			return nil
		}

		//生成并返回更新后的jwt
		reply.Token = common.CreateWebBusinessJWT(ss_struct.JwtDataWebBusiness{
			Account:        req.Account,
			AccountUid:     req.Uid,
			IdenNo:         req.IdenNo, //BusinessNo
			AccountType:    req.AccountType,
			LoginAccountNo: req.Uid,
			Email:          req.Email,
			Phone:          req.Phone,
			CountryCode:    req.CountryCode,
			//JumpIdenNo:     "",
			//JumpIdenType:   "",
			//MasterAccNo:    "",
			IsMasterAcc: "1",
		})
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (r *Auth) ModifyPayPWDByMailCode(ctx context.Context, req *go_micro_srv_auth.ModifyPayPWDByMailCodeRequest, reply *go_micro_srv_auth.ModifyPayPWDByMailCodeReply) error {
	email := dao.AccDaoInstance.GetAccountMailByUid(req.Uid)
	if email == "" {
		ss_log.Error("查询出商家账号uid[%v]的email为空.", req.Uid)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if !cache.CheckMailCode(constants.Paypwd_By_Mail, email, req.MailCode) {
		ss_log.Error("修改账号uid[%v]的支付密码时验证邮箱验证码不正确", req.Uid)
		reply.ResultCode = ss_err.ERR_Business_Verification_Code_FAILD
		return nil
	}

	var idenNo string
	switch req.AccountType {
	case constants.AccountType_PersonalBusiness:
		idenNo = dao.RelaAccIdenDaoInst.GetIdenFromAcc(req.Uid, constants.AccountType_PersonalBusiness)
	case constants.AccountType_EnterpriseBusiness:
		idenNo = dao.RelaAccIdenDaoInst.GetIdenFromAcc(req.Uid, constants.AccountType_EnterpriseBusiness)
	}

	if idenNo == "" {
		ss_log.Error("查询账号uid[%v]的关系[%v]身份no为空。", req.Uid, req.AccountType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if err := dao.BusinessDaoInstance.ModifyBusinessPayPwdByUid(idenNo, req.Password); err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err.Error(), "修改商家支付密码失败")
		reply.ResultCode = ss_err.ERR_MODIFY_PAY_PWD_FAILD
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (r *Auth) ModifyPWDByMailCode(ctx context.Context, req *go_micro_srv_auth.ModifyPWDByMailCodeRequest, reply *go_micro_srv_auth.ModifyPWDByMailCodeReply) error {
	//校验邮箱验证码
	if !cache.CheckMailCode(constants.Backpwd_By_Mail, req.Email, req.MailCode) {
		ss_log.Error("修改账号uid[%v]的支付密码时验证邮箱验证码不正确", req.Uid)
		reply.ResultCode = ss_err.ERR_ACCOUNT_SMS_MSG_FAILD
		return nil
	}

	uid := req.Uid
	if uid == "" { //如果是没有token情况下uid是为空的(忘记密码的找回密码，使用邮件设置新密码)，所以应该从前端传来的账户account查询出uid
		uid = dao.AccDaoInstance.GetAccNoFromAccount(req.Account)
		if uid == "" { //查询后仍等于空,说明账号查询不到
			ss_log.Error("查询账户uid出错,account=[%v]", req.Account)
			reply.ResultCode = ss_err.ERR_MSG_ACCOUNT_NOT_EXISTS
			return nil
		}
	}

	//修改登录密码
	if err := dao.AccDaoInstance.UpdateLoginPwdByUID(uid, req.Password); err != nil {
		ss_log.Error("修改登录密码错误，uid[%v],password[%v],err=[%v]", uid, req.Password, err)
		reply.ResultCode = ss_err.ERR_MODIFY_ACCOUNT_PWD_FAILD
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (r *Auth) ModifyPayPWDByOldPwd(ctx context.Context, req *go_micro_srv_auth.ModifyPayPWDByOldPwdRequest, reply *go_micro_srv_auth.ModifyPayPWDByOldPwdReply) error {
	//验证旧密码
	if errStr := common.CheckPayPWD(req.AccountType, req.IdenNo, req.NonStr, req.OldPayPwd); errStr != ss_err.ERR_SUCCESS {
		//如果是输入密码错误, 增加错误密码日志（关键操作日志）
		if errStr == ss_err.ERR_DB_PWD {
			description := "验证支付密码错误"
			if err := dao.LogDaoInstance.InsertAccountLog(description, req.AccountUid, req.AccountType, constants.MODIFYACCOUNTLOGTYPE); err != nil {
				ss_log.Error("err=[%v],missing key=[%v]", err, "插入验证支付密码错误的日志失败")
			}

			reply.ResultCode = ss_err.ERR_Business_OLD_PayPWD_FAILD
			return nil
		}
		reply.ResultCode = errStr
		return nil
	}

	if err := dao.BusinessDaoInstance.ModifyBusinessPayPwdByUid(req.IdenNo, req.NewPayPwd); err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err.Error(), "修改商家支付密码失败")
		reply.ResultCode = ss_err.ERR_MODIFY_PAY_PWD_FAILD
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (r *Auth) ModifyEmail(ctx context.Context, req *go_micro_srv_auth.ModifyEmailRequest, reply *go_micro_srv_auth.ModifyEmailReply) error {

	//验证登录密码
	loginPwdDb := dao.AccDaoInstance.QueryPWDByAccountUid(req.Uid)
	// 数据库取出的密码加盐(加的是和前端传来的盐一样)
	pwdMD5FixedDB := encrypt.DoMd5Salted(loginPwdDb, req.NonStr)
	if req.Password != pwdMD5FixedDB {
		ss_log.Error("验证账户uid[%v]的登录密码出错，req.Password[%v]---pwdMD5FixedDB[%v]", req.Uid, req.Password, pwdMD5FixedDB)
		reply.ResultCode = ss_err.ERR_Business_OLD_PWD_FAILD
		return nil
	}

	//验证新邮箱收到的验证码
	if !cache.CheckMailCode(constants.ModifyEmail_By_NewEmail, req.Email, req.MailCode) {
		ss_log.Error("修改账号[%v]的邮箱为[%v]时验证新邮箱验证码[%v]不正确", req.Uid, req.Email, req.MailCode)
		reply.ResultCode = ss_err.ERR_Business_Verification_Code_FAILD
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}
	defer ss_sql.Rollback(tx)

	//修改邮箱
	if err := dao.AccDaoInstance.ModifyAccountEmailTx(tx, req.Uid, req.Email); err != nil {
		ss_log.Error("修改账号uid[%v]邮箱为[%v]失败,err=[%v]", req.Uid, req.Email, err)
		reply.ResultCode = ss_err.ERR_Modify_ACCOUNT_Mail_FAILD
		return nil
	}

	//修改账号（企业账号就是邮箱）
	if req.AccountType == constants.AccountType_EnterpriseBusiness {
		//校验邮箱是否重复(现企业商家邮箱就是账号)
		cnt, errCheck := dao.AccDaoInstance.CheckAccount(req.Email)
		if errCheck != nil {
			ss_log.Error("确认账号唯一性失败，account=[%v],err=[%v]", req.Email, errCheck)
			reply.ResultCode = ss_err.ERR_ACCOUNT_Mail_FAILD
			return nil
		}

		if cnt != 0 {
			ss_log.Error("账号已存在[%v]", req.Email)
			reply.ResultCode = ss_err.ERR_ACCOUNT_Mail_FAILD
			return nil
		}

		//修改账号
		if err := dao.AccDaoInstance.ModifyAccountByUidTx(tx, req.Uid, req.Email); err != nil {
			ss_log.Error("修改账号uid[%v]的账号account为[%v]失败,err=[%v]", req.Uid, req.Email, err)
			reply.ResultCode = ss_err.ERR_Modify_ACCOUNT_Mail_FAILD
			return nil
		}

		reply.Token = ""
	} else if req.AccountType == constants.AccountType_PersonalBusiness {
		//生成并返回更新后的jwt
		reply.Token = common.CreateWebBusinessJWT(ss_struct.JwtDataWebBusiness{
			Account:        req.Account,
			AccountUid:     req.Uid,
			IdenNo:         req.IdenNo, //BusinessNo
			AccountType:    req.AccountType,
			LoginAccountNo: req.Uid,
			Email:          req.Email,
			Phone:          req.Phone,
			CountryCode:    req.CountryCode,
			//JumpIdenNo:     "",
			//JumpIdenType:   "",
			//MasterAccNo:    "",
			IsMasterAcc: "1",
		})
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (r *Auth) ModifyNickname(ctx context.Context, req *go_micro_srv_auth.ModifyNicknameRequest, reply *go_micro_srv_auth.ModifyNicknameReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	if req.Nickname == "" {
		ss_log.Error("nickname is null")
		reply.ResultCode = ss_err.ERR_MODIFY_NICKNAME_ISNULL_FAILD
		return nil
	}

	// 修改昵称
	if err := ss_sql.Exec(dbHandler, `update account set nickname = $2 where uid = $1 and is_delete = 0`, req.AccountNo, req.Nickname); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		reply.ResultCode = ss_err.ERR_MODIFY_NICKNAME_FAILD
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// pos机新增银行卡
func (r *Auth) AddCard(ctx context.Context, req *go_micro_srv_auth.AddCardRequest, reply *go_micro_srv_auth.AddCardReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	vaType := 0
	switch req.MoneyType {
	case "usd":
		vaType = constants.VaType_USD_DEBIT
	case "khr":
		vaType = constants.VaType_KHR_DEBIT
	default:
		ss_log.Error("金额类型出错[%v]", req.MoneyType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//  判断卡是否存在
	if cardNo := dao.CardDaoInstance.QueryCardNo(req.RecCarNum); cardNo != "" {
		ss_log.Error("err=[银行卡号已存在,卡号为--->%s]", req.RecCarNum)
		reply.ResultCode = ss_err.ERR_CARD_IS_EXIST
		return nil
	}

	// 获取pos机的accountNo
	var posAccountNo string
	var channelType string
	var queryChannelType string
	switch req.AccountType {
	case constants.AccountType_SERVICER: //服务商
		queryChannelType = "('" + constants.CHANNEL_POS + "','" + constants.CHANNEL_ALL + "')"
		channelType = constants.CHANNEL_POS
		posAccountNo = req.AccountUid
	case constants.AccountType_POS: // 收银员
		queryChannelType = "('" + constants.CHANNEL_POS + "','" + constants.CHANNEL_ALL + "')"
		channelType = constants.CHANNEL_POS
		// 获取服务商no,再获取accountNo
		serviceNo := dao.CashierDaoInstance.GetServicerNoFromOpAccNo(req.OpAccNo)
		posAccountNo = dao.ServicerDaoInstance.GetAccountNoFromServiceNo(serviceNo)
	case constants.AccountType_USER: // 用户
		queryChannelType = "('" + constants.CHANNEL_USE_HEADQUARTERS + "','" + constants.CHANNEL_ALL + "')"
		channelType = constants.CHANNEL_USE_HEADQUARTERS
	}
	if posAccountNo == "" {
		ss_log.Error("账号不存在")
		reply.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXIST
		return nil
	}

	channelNo := dao.ChannelDaoInstance.QeuryChannelNoFromNameAndType(req.ChannelName, queryChannelType)
	if channelNo == "" {
		ss_log.Error("err=[--->%s]", "channelNo为空,无法保存")
		reply.ResultCode = ss_err.ERR_SAVE_CARD_FAILD
		return nil
	}

	var isDefault int32
	// 判断该用户名下该币种是否有银行卡.如果没有的话设置为默认卡,有的话,判断是否存在卡号
	if cardNo := dao.CardDaoInstance.QueryIsNewAddCard(posAccountNo, req.MoneyType); cardNo != "" {

		if req.IsDefault == constants.IS_DEFAULT_CARD {
			// 找到默认卡
			if cardNo, _ := dao.CardDaoInstance.QueryIsDefaultCardNo(posAccountNo, req.MoneyType, "1"); cardNo != "" {
				// 设置默认卡为非默认
				if errStr := dao.CardDaoInstance.UpdateIsDefault(tx, constants.NO_DEFAULT_CARD, cardNo); errStr != ss_err.ERR_SUCCESS {
					reply.ResultCode = ss_err.ERR_SAVE_CARD_FAILD
					return nil
				}
			}
		}
		isDefault = req.IsDefault
	} else {
		isDefault = constants.IS_DEFAULT_CARD
	}

	if errStr := dao.CardDaoInstance.InsertCard(tx, req.AccountUid, channelNo, req.RecName, req.RecCarNum, req.MoneyType, isDefault); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = ss_err.ERR_SAVE_CARD_FAILD
		return nil
	}
	// 记录日志
	descript := "新增银行卡"
	if errStr := dao.LogCardDaoInstance.InsertLogCard(tx, req.RecCarNum, req.RecName, req.AccountUid, channelNo, channelType, descript, vaType); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = ss_err.ERR_SAVE_CARD_FAILD
		return nil
	}
	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (r *Auth) ModifyDefaultCard(ctx context.Context, req *go_micro_srv_auth.ModifyDefaultCardRequest, reply *go_micro_srv_auth.ModifyDefaultCardReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	vaType := 0
	switch req.MoneyType {
	case "usd":
		vaType = constants.VaType_USD_DEBIT
	case "khr":
		vaType = constants.VaType_KHR_DEBIT
	default:
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 判断是服务商,销售员,用户
	var posAccountNo string
	var channelType string
	switch req.AccountType {
	case constants.AccountType_SERVICER: //服务商
		channelType = "2"
		posAccountNo = req.AccountUid
	case constants.AccountType_POS: // 收银员
		channelType = "2"
		// 获取服务商no,再获取accountNo
		serviceNo := dao.CashierDaoInstance.GetServicerNoFromOpAccNo(req.OpAccNo)
		posAccountNo = dao.ServicerDaoInstance.GetAccountNoFromServiceNo(serviceNo)
	case constants.AccountType_USER: // 用户
		channelType = "1"
		posAccountNo = req.AccountUid
	}
	if posAccountNo == "" {
		reply.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXIST
		return nil
	}

	// 判断需要修改的卡是否存在
	num, isDefault, channelNo, accountNo := dao.CardDaoInstance.QueryCardFromNo(req.CardNo)
	if posAccountNo != accountNo {
		ss_log.Error("err=[绑定的银行卡不是自己的卡,需要绑定的卡 id 为----->%s,这张卡的account_id为----->%s]", req.CardNo, accountNo)
		reply.ResultCode = ss_err.ERR_CARD_NOT_EXIST
		return nil
	}

	if num == "" || isDefault == "" || channelNo == "" {
		ss_log.Error("err=[查询需要修改的卡的信息失败,卡号为----->%s]", req.CardNo)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	if isDefault == "1" {
		reply.ResultCode = ss_err.ERR_CURRENT_CARD_IS_DEFAULT
		return nil
	}

	// 找到默认卡
	cardNo, _ := dao.CardDaoInstance.QueryIsDefaultCardNo(posAccountNo, req.MoneyType, "1")
	if cardNo != "" {

		// 设置默认卡为非默认
		if errStr := dao.CardDaoInstance.UpdateIsDefault(tx, constants.NO_DEFAULT_CARD, cardNo); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = ss_err.ERR_MODIFY_DEFAULT_CARD_FAILD
			return nil
		}
	}

	//设置默认卡为默认
	if errStr := dao.CardDaoInstance.UpdateIsDefault(tx, constants.IS_DEFAULT_CARD, req.CardNo); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = ss_err.ERR_MODIFY_DEFAULT_CARD_FAILD
		return nil
	}

	// 记录日志
	descript := "修改默认卡"
	if errStr := dao.LogCardDaoInstance.InsertLogCard(tx, num, "-", req.AccountUid, channelNo, channelType, descript, vaType); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = ss_err.ERR_SAVE_CARD_FAILD
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (r *Auth) DeleteBindCard(ctx context.Context, req *go_micro_srv_auth.DeleteBindCardRequest, reply *go_micro_srv_auth.DeleteBindCardReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// 判断解绑的卡是否是默认卡
	isDefault, channelNo := dao.CardDaoInstance.QueryIsDefaultFromCardNo(req.CardNo)
	if isDefault == "1" {
		reply.ResultCode = ss_err.ERR_CURRENT_CARD_IS_DEFAULT
		return nil
	}

	vaType := 0
	switch req.MoneyType {
	case "usd":
		vaType = constants.VaType_USD_DEBIT
	case "khr":
		vaType = constants.VaType_KHR_DEBIT
	default:
		return nil
	}

	// 判断是服务商,销售员,用户
	var channelType string
	switch req.AccountType {
	case constants.AccountType_SERVICER: //服务商
		channelType = "2"
		//posAccountNo = req.AccountUid
	case constants.AccountType_POS: // 收银员
		channelType = "2"
	case constants.AccountType_USER: // 用户
		channelType = "1"
	}

	// 修改卡状态
	if errStr := dao.CardDaoInstance.UpdateIsDelete(tx, req.CardNo, "1"); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = ss_err.ERR_DELETE_CARD_DEFAULT
		return nil
	}
	// 记录日志
	descript := "解除绑定卡"
	if errStr := dao.LogCardDaoInstance.InsertLogCard(tx, req.CarNum, "-", req.AccountUid, channelNo, channelType, descript, vaType); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = ss_err.ERR_DELETE_CARD_DEFAULT
		return nil
	}
	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (r *Auth) MyData(ctx context.Context, req *go_micro_srv_auth.MyDataRequest, reply *go_micro_srv_auth.MyDataReply) error {
	servicerNo := ""
	switch req.AccountType {
	case constants.AccountType_SERVICER: //服务商
		servicerNo = req.OpAccNo
	case constants.AccountType_POS: // 收银员
		servicerNo = dao.CashierDaoInstance.GetServicerNoFromOpAccNo(req.OpAccNo)
	default:
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	openIdxT, contactPersonT, contactPhoneT, contactAddrT, addrT, incomeAuthorizationT, outgoAuthorizationT, createTimeT := dao.
		ServicerDaoInstance.GetServicerFromServiceNo(servicerNo)

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.OpenIdx = openIdxT
	reply.ContactPerson = contactPersonT
	reply.ContactPhone = contactPhoneT
	reply.ContactAddr = contactAddrT
	reply.Addr = addrT
	reply.IncomeAuthorization = incomeAuthorizationT
	reply.OutgoAuthorization = outgoAuthorizationT
	reply.CreateTime = createTimeT
	return nil
}

func (r *Auth) PerfectingInfo(ctx context.Context, req *go_micro_srv_auth.PerfectingInfoRequest, reply *go_micro_srv_auth.PerfectingInfoReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	if req.Nickname != "" {
		// 修改昵称
		if err := ss_sql.Exec(dbHandler, `update account set nickname = $2 where uid = $1 and is_delete = 0`, req.AccountNo, req.Nickname); err != nil {
			ss_log.Error("err=[%v]", err.Error())
			reply.ResultCode = ss_err.ERR_MODIFY_NICKNAME_FAILD
			return nil
		}
	}

	if req.ImageStr != "" {
		// 上传图片
		reply2, _ := i.CustHandlerInst.Client.UploadImage(context.TODO(), &go_micro_srv_cust.UploadImageRequest{
			ImageStr:   req.ImageStr,
			AccountUid: req.AccountNo,
			Type:       2,
		}, global.RequestTimeoutOptions)

		if reply2.ResultCode != ss_err.ERR_SUCCESS {
			ss_log.Error("err=[调用图片上传的rpc失败----->%s]", reply2.ResultCode)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		//修改头像
		if err := ss_sql.Exec(dbHandler, `update account set head_portrait_img_no = $2 where uid = $1 and is_delete = 0`, req.AccountNo, reply2.ImageId); err != nil {
			ss_log.Error("err=[%v]", err.Error())
			reply.ResultCode = ss_err.ERR_MODIFY_HEAD_PORTRAIT_FAILD
			return nil
		}
		reply.HeadPortraitImgNo = reply2.ImageId
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*Auth) GetAccountCollect(ctx context.Context, req *go_micro_srv_auth.GetAccountCollectRequest, resp *go_micro_srv_auth.GetAccountCollectReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//只显示最近一个月的
	nowTime := time.Now()
	endTime := nowTime.Format("2006-01-02") //获取当前 年月日

	getTime := nowTime.AddDate(0, -1, 0)      //年，月，日   获取一个月前的时间
	startTime := getTime.Format("2006-01-02") //获取的时间的格式

	var datas []*go_micro_srv_auth.AccountCollectData
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ac.account_no", Val: req.AccountNo, EqType: "="},
		{Key: "to_char(ac.modify_time,'yyyy-MM-dd')", Val: startTime, EqType: ">="},
		{Key: "to_char(ac.modify_time,'yyyy-MM-dd')", Val: endTime, EqType: "<="},
		{Key: "ac.is_delete", Val: "0", EqType: "="},
	})

	cntwhere := whereModel.WhereStr
	cntargs := whereModel.Args

	var total sql.NullString
	sqlCnt := "select count(1) " +
		" from account_collect ac " +
		" LEFT JOIN account acc ON acc.uid = ac.collect_account_no " + cntwhere
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, cntargs...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by ac.modify_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	where := whereModel.WhereStr
	args := whereModel.Args
	sqlStr := " SELECT ac.account_collect_no, ac.collect_phone, ac.modify_time, acc.account, acc.nickname, acc.head_portrait_img_no, acc.country_code " +
		" FROM account_collect ac " +
		" LEFT JOIN account acc ON acc.uid = ac.collect_account_no " + where

	Rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer Rows.Close()

	if err == nil {
		for Rows.Next() {
			var data go_micro_srv_auth.AccountCollectData
			var collectPhone, account, nickname, headPortraitImgNo, countryCode sql.NullString
			err = Rows.Scan(
				&data.CollectNo,
				&collectPhone,
				&data.ModifyTime,
				&account,
				&nickname,
				&headPortraitImgNo,
				&countryCode,
			)
			data.CollectPhone = collectPhone.String
			data.Account = account.String
			data.Nickname = nickname.String
			data.HeadPortraitImgNo = headPortraitImgNo.String
			data.CountryCode = countryCode.String

			datas = append(datas, &data)
		}
	} else {
		ss_log.Error("err=[%v]", err)
		resp.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	resp.Datas = datas
	resp.ResultCode = ss_err.ERR_SUCCESS
	resp.Total = strext.ToInt32(total.String)
	return nil
}

func (*Auth) AddAccountCollect(ctx context.Context, req *go_micro_srv_auth.AddAccountCollectRequest, resp *go_micro_srv_auth.AddAccountCollectReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//获取账号
	toAccountNo := dao.AccDaoInstance.GetAccNoFromPhone(req.ToPhone, req.CountryCode)
	if toAccountNo == "" {
		toAccountNo = "00000000-0000-0000-0000-000000000000"
		resp.ResultCode = ss_err.ERR_SUCCESS
		return nil
	}

	//检查该账号是否有该条最近转账人信息
	count, err := dao.AccDaoInstance.CheckAccountCollect(req.AccountNo, req.ToPhone)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	countSelf, err := dao.AccDaoInstance.GetAccountCollectCount(req.AccountNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		countSelf = 0
	}

	//如没有则添加
	if count == 0 {
		_, maxCnt, _ := cache.ApiDaoInstance.GetGlobalParam("friends_count_limit")
		if maxCnt == "" {
			maxCnt = "-1" // 不限制
		}
		if countSelf >= strext.ToInt(maxCnt) && maxCnt != "-1" {
			// 限制个数
			sqlSql := "insert into account_collect(account_no, collect_account_no, collect_phone, is_delete, create_time, modify_time, account_collect_no) " +
				"values($1,$2,$3,$4,current_timestamp,current_timestamp,$5)"
			err := ss_sql.Exec(dbHandler, sqlSql, req.AccountNo, toAccountNo, req.ToPhone, "0", strext.GetDailyId())
			if err != nil {
				ss_log.Error("err=[%v]", err)
				resp.ResultCode = ss_err.ERR_SYS_DB_ADD
				return nil
			}

			err = dao.AccDaoInstance.DelLastFriend(req.AccountNo)
			if err != nil {
				ss_log.Error("err=[%v]", err)
				resp.ResultCode = ss_err.ERR_SYS_DB_ADD
				return nil
			}
		} else {
			sqlSql := "insert into account_collect(account_no, collect_account_no, collect_phone, is_delete, create_time, modify_time, account_collect_no) " +
				"values($1,$2,$3,$4,current_timestamp,current_timestamp,$5)"

			err := ss_sql.Exec(dbHandler, sqlSql, req.AccountNo, toAccountNo, req.ToPhone, "0", strext.GetDailyId())
			if err != nil {
				ss_log.Error("err=[%v]", err)
				resp.ResultCode = ss_err.ERR_SYS_DB_ADD
				return nil
			}
		}
	} else { //如果有则更新最后修改时间
		sqlUpdate := "update account_collect set collect_account_no = $3, modify_time = current_timestamp where account_no = $1 and collect_phone = $2 and is_delete = '0' "
		err := ss_sql.Exec(dbHandler, sqlUpdate, req.AccountNo, req.ToPhone, toAccountNo)

		if err != nil {
			ss_log.Error("err=[%v]", err)
			resp.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}
	}
	resp.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//校验支付密码
func (*Auth) CheckPayPWD(ctx context.Context, req *go_micro_srv_auth.CheckPayPWDRequest, reply *go_micro_srv_auth.CheckPayPWDReply) error {
	ss_log.Info("req=[%+v]", req)
	if req.AccountUid == "" || req.AccountType == "" || req.IdenNo == "" || req.NonStr == "" || req.Password == "" {
		ss_log.Error("必要参数为空，req[%+v]", req)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	switch req.AccountType {
	case constants.AccountType_USER: //用户
		if resultCode, payPasswordErrTips := paymentPwdErrLimit2(req.AccountUid, req.AccountType, req.Password, req.NonStr, req.IdenNo); resultCode != ss_err.ERR_SUCCESS {
			if resultCode == ss_err.ERR_Payment_Pwd_Count_Limit {
				// 修改交易权限
				if err := dao.CustDaoInstance.UpdateTradingAuthority(req.IdenNo, constants.TradingAuthorityForbid); err != nil {
					ss_log.Error("Exchange UpdateTradingAuthority err: %s", err.Error())
					reply.ResultCode = ss_err.ERR_PARAM
					return nil
				}
			}

			reply.ResultCode = resultCode
			reply.ErrTips = payPasswordErrTips //提示还可以输入几次错误支付密码
			return nil
		}
	case constants.AccountType_POS: //店员
		fallthrough
	case constants.AccountType_SERVICER: //服务商
		if resultCode, payPasswordErrTips := paymentPwdErrLimit2(req.AccountUid, req.AccountType, req.Password, req.NonStr, req.IdenNo); resultCode != ss_err.ERR_SUCCESS {
			reply.ResultCode = resultCode
			reply.ErrTips = payPasswordErrTips //提示还可以输入几次错误支付密码
			return nil
		}
	default: //以下是通用的，只是单纯的验证支付密码的
		if errStr := common.CheckPayPWD(req.AccountType, req.IdenNo, req.NonStr, req.Password); errStr != ss_err.ERR_SUCCESS {
			//如果是输入密码错误, 增加错误密码日志（关键操作日志）
			if errStr == ss_err.ERR_DB_PWD {
				if err := dao.LogDaoInstance.InsertAccountLog("验证支付密码错误", req.AccountUid, req.AccountType, constants.MODIFYACCOUNTLOGTYPE); err != nil {
					ss_log.Error("err=[%v],missing key=[%v]", err, "插入验证支付密码错误的日志失败")
				}
			}
			reply.ResultCode = errStr
			return nil
		}
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//验证支付密码(第二版本,有密码限制版本)
//payPasswordErrTips : 提示还可以输入几次错误支付密码
func paymentPwdErrLimit2(accountUid, accountType, password, nonStr, idenNo string) (resultCode string, payPasswordErrTips string) {
	ss_log.Error("判断支付密码是否被限制,并验证支付密码。开始参数 accountType[%v], password[%v], nonStr[%v], idenNo[%v]",
		accountType, password, nonStr, idenNo)

	// 判断支付密码是否被限制
	limitKey := cache.GetPayPwdErrCountKey(cache.PrePaymentPwdErrCountKey, accountType, idenNo)
	result, errRedis := cache.RedisClient.Get(limitKey).Result()
	if errRedis != nil && errRedis.Error() != constants.RedisNilValue {
		ss_log.Error("paymentPwdErrLimit 判断交易密码错误是否超过规定次数,查询redis失败,err: %s", errRedis.Error())
		return ss_err.ERR_Payment_Pwd_Count_Limit, ""
	}

	count := strext.ToInt(result) //当前密码错误次数
	bool1, limitCount := payPwdisLimit(count)
	if bool1 {
		ss_log.Error("paymentPwdErrLimit 交易密码受限制,当前错误次数为 %d,限制的次数为: %d", count, limitCount)
		//reply.ResultCode = ss_err.ERR_Payment_Pwd_Count_Limit
		return ss_err.ERR_Payment_Pwd_Count_Limit, ""
	}

	//验证支付密码
	resultCode = common.CheckPayPWD(accountType, idenNo, nonStr, password)

	if resultCode == ss_err.ERR_DB_PWD {
		if err := dao.LogDaoInstance.InsertAccountLog("验证支付密码错误", accountUid, accountType, constants.MODIFYACCOUNTLOGTYPE); err != nil {
			ss_log.Error("err=[%v],missing key=[%v]", err, "插入验证支付密码错误的日志失败")
		}

		count = count + 1 //加上本次错误
		// 设置错误次数
		cache.RedisClient.Set(limitKey, count, 0)

		ss_log.Error(" count[%v],limitCount[%v]", count, limitCount)
		if count >= limitCount {
			ss_log.Error("paymentPwdErrLimit 支付密码连续错误次数超出限制，count[%v],limitCount[%v]", count, limitCount)
			return ss_err.ERR_Payment_Pwd_Count_Limit, ""
		}

		return ss_err.ERR_PAY_FAILED_COUNT, strext.ToStringNoPoint(limitCount - count)
	}

	// 密码正确清除密码错误次数限制
	delErr := cache.RedisClient.Del(limitKey).Err()
	if delErr != nil {
		ss_log.Error("Del err: %s", delErr.Error())
	}

	return ss_err.ERR_SUCCESS, ""
}

func payPwdisLimit(count int) (bool, int) {
	var limitCount int
	// 判断次数
	value := dao.GlobalParamDaoInstance.QeuryParamValue(constants.GlobalParamKeyPaymentPwdErrCount)
	if value == "" {
		limitCount = constants.ErrPwdLimitDefaultCount // 默认5条
	} else {
		limitCount = strext.ToInt(value)
	}
	return count >= limitCount, limitCount
}

func (r *Auth) GetPosRemain(ctx context.Context, req *go_micro_srv_auth.GetPosRemainRequest, reply *go_micro_srv_auth.GetPosRemainReply) error {

	// 获取 serAccNo
	var serAccNo string
	switch req.AccountType {
	case constants.AccountType_POS: // 销售员
		serAccNo = dao.CashierDaoInstance.GetSrvAccNoFromCaNo(req.AccountUid)
		if serAccNo == "" {
			ss_log.Error("err=[服务商账号id错误,销售员账号id为----->%s]", req.AccountUid)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

	case constants.AccountType_SERVICER: // 服务商
		serAccNo = req.AccountUid
	}

	// 获取虚账额度
	// usd 实时额度
	_, usdRealBalance := dao.VaccountDaoInst.GetBalanceFromAccNo(serAccNo, constants.VaType_QUOTA_USD_REAL)
	// usd 授权额度
	_, usdAuthBalance := dao.VaccountDaoInst.GetBalanceFromAccNo(serAccNo, constants.VaType_QUOTA_USD)
	// khr 实时额度
	_, khrRealBalance := dao.VaccountDaoInst.GetBalanceFromAccNo(serAccNo, constants.VaType_QUOTA_KHR_REAL)
	// khr 授权额度
	_, khrAuthBalance := dao.VaccountDaoInst.GetBalanceFromAccNo(serAccNo, constants.VaType_QUOTA_KHR)
	ss_log.Info("美金授权额度为--->%s ,使用额度为--->%s, 瑞尔授权额度为--->%s, 使用额度为--->%s", usdAuthBalance, usdRealBalance, khrAuthBalance, khrRealBalance)
	var useKhr, UseUsd string

	// 需要把授权额度加上去了
	//UseUsd = ss_count.Sub(usdAuthBalance, usdRealBalance).String()
	//useKhr = ss_count.Sub(khrAuthBalance, khrRealBalance).String()

	// 暂时不需要把授权额度加上去了
	UseUsd = ss_count.Sub("0", usdRealBalance).String()
	useKhr = ss_count.Sub("0", khrRealBalance).String()

	data := &go_micro_srv_auth.GetPosRemainData{
		UseKhr:  useKhr,
		AuthKhr: khrAuthBalance,
		UseUsd:  UseUsd,
		AuthUsd: usdAuthBalance,
	}
	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//查询用户的余额和银行卡(渠道logo图片)
func (*Auth) CustPayment(ctx context.Context, req *go_micro_srv_auth.CustPaymentRequest, reply *go_micro_srv_auth.CustPaymentReply) error {
	khr, usd := dao.AccDaoInstance.GetRemain(req.AccountUid)

	ss_log.Info("khr=[%v], usd=[%v]", khr, usd)
	cardDatas, err := dao.CardDaoInstance.GetCustPaymentCard(req.AccountUid)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v],查询银行卡出错", err)
	}

	ss_log.Info("cardDatas=[%v]", cardDatas)
	defPayNo := dao.AccDaoInstance.GetDefPayNo(req.CustNo)
	if defPayNo == "" {
		defPayNo = "usd_balance"
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DefPayNo = defPayNo
	reply.Data = &go_micro_srv_auth.CustPaymentData{
		KhrBalance: khr,
		UsdBalance: usd,
		CardDatas:  cardDatas,
	}
	return nil
}

//插入app打点日志
func (*Auth) InsertLogAppDot(ctx context.Context, req *go_micro_srv_auth.InsertLogAppDotRequest, reply *go_micro_srv_auth.InsertLogAppDotReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := dao.LogAppDotDaoInst.InsertLogAppDot(req.OpType, req.Uuid)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_ADD
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
 * 添加用户
 */
func (auth *Auth) AddCust(ctx context.Context, req *go_micro_srv_auth.AddCustRequest, reply *go_micro_srv_auth.AddCustReply) error {
	if req.Password == "" {
		_, defPass, _ := cache.ApiDaoInstance.GetGlobalParam("def_password")
		req.Password = defPass
	}
	//创建账号
	addAccReq := &go_micro_srv_auth.SaveAccountRequest{
		Nickname:    req.Nickname,
		Account:     req.Phone,
		UseStatus:   "1",
		Phone:       req.Phone,
		Password:    req.Password,
		AccountType: constants.AccountType_USER,
		CountryCode: req.CountryCode,
		UtmSource:   req.UtmSource,
	}
	addAccReply := &go_micro_srv_auth.SaveAccountReply{}

	err := AuthHandlerInst.SaveAccount(ctx, addAccReq, addAccReply)
	if err != nil || addAccReply.ResultCode != ss_err.ERR_SUCCESS {
		ss_log.Error("创建账号失败 err=[%v],ResultCode=[%v] ", err, addAccReply.ResultCode)
		reply.ResultCode = addAccReply.ResultCode
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	err2, custNo := dao.CustDaoInstance.AddCustTx(tx, addAccReply.Uid, req.Gender)
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		//删除刚创建的账号
		if err := dao.AccDaoInstance.DeleteAccount(addAccReply.Uid); err != nil {
			ss_log.Error("删除刚创建的账号失败accountUid[%v],err=[%v]", addAccReply.Uid, err)
		}

		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	errCode := dao.RelaAccIdenDaoInst.InsertRelaAccIdenTx(tx, addAccReply.Uid, custNo, constants.AccountType_USER)
	if errCode != ss_err.ERR_SUCCESS {
		//删除刚创建的账号
		if err := dao.AccDaoInstance.DeleteAccount(addAccReply.Uid); err != nil {
			ss_log.Error("删除刚创建的账号失败accountUid[%v],err=[%v]", addAccReply.Uid, err)
		}

		ss_log.Error("errCode=[%v]", errCode)
		reply.ResultCode = errCode
		return nil
	}
	ss_sql.Commit(tx)
	reply.AccUid = addAccReply.Uid
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
 * 添加用户当前设备指纹
 */
func (auth *Auth) AddAppFingerprint(ctx context.Context, req *go_micro_srv_auth.AddAppFingerprintRequest, reply *go_micro_srv_auth.AddAppFingerprintReply) error {
	//查询是否可录入
	if appFingerprintOn := dao.GlobalParamDaoInstance.QeuryParamValue("app_fingerprint_on"); appFingerprintOn != constants.AppFingerprintOn_True {
		ss_log.Error("指纹录入总开关未打开，不允许录入")
		reply.ResultCode = ss_err.ERR_AppFingerprintOn_False
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	if err := dao.AppFingerprintDaoInstance.SetUseStatusDisableByAccountTx(tx, req.AccountNo, req.DeviceUuid); err != nil {
		tx.Rollback()
		ss_log.Error("err[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_ADD
		return nil
	}

	signKey, err := dao.AppFingerprintDaoInstance.AddTx(tx, req.AccountNo, req.DeviceUuid)
	if err != nil {
		tx.Rollback()
		ss_log.Error("err[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_ADD
		return nil
	}

	tx.Commit()
	reply.SignKey = signKey
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
 * 关闭用户当前设备的指纹
 */
func (auth *Auth) CloseAppFingerprint(ctx context.Context, req *go_micro_srv_auth.CloseAppFingerprintRequest, reply *go_micro_srv_auth.CloseAppFingerprintReply) error {
	if err := dao.AppFingerprintDaoInstance.SetUseStatusDisableByAccount(req.AccountNo, req.DeviceUuid); err != nil {
		ss_log.Error("err[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
