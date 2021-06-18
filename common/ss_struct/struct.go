package ss_struct

// 头部公共信息结构体
type HeaderCommonData struct {
	Lang       string `json:"lang"`        // 语言类型
	AppVersion string `json:"app_version"` // app版本
	Timestamp  int64  `json:"timestamp"`   // 请求时间戳
	Platform   string `json:"platform"`    // 平台
	AppName    string `json:"app_name"`    // app名称(pos端或app端)
	UtmSource  string `json:"utm_source"`  // 来源
}

const (
	Platform_Ios     = "ios"
	Platform_Android = "android"

	// ios平台来源
	Utm_Source_App_Store_ = "app_store" // 苹果商店

	// android平台来源
	Utm_Source_Official_ = "official"          // 官方
	Utm_Source_Google_   = "google_play_store" // google商店

)

// app端-jwt内容
type JwtDataApp struct {
	Account        string
	AccountUid     string
	IdenNo         string
	AccountType    string
	LoginAccountNo string
	PubKey         string
	JumpIdenNo     string
	JumpIdenType   string
	MasterAccNo    string
	IsMasterAcc    string
	PosSn          string
}

// webAdmin端-jwt内容
type JwtDataWebAdmin struct {
	Account    string
	AccountUid string
	IdenNo     string
	//MerchantUid    string
	AccountType    string
	LoginAccountNo string
	JumpIdenNo     string
	JumpIdenType   string
	MasterAccNo    string
	IsMasterAcc    string
}

// webBusiness端-jwt内容
type JwtDataWebBusiness struct {
	Account        string
	AccountUid     string
	AccountType    string
	IdenNo         string
	Phone          string
	CountryCode    string
	LoginAccountNo string
	JumpIdenNo     string
	JumpIdenType   string
	MasterAccNo    string
	IsMasterAcc    string
	Email          string
}
