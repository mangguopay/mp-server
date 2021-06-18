package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	adminAuthProto "a.a/mp-server/common/proto/admin-auth"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type AdminRoleDao struct {
}

var AdminRoleDaoInst AdminRoleDao

/**
 * 获取Admin角色列表
 */
func (*AdminRoleDao) GetAdminRoleList(roleName, accountType string, pageSize, page int32) (errCode string, datas []*adminAuthProto.AdminRoleData) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "r.role_name", Val: roleName, EqType: "like"},
		{Key: "r.acc_type", Val: accountType, EqType: "="},
		{Key: "r.acc_type", Val: "('" + constants.AccountType_ADMIN + "','" + constants.AccountType_OPERATOR + "')", EqType: "in"},
		{Key: "r.is_delete", Val: "0", EqType: "="},
	})
	ss_sql.SsSqlFactoryInst.AppendWhereOrderBy(whereModel, "r.create_time", false)
	ss_sql.SsSqlFactoryInst.AppendWhereLimitI32(whereModel, pageSize, page)
	rows, stmt, err := ss_sql.SsSqlFactoryInst.Query(dbHandler, "SELECT r.role_no,r.role_name,r.create_time,r.modify_time,"+
		"r.acc_type,r.def_type,r.master_acc,acc.account FROM admin_role r left join admin_account acc on acc.uid=r.master_acc ", whereModel)
	if stmt != nil {
		defer stmt.Close()
		defer rows.Close()
	}
	ss_log.Error("err=[%v]", err)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_ADD, nil
	}

	datas = []*adminAuthProto.AdminRoleData{}
	for rows.Next() {
		data := adminAuthProto.AdminRoleData{}
		var masterAcc, masterAccount sql.NullString
		err := rows.Scan(
			&data.RoleNo,
			&data.RoleName,
			&data.CreateTime,
			&data.ModifyTime,
			&data.AccType,
			&data.DefType,
			&masterAcc,
			&masterAccount)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		data.MasterAcc = masterAcc.String
		data.MasterAccount = masterAccount.String
		datas = append(datas, &data)
	}
	return ss_err.ERR_SUCCESS, datas
}

/**
 * 获取Admin角色列表数量
 */
func (*AdminRoleDao) GetAdminRoleCnt(roleName, accountType string) (cnt int32) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "r.role_name", Val: roleName, EqType: "like"},
		{Key: "r.acc_type", Val: accountType, EqType: "="},
		{Key: "r.is_delete", Val: "0", EqType: "="},
	})

	var total sql.NullString
	err := ss_sql.SsSqlFactoryInst.QueryRow(dbHandler, "SELECT count(1) FROM admin_role r ", []*sql.NullString{&total}, whereModel)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		total.String = "0"
	}
	return strext.ToInt32(total.String)
}

func (AdminRoleDao) DeleteAdminRole(roleNo string) (retCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, "update admin_role set is_delete='1' where role_no=$1", roleNo)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}
