package ss_err

import (
	"a.a/mp-server/common/constants"
)

const (
	ACErrSuccess = "0" // 成功

	ACErrSysErr           = "AC000100" // 系统内部错误
	ACErrSysBusy          = "AC000101" // 系统繁忙
	ACErrSysMethodErr     = "AC000102" // 请求方法不支持
	ACErrSysRouteNotFound = "AC000103" // 请求路径错误
	ACErrVerifySignFail   = "AC000104" // 验签失败

	ACErrSysRetAbnormal = "AC000110" // 数据返回异常

	ACErrParamEmpty         = "AC000200" // 参数为空
	ACErrParamMiss          = "AC000201" // 缺少参数
	ACErrParamErr           = "AC000202" // 参数错误
	ACErrPostReadBodyErr    = "AC000203" // 读取body失败
	ACErrPostEmptyBody      = "AC000204" // 请求body为空
	ACErrPostBodyNotJson    = "AC000205" // 请求body格式不是json格式
	ACErrMissParamAppId     = "AC000206" // 缺少app_id参数
	ACErrMissParamSign      = "AC000207" // 缺少sign参数
	ACErrMissParamSignType  = "AC000208" // 缺少sign_type参数
	ACErrMissParamTimestamp = "AC000209" // 缺少timestamp参数
	ACErrSignTypeNotMatch   = "AC000210" // 参数sign_type与配置不匹配
	ACErrSignTypeErr        = "AC000211" // 参数sign_type错误
	ACErrIpWhiteListForbid  = "AC000212" // IP白名单限制
	ACErrParamTimestampErr  = "AC000213" // 参数timestamp错误

	ACErrAppNotExists     = "AC000301" // 应用不存在
	ACErrAppNotPutOn      = "AC000302" // 应用未上架
	ACErrAppConfigErr     = "AC000303" // 应用配置错误
	ACErrVerifySignFailed = "AC000305" // 验签失败
	ACErrRequestExpired   = "AC000306" // 请求已过期
)

var apiPayEnUSMap = make(map[string]string)
var apiPaykmKHMap = make(map[string]string)
var apiPayzhCNMap = make(map[string]string)

func init() {
	initApiPayEnUSMap()
	initApiPaykmKHMap()
	initApiPayzhCNMap()
}

func initApiPayEnUSMap() {
	apiPayEnUSMap[ACErrSuccess] = "Success"

	apiPayEnUSMap[ACErrSysErr] = "Internal system errors"
	apiPayEnUSMap[ACErrSysBusy] = "System Busy"
	apiPayEnUSMap[ACErrSysMethodErr] = "The request method is not supported."
	apiPayEnUSMap[ACErrSysRouteNotFound] = "Request path error"
	apiPayEnUSMap[ACErrSysRetAbnormal] = "Data Return Exception"
	apiPayEnUSMap[ACErrVerifySignFail] = "Sign verification failed"

	apiPayEnUSMap[ACErrParamEmpty] = "The parameter is empty."
	apiPayEnUSMap[ACErrParamMiss] = "Lack of parameters"
	apiPayEnUSMap[ACErrParamErr] = "parameter error"
	apiPayEnUSMap[ACErrPostReadBodyErr] = "Failed to read the body"
	apiPayEnUSMap[ACErrPostEmptyBody] = "The request body is empty."
	apiPayEnUSMap[ACErrPostBodyNotJson] = "Request body format not json format"

	apiPayEnUSMap[ACErrMissParamAppId] = "Lack of app_id parameter"
	apiPayEnUSMap[ACErrMissParamSign] = "Lack of sign parameter"
	apiPayEnUSMap[ACErrMissParamTimestamp] = "Lack of timestamp parameter"
	apiPayEnUSMap[ACErrMissParamSignType] = "Lack of sign_type parameter"
	apiPayEnUSMap[ACErrSignTypeNotMatch] = "The parameter sign_type does not match the configuration."
	apiPayEnUSMap[ACErrSignTypeErr] = "Parameter sign_type error"
	apiPayEnUSMap[ACErrIpWhiteListForbid] = "IP whitelisting restrictions"
	apiPayEnUSMap[ACErrParamTimestampErr] = "Parameter timestamp error"

	apiPayEnUSMap[ACErrAppNotExists] = "The application does not exist."
	apiPayEnUSMap[ACErrAppNotPutOn] = "Application not yet available"
	apiPayEnUSMap[ACErrAppConfigErr] = "Application configuration errors"

	apiPayEnUSMap[ACErrVerifySignFailed] = "Sign verification failed"
	apiPayEnUSMap[ACErrRequestExpired] = "The request has expired"
}

