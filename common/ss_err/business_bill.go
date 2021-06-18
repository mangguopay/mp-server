package ss_err

import (
	"a.a/mp-server/common/constants"
)

const (
	Success = "SUCCESS" //成功

	VerifySignFail = "VERIFY_SIGN_FAIL" //验签失败
	SystemErr      = "SYSTEM_ERROR"     //系统错误
	ParamErr       = "PARAM_ERROR"      //参数错误
	SceneDisabled  = "SCENE_DISABLED"   //产品已被禁用

	BusinessNotExist     = "MERCHANT_NOT_EXIST"     //商家不存在
	BusinessNotAvailable = "MERCHANT_NOT_AVAILABLE" //商家不可用
	AppNotExist          = "APP_NOT_EXIST"          //APP不存在
	AppNotPutOn          = "APP_NOT_PUT_ON"         //APP未上架
	ProductUnsigned      = "PRODUCT_NOT_SIGNED"     //产品未签约
	SignedExpired        = "PRODUCT_SIGNED_EXPIRED" //产品签约已过期
	UserHasNoRealName    = "USER_HAS_NO_REAL_NAME"  //未实名

	QrCodeNotExist     = "QR_CODE_NOT_EXIST"     //二维码不存在
	QrCodeNotAvailable = "QR_CODE_NOT_AVAILABLE" //二维码不可用
	QrCodeExpired      = "QR_CODE_EXPIRED"       //二维码已过期
	QrCodeNotInvalid   = "QR_CODE_INVALID"       //二维码无效
	PaymentCodeExpire  = "PAYMENT_CODE_EXPIRED"  //付款码已过期

	PlaceAnOrderFail   = "PLACE_AN_ORDER_FAIL"  //下单失败
	OrderAlreadyExist  = "ORDER_ALREADY_EXIST"  //订单已存在
	OrderNotExist      = "ORDER_NOT_EXIST"      //订单不存在
	OrderPaid          = "ORDER_PAID"           //订单已支付
	OrderUnpaid        = "ORDER_UNPAID"         //订单未支付
	OrderExpired       = "ORDER_EXPIRED"        //订单已过期
	OrderStatusUnknown = "ORDER_STATUS_UNKNOWN" //订单状态未知
	OrderFullRefund    = "ORDER_FULL_REFUND"    //订单已全额退款

	RefundNoNotExist          = "REFUND_ORDER_NOT_EXIST"       //退款订单不存在
	RefundAmountExcessBalance = "REFUND_AMOUNT_EXCESS_BALANCE" //退款金额超出可退金额
	BalanceNotEnough          = "BALANCE_NOT_ENOUGH"           //余额不足
	AmountNotLessThanOne      = "AMOUNT_NOT_LESS_THAN_1"       //发起金额不能小于1
	PaymentPwdError           = "PAYMENT_PWD_ERROR"            //支付密码错误
	OrderNotRefundable        = "ORDER_NOT_REFUNDABLE"         //订单不能进行退款操作
	RefundAmountDisagree      = "REFUND_AMOUNT_DISAGREE"       //退款金额与交易金额不一致
	TransferNoNotExist        = "TRANSFER_ORDER_NOT_EXIST"     //企业付款订单不存在

	AccountNoNotExist       = "ACCOUNT_NO_NOT_EXIST"        //账号不存在
	AccountTypeNotUser      = "ACCOUNT_TYPE_NOT_USER"       //账号不是用户账号
	PayeeAccountNoNotExist  = "PAYEE_ACCOUNT_NOT_EXIST"     //收款方账号不存在
	PayeeAccountErr         = "PAYEE_ACCOUNT_ERROR"         //收款方账号错误
	PayeeNotTradeForbid     = "PAYEE_NOT_TRADE_FORBID"      //收款方没有交易权限
	VirtualAccountNotExist  = "ACCOUNT_NOT_TRADE_FORBID"    //账户不存在
	AccountNoNotTradeForbid = "ACCOUNT_No_NOT_TRADE_FORBID" //账号没有交易权限

	OrderNoOrQrCodeIsEmpty = "ORDER_NO_OR_QR_CODE_PARAMETER_MISSING" //OrderNo或QrCode参数缺失
	AppIdIsEmpty           = "APP_ID_PARAMETER_MISSING"              //AppId参数缺失
	TradeTypeIsEmpty       = "TRADE_TYPE_PARAMETER_MISSING"          //TradeType参数缺失
	OutOrderNoIsEmpty      = "OUT_ORDER_NO_PARAMETER_MISSING"        //OutOrderNo参数缺失
	OrderNoIsEmpty         = "ORDER_NO_PARAMETER_MISSING"            //OrderNo参数缺失
	NotifyUrlIsEmpty       = "NOTIFY_URL_PARAMETER_MISSING"          //NotifyUrl参数缺失
	AmountIsEmpty          = "AMOUNT_PARAMETER_MISSING"              //Amount参数缺失
	CurrencyTypeIsEmpty    = "CURRENCY_TYPE_PARAMETER_MISSING"       //CurrencyType参数缺失
	QrCodeIdIsEmpty        = "QR_CODE_PARAMETER_MISSING"             //QrCode参数缺失
	AccountNoIsEmpty       = "ACCOUNT_NO_PARAMETER_MISSING"          //AccountNo参数缺失
	PaymentPwdIsEmpty      = "PAYMENT_PWD_PARAMETER_MISSING"         //PaymentPwd参数缺失
	NonStrIsEmpty          = "NONSTR_PARAMETER_MISSING"              //NonStr参数缺失
	AppPayContentIsEmpty   = "APP_PAY_CONTENT_PARAMETER_MISSING"     //AppPayContent参数缺失
	FixedQrCodeIsEmpty     = "FIXED_QR_CODE_PARAMETER_MISSING"       //FixedQrCode参数缺失
	PaymentCodeIsEmpty     = "PAYMENT_CODE_PARAMETER_MISSING"        //PaymentCode参数缺失
	BankCardNumberIsEmpty  = "BANK_CARD_NUMBER_PARAMETER_MISSING"    //BankCardNumber参数缺失
	PayeePhoneIsEmpty      = "PAYEE_PHONE_PARAMETER_MISSING"         //PayeePhone参数缺失
	CountryCodeIsEmpty     = "COUNTRY_CODE_PARAMETER_MISSING"        //CountryCode参数缺失
	RefundAmountIsEmpty    = "REFUND_AMOUNT_PARAMETER_MISSING"       //RefundAmount参数缺失

	TimeFormatErr                 = "TIME_FORMAT_ERROR"                                   //日期时间参数格式错误
	TradeTypeValueIsIllegality    = "TRADE_TYPE_VALUE_IS_ILLEGALITY"                      //交易类型错误
	CurrencyTypeValueIsIllegality = "CURRENCY_TYPE_VALUE_IS_ILLEGALITY"                   //币种错误
	OrderSettleFail               = "ORDER_SETTLE_FAIL"                                   //订单结算失败
	QueryOrderPaymentChannel      = "QUERY_ORDER_PAYMENT_CHANNEL_FAIL"                    //查询订单支付渠道失败
	OrderNotSupportedManualSettle = "MANUAL_SETTLEMENT_OF_ORDERS_IS_NOT_SUPPORTED"        //订单暂不支持手动结算
	BankCardNotSupported          = "BANK_CARD_NOT_SUPPORTED"                             //银行卡不支持
	TransactionAmountLimit        = "TRANSACTION_AMOUNT_EXCEEDS_MAXIMUM_OR_MINIMUM_LIMIT" //交易金额超出最大或最小限制

)

