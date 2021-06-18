package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_sql"
)

type BusinessSceneDao struct {
	SceneName         string
	ImageNo           string
	Notes             string
	ExampleImgNos     string
	ExampleImgNames   string
	TradeType         string
	BusinessChannelNo string
	Idx               int
	FloatRate         string
	ApplyType         string
	IsManualSigned    int32
}

var BusinessSceneDaoInst BusinessSceneDao

func (BusinessSceneDao) GetBusinessSceneCnt(whereList []*model.WhereSqlCond) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	var totalT sql.NullString
	sqlStr := "SELECT count(1) " +
		"FROM business_scene bs " + whereModel.WhereStr
	if err2 := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, whereModel.Args...); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return ""
	}

	return totalT.String
}

func (BusinessSceneDao) GetBusinessSceneList(whereList []*model.WhereSqlCond, page, pageSize string) (datasT []*go_micro_srv_cust.BusinessSceneData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY bs.idx asc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(pageSize), strext.ToInt(page))
	sqlStr := "SELECT bs.scene_no, bs.scene_name, bs.create_time, bs.update_time, bs.image_no," +
		" bs.notes, bs.idx, bs.is_delete, bs.float_rate, bs.apply_type, bs.is_manual_signed, " +
		" bc.business_channel_no, bc.channel_name " +
		" FROM business_scene bs " +
		" LEFT JOIN business_channel bc ON bc.business_channel_no = bs.business_channel_no "
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr+whereModel.WhereStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	var datas []*go_micro_srv_cust.BusinessSceneData
	for rows.Next() {
		data := go_micro_srv_cust.BusinessSceneData{}
		var isDelete, updateTime, imgNo, notes, idx, floatRate, applyType, isManualSigned, businessChannelNo, channelName sql.NullString
		err2 = rows.Scan(
			&data.SceneNo,
			&data.SceneName,
			&data.CreateTime,
			&updateTime,
			&imgNo,
			&notes,
			&idx,
			&isDelete,
			&floatRate,
			&applyType,
			&isManualSigned,
			&businessChannelNo,
			&channelName,
		)

		if err2 != nil {
			ss_log.Error("err=[%v]", err2)
			return nil, err2
		}

		data.UpdateTime = updateTime.String
		data.ImgNo = imgNo.String
		data.Notes = notes.String
		data.Idx = idx.String
		data.IsEnabled = isDelete.String
		data.FloatRate = floatRate.String
		data.ApplyType = applyType.String
		data.BusinessChannelNo = businessChannelNo.String
		data.ChannelName = channelName.String
		data.IsManualSigned = isManualSigned.String

		datas = append(datas, &data)
	}

	return datas, nil
}

