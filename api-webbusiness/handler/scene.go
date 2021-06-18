package handler

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-webbusiness/inner_util"
	"a.a/mp-server/common/constants"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_net"
	"context"
	"github.com/gin-gonic/gin"
)

func (*CustHandler) GetSceneSignedList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetSceneSignedListRequest{
				Page:          strext.ToInt32(params[0]),
				PageSize:      strext.ToInt32(params[1]),
				StartTime:     strext.ToStringNoPoint(params[2]),
				EndTime:       strext.ToStringNoPoint(params[3]),
				SceneName:     strext.ToStringNoPoint(params[4]),
				Status:        strext.ToStringNoPoint(params[5]),
				ApplyType:     strext.ToStringNoPoint(params[6]),
				SceneNo:       strext.ToStringNoPoint(params[7]),
				Lang:          ss_net.GetCommonData(c).Lang,
				BusinessAccNo: inner_util.GetJwtDataString(c, "account_uid"), //登陆账号的uid
			}

			if req.BusinessAccNo == "" {
				ss_log.Error("BusinessAccNo参数为空")
				return ss_err.ERR_ACCOUNT_NO_LOGIN, nil, 0, nil
			}

			if req.Lang == "" {
				req.Lang = constants.DefaultLang
			}

			reply, err := CustHandlerInst.Client.GetSceneSignedList(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用cust-srv.GetSceneSignedList()失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, gin.H{"data": reply.List, "show_add_button": reply.ShowAddButton}, reply.Total, nil
		}, "page", "page_size", "start_time", "end_time", "scene_name", "status", "apply_type", "scene_no")
	}
}