var en_US_Map = make(map[string]string)
var km_KH_Map = make(map[string]string)
var zh_CN_Map = make(map[string]string)

func init() {
	init_enUS()
	init_kmKH()
	init_zhCN()
}

func init_enUS() {
	en_US_Map[Success] = "success"
	en_US_Map[VerifySignFail] = "Sign verification failed"
	en_US_Map[SystemErr] = "System error"
	en_US_Map[ParamErr] = "Parameter error"
	en_US_Map[BusinessNotExist] = "Merchants don't exist."
	en_US_Map[BusinessNotAvailable] = "Unavailable to merchants"
	en_US_Map[SceneDisabled] = "Product has been disabled"

	en_US_Map[AppNotExist] = "The application does not exist."
	en_US_Map[AppNotPutOn] = "Application not yet available"
	en_US_Map[ProductUnsigned] = "Product not signed"
	en_US_Map[SignedExpired] = "Product contract has expired"
	en_US_Map[UserHasNoRealName] = "No real-name authentication"

	en_US_Map[QrCodeNotExist] = "QR codes don't exist."
	en_US_Map[QrCodeNotAvailable] = "The QR code is not available."
	en_US_Map[QrCodeExpired] = "The QR code has expired."
	en_US_Map[QrCodeNotInvalid] = "Invalid QR code"
	en_US_Map[PaymentCodeExpire] = "Payment code has expired"

	en_US_Map[PlaceAnOrderFail] = "Order failed"
	en_US_Map[OrderAlreadyExist] = "Order already exists"
	en_US_Map[OrderNotExist] = "The order does not exist."
	en_US_Map[OrderPaid] = "Order paid"
	en_US_Map[OrderUnpaid] = ""
	en_US_Map[OrderExpired] = "Order has expired"
	en_US_Map[OrderStatusUnknown] = "Order status unknown"
	en_US_Map[OrderFullRefund] = "The order has been fully refunded"

	en_US_Map[RefundNoNotExist] = "Refund order does not exist"
	en_US_Map[RefundAmountExcessBalance] = "The refund amount exceeds the refundable amount"
	en_US_Map[BalanceNotEnough] = "Insufficient balance"
	en_US_Map[AmountNotLessThanOne] = "The initiated amount cannot be less than 1"
	en_US_Map[PaymentPwdError] = "Wrong payment password"
	en_US_Map[OrderNotRefundable] = "Order cannot be refunded"
	en_US_Map[RefundAmountDisagree] = "The refund amount does not match the transaction amount"

	en_US_Map[AccountNoNotExist] = "User account does not exist"
	en_US_Map[AccountTypeNotUser] = "Account is not a user account"
	en_US_Map[PayeeAccountNoNotExist] = "Payee account does not exist"
	en_US_Map[PayeeAccountErr] = ""
	en_US_Map[PayeeNotTradeForbid] = ""
	en_US_Map[VirtualAccountNotExist] = "The account does not exist."
	en_US_Map[AccountNoNotTradeForbid] = "Account does not have transaction authority"

	en_US_Map[OrderNoOrQrCodeIsEmpty] = "order_no or qr_code parameter is missing"
	en_US_Map[AppIdIsEmpty] = "app_id parameter is missing"
	en_US_Map[TradeTypeIsEmpty] = "trade_type parameter is missing"
	en_US_Map[OutOrderNoIsEmpty] = "out_order_no parameter is missing"
	en_US_Map[OrderNoIsEmpty] = "order_no parameter is missing"
	en_US_Map[NotifyUrlIsEmpty] = "notify_url parameter is missing"
	en_US_Map[AmountIsEmpty] = "amount parameter is missing"
	en_US_Map[CurrencyTypeIsEmpty] = "currency_type parameter is missing"
	en_US_Map[QrCodeIdIsEmpty] = "qr_code parameter is missing"
	en_US_Map[AccountNoIsEmpty] = "account_no parameter is missing"
	en_US_Map[PaymentPwdIsEmpty] = "payment_pwd parameter is missing"
	en_US_Map[NonStrIsEmpty] = "non_str parameter is missing"
	en_US_Map[AppPayContentIsEmpty] = "app_pay_content parameter is missing"
	en_US_Map[FixedQrCodeIsEmpty] = "fixed_qr_code parameter is missing"
	en_US_Map[PaymentCodeIsEmpty] = "payment_code parameter is missing"
	en_US_Map[CountryCodeIsEmpty] = "country_code parameter is missing"
	en_US_Map[PayeePhoneIsEmpty] = "payee_phone parameter is missing"
	en_US_Map[BankCardNumberIsEmpty] = "bank_card_number parameter is missing"
	en_US_Map[RefundAmountIsEmpty] = "refund_amount parameter is missing"

	en_US_Map[TimeFormatErr] = "Time format error"
	en_US_Map[TradeTypeValueIsIllegality] = "Wrong transaction type"
	en_US_Map[CurrencyTypeValueIsIllegality] = "Currency error"
	en_US_Map[OrderSettleFail] = "Order settlement failed"
	en_US_Map[QueryOrderPaymentChannel] = "Failed to query order payment channel"
	en_US_Map[OrderNotSupportedManualSettle] = "Orders currently do not support manual settlement"
	en_US_Map[BankCardNotSupported] = "Bank card not supported"
	en_US_Map[TransactionAmountLimit] = "Transaction amount exceeds limit"
}

