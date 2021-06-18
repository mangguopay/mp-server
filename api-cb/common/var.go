package common

import "github.com/micro/go-micro/v2/broker"

const (
	// req里的data字段
	INNER_DATA = "inData"
	// 封装好的req
	INNER_REQ = "inReq"
	// 中转的错误信息
	INNER_ERR_REASON = "errReason"
	// ext
	INNER_ERR_REASON_EXTMSG = "errReasonExtMsg"
	// 是否终止
	INNER_IS_STOP = "isStop"
	// api模式
	INNER_API_MODE = "api_mode"
	// 跟踪号
	INNER_TRACE_NO = "traceNo"
	// 参数
	INNER_PARAM_MAP = "inParamMap"
	//
	INNER_FMT = "bodyFmt"
	// 签名方式
	INNER_SIGN_METHOD = "signMethod"
	// bool|验签是否通过 true/false
	INNER_SIGN_VERIFY = "signVerify"
	//
	INNER_RET_DATA = "retData"
	//=======================================
	// 商户号
	INNER_DATA_ACCNO = "acc_no"
	// 签名
	INNER_DATA_SIGN = "sign"

	RET_CODE         = "_ret_code"
	RET_MSG          = "_ret_msg"
	RET_DATA         = "_ret_data"
	RET_DATA_PRESEND = "_ret_data_presend"
	RET_TOKEN        = "_ret_token"
	RET_STATUS_CODE  = "_ret_statusCode"

	SignMethod_Rsa = "rsa"
	SignMethod_Md5 = "md5"
)

var (
	MqPushMsg *broker.Broker
)
