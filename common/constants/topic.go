package constants

const (
	// 核心服务nats
	Nats_Cluster_Primary = "primary"

	// 非核心服务nats
	Nats_Cluster_Secondary = "secondary"

	// --------- 核心服务下的主题-------------------------------
	// 清分
	Nats_Topic_Settle = "go.micro.topic.settle"

	// 风控
	Nats_Topic_Risk = "go.micro.topic.risk"

	// redis过期key监听
	Nats_Topic_Listen_Exp_key = "go.micro.srv.listen_exp_key"

	// ----------------------------------------------------------

	// --------- 非核心服务下的主题-------------------------------
	// 信息推送
	Nats_Topic_Push = "go.micro.topic.push"

	// 日志
	Nats_Topic_Statlog = "go.micro.topic.statlog"

	// ----------------------------------------------------------

	// 清分
	// Topic_Settle = "go.micro.topic.settle"
	// 风控
	//Topic_Risk = "go.micro.topic.risk"

	// 信息推送
	//Topic_Push = "go.micro.topic.push"
	// 日志
	//Topic_Statlog = "go.micro.topic.statlog"

	//ClusterId_Settle = "settle"
	Settle_Type           = "settle"
	BusinessSettle        = "business_settle"
	PaySystemResultNotify = "PaySystemResultNotify"

	Nats_Listen_Exp_key = "listenExpKey"

	Nats_Broker_Header_Reg_SMS       = "RegSms"
	Nats_Broker_Header_Write_Off     = "WriteOff"
	Nats_Broker_Header_Send_Push_Msg = "SendPushMsg"
	Nats_Broker_Header_Risk          = "risk"
)

const (
	Topic_Event_Srv_Gis = "srv_gis"
)
