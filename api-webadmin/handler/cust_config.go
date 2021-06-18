package handler

import (
	"context"

	"a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-webadmin/inner_util"
	"a.a/mp-server/api-webadmin/verify"
	"a.a/mp-server/common/constants"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_net"
	"github.com/gin-gonic/gin"
)

func (s *CustHandler) GetWithdrawConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			response, err := CustHandlerInst.Client.GetWithdrawConfig(context.TODO(), &custProto.GetWithdrawConfigRequest{
				ConfigType: container.GetValFromMapMaybe(params, "config_type").ToStringNoPoint(),
			})
			if err != nil {
				ss_log.Error("err=[%v]", err)
				return response.ResultCode, gin.H{}, nil
			}
			return ss_err.ERR_SUCCESS, gin.H{
				"facewithdrawdatasDatas":  response.FaceWithdrawDatas,
				"phonewithdrawdatasDatas": response.PhoneWithdrawDatas,
				"depositdatasDatas":       response.DepositDatas,
				"transferdatasDatas":      response.TransferDatas,
			}, nil
		}, "params")
	}
}

func (s *CustHandler) UpdateWithdrawConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, err := CustHandlerInst.Client.UpdateWithdrawConfig(context.TODO(), &custProto.UpdateWithdrawConfigRequest{
				ConfigType: container.GetValFromMapMaybe(params, "config_type").ToStringNoPoint(),
				LoginUid:   inner_util.GetJwtDataString(c, "account_uid"),

				UsdFaceWithdrawRate:   container.GetValFromMapMaybe(params, "usd_face_withdraw_rate").ToStringNoPoint(),
				UsdFaceMinWithdrawFee: container.GetValFromMapMaybe(params, "usd_face_min_withdraw_fee").ToStringNoPoint(),
				UsdFaceFreeFeePerYear: container.GetValFromMapMaybe(params, "usd_face_free_fee_per_year").ToStringNoPoint(),
				UsdFaceSingleMax:      container.GetValFromMapMaybe(params, "usd_face_single_max").ToStringNoPoint(),
				UsdFaceSingleMin:      container.GetValFromMapMaybe(params, "usd_face_single_min").ToStringNoPoint(),
				KhrFaceWithdrawRate:   container.GetValFromMapMaybe(params, "khr_face_withdraw_rate").ToStringNoPoint(),
				KhrFaceMinWithdrawFee: container.GetValFromMapMaybe(params, "khr_face_min_withdraw_fee").ToStringNoPoint(),
				KhrFaceFreeFeePerYear: container.GetValFromMapMaybe(params, "khr_face_free_fee_per_year").ToStringNoPoint(),
				KhrFaceSingleMax:      container.GetValFromMapMaybe(params, "khr_face_single_max").ToStringNoPoint(),
				KhrFaceSingleMin:      container.GetValFromMapMaybe(params, "khr_face_single_min").ToStringNoPoint(),

				UsdPhoneWithdrawRate:   container.GetValFromMapMaybe(params, "usd_phone_withdraw_rate").ToStringNoPoint(),
				UsdPhoneMinWithdrawFee: container.GetValFromMapMaybe(params, "usd_phone_min_withdraw_fee").ToStringNoPoint(),
				UsdPhoneFreeFeePerYear: container.GetValFromMapMaybe(params, "usd_phone_free_fee_per_year").ToStringNoPoint(),
				UsdPhoneSingleMax:      container.GetValFromMapMaybe(params, "usd_phone_single_max").ToStringNoPoint(),
				UsdPhoneSingleMin:      container.GetValFromMapMaybe(params, "usd_phone_single_min").ToStringNoPoint(),
				KhrPhoneWithdrawRate:   container.GetValFromMapMaybe(params, "khr_phone_withdraw_rate").ToStringNoPoint(),
				KhrPhoneMinWithdrawFee: container.GetValFromMapMaybe(params, "khr_phone_min_withdraw_fee").ToStringNoPoint(),
				KhrPhoneFreeFeePerYear: container.GetValFromMapMaybe(params, "khr_phone_free_fee_per_year").ToStringNoPoint(),
				KhrPhoneSingleMax:      container.GetValFromMapMaybe(params, "khr_phone_single_max").ToStringNoPoint(),
				KhrPhoneSingleMin:      container.GetValFromMapMaybe(params, "khr_phone_single_min").ToStringNoPoint(),

				KhrDepositRate:      container.GetValFromMapMaybe(params, "khr_deposit_rate").ToStringNoPoint(),
				KhrMinDepositFee:    container.GetValFromMapMaybe(params, "khr_min_deposit_fee").ToStringNoPoint(),
				KhrDepositSingleMin: container.GetValFromMapMaybe(params, "khr_deposit_single_min").ToStringNoPoint(),
				KhrDepositSingleMax: container.GetValFromMapMaybe(params, "khr_deposit_single_max").ToStringNoPoint(),
				UsdDepositRate:      container.GetValFromMapMaybe(params, "usd_deposit_rate").ToStringNoPoint(),
				UsdMinDepositFee:    container.GetValFromMapMaybe(params, "usd_min_deposit_fee").ToStringNoPoint(),
				UsdDepositSingleMin: container.GetValFromMapMaybe(params, "usd_deposit_single_min").ToStringNoPoint(),
				UsdDepositSingleMax: container.GetValFromMapMaybe(params, "usd_deposit_single_max").ToStringNoPoint(),

				KhrTransferRate:      container.GetValFromMapMaybe(params, "khr_transfer_rate").ToStringNoPoint(),
				KhrMinTransferFee:    container.GetValFromMapMaybe(params, "khr_min_transfer_fee").ToStringNoPoint(),
				KhrTransferSingleMin: container.GetValFromMapMaybe(params, "khr_transfer_single_min").ToStringNoPoint(),
				KhrTransferSingleMax: container.GetValFromMapMaybe(params, "khr_transfer_single_max").ToStringNoPoint(),
				UsdTransferRate:      container.GetValFromMapMaybe(params, "usd_transfer_rate").ToStringNoPoint(),
				UsdMinTransferFee:    container.GetValFromMapMaybe(params, "usd_min_transfer_fee").ToStringNoPoint(),
				UsdTransferSingleMin: container.GetValFromMapMaybe(params, "usd_transfer_single_min").ToStringNoPoint(),
				UsdTransferSingleMax: container.GetValFromMapMaybe(params, "usd_transfer_single_max").ToStringNoPoint(),
			})
			//ss_log.Info("err=[%v],resp=[%v]", err, response)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}

			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) GetExchangeRateConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			response, err := CustHandlerInst.Client.GetExchangeRateConfig(context.TODO(), &custProto.GetExchangeRateConfigRequest{})
			if err != nil {
				ss_log.Error("err=[%v]", err)
				return response.ResultCode, gin.H{}, nil
			}

			return ss_err.ERR_SUCCESS, gin.H{
				"usd_to_khr":     response.Datas.UsdToKhr,
				"khr_to_usd":     response.Datas.KhrToUsd,
				"usd_to_khr_fee": response.Datas.UsdToKhrFee,
				"khr_to_usd_fee": response.Datas.KhrToUsdFee,
			}, nil
		}, "params")
	}
}

