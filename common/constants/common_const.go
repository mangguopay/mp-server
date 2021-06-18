package constants

const (
	CB_STATUS_INIT  = 0
	CB_STATUS_RECV  = 1
	CB_STATUS_DOING = 2
	CB_STATUS_DONE  = 3

	CbMethod_PostBody = "0"
	CbMethod_PostForm = "1"
)

// 通道
const (
	ChannelCode_Unknown = "0"
)

type ChannelCodeName string

type RouteType int32

type ApiType string

const (
	//管理后台
	AccountType_ADMIN    = "1" // 管理员
	AccountType_OPERATOR = "2" //// 运营

	AccountType_SERVICER     = "3" // 服务商
	AccountType_USER         = "4" // 用户
	AccountType_POS          = "5" // 收银员
	AccountType_Headquarters = "6" // 总部

	//商家中心
	AccountType_PersonalBusiness   = "7" //个人商户
	AccountType_EnterpriseBusiness = "8" ////企业商户

)

const (
	OpAccType_Servicer = 1
	OpAccType_Pos      = 2
)

const (
	OpAccType_Count_Fee_Save     = 1
	OpAccType_Count_Fee_Withdraw = 2
)

const (
	Channel_Support_Type_Save          = "1"
	Channel_Support_Type_Withdraw      = "2"
	Channel_Support_Type_Save_Withdraw = "3"
)

const (
	Fee_Charge_Type_Rate  = "1" // 按比例收费
	Fee_Charge_Type_Count = "2" // 按单笔收取手续费
)

const (
	NO_DEFAULT_CARD int32 = 0
	IS_DEFAULT_CARD int32 = 1
)

const (
	TimeType_Second = "second"
	TimeType_Minute = "minute"
	TimeType_Hour   = "hour"
	TimeType_Day    = "day"
	TimeType_Week   = "week"
	TimeType_Month  = "month"
	TimeType_Year   = "year"
)

const (
	CheckStatus_NoChecked = "0"
	CheckStatus_Checked   = "1"
	CheckStatus_Exp       = "2"
)

const (
	RedisNilMsg   = "redigo: nil returned"
	RedisNilValue = "redis: nil"
)

//默认输入支付密码连续错误次数
const (
	PaymentPwdErrLimit      = "err_payment_pwd_count"
	ErrPwdLimitDefaultCount = 5
)

const (
	WalletStatus_NoIn = 0
	WalletStatus_In   = 1
)

const (
	LOGIN_CLI_WEB       = "1"
	LOGIN_CLI_MICROPRO  = "2"
	LOGIN_CLI_DEV       = "3"
	LOGIN_CLI_WEB_SCORE = "4"
)

const (
	// 电话
	LoginAccType_Phone = "1"
	// 邮件
	LoginAccType_Email = "2"
	// 账号
	LoginAccType_Acc = "3"
	// 微信
	LoginAccType_Wechat = "4"
)

const (
	// 内扣=直接扣款项
	Charging_Inner = "1" // 内扣
	// 外扣=支付全款,另外计算手续费
	Charging_Out = "2" // 外扣

	PayTime_Normal = "0" // 进账
	PayTime_Pre    = "1" // 预付
	PayTime_After  = "2" // 后付
	PayTime_Remain = "3" // 余额付
)

const (
	UNAUTHORIZED = "UNAUTHORIZED" //用户未提交认证信息
	AUTHORIZING  = "AUTHORIZING"  //认证中，已提交至上游
	DENIED       = "DENIED"       //认证失败, 可重新发起
	AUTHORIZED   = "AUTHORIZED"   //认证通过
)

const (
	UNPUBLISHED = "UNPUBLISHED" // 未发布
	PREPARING   = "PREPARING"   // 准备中
	PROGRESSING = "PROGRESSING" // 进行中
	ENDED       = "ENDED"       // 已结束
)

const (
	CREDIT = "CREDIT" //三级分润
	DEBIT  = "DEBIT"  //直推分润
)

