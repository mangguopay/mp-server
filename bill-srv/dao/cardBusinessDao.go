package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type CardBusinessDao struct {
}

var CardBusinessDaoInst CardBusinessDao

func (*CardBusinessDao) QueryNameAndNumFromNo(cardNo string) (string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var nameT, numT, balanceTypeT, channelNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select name, card_number, balance_type, channel_no from card_business where card_no=$1  and is_delete = '0' and audit_status = '1' and collect_status = '1' limit 1`,
		[]*sql.NullString{&nameT, &numT, &balanceTypeT, &channelNo}, cardNo)
	if err != nil {
		ss_log.Error("err=%v", err)
		return "", "", "", ""
	}
	return nameT.String, numT.String, balanceTypeT.String, channelNo.String
}
func (*CardBusinessDao) QueryNameFromNumAndChennel(cardNum, chennelName, channelType string) (string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var cardNoT, balanceTypeT sql.NullString
	sqlStr := `select c.card_no,c.balance_type 
				from card_business c
 				LEFT JOIN channel ch ON c.channel_no = ch.channel_no
				WHERE ch.channel_name= $1 and c.card_number=$2 
					and ch.channel_type in ('0',$3)	
					and c.collect_status='1'
					and c.audit_status='1'
					and c.is_delete = '0'
					limit 1`
	err := ss_sql.QueryRow(dbHandler, sqlStr,
		[]*sql.NullString{&cardNoT, &balanceTypeT}, chennelName, cardNum, channelType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return "", ""
	}
	return cardNoT.String, balanceTypeT.String
}
