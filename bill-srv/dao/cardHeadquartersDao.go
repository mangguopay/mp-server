package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type CardHeadquartersDao struct {
}

var CardHeadquartersDaoInst CardHeadquartersDao

func (*CardHeadquartersDao) QueryNameAndNumFromNo(cardNo, accountType string) (string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var nameT, numT, balanceTypeT, channelNoT sql.NullString
	//err := ss_sql.QueryRow(dbHandler, `select name,card_number,balance_type,channel_no from card_head where card_no=$1  and is_delete = '0'  and collect_status = '1' and account_type = $2 limit 1`,

	sqlStr := ""
	switch accountType { //总部卡是用于什么账号类型的
	case constants.AccountType_USER:
		sqlStr = "select ch.name,ch.card_number,ccc.currency_type,ccc.channel_no " +
			" from card_head  ch " +
			" LEFT JOIN channel_cust_config ccc ON ccc.id = ch.channel_cust_config_id" +
			" where ch.card_no= $1 and ch.is_delete = '0' and ch.collect_status = '1' and ch.account_type = $2 limit 1"
	case constants.AccountType_SERVICER:
		sqlStr = "select name,card_number,balance_type,channel_no " +
			" from card_head where card_no=$1  and is_delete = '0'  and collect_status = '1' and account_type = $2 limit 1"
	default:
		ss_log.Error("accountType错误")
		return "", "", "", ""
	}
	err := ss_sql.QueryRow(dbHandler, sqlStr,
		[]*sql.NullString{&nameT, &numT, &balanceTypeT, &channelNoT}, cardNo, accountType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return "", "", "", ""
	}
	return nameT.String, numT.String, balanceTypeT.String, channelNoT.String
}