const (
	// 存款码长度
	Len_SaveMoneyCode = 10
)

const (
	SwapType_Up   = "1"
	SwapType_Down = "2"
)

const (
	AppType_App = "0"
	AppType_Pos = "1"
)

const (
	FUNCTIONREG = "reg"
	BACKPWD     = "backpwd"
	PAYPWD      = "paypwd"
	CHECKPHONE  = "checkphone"
	MODIFYPHONE = "modifyphone"
)
const ( //商家使用的手机发送验证码
	BACKPWD_Business    = "backpwd_business"    //修改登录密码
	PAYPWD_Business     = "paypwd_business"     //修改支付密码
	CHECKPHONE_Business = "checkphone_business" //确认商家旧手机号

	MODIFYPHONE_Business = "modifyphone_business" //修改手机号
)

//邮箱功能类型FuntionType
const (
	Reg_By_Mail             = "reg_mail"         //注册
	Backpwd_By_Mail         = "backpwd_mail"     //修改密码
	ModifyEmail_By_NewEmail = "modifyemail_mail" //修改邮箱，发送邮件到新邮箱
	Paypwd_By_Mail          = "paypwd_mail"      //修改支付密码

	//下面暂时未用
	//Checkphone_By_Mail      = "checkphone_mail"  //验证邮箱
	//Modifyphone_By_Mail     = "modifyphone_mail" //修改手机号

)

const (
	MODIFYACCOUNTLOGTYPE = 0 //账号类型
	MODIFYPAYLOGTYPE     = 1 // 支付账号类型
	REMAINTYPE           = 2 // 余额变动
)

// 交易明细订单类型
const (
	BILL_TYPE_INCOME        = "1" // 用户存款
	BILL_TYPE_OUTGO         = "2" // 用户取款
	BILL_TYPE_PROFIT        = "3" // 收益(佣金)
	BILL_TYPE_RECHARGE      = "4" // 充值
	BILL_TYPE_WITHDRAWALS   = "5" // 提现
	BILL_TYPE_ChangeBalance = "6" // 平台修改余额
)

const (
	INIT_UUID  = "00000000-0000-0000-0000-000000000000"
	INIT_UUID1 = "00000000-1111-1111-1111-000000000000" // 出金时插入明细,merchantNo不为空,用这个做无效的uuid
)

// 收款方式 ,1-支票;2-现金;3-银行转账;4-其他
const (
	CHECK_COLLECTION_TYPE = "1" //支票
	CASH_COLLECTION_TYPE  = "2" //现金
	BANK_COLLECTION_TYPE  = "3" //银行转账
	OTHER_COLLECTION_TYPE = "4" //其他
)

// 渠道类型，0-通用,1-用户,总部;2-pos
const (
	CHANNEL_ALL              = "0" //通用
	CHANNEL_USE_HEADQUARTERS = "1" //用户,总部
	CHANNEL_POS              = "2" //pos
)

const (
	WITHDRAW_PHONE     int32 = 0
	WITHDRAW_SWEEP     int32 = 1
	WITHDRAW_SWEEP_ALL int32 = 2
)

const (
	FEES_COUNT_INIT  = 0 // 初始状态
	FEES_IS_COUNT    = 1 // 已统计 (后续可能需要做财务审核)
	FEES_COUNT_CLEAR = 2 // 财务已确认
	FEES_FAILE_COUNT = 3 // 统计失败
)

const (
	SCORE_SETTING_BUY_RATE  = "score_setting_buy_rate"  // 购买费率/万分比 1
	USD_TRANSFER_RATE       = "usd_transfer_rate"       // usd转账费率 2
	USD_RECV_RATE           = "usd_recv_rate"           // usd收款手续费 3
	USD_DEPOSIT_RATE        = "usd_deposit_rate"        // usd存款费率 4
	USD_PHONE_WITHDRAW_RATE = "usd_phone_withdraw_rate" // usd手机号取款费率 5
	USD_FACE_WITHDRAW_RATE  = "usd_face_withdraw_rate"  //	usd面对面取款费率 6
	KHR_TRANSFER_RATE       = "khr_transfer_rate"       //	khr转账手续费率 7
	KHR_RECV_RATE           = "khr_recv_rate"           //	khr收款手续费 8
	KHR_DEPOSIT_RATE        = "khr_deposit_rate"        //	khr存款手续费率 9
	KHR_PHONE_WITHDRAW_RATE = "khr_phone_withdraw_rate" //	khr手机号取款手续费率 10
	KHR_FACE_WITHDRAW_RATE  = "khr_face_withdraw_rate"  //	khr面对面取款手续费率 11
)

