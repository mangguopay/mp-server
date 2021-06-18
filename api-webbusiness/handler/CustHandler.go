package handler

import (
	"a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/api-webbusiness/common"
	"a.a/mp-server/api-webbusiness/dao"
	"a.a/mp-server/api-webbusiness/inner_util"
	"a.a/mp-server/common/ss_func"

	//	webUtil "a.a/mp-server/api-webbusiness/util"
	"context"
	"io/ioutil"
	"os"
	"strings"

	"a.a/mp-server/api-webbusiness/verify"
	"a.a/mp-server/common/aws_s3"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	authProto "a.a/mp-server/common/proto/auth"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_net"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/util/file"
	//"path"
)

type CustHandler struct {
	Client custProto.CustService
}

var (
	CustHandlerInst CustHandler
)

func (*CustHandler) GetBusinessAccountHome() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessAccountHome(context.TODO(), &custProto.GetBusinessAccountHomeRequest{
				//Page:        strext.ToString(params[0]),
				AccountUid: inner_util.GetJwtDataString(c, "account_uid"),
				IdenNo:     inner_util.GetJwtDataString(c, "iden_no"),
				AccType:    inner_util.GetJwtDataString(c, "account_type"),
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			return reply.ResultCode, gin.H{
				"business_data": reply.BusinessData, //商家详情
				"wallet_data":   reply.WalletData,   //商家钱包余额和冻结金额
				"sum_data":      reply.SumData,      //商家交易统计
			}, 0, nil
		}, "page")
	}

}

func (*CustHandler) GetBusinessBaseInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessBaseInfo(context.TODO(), &custProto.GetBusinessBaseInfoRequest{
				//Page:        strext.ToString(params[0]),
				AccountUid: inner_util.GetJwtDataString(c, "account_uid"),
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			return reply.ResultCode, reply.Data, 0, nil
		}, "page")
	}
}

func (*CustHandler) UpdateBusinessBaseInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.UpdateBusinessBaseInfoRequest{
				AccountUid:    inner_util.GetJwtDataString(c, "account_uid"),                     //登陆账号的uid
				MainIndustry:  container.GetValFromMapMaybe(params, "main_industry").ToString(),  //主要行业应用
				MainBusiness:  container.GetValFromMapMaybe(params, "main_business").ToString(),  //主营业务
				ContactPerson: container.GetValFromMapMaybe(params, "contact_person").ToString(), //联系人
				ContactPhone:  container.GetValFromMapMaybe(params, "contact_phone").ToString(),  //联系电话
				CountryCode:   container.GetValFromMapMaybe(params, "country_code").ToString(),   //联系电话的国家码
			}
			if errStr := verify.UpdateBusinessBaseInfoVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			if errStr := ss_func.CheckCountryCode(req.CountryCode); errStr != ss_err.ERR_SUCCESS {
				ss_log.Error("国家码[%v]不合法", req.CountryCode)
				return errStr, nil, nil
			}

			req.ContactPhone = ss_func.PrePhone(req.CountryCode, req.ContactPhone)

			var reply, err = CustHandlerInst.Client.UpdateBusinessBaseInfo(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}
			return reply.ResultCode, 0, nil
		})
	}
}

func (s *CustHandler) GetAuthMaterialDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			accountUid := inner_util.GetJwtDataString(c, "account_uid")
			accType := inner_util.GetJwtDataString(c, "account_type")

			if accountUid == "" {
				ss_log.Error("AccountUid参数为空")
				return ss_err.ERR_PARAM, nil, nil
			}

			switch accType { //根据账号类型 返回个人商家的认证材料，或企业商家认证材料
			case constants.AccountType_PersonalBusiness:
				req := &custProto.GetAuthMaterialBusinessDetailRequest{
					AccountUid: accountUid,
				}

				reply, err := CustHandlerInst.Client.GetAuthMaterialBusinessDetail(context.TODO(), req)
				if err != nil {
					ss_log.Error("api调用失败,err=[%v]", err)
					return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
				}

				return reply.ResultCode, gin.H{
					"data": reply.Data,
				}, err
			case constants.AccountType_EnterpriseBusiness:
				req := &custProto.GetAuthMaterialEnterpriseDetailRequest{
					AccountUid: accountUid,
				}

				reply, err := CustHandlerInst.Client.GetAuthMaterialEnterpriseDetail(context.TODO(), req)
				if err != nil {
					ss_log.Error("api调用失败,err=[%v]", err)
					return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
				}

				return reply.ResultCode, gin.H{
					"data": reply.Data,
				}, err
			default:
				ss_log.Error("账号类型不是个人商家、企业商家，accType[%v]", accType)
				return ss_err.ERR_PARAM, nil, nil
			}

			return ss_err.ERR_SYS_DB_GET, gin.H{}, nil
		}, "params")
	}
}

