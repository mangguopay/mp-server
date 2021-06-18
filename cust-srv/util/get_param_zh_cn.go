package util

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"fmt"
)

const (
	CollectionType           = "CollectionType"           //收款方式
	ExchangeType             = "ExchangeType"             //兑换类型
	AuditOrderStatus         = "AuditOrderStatus"         //需要审核的订单的状态
	OrderStatus              = "OrderStatus"              //需要审核的订单的状态
	AuthStatus               = "AuthStatus"               //账号的实名认证状态
	ChargeType               = "ChargeType"               //计算手续费类型(1-按比例收取手续费，2按单笔手续费收取)
	SupportType              = "SupportType"              //该渠道支持的业务类型（1-只支持存款，2只支持取款，3两个都支持）
	IsForce                  = "IsForce"                  //版本是否强制更新
	CurrencyType             = "CurrencyType"             //币种
	IsRecom                  = "IsRecom"                  //是否推荐(1-推荐，0-不推荐)
	Lang                     = "Lang"                     //语言
	AgreementType            = "AgreementType"            //协议类型
	InAuthorization          = "InAuthorization"          //用户充值权限
	OutAuthorization         = "OutAuthorization"         //用户提现权限
	InTransferAuthorization  = "InTransferAuthorization"  //用户可转账入权限
	OutTransferAuthorization = "OutTransferAuthorization" //用户可转账出权限
	IncomeAuthorization      = "IncomeAuthorization"      //服务商收款权限
	OutgoAuthorization       = "OutgoAuthorization"       //服务商取款权限
	UseStatus                = "UseStatus"                //使用状态
	System                   = "System"                   //版本系统Android、ios
	VsType                   = "VsType"                   //版本包类型moderpay app、moderpay pos、mangopay app、mangopay pos
	AppSignedStatus          = "AppSignedStatus"          //签约状态(1申请中,2未通过,3已通过,4已过期)
	ChannelType              = "ChannelType"              //用户、商家银行卡渠道类型（1-普通渠道，2-第三方渠道）
)

