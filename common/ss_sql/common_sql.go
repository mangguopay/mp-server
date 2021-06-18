package ss_sql

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"database/sql"
	"math/rand"
	"strconv"
	"time"
)

const (
	CARRY_BIT  = 10000
	RATE_BIT   = 2
	AMOUNT_BIT = 2
	UUID       = "00000000-0000-0000-0000-000000000000"

	AUDIT_STATUS_INIT    = "0"
	AUDIT_STATUS_SUCCESS = "1"
	AUDIT_STATUS_WAIT    = "2"
	AUDIT_STATUS_OUT     = "3"

	DB_NO_ROWS_MSG = "sql: no rows in result set"

	CURRENT_TIMESTAMP = "current_timestamp"
	CREATE_TIME       = "create_time"
	DbDuplicateKey    = "pq: duplicate key value violates unique constraint"
)

var (
	isLog = true
)

func BackSpaceTosting(x string, b int) string {
	x64 := strext.ToFloat64(x)
	if b == RATE_BIT {
		return strconv.FormatFloat(x64/strext.ToFloat64(CARRY_BIT), 'f', b, 64)
	} else {
		return strconv.FormatFloat(x64/strext.ToFloat64(CARRY_BIT), 'f', b, 64)
	}
}

func GenerateRangeNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max - min)
	randNum = randNum + min
	return randNum
}

func GetMerchantNumber(number_type int) string {
	return strconv.Itoa(number_type + GenerateRangeNum(10000, 99999))
}

// 组装更新sql
func MkUpdateSql(table string, values map[string]string, where string) (ret string, data []interface{}, idx int, err error) {
	if len(values) <= 0 {
		return "", []interface{}{}, 0, ss_err.ErrNoKeys
	}

	ret = "update " + table + " set "
	idx = 1
	for k, v := range values {
		if v == "" {
			continue
		}

		ret += " \"" + k + "\"=$" + strext.ToStringNoPoint(idx) + ","
		data = append(data, v)
		idx += 1
	}
	return ret[:len(ret)-1] + " " + where, data, idx, nil
}

func MkUpdateSqlMysql(table string, values map[string]string, where string) (ret string, data []interface{}, idx int, err error) {
	if len(values) <= 0 {
		return "", []interface{}{}, 0, ss_err.ErrNoKeys
	}

	ret = "update " + table + " set "
	idx = 1
	for k, v := range values {
		if v == "" {
			continue
		}

		ret += " `" + k + "`=?,"
		data = append(data, v)
		idx += 1
	}
	return ret[:len(ret)-1] + " " + where, data, idx, nil
}

// 组装更新sql
func MkUpdateSqlCustom(table string, values map[string]interface{}, where string) (ret string, data []interface{}, idx int, err error) {
	if len(values) <= 0 {
		return "", []interface{}{}, 0, ss_err.ErrNoKeys
	}

	ret = "update " + table + " set "
	idx = 1
	for k, v := range values {
		if v == "" {
			continue
		}

		ret += " \"" + k + "\"=$" + strext.ToStringNoPoint(idx) + ","
		data = append(data, v)
		idx += 1
	}
	return ret[:len(ret)-1] + " " + where, data, idx, nil
}

func MkUpdateSqlCustomMysql(table string, values map[string]interface{}, where string) (ret string, data []interface{}, idx int, err error) {
	if len(values) <= 0 {
		return "", []interface{}{}, 0, ss_err.ErrNoKeys
	}

	ret = "update " + table + " set "
	idx = 1
	for k, v := range values {
		if v == "" {
			continue
		}

		ret += " `" + k + "`=?,"
		data = append(data, v)
		idx += 1
	}
	return ret[:len(ret)-1] + " " + where, data, idx, nil
}

