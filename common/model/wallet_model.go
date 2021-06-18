package model

type Wallet struct {
	WalletNo     string
	Remain       string
	IdenNo       string
	IdenType     int32
	WalletStatus string
	CreateTime   string
	WalletType   int64
	TypeParam    string
	OwedAmount   string
}

type UpdateWalletRequest struct {
	InnerOrderNo string
}

type UpdateWalletAgentRequest struct {
	InnerOrderNo string
}

type SettleLog struct {
	LogNo         string
	InnerOrderNo  string
	IdenNo        string
	Profit        string
	AccType       string
	RateUpper     string
	RateSelf      string
	Min           string
	Amount        string
	SettleStatus  string
	RateSplit     string
	ProfitSplited string
	ProfitMix     string
	MinUpper      string
	CreateTime    string
	ChannelNo     string
	PayTime       string
	IsReal        string
}
type BillFile struct {
	CreateTime   string
	OutReqNo     string
	InnerOrderNo string
	Amount       string
	ProductType  string
	OrderStatus  string
	MercCode     string
	MercName     string
	SettleDate   string
	TermInnerNo  string
	CardNo       string
	SettingNo    string
	TerNo        string
	MercNo       string
}