// 手续费类型
const (
	Fees_Type_SCORE_SETTING_BUY_RATE  = 1  // 购买费率/万分比 1
	Fees_Type_USD_TRANSFER_RATE       = 2  // usd转账费率 2
	Fees_Type_USD_RECV_RATE           = 3  // usd收款手续费 3
	Fees_Type_USD_DEPOSIT_RATE        = 4  // usd存款费率 4
	Fees_Type_USD_PHONE_WITHDRAW_RATE = 5  // usd手机号取款费率 5
	Fees_Type_USD_FACE_WITHDRAW_RATE  = 6  //	usd面对面取款费率 6
	Fees_Type_KHR_TRANSFER_RATE       = 7  //	khr转账手续费率 7
	Fees_Type_KHR_RECV_RATE           = 8  //	khr收款手续费 8
	Fees_Type_KHR_DEPOSIT_RATE        = 9  //	khr存款手续费率 9
	Fees_Type_KHR_PHONE_WITHDRAW_RATE = 10 //	khr手机号取款手续费率 10
	Fees_Type_KHR_FACE_WITHDRAW_RATE  = 11 //	khr面对面取款手续费率 11
	Fees_Type_Usd_To_Khr_Count_Fee    = 12 //	usd-->khr 单笔手续费 12
	Fees_Type_Khr_To_Usd_Count_Fee    = 13 //	khr-->usd 单笔手续费 13
)

// 场景
const (
	Scene_Save     = 1 //存款
	Scene_Withdraw = 2 // 提现
	Scene_Transfer = 3 // 转账
)

const (
	KEY_USD_MIN_DEPOSIT_FEE        = "usd_min_deposit_fee"        //usd存款最低收取金额
	KEY_USD_PHONE_MIN_WITHDRAW_FEE = "usd_phone_min_withdraw_fee" //usd手机号取款最低收取金额
	KEY_USD_MIN_TRANSFER_FEE       = "usd_min_transfer_fee"       //usd转账最低收取金额
	KEY_USD_FACE_MIN_WITHDRAW_FEE  = "usd_face_min_withdraw_fee"  // usd面对面取款最低收取金额
	KEY_KHR_MIN_DEPOSIT_FEE        = "khr_min_deposit_fee"        // khr存款最低收取金额
	KEY_KHR_PHONE_MIN_WITHDRAW_FEE = "khr_phone_min_withdraw_fee" //khr手机号取款最低收取金额
	KEY_KHR_MIN_TRANSFER_FEE       = "khr_min_transfer_fee"       //khr转账最低收取金额
	KEY_KHR_FACE_MIN_WITHDRAW_FEE  = "khr_face_min_withdraw_fee"  //khr面对面取款最低收取金额
	KEY_USD_TO_KHR_FEE             = "usd_to_khr_fee"             // usd转khr单笔手续费
	KEY_KHR_TO_USD_FEE             = "khr_to_usd_fee"             // khr转usd单笔手续费
)

const (
	KEY_WriteOff_DurationDate = "write_off_duration_date" // 核销码有效期
)

const (
	KEY_STORE_UNAUTH_IMAGE_PATH = "store_unauth_image_path" // 不授权的图片路径
	KEY_STORE_AUTH_IMAGE_PATH   = "store_auth_image_path"   //需要授权的图片路径
)

