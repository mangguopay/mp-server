package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	"a.a/mp-server/common/proto/gis"
	"a.a/mp-server/common/ss_sql"
)

type ServiceDao struct {
}

var ServiceDaoInst ServiceDao

func (ServiceDao) GetSrvCoordinate() ([]*go_micro_srv_gis.NearbyServicerData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ser.is_delete", Val: "0", EqType: "="},
		{Key: "ser.use_status", Val: "1", EqType: "="},
	})
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by ser.create_time desc `)
	sqlStr := "SELECT ser.servicer_no, ser.servicer_name, ser.lat, ser.lng " +
		" FROM servicer ser " +
		" LEFT JOIN account acc ON acc.uid = ser.account_no " + whereModel.WhereStr

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		return nil, err
	}
	var datas []*go_micro_srv_gis.NearbyServicerData

	for rows.Next() {
		data := &go_micro_srv_gis.NearbyServicerData{}
		err = rows.Scan(
			&data.ServicerNo,
			&data.ServicerName,
			&data.Lat,
			&data.Lng,
		)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		if data.ServicerName == "" {
			ss_log.Error("服务商网点名称未配置------>ServicerNo:[%v]", data.ServicerNo)
			continue
		}
		if strext.ToStringNoPoint(data.Lng) == "0" && strext.ToStringNoPoint(data.Lat) == "0" {
			ss_log.Error("服务商经纬度未配置------>ServicerNo:[%v]", data.ServicerNo)
			continue
		}
		//与查询的经纬度的大圆距离
		//distance := ss_count.CountCircleDistance(strext.ToFloat64(data.Lat), strext.ToFloat64(data.Lng), strext.ToFloat64(req.Lat), strext.ToFloat64(req.Lng))
		//data.Distance = strext.ToStringNoPoint(distance)
		datas = append(datas, data)
	}
	return datas, nil
}
