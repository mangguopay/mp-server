package handler

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	go_micro_srv_gis "a.a/mp-server/common/proto/gis"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/common"
	"a.a/mp-server/cust-srv/dao"
	"a.a/mp-server/cust-srv/i"
	"a.a/mp-server/cust-srv/util"
	jsoniter "github.com/json-iterator/go"
	"github.com/micro/go-micro/v2"
)

/**
 *
 * 获取ModernPay商户列表
 */
func (*CustHandler) GetServicerList(ctx context.Context, req *go_micro_srv_cust.GetServicerListRequest, reply *go_micro_srv_cust.GetServicerListReply) error {
	if req.StartTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.StartTime) {
			ss_log.Error("参数错误:日期格式错误,StartTime:%s", req.StartTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	if req.EndTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.EndTime) {
			ss_log.Error("参数错误:日期格式错误,EndTime:%s", req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	// 检查结束时间是否大于等于开始时间
	if req.StartTime != "" && req.EndTime != "" {
		if cmp, _ := ss_time.CompareDate("2006/01/02 15:04:05", req.StartTime, req.EndTime); cmp > 0 {
			ss_log.Error("参数错误:开始时间大于结束时间,StartTime:%s,EndTime:%s", req.StartTime, req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	if req.Page < 1 || req.PageSize < 1 {
		ss_log.Error("参数错误:Page:%d,PageSize:%d", req.Page, req.PageSize)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var datas []*go_micro_srv_cust.ServicerData
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ser.servicer_no", Val: req.ServicerNo, EqType: "="},
		{Key: "ser.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "ser.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "ser.is_delete", Val: "0", EqType: "="},
		{Key: "acc.phone", Val: req.QueryPhone, EqType: "like"},
		{Key: "acc.nickname", Val: req.QueryName, EqType: "like"},
		{Key: "acc.account", Val: req.Account, EqType: "like"},
		{Key: "ser.servicer_name", Val: req.ServicerName, EqType: "like"},
	})
	//统计数目
	where := whereModel.WhereStr
	args := whereModel.Args
	var total sql.NullString
	sqlCnt := "SELECT count(1) " +
		"FROM servicer ser " +
		"LEFT JOIN account acc ON acc.uid = ser.account_no " + where
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	//添加limit和排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by ser.create_time desc `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr :=
		"SELECT ser.servicer_no, ser.account_no, ser.addr, ser.create_time" +
			", ser.use_status, ser.commission_sharing, ser.income_sharing, ser.income_authorization" +
			", ser.outgo_authorization, ser.lat, ser.lng, ser.scope, ser.scope_off, ser.servicer_name" +
			", ser.business_time, acc.nickname, acc.phone, acc.account, acc.country_code " +
			", vacc1.balance, vacc2.balance, vacc3.balance, vacc4.balance " +
			" FROM servicer ser " +
			" LEFT JOIN account acc ON acc.uid = ser.account_no " +
			" LEFT JOIN vaccount vacc1 ON vacc1.account_no = ser.account_no and vacc1.va_type = " + strext.ToStringNoPoint(constants.VaType_QUOTA_USD) +
			" LEFT JOIN vaccount vacc2 ON vacc2.account_no = ser.account_no and vacc2.va_type =  " + strext.ToStringNoPoint(constants.VaType_QUOTA_KHR) +
			" LEFT JOIN vaccount vacc3 ON vacc3.account_no = ser.account_no and vacc3.va_type =  " + strext.ToStringNoPoint(constants.VaType_QUOTA_USD_REAL) +
			" LEFT JOIN vaccount vacc4 ON vacc4.account_no = ser.account_no and vacc4.va_type =  " + strext.ToStringNoPoint(constants.VaType_QUOTA_KHR_REAL) +
			" " + where2

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	} else {
		for rows.Next() {
			data := &go_micro_srv_cust.ServicerData{}
			var servicerName, businessTime, countryCode sql.NullString
			var usdAuthCollectLimit, khrAuthCollectLimit sql.NullString
			var usdRealBalance, khrRealBalance sql.NullString
			err = rows.Scan(
				&data.ServicerNo,
				&data.AccountNo,
				&data.Addr,
				&data.CreateTime,
				&data.UseStatus,

				&data.CommissionSharing,
				&data.IncomeSharing,
				&data.IncomeAuthorization,
				&data.OutgoAuthorization,
				&data.Lat,

				&data.Lng,
				&data.Scope,
				&data.ScopeOff,
				&servicerName,
				&businessTime,

				&data.Nickname,
				&data.Phone,
				&data.Account,
				&countryCode,
				&usdAuthCollectLimit,
				&khrAuthCollectLimit,
				&usdRealBalance,
				&khrRealBalance,
			)

			if err != nil {
				ss_log.Error("err=[%v]", err)
				continue
			}
			data.ServicerName = servicerName.String
			data.BusinessTime = businessTime.String
			data.CountryCode = countryCode.String

			if usdAuthCollectLimit.String == "" {
				usdAuthCollectLimit.String = "0"
			}
			if khrAuthCollectLimit.String == "" {
				khrAuthCollectLimit.String = "0"
			}
			if usdRealBalance.String == "" {
				usdRealBalance.String = "0"
			}
			if khrRealBalance.String == "" {
				khrRealBalance.String = "0"
			}

			data.UsdAuthCollectLimit = usdAuthCollectLimit.String
			data.KhrAuthCollectLimit = khrAuthCollectLimit.String

			//查询服务商可用余额
			data.UsdBalance = ss_count.Sub("0", usdRealBalance.String).String()
			data.KhrBalance = ss_count.Sub("0", khrRealBalance.String).String()

			datas = append(datas, data)
		}
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 *
 * app获取附近商户列表
 */
func (*CustHandler) GetNearbyServicerList1(ctx context.Context, req *go_micro_srv_cust.GetNearbyServicerListRequest, reply *go_micro_srv_cust.GetNearbyServicerListReply) error {
	if req.Page < 1 || req.PageSize < 1 {
		ss_log.Error("参数错误:Page:%d,PageSize:%d", req.Page, req.PageSize)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	redisKey := "nearby_servicer_" + req.AccountUid
	//返回拍好序的
	datasT := []*go_micro_srv_cust.NearbyServicerData{}
	if value, _ := cache.RedisClient.Get(redisKey).Result(); value == "" { // 查询数据库并设置进redis
		dbHandler := db.GetDB(constants.DB_CRM)
		defer db.PutDB(constants.DB_CRM, dbHandler)
		whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
			{Key: "ser.is_delete", Val: "0", EqType: "="},
		})

		//添加limit和排序
		ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by ser.create_time desc `)
		//ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

		sqlStr := "SELECT ser.servicer_no, ser.servicer_name, ser.lat, ser.lng " +
			" FROM servicer ser " +
			" LEFT JOIN account acc ON acc.uid = ser.account_no " + whereModel.WhereStr

		rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
		if stmt != nil {
			defer stmt.Close()
		}
		defer rows.Close()

		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_GET
			return nil
		}
		var datas []*go_micro_srv_cust.NearbyServicerData

		for rows.Next() {
			data := &go_micro_srv_cust.NearbyServicerData{}
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
			//与查询的经纬度的大圆距离
			distance := ss_count.CountCircleDistance(strext.ToFloat64(data.Lat), strext.ToFloat64(data.Lng), strext.ToFloat64(req.Lat), strext.ToFloat64(req.Lng))
			data.Distance = strext.ToStringNoPoint(distance)
			datas = append(datas, data)
		}

		for _, data := range datas {
			if len(datasT) == 0 { //第一条直接加进来
				datasT = append(datasT, data)
				continue
			}
			for k, vv := range datasT {
				if strext.ToFloat64(data.Distance) < strext.ToFloat64(vv.Distance) {
					var datasTemp []*go_micro_srv_cust.NearbyServicerData
					datasTemp = append(datasTemp, datasT[:k]...) //当前位置之前的存起来
					datasTemp = append(datasTemp, data)          //在当前位置加上小的数据
					datasTemp = append(datasTemp, datasT[k:]...) // //将原来当前和之后的元素一个一个的添加回来
					datasT = datasTemp
					break
				}
			}
			//如果发现比前面都没有小的数,则直接在后面添加
			datasT = append(datasT, data)
		}

		//todo 设置到redis中去，请求第二页时候就从redis中拿数据。
		datasString := strext.ToJson(datasT)
		if err := cache.RedisClient.Set(redisKey, datasString, constants.CacheKeySecV2).Err(); err != nil {
			ss_log.Error("附近服务商存进redis失败,err=[%v],datas--->%s", err, datasT)
		}

	} else { //取出json转成对象，根据分页，返回适当的数据
		var serviceList []*go_micro_srv_cust.NearbyServicerData
		jsoniter.Unmarshal([]byte(value), &serviceList)
		//fmt.Printf("%+v\n",serviceList)
		startInt := strext.ToInt((req.Page - 1) * req.PageSize)
		endInt := strext.ToInt(req.Page * req.PageSize)
		datasT = serviceList[startInt:endInt]
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datasT
	//reply.Total = strext.ToInt32(total.String)
	return nil
}

// todo 代码不可删,会分离开来的
func (*CustHandler) GetNearbyServicerList(ctx context.Context, req *go_micro_srv_cust.GetNearbyServicerListRequest, reply *go_micro_srv_cust.GetNearbyServicerListReply) error {
	if req.Page < 1 || req.PageSize < 1 {
		ss_log.Error("参数错误:Page:%d,PageSize:%d", req.Page, req.PageSize)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 调用监听服务
	gisResp, _ := i.GisHandleInst.Client.GetNearbyServicerList(ctx, &go_micro_srv_gis.GetNearbyServicerListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Lat:      req.Lat,
		Lng:      req.Lng,
	})

	if gisResp.ResultCode != ss_err.ERR_SUCCESS {
		ss_log.Error("调用gis服务获取列表失败")
		reply.ResultCode = gisResp.ResultCode
		return nil
	}
	datas := make([]*go_micro_srv_cust.NearbyServicerData, 0)
	for _, value := range gisResp.Datas {
		datas = append(datas, &go_micro_srv_cust.NearbyServicerData{
			Lat:          strext.ToStringNoPoint(value.Lat),
			Lng:          strext.ToStringNoPoint(value.Lng),
			ServicerName: value.ServicerName,
			ServicerNo:   value.ServicerNo,
			Distance:     strext.ToStringNoPoint(value.Distance),
		})
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = gisResp.Total
	return nil
}

// 计算和排序周围的服务商坐标
func CountSrvCoordinate(lat, lng float64) []*go_micro_srv_cust.NearbyServicerData {
	serviceList := common.SrvCoordinates

	// 下面为测试是备用的代码
	//value := `[{"lat":"0","lng":"0","servicer_no":"a2d6eeea-9d6c-43c9-ab68-07a3b0370a79","servicer_name":"1"},{"lat":"0","lng":"0","servicer_no":"f2e63de6-4c50-43d5-964f-0a6af137b529","servicer_name":"网点2331"},{"lat":"0","lng":"0","servicer_no":"047087ba-88ba-40ef-a2f0-5c4adefbff5a","servicer_name":"服务商网点59"},{"lat":"0","lng":"0","servicer_no":"06091932-b2c0-4858-9ec3-2b71556d9dd8","servicer_name":"服务商网点55"},{"lat":"0","lng":"0","servicer_no":"28b92402-2ed8-4a75-97ea-582ce4a04444","servicer_name":"服务商网点62"},{"lat":"15.256","lng":"12.1153","servicer_no":"fea8a31f-5d7c-4b8f-ab33-9f4814e9f882","servicer_name":"服务商网点60"},{"lat":"0","lng":"0","servicer_no":"9a6b42a5-eb7f-4ed2-ab95-824faee657e1","servicer_name":"服务商网点56"},{"lat":"0","lng":"0","servicer_no":"064db502-f71c-408b-8ded-7b569c382ec2","servicer_name":"服务商网点49"},{"lat":"0","lng":"0","servicer_no":"7bf6cfed-374c-4084-be2c-fb33d0969a2c","servicer_name":"服务商网点54"},{"lat":"0","lng":"0","servicer_no":"78e67bf3-e43a-4e18-a208-d7216ba23403","servicer_name":"服务商网点58"},{"lat":"13.335438","lng":"103.872697","servicer_no":"26743ab9-d8b0-4a30-9528-2bcd18abd9c9","servicer_name":"服务商网点51"},{"lat":"0","lng":"0","servicer_no":"de24eb4d-4c88-46a1-a8d7-2435a293ff09","servicer_name":"服务商网点57"},{"lat":"22.547899","lng":"113.987186","servicer_no":"60f15170-c1db-41b0-bb3d-14185ab43d28","servicer_name":"服务商网点52"},{"lat":"0","lng":"0","servicer_no":"de9c27c9-c95d-4c88-a98f-7b3a8648d2a0","servicer_name":"服务商网点48"},{"lat":"0","lng":"0","servicer_no":"a74f4ce4-eccb-4446-8f98-f63227b363c6","servicer_name":"服务商网点47"},{"lat":"0","lng":"0","servicer_no":"b385deb9-ca70-4f95-a528-d76cf671bfe5","servicer_name":"服务商网点46"},{"lat":"0","lng":"0","servicer_no":"041eaa21-e4b5-4218-a7e2-df1bb9fcdcf5","servicer_name":"服务商网点45"},{"lat":"0","lng":"0","servicer_no":"f1236f63-2e29-47ef-9205-5af1fdbbe6ff","servicer_name":"服务商网点44"},{"lat":"0","lng":"0","servicer_no":"23f1dea1-0d74-4ed2-ad4a-c43a316fcdd2","servicer_name":"服务商网点43"},{"lat":"0","lng":"0","servicer_no":"68e133b1-4cdb-4c63-b088-4ce179d3183a","servicer_name":"服务商网点42"},{"lat":"11.543454","lng":"104.91116","servicer_no":"fa33f69e-ba91-4824-bb41-db649f8ac07d","servicer_name":"服务商网点61"},{"lat":"0","lng":"0","servicer_no":"96f7443f-2ea5-4eb3-b3e0-ccb65110f4d8","servicer_name":"服务商网点24"},{"lat":"0","lng":"0","servicer_no":"93ac0d7b-e4a8-44a6-893f-aeb992cf0448","servicer_name":"服务商网点12"},{"lat":"0","lng":"0","servicer_no":"5ff8f0e2-b8cd-49a4-be41-3c04d5eed6b4","servicer_name":"服务商网点21"},{"lat":"0","lng":"0","servicer_no":"e9263ec7-52fb-43dc-93af-27f1dcc586ee","servicer_name":"服务商网点20"},{"lat":"0","lng":"0","servicer_no":"fcee7e3e-1232-4800-83dd-b8a11a87b8aa","servicer_name":"服务商网点15"},{"lat":"23.0714778800","lng":"113.3125147400","servicer_no":"9825edb5-6c1c-4ae7-8603-ba754ff105c2","servicer_name":"服务商网点53"},{"lat":"0","lng":"0","servicer_no":"61cb0f07-312f-4e99-a392-37fe5123d0d1","servicer_name":"服务商网点14"},{"lat":"0","lng":"0","servicer_no":"d78905eb-2dd6-4bc7-9b66-f785fd32bb36","servicer_name":"服务商网点13"},{"lat":"0","lng":"0","servicer_no":"071d5a22-d9b3-4cc3-86b0-460efa920cdb","servicer_name":"服务商网点19"},{"lat":"81.87","lng":"29.85","servicer_no":"a7ec67a4-478c-4ea8-9661-05b7d725abb2","servicer_name":"服务商网点16"},{"lat":"69.11","lng":"56.52","servicer_no":"93458335-9317-48dc-a660-e8e617ec9ca9","servicer_name":"服务商网点9"},{"lat":"55.05","lng":"52.77","servicer_no":"3b522b46-6406-4212-86a3-937e43ab8a46","servicer_name":"服务商网点11"},{"lat":"1.02","lng":"4.52","servicer_no":"57be9ac1-ada1-4392-bd0a-f986fc4d2adf","servicer_name":"服务商网点8"},{"lat":"34.45","lng":"25.35","servicer_no":"20c1cf08-f7bf-43ed-8e63-aa635d3247d7","servicer_name":"服务商网点2"},{"lat":"95.16","lng":"47.04","servicer_no":"e357b129-d1a7-421c-9bf3-009e3834b94d","servicer_name":"服务商网点6"},{"lat":"5.1","lng":"1.1","servicer_no":"a2362aa0-1623-44e8-a353-eaf707ebb3a5","servicer_name":"服务商网点1"},{"lat":"51.23","lng":"28.19","servicer_no":"774b9e9d-2551-4ef0-8f07-b338b19a5430","servicer_name":"服务商网点5"},{"lat":"89.99","lng":"63.93","servicer_no":"d30582ae-c410-47f0-8e9c-45d155eb339f","servicer_name":"服务商网点4"},{"lat":"84.99","lng":"36.23","servicer_no":"2e3eeb08-e7be-4101-86be-863286fc0a2c","servicer_name":"服务商网点3"},{"lat":"63.25","lng":"18.16","servicer_no":"02a0c8c0-24df-49ee-a919-f7a4894295f9","servicer_name":"服务商网点18"},{"lat":"1.47","lng":"73.83","servicer_no":"c6780e20-49c9-4bea-985c-c61a94722ab4","servicer_name":"服务商网点32"},{"lat":"35.59","lng":"34.74","servicer_no":"c9739d99-f6d2-49d8-b2dd-16a2199d31d5","servicer_name":"服务商网点23"},{"lat":"66.16","lng":"44.83","servicer_no":"55f62e75-92a0-4e4c-810a-957591c82b43","servicer_name":"服务商网点36"},{"lat":"61.67","lng":"18.40","servicer_no":"961347b6-b226-46df-b627-550445cb58cd","servicer_name":"服务商网点17"},{"lat":"7.83","lng":"56.15","servicer_no":"47d03693-3818-4434-8c93-38551cc96606","servicer_name":"服务商网点26"},{"lat":"23.0714778800","lng":"113.3125147400","servicer_no":"99f047f6-30d0-4e76-853d-57d124c76cc0","servicer_name":"服务商网点50"},{"lat":"8.76","lng":"73.79","servicer_no":"e29bb1c8-9dfc-4de6-a6a8-479eafdb9518","servicer_name":"服务商网点25"},{"lat":"17.37","lng":"49.36","servicer_no":"c50b5b18-f668-4406-8ff2-20990baaba27","servicer_name":"服务商网点39"},{"lat":"19.94","lng":"89.51","servicer_no":"a95d4a11-d620-4e83-aabd-dcb03baf1c2b","servicer_name":"服务商网点7"},{"lat":"84.58","lng":"65.55","servicer_no":"f0e4c436-6acb-4eb2-99f8-961efec221eb","servicer_name":"服务商网点30"},{"lat":"91.74","lng":"64.59","servicer_no":"b153e70c-47ee-4ad2-9e9f-71d60c2f1f44","servicer_name":"服务商网点31"},{"lat":"83.71","lng":"36.32","servicer_no":"88c345c8-23a6-4c99-9ae9-39f592dc15f6","servicer_name":"服务商网点33"},{"lat":"55.15","lng":"66.83","servicer_no":"fb146e3d-4a77-4a01-a178-b9caa9e0a5da","servicer_name":"服务商网点27"},{"lat":"95.81","lng":"51.76","servicer_no":"58dd5aae-5238-4e2b-ab25-0019e9a00906","servicer_name":"服务商网点28"},{"lat":"19.42","lng":"21.37","servicer_no":"048060d7-132d-4b67-8de9-db2208337554","servicer_name":"服务商网点29"},{"lat":"65.96","lng":"45.09","servicer_no":"e7097e0e-4ea1-4d39-9695-6d38bdf7f876","servicer_name":"服务商网点38"},{"lat":"39.92","lng":"16.46","servicer_no":"abcd6641-d40e-45e7-8aaa-0af0663fe947","servicer_name":"服务商网点37"},{"lat":"31.25","lng":"87.63","servicer_no":"bc4713ca-a2bb-4730-a042-ce993921fa69","servicer_name":"服务商网点34"},{"lat":"117.10","lng":"39.10","servicer_no":"03707fce-2a62-40ce-b0f4-ee2e4e3d8121","servicer_name":"服务商网点41"},{"lat":"28.34","lng":"73.09","servicer_no":"5e8b51b5-574d-469f-835b-721532dd25e0","servicer_name":"服务商网点35"},{"lat":"94.35","lng":"53.98","servicer_no":"43a8ddda-0874-4c26-a6be-cf1a62aa7696","servicer_name":"服务商网点22"},{"lat":"74.55","lng":"85.09","servicer_no":"8ab67d7b-fea2-4643-a612-8ec8f5ead188","servicer_name":"服务商网点40"}]`
	//var serviceList1 []*go_micro_srv_cust.NearbyServicerData
	//jsoniter.Unmarshal([]byte(value), &serviceList1)
	//serviceList = serviceList1
	if len(serviceList) == 0 {
		ss_log.Error("CountSrvCoordinate, 服务商坐标列表为空")
		return serviceList
	}

	for _, data := range serviceList {
		//与查询的经纬度的大圆距离
		distance := ss_count.CountCircleDistance(strext.ToFloat64(data.Lat), strext.ToFloat64(data.Lng), strext.ToFloat64(lat), strext.ToFloat64(lng))
		data.Distance = strext.ToStringNoPoint(distance)
	}
	bf, _ := jsoniter.MarshalToString(serviceList)
	fmt.Println("排序前--------->", bf)
	//sort.Sort()
	quickSort(serviceList, 0, len(serviceList)-1)
	af, _ := jsoniter.MarshalToString(serviceList)
	fmt.Println("排序后--------->", af)
	return serviceList
}

// 快速排序
func quickSort(values []*go_micro_srv_cust.NearbyServicerData, left int, right int) {
	if left < right {
		// 设置分水岭
		temp := values[left]
		// 设置哨兵
		i, j := left, right
		for {
			// 从右向左找，找到第一个比分水岭小的数
			for strext.ToFloat64(values[j].Distance) >= strext.ToFloat64(temp.Distance) && i < j {
				j--
			}
			// 从左向右找，找到第一个比分水岭大的数
			for strext.ToFloat64(values[i].Distance) <= strext.ToFloat64(temp.Distance) && i < j {
				i++
			}
			// 如果哨兵相遇，则退出循环
			if i >= j {
				break
			}
			// 交换左右两侧的值
			values[i], values[j] = values[j], values[i]
		}
		// 将分水岭移到哨兵相遇点
		values[left] = values[i]
		values[i] = temp
		// 递归，左右两侧分别排序
		quickSort(values, left, i-1)
		quickSort(values, i+1, right)
	}
}

func (c *CustHandler) GetServicerPhoneList(ctx context.Context, req *go_micro_srv_cust.GetServicerPhoneListRequest, reply *go_micro_srv_cust.GetServicerPhoneListReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var datas []*go_micro_srv_cust.ServicerPhoneData
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "acc.is_delete", Val: "0", EqType: "="},
		{Key: "acc.phone", Val: req.QueryPhone, EqType: "like"},
		{Key: "rai.account_type", Val: constants.AccountType_SERVICER, EqType: "="},
	})

	// 排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by acc.create_time desc `)
	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "select acc.phone,acc.uid,rai.iden_no from account acc LEFT JOIN rela_acc_iden rai ON acc.uid = rai.account_no " + where2

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	} else {
		for rows.Next() {
			data := &go_micro_srv_cust.ServicerPhoneData{}
			err = rows.Scan(
				&data.Phone,
				&data.AccountNo,
				&data.ServicerNo,
			)

			datas = append(datas, data)
		}
	}
	reply.DataList = datas
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
获取指定商户信息
*/
func (c *CustHandler) GetServicerInfo(ctx context.Context, req *go_micro_srv_cust.GetServicerInfoRequest, reply *go_micro_srv_cust.GetServicerInfoReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	//获取服务商信息
	servicerData, servicerErr := dao.ServiceDaoInst.GetServicerInfo(dbHandler, req.ServicerNo)
	if servicerErr != ss_err.ERR_SUCCESS {
		ss_log.Error("获取服务商信息失败,err=[%v]", servicerErr)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	servicerCardDatas := []*go_micro_srv_cust.ServicerCardPackData{}
	//获取服务商对应的账号uid
	accountNo, GetAccountErr := dao.ServiceDaoInst.GetAccountNoByServicerNo(req.ServicerNo)
	if GetAccountErr != ss_err.ERR_SUCCESS {
		ss_log.Error("GetAccountErr=[%v]", GetAccountErr)
	} else {
		datas, getServicerCardErr := dao.ServiceDaoInst.GetServicerCards(accountNo)
		if getServicerCardErr != ss_err.ERR_SUCCESS {
			ss_log.Error("getServicerCardErr=[%v]", getServicerCardErr)
		}
		servicerCardDatas = datas
	}

	ids, errIds := dao.DictimagesDaoInst.GetImgIds(req.ServicerNo)
	if errIds != nil && errIds != sql.ErrNoRows {
		ss_log.Error("查询服务商[%v]的旧营业执照、营业场景失败", req.ServicerNo)
	}

	imgIds := strings.Split(ids, ",")

	imgUrl := dao.ImageDaoInstance.GetImgUrlsByImgIds(imgIds)
	servicerImgData := &go_micro_srv_cust.ServicerImgData{}
	if len(imgUrl) < 4 {
		imgUrl = append(imgUrl, "")
		imgUrl = append(imgUrl, "")
		imgUrl = append(imgUrl, "")
		imgUrl = append(imgUrl, "")
	}

	servicerImgData.ServicerImg1 = imgUrl[0]
	servicerImgData.ServicerImg2 = imgUrl[1]
	servicerImgData.ServicerImg3 = imgUrl[2]
	servicerImgData.ServicerImg4 = imgUrl[3]

	//获取服务商的终端信息
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "servicer_no", Val: req.ServicerNo, EqType: "="},
		{Key: "is_delete", Val: "0", EqType: "="},
	})
	servicerTerminalDatas, getServicerTerminalErr := dao.ServiceDaoInst.GetServicerTerminal(whereModel.WhereStr, whereModel.Args)
	if getServicerTerminalErr != ss_err.ERR_SUCCESS {
		ss_log.Error("getServicerTerminalErr=[%v]", getServicerTerminalErr)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = servicerData
	reply.CardList = servicerCardDatas
	reply.ServicerImgData = servicerImgData
	reply.ServicerTerminalDataList = servicerTerminalDatas
	return nil
}

/**
app获取指定商户信息
*/
func (*CustHandler) GetServicer(ctx context.Context, req *go_micro_srv_cust.GetServicerRequest, reply *go_micro_srv_cust.GetServicerReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var data go_micro_srv_cust.ServicerData

	servicerNo := ""
	if req.IdenNo != "" {
		servicerNo = req.IdenNo
	}

	if req.ServicerNo != "" {
		servicerNo = req.ServicerNo
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ser.servicer_no", Val: servicerNo, EqType: "="},
	})

	sqlStr := "SELECT ser.lat,ser.lng,ser.servicer_no,ser.contact_person,ser.contact_phone,ser.contact_addr,ser.addr,ser.servicer_name,ser.business_time" +
		" FROM servicer ser " + whereModel.WhereStr

	rows, stmt, err := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}

	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	var servicerName, businessTime, latT, lngT sql.NullString
	err = rows.Scan(
		&latT,
		&lngT,
		&data.ServicerNo,
		&data.ContactPerson,
		&data.ContactPhone,
		&data.ContactAddr,
		&data.Addr,
		&servicerName,
		&businessTime,
	)

	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}
	if servicerName.String != "" {
		data.ServicerName = servicerName.String
	}
	if businessTime.String != "" {
		data.BusinessTime = businessTime.String
	}
	if latT.String != "" {
		data.Lat = latT.String
	}
	if lngT.String != "" {
		data.Lng = lngT.String
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = &data
	return nil
}

/**
 * 获取ModernPay指定商户统计信息
 */
func (*CustHandler) GetServicerOrderCount(ctx context.Context, req *go_micro_srv_cust.GetServicerOrderCountRequest, reply *go_micro_srv_cust.GetServicerOrderCountReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "servicer_no", Val: req.ServicerNo, EqType: "="},
	})

	usdModel := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhere(usdModel, "currency_type", "usd", "=")
	usdData, errA := dao.ServiceDaoInst.GetServicerOrderCountDetail(usdModel.WhereStr, usdModel.Args)
	if errA != ss_err.ERR_SUCCESS {
		ss_log.Error("errA=[%v]", errA)
	}

	khrModel := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhere(khrModel, "currency_type", "khr", "=")
	khrData, errB := dao.ServiceDaoInst.GetServicerOrderCountDetail(khrModel.WhereStr, khrModel.Args)
	if errB != ss_err.ERR_SUCCESS {
		ss_log.Error("errB=[%v]", errB)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.UsdData = usdData
	reply.KhrData = khrData
	return nil
}

