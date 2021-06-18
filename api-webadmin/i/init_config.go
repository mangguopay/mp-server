package i

import (
	"fmt"
	"strings"
	"time"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-webadmin/common"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_struct"
	"github.com/micro/go-micro/v2/config"
)

// 初始化 aws s3 配置
func InitAwss3() {
	var s3Conf ss_struct.Awss3Conf

	// 获取配置信息
	err := config.Get("aws", "s3").Scan(&s3Conf)
	if err != nil {
		ss_log.Error("aws-s3初始化失败,err:%v", err)
	}

	// 验证配置信息
	if s3Conf.AccessKeyId == "" || s3Conf.SecretAccessKey == "" || s3Conf.Region == "" || s3Conf.Bucket == "" {
		ss_log.Error("aws-s3配置信息不完整,s3Conf:%+v", s3Conf)
		panic(fmt.Sprintf("aws-s3配置信息不完整,s3Conf:%+v", s3Conf))
	}

	ss_log.Info("s3Conf:%+v", s3Conf)

	// 初始化s3操作类
	common.InitUploadS3(s3Conf)
}

func DoInitBase() (host string, portFrom, portTo int) {
	m := map[string]interface{}{}
	err := config.Get("base", "base").Scan(&m)
	if err != nil {
		ss_log.Error("err=%v", err)
	}

	z, zErr := time.LoadLocation(strext.ToStringNoPoint(m["timezone"]))
	if zErr != nil {
		panic(fmt.Sprintf("解析时区出错,err: %v", zErr))
	}
	// 设置time包中的默认时区
	time.Local = z

	global.Tz = z
	p := strings.Split(strext.ToStringNoPoint(m["port"]), "-")

	return strext.ToStringNoPoint(m["host"]), strext.ToInt(p[0]), strext.ToInt(p[1])
}
