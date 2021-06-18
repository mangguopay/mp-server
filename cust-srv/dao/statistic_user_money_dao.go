package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
	"strings"
)

type StatisticUserMoneyDao struct{}

var StatisticUserMoneyDaoInst StatisticUserMoneyDao

type StatisticUserMoneyData struct {
	CreateTime           string
	UserUseBalance       int64
	UserKhrBalance       int64
	UserUseFrozenBalance int64
	UserKhrFrozenBalance int64
}

// 获取统计数据-图
func (s *StatisticUserMoneyDao) GetStatisticUserMoneyTime(startTime string, endTime string) ([]StatisticUserMoneyData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT create_time, user_usd_balance, user_khr_balance, user_khr_frozen_balance, user_usd_frozen_balance FROM statistic_user_money "
	sqlStr += " WHERE create_time >= $1 AND create_time <= $2  ORDER BY create_time ASC "

	rows, stmt, qErr := ss_sql.Query(dbHandler, sqlStr, startTime, endTime)
	if stmt != nil {
		defer stmt.Close()
		defer rows.Close()
	}

	if qErr != nil {
		if qErr.Error() != ss_sql.DB_NO_ROWS_MSG {
			return nil, qErr
		}
		return nil, nil
	}

	var list []StatisticUserMoneyData

	for rows.Next() {
		var createTime, userUsdBalance, userKhrBalance, userKhrFrozenBalance, userUsdFrozenBalance sql.NullString

		err := rows.Scan(
			&createTime,
			&userUsdBalance,
			&userKhrBalance,
			&userKhrFrozenBalance,
			&userUsdFrozenBalance,
		)
		if err != nil {
			return nil, err
		}
		//bbb := ss_time.FormatTimeStr(createTime.String,"2006-01-02T15:04:05Z", "2006-01-02 15:04:05",global.Tz)
		createTime.String = ss_time.StripPostTime(createTime.String)
		createTime.String = strings.Replace(createTime.String, "-", "/", -1)
		//ss_log.Info("createTime.String：%v",createTime.String)
		//ss_log.Info("createTime.String2：%v",bbb)
		list = append(list, StatisticUserMoneyData{
			CreateTime:           createTime.String,
			UserUseBalance:       strext.ToInt64(userUsdBalance.String),
			UserKhrBalance:       strext.ToInt64(userKhrBalance.String),
			UserUseFrozenBalance: strext.ToInt64(userUsdFrozenBalance.String),
			UserKhrFrozenBalance: strext.ToInt64(userKhrFrozenBalance.String),
		})
	}

	return list, nil
}

// 获取统计数据-列表
func (s StatisticUserMoneyDao) GetStatisticUserMoneyTimeList(req *go_micro_srv_cust.GetStatisticUserMoneyListRequest) ([]*go_micro_srv_cust.StatisticUserMoneyData, int32, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// 初始化where模型
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "create_time", Val: req.StartTime, EqType: ">="},
		{Key: "create_time", Val: req.EndTime, EqType: "<="},
	})

	//统计
	var totalStr sql.NullString
	sqlCnt := "SELECT count(1) FROM statistic_user_money " + whereModel.WhereStr
	if err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalStr}, whereModel.Args...); err != nil {
		return nil, 0, err
	}

	total := strext.ToInt32(totalStr.String)
	if total == 0 {
		return nil, 0, nil
	}

	// 添加order by 和limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `ORDER BY create_time DESC`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	sqlStr := "SELECT id, user_usd_balance, user_khr_balance, user_usd_frozen_balance, user_khr_frozen_balance, create_time, update_time FROM statistic_user_money " + whereModel.WhereStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}

	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		return nil, 0, err
	}

	var dataList []*go_micro_srv_cust.StatisticUserMoneyData

	for rows.Next() {
		data := &go_micro_srv_cust.StatisticUserMoneyData{}
		var updateTime sql.NullString
		err = rows.Scan(
			&data.Id,
			&data.UsdBalance,
			&data.KhrBalance,
			&data.UsdFrozenBalance,
			&data.KhrFrozenBalance,
			&data.CreateTime,
			&updateTime,
		)
		if err != nil {
			return nil, 0, err
		}

		data.UpdateTime = updateTime.String

		// 修改日期显示
		//data.Date = ss_time.StripPostDate(data.Date)
		dataList = append(dataList, data)
	}

	return dataList, total, nil
}

func (r *StatisticUserMoneyDao) Insert(wc *UserMoneyDataCount, timeT string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	id := strext.GetDailyId()

	return ss_sql.Exec(dbHandler, `
	INSERT INTO statistic_user_money(id, user_usd_balance, user_khr_balance, user_khr_frozen_balance, user_usd_frozen_balance, create_time, update_time)
	VALUES($1, $2, $3, $4, $5, $6, current_timestamp)
	ON CONFLICT ON CONSTRAINT unq_money_row 
	DO UPDATE SET user_usd_balance = EXCLUDED.user_usd_balance,
	user_khr_balance = EXCLUDED.user_khr_balance,
	user_khr_frozen_balance = EXCLUDED.user_khr_frozen_balance,
	user_usd_frozen_balance = EXCLUDED.user_usd_frozen_balance,
	update_time = current_timestamp`,
		id, wc.UserUsdBalance, wc.UserKhrBalance, wc.UserKhrFrozenBalance, wc.UserUsdFrozenBalance, timeT)
}

//

// UserMoneyDataCount 用户余额统计
type UserMoneyDataCount struct {
	UserUsdBalance       int64 //用户usd余额
	UserKhrBalance       int64 //用户khr余额
	UserKhrFrozenBalance int64 //用户当前khr冻结金额
	UserUsdFrozenBalance int64 //用户当前usd冻结金额
}
