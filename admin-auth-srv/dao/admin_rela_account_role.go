package dao

import (
	"database/sql"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

var AdminRelaAccountRoleDaoInst AdminRelaAccountRoleDao

type AdminRelaAccountRoleDao struct {
}

func (AdminRelaAccountRoleDao) DeleteAdminRelaAccountRoleTx(tx *sql.Tx, accountNo string) (retCode string) {
	err := ss_sql.ExecTx(tx, "delete from admin_rela_account_role where account_uid = $1 ", accountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_DELETE
	}
	return ss_err.ERR_SUCCESS
}

func (AdminRelaAccountRoleDao) InsertAdminRelaAccountRoleTx(tx *sql.Tx, accountNo, roleNo string) (retCode string) {
	err := ss_sql.ExecTx(tx, "insert into admin_rela_account_role (rela_uid,account_uid,role_uid) VALUES ($1,$2,$3)", strext.NewUUID(), accountNo, roleNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_ADD
	}
	return ss_err.ERR_SUCCESS
}

func (AdminRelaAccountRoleDao) GetAdminRoleAccTypeTx(tx *sql.Tx, roleNo string) (accType, retCode string) {
	var accTypeT sql.NullString
	err := ss_sql.QueryRowTx(tx, "select acc_type	from admin_role where role_no = $1 and is_delete = '0' ", []*sql.NullString{&accTypeT}, roleNo)
	if nil != err {
		ss_log.Error("根据角色id[%v]查询acc_type出错,err=%v", roleNo, err)
		return "", ss_err.ERR_PARAM
	}
	return accTypeT.String, ss_err.ERR_SUCCESS
}
