package ss_err

import "errors"

var (
	ErrWrongeCurrecyType = errors.New("参数币种类型错误")
	ErrLoginFailed       = errors.New("LoginFailed")
	ErrNoKeys            = errors.New("no keys")
	ErrBrokerIsNil       = errors.New("broker is nil")
	ErrArgsLen           = errors.New("参数个数不对")
)
