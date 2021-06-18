package handler

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_data"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/common"
	"a.a/mp-server/cust-srv/dao"
	"a.a/mp-server/cust-srv/util"
)

/**
 * 获取配置信息
 */
func (*CustHandler) GetWithdrawConfig(ctx context.Context, req *go_micro_srv_cust.GetWithdrawConfigRequest, reply *go_micro_srv_cust.GetWithdrawConfigReply) error {
	switch req.ConfigType {
	case "face_withdraw": //面对面取款配置
		_, usdFaceWithdrawRate, _ := cache.ApiDaoInstance.GetGlobalParam("usd_face_withdraw_rate") //usd取款费率
		_, usdFaceMinWithdrawFee, _ := cache.ApiDaoInstance.GetGlobalParam("usd_face_min_withdraw_fee")
		_, usdFaceFreeFeePerYear, _ := cache.ApiDaoInstance.GetGlobalParam("usd_face_free_fee_per_year")
		_, usdFaceSingleMax, _ := cache.ApiDaoInstance.GetGlobalParam("usd_face_single_max")
		_, usdFaceSingleMin, _ := cache.ApiDaoInstance.GetGlobalParam("usd_face_single_min")

		_, khrFaceWithdrawRate, _ := cache.ApiDaoInstance.GetGlobalParam("khr_face_withdraw_rate")
		_, khrFaceMinWithdrawFee, _ := cache.ApiDaoInstance.GetGlobalParam("khr_face_min_withdraw_fee")
		_, khrFaceFreeFeePerYear, _ := cache.ApiDaoInstance.GetGlobalParam("khr_face_free_fee_per_year")
		_, khrFaceSingleMax, _ := cache.ApiDaoInstance.GetGlobalParam("khr_face_single_max")
		_, khrFaceSingleMin, _ := cache.ApiDaoInstance.GetGlobalParam("khr_face_single_min")

		//if usdFaceWithdrawRate != "0" || usdFaceWithdrawRate != "" {
		//	usdFaceWithdrawRateD, _ := decimal.NewFromString(usdFaceWithdrawRate)
		//	//向上舍入
		//	usdFaceWithdrawRateD = ss_big.SsBigInst.ToRound(usdFaceWithdrawRateD.Div(decimal.NewFromInt(10000)), 2, ss_big.RoundingMode_CEILING)
		//	usdFaceWithdrawRate = usdFaceWithdrawRateD.String()
		//}
		//if khrFaceWithdrawRate != "0" || khrFaceWithdrawRate != "" {
		//	tempD, _ := decimal.NewFromString(khrFaceWithdrawRate)
		//	tempD = ss_big.SsBigInst.ToRound(tempD.Div(decimal.NewFromInt(10000)), 2, ss_big.RoundingMode_CEILING)
		//	khrFaceWithdrawRate = tempD.String()
		//}

		reply.FaceWithdrawDatas = &go_micro_srv_cust.FaceWithdrawConfigData{
			UsdFaceWithdrawRate:   usdFaceWithdrawRate,
			UsdFaceMinWithdrawFee: usdFaceMinWithdrawFee,
			UsdFaceFreeFeePerYear: usdFaceFreeFeePerYear,
			UsdFaceSingleMax:      usdFaceSingleMax,
			UsdFaceSingleMin:      usdFaceSingleMin,
			KhrFaceWithdrawRate:   khrFaceWithdrawRate,
			KhrFaceMinWithdrawFee: khrFaceMinWithdrawFee,
			KhrFaceFreeFeePerYear: khrFaceFreeFeePerYear,
			KhrFaceSingleMax:      khrFaceSingleMax,
			KhrFaceSingleMin:      khrFaceSingleMin,
		}
	case "phone_withdraw": //手机号取款配置
		_, usdPhoneWithdrawRate, _ := cache.ApiDaoInstance.GetGlobalParam("usd_phone_withdraw_rate") //usd取款费率
		_, usdPhoneMinWithdrawFee, _ := cache.ApiDaoInstance.GetGlobalParam("usd_phone_min_withdraw_fee")
		_, usdPhoneFreeFeePerYear, _ := cache.ApiDaoInstance.GetGlobalParam("usd_phone_free_fee_per_year")
		_, usdPhoneSingleMax, _ := cache.ApiDaoInstance.GetGlobalParam("usd_phone_single_max")
		_, usdPhoneSingleMin, _ := cache.ApiDaoInstance.GetGlobalParam("usd_phone_single_min")

		_, khrPhoneWithdrawRate, _ := cache.ApiDaoInstance.GetGlobalParam("khr_phone_withdraw_rate")
		_, khrPhoneMinWithdrawFee, _ := cache.ApiDaoInstance.GetGlobalParam("khr_phone_min_withdraw_fee")
		_, khrPhoneFreeFeePerYear, _ := cache.ApiDaoInstance.GetGlobalParam("khr_phone_free_fee_per_year")
		_, khrPhoneSingleMax, _ := cache.ApiDaoInstance.GetGlobalParam("khr_phone_single_max")
		_, khrPhoneSingleMin, _ := cache.ApiDaoInstance.GetGlobalParam("khr_phone_single_min")

		//if usdPhoneWithdrawRate != "0" || usdPhoneWithdrawRate != "" {
		//	usdPhoneWithdrawRateD, _ := decimal.NewFromString(usdPhoneWithdrawRate)
		//	usdPhoneWithdrawRateD = ss_big.SsBigInst.ToRound(usdPhoneWithdrawRateD.Div(decimal.NewFromInt(10000)), 2, ss_big.RoundingMode_CEILING)
		//	usdPhoneWithdrawRate = usdPhoneWithdrawRateD.String()
		//}
		//if khrPhoneWithdrawRate != "0" || khrPhoneWithdrawRate != "" {
		//	tempD, _ := decimal.NewFromString(khrPhoneWithdrawRate)
		//	tempD = ss_big.SsBigInst.ToRound(tempD.Div(decimal.NewFromInt(10000)), 2, ss_big.RoundingMode_CEILING)
		//	khrPhoneWithdrawRate = tempD.String()
		//}
		reply.PhoneWithdrawDatas = &go_micro_srv_cust.PhoneWithdrawConfigData{
			UsdPhoneWithdrawRate:   usdPhoneWithdrawRate,
			UsdPhoneMinWithdrawFee: usdPhoneMinWithdrawFee,
			UsdPhoneFreeFeePerYear: usdPhoneFreeFeePerYear,
			UsdPhoneSingleMax:      usdPhoneSingleMax,
			UsdPhoneSingleMin:      usdPhoneSingleMin,
			KhrPhoneWithdrawRate:   khrPhoneWithdrawRate,
			KhrPhoneMinWithdrawFee: khrPhoneMinWithdrawFee,
			KhrPhoneFreeFeePerYear: khrPhoneFreeFeePerYear,
			KhrPhoneSingleMax:      khrPhoneSingleMax,
			KhrPhoneSingleMin:      khrPhoneSingleMin,
		}
	case "deposit": //存款配置
		_, khrDepositRate, _ := cache.ApiDaoInstance.GetGlobalParam("khr_deposit_rate")
		_, khrMinDepositFee, _ := cache.ApiDaoInstance.GetGlobalParam("khr_min_deposit_fee")
		_, khr_deposit_single_min, _ := cache.ApiDaoInstance.GetGlobalParam("khr_deposit_single_min")
		_, khr_deposit_single_max, _ := cache.ApiDaoInstance.GetGlobalParam("khr_deposit_single_max")

		_, usd_deposit_rate, _ := cache.ApiDaoInstance.GetGlobalParam("usd_deposit_rate")
		_, usd_min_deposit_fee, _ := cache.ApiDaoInstance.GetGlobalParam("usd_min_deposit_fee")
		_, usd_deposit_single_min, _ := cache.ApiDaoInstance.GetGlobalParam("usd_deposit_single_min")
		_, usd_deposit_single_max, _ := cache.ApiDaoInstance.GetGlobalParam("usd_deposit_single_max")

		//if khrDepositRate != "0" || khrDepositRate != "" {
		//	tempD, _ := decimal.NewFromString(khrDepositRate)
		//	tempD = ss_big.SsBigInst.ToRound(tempD.Div(decimal.NewFromInt(10000)), 2, ss_big.RoundingMode_CEILING)
		//	khrDepositRate = tempD.String()
		//}
		//if usd_deposit_rate != "0" || usd_deposit_rate != "" {
		//	tempD, _ := decimal.NewFromString(usd_deposit_rate)
		//	tempD = ss_big.SsBigInst.ToRound(tempD.Div(decimal.NewFromInt(10000)), 2, ss_big.RoundingMode_CEILING)
		//	usd_deposit_rate = tempD.String()
		//}
		reply.DepositDatas = &go_micro_srv_cust.DepositConfigData{
			KhrDepositRate:      khrDepositRate,
			KhrMinDepositFee:    khrMinDepositFee,
			KhrDepositSingleMin: khr_deposit_single_min,
			KhrDepositSingleMax: khr_deposit_single_max,
			UsdDepositRate:      usd_deposit_rate,
			UsdMinDepositFee:    usd_min_deposit_fee,
			UsdDepositSingleMin: usd_deposit_single_min,
			UsdDepositSingleMax: usd_deposit_single_max,
		}
	case "transfer": //转账配置
		_, khr_transfer_rate, _ := cache.ApiDaoInstance.GetGlobalParam("khr_transfer_rate")
		_, khr_min_transfer_fee, _ := cache.ApiDaoInstance.GetGlobalParam("khr_min_transfer_fee")
		_, khr_transfer_single_min, _ := cache.ApiDaoInstance.GetGlobalParam("khr_transfer_single_min")
		_, khr_transfer_single_max, _ := cache.ApiDaoInstance.GetGlobalParam("khr_transfer_single_max")

		_, usd_transfer_rate, _ := cache.ApiDaoInstance.GetGlobalParam("usd_transfer_rate")
		_, usd_min_transfer_fee, _ := cache.ApiDaoInstance.GetGlobalParam("usd_min_transfer_fee")
		_, usd_transfer_single_min, _ := cache.ApiDaoInstance.GetGlobalParam("usd_transfer_single_min")
		_, usd_transfer_single_max, _ := cache.ApiDaoInstance.GetGlobalParam("usd_transfer_single_max")

		//if khr_transfer_rate != "0" || khr_transfer_rate != "" {
		//	tempD, _ := decimal.NewFromString(khr_transfer_rate)
		//	tempD = ss_big.SsBigInst.ToRound(tempD.Div(decimal.NewFromInt(10000)), 2, ss_big.RoundingMode_CEILING)
		//	khr_transfer_rate = tempD.String()
		//}
		//if usd_transfer_rate != "0" || usd_transfer_rate != "" {
		//	tempD, _ := decimal.NewFromString(usd_transfer_rate)
		//	tempD = ss_big.SsBigInst.ToRound(tempD.Div(decimal.NewFromInt(10000)), 2, ss_big.RoundingMode_CEILING)
		//	usd_transfer_rate = tempD.String()
		//}

		reply.TransferDatas = &go_micro_srv_cust.TransferConfigData{
			KhrTransferRate:      khr_transfer_rate,
			KhrMinTransferFee:    khr_min_transfer_fee,
			KhrTransferSingleMin: khr_transfer_single_min,
			KhrTransferSingleMax: khr_transfer_single_max,
			UsdTransferRate:      usd_transfer_rate,
			UsdMinTransferFee:    usd_min_transfer_fee,
			UsdTransferSingleMin: usd_transfer_single_min,
			UsdTransferSingleMax: usd_transfer_single_max,
		}
	default:
		ss_log.Error("ConfigType参数异常[%v]", req.ConfigType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) UpdateWithdrawConfig(ctx context.Context, req *go_micro_srv_cust.UpdateWithdrawConfigRequest, reply *go_micro_srv_cust.UpdateWithdrawConfigReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	description := ""
	sqlInsert := "insert into global_param(param_key,param_value,remark)values($1,$2,$3) on conflict(param_key) do update set param_value=$2,remark=$3"
	switch req.ConfigType {
	case "face_withdraw":
		l := [][]string{
			{"usd_face_withdraw_rate", req.UsdFaceWithdrawRate, "usd面对面取款费率"},
			{"usd_face_min_withdraw_fee", req.UsdFaceMinWithdrawFee, "usd面对面取款最低收取金额"},
			{"usd_face_free_fee_per_year", req.UsdFaceFreeFeePerYear, "usd面对面取款每年免手续费额度"},
			{"usd_face_single_max", req.UsdFaceSingleMax, "usd面对面取款单笔最大限额"},
			{"usd_face_single_min", req.UsdFaceSingleMin, "usd面对面取款单笔最小限额"},
			{"khr_face_withdraw_rate", req.KhrFaceWithdrawRate, "khr面对面取款费率"},
			{"khr_face_min_withdraw_fee", req.KhrFaceMinWithdrawFee, "khr面对面取款最低收取金额"},
			{"khr_face_free_fee_per_year", req.KhrFaceFreeFeePerYear, "khr面对面取款每年免手续费额度"},
			{"khr_face_single_max", req.KhrFaceSingleMax, "khr面对面取款单笔最大限额"},
			{"khr_face_single_min", req.KhrFaceSingleMin, "khr面对面取款单笔最小限额"},
		}
		description = fmt.Sprintf(
			"面对面取款配置修改为 "+
				"usd取款费率:[%v] "+
				"usd取款最低收取金额:[%v] "+
				"usd每年免手续费额度:[%v] "+
				"usd单笔最大限额:[%v] "+
				"usd单笔最小限额:[%v] "+
				"khr取款费率:[%v] "+
				"khr最低收取金额:[%v] "+
				"khr每年免手续费额度:[%v] "+
				"khr单笔最大限额:[%v] "+
				"khr单笔最小限额:[%v]",

			ss_count.Div(req.UsdFaceWithdrawRate, "100").String()+"%",
			ss_count.Div(req.UsdFaceMinWithdrawFee, "100").String(),
			ss_count.Div(req.UsdFaceFreeFeePerYear, "100").String(),
			ss_count.Div(req.UsdFaceSingleMax, "100").String(),
			ss_count.Div(req.UsdFaceSingleMin, "100").String(),

			ss_count.Div(req.KhrFaceWithdrawRate, "100").String()+"%",
			req.KhrFaceMinWithdrawFee,
			req.KhrFaceFreeFeePerYear,
			req.KhrFaceSingleMax,
			req.KhrFaceSingleMin,
		)
		for _, v := range l {
			err := ss_sql.Exec(dbHandler, sqlInsert, v[0], v[1], v[2])
			if err != nil {
				reply.ResultCode = ss_err.ERR_PARAM
				ss_log.Error("err=[%v]", err)
			}
			// 清缓存
			//cache.RedisCli.Del(constants.DefPoolName, ss_data.MkGlobalParamValue(v[0]))
			cache.RedisClient.Del(ss_data.MkGlobalParamValue(v[0]))
		}
		//description = fmt.Sprintf("修改手机号[%v]", typeStr, req.Text, langStr, statusStr)
	case "phone_withdraw":
		l := [][]string{
			{"usd_phone_withdraw_rate", req.UsdPhoneWithdrawRate, "usd手机号取款费率"},
			{"usd_phone_min_withdraw_fee", req.UsdPhoneMinWithdrawFee, "usd手机号取款最低收取金额"},
			{"usd_phone_free_fee_per_year", req.UsdPhoneFreeFeePerYear, "usd手机号取款每年免手续费额度"},
			{"usd_phone_single_max", req.UsdPhoneSingleMax, "usd手机号取款单笔最大限额"},
			{"usd_phone_single_min", req.UsdPhoneSingleMin, "usd手机号取款单笔最小限额"},
			{"khr_phone_withdraw_rate", req.KhrPhoneWithdrawRate, "khr手机号取款费率"},
			{"khr_phone_min_withdraw_fee", req.KhrPhoneMinWithdrawFee, "khr手机号取款最低收取金额"},
			{"khr_phone_free_fee_per_year", req.KhrPhoneFreeFeePerYear, "khr手机号取款每年免手续费额度"},
			{"khr_phone_single_max", req.KhrPhoneSingleMax, "khr手机号取款单笔最大限额"},
			{"khr_phone_single_min", req.KhrPhoneSingleMin, "khr手机号取款单笔最小限额"},
		}
		description = fmt.Sprintf(
			"修改手机号取款配置 "+
				"usd取款费率:[%v] "+
				"usd取款最低收取金额:[%v] "+
				"usd每年免手续费额度:[%v] "+
				"usd单笔最小限额:[%v] "+
				"usd单笔最大限额:[%v] "+
				"khr取款费率:[%v] "+
				"khr最低收取金额:[%v] "+
				"khr每年免手续费额度:[%v] "+
				"khr单笔最小限额:[%v]"+
				"khr单笔最大限额:[%v] ",
			ss_count.Div(req.UsdPhoneWithdrawRate, "100").String()+"%",
			ss_count.Div(req.UsdPhoneMinWithdrawFee, "100").String(),
			ss_count.Div(req.UsdPhoneFreeFeePerYear, "100").String(),
			ss_count.Div(req.UsdPhoneSingleMin, "100").String(),
			ss_count.Div(req.UsdPhoneSingleMax, "100").String(),

			ss_count.Div(req.KhrPhoneWithdrawRate, "100").String()+"%",
			req.KhrPhoneMinWithdrawFee,
			req.KhrPhoneFreeFeePerYear,
			req.KhrPhoneSingleMin,
			req.KhrPhoneSingleMax,
		)
		for _, v := range l {
			err := ss_sql.Exec(dbHandler, sqlInsert, v[0], v[1], v[2])
			if err != nil {
				reply.ResultCode = ss_err.ERR_PARAM
				ss_log.Error("err=[%v]", err)
			}
			// 清缓存
			//cache.RedisCli.Del(constants.DefPoolName, ss_data.MkGlobalParamValue(v[0]))
			cache.RedisClient.Del(ss_data.MkGlobalParamValue(v[0]))
		}
	case "deposit":
		l := [][]string{
			{"usd_deposit_rate", req.UsdDepositRate, "usd存款手续费率"},
			{"usd_min_deposit_fee", req.UsdMinDepositFee, "usd存款最低收取金额"},
			{"usd_deposit_single_max", req.UsdDepositSingleMax, "usd存款单笔最大限额"},
			{"usd_deposit_single_min", req.UsdDepositSingleMin, "usd存款单笔最小限额"},
			{"khr_deposit_rate", req.KhrDepositRate, "khr存款手续费率"},
			{"khr_min_deposit_fee", req.KhrMinDepositFee, "khr存款最低收取金额"},
			{"khr_deposit_single_max", req.KhrDepositSingleMax, "khr存款单笔最大限额"},
			{"khr_deposit_single_min", req.KhrDepositSingleMin, "khr存款单笔最小限额"},
		}

		description = fmt.Sprintf(
			"存款配置修改为 "+
				"usd存款手续费率:[%v] "+
				"usd存款最低收取金额:[%v] "+
				"usd存款单笔最小限额:[%v] "+
				"usd存款单笔最大限额:[%v] "+
				"khr存款手续费率:[%v] "+
				"khr存款最低收取金额:[%v] "+
				"khr存款单笔最小限额:[%v] "+
				"khr存款单笔最大限额:[%v] ",

			ss_count.Div(req.UsdDepositRate, "100").String()+"%",
			ss_count.Div(req.UsdMinDepositFee, "100").String(),
			ss_count.Div(req.UsdDepositSingleMin, "100").String(),
			ss_count.Div(req.UsdDepositSingleMax, "100").String(),

			ss_count.Div(req.KhrDepositRate, "100").String()+"%",
			req.KhrMinDepositFee,
			req.KhrDepositSingleMin,
			req.KhrDepositSingleMax,
		)
		for _, v := range l {
			err := ss_sql.Exec(dbHandler, sqlInsert, v[0], v[1], v[2])
			if err != nil {
				reply.ResultCode = ss_err.ERR_PARAM
				ss_log.Error("err=[%v]", err)
			}
			// 清缓存
			//cache.RedisCli.Del(constants.DefPoolName, ss_data.MkGlobalParamValue(v[0]))
			cache.RedisClient.Del(ss_data.MkGlobalParamValue(v[0]))
		}
	case "transfer":
		l := [][]string{
			{"usd_transfer_rate", req.UsdTransferRate, "usd转账手续费率"},
			{"usd_min_transfer_fee", req.UsdMinTransferFee, "usd转账最低收取金额"},
			{"usd_transfer_single_max", req.UsdTransferSingleMax, "usd转账单笔最大限额"},
			{"usd_transfer_single_min", req.UsdTransferSingleMin, "usd转账单笔最小限额"},
			{"khr_transfer_rate", req.KhrTransferRate, "khr转账手续费率"},
			{"khr_min_transfer_fee", req.KhrMinTransferFee, "khr转账最低收取金额"},
			{"khr_transfer_single_max", req.KhrTransferSingleMax, "khr转账单笔最大限额"},
			{"khr_transfer_single_min", req.KhrTransferSingleMin, "khr转账单笔最小限额"},
		}

		description = fmt.Sprintf(
			"转账配置修改为 "+
				"usd转账手续费率:[%v] "+
				"usd转账最低收取金额:[%v] "+
				"usd转账单笔最小限额:[%v] "+
				"usd转账单笔最大限额:[%v] "+
				"khr转账手续费率:[%v] "+
				"khr转账最低收取金额:[%v] "+
				"khr转账单笔最小限额:[%v] "+
				"khr转账单笔最大限额:[%v] ",

			ss_count.Div(req.UsdTransferRate, "100").String()+"%",
			ss_count.Div(req.UsdMinTransferFee, "100").String(),
			ss_count.Div(req.UsdTransferSingleMin, "100").String(),
			ss_count.Div(req.UsdTransferSingleMax, "100").String(),

			ss_count.Div(req.KhrTransferRate, "100").String()+"%",
			req.KhrMinTransferFee,
			req.KhrTransferSingleMin,
			req.KhrTransferSingleMax,
		)
		for _, v := range l {
			err := ss_sql.Exec(dbHandler, sqlInsert, v[0], v[1], v[2])
			if err != nil {
				reply.ResultCode = ss_err.ERR_PARAM
				ss_log.Error("err=[%v]", err)
			}
			// 清缓存
			//cache.RedisCli.Del(constants.DefPoolName, ss_data.MkGlobalParamValue(v[0]))
			cache.RedisClient.Del(ss_data.MkGlobalParamValue(v[0]))
		}
	default:
		ss_log.Error("ConfigType参数异常")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetExchangeRateConfig(ctx context.Context, req *go_micro_srv_cust.GetExchangeRateConfigRequest, reply *go_micro_srv_cust.GetExchangeRateConfigReply) error {
	_, usdToKhr, _ := cache.ApiDaoInstance.GetGlobalParam("usd_to_khr")
	_, khrToUsd, _ := cache.ApiDaoInstance.GetGlobalParam("khr_to_usd")

	_, usdToKhrFee, _ := cache.ApiDaoInstance.GetGlobalParam("usd_to_khr_fee")
	_, khrToUsdFee, _ := cache.ApiDaoInstance.GetGlobalParam("khr_to_usd_fee")

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = &go_micro_srv_cust.ExchangeRateConfigData{
		UsdToKhr:    usdToKhr,
		KhrToUsd:    khrToUsd,
		UsdToKhrFee: usdToKhrFee,
		KhrToUsdFee: khrToUsdFee,
	}
	return nil
}

func (*CustHandler) UpdateExchangeRateConfig(ctx context.Context, req *go_micro_srv_cust.UpdateExchangeRateConfigRequest, reply *go_micro_srv_cust.UpdateExchangeRateConfigReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlInsert := "insert into global_param(param_key,param_value,remark)values($1,$2,$3) on conflict(param_key) do update set param_value=$2,remark=$3"
	l := [][]string{
		{"usd_to_khr", req.UsdToKhr, "USD兑换KHR"},
		{"khr_to_usd", req.KhrToUsd, "KHR兑换USD"},

		{"usd_to_khr_fee", req.UsdToKhrFee, "USD兑换KHR单笔手续费"},
		{"khr_to_usd_fee", req.KhrToUsdFee, "KHR兑换USD单笔手续费"},
	}

	description := fmt.Sprintf(
		"兑换配置修改为 "+
			"USD兑换KHR:[1USD=%vKHR] "+
			"USD兑换KHR单笔手续费:[%v] "+
			"KHR兑换USD:[%vKHR=1USD] "+
			"KHR兑换USD单笔手续费:[%v] ",

		req.UsdToKhr,
		ss_count.Div(req.UsdToKhrFee, "100").String(),

		req.KhrToUsd,
		req.KhrToUsdFee,
	)

	for _, v := range l {
		err := ss_sql.Exec(dbHandler, sqlInsert, v[0], v[1], v[2])
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
		// 清缓存
		//cache.RedisCli.Del(constants.DefPoolName, ss_data.MkGlobalParamValue(v[0]))
		cache.RedisClient.Del(ss_data.MkGlobalParamValue(v[0]))
	}

	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetBusinessConfig(ctx context.Context, req *go_micro_srv_cust.GetBusinessConfigRequest, reply *go_micro_srv_cust.GetBusinessConfigReply) error {
	_, usdAmountLimit, _ := cache.ApiDaoInstance.GetGlobalParam(constants.GlobalParamKeyBusinessTransferUSDAmount)
	stringArr1 := strings.Split(usdAmountLimit, "/")
	usdAmountMinLimit := stringArr1[0]
	usdAmountMaxLimit := stringArr1[1]

	_, khrAmountLimit, _ := cache.ApiDaoInstance.GetGlobalParam(constants.GlobalParamKeyBusinessTransferKHRAmount)
	stringArr2 := strings.Split(khrAmountLimit, "/")
	khrAmountMinLimit := stringArr2[0]
	khrAmountMaxLimit := stringArr2[1]

	_, usdRate, _ := cache.ApiDaoInstance.GetGlobalParam(constants.GlobalParamKeyBusinessTransferUSDRate)
	_, khrRate, _ := cache.ApiDaoInstance.GetGlobalParam(constants.GlobalParamKeyBusinessTransferKHRRate)
	_, usdMinFee, _ := cache.ApiDaoInstance.GetGlobalParam(constants.GlobalParamKeyBusinessTransferUSDMinFee)
	_, khrMinFee, _ := cache.ApiDaoInstance.GetGlobalParam(constants.GlobalParamKeyBusinessTransferKHRMinFee)

	_, batchNumber, _ := cache.ApiDaoInstance.GetGlobalParam(constants.GlobalParamKeyBusinessTransferBatchNum)

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.TransferConfigData = &go_micro_srv_cust.BusinessTransferConfigData{
		UsdAmountMinLimit: usdAmountMinLimit,
		UsdAmountMaxLimit: usdAmountMaxLimit,
		KhrAmountMinLimit: khrAmountMinLimit,
		KhrAmountMaxLimit: khrAmountMaxLimit,
		UsdRate:           usdRate,
		KhrRate:           khrRate,
		BatchNumber:       batchNumber,
		UsdMinFee:         usdMinFee,
		KhrMinFee:         khrMinFee,
	}
	return nil
}

func (*CustHandler) UpdateBusinessConfig(ctx context.Context, req *go_micro_srv_cust.UpdateBusinessConfigRequest, reply *go_micro_srv_cust.UpdateBusinessConfigReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlInsert := "insert into global_param(param_key,param_value)values($1,$2) on conflict(param_key) do update set param_value=$2"

	usdAmountLimit := req.UsdAmountMinLimit + "/" + req.UsdAmountMaxLimit
	khrAmountLimit := req.KhrAmountMinLimit + "/" + req.KhrAmountMaxLimit

	l := [][]string{
		{constants.GlobalParamKeyBusinessTransferUSDAmount, usdAmountLimit}, //"商家转账USD最小/最大转账金额限制"
		{constants.GlobalParamKeyBusinessTransferKHRAmount, khrAmountLimit}, // "商家转账KHR最小/最大转账金额限制"

		{constants.GlobalParamKeyBusinessTransferUSDRate, req.UsdRate}, // "商家转账USD手续费比率(万分比, 四舍五入)"
		{constants.GlobalParamKeyBusinessTransferKHRRate, req.KhrRate}, // "商家转账KHR手续费比率(万分比, 四舍五入)"

		{constants.GlobalParamKeyBusinessTransferUSDMinFee, req.UsdMinFee}, // "商家转账最低收取usd手续费"
		{constants.GlobalParamKeyBusinessTransferKHRMinFee, req.KhrMinFee}, // "商家转账最低收取khr手续费"

		{constants.GlobalParamKeyBusinessTransferBatchNum, req.BatchNumber}, // "商家转账批量付款最大总人数(包括USD和KHR)"
	}

	description := fmt.Sprintf(
		"商家付款配置修改为 "+
			"商家转账USD最小/最大转账金额限制:[%v] "+
			"商家转账KHR最小/最大转账金额限制:[%v] "+
			"商家转账USD手续费比率:[%v] "+
			"商家转账KHR手续费比率:[%v] "+
			"商家转账最低收取usd手续费:[%v] "+
			"商家转账最低收取khr手续费:[%v] "+
			"商家转账批量付款最大总人数(包括USD和KHR):[%v] ",
		ss_count.Div(usdAmountLimit, "100").String(),
		khrAmountLimit,
		ss_count.Div(req.UsdRate, "100").String(),
		ss_count.Div(req.KhrRate, "100").String(),
		ss_count.Div(req.UsdMinFee, "100").String(),
		ss_count.Div(req.KhrMinFee, "100").String(),
		req.BatchNumber,
	)

	for _, v := range l {
		err := ss_sql.Exec(dbHandler, sqlInsert, v[0], v[1])
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
		// 清缓存
		//cache.RedisCli.Del(constants.DefPoolName, ss_data.MkGlobalParamValue(v[0]))
		cache.RedisClient.Del(ss_data.MkGlobalParamValue(v[0]))
	}

	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetFuncConfig(ctx context.Context, req *go_micro_srv_cust.GetFuncConfigRequest, reply *go_micro_srv_cust.GetFuncConfigReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
		{Key: "application_type", Val: req.ApplicationType, EqType: "="},
		{Key: "func_name", Val: req.FuncName, EqType: "like"},
	})

	sqlCnt := " select count(1) from func_config " + whereModel.WhereStr
	var total sql.NullString
	errCnt := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...)
	if errCnt != nil {
		ss_log.Error("err=[%v]", errCnt)
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by idx `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT func_no, func_name, idx, use_status, img, jump_url, application_type " +
		" FROM func_config " + where2
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var datas []*go_micro_srv_cust.FuncConfigData

	imageBaseUrl := dao.GlobalParamDaoInstance.QeuryParamValue("image_base_url")

	if err == nil {
		for rows.Next() {
			var data go_micro_srv_cust.FuncConfigData
			err = rows.Scan(
				&data.FuncNo,
				&data.FuncName,
				&data.Idx,
				&data.UseStatus,
				&data.Img,
				&data.JumpUrl,
				&data.ApplicationType,
			)
			if err != nil {
				ss_log.Error("err=[%v]", err)
				continue
			}

			data.Img = imageBaseUrl + "/" + data.Img

			datas = append(datas, &data)
		}
	} else {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

func (c *CustHandler) UpdateFuncConfig(ctx context.Context, req *go_micro_srv_cust.UpdateFuncConfigRequest, reply *go_micro_srv_cust.UpdateFuncConfigReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	imgId := ""

	//根据上传的图片base64字符串保存图片
	upReq := &go_micro_srv_cust.UploadImageRequest{
		ImageStr:   req.ImgBase64,
		AccountUid: "00000000-0000-0000-0000-000000000000",
		Type:       constants.UploadImage_UnAuth, //不需要授权的
	}
	upReply := &go_micro_srv_cust.UploadImageReply{}
	c.UploadImage(ctx, upReq, upReply)
	if upReply.ResultCode != ss_err.ERR_SUCCESS {
		ss_log.Error("保存图片失败")
		reply.ResultCode = ss_err.ERR_SAVE_IMAGE_FAILD
		return nil
	}

	imgId = upReply.ImageId

	imgDaoData, errGet := dao.ImageDaoInstance.GetImageUrlById(imgId)
	if errGet != nil || imgDaoData.ImageId == "" {
		ss_log.Error("获取图片ImageUrl失败，id[%v]", imgId)
		reply.ResultCode = ss_err.ERR_SAVE_IMAGE_FAILD
		return nil
	}

	if req.FuncNo == "" {
		// 添加
		errCode := dao.FuncConfigDaoInst.AddFuncConfig(req.FuncName, req.JumpUrl, imgDaoData.ImageUrl, imgId, req.ApplicationType)
		if errCode != ss_err.ERR_SUCCESS {
			ss_log.Error("err=[%v]", errCode)
			reply.ResultCode = errCode
			return nil
		}
	} else {

		imageURL, err := dao.FuncConfigDaoInst.GetImgURLFromNo(req.FuncNo)
		if err != nil {
			ss_log.Error("GetImgURLFromNo 查询图片id失败,req.FuncNo: %s,err: %s", req.FuncNo, err.Error())
			return nil
		}
		// 删除s3 上面的图片
		if _, err := common.UploadS3.DeleteOne(imageURL); err != nil {
			notes := fmt.Sprintf("UpdateFuncConfig s3上删除图片失败,图片路劲为: %+v,err: %s", imageURL, err.Error())
			ss_log.Error(notes)
			dao.DictimagesDaoInst.AddDelFaildLog(notes)
		}
		// 删除数据库的图片
		if err := dao.DictimagesDaoInst.Delete(imageURL); err != nil {
			notes := fmt.Sprintf("UpdateFuncConfig 删除图片记录失败,图片路劲为: %+v,err: %s", imageURL, err.Error())
			ss_log.Error(notes)
			dao.DictimagesDaoInst.AddDelFaildLog(notes)
		}

		// 修改
		errCode := dao.FuncConfigDaoInst.UpdateFuncConfig(tx, req.FuncNo, req.FuncName, req.JumpUrl, imgDaoData.ImageUrl, imgId)
		if errCode != ss_err.ERR_SUCCESS {
			ss_log.Error("err=[%v]", errCode)
			reply.ResultCode = errCode
			return nil
		}
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeleteFuncConfig(ctx context.Context, req *go_micro_srv_cust.DeleteFuncConfigRequest, reply *go_micro_srv_cust.DeleteFuncConfigReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	id := req.FuncNo

	//获取当前要删除记录的idx
	startIdx, applicationType, errGet1 := dao.FuncConfigDaoInst.GetIdxAndApplicationTypeById(tx, id)
	if errGet1 != ss_err.ERR_SUCCESS {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	endIdx, errGet1 := dao.FuncConfigDaoInst.GetMaxidx(tx, applicationType)
	if errGet1 != ss_err.ERR_SUCCESS {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//将其要删除的排序序号后面的元素往前移
	for i := startIdx + 1; i <= endIdx; i++ {
		errUp := dao.FuncConfigDaoInst.ReplaceIdx(tx, i, applicationType)
		if errUp != ss_err.ERR_SUCCESS {
			ss_log.Error("errUp=[%v],i=[%v]", errUp, i)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}
	imageURL, err := dao.FuncConfigDaoInst.GetImgURLFromNo(req.FuncNo)
	if err != nil {
		ss_log.Error("查询funcConfig的图片url失败,err: %s", err.Error())
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	errCode := dao.FuncConfigDaoInst.DeleteFuncConfig(tx, req.FuncNo)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", errCode)
		reply.ResultCode = errCode
		return nil
	}

	tx.Commit()
	// 从s3上删除图片
	if imageURL != "" {
		if _, err := common.UploadS3.DeleteOne(imageURL); err != nil {
			notes := fmt.Sprintf("DeleteFuncConfig s3上删除图片失败,图片路劲为: %s,err: %s", imageURL, err.Error())
			ss_log.Error(notes)
			dao.DictimagesDaoInst.AddDelFaildLog(notes)
		}

		// 删除数据库的图片
		if err := dao.DictimagesDaoInst.Delete(imageURL); err != nil {
			notes := fmt.Sprintf("DeleteFuncConfig 删除图片记录失败,图片路劲为: %s,err: %s", imageURL, err.Error())
			ss_log.Error(notes)
			dao.DictimagesDaoInst.AddDelFaildLog(notes)
		}
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) ModifyUseStatusFuncConfig(ctx context.Context, req *go_micro_srv_cust.ModifyUseStatusFuncConfigRequest, reply *go_micro_srv_cust.ModifyUseStatusFuncConfigReply) error {
	errCode := dao.FuncConfigDaoInst.ModifyUseStatusFuncConfig(req.FuncNo, req.UseStatus)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", errCode)
		reply.ResultCode = errCode
		return nil
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetTransferSecurityConfig(ctx context.Context, req *go_micro_srv_cust.GetTransferSecurityConfigRequest, reply *go_micro_srv_cust.GetTransferSecurityConfigReply) error {
	_, payPwdErrCount, _ := cache.ApiDaoInstance.GetGlobalParam(constants.GlobalParamKeyPaymentPwdErrCount)
	_, continuousErrPassword, _ := cache.ApiDaoInstance.GetGlobalParam(constants.GlobalParamKeyLoginPwdErrCount)
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = &go_micro_srv_cust.TransferSecurityConfigData{
		ErrPaymentPwdCount:    payPwdErrCount,
		ContinuousErrPassword: continuousErrPassword,
	}
	return nil
}

func (*CustHandler) UpdateTransferSecurityConfig(ctx context.Context, req *go_micro_srv_cust.UpdateTransferSecurityConfigRequest, reply *go_micro_srv_cust.UpdateTransferSecurityConfigReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlInsert := "insert into global_param(param_key,param_value,remark)values($1,$2,$3) on conflict(param_key) do update set param_value=$2,remark=$3"
	l := [][]string{
		{constants.GlobalParamKeyLoginPwdErrCount, req.ContinuousErrPassword, "连续输错登录密码次数"},
		{constants.GlobalParamKeyPaymentPwdErrCount, req.ErrPaymentPwdCount, "连续输错支付密码次数"},
		//{"defrozen_trading_rights_time", req.DefrozenTradingRightsTime, "冻结交易权限时间"},
	}

	description := fmt.Sprintf(
		"交易安全配置修改为 "+
			"连续输错登录密码次数:[%v] "+
			"连续输错支付密码次数:[%v] ",
		//"冻结交易权限时间:[%v] ",

		req.ContinuousErrPassword,
		req.ErrPaymentPwdCount,
	)

	for _, v := range l {
		err := ss_sql.Exec(dbHandler, sqlInsert, v[0], v[1], v[2])
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
		// 清缓存
		//cache.RedisCli.Del(constants.DefPoolName, ss_data.MkGlobalParamValue(v[0]))
		cache.RedisClient.Del(ss_data.MkGlobalParamValue(v[0]))
	}

	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetWriteOffDurationDateConfig(ctx context.Context, req *go_micro_srv_cust.GetWriteOffDurationDateConfigRequest, reply *go_micro_srv_cust.GetWriteOffDurationDateConfigReply) error {
	_, writeOffDurationDate, _ := cache.ApiDaoInstance.GetGlobalParam("write_off_duration_date")
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DurationDate = writeOffDurationDate
	return nil
}

func (*CustHandler) UpdateWriteOffDurationDateConfig(ctx context.Context, req *go_micro_srv_cust.UpdateWriteOffDurationDateConfigRequest, reply *go_micro_srv_cust.UpdateWriteOffDurationDateConfigReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlInsert := "insert into global_param(param_key,param_value,remark)values($1,$2,$3) on conflict(param_key) do update set param_value=$2,remark=$3"
	if err := ss_sql.Exec(dbHandler, sqlInsert, constants.KEY_WriteOff_DurationDate, req.DurationDate, "核销码有效期(单位: 天)"); err != nil {
		ss_log.Error("修改核销码有效期失败，err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}
	cache.RedisClient.Del(ss_data.MkGlobalParamValue(constants.KEY_WriteOff_DurationDate))

	description := fmt.Sprintf("核销码有效期修改为 [%v] ", req.DurationDate)
	if errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config); errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetIncomeOugoConfig(ctx context.Context, req *go_micro_srv_cust.GetIncomeOugoConfigRequest, reply *go_micro_srv_cust.GetIncomeOugoConfigReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var datas []*go_micro_srv_cust.IncomeOugoConfigData

	switch req.ConfigType {
	case "1": //类型（1.充值方式。2.提现方式）
	case "2":
	default:
		ss_log.Error("ConfigType is no in (1,2)")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
		{Key: "config_type", Val: req.ConfigType, EqType: "="},
		{Key: "name", Val: req.Name, EqType: "like"},
		{Key: "currency_type", Val: req.CurrencyType, EqType: "="},
	})

	sqlCnt := "SELECT count(1) " +
		" FROM income_ougo_config " + whereModel.WhereStr
	var total sql.NullString
	errCnt := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...)
	if errCnt != nil {
		ss_log.Error("errCnt=[%v]", errCnt)
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by idx `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT income_ougo_config_no, currency_type, name, use_status, idx, config_type " +
		" FROM income_ougo_config " + where2
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err == nil {
		for rows.Next() {
			var data go_micro_srv_cust.IncomeOugoConfigData
			err = rows.Scan(
				&data.IncomeOugoConfigNo,
				&data.CurrencyType,
				&data.Name,
				&data.UseStatus,
				&data.Idx,
				&data.ConfigType,
			)
			if err != nil {
				ss_log.Error("err=[%v]", err)
				continue
			}
			datas = append(datas, &data)
		}
	} else {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

// 交换位置
func (c *CustHandler) SwapFuncConfigIdx(ctx context.Context, req *go_micro_srv_cust.SwapFuncConfigIdxRequest, reply *go_micro_srv_cust.SwapFuncConfigIdxReply) error {
	switch req.AppType {
	case constants.AppType_App:
	case constants.AppType_Pos:
	default:
		ss_log.Error("appType错误")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	switch req.SwapType {
	case constants.SwapType_Up: // up
	case constants.SwapType_Down: // down
	default:
		ss_log.Error("swapType错误")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.Idx == "1" && req.SwapType == constants.SwapType_Up {
		ss_log.Error("swapType错误")
		reply.ResultCode = ss_err.ERR_MERC_IS_UPTOP
		return nil
	}

	funcNoFrom := dao.FuncConfigDaoInst.GetFuncNo(req.AppType, req.Idx)
	if funcNoFrom == "" {
		ss_log.Error("获取funcNo失败")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	funcNoTo := dao.FuncConfigDaoInst.GetNearIdxFuncNo(req.AppType, req.Idx, req.SwapType)
	if funcNoTo == "" {
		ss_log.Error("获取funcNo失败")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	switch req.SwapType {
	case constants.SwapType_Up:
		// 向上需要交换一下
		x := funcNoTo
		funcNoTo = funcNoFrom
		funcNoFrom = x
	}

	errCode := dao.FuncConfigDaoInst.ExchangeIdx(funcNoFrom, funcNoTo)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", errCode)
		reply.ResultCode = errCode
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) UpdateIncomeOugoConfig(ctx context.Context, req *go_micro_srv_cust.UpdateIncomeOugoConfigRequest, reply *go_micro_srv_cust.UpdateIncomeOugoConfigReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	switch req.CurrencyType {
	case "":
	case "usd":
	case "khr":
	default:
		ss_log.Error("UpdateIncomeOugoConfig |  err = CurrencyType is no in (usd,khr)")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	switch req.ConfigType {
	case "":
	case "1":
	case "2":
	default:
		ss_log.Error("ConfigType is no in (1,2)")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.IncomeOugoConfigNo == "" { //插入
		var maxIdx sql.NullString //查询最大的idx数
		sqlStr := "select idx from income_ougo_config where config_type = $1 order by idx desc limit 1 "
		selectErr := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&maxIdx}, req.ConfigType)
		if selectErr != nil {
			ss_log.Error("selectErr=[%v]", selectErr)
			reply.ResultCode = ss_err.ERR_PARAM
		}

		var maxIdx2 int64
		if maxIdx.String != "" {
			maxIdx2 = strext.ToInt64(maxIdx.String) + 1
		} else {
			maxIdx2 = 1
		}

		sqlInsert := "insert into income_ougo_config(income_ougo_config_no,currency_type, name, use_status, idx, config_type) values($1,$2,$3,$4,$5,$6)"
		insertErr := ss_sql.Exec(dbHandler, sqlInsert, strext.NewUUID(), req.CurrencyType, req.Name, "1", maxIdx2, req.ConfigType)
		if insertErr != nil {
			ss_log.Error("insertErr=[%v]", insertErr)
			reply.ResultCode = ss_err.ERR_PARAM
		}
	} else { //修改
		sqlInsert := "update income_ougo_config set currency_type=$2,name=$3 where income_ougo_config_no = $1 and config_type = $4"
		updateErr := ss_sql.Exec(dbHandler, sqlInsert, req.IncomeOugoConfigNo, req.CurrencyType, req.Name, req.ConfigType)
		if updateErr != nil {
			ss_log.Error("updateErr=[%v]", updateErr)
			reply.ResultCode = ss_err.ERR_PARAM
		}

	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) UpdateIncomeOugoConfigUseStatus(ctx context.Context, req *go_micro_srv_cust.UpdateIncomeOugoConfigUseStatusRequest, reply *go_micro_srv_cust.UpdateIncomeOugoConfigUseStatusReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	switch req.UseStatus {
	case "1":
	case "0":
	default:
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	sqlUpdate := "update income_ougo_config set use_status=$1 where income_ougo_config_no = $2 and is_delete = '0' "
	err := ss_sql.Exec(dbHandler, sqlUpdate, req.UseStatus, req.IncomeOugoConfigNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) UpdateIncomeOugoConfigIdx(ctx context.Context, req *go_micro_srv_cust.UpdateIncomeOugoConfigIdxRequest, reply *go_micro_srv_cust.UpdateIncomeOugoConfigIdxReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	switch req.ConfigType {
	case "1": //1.充值方式。2.提现方式
	case "2":
	default:
		ss_log.Error("config_type错误 [%v]", req.ConfigType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	switch req.SwapType {
	case constants.SwapType_Up: // up
	case constants.SwapType_Down: // down
	default:
		ss_log.Error("swapType错误 [%v]", req.SwapType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.Idx == "1" && req.SwapType == constants.SwapType_Up {
		ss_log.Error("swapType错误")
		reply.ResultCode = ss_err.ERR_MERC_IS_UPTOP
		return nil
	}

	maxIdx, errGet1 := dao.IncomeOutgoConfigDaoInst.GetMaxidx(tx, req.ConfigType)
	if errGet1 != ss_err.ERR_SUCCESS {
		ss_log.Error("查询最大idx出错")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.Idx == strext.ToStringNoPoint(maxIdx) && req.SwapType == constants.SwapType_Down {
		ss_log.Error("swapType错误")
		reply.ResultCode = ss_err.ERR_MERC_IS_DOWN
		return nil
	}

	incomeOutgoNoFrom := dao.IncomeOutgoConfigDaoInst.GetIncomeOutgoConfigId(tx, req.ConfigType, req.Idx)
	if incomeOutgoNoFrom == "" {
		ss_log.Error("获取incomeOutgoNo失败")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	incomeOutgoNoTo := dao.IncomeOutgoConfigDaoInst.GetNearIdxIncomeOugoNo(tx, req.ConfigType, req.Idx, req.SwapType)
	if incomeOutgoNoTo == "" {
		ss_log.Error("获取incomeOutgoNoTo失败")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	switch req.SwapType {
	case constants.SwapType_Up:
		// 向上需要交换一下
		x := incomeOutgoNoTo
		incomeOutgoNoTo = incomeOutgoNoFrom
		incomeOutgoNoFrom = x
	}

	errCode := dao.IncomeOutgoConfigDaoInst.ExchangeIdx(tx, incomeOutgoNoFrom, incomeOutgoNoTo)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", errCode)
		reply.ResultCode = errCode
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeleteIncomeOugoConfig(ctx context.Context, req *go_micro_srv_cust.DeleteIncomeOugoConfigRequest, reply *go_micro_srv_cust.DeleteIncomeOugoConfigReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	id := req.IncomeOugoConfigNo

	//获取当前要删除记录的idx
	startIdx, configType, errGet1 := dao.IncomeOutgoConfigDaoInst.GetIdxAndConfigTypeById(tx, id)
	if errGet1 != ss_err.ERR_SUCCESS {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	endIdx, errGet1 := dao.IncomeOutgoConfigDaoInst.GetMaxidx(tx, configType)
	if errGet1 != ss_err.ERR_SUCCESS {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//将其要删除的排序序号后面的元素往前移
	for i := startIdx + 1; i <= endIdx; i++ {
		errUp := dao.IncomeOutgoConfigDaoInst.ReplaceIdx(tx, i, configType)
		if errUp != ss_err.ERR_SUCCESS {
			ss_log.Error("errUp=[%v],i=[%v]", errUp, i)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	err := ss_sql.ExecTx(tx, `update income_ougo_config set is_delete='1', idx='-1' where income_ougo_config_no = $1 and is_delete='0' `, id)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	tx.Commit()

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (c *CustHandler) GetFuncConfigDetail(ctx context.Context, req *go_micro_srv_cust.GetFuncConfigDetailRequest, reply *go_micro_srv_cust.GetFuncConfigDetailReply) error {
	//data := dao.FuncConfigDaoInst.GetFuncData(req.FuncNo)
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select func_no, func_name, img, img_id, jump_url, idx, use_status, application_type " +
		" from func_config " +
		" where func_no=$1 and is_delete='0' limit 1"
	row, stmt, err := ss_sql.QueryRowN(dbHandler, sqlStr, req.FuncNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil
	}
	defer stmt.Close()

	data := &go_micro_srv_cust.FuncConfigData{}
	err = row.Scan(
		&data.FuncNo,
		&data.FuncName,
		&data.Img,
		&data.ImgId,
		&data.JumpUrl,
		&data.Idx,
		&data.UseStatus,
		&data.ApplicationType,
	)
	imageBaseUrl := dao.GlobalParamDaoInstance.QeuryParamValue("image_base_url")
	data.Img = imageBaseUrl + "/" + data.Img

	if data.ImgId != "" {
		reqImg := &go_micro_srv_cust.UnAuthDownloadImageBase64Request{
			ImageId: data.ImgId,
		}
		replyImg := &go_micro_srv_cust.UnAuthDownloadImageBase64Reply{}
		c.UnAuthDownloadImageBase64(ctx, reqImg, replyImg)
		if replyImg.ResultCode != ss_err.ERR_SUCCESS {
			ss_log.Error("获取图片base64字符串失败")
		} else {
			data.ImgBase64 = replyImg.ImageBase64
		}
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

func (c *CustHandler) GetLangs(ctx context.Context, req *go_micro_srv_cust.GetLangsRequest, reply *go_micro_srv_cust.GetLangsReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
		{Key: "type", Val: req.Type, EqType: "="},
		{Key: "key", Val: req.Key, EqType: "like"},
		{Key: "lang_ch", Val: req.LangCh, EqType: "like"},
		{Key: "lang_en", Val: req.LangEn, EqType: "like"},
		{Key: "lang_km", Val: req.LangKm, EqType: "like"},
		{Key: "key", Val: req.SearchKey, EqType: "like"},
	})
	var total sql.NullString
	sqlCnt := "SELECT count(1) FROM lang " + whereModel.WhereStr
	errCnt := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...)
	if errCnt != nil {
		ss_log.Error("errCnt=[%v]", errCnt)
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by key asc `)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT key, type, lang_km, lang_en, lang_ch " +
		" FROM lang " + where2
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var datas []*go_micro_srv_cust.LangData

	if err == nil {
		for rows.Next() {
			var data go_micro_srv_cust.LangData
			var langKm, langEn, langCh sql.NullString
			err = rows.Scan(
				&data.Key,
				&data.Type,
				&langKm,
				&langEn,
				&langCh,
			)
			if err != nil {
				ss_log.Error("err=[%v]", err)
				continue
			}
			if data.Type == constants.LANG_TYPE_IMG {
				langImgIds := []string{langKm.String, langEn.String, langCh.String}
				var imgUrls []string
				for _, langImgId := range langImgIds {
					reqImg := &go_micro_srv_cust.UnAuthDownloadImageRequest{
						ImageId: langImgId,
					}
					replyImg := &go_micro_srv_cust.UnAuthDownloadImageReply{}
					c.UnAuthDownloadImage(ctx, reqImg, replyImg)
					if replyImg.ResultCode != ss_err.ERR_SUCCESS {
						ss_log.Error("获取图片url失败")
					}
					imgUrls = append(imgUrls, replyImg.ImageUrl)
				}

				data.LangKm = imgUrls[0]
				data.LangEn = imgUrls[1]
				data.LangCh = imgUrls[2]
			} else {
				data.LangKm = langKm.String
				data.LangEn = langEn.String
				data.LangCh = langCh.String
			}

			datas = append(datas, &data)
		}
	} else {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.DataList = datas
	reply.Total = strext.ToInt32(total.String)
	return nil
}

func (c *CustHandler) UpdateLang(ctx context.Context, req *go_micro_srv_cust.UpdateLangRequest, reply *go_micro_srv_cust.UpdateLangReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	if req.Type == constants.LANG_TYPE_IMG && req.Key != "" {
		// 删除图片
		langKm, langEn, langCh, err := dao.LangDaoInst.GetLangByKey(req.Key, req.Type)
		if err != nil {
			ss_log.Error("查询多语言失败,req.Key: %s,req.Type: %s,err: %s", req.Key, req.Type, err.Error())
		}

		var imgPaths []string
		// 找到对应的图片路径
		if req.LangKm != "" && langKm != "" {
			path, err := dao.DictimagesDaoInst.GetKmImgPath(langKm)
			if err != nil {
				ss_log.Error("GetKmImgPath 失败,langKm: %s,err: %s", langKm, err.Error())
			}
			if path != "" {
				imgPaths = append(imgPaths, path)
			}
		}

		if req.LangEn != "" && langEn != "" {
			path, err := dao.DictimagesDaoInst.GetEnImgPath(langEn)
			if err != nil {
				ss_log.Error("GetEnImgPath 失败,langEn: %s,err: %s", langKm, err.Error())
			}
			if path != "" {
				imgPaths = append(imgPaths, path)
			}
		}
		if req.LangCh != "" && langCh != "" {
			path, err := dao.DictimagesDaoInst.GetChImgPath(langCh)
			if err != nil {
				ss_log.Error("GetChImgPath 失败,langCh: %s,err: %s", langKm, err.Error())
			}
			if path != "" {
				imgPaths = append(imgPaths, path)
			}
		}
		ss_log.Info("s3上图片的路劲集合为: %+v", imgPaths)
		// 从s3上删除
		if len(imgPaths) > 0 {
			if _, err := common.UploadS3.DeleteMulti(imgPaths); err != nil {
				notes := fmt.Sprintf("UpdateLang s3上删除图片失败,图片路劲集合为: %+v,err: %s", imgPaths, err.Error())
				ss_log.Error(notes)
				dao.DictimagesDaoInst.AddDelFaildLog(notes)
			}
			// 删除image表中的图片
			for _, v := range imgPaths {
				if err := dao.DictimagesDaoInst.Delete(v); err != nil {
					notes := fmt.Sprintf("UpdateLang 删除图片记录失败,图片路劲为: %s,err: %s", v, err.Error())
					ss_log.Error(notes)
					dao.DictimagesDaoInst.AddDelFaildLog(notes)
				}
			}
		}

	}

	if err := dao.LangDaoInst.Insert(tx, req.Key, req.Type, req.LangKm, req.LangEn, req.LangCh); err != nil {
		ss_log.Error("更新或插入语言失败,key: %s,type: %s,langKm: %s,langEn: %s,langCh: %s,err: %s",
			req.Key, req.Type, req.LangKm, req.LangEn, req.LangCh, err.Error())
		reply.ResultCode = ss_err.ERR_SYS_DB_ADD
		return nil
	}
	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeleteLang(ctx context.Context, req *go_micro_srv_cust.DeleteLangRequest, reply *go_micro_srv_cust.DeleteLangReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, `update lang set is_delete='1' where key = $1`, req.Key)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) GetIncomeOugoConfigDetail(ctx context.Context, req *go_micro_srv_cust.GetIncomeOugoConfigDetailRequest, reply *go_micro_srv_cust.GetIncomeOugoConfigDetailReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ioc.is_delete", Val: "0", EqType: "="},
		{Key: "ioc.income_ougo_config_no", Val: req.IncomeOugoConfigNo, EqType: "="},
	})

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT ioc.income_ougo_config_no, ioc.currency_type, ioc.name, ioc.use_status, ioc.idx, ioc.config_type " +
		" FROM income_ougo_config ioc " + where2
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	data := &go_micro_srv_cust.IncomeOugoConfigData{}
	if err == nil {
		for rows.Next() {
			err = rows.Scan(
				&data.IncomeOugoConfigNo,
				&data.CurrencyType,
				&data.Name,
				&data.UseStatus,
				&data.Idx,
				&data.ConfigType,
			)

			if err != nil {
				ss_log.Error("err=[%v]", err)
				continue
			}
		}
	} else {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.Data = &go_micro_srv_cust.IncomeOugoConfigData{}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

func (c *CustHandler) GetLangDetail(ctx context.Context, req *go_micro_srv_cust.GetLangDetailRequest, reply *go_micro_srv_cust.GetLangDetailReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
		{Key: "key", Val: req.Key, EqType: "="},
	})

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT key, type, lang_km, lang_en,lang_ch " +
		" FROM lang " + where2
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	data := &go_micro_srv_cust.LangData{}
	if err == nil {
		for rows.Next() {
			var langKm, langEn, langCh sql.NullString
			err = rows.Scan(
				&data.Key,
				&data.Type,
				&langKm,
				&langEn,
				&langCh,
			)
			if err != nil {
				ss_log.Error("err=[%v]", err)
				continue
			}
			if data.Type == constants.LANG_TYPE_IMG {
				imgIds := []string{
					langKm.String,
					langEn.String,
					langCh.String,
				}
				var imgBase64s []string
				for _, imgId := range imgIds {
					//由图片id获取base64字符串(用于修改处回显图片)
					if imgId != "" {
						reqImg := &go_micro_srv_cust.UnAuthDownloadImageBase64Request{
							ImageId: imgId,
						}
						replyImg := &go_micro_srv_cust.UnAuthDownloadImageBase64Reply{}
						c.UnAuthDownloadImageBase64(ctx, reqImg, replyImg)
						if replyImg.ResultCode != ss_err.ERR_SUCCESS {
							ss_log.Error("获取图片base64字符串失败")
						}
						imgBase64s = append(imgBase64s, replyImg.ImageBase64)
					}
				}
				data.LangKm = imgBase64s[0]
				data.LangEn = imgBase64s[1]
				data.LangCh = imgBase64s[2]
			} else {
				data.LangKm = langKm.String
				data.LangEn = langEn.String
				data.LangCh = langCh.String
			}

		}
	} else {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

func (c *CustHandler) GetChannelDetail(ctx context.Context, req *go_micro_srv_cust.GetChannelDetailRequest, reply *go_micro_srv_cust.GetChannelDetailReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
		{Key: "channel_no", Val: req.ChannelNo, EqType: "="},
	})

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT channel_no, channel_name, create_time, use_status, logo_img_no, logo_img_no_grey, color_begin, color_end, channel_type  " +
		" FROM channel " + where2
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	data := &go_micro_srv_cust.ChannelDetailData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}
	for rows.Next() {
		var logoImgNo, logoImgNoGrey sql.NullString
		var colorBegin, colorEnd, channelType sql.NullString
		err = rows.Scan(
			&data.ChannelNo,
			&data.ChannelName,
			&data.CreateTime,
			&data.UseStatus,
			&logoImgNo,
			&logoImgNoGrey,
			&colorBegin,
			&colorEnd,
			&channelType,
		)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}

		data.ColorBegin = colorBegin.String
		data.ColorEnd = colorEnd.String
		data.ChannelType = channelType.String

		//由图片id获取base64字符串(用于修改处回显图片)
		if logoImgNo.String != "" {
			reqImg := &go_micro_srv_cust.UnAuthDownloadImageBase64Request{
				ImageId: logoImgNo.String,
			}
			replyImg := &go_micro_srv_cust.UnAuthDownloadImageBase64Reply{}
			c.UnAuthDownloadImageBase64(ctx, reqImg, replyImg)
			if replyImg.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取图片base64字符串失败")
			} else {
				data.LogoImg = replyImg.ImageBase64
			}
		}

		if logoImgNoGrey.String != "" {
			data.LogoImgNoGrey = logoImgNoGrey.String
			reqImg := &go_micro_srv_cust.UnAuthDownloadImageBase64Request{
				ImageId: logoImgNoGrey.String,
			}
			replyImg := &go_micro_srv_cust.UnAuthDownloadImageBase64Reply{}
			c.UnAuthDownloadImageBase64(ctx, reqImg, replyImg)
			if replyImg.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取图片base64字符串失败")
			} else {
				data.LogoImgGrey = replyImg.ImageBase64
			}
		}

	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

func (*CustHandler) InsertOrUpdateChannel(ctx context.Context, req *go_micro_srv_cust.InsertOrUpdateChannelRequest, reply *go_micro_srv_cust.InsertOrUpdateChannelReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	description := ""
	if req.ChannelNo == "" { //插入
		str1, legal1 := util.GetParamZhCn(req.ChannelType, util.ChannelType)
		if !legal1 {
			ss_log.Error("ChannelType %v", str1)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		channelNo, err := dao.ChannelDaoInst.AddChannel(tx, req.ChannelName, req.LogoImgNo, req.LogoImgNoGrey, req.ColorBegin, req.ColorEnd, req.ChannelType)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		description = fmt.Sprintf("渠道仓库添加渠道[%v],名称为[%v],渠道类型[%v]", channelNo, req.ChannelName, str1)
	} else { //修改
		// 删除s3上面的图片
		var imageURLs []string
		if req.LogoImgNo != "" {
			imageURL, err := dao.ChannelDaoInst.GetLogoImageUrlByChannelNo(tx, req.ChannelNo)
			if err != nil {
				ss_log.Error("GetLogoImageUrlByChannelNo 失败,channelNO: %s, err: %s", req.ChannelNo, err.Error())
			}
			if imageURL != "" {
				imageURLs = append(imageURLs, imageURL)
			}
		}

		if req.LogoImgNoGrey != "" {
			greyURL, err := dao.ChannelDaoInst.GetLogoImageGreyUrlByChannelNo(tx, req.ChannelNo)
			if err != nil {
				ss_log.Error("GetLogoImageGreyUrlByChannelNo 失败,channelNO: %s, err: %s", req.ChannelNo, err.Error())
			}
			if greyURL != "" {
				imageURLs = append(imageURLs, greyURL)
			}
		}

		if len(imageURLs) > 0 {
			// 删除多张
			if _, err := common.UploadS3.DeleteMulti(imageURLs); err != nil {
				notes := fmt.Sprintf("InsertOrUpdateChannel s3上删除图片失败,图片路劲集合为: %+v,err: %s", imageURLs, err.Error())
				ss_log.Error(notes)
				dao.DictimagesDaoInst.AddDelFaildLog(notes)
			}

			for _, v := range imageURLs {
				if err := dao.DictimagesDaoInst.Delete(v); err != nil {
					notes := fmt.Sprintf("InsertOrUpdateChannel 删除图片记录失败,图片路劲为: %s,err: %s", v, err.Error())
					ss_log.Error(notes)
					dao.DictimagesDaoInst.AddDelFaildLog(notes)
				}
			}

		}

		err2 := dao.ChannelDaoInst.UpdateChannel(tx, req.ChannelNo, req.ChannelName, req.LogoImgNo, req.LogoImgNoGrey, req.ColorBegin, req.ColorEnd)
		if err2 != ss_err.ERR_SUCCESS {
			ss_log.Error("err2=[%v]", err2)
			reply.ResultCode = ss_err.ERR_SYS_DB_UPDATE
			return nil
		}
		description = fmt.Sprintf("渠道仓库修改渠道[%v]信息为，名称[%v]", req.ChannelNo, req.ChannelName)
	}

	if errAddLog := dao.LogDaoInstance.InsertWebAccountLogTx(tx, description, req.LoginUid, constants.LogAccountWebType_Config); errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeleteChannel(ctx context.Context, req *go_micro_srv_cust.DeleteChannelRequest, reply *go_micro_srv_cust.DeleteChannelReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// 删除s3上面的图片
	var imageURLs []string
	imageURL, err := dao.ChannelDaoInst.GetLogoImageUrlByChannelNo(tx, req.ChannelNo)
	if err != nil {
		ss_log.Error("DeleteChannel GetLogoImageUrlByChannelNo 失败,channelNO: %s, err: %s", req.ChannelNo, err.Error())
	}
	if imageURL != "" {
		imageURLs = append(imageURLs, imageURL)
	}
	greyURL, err := dao.ChannelDaoInst.GetLogoImageGreyUrlByChannelNo(tx, req.ChannelNo)
	if err != nil {
		ss_log.Error("DeleteChannel GetLogoImageGreyUrlByChannelNo 失败,channelNO: %s, err: %s", req.ChannelNo, err.Error())
	}
	if greyURL != "" {
		imageURLs = append(imageURLs, greyURL)
	}

	if err := dao.ChannelDaoInst.ModifyChannelStatus(tx, req.ChannelNo); err != nil {
		ss_log.Error("DeleteChannel 失败,err: %s", err.Error())
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	description := fmt.Sprintf("渠道仓库 删除渠道[%v]", req.ChannelNo)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	if len(imageURLs) > 0 {
		// 删除多张
		if _, err := common.UploadS3.DeleteMulti(imageURLs); err != nil {
			notes := fmt.Sprintf("DeleteChannel s3上删除图片失败,图片路劲集合为: %+v,err: %s", imageURLs, err.Error())
			ss_log.Error(notes)
			dao.DictimagesDaoInst.AddDelFaildLog(notes)
		}
		for _, v := range imageURLs {
			// 删除数据库的图片
			if err := dao.DictimagesDaoInst.Delete(v); err != nil {
				notes := fmt.Sprintf("DeleteChannel 删除图片记录失败,图片路劲为: %s,err: %s", v, err.Error())
				ss_log.Error(notes)
				dao.DictimagesDaoInst.AddDelFaildLog(notes)
			}
		}
	}
	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) ModifyChannelStatus(ctx context.Context, req *go_micro_srv_cust.ModifyChannelStatusRequest, reply *go_micro_srv_cust.ModifyChannelStatusReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	str1, legal1 := util.GetParamZhCn(req.UseStatus, util.UseStatus)
	if !legal1 {
		ss_log.Error("%v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	sqlUpdate := "update channel set use_status = $2 where channel_no = $1 "
	updateErr := ss_sql.Exec(dbHandler, sqlUpdate, req.ChannelNo, req.UseStatus)
	if updateErr != nil {
		ss_log.Error("updateErr=[%v]", updateErr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	description := fmt.Sprintf("修改渠道仓库 渠道[%v]的使用状态为[%v]", req.ChannelNo, str1)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//
func (*CustHandler) ModifyChannelPosStatus(ctx context.Context, req *go_micro_srv_cust.ModifyChannelPosStatusRequest, reply *go_micro_srv_cust.ModifyChannelPosStatusReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	str1, legal1 := util.GetParamZhCn(req.UseStatus, util.UseStatus)
	if !legal1 {
		ss_log.Error("UseStatus %v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	sqlUpdate := "update channel_servicer set use_status = $2 where id = $1 "
	updateErr := ss_sql.Exec(dbHandler, sqlUpdate, req.Id, req.UseStatus)
	if updateErr != nil {
		ss_log.Error("updateErr=[%v]", updateErr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	description := fmt.Sprintf("修改服务商渠道id[%v],使用状态修改为[%v]", req.Id, str1)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Servicer)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) InsertPosChannel(ctx context.Context, req *go_micro_srv_cust.InsertPosChannelRequest, reply *go_micro_srv_cust.InsertPosChannelReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//是否推荐(1-推荐，0-不推荐)
	str1, legal1 := util.GetParamZhCn(req.IsRecom, util.IsRecom)
	if !legal1 {
		ss_log.Error("IsRecom %v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	str2, legal2 := util.GetParamZhCn(req.CurrencyType, util.CurrencyType)
	if !legal2 {
		ss_log.Error("CurrencyType %v", str2)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	if req.IsRecom == "1" { //如果插入的是推荐卡，一个币种的推荐卡只能有一张
		//将推荐渠道设置不推荐渠道（推荐渠道只能一个币种只有一个）
		err := dao.ChannelDaoInst.ModifyPosChannelIsRecom(tx, req.CurrencyType)
		if err != ss_err.ERR_SUCCESS {
			ss_log.Error("updateIsRecomErr=[%v]", err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	//确认是否有同样channelNo与CurrencyType的记录
	if dao.ChannelDaoInst.CheckPosChannel(req.ChannelNo, req.CurrencyType) {
		ss_log.Error("pos渠道存在相同币种的渠道")
		reply.ResultCode = ss_err.ERR_PosChannel_FAILD
		return nil
	}

	id, err2 := dao.ChannelDaoInst.AddPosChannel(tx, req.ChannelNo, req.IsRecom, req.CurrencyType)
	if err2 != nil {
		ss_log.Error("添加pos渠道出错，err2=[%v]", err2)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//添加关键操作日志
	description := fmt.Sprintf("添加pos渠道，id[%v] ,渠道No[%v],是否推荐使用[%v],币种[%v]", id, req.ChannelNo, str1, str2)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Servicer)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeletePosChannel(ctx context.Context, req *go_micro_srv_cust.DeletePosChannelRequest, reply *go_micro_srv_cust.DeletePosChannelReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlUpdate := "update channel_servicer set is_delete='1' where id = $1 and is_delete = '0' "
	updateErr := ss_sql.Exec(dbHandler, sqlUpdate, req.Id)
	if updateErr != nil {
		ss_log.Error("updateErr=[%v]", updateErr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	description := fmt.Sprintf("删除服务商渠道 id:[%v]", req.Id)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Servicer)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) ModifyChannelPosIsRecom(ctx context.Context, req *go_micro_srv_cust.ModifyChannelPosIsRecomRequest, reply *go_micro_srv_cust.ModifyChannelPosIsRecomReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	str1, legal1 := util.GetParamZhCn(req.IsRecom, util.IsRecom)
	if !legal1 {
		ss_log.Error("IsRecom %v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.IsRecom == "1" { //如果插入的是推荐卡，一个币种的推荐卡只能有一张
		currencyType := dao.ChannelDaoInst.GetCurrencyByPosChannelIdTx(tx, req.Id)
		if currencyType == "" {
			ss_log.Error("获取币种失败")
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		//将推荐渠道设置不推荐渠道（推荐渠道只能一个币种只有一个）
		err := dao.ChannelDaoInst.ModifyPosChannelIsRecom(tx, currencyType)
		if err != ss_err.ERR_SUCCESS {
			ss_log.Error("updateIsRecomErr=[%v]", err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	sqlUpdate := "update channel_servicer set is_recom = $2 where id = $1 and is_delete = '0' "
	updateErr := ss_sql.ExecTx(tx, sqlUpdate, req.Id, req.IsRecom)
	if updateErr != nil {
		ss_log.Error("updateErr=[%v]", updateErr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	description := fmt.Sprintf("修改服务商渠道 id:[%v],设置是否推荐为[%v]", req.Id, str1)

	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Servicer)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) InsertUseChannel(ctx context.Context, req *go_micro_srv_cust.InsertUseChannelRequest, reply *go_micro_srv_cust.InsertUseChannelReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	str1, legal1 := util.GetParamZhCn(req.CurrencyType, util.CurrencyType)
	if !legal1 {
		ss_log.Error("CurrencyType %v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	str2, legal2 := util.GetParamZhCn(req.SupportType, util.SupportType)
	if !legal2 {
		ss_log.Error("SupportType %v", str2)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	str3, legal3 := util.GetParamZhCn(req.SaveChargeType, util.ChargeType)
	if !legal3 {
		ss_log.Error("SaveChargeType %v", str3)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	str4, legal4 := util.GetParamZhCn(req.WithdrawChargeType, util.ChargeType)
	if !legal4 {
		ss_log.Error("WithdrawChargeType %v", str4)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	channelData, errGet := dao.ChannelDaoInst.GetChannelDetail(req.ChannelNo)
	if errGet != nil {
		ss_log.Error("[%v]渠道异常，查询渠道名称失败  err=[%v]", req.ChannelNo, errGet)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	str5, legal5 := util.GetParamZhCn(channelData.ChannelType, util.ChannelType)
	if !legal5 {
		ss_log.Error("ChannelType %v", str5)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	saveRate := ss_count.Div(req.SaveRate, "100").String() + "%"
	saveSingleMinFee := req.SaveSingleMinFee
	saveMaxAmount := req.SaveMaxAmount
	withdrawRate := ss_count.Div(req.WithdrawRate, "100").String() + "%"
	withdrawSingleMinFee := req.WithdrawSingleMinFee
	withdrawMaxAmount := req.WithdrawMaxAmount

	if req.CurrencyType == "usd" {
		saveSingleMinFee = ss_count.Div(req.SaveSingleMinFee, "100").String()
		saveMaxAmount = ss_count.Div(req.SaveMaxAmount, "100").String()
		withdrawSingleMinFee = ss_count.Div(req.WithdrawSingleMinFee, "100").String()
		withdrawMaxAmount = ss_count.Div(req.WithdrawMaxAmount, "100").String()
	}

	description := fmt.Sprintf(" 渠道名称[%v],币种[%v],渠道业务类型[%v],存款计算手续费类型[%v],取款计算手续费类型[%v],渠道类型[%v]", channelData.ChannelName, str1, str2, str3, str4, str5)
	description = fmt.Sprintf("%v,存款手续费率[%v],存款单笔手续费[%v],存款单笔最大金额[%v]", description, saveRate, saveSingleMinFee, saveMaxAmount)
	description = fmt.Sprintf("%v,取款手续费率[%v],取款单笔手续费[%v],取款单笔最大金额[%v]", description, withdrawRate, withdrawSingleMinFee, withdrawMaxAmount)

	if req.Id == "" {
		//确认是否有同样channelNo与CurrencyType的记录
		if dao.ChannelDaoInst.CheckUseChannel(req.ChannelNo, req.CurrencyType) {
			ss_log.Error("银行卡存取款渠道存在相同币种的渠道")
			reply.ResultCode = ss_err.ERR_UseChannel_FAILD
			return nil
		}

		id, err2 := dao.ChannelDaoInst.AddUseChannel(tx, dao.UseChannelData{
			ChannelNo:    req.ChannelNo,
			CurrencyType: req.CurrencyType,
			SupportType:  req.SupportType,
			ChannelType:  channelData.ChannelType,

			SaveRate:         req.SaveRate,
			SaveSingleMinFee: req.SaveSingleMinFee,
			SaveMaxAmount:    req.SaveMaxAmount,
			SaveChargeType:   req.SaveChargeType,

			WithdrawRate:         req.WithdrawRate,
			WithdrawSingleMinFee: req.WithdrawSingleMinFee,
			WithdrawMaxAmount:    req.WithdrawMaxAmount,
			WithdrawChargeType:   req.WithdrawChargeType,
		})
		if err2 != nil {
			ss_log.Error("err2=[%v]", err2)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		description = fmt.Sprintf("插入新银行卡存取款渠道 id[%v], %v ", id, description)
	} else {
		//确认是否有同样channelNo与CurrencyType的记录
		if dao.ChannelDaoInst.CheckUseChannel(req.ChannelNo, req.CurrencyType) {
			if id := dao.ChannelDaoInst.GetUseChannelId(req.ChannelNo, req.CurrencyType); id != "" && id != req.Id {
				ss_log.Error("银行卡存取款渠道存在相同币种的渠道")
				reply.ResultCode = ss_err.ERR_UseChannel_FAILD
				return nil
			}
		}

		err2 := dao.ChannelDaoInst.ModifyUseChannel(tx, req.Id, req.SupportType, req.SaveRate, req.SaveSingleMinFee,
			req.SaveMaxAmount, req.SaveChargeType, req.WithdrawRate, req.WithdrawSingleMinFee, req.WithdrawMaxAmount, req.WithdrawChargeType)
		if err2 != ss_err.ERR_SUCCESS {
			ss_log.Error("err2=[%v]", err2)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		description = fmt.Sprintf("修改旧银行卡存取款渠道 id[%v], %v", req.Id, description)
	}

	errAddLog := dao.LogDaoInstance.InsertWebAccountLogTx(tx, description, req.LoginUid, constants.LogAccountWebType_Config)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) DeleteUseChannel(ctx context.Context, req *go_micro_srv_cust.DeleteUseChannelRequest, reply *go_micro_srv_cust.DeleteUseChannelReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	if req.Id == "" {
		ss_log.Error("参数Id为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	sqlUpdate := "update channel_cust_config set is_delete='1' where id = $1 and is_delete = '0' "
	updateErr := ss_sql.Exec(dbHandler, sqlUpdate, req.Id)
	if updateErr != nil {
		ss_log.Error("updateErr=[%v]", updateErr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	description := fmt.Sprintf("删除旧银行卡存取款渠道 id[%v]", req.Id)

	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (*CustHandler) ModifyUseChannelStatus(ctx context.Context, req *go_micro_srv_cust.ModifyUseChannelStatusRequest, reply *go_micro_srv_cust.ModifyUseChannelStatusReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	if req.Id == "" {
		ss_log.Error("Id参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	str1, legal1 := util.GetParamZhCn(req.UseStatus, util.UseStatus)
	if !legal1 {
		ss_log.Error("UseStatus%v", str1)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	sqlUpdate := "update channel_cust_config set use_status = $2 where id = $1 "
	updateErr := ss_sql.Exec(dbHandler, sqlUpdate, req.Id, req.UseStatus)
	if updateErr != nil {
		ss_log.Error("updateErr=[%v]", updateErr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	description := fmt.Sprintf("修改旧银行卡存取款渠道 id[%v] 的渠道状态为[%v]", req.Id, str1)
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.LoginUid, constants.LogAccountWebType_Config)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