func (s *CustHandler) UpdateExchangeRateConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, err := CustHandlerInst.Client.UpdateExchangeRateConfig(context.TODO(), &custProto.UpdateExchangeRateConfigRequest{
				UsdToKhr: container.GetValFromMapMaybe(params, "usd_to_khr").ToStringNoPoint(),
				KhrToUsd: container.GetValFromMapMaybe(params, "khr_to_usd").ToStringNoPoint(),

				UsdToKhrFee: container.GetValFromMapMaybe(params, "usd_to_khr_fee").ToStringNoPoint(),
				KhrToUsdFee: container.GetValFromMapMaybe(params, "khr_to_usd_fee").ToStringNoPoint(),

				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),
			})
			if err != nil {
				ss_log.Error("err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) GetBusinessConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			reply, err := CustHandlerInst.Client.GetBusinessConfig(context.TODO(), &custProto.GetBusinessConfigRequest{})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, gin.H{
				"transfer_config_data": reply.TransferConfigData,
			}, nil
		}, "params")
	}
}

func (s *CustHandler) UpdateBusinessConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, err := CustHandlerInst.Client.UpdateBusinessConfig(context.TODO(), &custProto.UpdateBusinessConfigRequest{
				UsdAmountMinLimit: container.GetValFromMapMaybe(params, "usd_amount_min_limit").ToStringNoPoint(),
				UsdAmountMaxLimit: container.GetValFromMapMaybe(params, "usd_amount_max_limit").ToStringNoPoint(),

				KhrAmountMinLimit: container.GetValFromMapMaybe(params, "khr_amount_min_limit").ToStringNoPoint(),
				KhrAmountMaxLimit: container.GetValFromMapMaybe(params, "khr_amount_max_limit").ToStringNoPoint(),

				UsdRate: container.GetValFromMapMaybe(params, "usd_rate").ToStringNoPoint(),
				KhrRate: container.GetValFromMapMaybe(params, "khr_rate").ToStringNoPoint(),

				UsdMinFee: container.GetValFromMapMaybe(params, "usd_min_fee").ToStringNoPoint(),
				KhrMinFee: container.GetValFromMapMaybe(params, "khr_min_fee").ToStringNoPoint(),

				BatchNumber: container.GetValFromMapMaybe(params, "batch_number").ToStringNoPoint(),

				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) GetFuncConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetFuncConfig(context.TODO(), &custProto.GetFuncConfigRequest{
				Page:            strext.ToInt32(params[0]),
				PageSize:        strext.ToInt32(params[1]),
				ApplicationType: strext.ToString(params[2]),
				FuncName:        strext.ToString(params[3]),
			})
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "application_type", "func_name")
	}
}