// 组装插入sql
func MkInsertSql(table string, values map[string]string) (ret string, data []interface{}, idx int, err error) {
	if len(values) <= 0 {
		return "", []interface{}{}, 0, ss_err.ErrNoKeys
	}

	ret = "insert into " + table
	inertSql := ""
	valuesSql := ""
	idx = 1
	for k, v := range values {
		if v == "" {
			continue
		}
		inertSql += k + ","
		if k == CREATE_TIME && v == CURRENT_TIMESTAMP {
			valuesSql += CURRENT_TIMESTAMP + ","
		} else {
			valuesSql += "$" + strext.ToStringNoPoint(idx) + ","
			data = append(data, v)
			idx += 1
		}
	}
	sqlStr := ret + "(" + inertSql[:len(inertSql)-1] + ") values(" + valuesSql[:len(valuesSql)-1] + ")"
	return sqlStr, data, idx, nil
}

func GetPayType(payType string) string {
	var productType string
	switch payType {
	case "101":
		fallthrough
	case "102":
		fallthrough
	case "103":
		fallthrough
	case "104":
		fallthrough
	case "105":
		productType = "微信"
	case "201":
		fallthrough
	case "202":
		fallthrough
	case "203":
		fallthrough
	case "204":
		fallthrough
	case "205":
		productType = "支付宝"
	case "301":
		productType = "信用卡"
	case "311":
		productType = "快捷支付"
	case "302":
		productType = "网管支付"
	case "303":
		productType = "代收"
	case "304":
		productType = "代付"
	case "305":
		productType = "银联钱包"
	case "306":
		productType = "银联扫码"
	case "307":
		productType = "银联h5"
	case "308":
		productType = "银联控件"
	case "309":
		productType = "信用卡代还"
	case "310":
		productType = "网银"

		//case "311":
		//	productType="云闪付"
	case "312":
		productType = "企业网银"
	case "316":
		productType = "银联条码"
	case "401":
		productType = "京东扫码支付"
	case "402":
		productType = "京东h5"
	case "501":
		productType = "蓝牙卡头"
	case "601":
		productType = "qq支付"
	case "602":
		productType = "QQh5支付"
	case "701":
		productType = "后台入金"
	case "801":
		productType = "农行网银"
	case "0":
		productType = "银行卡"
	case "1":
		productType = "社保卡"
	default:
		productType = "银行卡"
	}
	return productType
}

func GetPayTypeDetailed(payType string) string {
	var productType string
	switch payType {
	case "101":
		productType = "微信扫码支付"
	case "102":
		productType = "微信公众号支付"
	case "103":
		productType = "微信H5支付"
	case "104":
		productType = "微信条码支付"
	case "105":
		productType = "微信app支付"
	case "201":
		productType = "支付宝扫码支付"
	case "202":
		productType = "支付宝服务窗支付"
	case "203":
		productType = "支付宝WAP支付"
	case "204":
		productType = "支付宝条码"
	case "205":
		productType = "支付宝H5"
	case "301":
		productType = "信用卡"
	case "311":
		productType = "快捷支付"
	case "302":
		productType = "网管支付"
	case "303":
		productType = "代收"
	case "304":
		productType = "代付"
	case "305":
		productType = "银联钱包"
	case "306":
		productType = "银联扫码"
	case "307":
		productType = "银联h5"
	case "308":
		productType = "银联控件"
	case "309":
		productType = "信用卡代还"
	case "310":
		productType = "网银"

		//case "311":
		//	productType="云闪付"
	case "312":
		productType = "企业网银"
	case "316":
		productType = "银联条码"
	case "401":
		productType = "京东扫码支付"
	case "402":
		productType = "京东h5"
	case "501":
		productType = "蓝牙卡头"
	case "601":
		productType = "qq支付"
	case "602":
		productType = "QQh5支付"
	case "701":
		productType = "后台入金"
	case "801":
		productType = "农行网银"
	case "0":
		productType = "银行卡"
	case "1":
		productType = "社保卡"
	default:
		productType = "银行卡"
	}
	return productType
}

