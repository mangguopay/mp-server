package handler

import (
	"fmt"

	"a.a/cu/db"
	"a.a/cu/encrypt"
	"a.a/cu/jwt"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/admin-auth-srv/common"
	"a.a/mp-server/admin-auth-srv/dao"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_struct"

	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"a.a/mp-server/admin-auth-srv/util"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	adminAuthProto "a.a/mp-server/common/proto/admin-auth"

	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

/**
 * 获取账户
 */
func (*AdminAuth) GetAccount(ctx context.Context, req *adminAuthProto.GetAccountRequest, resp *adminAuthProto.GetAccountReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	row, stmt, err := ss_sql.QueryRowN(dbHandler, "SELECT ac.uid, ac.account, ac.use_status, ac.create_time, "+
		"ac.modify_time, ac.phone "+
		"FROM account ac "+
		"WHERE ac.uid = $1 LIMIT 1", req.AccountUid)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
		return nil
	}

	err = row.Scan(
		&resp.Uid,
		&resp.Account,
		&resp.UseStatus,
		&resp.CreateTime,
		&resp.ModifyTime,
		&resp.Phone,
	)
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
				data := adminAuthProto.RouteData{}
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

	//idenNo := dao.AccDaoInstance.GetIdenNoFromAcc(req.AccountUid, req.AccountType)

	//if idenNo == "" {
	//	ss_log.Info("获取账户角色信息失败,uid为-------->%s", req.AccountUid)
	//	resp.ResultCode = ss_err.ERR_PARAM
	//	return nil
	//}

	ss_log.Info("获取账户[%v]角色信息成功", req.AccountUid)
	resp.ResultCode = ss_err.ERR_SUCCESS
	resp.AccountType = req.AccountType
	return nil
}

/**
 * 获取账户
 */
