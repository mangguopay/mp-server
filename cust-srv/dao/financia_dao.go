package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"context"
)

type FinanciaDao struct {
}

var FinanciaDaoInst FinanciaDao

//todo 再传多一个参数，创建的是服务商还是用户的收款账户
func (FinanciaDao) InsertChannel(channelNo, name, cardNumber, note, balanceType, isDefalut string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(context.TODO(), nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		return ss_err.ERR_SYS_DB_OP
	}
	defer ss_sql.Rollback(tx)

	//查询平台总账号
	_, accPlat, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")
	if isDefalut == "1" { //只能有一张默认推荐收款卡
		//修改原来默认的卡为不默认推荐
		//todo 此处应该加上是服务商的还是用户的
		sqlStr2 := "update card set is_defalut = '0' where account_no = $1 and balance_type = $2 "
		err2 := ss_sql.ExecTx(tx, sqlStr2, accPlat, balanceType)
		if err2 != nil {
			ss_log.Error("FinanciaDao |UpdateCard err2=[%v]", err2)
			return ss_err.ERR_PARAM
		}
	}

	sqlStr3 := "insert into card(card_no, account_no, channel_no, name, create_time, is_delete, card_number, balance_type, is_defalut, collect_status, audit_status, note) " +
		" values ($1,$2,$3,$4,current_timestamp,$5,$6,$7,$8,$9,$10,$11)"
	err3 := ss_sql.ExecTx(tx, sqlStr3, strext.NewUUID(), accPlat, channelNo, name, "0", cardNumber, balanceType, isDefalut, "1", "0", note)
	if err3 != nil {
		ss_log.Error("FinanciaDao |InsertChannel err3=[%v]", err3)
		return ss_err.ERR_PARAM
	}

	tx.Commit()
	return ss_err.ERR_SUCCESS
}

//
func (FinanciaDao) UpdateCard(cardNo, channelNo, name, cardNumber, note, balanceType, isDefalut string) string {
	ss_log.Info("FinanciaDao | UpdateCard")
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(context.TODO(), nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		return ss_err.ERR_SYS_DB_OP
	}
	defer ss_sql.Rollback(tx)

	//查询平台总账号
	_, accPlat, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")
	if isDefalut == "1" {
		//修改原来默认的卡为不默认推荐
		//todo 此处应该加上是服务商还是用户的
		sqlStr2 := "update card set is_defalut = '0' where account_no = $1 and balance_type = $2 and is_defalut ='1' "
		err2 := ss_sql.ExecTx(tx, sqlStr2, accPlat, balanceType)
		if err2 != nil {
			ss_log.Error("FinanciaDao |UpdateCard err2=[%v]", err2)
			return ss_err.ERR_PARAM
		}
	}

	//设置默认卡
	sqlStr3 := "update card set channel_no = $2,name =$3,card_number = $4,note = $5,balance_type=$6,is_defalut=$7 where card_no = $1 "
	err3 := ss_sql.ExecTx(tx, sqlStr3, cardNo, channelNo, name, cardNumber, note, balanceType, isDefalut)
	if err3 != nil {
		ss_log.Error("FinanciaDao |UpdateCard err3=[%v]", err3)
		return ss_err.ERR_PARAM
	}

	tx.Commit()
	return ss_err.ERR_SUCCESS
}