func initApiPaykmKHMap() {
	apiPaykmKHMap[ACErrSuccess] = "Success"

	apiPaykmKHMap[ACErrSysErr] = "កំហុសប្រព័ន្ធខាងក្នុង"
	apiPaykmKHMap[ACErrSysBusy] = "ប្រព័ន្ធទំនាក់ទំនងកំពុងរវល់"
	apiPaykmKHMap[ACErrSysMethodErr] = "មិនគាំទ្រវិធីសាស្ត្រស្នើសុំ"
	apiPaykmKHMap[ACErrSysRouteNotFound] = "ស្នើសុំផ្លូវមានកំហុស"
	apiPaykmKHMap[ACErrSysRetAbnormal] = "ទិន្នន័យត្រឡប់មកវិញមិនធម្មតា"
	apiPaykmKHMap[ACErrVerifySignFail] = "Sign verification failed"

	apiPaykmKHMap[ACErrParamEmpty] = "ការបញ្ជូលេខគឺទទេ"
	apiPaykmKHMap[ACErrParamMiss] = "ខ្វះខាតការបញ្ជូលេខ"
	apiPaykmKHMap[ACErrParamErr] = "ការបញ្ជូលេខខុស"
	apiPaykmKHMap[ACErrPostReadBodyErr] = "អានbodyទទេ"
	apiPaykmKHMap[ACErrPostEmptyBody] = "ស្ថាប័នស្នើសុំbodyទទេ"
	apiPaykmKHMap[ACErrPostBodyNotJson] = "ទំរង់បែបបទស្នើសុំbodyមិនមែនជាទម្រង់jsonទេ"

	apiPaykmKHMap[ACErrMissParamAppId] = "ខ្វះខាតតួរអក្សរapp_id"
	apiPaykmKHMap[ACErrMissParamSign] = "ខ្វះខាតតួរអក្សរsign"
	apiPaykmKHMap[ACErrMissParamTimestamp] = "ខ្វះខាតតួរអក្សរtimestamp"
	apiPaykmKHMap[ACErrMissParamSignType] = "ខ្វះខាតតួរអក្សរsign_type"
	apiPaykmKHMap[ACErrSignTypeNotMatch] = "ខ្វះការបញ្ជូលេខsign_typeមិនត្រូវនឹងការកំណត់រចនាសម្ព័ន្ធ"
	apiPaykmKHMap[ACErrSignTypeErr] = "ខ្វះការបញ្ជូលេខsign_typeមានកំហុស"
	apiPaykmKHMap[ACErrIpWhiteListForbid] = "IPការរឹតបន្តឹងបញ្ជីស"
	apiPaykmKHMap[ACErrParamTimestampErr] = "ខ្វះការបញ្ជូលេខtimestampមានកំហុស"

	apiPaykmKHMap[ACErrAppNotExists] = "មិនមានពាក្យសុំកម្មវិធីទេ"
	apiPaykmKHMap[ACErrAppNotPutOn] = "កម្មវិធីមិនមានលើកឡើងទេ"
	apiPaykmKHMap[ACErrAppConfigErr] = "កំហុសក្នុងការកំណត់រចនាសម្ព័ន្ធកម្មវិធី"

	apiPaykmKHMap[ACErrVerifySignFailed] = "Sign verification failed"
	apiPaykmKHMap[ACErrRequestExpired] = "សំណើបានផុតកំណត់"
}

func initApiPayzhCNMap() {
	apiPayzhCNMap[ACErrSuccess] = "Success"

	apiPayzhCNMap[ACErrSysErr] = "系统内部错误"
	apiPayzhCNMap[ACErrSysBusy] = "系统繁忙"
	apiPayzhCNMap[ACErrSysMethodErr] = "请求方法不支持"
	apiPayzhCNMap[ACErrSysRouteNotFound] = "请求路径错误"
	apiPayzhCNMap[ACErrSysRetAbnormal] = "数据返回异常"
	apiPayzhCNMap[ACErrVerifySignFail] = "验签失败"

	apiPayzhCNMap[ACErrParamEmpty] = "参数为空"
	apiPayzhCNMap[ACErrParamMiss] = "缺少参数"
	apiPayzhCNMap[ACErrParamErr] = "参数错误"
	apiPayzhCNMap[ACErrPostReadBodyErr] = "读取body失败"
	apiPayzhCNMap[ACErrPostEmptyBody] = "请求body为空"
	apiPayzhCNMap[ACErrPostBodyNotJson] = "请求body格式不是json格式"

	apiPayzhCNMap[ACErrMissParamAppId] = "缺少app_id参数"
	apiPayzhCNMap[ACErrMissParamSign] = "缺少sign参数"
	apiPayzhCNMap[ACErrMissParamTimestamp] = "缺少timestamp参数"
	apiPayzhCNMap[ACErrMissParamSignType] = "缺少sign_type参数"
	apiPayzhCNMap[ACErrSignTypeNotMatch] = "参数sign_type与配置不匹配"
	apiPayzhCNMap[ACErrSignTypeErr] = "参数sign_type错误"
	apiPayzhCNMap[ACErrIpWhiteListForbid] = "IP白名单限制"
	apiPayzhCNMap[ACErrParamTimestampErr] = "参数timestamp错误"

	apiPayzhCNMap[ACErrAppNotExists] = "应用不存在"
	apiPayzhCNMap[ACErrAppNotPutOn] = "应用未上架"
	apiPayzhCNMap[ACErrAppConfigErr] = "应用配置错误"

	apiPayzhCNMap[ACErrVerifySignFailed] = "验签失败"
	apiPayzhCNMap[ACErrRequestExpired] = "请求已过期"
}

func GetPayApiErrMsg(retCode, lang string) string {
	msg := ""

	switch lang {
	case constants.LangEnUS:
		msg = apiPayEnUSMap[retCode]
	case constants.LangKmKH:
		msg = apiPaykmKHMap[retCode]
	case constants.LangZhCN:
		msg = apiPayzhCNMap[retCode]
	default:
		msg = retCode
	}

	return msg
}
