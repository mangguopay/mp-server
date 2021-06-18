package dao

import (
	_ "a.a/mp-server/auth-srv/test"
	"fmt"
	"testing"
)

func TestInset(t *testing.T) {
	err := CountryCodePhoneDaoInst.Insert(nil, "86", "10086")
	if err != nil {
		fmt.Println("=================", err.Error())
		return
	}
	//pq: duplicate key value violates unique constraint "country_code_phone_country_code_phone_key"

	fmt.Println("0------------")
}
func TestDelete(t *testing.T) {
	//err := CountryCodePhoneDaoInst.Delete(nil, "86", "10086")
	//if err != nil {
	//	fmt.Println("=================", err.Error())
	//	return
	//}
	fmt.Println("0------------", fmt.Sprintf("%04d", 855))
}
