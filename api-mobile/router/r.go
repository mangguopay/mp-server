package router

import (
	"a.a/mp-server/api-mobile/common"
	"a.a/mp-server/api-mobile/handler"
	mw "a.a/mp-server/api-mobile/middleware"
	"a.a/mp-server/common/ss_err"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode("release")
	router := gin.New()
	router.Use(gin.Logger(),
		mw.GenTraceNoMwInst.GenTraceNo(), // 生成跟踪号

		// 解析头部公共信息
		mw.ParseCommonDataHeaderMwInst.ParseCommonDataHeader(),

		mw.GetParamsMwInst.FetchGetParams(),          // 读取get
		mw.GetParamsMwInst.FetchPostJsonBodyParams(), // 读取post/body-json
		mw.JwtVerifyMwInst.VerifyToken(),             // jwt验签
		// 获取公钥
		mw.GetPriPubKeyMwInst.GetPubKeyMiddleWare(),
		// 签名
		mw.RsaMwInst.DecodeRsa(), // 解密
		mw.RsaMwInst.Verify(),    // 验签
		// 换公钥
		mw.GetPriPubKeyMwInst.UpdatePubkeyWhenLogout(),
		// 围栏中间件
		mw.FenceMwInst.FenceMw(),

		mw.AuthMwInst.DoAuth(), // 验证

		// 逆序
		mw.RespMwInst.Resp(),  // 返回封装
		mw.RsaMwInst.Sign(),   // 签名
		mw.RsaMwInst.Encode(), // 加密
		mw.RespMwInst.Pack(),  // json封装

		mw.RecoveryMiddleWareInst.Recovery(), //恢复
	)

	router.NoRoute(func(c *gin.Context) {
		c.Set(common.RET_CODE, ss_err.ERR_SYS_NO_ROUTE)
		return
	})
	// 不需要认证的组
	unauthorized := router.Group("/mobile/auth", mw.GetPriPubKeyMwInst.GetPubKeyRetMiddleWare())
	unauthorized.POST("/login", handler.AuthHandlerInst.MobileLogin())
	unauthorized.POST("/reg_sms", handler.AuthHandlerInst.RegSms())     //todo 文档已写
	unauthorized.POST("/back_pwd", handler.AuthHandlerInst.BackPWD())   // 找回密码 todo 文档已写
	unauthorized.POST("/check_sms", handler.AuthHandlerInst.CheckSms()) // 验证手机验证码 todo 文档已写
	unauthorized.POST("/reg_login", handler.AuthHandlerInst.RegLogin()) // 注册和登录  todo 文档已写
	unauthorized.POST("/download", handler.AuthHandlerInst.UnAuthDownloadImage())
	unauthorized.GET("/get_common_helps", handler.AuthHandlerInst.GetAppCommonHelps())                     // 获取帮助列表 todo 文档已写
	unauthorized.GET("/get_common_help", handler.AuthHandlerInst.GetCommonHelpDetail())                    // 获取帮助列表(单个) todo 文档已写
	unauthorized.GET("/get_agreement", handler.AuthHandlerInst.GetAgreementAppDetail())                    // 获取协议内容
	unauthorized.GET("/get_consultation_config", handler.AuthHandlerInst.GetAppConsultationConfigDetail()) // 获取联系我们的信息（现是返回一个列表）
	unauthorized.POST("/upload_client_info", handler.AuthHandlerInst.UploadClientInfo())                   // 上传客户端信息
	unauthorized.POST("/log_app_dot", handler.AuthHandlerInst.InsertLogAppDot())                           // 添加app打点日志

	authorized := router.Group("/mobile", mw.LoginOutMwInst.DoLoginOut(), mw.AuthMwInst.DoAuthJwt())
	authorized.GET("/userinfo", handler.AuthHandlerInst.Userinfo())                               // 获取userinfo信息 todo 文档已写
	authorized.GET("/version_info", handler.AuthHandlerInst.VersionInfo())                        // app 版本信息 todo 文档已写
	authorized.POST("/info/send_sms", handler.AuthHandlerInst.RegSms())                           // 没有手机号码的发送短信
	authorized.GET("/info/funcs", handler.AuthHandlerInst.GetFuncList())                          // todo 文档已写
	authorized.GET("/info/remains", handler.AuthHandlerInst.GetRemain())                          // 获取余额 todo 文档已写
	authorized.GET("/info/pos_remains", handler.AuthHandlerInst.GetPosRemain())                   // 获取 pos 余额 todo 文档已写
	authorized.GET("/info/exchange_rate", handler.AuthHandlerInst.GetExchangeRate())              // 获取兑换汇率 todo 文档已写
	authorized.GET("/info/servicer", handler.AuthHandlerInst.GetServicer())                       // 获取服务商 todo 文档已写
	authorized.POST("/info/recv_card", handler.AuthHandlerInst.AddRecvCard())                     // 添加收款人的银行卡 todo 文档已写
	authorized.GET("/info/headquarters_cards", handler.AuthHandlerInst.GetHeadquartersCards())    // 获取总部卡的列表 todo 文档已写
	authorized.GET("/info/servicer_cust_cards", handler.AuthHandlerInst.GetServicerOrCustCards()) // 查询绑定的卡列表信息 todo 文档已写
	authorized.GET("/info/card_detail", handler.AuthHandlerInst.GetServicerOrCustCardDetail())    // 查询绑定的银行卡详细信息

	authorized.GET("/info/cust_payment", handler.AuthHandlerInst.CustPayment()) // 用户余额与银行卡列表
	authorized.POST("/info/fence", handler.AuthHandlerInst.PosFence())          // pos机围栏

	//authorized.GET("/info/modify_cards_is_defalut", handler.AuthHandlerInst.ModifyCardsDefalut()) 		//设置服务商或用户的卡为默认、推荐收款卡
	authorized.GET("/info/channels", handler.AuthHandlerInst.GetPosChannelList())     // 获取pos渠道列表 todo 文档已写
	authorized.GET("/info/use_channels", handler.AuthHandlerInst.GetUseChannelList()) // 获取用户可用渠道列表
	//authorized.POST("/info/add_card", handler.AuthHandlerInst.AddCard())                                    // pos 新增银行卡 todo 文档已写
	authorized.POST("/info/modify_default_card", handler.AuthHandlerInst.ModifyDefaultCard())               // 修改默认卡 todo 文档已写
	authorized.POST("/info/delete_bind_card", handler.AuthHandlerInst.DeleteBindCard())                     // 解除绑定卡 todo 文档已写
	authorized.GET("/info/my_data", handler.AuthHandlerInst.MyData())                                       // pos机  我的资料 todo 文档已写
	authorized.GET("/info/log_app_messages", handler.AuthHandlerInst.GetLogAppMessages())                   //查询app消息中心推送消息多个 todo 文档已写
	authorized.GET("/info/log_app_messages_cnt", handler.AuthHandlerInst.GetLogAppMessagesCnt())            //查询app未读的推送消息数量
	authorized.POST("/info/modify_app_messages_is_read", handler.AuthHandlerInst.ModifyAppMessagesIsRead()) // 批量修改推送消息的已读状态
	authorized.GET("/info/servicer_pos", handler.AuthHandlerInst.GetSerPos())                               // 获取pos终端号
	authorized.POST("/info/modify_pos_status", handler.AuthHandlerInst.ModifySerPosStatus())                // 修改pos终端状态

	account := authorized.Group("/account")
	account.POST("/modify_pwd", handler.AuthHandlerInst.ModifyPWD())           // 修改登录密码 todo 文档已写
	account.POST("/modify_pay_pwd", handler.AuthHandlerInst.ModifyPayPWD())    // 修改支付密码 todo 文档已写
	account.POST("/modify_phone", handler.AuthHandlerInst.ModifyPhone())       // 修改手机号 todo 文档已写
	account.POST("/modify_nickname", handler.AuthHandlerInst.ModifyNickname()) //修改昵称 todo 文档已写
	account.POST("/perfecting_info", handler.AuthHandlerInst.PerfectingInfo()) //完善信息(修改头像，昵称)  todo 文档已写

	account.POST("/auth_material", handler.AuthHandlerInst.AddAuthMaterialInfo()) //上传实名认证信息
	account.GET("/auth_material", handler.AuthHandlerInst.GetAuthMaterialInfo())  //用户查询实名认证信息

	account.POST("/check_pay_pwd", handler.AuthHandlerInst.CheckPayPWD()) //验证支付密码 todo 文档已写

	authorized.POST("/bill/exchange", handler.AuthHandlerInst.Exchange())                                        // 兑换  todo 文档已写
	authorized.POST("/bill/transfer", handler.AuthHandlerInst.Transfer())                                        // 转账 todo 文档已写
	authorized.GET("/bill/get_account_collect", handler.AuthHandlerInst.GetAccountCollect())                     // 获取用户最近转账的人（账号、昵称、头像）  todo 文档已写
	authorized.POST("/bill/gen_recv_code", handler.AuthHandlerInst.GenRecvCode())                                // 产生收款码 todo 文档已写
	authorized.GET("/bill/scan_recv_code", handler.AuthHandlerInst.ScanRecvCode())                               // 扫收款码获取信息,但还没发起 todo 文档已写
	authorized.POST("/bill/collection", handler.AuthHandlerInst.Collection())                                    // 收款 todo 文档已写
	authorized.GET("/bill/gen_withdraw_code", handler.AuthHandlerInst.GenWithdrawCode())                         // pos机给用户的扫一扫取款吗,没金额的码 todo 文档已写
	authorized.POST("/bill/modify_code_status", handler.AuthHandlerInst.ModityGenCodeStatus())                   // 修改用户扫一扫取款吗状态 todo 文档已写
	authorized.GET("/bill/query_code_status", handler.AuthHandlerInst.QuerySweepCodeStatus())                    // 查看用户扫一扫的码的状态 todo 文档已写
	authorized.POST("/bill/sweep_withdraw", handler.AuthHandlerInst.SweepWithdrawal())                           // 扫一扫取款 todo 文档已写
	authorized.GET("/bill/sweep_withdraw_detail", handler.AuthHandlerInst.SweepWithdrawDetail())                 // pos 被扫码取款后的订单详情
	authorized.GET("/bill/save_detail", handler.AuthHandlerInst.SaveMoneyDetail())                               // pos 存款后的订单详情
	authorized.POST("/bill/confirm_withdraw", handler.AuthHandlerInst.ConfirmpWithdrawal())                      // pos 确认取款 todo 文档已写
	authorized.POST("/bill/cancel_withdraw", handler.AuthHandlerInst.CancelWithdrawal())                         // pos 取消取款 todo 文档已写
	authorized.POST("/bill/transfer_to_headquarters", handler.AuthHandlerInst.TransferToHeadquarters())          // 充值(转账到总部) todo 文档已写
	authorized.POST("/bill/cust_transfer_to_headquarters", handler.AuthHandlerInst.CustTransferToHeadquarters()) // app端 客户充值(转账到总部) todo 文档已写
	authorized.GET("/bill/get_to_headquarters", handler.AuthHandlerInst.GetTransferToHeadquartersLog())          // 获取转账到总部的转账记录 todo 文档已写
	authorized.GET("/bill/get_to_servicers", handler.AuthHandlerInst.GetTransferToServicerLogs())                // 获取转账到服务商的转账记录 todo 文档已写,但没数据
	authorized.POST("/bill/apply_money", handler.AuthHandlerInst.ApplyMoney())                                   // 提现 todo 文档已写
	authorized.POST("/bill/upload", handler.AuthHandlerInst.UploadImage())                                       // 上传图片 todo 文档已写

	// 前端没有用到此接口，xiaoyanchun 注释 2020-05-12
	//authorized.POST("/bill/download", handler.AuthHandlerInst.DownloadImage()) // 下载图片 todo 文档已写

	authorized.GET("/bill/cust_bills", handler.AuthHandlerInst.CustBills())                          // app端获取用户账单 todo 文档已写
	authorized.GET("/bill/cust_order_bill", handler.AuthHandlerInst.CustOrderBillDetail())           // app端获取用户账单详情查询 todo 文档已写
	authorized.POST("/bill/cust_withdraw", handler.AuthHandlerInst.CustWithdraw())                   // app端,客户直接从总部提现
	authorized.GET("/bill/to_custs", handler.AuthHandlerInst.GetLogToCusts())                        // 用户获取银行卡提现账单
	authorized.GET("/bill/cust_to_headquarters", handler.AuthHandlerInst.GetLogCustToHeadquarters()) // 用户获取银行卡充值账单

	authorized.GET("/bill/cust_outgo_bills", handler.AuthHandlerInst.CustOutgoBills())   // 用户获取网点提现账单（一个账单只有一条，不是CustBills那样账单和手续费分开）
	authorized.GET("/bill/cust_income_bills", handler.AuthHandlerInst.CustIncomeBills()) // 用户获取网点充值账单（一个账单只有一条，不是CustBills那样账单和手续费分开）

	//authorized.GET("/bill/cust_income_bill", handler.AuthHandlerInst.CustIncomeBillsDetail())  // app端获取用户充值账单查询
	//authorized.GET("/bill/cust_outgo_bill", handler.AuthHandlerInst.CustOutgoBillsDetail())             // app端获取用户提现账单查询
	//authorized.GET("/bill/cust_transfer_bill", handler.AuthHandlerInst.CustTransferBillsDetail())       // app端获取用户转账账单查询
	//authorized.GET("/bill/cust_collection_bill", handler.AuthHandlerInst.CustCollectionBillsDetail())   // app端获取用户收款账单查询

	authorized.GET("/bill/servicer_bills", handler.AuthHandlerInst.ServicerBills())            //todo   pos端获取交易明细（只返回存款、取款） todo 文档已写
	authorized.POST("/bill/query_rate", handler.AuthHandlerInst.QueryRate())                   // 查询和计算费率 todo 文档已写
	authorized.POST("/bill/cust_query_rate", handler.AuthHandlerInst.CustQueryRate())          // 客户跟平台之间的操作查询和计算费率 todo 文档已写
	authorized.POST("/bill/query_min_max_amount", handler.AuthHandlerInst.QueryMinMaxAmount()) // 查询最小最大金额 todo 文档已写
	authorized.GET("/bill/query_save_receipt", handler.AuthHandlerInst.QuerySaveReceipt())     // 存款打印小票查询 todo 文档已写
	authorized.GET("/bill/withdraw_receipt", handler.AuthHandlerInst.WithdrawReceipt())        // pos机 手机号,扫一扫取款打印小票查询 todo 文档已写
	authorized.GET("/bill/save_withdraw_detail", handler.AuthHandlerInst.SaveWithdrawDetail()) // pos 端存款取款明细 todo 文档已写
	authorized.GET("/bill/exchange_amount", handler.AuthHandlerInst.ExchangeAmount())          // app 兑换金额 todo 文档已写

	authorized.GET("/servicer/servicer_collect_limit", handler.AuthHandlerInst.ServicerCollectLimit()) //pos获取服务商的收款额度 // todo 文档已写

	authorized.GET("/bill/get_servicer_billing_details", handler.AuthHandlerInst.GetServicerBillingDetails()) // POS获取服务商账单明细 todo 文档已写

	// 对账列表
	authorized.GET("/bill/get_servicer_check_list", handler.AuthHandlerInst.GetServicerCheckList()) // POS获取服务商对账单 todo 文档已写

	authorized.GET("/bill/get_servicer_profit_ledgers", handler.AuthHandlerInst.GetServicerProfitLedgers())     // POS获取服务商佣金统计 todo 文档已写
	authorized.GET("/bill/get_servicer_profit_ledger", handler.AuthHandlerInst.GetServicerProfitLedgerDetail()) // POS获取服务商佣金详情 todo 文档已写

	authorized.GET("/bill/real_time_count", handler.AuthHandlerInst.RealTimeCount()) // POS获取服务商当天统计（实时报表）  todo 文档已写

	authorized.POST("/cust/def_card", handler.AuthHandlerInst.UpdateDefCard()) // 修改用户默认卡

	authorized.GET("/get_nearby_servicers", handler.AuthHandlerInst.GetNearbyServicerList()) // 获取附近服务商列表

	authorized.GET("/refresh_token", handler.AuthHandlerInst.RefreshToken()) // 刷新token

	// 商家扫码用户付款二维码支付
	authorized.GET("/pay/payment_code", handler.BusinessBillHandlerInst.QueryPaymentCode())            //获取付款码
	authorized.GET("/pay/query_pending_order", handler.BusinessBillHandlerInst.QueryPendingPayOrder()) //查询用户待支付订单
	authorized.POST("/pay/order_pay", handler.BusinessBillHandlerInst.OrderPay())                      //用户订单支付

	// 用户扫码商家带金额的二维码支付
	authorized.GET("/pay/get_trans_info", handler.BusinessBillHandlerInst.QueryOrderTransInfo()) //获取订单交易信息(金额, 订单标题, 商家名称)
	authorized.POST("/pay/qr_amount_pay", handler.BusinessBillHandlerInst.QrCodeAmountPay())     // 扫码支付-扫临时二维码

	// 用户扫码商家固定二维码支付
	authorized.GET("/pay/get_business_app_info", handler.BusinessBillHandlerInst.GetBusinessAppInfo()) // 获取商家的应用信息(用户扫码付款时)
	authorized.POST("/pay/qr_fixed_prepay", handler.BusinessBillHandlerInst.QrCodeFixedPrePay())       //用户APP下单
	authorized.POST("/pay/qr_fixed_pay", handler.BusinessBillHandlerInst.QrCodeFixedPay())             // 扫码支付-扫固定二维码

	//app支付
	authorized.POST("/pay/app_pay", handler.BusinessBillHandlerInst.AppPay()) // app支付

	//扫个人商家支付
	authorized.GET("/pay/get_business_info", handler.BusinessBillHandlerInst.GetPersonalBusinessInfo())  // 获取个人商家的信息
	authorized.POST("/pay/fixed_code_prepay", handler.BusinessBillHandlerInst.PersonalFixedCodePrePay()) //用户扫个人商家固码下单

	//添加指纹支付
	authorized.POST("/fingerprint/add", handler.AuthHandlerInst.AddFingerprint())     // 添加指纹支付
	authorized.POST("/fingerprint/close", handler.AuthHandlerInst.CloseFingerprint()) // 关闭指纹支付

	//商家服务
	merchant := router.Group("/mobile/merchant", mw.AuthMwInst.DoAuthJwt())
	merchant.GET("/industry/get_list", handler.AuthHandlerInst.GetMainIndustryCascaderDatas())            //获取主要行业级联器所需数据
	merchant.POST("/individual_business_auth/add", handler.AuthHandlerInst.AddAuthMaterialBusiness())     //添加个人商家认证信息
	merchant.GET("/individual_business_auth/get_detail", handler.AuthHandlerInst.GetAuthMaterialDetail()) //获取个人商家认证信息

	merchant.GET("/today_trading", handler.BusinessInst.GetPersonalBusinessBalance())    //个人商家今日交易统计
	merchant.GET("/bills", handler.BusinessInst.GetPersonalBusinessBills())              //个人商家订单列表
	merchant.GET("/bill_detail", handler.BusinessInst.GetPersonalBusinessBillDetail())   //个人商家订单详情
	merchant.GET("/pay/fixed_code", handler.BusinessInst.GetPersonalBusinessFixedCode()) //获取个人商家固定收款码
	merchant.POST("/pay/personal_prepay", handler.BusinessInst.PersonalBusinessPerPay()) //个人商家设置金额下单

	merchant.GET("/personal_info", handler.BusinessInst.GetPersonalBusinessInfo())

	merchant.GET("/card/channels", handler.BusinessInst.GetChannelList())      //查询支持的渠道列表
	merchant.GET("/card/list", handler.BusinessInst.GetBusinessCards())        //查询个人商家银行卡列表
	merchant.GET("/card/detail", handler.BusinessInst.GetBusinessCardDetail()) //查询个人商家银行卡详情
	merchant.POST("/card/add", handler.BusinessInst.AddBusinessCard())         //添加个人商家银行卡
	merchant.POST("/card/del", handler.BusinessInst.DelBusinessCard())         //删除个人商家银行卡

	merchant.GET("/card/get_heads", handler.BusinessInst.GetHeadCards())                           // 获取总部卡的列表
	merchant.POST("/apply_business/to_head/add", handler.BusinessInst.AddBusinessToHead())         //个人商家充值
	merchant.GET("/apply_business/to_head/list", handler.BusinessInst.GetBusinessToHeadList())     //商家充值列表
	merchant.GET("/apply_business/to_head/detail", handler.BusinessInst.GetBusinessToHeadDetail()) //商家充值详情

	return router
}