func (*AdminAuth) GetAdminAccount(ctx context.Context, req *adminAuthProto.GetAdminAccountRequest, resp *adminAuthProto.GetAdminAccountReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	row, stmt, err := ss_sql.QueryRowN(dbHandler, "SELECT ac.uid, ac.account, ac.use_status, ac.create_time, "+
		"ac.modify_time, ac.phone "+
		"FROM admin_account ac "+
		"WHERE ac.uid = $1 LIMIT 1", req.AccountUid)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
		return nil
	}

	err = row.Scan(
		&resp.Uid,
		&resp.Account,
		&resp.UseStatus,
		&resp.CreateTime,
		&resp.ModifyTime,
		&resp.Phone,
	)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
		return nil
	}

	rowsRole, stmt, err2 := ss_sql.Query(dbHandler, "SELECT role_uid FROM admin_rela_account_role WHERE account_uid = $1", req.AccountUid)
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
		rowsUrl, stmt, errSel := ss_sql.Query(dbHandler, "SELECT url_uid,url_name,url,parent_uid,title,icon,component_name,component_path,redirect,idx,is_hidden FROM admin_url "+
			"WHERE url_uid in (SELECT url_uid FROM admin_rela_role_url WHERE role_uid in ('"+strings.Join(roleUids, "','")+"'))")
		if stmt != nil {
			defer stmt.Close()
		}
		defer rowsUrl.Close()
		if errSel == nil {
			for rowsUrl.Next() {
				data := adminAuthProto.RouteData{}
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

	ss_log.Info("获取账户[%v]角色信息成功", req.AccountUid)
	resp.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//修改用户使用状态（冻结）
func (*AdminAuth) ModifyUserStatus(ctx context.Context, req *adminAuthProto.ModifyUserStatusRequest, resp *adminAuthProto.ModifyUserStatusReply) error {
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

func (r *AdminAuth) Login(ctx context.Context, req *adminAuthProto.LoginRequest, resp *adminAuthProto.LoginReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// XXX 测试
	//ret, err := cache.RedisCli.Get(constants.DefPoolName, "verify_"+req.Verifyid)
	ret, err := cache.RedisClient.Get("verify_" + req.Verifyid).Result()

	if ret == "" || err != nil {
		ss_log.Error("----------->crypt%s", err.Error())
		resp.ResultCode = ss_err.ERR_ACCOUNT_SMS_CODE
		return nil
	}

	if strings.ToLower(strext.ToStringNoPoint(ret)) != strings.ToLower(req.Verifynum) {
		//--------------v1版本----------------------
		//_, err = cache.RedisCli.Del(constants.DefPoolName, "verify_"+req.Verifyid)

		err = cache.RedisClient.Del("verify_" + req.Verifyid).Err()
		resp.ResultCode = ss_err.ERR_ACCOUNT_LOGIN_CODE
		return nil
	}
	//--------------v1版本----------------------
	//_, err = cache.RedisCli.Del(constants.DefPoolName, "verify_"+req.Verifyid)

	err = cache.RedisClient.Del("verify_" + req.Verifyid).Err()
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	var accUid, pwdMD5 sql.NullString
	err2 := ss_sql.QueryRow(dbHandler, "SELECT uid,password FROM admin_account WHERE account=$1 LIMIT 1", []*sql.NullString{&accUid, &pwdMD5}, req.Account)
	if err2 != nil {
		ss_log.Error("login failed|acc=[%v]|password=[%v],err=[%v]", req.Account, pwdMD5.String, err2)
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
		err = ss_sql.Exec(dbHandler, "insert into admin_log_login(log_time,acc_no,ip,result,log_no) values (current_timestamp,$1,$2,$3,$4)",
			accUid.String, "", ss_err.ERR_ACCOUNT_NOT_EXISTS, strext.GetDailyId())
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
		resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
		return nil
	}

	getAcc := adminAuthProto.GetAdminAccountReply{}
	_ = r.GetAdminAccount(ctx, &adminAuthProto.GetAdminAccountRequest{
		AccountUid: accUid.String,
	}, &getAcc)
	if getAcc.UseStatus != "1" {
		err = ss_sql.Exec(dbHandler, "insert into admin_log_login(log_time,acc_no,ip,result,log_no) values (current_timestamp,$1,$2,$3,$4)",
			accUid.String, "", ss_err.ERR_ACCOUNT_NO_PERMISSION, strext.GetDailyId())
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
		resp.ResultCode = ss_err.ERR_ACCOUNT_NO_PERMISSION
		return nil
	}

	accountType := dao.AdminAccDaoInstance.GetAccountTypeFromAccNoAdminOrOp(accUid.String)
	if accountType == "" {
		err = ss_sql.Exec(dbHandler, "insert into admin_log_login (log_time,acc_no,ip,result,log_no) values (current_timestamp,$1,$2,$3,$4)",
			accUid.String, "", ss_err.ERR_ACCOUNT_NO_PERMISSION, strext.GetDailyId())
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
		resp.ResultCode = ss_err.ERR_ACCOUNT_NO_PERMISSION
		return nil
	}

	retMap := common.JwtStructToMapWebAdmin(ss_struct.JwtDataWebAdmin{
		Account:        req.Account,
		AccountUid:     accUid.String,
		IdenNo:         accUid.String,
		AccountType:    accountType,
		LoginAccountNo: accUid.String,
		JumpIdenNo:     "",
		JumpIdenType:   "",
		MasterAccNo:    "",
		IsMasterAcc:    "1",
	})

	k1, loginSignKey, err := cache.ApiDaoInstance.GetGlobalParam("login_sign_key")
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k1)
	}
	k2, loginAesKey, err := cache.ApiDaoInstance.GetGlobalParam("login_aes_key")
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k2)
	}
	jwt2 := jwt.GetNewEncryptedJWTToken(-1, retMap, loginAesKey, loginSignKey)
	err = ss_sql.Exec(dbHandler, "insert into admin_log_login(log_time,acc_no,ip,result,log_no) values (current_timestamp,$1,$2,$3,$4)",
		accUid.String, req.Ip, ss_err.ERR_SUCCESS, strext.GetDailyId())
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

