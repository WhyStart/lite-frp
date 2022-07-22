package main

import (
	"flag"
	"lite-frp/config"
	"lite-frp/server"
	"lite-frp/tools/log"
	"os"
)

var serverConfigPath string

func main() {
	flag.StringVar(&serverConfigPath, "f", "./conf/server.int", "server config file path")
	flag.Parse()
	cfg, err := config.ReadServerConfigFile(serverConfigPath)
	if err != nil {
		log.Error("读取配置文件失败 [%s]", serverConfigPath)
		os.Exit(1)
	}
	server.NewService(cfg).Run()
}
