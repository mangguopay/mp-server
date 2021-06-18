package handler

import (
	"context"
	"database/sql"
	"runtime/debug"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/admin-auth-srv/dao"
	"a.a/mp-server/common/constants"
	adminAuthProto "a.a/mp-server/common/proto/admin-auth"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	_ "github.com/lib/pq"
)

// 获取账户列表
func (a *AdminAuth) GetAdminRoleUrlList(ctx context.Context, req *adminAuthProto.GetAdminRoleUrlListRequest, reply *adminAuthProto.GetAdminRoleUrlListReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	datas := getAdminRoleUrlData(ss_sql.UUID, req.RoleNo)
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	return nil
}

func (a *AdminAuth) GetAdminRoleList(ctx context.Context, req *adminAuthProto.GetAdminRoleListRequest, reply *adminAuthProto.GetAdminRoleListReply) error {
	cnt := dao.AdminRoleDaoInst.GetAdminRoleCnt(req.Search, req.AccType)
	errCode, datas := dao.AdminRoleDaoInst.GetAdminRoleList(req.Search, req.AccType, req.PageSize, req.Page)
	reply.ResultCode = errCode
	if errCode == ss_err.ERR_SUCCESS {
		reply.DataList = datas
	}
	reply.Len = cnt
	return nil
}

func getAdminRoleUrlData(urlUid, roleNo string) []*adminAuthProto.AdminRoleUrlListData {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlstr := "SELECT u.url_uid,u.url_name,u.parent_uid,COALESCE(rru.role_uid, '00000000-0000-0000-0000-000000000000') as role_no FROM  admin_url u " +
		"LEFT JOIN admin_rela_role_url rru ON rru.url_uid = u.url_uid AND rru.role_uid = $1 " +
		"WHERE parent_uid = $2 "
	rows, stmt, err := ss_sql.Query(dbHandler, sqlstr, roleNo, urlUid)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datas := []*adminAuthProto.AdminRoleUrlListData{}
	if err == nil {
		for rows.Next() {
			data := adminAuthProto.AdminRoleUrlListData{}
			rows.Scan(&data.Id, &data.Name, &data.ParentUid, &data.RoleNo)
			data.Children = getAdminRoleUrlData(data.Id, roleNo)
			datas = append(datas, &data)
		}
	} else {
		ss_log.Error("err=[%v]", err)
	}
	return datas
}

// 获取账户列表
func (a *AdminAuth) GetAdminRoleInfo(ctx context.Context, req *adminAuthProto.GetAdminRoleInfoRequest, reply *adminAuthProto.GetAdminRoleInfoReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var data adminAuthProto.RoleData
	row, stmt, err := ss_sql.QueryRowN(dbHandler, "SELECT r.role_no,r.role_name,r.create_time,r.modify_time,"+
		"r.acc_type,r.master_acc,acc.account FROM admin_role r left join admin_account acc on acc.uid=r.master_acc WHERE r.role_no = $1 LIMIT 1",
		req.RoleNo)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		ss_log.Error("err=%v", err.Error())
		reply.ResultCode = ss_err.ERR_ACCOUNT_ROLE_NOT_EXISTS
		reply.Data = &data
		return nil
	}
	var masterAcc, masterAccount sql.NullString
	err = row.Scan(&data.RoleNo, &data.RoleName, &data.CreateTime, &data.ModifyTime, &data.AccType, &masterAcc, &masterAccount)
	data.MasterAcc = masterAcc.String
	data.MasterAccount = masterAccount.String
	if err != nil {
		ss_log.Error("err=%v", err.Error())
		reply.ResultCode = ss_err.ERR_ACCOUNT_ROLE_NOT_EXISTS
		reply.Data = &data
		return nil
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = &data
	return nil
}

