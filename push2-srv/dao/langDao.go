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

/*
func (r *LangDao) GetLangTextByKeys(dbHandler *sql.DB, keys []string, lang string) (langText []string) {
	var langText1 []string
	for _, v := range keys {
		temp := LangDaoInstance.GetLangTextByKey(dbHandler, v, lang)
		langText1 = append(langText1, temp)
	}
	return langText1
}
*/
//只返回多种语言的文字
func (r *LangDao) GetLangTextByKey(key string) (langKm, langEn, langCh string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select lang_km, lang_en, lang_ch from lang where key = $1 and is_delete = '0' and type ='1' "
	var langKmT, langEnT, langChT sql.NullString
	err = ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&langKmT, &langEnT, &langChT}, key)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return langKm, langEn, langCh, err
	}

	return langKmT.String, langEnT.String, langChT.String, err
}

// 获取对应语言的文本
func (r *LangDao) GetText(key, langStr string) (str string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select lang_en from lang where key = $1 and is_delete = '0' and type ='1' "
	switch langStr {
	case constants.LangZhCN:
		sqlStr = "select lang_ch from lang where key = $1 and is_delete = '0' and type ='1' "
	case constants.LangEnUS:
	case constants.LangKmKH:
		sqlStr = "select lang_km from lang where key = $1 and is_delete = '0' and type ='1' "
	}

	var strT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&strT}, key)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return strT.String
	}

	return strT.String
}