func (*CustHandler) AddAuthMaterialEnterprise() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.AddAuthMaterialEnterpriseRequest{
				ImgBase64:    container.GetValFromMapMaybe(params, "img_base64").ToString(),    //营业执照图片base64
				AuthName:     container.GetValFromMapMaybe(params, "auth_name").ToString(),     //公司名称
				AuthNumber:   container.GetValFromMapMaybe(params, "auth_number").ToString(),   //注册号/组织机构代码
				StartDate:    container.GetValFromMapMaybe(params, "start_date").ToString(),    //营业期限起始时间
				EndDate:      container.GetValFromMapMaybe(params, "end_date").ToString(),      //营业期限结束时间
				TermType:     container.GetValFromMapMaybe(params, "term_type").ToString(),     //营业期限类型（1-短期，2长期）
				SimplifyName: container.GetValFromMapMaybe(params, "simplify_name").ToString(), //简称
				//Addr:       container.GetValFromMapMaybe(params, "addr").ToString(),        //公司地址
				AccountUid: inner_util.GetJwtDataString(c, "account_uid"), //登陆账号的uid
			}

			if errStr := verify.AddAuthMaterialEnterpriseVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			req.AuthName = strings.Trim(req.AuthName, " ")
			req.AuthNumber = strings.Trim(req.AuthNumber, " ")
			req.SimplifyName = strings.Trim(req.SimplifyName, " ")

			var reply, err = CustHandlerInst.Client.AddAuthMaterialEnterprise(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}
			return reply.ResultCode, 0, nil
		})
	}
}

func (s *CustHandler) GetChannelList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessChannels(context.TODO(), &custProto.GetBusinessChannelsRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				CurrencyType: strext.ToStringNoPoint(params[2]),
				ChannelName:  strext.ToStringNoPoint(params[3]),
				UseStatus:    constants.Status_Enable, //只要启用的
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "currency_type", "channel_name")
	}
}

func (*CustHandler) GetCards() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessCards(context.TODO(), &custProto.GetBusinessCardsRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				BalanceType: strext.ToStringNoPoint(params[2]),
				AccountNo:   inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			return reply.ResultCode, reply.DataList, reply.Total, nil
		}, "page", "page_size", "money_type")
	}

}
func (*CustHandler) GetCardDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			reply, err := CustHandlerInst.Client.GetBusinessCardDetail(context.TODO(), &custProto.GetBusinessCardDetailRequest{
				AccountNo:   inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				CardNo:      strext.ToStringNoPoint(params),
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, gin.H{
				"data": reply.Data,
			}, nil
		}, "card_no")
	}

}

func (*CustHandler) AddCard() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			//验证支付密码
			replyCheckPwd, errCheckPwd := AuthHandlerInst.Client.CheckPayPWD(context.TODO(), &authProto.CheckPayPWDRequest{
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				IdenNo:      inner_util.GetJwtDataString(c, "iden_no"),
				Password:    container.GetValFromMapMaybe(params, "pay_password").ToStringNoPoint(),
				NonStr:      container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),
			})

			if errCheckPwd != nil {
				ss_log.Error("api调用失败,err=[%v]", errCheckPwd)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			if replyCheckPwd.ResultCode != ss_err.ERR_SUCCESS {
				return replyCheckPwd.ResultCode, nil, nil
			}

			req := &custProto.InsertOrUpdateBusinessCardRequest{
				CardNo:      container.GetValFromMapMaybe(params, "card_no").ToStringNoPoint(),
				AccountNo:   inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				ChannelId:   container.GetValFromMapMaybe(params, "channel_id").ToStringNoPoint(),
				Name:        container.GetValFromMapMaybe(params, "name").ToStringNoPoint(),
				CardNumber:  container.GetValFromMapMaybe(params, "card_number").ToStringNoPointReg(`^[\d]+$`),
				IsDefault:   container.GetValFromMapMaybe(params, "is_default").ToStringNoPoint(),
			}

			if errStr := verify.InsertOrUpdateBusinessCardVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			reply, err := CustHandlerInst.Client.InsertOrUpdateBusinessCard(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}
			return reply.ResultCode, 0, nil
		})
	}
}

