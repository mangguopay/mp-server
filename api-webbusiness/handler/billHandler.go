package handler

import (
	"a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/mp-server/api-webbusiness/inner_util"
	"a.a/mp-server/api-webbusiness/verify"
	billProto "a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/ss_err"

	"a.a/mp-server/common/ss_net"
	"context"
	"github.com/gin-gonic/gin"
)

type BillHandler struct {
	Client billProto.BillService
}

var (
	BillHandlerInst BillHandler
)

/**
提现
*/
func (*BillHandler) AddBusinessWithdraw() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &billProto.BusinessWithdrawRequest{
				AccountUid:   inner_util.GetJwtDataString(c, "account_uid"),
				AccountType:  inner_util.GetJwtDataString(c, "account_type"),
				IdenNo:       inner_util.GetJwtDataString(c, "iden_no"),
				WithdrawType: container.GetValFromMapMaybe(params, "withdraw_type").ToStringNoPoint(), //1-普通提现;2-全部提现
				CardNo:       container.GetValFromMapMaybe(params, "card_no").ToStringNoPoint(),
				PayPwd:       container.GetValFromMapMaybe(params, "pay_pwd").ToStringNoPoint(),
				NonStr:       container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),
				Amount:       container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
				MoneyType:    container.GetValFromMapMaybe(params, "money_type").ToStringNoPoint(),
			}

			if errStr := verify.BusinessWithdrawVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			reply, err := BillHandlerInst.Client.BusinessWithdraw(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}
			return reply.ResultCode, gin.H{
				"log_no": reply.LogNo,
			}, nil
		})
	}
}