func init_kmKH() {
	km_KH_Map[Success] = "ជោគជ័យ"
	km_KH_Map[VerifySignFail] = "ការផ្ទៀងផ្ទាត់បានបរាជ័យ"
	km_KH_Map[SystemErr] = "កំហុសប្រព័ន្"
	km_KH_Map[ParamErr] = "ការបញ្ជូលេខខុស"
	km_KH_Map[BusinessNotExist] = "មិនមានអាជីវកម្មទេ"
	km_KH_Map[BusinessNotAvailable] = "អាជីវកម្មនេះមិនអាចប្រើបានទេ"
	km_KH_Map[SceneDisabled] = ""

	km_KH_Map[AppNotExist] = "មិនមានពាក្យសុំកម្មវិធីទេ"
	km_KH_Map[AppNotPutOn] = "កម្មវិធីមិនមានលើកឡើងទេ"
	km_KH_Map[ProductUnsigned] = "ផលិតផលមិនបានចុះហត្ថលេខា"
	km_KH_Map[SignedExpired] = "កិច្ចសន្យាផលិតផលបានផុតកំណត់ហើយ"
	km_KH_Map[UserHasNoRealName] = "គ្មានការផ្ទៀងផ្ទាត់ឈ្មោះពិតទេ"

	km_KH_Map[QrCodeNotExist] = "មិនមានQRកូដ"
	km_KH_Map[QrCodeNotAvailable] = "QRកូដប្រើមិន​បាន"
	km_KH_Map[QrCodeExpired] = "QRកូដបានផុតកំណត់"
	km_KH_Map[QrCodeNotInvalid] = "លេខកូដ QR មិនត្រឹមត្រូវ"
	km_KH_Map[PaymentCodeExpire] = "ការផ្ទេរប្រាក់បានផុតកំណត់ហើយ"

	km_KH_Map[PlaceAnOrderFail] = "ការបញ្ជាទិញបានបរាជ័យ"
	km_KH_Map[OrderAlreadyExist] = "ការបញ្ជាទិញមានរួចហើយ"
	km_KH_Map[OrderNotExist] = "ការបញ្ជាទិញមិនមានទេ"
	km_KH_Map[OrderPaid] = "ការបញ្ជាទិញបានបង់"
	km_KH_Map[OrderUnpaid] = ""
	km_KH_Map[OrderExpired] = "ការបញ្ជាទិញបានផុតកំណត់ហើយ"
	km_KH_Map[OrderStatusUnknown] = "មិនស្គាល់ស្ថានភាពបញ្ជាទិញ"
	km_KH_Map[OrderFullRefund] = "ការបញ្ជាទិញត្រូវបានសងប្រាក់វិញទាំងស្រុង"

	km_KH_Map[RefundNoNotExist] = "ការបញ្ជាទិញសងប្រាក់វិញមិនមានទេ"
	km_KH_Map[RefundAmountExcessBalance] = "ចំនួនទឹកប្រាក់សងវិញលើសចំនួនដែលអាចសងវិញ"
	km_KH_Map[BalanceNotEnough] = "តុល្យភាពមិនគ្រប់គ្រាន់"
	km_KH_Map[AmountNotLessThanOne] = "ចំនួនទឹកប្រាក់ចាប់ផ្តើមមិនតិចជាង ១ ទេ"
	km_KH_Map[PaymentPwdError] = "កំហុសក្នុងការបញ្ចូលពាក្យសម្ងាត់ទូទាត់"
	km_KH_Map[OrderNotRefundable] = "ការបញ្ជាទិញមិនអាចត្រូវបានសងប្រាក់វិញទេ"
	km_KH_Map[RefundAmountDisagree] = "ចំនួនទឹកប្រាក់សងប្រាក់វិញមិនត្រូវនឹងចំនួនប្រតិបត្តិការទេ"

	km_KH_Map[AccountNoNotExist] = "មិនមានគណនីអ្នកប្រើទេ"
	km_KH_Map[AccountTypeNotUser] = "គណនីមិនមែនជាគណនីអ្នកប្រើទេ"
	km_KH_Map[PayeeAccountNoNotExist] = "គណនីអ្នកបង់ប្រាក់មិនមានទេ"
	km_KH_Map[PayeeAccountErr] = ""
	km_KH_Map[PayeeNotTradeForbid] = ""
	km_KH_Map[VirtualAccountNotExist] = "មិនមានគណនីទេ"
	km_KH_Map[AccountNoNotTradeForbid] = "គណនីមិនមានសិទ្ធិប្រតិបត្តិការទេ"

	km_KH_Map[OrderNoOrQrCodeIsEmpty] = "ប៉ារ៉ាម៉ែត្រ __ ឬ qr_code គឺទទេ"
	km_KH_Map[AppIdIsEmpty] = "ប៉ារ៉ាម៉ែត្រ app_id គឺទទេ"
	km_KH_Map[TradeTypeIsEmpty] = "ប៉ារ៉ាម៉ែត្រនៃការធ្វើពាណិជ្ជកម្មគឺទទេ"
	km_KH_Map[OutOrderNoIsEmpty] = "out_order_nមិនអាចទទេបានទេ"
	km_KH_Map[OrderNoIsEmpty] = "order_noមិនអាចទទេបានទេ"
	km_KH_Map[NotifyUrlIsEmpty] = "notify_urlមិនអាចទទេបានទេ"
	km_KH_Map[AmountIsEmpty] = "amountមិនអាចទទេបានទេ"
	km_KH_Map[CurrencyTypeIsEmpty] = "currency_typeមិនអាចទទេបានទេ"
	km_KH_Map[QrCodeIdIsEmpty] = "ប៉ារ៉ាម៉ែត្រ qr_code គឺទទេ"
	km_KH_Map[AccountNoIsEmpty] = "account_noមិនអាចទទេបានទេ"
	km_KH_Map[PaymentPwdIsEmpty] = "ប៉ារ៉ាម៉ែត្រនៃការបង់ប្រាក់គឺទទេ"
	km_KH_Map[NonStrIsEmpty] = "ប៉ារ៉ាម៉ែត្រ non_str គឺទទេ"
	km_KH_Map[AppPayContentIsEmpty] = "app_pay_content ប៉ារ៉ាម៉ែត្រគឺទទេ"
	km_KH_Map[FixedQrCodeIsEmpty] = "ប៉ារ៉ាម៉ែត្រថេរ _qr_code គឺទទេ"
	km_KH_Map[PaymentCodeIsEmpty] = "ប៉ារ៉ាម៉ែត្រនៃការទូទាត់ - កូដគឺទទេ"
	km_KH_Map[CountryCodeIsEmpty] = ""
	km_KH_Map[PayeePhoneIsEmpty] = ""
	km_KH_Map[BankCardNumberIsEmpty] = ""
	km_KH_Map[RefundAmountIsEmpty] = ""

	km_KH_Map[TimeFormatErr] = "ទ្រង់ទ្រាយពេលវេលាខុស"
	km_KH_Map[TradeTypeValueIsIllegality] = "ប្រភេទប្រតិបត្តិការមិនត្រឹមត្រូវ"
	km_KH_Map[CurrencyTypeValueIsIllegality] = "កំហុសរូបិយប័ណ្ណ"
	km_KH_Map[OrderSettleFail] = "ការទូទាត់ការបញ្ជាទិញបានបរាជ័យ"
	km_KH_Map[QueryOrderPaymentChannel] = "បានបរាជ័យក្នុងការសាកសួរឆានែលទូទាត់លំដាប់"
	km_KH_Map[OrderNotSupportedManualSettle] = "ការបញ្ជាទិញបច្ចុប្បន្នមិនគាំទ្រការទូទាត់ដោយដៃទេ"
	km_KH_Map[BankCardNotSupported] = "កាតធនាគារមិនគាំទ្រទេ"
	km_KH_Map[TransactionAmountLimit] = "ចំនួនប្រតិបត្តិការលើសពីដែនកំណត់"
}

