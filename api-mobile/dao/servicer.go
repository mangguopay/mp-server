package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type ServiceDao struct {
}

var ServiceDaoInst ServiceDao

func (ServiceDao) GetServicerNoByCashierNo(cid string) (servicerNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var Sno sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select servicer_no from cashier where uid =$1 and is_delete='0' limit 1`, []*sql.NullString{&Sno}, cid)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return Sno.String
}

func (ServiceDao) GetLatLngInfoFromNo(servicerNo string) (string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var latT, lngT, scopeT sql.NullString
	sqlStr := "select  lat, lng, scope from servicer where servicer_no =$1 and is_delete = $2 and use_status = $3"
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&latT, &lngT, &scopeT}, servicerNo, 0, 1)
	if err != nil {
		ss_log.Error("ServiceDao | GetLatLngInfoFromNo |  err= [%v]", err)
		return "", "", ""
	}
	return latT.String, lngT.String, scopeT.String

}

func (ServiceDao) GetScopeOffNoBySrvNo(srvNo string) (string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select s.scope_off,st.pos_sn from servicer s left join servicer_terminal st on s.servicer_no = st.servicer_no" +
		"  where st.servicer_no = $1 and st.is_delete = '0' "
	var scopeOff, posSn sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&scopeOff, &posSn}, srvNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "", ""
	}

	return scopeOff.String, posSn.String
}
