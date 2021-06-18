package constants

const (
	//全局参数key
	GlobalParamKeyAccPlat                   = "acc_plat"                           //平台总账号
	GlobalParamKeyAppSignedTerm             = "service_term"                       //商户签约期限(天)
	GlobalParamKeyPaymentPwdErrCount        = "err_payment_pwd_count"              //支付密码错误次数
	GlobalParamKeyLoginPwdErrCount          = "continuous_err_password"            //登录密码错误次数
	GlobalParamKeyBusinessTransferUSDAmount = "business_transfer_usd_amount_limit" //商家转账USD最小/最大转账金额限制
	GlobalParamKeyBusinessTransferKHRAmount = "business_transfer_khr_amount_limit" //商家转账KHR最小/最大转账金额限制
	GlobalParamKeyBusinessTransferUSDRate   = "business_transfer_usd_rate"         //商家转账USD手续费比率(万分比, 四舍五入)
	GlobalParamKeyBusinessTransferKHRRate   = "business_transfer_khr_rate"         //商家转账KHR手续费比率(万分比, 四舍五入)
	GlobalParamKeyBusinessTransferUSDMinFee = "business_transfer_usd_min_fee"      //商家转账最低收取USD手续费
	GlobalParamKeyBusinessTransferKHRMinFee = "business_transfer_khr_min_fee"      //商家转账最低收取KHR手续费
	GlobalParamKeyBusinessTransferBatchNum  = "business_transfer_batch_number"     //商家转账批量付款最大总人数(包括USD和KHR)

	GlobalParamKeyBatchTransferBaseFile = "batch_transfer_base_file" //批量转账模板文件

	//App支付签名key
	AppPaySignKey = "Y0zEyRTW4RIVMUHTvRsWVcvPoIv50dpL"

	TransferToUnRegisteredUser = "transfer_to_un_registered_user" //是否允许转账至未注册用户
)

const (
	AppFingerprintKey = "app_fingerprint_on" //是否开启用户指纹支付的录入参数key
)