//用户消息中心消息类型
const (
	LOG_APP_MESSAGES_ORDER_TYPE_EXCHANGE       = "1" // 兑换
	LOG_APP_MESSAGES_ORDER_TYPE_EXCHANGE_Fail  = "2" // 兑换失败
	LOG_APP_MESSAGES_ORDER_TYPE_EXCHANGE_Apply = "3" // 兑换成功

	LOG_APP_MESSAGES_ORDER_TYPE_INCOME       = "4" // 存款、充值
	LOG_APP_MESSAGES_ORDER_TYPE_INCOME_Fail  = "5" // 存款、充值失败
	LOG_APP_MESSAGES_ORDER_TYPE_INCOME_Apply = "6" // 存款、充值成功

	LOG_APP_MESSAGES_ORDER_TYPE_OUTGO       = "7" //取款、提现
	LOG_APP_MESSAGES_ORDER_TYPE_OUTGO_Fail  = "8" //取款、提现失败
	LOG_APP_MESSAGES_ORDER_TYPE_OUTGO_Apply = "9" //取款、提现成功

	LOG_APP_MESSAGES_ORDER_TYPE_TRANSFER       = "10" // 转账
	LOG_APP_MESSAGES_ORDER_TYPE_TRANSFER_Fail  = "11" // 转账失败
	LOG_APP_MESSAGES_ORDER_TYPE_TRANSFER_Apply = "12" // 转账成功

	LOG_APP_MESSAGES_ORDER_TYPE_COLLECTION       = "13" // 收款
	LOG_APP_MESSAGES_ORDER_TYPE_COLLECTION_Fail  = "14" // 收款失败
	LOG_APP_MESSAGES_ORDER_TYPE_COLLECTION_Apply = "15" // 收款成功

	LOG_APP_MESSAGES_ORDER_TYPE_Cust_Withdraw       = "16" // 用户线上提现申请
	LOG_APP_MESSAGES_ORDER_TYPE_Cust_Withdraw_Fail  = "17" // 用户线上提现失败
	LOG_APP_MESSAGES_ORDER_TYPE_Cust_Withdraw_Apply = "18" // 用户线上提现成功

	LOG_APP_MESSAGES_ORDER_TYPE_Cust_Save       = "19" // 用户线上充值申请
	LOG_APP_MESSAGES_ORDER_TYPE_Cust_Save_Fail  = "20" // 用户线上充值失败
	LOG_APP_MESSAGES_ORDER_TYPE_Cust_Save_Apply = "21" // 用户线上充值成功

	LOG_APP_MESSAGES_ORDER_TYPE_Business_Refund_Cust = "22" // 商家退用户钱成功

	LOG_APP_MESSAGES_ORDER_TYPE_Cust_Pay = "23" // 用户支付成功
)

//多语言配置类型
const (
	LANG_TYPE_WORD = "1" //文字
	LANG_TYPE_IMG  = "2" //图片
	LANG_TYPE_Err  = "3" //错误码
)

const (
	//Lang_En = "en"
	//Lang_Km = "km"
	//Lang_Cn = "zh_CN"

	// 规范化语言类型
	LangEnUS = "en_US"
	LangKmKH = "km_KH"
	LangZhCN = "zh_CN"

	// 系统默认语言为英语
	DefaultLang = LangEnUS
)

//xlsx文件类型
const (
	XLSX_FILE_TYPE_EXCHANGE        = "1" //兑换
	XLSX_FILE_TYPE_INCOME          = "2" //存
	XLSX_FILE_TYPE_OUTGO           = "3" //取
	XLSX_FILE_TYPE_TRANSFER        = "4" //转
	XLSX_FILE_TYPE_COLLECTION      = "5" //收
	XLSX_FILE_TYPE_TO_HEADQUARTERS = "6" //服务商充值
	XLSX_FILE_TYPE_TO_SERVICER     = "7" //服务商提款
	XLSX_FILE_TYPE_VACCOUNT_LOG    = "8" //虚拟账户日志流水
)

