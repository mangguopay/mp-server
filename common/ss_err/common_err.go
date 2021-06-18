package ss_err

import (
	"fmt"

	"a.a/cu/ss_lang"
	"a.a/cu/strext"
)

var (
	I18nInstance ss_lang.SsI18n
)

const (
	ERR_SUCCESS           = "0"
	ERR_ARGS              = "AA2"
	ERR_PARAM             = "AA3" //参数错误
	ERR_SYSTEM            = "AA5" //系统错误
	ERR_TIMEFORMAT        = "AA6" // 时间格式错误
	ERR_UNKNOW_ERR        = "AA7" // 未知错误
	ERR_PARAM_MISSING     = "AA8" //参数缺失
	ERR_PERMISSION_DENIED = "AA9" //权限不足，拒绝访问

	ERR_SYS_DB_INIT        = "AA100001"
	ERR_SYS_DB_OP          = "AA100002" //数据库操作失败
	ERR_SYS_NETWORK        = "AA100003"
	ERR_SYS_REMOTE_API_ERR = "AA100004" //调用api失败
	ERR_SYS_IO_ERR         = "AA100005"
	ERR_SYS_DB_GET         = "AA100006"
	ERR_SYS_DB_ADD         = "AA100007"
	ERR_SYS_DB_UPDATE      = "AA100008"
	ERR_SYS_DB_DELETE      = "AA100009"
	ERR_SYS_DB_SAVE        = "AA100010"
	ERR_SYS_SIGN           = "AA100012"
	ERR_SYS_NO_API_AUTH    = "AA100013"
	ERR_SYS_EMPTY_BODY     = "AA100016"
	ERR_SYS_BODY_NOT_JSON  = "AA100017"

	ERR_SYS_NO_ROUTE             = "AA100028"
	ERR_SYS_DECODE               = "AA100029"
	ERR_SYS_UNKNOWN              = "AA100030"
	ERR_SYS_SIGN_JWT             = "AA100031"
	ERR_SYS_Common_Data          = "AA100032"
	ERR_ACCOUNT_ALREADY_EXISTS   = "AA101001"
	ERR_UPLOAD                   = "AA101002"
	ERR_ACCOUNT_NOT_EXISTS       = "AA101003"
	ERR_ACCOUNT_WRONG_PASSWORD   = "AA101004"
	ERR_ACCOUNT_SMS_CODE         = "AA101005"
	ERR_ACCOUNT_JWT_OUTDATED     = "AA101008"
	ERR_ACCOUNT_NO_PERMISSION    = "AA101013"
	ERR_ACCOUNT_NOT_PASSWORD     = "AA101015"
	ERR_ACCOUNT_INIT_ACCOUNT_ERR = "AA101017"
	ERR_ACCOUNT_ROLE_NOT_EXISTS  = "AA101018"
	ERR_ACCOUNT_MENU_NOT_EXISTS  = "AA101019"
	ERR_ACCOUNT_LOGIN_CODE       = "AA101021"
	ERR_MSG_ACCOUNT_NOT_EXISTS   = "AA101025" //账号不存在
	ERR_WALLET_PAY_PWD_ERR       = "AA101028"
	ERR_WALLET_AMOUNT_NULL       = "AA101029" //金额不能为空

	ERR_ACCOUNT_NF_IDEN_NO = "AA101040"

	ERR_ACCOUNT_LOGINED = "AA101043"
	ERR_ACCOUNT_STATUS  = "AA101051"

	ERR_ACCOUNT_SMS_MSG_FAILD                     = "AA101054" //"短信验证码错误"
	ERR_ACCOUNT_OLD_PWD_FAILD                     = "AA101055" //"原密码不正确"
	ERR_MODIFY_ACCOUNT_PWD_FAILD                  = "AA101056" //"修改密码失败"
	ERR_MODIFY_PAY_PWD_FAILD                      = "AA101057" //"修改支付密码失败"
	ERR_MODIFY_PHONE_FAILD                        = "AA101058" //"修改手机号码失败"
	ERR_MODIFY_NICKNAME_FAILD                     = "AA101059" //"修改昵称失败"
	ERR_MODIFY_GEN_KEY_FAILD                      = "AA101060" //"修改二维码状态失败"
	ERR_QUERY_SWEEP_CODE_STATUS_FAILD             = "AA101061" //"获取扫一扫码的状态失败"
	ERR_SAVE_IMAGE_FAILD                          = "AA101063" //"保存图片失败"
	ERR_REC_CARD_NUM_FAILD                        = "AA101064" // "收款人卡号不存在或不正确"
	ERR_TRANSFER_TO_HEAD_QUARTERS_FAILD           = "AA101065" // "转账到总部失败"
	ERR_LOG_TO_SERVICE_FAILD                      = "AA101066" // "向总部请款失败"
	ERR_SAVE_CARD_FAILD                           = "AA101067" // "绑定银行卡失败"
	ERR_CARD_IS_EXIST                             = "AA101068" // "银行卡已存在"
	ERR_MODIFY_DEFAULT_CARD_FAILD                 = "AA101069" // "修改默认卡失败"
	ERR_CARD_NOT_EXIST                            = "AA101070" // "银行卡不存在"
	ERR_CURRENT_CARD_IS_DEFAULT                   = "AA101071" // "当前卡是默认卡"
	ERR_DELETE_CARD_DEFAULT                       = "AA101072" // "解除银行卡失败"
	ERR_RATE_FAILD                                = "AA101074" // "计算费率结果失败"
	ERR_NOT_ROLE                                  = "AA101076" // 没有转入或转出的操作权限
	ERR_OPERATE_FAILD                             = "AA101078" // "操作失败"
	ERR_ORDER_STATUS_NO_INIT                      = "AA101079" // "订单不是待审核状态"
	ERR_PAY_VACC_ACCOUNT_NO_EXIST                 = "AA101080" // "付款人虚拟账户不存在"
	ERR_COLLECTION_VACC_ACCOUNT_NO_EXIST          = "AA101081" // "收款人虚拟账户不存在"
	ERR_MODIFY_HEAD_PORTRAIT_FAILD                = "AA101082" // "修改头像失败"
	ERR_ACCOUNT_IMAGE_BIG                         = "AA101083" // 图片太大了
	ERR_ACCOUNT_TRANSFER_TO_SELF                  = "AA101084" // 自己不能给自己转账
	ERR_PAY_CODE_STATUS                           = "AA101085" // 当前码的状态可能已被扫
	ERR_IMAGE_OP_FAILD                            = "AA101086" // 图片不存在
	ERR_ACC_NO_SER_FAILD                          = "AA101088" // 账号无对应服务商
	ERR_MODIFY_NICKNAME_ISNULL_FAILD              = "AA101089" //"修改的昵称不能为空"
	ERR_CHECK_PAY_PWD_FAILD                       = "AA101090" //"验证支付密码失败"
	ERR_TERMINAL_SN_FAILD                         = "AA101091" //"终端编号或posSn码为空"
	ERR_TERMINALNUM_UNIQUE_FAILD                  = "AA101092" //"终端编号已被添加过"
	ERR_TERMINALPOSSN_UNIQUE_FAILD                = "AA101093" //"终端posSn码已被添加过"
	ERR_MODIFY_TERMINAL_STATUS_FAILD              = "AA101094" //"无权限修改他人的pos状态"
	ERR_SAVE_ACCOUNT_LENGTH_FAILD                 = "AA101095" //"创建账号非法,长度小于6位"
	ERR_REC_ACCOUNT_IS_EXISTS_NO_PHONE_WITHDRAWAL = "AA101096" //"账号已存在,请通过扫码途径取款"
	ERR_AGREEMENT_BEING_USE                       = "AA101098" //"协议使用中，不可删除或修改为不使用"
	ERR_PROFITC_ASHABLE_FAILD                     = "AA101099" //"平台盈利可提现余额不足此次提现"

	ERR_WALLET_AMOUNT_FAILD                  = "AA101100" //"请输入取款金额"
	ERR_ACCOUNT_NO_LOGIN                     = "AA101101" //未登录或认证已过期
	ERR_ORDER_IS_NO_PENDING                  = "AA101102" // 订单不是待确认状态,二维码可能过期
	ERR_ACCOUNT_X_SIGN_FAILD                 = "AA101103" //x-sign认证失败
	ERR_POS_OUT_OF_RANGE                     = "AA101104" // POS机超出范围
	ERR_FILE_OP_FAILD                        = "AA101105" // 文件不存在
	ERR_PosChannel_FAILD                     = "AA101106" // pos渠道存在相同币种的渠道
	ERR_ACCOUNT_IS_RELA                      = "AA101107" // 该账号已被其他服务商添加了
	ERR_ACCOUNT_IS_SERVICER                  = "AA101108" // 不允许添加服务商为店员
	ERR_UseChannel_FAILD                     = "AA101109" // 银行卡存取款渠道存在相同币种的渠道
	ERR_REFRESH_TOKEN_FREQUENTLY             = "AA101110" // 刷新token太过频繁
	ERR_ACCOUNT_NOT_REAL_AUTH                = "AA101111" // 账号未实名认证
	ERR_ACCOUNT_REAL_NAME_NOT_SAME           = "AA101112" // 实名制姓名和持卡人姓名不一致
	ERR_ACCOUNT_ERR_PWD_LIMIT                = "AA101113" // 密码错误次数太多,请第二天再登录
	ERR_Card_Number_FAILD                    = "AA101114" // 卡号不正确
	ERR_Account_UseStatus_Frozen_FAILD       = "AA101115" // 账号已被冻结,无访问权限
	ERR_Payment_Pwd_Count_Limit              = "AA101116" // 支付密码错误超出今日限制,请第二天再试
	ERR_SYS_NO_ACCNO                         = "AA101117" // 获取账号失败
	ERR_SYS_API_SIGN_INFO                    = "AA101118" // 获取签名信息失败
	ERR_ACCOUNT_MailCode_FAILD               = "AA101119" //"验证码不正确，请重试"
	ERR_ACCOUNT_Mail_FAILD                   = "AA101120" //"邮箱已存在"
	ERR_Modify_ACCOUNT_Mail_FAILD            = "AA101121" //"修改邮箱失败"
	ERR_CardNumberLength_FAILD               = "AA101123" //"银行卡长度不符合规则"
	ERR_WRITE_OFF_CODE_Expired               = "AA101124" //"核销码已过期"
	ERR_TERMINALNUM_Account_NoRelation_FAILD = "AA101125" //"账号未绑定该pos机,无权登录"
	ERR_Audit_FAILD                          = "AA101126" //"审核失败"
	ERR_LOG_TO_BUSINESS_FAILD                = "AA101128" // "商家提现失败"
	ERR_BUSINESSApp_Status_FAILD             = "AA101129" // "应用不是审核不通过状态,不允许编辑"

	ERR_HasPayPwd_FAILD = "AA101131" //支付密码已初始化过，不允许再次初始化

	ERR_BulletinUseStatus_FAILD                = "AA101132" //无法编辑不是未发布状态的公告
	ERR_BusinessPublicKeyLength_FAILD          = "AA101133" //商家公钥不合法，长度小于350
	ERR_BusinessBatchTransferOrderStatus_FAILD = "AA101134" //批量转账已支付，无法重复支付
	ERR_BusinessBatchTransNumber_FAILD         = "AA101135" //批量转账超出限额笔数
	ERR_BusinessBatchTransAmount_FAILD         = "AA101136" //批量转账付款金额为0,不允许支付
	ERR_CountryCode_FAILD                      = "AA101137" //未知国家码
	ERR_PhoneISNull_FAILD                      = "AA101138" //手机号不能为空
	ERR_AddCard_REAL_NAME_NOT_SAME_FAILD       = "AA101139" // 银行卡开户名与认证信息不一致，请重新输入
	ERR_AddCard_No_BusinessRealName_FAILD      = "AA101140" // 请先进行实名认证，再添加银行卡
	ERR_HaveNotPass_BusinessRealName_FAILD     = "AA101141" // 该账号尚未通过实名认证，请前往实名认证处进行认证操作
	ERR_Account_Registered_FAILD               = "AA101142" // 该账号已注册
	ERR_BusinessPayPwd_FAILD                   = "AA101143" // 支付密码输入错误，请重试
	ERR_Business_OLD_PWD_FAILD                 = "AA101144" //当前密码输入错误，请重试
	ERR_Business_OLD_PayPWD_FAILD              = "AA101145" //原支付密码输入错误，请重试
	ERR_Business_Verification_Code_FAILD       = "AA101146" //"验证码错误，请重试"
	ERR_Business_App_Scene_Unique_FAILD        = "AA101147" //"应用已添加过产品签约,不允许重复签约"
	ERR_SignMethod_FAILD                       = "AA101148" //未知签名方式
	ERR_BusinessAuthName_Unique_FAILD          = "AA101149" //该公司名称已认证
	ERR_BusinessAuthNumber_Unique_FAILD        = "AA101150" //该注册号/机构组织代码已认证
	ERR_Servicer_Have_Cashier_FAILD            = "AA101151" //不允许重复添加店员
	ERR_BusinessSimplifyName_Unique_FAILD      = "AA101152" //该商家简称已认证
	ERR_TERMINALNUM_IN_USE                     = "AA101153" //终端编号被使用中
	ERR_TERMINALPOSSN_IN_USE                   = "AA101154" //POS码被使用中
	ERR_BusinessPhone_Unique_FAILD             = "AA101155" //该手机号已注册
	ERR_BusinessIndustryRateCycle_Unique_FAILD = "AA101156" //该行业的渠道已设置过费率和结算周期
	ERR_NoBusinessIndustryRateCycle_FAILD      = "AA101157" //未添加该行业
	ERR_AuthName_FAILD                         = "AA101158" //认证名称错误
	ERR_UnFilledAuthName_FAILD                 = "AA101159" //未填写认证名称
	ERR_UnFilledCurrencyType_FAILD             = "AA101160" //未填写币种
	ERR_Business_Scene_Unique_FAILD            = "AA101161" //不允许重复申请产品签约
	ERR_AppFingerprintOn_False                 = "AA101162" //指纹录入功能暂时已关闭
	ERR_PhoneUnRegisteredUser                  = "AA101163" //手机号未注册
	ERR_AppFingerprint_FAILD                   = "AA101164" //指纹支付失败,请前往设置处重新开启指纹支付

	ERR_MERC_IS_UPTOP = "AA200041"
	ERR_MERC_IS_DOWN  = "AA200042"

	ERR_PAY_UNKOWN_PRODUCT              = "AA300001" //产品不存在
	ERR_PAY_QUERY_ORDER_ERROR           = "AA300012"
	ERR_PAY_NO_THIS_ORDER               = "AA300013"
	ERR_PAY_REFUND_AMOUNT               = "AA300025"
	ERR_PAY_AMT_TOO_LOW                 = "AA300038"
	ERR_PAY_AMT_NOT_ENOUGH              = "AA300040"
	ERR_PAY_NO_QRCODE                   = "AA300042"
	ERR_PAY_TIMEOUT                     = "AA300056"
	ERR_PAY_ORDER_STATUS_MISTAKE        = "AA300058"
	ERR_PAY_PAY_TYPE_NOT_SUPPORT        = "AA300059"
	ERR_PAY_MISSING_EXCHANGE_RATE       = "AA300062"
	ERR_PAY_CANNOT_PAY_CODE             = "AA300063"
	ERR_PAY_VACCOUNT_OP_MISSING         = "AA300064"
	ERR_PAY_SAVE_MONEY                  = "AA300065"
	ERR_DB_PWD                          = "AA300067" // 密码错误
	ERR_WRONG_AMOUNT                    = "AA300068" // 取款金额不对
	ERR_PAY_OUT_MONEY                   = "AA300069" // 取款失败
	ERR_WRITE_OFF_CODE_FAILD            = "AA300070" // 核销码不正确
	ERR_QR_CODE_EXPIRED                 = "AA300072" // 二维码已过期
	ERR_DB_OP_SER                       = "AA300074" // 查询收银员的服务商失败
	ERR_WITHDRAW_AMT_NOT_ENOUGH         = "AA300075" // 提现余额不足
	ERR_MONEY_TYPE_FAILD                = "AA300077" // 币种不符合
	ERR_ACCOUNT_NOT_EXIST               = "AA300079" // 账号不存在
	ERR_PAY_QUOTA_NOT_ENOUGH            = "AA300081" // 额度不足
	ERR_PAY_AMOUNT_NOT_MIN_ZERO         = "AA300082" // 金额不能小于0
	ERR_PAY_NO_IN_COME_PERMISSION       = "AA300083" // 没有转入权限
	ERR_PAY_NO_OUT_GO_PERMISSION        = "AA300084" // 没有转出权限
	ERR_PAY_AMOUNT_IS_NO_INTEGER        = "AA300085" // usd 取款金额应为整数
	ERR_PAY_PWD_IS_NULL                 = "AA300086" // 请设置支付密码
	ERR_CreateFileDataNull              = "AA300087" // 查询的数据为空,已阻止生成xlsx文件
	ERR_PAY_AMOUNT_IS_LIMIT             = "AA300088" // 操作金额超出金额限制
	ERR_PAY_AMOUNT_FAILD                = "AA300089" // 操作金额不正确
	ERR_PAY_Channel_No_Support_Save     = "AA300090" // 该渠道不支持存款
	ERR_PAY_Channel_No_Support_Withdraw = "AA300091" // 该渠道不支持取款
	ERR_PAY_EXCHANGE_MIN_AMOUNT         = "AA300092" // 最低兑换金额为0.01USD
	ERR_PAY_QUERY_FEE_FAILD             = "AA300093" // 获取计算费率的数据信息失败
	ERR_PAY_DUP_ORDER                   = "AA300094" //订单号重复
	ERR_PAY_UPDATE_ORDER                = "AA300095" //更新订单失败
	ERR_MERC_NO_USE                     = "AA300098" // 商户被禁用状态
	ERR_Bank_Card_Not_Supported         = "AA300099" //银行卡不支持
	ERR_QueryOrderPaymentChannel        = "AA300100" //查询订单支付渠道失败
	ERR_OrderSettleFail                 = "AA300101" //订单结算失败

	ERR_ORDER_SUBMITTED_FREQUENTLY       = "AA300112" // 不能频繁提交订单
	ERR_PAY_FAILED_COUNT                 = "AA300113" // 密码错误,还可输入 %s 次
	ERR_Menu_Have_Child_ERR              = "AA300114" //菜单拥有子菜单，不允许删除
	ERR_PAY_ORDERSTATUS_NOT_PAY          = "AA300117" //订单无法退款
	ERR_WRITE_OFF_CODE_Cancelled         = "AA300119" //核销码已注销
	ERR_ACCOUNT_Actived_NOT_PERFORMED_OP = "AA300120" //账号已激活不能进行此操作
	ERR_WRITE_OFF_NOT_Operate            = "AA300121" //无法对核销码进行当前操作

	ERR_LOCAL_RULE_EXCEED_AMOUNT = "AA700024" //超出单笔交易额度

	ERR_VERSION_IsForce_Faile  = "AA1000005" //最新版本不可设置为强制更新
	ERR_CUST_NOT_EXISTS        = "AA1200011" //查询的用户不存在
	ERR_RISK_IS_RISK           = "AA1400001" // 被风控了
	ERR_PUSH_ACCOUNT_IS_NIL    = "AA1500001" // accToken 为空
	ERR_PUSH_SERVER_KEY_IS_NIL = "AA1500002" // server_key 为空
	ERR_PUSH_PHONE_IS_NIL      = "AA1500003" // 推送消息目标手机号为空
	ERR_PUSH_FAIL              = "AA1500004" // 推送消息失败

	ERR_PayAccountNo_NotTradeForbid = "AA1800001"
	ERR_AccountType                 = "AA1800002"
	ERR_PayOrderNoNotExist          = "AA1800003"
	ERR_PayOrderStatusErr           = "AA1800004"
	ERR_PayQrCodeNotExist           = "AA1800005"
	ERR_PayQrCodeExpire             = "AA1800006"
	ERR_PayOrderRefund              = "AA1800007"
	ERR_PayOrderStatusUnknown       = "AA1800008"
	ERR_PayQrCodeNotAvailable       = "AA1800009"
	ERR_PayPaymentCodeExpire        = "AA1800011"
	ERR_PayPaymentOrderFail         = "AA1800012" // 下单失败
	ERR_AppUnusable                 = "AA1800013"
	ERR_AppNotExist                 = "AA1800014"
	ERR_BusinessNotExist            = "AA1800016" // 商家不存在
	ERR_BusinessNotAvailable        = "AA1800017" // 商家不可用
	ERR_RefundNoNotExist            = "AA1800018" // 退款单号不存在
	ERR_SceneDisabled               = "AA1800019" // 该产品已禁用
	ERR_OrderCanNotManualSettle     = "AA1800020" // 订单不能手动结算
	ERR_OrderUnpaid                 = "AA1800021" // 订单未支付
	ERR_NOT_PERSONAL_BUSINESS       = "AA1800022" // 不是个人商户

	ERR_PayeeNotExist          = "AA1900001" //该账号不存在(收款方不存在)
	ERR_VirtualAccountNotExist = "AA1900002" //账户不存在
)