func (s *CustHandler) DelCard() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DelBusinessCardRequest{
				CardNo: container.GetValFromMapMaybe(params, "card_no").ToStringNoPoint(),
			}
			if req.CardNo == "" {
				return ss_err.ERR_PARAM, nil, nil
			}

			reply, err := CustHandlerInst.Client.DelBusinessCard(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}

			return reply.ResultCode, 0, err
		})
	}
}

func (*CustHandler) GetMainIndustryCascaderDatas() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetMainIndustryCascaderDatas(context.TODO(), &custProto.GetMainIndustryCascaderDatasRequest{
				//Page:        strext.ToString(params[0]),
				Lang: ss_net.GetCommonData(c).Lang,
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			return reply.ResultCode, reply.Datas, 0, nil
		}, "page")
	}
}

func (*CustHandler) AddBusinessToHead() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			imgStr := container.GetValFromMapMaybe(params, "image_base64").ToStringNoPoint()

			if len(imgStr) > constants.UploadImgBase64LengthMax {
				return ss_err.ERR_ACCOUNT_IMAGE_BIG, nil, nil
			}

			req := &custProto.AddBusinessToHeadRequest{
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				IdenNo:      inner_util.GetJwtDataString(c, "iden_no"),
				CardNo:      container.GetValFromMapMaybe(params, "card_no").ToStringNoPoint(),
				Amount:      container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
				PayPwd:      container.GetValFromMapMaybe(params, "pay_pwd").ToStringNoPoint(),
				NonStr:      container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),
				MoneyType:   container.GetValFromMapMaybe(params, "money_type").ToStringNoPoint(),
				ImageBase64: imgStr,
			}

			if errStr := verify.AddBusinessToHeadVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			var reply, err = CustHandlerInst.Client.AddBusinessToHead(context.TODO(), req)
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

/**
 *
 */
func (s *CustHandler) IdentityVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			req := &custProto.IdentityVerifyRequest{
				Account:   container.GetValFromMapMaybe(params, "account").ToStringNoPoint(),
				Verifyid:  container.GetValFromMapMaybe(params, "verifyid").ToStringNoPoint(),
				Verifynum: container.GetValFromMapMaybe(params, "verifynum").ToStringNoPoint(),
			}

			if errStr := verify.IdentityVerifyVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.IdentityVerify(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			return strext.ToString(reply.ResultCode), gin.H{
				"phone": reply.Phone,
				"email": reply.Email,
			}, nil
		}, "params")
	}
}

func (*CustHandler) GetBusinessSceneList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessSceneListRequest{
				Page:      strext.ToStringNoPoint(params[0]),
				PageSize:  strext.ToStringNoPoint(params[1]),
				StartTime: strext.ToStringNoPoint(params[2]),
				EndTime:   strext.ToStringNoPoint(params[3]),
				SceneName: strext.ToStringNoPoint(params[4]),
				IsDelete:  "0", //商家中心只查询启用中的产品
				Lang:      ss_net.GetCommonData(c).Lang,
			}

			if req.Lang == "" {
				req.Lang = constants.DefaultLang
			}

			reply, err := CustHandlerInst.Client.GetBusinessSceneList(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			return reply.ResultCode, reply.Datas, reply.Total, nil
		}, "page", "page_size", "start_time", "end_time", "scene_name")
	}
}

func (*CustHandler) GetBusinessSceneDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessSceneDetailRequest{
				SceneNo: strext.ToStringNoPoint(params[0]),
				Lang:    ss_net.GetCommonData(c).Lang,
			}

			if req.Lang == "" {
				req.Lang = constants.DefaultLang
			}

			reply, err := CustHandlerInst.Client.GetBusinessSceneDetail(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Data, 0, nil
		}, "scene_no")
	}
}

func (*CustHandler) GetBusinessSignedList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessSignedListRequest{
				Page:      strext.ToStringNoPoint(params[0]),
				PageSize:  strext.ToStringNoPoint(params[1]),
				StartTime: strext.ToStringNoPoint(params[2]),
				EndTime:   strext.ToStringNoPoint(params[3]),
				//Status:     strext.ToStringNoPoint(params[4]),
				AppId:         strext.ToStringNoPoint(params[4]),
				AccountUid:    inner_util.GetJwtDataString(c, "account_uid"),
				IsBusinessReq: true, //是否是商家前端的请求（商家前端只显示通过和已过期的）
				Lang:          ss_net.GetCommonData(c).Lang,
			}

			if req.AppId == "" || req.AccountUid == "" {
				ss_log.Error("必要参数AppId、AccountUid为空.")
				return ss_err.ERR_PARAM, nil, 0, nil
			}

			if req.Lang == "" {
				req.Lang = constants.DefaultLang
			}

			reply, err := CustHandlerInst.Client.GetBusinessSignedList(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, nil
		}, "page", "page_size", "start_time", "end_time", "app_id")
	}
}

