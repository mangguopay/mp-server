package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type UserLoginLogDao struct {
	Uid       string
	LoginTime string
	Result    int32
	DeviceId  string
	Ip        string
	Lat       string
	Lng       string
	Client    string
}

const (
	// 登录结果
	LoginResultSuccess   = 1 // 成功
	LoginResultPassWrong = 2 // 密码错误

	// 登录的客户端
	LoginClientApp = "app" // app端
	LoginClientPos = "pos" // pos端
)

var UserLoginLogDaoInstance UserLoginLogDao

func (*UserLoginLogDao) GetLastByTime(uid string, lastTime string) ([]UserLoginLogDao, error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	sqlStr := "SELECT login_time, result, device_id, ip, lat, lng, client "
	sqlStr += " FROM user_login_log WHERE uid=$1 AND login_time >=$2 ORDER BY login_time DESC"

	rows, stmt, qErr := ss_sql.Query(dbHandler, sqlStr, uid, lastTime)
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

	var list []UserLoginLogDao

	for rows.Next() {
		var loginTime, result, deviceId, ip, lat, lng, client sql.NullString

		err := rows.Scan(&loginTime, &result, &deviceId, &ip, &lat, &lng, &client)
		if err != nil {
			return nil, err
		}
		list = append(list, UserLoginLogDao{
			LoginTime: loginTime.String,
			Result:    strext.ToInt32(result.String),
			DeviceId:  deviceId.String,
			Ip:        ip.String,
			Lat:       lat.String,
			Lng:       lng.String,
			Client:    client.String,
		})
	}

	return list, nil
}