func (s *CustHandler) UpdateFuncConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.UpdateFuncConfigRequest{
				FuncName:        container.GetValFromMapMaybe(params, "func_name").ToStringNoPoint(),
				JumpUrl:         container.GetValFromMapMaybe(params, "jump_url").ToStringNoPoint(),
				ImgBase64:       container.GetValFromMapMaybe(params, "img_base64").ToStringNoPoint(),
				ApplicationType: container.GetValFromMapMaybe(params, "application_type").ToStringNoPoint(),
				FuncNo:          container.GetValFromMapMaybe(params, "func_no").ToStringNoPoint(),
			}

			if errStr := verify.UpdateFuncConfigReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			response, _ := CustHandlerInst.Client.UpdateFuncConfig(context.TODO(), req)
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) SwapFuncConfigIdx() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, _ := CustHandlerInst.Client.SwapFuncConfigIdx(context.TODO(), &custProto.SwapFuncConfigIdxRequest{
				Idx:      container.GetValFromMapMaybe(params, "idx").ToStringNoPoint(),
				SwapType: container.GetValFromMapMaybe(params, "swap_type").ToStringNoPoint(),
				AppType:  container.GetValFromMapMaybe(params, "application_type").ToStringNoPoint(),
			})
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) DeleteFuncConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, _ := CustHandlerInst.Client.DeleteFuncConfig(context.TODO(), &custProto.DeleteFuncConfigRequest{
				FuncNo: container.GetValFromMapMaybe(params, "func_no").ToStringNoPoint(),
			})
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) ModifyUseStatusFuncConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, _ := CustHandlerInst.Client.ModifyUseStatusFuncConfig(context.TODO(), &custProto.ModifyUseStatusFuncConfigRequest{
				FuncNo:    container.GetValFromMapMaybe(params, "func_no").ToStringNoPoint(),
				UseStatus: container.GetValFromMapMaybe(params, "use_status").ToStringNoPoint(),
			})
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) GetTransferSecurityConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetTransferSecurityConfig(context.TODO(), &custProto.GetTransferSecurityConfigRequest{})
			return reply.ResultCode, gin.H{
				"continuous_err_password": reply.Datas.ContinuousErrPassword,
				"err_payment_pwd_count":   reply.Datas.ErrPaymentPwdCount,
			}, 0, err
		}, "params")
	}
}

