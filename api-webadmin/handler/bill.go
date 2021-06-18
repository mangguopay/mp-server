package handler

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_err"
	"context"

	"a.a/cu/container"
	"a.a/mp-server/api-webadmin/inner_util"
	"a.a/mp-server/api-webadmin/verify"
	billProto "a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/ss_net"
	"github.com/gin-gonic/gin"
)

type BillHandler struct {
	Client billProto.BillService
}

var BillHandlerInst BillHandler

func (b *BillHandler) InsertHeadquartersProfitWithdraw() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {

			req := &billProto.InsertHeadquartersProfitWithdrawRequest{
				AccountNo:    inner_util.GetJwtDataString(c, "account_uid"),
				CurrencyType: container.GetValFromMapMaybe(params, "currency_type").ToStringNoPoint(),
				Amount:       container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
				Note:         container.GetValFromMapMaybe(params, "note").ToStringNoPoint(),
			}

			// 参数校验
			if errStr := verify.CheckInsertHeadquartersProfitWithdrawVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := BillHandlerInst.Client.InsertHeadquartersProfitWithdraw(context.TODO(), req)
			return reply.ResultCode, 0, err
		})
	}
}

//
func (*BillHandler) AddCashRecharge() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := BillHandlerInst.Client.AddCashRecharge(context.TODO(), &billProto.AddCashRechargeRequest{
				AccAccount:   container.GetValFromMapMaybe(params, "acc_account").ToStringNoPoint(),
				Amount:       container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
				CurrencyType: container.GetValFromMapMaybe(params, "currency_type").ToStringNoPoint(),
				Notes:        container.GetValFromMapMaybe(params, "notes").ToStringNoPoint(),
				Password:     container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				NonStr:       container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),
				Uid:          inner_util.GetJwtDataString(c, "account_uid"),
			})
			//}, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

//
func (*BillHandler) UpdateBusinessToHeadStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := BillHandlerInst.Client.UpdateBusinessToHeadStatus(context.TODO(), &billProto.UpdateBusinessToHeadStatusRequest{
				LogNo:      container.GetValFromMapMaybe(params, "log_no").ToStringNoPoint(),
				Status:     container.GetValFromMapMaybe(params, "order_status").ToStringNoPoint(),
				Notes:      container.GetValFromMapMaybe(params, "notes").ToStringNoPoint(),
				AccountUid: inner_util.GetJwtDataString(c, "account_uid"),
			})
			//}, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

//
func (*BillHandler) UpdateToBusinessStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &billProto.UpdateToBusinessStatusRequest{
				LogNo:       container.GetValFromMapMaybe(params, "log_no").ToStringNoPoint(),
				OrderStatus: container.GetValFromMapMaybe(params, "order_status").ToStringNoPoint(),
				ImgStr:      container.GetValFromMapMaybe(params, "base64_img").ToStringNoPoint(),
				Notes:       container.GetValFromMapMaybe(params, "notes").ToStringNoPoint(),
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}

			if errStr := verify.UpdateToBusinessStatusVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := BillHandlerInst.Client.UpdateToBusinessStatus(context.TODO(), req, global.RequestTimeoutOptions)
			//}, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

//
func (*BillHandler) AddChangeBalanceOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &billProto.AddChangeBalanceOrderRequest{
				AccUid:       container.GetValFromMapMaybe(params, "acc_uid").ToString(),
				Amount:       container.GetValFromMapMaybe(params, "amount").ToString(),
				CurrencyType: container.GetValFromMapMaybe(params, "currency_type").ToString(), //
				OpType:       container.GetValFromMapMaybe(params, "op_type").ToString(),       //
				ChangeReason: container.GetValFromMapMaybe(params, "change_reason").ToString(), //
				NonStr:       container.GetValFromMapMaybe(params, "non_str").ToString(),       //
				LoginPwd:     container.GetValFromMapMaybe(params, "password").ToString(),      //
				LoginUid:     inner_util.GetJwtDataString(c, "account_uid"),                    //后台登陆账号的uid
				AccountType:  container.GetValFromMapMaybe(params, "account_type").ToString(),  //账号类型（修改的是什么身份的余额）
			}

			if errStr := verify.AddChangeBalanceOrderVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			var reply, err = BillHandlerInst.Client.AddChangeBalanceOrder(context.TODO(), req)
			return reply.ResultCode, 0, err
		})
	}
}
