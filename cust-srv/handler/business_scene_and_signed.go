package handler

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/dao"
	"context"
	"database/sql"
)

func (*CustHandler) GetSceneSignedList(ctx context.Context, req *custProto.GetSceneSignedListRequest, reply *custProto.GetSceneSignedListReply) error {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "bs.scene_name", Val: req.SceneName, EqType: "like"},
		{Key: "ss.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "ss.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "ss.status", Val: req.Status, EqType: "="},
		{Key: "bs.apply_type", Val: req.ApplyType, EqType: "="},
		{Key: "bs.scene_no", Val: req.SceneNo, EqType: "="},
		{Key: "ss.business_account_no", Val: req.BusinessAccNo, EqType: "="},
	})

	//查询是否要显示添加申请按钮（没有申请中的记录并且产品是可以手动签约的则为true）
	if req.SceneNo != "" {
		whereModel2 := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
			{Key: "bs.scene_no", Val: req.SceneNo, EqType: "="},
			{Key: "bs.is_manual_signed", Val: constants.ProductIsManualSigned_True, EqType: "="},
			{Key: "ss.status", Val: constants.SignedStatusPending, EqType: "="},
			{Key: "ss.business_account_no", Val: req.BusinessAccNo, EqType: "="},
		})
		total, err := dao.SceneSignedDaoInst.Count(whereModel2.WhereStr, whereModel2.Args)
		if err != nil && err != sql.ErrNoRows {
			ss_log.Error("查询是否有申请中的数据失败，SceneNo[%v], request=%v, err=%v", req.SceneNo, strext.ToJson(req), err)
			reply.ResultCode = ss_err.ERR_SYSTEM
			return nil
		}

		if total == "0" {
			reply.ShowAddButton = true
		}

	}

	total, err := dao.SceneSignedDaoInst.Count(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("统计核销码数量失败，request=%v, err=%v", strext.ToJson(req), err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "ORDER BY ss.create_time DESC ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
	list, err := dao.SceneSignedDaoInst.GetList(whereModel.WhereStr, req.Lang, whereModel.Args)
	if err != nil {
		ss_log.Error("查询核销码列表失败，request=%v, err=%v", strext.ToJson(req), err)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	var keys []string
	var langDatas []*custProto.LangData
	keyMap := make(map[string]string) //用于去重，不用重复查询一些key
	for k, data := range list {
		if data.SceneName != "" { //产品名称记录的是多语言的key
			if _, ok := keyMap[data.SceneName]; !ok { //只有没添加过的才去查询
				keyMap[data.SceneName] = data.SceneName
				keys = append(keys, data.SceneName)
			}
		}

		//一次最多查30个key对应的语言
		if len(keys) == 30 || k == len(list)-1 {
			//读取多语言
			langDatas2, errLang := dao.LangDaoInst.GetLangTextsByKeys(keys)
			if errLang != nil {
				ss_log.Error("查询多语言出错,keys[%v]", keys)
				reply.ResultCode = ss_err.ERR_SYS_DB_GET
				return nil
			}
			langDatas = append(langDatas, langDatas2...)
			keys = keys[0:0]
		}

	}

	var dataList []*custProto.SceneSigned
	for _, v := range list {
		switch req.Lang {
		case constants.LangZhCN:
			for _, langData := range langDatas {
				if v.SceneName == langData.Key {
					v.SceneName = langData.LangCh
					break
				}
			}
		case constants.LangEnUS:
			for _, langData := range langDatas {
				if v.SceneName == langData.Key {
					v.SceneName = langData.LangEn
					break
				}
			}

		case constants.LangKmKH:
			for _, langData := range langDatas {
				if v.SceneName == langData.Key {
					v.SceneName = langData.LangKm
					break
				}
			}
		default:
			for _, langData := range langDatas {
				if v.SceneName == langData.Key {
					v.SceneName = langData.LangEn
					break
				}
			}
		}

		data := new(custProto.SceneSigned)
		data.SignedNo = v.SignedNo
		data.SceneName = v.SceneName
		data.IndustryName = v.IndustryName
		data.StartTime = v.StartTime
		data.EndTime = v.EndTime
		data.Status = v.Status
		data.CreateTime = v.CreateTime
		data.Rate = strext.ToInt32(v.Rate)
		data.Cycle = v.Cycle
		dataList = append(dataList, data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Total = strext.ToInt32(total)
	reply.List = dataList
	return nil
}
