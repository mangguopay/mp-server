package handler

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/auth-srv/dao"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_func"
	"database/sql"
)

type VAccountHandler struct {
}

var VAccountHandlerInst VAccountHandler

//初始化企业商家虚账
func (VAccountHandler) InitVAccount(tx *sql.Tx, accountNo string) error {
	//查询账号是否已存在虚账
	vAccList, vAccErr := dao.VaccountDaoInst.GetVAccNoByAccountNo(accountNo)
	if vAccErr != nil && vAccErr != sql.ErrNoRows {
		ss_log.Error("查询账号的已有虚账失败, accountNo=%v, err=%v", accountNo, vAccErr)
		return vAccErr
	}

	if len(vAccList) >= 1 {
		ss_log.Info("企业商家账号(%v)已有虚账，vAccountLsit=%v", accountNo, strext.ToJson(vAccList))
		return nil
	}
	//创建币种为USD的相关账号
	//美金--商家已结算
	usdSettled, err := dao.VaccountDaoInst.InitVAccountNoTx(tx, accountNo, constants.VaType_USD_BUSINESS_SETTLED, constants.CURRENCY_USD)
	if err != nil {
		ss_log.Error("初始化企业商家虚账失败，accountNo=%v, vaccType=%v, err=%v", accountNo, constants.VaType_USD_BUSINESS_SETTLED, err)
		return err
	}
	//美金--商家未结算
	usdUnSettled, err := dao.VaccountDaoInst.InitVAccountNoTx(tx, accountNo, constants.VaType_USD_BUSINESS_UNSETTLED, constants.CURRENCY_USD)
	if err != nil {
		ss_log.Error("初始化企业商家虚账失败，accountNo=%v, vaccType=%v, err=%v", accountNo, constants.VaType_USD_BUSINESS_UNSETTLED, err)
		return err
	}

	//创建币种为KHR的相关账号
	//瑞尔--商家已结算
	khrSettled, err := dao.VaccountDaoInst.InitVAccountNoTx(tx, accountNo, constants.VaType_KHR_BUSINESS_SETTLED, constants.CURRENCY_KHR)
	if err != nil {
		ss_log.Error("初始化企业商家虚账失败，accountNo=%v, vaccType=%v, err=%v", accountNo, constants.VaType_KHR_BUSINESS_SETTLED, err)
		return err
	}
	//瑞尔--商家未结算
	khrUnSettled, err := dao.VaccountDaoInst.InitVAccountNoTx(tx, accountNo, constants.VaType_KHR_BUSINESS_UNSETTLED, constants.CURRENCY_KHR)
	if err != nil {
		ss_log.Error("初始化企业商家虚账失败，accountNo=%v, vaccType=%v, err=%v", accountNo, constants.VaType_KHR_BUSINESS_UNSETTLED, err)
		return err
	}

	logInfo := map[string]string{
		"美金--商家已结算虚账": usdSettled,
		"美金--商家未结算虚账": usdUnSettled,
		"瑞尔--商家已结算虚账": khrSettled,
		"瑞尔--商家未结算虚账": khrUnSettled,
	}
	ss_log.Info("初始化企业商家虚账完成；accountNo=%v, vAccountList=%v", accountNo, strext.ToJson(logInfo))
	return nil
}

//同步个人用户虚账余额
func (VAccountHandler) SyncFreezeVAccountBalance(tx *sql.Tx, accountNo string) error {
	//查询是否有未激活账户
	vAccNoList, vAccErr := dao.VaccountDaoInst.GetFreezeVAccountNo(accountNo)
	if vAccErr != nil {
		ss_log.Error("查询个人冻结存款虚账失败, accountNo=%v, err=%v", accountNo, vAccErr)
		return vAccErr
	}

	if len(vAccNoList) <= 0 {
		ss_log.Info("该账号没有个人冻结存款虚账， accountNo=%v", accountNo)
		return nil
	}

	//创建对应币种的个人虚账
	for _, v := range vAccNoList {
		var vAccNo string
		var err error
		//创建个人虚账
		if strext.ToInt(v.VAccType) == constants.VaType_FREEZE_USD_DEBIT { //美金
			vAccNo, err = dao.VaccountDaoInst.InitVAccountNoTx(tx, accountNo, constants.VaType_USD_DEBIT, constants.CURRENCY_USD)
			if err != nil {
				ss_log.Error("初始化个人虚账失败，accountNo=%v, vaccType=%v, err=%v", accountNo, constants.VaType_USD_DEBIT, err)
				return err
			}

		} else if strext.ToInt(v.VAccType) == constants.VaType_FREEZE_KHR_DEBIT { //瑞尔
			vAccNo, err = dao.VaccountDaoInst.InitVAccountNoTx(tx, accountNo, constants.VaType_KHR_DEBIT, constants.CURRENCY_KHR)
			if err != nil {
				ss_log.Error("初始化个人虚账失败，accountNo=%v, vaccType=%v, err=%v", accountNo, constants.VaType_KHR_DEBIT, err)
				return err
			}
		}

		//这里的金额变动只是把未激活虚账的钱同步到激活虚账，不需要记录账户金额变动日志
		//减少冻结虚账金额
		frozenVAccbalance, _, err := dao.VaccountDaoInst.MinusBalance(tx, v.VAccountNo, v.Balance)
		if err != nil {
			ss_log.Error("减少用户冻结虚账余额失败，accountNo=%v, err=%v", accountNo, err)
			return err
		}

		// 判断余额是否为负数
		if r, err := ss_func.JudgeAmountPositiveOrNegative(frozenVAccbalance); err != nil || r < 0 {
			ss_log.Error("用户账号余额不足,err:%v,custBalance:%v, result:%v", err, frozenVAccbalance, r)
			return err
		}

		//增加个人虚账余额
		balance, _, err := dao.VaccountDaoInst.PlusBalance(tx, vAccNo, v.Balance)
		if err != nil {
			ss_log.Error("增加用户冻结虚账余额失败，accountNo=%v, err=%v", accountNo, err)
			return err
		}

		ss_log.Info("个人虚账余额(%v)：%v", vAccNo, balance)
	}

	//err := dao.VaccountDaoInst.SyncAccRemain(tx, accountNo)
	//if err != nil {
	//	ss_log.Error("同步用户账户，虚账金额失败, accountNo=%v, err=%v", accountNo, err)
	//	return err
	//}

	return nil
}