func (*CustHandler) GetBusinessAppList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessAppList(context.TODO(), &custProto.GetBusinessAppListRequest{
				Page:      strext.ToStringNoPoint(params[0]),
				PageSize:  strext.ToStringNoPoint(params[1]),
				StartTime: strext.ToStringNoPoint(params[2]),
				EndTime:   strext.ToStringNoPoint(params[3]),
				IdenNo:    inner_util.GetJwtDataString(c, "iden_no"),
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, nil
		}, "page", "page_size", "start_time", "end_time")
	}
}

func (*CustHandler) GetBusinessAppDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessAppDetailRequest{
				IdenNo: inner_util.GetJwtDataString(c, "iden_no"),
				AppId:  strext.ToStringNoPoint(params[0]),
			}

			if req.IdenNo == "" || req.AppId == "" {
				ss_log.Error("必要参数IdenNo、AppId为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetBusinessAppDetail(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Data, 0, nil
		}, "app_id")
	}
}

func (*CustHandler) InsertOrUpdateBusinessApp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			imgStr1 := container.GetValFromMapMaybe(params, "img_str1").ToStringNoPoint()
			imgStr2 := container.GetValFromMapMaybe(params, "img_str2").ToStringNoPoint()
			loginUid := inner_util.GetJwtDataString(c, "account_uid")

			if imgStr1 == "" || imgStr2 == "" || loginUid == "" {
				ss_log.Error("参数为空")
				return ss_err.ERR_PARAM, nil, nil
			}

			req := &custProto.InsertOrUpdateBusinessAppRequest{
				IdenNo:      inner_util.GetJwtDataString(c, "iden_no"),
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),                       //
				AppId:       container.GetValFromMapMaybe(params, "app_id").ToStringNoPoint(),     //应用id
				ApplyType:   container.GetValFromMapMaybe(params, "apply_type").ToStringNoPoint(), //应用类型
				AppName:     container.GetValFromMapMaybe(params, "app_name").ToStringNoPoint(),   //应用名称
				Describe:    container.GetValFromMapMaybe(params, "describe").ToStringNoPoint(),   //应用描述
				//BusinessPubKey: container.GetValFromMapMaybe(params, "business_pub_key").ToStringNoPoint(), //商家公钥
				//SignMethod:     container.GetValFromMapMaybe(params, "sign_method").ToStringNoPoint(),      //签名方式
				//IpWhiteList: container.GetValFromMapMaybe(params, "ip_white_list").ToStringNoPoint(), //ip白名单列表（使用逗号隔开）
			}

			imgStrs := []string{
				imgStr1,
				imgStr2,
			}

			var imgIds []string
			for k, v := range imgStrs {
				ss_log.Info("k=[%v]", k)
				upReply, errU := CustHandlerInst.Client.UploadImage(c, &custProto.UploadImageRequest{
					ImageStr:   v,
					AccountUid: loginUid,
					Type:       constants.UploadImage_UnAuth,
				})
				if errU != nil || upReply.ResultCode != ss_err.ERR_SUCCESS {
					ss_log.Error("k[%v],addErr=[%v]", k, errU)
					return ss_err.ERR_SAVE_IMAGE_FAILD, nil, nil
				}
				imgIds = append(imgIds, upReply.ImageId)
			}

			req.SmallImgNo = imgIds[0] //应用小图标
			req.BigImgNo = imgIds[1]   //应用大图标

			if errStr := verify.AddBusinessAppVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			reply, err := CustHandlerInst.Client.InsertOrUpdateBusinessApp(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}
			return reply.ResultCode, 0, nil
		})
	}
}

func (*CustHandler) DelBusinessApp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DelBusinessAppRequest{
				IdenNo: inner_util.GetJwtDataString(c, "iden_no"),
				AppId:  container.GetValFromMapMaybe(params, "app_id").ToStringNoPoint(),
			}

			if req.IdenNo == "" || req.AppId == "" {
				return ss_err.ERR_PARAM, 0, nil
			}

			var reply, err = CustHandlerInst.Client.DelBusinessApp(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}
			return reply.ResultCode, 0, nil
		})
	}
}

