package router

import (
	"a.a/mp-server/api-webadmin/handler"
	mw "a.a/mp-server/api-webadmin/middleware"
	"github.com/gin-gonic/gin"
)

var (
	respMw       mw.RespMiddleWare
	genTraceNoMw mw.GenTraceNoMw
)

func InitRouter() *gin.Engine {
	gin.SetMode("release")
	router := gin.New()
	router.Use(gin.Logger(),
		mw.GenTraceNoMwInst.GenTraceNo(),    // 生成跟踪号
		mw.GetParamsMwInst.FetchGetParams(), // 读取get
		mw.GetParamsMwInst.FetchPostJsonBodyParams([]string{
			//"/webadmin/app_version/app_version",
			"/webadmin/upload/upload_file",
		}, []string{
			//"/webadmin/app_version/app_version",
			"/webadmin/upload/upload_file",
		}), // 读取post/body-json

		// 逆序
		mw.RespMwInst.Resp(), // 返回封装
		mw.RespMwInst.Pack(), // json封装

		mw.RecoveryMiddleWareInst.Recovery(), //恢复
	)

	// 不需要认证的组
	unauthorized := router.Group("/webadmin/auth", mw.XPathVerifyMwInst.Verify(false))
	unauthorized.POST("/login", handler.AdminAuthHandlerInst.Login())
	unauthorized.GET("/get_captcha", handler.AdminAuthHandlerInst.GetCaptcha())

	// 需要认证的组
	authorized := router.Group("/webadmin", mw.XPathVerifyMwInst.Verify(true), mw.JwtVerifyMwInst.VerifyToken(), mw.AuthMwInst.DoAuthJwt())
	authorized.POST("/auth/logout", handler.AdminAuthHandlerInst.Logout())

	authorized.GET("/account/accounts", handler.AdminAuthHandlerInst.GetAccountList())
	authorized.GET("/account/accounts2", handler.AdminAuthHandlerInst.GetAccountList2())                //获取用户列表 、获取用户账户明细
	authorized.GET("/account/un_actived_accounts", handler.AdminAuthHandlerInst.GetUnActivedAccounts()) //获取未激活用户列表
	authorized.GET("/account/account", handler.AdminAuthHandlerInst.GetAccount())
	authorized.POST("/account/account", handler.AdminAuthHandlerInst.SaveAccount())
	authorized.GET("/account/account_by_nickname", handler.AdminAuthHandlerInst.GetAccountByNickname())
	authorized.PATCH("/account/account_auth", handler.AdminAuthHandlerInst.UpdateOrInsertAccountAuth())
	authorized.DELETE("/account/account", handler.AdminAuthHandlerInst.DeleteAccount()) //现在不允许删除账号，此接口将直接返回错误
	authorized.POST("/account/modify_password", handler.AdminAuthHandlerInst.ModifyPw())
	authorized.POST("/account/reset_password", handler.AdminAuthHandlerInst.ResetPw())
	//authorized.GET("/log/login", handler.AccountHandlerInst.GetLogLoginList())
	authorized.GET("/account/role_from_acc", handler.AdminAuthHandlerInst.GetRoleFromAcc())
	authorized.POST("/account/modify_user_status", handler.AdminAuthHandlerInst.ModifyUserStatus())
	authorized.GET("/log/get_log_accounts", handler.CustHandlerInst.GetLogAccounts()) //获取操作日志（WEB、app的）

	authorized.POST("/account/add_cashier", handler.AuthHandlerInst.AddCashier())       //添加店员
	authorized.GET("/cashier/cashiers", handler.CustHandlerInst.GetCashiers())          //获取店员列表信息
	authorized.GET("/cashier/cashier", handler.CustHandlerInst.GetCashierDetail())      //获取单个店员信息
	authorized.DELETE("/cashier/cashier", handler.CustHandlerInst.DeleteCashier())      //删除店员
	authorized.POST("/cashier/modify_cashier", handler.CustHandlerInst.ModifyCashier()) //修改店员 (目前该接口是错误的设计，后台只允许帮服务商创建和删除店员，不允许更改店员的手机号)

	//检测账户是否存在
	authorized.GET("/account/check", handler.AdminAuthHandlerInst.CheckAccount())

	//菜单管理
	authorized.GET("/menu/menus", handler.AdminAuthHandlerInst.GetMenuList())
	authorized.GET("/menu/menu", handler.AdminAuthHandlerInst.GetMenu())
	authorized.POST("/menu/menu", handler.AdminAuthHandlerInst.SaveOrInsertMenu())
	authorized.PATCH("/menu/menu", handler.AdminAuthHandlerInst.SaveOrInsertMenu())
	authorized.DELETE("/menu/menu", handler.AdminAuthHandlerInst.DeleteMenu())
	authorized.POST("/menu/refresh_child", handler.AdminAuthHandlerInst.MenuRefreshChild())

	//角色相关
	authorized.GET("/role/roles", handler.AdminAuthHandlerInst.GetRoleList())              // 角色列表
	authorized.GET("/role/role_info", handler.AdminAuthHandlerInst.GetRoleInfo())          // 获取角色信息
	authorized.GET("/role/role_auth_table", handler.AdminAuthHandlerInst.GetRoleUrlList()) // 获取角色菜单列表
	authorized.GET("/role/role", handler.AdminAuthHandlerInst.GetRole())
	authorized.POST("/role/role", handler.AdminAuthHandlerInst.UpdateOrInsertRole())
	authorized.PATCH("/role/role_auth", handler.AdminAuthHandlerInst.UpdateOrInsertRoleAuth()) // 授权给角色
	authorized.DELETE("/role/role", handler.AdminAuthHandlerInst.DeleteRole())
	authorized.POST("/role/role_def", handler.AdminAuthHandlerInst.AuthRole())

	//管理员菜单管理
	authorized.GET("/admin_menu/menus", handler.AdminAuthHandlerInst.GetAdminMenuList())
	authorized.GET("/admin_menu/menu", handler.AdminAuthHandlerInst.GetAdminMenu())
	authorized.POST("/admin_menu/menu", handler.AdminAuthHandlerInst.SaveOrInsertAdminMenu())
	authorized.PATCH("/admin_menu/menu", handler.AdminAuthHandlerInst.SaveOrInsertAdminMenu())
	authorized.DELETE("/admin_menu/menu", handler.AdminAuthHandlerInst.DeleteAdminMenu())
	authorized.POST("/admin_menu/refresh_child", handler.AdminAuthHandlerInst.AdminMenuRefreshChild())

	//管理员角色相关
	authorized.GET("/admin_role/roles", handler.AdminAuthHandlerInst.GetAdminRoleList())              // 角色列表
	authorized.GET("/admin_role/role_info", handler.AdminAuthHandlerInst.GetAdminRoleInfo())          // 获取角色信息
	authorized.GET("/admin_role/role_auth_table", handler.AdminAuthHandlerInst.GetAdminRoleUrlList()) // 获取角色菜单列表
	authorized.GET("/admin_role/role", handler.AdminAuthHandlerInst.GetAdminRole())
	authorized.POST("/admin_role/role", handler.AdminAuthHandlerInst.UpdateOrInsertAdminRole())
	authorized.PATCH("/admin_role/role_auth", handler.AdminAuthHandlerInst.UpdateOrInsertAdminRoleAuth()) // 授权给角色
	authorized.DELETE("/admin_role/role", handler.AdminAuthHandlerInst.DeleteAdminRole())                 // 删除角色-软删除
	authorized.POST("/admin_role/role_def", handler.AdminAuthHandlerInst.AuthAdminRole())

	//管理员账号相关操作
	authorized.GET("/admin_account/accounts", handler.AdminAuthHandlerInst.GetAdminAccountList()) //管理员账号列表
	authorized.GET("/admin_account/account", handler.AdminAuthHandlerInst.GetAdminAccount())      //获取单个管理员账号信息
	authorized.POST("/admin_account/account", handler.AdminAuthHandlerInst.SaveAdminAccount())    //添加或修改账号信息
	authorized.POST("/admin_account/modify_password", handler.AdminAuthHandlerInst.ModifyAdminPw())
	authorized.POST("/admin_account/reset_password", handler.AdminAuthHandlerInst.ResetAdminPw())
	authorized.PATCH("/admin_account/account_auth", handler.AdminAuthHandlerInst.UpdateOrInsertAdminAccountAuth()) //给管理员账号添加角色
	authorized.GET("/admin_account/role_from_acc", handler.AdminAuthHandlerInst.GetAdminRoleFromAcc())             //获取账号的角色列表

	//检测Admin账户是否存在
	authorized.GET("/admin_account/check", handler.AdminAuthHandlerInst.CheckAdminAccount())

	//-------------用户中心
	//用户管理
	//authorized.GET("/cust/get_custlist", rh.HCust.GetCustList())
	authorized.GET("/cust/get_custinfo", handler.CustHandlerInst.GetCustInfo())        //获取用户详细信息
	authorized.POST("/cust/modify_custinfo", handler.CustHandlerInst.ModifyCustInfo()) //修改用户权限

	authorized.GET("/bill/get_cust_bills", handler.CustHandlerInst.GetCustBills())            //WEB获取指定用户账单明细列表
	authorized.GET("/bill/get_servicer_bills", handler.CustHandlerInst.GetSerBills())         //WEB获取指定服务商账单明细列表
	authorized.GET("/bill/get_business_bills", handler.CustHandlerInst.GetBusinessBillList()) //WEB获取指定商家账单明细列表

	//钱包管理
	authorized.GET("/wallet/get_change_balance_orders", handler.CustHandlerInst.GetChangeBalanceOrders()) //WEB获取改变余额流水日志
	authorized.POST("/wallet/add_change_balance_order", handler.BillHandlerInst.AddChangeBalanceOrder())  //添加修改余额订单

	//认证管理
	authorized.GET("/auth_material/auth_materials", handler.CustHandlerInst.GetAuthMaterials())                                //获取个人身份认证信息列表
	authorized.POST("/auth_material/auth_material_status", handler.CustHandlerInst.ModifyAuthMaterialStatus())                 //修改个人身份认证信息的认证状态
	authorized.GET("/auth_material/personal_business_list", handler.CustHandlerInst.GetAuthMaterialBusinessList())             //获取个人商家认证信息列表
	authorized.POST("/auth_material/personal_business_status", handler.CustHandlerInst.ModifyAuthMaterialBusinessStatus())     //修改个人商家认证信息的认证状态
	authorized.GET("/auth_material/enterprise_business_list", handler.CustHandlerInst.GetAuthMaterialEnterpriseList())         //获取企业商家认证信息列表
	authorized.POST("/auth_material/enterprise_business_status", handler.CustHandlerInst.ModifyAuthMaterialEnterpriseStatus()) //修改企业商家认证信息的认证状态

	authorized.GET("/auth_material/business/update_info_list", handler.CustHandlerInst.GetAuthMaterialBusinessUpdateList())         //获取修改商家认证信息审核列表
	authorized.POST("/auth_material/business/update_info_status", handler.CustHandlerInst.ModifyAuthMaterialBusinessUpdateStatus()) //更新修改商家认证信息的审核状态

	//服务商管理
	authorized.GET("/servicer/get_servicer_list", handler.CustHandlerInst.GetServicerList())            //获取服务商列表
	authorized.GET("/servicer/get_servicer_phone_list", handler.CustHandlerInst.GetServicerPhoneList()) //获取服务商手机号列表
	authorized.GET("/servicer/get_servicer_info", handler.CustHandlerInst.GetServicerInfo())            //获取指定服务商信息（查看）
	authorized.POST("/servicer/modify_servicer_status", handler.CustHandlerInst.ModifyServicerStatus()) //修改服务商状态
	authorized.POST("/servicer/modify_servicer_config", handler.CustHandlerInst.ModifyServicerConfig()) //修改服务商配置
	authorized.POST("/servicer/modify_servicer_info", handler.CustHandlerInst.ModifyServicerInfo())     //修改服务商部分信息（图片（营业执照、经营场所照片）、基本信息）

	authorized.GET("/servicer/get_service_transactions", handler.CustHandlerInst.GetServiceTransactions())    //服务商交易明细查询(只查存取款)
	authorized.GET("/servicer/get_profit_ledger_list", handler.CustHandlerInst.GetServicerProfitLedgerList()) //服务商收益明细查询
	authorized.GET("/servicer/get_servicer_accounts", handler.CustHandlerInst.GetServicerAccounts())          //服务商账户查询

	authorized.GET("/servicer/get_servicer_order_count", handler.CustHandlerInst.GetServicerOrderCount()) //获取服务商收取款统计信息

	//POS机管理
	authorized.GET("/terminal/get_list", handler.CustHandlerInst.GetTerminalList())     //获取终端列表
	authorized.POST("/terminal/add", handler.CustHandlerInst.AddTerminal())             //添加终端
	authorized.POST("/terminal/enable", handler.CustHandlerInst.UpdateTerminalStatus()) //修改终端使用状态
	authorized.POST("/terminal/del", handler.CustHandlerInst.DelTerminal())             //删除终端

	//商家管理
	authorized.GET("/business/get_business_accounts", handler.CustHandlerInst.GetBusinessAccounts()) //商家账户查询
	authorized.GET("/business/profit", handler.CustHandlerInst.GetBusinessAccountsProfit())          //商家账户收益查询
	authorized.GET("/business/profit_detail_list", handler.CustHandlerInst.GetBusinessProfitList())  //商家账户收益明细查询

	authorized.GET("/business/get_business_list", handler.CustHandlerInst.GetBusinessList())            //获取商家列表
	authorized.POST("/business/modify_business_status", handler.CustHandlerInst.ModifyBusinessStatus()) //修改状态（商家状态、商家收款权限、商家出款权限）

	authorized.GET("/business/get_business_industrys", handler.CustHandlerInst.GetBusinessIndustryList())           //获取主要行业应用列表
	authorized.GET("/business/get_business_industry_detail", handler.CustHandlerInst.GetBusinessIndustryDetail())   //获取主要行业应用详情
	authorized.POST("/business/update_business_industry", handler.CustHandlerInst.InsertOrUpdateBusinessIndustry()) //添加主要行业应用
	authorized.POST("/business/del_business_industry", handler.CustHandlerInst.DelBusinessIndustry())               //删除主要行业应用

	//行业渠道费率与结算周期
	authorized.GET("/business/get_industrys_rate_cycle_list", handler.CustHandlerInst.GetBusinessIndustryRateCycleList())      //获取主要行业应用列表
	authorized.POST("/business/update_industry_rate_cycle", handler.CustHandlerInst.InsertOrUpdateBusinessIndustryRateCycle()) //添加主要行业应用
	authorized.POST("/business/del_industry_rate_cycle", handler.CustHandlerInst.DelBusinessIndustryRateCycle())               //删除主要行业应用

	//产品管理
	authorized.GET("/business/scene/list", handler.CustHandlerInst.GetBusinessSceneList())           //获取产品列表
	authorized.GET("/business/scene/detail", handler.CustHandlerInst.GetBusinessSceneDetail())       //获取产品详情
	authorized.POST("/business/scene/update", handler.CustHandlerInst.InsertOrUpdateBusinessScene()) //添加或修改产品
	authorized.POST("/business/scene/is_enabled", handler.CustHandlerInst.IsEnabled())               // 启用/禁用产品
	authorized.POST("/business/scene/update_idx", handler.CustHandlerInst.UpdateBusinessSceneIdx())  //修改产品的序号

	//商家签约管理(这是旧的逻辑，应用签产品、行业)
	//authorized.GET("/business/signed/list", handler.CustHandlerInst.GetBusinessSignedList())                //获取商家签约列表
	//authorized.POST("/business/signed/update_info", handler.CustHandlerInst.UpdateBusinessSignedInfo())     //修改商家应用签约产品的结算周期、费率
	//authorized.POST("/business/signed/update_status", handler.CustHandlerInst.UpdateBusinessSignedStatus()) //修改商家应用签约状态

	//商家签约管理
	authorized.GET("/business/signed/list", handler.CustHandlerInst.GetBusinessSceneSignedList())                //获取商家产品签约列表
	authorized.POST("/business/signed/update_info", handler.CustHandlerInst.UpdateBusinessSceneSignedInfo())     //修改商家产品签约的结算周期、费率
	authorized.POST("/business/signed/update_status", handler.CustHandlerInst.UpdateBusinessSceneSignedStatus()) //修改商家产品签约状态

	//商家应用管理
	authorized.GET("/business/app/list", handler.CustHandlerInst.GetBusinessAppList())                //获取商家应用审核列表
	authorized.POST("/business/app/update_status", handler.CustHandlerInst.UpdateBusinessAppStatus()) //修改商家应用审核状态

	//财务管理
	authorized.GET("/financial/get_financial_servicer_check", handler.CustHandlerInst.GetFinancialServicerCheckList()) //获取服务商的对账单
	authorized.GET("/financial/get_billing_details_results", handler.CustHandlerInst.GetBillingDetailsResultsList())   //查看指定服务商的某天账单明细
	authorized.GET("/financial/collection_management_list", handler.CustHandlerInst.CollectionManagementList())        //平台收款账户管理
	authorized.POST("/financial/modify_collect_status", handler.CustHandlerInst.ModifyCollectStatus())                 //平台收款账户收款状态修改
	authorized.DELETE("/financial/delect_card", handler.CustHandlerInst.DelectCard())                                  //删除平台收款账户
	authorized.GET("/financial/get_card_info", handler.CustHandlerInst.GetCardlInfo())                                 //获取平台收款账户(单个)
	authorized.POST("/financial/card_edit", handler.CustHandlerInst.UpdateOrInsertCard())                              //添加或修改平台收款账户

	authorized.GET("/financial/get_headquarters_profit_list", handler.CustHandlerInst.GetHeadquartersProfitList())         //平台盈利统计
	authorized.POST("/financial/headquarters_profit_withdraw", handler.BillHandlerInst.InsertHeadquartersProfitWithdraw()) //平台盈利提现
	authorized.GET("/financial/headquarters_profit_withdraws", handler.CustHandlerInst.HeadquartersProfitWithdraws())      //获取平台盈利提现订单列表

	//交易管理
	authorized.GET("/trading/get_to_headquarters", handler.CustHandlerInst.GetToHeadquartersList())      //获取服务商充值订单(转账至总部流水)
	authorized.POST("/trading/update_to_headquarters", handler.CustHandlerInst.UpdateToHeadquarters())   //更改转账至总部流水订单状态
	authorized.GET("/trading/get_to_servicers", handler.CustHandlerInst.GetToServicerList())             //获取服务商提现订单(转账至服务商流水)
	authorized.POST("/trading/add_to_servicer", handler.CustHandlerInst.AddToServicer())                 //增加转账 转账至总部流水
	authorized.GET("/trading/get_transfer_orders", handler.CustHandlerInst.GetTransferOrderList())       //获取转账流水
	authorized.GET("/trading/get_outgo_orders", handler.CustHandlerInst.GetOutgoOrderList())             //获取取款流水
	authorized.GET("/trading/get_income_orders", handler.CustHandlerInst.GetIncomeOrderList())           //获取存款流水
	authorized.GET("/trading/get_exchange_orders", handler.CustHandlerInst.GetExchangeOrderList())       //获取兑换流水
	authorized.GET("/trading/get_collection_orders", handler.CustHandlerInst.GetCollectionOrders())      //获取收款流水
	authorized.GET("/trading/get_log_vaccounts", handler.CustHandlerInst.GetLogVaccounts())              //获取虚拟账户日志流水列表
	authorized.POST("/trading/update_service_withdraw", handler.CustHandlerInst.UpdateServiceWithdraw()) // 修改服务商提现订单状态(向服务商取款修改订单状态)

	authorized.GET("/trading/write_off/list", handler.CustHandlerInst.GetWriteOffCodeList())     //获取核销码列表
	authorized.POST("/trading/write_off/dispose", handler.CustHandlerInst.DisposeWriteOffCode()) //处理核销码（冻结、解冻、注销）

	authorized.GET("/trading/get_to_custs", handler.CustHandlerInst.GetToCustList())               //获取用户提现申请订单列表
	authorized.POST("/trading/update_cust_withdraw", handler.CustHandlerInst.UpdateCustWithdraw()) // 修改用户向总部提现订单状态

	authorized.GET("/trading/cust_to_headquarters", handler.CustHandlerInst.GetCustToHeadquartersList()) //获取用户向总部存款订单列表
	authorized.POST("/trading/update_cust_save", handler.CustHandlerInst.UpdateCustSave())               // 修改用户向总部存款订单状态

	authorized.GET("/trading/get_cash_recharge_orders", handler.CustHandlerInst.GetCashRechargeOrderList()) //获取服务商现金充值流水
	authorized.POST("/trading/add_cash_recharge_order", handler.BillHandlerInst.AddCashRecharge())          // 添加服务商现金充值

	authorized.GET("/trading/business_to_headquarters", handler.CustHandlerInst.GetBusinessToHeadList())             //获取商家向总部充值订单列表
	authorized.POST("/trading/update_business_to_head_status", handler.BillHandlerInst.UpdateBusinessToHeadStatus()) //审核商家充值

	authorized.GET("/trading/to_business", handler.CustHandlerInst.GetToBusinessList())                     //获取商家提现订单列表
	authorized.POST("/trading/update_to_business_status", handler.BillHandlerInst.UpdateToBusinessStatus()) //审核商家提现

	authorized.GET("/trading/business/transfer_order_list", handler.CustHandlerInst.GetBusinessTransferOrderList())       //商家转账订单列表
	authorized.GET("/trading/business/batch_transfer_order_list", handler.CustHandlerInst.GetBusinessTransferBatchList()) //商家批量转账列表

	authorized.POST("/trading/create_bills_file", handler.CustHandlerInst.CreateXlsxFile()) // 生成并下载对应订单信息的xlsx文件

	authorized.GET("/trading/mod/get_bills", handler.CustHandlerInst.GetModPayBills())                   //获取商家交易订单列表
	authorized.GET("/trading/mod/get_refund_bills", handler.CustHandlerInst.GetModRefundBills())         //获取商家退款订单列表
	authorized.GET("/trading/channel/get_bills", handler.CustHandlerInst.GetChannelBills())              //获取上游渠道交易订单列表
	authorized.GET("/trading/channel/get_refund_bills", handler.CustHandlerInst.GetChannelRefundBills()) //获取上游渠道退款订单列表
	authorized.GET("/trading/get_bills_detail", handler.CustHandlerInst.GetBillsDetail())                //获取交易订单详情
	authorized.GET("/trading/get_refund_detail", handler.CustHandlerInst.GetRefundBillsDetail())         //获取退款订单详情

	authorized.POST("/trading/channel_bill_settle", handler.BusinessBillHandlerInst.ManualSettle()) //上游渠道订单结算

	//支付渠道管理
	authorized.GET("/payment_channel/list", handler.CustHandlerInst.GetAllPaymentChannel())    //获取全部支付渠道
	authorized.POST("/payment_channel/add", handler.CustHandlerInst.AddPaymentChannel())       //添加支付渠道
	authorized.POST("/payment_channel/update", handler.CustHandlerInst.UpdatePaymentChannel()) //修改支付渠道

	//配置管理
	authorized.GET("/withdraw_config/withdraw_config", handler.CustHandlerInst.GetWithdrawConfig())     //获取存取转配置
	authorized.POST("/withdraw_config/withdraw_config", handler.CustHandlerInst.UpdateWithdrawConfig()) //

	authorized.GET("/exchange_rate_config/exchange_rate_config", handler.CustHandlerInst.GetExchangeRateConfig())     //兑换率
	authorized.POST("/exchange_rate_config/exchange_rate_config", handler.CustHandlerInst.UpdateExchangeRateConfig()) //

	authorized.GET("/business_config/list", handler.CustHandlerInst.GetBusinessConfig())       //商家付款配置列表
	authorized.POST("/business_config/update", handler.CustHandlerInst.UpdateBusinessConfig()) //修改

	authorized.GET("/func_config/func_configs", handler.CustHandlerInst.GetFuncConfig())                        //获取钱包功能入口配置
	authorized.GET("/func_config/func_config", handler.CustHandlerInst.GetFuncConfigDetail())                   //钱包功能入口配置
	authorized.POST("/func_config/func_config", handler.CustHandlerInst.UpdateFuncConfig())                     //
	authorized.POST("/func_config/swap_idx", handler.CustHandlerInst.SwapFuncConfigIdx())                       //
	authorized.DELETE("/func_config/func_config", handler.CustHandlerInst.DeleteFuncConfig())                   //
	authorized.POST("/func_config/modify_func_config_use", handler.CustHandlerInst.ModifyUseStatusFuncConfig()) //

	authorized.GET("/transfer_security_config/transfer_security_config", handler.CustHandlerInst.GetTransferSecurityConfig())     //交易安全配置
	authorized.POST("/transfer_security_config/transfer_security_config", handler.CustHandlerInst.UpdateTransferSecurityConfig()) //更新交易安全配置

	authorized.GET("/write_off/get_duration_date_config", handler.CustHandlerInst.GetWriteOffDurationDateConfig())        //核销码有效期
	authorized.POST("/write_off/update_duration_date_config", handler.CustHandlerInst.UpdateWriteOffDurationDateConfig()) //更新核销码有效期

	authorized.GET("/income_ougo_config/income_ougo_configs", handler.CustHandlerInst.GetIncomeOugoConfig())            //充值、提现方式配置列表
	authorized.GET("/income_ougo_config/income_ougo_config", handler.CustHandlerInst.GetIncomeOugoConfigDetail())       //充值、提现方式配置（单个）
	authorized.POST("/income_ougo_config/income_ougo_config", handler.CustHandlerInst.UpdateIncomeOugoConfig())         //插入与更新充值、提现方式配置
	authorized.POST("/income_ougo_config/update_use_status", handler.CustHandlerInst.UpdateIncomeOugoConfigUseStatus()) //修改充值、提现方式配置状态
	authorized.POST("/income_ougo_config/update_idx", handler.CustHandlerInst.UpdateIncomeOugoConfigIdx())              //修改充值、提现方式排序
	authorized.DELETE("/income_ougo_config/delete_config", handler.CustHandlerInst.DeleteIncomeOugoConfig())            //删除充值、提现方式配置

	authorized.GET("/lang/langs", handler.CustHandlerInst.GetLangs())           //多语言配置
	authorized.GET("/lang/lang", handler.CustHandlerInst.GetLangDetail())       //多语言配置(单个)
	authorized.POST("/lang/lang", handler.CustHandlerInst.InsertOrUpdateLang()) //插入与更新多语言配置
	authorized.DELETE("/lang/lang", handler.CustHandlerInst.DeleteLang())       //删除多语言配置

	//用户仓库接口
	authorized.GET("/channel_use/channels", handler.CustHandlerInst.GetUseChannels())               //获取渠道列表
	authorized.GET("/channel_use/channel", handler.CustHandlerInst.GetUseChannelDetail())           //获取单个渠道
	authorized.POST("/channel_use/update_channel", handler.CustHandlerInst.InsertUseChannel())      //插入或更新银行卡收取款渠道配置
	authorized.DELETE("/channel_use/delete_channel", handler.CustHandlerInst.DeleteUseChannel())    //删除渠道
	authorized.POST("/channel_use/modify_status", handler.CustHandlerInst.ModifyUseChannelStatus()) //修改渠道状态

	//服务商仓库接口
	authorized.GET("/channel_pos/channels", handler.CustHandlerInst.GetPosChannels())                  //获取渠道列表
	authorized.POST("/channel_pos/update_channel", handler.CustHandlerInst.InsertPosChannel())         //插入服务商渠道配置
	authorized.DELETE("/channel_pos/delete_channel", handler.CustHandlerInst.DeletePosChannel())       //删除渠道
	authorized.POST("/channel_pos/modify_status", handler.CustHandlerInst.ModifyChannelPosStatus())    //修改渠道状态
	authorized.POST("/channel_pos/modify_is_recom", handler.CustHandlerInst.ModifyChannelPosIsRecom()) //修改渠道是否推荐

	//商家仓库接口
	authorized.GET("/channel_business/channels", handler.CustHandlerInst.GetBusinessChannels())               //获取渠道列表
	authorized.GET("/channel_business/channel", handler.CustHandlerInst.GetBusinessChannelDetail())           //获取单个渠道
	authorized.POST("/channel_business/update_channel", handler.CustHandlerInst.InsertBusinessChannel())      //插入或更新渠道配置
	authorized.DELETE("/channel_business/delete_channel", handler.CustHandlerInst.DeleteBusinessChannel())    //删除渠道
	authorized.POST("/channel_business/modify_status", handler.CustHandlerInst.ModifyBusinessChannelStatus()) //修改渠道状态

	//渠道仓库接口
	authorized.GET("/channel/channels", handler.CustHandlerInst.GetChannels())                  //获取渠道仓库列表
	authorized.GET("/channel/channel", handler.CustHandlerInst.GetChannelDetail())              //获取渠道仓库(单个)
	authorized.POST("/channel/update_channel", handler.CustHandlerInst.InsertOrUpdateChannel()) //插入与更新渠道仓库配置
	authorized.DELETE("/channel/delete_channel", handler.CustHandlerInst.DeleteChannel())       //删除渠道
	authorized.POST("/channel/modify_status", handler.CustHandlerInst.ModifyChannelStatus())    //修改渠道状态

	//客服管理
	authorized.GET("/help/helps_count", handler.CustHandlerInst.GetCommonHelpCount())  // 获取帮助统计
	authorized.GET("/help/helps", handler.CustHandlerInst.GetCommonHelps())            // 获取帮助列表
	authorized.GET("/help/help", handler.CustHandlerInst.GetCommonHelpDetail())        // 获取帮助（单个）
	authorized.POST("/help/help", handler.CustHandlerInst.InsertOrUpdateHelp())        //插入或更新帮助
	authorized.DELETE("/help/help", handler.CustHandlerInst.DeleteHelp())              //删除帮助
	authorized.POST("/help/modify_status", handler.CustHandlerInst.ModifyHelpStatus()) //修改帮助的使用状态
	authorized.POST("/help/swap_idx", handler.CustHandlerInst.SwapHelpIdx())           //修改帮助的排序

	authorized.GET("/consultation_config/consultation_configs", handler.CustHandlerInst.GetConsultationConfigs())           // 获取咨询服务列表
	authorized.GET("/consultation_config/consultation_config", handler.CustHandlerInst.GetConsultationConfigDetail())       // 获取咨询服务（单个）
	authorized.POST("/consultation_config/consultation_config", handler.CustHandlerInst.InsertOrUpdateConsultationConfig()) //插入或更新咨询服务
	authorized.DELETE("/consultation_config/consultation_config", handler.CustHandlerInst.DeleteConsultationConfig())       //删除咨询服务
	authorized.POST("/consultation_config/modify_status", handler.CustHandlerInst.ModifyConsultationConfigStatus())         //修改咨询服务的使用状态
	authorized.POST("/consultation_config/swap_idx", handler.CustHandlerInst.SwapConsultationIdx())                         //移动优先级

	//协议
	authorized.GET("/agreement/agreements", handler.CustHandlerInst.GetAgreements())             // 获取协议列表
	authorized.GET("/agreement/agreement", handler.CustHandlerInst.GetAgreementDetail())         // 获取协议（单个）
	authorized.POST("/agreement/agreement", handler.CustHandlerInst.InsertOrUpdateAgreement())   //插入或更新协议
	authorized.DELETE("/agreement/agreement", handler.CustHandlerInst.DeleteAgreement())         //删除协议
	authorized.POST("/agreement/modify_status", handler.CustHandlerInst.ModifyAgreementStatus()) //修改协议的使用状态

	//版本
	authorized.GET("/app_version/app_version_count", handler.CustHandlerInst.GetAppVersionsCount())    //获取版本统计信息。
	authorized.GET("/app_version/app_versions", handler.CustHandlerInst.GetAppVersions())              // 获取app、pos版本列表
	authorized.GET("/app_version/app_version", handler.CustHandlerInst.GetAppVersion())                //获取单个版本的信息
	authorized.GET("/app_version/get_new_version", handler.CustHandlerInst.GetNewVersion())            //获取最新版本
	authorized.POST("/app_version/modify_status", handler.CustHandlerInst.ModifyAppVersionStatus())    //修改版本状态
	authorized.POST("/app_version/modify_is_force", handler.CustHandlerInst.ModifyAppVersionIsForce()) //修改版本强制更新状态
	authorized.POST("/app_version/app_version", handler.CustHandlerInst.InsertOrUpdateAppVersion())    //插入或更新版本

	authorized.POST("/upload/upload_file", handler.CustHandlerInst.UploadFile()) //上传App文件

	authorized.GET("/img/get_img_url", handler.CustHandlerInst.GetImgUrl()) //根据图片id获取url

	authorized.GET("/client_info/client_infos", handler.CustHandlerInst.GetClientInfos()) // 获取收集的用户设备信息列表

	//推送页面
	authorized.GET("/push/push_confs", handler.CustHandlerInst.GetPushConfs())                //
	authorized.GET("/push/push_conf", handler.CustHandlerInst.GetPushConf())                  //
	authorized.POST("/push/push_conf", handler.CustHandlerInst.InsertOrUpdatePushConf())      //
	authorized.DELETE("/push/del_push_conf", handler.CustHandlerInst.DeletePushConf())        //
	authorized.POST("/push/push_conf_status", handler.CustHandlerInst.ModifyPushConfStatus()) //

	authorized.GET("/push/push_records", handler.CustHandlerInst.GetPushRecords())       //
	authorized.GET("/push/push_temps", handler.CustHandlerInst.GetPushTemps())           //
	authorized.GET("/push/push_temp", handler.CustHandlerInst.GetPushTemp())             // 推送模板列表
	authorized.POST("/push/push_temp", handler.CustHandlerInst.InsertOrUpdatePushTemp()) //插入或更新推送模板
	authorized.DELETE("/push/del_push_temp", handler.CustHandlerInst.DeletePushTemp())   //

	//商家公告
	authorized.GET("/bulletin/list", handler.CustHandlerInst.GetBulletins())                   //公告列表
	authorized.POST("/bulletin/update", handler.CustHandlerInst.InsertOrUpdateBulletin())      //插入或更新公告
	authorized.POST("/bulletin/delete", handler.CustHandlerInst.DelBulletin())                 //删除公告
	authorized.POST("/bulletin/update_status", handler.CustHandlerInst.UpdateBulletinStatus()) //发布公告、置顶公告

	//商家消息
	authorized.GET("/business_messages/list", handler.CustHandlerInst.GetBusinessMessagesList()) //获取商家消息列表

	//风控配置界面接口
	//event
	authorized.GET("/risk/events", handler.CustHandlerInst.GetEvents())           //获取列表
	authorized.GET("/risk/event", handler.CustHandlerInst.GetEvent())             //获取单个
	authorized.POST("/risk/event", handler.CustHandlerInst.InsertOrUpdateEvent()) //插入或更新
	authorized.DELETE("/risk/del_event", handler.CustHandlerInst.DeleteEvent())   //删除

	//风控中心-指纹支付配置管理
	authorized.GET("/risk/app_fingerprint_on", handler.CustHandlerInst.GetAppFingerprintOn())                    //获取用户设备指纹列表
	authorized.GET("/risk/fingerprint_list", handler.CustHandlerInst.GetFingerprintList())                       //获取用户设备指纹列表
	authorized.POST("/risk/close_fingerprint_function", handler.CustHandlerInst.CloseFingerprintInputFunction()) //关闭指纹录入功能
	authorized.POST("/risk/clean_fingerprint", handler.CustHandlerInst.CleanFingerprintData())                   //清除数据

	//eva_param
	authorized.GET("/risk/eva_params", handler.CustHandlerInst.GetEvaParams())           //
	authorized.POST("/risk/eva_param", handler.CustHandlerInst.InsertOrUpdateEvaParam()) //插入或更新
	authorized.DELETE("/risk/del_eva_param", handler.CustHandlerInst.DeleteEvaParam())   //删除

	//global_param
	authorized.GET("/risk/global_params", handler.CustHandlerInst.GetGlobalParams())           //
	authorized.POST("/risk/global_param", handler.CustHandlerInst.InsertOrUpdateGlobalParam()) //插入或更新
	authorized.DELETE("/risk/del_global_param", handler.CustHandlerInst.DeleteGlobalParam())   //删除

	//log_result
	authorized.GET("/risk/log_results", handler.CustHandlerInst.GetLogResults()) //

	//op
	authorized.GET("/risk/ops", handler.CustHandlerInst.GetOps())           //
	authorized.GET("/risk/op", handler.CustHandlerInst.GetOp())             //获取单个
	authorized.POST("/risk/op", handler.CustHandlerInst.InsertOrUpdateOp()) //插入或更新
	authorized.DELETE("/risk/del_op", handler.CustHandlerInst.DeleteOp())   //删除

	//rela_api_event
	authorized.GET("/risk/rela_api_events", handler.CustHandlerInst.GetRelaApiEvents())           //
	authorized.POST("/risk/rela_api_event", handler.CustHandlerInst.InsertOrUpdateRelaApiEvent()) //插入或更新
	authorized.DELETE("/risk/del_rela_api_event", handler.CustHandlerInst.DeleteRelaApiEvent())   //删除

	//rela_event_rule
	authorized.GET("/risk/rela_event_rules", handler.CustHandlerInst.GetRelaEventRules())           //
	authorized.POST("/risk/rela_event_rule", handler.CustHandlerInst.InsertOrUpdateRelaEventRule()) //插入或更新
	authorized.DELETE("/risk/del_rela_event_rule", handler.CustHandlerInst.DeleteRelaEventRule())   //删除

	//risk_threshold
	authorized.GET("/risk/risk_thresholds", handler.CustHandlerInst.GetRiskThresholds())           //
	authorized.POST("/risk/risk_threshold", handler.CustHandlerInst.InsertOrUpdateRiskThreshold()) //插入或更新
	authorized.DELETE("/risk/del_risk_threshold", handler.CustHandlerInst.DeleteRiskThreshold())   //删除

	authorized.GET("/risk/rules", handler.CustHandlerInst.GetRules())           //
	authorized.GET("/risk/rule", handler.CustHandlerInst.GetRule())             //获取单个
	authorized.POST("/risk/rule", handler.CustHandlerInst.InsertOrUpdateRule()) //插入或更新
	authorized.DELETE("/risk/del_rule", handler.CustHandlerInst.DeleteRule())   //删除

	//用户提现统计
	authorized.GET("/statistic_user/get_withdraw_chart_datas", handler.CustHandlerInst.GetStatisticUserWithdraws()) //
	authorized.GET("/statistic_user/get_withdraw_datas", handler.CustHandlerInst.GetStatisticUserWithdrawList())    //

	//用户充值统计
	authorized.GET("/statistic_user/get_recharge_chart_datas", handler.CustHandlerInst.GetStatisticUserRecharges()) //
	authorized.GET("/statistic_user/get_recharge_datas", handler.CustHandlerInst.GetStatisticUserRechargeList())    //

	//用户兑换统计
	authorized.GET("/statistic_user/get_exchange_chart_datas", handler.CustHandlerInst.GetStatisticUserExchanges()) //
	authorized.GET("/statistic_user/get_exchange_datas", handler.CustHandlerInst.GetStatisticUserExchangeList())    //

	//用户转账统计
	authorized.GET("/statistic_user/get_transfer_chart_datas", handler.CustHandlerInst.GetStatisticUserTransfers()) //
	authorized.GET("/statistic_user/get_transfer_datas", handler.CustHandlerInst.GetStatisticUserTransferList())    //

	//按天统计统计
	authorized.GET("/statistic_user/get_day_chart_datas", handler.CustHandlerInst.GetStatisticDates()) //
	authorized.GET("/statistic_user/get_day_datas", handler.CustHandlerInst.GetStatisticDateList())    //

	// 对统计数据进行重新统计
	authorized.POST("/statistic_user/re_statistic", handler.CustHandlerInst.ReStatistic()) //

	//用户资金总留存统计
	authorized.GET("/statistic_user/get_user_money_chart_datas", handler.CustHandlerInst.GetStatisticUserMoneyDates()) //
	authorized.GET("/statistic_user/get_user_money_datas", handler.CustHandlerInst.GetStatisticUserMoneyList())        //

	//支付系统日志
	authorized.GET("/log/api_pay", handler.CustHandlerInst.GetApiPayLogList()) //

	return router
}
