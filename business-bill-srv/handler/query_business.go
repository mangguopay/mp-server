package handler

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/business-bill-srv/dao"
	"a.a/mp-server/common/constants"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"a.a/mp-server/common/ss_err"
	"context"
	"database/sql"
)

// 通过固定收款二维码获取商户应用信息
func (*BusinessBillHandler) GetBusinessAppInfo(ctx context.Context, req *businessBillProto.GetBusinessAppInfoRequest, reply *businessBillProto.GetBusinessAppInfoReply) error { // 0.检查
	if req.FixedQrcode == "" {
		ss_log.Error("FixedQrcode参数为空")
		reply.ResultCode = ss_err.FixedQrCodeIsEmpty
		reply.Msg = ss_err.GetMsg(ss_err.FixedQrCodeIsEmpty, req.Lang)
		return nil
	}

	// 查询应用信息
	app, err := dao.BusinessAppDaoInst.GetAppInfoByFixedQrCode(req.FixedQrcode)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("二维码不存在,FixedQrcode:%v", req.FixedQrcode)
			reply.ResultCode = ss_err.QrCodeNotInvalid
			reply.Msg = ss_err.GetMsg(ss_err.QrCodeNotInvalid, req.Lang)
			return nil
		}
		ss_log.Error("二维码查询应用失败,FixedQrcode=%v, err=%v", req.FixedQrcode, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	if app.Status != constants.BusinessAppStatus_Up {
		ss_log.Error("应用状态不是上架状态,FixedQrcode:%s, status:%v", req.FixedQrcode, app.Status)
		reply.ResultCode = ss_err.QrCodeNotAvailable
		reply.Msg = ss_err.GetMsg(ss_err.QrCodeNotAvailable, req.Lang)
		return nil
	}

	reply.ResultCode = ss_err.Success
	reply.Msg = ss_err.GetMsg(ss_err.Success, req.Lang)
	reply.AppName = app.SimplifyName // 现在显示商家简称
	return nil
}

// 通过固定收款二维码获取商户应用信息
func (*BusinessBillHandler) GetPersonalBusinessInfo(ctx context.Context, req *businessBillProto.GetPersonalBusinessInfoRequest, reply *businessBillProto.GetPersonalBusinessInfoReply) error { // 0.检查
	if req.FixedCode == "" {
		ss_log.Error("FixedQrcode参数为空")
		reply.ResultCode = ss_err.FixedQrCodeIsEmpty
		reply.Msg = ss_err.GetMsg(ss_err.FixedQrCodeIsEmpty, req.Lang)
		return nil
	}

	businessNo, err := dao.BusinessFixedCodeDaoInst.GetBusinessByCode(req.FixedCode)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("二维码不存在, FixedCode=%v", req.FixedCode)
			reply.ResultCode = ss_err.QrCodeNotExist
			reply.Msg = ss_err.GetMsg(ss_err.QrCodeNotExist, req.Lang)
			return nil
		}
		ss_log.Error("查询二维码失败, FixedCode=%v, err=%v", req.FixedCode, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	business, err := dao.BusinessDaoInst.GetTransConfig(businessNo)
	if err != nil {
		ss_log.Error("查询商户信息失败, BusinessNo=%v, err=%v", businessNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	resultCode, err := CheckBusinessIncomeAuth(business)
	if err != nil {
		ss_log.Error("商户检测不通过, BusinessNo=%v, err=%v", businessNo, err)
		reply.ResultCode = resultCode
		reply.Msg = ss_err.GetMsg(resultCode, req.Lang)
		return nil
	}

	reply.ResultCode = ss_err.Success
	reply.Msg = ss_err.GetMsg(ss_err.Success, req.Lang)
	reply.BusinessName = business.SimplifyName
	return nil
}
