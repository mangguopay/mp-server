package constants

// 虚拟账户类型
const (
	// ----------------------用户才有---------------------------
	// 1.正常注册时初始化1和2
	// 2.非正注册时初始化3和4，待激活时再初始化1和2

	// 美金存款 -- 个人
	VaType_USD_DEBIT = 1
	// 瑞尔存款 -- 个人
	VaType_KHR_DEBIT = 2

	// 美金冻结存款 -- 个人(用户未激活时)
	VaType_FREEZE_USD_DEBIT = 3
	// 瑞尔冻结存款 -- 个人(用户未激活时)
	VaType_FREEZE_KHR_DEBIT = 4
	// --------------------------------------------------------

	// ----------------------服务商才有-------------------------
	// 管理后台初次授权账号的角色是服务商时 初始化5、6、7、8
	// 美金额度
	VaType_QUOTA_USD = 5
	// 瑞尔额度
	VaType_QUOTA_KHR = 6

	// 美金实时AA300088额度(此账号的余额可以为负数)
	VaType_QUOTA_USD_REAL = 7
	// 瑞尔实时额度(此账号的余额可以为负数)
	VaType_QUOTA_KHR_REAL = 8

	// 美金存款 -- 服务商
	// VaType_USD_DEBIT_SRV = 9 // 已经没有使用
	// 瑞尔存款 -- 服务商
	// VaType_KHR_DEBIT_SRV = 10 // 已经没有使用
	// ---------------------------------------------------------

	// ----------------------总部账户和服务商都有？待确定------------------------
	// 因为历史原因，默认是使用到账户时才去初始化的
	// usd 手续费类型
	VaType_USD_FEES = 11

	// khr 手续费类型
	VaType_KHR_FEES = 12
	// ---------------------------------------------------------

	// ----------------------商家才有----------------------------
	// 1.企业商家在注册时初始化
	// 2.个人商家在申请为商家通过时初始化

	// 美金--商家已结算
	VaType_USD_BUSINESS_SETTLED = 13

	// 瑞尔--商家已结算
	VaType_KHR_BUSINESS_SETTLED = 14

	// 美金--商家未结算
	VaType_USD_BUSINESS_UNSETTLED = 15

	// 瑞尔--商家未结算
	VaType_KHR_BUSINESS_UNSETTLED = 16
	// --------------------------------------------------------
)

const (
	VaReason_Exchange                        = "1"  // 兑换
	VaReason_INCOME                          = "2"  //充值
	VaReason_OUTGO                           = "3"  //提现
	VaReason_TRANSFER                        = "4"  //转账
	VaReason_COLLECTION                      = "5"  //收款
	VaReason_FEES                            = "6"  // 手续费
	VaReason_Cancel_withdraw                 = "7"  // pos 端取消提现
	VaReason_PROFIT_OUTGO                    = "8"  //平台盈利提现
	VaReason_Cust_Withdraw                   = "9"  // 客户向总部提现
	VaReason_Cust_Cancel_Withdraw            = "10" // 驳回客户向总部提现
	VaReason_Cust_Save                       = "11" // 客户向总部充值
	VaReason_Srv_Save                        = "12" // 服务商向总部充值
	VaReason_Srv_Withdraw                    = "13" // 服务商向总部提现
	VaReason_Cust_Cancel_Save                = "14" // 驳回客户向总部存款
	VaReason_Cust_Pay_Order                  = "15" // 客户支付商家的订单
	VaReason_Business_Payee                  = "16" // 商家收到客户的订单
	VaReason_Srv_CashRecharge                = "17" // 服务商现金充值
	VaReason_Business_Settle                 = "18" // 商家结算
	VaReason_Business_Save                   = "19" // 商家充值
	VaReason_Business_Withdraw               = "20" // 商家提现
	VaReason_Business_Cancel_Withdraw        = "21" // 驳回商家提现
	VaReason_ChangeCustBalance               = "22" // 改变用户余额
	VaReason_ChangeSrvBalance                = "23" // 改变服务商余额
	VaReason_BusinessTransferToBusiness      = "24" // 商家转账(注意现在商家转个人 是转用户账户而不是商家账户)
	VaReason_BusinessRefund                  = "25" // 商家退款
	VaReason_BusinessBatchTransferToBusiness = "26" // 商家付款批量转账
	VaReason_PlatformFreeze                  = "27" // 平台冻结
)