func (*CustHandler) BusinessUpdateAppStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.BusinessUpdateAppStatusRequest{
				IdenNo: inner_util.GetJwtDataString(c, "iden_no"),
				AppId:  container.GetValFromMapMaybe(params, "app_id").ToStringNoPoint(),
				Status: container.GetValFromMapMaybe(params, "status").ToStringNoPoint(),
			}

			if req.IdenNo == "" || req.AppId == "" {
				return ss_err.ERR_PARAM, 0, nil
			}

			var reply, err = CustHandlerInst.Client.BusinessUpdateAppStatus(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}
			return reply.ResultCode, 0, nil
		})
	}
}

/**
 * 获取总部卡列表
 */
func (s *AuthHandler) GetHeadCards() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetHeadquartersCardsRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				MoneyType:   strext.ToStringNoPoint(params[2]),
				AccountType: inner_util.GetJwtDataString(c, "account_type"), //7个人商家。8企业商家
			}
			if errStr := verify.GetHeadquartersCardsReqVerify(req); errStr != "" {
				return errStr, nil, 0, nil
			}

			if req.AccountType != constants.AccountType_EnterpriseBusiness && req.AccountType != constants.AccountType_PersonalBusiness {
				ss_log.Error("登录的账号类型[%v]错误,无权限调用接口", req.AccountType)
				return ss_err.ERR_PARAM, nil, 0, nil
			}
			reply, err := CustHandlerInst.Client.GetHeadquartersCards(c, req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "money_type")
	}
}

func (*CustHandler) GetBusinessToHeadList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessToHeadListRequest{
				Page:        strext.ToStringNoPoint(params[0]),
				PageSize:    strext.ToStringNoPoint(params[1]),
				StartTime:   strext.ToStringNoPoint(params[2]),
				EndTime:     strext.ToStringNoPoint(params[3]),
				MoneyType:   strext.ToStringNoPoint(params[4]),
				LogNo:       strext.ToStringNoPoint(params[5]),
				OrderStatus: strext.ToStringNoPoint(params[6]),
				IdenNo:      inner_util.GetJwtDataString(c, "iden_no"),
			}
			if req.IdenNo == "" {
				ss_log.Error("IdenNo参数为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetBusinessToHeadList(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, nil
		}, "page", "page_size", "start_time", "end_time", "money_type", "log_no", "order_status")
	}
}

func (*CustHandler) GetBusinessToHeadDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessToHeadDetailRequest{
				LogNo: strext.ToStringNoPoint(params[0]),
			}

			if req.LogNo == "" {
				ss_log.Error("LogNo为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetBusinessToHeadDetail(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Data, 0, nil
		}, "log_no")
	}
}

//提现列表
func (*CustHandler) GetToBusinessList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetToBusinessListRequest{
				Page:        strext.ToStringNoPoint(params[0]),
				PageSize:    strext.ToStringNoPoint(params[1]),
				StartTime:   strext.ToStringNoPoint(params[2]),
				EndTime:     strext.ToStringNoPoint(params[3]),
				MoneyType:   strext.ToStringNoPoint(params[4]),
				LogNo:       strext.ToStringNoPoint(params[5]),
				OrderStatus: strext.ToStringNoPoint(params[6]),
				IdenNo:      inner_util.GetJwtDataString(c, "iden_no"),
			}
			if req.IdenNo == "" {
				ss_log.Error("IdenNo参数为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetToBusinessList(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, nil
		}, "page", "page_size", "start_time", "end_time", "money_type", "log_no", "order_status")
	}
}

//提现详情
func (*CustHandler) GetBusinessToWithdrawDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessToWithdrawDetailRequest{
				LogNo: strext.ToStringNoPoint(params[0]),
			}

			if req.LogNo == "" {
				ss_log.Error("LogNo为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetBusinessToWithdrawDetail(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Data, 0, nil
		}, "log_no")
	}
}

//获取商家余额
func (*CustHandler) GetBusinessVAccBalance() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessVAccBalanceRequest{
				MoneyType:         strext.ToStringNoPoint(params[0]),
				BusinessAccountNo: inner_util.GetJwtDataString(c, "account_uid"),
			}
			if req.BusinessAccountNo == "" {
				ss_log.Error("account_uid参数为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetBusinessVAccBalance(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			data := gin.H{
				"balance":            reply.Balance,
				"frozen_balance":     reply.FrozenBalance,
				"recorded_amount":    reply.RecordedAmount,
				"expenditure_amount": reply.ExpenditureAmount,
				"money_type":         reply.MoneyType,
			}
			return reply.ResultCode, data, 0, nil
		}, "money_type")
	}
}

