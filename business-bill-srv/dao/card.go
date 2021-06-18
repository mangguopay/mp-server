package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type Card struct {
	CardNo       string
	Name         string
	CardNumber   string
	CurrencyType string
	AuditStatus  string
}

var CardDao Card

func (Card) GetCardBaseInfo(accountNo, bankCardNo string) (*Card, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT ca.card_no, ca.name, ca.card_number, ca.balance_type, ca.audit_status " +
		"FROM card ca " +
		"LEFT JOIN channel c ON c.channel_no = ca.channel_no " +
		"WHERE ca.account_no = $1 AND ca.card_no = $2 " +
		"AND ca.is_delete != 1 AND c.use_status = 1 LIMIT 1"

	var cardNo, name, cardNumber, currencyType, auditStatus sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cardNo, &name, &cardNumber, &currencyType, &auditStatus},
		accountNo, bankCardNo)
	if err != nil {
		return nil, err
	}

	obj := new(Card)
	obj.CardNo = cardNo.String
	obj.CardNumber = cardNumber.String
	obj.Name = name.String
	obj.CurrencyType = currencyType.String
	obj.AuditStatus = auditStatus.String

	return obj, nil
}