const (
	// 1:+;2:-;3:冻结;4:解冻;

	VaOpType_Add = "1" // 增加balance字段的金额

	VaOpType_Minus = "2" // 减少balance字段的金额

	VaOpType_Freeze = "3" // 减少balance字段的金额，增加到frozen_balance字段

	VaOpType_Defreeze = "4" // 减少frozen_balance字段的金额

	// 解冻并扣减
	VaOpType_Defreeze_Minus = "5" // 减少frozen_balance字段的金额

	// 解冻并增加
	VaOpType_Defreeze_Add = "6" // 增加balance字段的金额, 减少frozen_balance字段的金额

	// 解冻不减
	VaOpType_Defreeze_But_Minus = "7" // 增加frozen_balance字段的金额

	// 服务商提现,冻结+,实时+
	VaOpType_Balance_Frozen_Add = "8" // 增加balance字段的金额, 增加frozen_balance字段的金额, 金额一致

	// 服务商提现驳回,冻结-,实时+
	VaOpType_Balance_Defreeze_Add = "9" // 减少balance字段的金额, 减少frozen_balance字段的金额, 金额一致
)

const (
	CodeType_Recv  = "1"
	CODETYPE_SWEEP = "2"
)
const (
	CODE_USE_STATUS_IS_NO_SWEEP = "1" // 初始化
	CODE_USE_STATUS_IS_SWEEP    = "2" // 已扫码
	CODE_USE_STATUS_IS_PAY      = "3" // 已支付
	CODE_EXP                    = "4" // 已过期
	CODE_Pendding_Confirm       = "5" // 待确认
	CODE_CANCEL                 = "6" // 取消
)

const (
	// 预存
	QuotaOp_PreSave = "1"
	// 存
	QuotaOp_Save = "2"
	// 取
	QuotaOp_Withdraw = "3"
	// 服务商预存
	QuotaOp_SvrPreSave = "4" //现服务商发起充值，不需对冻结金额改变了，所以此字段暂时不需要了。。
	// 服务商存
	QuotaOp_SvrSave = "5"
	// 服务商取
	QuotaOp_SvrWithdraw = "6"
	// 回滚
	QuotaOp_Rollback = "7"
	// 服务商取款取消
	QuotaOp_SvrWithdraw_Cancel = "8"
	// 服务商预取
	QuotaOp_SvrPreWithdraw  = "9"
	QuotaOp_CustPreSave     = "10" // 客户向总部预存款
	QuotaOp_CustSave        = "11" // 客户向总部存款
	QuotaOp_CustSave_Cancel = "12" // 后台对客户存款操作进行驳回操作

	QuotaOp_SvrCashRecharge = "13" // 服务商现金充值
	QuotaOp_BusinessSave    = "14" // 商家充值

	QuotaOp_ChangeCustBalanceAdd   = "15" // 改变用户余额(增加)
	QuotaOp_ChangeCustBalanceMinus = "16" // 改变用户余额(减少)
	QuotaOp_ChangeSvrBalanceAdd    = "17" // 改变服务商余额(增加)
	QuotaOp_ChangeSvrBalanceMinus  = "18" // 改变服务商余额(减少)

)

// 转账类型
const (
	TRANSFER_TYPE_BILL  = 1 // 交易转账
	TRANSFER_TYPE_CLEAR = 2 // 结算转账
)

const (
	FEES_TYPE_EXCHANGE   = "1"
	FEES_TYPE_TRANSFER   = "2"
	FEES_TYPE_COLLECTION = "3"
	FEES_TYPE_SAVEMONEY  = "4"
	FEES_TYPE_WITHDRAW   = "5"
	//FEES_TYPE_MOBILE_NUM_WITHDRAW = "5"
	//FEES_TYPE_SWEEP_WITHDRAW      = "6"
)
