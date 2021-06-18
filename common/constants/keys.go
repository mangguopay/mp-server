package constants

import "time"

const (
	WebSignKey = "at78r23fas8fhg8^TasuSdfs782"
	WebHtmlKey = "asdfu892jkDHfasdfha"

	MobileKeyJwtSign = "hgwrur13bTEcqw#e6r5bfqewdh"
	MobileKeyJwtAes  = "123456789_mobile"

	ApikeyEncodeKey        = "encode_key"
	ApikeyEncodeKey2       = "encode_key2"
	ApikeyEncodeKey3       = "encode_key3"
	ApikeyFirstClassMercNo = "fc_merc_no"
	ApikeyOrgNo            = "org_no"
	ApikeyReqsystem        = "req_system"
	ApikeyDecodeKey        = "decode_key"
	ApikeyVerifyKey        = "verify_key"
	ApikeySignKey          = "sign_key"
	ApikeyURLType          = "url_type"
	ApikeyUrlTypeQuery     = "url_type_query"
	ApikeyCallingHost      = "calling_host"
	ApikeyBiz              = "biz"
	ApikeyCertId           = "cert_id"
	ApikeyGetUrlHost       = "get_url_host"
	ApikeyFtpHost          = "ftp_host"
	ApikeyFtpPort          = "ftp_port"
	ApikeyFtpUser          = "ftp_user"
	ApikeyFtpPassword      = "ftp_password"
	ApikeyDownloadToPath   = "download_to_path"
	ApikeyOpMode           = "op_mode"
	ApikeyBizNo            = "bizno"
	ApikeyRateId           = "rate_id"
	ApikeyIp               = "ip"
	ApikeyBizRefund        = "biz_refund"
	ApikeyOperatorId       = "operatorId"
	ApikeyBiz101           = "biz_101"
	ApikeyBiz102           = "biz_102"
	ApikeyBiz103           = "biz_103"
	ApikeyBiz201           = "biz_201"
	ApikeyBiz307           = "biz_307"
	ApikeyBiz308           = "biz_308"
	ApikeyBiz203           = "biz_203"
	ApikeyBiz205           = "biz_205"
	ApikeyOfflin           = "offline_agreement_path"
	ApiKeyCommno           = "recommend_no"
	ApikeyStoreId          = "storeId"
	ApikeyPubKey           = "pub_key"
	ApikeyPackageCode      = "package_code"
	ApikeyPassword         = "transaction_password"

	DB_STATIS = "statistics"
	DB_CRM    = "crm"
	DbStat    = "stat"
	//DB_RISK   = "risk"
	//DB_CRM_REPL = "chain_crm"

	DB_CRM_REPL = "crm_repl"

	DbTask     = "task"
	DbBill     = "crm"
	DbCallback = "paycb"
	DB_RISK    = "risk"
	DbMobile   = "mobile"

	// 业务类型
	InterfaceBizMerc        = 1
	InterfaceBizPay         = 2
	InterfaceBizSettle      = 3
	InterfaceBizAgentPay    = 4
	InterfaceBizAgentpayPub = 5

	PreInterfaceParam       = "iter_param"
	PreMercPoolNo           = "merc_pool_no"
	PreChannelStrategy      = "channel_strategy"
	PreAgentChannelStrategy = "agent_channel_strategy"
	PreMercKey              = "merc_key"
	PreAgencyKey            = "agency_key"
	PreCallbackURL          = "cb_url"
	PreChannelInterfaceKey  = "channel_inter"
	PreErrCode              = "err_code"
	PreMercCodeMercNo       = "merc_code_to_no"
	PreRelaApi              = "rela_api"
	PreGlobalParam          = "gl_param"
	PreRefundChkWallet      = "refund_chk_wallet"
	PreScoreUserKey         = "score_user_key"
	PreMercApiMode          = "merc_api_mode"
	PRE_CHANNEL_SUPPLIER    = "channel_supplier"
	PreRiskBwlist           = "r_bwlist"
	PreRiskPerorder         = "r_perOrder"
	PreRiskStat             = "r_stat"
	PreRiskStatInstance     = "r_istat"
	PRE_CHANNEL_PARAM       = "channel_param"
	PreSms                  = "sms"
	PrePicToken             = "pic_token"
	PreMail                 = "mail"

	CacheKeySec   = "600"
	CacheKeySecV2 = time.Second * 600
	SmsKeySec     = "45"
	SmsKeySecV2   = time.Second * 45
	PosNoKeySec   = "3600" // 一小时
	PosNoKeySecV2 = time.Minute * 15
	MailKeySec    = "45"
	MailKeySecV2  = time.Second * 45

	//
	IdxMercPubKey = 0
	IdxPlatPubKey = 1
	IdxPlatPriKey = 2
	IdxMercMD5Key = 3

	// global_param
	GlparamSettleDate       = "settle_date"
	GlparamDatecutLock      = "datecut_lock"
	GlparamDatecutTime      = "datecut_time"
	GlparamOperationNo      = "operation_no"
	GlparamNoSmsRegChk      = "no_sms_reg_chk"
	GlparamGxyRoleNo        = "gxy_role_no"
	GlparamYyRoleNo         = "yy_role_no"
	GlparamT1agentRoleNo    = "t1agent_role_no"
	GlparamScoerPayTypeList = "score_pay_type_list"

	// 进件信息
	AuditInfoBizLicense  = "1"  //"营业执照"
	AuditInfoAgentPic    = "2"  // "拓展人与门店负责人商户门头合照"
	AuditInfoBizPlace    = "3"  // "经营场景照"
	AuditInfoIdCardPic   = "4"  // 法人身份证正面照片
	AuditInfoIdCardPic1  = "5"  // 法人身份证反面照片
	AuditInfoTestVideo   = "6"  // 验证视频
	AuditInfoBankCardPic = "7"  // 结账卡
	AuditInfoOther       = "8"  // 其他
	AuditInfoStore       = "9"  // 门店
	AuditInfoCashDesk    = "10" // 收银台
	AuditInfoRoomInner   = "11" // 室内

	// api服务
	ServerNameApiMobile      = "go.micro.api.mobile"
	ServerNameApiPos         = "go.micro.api.pos"
	ServerNameApiWebadmin    = "go.micro.api.webadmin"
	ServerNameApiWebbusiness = "go.micro.api.webbusiness"
	ServerNameApiApiCb       = "go.micro.api.api-cb"
	ServerNameApiPay         = "go.micro.api.api"

	// 内部服务
	ServerNameBill           = "go.micro.srv.bill"
	ServerNameBusinessBill   = "go.micro.srv.business-bill"
	ServerNameBusinessSettle = "go.micro.srv.business-settle"
	ServerNamePayNotify      = "go.micro.srv.notify-srv"
	ServerNameAuth           = "go.micro.srv.auth"
	ServerNameCust           = "go.micro.srv.cust"
	ServerNameGis            = "go.micro.srv.gis"
	ServerNameListenExpKey   = "go.micro.srv.listen_exp_key"
	ServerNamePush           = "go.micro.srv.push"
	ServerNameTm             = "go.micro.srv.tm"
	ServerNameSettle         = "go.micro.srv.settle"
	ServerNameRiskctrl       = "go.micro.srv.riskctrl"
	ServerNameQuota          = "go.micro.srv.quota"
	ServerNameAdminAuth      = "go.micro.srv.admin-auth"

	BrokerTypeHTTP = "http"
	BrokerTypeNats = "nats"

	//秘钥对格式
	SecretKeyPKCS1 = "PKCS1"
	SecretKeyPKCS8 = "PKCS8"
)