/**
重置密码
*/
func (r *AdminAuth) ResetPw(ctx context.Context, req *adminAuthProto.ResetPwRequest, reply *adminAuthProto.ResetPwReply) error {
	// 检查是否存在
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	if req.Cli == constants.LOGIN_CLI_WEB {
		redisKey, defPass, err := cache.ApiDaoInstance.GetGlobalParam("def_password")
		if err != nil {
			ss_log.Error("redisKey=[%v],err=[%v]", redisKey, err)
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}
		req.NewPw = defPass

		k1, passwordSalt, err := cache.ApiDaoInstance.GetGlobalParam("password_salt")
		if err != nil {
			ss_log.Error("err=[%v],missing key=[%v]", err, k1)
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}
		var tmp sql.NullString
		pwdMD5 := encrypt.DoMd5Salted(req.LoginPw, passwordSalt)
		err = ss_sql.QueryRow(dbHandler, "select 1 from admin_account where uid = $1 and password=$2 limit 1", []*sql.NullString{&tmp}, req.LoginAccNo, pwdMD5)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_DB_PWD
			return nil
		}
		if tmp.String != "1" {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_DB_PWD
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
		reply.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXIST
		return nil
	}

	// 已存在
	if tmp.String != "" {
		redisKey, defPass, err := cache.ApiDaoInstance.GetGlobalParam("def_password")
		if err != nil {
			ss_log.Error("redisKey=[%v],err=[%v]", redisKey, err)
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
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
			reply.ResultCode = ss_err.ERR_DB_PWD
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
重置密码
*/
func (r *AdminAuth) ResetAdminPw(ctx context.Context, req *adminAuthProto.ResetAdminPwRequest, reply *adminAuthProto.ResetAdminPwReply) error {
	// 检查是否存在
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	redisKey, defPass, err := cache.ApiDaoInstance.GetGlobalParam("def_password")
	if err != nil {
		ss_log.Error("redisKey=[%v],err=[%v]", redisKey, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	k1, passwordSalt, err := cache.ApiDaoInstance.GetGlobalParam("password_salt")
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k1)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}
	var tmp sql.NullString
	pwdMD5 := encrypt.DoMd5Salted(req.LoginPw, passwordSalt)
	err = ss_sql.QueryRow(dbHandler, "select 1 from admin_account where uid = $1 and password=$2 limit 1", []*sql.NullString{&tmp}, req.LoginAccNo, pwdMD5)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_ACCOUNT_WRONG_PASSWORD
		return nil
	}
	if tmp.String != "1" {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
		return nil
	}

	var tmp2 sql.NullString
	err2 := ss_sql.QueryRow(dbHandler, "select 1 from admin_account where uid = $1 limit 1", []*sql.NullString{&tmp2}, req.Account)
	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		reply.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXIST
		return nil
	}

	// 已存在
	if tmp2.String != "" {
		pwdMD5 := encrypt.DoMd5Salted(defPass, passwordSalt)
		err = ss_sql.Exec(dbHandler, "update admin_account set password=$1 where uid=$2", pwdMD5, req.Account)
		if err != nil {
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}

		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	}

	reply.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
	return nil
}

func (*AdminAuth) ModifyPw(ctx context.Context, req *adminAuthProto.ModifyPwRequest, reply *adminAuthProto.ModifyPwReply) error {
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
			tmp, err := dao.AdminAccDaoInstance.GetAccountCnt(req.Account)
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
		err = ss_sql.QueryRow(dbHandler, `select 1 from admin_account where account = $1 and "password" = $2 and is_delete='0' limit 1`, []*sql.NullString{&tmp}, req.Account, oldPwdMD5)
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
			err = ss_sql.Exec(dbHandler, `update admin_account set "password"=$1 where account=$2 and is_delete='0'`, pwdMD5, req.Account)

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
		err = ss_sql.QueryRow(dbHandler, `select 1 from admin_account where uid = $1 and "password" = $2 limit 1`, []*sql.NullString{&tmp}, req.Account, oldPwdMD5)
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
			err = ss_sql.Exec(dbHandler, `update admin_account set "password"=$1 where uid=$2`, pwdMD5, req.Account)

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

func (*AdminAuth) ModifyAdminPw(ctx context.Context, req *adminAuthProto.ModifyAdminPwRequest, reply *adminAuthProto.ModifyAdminPwReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var tmp sql.NullString

	k1, passwordSalt, err := cache.ApiDaoInstance.GetGlobalParam("password_salt")
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k1)
	}
	oldPwdMD5 := encrypt.DoMd5Salted(req.OldPw, passwordSalt)
	err = ss_sql.QueryRow(dbHandler, `select 1 from admin_account where uid = $1 and "password" = $2 limit 1`, []*sql.NullString{&tmp}, req.AccountUid, oldPwdMD5)
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
		err = ss_sql.Exec(dbHandler, `update admin_account set "password"=$1 where uid=$2`, pwdMD5, req.AccountUid)

		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}

		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	}
	reply.ResultCode = ss_err.ERR_ACCOUNT_NOT_EXISTS
	return nil
}

