package main

import (
	"a.a/mp-server/conf-srv/i"
	_ "database/sql"
	"flag"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	// 运行时指定配置文件
	var configFile = flag.String("config", "./config/run.yml", "指定配置文件")
	flag.Parse()

	log.Printf("configFile:[%v]", *configFile)
	i.LoadAndStoreConfig(*configFile)
	log.Println("已同步配置...")
}
