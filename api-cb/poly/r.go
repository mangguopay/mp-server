package poly

import (
	"a.a/mp-server/api-cb/poly/supplier/fakebank"
	"a.a/mp-server/api-cb/poly/supplier/p66"
	"a.a/mp-server/common/constants"
)

func getTargetApi(supplierCode string) Ipoly {
	switch supplierCode {
	case constants.SupplierCodeFakebank:
		return fakebank.PolyFakebankInst
	case constants.SupplierCodeP66:
		return p66.PolyP66Inst
	default:
		// 没有找到对应供应商
	}
	return nil
}
