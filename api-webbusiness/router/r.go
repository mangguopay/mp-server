package router

import (
	"a.a/mp-server/api-webbusiness/handler"
	mw "a.a/mp-server/api-webbusiness/middleware"
	"github.com/gin-gonic/gin"
)

var (
	respMw       mw.RespMiddleWare
	genTraceNoMw mw.GenTraceNoMw
)

type RouterHandler struct {
	HCust     handler.CustHandler
	HTransfer handler.BusinessTransferHandler
	HBill     handler.BillHandler
}

func InitRouter() *gin.Engine {
	gin.SetMode("release")
	router := gin.New()
	router.Use(gin.Logger(),
		mw.GenTraceNoMwInst.GenTraceNo(),    // 生成跟踪号
		mw.GetParamsMwInst.FetchGetParams(), // 读取get
		mw.GetParamsMwInst.FetchPostJsonBodyParams([]string{
			"/webbusiness/upload/upload_file",
		}), // 读取post/body-json

		// 逆序
		mw.RespMwInst.Resp(), // 返回封装
		mw.RespMwInst.Pack(), // json封装

		mw.RecoveryMiddleWareInst.Recovery(), //恢复
	)

	rh := RouterHandler{}

	// 不需要认证的组
	unauthorized := router.Group("/webbusiness/auth", mw.XPathVerifyMwInst.Verify(false))
	unauthorized.GET("/get_captcha", handler.AuthHandlerInst.GetCaptcha())                                   //获取验证码
	unauthorized.POST("/login", handler.AuthHandlerInst.Login())                                             //登录(目前是企业商家使用账号登录，个人商家的扫码登录还未开发)
	unauthorized.POST("/identity_verify", rh.HCust.IdentityVerify())                                         //忘记密码里的身份验证，返回脱敏的手机号和邮箱
	unauthorized.POST("/account/modify_loginpwd_sms", handler.AuthHandlerInst.NoTokenModifyPWD())            // 使用手机验证码方式修改登录密码(没有token情况下)
	unauthorized.POST("/account/modify_loginpwd_mail", handler.AuthHandlerInst.NoTokenModifyPWDByMailCode()) // 使用邮箱验证码方式修改登录密码(没有token情况下)

	unauthorized.GET("/account/check", handler.AccountHandlerInst.CheckAccount())    //检测账户是否存在
	unauthorized.POST("/account/add", handler.AuthHandlerInst.SaveBusinessAccount()) //注册
	unauthorized.POST("/account/send_mail", handler.AuthHandlerInst.SendMail())      //向邮箱发送验证码（无需jwt）
	unauthorized.POST("/account/send_sms", handler.AuthHandlerInst.SendSms())        //向手机号发送验证码（无需jwt）

	unauthorized.POST("/check/sms", handler.AuthHandlerInst.CheckSms())      // 验证手机验证码 是否正确（无需jwt的）
	unauthorized.GET("/check/email", handler.AccountHandlerInst.CheckMail()) // 确认email是否唯一

	// 需要认证的组
	authorized := router.Group("/webbusiness", mw.XPathVerifyMwInst.Verify(true), mw.JwtVerifyMwInst.VerifyToken(), mw.AuthMwInst.DoAuthJwt())
	authorized.POST("/auth/logout", handler.AuthHandlerInst.Logout())              //登出
	authorized.POST("/account/init_pay_pwd", handler.AuthHandlerInst.InitPayPwd()) //初次登录修改支付密码

	authorized.POST("/account/send_sms", handler.AuthHandlerInst.SendSms())   //向手机号发送验证码（需jwt）
	authorized.POST("/account/send_mail", handler.AuthHandlerInst.SendMail()) //向邮箱发送验证码（需jwt）

	authorized.POST("/check/sms", handler.AuthHandlerInst.CheckSms()) // 验证手机验证码 是否正确（需jwt的）

	//首页
	authorized.GET("/account/get_home", rh.HCust.GetBusinessAccountHome()) //获取商家首页信息

	//实名认证
	authorized.GET("/auth_material/get_detail", rh.HCust.GetAuthMaterialDetail()) //获取商家认证信息（个人或企业）
	authorized.POST("/auth_material/add", rh.HCust.AddAuthMaterialEnterprise())   //添加企业商家认证信息
	authorized.POST("/auth_material/update", rh.HCust.UpdateAuthMaterialInfo())   //修改商家认证信息（目前只有修改简称信息）

	//基本资料
	authorized.GET("/account/business_info", rh.HCust.GetBusinessBaseInfo())                //获取商家信息
	authorized.POST("/account/update_business_info", rh.HCust.UpdateBusinessBaseInfo())     //修改商家信息
	authorized.GET("/account/mail_industry_datas", rh.HCust.GetMainIndustryCascaderDatas()) //获取主要行业级联器所需数据

	//安全配置
	authorized.POST("/account/modify_login_pwd", handler.AuthHandlerInst.ModifyLoginPWD())          // 修改登录密码
	authorized.POST("/account/modify_loginpwd_sms", handler.AuthHandlerInst.ModifyPWD())            // 使用手机验证码方式修改登录密码(有token情况下)
	authorized.POST("/account/modify_loginpwd_mail", handler.AuthHandlerInst.ModifyPWDByMailCode()) // 使用邮箱验证码方式修改登录密码(有token情况下)

	authorized.POST("/account/modify_paypwd_sms", handler.AuthHandlerInst.ModifyPayPWD())             // 使用手机验证码方式修改支付密码
	authorized.POST("/account/modify_paypwd_mail", handler.AuthHandlerInst.ModifyPayPWDByMailCode())  // 使用邮箱验证码方式修改支付密码
	authorized.POST("/account/modify_paypwd_old_pay", handler.AuthHandlerInst.ModifyPayPWDByOldPwd()) // 使用旧支付密码方式修改支付密码

	authorized.POST("/account/modify_email", handler.AuthHandlerInst.ModifyEmail()) // 使用登录密码修改绑定邮箱
	authorized.POST("/account/modify_phone", handler.AuthHandlerInst.ModifyPhone()) // 修改绑定的手机号

	//银行卡管理
	authorized.GET("/card/channels", rh.HCust.GetChannelList()) //查询支持的渠道列表
	authorized.GET("/card/list", rh.HCust.GetCards())           //查询银行卡列表
	authorized.GET("/card/detail", rh.HCust.GetCardDetail())    //查询银行卡详情
	authorized.POST("/card/add", rh.HCust.AddCard())            //添加商家银行卡
	authorized.POST("/card/del", rh.HCust.DelCard())            //删除银行卡

	authorized.GET("/card/get_heads", handler.AuthHandlerInst.GetHeadCards()) // 获取总部卡的列表

	//资金管理
	authorized.GET("/apply_business/business/v_acc_log", rh.HCust.GetBusinessVAccLogList())          //商家账户流水
	authorized.GET("/apply_business/business/v_acc_log_detail", rh.HCust.GetBusinessVAccLogDetail()) //商家账户流水(转账详情)

	authorized.POST("/apply_business/business_to_head/add", rh.HCust.AddBusinessToHead())         //商家充值
	authorized.GET("/apply_business/business_to_head/list", rh.HCust.GetBusinessToHeadList())     //商家充值列表
	authorized.GET("/apply_business/business_to_head/detail", rh.HCust.GetBusinessToHeadDetail()) //商家充值详情

	authorized.POST("/apply_business/to_business/add", handler.BillHandlerInst.AddBusinessWithdraw()) //商家提现
	authorized.GET("/apply_business/to_business/list", rh.HCust.GetToBusinessList())                  //商家提现列表
	authorized.GET("/apply_business/to_business/detail", rh.HCust.GetBusinessToWithdrawDetail())      //商家提现详情
	authorized.GET("/apply_business/to_business/balance", rh.HCust.GetBusinessVAccBalance())          //商家账户金额

	authorized.POST("/apply_business/transfer/add", rh.HTransfer.AddBusinessTransfer())                    //商家转账
	authorized.GET("/apply_business/transfer/order_list", rh.HTransfer.GetBusinessTransferOrderList())     //商家转账订单列表
	authorized.GET("/apply_business/transfer/order_detail", rh.HTransfer.GetBusinessTransferOrderDetail()) //商家转账订单详情

	authorized.GET("/apply_business/transfer/batch_list", rh.HTransfer.GetBusinessTransferBatchList())     //商家转账批次列表
	authorized.GET("/apply_business/transfer/batch_detail", rh.HTransfer.GetBusinessTransferBatchDetail()) //商家转账批次详情

	authorized.POST("/download/batch_transfer_base_file", rh.HCust.DownloadBatchTransferBaseFile()) //下载批量转账模板文件

	//新增商家批量转账(付款)
	//1.上传文件
	authorized.POST("/upload/upload_file", rh.HCust.UploadFile()) //上传文件
	//2.分析信息、获取分析结果
	authorized.GET("/apply_business/batch_transfer/analysis_result", rh.HTransfer.GetBatchAnalysisResult()) //分析上传的批量转账文件，并返回结果
	//3.确认结果，并输入支付密码，根据结果进行转账
	authorized.POST("/apply_business/batch_transfer/confirm", rh.HTransfer.BatchConfirm()) //确认结果，开始执行转账

	//订单
	authorized.GET("/bill/get_bills", rh.HCust.GetBills())            //账单列表
	authorized.GET("/bill/get_bill_detail", rh.HCust.GetBillDetail()) //账单详情
	authorized.POST("/bill/download_bills", rh.HCust.DownloadBills()) //下载账单

	authorized.POST("/bill/refund", rh.HBill.BusinessBillRefund())           //退款
	authorized.GET("/bill/refund_bills", rh.HBill.GetRefundBills())          //退款账单列表
	authorized.GET("/bill/refund_detail", rh.HBill.GetRefundDetail())        //退款账单详情
	authorized.POST("/bill/download_refund", rh.HBill.DownloadRefundBills()) //下载账单

	//产品
	authorized.GET("/scene/list", rh.HCust.GetBusinessSceneList())     //获取产品列表
	authorized.GET("/scene/detail", rh.HCust.GetBusinessSceneDetail()) //获取产品详情

	//产品签约
	authorized.GET("/scene_signed/list", rh.HCust.GetSceneSignedList())     //获取商家产品签约列表
	authorized.POST("/scene_signed/add", rh.HCust.AddBusinessSceneSigned()) //添加产品签约

	//应用
	authorized.GET("/app/list", rh.HCust.GetBusinessAppList())                    //获取我的应用列表
	authorized.GET("/app/detail", rh.HCust.GetBusinessAppDetail())                //获取我的应用详情
	authorized.POST("/app/add_or_update", rh.HCust.InsertOrUpdateBusinessApp())   //添加或修改应用
	authorized.POST("/app/update_app_status", rh.HCust.BusinessUpdateAppStatus()) //应用的上下架
	authorized.POST("/app/del", rh.HCust.DelBusinessApp())                        //删除应用
	authorized.POST("/app/update_partial", rh.HCust.BusinessUpdatePartial())      //修改应用可修改部分的内容（审核通过后的应用）

	authorized.POST("/generate_keys", rh.HCust.GenerateKeys()) //生成秘钥对

	//公告
	authorized.GET("/bulletin/list", rh.HCust.GetBulletins())        //公告列表
	authorized.GET("/bulletin/detail", rh.HCust.GetBulletinDetail()) //公告详情

	//商家消息
	authorized.GET("/business_messages/list", rh.HCust.GetBusinessMessagesList())         //获取商家消息列表
	authorized.GET("/business_messages/unread_num", rh.HCust.GetBusinessMessagesUnRead()) //获取商家未读消息数量
	authorized.POST("/business_messages/read_all", rh.HCust.ReadAllBusinessMessages())    //修改商家未读消息为已读

	//查看图片
	authorized.GET("/img/get_img", rh.HCust.GetImgUrl()) //根据图片id获取url或base64字符串

	return router
}
