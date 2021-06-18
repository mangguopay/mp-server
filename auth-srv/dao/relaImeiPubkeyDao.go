package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

var (
	RelaImeiPubkeyDaoInst RelaImeiPubkeyDao
)

type RelaImeiPubkeyDao struct {
}

func (r *RelaImeiPubkeyDao) InsertRelaImeiPubKey(imei, pubKey string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	relaNo := strext.NewUUID()
	err := ss_sql.Exec(dbHandler, `insert into rela_imei_pubkey(rela_no,imei,pub_key,create_time)values($1,$2,$3,current_timestamp)`,
		relaNo, imei, pubKey)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_ADD
	}

	return ss_err.ERR_SUCCESS
}
