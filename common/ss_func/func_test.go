package ss_func

import "testing"

func TestPreCountryCode(t *testing.T) {
	countryCode := "86"
	t.Logf("国家码：%v", PreCountryCode(countryCode))

}