func (s *CustHandler) GetWriteOffDurationDateConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetWriteOffDurationDateConfig(context.TODO(), &custProto.GetWriteOffDurationDateConfigRequest{})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, gin.H{
				"duration_date": reply.DurationDate,
			}, 0, nil
		}, "params")
	}
}

func (s *CustHandler) UpdateWriteOffDurationDateConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.UpdateWriteOffDurationDateConfigRequest{
				DurationDate: container.GetValFromMapMaybe(params, "duration_date").ToStringNoPointReg(`^[\d]+$`),
				LoginUid:     inner_util.GetJwtDataString(c, "account_uid"),
			}

			// 参数校验
			if errStr := verify.UpdateWriteOffDurationDateConfigVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			response, err := CustHandlerInst.Client.UpdateWriteOffDurationDateConfig(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return response.ResultCode, "", err
		})
	}
}

func (s *CustHandler) GetFuncConfigDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			response, _ := CustHandlerInst.Client.GetFuncConfigDetail(context.TODO(), &custProto.GetFuncConfigDetailRequest{
				FuncNo: container.GetValFromMapMaybe(params, "func_no").ToStringNoPoint(),
			})

			return response.ResultCode, gin.H{
				"data": response.Data,
			}, nil
		}, "params")
	}
}

func (s *CustHandler) UpdateTransferSecurityConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.UpdateTransferSecurityConfigRequest{
				ContinuousErrPassword: container.GetValFromMapMaybe(params, "continuous_err_password").ToStringNoPointReg(`^[\d]+$`), //登录密码错误次数
				ErrPaymentPwdCount:    container.GetValFromMapMaybe(params, "err_payment_pwd_count").ToStringNoPointReg(`^[\d]+$`),   //支付密码错误次数
				LoginUid:              inner_util.GetJwtDataString(c, "account_uid"),
			}

			// 参数校验
			if errStr := verify.UpdateTransferSecurityConfigVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			response, err := CustHandlerInst.Client.UpdateTransferSecurityConfig(context.TODO(), req)

			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return response.ResultCode, "", err
		})
	}
}

func (s *CustHandler) GetIncomeOugoConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetIncomeOugoConfig(context.TODO(), &custProto.GetIncomeOugoConfigRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				ConfigType:   strext.ToString(params[2]),
				Name:         strext.ToString(params[3]),
				CurrencyType: strext.ToString(params[4]),
			})
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "config_type", "name", "currency_type")
	}
}

func (s *CustHandler) UpdateIncomeOugoConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, err := CustHandlerInst.Client.UpdateIncomeOugoConfig(context.TODO(), &custProto.UpdateIncomeOugoConfigRequest{
				IncomeOugoConfigNo: container.GetValFromMapMaybe(params, "income_ougo_config_no").ToStringNoPoint(),
				CurrencyType:       container.GetValFromMapMaybe(params, "currency_type").ToStringNoPoint(),
				Name:               container.GetValFromMapMaybe(params, "name").ToStringNoPoint(),
				ConfigType:         container.GetValFromMapMaybe(params, "config_type").ToStringNoPoint(),
			})

			return response.ResultCode, "", err
		})
	}
}

