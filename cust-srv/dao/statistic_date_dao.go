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

type StatisticDateDao struct{}

var StatisticDateDaoInst StatisticDateDao

type StatisticDateData struct {
	Day            string
	RegUserNum     int64
	RegServicerNum int64
}

// 获取按天统计数据
func (s *StatisticDateDao) GetStatisticData(startDate string, endDate string) ([]StatisticDateData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	endDate, tErr := ss_time.TimeAfter(endDate, ss_time.DateFormat, time.Hour*24) // 下一天的日期
	if tErr != nil {
		return nil, tErr
	}

	sqlStr := "SELECT day, SUM(reg_user_num) AS reg_user_num, SUM(reg_servicer_num) AS reg_servicer_num FROM statistic_date "
	sqlStr += " WHERE day >= $1 AND day < $2  GROUP BY day ORDER BY day ASC "

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

	var list []StatisticDateData

	for rows.Next() {
		var day, regUserNum, regServicerNum sql.NullString

		err := rows.Scan(&day, &regUserNum, &regServicerNum)
		if err != nil {
			return nil, err
		}

		list = append(list, StatisticDateData{
			Day:            ss_time.StripPostDate(day.String),
			RegUserNum:     strext.ToInt64(regUserNum.String),
			RegServicerNum: strext.ToInt64(regServicerNum.String),
		})
	}

	return list, nil
}

// 获取按天统计数据-列表
func (s StatisticDateDao) GetStatisticDataList(req *go_micro_srv_cust.GetStatisticDateListRequest) ([]*go_micro_srv_cust.StatisticDateData, int32, error) {
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
	sqlCnt := "SELECT count(1) FROM statistic_date " + whereModel.WhereStr
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

	sqlStr := "SELECT id, day, reg_user_num, reg_servicer_num, create_time, update_time FROM statistic_date " + whereModel.WhereStr
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

	var dataList []*go_micro_srv_cust.StatisticDateData

	for rows.Next() {
		var data go_micro_srv_cust.StatisticDateData
		err = rows.Scan(
			&data.Id,
			&data.Date,
			&data.RegUserNum,
			&data.RegServicerNum,
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

func (r *StatisticDateDao) Insert(wc *DataCount) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	id := strext.GetDailyId()
	return ss_sql.Exec(dbHandler, `INSERT INTO statistic_date(id,day,reg_user_num,reg_servicer_num,create_time,update_time)
	VALUES($1,$2, $3, $4,current_timestamp,current_timestamp)
	ON CONFLICT ON CONSTRAINT unq_day
	DO UPDATE SET reg_user_num = EXCLUDED.reg_user_num,  reg_servicer_num = EXCLUDED.reg_servicer_num,update_time = current_timestamp`,
		id, wc.Day, wc.RegNum, wc.ServerNum)
}