func init_zhCN() {
	zh_CN_Map[Success] = "成功"
	zh_CN_Map[VerifySignFail] = "验签失败"
	zh_CN_Map[SystemErr] = "系统错误"
	zh_CN_Map[ParamErr] = "参数错误"
	zh_CN_Map[BusinessNotExist] = "商家不存在"
	zh_CN_Map[BusinessNotAvailable] = "商家不可用"
	zh_CN_Map[SceneDisabled] = "产品已被禁用，暂不可交易"

	zh_CN_Map[AppNotExist] = "应用不存在"
	zh_CN_Map[AppNotPutOn] = "应用未上架"
	zh_CN_Map[ProductUnsigned] = "产品未签约"
	zh_CN_Map[SignedExpired] = "产品签约已过期"
	zh_CN_Map[UserHasNoRealName] = "未实名"

	zh_CN_Map[QrCodeNotExist] = "二维码不存在"
	zh_CN_Map[QrCodeNotAvailable] = "二维码不可用"
	zh_CN_Map[QrCodeExpired] = "二维码已过期"
	zh_CN_Map[QrCodeNotInvalid] = "二维码无效"
	zh_CN_Map[PaymentCodeExpire] = "付款码已过期"

	zh_CN_Map[PlaceAnOrderFail] = "下单失败"
	zh_CN_Map[OrderAlreadyExist] = "订单已存在"
	zh_CN_Map[OrderNotExist] = "订单不存在"
	zh_CN_Map[OrderPaid] = "订单已支付"
	zh_CN_Map[OrderUnpaid] = "订单未支付"
	zh_CN_Map[OrderExpired] = "订单已过期"
	zh_CN_Map[OrderStatusUnknown] = "订单状态未知"
	zh_CN_Map[OrderFullRefund] = "订单已全额退款"

	zh_CN_Map[RefundNoNotExist] = "退款订单不存在"
	zh_CN_Map[RefundAmountExcessBalance] = "退款金额超出可退金额"
	zh_CN_Map[BalanceNotEnough] = "余额不足"
	zh_CN_Map[AmountNotLessThanOne] = "发起金额不能小于1"
	zh_CN_Map[PaymentPwdError] = "支付密码错误"
	zh_CN_Map[OrderNotRefundable] = "订单不能进行退款操作"
	zh_CN_Map[RefundAmountDisagree] = "退款金额与交易金额不一致"

	zh_CN_Map[AccountNoNotExist] = "账号不存在"
	zh_CN_Map[AccountTypeNotUser] = "账号不是用户账号"
	zh_CN_Map[PayeeAccountNoNotExist] = "收款方账号不存在"
	zh_CN_Map[PayeeAccountErr] = "收款方账号错误"
	zh_CN_Map[PayeeNotTradeForbid] = "收款方没有交易权限"
	zh_CN_Map[VirtualAccountNotExist] = "账户不存在"
	zh_CN_Map[AccountNoNotTradeForbid] = "账号没有交易权限"

	zh_CN_Map[OrderNoOrQrCodeIsEmpty] = "order_no或qr_code参数缺失"
	zh_CN_Map[AppIdIsEmpty] = "app_id参数缺失"
	zh_CN_Map[TradeTypeIsEmpty] = "trade_type参数缺失"
	zh_CN_Map[OutOrderNoIsEmpty] = "out_order_no参数缺失"
	zh_CN_Map[OrderNoIsEmpty] = "order_no参数缺失"
	zh_CN_Map[NotifyUrlIsEmpty] = "notify_url参数缺失"
	zh_CN_Map[AmountIsEmpty] = "amount参数缺失"
	zh_CN_Map[CurrencyTypeIsEmpty] = "currency_type参数缺失"
	zh_CN_Map[QrCodeIdIsEmpty] = "qr_code参数缺失"
	zh_CN_Map[AccountNoIsEmpty] = "account_no参数缺失"
	zh_CN_Map[PaymentPwdIsEmpty] = "payment_pwd参数缺失"
	zh_CN_Map[NonStrIsEmpty] = "non_str参数缺失"
	zh_CN_Map[AppPayContentIsEmpty] = "app_pay_content参数缺失"
	zh_CN_Map[FixedQrCodeIsEmpty] = "fixed_qr_code参数缺失"
	zh_CN_Map[PaymentCodeIsEmpty] = "payment_code参数缺失"
	zh_CN_Map[CountryCodeIsEmpty] = "country_code参数缺失"
	zh_CN_Map[PayeePhoneIsEmpty] = "payee_phone参数缺失"
	zh_CN_Map[BankCardNumberIsEmpty] = "bank_card_number参数缺失"
	zh_CN_Map[RefundAmountIsEmpty] = "refund_amount参数缺失"

	zh_CN_Map[TimeFormatErr] = "时间格式错误"
	zh_CN_Map[TradeTypeValueIsIllegality] = "交易类型错误"
	zh_CN_Map[CurrencyTypeValueIsIllegality] = "币种错误"
	zh_CN_Map[OrderSettleFail] = "订单结算失败"
	zh_CN_Map[QueryOrderPaymentChannel] = "查询订单支付渠道失败"
	zh_CN_Map[OrderNotSupportedManualSettle] = "订单暂不支持手动结算"
	zh_CN_Map[BankCardNotSupported] = "银行卡不支持"
	zh_CN_Map[TransactionAmountLimit] = "交易金额超出最大或最小限制"
}

