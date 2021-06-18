package constants

const (
	NATS_POOLNAME_MAIN   = "main_server"
	NATS_POOLNAME_SETTLE = "settle_mq"
	NATS_POOLNAME_LOG    = "log_mq"
	NATS_POOLNAME_RISK   = "risk_mq"
	NATS_POOLNAME_WALLET = "wallet_mq"

	NATS_SUBJECT_CLEARING = "subject_clearing"
	NATS_SUBJECT_LOG      = "subject_log"

	NATS_QGROUP_SETTLE_1 = "settle_1"
	NATS_QGROUP_LOG_1    = "log_1"
	NATS_QGROUP_RISK_1   = "risk_1"
	NATS_QGROUP_WALLET_1 = "wallet_1"

	NATS_ARGS_KEY_REQ  = "req"
	NATS_ARGS_KEY_TYPE = "type"

	NATS_REQ_CLEARING = "1"
	NATS_REQ_AGENT    = "2"

	NATS_REQ_LOG_LOG = "1"
)
