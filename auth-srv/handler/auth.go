package handler

import (
	"context"
	"time"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/common/cache"
	auth "a.a/mp-server/common/proto/auth"
	"a.a/mp-server/common/ss_err"
)

type Auth struct{}

var AuthHandlerInst Auth

// 获取验证码
func (*Auth) GetCaptcha(ctx context.Context, req *auth.GetCaptchaRequet, resp *auth.GetCaptchaReply) error {
	strCode := util.RandomDigitStrOnlyAlphabetUpper(req.Strlen)
	_uuid := strext.NewUUIDNoSplit()

	err := cache.RedisClient.Set("verify_"+_uuid, strCode, time.Second*600).Err()
	ss_log.Debug("_uuid=[%v]|code=[%v],err: %v", _uuid, strCode, err)
	resp.ResultCode = ss_err.ERR_SUCCESS
	resp.Verifyid = _uuid

	resp.Base64Png = strCode
	return nil
}
