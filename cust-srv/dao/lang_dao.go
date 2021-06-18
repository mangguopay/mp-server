package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type LangDao struct {
}

var LangDaoInst LangDao

func (LangDao) CheckLangKey(dbHandler *sql.DB, key, typeB string) (count int, err string) {
	var countT sql.NullString
	errT := ss_sql.QueryRow(dbHandler, "SELECT COUNT(1) FROM lang WHERE key = $1 and type = $2  and is_delete='0'", []*sql.NullString{&countT}, key, typeB)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return 0, ss_err.ERR_PARAM
	}

	return strext.ToInt(countT.String), ss_err.ERR_SUCCESS
}

func (LangDao) GetLangId(dbHandler *sql.DB, key, typeB string) (id string, err string) {
	var idT sql.NullString
	errT := ss_sql.QueryRow(dbHandler, "SELECT id FROM lang WHERE key = $1 and type = $2  and is_delete='0'", []*sql.NullString{&idT}, key, typeB)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", ss_err.ERR_PARAM
	}
	return idT.String, ss_err.ERR_SUCCESS
}

//只返回一种语言对应的文字
func (LangDao) GetLangTextByKey(dbHandler *sql.DB, key, lang string) (langText string) {
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

	//默认返回英文
	return langEn.String
}

//只返回一种语言对应的文字
func (LangDao) GetLangByKey(key, typeT string) (string, string, string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select lang_km, lang_en, lang_ch from lang where key = $1 and is_delete = '0' and type =$2 "
	var langKm, langEn, langCh sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&langKm, &langEn, &langCh}, key, typeT)
	if nil != err {
		return "", "", "", err
	}
	return langKm.String, langEn.String, langCh.String, nil
}

func (LangDao) Insert(tx *sql.Tx, key, typeT, langKm, langEn, langCh string) error {

	sqlStr := "insert into lang(key, type, lang_km, lang_en, lang_ch) " +
		"values($1,$2,$3,$4,$5) on conflict (key) do update set type=$2,lang_km=$3,lang_en=$4,lang_ch=$5 "
	return ss_sql.ExecTx(tx, sqlStr, key, typeT, langKm, langEn, langCh)
}

//返回多个key对应的语言(注意图片只会得到id)
func (LangDao) GetLangTextsByKeys(keys []string) (datas []*go_micro_srv_cust.LangData, err error) {
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
		data := go_micro_srv_cust.LangData{}
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
