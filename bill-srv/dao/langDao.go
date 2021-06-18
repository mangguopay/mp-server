package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

var (
	LangDaoInstance LangDao
)

type LangDao struct {
	LangKm string
	LangEn string
	LangCh string
	Key    string
	Type   string
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

//返回多个key对应的语言(注意图片只会得到id)
func (LangDao) GetLangTextsByKeys(keys []string) (datas []*LangDao, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	if len(keys) == 0 {
		return nil, nil
	}

	inStr := "("

	for k, v := range keys {
		if k != 0 { //不处于第一个则先加,再加'key'
			inStr += ","
		}
		inStr += "'" + v + "'"
	}

	inStr += ")"

	whereList := []*model.WhereSqlCond{
		{Key: "key", Val: inStr, EqType: "in"},
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlStr := "select key, lang_km, lang_en, lang_ch, type from lang "
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr+whereModel.WhereStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	for rows.Next() {
		data := LangDao{}
		var key, langKm, langEn, langCh, typeStr sql.NullString
		err2 = rows.Scan(
			&key,
			&langKm,
			&langEn,
			&langCh,
			&typeStr,
		)
		if err2 != nil {
			ss_log.Error("err=[%v]", err2)
			return nil, err2
		}

		data.Key = key.String
		data.LangKm = langKm.String
		data.LangEn = langEn.String
		data.LangCh = langCh.String
		data.Type = typeStr.String

		datas = append(datas, &data)
	}
	return datas, nil
}
