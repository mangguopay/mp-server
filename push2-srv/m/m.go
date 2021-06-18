package m

type AdapterConfig struct {
}

type SendReq struct {
	Title    string
	Content  string
	Config   map[string]interface{}
	AccToken string // 与上游绑定的token
}

type SendResult struct {
}