// 获取账户列表
func (a *AdminAuth) GetAdminRole(ctx context.Context, req *adminAuthProto.GetAdminRoleRequest, reply *adminAuthProto.GetAdminRoleReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	row, stmt, err := ss_sql.QueryRowN(dbHandler, "SELECT r.role_no,r.role_name,r.create_time,r.modify_time,"+
		"r.master_acc,acc.account,r.top_agency FROM admin_role r left join admin_account acc on acc.uid=r.master_acc WHERE r.role_name = $1 LIMIT 1",
		req.RoleName)
	if stmt != nil {
		defer stmt.Close()
	}
	if nil != err {
		ss_log.Error("GetAdminRole|err=%v", err.Error())
		debug.PrintStack()
		reply.ResultCode = ss_err.ERR_ARGS
		return nil
	}
	var masterAcc, masterAccount sql.NullString
	err = row.Scan(&reply.RoleNo, &reply.RoleName, &reply.CreateTime, &reply.ModifyTime, &masterAcc, &masterAccount, &reply.TopAgency)
	reply.MasterAcc = masterAcc.String
	reply.MasterAccount = masterAccount.String
	if nil != err {
		ss_log.Error("GetAdminRole|err=%v", err.Error())
		debug.PrintStack()
		reply.ResultCode = ss_err.ERR_ARGS
		return nil
	}

	rowsUrl, stmt, err := ss_sql.Query(dbHandler, "select url_uid,url_name,url from admin_url where url_uid in (select url_uid from admin_rela_role_url where role_uid = $1)",
		reply.RoleNo)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rowsUrl.Close()
	if nil != err {
		ss_log.Error("GetRole|err=%v", err)
		debug.PrintStack()
		reply.ResultCode = ss_err.ERR_ACCOUNT_MENU_NOT_EXISTS
		return nil
	}

	rowsUrl2, stmt, err := ss_sql.Query(dbHandler, "select url_uid,url_name,url from admin_url")
	if stmt != nil {
		defer stmt.Close()
	}
	defer rowsUrl2.Close()

	if nil == err {
		reply.ResultCode = ss_err.ERR_SUCCESS
		for rowsUrl.Next() {
			data := adminAuthProto.GetRoleUrlData{}
			rowsUrl.Scan(&data.UrlUid, &data.Name, &data.Url)
			reply.UrlData = append(reply.UrlData, &data)
		}

		for rowsUrl2.Next() {
			data := adminAuthProto.GetRoleUrlData{}
			rowsUrl2.Scan(&data.UrlUid, &data.Name, &data.Url)
			reply.UrlData_2 = append(reply.UrlData_2, &data)
		}
	} else {
		debug.PrintStack()
		ss_log.Error("err=%v", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
	}
	return nil
}

// 删除账户
func (a *AdminAuth) DeleteAdminRole(ctx context.Context, req *adminAuthProto.DeleteAdminRoleRequest, reply *adminAuthProto.DeleteAdminRoleReply) error {
	errCode := dao.AdminRoleDaoInst.DeleteAdminRole(req.RoleNo)
	reply.ResultCode = errCode
	return nil
}

// 更新或者插入授权信息
func (a *AdminAuth) UpdateOrInsertAdminRoleAuth(ctx context.Context, req *adminAuthProto.UpdateOrInsertAdminRoleAuthRequest, reply *adminAuthProto.UpdateOrInsertAdminRoleAuthReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, err := dbHandler.Begin()
	if nil != err {
		ss_log.Error("UpdateOrInsertRoleAuth|begin|err=[%v]\n", err)
		debug.PrintStack()
		tx.Rollback()
		reply.ResultCode = ss_err.ERR_SYS_DB_INIT
		return nil
	}

	var oldUrls []string
	var delUrls []string
	var addUrls []string

	// 把原有的先全部查出来
	rowsOldUrl, stmt, errSel := ss_sql.Query(dbHandler, "select url_uid from admin_rela_role_url where role_uid = $1", req.RoleNo)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rowsOldUrl.Close()
	if nil != errSel {
		tx.Rollback()
		ss_log.Error("err=%v", errSel)
		debug.PrintStack()
		reply.ResultCode = ss_err.ERR_ACCOUNT_ROLE_NOT_EXISTS
		return nil
	}

	oldUrls = db.ToStringList(rowsOldUrl)

	urlUidList := req.Urls
	for _, v := range oldUrls {
		// 原先有，现在没有
		idx := getKey(v, urlUidList)
		if idx < 0 {
			delUrls = append(delUrls, v)
		} else {
			urlUidList = remove(idx, urlUidList)
		}
	}

	// 原先没有，现在有
	addUrls = urlUidList

	// 实际执行
	pDel, errDelete := dbHandler.Prepare("delete from admin_rela_role_url where role_uid = $1 and url_uid = $2")
	defer pDel.Close()

	if nil != errDelete {
		tx.Rollback()
		ss_log.Error("err=%v", errDelete)
		debug.PrintStack()
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}

	for _, v := range delUrls {
		_, errDelete := pDel.Exec(req.RoleNo, v)

		if nil != errDelete {
			tx.Rollback()
			ss_log.Error("err=%v", errDelete)
			debug.PrintStack()
			reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
			return nil
		}
	}

	pInsert, errInsert := dbHandler.Prepare("insert into admin_rela_role_url(rela_uid,url_uid,role_uid) VALUES ($1,$2,$3)")
	defer pInsert.Close()
	if nil != errInsert {
		tx.Rollback()
		ss_log.Error("err=%v", errInsert)
		debug.PrintStack()
		reply.ResultCode = ss_err.ERR_SYS_DB_ADD
		return nil
	}

	for _, v := range addUrls {
		_, errInsert := pInsert.Exec(strext.NewUUID(), v, req.RoleNo)

		if nil != errInsert {
			tx.Rollback()
			ss_log.Error("err=%v", errInsert)
			debug.PrintStack()
			reply.ResultCode = ss_err.ERR_SYS_DB_ADD
			return nil
		}
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 更新或者插入新主商户
func (a *AdminAuth) UpdateOrInsertAdminRole(ctx context.Context, req *adminAuthProto.UpdateOrInsertAdminRoleRequest, reply *adminAuthProto.UpdateOrInsertAdminRoleReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, err := dbHandler.Begin()

	if nil != err {
		ss_log.Error("err=[%v]", err)
		tx.Rollback()
		reply.ResultCode = ss_err.ERR_SYS_DB_INIT
		return nil
	}

	if req.MasterAcc == "" {
		req.MasterAcc = ss_sql.UUID
	}

	if "" == req.RoleNo {
		// 如果没有uid，那就尝试插入新记录
		roleNo := strext.NewUUID()
		errInsert := ss_sql.Exec(dbHandler, "INSERT INTO admin_role (role_no,role_name,create_time,acc_type,master_acc,is_delete) "+
			"VALUES ($1,$2,current_timestamp,$3,$4,'0')",
			roleNo, req.RoleName, req.AccType, req.MasterAcc)
		if nil != errInsert {
			ss_log.Error("err=[%v]", errInsert)
			tx.Rollback()
			reply.ResultCode = ss_err.ERR_SYS_DB_SAVE
			return nil
		}
	} else {
		// 有uid，则更新
		var cnt sql.NullString
		err := ss_sql.QueryRow(dbHandler, "SELECT count(1) FROM admin_role "+
			"WHERE role_no = $1 LIMIT 1", []*sql.NullString{&cnt}, req.RoleNo)
		if !cnt.Valid || strext.ToInt32(cnt.String) <= 0 || err != nil {
			tx.Rollback()
			ss_log.Error("err|UpdateOrInsertRole|no this roleUid=%v\n", req.RoleNo)
			reply.ResultCode = ss_err.ERR_ACCOUNT_ROLE_NOT_EXISTS
			return nil
		}

		errUpdate := ss_sql.Exec(dbHandler, "UPDATE admin_role SET role_name = $1, modify_time = current_timestamp, "+
			"acc_type=$2,master_acc=$4 WHERE role_no = $3", req.RoleName, req.AccType, req.RoleNo, req.MasterAcc)
		//if req.AccType == constants.AccountType_AUTH {
		//	err := ss_sql.ExecTx(tx, "insert  into  admin_rela_role_account_auth (uid,account_uid,role_uid) values ($1,$2,$3)", strext.NewUUID(), req.AccUid, req.RoleNo)
		//	if nil != err {
		//		ss_log.Error("UpdateOrInsertRole|err=[%v]", errUpdate)
		//		// do nothing...
		//	}
		//}
		if nil != errUpdate {
			tx.Rollback()
			ss_log.Error("UpdateOrInsertRole|err=[%v]", errUpdate)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (a *AdminAuth) AuthAdminRole(ctx context.Context, req *adminAuthProto.AuthAdminRoleRequest, reply *adminAuthProto.AuthAdminRoleReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, err := dbHandler.Begin()
	if nil != err {
		ss_log.Error("err=[%v]", err)
		tx.Rollback()
		reply.ResultCode = ss_err.ERR_SYS_DB_INIT
		return nil
	}
	sqlU := "update admin_role set def_type=$1 where role_no=$2"
	err = ss_sql.ExecTx(tx, sqlU, req.DefType, req.RoleNo)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		tx.Rollback()
		reply.ResultCode = ss_err.ERR_SYS_DB_INIT
		return nil
	}
	//var topAgency sql.NullString
	//sqlS := "select top_agency from admin_role where role_no=$1 limit 1"
	//err = ss_sql.QueryRowTx(tx, sqlS, []*sql.NullString{&topAgency}, req.RoleNo)
	//if nil != err {
	//	ss_log.Error("err=[%v]", err)
	//	tx.Rollback()
	//	reply.ResultCode = ss_err.ERR_SYS_DB_INIT
	//	return nil
	//}
	//if req.DefType == "1" {
	//	sqlStr := "update admin_role set def_type='0' where acc_type=$1 and role_no <> $2 and top_agency=$3"
	//	err := ss_sql.ExecTx(tx, sqlStr, req.AccType, req.RoleNo, topAgency.String)
	//	if nil != err {
	//		ss_log.Error("err=[%v]", err)
	//		tx.Rollback()
	//		reply.ResultCode = ss_err.ERR_SYS_DB_INIT
	//		return nil
	//	}
	//}
	if req.DefType == "1" {
		sqlStr := "update admin_role set def_type='0' where acc_type=$1 and role_no <> $2 "
		err := ss_sql.ExecTx(tx, sqlStr, req.AccType, req.RoleNo)
		if nil != err {
			ss_log.Error("err=[%v]", err)
			tx.Rollback()
			reply.ResultCode = ss_err.ERR_SYS_DB_INIT
			return nil
		}
	}
	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