func (s *CustHandler) UpdateIncomeOugoConfigUseStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, err := CustHandlerInst.Client.UpdateIncomeOugoConfigUseStatus(context.TODO(), &custProto.UpdateIncomeOugoConfigUseStatusRequest{
				IncomeOugoConfigNo: container.GetValFromMapMaybe(params, "income_ougo_config_no").ToStringNoPoint(),
				UseStatus:          container.GetValFromMapMaybe(params, "use_status").ToStringNoPoint(),
			})

			return response.ResultCode, "", err
		})
	}
}

func (s *CustHandler) UpdateIncomeOugoConfigIdx() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, _ := CustHandlerInst.Client.UpdateIncomeOugoConfigIdx(context.TODO(), &custProto.UpdateIncomeOugoConfigIdxRequest{
				Idx:        container.GetValFromMapMaybe(params, "idx").ToStringNoPoint(),
				SwapType:   container.GetValFromMapMaybe(params, "swap_type").ToStringNoPoint(),
				ConfigType: container.GetValFromMapMaybe(params, "config_type").ToStringNoPoint(),
			})

			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) DeleteIncomeOugoConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, _ := CustHandlerInst.Client.DeleteIncomeOugoConfig(context.TODO(), &custProto.DeleteIncomeOugoConfigRequest{
				IncomeOugoConfigNo: container.GetValFromMapMaybe(params, "income_ougo_config_no").ToStringNoPoint(),
			})
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) GetLangs() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetLangs(context.TODO(), &custProto.GetLangsRequest{
				Page:      strext.ToInt32(params[0]),
				PageSize:  strext.ToInt32(params[1]),
				Key:       strext.ToStringNoPoint(params[2]),
				Type:      strext.ToStringNoPoint(params[3]),
				SearchKey: strext.ToStringNoPoint(params[4]),
				LangCh:    strext.ToStringNoPoint(params[5]),
				LangEn:    strext.ToStringNoPoint(params[6]),
				LangKm:    strext.ToStringNoPoint(params[7]),
			})
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "key", "type", "search", "lang_ch", "lang_en", "lang_km")
	}
}

func (s *CustHandler) InsertOrUpdateLang() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.UpdateLangRequest{
				//Id:     container.GetValFromMapMaybe(params, "id").ToStringNoPoint(),
				Key:    container.GetValFromMapMaybe(params, "key").ToStringNoPoint(),
				Type:   container.GetValFromMapMaybe(params, "type").ToStringNoPoint(),
				LangKm: container.GetValFromMapMaybe(params, "lang_km").ToStringNoPoint(),
				LangEn: container.GetValFromMapMaybe(params, "lang_en").ToStringNoPoint(),
				LangCh: container.GetValFromMapMaybe(params, "lang_ch").ToStringNoPoint(),
			}

			// 参数校验
			if errStr := verify.CheckUpdateLangRequestVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			switch req.Type {
			case constants.LANG_TYPE_WORD: //如果是文字直接添加
			case constants.LANG_TYPE_IMG: //如果是图片则存到服务商后，添加的是返回的图片id
				replyImgKm, _ := CustHandlerInst.Client.UploadImage(context.TODO(), &custProto.UploadImageRequest{
					ImageStr:   req.LangKm,
					AccountUid: "00000000-0000-0000-0000-000000000000",
					Type:       constants.UploadImage_UnAuth,
				})

				replyImgEn, _ := CustHandlerInst.Client.UploadImage(context.TODO(), &custProto.UploadImageRequest{
					ImageStr:   req.LangEn,
					AccountUid: "00000000-0000-0000-0000-000000000000",
					Type:       constants.UploadImage_UnAuth,
				})

				replyImgCh, _ := CustHandlerInst.Client.UploadImage(context.TODO(), &custProto.UploadImageRequest{
					ImageStr:   req.LangCh,
					AccountUid: "00000000-0000-0000-0000-000000000000",
					Type:       constants.UploadImage_UnAuth,
				})

				if replyImgKm.ResultCode != ss_err.ERR_SUCCESS || replyImgEn.ResultCode != ss_err.ERR_SUCCESS || replyImgCh.ResultCode != ss_err.ERR_SUCCESS {
					ss_log.Error("保存多语言图片失败")
					return ss_err.ERR_SAVE_IMAGE_FAILD, "", nil
				} else {
					ss_log.Info("保存图片成功")
					req.LangKm = replyImgKm.ImageId
					req.LangEn = replyImgEn.ImageId
					req.LangCh = replyImgCh.ImageId
				}

			}

			response, err := CustHandlerInst.Client.UpdateLang(context.TODO(), req)
			return response.ResultCode, "", err
		})
	}
}

