package ss_err

import (
	"testing"

	"a.a/mp-server/common/constants"
)

func TestGetPayApiErrMsg(t *testing.T) {

	retCode := ACErrSysRouteNotFound
	lang := constants.LangEnUS
	msg := GetPayApiErrMsg(retCode, lang)

	t.Logf("msg: %v", msg)
}
