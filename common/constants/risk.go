package constants

const (
	// 黑白名单
	RiskFilter_BWList = 1
	// 每单
	RiskFilter_PerOrder = 2
	// 统计量
	RiskFilter_Stat = 3
)

const (
	RiskEffType_MercNo    = 1
	RiskEffType_ChannelNo = 2
)

const (
	RiskBw_Black = 0
	RiskBw_White = 1
)

const (
	// 未命中
	RiskPassReason_Noop = 0
	// 命中白名单
	RiskPassReason_WhiteList = 1
	// 命中黑名单
	RiskPassReason_BlackList = 2
	// 获取规则失败
	RiskPassReason_GetRuleFailed = 3
	// 单笔值不满足
	RiskPassReason_NotFitPerOrder = 4
	// 统计值不满足
	RiskPassReason_NotFitStat = 5
)

const (
	RiskBwRuleType_RecvPhone = 1
	RiskBwRuleType_RecvAccNo = 2
	RiskBwRuleType_PayPhone  = 3
	RiskBwRuleType_PayAccNo  = 4
	RiskBwRuleType_MercNo    = 5
	RiskBwRuleType_ChannelNo = 6
)

const (
	RiskRuleOp_Bigger         = 1
	RiskRuleOp_Smaller        = 2
	RiskRuleOp_Equal          = 3
	RiskRuleOp_BiggerOrEqual  = 4
	RiskRuleOp_SmallerOrEqual = 5
	RiskRuleOp_NotEqual       = 6
)

const (
	RiskRuleType_Amount = 1
	RiskRuleType_Time   = 2
)

const (
	RiskStatRuleType_Amount = 1
	RiskStatRuleType_Cnt    = 2
)

const (
	RiskTimeType_All      = 0
	RiskTimeType_Min      = 1
	RiskTimeType_HalfHour = 2
	RiskTimeType_Hour     = 3
	RiskTimeType_Day      = 4
)

const (
	// 初始化
	RiskStatStatus_Init = 1
	// 已统计
	RiskStatStatus_Handled = 2
	RiskStatStatus_Doing   = 3
)

const (
	RiskExecuteType_Offline  = "offline"  // 离线
	RiskExecuteType_Half     = "half"     // 半实时
	RiskExecuteType_Online   = "online"   // 实时
	RiskExecuteType_Hesitate = "hesitate" // 犹豫中
)

const (
	RiskEvaluator_Amount = "amount"
)
