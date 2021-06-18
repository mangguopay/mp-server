package global

import (
	"github.com/micro/go-micro/v2/client"
	"time"
)

// service增加CallOption并设置超时时间
var RequestTimeoutOptions client.CallOption = func(o *client.CallOptions) {
	o.RequestTimeout = time.Second * 60
	o.DialTimeout = time.Second * 60
}