func (BusinessSceneDao) GetBusinessSceneDetail(sceneNo string) (data *go_micro_srv_cust.BusinessSceneData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereList := []*model.WhereSqlCond{
		{Key: "scene_no", Val: sceneNo, EqType: "="},
		//{Key: "is_delete", Val: "0", EqType: "="},
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	sqlStr := "SELECT bs.scene_no, bs.scene_name, bs.create_time, bs.update_time, bs.image_no," +
		" bs.notes, bs.example_img_nos, bs.example_img_names, bs.trade_type, bs.idx, bs.is_delete, " +
		" bs.float_rate, bs.apply_type, bs.is_manual_signed, bc.business_channel_no, bc.channel_name " +
		" FROM business_scene bs " +
		" LEFT JOIN business_channel bc ON bc.business_channel_no = bs.business_channel_no "
	rows, stmt, err2 := ss_sql.QueryRowN(dbHandler, sqlStr+whereModel.WhereStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	dataT := &go_micro_srv_cust.BusinessSceneData{}
	var (
		exampleImgNos, exampleImgNames, tradeType, idx, updateTime,
		imgNo, notes, isEnabled, floatRate, applyType, isManualSigned,
		businessChannelNo, channelName sql.NullString
	)
	err2 = rows.Scan(
		&dataT.SceneNo,
		&dataT.SceneName,
		&dataT.CreateTime,
		&updateTime,
		&imgNo,

		&notes,
		&exampleImgNos,
		&exampleImgNames,
		&tradeType,
		&idx,

		&isEnabled,
		&floatRate,
		&applyType,
		&isManualSigned,
		&businessChannelNo,
		&channelName,
	)

	if err2 != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err2
	}

	dataT.ExampleImgNos = exampleImgNos.String
	dataT.ExampleImgNames = exampleImgNames.String
	dataT.TradeType = tradeType.String
	dataT.Idx = idx.String
	dataT.UpdateTime = updateTime.String
	dataT.ImgNo = imgNo.String
	dataT.Notes = notes.String
	dataT.IsEnabled = isEnabled.String
	dataT.FloatRate = floatRate.String
	dataT.ApplyType = applyType.String
	dataT.IsManualSigned = isManualSigned.String
	dataT.BusinessChannelNo = businessChannelNo.String
	dataT.ChannelName = channelName.String

	return dataT, nil
}

func (b BusinessSceneDao) AddBusinessScene(d *BusinessSceneDao) (sceneNoT string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sceneNo := strext.NewUUID()
	sqlStr := "insert into business_scene(scene_no, scene_name, image_no, notes, example_img_nos, example_img_names, " +
		"trade_type, idx, business_channel_no, float_rate, apply_type, is_manual_signed, create_time, update_time)" +
		" values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, current_timestamp,current_timestamp) "
	if err2 := ss_sql.Exec(dbHandler, sqlStr, sceneNo, d.SceneName, d.ImageNo, d.Notes, d.ExampleImgNos, d.ExampleImgNames,
		d.TradeType, d.Idx, d.BusinessChannelNo, d.FloatRate, d.ApplyType, d.IsManualSigned); err2 != nil {
		return "", err2
	}

	return sceneNo, nil
}

func (b BusinessSceneDao) UpdateBusinessScene(sceneNo, sceneName, imageNo, notes, exampleImgNos, exampleImgNames, tradeType, floatRate, applyType, businessChannelNo string, isManualSigned int32) (err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update business_scene set scene_name =$2, image_no = $3, notes = $4, example_img_nos = $5," +
		" example_img_names = $6, trade_type = $7, float_rate = $8, apply_type = $9, business_channel_no = $10, is_manual_signed = $11, update_time = current_timestamp " +
		" where scene_no = $1 "
	if err2 := ss_sql.Exec(dbHandler, sqlStr, sceneNo, sceneName, imageNo, notes, exampleImgNos, exampleImgNames, tradeType, floatRate, applyType, businessChannelNo, isManualSigned); err2 != nil {
		return err2
	}

	return nil
}

func (b BusinessSceneDao) EnabledScene(isEnabled, sceneNo, isDelete string) (err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlStr := "update business_scene set is_delete = $1 where scene_no = $2 and is_delete = $3 "
	if err2 := ss_sql.Exec(dbHandler, sqlStr, isEnabled, sceneNo, isDelete); err2 != nil {
		return err2
	}

	return nil
}

func (BusinessSceneDao) BusinessGetSceneCnt(whereList []*model.WhereSqlCond, appId string) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	var totalT sql.NullString
	sqlStr := "SELECT count(1) " +
		" FROM business_scene bs" +
		" LEFT JOIN business_signed bsi ON bsi.scene_no = bs.scene_no and bsi.status = '1' AND bsi.app_id = '" + appId + "' " + whereModel.WhereStr
	if err2 := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, whereModel.Args...); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return ""
	}

	return totalT.String
}

