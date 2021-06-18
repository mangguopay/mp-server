package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
	"strings"
)

type PushMsgTypeDao struct {
}

var (
	PushMsgTypeDaoInst PushMsgTypeDao
)

// 获取模板
func (r *PushMsgTypeDao) GetTemplate(tempNo string) (pushNoList []string, titleKey, contentKey string, lenArgs int32, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select push_nos,title_key,content_key,len_args from push_temp where temp_no=$1 and is_delete = '0' limit 1"
	var pushNos, titleKeyT, contentKeyT, lenArgsT sql.NullString
	err = ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&pushNos, &titleKeyT, &contentKeyT, &lenArgsT}, tempNo)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return nil, "", "", 0, err
	}

	pushNoList = strings.Split(pushNos.String, ",")
	return pushNoList, titleKeyT.String, contentKeyT.String, strext.ToInt32(lenArgsT.String), nil
}
