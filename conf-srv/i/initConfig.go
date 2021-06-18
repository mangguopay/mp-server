package i

import (
	"a.a/mp-server/common/ss_etcd"
	"errors"
	"fmt"
	"github.com/micro/go-micro/config"
	"os"
)

func LoadAndStoreConfig(configFile string) {
	// 加载配置文件
	if err := LoadConfigFile(configFile); err != nil {
		panic("加载配置文件出错: " + err.Error())
	}

	// 从环境变量中读取etcd地址
	addrList := ss_etcd.GetEtcdAddr()

	// 将配置信息载入etcd中
	if err := StoreToEtcd(addrList); err != nil {
		panic("配置加载到etcd失败: " + err.Error())
	}
}

// 加载配置文件
func LoadConfigFile(configFile string) error {
	// 检查配置文件
	if err := CheckconfigFile(configFile); err != nil {
		return err
	}

	// 解析和加载配置
	if err := config.LoadFile(configFile); err != nil {
		return err
	}

	return nil
}

// 检查文件是否存在
func CheckconfigFile(configFile string) error {
	file, err := os.Stat(configFile)
	if err != nil {
		return err
	}

	if file.IsDir() {
		return errors.New(fmt.Sprintf("指定的配置文件是一个目录:%s", configFile))
	}

	return nil
}