//商家虚账流水
func (*CustHandler) GetBusinessVAccLogList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessVAccLogListRequest{
				Page:          strext.ToStringNoPoint(params[0]),
				PageSize:      strext.ToStringNoPoint(params[1]),
				MoneyType:     strext.ToStringNoPoint(params[2]),
				BusinessAccNo: inner_util.GetJwtDataString(c, "account_uid"),
			}
			if req.BusinessAccNo == "" {
				ss_log.Error("IdenNo参数为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetBusinessVAccLogList(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.LogList, reply.Total, nil
		}, "page", "page_size", "money_type")
	}
}

//商家虚账流水(转账详情)
func (*CustHandler) GetBusinessVAccLogDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessVAccLogDetailRequest{
				LogNo:  strext.ToStringNoPoint(params[0]),
				Reason: strext.ToStringNoPoint(params[1]),
			}

			reply, err := CustHandlerInst.Client.GetBusinessVAccLogDetail(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Data, 0, nil
		}, "log_no", "reason")
	}
}

func (s *CustHandler) GetBulletins() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBulletins(context.TODO(), &custProto.GetBulletinsRequest{
				Page:        strext.ToStringNoPoint(params[0]),
				PageSize:    strext.ToStringNoPoint(params[1]),
				AccountType: inner_util.GetJwtDataString(c, "account_type"), //
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size")
	}
}

func (s *CustHandler) GetBulletinDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBulletinDetailRequest{
				BulletinId: strext.ToStringNoPoint(params[0]),
			}
			if req.BulletinId == "" {
				ss_log.Error("参数BulletinId为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetBulletinDetail(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			return reply.ResultCode, reply.Data, 0, err
		}, "bulletin_id")
	}
}

func (s *CustHandler) GetBusinessMessagesList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessMessagesList(context.TODO(), &custProto.GetBusinessMessagesListRequest{
				Page:        strext.ToStringNoPoint(params[0]),
				PageSize:    strext.ToStringNoPoint(params[1]),
				AccountType: inner_util.GetJwtDataString(c, "account_type"), //
				AccountNo:   inner_util.GetJwtDataString(c, "account_uid"),  //
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size")
	}
}

func (s *CustHandler) GetBusinessMessagesUnRead() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessMessagesUnRead(context.TODO(), &custProto.GetBusinessMessagesUnReadRequest{
				AccountType: inner_util.GetJwtDataString(c, "account_type"), //
				AccountNo:   inner_util.GetJwtDataString(c, "account_uid"),  //
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, gin.H{
				"unread_number": reply.Number,
			}, 0, err
		})
	}
}

func (*CustHandler) ReadAllBusinessMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.ReadAllBusinessMessagesRequest{
				AccountNo: inner_util.GetJwtDataString(c, "account_uid"),
			}

			if req.AccountNo == "" {
				return ss_err.ERR_PARAM, 0, nil
			}

			var reply, err = CustHandlerInst.Client.ReadAllBusinessMessages(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}
			return reply.ResultCode, 0, nil
		})
	}
}

func (*CustHandler) BusinessUpdatePartial() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.BusinessUpdatePartialRequest{
				IdenNo:      inner_util.GetJwtDataString(c, "iden_no"),
				AppId:       container.GetValFromMapMaybe(params, "app_id").ToStringNoPoint(),
				PubKey:      container.GetValFromMapMaybe(params, "business_pub_key").ToStringNoPoint(),
				SignMethod:  container.GetValFromMapMaybe(params, "sign_method").ToStringNoPoint(), //签名方式
				IpWhiteList: container.GetValFromMapMaybe(params, "ip_white_list").ToStringNoPoint(),
			}

			if req.IdenNo == "" || req.AppId == "" {
				return ss_err.ERR_PARAM, 0, nil
			}

			// 签名方式校验
			if errStr := ss_func.CheckSignMethod(req.SignMethod); errStr != ss_err.ERR_SUCCESS {
				return errStr, 0, nil
			}

			if len(req.PubKey) < 350 {
				ss_log.Error("商家公钥长度小于设定值350")
				return ss_err.ERR_BusinessPublicKeyLength_FAILD, 0, nil
			}

			var reply, err = CustHandlerInst.Client.BusinessUpdatePartial(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}
			return reply.ResultCode, 0, nil
		})
	}
}

