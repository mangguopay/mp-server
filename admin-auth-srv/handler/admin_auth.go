package handler

import (
	"context"
	"time"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/common/cache"
	adminAuthProto "a.a/mp-server/common/proto/admin-auth"
	"a.a/mp-server/common/ss_err"
)

type AdminAuth struct{}

// 获取验证码
func (a *AdminAuth) GetCaptcha(ctx context.Context, req *adminAuthProto.GetCaptchaRequet, resp *adminAuthProto.GetCaptchaReply) error {
	strCode := util.RandomDigitStrOnlyAlphabetUpper(req.Strlen)
	_uuid := strext.NewUUIDNoSplit()

	err := cache.RedisClient.Set("verify_"+_uuid, strCode, time.Second*600).Err()
	ss_log.Debug("_uuid=[%v]|code=[%v],err: %v", _uuid, strCode, err)
	resp.ResultCode = ss_err.ERR_SUCCESS
	resp.Verifyid = _uuid

	resp.Base64Png = strCode
	return nil
}
