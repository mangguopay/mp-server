package router

import (
	"a.a/mp-server/api-pos/common"
	"a.a/mp-server/api-pos/handler"
	mw "a.a/mp-server/api-pos/middleware"
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
	posUnauthorized(router)
	// 需要认证的组
	posAuthorized(router)
	return router
}

func posUnauthorized(router *gin.Engine) {
	unauthorized := router.Group("/pos/auth", mw.GetPriPubKeyMwInst.GetPubKeyRetMiddleWare())
	unauthorized.POST("/login", handler.AuthHandlerInst.MobileLogin())
	unauthorized.POST("/reg_sms", handler.AuthHandlerInst.RegSms())     //todo 文档已写
	unauthorized.POST("/back_pwd", handler.AuthHandlerInst.BackPWD())   // 找回密码 todo 文档已写
	unauthorized.POST("/check_sms", handler.AuthHandlerInst.CheckSms()) // 验证手机验证码 todo 文档已写
	unauthorized.POST("/download", handler.AuthHandlerInst.UnAuthDownloadImage())
	unauthorized.POST("/upload_client_info", handler.AuthHandlerInst.UploadClientInfo())                   // 上传客户端信息
	unauthorized.GET("/get_common_helps", handler.AuthHandlerInst.GetAppCommonHelps())                     // 获取帮助列表 todo 文档已写
	unauthorized.GET("/get_common_help", handler.AuthHandlerInst.GetCommonHelpDetail())                    // 获取帮助列表(单个) todo 文档已写
	unauthorized.GET("/get_consultation_config", handler.AuthHandlerInst.GetAppConsultationConfigDetail()) // 获取联系我们的信息（现是返回一个列表）
}