func (s *CustHandler) DeleteLang() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, _ := CustHandlerInst.Client.DeleteLang(context.TODO(), &custProto.DeleteLangRequest{
				Key: container.GetValFromMapMaybe(params, "key").ToStringNoPoint(),
			})
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) GetIncomeOugoConfigDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			response, _ := CustHandlerInst.Client.GetIncomeOugoConfigDetail(context.TODO(), &custProto.GetIncomeOugoConfigDetailRequest{
				IncomeOugoConfigNo: container.GetValFromMapMaybe(params, "income_ougo_config_no").ToStringNoPoint(),
			})
			return response.ResultCode, gin.H{
				"data": response.Data,
			}, nil
		}, "params")
	}
}

func (s *CustHandler) GetLangDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			response, err := CustHandlerInst.Client.GetLangDetail(context.TODO(), &custProto.GetLangDetailRequest{
				Key: container.GetValFromMapMaybe(params, "key").ToStringNoPoint(),
			})
			return response.ResultCode, gin.H{
				"data": response.Data,
			}, err
		}, "params")
	}
}

func (s *CustHandler) GetChannelDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			response, err := CustHandlerInst.Client.GetChannelDetail(context.TODO(), &custProto.GetChannelDetailRequest{
				ChannelNo: container.GetValFromMapMaybe(params, "channel_no").ToStringNoPoint(),
			})
			return response.ResultCode, gin.H{
				"data": response.Data,
			}, err
		}, "params")
	}
}

