package constants

const (
	// 初始化
	OrderStatus_Init = "1"
	// 等待
	OrderStatus_Pending = "2"
	// 已支付
	OrderStatus_Paid = "3"
	// 失败
	OrderStatus_Err             = "4"
	OrderStatus_Pending_Confirm = "5" // 待确认
	OrderStatus_Cancel          = "6" // 取消
)

//付款方式
const (
	ORDER_PAYMENT_TYPE_CASH     = "1" //现金
	ORDER_PAYMENT_BALANCE       = "2" //余额
	ORDER_PAYMENT_BANK_TRANSFER = "3" // 银行卡转账
	ORDER_PAYMENT_BANK_WITHDRAW = "4" // 银行卡提现
)

const (
	AuditOrderStatus_Pending = "0" //未审核
	AuditOrderStatus_Passed  = "1" //通过
	AuditOrderStatus_Deny    = "2" //不通过
)

const (
	//商家提现订单状态
	WithdrawalOrderStatusPending = "0" //审核中，未受理
	WithdrawalOrderStatusPassed  = "1" //审核通过，已受理
	WithdrawalOrderStatusDeny    = "2" //审核不通过，驳回
	WithdrawalOrderStatusSuccess = "3" //完成转账， 提现成功
	WithdrawalOrderStatusFail    = "4" //转账失败， 提现失败
)

//账号的实名认证状态、个人商家、企业商家认证状态
const (
	AuthMaterialStatus_Pending       = "0" //未审核
	AuthMaterialStatus_Passed        = "1" //通过
	AuthMaterialStatus_Deny          = "2" //不通过
	AuthMaterialStatus_UnAuth        = "3" //未认证
	AuthMaterialStatus_Appeal_Passed = "4" //申诉作废通过的实名认证
)

const (
	Risk_Pay_Type_Transfer              = "1" // 转账
	Risk_Pay_Type_Exchange              = "2" // 兑换
	Risk_Pay_Type_Save_Money            = "3" // 存款
	Risk_Pay_Type_Mobile_Num_Withdrawal = "4" // 手机号取款
	Risk_Pay_Type_Sweep_Withdrawal      = "5" // 扫一扫取款
	Risk_Pay_Type_Collection            = "6" // 收款
)
const (
	Risk_Result_Pass    = "0" // 风控通过
	Risk_Result_No_Pass = "1" // 风控不通过
)
const (
	Risk_Result_Pass_Str    = "passed"    // 风控通过
	Risk_Result_No_Pass_Str = "notpassed" // 风控不通过
	Risk_Result_Pending_Str = "pending"   // 异步状态的结果
)

//收款方式
const ( //to_headquarters,to_servicer,to_business
	COLLECTION_TYPE_CHECK         = "1" // 支票
	COLLECTION_TYPE_CASH          = "2" // 现金
	COLLECTION_TYPE_BANK_TRANSFER = "3" // 银行转账
	COLLECTION_TYPE_OTHER         = "4" // 其他
)

const (
	Order_Type_In_Come         = "1"
	Order_Type_Mobile_Withdraw = "2"
	Order_Type_Sweep_Withdraw  = "3"
)

const (
	WITHDRAWAL_TYPE_ORDINARY = "1" //普通提现
	WITHDRAWAL_TYPE_ALL      = "2" //全部提现
)

const (
	TRANSACTION_TYPE_MOBILE_PHONE_WITHDRAW = "1" // 手机号取款
	TRANSACTION_TYPE_SAVE_MONEY            = "2" // 存款
	TRANSACTION_TYPE_TRANSFER              = "3" // 转账
	TRANSACTION_TYPE_SWEEP_WITHDRAW        = "4" // 扫码提现

)

const (
	Charge_Type_Rate  = "1" // 按比例计算手续费
	Charge_Type_Count = "2" // 按单笔计算手续费
)
const (
	Cust_Op_Type_Save     = 1
	Cust_Op_Type_Withdraw = 2
)

//订单冻结状态
const (
	SettlementMethodWithdraw = "1" //手动
	SettlementMethodTransfer = "2" //代付
)

const (
	// business_bill表的订单状态
	BusinessOrderStatusPending       = "1" // 待支付
	BusinessOrderStatusPay           = "2" // 支付成功
	BusinessOrderStatusPayTimeOut    = "3" // 支付失败
	BusinessOrderStatusRefund        = "4" // 已全额退款退款
	BusinessOrderStatusRebatesRefund = "5" // 部分退款退款

	//订单通知状态
	NotifyStatusNOT     = "0" //0未通知
	NotifyStatusSuccess = "1" //1成功
	NotifyStatusDoing   = "2" //2通知进行中
	NotifyStatusTimeout = "3" //3超时

	//订单过期时间
	BusinessOrderExpireTime = 30 //30分钟
)

const (
	//取款类型 0-手机号提现;1-扫码提现;2-全部提现(现是扫码提现,手机号提现是一个核销码对应的订单金额必需一次性全取出来)
	OutgoOrderPaymentType_MobileNum = "0" // 手机号提现
	OutgoOrderPaymentType_Sweep     = "1" // 扫码
	OutgoOrderPaymentType_SweepAll  = "2" // 扫码全部提现
)

const (
	//商家转账订单类型
	BusinessTransferOrderTypeOrdinary   = "1" //普通付款
	BusinessTransferOrderTypeEnterprise = "2" //企业付款

	//商家转账订单状态
	BusinessTransferOrderStatusPending = "0" //处理中
	BusinessTransferOrderStatusSuccess = "1" //1成功
	BusinessTransferOrderStatusFail    = "2" //2失败
)

const (
	//商家退款订单状态
	BusinessRefundStatusPending = "0" //处理中
	BusinessRefundStatusSuccess = "1" //1成功
	BusinessRefundStatusFail    = "2" //2失败
)

const (
	//商家批量转账订单状态
	BusinessBatchTransferOrderStatusPending    = "0" //等待，初始化
	BusinessBatchTransferOrderStatusSuccess    = "1" //1成功
	BusinessBatchTransferOrderStatusFail       = "2" //2失败
	BusinessBatchTransferOrderStatusPaySuccess = "3" //3处理中（钱已支付才会处理中）
)

// 交易的订单类型
const (
	TradeTypeModernpayAPP        = "MANGOPAY_APP"            // APP支付
	TradeTypeModernpayMWEB       = "MANGOPAY_MWEB"           // H5支付
	TradeTypeModernpayFaceToFace = "MANGOPAY_FACE_TO_FACE"   // 面对面扫码支付
	TradeTypeEnterprisePay       = "MANGOPAY_ENTERPRISE_PAY" // 企业支付

	TradeTypeWeiXinFaceToFace = "WEIXIN_FACE_TO_FACE" // 微信当面付
	TradeTypeWeiXinApp        = "WEIXIN_APP"          // 微信APP支付

	TradeTypeAlipayFaceToFace = "ALIPAY_FACE_TO_FACE" // 支付宝当面付
	TradeTypeAlipayAPP        = "ALIPAY_APP"          // 支付宝APP支付

	TradeTypeBusinessPay = "MANGOPAY_PERSONAL_BUSINESS_PAY" //商家收款

)

//自动签约的产品名称(产品存的是多语言的key)
const (
	AutoSignedSceneName = "商家收款" // 自动签约的产品名称

)

//支付方式
const (
	PayMethodBalance  = "BALANCE"   //modernPay余额
	PayMethodBankCard = "BANK_CARD" //银行卡
)
