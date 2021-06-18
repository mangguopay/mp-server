package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

var (
	MenuDaoInstance MenuDao
)

type MenuDao struct {
}

//有子菜单返回true,没有则返回false
func (r *MenuDao) CheckAdminMenuHaveChild(parentUid string) (bool, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var tem sql.NullString
	sqlStr := "select count(1) from admin_url where parent_uid = $1 "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&tem}, parentUid)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return false, err
	}

	return tem.String != "0", nil
}

func (r *MenuDao) DelAdminMenuHa(parentUid string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := " DELETE FROM admin_url WHERE url_uid = $1 "
	err := ss_sql.Exec(dbHandler, sqlStr, parentUid)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}

	return nil
}

//有子菜单返回true,没有则返回false
func (r *MenuDao) CheckMenuHaveChild(parentUid string) (bool, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var tem sql.NullString
	sqlStr := "select count(1) from url where parent_uid = $1 "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&tem}, parentUid)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return false, err
	}

	return tem.String != "0", nil
}

func (r *MenuDao) DelMenuHa(parentUid string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := " DELETE FROM url WHERE url_uid = $1 "
	err := ss_sql.Exec(dbHandler, sqlStr, parentUid)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}

	return nil
}