func GetMsg(code, lang string) string {
	msg := ""

	switch lang {
	case constants.LangZhCN:
		msg = zh_CN_Map[code]
	case constants.LangEnUS:
		msg = en_US_Map[code]
	case constants.LangKmKH:
		msg = km_KH_Map[code]
	default:
		msg = code
	}

	return msg
}

// auth-srv服务接口的错误code转换
func AuthSrvRetCode(retCode string) string {
	switch retCode {
	case ERR_SUCCESS:
		return Success
	case ERR_PAY_PWD_IS_NULL:
		return PaymentPwdIsEmpty
	case ERR_DB_PWD:
		return PaymentPwdError
	case ERR_PARAM:
		return ParamErr
	}

	return SystemErr
}

// bill-srv服务接口的错误code转换
func BillSrvRetCode(retCode string) string {
	switch retCode {
	case ERR_SUCCESS:
		return Success
	case ERR_PARAM:
		return ParamErr
	case ERR_SYSTEM:
		return SystemErr
	case ERR_PAY_AMOUNT_IS_LIMIT:
		return TransactionAmountLimit
	case ERR_WALLET_AMOUNT_NULL:
		return AmountIsEmpty
	case ERR_MERC_NO_USE:
		return BusinessNotAvailable
	case ERR_PAY_NO_OUT_GO_PERMISSION:
		return AccountNoNotTradeForbid
	case ERR_PAY_AMT_NOT_ENOUGH:
		return BalanceNotEnough
	case ERR_PayeeNotExist:
		return PayeeAccountNoNotExist
	case ERR_PAY_NO_IN_COME_PERMISSION:
		return PayeeNotTradeForbid
	case ERR_PAY_VACCOUNT_OP_MISSING:
		return VirtualAccountNotExist
	case ERR_WITHDRAW_AMT_NOT_ENOUGH:
		return BalanceNotEnough
	default:
		return ParamErr
	}

	return SystemErr
}
