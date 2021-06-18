package common

import "github.com/wiwii/base64Captcha"

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
	INNER_IS_JWT_CHECKED  = "isJwtChecked"
	INNER_PRESEND_ENCODED = "presend_encoded"
	//=====================================
	RET_CODE         = "_ret_code"
	RET_MSG          = "_ret_msg"
	RET_DATA         = "_ret_data"
	RET_DATA_PRESEND = "_ret_data_presend"
	Login_URI        = "/webadmin/auth/login"
)

var ConfigC = base64Captcha.ConfigCharacter{
	Height: 44,
	Width:  156,
	//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
	Mode:                 base64Captcha.CaptchaModeAlphabetUpper,
	ComplexOfNoiseText:   base64Captcha.CaptchaComplexLower,
	ComplexOfNoiseDot:    base64Captcha.CaptchaComplexLower,
	IsShowHollowLine:     false,
	IsShowNoiseDot:       false,
	IsShowNoiseText:      false,
	IsShowSlimeLine:      false,
	IsShowSineLine:       false,
	IsUseCustomFontColor: true,
	FontColorZone: [][]int{
		[]int{255, 0},
		[]int{255, 0},
		[]int{255, 0},
	},
	FontSizeMin:   0.2,
	CaptchaLen:    4,
	IsUseFontName: true,
	FontName:      "fonts/DeborahFancyDress.ttf",
}
