package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

var (
	LangDaoInstance LangDao
)

type LangDao struct {
}

//只返回一种语言对应的文字
func (r *LangDao) GetLangTextByKey(key, lang string) (langText string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select lang_km, lang_en, lang_ch from lang where key = $1 and is_delete = '0' and type ='1' "
	var langKm, langEn, langCh sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&langKm, &langEn, &langCh}, key)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	switch lang {
	case constants.LangKmKH:
		return langKm.String
	case constants.LangEnUS:
		return langEn.String
	case constants.LangZhCN:
		return langCh.String
	}

	//默认返回英语
	return langEn.String
}
