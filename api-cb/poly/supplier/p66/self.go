package p66

import (
	"a.a/cu/encrypt"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/api-cb/dao"
	"a.a/mp-server/api-cb/m"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"encoding/base64"
	"strings"
)

var PolyP66Inst PolyP66

type PolyP66 struct {
}

func (s PolyP66) Callback(req *m.PolyCallbackReq) *m.PolyCallbackResp {
	m2 := req.RecvMap
	sign := strext.ToStringNoPoint(m2["sign"])
	sign2, _ := base64.StdEncoding.DecodeString(sign)
	delete(m2, "sign")

	channelNo, amountT := dao.BillDaoInst.GetBillChannelNo(strext.ToStringNoPoint(m2["out_trade_no"]))
	param := dao.ChannelDaoInst.GetChannelParam(channelNo)
	signBefore := encrypt.Map2FormStr(m2, param.Key1, "&key=",
		encrypt.FIELD_ENCODED_NONE, nil, "", false)
	md5Str := strings.ToUpper(encrypt.DoMd5(signBefore))
	_, err := encrypt.DoRsa(encrypt.HANDLE_VERIFY, encrypt.KEYTYPE_PKCS8, encrypt.RETFMT_NONE, encrypt.HASHLENTYPE_SHA256,
		map[string]interface{}{
			"sign": []byte(sign2),
			"data": []byte(md5Str),
		}, param.Key2,
		encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_SIGN_MAP)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil
	}

	// ??? 协议没有定义失败
	orderStatus := constants.OrderStatus_Paid
	if m2["status"] == "success" {
		orderStatus = constants.OrderStatus_Paid
	}

	amount := strext.ToStringNoPoint((strext.ToFloat64(m2["origin_amount"]) + strext.ToFloat64(m2["rand"])) * 100)
	if amount != amountT {
		ss_log.Error("err=[%v|金额对不上]", m2["out_trade_no"])
		return nil
	}

	return &m.PolyCallbackResp{
		Amount:       amount,
		InnerOrderNo: strext.ToStringNoPoint(m2["out_trade_no"]),
		UpperOrderNo: strext.ToStringNoPoint(m2["trade_no"]),
		UpdateTime:   util.NowWithFmt("2006-01-02 15:04:05"),
		RetCode:      ss_err.ERR_SUCCESS,
		RetMsg:       "成功",
		OrderStatus:  orderStatus,
		RetBody:      "success",
	}
}

func (s PolyP66) TransferCallback(req *m.PolyTransferCallbackReq) *m.PolyTransferCallbackResp {
	m2 := req.RecvMap
	sign := strext.ToStringNoPoint(m2["sign"])
	sign2, _ := base64.StdEncoding.DecodeString(sign)
	delete(m2, "sign")

	channelNo, amountT := dao.BillDaoInst.GetBillChannelNo(strext.ToStringNoPoint(m2["out_trade_no"]))
	param := dao.ChannelDaoInst.GetChannelParam(channelNo)
	signBefore := encrypt.Map2FormStr(m2, param.Key1, "&key=",
		encrypt.FIELD_ENCODED_NONE, nil, "", false)
	md5Str := strings.ToUpper(encrypt.DoMd5(signBefore))
	_, err := encrypt.DoRsa(encrypt.HANDLE_VERIFY, encrypt.KEYTYPE_PKCS8, encrypt.RETFMT_NONE, encrypt.HASHLENTYPE_SHA256,
		map[string]interface{}{
			"sign": []byte(sign2),
			"data": []byte(md5Str),
		}, param.Key2,
		encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_SIGN_MAP)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil
	}

	// ??? 协议没有定义失败
	orderStatus := constants.OrderStatus_Paid
	if m2["status"] == "success" {
		orderStatus = constants.OrderStatus_Paid
	}

	amount := strext.ToStringNoPoint((strext.ToFloat64(m2["origin_amount"]) + strext.ToFloat64(m2["rand"])) * 100)
	if amount != amountT {
		ss_log.Error("err=[%v|金额对不上]", m2["out_trade_no"])
		return nil
	}

	return &m.PolyTransferCallbackResp{
		Amount:       amount,
		InnerOrderNo: strext.ToStringNoPoint(m2["out_trade_no"]),
		UpperOrderNo: strext.ToStringNoPoint(m2["trade_no"]),
		UpdateTime:   util.NowWithFmt("2006-01-02 15:04:05"),
		RetCode:      ss_err.ERR_SUCCESS,
		RetMsg:       "成功",
		OrderStatus:  orderStatus,
		RetBody:      "success",
	}
}
