package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/api-pos/util"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
	"encoding/json"
	"time"
)

var (
	AccDaoInstance AccDao
)

type AccDao struct {
}

/**
 * 登录记录
 */
func (*AccDao) InsertLoginToken(accNo, routes, token string, isForce int32, ip string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	switch isForce {
	case 1:
		sqlInsert := `insert into login_token(acc_no,routes,token,login_time,ip)values($1,$2,$3,current_timestamp,$4) on conflict(acc_no) do update set routes=$2, token=$3, login_time=current_timestamp,ip=$4 `
		err := ss_sql.Exec(dbHandler, sqlInsert, accNo, routes, token, ip)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return ss_err.ERR_SYS_DB_OP
		}
	default:
		var cnt sql.NullString
		sqlGet := `select count(1) from login_token where acc_no=$1 and login_time < now() - interval '1 H' `
		err := ss_sql.QueryRow(dbHandler, sqlGet, []*sql.NullString{&cnt}, accNo)
		if strext.ToInt64(cnt.String) > 0 {
			sqlDel := `delete from login_token where acc_no=$1`
			err := ss_sql.Exec(dbHandler, sqlDel, accNo)
			if err != nil {
				ss_log.Error("err=[%v]", err)
			}
		}

		sqlGet2 := `select count(1) from login_token where acc_no=$1 `
		err = ss_sql.QueryRow(dbHandler, sqlGet2, []*sql.NullString{&cnt}, accNo)
		if strext.ToInt64(cnt.String) <= 0 {
			sqlInsert := `insert into login_token(acc_no,routes,token,login_time,ip)values($1,$2,$3,current_timestamp,$4) `
			err = ss_sql.Exec(dbHandler, sqlInsert, accNo, routes, token, ip)
			if err != nil {
				// 被别人抢了
				ss_log.Error("err=[%v]", err)
				return ss_err.ERR_ACCOUNT_LOGINED
			}
		} else {
			return ss_err.ERR_ACCOUNT_LOGINED
		}
	}

	return ss_err.ERR_SUCCESS
}
func (*AccDao) InsertPosLoginToken(accNo, routes, token string, isForce int32, ip string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	switch isForce {
	case 1:
		sqlInsert := `insert into pos_login_token(acc_no,routes,token,login_time,ip)values($1,$2,$3,current_timestamp,$4) on conflict(acc_no) do update set routes=$2, token=$3, login_time=current_timestamp,ip=$4 `
		err := ss_sql.Exec(dbHandler, sqlInsert, accNo, routes, token, ip)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return ss_err.ERR_SYS_DB_OP
		}
	default:
		var cnt sql.NullString
		sqlGet := `select count(1) from pos_login_token where acc_no=$1 and login_time < now() - interval '1 H' `
		err := ss_sql.QueryRow(dbHandler, sqlGet, []*sql.NullString{&cnt}, accNo)
		if strext.ToInt64(cnt.String) > 0 {
			sqlDel := `delete from pos_login_token where acc_no=$1`
			err := ss_sql.Exec(dbHandler, sqlDel, accNo)
			if err != nil {
				ss_log.Error("err=[%v]", err)
			}
		}

		sqlGet2 := `select count(1) from pos_login_token where acc_no=$1 `
		err = ss_sql.QueryRow(dbHandler, sqlGet2, []*sql.NullString{&cnt}, accNo)
		if strext.ToInt64(cnt.String) <= 0 {
			sqlInsert := `insert into pos_login_token(acc_no,routes,token,login_time,ip)values($1,$2,$3,current_timestamp,$4) `
			err = ss_sql.Exec(dbHandler, sqlInsert, accNo, routes, token, ip)
			if err != nil {
				// 被别人抢了
				ss_log.Error("err=[%v]", err)
				return ss_err.ERR_ACCOUNT_LOGINED
			}
		} else {
			return ss_err.ERR_ACCOUNT_LOGINED
		}
	}

	return ss_err.ERR_SUCCESS
}

func (*AccDao) DeleteLoginToken(accNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlInsert := `delete from login_token where acc_no=$1 `
	err := ss_sql.Exec(dbHandler, sqlInsert, accNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}
}