func QueryRowTx(tx *sql.Tx, sqlStr string, retStr []*sql.NullString, args ...interface{}) error {
	return db.QueryRowTx(tx, sqlStr, ss_log.Error, retStr, args...)
}

func QueryRowNTx(tx *sql.Tx, sqlStr string, args ...interface{}) (*sql.Row, *sql.Stmt, error) {
	return db.QueryRowTxN(tx, sqlStr, ss_log.Error, args...)
}

func QueryTx(tx *sql.Tx, sqlStr string, args ...interface{}) (*sql.Rows, *sql.Stmt, error) {
	return db.QueryTx(tx, sqlStr, ss_log.Error, args...)
}

func ExecTx(tx *sql.Tx, sqlStr string, args ...interface{}) error {
	if isLog {
		return db.ExecTx(tx, sqlStr, ss_log.Error, args...)
	}
	return db.ExecTx(tx, sqlStr, nil, args...)
}
func ExecTxWithId(tx *sql.Tx, sqlStr string, args ...interface{}) (idNew string, err error) {
	return db.ExecTxWithId(tx, sqlStr, ss_log.Error, args...)
}

func QueryRow(dbHandler *sql.DB, sqlStr string, retStr []*sql.NullString, args ...interface{}) error {
	return db.QueryRow(dbHandler, sqlStr, ss_log.Error, retStr, args...)
}

func QueryRowN(dbHandler *sql.DB, sqlStr string, args ...interface{}) (*sql.Row, *sql.Stmt, error) {
	return db.QueryRowN(dbHandler, sqlStr, ss_log.Error, args...)
}

func Query(dbHandler *sql.DB, sqlStr string, args ...interface{}) (*sql.Rows, *sql.Stmt, error) {
	return db.Query(dbHandler, sqlStr, ss_log.Error, args...)
}

func Exec(dbHandler *sql.DB, sqlStr string, args ...interface{}) error {
	return db.Exec(dbHandler, sqlStr, ss_log.Error, args...)
}
func ExecWithId(dbHandler *sql.DB, sqlStr string, args ...interface{}) (idNew string, err error) {
	return db.ExecWithId(dbHandler, sqlStr, ss_log.Error, args...)
}
func ExecNoLog(dbHandler *sql.DB, sqlStr string, args ...interface{}) error {
	return db.Exec(dbHandler, sqlStr, nil, args...)
}
func ExecWithResult(dbHandler *sql.DB, sqlStr string, args ...interface{}) (result sql.Result, err error) {
	return db.ExecWtihResult(dbHandler, sqlStr, ss_log.ExecSqlLog, args...)
}
func GetBestPayFileName(paramType string) string {
	var paramName string
	switch paramType {
	case "ID_BACK":
		paramName = "法人身份证反面"
	case "ID":
		paramName = "法人身份证正面"
	case "AGENT_Z_CARD":
		paramName = "代理人身份证正面"
	case "AGENT_F_CARD":
		paramName = "代理人身份证反面"
	case "SOCIETY_CREDIT_LICENSE":
		paramName = "多证合一照"
	case "PROXY_PROVE":
		paramName = "授权委托书"
	case "BANK_OPEN_PROVE":
		paramName = "银行开户证明"
	}
	return paramName
}
func GetBankCardType(paramType string) string {
	var paramName string
	switch paramType {
	case "0":
		paramName = "CREDIT" //信用卡
	case "1":
		paramName = "DEBIT" //借记卡
	case "2":
		paramName = "BA" //银行账户
	case "3":
		paramName = "PB" //存折
	}
	return paramName
}
func GetMerchantApplyType(paramType string) string {
	var paramName string
	switch paramType {
	case "0":
		paramName = "COMPANY" //企业
	case "1":
		paramName = "EVIDENCE" //事业单位
	case "2":
		paramName = "CIVILIAN" //民办非企业单位
	case "3":
		paramName = "SOCIETY" //社会团体
	case "4":
		paramName = "PARTY" //党组织
	case "5":
		paramName = "OVERSEAS_ENT" //海外企业
	}
	return paramName
}
func GetCertificateType(paramType string) string {
	var paramName string
	switch paramType {
	case "0":
		paramName = "LICENSE" //营业执照
	case "1":
		paramName = "MLICENSE" //多证合一营业执照
	case "2":
		paramName = "PARTY_PROVE" //党组织成立证明
	case "3":
		paramName = "SOCIETY" //社证
	case "4":
		paramName = "CIVIL" //民证
	case "5":
		paramName = "EVIDENCE" //事证
	case "6":
		paramName = "CERTOFCORP" //公司注册证书
	}
	return paramName
}
func GetIDType(paramType string) string {
	var paramName string
	switch paramType {
	case "0":
		paramName = "ID" //身份证
	case "1":
		paramName = "PASSPORT" //护照
	case "2":
		paramName = "MTPS" //台胞证
	case "3":
		paramName = "" //港澳通行证

	}
	return paramName
}
func GetContactType(paramType string) string {
	var paramName string
	switch paramType {
	case "0":
		paramName = "LEGAL" //身份证
	case "1":
		paramName = "AGENT" //护照
	}
	return paramName
}