func (*AdminAuth) UpdateOrInsertAccountAuth(ctx context.Context, req *adminAuthProto.UpdateOrInsertAccountAuthRequest, resp *adminAuthProto.UpdateOrInsertAccountAuthReply) error {
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
					errCode := dao.RelaAccIdenDaoInst.InsertRelaAccIden(tx, req.Uid, servicerNo, constants.AccountType_SERVICER)
					if errCode != ss_err.ERR_SUCCESS {
						ss_log.Info("errCode==[%v]", errCode)
						resp.ResultCode = errCode
						return nil
					}

					//初始化钱包
					dao.VaccountDaoInst.InitVaccountNo(req.Uid, constants.CURRENCY_USD, constants.VaType_QUOTA_USD)
					dao.VaccountDaoInst.InitVaccountNo(req.Uid, constants.CURRENCY_USD, constants.VaType_QUOTA_USD_REAL)
					dao.VaccountDaoInst.InitVaccountNo(req.Uid, constants.CURRENCY_KHR, constants.VaType_QUOTA_KHR)
					dao.VaccountDaoInst.InitVaccountNo(req.Uid, constants.CURRENCY_KHR, constants.VaType_QUOTA_KHR_REAL)

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

func (*AdminAuth) UpdateOrInsertAdminAccountAuth(ctx context.Context, req *adminAuthProto.UpdateOrInsertAdminAccountAuthRequest, resp *adminAuthProto.UpdateOrInsertAdminAccountAuthReply) error {
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
	errDelete := dao.AdminRelaAccountRoleDaoInst.DeleteAdminRelaAccountRoleTx(tx, req.Uid)
	if ss_err.ERR_SUCCESS != errDelete {
		ss_log.Error("err=%v", errDelete)
		resp.ResultCode = errDelete
		return nil
	}

	account := dao.AdminAccDaoInstance.GetAdminAccountByUid(req.Uid)
	if account == "" {
		ss_log.Error("获取账号失败，uid:[%v]", req.Uid)
		account = req.Uid
	}
	description := fmt.Sprintf("授权账号[%v]的角色为[", account)

	if req.Roles != "" {
		roles := strings.Split(req.Roles, ",")
		for _, role := range roles {
			accType, err := dao.AdminRelaAccountRoleDaoInst.GetAdminRoleAccTypeTx(tx, role)
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
				serNo := dao.AdminRelaAccIdenDaoInst.GetAdminIdenFromAcc(req.Uid, constants.AccountType_SERVICER)
				if serNo == "" { //如果没有关联关系,说明该账号未创建过服务商
					//创建服务商(servicer)
					servicerNo, err := dao.ServicerDaoInstance.InsertInitService(tx, req.Uid)
					if err != nil {
						ss_log.Error("添加账号角色为服务商,插入服务商表失败,err: %s", err.Error())
						resp.ResultCode = ss_err.ERR_SYS_DB_OP
						return nil
					}
					// 建立关联关系
					errCode := dao.AdminRelaAccIdenDaoInst.InsertAdminRelaAccIden(tx, req.Uid, servicerNo, constants.AccountType_SERVICER)
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
			errCode := dao.AdminRelaAccountRoleDaoInst.InsertAdminRelaAccountRoleTx(tx, req.Uid, role)
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
func (*AdminAuth) GetAccountList(ctx context.Context, req *adminAuthProto.GetAccountListRequest, resp *adminAuthProto.GetAccountListReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	condStr, _ := json.Marshal([]*model.WhereSqlCond{
		{Key: "acc.account", Val: req.Search, EqType: "like"},
		{Key: "acc.nickname", Val: req.Search, EqType: "like"},
	})

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		//{Key: "acc.use_status", Val: "1", EqType: "="},//这里是筛选只要使用状态为正常的用户
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
	datas := []*adminAuthProto.Account{}
	if err == nil {
		for rows.Next() {
			data := adminAuthProto.Account{}
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
 * 获取账户列表
 */
func (*AdminAuth) GetAdminAccountList(ctx context.Context, req *adminAuthProto.GetAdminAccountListRequest, resp *adminAuthProto.GetAdminAccountListReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	condStr, _ := json.Marshal([]*model.WhereSqlCond{
		{Key: "acc.account", Val: req.Search, EqType: "like"},
	})

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "acc.use_status", Val: "1", EqType: "="},
		{Key: "", Val: string(condStr), EqType: "or_group"},
		{Key: "acc.is_delete", Val: "0", EqType: "="},
	})
	where := whereModel.WhereStr
	args := whereModel.Args

	var total sql.NullString
	sqlCnt := "SELECT count(1) FROM admin_account acc " + where
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by acc.create_time desc")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
	where = whereModel.WhereStr
	args = whereModel.Args
	rows, stmt, err := ss_sql.Query(dbHandler, "SELECT acc.uid, acc.account, acc.use_status, acc.create_time, acc.modify_time "+
		" FROM admin_account acc  "+where, args...)
	if stmt != nil {
		defer stmt.Close()
		defer rows.Close()
	}
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}
	datas := []*adminAuthProto.Account{}
	if err == nil {
		for rows.Next() {
			data := adminAuthProto.Account{}
			err = rows.Scan(&data.Uid, &data.Account, &data.UseStatus, &data.CreateTime, &data.ModifyTime)
			if err != nil {
				ss_log.Error("err=[%v]", err)
				continue
			}

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
 * 获取账户列表2
 */
func (*AdminAuth) GetAccountList2(ctx context.Context, req *adminAuthProto.GetAccountListRequest2, resp *adminAuthProto.GetAccountListReply2) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	if req.StartTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.StartTime) {
			ss_log.Error("GetAccountList2 StartTime格式不正确,StartTime: %s", req.StartTime)
			resp.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}
	if req.EndTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.EndTime) {
			ss_log.Error("GetAccountList2 EndTime格式不正确,StartTime: %s", req.EndTime)
			resp.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	var datas []*adminAuthProto.Account2
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "c.is_delete", Val: "0", EqType: "="},
		{Key: "acc.is_delete", Val: "0", EqType: "="},
		{Key: "acc.phone", Val: req.QueryPhone, EqType: "like"},
		{Key: "acc.account", Val: req.Account, EqType: "like"},
		{Key: "acc.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "acc.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "acc.nickname", Val: req.QueryNickname, EqType: "like"},
		{Key: "acc.individual_auth_status", Val: req.AuthStatus, EqType: "="},
		{Key: "acc.use_status", Val: req.UseStatus, EqType: "="},
	})

	cntwhere := whereModel.WhereStr
	cntargs := whereModel.Args

	var total sql.NullString
	sqlCnt := "select count(1) " +
		" from cust c " +
		" LEFT JOIN account acc ON acc.uid = c.account_no " + cntwhere
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, cntargs...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	switch req.SortType {
	case "usd_up": //usd余额正向排序(为null的最上面，余额按升序排，余额相同再按创建时间反序。)
		ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by vacc.balance is null desc,vacc.balance ASC, acc.create_time desc`)
	case "usd_down": //usd余额反向排序(为null的最下面，余额按降序排，余额相同再按创建时间反序。)
		ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by vacc.balance is null asc,vacc.balance desc, acc.create_time desc`)
	case "khr_up": //khr余额正向排序
		ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by vacc2.balance is null desc,vacc2.balance ASC, acc.create_time desc`)
	case "khr_down": //khr余额反向排序
		ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by vacc2.balance is null asc,vacc2.balance desc, acc.create_time desc`)
	default:
		ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by acc.create_time desc`)
	}
	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	sqlStr := "SELECT acc.uid, acc.nickname, acc.use_status, acc.phone, acc.create_time" +
		", vacc.balance, vacc2.balance, acc.account, acc.individual_auth_status, acc.country_code" +
		", acc.is_actived  " +
		" FROM cust c " +
		" LEFT JOIN account acc ON acc.uid = c.account_no " +
		" LEFT JOIN vaccount vacc ON vacc.account_no = acc.uid and vacc.va_type = " + strext.ToStringNoPoint(constants.VaType_USD_DEBIT) +
		" LEFT JOIN vaccount vacc2 ON vacc2.account_no = acc.uid and vacc2.va_type = " + strext.ToStringNoPoint(constants.VaType_KHR_DEBIT) + whereModel.WhereStr

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("err=[%v]", err)
		resp.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	for rows.Next() {
		data := &adminAuthProto.Account2{}
		var phone, usdBalance, khrBalance, authStatus, countryCode, isActived sql.NullString
		err = rows.Scan(
			&data.Uid,
			&data.Nickname,
			&data.UseStatus,
			&phone,
			&data.CreateTime,
			&usdBalance,
			&khrBalance,
			&data.Account,
			&authStatus,
			&countryCode,
			&isActived,
		)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		data.UsdBalance = usdBalance.String
		if usdBalance.String == "" {
			data.UsdBalance = "0"
		}

		data.KhrBalance = khrBalance.String
		if khrBalance.String == "" {
			data.KhrBalance = "0"
		}

		data.CountryCode = countryCode.String
		data.AuthStatus = authStatus.String
		data.IsActived = isActived.String

		data.Phone = phone.String
		datas = append(datas, data)
	}

	resp.ResultCode = ss_err.ERR_SUCCESS
	resp.AccountData = datas
	resp.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * 添加账户
 */
func (*AdminAuth) SaveAccount(ctx context.Context, req *adminAuthProto.SaveAccountRequest, resp *adminAuthProto.SaveAccountReply) error {
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
	account := fmt.Sprintf("%s%s", common.PreCountryCode(req.CountryCode), req.Account)
	if req.AccountUid == "" {
		if req.Password == "" {
			resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_PASSWORD
			return nil
		}
		accountCount, err := dao.AdminAccDaoInstance.CheckAccount(tx, req.Account)
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
		req.AccountUid, err = dao.AdminAccDaoInstance.AddAccount(tx, req.Nickname, account, req.Password, req.UseStatus,
			req.MasterAcc, req.Phone, req.CountryCode, req.UtmSource)
		if nil != err {
			ss_log.Error("err=%v", err.Error())
			resp.ResultCode = ss_err.ERR_ACCOUNT_INIT_ACCOUNT_ERR
			return nil
		}
		//管理员和运营才添加关联关系
		if req.AccountType == constants.AccountType_ADMIN || req.AccountType == constants.AccountType_OPERATOR {
			errCode := dao.RelaAccIdenDaoInst.InsertRelaAccIden(tx, req.AccountUid, "00000000-0000-0000-0000-000000000000", req.AccountType)
			if errCode != ss_err.ERR_SUCCESS {
				ss_log.Info("errCode==[%v]", errCode)
				resp.ResultCode = errCode
				return nil
			}
			ss_log.Info("AccountUid == [%v],AccountType == [%v]", req.AccountUid, req.AccountType)

			errCode = dao.AdminAccDaoInstance.AuthAccountRetCode(tx, req.AccountType, req.AccountUid)
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
		accountCount, err := dao.AdminAccDaoInstance.CheckAccountUpdate(tx, req.AccountUid)
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
		oldCountryCode, oldPhone, err := dao.AdminAccDaoInstance.GetCountryCodePhoneByUidTx(tx, req.AccountUid)
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
		err = dao.AdminAccDaoInstance.UpdateAccount(tx, req.Nickname, account, req.Password, req.UseStatus,
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
 * 添加账户
 */
func (*AdminAuth) SaveAdminAccount(ctx context.Context, req *adminAuthProto.SaveAdminAccountRequest, resp *adminAuthProto.SaveAdminAccountReply) error {
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

	if req.AccountUid == "" {
		if req.Password == "" {
			resp.ResultCode = ss_err.ERR_ACCOUNT_NOT_PASSWORD
			return nil
		}
		accountCount, err := dao.AdminAccDaoInstance.CheckAccount(tx, req.Account)
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
		req.AccountUid, err = dao.AdminAccDaoInstance.AddAdminAccount(tx, req.Account, req.Password, req.UseStatus, req.Phone)
		if nil != err {
			ss_log.Error("err=%v", err.Error())
			resp.ResultCode = ss_err.ERR_ACCOUNT_INIT_ACCOUNT_ERR
			return nil
		}
		//管理员和运营才添加关联关系
		if req.AccountType == constants.AccountType_ADMIN || req.AccountType == constants.AccountType_OPERATOR {
			errCode := dao.AdminRelaAccIdenDaoInst.InsertAdminRelaAccIden(tx, req.AccountUid, "00000000-0000-0000-0000-000000000000", req.AccountType)
			if errCode != ss_err.ERR_SUCCESS {
				ss_log.Info("errCode==[%v]", errCode)
				resp.ResultCode = errCode
				return nil
			}
			ss_log.Info("AccountUid == [%v],AccountType == [%v]", req.AccountUid, req.AccountType)

			errCode = dao.AdminAccDaoInstance.AuthAccountRetCode(tx, req.AccountType, req.AccountUid)
			if errCode != ss_err.ERR_SUCCESS {
				ss_log.Info("errCode==[%v]", errCode)
				resp.ResultCode = errCode
				return nil
			}
		}
	} else {
		accountCount, err := dao.AdminAccDaoInstance.CheckAccountUpdate(tx, req.AccountUid)
		if nil != err {
			ss_log.Error("SaveAccount|CheckAccount|err=[%v]", err.Error())
			resp.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}
		if accountCount == 0 {
			resp.ResultCode = ss_err.ERR_MSG_ACCOUNT_NOT_EXISTS
			return nil
		}
		// 更新
		err = dao.AdminAccDaoInstance.UpdateAdminAccount(tx, req.Account, req.Password, req.UseStatus,
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

func (*AdminAuth) CheckAccount(ctx context.Context, req *adminAuthProto.CheckAccountRequest, resp *adminAuthProto.CheckAccountReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	data := &adminAuthProto.CheckAccountData{}

	cnt, err := dao.AccDaoInstance.GetAccountCnt(req.Account)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}
	data.Count = strext.ToStringNoPoint(cnt)
	resp.ResultCode = ss_err.ERR_SUCCESS
	resp.Data = data
	return nil
}

func (*AdminAuth) CheckAdminAccount(ctx context.Context, req *adminAuthProto.CheckAdminAccountRequest, resp *adminAuthProto.CheckAdminAccountReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	data := &adminAuthProto.CheckAdminAccountData{}

	cnt, err := dao.AdminAccDaoInstance.GetAccountCnt(req.Account)
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
func (*AdminAuth) DeleteAccount(ctx context.Context, req *adminAuthProto.DeleteAccountRequest, resp *adminAuthProto.DeleteAccountReply) error {
	dao.AccDaoInstance.DeleteAccountList(req.AccountUids)
	// 不管是否成功
	resp.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
 * 获取账户
 */
func (*AdminAuth) GetAccountByNickname(ctx context.Context, req *adminAuthProto.GetAccountByNicknameRequest, resp *adminAuthProto.GetAccountByNicknameReply) error {
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
func (*AdminAuth) GetRoleFromAcc(ctx context.Context, req *adminAuthProto.GetRoleFromAccRequest, reply *adminAuthProto.GetRoleFromAccReply) error {
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

	ds := []*adminAuthProto.RoleSimpleData{}
	for rowsRole.Next() {
		d := adminAuthProto.RoleSimpleData{}
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

/**
 * 获取账户列表
 */
func (*AdminAuth) GetAdminRoleFromAcc(ctx context.Context, req *adminAuthProto.GetAdminRoleFromAccRequest, reply *adminAuthProto.GetAdminRoleFromAccReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	rowsRole, stmt, err2 := ss_sql.Query(dbHandler, "SELECT r.role_no,r.role_name FROM admin_rela_account_role rela "+
		" left join admin_role r on r.role_no=rela.role_uid WHERE account_uid = $1", req.AccNo)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rowsRole.Close()
	if nil != err2 {
		ss_log.Error("err|2=%v", err2)
		reply.ResultCode = ss_err.ERR_SYS_IO_ERR
		return nil
	}

	ds := []*adminAuthProto.RoleSimpleData{}
	for rowsRole.Next() {
		d := adminAuthProto.RoleSimpleData{}
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

func (*AdminAuth) GetLogLoginList(ctx context.Context, req *adminAuthProto.GetLogLoginListRequest, reply *adminAuthProto.GetLogLoginListReply) error {
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
	var datas []*adminAuthProto.LogLogin
	if err == nil {
		for rows.Next() {
			var data adminAuthProto.LogLogin
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

/**
 * 获取未激活账户列表
 */
func (*AdminAuth) GetUnActivedAccounts(ctx context.Context, req *adminAuthProto.GetUnActivedAccountsRequest, resp *adminAuthProto.GetUnActivedAccountsReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	if req.StartTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.StartTime) {
			ss_log.Error("GetAccountList2 StartTime格式不正确,StartTime: %s", req.StartTime)
			resp.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}
	if req.EndTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.EndTime) {
			ss_log.Error("GetAccountList2 EndTime格式不正确,StartTime: %s", req.EndTime)
			resp.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "acc.is_delete", Val: "0", EqType: "="},
		{Key: "acc.is_actived", Val: "0", EqType: "="}, //未激活的
		{Key: "acc.account", Val: req.Account, EqType: "like"},
		{Key: "acc.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "acc.create_time", Val: req.EndTime, EqType: "<="},
	})

	cntwhere := whereModel.WhereStr
	cntargs := whereModel.Args

	var total sql.NullString
	sqlCnt := "select count(1) " +
		" from account acc " + cntwhere
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, cntargs...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by acc.create_time desc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	where := whereModel.WhereStr
	args := whereModel.Args
	sqlStr := "SELECT acc.uid, acc.phone, acc.account, acc.country_code, acc.create_time" +
		" FROM account acc " + where

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	var datas []*adminAuthProto.Account3
	if err == nil {
		for rows.Next() {
			var data adminAuthProto.Account3
			var phone, account, countryCode, createTime sql.NullString
			err = rows.Scan(
				&data.Uid,
				&phone,
				&account,
				&countryCode,
				&createTime,
			)
			if err != nil {
				ss_log.Error("err=[%v]", err)
				continue
			}
			data.Phone = phone.String
			data.Account = account.String
			data.CountryCode = countryCode.String
			data.CreateTime = createTime.String

			datas = append(datas, &data)
		}
	} else {
		ss_log.Error("err=[%v]", err)
		resp.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	if datas == nil {
		resp.ResultCode = ss_err.ERR_CUST_NOT_EXISTS
		resp.Total = strext.ToInt32("0")
		return nil
	}

	resp.ResultCode = ss_err.ERR_SUCCESS
	resp.AccountData = datas
	resp.Total = strext.ToInt32(total.String)
	return nil
}