func (*CustHandler) GenerateKeys() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			keyType := container.GetValFromMapMaybe(params, "key_type").ToStringNoPoint()
			if keyType == "" {
				keyType = constants.SecretKeyPKCS1
			}

			if !util.InSlice(keyType, []string{constants.SecretKeyPKCS1, constants.SecretKeyPKCS8}) {
				ss_log.Error("秘钥格式错误,keyType=[%v]", keyType)
				return ss_err.ERR_PARAM, 0, nil
			}

			req := &custProto.GenerateKeysRequest{
				KeyType: container.GetValFromMapMaybe(params, "key_type").ToStringNoPoint(),
			}

			var reply, err = CustHandlerInst.Client.GenerateKeys(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			if reply.ResultCode != ss_err.ERR_SUCCESS {
				return reply.ResultCode, nil, nil
			}

			keyFile := gin.H{
				"file_name":    reply.FileName,
				"file_content": reply.FileContent,
			}

			return ss_err.ERR_SUCCESS, keyFile, nil
		})
	}
}

func (s *CustHandler) UploadFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			accNo := inner_util.GetJwtDataString(c, "account_uid")
			accType := inner_util.GetJwtDataString(c, "account_type")
			baseName := container.GetValFromMapMaybe(params, "filename").ToStringNoPoint()
			filenameWithSuffix := strings.ToLower(container.GetValFromMapMaybe(params, "filename_with_suffix").ToStringNoPoint())
			uploadPath := container.GetValFromMapMaybe(params, "upload_path").ToStringNoPoint()

			ss_log.Error("获取到的数据 accNo=[%v],baseName=[%v],filenameWithSuffix=[%v]", accNo, baseName, filenameWithSuffix)
			if accNo == "" || baseName == "" || filenameWithSuffix == "" {
				ss_log.Error("必要参数为空为空 accNo=[%v],baseName=[%v],filenameWithSuffix=[%v]", accNo, baseName, filenameWithSuffix)
				return ss_err.ERR_PARAM, nil, nil
			}

			//确认一遍是否文件存在
			boolean, _ := file.Exists(uploadPath)
			if !boolean {
				ss_log.Error("上传的文件不存在,[%v]", uploadPath)
				return ss_err.ERR_FILE_OP_FAILD, nil, nil
			}

			fileName := ""
			fileType := ""
			//转存到该存的地方
			switch filenameWithSuffix {
			case ".xlsx":
				fileType = constants.UploadFileType_XLSX
				fileName = aws_s3.Xlsx_Dir + "/" + baseName
			default:
				ss_log.Error("文件后缀名[%v]错误", filenameWithSuffix)
				return ss_err.ERR_PARAM, nil, nil
			}

			_, errAwsS3 := common.UploadS3.UploadFile(uploadPath, fileName, true)
			if errAwsS3 != nil {
				ss_log.Error("上传到AwsS3失败，errAwsS3:[%v]", errAwsS3)
				return ss_err.ERR_UPLOAD, nil, nil
			}

			//删除临时文件
			if err := os.Remove(uploadPath); err != nil {
				ss_log.Error("删除临时文件失败,err[%v]", err)
			}

			//添加上传app版本日志
			fileLogId, err := dao.UploadFileLogDaoInstance.AddUploadFileLog(accNo, accType, fileName, fileType)
			if err != ss_err.ERR_SUCCESS {
				ss_log.Error("添加APP文件日志失败")
				return ss_err.ERR_PARAM, nil, nil
			}

			return ss_err.ERR_SUCCESS, gin.H{
				"file_id": fileLogId,
			}, nil
		}, "params")

	}
}

func (*CustHandler) DownloadBatchTransferBaseFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			//lang := container.GetValFromMapMaybe(params, "lang").ToStringNoPoint()

			//获取批量转账模板文件
			name := dao.GlobalParamDaoInstance.QeuryParamValue(constants.GlobalParamKeyBatchTransferBaseFile)
			fileName := aws_s3.Misc_Dir + "/" + name
			ss_log.Info("fileName:[%v]", fileName)

			//从s3获取文件
			result, s3Err := common.UploadS3.GetObject(fileName)
			if s3Err != nil {
				ss_log.Error("从s3获取文件失败,FileName:%s, err:%v", fileName, s3Err)
				return ss_err.ERR_FILE_OP_FAILD, nil, nil
			}

			// 读取body内容
			buff, err := ioutil.ReadAll(result.Body)
			if err != nil {
				ss_log.Error("err=[%v]\n", err)
				c.Set(ss_net.RET_CODE, ss_err.ERR_FILE_OP_FAILD)
				return ss_err.ERR_FILE_OP_FAILD, nil, nil
			}
			_, err2 := c.Writer.Write(buff)
			if err2 != nil {
				ss_log.Error("err=[%v]\n", err2)
				c.Set(ss_net.RET_CODE, ss_err.ERR_FILE_OP_FAILD)
				return ss_err.ERR_FILE_OP_FAILD, nil, nil
			}

			ss_log.Info("成功返回文件数据")

			c.Set(ss_net.RET_CODE, ss_err.ERR_SUCCESS)
			return ss_err.ERR_SUCCESS, nil, nil
		})

	}
}

