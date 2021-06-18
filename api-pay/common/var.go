package common

const (
	// 是否终止
	INNER_IS_STOP = "isStop"
	// api模式
	INNER_API_MODE = "api_mode"
	// 跟踪号
	INNER_TRACE_NO = "traceNo"
	// 参数
	INNER_PARAM_MAP = "inParamMap"
	//请求body
	REQUEST_BODY = "request_body"

	// 应用的签名方式
	InnerAppSignMethod = "appSignMethod"

	// 应用的ip白名单
	InnerAppIpWhiteList = "appIpWhiteList"

	// 应用的商户的公钥
	InnerAppBusinessPublickKey = "appBusinessPublickKey"

	// 应用的平台的公钥
	InnerAppPlatformPrivateKey = "appPlatformPrivateKey"

	InnerLanguage = "language"

	// 异常时-跳过返回签名步骤
	InnerSkipSign      = "innerSkipSign"
	InnerSkipSignValue = "1"

	RET_CODE         = "_ret_code"
	RET_DATA         = "_ret_data"
	RET_DATA_PRESEND = "_ret_data_presend"

	SignMethod_RSA2 = "RSA2"
	SignMethod_MD5  = "MD5"

	// 请求的有效时间 (单位: 秒)
	RequestExpireTime int64 = 120

	// 请求参数
	ParamsSignType  = "sign_type" // 签名类型
	ParamsSign      = "sign"      // 签名
	ParamsLang      = "lang"      // 请求的语言类型
	ParamsTimestamp = "timestamp" // 请求的时间戳
	ParamsVersion   = "version"   // 接口的版本号

	ParamsAppId = "app_id" // 应用id

	// 固定返回的参数
	RetFieldCode = "code"
	RetFieldMsg  = "msg"
	RetFieldSgin = "sign"
)
