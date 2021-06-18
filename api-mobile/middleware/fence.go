package middleware

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-mobile/common"
	"a.a/mp-server/api-mobile/dao"
	"a.a/mp-server/api-mobile/inner_util"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type FenceMw struct {
}

var FenceMwInst FenceMw

/**
 * 围栏中间件
 */
func (*FenceMw) FenceMw() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(common.INNER_IS_STOP) {
			ss_log.Info("skip")
			return
		}

		//posSn := inner_util.GetJwtDataString(c, "pos_sn")
		//if posSn == "" { // 不是 pos 机不处理
		//	return
		//}
		accountType := inner_util.GetJwtDataString(c, "account_type")
		if accountType != constants.AccountType_POS && accountType != constants.AccountType_SERVICER {
			return
		}
		opAccNo := inner_util.GetJwtDataString(c, "iden_no")

		var servicerNo string
		switch accountType {
		case constants.AccountType_POS:
			servicerNo = dao.ServiceDaoInst.GetServicerNoByCashierNo(opAccNo)
		case constants.AccountType_SERVICER:
			servicerNo = opAccNo
		}

		// 获取服务商的围栏开关
		scopeOff, posSn := dao.ServiceDaoInst.GetScopeOffNoBySrvNo(servicerNo)
		if scopeOff == "" || scopeOff == "0" {
			ss_log.Error("pos机未开启围栏控制 scopeOff---> %s", scopeOff)
			return
		}

		uri := c.Request.RequestURI
		// 二维码,存款,手机号取款,pos确认,pos取消
		if !strings.HasPrefix(uri, "/mobile/bill/gen_withdraw_code") && uri != "/mobile/bill/save_money" && uri != "/mobile/bill/mobile_num_withdrawal" && uri != "/mobile/bill/confirm_withdraw" && uri != "/mobile/bill/cancel_withdraw" {
			return
		}

		paramMap, _ := c.Get(common.INNER_PARAM_MAP)
		switch p2 := paramMap.(type) {
		case map[string]interface{}:
			latReq := strext.ToStringNoPoint(p2["lat"])
			lngReq := strext.ToStringNoPoint(p2["lng"])
			if latReq == "" || lngReq == "" {
				ss_log.Error("中间件 经纬度参数为空,lat---> %s,lng---> %s", latReq, lngReq)
				c.Set(common.RET_CODE, ss_err.ERR_POS_OUT_OF_RANGE)
				c.Set(common.INNER_IS_STOP, true)
				RespMwInst.PackInner(c)
				RsaMwInst.EncodeInner(c)
				RsaMwInst.SignInner(c)
				RespMwInst.RespInner(c)
				c.Abort()
				return
			}
			// 计算围栏
			// 先去redis里面查找看有没有此数据
			redisKey := "pos_sn_" + posSn
			var lat, lng, scope string
			var err error
			// ------------第一版本的redis------------
			//if value, _ := ss_data.GetPosSnFromCache(redisKey, cache.RedisCli, constants.DefPoolName); value == "" { // 查询数据库并设置进redis

			if value, _ := cache.RedisClient.Get(redisKey).Result(); value == "" { // 查询数据库并设置进redis
				lat, lng, scope, err = GetPosLatLng(posSn)
				if err != nil {
					ss_log.Error("%s", err.Error())
					c.Set(common.RET_CODE, ss_err.ERR_POS_OUT_OF_RANGE)
					c.Set(common.INNER_IS_STOP, true)
					RespMwInst.PackInner(c)
					RsaMwInst.EncodeInner(c)
					RsaMwInst.SignInner(c)
					RespMwInst.RespInner(c)
					c.Abort()
					return
				}
				redisValue := fmt.Sprintf("%s,%s,%s", lat, lng, scope)

				// 设进redis
				// ------------第一版本的redis------------
				//if err := ss_data.SetPosSnToCache(redisKey, redisValue, cache.RedisCli, constants.DefPoolName); err != nil {

				if err := cache.RedisClient.Set(redisKey, redisValue, constants.PosNoKeySecV2).Err(); err != nil {
					ss_log.Error("经纬度存进redis失败,posNo--->%s,lat--->%s,lng--->%s,scope--->%s", posSn, lat, lng, scope)
				}
			} else {
				split := strings.Split(value, ",")
				lat = split[0]
				lng = split[1]
				scope = split[2]
			}
			// 计算距离
			distance := ss_count.CountCircleDistance(strext.ToFloat64(lat), strext.ToFloat64(lng), strext.ToFloat64(latReq), strext.ToFloat64(lngReq))
			ss_log.Info("中间件 pos 计算的范围是--->%v,限定的范围是--->v", distance, scope)
			if distance > strext.ToFloat64(scope) {
				ss_log.Error("中间件 pos机超出使用范围,计算的范围是--->%v,限定的范围是--->v", distance, scope)

				c.Set(common.RET_CODE, ss_err.ERR_POS_OUT_OF_RANGE)
				c.Set(common.INNER_IS_STOP, true)
				RespMwInst.PackInner(c)
				RsaMwInst.EncodeInner(c)
				RsaMwInst.SignInner(c)
				RespMwInst.RespInner(c)
				c.Abort()
				return

			}
			return
		default:
			ss_log.Error("body格式错误")
			c.Set(common.INNER_IS_STOP, true)
			c.Set(common.RET_CODE, ss_err.ERR_SYS_BODY_NOT_JSON)
			return
		}

	}
}

func GetPosLatLng(posSn string) (lat, lng, scope string, err error) {
	// 根据pos_sn找服务商no
	serverNo := dao.ServicerTerminalDaoInst.GetSerPosServicerNoByPosNo(posSn)
	if serverNo == "" {
		//ss_log.Error("%s", "根据posNo查找服务商 id 失败")
		return "", "", "", errors.New("根据posNo查找服务商 id 失败")
	}
	lat, lng, scope = dao.ServiceDaoInst.GetLatLngInfoFromNo(serverNo)
	if lat == "" || lng == "" || scope == "" {
		//ss_log.Error("根据 serverNo 查找服务商范围失败 lat--->%s,lng--->%s,scope--->%s", lat, lng, scope)
		return "", "", "", errors.New("根据 serverNo 查找服务商范围失败")
	}
	return lat, lng, scope, nil
}