//商家获取产品列表(拥有产品是否显示签约按钮的返回字段)
func (BusinessSceneDao) BusinessGetSceneList(whereList []*model.WhereSqlCond, appId string) (datasT []*go_micro_srv_cust.BusinessGetSceneData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY bs.create_time desc ")

	sqlStr := "SELECT bs.scene_no, bs.scene_name " +
		" FROM business_scene bs " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	var datas []*go_micro_srv_cust.BusinessGetSceneData
	for rows.Next() {
		data := go_micro_srv_cust.BusinessGetSceneData{}
		//var signedNo, status, notes sql.NullString
		err2 = rows.Scan(
			&data.SceneNo,
			&data.SceneName,
			//&signedNo,
			//&status,
			//&notes,
		)

		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		status, notes := getSignedStatusAndNotes(dbHandler, appId, data.SceneNo)
		data.Notes = notes //此备注是审核备注

		data.Status = status
		switch status {
		case "":
			fallthrough
		case constants.SignedStatusDeny:
			fallthrough
		case constants.SignedStatusInvalid:
			data.ShowSignedBtn = true
		default:
			data.ShowSignedBtn = false
		}

		datas = append(datas, &data)
	}

	return datas, nil
}

func getSignedStatusAndNotes(dbHandler *sql.DB, appId, sceneNo string) (statusT, notesT string) {
	sqlStr := ` SELECT DISTINCT status, notes, create_time FROM business_signed WHERE app_id = $1 AND scene_no = $2 ORDER BY create_time DESC LIMIT 1`
	rows, stmt, err2 := ss_sql.QueryRowN(dbHandler, sqlStr, appId, sceneNo)
	if stmt != nil {
		defer stmt.Close()
	}
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return "", ""
	}

	var status, notes, createTime sql.NullString
	err2 = rows.Scan(
		&status,
		&notes,
		&createTime,
	)

	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return "", ""
	}
	return status.String, notes.String
}

func (BusinessSceneDao) SetSceneIdx(tx *sql.Tx, sceneNo string, idx int) error {
	sqlStr := "update business_scene set idx = $2 where scene_no = $1"
	return ss_sql.ExecTx(tx, sqlStr, sceneNo, idx)
}

//往上
func (BusinessSceneDao) SwapSceneUp(tx *sql.Tx, idx int) error {
	//sqlStr := "update business_scene set idx = idx-1 where idx = $1 and is_delete = '0' "
	sqlStr := "update business_scene set idx = idx-1 where idx = $1 "
	return ss_sql.ExecTx(tx, sqlStr, idx)
}

//往下
func (BusinessSceneDao) SwapSceneDown(tx *sql.Tx, idx int) error {
	//sqlStr := "update business_scene set idx = idx+1 where idx = $1 and is_delete = '0' "
	sqlStr := "update business_scene set idx = idx+1 where idx = $1 "
	return ss_sql.ExecTx(tx, sqlStr, idx)
}
func (BusinessSceneDao) GetNowSceneMaxIdx() (maxIdx string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//sqlStr := "select MAX(idx) from business_scene where is_delete = '0' "
	sqlStr := "select MAX(idx) from business_scene "
	var maxIdxT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&maxIdxT}); err != nil {
		ss_log.Error("err=[%v]", err)
		return "1"
	}
	return maxIdxT.String

}

//获取商家指定自动签约产品
func (BusinessSceneDao) GetBusinessAutoSceneDetail() (data *go_micro_srv_cust.BusinessSceneData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereList := []*model.WhereSqlCond{
		{Key: "scene_name", Val: constants.AutoSignedSceneName, EqType: "="},
		{Key: "is_manual_signed", Val: constants.ProductIsManualSigned_False, EqType: "="}, //不可手动签约的
		{Key: "is_delete", Val: "0", EqType: "="},                                          //
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	sqlStr := "SELECT bs.scene_no, bs.float_rate, bs.business_channel_no " +
		" FROM business_scene bs "
	rows, stmt, err2 := ss_sql.QueryRowN(dbHandler, sqlStr+whereModel.WhereStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	dataT := &go_micro_srv_cust.BusinessSceneData{}
	var (
		sceneNo, floatRate, businessChannelNo sql.NullString
	)
	err2 = rows.Scan(
		&sceneNo,
		&floatRate,
		&businessChannelNo,
	)

	if err2 != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err2
	}

	dataT.SceneNo = sceneNo.String
	dataT.FloatRate = floatRate.String
	dataT.BusinessChannelNo = businessChannelNo.String

	return dataT, nil
}
