package handler

import (
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/cust-srv/dao"
	"context"
)

//提现订单列表
func (c *CustHandler) GetToBusinessList(ctx context.Context, req *go_micro_srv_cust.GetToBusinessListRequest, reply *go_micro_srv_cust.GetToBusinessListReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "tb.log_no", Val: req.LogNo, EqType: "="},
		{Key: "tb.create_time", Val: req.StartTime, EqType: ">="},
		{Key: "tb.create_time", Val: req.EndTime, EqType: "<="},
		{Key: "tb.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "tb.business_no", Val: req.IdenNo, EqType: "="},
		{Key: "tb.currency_type", Val: req.MoneyType, EqType: "="},
		{Key: "acc.account", Val: req.Account, EqType: "like"}, //商家账号
	}

	total := dao.LogToBusinessDaoInst.GetToBusinessCnt(whereList)

	if total == "" || total == "0" {
		reply.Total = strext.ToInt32(total)
		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	}

	datas, err := dao.LogToBusinessDaoInst.GetToBusinessList(whereList, req.Page, req.PageSize)
	if err != nil {
		ss_log.Error("查询商家提现数据列表失败，req=[%+v],err=[%v]", req, err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	//获取图片的url
	for _, data := range datas {
		ids := []string{
			data.ImageId,
		}
		imgUrls := dao.ImageDaoInstance.GetImgUrlsByImgIds(ids)
		data.ImageUrl = imgUrls[0]
	}

	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//提现订单详情
func (c *CustHandler) GetBusinessToWithdrawDetail(ctx context.Context, req *go_micro_srv_cust.GetBusinessToWithdrawDetailRequest, reply *go_micro_srv_cust.GetBusinessToWithdrawDetailReply) error {
	if req.LogNo == "" {
		ss_log.Error("LogNo参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	data, err := dao.LogToBusinessDaoInst.GetToBusinessDetail(req.LogNo)
	if err != nil {
		ss_log.Error("查询商家提现订单详情失败，logNo=%v, err=%v", req.LogNo, err)
		reply.ResultCode = ss_err.ERR_SYSTEM
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

//获取商家余额
func (c *CustHandler) GetBusinessVAccBalance(ctx context.Context, req *go_micro_srv_cust.GetBusinessVAccBalanceRequest, reply *go_micro_srv_cust.GetBusinessVAccBalanceReply) error {
	if req.BusinessAccountNo == "" {
		ss_log.Error("BusinessAccountNo参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.MoneyType == "" {
		ss_log.Error("MoneyType参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//查询商家余额和冻结金额
	businessVAccType := global.GetBusinessVAccType(req.MoneyType, true)
	errCode, balance := dao.VaccountDaoInst.GetBalanceFromAccNo(req.BusinessAccountNo, businessVAccType)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("查询商家可用金额失败，errCode=%v", errCode)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	errCode2, frozenBalance := dao.VaccountDaoInst.GetFrozenBalanceFromAccNo(req.BusinessAccountNo, businessVAccType)
	if errCode2 != ss_err.ERR_SUCCESS {
		ss_log.Error("查询商家冻结金额失败，errCode=%v", errCode2)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//获取商户虚账
	businessVAccNo, err := dao.VaccountDaoInst.GetVaccountNo(req.BusinessAccountNo, strext.ToInt32(businessVAccType))
	if err != nil {
		ss_log.Error("查询商户虚账失败，accountNo=%v, moneyType=%v", req.BusinessAccountNo, req.MoneyType)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//昨日
	startTime := ss_time.Now(global.Tz).AddDate(0, 0, -1).Format(ss_time.DateFormat) + " 00:00:00"
	endTime := ss_time.Now(global.Tz).AddDate(0, 0, -1).Format(ss_time.DateFormat) + " 23:59:59"

	//入账
	recordedAmount, err := dao.LogVaccountDaoInst.GetBusinessRecordedAmount(businessVAccNo, startTime, endTime)
	if err != nil {
		ss_log.Error("查询商户昨日入账金额失败，accountNo=%v, moneyType=%v", req.BusinessAccountNo, req.MoneyType)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	//出账
	expenditureAmount, err := dao.LogVaccountDaoInst.GetBusinessExpenditureAmount(businessVAccNo, startTime, endTime)
	if err != nil {
		ss_log.Error("查询商户昨日出账金额失败，accountNo=%v, moneyType=%v", req.BusinessAccountNo, req.MoneyType)
		reply.ResultCode = ss_err.ERR_SYSTEM
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.MoneyType = req.MoneyType
	reply.Balance = balance
	reply.FrozenBalance = frozenBalance
	reply.RecordedAmount = recordedAmount
	reply.ExpenditureAmount = expenditureAmount
	return nil
}
