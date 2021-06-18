package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_auth "a.a/mp-server/common/proto/auth"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type CardDao struct{}

var CardDaoInstance CardDao

func (*CardDao) InsertCard(tx *sql.Tx, accountNo, channelNo, name, cardNum, balanceType string, isDefault int32) string {
	cardNoT := strext.NewUUID()
	err := ss_sql.ExecTx(tx, `insert into card(card_no,account_no,channel_no,name,is_delete,`+
		`card_number,balance_type,is_defalut,collect_status,audit_status,note,create_time) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,current_timestamp)`,
		cardNoT, accountNo, channelNo, name, 0, cardNum, balanceType, isDefault, 0, 0, "0")
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}

	return ss_err.ERR_SUCCESS
}

func (*CardDao) QueryCardNo(carNum string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var cardNOT sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT card_no FROM card WHERE  card_number= $1 ",
		[]*sql.NullString{&cardNOT}, carNum)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}
	return cardNOT.String
}

// 查询名字
func (*CardDao) QueryCardFromNo(cardNo string) (string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var numT, isDefaultT, channelNoT, accountNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT name,is_defalut,channel_no,account_no FROM card WHERE card_no=$1 and audit_status='1'",
		[]*sql.NullString{&numT, &isDefaultT, &channelNoT, &accountNoT}, cardNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "", "", "", ""
	}
	return numT.String, isDefaultT.String, channelNoT.String, accountNoT.String
}

func (*CardDao) QueryIsNewAddCard(accountNo, balanceType string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var cardNOT sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT card_no  FROM card WHERE account_no=$1 and balance_type=$2 and is_delete='0'",
		[]*sql.NullString{&cardNOT}, accountNo, balanceType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}
	return cardNOT.String
}

func (*CardDao) QueryIsDefaultCardNo(accountNo, balanceType, isDefalt string) (string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var cardNOT, cardNumberT sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT card_no,card_number  FROM card WHERE account_no=$1   and balance_type=$2 and is_defalut = $3 and is_delete='0' limit 1",
		[]*sql.NullString{&cardNOT, &cardNumberT}, accountNo, balanceType, isDefalt)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "", ""
	}
	return cardNOT.String, cardNumberT.String
}

func (*CardDao) UpdateIsDefault(tx *sql.Tx, status int32, cardNo string) string {
	err := ss_sql.ExecTx(tx, `update card set is_defalut = $1,modify_time=current_timestamp where card_no=$2`,
		status, cardNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SAVE_CARD_FAILD
	}
	return ss_err.ERR_SUCCESS
}

func (*CardDao) UpdateIsDefaultByCardNum(tx *sql.Tx, accountNo, balanceType, cardNum string, status int32) string {
	err := ss_sql.ExecTx(tx, `update card set is_defalut = $1,modify_time=current_timestamp where card_number=$2 and balance_type = $3 and account_no = $4 and is_delete='0'`,
		status, cardNum, balanceType, accountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_MODIFY_DEFAULT_CARD_FAILD
	}
	return ss_err.ERR_SUCCESS
}

func (*CardDao) QueryIsDefaultFromCardNo(cardNo string) (string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var isDefaultT, channelNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT is_defalut,channel_no FROM card WHERE card_no=$1 and is_delete = '0'",
		[]*sql.NullString{&isDefaultT, &channelNoT}, cardNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "", ""
	}
	return isDefaultT.String, channelNoT.String
}

func (*CardDao) UpdateIsDelete(tx *sql.Tx, cardNo, isDelete string) string {
	err := ss_sql.ExecTx(tx, `update card set is_delete = $1,modify_time=current_timestamp where card_no=$2 and is_delete='0'`,
		isDelete, cardNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return ss_err.ERR_SUCCESS
}

func (*CardDao) GetCustPaymentCard(accountNo string) (datasR []*go_micro_srv_auth.CardData, errR string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var datas []*go_micro_srv_auth.CardData

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ca.account_no", Val: accountNo, EqType: "="},
		{Key: "ca.is_delete", Val: "0", EqType: "="},
	})

	where2 := whereModel.WhereStr
	args2 := whereModel.Args

	sqlStr := "SELECT ca.card_no, ch.logo_img_no, ch.logo_img_no_grey, ch.color_begin, ch.color_end, ch.channel_name, ca.card_number,ca.is_defalut, ca.balance_type " +
		" FROM card ca " +
		" LEFT JOIN channel ch ON ch.channel_no = ca.channel_no " + where2

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err == nil {
		for rows.Next() {
			data := &go_micro_srv_auth.CardData{}
			var cardNumber, logoImgNoGrey, colorBegin, colorEnd sql.NullString
			err = rows.Scan(
				&data.CardNo,
				&data.LogoImgNo,
				&logoImgNoGrey,
				&colorBegin,
				&colorEnd,
				&data.ChannelName,
				&cardNumber,
				&data.IsDefalut,
				&data.CurrencyType,
			)
			if err != nil {
				ss_log.Error("err=%v", err)
				continue
			}

			data.LogoImgNoGrey = logoImgNoGrey.String
			data.ColorBegin = colorBegin.String
			data.ColorEnd = colorEnd.String

			if cardNumber.String != "" {
				//cardNumberStr :=cardNumber.String
				//data.CardNumber = cardNumberStr[len(cardNumberStr)-4:]
				data.CardNumber = cardNumber.String
			}

			datas = append(datas, data)
		}
	} else {
		ss_log.Error("err=[%v]", err)
		return nil, ss_err.ERR_SYS_DB_GET
	}
	return datas, ss_err.ERR_SUCCESS
}
