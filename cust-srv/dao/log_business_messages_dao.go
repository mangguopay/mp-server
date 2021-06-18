package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_sql"
)

type LogBusinessMessagesDao struct {
}

var LogBusinessMessagesDaoInst LogBusinessMessagesDao

//设置账号的所有消息为已读
func (*LogBusinessMessagesDao) ModiftAllRead(accNo string) (errR error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update log_business_messages set is_read = '1' where account_no = $1 and is_read = '0' "
	err := ss_sql.Exec(dbHandler, sqlStr, accNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

func (*LogBusinessMessagesDao) GetBusinessMessagesList(whereList []*model.WhereSqlCond, page, pageSize int) (datas []*go_micro_srv_cust.BusinessMessagesData, returnErr error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by bm.create_time desc, bm.is_read asc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, pageSize, page)

	sqlStr := "SELECT bm.log_no, bm.is_read, bm.account_no, bm.create_time, bm.content, bm.account_type, acc.account " +
		" FROM log_business_messages bm " +
		" left join account acc on acc.uid = bm.account_no " + whereModel.WhereStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	for rows.Next() {
		data := &go_micro_srv_cust.BusinessMessagesData{}
		err = rows.Scan(
			&data.LogNo,
			&data.IsRead,
			&data.AccountNo,
			&data.CreateTime,
			&data.Content,
			&data.AccountType,
			&data.Account,
		)

		if err != nil {
			ss_log.Error("err=%v", err)
			return nil, err
		}
		datas = append(datas, data)
	}
	return datas, nil
}

func (*LogBusinessMessagesDao) GetCnt(whereList []*model.WhereSqlCond) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlStr := "SELECT count(1) " +
		" FROM log_business_messages bm " +
		" left join account acc on acc.uid = bm.account_no " + whereModel.WhereStr
	var totalT sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, whereModel.Args...)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return "", errT
	}

	return totalT.String, nil
}
