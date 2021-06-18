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
	// 是否出错了
	INNER_IS_FALSE = "isFailed"
	// 是否终止
	INNER_IS_STOP = "isStop"
	// api模式
	INNER_API_MODE = "api_mode"
	// 跟踪号
	INNER_TRACE_NO = "traceNo"
	// 参数
	INNER_PARAM_MAP = "params"
	//
	INNER_FMT = "bodyFmt"
	// 签名方式
	INNER_SIGN_METHOD = "signMethod"
	// bool|验签是否通过 true/false
	INNER_SIGN_VERIFY    = "signVerify"
	INNER_IS_ENCODED     = "isEncoded"
	INNER_POST_RSA_PARAM = "postRsaParam"
	// 签名
	INNER_SIGN     = "sign"
	INNER_JWT_DATA = "decodedJwt"
	// jwt合法？
	INNER_IS_JWT_CHECKED   = "isJwtChecked"
	INNER_PRESEND_ENCODED  = "presend_encoded"
	INNER_PRESEND_ENCODED2 = "presend_encoded2"
	//=====================================
	RET_CODE         = "_ret_code"
	RET_MSG          = "_ret_msg"
	RET_DATA         = "_ret_data"
	RET_DATA_PRESEND = "_ret_data_presend"

	Login_URI = "/mobile/auth/login"

	Pub_Key = "pub_key"

	SKEY_PlatPri = `pri_key`         // 平台私钥
	SKEY_PlatPub = `pub_key`         // 平台公钥
	SKEY_DefPri  = `default_pri_key` // 默认私钥
	SKEY_DefPub  = `default_pub_key` // 默认公钥
)

var (
	MqPushMsg  *broker.Broker
	EncryptMap map[string]interface{}
)