// 校验参数是否合法,合法返回中文和true,不合法则返回false
func GetParamZhCn(paramStr, typeStr string) (ParamZhCn string, legal bool) {
	if paramStr == "" {
		return " 参数为空", false
	}

	switch typeStr {
	case CollectionType: //收款方式
		switch paramStr { //收款方式,1-支票;2-现金;3-银行转账;4-其他
		case constants.COLLECTION_TYPE_CHECK:
			ParamZhCn = "支票"
		case constants.COLLECTION_TYPE_CASH:
			ParamZhCn = "现金"
		case constants.COLLECTION_TYPE_BANK_TRANSFER:
			ParamZhCn = "银行转账"
		case constants.COLLECTION_TYPE_OTHER:
			ParamZhCn = "其他"
		default:
			ParamZhCn = fmt.Sprintf(" 参数错误[%v]", paramStr)
			return ParamZhCn, false
		}
	case ExchangeType: //兑换类型
		switch paramStr {
		case constants.Exchange_Usd_To_Khr:
			ParamZhCn = "USD兑换KHR"
		case constants.Exchange_Khr_To_Usd:
			ParamZhCn = "KHR兑换USD"
		default:
			ParamZhCn = fmt.Sprintf(" 参数错误[%v]", paramStr)
			return ParamZhCn, false
		}
	case AuditOrderStatus: //需要审核的订单的状态
		switch paramStr {
		case constants.AuditOrderStatus_Pending:
			ParamZhCn = "等待"
		case constants.AuditOrderStatus_Passed:
			ParamZhCn = "通过"
		case constants.AuditOrderStatus_Deny:
			ParamZhCn = "不通过"
		default:
			ParamZhCn = fmt.Sprintf(" 参数错误[%v]", paramStr)
			return ParamZhCn, false
		}
	case OrderStatus:
		switch paramStr {
		case constants.OrderStatus_Init:
			ParamZhCn = "初始化"
		case constants.OrderStatus_Pending:
			ParamZhCn = "等待"
		case constants.OrderStatus_Paid:
			ParamZhCn = "已支付"
		case constants.OrderStatus_Err:
			ParamZhCn = "失败"
		case constants.OrderStatus_Pending_Confirm:
			ParamZhCn = "待确认"
		case constants.OrderStatus_Cancel:
			ParamZhCn = "取消"
		default:
			ParamZhCn = fmt.Sprintf(" 参数错误[%v]", paramStr)
			return ParamZhCn, false
		}
	case AuthStatus: //	账号的实名认证状态
		switch paramStr {
		case constants.AuthMaterialStatus_Pending:
			ParamZhCn = "未审核"
		case constants.AuthMaterialStatus_Passed:
			ParamZhCn = "通过"
		case constants.AuthMaterialStatus_Deny:
			ParamZhCn = "不通过"
		case constants.AuthMaterialStatus_UnAuth:
			ParamZhCn = "未认证"
		case constants.AuthMaterialStatus_Appeal_Passed:
			ParamZhCn = "通过的认证作废"
		default:
			ParamZhCn = fmt.Sprintf(" 参数错误[%v]", paramStr)
			return ParamZhCn, false
		}
	case ChargeType: //	计算手续费类型(1-按比例收取手续费，2按单笔手续费收取)
		switch paramStr {
		case constants.Fee_Charge_Type_Rate:
			ParamZhCn = "按比例收取手续费"
		case constants.Fee_Charge_Type_Count:
			ParamZhCn = "按单笔手续费收取"
		default:
			ParamZhCn = fmt.Sprintf(" 参数错误[%v]", paramStr)
			return ParamZhCn, false
		}
	case SupportType: //	该渠道支持的业务类型（1-只支持存款，2只支持取款，3两个都支持）
		switch paramStr {
		case constants.SupportType_In:
			ParamZhCn = "只支持存款"
		case constants.SupportType_Out:
			ParamZhCn = "只支持取款"
		case constants.SupportType_Common:
			ParamZhCn = "通用"
		default:
			ParamZhCn = fmt.Sprintf(" 参数错误[%v]", paramStr)
			return ParamZhCn, false
		}
	case IsForce: //	版本是否强制更新
		switch paramStr {
		case constants.AppVersionIsForce_False:
			ParamZhCn = "不强制"
		case constants.AppVersionIsForce_True:
			ParamZhCn = "强制"
		default:
			ParamZhCn = fmt.Sprintf(" 参数错误[%v]", paramStr)
			return ParamZhCn, false
		}
	case CurrencyType: //	币种
		switch paramStr {
		case constants.CURRENCY_USD:
			ParamZhCn = "USD"
		case constants.CURRENCY_KHR:
			ParamZhCn = "KHR"
		case constants.CURRENCY_UP_USD:
			ParamZhCn = "USD"
		case constants.CURRENCY_UP_KHR:
			ParamZhCn = "KHR"
		default:
			ParamZhCn = fmt.Sprintf(" 参数错误[%v]", paramStr)
			return ParamZhCn, false
		}
	case IsRecom: //	是否推荐(1-推荐，0-不推荐)
		switch paramStr {
		case constants.IsRecom_True:
			ParamZhCn = "是"
		case constants.IsRecom_False:
			ParamZhCn = "否"
		default:
			ParamZhCn = fmt.Sprintf(" 参数错误[%v]", paramStr)
			return ParamZhCn, false
		}
	case Lang:
		switch paramStr {
		case constants.LangZhCN:
			ParamZhCn = "中文"
		case constants.LangKmKH:
			ParamZhCn = "柬埔寨"
		case constants.LangEnUS:
			ParamZhCn = "英语"
		default:
			ParamZhCn = fmt.Sprintf(" 参数错误[%v]", paramStr)
			return ParamZhCn, false
		}
	case AgreementType:
		switch paramStr {
		case constants.AgreementType_Use: //用户协议
			ParamZhCn = "用户协议"
		case constants.AgreementType_Privacy: //隐私协议
			ParamZhCn = "隐私协议"
		case constants.AgreementType_Auth_Material: //实名认证协议
			ParamZhCn = "实名认证协议"
		default:
			ParamZhCn = fmt.Sprintf(" 参数错误[%v]", paramStr)
			return ParamZhCn, false
		}
	case InAuthorization: //用户充值权限
		fallthrough
	case OutAuthorization: //用户提现权限
		fallthrough
	case InTransferAuthorization: //用户可转账入权限
		fallthrough
	case OutTransferAuthorization: //用户可转账出权限
		fallthrough
	case IncomeAuthorization: //服务商收款权限
		fallthrough
	case OutgoAuthorization: //服务商取款权限
		fallthrough
	case UseStatus:
		switch paramStr {
		case constants.Status_Disable: //0禁用
			ParamZhCn = "禁用"
		case constants.Status_Enable: //1启用
			ParamZhCn = "启用"
		default:
			ParamZhCn = fmt.Sprintf("参数错误 %v:[%v]", typeStr, paramStr)
			return ParamZhCn, false
		}
	case System:
		switch paramStr {
		case constants.AppVersionSystem_Ios: //ios
			ParamZhCn = "ios"
		case constants.AppVersionSystem_Android: //Android
			ParamZhCn = "android"
		default:
			ParamZhCn = fmt.Sprintf("参数错误[%v]", paramStr)
			return ParamZhCn, false
		}
	case VsType:
		switch paramStr {
		case constants.AppVersionVsType_app: //app
			ParamZhCn = "moderpay-app"
		case constants.AppVersionVsType_pos: //pos
			ParamZhCn = "moderpay-pos"
		case constants.APPVERSIONVSTYPE_MANGOPAY_APP:
			ParamZhCn = "mangopay-app"
		case constants.APPVERSIONVSTYPE_MANGOPAY_POS:
			ParamZhCn = "mangopay-pos"
		default:
			ParamZhCn = fmt.Sprintf("参数错误[%v]", paramStr)
			return ParamZhCn, false
		}
	case AppSignedStatus:
		switch paramStr {
		case constants.SignedStatusPending: //申请中
			ParamZhCn = "申请中"
		case constants.SignedStatusDeny: //驳回
			ParamZhCn = "驳回"
		case constants.SignedStatusPassed: //通过
			ParamZhCn = "通过"
		case constants.SignedStatusInvalid: //过期
			ParamZhCn = "过期"
		default:
			ParamZhCn = fmt.Sprintf("参数错误[%v]", paramStr)
			return ParamZhCn, false
		}
	case ChannelType:
		switch paramStr {
		case constants.CHANNELTYPE_ORDINARY: //普通渠道
			ParamZhCn = "普通渠道"
		case constants.CHANNELTYPE_THIRDPARTY: //第三方渠道
			ParamZhCn = "第三方渠道"
		default:
			ParamZhCn = fmt.Sprintf("参数错误[%v]", paramStr)
			return ParamZhCn, false
		}
	default:
		ss_log.Error("typeStr参数错误")
		ParamZhCn = fmt.Sprintf("传入参数typeStr[%v]错误", typeStr)
		return ParamZhCn, false
	}
	return ParamZhCn, true
}
