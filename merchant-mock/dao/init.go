package dao

import (
	"a.a/cu/db"
	_ "github.com/lib/pq"
)

const (
	DbMerchantMock = "merchant-mock"
)

func InitDB() {
	alias := DbMerchantMock
	host := "10.41.1.241"
	port := "5432"
	user := "postgres"
	password := "123"
	name := "merchant-mock"

	db.DoDBInitPostgres(alias, host, port, user, password, name)
}