func (s *CustHandler) GetImgUrl() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			imageId := container.GetValFromMapMaybe(params, "image_id").ToStringNoPoint()
			imgType := container.GetValFromMapMaybe(params, "img_type").ToStringNoPoint()
			if imageId == "" {
				ss_log.Error("参数ImageId为空")
				return ss_err.ERR_PARAM, nil, nil
			}
			switch imgType {
			case "1": // 需要授权

				req := &custProto.AuthDownloadImageRequest{ImageId: imageId}
				reply, err := CustHandlerInst.Client.AuthDownloadImage(context.TODO(), req, global.RequestTimeoutOptions)
				ss_log.Info("reply=[%v],err=[%v]", err)
				return reply.ResultCode, gin.H{
					"image_str": reply.ImageStr,
				}, nil

			case "2": // 不需要授权
				req := &custProto.UnAuthDownloadImageRequest{
					ImageId: imageId,
				}
				reply, err := CustHandlerInst.Client.UnAuthDownloadImage(context.TODO(), req)
				ss_log.Info("reply=[%v],err=[%v]", err)
				return reply.ResultCode, gin.H{
					"image_url": reply.ImageUrl,
				}, nil
			default:
				ss_log.Error("imgType参数错误")
				return ss_err.ERR_PARAM, nil, nil
			}

		}, "params")
	}
}

func (*CustHandler) AddBusinessSigned() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.AddBusinessSignedRequest{
				AccUid:     inner_util.GetJwtDataString(c, "account_uid"),
				BusinessNo: inner_util.GetJwtDataString(c, "iden_no"),
				AppId:      container.GetValFromMapMaybe(params, "app_id").ToStringNoPoint(),
				IndustryNo: container.GetValFromMapMaybe(params, "industry_no").ToStringNoPoint(),
				SceneNo:    container.GetValFromMapMaybe(params, "scene_no").ToStringNoPoint(),
			}

			if errStr := verify.AddBusinessSignedVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			var reply, err = CustHandlerInst.Client.AddBusinessSigned(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}
			return reply.ResultCode, nil, nil
		})
	}
}

func (*CustHandler) BusinessGetSceneList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.BusinessGetSceneListRequest{
				AppId:     strext.ToStringNoPoint(params[0]),
				ApplyType: strext.ToStringNoPoint(params[1]),
				Lang:      ss_net.GetCommonData(c).Lang,
			}

			if req.Lang == "" {
				req.Lang = constants.DefaultLang
			}

			reply, err := CustHandlerInst.Client.BusinessGetSceneList(context.TODO(), req)

			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, nil
		}, "app_id", "apply_type", "page_size", "start_time", "end_time", "scene_name")
	}
}

func (*CustHandler) UpdateAuthMaterialInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.UpdateBusinessAuthMaterialInfoRequest{
				SimplifyName: container.GetValFromMapMaybe(params, "simplify_name").ToString(), //简称
				AccountUid:   inner_util.GetJwtDataString(c, "account_uid"),                    //登陆账号的uid
				AccountType:  inner_util.GetJwtDataString(c, "account_type"),                   //登陆账号的角色
			}

			if errStr := verify.UpdateBusinessAuthMaterialInfoVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			req.SimplifyName = strings.Trim(req.SimplifyName, " ") //去掉前后空格

			var reply, err = CustHandlerInst.Client.UpdateBusinessAuthMaterialInfo(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}
			return reply.ResultCode, 0, nil
		})
	}
}

func (*CustHandler) AddBusinessSceneSigned() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.AddBusinessSceneSignedRequest{
				AccUid:     inner_util.GetJwtDataString(c, "account_uid"),
				BusinessNo: inner_util.GetJwtDataString(c, "iden_no"),
				IndustryNo: container.GetValFromMapMaybe(params, "industry_no").ToStringNoPoint(),
				SceneNo:    container.GetValFromMapMaybe(params, "scene_no").ToStringNoPoint(),
			}

			if errStr := verify.AddBusinessSceneSignedVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			reply, err := CustHandlerInst.Client.AddBusinessSceneSigned(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}
			return reply.ResultCode, nil, nil
		})
	}
}
