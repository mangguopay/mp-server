package handler

import (
	"context"
	"database/sql"
	_ "database/sql"
	"runtime/debug"
	"strings"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/auth-srv/dao"
	"a.a/mp-server/common/constants"
	auth "a.a/mp-server/common/proto/auth"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	_ "github.com/lib/pq"
)

/**
 * 获取账户列表
 */
func (*Auth) GetRoleUrlList(ctx context.Context, req *auth.GetRoleUrlListRequest, reply *auth.GetRoleUrlListReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	masterAcc := dao.RoleDaoInst.GetRoleMasterAcc(req.RoleNo)
	if masterAcc == "" {
		reply.ResultCode = ss_err.ERR_ACCOUNT_ROLE_NOT_EXISTS
		return nil
	}

	datas := []*auth.RoleUrlListData{}
	if masterAcc != "" && masterAcc != ss_sql.UUID {
		rowsRole, stmt, err2 := ss_sql.Query(dbHandler, "SELECT role_uid FROM rela_account_role WHERE account_uid = $1 and is_delete='0' ", masterAcc)
		if stmt != nil {
			defer stmt.Close()
		}
		defer rowsRole.Close()
		roleUids := db.ToStringList(rowsRole)
		ss_log.Info("账号[%v]关联的角色ids[%v]", masterAcc, roleUids)
		if nil != err2 {
			ss_log.Error("err|2=%v", err2)
			reply.ResultCode = ss_err.ERR_ACCOUNT_MENU_NOT_EXISTS
			return nil
		}

		urids := []string{}

		// 读取菜单列表
		if len(roleUids) > 0 {
			rowsUrl, stmt, errSel := ss_sql.Query(dbHandler, "SELECT url_uid from url WHERE url_uid in (SELECT url_uid FROM rela_role_url WHERE role_uid in ('"+strings.Join(roleUids, "','")+"'))")
			if stmt != nil {
				defer stmt.Close()
			}
			defer rowsUrl.Close()
			if errSel == nil {
				for rowsUrl.Next() {
					var urlUid string
					err := rowsUrl.Scan(&urlUid)
					if err != nil {
						ss_log.Error("err=[%v]", err)
						continue
					}
					urids = append(urids, urlUid)
				}
			}

			if nil != errSel {
				ss_log.Error("err|3=%v", errSel)
				reply.ResultCode = ss_err.ERR_ACCOUNT_MENU_NOT_EXISTS
				return nil
			}
		}
		datas = getRoleUrlData2(ss_sql.UUID, req.RoleNo, strings.Join(urids, "','"))
	} else {
		datas = getRoleUrlData(ss_sql.UUID, req.RoleNo)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	reply.MasterAcc = masterAcc
	return nil
}

func (*Auth) GetRoleList(ctx context.Context, req *auth.GetRoleListRequest, reply *auth.GetRoleListReply) error {
	errCode, datas := dao.RoleDaoInst.GetRoleList(req.Search, req.AccType, req.MasterAcc, req.PageSize, req.Page)
	reply.ResultCode = errCode
	if errCode == ss_err.ERR_SUCCESS {
		reply.DataList = datas
	}
	cnt := dao.RoleDaoInst.GetRoleCnt(req.Search, req.AccType, req.MasterAcc)
	reply.Len = cnt
	return nil
}

func getRoleUrlData(urlUid, roleNo string) []*auth.RoleUrlListData {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlstr := "SELECT u.url_uid,u.url_name,u.parent_uid,COALESCE(rru.role_uid, '00000000-0000-0000-0000-000000000000') as role_no FROM  url u " +
		"LEFT JOIN rela_role_url rru ON rru.url_uid = u.url_uid AND rru.role_uid = $1 " +
		"WHERE parent_uid = $2 "
	rows, stmt, err := ss_sql.Query(dbHandler, sqlstr, roleNo, urlUid)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datas := []*auth.RoleUrlListData{}
	if err == nil {
		for rows.Next() {
			data := auth.RoleUrlListData{}
			rows.Scan(&data.Id, &data.Name, &data.ParentUid, &data.RoleNo)
			data.Children = getRoleUrlData(data.Id, roleNo)
			datas = append(datas, &data)
		}
	} else {
		ss_log.Error("err=[%v]", err)
	}
	return datas
}

func getRoleUrlData2(urlUid, roleNo, masterAccs string) []*auth.RoleUrlListData {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlstr := "SELECT u.url_uid,u.url_name,u.parent_uid,COALESCE(rru.role_uid, '00000000-0000-0000-0000-000000000000') as role_no FROM  url u " +
		"LEFT JOIN rela_role_url rru ON rru.url_uid = u.url_uid AND rru.role_uid = $1 " +
		"WHERE parent_uid = $2 and u.url_uid in ('" + masterAccs + "')"
	rows, stmt, err := ss_sql.Query(dbHandler, sqlstr, roleNo, urlUid)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datas := []*auth.RoleUrlListData{}
	if err == nil {
		for rows.Next() {
			data := auth.RoleUrlListData{}
			rows.Scan(&data.Id, &data.Name, &data.ParentUid, &data.RoleNo)
			data.Children = getRoleUrlData2(data.Id, roleNo, masterAccs)
			datas = append(datas, &data)
		}
	} else {
		ss_log.Error("err=[%v]", err)
	}
	return datas
}

/**
 * 获取账户列表
 */
func (*Auth) GetRoleInfo(ctx context.Context, req *auth.GetRoleInfoRequest, reply *auth.GetRoleInfoReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var data auth.RoleData
	row, stmt, err := ss_sql.QueryRowN(dbHandler, "SELECT r.role_no,r.role_name,r.create_time,r.modify_time,"+
		"r.acc_type,r.master_acc,acc.account FROM role r left join account acc on acc.uid=r.master_acc WHERE r.role_no = $1 LIMIT 1",
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

/**
 * 获取账户列表
 */
func (*Auth) GetRole(ctx context.Context, req *auth.GetRoleRequest, reply *auth.GetRoleReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	row, stmt, err := ss_sql.QueryRowN(dbHandler, "SELECT r.role_no,r.role_name,r.create_time,r.modify_time,"+
		"r.master_acc,acc.account,r.top_agency FROM role r left join account acc on acc.uid=r.master_acc WHERE r.role_name = $1 LIMIT 1",
		req.RoleName)
	if stmt != nil {
		defer stmt.Close()
	}
	if nil != err {
		ss_log.Error("GetRole|err=%v", err.Error())
		debug.PrintStack()
		reply.ResultCode = ss_err.ERR_ARGS
		return nil
	}
	var masterAcc, masterAccount sql.NullString
	err = row.Scan(&reply.RoleNo, &reply.RoleName, &reply.CreateTime, &reply.ModifyTime, &masterAcc, &masterAccount, &reply.TopAgency)
	reply.MasterAcc = masterAcc.String
	reply.MasterAccount = masterAccount.String
	if nil != err {
		ss_log.Error("GetRole|err=%v", err.Error())
		debug.PrintStack()
		reply.ResultCode = ss_err.ERR_ARGS
		return nil
	}

	rowsUrl, stmt, err := ss_sql.Query(dbHandler, "select url_uid,url_name,url from url where url_uid in (select url_uid from rela_role_url where role_uid = $1)",
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

	rowsUrl2, stmt, err := ss_sql.Query(dbHandler, "select url_uid,url_name,url from url")
	if stmt != nil {
		defer stmt.Close()
	}
	defer rowsUrl2.Close()

	if nil == err {
		reply.ResultCode = ss_err.ERR_SUCCESS
		for rowsUrl.Next() {
			data := auth.GetRoleUrlData{}
			rowsUrl.Scan(&data.UrlUid, &data.Name, &data.Url)
			reply.UrlData = append(reply.UrlData, &data)
		}

		for rowsUrl2.Next() {
			data := auth.GetRoleUrlData{}
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

/**
 * 删除账户
 */
func (*Auth) DeleteRole(ctx context.Context, req *auth.DeleteRoleRequest, reply *auth.DeleteRoleReply) error {
	errCode := dao.RoleDaoInst.DeleteRole(req.RoleNo)
	reply.ResultCode = errCode
	return nil
}

func getKey(key string, lList []string) int {
	for k, v := range lList {
		if v == key {
			return k
		}
	}
	return -1
}

func remove(idx int, lList []string) []string {
	return append(lList[:idx], lList[(idx+1):]...)
}

/**
 * 更新或者插入授权信息
 */
func (*Auth) UpdateOrInsertRoleAuth(ctx context.Context, req *auth.UpdateOrInsertRoleAuthRequest, reply *auth.UpdateOrInsertRoleAuthReply) error {
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
	rowsOldUrl, stmt, errSel := ss_sql.Query(dbHandler, "select url_uid from rela_role_url where role_uid = $1", req.RoleNo)
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
	pDel, errDelete := dbHandler.Prepare("delete from rela_role_url where role_uid = $1 and url_uid = $2")
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

	pInsert, errInsert := dbHandler.Prepare("insert into rela_role_url(rela_uid,url_uid,role_uid) VALUES ($1,$2,$3)")
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

/**
 * 更新或者插入新主商户
 */
func (*Auth) UpdateOrInsertRole(ctx context.Context, req *auth.UpdateOrInsertRoleRequest, reply *auth.UpdateOrInsertRoleReply) error {
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
		errInsert := ss_sql.Exec(dbHandler, "INSERT INTO role(role_no,role_name,create_time,acc_type,master_acc,is_delete) "+
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
		err := ss_sql.QueryRow(dbHandler, "SELECT count(1) FROM role "+
			"WHERE role_no = $1 LIMIT 1", []*sql.NullString{&cnt}, req.RoleNo)
		if !cnt.Valid || strext.ToInt32(cnt.String) <= 0 || err != nil {
			tx.Rollback()
			ss_log.Error("err|UpdateOrInsertRole|no this roleUid=%v\n", req.RoleNo)
			reply.ResultCode = ss_err.ERR_ACCOUNT_ROLE_NOT_EXISTS
			return nil
		}

		errUpdate := ss_sql.Exec(dbHandler, "UPDATE role SET role_name = $1, modify_time = current_timestamp, "+
			"acc_type=$2,master_acc=$4 WHERE role_no = $3", req.RoleName, req.AccType, req.RoleNo, req.MasterAcc)
		//if req.AccType == constants.AccountType_AUTH {
		//	err := ss_sql.ExecTx(tx, "insert  into  rela_role_account_auth (uid,account_uid,role_uid) values ($1,$2,$3)", strext.NewUUID(), req.AccUid, req.RoleNo)
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

func (*Auth) AuthRole(ctx context.Context, req *auth.AuthRoleRequest, reply *auth.AuthRoleReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, err := dbHandler.Begin()
	if nil != err {
		ss_log.Error("err=[%v]", err)
		tx.Rollback()
		reply.ResultCode = ss_err.ERR_SYS_DB_INIT
		return nil
	}
	sqlU := "update role set def_type=$1 where role_no=$2"
	err = ss_sql.ExecTx(tx, sqlU, req.DefType, req.RoleNo)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		tx.Rollback()
		reply.ResultCode = ss_err.ERR_SYS_DB_INIT
		return nil
	}
	//var topAgency sql.NullString
	//sqlS := "select top_agency from role where role_no=$1 limit 1"
	//err = ss_sql.QueryRowTx(tx, sqlS, []*sql.NullString{&topAgency}, req.RoleNo)
	//if nil != err {
	//	ss_log.Error("err=[%v]", err)
	//	tx.Rollback()
	//	reply.ResultCode = ss_err.ERR_SYS_DB_INIT
	//	return nil
	//}
	//if req.DefType == "1" {
	//	sqlStr := "update role set def_type='0' where acc_type=$1 and role_no <> $2 and top_agency=$3"
	//	err := ss_sql.ExecTx(tx, sqlStr, req.AccType, req.RoleNo, topAgency.String)
	//	if nil != err {
	//		ss_log.Error("err=[%v]", err)
	//		tx.Rollback()
	//		reply.ResultCode = ss_err.ERR_SYS_DB_INIT
	//		return nil
	//	}
	//}
	if req.DefType == "1" {
		sqlStr := "update role set def_type='0' where acc_type=$1 and role_no <> $2 "
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
