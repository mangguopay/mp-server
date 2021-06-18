package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type BusinessSceneSignedDao struct {
	SignedNo          string
	StartTime         string
	EndTime           string
	BusinessNo        string
	BusinessAccountNo string
	Status            string
	createTime        string
	SceneNo           string
	Rate              string
	Cycle             string
	IndustryNo        string

	//外部字段
	TradeType         string
	BusinessChannelNo string
}

var BusinessSceneSignedDaoInst BusinessSceneSignedDao

type UpstreamRate struct {
	UpWxRate     string
	UpAliPayRate string
}

//查询商家产品签约的上游费率
func (BusinessSceneSignedDao) GetBusinessUpstreamRate(businessNo, sceneNo string) (*UpstreamRate, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var upWxRate, upAliPayRate sql.NullString

	sqlStr := "SELECT bi.up_wx_rate, bi.up_alipay_rate " +
		"FROM business_scene_signed bs " +
		"LEFT JOIN business_scene AS scene ON scene.scene_no = bs.scene_no " +
		"LEFT JOIN business_industry bi ON bi.code = bs.industry_no " +
		"WHERE bs.app_id=$1 AND scene.scene_no=$2 ORDER BY bs.create_time DESC LIMIT 1 "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&upWxRate, &upAliPayRate}, businessNo, sceneNo)
	if err != nil {
		return nil, err
	}

	obj := new(UpstreamRate)
	obj.UpWxRate = upWxRate.String
	obj.UpAliPayRate = upAliPayRate.String
	return obj, nil
}

//根据tradeType查询产品签约信息
func (BusinessSceneSignedDao) GetSignedByTradeType(businessNo, tradeType string) (*BusinessSceneSignedDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var signedNo, businessNoT, status, rate, sceneNo, endTime, tradeTypeT, channelNo sql.NullString
	sqlStr := "SELECT bs.signed_no, bs.business_no, bs.status, bs.rate, bs.scene_no, " +
		"bs.end_time, sc.trade_type, sc.business_channel_no " +
		"FROM business_scene_signed bs " +
		"LEFT JOIN business_scene AS sc ON sc.scene_no = bs.scene_no " +
		"WHERE sc.trade_type=$1 AND bs.status=$2 AND bs.business_no=$3 ORDER BY bs.create_time DESC LIMIT 1 "
	err := ss_sql.QueryRow(dbHandler, sqlStr,
		[]*sql.NullString{&signedNo, &businessNoT, &status, &rate, &sceneNo, &endTime, &tradeTypeT, &channelNo},
		tradeType, constants.SignedStatusPassed, businessNo)
	if err != nil {
		return nil, err
	}

	obj := new(BusinessSceneSignedDao)
	obj.SignedNo = signedNo.String
	obj.BusinessNo = businessNoT.String
	obj.Status = status.String
	obj.Rate = rate.String
	obj.SceneNo = sceneNo.String
	obj.EndTime = endTime.String
	obj.TradeType = tradeTypeT.String
	obj.BusinessChannelNo = channelNo.String

	return obj, nil
}

//获取商户使用中的产品签约的结算周期
func (BusinessSceneSignedDao) GetBusinessSettleCycle(businessNo, sceneNo, signedStatus string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT cycle " +
		"FROM business_scene_signed bs " +
		"LEFT JOIN business_scene sc ON sc.scene_no = bs.scene_no " +
		"WHERE bs.business_no=$1 AND bs.scene_no=$2 AND bs.status=$3 "
	var cycle sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cycle}, businessNo, sceneNo, signedStatus)
	if err != nil {
		return "", err
	}

	return cycle.String, nil
}