func GetAuditInfoName(paramType string) string {
	var paramName string
	switch paramType {
	case constants.AuditInfoBizLicense:
		paramName = "营业执照"
	case constants.AuditInfoAgentPic:
		paramName = "拓展人与法人合照"
	case constants.AuditInfoBizPlace:
		paramName = "经营场景照"
	case constants.AuditInfoIdCardPic:
		paramName = "法人身份证正面"
	case constants.AuditInfoIdCardPic1:
		paramName = "法人身份证反面"
	case constants.AuditInfoTestVideo:
		paramName = "视频认证"
	case constants.AuditInfoBankCardPic:
		paramName = "结算卡照"
	case constants.AuditInfoOther:
		paramName = "其他"
	case constants.AuditInfoStore:
		paramName = "门店"
	case constants.AuditInfoCashDesk:
		paramName = "其他" //收银台
	case constants.AuditInfoRoomInner:
		paramName = "其他" //室内
	}
	return paramName
}

func BackSpace(x interface{}) string {
	x64 := strext.ToFloat64(x)
	return strconv.FormatFloat(x64/100, 'f', 2, 64)
}

func PackSqlNullString(val string) sql.NullString {
	return sql.NullString{
		String: val,
		Valid:  true,
	}
}

func Commit(tx *sql.Tx) {
	if tx != nil {
		err := tx.Commit()
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
		tx = nil
	}
}

// 这里由于早就把参数传进来，所以tx不会nil的，所以只能无视错误信息
func Rollback(tx *sql.Tx) {
	if tx != nil {
		err := tx.Rollback()
		if err != nil && err != sql.ErrTxDone {
			ss_log.Error("err=[%v]", err)
		}
		tx = nil
	}
}

func BeginTx(dbHandler *sql.DB) *sql.Tx {
	tx, err := dbHandler.Begin()
	if err != nil {
		ss_log.Error("err=[%v]", err)
		Rollback(tx)
		return nil
	}
	return tx
}

func HandleList(dbname, sqlStr string, reqList []interface{}, recvList []interface{}, do func(...interface{}) error) error {
	dbHandler := db.GetDB(dbname)
	defer db.PutDB(dbname, dbHandler)

	rows, stmt, err := Query(dbHandler, sqlStr, reqList...)
	if err != nil {
		ss_log.Error("err=[%v]\n", err)
		return err
	}
	defer func() {
		stmt.Close()
		rows.Close()
	}()

	for rows.Next() {
		err := rows.Scan(recvList...)
		if err != nil {
			ss_log.Error("err=[%v]\n", err)
			continue
		}
		err = do(recvList...)
		if err != nil {
			ss_log.Error("err=[%v]\n", err)
			continue
		}
	}

	return nil
}

func ToggleSqlLog(b bool) {
	isLog = b
}
