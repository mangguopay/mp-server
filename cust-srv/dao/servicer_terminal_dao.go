package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type ServicerTerminal struct {
	TerminalNo     string
	TerminalNumber string
	PosSn          string
	UseStatus      string
	ServicerNo     string
	ServiceAccount string
	CreateTime     string
	UpdateTime     string
}

var ServicerTerminalDao ServicerTerminal

func (ServicerTerminal) Cnt(whereStr string, args []interface{}) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select count(1) " +
		"from servicer_terminal st " +
		"left join servicer s on s.servicer_no = st.servicer_no " +
		"left join account acc on acc.uid = s.account_no " + whereStr

	var total sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&total}, args...); err != nil {
		return "", err
	}
	return total.String, nil

}

func (ServicerTerminal) GetTerminalList(whereStr string, args []interface{}) ([]*ServicerTerminal, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select st.terminal_no, st.terminal_number, st.pos_sn, st.use_status, st.create_time, st.update_time, " +
		"acc.account " +
		"from servicer_terminal st " +
		"left join servicer s on s.servicer_no = st.servicer_no " +
		"left join account acc on acc.uid = s.account_no " + whereStr

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	var dataList []*ServicerTerminal
	for rows.Next() {
		var terminalNo, terminalNumber, posSn, useStatus, serviceAccount, createTime, updateTime sql.NullString
		err := rows.Scan(&terminalNo, &terminalNumber, &posSn, &useStatus, &createTime, &updateTime, &serviceAccount)
		if err != nil {
			return nil, err
		}
		obj := new(ServicerTerminal)
		obj.TerminalNo = terminalNo.String
		obj.TerminalNumber = terminalNumber.String
		obj.PosSn = posSn.String
		obj.UseStatus = useStatus.String
		obj.ServiceAccount = serviceAccount.String
		obj.CreateTime = createTime.String
		obj.UpdateTime = updateTime.String

		dataList = append(dataList, obj)
	}

	return dataList, nil
}

//查看pos属于哪个服务商
func (ServicerTerminal) GetSerPosServicerNoByPosNo(posSn string) (servicerNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select servicer_no from servicer_terminal where pos_sn = $1 and is_delete= $2 and use_status = $3"
	var servicerNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&servicerNoT}, posSn, 0, 1)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return servicerNoT.String

}

func (ServicerTerminal) DeleteTerminalByIdTx(tx *sql.Tx, terminalNo string) error {
	sqlStr := "update servicer_terminal set is_delete='1' where terminal_no = $1 and is_delete = '0' "
	if err := ss_sql.ExecTx(tx, sqlStr, terminalNo); err != nil {
		return err
	}
	return nil
}

func (ServicerTerminal) Insert(servicerNo, terNumber, posSn, useStatus string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	terminalNo := strext.GetDailyId()
	sqlStr := "INSERT INTO servicer_terminal(terminal_no, servicer_no, terminal_number, pos_sn, use_status, create_time) " +
		"VALUES ($1,$2,$3,$4,$5, CURRENT_TIMESTAMP)"
	err := ss_sql.Exec(dbHandler, sqlStr, terminalNo, servicerNo, terNumber, posSn, useStatus)
	if err != nil {
		return "", err
	}

	return terminalNo, nil

}
func (ServicerTerminal) InsertTx(tx *sql.Tx, servicerNo, terNumber, posSn, useStatus string) (string, error) {
	terminalNo := strext.GetDailyId()
	sqlStr := "INSERT INTO servicer_terminal(terminal_no, servicer_no, terminal_number, pos_sn, use_status, create_time) " +
		"VALUES ($1,$2,$3,$4,$5, CURRENT_TIMESTAMP)"
	err := ss_sql.ExecTx(tx, sqlStr, terminalNo, servicerNo, terNumber, posSn, useStatus)
	if err != nil {
		return "", err
	}

	return terminalNo, nil

}

func (ServicerTerminal) UpdateUseStatusTx(tx *sql.Tx, terminalNo, useStatus string) error {
	sqlStr := "update servicer_terminal set use_status = $1 where terminal_no = $2  "
	return ss_sql.ExecTx(tx, sqlStr, useStatus, terminalNo)
}

func (ServicerTerminal) GetTerminalById(terminalNo string) (*ServicerTerminal, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select terminal_no, terminal_number, pos_sn, use_status from servicer_terminal where terminal_no=$1 and is_delete = $2 "

	var terminalId, terminalNumber, posSn, useStatus sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&terminalId, &terminalNumber, &posSn, &useStatus},
		terminalNo, constants.IsDel_No)
	if err != nil {
		return nil, err
	}

	obj := new(ServicerTerminal)
	obj.TerminalNo = terminalId.String
	obj.TerminalNumber = terminalNumber.String
	obj.PosSn = posSn.String
	obj.UseStatus = useStatus.String

	return obj, nil
}

func (ServicerTerminal) GetTerminalByNumber(terminalNumber, useStatus string) ([]*ServicerTerminal, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select terminal_no, terminal_number, pos_sn, use_status " +
		"from servicer_terminal " +
		"where terminal_number=$1 and is_delete = $2 and use_status = $3 "
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, terminalNumber, constants.IsDel_No, useStatus)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var list []*ServicerTerminal
	for rows.Next() {
		var terminalId, terminalNumberT, posSn, useStatus sql.NullString
		if err := rows.Scan(&terminalId, &terminalNumberT, &posSn, &useStatus); err != nil {
			return nil, err
		}

		obj := new(ServicerTerminal)
		obj.TerminalNo = terminalId.String
		obj.TerminalNumber = terminalNumberT.String
		obj.PosSn = posSn.String
		obj.UseStatus = useStatus.String
		list = append(list, obj)
	}
	return list, nil
}

func (ServicerTerminal) GetTerminalUseStatus(terminalNo string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select use_status from servicer_terminal where terminal_no=$1  and is_delete = $2 "

	var useStatus sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&useStatus}, terminalNo, constants.IsDel_No)
	if err != nil {
		return "", err
	}
	return useStatus.String, nil
}
