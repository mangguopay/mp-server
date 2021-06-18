package adapter

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	go_micro_srv_push "a.a/mp-server/common/proto/push"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/push2-srv/adapter/common"
	"a.a/mp-server/push2-srv/dao"
	"a.a/mp-server/push2-srv/m"
	"context"
	"fmt"
)

const (
	conditionTypeCountryCode = 1 // 国家码条件
	noCondition              = 0 // 没有限制
)

func Push(ctx context.Context, req *go_micro_srv_push.PushReqest) error {
	ss_log.Info("recv=[%v]", req)
	errFlag := false

	// 获取模板
	pushNoList, titleKey, contentKey, lenArgs, err := dao.PushMsgTypeDaoInst.GetTemplate(req.TempNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}

	// 发给谁
	for _, v := range req.Accounts {
		title, content := getMsg(v, titleKey, contentKey)
		// 直接替换
		str := []interface{}{}
		for _, v := range req.Args {
			str = append(str, v)
		}
		// 长度不对，直接屏蔽
		if strext.HasN(content, "%s") != len(req.Args) {
			ss_log.Error("content=[%v], err=[%v],lenArgs=[%v],len(req.Args)=[%v]", content, ss_err.ErrArgsLen, lenArgs, len(req.Args))
			errFlag = true
		}
		//填充内容
		content = fmt.Sprintf(content, str...)

		if req.TitleArgs != nil {
			str2 := []interface{}{}
			for _, v := range req.TitleArgs {
				str2 = append(str2, v)
			}
			// 长度不对，直接屏蔽
			if strext.HasN(title, "%s") != len(req.TitleArgs) {
				ss_log.Error("title=[%v], err=[%v],lenTitleArgs=[%v],len(req.TitleArgs)=[%v]", title, ss_err.ErrArgsLen, lenArgs, len(req.TitleArgs))
				errFlag = true
			}

			//填充标题
			title = fmt.Sprintf(title, str2...)
		}

		//ss_log.Info("内容---------------> %s", content)
		// TODO 处理lenArgs，额外参数是否需要也转码?
		sendSingle(v, title, content, pushNoList, req.TempNo, errFlag)
	}
	return nil
}

func getMsg(accReq *go_micro_srv_push.PushAccout, titleKey, contentKey string) (title, content string) {
	langStr := accReq.Lang
	if langStr == "" {
		langStr = dao.AccDaoInst.GetAccLang(accReq.AccountNo, accReq.AccountType)
		if langStr == "" {
			langStr = constants.LangEnUS
		}
	}
	// 调整模板
	title = dao.LangDaoInstance.GetText(titleKey, langStr)
	content = dao.LangDaoInstance.GetText(contentKey, langStr)
	return title, content
}

func sendSingle(accReq *go_micro_srv_push.PushAccout, title, content string, pushNoList []string, tempNo string, errFlag bool) {
	//// 只给用户发
	//if accReq.AccountType != constants.AccountType_USER {
	//	return
	//}
	// 怎么发
	for _, v := range pushNoList {
		configInfo, err := dao.PushConfigDaoInst.GetPushConfig(v)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}

		// 判断是否符合逻辑
		if !checkCondition(accReq, configInfo.ConditionType, configInfo.ConditionValue) {
			// 条件不符,跳过
			continue
		}

		message := ""
		ret := ss_err.ERR_SUCCESS
		accToken, errMsg := common.PushAdapterWrapperInst.GetAccToken(configInfo.Pusher, accReq)
		if accToken == "" {
			errFlag = true
			message = errMsg
			if message == "" {
				message = `账号对应信息不存在或不完整`
			}
			ret = ss_err.ERR_PARAM
		}

		if !errFlag {
			message, ret = common.PushAdapterWrapperInst.Send(configInfo.Pusher, m.SendReq{
				Title:    title,
				Content:  content,
				Config:   configInfo.Config,
				AccToken: accToken,
			})
		} else if message == "" {
			message = `配置错误，参数长度和入参数量不一致`
			ret = ss_err.ERR_PARAM
		}

		ss_log.Info("调用推送最终返回=[%v]", ret)
		// 记录日志   business, phone, content string, status,pushNO
		var status int
		if ret != ss_err.ERR_SUCCESS {
			status = 1
		} else {
			status = 0
		}
		if err := dao.PushRecordDaoInst.Insert(configInfo.Pusher, accToken, content, v, tempNo, message, status); err != nil {
			ss_log.Error("插入推送记录失败,err: %s", err.Error())
		}
	}
	return
}

// CheckCondition 检验是否限制发送短信
func checkCondition(pushAcc *go_micro_srv_push.PushAccout, conditionType int, conditionValue string) bool {
	switch conditionType {
	case conditionTypeCountryCode:
		return pushAcc.CountryCode == conditionValue
	default:
		// 无条件默认通过
		return true
	}
}
