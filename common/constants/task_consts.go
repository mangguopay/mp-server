package constants

const (
	CBTASK_BILL     = "1" // 支付
	CBTASK_TRANSFER = "2" // 代付
	CBTASK_WITHDRAW = "3" // 提现

	// 未执行
	AgentTaskStatus_Init = "0"
	// 已执行
	AgentTaskStatus_Ran = "1"
	// 挂起
	AgentTaskStatus_Hold = "2"

	// 未审核
	AgentTaskAuditStatus_Init = "0"
	// 不通过
	AgentTaskAuditStatus_NotPass = "1"
	// 通过
	AgentTaskAuditStatus_Passed = "2"
	// 审核中
	AgentTaskAuditStatus_Auditing = "3"

	ExpBillStatus_Long  = "1"
	ExpBillStatus_Short = "2"
)

const (
	//初始化
	TASK_STATUS_INIT = "0"
	//成功
	TASK_STATUS_SUCCESS = "1"
	//异常挂起
	TASK_STATUS_ERR = "2"
	//执行中
	TASK_STATUS_EXECUTE       = "3"
	TASK_STATUS_PARTLY_FAILED = "4"
	TASK_STATUS_FAILED        = "5"
	TASK_STATUS_UNINIT        = "6"

	TASK_STATUS_MSG_SUCCESS       = "成功"
	TASK_STATUS_MSG_EXECUTE       = "执行中"
	TASK_STATUS_MSG_PARTLY_FAILED = "部分失败"
	TASK_STATUS_MSG_FAILED        = "失败"
)

const (
	//日结
	AGENT_TASK_TYPE_DAY = "1"

	//月结
	AGENT_TASK_TYPE_MONTH = "2"

	//渠道出金
	AGENT_TASK_TYPE_OUT = "3"
)

const (
	T1Audit_Passed = "1"
	T1Audit_Audit  = "2"
)
