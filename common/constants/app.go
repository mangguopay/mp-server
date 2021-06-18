package constants

// 商家应用审核状态
const (
	BusinessAppStatus_Pending = "0" //未审核
	BusinessAppStatus_Passed  = "1" //通过
	BusinessAppStatus_Deny    = "2" //不通过
	BusinessAppStatus_Invalid = "3" //作废
	BusinessAppStatus_Up      = "4" //上架
	BusinessAppStatus_Delete  = "5" //删除
)

//应用签约状态
const (
	SignedStatusPending = "1" //1申请中
	SignedStatusDeny    = "2" //2未通过
	SignedStatusPassed  = "3" //3已通过
	SignedStatusInvalid = "4" //4已过期
)