func posAuthorized(router *gin.Engine) {
	authorized := router.Group("/pos", mw.LoginOutMwInst.DoLoginOut(), mw.AuthMwInst.DoAuthJwt())
	authorized.GET("/userinfo", handler.AuthHandlerInst.Userinfo())                               // 获取userinfo信息 todo 文档已写
	authorized.GET("/version_info", handler.AuthHandlerInst.VersionInfo())                        // app 版本信息 todo 文档已写
	authorized.POST("/info/send_sms", handler.AuthHandlerInst.RegSms())                           // 没有手机号码的发送短信
	authorized.GET("/info/funcs", handler.AuthHandlerInst.GetFuncList())                          // todo 文档已写
	authorized.GET("/info/pos_remains", handler.AuthHandlerInst.GetPosRemain())                   // 获取 pos 余额 todo 文档已写
	authorized.GET("/info/servicer", handler.AuthHandlerInst.GetServicer())                       // 获取服务商 todo 文档已写
	authorized.POST("/info/recv_card", handler.AuthHandlerInst.AddRecvCard())                     // 添加收款人的银行卡 todo 文档已写
	authorized.GET("/info/headquarters_cards", handler.AuthHandlerInst.GetHeadquartersCards())    // 获取总部卡的列表 todo 文档已写
	authorized.GET("/info/servicer_cust_cards", handler.AuthHandlerInst.GetServicerOrCustCards()) //查询绑定的卡列表信息 todo 文档已写
	authorized.POST("/info/fence", handler.AuthHandlerInst.PosFence())                            // pos机围栏

	authorized.GET("/info/channels", handler.AuthHandlerInst.GetPosChannelList())             // 获取pos渠道列表 todo 文档已写
	authorized.POST("/info/add_card", handler.AuthHandlerInst.AddCard())                      // pos 新增银行卡 todo 文档已写
	authorized.POST("/info/modify_default_card", handler.AuthHandlerInst.ModifyDefaultCard()) // 修改默认卡 todo 文档已写
	authorized.POST("/info/delete_bind_card", handler.AuthHandlerInst.DeleteBindCard())       // 解除绑定卡 todo 文档已写
	authorized.GET("/info/my_data", handler.AuthHandlerInst.MyData())                         // pos机  我的资料 todo 文档已写
	authorized.GET("/info/servicer_pos", handler.AuthHandlerInst.GetSerPos())                 // 获取pos终端号
	authorized.POST("/info/modify_pos_status", handler.AuthHandlerInst.ModifySerPosStatus())  // 修改pos终端状态
	account := authorized.Group("/account")
	account.POST("/modify_pwd", handler.AuthHandlerInst.ModifyPWD())           // 修改登录密码 todo 文档已写
	account.POST("/modify_pay_pwd", handler.AuthHandlerInst.ModifyPayPWD())    // 修改支付密码 todo 文档已写
	account.POST("/modify_phone", handler.AuthHandlerInst.ModifyPhone())       // 修改手机号 todo 文档已写
	account.POST("/modify_nickname", handler.AuthHandlerInst.ModifyNickname()) //修改昵称 todo 文档已写
	account.POST("/perfecting_info", handler.AuthHandlerInst.PerfectingInfo()) //完善信息(修改头像，昵称)  todo 文档已写

	account.POST("/check_pay_pwd", handler.AuthHandlerInst.CheckPayPWD()) //验证支付密码 todo 文档已写

	authorized.GET("/bill/get_account_collect", handler.AuthHandlerInst.GetAccountCollect())      // 获取用户最近转账的人（账号、昵称、头像）  todo 文档已写
	authorized.POST("/bill/save_money", handler.AuthHandlerInst.SaveMoney())                      // pos 存款,线下现金存款 todo 文档已写
	authorized.POST("/bill/mobile_num_withdrawal", handler.AuthHandlerInst.MobileNumWithdrawal()) // 手机号取款 todo 文档已写
	authorized.GET("/bill/gen_withdraw_code", handler.AuthHandlerInst.GenWithdrawCode())          // pos机给用户的扫一扫取款吗,没金额的码 todo 文档已写
	// 此接口pos端没有使用到, xiaoyanchun注释 2020-07-13
	// authorized.POST("/bill/modify_code_status", handler.AuthHandlerInst.ModityGenCodeStatus())          // 修改用户扫一扫取款吗状态 todo 文档已写
	authorized.GET("/bill/query_code_status", handler.AuthHandlerInst.QuerySweepCodeStatus())    // 查看用户扫一扫的码的状态 todo 文档已写
	authorized.GET("/bill/sweep_withdraw_detail", handler.AuthHandlerInst.SweepWithdrawDetail()) // pos 被扫码取款后的订单详情(交易明细内的账单详情也是调这)
	authorized.GET("/bill/save_detail", handler.AuthHandlerInst.SaveMoneyDetail())               // pos 存款后的订单详情(交易明细内的账单详情也是调这)
	authorized.GET("/bill/change_balance_detail", handler.AuthHandlerInst.ChangeBalanceDetail()) // 平台修改服务商余额订单详情

	authorized.POST("/bill/confirm_withdraw", handler.AuthHandlerInst.ConfirmpWithdrawal())             // pos 确认取款 todo 文档已写
	authorized.POST("/bill/cancel_withdraw", handler.AuthHandlerInst.CancelWithdrawal())                // pos 取消取款 todo 文档已写
	authorized.POST("/bill/transfer_to_headquarters", handler.AuthHandlerInst.TransferToHeadquarters()) // 充值(转账到总部) todo 文档已写
	authorized.GET("/bill/get_to_headquarters", handler.AuthHandlerInst.GetTransferToHeadquartersLog()) // 获取转账到总部的转账记录 todo 文档已写
	authorized.GET("/bill/get_to_servicers", handler.AuthHandlerInst.GetTransferToServicerLogs())       // 获取转账到服务商的转账记录 todo 文档已写,但没数据
	authorized.POST("/bill/apply_money", handler.AuthHandlerInst.ApplyMoney())                          // 提现 todo 文档已写
	authorized.POST("/bill/upload", handler.AuthHandlerInst.UploadImage())                              // 上传图片 todo 文档已写

	// 前端没有用到此接口，xiaoyanchun 注释 2020-05-12
	//authorized.POST("/bill/download", handler.AuthHandlerInst.DownloadImage()) // 下载图片 todo 文档已写

	authorized.GET("/bill/servicer_bills", handler.AuthHandlerInst.ServicerBills())            //  pos端获取交易明细
	authorized.POST("/bill/query_rate", handler.AuthHandlerInst.QueryRate())                   // 查询和计算费率 todo 文档已写
	authorized.POST("/bill/query_min_max_amount", handler.AuthHandlerInst.QueryMinMaxAmount()) // 查询最小最大金额 todo 文档已写
	authorized.GET("/bill/query_save_receipt", handler.AuthHandlerInst.QuerySaveReceipt())     // 存款打印小票查询 todo 文档已写
	authorized.GET("/bill/withdraw_receipt", handler.AuthHandlerInst.WithdrawReceipt())        // pos机 手机号,扫一扫取款打印小票查询 todo 文档已写
	authorized.GET("/bill/save_withdraw_detail", handler.AuthHandlerInst.SaveWithdrawDetail()) // pos 端存款取款明细 （目前未知是否有调用该接口）

	authorized.GET("/servicer/servicer_collect_limit", handler.AuthHandlerInst.ServicerCollectLimit()) //pos获取服务商的收款额度 // todo 文档已写

	authorized.GET("/bill/get_servicer_billing_details", handler.AuthHandlerInst.GetServicerBillingDetails()) // POS获取服务商账单明细 todo 文档已写

	// 对账列表
	authorized.GET("/bill/get_servicer_check_list", handler.AuthHandlerInst.GetServicerCheckList()) // POS获取服务商对账单 todo 文档已写

	authorized.GET("/bill/get_servicer_profit_ledgers", handler.AuthHandlerInst.GetServicerProfitLedgers())     // POS获取服务商佣金统计 todo 文档已写
	authorized.GET("/bill/get_servicer_profit_ledger", handler.AuthHandlerInst.GetServicerProfitLedgerDetail()) // POS获取服务商佣金详情 todo 文档已写

	authorized.GET("/bill/real_time_count", handler.AuthHandlerInst.RealTimeCount()) // POS获取服务商当天统计（实时报表）  todo 文档已写

	authorized.POST("/cust/def_card", handler.AuthHandlerInst.UpdateDefCard()) // 修改用户默认卡

	authorized.POST("/cashier/add_cashier", handler.AuthHandlerInst.AddCashier())       // 添加店员
	authorized.GET("/cashier/cashiers", handler.AuthHandlerInst.GetCashiers())          // 获取店员列表
	authorized.GET("/cashier/cashier", handler.AuthHandlerInst.GetCashierDetail())      // 获取单个店员信息
	authorized.POST("/cashier/cashiers", handler.AuthHandlerInst.DeleteCashier())       // 删除店员
	authorized.POST("/cashier/modify_cashier", handler.AuthHandlerInst.ModifyCashier()) // 修改店员手机号

	authorized.GET("/refresh_token", handler.AuthHandlerInst.RefreshToken()) // 刷新token
}