func (s *CustHandler) InsertOrUpdateChannel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			//====================上传图片====================
			imageStr := container.GetValFromMapMaybe(params, "logo_img").ToStringNoPoint()
			imageGreyStr := container.GetValFromMapMaybe(params, "logo_img_grey").ToStringNoPoint()
			if len(imageStr) > constants.UploadImgBase64LengthMax {
				ss_log.Error("彩色logo图片太大了")
				return ss_err.ERR_ACCOUNT_IMAGE_BIG, nil, nil
			}
			if len(imageGreyStr) > constants.UploadImgBase64LengthMax {
				ss_log.Error("灰色logo图片太大了")
				return ss_err.ERR_ACCOUNT_IMAGE_BIG, nil, nil
			}

			req := &custProto.UploadImageRequest{
				ImageStr:   imageStr,
				AccountUid: inner_util.GetJwtDataString(c, "account_uid"),
				//Type:       container.GetValFromMapMaybe(params, "type").ToInt32(),
				Type: constants.UploadImage_UnAuth, //类型是不需要授权的图片
			}

			if errStr := verify.UploadImageReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}
			reply, err := CustHandlerInst.Client.UploadImage(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用保存图片服务失败,err=[%v]", err)
				return ss_err.ERR_SAVE_IMAGE_FAILD, "", nil
			}

			//保存灰色图片
			reqGrey := &custProto.UploadImageRequest{
				ImageStr:   imageGreyStr,
				AccountUid: inner_util.GetJwtDataString(c, "account_uid"),
				//Type:       container.GetValFromMapMaybe(params, "type").ToInt32(),
				Type: constants.UploadImage_UnAuth, //类型是不需要授权的图片
			}

			if errStr := verify.UploadImageReqVerify(reqGrey); errStr != "" {
				return errStr, nil, nil
			}
			reply2, err2 := CustHandlerInst.Client.UploadImage(context.TODO(), reqGrey)
			if err2 != nil {
				ss_log.Error("调用保存图片服务失败,err=[%v]", err2)
				return ss_err.ERR_SAVE_IMAGE_FAILD, "", nil
			}

			//=====================插入信息==================
			channelReq := &custProto.InsertOrUpdateChannelRequest{
				ChannelNo:     container.GetValFromMapMaybe(params, "channel_no").ToStringNoPoint(),
				ChannelName:   container.GetValFromMapMaybe(params, "channel_name").ToStringNoPoint(),
				ColorBegin:    container.GetValFromMapMaybe(params, "color_begin").ToStringNoPoint(),
				ColorEnd:      container.GetValFromMapMaybe(params, "color_end").ToStringNoPoint(),
				ChannelType:   container.GetValFromMapMaybe(params, "channel_type").ToStringNoPoint(), //插入时传，更新时不传（不允许更新渠道类型）
				LoginUid:      inner_util.GetJwtDataString(c, "account_uid"),
				LogoImgNo:     reply.ImageId,
				LogoImgNoGrey: reply2.ImageId,
			}

			//验证数据
			if errStr := verify.InsertOrUpdateChannelVerify(channelReq); errStr != "" {
				return errStr, nil, nil
			}

			//开始添加或更新渠道
			replyC, errC := CustHandlerInst.Client.InsertOrUpdateChannel(context.TODO(), channelReq)
			if errC != nil {
				ss_log.Error("err=[%v]", errC)
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return replyC.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) DeleteChannel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, _ := CustHandlerInst.Client.DeleteChannel(context.TODO(), &custProto.DeleteChannelRequest{
				ChannelNo: container.GetValFromMapMaybe(params, "channel_no").ToStringNoPoint(),
				LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			})
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) ModifyChannelStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, _ := CustHandlerInst.Client.ModifyChannelStatus(context.TODO(), &custProto.ModifyChannelStatusRequest{
				ChannelNo: container.GetValFromMapMaybe(params, "channel_no").ToStringNoPoint(),
				UseStatus: container.GetValFromMapMaybe(params, "use_status").ToStringNoPoint(),
				LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			})
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) ModifyChannelPosStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.ModifyChannelPosStatusRequest{
				Id:        container.GetValFromMapMaybe(params, "id").ToStringNoPoint(),
				UseStatus: container.GetValFromMapMaybe(params, "use_status").ToStringNoPoint(),
				LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}
			if errStr := verify.ModifyChannelPosStatusVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			response, err := CustHandlerInst.Client.ModifyChannelPosStatus(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) InsertPosChannel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {

			req := &custProto.InsertPosChannelRequest{
				ChannelNo:    container.GetValFromMapMaybe(params, "channel_no").ToStringNoPoint(),
				IsRecom:      container.GetValFromMapMaybe(params, "is_recom").ToStringNoPoint(),
				CurrencyType: container.GetValFromMapMaybe(params, "currency_type").ToStringNoPoint(),
				LoginUid:     inner_util.GetJwtDataString(c, "account_uid"),
			}
			if errStr := verify.InsertPosChannelVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			//开始添加或更新渠道
			response, err := CustHandlerInst.Client.InsertPosChannel(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) DeletePosChannel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DeletePosChannelRequest{
				Id:       container.GetValFromMapMaybe(params, "id").ToStringNoPoint(),
				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),
			}

			if errStr := verify.DeletePosChannelVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			response, err := CustHandlerInst.Client.DeletePosChannel(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) ModifyChannelPosIsRecom() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {

			req := &custProto.ModifyChannelPosIsRecomRequest{
				Id:       container.GetValFromMapMaybe(params, "id").ToStringNoPoint(),
				IsRecom:  container.GetValFromMapMaybe(params, "is_recom").ToStringNoPoint(),
				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),
			}
			if errStr := verify.ModifyChannelPosIsRecomVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			response, err := CustHandlerInst.Client.ModifyChannelPosIsRecom(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}

			return response.ResultCode, "", nil
		})
	}
}
