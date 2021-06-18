package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	"a.a/mp-server/common/ss_sql"
	_ "a.a/mp-server/cust-srv/test"
	"database/sql"
	"testing"
)

func TestWriteoffDao_GetCodeList(t *testing.T) {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "code", Val: "1544826759", EqType: "="},
	})
	got, err := WriteoffDaoInst.GetCodeList(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		t.Errorf("GetCodeList() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(got))
}

func TestWriteoffDao_GetCodeDetailByCode(t *testing.T) {
	code := "3387261752"
	got, err := WriteoffDaoInst.GetCodeDetailByCode(code)
	if err != nil {
		t.Errorf("GetCodeDetailByCode() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(got))
}

func TestWriteoffDao_GetExpiredCodeArr(t *testing.T) {
	gotCodes := WriteoffDaoInst.GetExpiredCodeArr()
	t.Logf("过期核销码：%v", strext.ToJson(gotCodes))

}

func TestCustHandler_InComeOrderFilling(t *testing.T) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr1 := `SELECT io.act_acc_no, io.recv_acc_no, wo.income_order_no
		FROM writeoff wo
		LEFT JOIN income_order io ON io.log_no = wo.income_order_no
		WHERE wo.income_order_no != '' `

	sqlStr2 := `UPDATE writeoff SET send_acc_no = $1, recv_acc_no = $2 WHERE income_order_no = $3 `

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr1)
	if err != nil {
		return
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var i int
	for rows.Next() {
		var sendAccountNo, recvAccountNo, incomeOrderNo sql.NullString
		err := rows.Scan(&sendAccountNo, &recvAccountNo, &incomeOrderNo)
		if err != nil {
			t.Logf("err=%v", err)
			continue
		}
		err = ss_sql.Exec(dbHandler, sqlStr2, sendAccountNo.String, recvAccountNo.String, incomeOrderNo.String)
		if err != nil {
			t.Logf("update err=%v", err)
		}
		i++
	}
	t.Logf("共计填充%v条记录", i)
}

func TestCustHandler_TransferOrderFilling(t *testing.T) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr1 := `SELECT acc1.uid from_account_no, acc2.uid to_account_no, wo.transfer_order_no
		FROM writeoff wo
		LEFT JOIN transfer_order tro ON tro.log_no = wo.transfer_order_no
		LEFT JOIN vaccount vacc1 ON vacc1.vaccount_no = tro.from_vaccount_no
		LEFT JOIN vaccount vacc2 ON vacc2.vaccount_no = tro.to_vaccount_no
		LEFT JOIN account acc1 ON acc1.uid = vacc1.account_no 
		LEFT JOIN account acc2 ON acc2.uid = vacc2.account_no 
		WHERE wo.transfer_order_no != '' `

	sqlStr2 := `UPDATE writeoff SET send_acc_no = $1, recv_acc_no = $2 WHERE transfer_order_no = $3 `

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr1)
	if err != nil {
		return
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var i int
	for rows.Next() {
		var fromAccountNo, toAccountNo, transferOrderNo sql.NullString
		err := rows.Scan(&fromAccountNo, &toAccountNo, &transferOrderNo)
		if err != nil {
			t.Logf("err=%v", err)
			continue
		}
		err = ss_sql.Exec(dbHandler, sqlStr2, fromAccountNo.String, toAccountNo.String, transferOrderNo.String)
		if err != nil {
			t.Logf("update err=%v", err)
		}
		i++
	}
	t.Logf("共计填充%v条记录", i)
}
