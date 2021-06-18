package ss_serv

import (
	"testing"

	"a.a/cu/db"
	"a.a/mp-server/common/constants"
)

func initDB() {
	alias := constants.DB_CRM
	host := "10.41.1.241"
	port := "5432"
	user := "postgres"
	password := "123"
	name := "mp_crm"

	db.DoDBInitPostgres(alias, host, port, user, password, name)
}

func TestDoInitMultiFromDB(t *testing.T) {
	initDB()
	DoInitMultiFromDB(constants.DB_CRM)
}