// 支付系统的错误code转换
func PayRetCode(code string) string {
	switch code {
	case Success:
		return ERR_SUCCESS
	case VerifySignFail:
		return ERR_SYS_SIGN
	case SystemErr:
		return ERR_SYSTEM
	case ParamErr:
		return ERR_PARAM
	case AppNotExist:
		return ERR_AppNotExist
	case AppNotPutOn:
		return ERR_AppUnusable
	case ProductUnsigned:
		return ERR_AppUnusable
	case AccountNoNotExist: //账号不存在
		return ERR_ACCOUNT_NOT_EXIST
	case OrderAlreadyExist: //订单已存在
		return ERR_PAY_DUP_ORDER
	case OrderNotExist: //订单不存在
		return ERR_PayOrderNoNotExist
	case PaymentCodeExpire: //付款码已过期
		return ERR_PayPaymentCodeExpire
	case PlaceAnOrderFail: //下单失败
		return ERR_PayPaymentOrderFail
	case OrderNoOrQrCodeIsEmpty:
		return ERR_PARAM_MISSING
	case AppIdIsEmpty:
		return ERR_PARAM_MISSING
	case TradeTypeIsEmpty:
		return ERR_PARAM_MISSING
	case OutOrderNoIsEmpty:
		return ERR_PARAM_MISSING
	case OrderNoIsEmpty:
		return ERR_PARAM_MISSING
	case NotifyUrlIsEmpty:
		return ERR_PARAM_MISSING
	case AmountIsEmpty:
		return ERR_PARAM_MISSING
	case CurrencyTypeIsEmpty:
		return ERR_PARAM_MISSING
	case QrCodeIdIsEmpty:
		return ERR_PARAM_MISSING
	case PaymentPwdIsEmpty:
		return ERR_PARAM_MISSING
	case NonStrIsEmpty:
		return ERR_PARAM_MISSING
	case AppPayContentIsEmpty:
		return ERR_PARAM_MISSING
	case TradeTypeValueIsIllegality: //交易类型错误
		return ERR_PARAM
	case CurrencyTypeValueIsIllegality: //币种错误
		return ERR_PARAM
	case AmountNotLessThanOne: //金额不能小于1
		return ERR_PAY_AMT_TOO_LOW
	case TimeFormatErr: //时间格式错误
		return ERR_TIMEFORMAT
	case AccountTypeNotUser: //账户不是用户
		return ERR_AccountType
	case QrCodeNotInvalid: //无效二维码
		return ERR_PayQrCodeExpire
	case QrCodeNotExist: //二维码不存在
		return ERR_PayQrCodeNotExist
	case QrCodeNotAvailable: //二维码不可用
		return ERR_PayQrCodeNotAvailable
	case BalanceNotEnough: //余额不足
		return ERR_PAY_AMT_NOT_ENOUGH
	case PaymentPwdError: //支付密码错误
		return ERR_WALLET_PAY_PWD_ERR
	case OrderNotRefundable: //订单不可退款
		return ERR_PAY_ORDERSTATUS_NOT_PAY
	case RefundAmountDisagree: //退款金额与交易金额不一致
		return ERR_PAY_REFUND_AMOUNT
	case PayeeAccountNoNotExist: //收款方账号不存在
		return ERR_PayeeNotExist
	case OrderPaid:
		return ERR_PayOrderStatusErr //订单已支付
	case OrderUnpaid:
		return ERR_OrderUnpaid
	case OrderExpired: //订单已过期
		return ERR_ORDER_IS_NO_PENDING
	case AccountNoNotTradeForbid: //账号没有交易权限
		return ERR_PayAccountNo_NotTradeForbid
	case VirtualAccountNotExist: //虚账不存在
		return ERR_VirtualAccountNotExist
	case RefundAmountExcessBalance:
		return ERR_PAY_REFUND_AMOUNT
	case BusinessNotExist:
		return ERR_BusinessNotExist
	case BusinessNotAvailable:
		return ERR_BusinessNotAvailable
	case OrderFullRefund:
		return ERR_PayOrderRefund
	case OrderStatusUnknown:
		return ERR_PayOrderStatusUnknown
	case QrCodeExpired:
		return ERR_PayQrCodeExpire
	case RefundNoNotExist:
		return ERR_RefundNoNotExist
	case SceneDisabled:
		return ERR_SceneDisabled
	case BankCardNotSupported:
		return ERR_Bank_Card_Not_Supported
	case BankCardNumberIsEmpty:
		return ERR_PARAM_MISSING
	case OrderNotSupportedManualSettle:
		return ERR_OrderCanNotManualSettle
	case QueryOrderPaymentChannel:
		return ERR_QueryOrderPaymentChannel
	case OrderSettleFail:
		return ERR_OrderSettleFail
	case UserHasNoRealName:
		return ERR_ACCOUNT_NOT_REAL_AUTH

	}
	return ERR_SYSTEM
}

func GetErrMsgMulti(langType ss_lang.SsLang, code string, extstr ...interface{}) string {
	extStr2 := []interface{}{}
	for _, v := range extstr {
		extStr2 = append(extStr2, I18nInstance.GetMsg(langType, strext.ToStringNoPoint(v)))
	}
	return fmt.Sprintf(I18nInstance.GetMsg(langType, code), extStr2...)
}

func GetMsgAddArgs(langType string, code string, args ...interface{}) string {
	return fmt.Sprintf(I18nInstance.GetMsg(ss_lang.SsLang(langType), code), args...)
}