const (
	Risk_Ctrl_Sweep_Withdrawal      = "sweep_withdrawal"
	Risk_Ctrl_Save_Money            = "save_money"
	Risk_Ctrl_Mobile_Num_Withdrawal = "mobile_num_withdrawal"
	Risk_Ctrl_Exchange              = "exchange"
	Risk_Ctrl_Transfer              = "transfer"
	Risk_Ctrl_Collection            = "collection"
)

const (
	Exchange_Khr_To_Usd = "khr_to_usd"
	Exchange_Usd_To_Khr = "usd_to_khr"
)

const (
	AppVersionVsType_app          = "0" //moderpay app
	AppVersionVsType_pos          = "1" //moderpay pos
	APPVERSIONVSTYPE_MANGOPAY_APP = "2" //mangopay app
	APPVERSIONVSTYPE_MANGOPAY_POS = "3" //mangopay pos
)
const (
	AppVersionSystem_Android = "0"
	AppVersionSystem_Ios     = "1"
)

//
const (
	Status_Disable = "0" //禁用
	Status_Enable  = "1" //启用

	IsDel_No  = "0" //未删除
	IsDel_Yes = "1" //已删除
)

//是否强制更新
const (
	AppVersionIsForce_False = "0"
	AppVersionIsForce_True  = "1"
)

//选择版本更新
const (
	AppVersionUpType_Big   = "1" //大版本更新
	AppVersionUpType_Small = "2" //小版本更新
	AppVersionUpType_Bug   = "3" //bug版本更新
)

//初始化版本
const (
	AppVersion_Init = "1.0.0"
)

//WEB关键操作日志type
const (
	//用户中心
	LogAccountWebType_Account   = "1" //用户相关
	LogAccountWebType_Servicer  = "2" //服务商相关
	LogAccountWebType_Financial = "3" //财务管理相关
	LogAccountWebType_Config    = "4" //配置管理相关
	//账号中心
	LogAccountWebType_Account_Menu = "5" //账号与菜单相关
	//交易中心
	LogAccountWebType_Trading_Order = "6" //订单审核相关

	LogAccountWebType_Business = "7" //商家相关
)

//协议类型
const (
	AgreementType_Use           = "0" //用户协议
	AgreementType_Privacy       = "1" //隐私协议
	AgreementType_Auth_Material = "2" //实名认证协议
)

//卡
const (
	CardUsable_Disable = "0" //禁用
	CardUsable_Enable  = "1" //启用
)

const (
	UploadImage_Auth   = 1 //需要授权认证的图片
	UploadImage_UnAuth = 2 //不需授权认证的图片
)

//上传的图片base64字符串最大值
const (
	UploadImgBase64LengthMax = 500000 * 4 / 3 //最大值
)

//是否添加水印
const (
	AddWatermark_True = "1" //是
)

//是否是推荐的
const (
	IsRecom_True  = "1" //是
	IsRecom_False = "0" //否
)

const (
	Template_Reg                  = "reg"
	Template_ExchangeSuccess      = "exchangeSuccess"      // 兑换成功
	Template_TransferSuccess      = "transferSuccess"      // 转账成功
	Template_SmsWriteOff          = "smsWriteOff"          // 短信核销码
	Template_WithdrawApplySuccess = "withdrawApplySuccess" // 提现申请成功
	Template_WithdrawSuccess      = "withdrawSuccess"      // 提现成功
	Template_WithdrawFail         = "withdrawFail"         //  提现失败
	Template_AddSuccess           = "addSuccess"           // 充值成功
	Template_AddFail              = "addFail"              // 充值失败
	Template_CollectSuccess       = "collectSuccess"       // 收款成功
	Template_VerifyEmail          = "verifyEmail"          // 验证你的电子邮件地址
	Template_RefundSuccess        = "refundSuccess"        // 退款成功
	Template_PaySuccess           = "paySuccess"           // 付款成功
)

//银行卡存取款渠道的卡所支持业务类型
const (
	SupportType_In     = "1" //只支持存款
	SupportType_Out    = "2" //只支持取款
	SupportType_Common = "3" //通用
)

