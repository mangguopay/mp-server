package handler

import (
	"context"
	"database/sql"
	"runtime/debug"

	"a.a/mp-server/admin-auth-srv/dao"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	adminAuthProto "a.a/mp-server/common/proto/admin-auth"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

//==============================================================|菜单

// 获取菜单详情
func (a *AdminAuth) GetAdminMenu(ctx context.Context, req *adminAuthProto.GetAdminMenuRequest, reply *adminAuthProto.GetAdminMenuReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlstr := "SELECT url_uid,url_name,url,parent_uid,title,icon,component_name,component_path,redirect,idx,is_hidden FROM admin_url WHERE url_uid = $1 LIMIT 1 "
	row, stmt, err := ss_sql.QueryRowN(dbHandler, sqlstr, req.UrlUid)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		ss_log.Error("GetAdminMenu|err=[%v]\n", err)
		reply.ResultCode = ss_err.ERR_ACCOUNT_MENU_NOT_EXISTS
		return nil
	}
	var data adminAuthProto.AdminRouteData
	err = row.Scan(&data.UrlUid,
		&data.UrlName,
		&data.Url,
		&data.ParentUid,
		&data.Title,
		&data.Icon,
		&data.ComponentName,
		&data.ComponentPath,
		&data.Redirect,
		&data.Idx,
		&data.IsHidden)
	if err != nil {
		ss_log.Error("GetAdminMenu|err=[%v]\n", err)
		reply.ResultCode = ss_err.ERR_ACCOUNT_MENU_NOT_EXISTS
		reply.Data = &data
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = &data
	return nil
}

// 获取admin菜单列表
func (a *AdminAuth) GetAdminMenuList(ctx context.Context, req *adminAuthProto.GetAdminMenuListRequest, reply *adminAuthProto.GetAdminMenuListReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "url", Val: req.Search, EqType: "like"},
		//{Key: "url_name", Val: req.Search, EqType: "like"},
	})

	sqlCnt := "SELECT count(1) " +
		" FROM admin_url " + whereModel.WhereStr
	var total sql.NullString
	errCnt := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...)
	if errCnt != nil {
		ss_log.Error("errCnt=[%v]", errCnt)
	}

	//where := " WHERE url like '%" + req.Search + "%' OR url_name like '%" + req.Search + "%' "
	//sqlstr := "SELECT url_uid,url_name,url,parent_uid,title,icon,component_name,component_path,redirect,idx,is_hidden FROM admin_url  " +
	//	where + " order by create_time desc,url_uid LIMIT $1 OFFSET $2"

	ss_sql.SsSqlFactoryInst.AppendWhereOrGroup(whereModel, "url_name", req.Search, "like", "")
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by create_time desc,url_uid `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
	sqlStr := "SELECT url_uid,url_name,url,parent_uid,title,icon,component_name,component_path,redirect,idx,is_hidden " +
		" FROM admin_url  " + whereModel.WhereStr

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("GetAdminMenuList|err=[%v]\n", err)
		reply.ResultCode = ss_err.ERR_ACCOUNT_MENU_NOT_EXISTS
		reply.Total = int32(0)
		return nil
	}
	var datas []*adminAuthProto.AdminRouteData
	for rows.Next() {
		data := adminAuthProto.AdminRouteData{}
		err = rows.Scan(
			&data.UrlUid,
			&data.UrlName,
			&data.Url,
			&data.ParentUid,
			&data.Title,
			&data.Icon,
			&data.ComponentName,
			&data.ComponentPath,
			&data.Redirect,
			&data.Idx,
			&data.IsHidden,
		)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		datas = append(datas, &data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Total = strext.ToInt32(total.String)
	reply.DataList = datas
	return nil
}

// 更新或者插入新主商户
func (a *AdminAuth) SaveOrInsertAdminMenu(ctx context.Context, req *adminAuthProto.SaveOrInsertAdminMenuRequest, reply *adminAuthProto.SaveOrInsertAdminMenuReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, err := dbHandler.Begin()
	if nil != err {
		tx.Rollback()
		ss_log.Error("SaveOrInsertAdminMenu|Begin|err=[%v]\n", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_INIT
		return nil
	}

	if "" == req.UrlUid {
		// 如果没有uid，那就尝试插入新记录
		errInsert := ss_sql.Exec(dbHandler, "INSERT INTO admin_url (url_uid,url_name,url,parent_uid,title,icon,component_name,component_path,redirect,idx,is_hidden, create_time) "+
			"VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11, current_timestamp)",
			strext.NewUUID(), req.UrlName, req.Url, req.ParentUid, req.Title, req.Icon, req.ComponentName, req.ComponentPath, req.Redirect, req.Idx, req.IsHidden)
		if nil != errInsert {
			tx.Rollback()
			debug.PrintStack()
			ss_log.Error("SaveOrInsertAdminMenu|INSERT|err=[%v]", errInsert)
			reply.ResultCode = ss_err.ERR_SYS_DB_ADD
			return nil
		}
	} else {
		// 有uid，则更新
		var tmp sql.NullString
		err = ss_sql.QueryRow(dbHandler, "SELECT 1 FROM admin_url "+
			"WHERE url_uid = $1 LIMIT 1 FOR UPDATE SKIP LOCKED", []*sql.NullString{&tmp}, req.UrlUid)
		if "" == tmp.String || err != nil {
			ss_log.Error("err=[%v]", err)
			tx.Rollback()
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}

		var errUpdate error
		errUpdate = ss_sql.Exec(dbHandler, "UPDATE admin_url SET url_name = $1, url = $2, parent_uid = $3, title = $4, icon = $5, "+
			"component_name = $6, component_path = $7, idx = $8, is_hidden = $9 WHERE url_uid = $10",
			req.UrlName, req.Url, req.ParentUid, req.Title, req.Icon, req.ComponentName, req.ComponentPath, req.Idx, req.IsHidden, req.UrlUid)

		if nil != errUpdate {
			tx.Rollback()
			ss_log.Error("SaveOrInsertAdminMenu|err=[%v]\n", errUpdate)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 删除菜单
func (a *AdminAuth) DeleteAdminMenu(ctx context.Context, req *adminAuthProto.DeleteAdminMenuRequest, reply *adminAuthProto.DeleteAdminMenuReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	//查询该菜单下是否有子菜单
	boolA, err := dao.MenuDaoInstance.CheckAdminMenuHaveChild(req.UrlUid)
	if err != nil {
		ss_log.Error("查询菜单[%v]是否拥有子菜单失败,err=[%v]", req.UrlUid, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}

	if boolA {
		ss_log.Error("该菜单[%v]拥有子菜单，不允许删除", req.UrlUid)
		reply.ResultCode = ss_err.ERR_Menu_Have_Child_ERR
		return nil
	}

	if err := dao.MenuDaoInstance.DelAdminMenuHa(req.UrlUid); err != nil {
		ss_log.Error("删除菜单[%v]失败，err=[%v]", req.UrlUid, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (a *AdminAuth) AdminMenuRefreshChild(ctx context.Context, req *adminAuthProto.AdminMenuRefreshChildRequest, reply *adminAuthProto.AdminMenuRefreshChildReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	idx := 1
	for _, v := range req.UrlNo {
		err := ss_sql.Exec(dbHandler, "update admin_url set idx=$1 where url_uid=$2", idx, v)
		if nil != err {
			ss_log.Error("err=[%v]", err)
			continue
		}
		idx++
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 获取菜单详情
func (a *AdminAuth) GetMenu(ctx context.Context, req *adminAuthProto.GetMenuRequest, reply *adminAuthProto.GetMenuReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlstr := "SELECT url_uid,url_name,url,parent_uid,title,icon,component_name,component_path,redirect,idx,is_hidden " +
		" FROM url WHERE url_uid = $1 LIMIT 1 "
	row, stmt, err := ss_sql.QueryRowN(dbHandler, sqlstr, req.UrlUid)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		ss_log.Error("GetMenu|err=[%v]\n", err)
		reply.ResultCode = ss_err.ERR_ACCOUNT_MENU_NOT_EXISTS
		return nil
	}
	var data adminAuthProto.RouteData
	err = row.Scan(&data.UrlUid,
		&data.UrlName,
		&data.Url,
		&data.ParentUid,
		&data.Title,
		&data.Icon,
		&data.ComponentName,
		&data.ComponentPath,
		&data.Redirect,
		&data.Idx,
		&data.IsHidden)
	if err != nil {
		ss_log.Error("GetAdminMenu|err=[%v]\n", err)
		reply.ResultCode = ss_err.ERR_ACCOUNT_MENU_NOT_EXISTS
		reply.Data = &data
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = &data
	return nil
}

// 获取菜单列表
func (a *AdminAuth) GetMenuList(ctx context.Context, req *adminAuthProto.GetMenuListRequest, reply *adminAuthProto.GetMenuListReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "url", Val: req.Search, EqType: "like"},
		//{Key: "url_name", Val: req.Search, EqType: "like"},
	})

	sqlCnt := "SELECT count(1) " +
		" FROM url " + whereModel.WhereStr
	var total sql.NullString
	errCnt := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...)
	if errCnt != nil {
		ss_log.Error("errCnt=[%v]", errCnt)
	}

	//where := " WHERE url like '%" + req.Search + "%' OR url_name like '%" + req.Search + "%' "
	//sqlstr := "SELECT url_uid,url_name,url,parent_uid,title,icon,component_name,component_path,redirect,idx,is_hidden FROM admin_url  " +
	//	where + " order by create_time desc,url_uid LIMIT $1 OFFSET $2"

	ss_sql.SsSqlFactoryInst.AppendWhereOrGroup(whereModel, "url_name", req.Search, "like", "")
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by create_time desc,url_uid `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
	sqlStr := "SELECT url_uid,url_name,url,parent_uid,title,icon,component_name,component_path,redirect,idx,is_hidden " +
		" FROM url  " + whereModel.WhereStr

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("GetMenuList|err=[%v]\n", err)
		reply.ResultCode = ss_err.ERR_ACCOUNT_MENU_NOT_EXISTS
		reply.Total = int32(0)
		return nil
	}
	var datas []*adminAuthProto.RouteData
	for rows.Next() {
		data := adminAuthProto.RouteData{}
		err = rows.Scan(
			&data.UrlUid,
			&data.UrlName,
			&data.Url,
			&data.ParentUid,
			&data.Title,
			&data.Icon,
			&data.ComponentName,
			&data.ComponentPath,
			&data.Redirect,
			&data.Idx,
			&data.IsHidden,
		)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		datas = append(datas, &data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Total = strext.ToInt32(total.String)
	reply.DataList = datas
	return nil
}

// 更新或者插入新主商户
func (a *AdminAuth) SaveOrInsertMenu(ctx context.Context, req *adminAuthProto.SaveOrInsertMenuRequest, reply *adminAuthProto.SaveOrInsertMenuReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, err := dbHandler.Begin()
	if nil != err {
		tx.Rollback()
		ss_log.Error("SaveOrInsertMenu|Begin|err=[%v]\n", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_INIT
		return nil
	}

	if "" == req.UrlUid {
		// 如果没有uid，那就尝试插入新记录
		errInsert := ss_sql.Exec(dbHandler, "INSERT INTO url (url_uid,url_name,url,parent_uid,title,icon,component_name,component_path,redirect,idx,is_hidden, create_time) "+
			"VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11, current_timestamp)",
			strext.NewUUID(), req.UrlName, req.Url, req.ParentUid, req.Title, req.Icon, req.ComponentName, req.ComponentPath, req.Redirect, req.Idx, req.IsHidden)
		if nil != errInsert {
			tx.Rollback()
			debug.PrintStack()
			ss_log.Error("SaveOrInsertMenu|INSERT|err=[%v]", errInsert)
			reply.ResultCode = ss_err.ERR_SYS_DB_ADD
			return nil
		}
	} else {
		// 有uid，则更新
		var tmp sql.NullString
		err = ss_sql.QueryRow(dbHandler, "SELECT 1 FROM url "+
			"WHERE url_uid = $1 LIMIT 1 FOR UPDATE SKIP LOCKED", []*sql.NullString{&tmp}, req.UrlUid)
		if "" == tmp.String || err != nil {
			ss_log.Error("err=[%v]", err)
			tx.Rollback()
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}

		var errUpdate error
		errUpdate = ss_sql.Exec(dbHandler, "UPDATE url SET url_name = $1, url = $2, parent_uid = $3, title = $4, icon = $5, "+
			"component_name = $6, component_path = $7, idx = $8, is_hidden = $9 WHERE url_uid = $10",
			req.UrlName, req.Url, req.ParentUid, req.Title, req.Icon, req.ComponentName, req.ComponentPath, req.Idx, req.IsHidden, req.UrlUid)

		if nil != errUpdate {
			tx.Rollback()
			ss_log.Error("SaveOrInsertMenu|err=[%v]\n", errUpdate)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 删除菜单
func (a *AdminAuth) DeleteMenu(ctx context.Context, req *adminAuthProto.DeleteMenuRequest, reply *adminAuthProto.DeleteMenuReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	//查询该菜单下是否有子菜单
	boolA, err := dao.MenuDaoInstance.CheckMenuHaveChild(req.UrlUid)
	if err != nil {
		ss_log.Error("查询菜单[%v]是否拥有子菜单失败,err=[%v]", req.UrlUid, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}

	if boolA {
		ss_log.Error("该菜单[%v]拥有子菜单，不允许删除", req.UrlUid)
		reply.ResultCode = ss_err.ERR_Menu_Have_Child_ERR
		return nil
	}

	if err := dao.MenuDaoInstance.DelMenuHa(req.UrlUid); err != nil {
		ss_log.Error("删除菜单[%v]失败，err=[%v]", req.UrlUid, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (a *AdminAuth) MenuRefreshChild(ctx context.Context, req *adminAuthProto.MenuRefreshChildRequest, reply *adminAuthProto.MenuRefreshChildReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	idx := 1
	for _, v := range req.UrlNo {
		err := ss_sql.Exec(dbHandler, "update url set idx=$1 where url_uid=$2", idx, v)
		if nil != err {
			ss_log.Error("err=[%v]", err)
			continue
		}
		idx++
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
