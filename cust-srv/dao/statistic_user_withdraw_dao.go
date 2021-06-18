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

type StatisticUserWithdrawDao struct{}

var StatisticUserWithdrawDaoInst StatisticUserWithdrawDao

const (
	// 提现类型
	StatisticWithdrawTypeCard     = "1" // 银行卡提现(来自log_to_cust表)
	StatisticRechargeTypeWriteoff = "2" // 核销码取款(来自outgo_order表withdraw_type=0)
	StatisticRechargeTypeScan     = "3" // 扫一扫取款(来自outgo_order表withdraw_type=1 或withdraw_type=2)
)

func (r *StatisticUserWithdrawDao) Insert(wc *DataCount) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	id := strext.GetDailyId()
	return ss_sql.Exec(dbHandler, `INSERT INTO statistic_user_withdraw(id,day,type,ctype,num,amount,fee,create_time,update_time)
			VALUES($1,$2,$3,$4,$5,$6,$7,current_timestamp,current_timestamp)
			ON CONFLICT ON CONSTRAINT unq_withdraw_row
			DO UPDATE SET type = EXCLUDED.type,ctype = EXCLUDED.ctype,num = EXCLUDED.num,
			amount = EXCLUDED.amount,fee = EXCLUDED.fee,update_time = current_timestamp`,
		id, wc.Day, wc.Type, wc.CType, wc.Num, wc.Amount, wc.Fee)
}

type StatisticUserWithdrawData struct {
	Day          string
	CurrencyType string
	TotalNum     int64
	TotalAmount  int64
	TotalFee     int64
}

// 获取提现的统计数据
func (s *StatisticUserWithdrawDao) GetStatisticData(startDate string, endDate string, currencyType string) ([]StatisticUserWithdrawData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	endDate, tErr := ss_time.TimeAfter(endDate, ss_time.DateFormat, time.Hour*24) // 下一天的日期
	if tErr != nil {
		return nil, tErr
	}

	sqlStr := "SELECT day, SUM(num) AS total_num, SUM(amount) AS total_amount, SUM(fee) AS total_fee FROM statistic_user_withdraw "
	sqlStr += " WHERE day >= $1 and day < $2 AND ctype=$3 GROUP BY day ORDER BY day ASC "

	rows, stmt, qErr := ss_sql.Query(dbHandler, sqlStr, startDate, endDate, currencyType)
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

	list := []StatisticUserWithdrawData{}

	for rows.Next() {
		var day, totalNum, totalAmount, totalFee sql.NullString

		err := rows.Scan(&day, &totalNum, &totalAmount, &totalFee)
		if err != nil {
			return nil, err
		}

		list = append(list, StatisticUserWithdrawData{
			Day:          ss_time.StripPostDate(day.String),
			CurrencyType: currencyType,
			TotalNum:     strext.ToInt64(totalNum.String),
			TotalAmount:  strext.ToInt64(totalAmount.String),
			TotalFee:     strext.ToInt64(totalFee.String),
		})
	}

	return list, nil
}

// 获取提现的统计数据-列表
func (s StatisticUserWithdrawDao) GetStatisticDataList(req *go_micro_srv_cust.GetStatisticUserWithdrawListRequest) ([]*go_micro_srv_cust.StatisticUserWithdrawData, int32, error) {
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

	// 提现类型条件
	if req.WithdrawType != "" {
		whereList = append(whereList, &model.WhereSqlCond{Key: "type", Val: req.WithdrawType, EqType: "="})
	}

	// 币种类型条件
	if req.CurrencyType != "" {
		whereList = append(whereList, &model.WhereSqlCond{Key: "ctype", Val: req.CurrencyType, EqType: "="})
	}

	// 初始化where模型
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	//统计
	var totalStr sql.NullString
	sqlCnt := "SELECT count(1) FROM statistic_user_withdraw " + whereModel.WhereStr
	if err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalStr}, whereModel.Args...); err != nil {
		return nil, 0, err
	}

	total := strext.ToInt32(totalStr.String)
	if total == 0 {
		return nil, 0, nil
	}

	// 添加order by 和limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `ORDER BY day DESC, ctype ASC, type ASC `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	sqlStr := "SELECT id, day, type, ctype, num, amount, fee, create_time, update_time FROM statistic_user_withdraw " + whereModel.WhereStr
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

	var dataList []*go_micro_srv_cust.StatisticUserWithdrawData

	for rows.Next() {
		var data go_micro_srv_cust.StatisticUserWithdrawData
		err = rows.Scan(
			&data.Id,
			&data.Date,
			&data.WithdrawType,
			&data.CurrencyType,
			&data.Num,
			&data.Amount,
			&data.Fee,
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
