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
	"time"
)

type StatisticUserExchangeDao struct{}

var StatisticUserExchangeDaoInst StatisticUserExchangeDao

type StatisticUserExchangeData struct {
	Day           string
	Usd2khrNum    int64
	Usd2khrAmount int64
	Usd2khrFee    int64
	Khr2usdNum    int64
	Khr2usdAmount int64
	Khr2usdFee    int64
}

// 获取兑换的统计数据
func (s *StatisticUserExchangeDao) GetStatisticData(startDate string, endDate string) ([]StatisticUserExchangeData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	endDate, tErr := ss_time.TimeAfter(endDate, ss_time.DateFormat, time.Hour*24) // 下一天的日期
	if tErr != nil {
		return nil, tErr
	}

	sqlStr := "SELECT day, SUM(usd2khr_num) AS usd2khr_num, SUM(usd2khr_amount) AS usd2khr_amount, SUM(usd2khr_fee) AS usd2khr_fee, "
	sqlStr += " SUM(khr2usd_num) AS khr2usd_num, SUM(khr2usd_amount) AS khr2usd_amount, SUM(khr2usd_fee) AS khr2usd_fee "
	sqlStr += " FROM statistic_user_exchange WHERE day >= $1 and day < $2 GROUP BY day ORDER BY day ASC "

	rows, stmt, qErr := ss_sql.Query(dbHandler, sqlStr, startDate, endDate)
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

	var list []StatisticUserExchangeData

	for rows.Next() {
		var day, usd2khrNum, usd2khrAmount, usd2khrFee, khr2usdNum, khr2usdAmount, khr2usdFee sql.NullString

		err := rows.Scan(&day, &usd2khrNum, &usd2khrAmount, &usd2khrFee, &khr2usdNum, &khr2usdAmount, &khr2usdFee)
		if err != nil {
			return nil, err
		}

		list = append(list, StatisticUserExchangeData{
			Day:           ss_time.StripPostDate(day.String),
			Usd2khrNum:    strext.ToInt64(usd2khrNum.String),
			Usd2khrAmount: strext.ToInt64(usd2khrAmount.String),
			Usd2khrFee:    strext.ToInt64(usd2khrFee.String),
			Khr2usdNum:    strext.ToInt64(khr2usdNum.String),
			Khr2usdAmount: strext.ToInt64(khr2usdAmount.String),
			Khr2usdFee:    strext.ToInt64(khr2usdFee.String),
		})
	}

	return list, nil
}

// 获取兑换的统计数据-列表
func (s StatisticUserExchangeDao) GetStatisticDataList(req *go_micro_srv_cust.GetStatisticUserExchangeListRequest) ([]*go_micro_srv_cust.StatisticUserExchangeData, int32, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// 获取总记录数
	// 日期往后加1天，where条件中用小于
	endDate, retErr := ss_time.TimeAfter(req.EndDate, ss_time.DateFormat, time.Hour*24) // 下一天的日期
	if retErr != nil {
		return nil, 0, retErr
	}

	// 组合条件
	whereList := []*model.WhereSqlCond{
		{Key: "day", Val: req.StartDate, EqType: ">="},
		{Key: "day", Val: endDate, EqType: "<"},
	}

	// 初始化where模型
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	//统计
	var totalStr sql.NullString
	sqlCnt := "SELECT count(1) FROM statistic_user_exchange " + whereModel.WhereStr
	if err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalStr}, whereModel.Args...); err != nil {
		return nil, 0, err
	}

	total := strext.ToInt32(totalStr.String)
	if total == 0 {
		return nil, 0, nil
	}

	// 添加order by 和limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `ORDER BY day DESC`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	sqlStr := "SELECT id, day, usd2khr_num, usd2khr_amount, usd2khr_fee, khr2usd_num, khr2usd_amount, khr2usd_fee, create_time, update_time "
	sqlStr += " FROM statistic_user_exchange " + whereModel.WhereStr
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

	var dataList []*go_micro_srv_cust.StatisticUserExchangeData

	for rows.Next() {
		var data go_micro_srv_cust.StatisticUserExchangeData
		err = rows.Scan(
			&data.Id,
			&data.Date,
			&data.Usd2KhrNum,
			&data.Usd2KhrAmount,
			&data.Usd2KhrFee,
			&data.Khr2UsdNum,
			&data.Khr2UsdAmount,
			&data.Khr2UsdFee,
			&data.CreateTime,
			&data.UpdateTime,
		)
		if err != nil {
			return nil, 0, err
		}

		// 修改日期显示
		data.Date = ss_time.StripPostDate(data.Date)
		dataList = append(dataList, &data)
	}

	return dataList, total, nil
}

func (r *StatisticUserExchangeDao) Insert(wc *DataCount) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	id := strext.GetDailyId()
	return ss_sql.Exec(dbHandler, `INSERT INTO statistic_user_exchange(id,day,usd2khr_num,usd2khr_amount,usd2khr_fee,
		khr2usd_num,khr2usd_amount,khr2usd_fee,create_time,update_time)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,current_timestamp,current_timestamp)
		ON CONFLICT ON CONSTRAINT unq_exchange_row
		DO UPDATE SET usd2khr_num = EXCLUDED.usd2khr_num,usd2khr_amount = EXCLUDED.usd2khr_amount,usd2khr_fee = EXCLUDED.usd2khr_fee,
		khr2usd_num = EXCLUDED.khr2usd_num,khr2usd_amount = EXCLUDED.khr2usd_amount,khr2usd_fee = EXCLUDED.khr2usd_fee,
		update_time = current_timestamp`,
		id, wc.Day, wc.Usd2khrNum, wc.Usd2khrAmount, wc.Usd2khrFee, wc.Khr2usdNum, wc.Khr2usdAmount, wc.Khr2usdFee)
}
