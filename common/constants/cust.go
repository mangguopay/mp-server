package constants

const (
	//用户转账入权限
	CustInTransferAuthorizationDisabled = 0 //禁用
	CustInTransferAuthorizationEnable   = 1 //启用

	//用户转账出权限
	CustOutTransferAuthorizationDisabled = 0 //禁用
	CustOutTransferAuthorizationEnable   = 1 //启用

	//核销码相关操作
	WrittenOffCodeOpFreeze   = "freeze"   //冻结
	WrittenOffCodeOpUnFreeze = "unfreeze" //解冻
	WrittenOffCodeOpCancel   = "cancel"   //注销

	//核销码状态状态1-.2-,3,4,5
	WriteOffCodeWaitUse   = "1" //初始状态
	WriteOffCodeIsUse     = "2" //已使用状态
	WriteOffCodeExpired   = "3" //已过期
	WriteOffCodeFrozen    = "4" //已冻结
	WriteOffCodeCancelled = "5" //已注销
)