/**
 * 修改ModernPay指定商户状态
 */
func (*CustHandler) ModifyServicerStatus(ctx context.Context, req *go_micro_srv_cust.ModifyServicerStatusRequest, reply *go_micro_srv_cust.ModifyServicerStatusReply) error {
	//参数校验
	useStatusStr, legal1 := util.GetParamZhCn(req.UseStatus, util.UseStatus)
	if !legal1 {
		ss_log.Error("UseStatus %v", useStatusStr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	oldUseStatus, err1 := dao.ServiceDaoInst.GetStatusBySerNo(req.ServicerNo)
	if err1 != nil {
		ss_log.Error("查询旧数据出错，servicerNo[%v]", req.ServicerNo)
		return nil
	}
	if oldUseStatus == req.UseStatus {
		ss_log.Error("要设置的UseStatus值与旧的UseStatus值相同")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if errStr := dao.ServiceDaoInst.ModifyServiceStatus(tx, req.UseStatus, req.ServicerNo); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
	}

	account, errGet := dao.ServiceDaoInst.GetAccountBySerNo(req.ServicerNo)
	if errGet != nil {
		ss_log.Error("查询服务商的账号失败")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	oldUseStatus, _ = util.GetParamZhCn(oldUseStatus, util.UseStatus)
	description := fmt.Sprintf("修改服务商账号[%v],将商户状态[%v]更改为[%v]", account, oldUseStatus, useStatusStr)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Servicer)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
 * 修改ModernPay指定商户配置
 */
func (*CustHandler) ModifyServicerConfig(ctx context.Context, req *go_micro_srv_cust.ModifyServicerConfigRequest, reply *go_micro_srv_cust.ModifyServicerConfigReply) error {
	//参数校验
	str1, legal1 := util.GetParamZhCn(req.IncomeAuthorization, util.IncomeAuthorization)
	if !legal1 {
		ss_log.Error("IncomeAuthorization %v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	str2, legal2 := util.GetParamZhCn(req.OutgoAuthorization, util.OutgoAuthorization)
	if !legal2 {
		ss_log.Error("OutgoAuthorization %v", str2)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	if errStr := dao.ServiceDaoInst.ModifyServicerConfig(tx, req.ServicerNo, req.IncomeAuthorization, req.OutgoAuthorization, req.CommissionSharing, req.IncomeSharing,
		req.Lat, req.Lng, req.Scope, req.ScopeOff, req.ServicerName,
		req.BusinessTime, req.Addr); errStr != ss_err.ERR_SUCCESS {
		ss_log.Error("更新服务商[%v]信息失败", req.ServicerNo)
		reply.ResultCode = errStr
		return nil
	}

	//获取服务商的账号id
	accountNo, errGet := dao.ServiceDaoInst.GetAccountNoByServicerNo(req.ServicerNo)
	if errGet != ss_err.ERR_SUCCESS {
		ss_log.Error("获取服务商[%v]的账号失败", req.ServicerNo)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	errUp := dao.ServiceDaoInst.ModifySerAuthCollectLimit(tx, accountNo, "usd", req.UsdAuthCollectLimit)
	if errUp != ss_err.ERR_SUCCESS {
		ss_log.Error("更新服务商[%v]授权收款usd额度失败", req.ServicerNo)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	errUp = dao.ServiceDaoInst.ModifySerAuthCollectLimit(tx, accountNo, "khr", req.KhrAuthCollectLimit)
	if errUp != ss_err.ERR_SUCCESS {
		ss_log.Error("更新服务商[%v]授权收款khr额度失败", req.ServicerNo)
		//更新服务商授权收款额度
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	account, errGet1 := dao.ServiceDaoInst.GetAccountBySerNo(req.ServicerNo)
	if errGet1 != nil {
		ss_log.Error("获取服务商的账号失败")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	commissionSharingStr := ss_count.Div(req.CommissionSharing, "100").String() + "%"
	incomeSharingStr := ss_count.Div(req.IncomeSharing, "100").String() + "%"
	usdAuthCollectLimitStr := ss_count.Div(req.UsdAuthCollectLimit, "100").String()

	//插入关键操作日志的描述
	description := fmt.Sprintf("修改服务商账号[%v]的收款权限为[%v],取款权限[%v],取款手续费分成:[%v],存款手续费分成:[%v]", account, str1, str2, commissionSharingStr, incomeSharingStr)
	description = fmt.Sprintf("%v,经度:[%v],纬度:[%v],服务范围:[%v],围栏开关:[%v]", description, req.Lng, req.Lat, req.Scope, req.ScopeOff)
	description = fmt.Sprintf("%v,网点名称:[%v],网点地址:[%v],营业时间:[%v]", description, req.ServicerName, req.Addr, req.BusinessTime)
	description = fmt.Sprintf("%v,服务商授权收款usd额度:[%v],服务商授权收款khr额度:[%v]", description, usdAuthCollectLimitStr, req.KhrAuthCollectLimit)

	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Account)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	tx.Commit()

	//=============================================================
	//  推送同步数据的时间
	if req.Lat != "" && req.Lng != "" {
		sendSrvGis(constants.Topic_Event_Srv_Gis, common.SrvGisPub)
		ss_log.Info("更新了服务商坐标信息,推送事件成功")
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func sendSrvGis(topic string, p micro.Publisher) {
	// create new event
	ev := &go_micro_srv_gis.ListenEvenRequest{
		IsSync: true,
	}
	ss_log.Info("publishing %+v\n", ev)
	// publish an event
	if err := p.Publish(context.TODO(), ev); err != nil {
		ss_log.Error("error publishing: %v", err)
	}
}

/**
 * 获取服务商交易信息列表
 */
func (*CustHandler) GetServiceTransactions(ctx context.Context, req *go_micro_srv_cust.GetServiceTransactionsRequest, reply *go_micro_srv_cust.GetServiceTransactionsReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	if req.StartTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.StartTime) {
			ss_log.Error("GetServiceTransactions StartTime格式不正确,StartTime: %s", req.StartTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}
	if req.EndTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.EndTime) {
			ss_log.Error("GetServiceTransactions EndTime格式不正确,StartTime: %s", req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	var datas []*go_micro_srv_cust.ServiceTransactionsData

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "acc.nickname", Val: req.Nickname, EqType: "like"},
		{Key: "acc.phone", Val: req.Phone, EqType: "like"},
		{Key: "acc.account", Val: req.Account, EqType: "like"},
		{Key: "bdr.order_no", Val: req.OrderNo, EqType: "like"},
		{Key: "bdr.currency_type", Val: req.CurrencyType, EqType: "="},
		{Key: "bdr.bill_type", Val: "('" + constants.BILL_TYPE_INCOME + "','" + constants.BILL_TYPE_OUTGO + "')", EqType: "in"},
		{Key: "bdr.bill_type", Val: req.BillType, EqType: "="},
		{Key: "bdr.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "bdr.create_time", Val: req.EndTime, EqType: "<="},
	})
	whereCnt := whereModel.WhereStr
	argsCnt := whereModel.Args

	sqlCnt := "select count(1) " +
		" from billing_details_results bdr " +
		" left join account acc on acc.uid = bdr.account_no " + whereCnt

	var total sql.NullString
	cntErr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, argsCnt...)
	if cntErr != nil {
		ss_log.Error("cntErr=[%v]", cntErr)
	}
	//添加limit和排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by bdr.create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	where := whereModel.WhereStr
	args := whereModel.Args
	sqlStr := "SELECT bdr.order_no, acc.nickname, acc.phone, acc.account, bdr.amount, bdr.currency_type, bdr.bill_type, bdr.create_time " +
		" from billing_details_results bdr " +
		" left join account acc on acc.uid = bdr.account_no " + where
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("err=[%v]\nreq=[%v]\nsql=[%v]", err, req, sqlStr)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}
	for rows.Next() {
		var data go_micro_srv_cust.ServiceTransactionsData
		if err = rows.Scan(
			&data.OrderNo,
			&data.Nickname,
			&data.Phone,
			&data.Account,
			&data.Amount,

			&data.CurrencyType,
			&data.BillType,
			&data.CreateTime,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		datas = append(datas, &data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

/**
 * 获取服务商收益信息列表
 */
func (*CustHandler) GetServicerProfitLedgerList(ctx context.Context, req *go_micro_srv_cust.GetServicerProfitLedgerListRequest, reply *go_micro_srv_cust.GetServicerProfitLedgerListReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	if req.StartTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.StartTime) {
			ss_log.Error("GetServicerProfitLedgerList StartTime格式不正确,StartTime: %s", req.StartTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}
	if req.EndTime != "" {
		if !ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, req.EndTime) {
			ss_log.Error("GetServicerProfitLedgerList EndTime格式不正确,StartTime: %s", req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	var datas []*go_micro_srv_cust.ServicerProfitLedgerData

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "acc.nickname", Val: req.Nickname, EqType: "like"},
		{Key: "acc.phone", Val: req.Phone, EqType: "like"},
		{Key: "acc.account", Val: req.Account, EqType: "like"},
		{Key: "spl.currency_type", Val: req.CurrencyType, EqType: "="},
		{Key: "spl.payment_time", Val: req.StartTime, EqType: ">="},
		{Key: "spl.payment_time", Val: req.EndTime, EqType: "<="},
		{Key: "spl.log_no", Val: req.LogNo, EqType: "like"},
		{Key: "spl.order_type", Val: req.OrderType, EqType: "="},
	})
	whereCnt := whereModel.WhereStr
	argsCnt := whereModel.Args

	sqlCnt := "select count(1) " +
		"from servicer_profit_ledger spl " +
		" left join servicer ser on ser.servicer_no = spl.servicer_no " +
		" left join account acc on acc.uid = ser.account_no " + whereCnt
	var total sql.NullString
	cntErr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, argsCnt...)
	if cntErr != nil {
		ss_log.Error("cntErr=[%v]", cntErr)
	}

	//统计usd、khr总数与统计
	var usdCount, usdSum, khrCount, khrSum sql.NullString
	usdM := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhere(usdM, "spl.currency_type", "usd", "=")
	sqlUsd := "select count(1),sum(spl.actual_income) " +
		" from servicer_profit_ledger spl" +
		" left join servicer ser on ser.servicer_no = spl.servicer_no " +
		" left join account acc on acc.uid = ser.account_no " + usdM.WhereStr
	usdErr := ss_sql.QueryRow(dbHandler, sqlUsd, []*sql.NullString{&usdCount, &usdSum}, usdM.Args...)
	if usdErr != nil {
		ss_log.Error("cntErr=[%v]", usdErr)
	}

	khrM := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhere(khrM, "spl.currency_type", "khr", "=")
	sqlKhr := "select count(1),sum(spl.actual_income) " +
		" from servicer_profit_ledger spl " +
		" left join servicer ser on ser.servicer_no = spl.servicer_no " +
		" left join account acc on acc.uid = ser.account_no " + khrM.WhereStr
	khrErr := ss_sql.QueryRow(dbHandler, sqlKhr, []*sql.NullString{&khrCount, &khrSum}, khrM.Args...)
	if khrErr != nil {
		ss_log.Error("cntErr=[%v]", khrErr)
	}

	//添加limit和排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by spl.payment_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	where := whereModel.WhereStr
	args := whereModel.Args
	sqlStr := "SELECT spl.log_no, spl.amount_order, spl.servicefee_amount_sum, spl.split_proportion" +
		", spl.actual_income, spl.payment_time, spl.servicer_no, spl.currency_type, spl.order_type " +
		", acc.nickname, acc.phone, acc.account " +
		" from servicer_profit_ledger spl" +
		" left join servicer ser on ser.servicer_no = spl.servicer_no " +
		" left join account acc on acc.uid = ser.account_no " + where
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}
	for rows.Next() {
		var data go_micro_srv_cust.ServicerProfitLedgerData
		var orderType sql.NullString
		if err = rows.Scan(
			&data.LogNo,
			&data.AmountOrder,
			&data.ServicefeeAmountSum,
			&data.SplitProportion,
			&data.ActualIncome,

			&data.PaymentTime,
			&data.ServicerNo,
			&data.CurrencyType,
			&orderType,

			&data.Nickname,
			&data.Phone,
			&data.Account,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		if orderType.String != "" {
			data.OrderType = orderType.String
		}

		datas = append(datas, &data)
	}

	reply.CountData = &go_micro_srv_cust.ServicerProfitLedgerCountData{}
	if usdCount.String == "" {
		reply.CountData.UsdCount = "0"
	} else {
		reply.CountData.UsdCount = usdCount.String
	}

	if usdSum.String == "" {
		reply.CountData.UsdSum = "0"
	} else {
		reply.CountData.UsdSum = usdSum.String
	}

	if khrCount.String == "" {
		reply.CountData.KhrCount = "0"
	} else {
		reply.CountData.KhrCount = khrCount.String
	}

	if khrSum.String == "" {
		reply.CountData.KhrSum = "0"
	} else {
		reply.CountData.KhrSum = khrSum.String
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

func (*CustHandler) GetServicerAccounts(ctx context.Context, req *go_micro_srv_cust.GetServicerAccountsRequest, reply *go_micro_srv_cust.GetServicerAccountsReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "acc.account", Val: req.Account, EqType: "like"},
	}

	total, errCnt := dao.ServiceDaoInst.GetServicerAccountCnt(whereList)
	if errCnt != nil {
		ss_log.Error("查询服务商账户数量出错,err[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	datas, err := dao.ServiceDaoInst.GetServicerAccount(whereList, req.Page, req.PageSize, req.SortType)
	if err != nil {
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	return nil
}

/**
 * WEB获取指定服务商账单明细列表
 */
func (*CustHandler) GetServicerBills(ctx context.Context, req *go_micro_srv_cust.GetServicerBillsRequest, reply *go_micro_srv_cust.GetServicerBillsReply) error {

	if req.StartTime != "" {
		if !ss_time.CheckDateTimeIsRight("2006/01/02 15:04:05", req.StartTime) {
			ss_log.Error("参数错误:日期格式错误,StartTime:%s", req.StartTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	if req.EndTime != "" {
		if !ss_time.CheckDateTimeIsRight("2006/01/02 15:04:05", req.EndTime) {
			ss_log.Error("参数错误:日期格式错误,EndTime:%s", req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	// 检查结束时间是否大于等于开始时间
	if req.StartTime != "" && req.EndTime != "" {
		if cmp, _ := ss_time.CompareDate("2006/01/02 15:04:05", req.StartTime, req.EndTime); cmp > 0 {
			ss_log.Error("参数错误:开始时间大于结束时间,StartTime:%s,EndTime:%s", req.StartTime, req.EndTime)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	if req.Page < 1 || req.PageSize < 1 {
		ss_log.Error("参数错误:Page:%d,PageSize:%d", req.Page, req.PageSize)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "lv.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "lv.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "lv.biz_log_no", Val: req.BizLogNo, EqType: "like"},
		{Key: "vacc.balance_type", Val: req.BalanceType, EqType: "="},
		{Key: "lv.reason", Val: req.Reason, EqType: "="},
		{Key: "vacc.account_no", Val: req.Uid, EqType: "="},
	})
	//只要服务商的
	strs := " and vacc.va_type in (" +
		"'" + strext.ToStringNoPoint(constants.VaType_QUOTA_USD_REAL) + "'" +
		",'" + strext.ToStringNoPoint(constants.VaType_QUOTA_KHR_REAL) + "')"
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, strs)

	extraStr := getShowServicerBillExtraStr()
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, extraStr)

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	total := dao.LogVaccountDaoInst.GetCnt(whereModel.WhereStr, whereModel.Args)

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by lv.create_time desc,case reason when "+constants.VaReason_FEES+" then 1 end")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.LogVaccountDaoInst.GetServicerBills(dbHandler, whereModel.WhereStr, whereModel.Args, req.Uid)
	if err != nil {
		ss_log.Error("查询服务商账户uid:[%v]的虚帐日志失败", req.Uid)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	return nil
}

//获取服务商账单要显示的白名单
func getShowServicerBillExtraStr() string {

	//黑名单
	//银行卡提现成功、手续费(9-4、6-4)

	//白名单
	extraStr := " AND ( "

	//手续费(6-2、6-3、6-6)
	//extraStr += " ( lv.reason = '" + constants.VaReason_FEES + "' AND lv.op_type in( '" + constants.VaOpType_Minus + "','" + constants.VaOpType_Freeze + "','" + constants.VaOpType_Defreeze_Add + "' )) "

	//用户成功存款(2-6)
	extraStr += "  ( lv.reason = '" + constants.VaReason_INCOME + "' AND lv.op_type = '" + constants.VaOpType_Defreeze_Add + "' )"

	//用户成功取款、手续费(3-2、6-2)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_OUTGO + "' AND lv.op_type = '" + constants.VaOpType_Minus + "' )"
	extraStr += " OR ( lv.reason = '" + constants.VaReason_FEES + "' AND lv.op_type = '" + constants.VaOpType_Minus + "' )"

	//平台修改余额(23-1,23-2)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_ChangeSrvBalance + "' AND lv.op_type = '" + constants.VaOpType_Add + "' )"
	extraStr += " OR ( lv.reason = '" + constants.VaReason_ChangeSrvBalance + "' AND lv.op_type = '" + constants.VaOpType_Minus + "' )"

	//服务商充值成功(12-1)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_Srv_Save + "' AND lv.op_type = '" + constants.VaOpType_Add + "' )"

	//服务商提现申请(13-8)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_Srv_Withdraw + "' AND lv.op_type = '" + constants.VaOpType_Balance_Frozen_Add + "' )"
	//服务商提现驳回(13-9)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_Srv_Withdraw + "' AND lv.op_type = '" + constants.VaOpType_Balance_Defreeze_Add + "' )"
	//服务商提现成功(13-4)
	extraStr += " OR ( lv.reason = '" + constants.VaReason_Srv_Withdraw + "' AND lv.op_type = '" + constants.VaOpType_Defreeze + "' )"

	extraStr += " )"

	return extraStr
}

func (*CustHandler) GetSerPos(ctx context.Context, req *go_micro_srv_cust.GetSerPosRequest, reply *go_micro_srv_cust.GetSerPosReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	servicerNo := ""
	switch req.AccountType {
	case constants.AccountType_SERVICER: //服务商
		//查询账号对应的服务商id
		servicerNoT := dao.ServiceDaoInst.GetServicerNoByAccNo(req.AccountUid)
		if servicerNoT == "" {
			ss_log.Error("服务商账号查询服务商id出错,服务商账号uid---------------->%s", req.AccountUid)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		servicerNo = servicerNoT
	case constants.AccountType_POS: // 收银员
		// 获取店员no,再获取服务商no
		cashierNo := dao.RelaAccIdenDaoInst.GetIdenFromAcc(req.AccountUid, constants.AccountType_POS)
		if cashierNo == "" {
			ss_log.Error("店员账号查询店员id出错,店员账号uid---------------->%s", req.AccountUid)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		err, servicerNoT := dao.ServiceDaoInst.GetServiceByCashierNo(cashierNo)
		if err != nil {
			ss_log.Error("店员id查询服务商id出错,店员账号uid---------------->%s", req.AccountUid)
			reply.ResultCode = ss_err.ERR_DB_OP_SER
			return nil
		}
		servicerNo = servicerNoT
	default:
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "servicer_no", Val: servicerNo, EqType: "="},
		{Key: "is_delete", Val: "0", EqType: "="},
	})

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by use_status desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err2 := dao.ServiceDaoInst.GetServicerTerminal(whereModel.WhereStr, whereModel.Args)
	if err2 != ss_err.ERR_SUCCESS {
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	//reply.Total = strext.ToInt32(total)
	return nil
}

func (*CustHandler) ModifySerPosStatus(ctx context.Context, req *go_micro_srv_cust.ModifySerPosStatusRequest, reply *go_micro_srv_cust.ModifySerPosStatusReply) error {
	//dbHandler := db.GetDB(constants.DB_CRM)
	//defer db.PutDB(constants.DB_CRM, dbHandler)

	if req.UseStatus != "1" && req.UseStatus != "0" {
		ss_log.Error("UseStatus is no in(0,1)")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	servicerNo := ""
	switch req.AccountType {
	case constants.AccountType_SERVICER: //服务商
		//查询账号对应的服务商id
		servicerNoT := dao.ServiceDaoInst.GetServicerNoByAccNo(req.AccountUid)
		if servicerNoT == "" {
			ss_log.Error("服务商账号查询服务商id出错,服务商账号uid---------------->%s", req.AccountUid)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		servicerNo = servicerNoT
	case constants.AccountType_POS: // 收银员
		// 获取店员no,再获取服务商no
		cashierNo := dao.RelaAccIdenDaoInst.GetIdenFromAcc(req.AccountUid, constants.AccountType_POS)
		if cashierNo == "" {
			ss_log.Error("店员账号查询店员id出错,店员账号uid---------------->%s", req.AccountUid)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		err, servicerNoT := dao.ServiceDaoInst.GetServiceByCashierNo(cashierNo)
		if err != nil {
			ss_log.Error("店员id查询服务商id出错,店员账号uid---------------->%s", req.AccountUid)
			reply.ResultCode = ss_err.ERR_DB_OP_SER
			return nil
		}
		servicerNo = servicerNoT
	default:
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//要修改的pos机的服务商id
	posServicerNo := dao.ServiceDaoInst.GetSerPosServicerNoByTerminalNo(req.TerminalNo)

	if servicerNo != posServicerNo {
		ss_log.Error("无权限修改他人的pos")
		reply.ResultCode = ss_err.ERR_MODIFY_TERMINAL_STATUS_FAILD
		return nil
	}

	if errStr := dao.ServiceDaoInst.ModifySerPosStatus(req.TerminalNo, req.UseStatus); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//修改服务商信息（营业场所、营业执照，服务商基本信息）
func (c *CustHandler) ModifyServicerInfo(ctx context.Context, req *go_micro_srv_cust.ModifyServicerInfoRequest, reply *go_micro_srv_cust.ModifyServicerInfoReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	//服务商账号
	account, err := dao.ServiceDaoInst.GetAccountBySerNo(req.ServicerNo)
	if err != nil {
		ss_log.Error("服务商账号查询不到")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//修改的应当是服务商的联系人电话、地址、名称
	sqlStr2 := "update servicer set contact_person = $2, contact_phone = $3, contact_addr = $4 where servicer_no = $1 and is_delete = '0' "
	errUp := ss_sql.ExecTx(tx, sqlStr2, req.ServicerNo, req.Nickname, req.Phone, req.Addr)
	if errUp != nil {
		ss_log.Error("修改服务商的联系人、联系人手机号、联系人地址。errUp=[%v]", errUp)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	description := fmt.Sprintf("修改[%v]服务商联系人为[%v],联系人手机号:[%v]、联系人地址:[%v]", account, req.Nickname, req.Phone, req.Addr)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Account)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	//获取旧的图片url集合
	ids, errIds := dao.DictimagesDaoInst.GetImgIds(req.ServicerNo)
	if errIds != nil && errIds != sql.ErrNoRows {
		ss_log.Error("查询服务商[%v]的旧营业执照、营业场景失败", req.ServicerNo)
	}

	oldImgIds := strings.Split(ids, ",")
	imageBaseUrls := dao.ImageDaoInstance.GetImgUrlsByImgIds(oldImgIds)

	var imgUrls []string
	for _, v := range oldImgIds {
		if imageData, errImg := dao.ImageDaoInstance.GetImageUrlById(v); errImg != nil {
			ss_log.Error("err=[%v]", errImg)
			imgUrls = append(imgUrls, "")
		} else {
			imgUrls = append(imgUrls, imageData.ImageUrl)
		}
	}

	if len(imageBaseUrls) < 4 {
		imageBaseUrls = append(imageBaseUrls, "")
		imageBaseUrls = append(imageBaseUrls, "")
		imageBaseUrls = append(imageBaseUrls, "")
		imageBaseUrls = append(imageBaseUrls, "")
	}
	imgStrs := []string{req.ServicerImg1, req.ServicerImg2, req.ServicerImg3, req.ServicerImg4}

	//最新的图片ids
	var imgIds []string
	var delImgUrl []string

	for k, imgStr := range imgStrs {
		if imgStr == "" {
			imgIds = append(imgIds, "")
			if imageBaseUrls[k] != "" { //原来有图片，则要删掉原来的图片
				delImgUrl = append(delImgUrl, imgUrls[k])
			}
			continue
		}

		if imgStr != imageBaseUrls[k] { //是新图片则上传
			if imageBaseUrls[k] != "" { //原来有图片，则要删掉原来的图片
				delImgUrl = append(delImgUrl, imgUrls[k])
			}

			upReg := &go_micro_srv_cust.UploadImageRequest{
				ImageStr:   imgStr,
				AccountUid: req.AccountNo,
				Type:       constants.UploadImage_UnAuth,
			}
			upReply := &go_micro_srv_cust.UploadImageReply{}
			c.UploadImage(ctx, upReg, upReply)
			if upReply.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("addImg1Err=[%v]", upReply.ResultCode)
				reply.ResultCode = ss_err.ERR_SAVE_IMAGE_FAILD
				return nil
			} else {
				imgIds = append(imgIds, upReply.ImageId)
			}
		} else { //和数据库的相同则不用修改了
			imgIds = append(imgIds, oldImgIds[k])
		}
	}

	newIdsStr := ""
	for k, id := range imgIds {
		if k == 0 {
			newIdsStr = id
		} else {
			newIdsStr = newIdsStr + "," + id
		}
	}

	//修改或插入服务商图片ids
	if errStr := dao.ServiceDaoInst.UpdateServicerImg(tx, newIdsStr, req.ServicerNo); errStr != ss_err.ERR_SUCCESS {
		ss_log.Error("err[%v]", errStr)
		reply.ResultCode = errStr
		return nil
	}

	//删除旧图片
	if len(delImgUrl) != 0 {
		if _, err := common.UploadS3.DeleteMulti(delImgUrl); err != nil {
			notes := fmt.Sprintf("ModifyServicerInfo s3上删除图片失败,图片路劲集合为: %+v,err: %s", delImgUrl, err.Error())
			ss_log.Error(notes)
			dao.DictimagesDaoInst.AddDelFaildLog(notes)
		}

		// 删除image表中的图片
		for _, v := range delImgUrl {
			if err := dao.DictimagesDaoInst.Delete(v); err != nil {
				notes := fmt.Sprintf("ModifyServicerInfo 删除图片记录失败,图片路劲为: %s,err: %s", v, err.Error())
				ss_log.Error(notes)
				dao.DictimagesDaoInst.AddDelFaildLog(notes)
			}
		}
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