func (r *AccDao) GetLoginToken(accNo, loginToken, urlPath, xRealIp string) (isTimeout bool, isReset bool) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var routes, ip, lastOpTime, loginTime sql.NullString
	sqlInsert := `select routes,ip,last_op_time,login_time from login_token where acc_no=$1 and token=$2 `
	err := ss_sql.QueryRow(dbHandler, sqlInsert, []*sql.NullString{&routes, &ip, &lastOpTime, &loginTime}, accNo, loginToken)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return false, false
	}

	routesMap := []map[string]interface{}{}
	err = json.Unmarshal([]byte(routes.String), &routesMap)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return false, false
	}
	isExist := util.IsContainMenu(routesMap, urlPath)
	if ip.String != xRealIp && ip.String != "" {
		ss_log.Error("登录ip改变")
		//return false, false
	}

	isReset = false
	if lastOpTime.String != "" && loginTime.String != "" {
		lastOpTimeT, err := time.Parse("2006-01-02T15:04:05Z", lastOpTime.String)
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
		loginTimeT, err := time.Parse("2006-01-02T15:04:05Z", loginTime.String)
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
		curTimeT := ss_time.Now(global.Tz)
		if curTimeT.After(loginTimeT.Add(55*time.Minute)) && curTimeT.Before(lastOpTimeT.Add(5*time.Minute)) {
			isReset = true
		}
	}

	if isExist {
		r.UpdateLastOpTime(accNo)
	}

	return isExist, isReset
}

func (r *AccDao) GetLoginToken1(accNo, loginToken, xRealIp string) (isReset bool) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var ip, lastOpTime, loginTime sql.NullString
	sqlInsert := `select ip,last_op_time,login_time from login_token where acc_no=$1 and token=$2 `
	err := ss_sql.QueryRow(dbHandler, sqlInsert, []*sql.NullString{&ip, &lastOpTime, &loginTime}, accNo, loginToken)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return false
	}

	if ip.String != xRealIp && ip.String != "" {
		ss_log.Error("登录ip改变")
		return false
	}

	//isReset = false
	//if lastOpTime.String != "" && loginTime.String != "" {
	//	lastOpTimeT, err := time.Parse("2006-01-02T15:04:05Z", lastOpTime.String)
	//	if err != nil {
	//		ss_log.Error("err=[%v]", err)
	//	}
	//	loginTimeT, err := time.Parse("2006-01-02T15:04:05Z", loginTime.String)
	//	if err != nil {
	//		ss_log.Error("err=[%v]", err)
	//	}
	//	curTimeT := util2.Now()
	//	if curTimeT.After(loginTimeT.Add(55*time.Minute)) && curTimeT.Before(lastOpTimeT.Add(5*time.Minute)) {
	//		isReset = true
	//	}
	//}
	return true
}
func (r *AccDao) GetPosLoginToken(accNo, loginToken, xRealIp string) (isReset bool) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var ip, lastOpTime, loginTime sql.NullString
	sqlInsert := `select ip,last_op_time,login_time from pos_login_token where acc_no=$1 and token=$2 `
	err := ss_sql.QueryRow(dbHandler, sqlInsert, []*sql.NullString{&ip, &lastOpTime, &loginTime}, accNo, loginToken)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return false
	}

	if ip.String != xRealIp && ip.String != "" {
		ss_log.Error("登录ip改变")
		//return false
	}

	//isReset = false
	//if lastOpTime.String != "" && loginTime.String != "" {
	//	lastOpTimeT, err := time.Parse("2006-01-02T15:04:05Z", lastOpTime.String)
	//	if err != nil {
	//		ss_log.Error("err=[%v]", err)
	//	}
	//	loginTimeT, err := time.Parse("2006-01-02T15:04:05Z", loginTime.String)
	//	if err != nil {
	//		ss_log.Error("err=[%v]", err)
	//	}
	//	curTimeT := util2.Now()
	//	if curTimeT.After(loginTimeT.Add(55*time.Minute)) && curTimeT.Before(lastOpTimeT.Add(5*time.Minute)) {
	//		isReset = true
	//	}
	//}
	return true
}

func (*AccDao) UpdateLoginRoutes(accNo, routes, ip string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlUpdate := `update login_token set routes=$1,ip=$3 where acc_no=$2 `
	err := ss_sql.Exec(dbHandler, sqlUpdate, routes, accNo, ip)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_OP
	}

	return ss_err.ERR_SUCCESS
}

func (*AccDao) UpdateLastOpTime(accNo string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlUpdate := `update login_token set last_op_time=$1 where acc_no=$2 `
	err := ss_sql.Exec(dbHandler, sqlUpdate, ss_time.NowWithFmt("2006-01-02 15:04:05", global.Tz), accNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_OP
	}

	return ss_err.ERR_SUCCESS
}

func (*AccDao) UpdateLoginTime(accNo string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlUpdate := `update login_token set login_time=$1 where acc_no=$2 `
	err := ss_sql.Exec(dbHandler, sqlUpdate, ss_time.NowWithFmt("2006-01-02 15:04:05", global.Tz), accNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_OP
	}

	return ss_err.ERR_SUCCESS
}
