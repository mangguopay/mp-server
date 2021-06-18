package handler

import (
	"a.a/mp-server/common/ss_count"
	"context"
	"fmt"
	"strings"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_rsa"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/common"
	"a.a/mp-server/cust-srv/dao"
	"a.a/mp-server/cust-srv/util"
)

//获取商户主页信息
func (c *CustHandler) GetBusinessAccountHome(ctx context.Context, req *custProto.GetBusinessAccountHomeRequest, reply *custProto.GetBusinessAccountHomeReply) error {

	//查询商家信息
	whereList := []*model.WhereSqlCond{
		{Key: "bu.is_delete", Val: "0", EqType: "="},
		{Key: "bu.account_no", Val: req.AccountUid, EqType: "="},
	}

	bData, err1 := dao.BusinessDaoInst.GetBusinessDetail(whereList)
	if err1 != nil {
		ss_log.Error("err=[%v]", err1)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}
	reply.BusinessData = bData

	//查询商家余额和冻结金额
	_, usdBalance := dao.VaccountDaoInst.GetBalanceFromAccNo(req.AccountUid, constants.VaType_USD_BUSINESS_SETTLED)
	_, usdFrozenBalance := dao.VaccountDaoInst.GetFrozenBalanceFromAccNo(req.AccountUid, constants.VaType_USD_BUSINESS_SETTLED)
	_, khrBalance := dao.VaccountDaoInst.GetBalanceFromAccNo(req.AccountUid, constants.VaType_KHR_BUSINESS_SETTLED)
	_, khrFrozenBalance := dao.VaccountDaoInst.GetFrozenBalanceFromAccNo(req.AccountUid, constants.VaType_KHR_BUSINESS_SETTLED)

	reply.WalletData = &custProto.BusinessWalletData{
		UsdBalance:       usdBalance,
		UsdFrozenBalance: usdFrozenBalance,
		KhrBalance:       khrBalance,
		KhrFrozenBalance: khrFrozenBalance,
	}

	//交易统计
	whereList2 := []*model.WhereSqlCond{
		{Key: "bb.business_no", Val: req.IdenNo, EqType: "="},
	}

	//已支付的
	successSumWhereList := append(whereList2, &model.WhereSqlCond{Key: "bb.order_status", Val: constants.BusinessOrderStatusPay, EqType: "="})

	//获取usd统计信息(成功的)
	usdCnt := dao.BusinessBillDaoInst.GetCnt(append(successSumWhereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_USD, EqType: "="}))
	usdSum := dao.BusinessBillDaoInst.GetSum(append(successSumWhereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_USD, EqType: "="}))

	//获取khr统计信息(成功的)
	khrCnt := dao.BusinessBillDaoInst.GetCnt(append(successSumWhereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_KHR, EqType: "="}))
	khrSum := dao.BusinessBillDaoInst.GetSum(append(successSumWhereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_KHR, EqType: "="}))

	//已退款的
	refundSumWhereList := append(whereList2, &model.WhereSqlCond{Key: "bb.order_status", Val: constants.BusinessOrderStatusRefund, EqType: "="})
	//退款金额
	refundUsdSum := dao.BusinessBillDaoInst.GetSum(append(refundSumWhereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_USD, EqType: "="}))
	refundKhrSum := dao.BusinessBillDaoInst.GetSum(append(refundSumWhereList, &model.WhereSqlCond{Key: "bb.currency_type", Val: constants.CURRENCY_UP_KHR, EqType: "="}))

	reply.SumData = &custProto.BusinessBillsSumData{
		UsdCnt:       usdCnt,
		UsdSum:       usdSum,
		KhrCnt:       khrCnt,
		KhrSum:       khrSum,
		RefundUsdSum: refundUsdSum,
		RefundKhrSum: refundKhrSum,
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//获取商户基本信息
func (c *CustHandler) GetBusinessBaseInfo(ctx context.Context, req *custProto.GetBusinessBaseInfoRequest, reply *custProto.GetBusinessBaseInfoReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "bu.is_delete", Val: "0", EqType: "="},
		{Key: "bu.account_no", Val: req.AccountUid, EqType: "="},
	}

	data, err1 := dao.BusinessDaoInst.GetBusinessDetail(whereList)
	if err1 != nil {
		ss_log.Error("err=[%v]", err1)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//修改商户基本信息
func (c *CustHandler) UpdateBusinessBaseInfo(ctx context.Context, req *custProto.UpdateBusinessBaseInfoRequest, reply *custProto.UpdateBusinessBaseInfoReply) error {
	data, err := dao.BusinessIndustryDaoInst.GetBusinessIndustryDetail(req.MainIndustry)
	if err != nil {
		ss_log.Error("主要行业应用[%v]查询失败,err[%v]", req.MainIndustry, err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if data == nil {
		ss_log.Error("主要行业应用[%v]不存在", req.MainIndustry)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if data.Level == constants.Businesslevel_One {
		ss_log.Error("主要行业应用[%v]的等级是一级，属性分类，不具备费率和结算周期。不允许选择。", req.MainIndustry)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	err1 := dao.BusinessDaoInst.UpdateBusinessInfo(req.AccountUid, req.MainIndustry, req.MainBusiness, req.ContactPerson, req.ContactPhone, req.CountryCode)
	if err1 != nil {
		ss_log.Error("err=[%v]", err1)
		reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//查询绑定的卡列表信息（个人商家、企业商家）
func (c *CustHandler) GetBusinessCards(ctx context.Context, req *custProto.GetBusinessCardsRequest, reply *custProto.GetBusinessCardsReply) error {
	switch req.AccountType {
	case constants.AccountType_PersonalBusiness:
	case constants.AccountType_EnterpriseBusiness:
	default:
		ss_log.Error("AccountType类型错误[%v]", req.AccountType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	datas, total, err := dao.CardBusinessDaoInst.GetBusinessCards(req.AccountNo, req.BalanceType, req.AccountType)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	for _, v := range datas {
		if v.LogoImgNo != "" { //查询图片对应的url
			reqImg := &custProto.UnAuthDownloadImageRequest{
				ImageId: v.LogoImgNo,
			}
			replyImg := &custProto.UnAuthDownloadImageReply{}
			err2 := c.UnAuthDownloadImage(ctx, reqImg, replyImg)
			if err2 != nil || replyImg.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取图片url失败")
			} else {
				v.LogoImgUrl = replyImg.ImageUrl
			}
		}
		if v.LogoImgNoGrey != "" { //查询图片对应的url
			reqImg := &custProto.UnAuthDownloadImageRequest{
				ImageId: v.LogoImgNoGrey,
			}
			replyImg := &custProto.UnAuthDownloadImageReply{}
			err2 := c.UnAuthDownloadImage(ctx, reqImg, replyImg)
			if err2 != nil || replyImg.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取图片url失败")
			} else {
				v.LogoImgUrlGrey = replyImg.ImageUrl
			}
		}
	}

	reply.DataList = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//查询绑定的银行卡详情信息（个人商家、企业商家）
func (c *CustHandler) GetBusinessCardDetail(ctx context.Context, req *custProto.GetBusinessCardDetailRequest, reply *custProto.GetBusinessCardDetailReply) error {
	switch req.AccountType {
	case constants.AccountType_PersonalBusiness:
	case constants.AccountType_EnterpriseBusiness:
	default:
		ss_log.Error("AccountType类型错误[%v]", req.AccountType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	data, err := dao.CardBusinessDaoInst.GetBusinessCardDetail(req.CardNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	if data.LogoImgNo != "" { //查询图片对应的url
		reqImg := &custProto.UnAuthDownloadImageRequest{
			ImageId: data.LogoImgNo,
		}
		replyImg := &custProto.UnAuthDownloadImageReply{}
		err2 := c.UnAuthDownloadImage(ctx, reqImg, replyImg)
		if err2 != nil || replyImg.ResultCode != ss_err.ERR_SUCCESS {
			ss_log.Error("获取图片url失败")
		} else {
			data.LogoImgUrl = replyImg.ImageUrl
		}
	}

	if data.LogoImgNoGrey != "" { //查询图片对应的url
		reqImg := &custProto.UnAuthDownloadImageRequest{
			ImageId: data.LogoImgNoGrey,
		}
		replyImg := &custProto.UnAuthDownloadImageReply{}
		err2 := c.UnAuthDownloadImage(ctx, reqImg, replyImg)
		if err2 != nil || replyImg.ResultCode != ss_err.ERR_SUCCESS {
			ss_log.Error("获取图片url失败")
		} else {
			data.LogoImgUrlGrey = replyImg.ImageUrl
		}
	}

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//删除银行卡
func (c *CustHandler) DelBusinessCard(ctx context.Context, req *custProto.DelBusinessCardRequest, reply *custProto.DelBusinessCardReply) error {
	if err := dao.CardBusinessDaoInst.DeleteCard(req.CardNo); err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (c *CustHandler) GetMainIndustryCascaderDatas(ctx context.Context, req *custProto.GetMainIndustryCascaderDatasRequest, reply *custProto.GetMainIndustryCascaderDatasReply) error {

	whereList := []*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
	}

	datas, err := dao.BusinessIndustryDaoInst.GetBusinessIndustryDatas(whereList, "", "")

	if err != nil {
		ss_log.Error("查询主要业务级联数据失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//处理数据，使数据符合ElementUI组件Cascader级联选择器的数据格式
	datas2 := dao.BusinessIndustryDaoInst.TreatMainIndustryData(datas, req.Lang)

	reply.Datas = datas2
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (c *CustHandler) GetBusinessIndustryList(ctx context.Context, req *custProto.GetBusinessIndustryListRequest, reply *custProto.GetBusinessIndustryListReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "create_time", Val: req.StartTime, EqType: ">="},
		{Key: "create_time", Val: req.EndTime, EqType: "<="},
		{Key: "code", Val: req.Code, EqType: "="},
		{Key: "name_ch", Val: req.NameCh, EqType: "like"},
		{Key: "name_en", Val: req.NameEn, EqType: "like"},
		{Key: "level", Val: req.Level, EqType: "="},
		{Key: "up_code", Val: req.UpCode, EqType: "="},
		{Key: "name_ch", Val: req.Search, EqType: "like"},
		{Key: "is_delete", Val: "0", EqType: "="},
	}

	total := dao.BusinessIndustryDaoInst.GetBusinessIndustryCnt(whereList)

	datas, err := dao.BusinessIndustryDaoInst.GetBusinessIndustryDatas(whereList, req.Page, req.PageSize)
	if err != nil {
		ss_log.Error("查询主要业务数据列表失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (c *CustHandler) GetBusinessIndustryDetail(ctx context.Context, req *custProto.GetBusinessIndustryDetailRequest, reply *custProto.GetBusinessIndustryDetailReply) error {

	data, err := dao.BusinessIndustryDaoInst.GetBusinessIndustryDetail(req.Code)
	if err != nil {
		ss_log.Error("查询主要业务数据详情失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (c *CustHandler) InsertOrUpdateBusinessIndustry(ctx context.Context, req *custProto.InsertOrUpdateBusinessIndustryRequest, reply *custProto.InsertOrUpdateBusinessIndustryReply) error {

	switch req.Level {
	case "1":
	case "2":
	default:
		ss_log.Error("参数level不合法")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.Code == "" {
		data := dao.BusinessIndustryDao{
			NameCh: req.NameCh,
			NameEn: req.NameEn,
			NameKm: req.NameKm,
			Level:  req.Level,
			UpCode: req.UpCode,
		}
		if err := dao.BusinessIndustryDaoInst.AddBusinessIndustry(data); err != nil {
			ss_log.Error("添加主要业务应用失败，err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_ADD
			return nil
		}
	} else {
		data := dao.BusinessIndustryDao{
			Code:   req.Code,
			NameCh: req.NameCh,
			NameEn: req.NameEn,
			NameKm: req.NameKm,
			Level:  req.Level,
			UpCode: req.UpCode,
		}
		if err := dao.BusinessIndustryDaoInst.UpdateBusinessIndustry(data); err != nil {
			ss_log.Error("修改主要业务应用失败，err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (c *CustHandler) DelBusinessIndustry(ctx context.Context, req *custProto.DelBusinessIndustryRequest, reply *custProto.DelBusinessIndustryReply) error {
	if req.Code == "" {
		ss_log.Error("参数Code为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//确认该Code下没有子节点
	if total := dao.BusinessIndustryDaoInst.CountByUpCode(req.Code); total != "0" {
		ss_log.Error("该主要业务应用[%v]下有子节点，不允许删除", req.Code)
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}

	if err := dao.BusinessIndustryDaoInst.DelBusinessIndustry(req.Code); err != nil {
		ss_log.Error("删除主要业务应用[%v]失败，err=[%v]", req.Code, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//行业费率和结算周期列表
func (c *CustHandler) GetBusinessIndustryRateCycleList(ctx context.Context, req *custProto.GetBusinessIndustryRateCycleListRequest, reply *custProto.GetBusinessIndustryRateCycleListReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "birc.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "birc.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "birc.code", Val: req.Code, EqType: "="},
		{Key: "birc.is_delete", Val: "0", EqType: "="},
		{Key: "bi.up_code", Val: req.UpCode, EqType: "="},
		{Key: "bi.name_ch", Val: req.IndustryName, EqType: "like"},
		{Key: "bc.channel_name", Val: req.ChannelName, EqType: "like"},
	}

	total := dao.BusinessIndustryRateCycleDaoInst.GetCnt(whereList)

	datas, err := dao.BusinessIndustryRateCycleDaoInst.GetDatas(whereList, req.Page, req.PageSize)
	if err != nil {
		ss_log.Error("查询主要业务数据列表失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//插入行业费率和结算周期
func (c *CustHandler) InsertOrUpdateBusinessIndustryRateCycle(ctx context.Context, req *custProto.InsertOrUpdateBusinessIndustryRateCycleRequest, reply *custProto.InsertOrUpdateBusinessIndustryRateCycleReply) error {

	data := dao.BusinessIndustryRateCycleDao{
		Id:                req.Id,
		Code:              req.Code,
		BusinessChannelNo: req.BusinessChannelNo,
		Rate:              req.Rate,
		Cycle:             req.Cycle,
	}

	//确认要添加或修改成的行业、支付渠道组合只有一个（删除的不算）
	if !dao.BusinessIndustryRateCycleDaoInst.CheckUnique(data.Id, data.Code, data.BusinessChannelNo) {
		ss_log.Error("经营类目Code[%v]、支付渠道BusinessChannelNo[%v]组合不唯一", data.Code, data.BusinessChannelNo)
		reply.ResultCode = ss_err.ERR_BusinessIndustryRateCycle_Unique_FAILD
		return nil
	}

	//查询支付渠道名称，顺便确认支付渠道是存在的
	channelName, errGetName := dao.BusinessChannelDao.GetChannelName(data.BusinessChannelNo)
	if errGetName != nil {
		ss_log.Error("支付渠道[%v]不存在", data.BusinessChannelNo)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//查询经营类目中文名称，顺便确认经营类目是存在的
	IndustryData, errGetdata := dao.BusinessIndustryDaoInst.GetBusinessIndustryDetail(data.Code)
	if errGetdata != nil {
		ss_log.Error("经营类目[%v]不存在", data.Code)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if IndustryData.Level == constants.Businesslevel_One {
		ss_log.Error("经营类目[%v]的等级是一级，属于分类，不允许添加费率和结算周期。", req.Code)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	description := "" //关键操作日志描述
	if req.Id == "" {
		id, err := dao.BusinessIndustryRateCycleDaoInst.Add(data)
		if err != nil {
			ss_log.Error("添加行业费率结算周期失败，err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_ADD
			return nil
		}

		description = fmt.Sprintf("添加新的行业渠道费率和结算周期,id[%v], 行业[%v],支付渠道[%v],费率[%v],结算周期[%v]",
			id, IndustryData.NameCh, channelName, data.Rate, data.Cycle,
		)

	} else {
		//获取旧数据，如果没有则说明要修改的记录不存在
		oldData, errGetOldData := dao.BusinessIndustryRateCycleDaoInst.GetDetail([]*model.WhereSqlCond{
			{Key: "birc.id", Val: data.Id, EqType: "="},
		})

		if errGetOldData != nil {
			ss_log.Error("查询旧行业费率结算周期失败，id[%v],err=[%v]", data.Id, errGetOldData)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}

		//更新
		if err := dao.BusinessIndustryRateCycleDaoInst.Update(data); err != nil {
			ss_log.Error("修改行业费率结算周期失败，err=[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}

		description = fmt.Sprintf("修改旧的行业费率和结算周期,id[%v], 行业[%v],支付渠道[%v],费率[%v]修改为[%v],结算周期[%v]修改为[%v]",
			data.Id,
			oldData.IndustryName,
			oldData.ChannelName,
			oldData.Rate, data.Rate,
			oldData.Cycle, data.Cycle,
		)
	}

	//关键操作日志
	if errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config); errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (c *CustHandler) DelBusinessIndustryRateCycle(ctx context.Context, req *custProto.DelBusinessIndustryRateCycleRequest, reply *custProto.DelBusinessIndustryRateCycleReply) error {

	if err := dao.BusinessIndustryRateCycleDaoInst.Delete(req.Id); err != nil {
		ss_log.Error("删除行业渠道费率、结算周期失败,Id=[%v], err=[%v]", req.Id, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}

	oldData, errGetOldData := dao.BusinessIndustryRateCycleDaoInst.GetDetail([]*model.WhereSqlCond{
		{Key: "birc.id", Val: req.Id, EqType: "="},
	})

	if errGetOldData != nil {
		ss_log.Error("查询旧行业费率结算周期失败，id[%v],err=[%v]", req.Id, errGetOldData)
		reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
		return nil
	}
	description := fmt.Sprintf("删除行业渠道费率、结算周期，id[%v],行业[%v],支付渠道[%v],费率[%v],结算周期[%v]",
		req.Id,
		oldData.IndustryName,
		oldData.ChannelName,
		oldData.Rate,
		oldData.Cycle,
	)

	//关键操作日志
	if errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config); errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//产品列表
func (c *CustHandler) GetBusinessSceneList(ctx context.Context, req *custProto.GetBusinessSceneListRequest, reply *custProto.GetBusinessSceneListReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "bs.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "bs.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "bs.is_delete", Val: req.IsDelete, EqType: "="},
		{Key: "bs.scene_name", Val: req.SceneName, EqType: "like"},
	}

	total := dao.BusinessSceneDaoInst.GetBusinessSceneCnt(whereList)

	datas, err := dao.BusinessSceneDaoInst.GetBusinessSceneList(whereList, req.Page, req.PageSize)
	if err != nil {
		ss_log.Error("查询产品数据列表失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	var keys []string
	var langDatas []*custProto.LangData
	keyMap := make(map[string]string) //用于去重，不用重复查询一些key
	for k, data := range datas {
		if data.SceneName != "" { //产品名称记录的是多语言的key
			if _, ok := keyMap[data.SceneName]; !ok { //只有没添加过的才去查询
				keyMap[data.SceneName] = data.SceneName
				keys = append(keys, data.SceneName)
			}
		}
		if data.Notes != "" { //产品名称记录的是多语言的key
			if _, ok := keyMap[data.Notes]; !ok { //只有没添加过的才去查询
				keyMap[data.Notes] = data.Notes
				keys = append(keys, data.Notes)
			}
		}

		//一次最多查30个key对应的语言
		if len(keys) == 30 || k == len(datas)-1 {
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

	for _, data := range datas {
		if data.ImgNo != "" {
			reqImg := &custProto.UnAuthDownloadImageRequest{
				ImageId: data.ImgNo,
			}
			replyImg := &custProto.UnAuthDownloadImageReply{}
			if errImg := c.UnAuthDownloadImage(ctx, reqImg, replyImg); errImg != nil {
				ss_log.Error("获取图片url失败,图片id[%v],err=[%v]", data.ImgNo, errImg)
			} else {
				data.ImgUrl = replyImg.ImageUrl
			}
		}

		for _, langData := range langDatas {
			switch req.Lang {
			case constants.LangZhCN:
				switch langData.Key {
				case data.SceneName:
					data.SceneName = langData.LangCh
				case data.Notes:
					data.Notes = langData.LangCh
				}

			case constants.LangEnUS:
				switch langData.Key {
				case data.SceneName:
					data.SceneName = langData.LangEn
				case data.Notes:
					data.Notes = langData.LangEn
				}
			case constants.LangKmKH:
				switch langData.Key {
				case data.SceneName:
					data.SceneName = langData.LangKm
				case data.Notes:
					data.Notes = langData.LangKm
				}
			default:

			}
		}
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//商家产品列表
func (c *CustHandler) BusinessGetSceneList(ctx context.Context, req *custProto.BusinessGetSceneListRequest, reply *custProto.BusinessGetSceneListReply) error {
	if req.AppId == "" {
		ss_log.Error("参数AppId为空")
	}

	var whereList []*model.WhereSqlCond
	whereList = append(whereList, &model.WhereSqlCond{Key: "bs.apply_type", Val: req.ApplyType, EqType: "="})
	//只查询启用状态的产品
	whereList = append(whereList, &model.WhereSqlCond{Key: "bs.is_delete", Val: "1", EqType: "!="})

	total := dao.BusinessSceneDaoInst.BusinessGetSceneCnt(whereList, req.AppId)

	datas, err := dao.BusinessSceneDaoInst.BusinessGetSceneList(whereList, req.AppId)
	if err != nil {
		ss_log.Error("查询主要业务数据列表失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	var keys []string
	for _, data := range datas {
		keys = append(keys, data.SceneName)
		if data.Notes != "" {
			keys = append(keys, data.Notes)
		}
	}

	langDatas, errLang := dao.LangDaoInst.GetLangTextsByKeys(keys)
	if errLang != nil {
		ss_log.Error("查询多语言出错,keys[%v]", keys)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	for _, data := range datas {
		for _, langData := range langDatas {
			switch req.Lang {
			case constants.LangZhCN:
				if data.SceneName == langData.Key {
					data.SceneName = langData.LangCh
				}

			case constants.LangEnUS:
				if data.SceneName == langData.Key {
					data.SceneName = langData.LangEn
				}
			case constants.LangKmKH:
				if data.SceneName == langData.Key {
					data.SceneName = langData.LangKm
				}
			default:

			}
		}
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (c *CustHandler) GetBusinessSceneDetail(ctx context.Context, req *custProto.GetBusinessSceneDetailRequest, reply *custProto.GetBusinessSceneDetailReply) error {
	if req.SceneNo == "" {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	data, err := dao.BusinessSceneDaoInst.GetBusinessSceneDetail(req.SceneNo)
	if err != nil {
		ss_log.Error("查询产品详细数据失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}
	splitStr := "&&&&" //数据库保存使用案例图片、使用案例描述的分割符（数据库是一个字段保存多个使用案例图片id）

	var keys []string
	keys = append(keys, data.SceneName)
	if data.Notes != "" {
		keys = append(keys, data.Notes)
	}
	for _, v := range strings.Split(data.ExampleImgNames, splitStr) {
		if v != "" {
			keys = append(keys, v)
		}
	}

	//读取多语言
	langDatas, errLang := dao.LangDaoInst.GetLangTextsByKeys(keys)
	if errLang != nil {
		ss_log.Error("查询多语言出错,keys[%v]", keys)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	for _, langData := range langDatas {
		switch req.Lang {
		case constants.LangZhCN:
			switch langData.Key {
			case data.SceneName:
				data.SceneName = langData.LangCh
			case data.Notes:
				data.Notes = langData.LangCh
			}

			var newStr []string
			for _, v := range strings.Split(data.ExampleImgNames, splitStr) {
				if v == langData.Key {
					newStr = append(newStr, langData.LangCh)
				} else {
					newStr = append(newStr, v)
				}
			}
			data.ExampleImgNames = strings.Join(newStr, splitStr)

		case constants.LangEnUS:
			switch langData.Key {
			case data.SceneName:
				data.SceneName = langData.LangEn
			case data.Notes:
				data.Notes = langData.LangEn
			}

			var newStr []string
			for _, v := range strings.Split(data.ExampleImgNames, splitStr) {
				if v == langData.Key {
					newStr = append(newStr, langData.LangEn)
				} else {
					newStr = append(newStr, v)
				}
			}
			data.ExampleImgNames = strings.Join(newStr, splitStr)
		case constants.LangKmKH:
			switch langData.Key {
			case data.SceneName:
				data.SceneName = langData.LangKm
			case data.Notes:
				data.Notes = langData.LangKm
			}

			var newStr []string
			for _, v := range strings.Split(data.ExampleImgNames, splitStr) {
				if v == langData.Key {
					newStr = append(newStr, langData.LangKm)
				} else {
					newStr = append(newStr, v)
				}
			}
			data.ExampleImgNames = strings.Join(newStr, splitStr)
		default:

		}
	}

	var imgNos []string
	imgNos = append(imgNos, data.ImgNo)
	imgNos = append(imgNos, strings.Split(data.ExampleImgNos, splitStr)...)

	urlStr := dao.ImageDaoInstance.GetImgUrlsByImgIds(imgNos)

	exampleImgUrlStr := ""
	for k, v := range urlStr[1:] { //第一是产品图标url，所以不放到使用案例图片里
		if k == 0 {
			exampleImgUrlStr = v
		} else {
			exampleImgUrlStr = exampleImgUrlStr + splitStr + v
		}
	}
	data.ExampleImgUrls = exampleImgUrlStr

	data.ImgUrl = urlStr[0]

	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//启用/禁用产品
func (c *CustHandler) IsEnabledScene(ctx context.Context, req *custProto.IsEnabledSceneRequest, reply *custProto.IsEnabledSceneReply) error {
	if req.SceneNo == "" {
		ss_log.Error("参数SceneNo为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.IsEnabled == "" {
		ss_log.Error("参数IsEnabled为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//req.IsEnabled：0-启用，1-禁用
	//isDelete 参数为之前的状态
	isDelete := "0"
	if req.IsEnabled == "0" {
		isDelete = "1"
	} else if req.IsEnabled == "1" {
		isDelete = "0"
	} else {
		ss_log.Error("参数IsEnabled值错误，req.IsEnabled=%v", req.IsEnabled)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if err := dao.BusinessSceneDaoInst.EnabledScene(req.IsEnabled, req.SceneNo, isDelete); err != nil {
		ss_log.Error("禁用产品[%v]失败，err=[%v]", req.SceneNo, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_DELETE
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (c *CustHandler) InsertOrUpdateBusinessScene(ctx context.Context, req *custProto.InsertOrUpdateBusinessSceneRequest, reply *custProto.InsertOrUpdateBusinessSceneReply) error {
	splitStr := "&&&&" //数据库保存使用案例图片、使用案例描述的分割符（数据库是一个字段保存多个使用案例图片id）

	if req.SceneNo == "" {
		if req.PaymentChannel == "" {
			ss_log.Error("PaymentChannel参数为空")
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		//上传全部图片
		var imgStrs []string
		imgStrs = append(imgStrs, req.ImageStr)
		imgStrs = append(imgStrs, strings.Split(req.ExampleImgs, splitStr)...)

		var imgIds []string
		for k, v := range imgStrs {
			ss_log.Info("k=[%v]", k)
			if v == "" {
				imgIds = append(imgIds, "")
				ss_log.Error("上传图片有空的,len=[%v],k=[%v]", len(imgStrs), k)
				continue
			}

			imgReq := &custProto.UploadImageRequest{
				ImageStr:   v,
				AccountUid: req.LoginUid,
				Type:       constants.UploadImage_UnAuth,
			}
			imgReply := &custProto.UploadImageReply{}
			errU := c.UploadImage(ctx, imgReq, imgReply)
			if errU != nil || imgReply.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("k[%v],addErr=[%v]", k, imgReply.ResultCode)
				reply.ResultCode = ss_err.ERR_SAVE_IMAGE_FAILD
				return nil
			}
			imgIds = append(imgIds, imgReply.ImageId)
		}

		exampleImgNos := ""
		for k, v := range imgIds[1:] { //第一个id是产品图标id,这里拼的是使用案例图片ids
			if k == 0 {
				exampleImgNos = v
			} else {
				exampleImgNos = exampleImgNos + splitStr + v
			}
		}

		//获取当前最大的序号
		maxIdx := dao.BusinessSceneDaoInst.GetNowSceneMaxIdx()
		addIdx := strext.ToInt(maxIdx) + 1

		d := &dao.BusinessSceneDao{
			SceneName:         req.SceneName,
			ImageNo:           imgIds[0],
			Notes:             req.Notes,
			ExampleImgNos:     exampleImgNos,
			ExampleImgNames:   req.ExampleNames,
			TradeType:         req.TradeType,
			BusinessChannelNo: req.PaymentChannel,
			Idx:               addIdx,
			FloatRate:         req.FloatRate,
			ApplyType:         req.ApplyType,
			IsManualSigned:    req.IsManualSigned,
		}
		if _, err := dao.BusinessSceneDaoInst.AddBusinessScene(d); err != nil {
			ss_log.Error("插入新产品失败,req[%+v],err[%v]", req, err)
			reply.ResultCode = ss_err.ERR_SYS_DB_ADD
			return nil
		}
	} else {
		//获取修改前的产品数据，产品图标、使用案例图片ids集合，与url集合
		imgNos, imgUrls, errStr := getBusinessSceneOldImgNosAndImgUrlsData(req.SceneNo, splitStr)
		if errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}

		//上传的图片字符串集合，有修改过的则是base64,否则是url
		var imgStrs []string
		imgStrs = append(imgStrs, req.ImageStr)
		imgStrs = append(imgStrs, strings.Split(req.ExampleImgs, splitStr)...)

		//处理数据，获取最新的图片ids集合与要删除的原图片urls集合
		newImgNos, delImgUrls, errStr2 := c.processBusinessSceneImg(imgStrs, imgNos, imgUrls, req.LoginUid)
		if errStr2 != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr2
			return nil
		}

		//删除原来s3上的图片
		if len(delImgUrls) > 0 {
			if _, err := common.UploadS3.DeleteMulti(delImgUrls); err != nil {
				notes := fmt.Sprintf("InsertOrUpdateBusinessScene s3上删除图片失败,图片路劲集合为: %+v,err: %s", delImgUrls, err.Error())
				ss_log.Error(notes)
				dao.DictimagesDaoInst.AddDelFaildLog(notes)
			}
			// 删除image表中的图片
			for _, v := range delImgUrls {
				if err := dao.DictimagesDaoInst.Delete(v); err != nil {
					notes := fmt.Sprintf("InsertOrUpdateBusinessScene 删除图片记录失败,图片路劲为: %s,err: %s", v, err.Error())
					ss_log.Error(notes)
					dao.DictimagesDaoInst.AddDelFaildLog(notes)
				}
			}
		}

		newImgNosStr := ""
		for k, v := range newImgNos[1:] {
			if k == 0 {
				newImgNosStr = v
			} else {
				newImgNosStr = newImgNosStr + splitStr + v
			}
		}

		err := dao.BusinessSceneDaoInst.UpdateBusinessScene(req.SceneNo, req.SceneName, newImgNos[0], req.Notes, newImgNosStr,
			req.ExampleNames, req.TradeType, req.FloatRate, req.ApplyType, req.PaymentChannel, req.IsManualSigned)
		if err != nil {
			ss_log.Error("修改产品失败，req[%+v],err=[%v]", req, err)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//imgStrs上传的图片字符串集合，有修改过的则是base64,否则是url
func (c *CustHandler) processBusinessSceneImg(imgStrs, imgNos, imgUrls []string, loginUid string) (newImgNosT []string, delImgUrlsT []string, errT string) {
	var newImgNos []string  //最新的图片ids集合
	var delImgUrls []string //要删除的原图片urls
	//判断哪个是修改过的
	for k, v := range imgStrs {
		if k < len(imgUrls) { //小于数据库图片张数
			if imgUrls[k] != "" {
				if strings.Contains(v, imgUrls[k]) { //有和数据库的该位置一样的url说明图片未改变
					newImgNos = append(newImgNos, imgNos[k])
					continue
				} else { //当改变的时候,添加删除原有图片url,后面添加新图片
					delImgUrls = append(delImgUrls, imgUrls[k])
				}
			}
		}
		//传的图片比数据库的多了、有图片改变的时候,都要添加新图片
		if v != "" { //不为空
			imgReq := &custProto.UploadImageRequest{
				ImageStr:   v,
				AccountUid: loginUid,
				Type:       constants.UploadImage_UnAuth,
			}
			imgReply := &custProto.UploadImageReply{}
			//imgReply, errU := CustHandlerInst.Client.UploadImage(context.TODO(), imgReq)
			errU := c.UploadImage(context.TODO(), imgReq, imgReply)
			if errU != nil || imgReply.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("k[%v],addErr=[%v]", k, imgReply.ResultCode)
				return nil, nil, ss_err.ERR_SAVE_IMAGE_FAILD
			}
			newImgNos = append(newImgNos, imgReply.ImageId)
		} else {
			newImgNos = append(newImgNos, "")
		}
	}

	return newImgNos, delImgUrls, ss_err.ERR_SUCCESS
}
func getBusinessSceneOldImgNosAndImgUrlsData(sceneNo, splitStr string) (imgNosT, imgUrlsT []string, errT string) {
	//查询旧数据
	data, err := dao.BusinessSceneDaoInst.GetBusinessSceneDetail(sceneNo)
	if err != nil {
		ss_log.Error("查询不到要修改的产品[%v]，err=[%v]", sceneNo, err)
		return nil, nil, ss_err.ERR_SYS_DB_GET
	}

	var imgNos []string
	imgNos = append(imgNos, data.ImgNo)
	imgNos = append(imgNos, strings.Split(data.ExampleImgNos, splitStr)...)

	//获取旧url集合
	var imgUrls []string
	for _, v := range imgNos {
		if imageData, errImg := dao.ImageDaoInstance.GetImageUrlById(v); errImg != nil {
			ss_log.Error("err=[%v]", errImg)
			imgUrls = append(imgUrls, "")
		} else {
			imgUrls = append(imgUrls, imageData.ImageUrl)
		}
	}

	ss_log.Info("imgUrls:[%+v]", imgUrls)
	return imgNos, imgUrls, ss_err.ERR_SUCCESS
}

//这是旧的签约逻辑所查询签约的接口（应用详情选产品、行业签）
func (c *CustHandler) GetBusinessSignedList(ctx context.Context, req *custProto.GetBusinessSignedListRequest, reply *custProto.GetBusinessSignedListReply) error {

	status := ""
	if req.IsBusinessReq { //是否是商家前端的请求（商家前端只显示通过和已过期的）
		status = "('" + constants.SignedStatusPassed + "','" + constants.SignedStatusInvalid + "')"
	}

	whereList := []*model.WhereSqlCond{
		{Key: "bs.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "bs.create_time", Val: req.EndTime, EqType: "<="},
		//{Key: "bs.is_delete", Val: "0", EqType: "="},
		{Key: "bs.business_account_no", Val: req.AccountUid, EqType: "="},
		{Key: "bs.app_id", Val: req.AppId, EqType: "="},
		{Key: "bs.status", Val: req.Status, EqType: "="},
		{Key: "bs.status", Val: status, EqType: "in"},
		{Key: "acc.account", Val: req.Account, EqType: "like"},
	}

	total := dao.BusinessSignedDaoInst.GetBusinessSignedCnt(whereList)

	datas, err := dao.BusinessSignedDaoInst.GetBusinessSignedList(whereList, req.Page, req.PageSize)
	if err != nil {
		ss_log.Error("查询商家签约数据列表失败，req=[%+v],err=[%v]", req, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	var keys []string
	for _, data := range datas {
		keys = append(keys, data.SceneName)
	}

	//读取多语言
	langDatas, errLang := dao.LangDaoInst.GetLangTextsByKeys(keys)
	if errLang != nil {
		ss_log.Error("查询多语言出错,keys[%v]", keys)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	for _, data := range datas {
		for _, langData := range langDatas {
			switch req.Lang {
			case constants.LangZhCN:
				if data.SceneName == langData.Key {
					data.SceneName = langData.LangCh
				}

			case constants.LangEnUS:
				if data.SceneName == langData.Key {
					data.SceneName = langData.LangEn
				}

			case constants.LangKmKH:
				if data.SceneName == langData.Key {
					data.SceneName = langData.LangKm
				}

			default:

			}
		}
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//这是旧的签约逻辑所添加签约的接口（应用详情选产品、行业签）
func (c *CustHandler) AddBusinessSigned(ctx context.Context, req *custProto.AddBusinessSignedRequest, reply *custProto.AddBusinessSignedReply) error {
	if !dao.BusinessAppDaoInst.CheckBusinessApp(req.AppId, req.BusinessNo) {
		ss_log.Error("应用编号[%v]不属于商家[%v],无权限添加签约。", req.AppId, req.BusinessNo)
		reply.ResultCode = ss_err.ERR_SYS_NO_API_AUTH
		return nil
	}

	//确认产品是否存在
	sceneData, errGetScene := dao.BusinessSceneDaoInst.GetBusinessSceneDetail(req.SceneNo)
	if errGetScene != nil {
		ss_log.Error("产品[%v]不存在", req.SceneNo)
		reply.ResultCode = ss_err.ERR_PAY_UNKOWN_PRODUCT
		return nil
	}

	//一个应用可有多个产品,但产品只能签一次
	if dao.BusinessSignedDaoInst.CheckAppAndSceneUnique(req.AppId, req.SceneNo) {
		ss_log.Error("应用编号[%v]已添加过产品[%v],不允许重复添加。", req.AppId, req.SceneNo)
		reply.ResultCode = ss_err.ERR_Business_App_Scene_Unique_FAILD
		return nil
	}

	//检验经营类目
	indData, err := dao.BusinessIndustryDaoInst.GetBusinessIndustryDetail(req.IndustryNo)
	if err != nil {
		ss_log.Error("查询经营类目的费率和结算周期失败.")
		reply.ResultCode = ss_err.ERR_NoBusinessIndustryRateCycle_FAILD
		return nil
	}

	if indData == nil {
		ss_log.Error("经营类目[%v]不存在", req.IndustryNo)
		reply.ResultCode = ss_err.ERR_NoBusinessIndustryRateCycle_FAILD
		return nil
	}

	if indData.Level == constants.Businesslevel_One {
		ss_log.Error("经营类目[%v]的等级是一级，属于分类，不具备费率和结算周期。不允许选择。", req.IndustryNo)
		reply.ResultCode = ss_err.ERR_NoBusinessIndustryRateCycle_FAILD
		return nil
	}

	//根据code和business_channel_no查询经营类目的费率和结算周期
	rateData, errGetRate := dao.BusinessIndustryRateCycleDaoInst.GetDetail([]*model.WhereSqlCond{
		{Key: "birc.business_channel_no", Val: sceneData.BusinessChannelNo, EqType: "="},
		{Key: "birc.code", Val: req.IndustryNo, EqType: "="},
		{Key: "birc.is_delete", Val: "0", EqType: "="},
	})

	if errGetRate != nil {
		ss_log.Error("获取行业费率、结算周期失败, IndustryNo[%v], BusinessChannelNo[%v], err[%v]", req.IndustryNo, sceneData.BusinessChannelNo, errGetRate)
		reply.ResultCode = ss_err.ERR_NoBusinessIndustryRateCycle_FAILD
		return nil
	}

	if rateData.Id == "" {
		reply.ResultCode = ss_err.ERR_NoBusinessIndustryRateCycle_FAILD
		return nil
	}

	//签约费率 = 基础费率（行业费率） + 产品浮动费率
	rate := ss_count.Add(rateData.Rate, sceneData.FloatRate)
	cycle := rateData.Cycle

	if strext.ToInt(rate) < 0 { //小于0一律按0处理
		rate = "0"
	}

	_, errAdd1 := dao.BusinessSignedDaoInst.AddBusinessSigned(req.AppId, req.AccUid, req.BusinessNo, req.SceneNo, rate, cycle, req.IndustryNo)
	if errAdd1 != nil {
		ss_log.Error("添加签约失败,errAdd[%v]", errAdd1)
		reply.ResultCode = ss_err.ERR_SYS_DB_ADD
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (c *CustHandler) GetBusinessList(ctx context.Context, req *custProto.GetBusinessListRequest, reply *custProto.GetBusinessListReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "bu.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "bu.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "bu.use_status", Val: req.UseStatus, EqType: "="},
		{Key: "bu.is_delete", Val: "0", EqType: "="},
		{Key: "bu.business_id", Val: req.BusinessId, EqType: "like"},        //商家id
		{Key: "acc.business_auth_status", Val: req.AuthStatus, EqType: "="}, //商家认证状态
		{Key: "acc.account", Val: req.Account, EqType: "like"},              //商家账号
	}

	total := dao.BusinessDaoInst.GetCnt(whereList)

	datas, err := dao.BusinessDaoInst.GetBusinessList(whereList, strext.ToInt(req.Page), strext.ToInt(req.PageSize))
	if err != nil {
		ss_log.Error("查询商家数据列表失败，req=[%+v],err=[%v]", req, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (c *CustHandler) ModifyBusinessStatus(ctx context.Context, req *custProto.ModifyBusinessStatusRequest, reply *custProto.ModifyBusinessStatusReply) error {
	switch req.Status {
	case constants.Status_Disable:
	case constants.Status_Enable:
	default:
		ss_log.Error("参数UseStatus不合法")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	switch req.StatusType {
	case "use_status": //商家状态
		if err := dao.BusinessDaoInst.UpdateBusinessStatus(req.BusinessNo, req.Status); err != nil {
			ss_log.Error("修改商家状态失败,BusinessNo[%v],err=[%v]", req.BusinessNo, err)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}
	case "income_authorization": //商家收款状态
		if err := dao.BusinessDaoInst.UpdateBusinessIncomeAuthorizationStatus(req.BusinessNo, req.Status); err != nil {
			ss_log.Error("修改商家收款状态失败,BusinessNo[%v],err=[%v]", req.BusinessNo, err)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}

	case "outgo_authorization": //商家出款状态
		if err := dao.BusinessDaoInst.UpdateBusinessOutgoAuthorizationStatus(req.BusinessNo, req.Status); err != nil {
			ss_log.Error("修改商家出款状态失败,BusinessNo[%v],err=[%v]", req.BusinessNo, err)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}
	default:
		ss_log.Error("状态类型[%v]错误 no in ( use_status, income_authorization, outgo_authorization)", req.StatusType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//这是旧的签约逻辑所修改的费率和结算周期（应用详情选产品、行业签）
func (c *CustHandler) UpdateBusinessSignedInfo(ctx context.Context, req *custProto.UpdateBusinessSignedInfoRequest, reply *custProto.UpdateBusinessSignedInfoReply) error {
	if err := dao.BusinessSignedDaoInst.UpdateInfo(req.SignedId, req.Cycle, req.Rate); err != nil {
		ss_log.Error("err[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
		return nil
	}

	// 添加关键操作日志
	description := fmt.Sprintf("将商家应用签约产品的签约[%v]的结算周期更改为[%v],费率修改为[%v]", req.SignedId, req.Cycle, req.Rate)
	if errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Business); errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//这是旧的签约逻辑，审核通过或不通过（应用详情选产品、行业签）
func (c *CustHandler) UpdateBusinessSignedStatus(ctx context.Context, req *custProto.UpdateBusinessSignedStatusRequest, reply *custProto.UpdateBusinessSignedStatusReply) error {

	switch req.Status {
	case constants.SignedStatusPassed:
	case constants.SignedStatusDeny:
	default:
		ss_log.Error("Status参数[%v]不合法", req.Status)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	if err := dao.BusinessSignedDaoInst.UpdateStatus(req.SignedId, req.Status, req.Notes); err != nil {
		ss_log.Error("err[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
		return nil
	}

	str, _ := util.GetParamZhCn(req.Status, util.AppSignedStatus)

	// 添加关键操作日志
	description := fmt.Sprintf("将商家应用签约产品的签约[%v]状态更改为[%v]", req.SignedId, str)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Business)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//商家消息列表
func (c *CustHandler) GetBusinessMessagesList(ctx context.Context, req *custProto.GetBusinessMessagesListRequest, reply *custProto.GetBusinessMessagesListReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "bm.account_type", Val: req.AccountType, EqType: "="},
		{Key: "bm.account_no", Val: req.AccountNo, EqType: "="},
		{Key: "acc.account", Val: req.Account, EqType: "like"},
	}

	total, errGet := dao.LogBusinessMessagesDaoInst.GetCnt(whereList)
	if errGet != nil {
		ss_log.Error("获取数量出错,err=[%v]", errGet)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	datas, err := dao.LogBusinessMessagesDaoInst.GetBusinessMessagesList(whereList, strext.ToInt(req.Page), strext.ToInt(req.PageSize))
	if err != nil {
		ss_log.Error("获取商家消息出错 err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	return nil
}

//商家未读消息数量
func (c *CustHandler) GetBusinessMessagesUnRead(ctx context.Context, req *custProto.GetBusinessMessagesUnReadRequest, reply *custProto.GetBusinessMessagesUnReadReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "bm.account_type", Val: req.AccountType, EqType: "="},
		{Key: "bm.account_no", Val: req.AccountNo, EqType: "="},
		{Key: "bm.is_read", Val: "0", EqType: "="}, //未读消息
	}

	total, errGet := dao.LogBusinessMessagesDaoInst.GetCnt(whereList)
	if errGet != nil {
		ss_log.Error("获取未读消息数量出错,err=[%v]", errGet)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Number = total
	return nil
}

//修改商家未读消息为已读
func (c *CustHandler) ReadAllBusinessMessages(ctx context.Context, req *custProto.ReadAllBusinessMessagesRequest, reply *custProto.ReadAllBusinessMessagesReply) error {

	if err := dao.LogBusinessMessagesDaoInst.ModiftAllRead(req.AccountNo); err != nil {
		ss_log.Error("修改商家未读消息为已读失败。AccountNo[%v],err=[%v]", req.AccountNo, err)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//修改商家应用部分信息
func (c *CustHandler) BusinessUpdatePartial(ctx context.Context, req *custProto.BusinessUpdatePartialRequest, reply *custProto.BusinessUpdatePartialReply) error {
	if !dao.BusinessAppDaoInst.CheckBusinessApp(req.AppId, req.IdenNo) {
		ss_log.Error("应用编号[%v]不属于商家[%v],无权限修改他人应用的信息。", req.AppId, req.IdenNo)
		reply.ResultCode = ss_err.ERR_SYS_NO_API_AUTH
		return nil
	}

	// 商户公钥需去头去尾
	req.PubKey = ss_rsa.StripRSAKey(req.PubKey)

	if err := dao.BusinessAppDaoInst.ModifyAppBusinessPartial(req.AppId, req.PubKey, req.SignMethod, req.IpWhiteList); err != nil {
		ss_log.Error("修改商家应用部分信息失败。AppId[%v],err=[%v]", req.AppId, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//修改产品序号
func (c *CustHandler) UpdateBusinessSceneIdx(ctx context.Context, req *custProto.UpdateBusinessSceneIdxRequest, reply *custProto.UpdateBusinessSceneIdxReply) error {
	swapIdx := strext.ToInt(req.SwapIdx)
	if swapIdx <= 0 {
		ss_log.Error("不允许小于或等于0，SwapIdx[%v]", req.SwapIdx)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	maxIdx := dao.BusinessSceneDaoInst.GetNowSceneMaxIdx()
	if swapIdx > strext.ToInt(maxIdx) {
		ss_log.Error("不允许调到比最大位置的数还大，maxIdx[%v],SwapIdx[%v]", maxIdx, req.SwapIdx)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	data, err := dao.BusinessSceneDaoInst.GetBusinessSceneDetail(req.SceneNo)
	if err != nil {
		ss_log.Error("查询产品[%v]信息出错, err=[%v]", req.SceneNo, err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}
	defer ss_sql.Rollback(tx)

	nowIdx := strext.ToInt(data.Idx)
	//往前移
	if nowIdx > swapIdx { //从该位置到产品位置的原产品位置都+1，再把产品移动swapIdx位置
		for i := nowIdx - 1; i >= swapIdx; i-- {
			if err := dao.BusinessSceneDaoInst.SwapSceneDown(tx, i); err != nil {
				ss_log.Error("位置往后移出错，i[%v]", i)
				reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
				return nil
			}
		}

		if err := dao.BusinessSceneDaoInst.SetSceneIdx(tx, req.SceneNo, swapIdx); err != nil {
			ss_log.Error("目标产品[%v]移动到位置[%v]出错。", req.SceneNo, swapIdx)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}

	} else if nowIdx < swapIdx { //从该位置到产品位置的原产品位置都-1，把产品移动到swapIdx位置
		for i := nowIdx + 1; i <= swapIdx; i++ {
			if err := dao.BusinessSceneDaoInst.SwapSceneUp(tx, i); err != nil {
				ss_log.Error("位置往前移出错，i[%v]", i)
				reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
				return nil
			}
		}

		if err := dao.BusinessSceneDaoInst.SetSceneIdx(tx, req.SceneNo, swapIdx); err != nil {
			ss_log.Error("目标产品[%v]移动到位置[%v]出错。", req.SceneNo, swapIdx)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}
	}
	//如果相等，不做处理

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//商家账户列表
func (*CustHandler) GetBusinessAccounts(ctx context.Context, req *custProto.GetBusinessAccountsRequest, reply *custProto.GetBusinessAccountsReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "acc.account", Val: req.Account, EqType: "like"},
	}

	total, errCnt := dao.BusinessDaoInst.GetBusinessAccountCnt(whereList)
	if errCnt != nil {
		ss_log.Error("查询服务商账户数量出错,err[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	datas, err := dao.BusinessDaoInst.GetBusinessAccount(whereList, req.Page, req.PageSize, req.SortType)
	if err != nil {
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	return nil
}

//商家账户收益
func (*CustHandler) GetBusinessAccountsProfit(ctx context.Context, req *custProto.GetBusinessAccountsProfitRequest, reply *custProto.GetBusinessAccountsProfitReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "acc.account", Val: req.Account, EqType: "like"},
	}

	//统计商户数量
	total, errCnt := dao.BusinessDaoInst.GetBusinessAccountCnt(whereList)
	if errCnt != nil {
		ss_log.Error("查询服务商账户数量出错,err[%v]", errCnt)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if total == "0" {
		ss_log.Info("统计出商家数量为0")
		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	}

	//商户收益：订单结算，他人转入
	profitReason := []string{constants.VaReason_Business_Settle, constants.VaReason_BusinessTransferToBusiness}
	profit, err := dao.BusinessDaoInst.GetBusinessProfit(whereList, req.Page, req.PageSize, constants.VaOpType_Add, profitReason)
	if err != nil {
		ss_log.Error("查询商户收益失败, err=%v", err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	//商户支出：商家退款
	expReason := []string{constants.VaReason_BusinessRefund}
	expenditure, err := dao.BusinessDaoInst.GetBusinessProfit(whereList, req.Page, req.PageSize, constants.VaOpType_Minus, expReason)
	if err != nil {
		ss_log.Error("查询商户收益支出失败, err=%v", err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	var list []*custProto.BusinessAccountsProfitData
	for _, p := range profit {
		var obj custProto.BusinessAccountsProfitData
		obj.BusinessAccount = p.BusinessAccount
		obj.BusinessNo = p.BusinessNo
		obj.BusinessAccountNo = p.BusinessAccNo
		obj.BusinessType = p.BusinessType
		obj.UsdProfit = p.UsdAmount
		obj.KhrProfit = p.KhrAmount
		list = append(list, &obj)
	}

	for _, v := range list {
		for _, e := range expenditure {
			if v.BusinessAccountNo == e.BusinessAccNo {
				v.UsdExpenditure = e.UsdAmount
				v.KhrExpenditure = e.KhrAmount
			}
		}
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Total = strext.ToInt32(total)
	reply.List = list
	return nil
}

//商家收益明细列表
func (*CustHandler) GetBusinessProfitList(ctx context.Context, req *custProto.GetBusinessProfitListRequest, reply *custProto.GetBusinessProfitListReply) error {
	if req.BusinessAccountNo == "" {
		ss_log.Error("BusinessAccountNo参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.OpType != constants.VaOpType_Add && req.OpType != constants.VaOpType_Minus {
		ss_log.Error("参数错误:OpType:%v", req.OpType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

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

	whereList := []*model.WhereSqlCond{
		{Key: "lv.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "lv.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "lv.reason", Val: req.Reason, EqType: "="},
		{Key: "lv.op_type", Val: req.OpType, EqType: "="},
		{Key: "vacc.account_no", Val: req.BusinessAccountNo, EqType: "="},
		{Key: "vacc.balance_type", Val: strings.ToLower(req.CurrencyType), EqType: "="},
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	if req.Reason == "" {
		if req.OpType == constants.VaOpType_Add {
			reason := fmt.Sprintf("lv.reason in (%v,%v)", constants.VaReason_Business_Settle, constants.VaReason_BusinessTransferToBusiness)
			ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, fmt.Sprintf("and %v ", reason))
		} else {
			reason := fmt.Sprintf("lv.reason=%v", constants.VaReason_BusinessRefund)
			ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, fmt.Sprintf("and %v ", reason))
		}
	}

	total := dao.LogVaccountDaoInst.GetCnt(whereModel.WhereStr, whereModel.Args)

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by lv.create_time desc")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	dataList, err := dao.LogVaccountDaoInst.GetBusinessProfitList(whereModel.WhereStr, whereModel.Args)
	if err != nil {
		ss_log.Error("查询商家账户uid:[%v]的虚帐日志失败", req.BusinessAccountNo)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	var list []*custProto.GetBusinessBillData
	for _, v := range dataList {
		data := &custProto.GetBusinessBillData{
			CreateTime:   v.CreateTime,
			LogNo:        v.BizLogNo,
			Amount:       v.Amount,
			Balance:      v.Balance,
			OpType:       v.OpType,
			Reason:       v.Reason,
			CurrencyType: v.CurrencyType,
			VaType:       v.VAccType,
			OrderType:    v.Reason,
		}
		list = append(list, data)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Total = strext.ToInt32(total)
	reply.List = list
	return nil
}

/**
 * WEB获取指定商家账单明细列表
 */
func (*CustHandler) GetBusinessBillList(ctx context.Context, req *custProto.GetBusinessBillListRequest, reply *custProto.GetBusinessBillListReply) error {

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
		{Key: "lv.biz_log_no", Val: req.LogNo, EqType: "like"},
		{Key: "vacc.balance_type", Val: req.CurrencyType, EqType: "="},
		{Key: "vacc.va_type", Val: req.VaType, EqType: "="},
		{Key: "lv.reason", Val: req.Reason, EqType: "="},
		{Key: "vacc.account_no", Val: req.Uid, EqType: "="},
	})
	//只要服务商的
	strs := " and vacc.va_type in (" +
		"'" + strext.ToStringNoPoint(constants.VaType_USD_BUSINESS_SETTLED) + "'," +
		"'" + strext.ToStringNoPoint(constants.VaType_KHR_BUSINESS_SETTLED) + "'," +
		"'" + strext.ToStringNoPoint(constants.VaType_USD_BUSINESS_UNSETTLED) + "'," +
		"'" + strext.ToStringNoPoint(constants.VaType_KHR_BUSINESS_UNSETTLED) + "')"
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, strs)

	//extraStr := getShowBusinessBillExtraStr()
	//ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, extraStr)

	total := dao.LogVaccountDaoInst.GetCnt(whereModel.WhereStr, whereModel.Args)

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by lv.create_time desc,case reason when "+constants.VaReason_FEES+" then 1 end")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.LogVaccountDaoInst.GetBusinessBills(whereModel.WhereStr, whereModel.Args, req.Uid)
	if err != nil {
		ss_log.Error("查询商家账户uid:[%v]的虚帐日志失败", req.Uid)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	return nil
}

//获取服务商账单要显示的白名单
func getShowBusinessBillExtraStr() string {

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

//添加产品签约
func (c *CustHandler) AddBusinessSceneSigned(ctx context.Context, req *custProto.AddBusinessSceneSignedRequest, reply *custProto.AddBusinessSceneSignedReply) error {

	//确认是否已经存在该产品签约正在审核中的签约（只允许有一条正审核中的）
	if !dao.BusinessSceneSignedDaoInst.CheckSceneSignedUnique(req.BusinessNo, req.SceneNo, constants.SignedStatusPending) {
		ss_log.Error("已存在正审核中的产品签约，不允许再申请")
		reply.ResultCode = ss_err.ERR_Business_Scene_Unique_FAILD
		return nil
	}

	//确认产品是否存在
	sceneData, errGetScene := dao.BusinessSceneDaoInst.GetBusinessSceneDetail(req.SceneNo)
	if errGetScene != nil {
		ss_log.Error("产品[%v]不存在", req.SceneNo)
		reply.ResultCode = ss_err.ERR_PAY_UNKOWN_PRODUCT
		return nil
	}

	//判断产品是否可手动签约
	if sceneData.IsManualSigned == constants.ProductIsManualSigned_False {
		ss_log.Error("产品[%v]不可手动签约", req.SceneNo)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//检验经营类目
	indData, err := dao.BusinessIndustryDaoInst.GetBusinessIndustryDetail(req.IndustryNo)
	if err != nil {
		ss_log.Error("查询经营类目的费率和结算周期失败.")
		reply.ResultCode = ss_err.ERR_NoBusinessIndustryRateCycle_FAILD
		return nil
	}

	if indData == nil {
		ss_log.Error("经营类目[%v]不存在", req.IndustryNo)
		reply.ResultCode = ss_err.ERR_NoBusinessIndustryRateCycle_FAILD
		return nil
	}

	if indData.Level == constants.Businesslevel_One {
		ss_log.Error("经营类目[%v]的等级是一级，属于分类，不具备费率和结算周期。不允许选择。", req.IndustryNo)
		reply.ResultCode = ss_err.ERR_NoBusinessIndustryRateCycle_FAILD
		return nil
	}

	//根据code和business_channel_no查询经营类目的费率和结算周期
	rateData, errGetRate := dao.BusinessIndustryRateCycleDaoInst.GetDetail([]*model.WhereSqlCond{
		{Key: "birc.business_channel_no", Val: sceneData.BusinessChannelNo, EqType: "="},
		{Key: "birc.code", Val: req.IndustryNo, EqType: "="},
		{Key: "birc.is_delete", Val: "0", EqType: "="},
	})

	if errGetRate != nil {
		ss_log.Error("获取行业费率、结算周期失败, IndustryNo[%v], BusinessChannelNo[%v], err[%v]", req.IndustryNo, sceneData.BusinessChannelNo, errGetRate)
		reply.ResultCode = ss_err.ERR_NoBusinessIndustryRateCycle_FAILD
		return nil
	}

	if rateData.Id == "" {
		reply.ResultCode = ss_err.ERR_NoBusinessIndustryRateCycle_FAILD
		return nil
	}

	//签约费率 = 基础费率（行业费率） + 产品浮动费率
	rate := ss_count.Add(rateData.Rate, sceneData.FloatRate)
	cycle := rateData.Cycle

	if strext.ToInt(rate) < 0 { //小于0一律按0处理
		rate = "0"
	}

	_, errAdd1 := dao.BusinessSceneSignedDaoInst.AddBusinessSigned(req.AccUid, req.BusinessNo, req.SceneNo, req.IndustryNo, rate, cycle)
	if errAdd1 != nil {
		ss_log.Error("添加签约失败,errAdd[%v]", errAdd1)
		reply.ResultCode = ss_err.ERR_SYS_DB_ADD
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//产品签约的审核（审核通过或不通过）
func (c *CustHandler) UpdateBusinessSceneSignedStatus(ctx context.Context, req *custProto.UpdateBusinessSceneSignedStatusRequest, reply *custProto.UpdateBusinessSceneSignedStatusReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}
	defer ss_sql.Rollback(tx)

	str, _ := util.GetParamZhCn(req.Status, util.AppSignedStatus)

	// 添加关键操作日志
	description := fmt.Sprintf("将商家应用签约产品的签约[%v]状态更改为[%v] ", req.SignedId, str)

	switch req.Status {
	case constants.SignedStatusPassed:
		//通过的则要废除原来的签约（如果有的话），将使用新的产品签约
		data, err := dao.BusinessSceneSignedDaoInst.GetBusinessSceneSignedDetail(req.SignedId)
		if err != nil {
			ss_log.Error("查询要审核的签约信息失败，SignedId[%v] err[%v]", req.SignedId, err)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}

		//产品只允许一个通过的签约。修改以前通过的为作废(如果有的话)
		if !dao.BusinessSceneSignedDaoInst.CheckSceneSignedUnique(data.BusinessNo, data.SceneNo, constants.SignedStatusPassed) {
			if err2 := dao.BusinessSceneSignedDaoInst.SetStatusInvalidTx(tx, data.BusinessNo, data.SceneNo); err2 != nil {
				ss_log.Error("产品只允许一个通过的签约。修改以前通过的为作废失败，err2[%v]", err2)
				reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
				return nil
			}
			//
			description += "并将和该签约一样的商家、产品原来通过的产品签约作废"
		}

	case constants.SignedStatusDeny:
	default:
		ss_log.Error("Status参数[%v]不合法", req.Status)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if err := dao.BusinessSceneSignedDaoInst.UpdateStatusTx(tx, req.SignedId, constants.SignedStatusPending, req.Status, req.Notes); err != nil {
		ss_log.Error("err[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
		return nil
	}

	errAddLog := dao.LogDaoInstance.InsertWebAccountLogTx(tx, description, req.LoginUid, constants.LogAccountWebType_Business)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//修改产品签约的费率、结算周期
func (c *CustHandler) UpdateBusinessSceneSignedInfo(ctx context.Context, req *custProto.UpdateBusinessSceneSignedInfoRequest, reply *custProto.UpdateBusinessSceneSignedInfoReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}
	defer ss_sql.Rollback(tx)

	oldData, errGet := dao.BusinessSceneSignedDaoInst.GetBusinessSceneSignedDetail(req.SignedId)
	if errGet != nil {
		ss_log.Error("查询产品签约的旧数据失败，SignedNo[%v]", req.SignedId)
		reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
		return nil
	}

	if err := dao.BusinessSceneSignedDaoInst.UpdateInfoTx(tx, req.SignedId, req.Cycle, req.Rate); err != nil {
		ss_log.Error("修改产品签约的周期和费率失败，err[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
		return nil
	}

	// 添加关键操作日志
	rate := ss_count.Div(req.Rate, "100").String() + "%"
	oldRate := ss_count.Div(oldData.Rate, "100").String() + "%"
	description := fmt.Sprintf("将商家产品签约[%v]的结算周期由[%v]更改为[%v],费率由[%v]修改为[%v]", req.SignedId, oldData.Cycle, req.Cycle, oldRate, rate)
	if errAddLog := dao.LogDaoInstance.InsertWebAccountLogTx(tx, description, req.LoginUid, constants.LogAccountWebType_Business); errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (c *CustHandler) GetBusinessSceneSignedList(ctx context.Context, req *custProto.GetBusinessSceneSignedListRequest, reply *custProto.GetBusinessSceneSignedListReply) error {

	whereList := []*model.WhereSqlCond{
		{Key: "bss.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "bss.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "bss.status", Val: req.Status, EqType: "="},
		{Key: "acc.account", Val: req.Account, EqType: "like"},
	}

	total := dao.BusinessSceneSignedDaoInst.GetCnt(whereList)

	datas, err := dao.BusinessSceneSignedDaoInst.GetList(whereList, req.Page, req.PageSize)
	if err != nil {
		ss_log.Error("查询商家产品签约数据列表失败，req=[%+v],err=[%v]", req, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	var keys []string
	var langDatas []*custProto.LangData
	keyMap := make(map[string]string) //用于去重，不用重复查询一些key
	for k, data := range datas {
		if data.SceneName != "" { //产品名称记录的是多语言的key
			if _, ok := keyMap[data.SceneName]; !ok { //只有没添加过的才去查询
				keyMap[data.SceneName] = data.SceneName
				keys = append(keys, data.SceneName)
			}
		}

		//一次最多查30个key对应的语言
		if len(keys) == 30 || k == len(datas)-1 {
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

	var datas2 []*custProto.BusinessSceneSignedData
	for _, data := range datas {
		industryName := ""

		switch req.Lang {
		case constants.LangZhCN:
			for _, langData := range langDatas {
				if data.SceneName == langData.Key {
					data.SceneName = langData.LangCh
					break
				}
			}
			industryName = data.IndustryNameCh
		case constants.LangEnUS:
			for _, langData := range langDatas {
				if data.SceneName == langData.Key {
					data.SceneName = langData.LangEn
					break
				}
			}
			industryName = data.IndustryNameEn

		case constants.LangKmKH:
			for _, langData := range langDatas {
				if data.SceneName == langData.Key {
					data.SceneName = langData.LangKm
					break
				}
			}
			industryName = data.IndustryNameKm
		default:
			for _, langData := range langDatas {
				if data.SceneName == langData.Key {
					data.SceneName = langData.LangEn
					break
				}
			}
			industryName = data.IndustryNameEn
		}

		datas2 = append(datas2, &custProto.BusinessSceneSignedData{
			SignedNo:     data.SignedNo,
			SceneName:    data.SceneName,
			IndustryName: industryName,
			Cycle:        data.Cycle,
			Rate:         data.Rate,
			Account:      data.Account,
			Status:       data.Status,
			CreateTime:   data.CreateTime,
			StartTime:    data.StartTime,
			EndTime:      data.EndTime,
			Notes:        data.Notes,
		})
	}

	reply.Datas = datas2
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
