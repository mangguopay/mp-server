package handler

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"fmt"
	"testing"
)

func DoInitDb() {
	l := []string{constants.DB_CRM, "risk"}
	names := []string{"mp_crm", "mp_risk"}
	for i := 0; i < len(l); i++ {
		host := "10.41.1.241"
		port := "5432"
		user := "postgres"
		password := "123"
		name := names[i]
		alias := strext.ToStringNoPoint(l[i])
		driver := "postgres"
		switch driver {
		case "postgres":
			db.DoDBInitPostgres(alias, host, port, user, password, name)
		default:
			fmt.Printf("not support database|driver=[%v]\n", driver)
		}
	}
}

func init() {
	DoInitDb()

}

var testSwitchJson = `
{
    "steps":[
        {
            "type":"switch",
            "steps":"133bb506-0be9-403c-be20-4f2507f6a5fb",
            "ops":{
                "default":[
                    {
                        "type":"pipeline",
                        "steps":"3512aab0-8647-4468-aaef-ca2dccdb6dc5"
                    }
                ]
            }
        }
    ]
}
`
var testIfJson = `
{
    "steps":[
        {
            "type":"if",
            "condition":"3512aab0-8647-4468-aaef-ca2dccdb6dc5",
            "rollback":"",
            "y":[

            ],
            "n":[

            ]
        }
    ]
}
`

func TestParsRule(t *testing.T) {
	//ops := unmarshalRuleStr(testIfJson)
	//result, threshold := ParsRule(ops)
	//t.Logf("结果为----->%s,阈值为----->%s", result, threshold)
	//req := &go_micro_srv_riskctrl.GetRiskCtrlResultRequest{
	//	ApiType:    "getLoginTimeInNSec",
	//	PayerAccNo: "e9586425-bfb7-4054-88b2-f1dfa47bdfa3",
	//	ActionTime: time.Now().String(),
	//	Amount:     "1000",
	//	Ip:         "127.0.0.1",
	//	PayType:    "余额",
	//	// 收款人账号
	//	PayeeAccNo:  "c0b497c9-8f5f-4903-9a27-098a2fc5af48",
	//	ProductType: "pos",
	//	// 币种
	//	MoneyType: "usd",
	//	// 订单号
	//	OrderNo: "2020011017504357966221",
	//}

	//result, threshold, riskNo, errStr := riskResult(req.ApiType, req.PayerAccNo, req.ActionTime, "", "", req.MoneyType, req.OrderNo)
	//if errStr != ss_err.ERR_SUCCESS {
	//	t.Errorf("------->%s", "!=success")
	//}
	//t.Logf("result-----> %s,threshold-----> %s,riskNo-----> %s", result, threshold, riskNo)
}
