package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	tmProto "a.a/mp-server/common/proto/tm"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/common/ss_struct"
)

type Writeoff struct {
	Code            string
	IncomeOrderNo   string
	OutGoOrder      string
	TransferOrderNo string
	UseStatus       string
	SendPhone       string
	RecvPhone       string
	CreateTime      string
	DurationTime    string
	SendAccNo       string
	RecvAccNo       string
}

var (
	WriteoffInst Writeoff
)

func (*Writeoff) InitWriteoff(tx *sql.Tx, log *Writeoff) (errCode string) {
	sqlStr := `insert into writeoff(code, income_order_no, outgo_order_no, transfer_order_no, use_status,
		send_phone, recv_phone, create_time, duration_time, send_acc_no, recv_acc_no)
		values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	err := ss_sql.ExecTx(tx, sqlStr, log.Code, log.IncomeOrderNo, log.OutGoOrder, log.TransferOrderNo,
		log.UseStatus, log.SendPhone, log.RecvPhone, log.CreateTime, log.DurationTime, log.SendAccNo, log.RecvAccNo)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}
func (*Writeoff) InitWriteoffV2(tmProxy *ss_struct.TmServerProxy, log *Writeoff) error {

	sql := `insert into writeoff(code, income_order_no, outgo_order_no, transfer_order_no, use_status,
		send_phone, recv_phone, create_time, duration_time, send_acc_no, recv_acc_no)
		values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	args := []string{log.Code, log.IncomeOrderNo, log.OutGoOrder, log.TransferOrderNo, log.UseStatus, log.SendPhone,
		log.RecvPhone, log.CreateTime, log.DurationTime, log.SendAccNo, log.RecvAccNo}
	rsp, err := tmProxy.GetTmServer().TxExec(context.TODO(),
		&tmProto.TxExecRequest{FromServerId: tmProxy.GetFromServerId(), TxNo: tmProxy.GetTxNo(), Sql: sql, Args: args})
	if err != nil {
		return err
	}
	if rsp.Err != "" {
		return errors.New(rsp.Err)
	}
	return nil
}

func (*Writeoff) QueryIncomeOrderNo(code, recvPhone string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var incomeOrderNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select income_order_no from writeoff where recv_phone= $1 and code = $2 and use_status='1' limit 1`,
		[]*sql.NullString{&incomeOrderNo}, recvPhone, code)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return incomeOrderNo.String
}

func (*Writeoff) QueryOrderNo(code, recvPhone string) (*Writeoff, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var incomeOrderNo, outgoOrderNo, transferOrderNo, durationTime sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select income_order_no,outgo_order_no,transfer_order_no,duration_time from writeoff `+
		`where recv_phone= $1 and code = $2 and use_status='1' limit 1`,
		[]*sql.NullString{&incomeOrderNo, &outgoOrderNo, &transferOrderNo, &durationTime}, recvPhone, code)
	if nil != err {
		ss_log.Error("err=%v", err)
		return nil, err
	}

	data := new(Writeoff)
	data.IncomeOrderNo = incomeOrderNo.String
	data.OutGoOrder = outgoOrderNo.String
	data.TransferOrderNo = transferOrderNo.String
	data.DurationTime = durationTime.String

	return data, nil
}
func (*Writeoff) QueryTransferOrderNo(code, recvPhone string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var transferOrderNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select transfer_order_no from writeoff where recv_phone= $1 and code = $2 and use_status='1' limit 1`,
		[]*sql.NullString{&transferOrderNo}, recvPhone, code)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return transferOrderNo.String
}

func (*Writeoff) UpdateWriteoffStatus(tx *sql.Tx, code, status, outGoNo string) (errCode string) {
	sqlStr := `update writeoff set use_status=$1, outgo_order_no=$2, modify_time=current_timestamp, finish_time=current_timestamp  where code=$3`
	err := ss_sql.ExecTx(tx, sqlStr,
		status, outGoNo, code)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_UPDATE
	}

	return ss_err.ERR_SUCCESS
}

func (*Writeoff) GetCode(dbHandler *sql.DB, logNo, orderType string) (code, err string) {
	sqlStr := ""
	switch orderType {
	case constants.VaReason_INCOME:
		sqlStr = "select code from writeoff where income_order_no = $1  "
	case constants.VaReason_OUTGO:
		sqlStr = "select code from writeoff where outgo_order_no = $1  "
	case constants.VaReason_TRANSFER:
		sqlStr = "select code from writeoff where transfer_order_no = $1  "
	default:
		ss_log.Error("orderType 类型错误 [%v]", orderType)
		return "", ss_err.ERR_PARAM
	}

	var codeT sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&codeT}, logNo)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", ss_err.ERR_PARAM
	}

	return codeT.String, ss_err.ERR_SUCCESS
}

func (*Writeoff) GetCodeByIncomeOrderNo(incomeOrderNo string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return "", errors.New("获取数据库连接失败")
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select code from writeoff where income_order_no = $1 "
	var code sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&code}, incomeOrderNo)
	return code.String, err
}

//根据当前时间获取核销码最后有效时间
func (*Writeoff) GetCodeEndTimeStr(nowTimeStr string) (codeEndTimeStr string, err error) {
	//获取核销码有效期(天)
	durationDate := GlobalParamDaoInstance.QeuryParamValue(constants.KEY_WriteOff_DurationDate)
	addHour := time.Hour * 24 * time.Duration(strext.ToInt64(durationDate))
	endTimeStr, errT := ss_time.TimeAfter(nowTimeStr, ss_time.DateTimeSlashFormat, addHour) // 核销码期限时间
	if errT != nil {
		ss_log.Error("获取核销码最后期限时间失败,err=[%v]", errT)
		return "", err
	}
	return endTimeStr, nil
}