const (
	// 公共数据
	Common_Data = "common_data"
)

const (
	//账号状态
	AccountUseStatusAlwaysDisabled    = 0 // 永久禁用
	AccountUseStatusNormal            = 1 // 正常使用
	AccountUseStatusTemporaryDisabled = 2 // 临时禁用
)

const (
	TradingAuthorityForbid = 0 // 禁止交易
	TradingAuthorityAllow  = 1 // 允许交易
)

const (
	AccountNoActived = "0" // 未激活
	AccountActived   = "1" // 激活
)

const (
	ExchangeAmountUsdToKhr = 1
	ExchangeAmountKhrToUsd = 2
)

const (
	RoleType_Merc   = "1" // 商户
	RoleType_Agency = "2" // 代理
	RoleType_Oper   = "3" // 运营
	RoleType_Admin  = "4" // 管理员
)

const (
	SignMethodRsa = "rsa"
	SignMethodMd5 = "md5"
)

const (
	CHANNEL_STATUS_FORBID = "0" // 禁止
	CHANNEL_STATUS_NORMAL = "1" // 正常
)

//经营期限类型
const (
	TermType_Short = "1"
	TermType_Long  = "2"
)

//是否已初始化支付密码(0否，1是)
const (
	InitPayPwdStatus_false = "0"
	InitPayPwdStatus_true  = "1"
)

//平台收益来源headquarters_profit
const ( //收益来源
	ProfitSource_WithdrawFee          = "1"  //客户提现手续费
	ProfitSource_TRANSFERFee          = "2"  //客户转账手续费
	ProfitSource_Exchange             = "3"  //客户兑换手续费
	ProfitSource_COLLECTION           = "4"  //客户收款手续费
	ProfitSource_INCOME               = "5"  //客户存款手续费
	ProfitSource_ToBusinessFee        = "6"  //商家提现手续费
	ProfitSource_BusinessToFee        = "7"  //商家充值手续费
	ProfitSource_BusinessTransferFee  = "8"  //商家转账手续费
	ProfitSource_ModernPayOrderFee    = "9"  //商家ModernPay交易手续费
	ProfitSource_ModernPayOrderRefund = "10" //商家ModernPay退款
)

const (
	//平台收益操作类型
	PlatformProfitAdd   = "1" // +
	PlatformProfitMinus = "2" // -
)

//公告发布状态
const (
	BulletinUseStatus_UnBulletin = "0" //未发布
	BulletinUseStatus_All        = "1" //发布所有
	BulletinUseStatus_personal   = "2" //只发布个人商家
	BulletinUseStatus_Enterprise = "3" //只发布企业商家
)

//公告置顶状态
const (
	BulletinTopStatus_False = "0" //否
	BulletinTopStatus_True  = "1" //是
)

//上传文件类型状态
const (
	UploadFileType_IPA  = "1" //ipa
	UploadFileType_APK  = "2" //apk
	UploadFileType_HTML = "3" //html
	UploadFileType_XLSX = "4" //xlsx
)

//签名方式
const (
	SignMethod_RSA2 = "RSA2"
)

//银行卡卡号长度限制(普通渠道的)
const (
	CardNumberLengthMax = 25 //最大值
	CardNumberLengthMin = 9  //最小值
)

//用户指纹状态
const (
	AppFingerprintUseStatus_Disable = "0" //禁用
	AppFingerprintUseStatus_Enable  = "1" //启用
)

//是否允许用户录入指纹获取标识
const (
	AppFingerprintOn_True  = "true"  //允许
	AppFingerprintOn_False = "false" //禁止
)

//用户、商家使用的银行卡渠道类型
const (
	CHANNELTYPE_ORDINARY   = "1" //普通渠道
	CHANNELTYPE_THIRDPARTY = "2" //第三方渠道
)

//银行卡卡号长度限制(第三方渠道的)
const (
	ThirdPartyCardNumberLengthMin = 8  //最小值
	ThirdPartyCardNumberLengthMax = 12 //最大值
)
