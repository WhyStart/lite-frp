package main

import (
	"flag"
	"github.com/beego/beego/v2/core/logs"
	"lite-frp/client"
	"lite-frp/config"
	"os"
)

var clientConfigPath string

func main() {
	flag.StringVar(&clientConfigPath, "f", "./conf/client.int", "client config file path")
	flag.Parse()
	cfg, err := config.ReadClientConfigFile(clientConfigPath)
	if err != nil {
		logs.Error("读取配置文件失败 [%s]", clientConfigPath)
		os.Exit(1)
	}
	client.NewService(cfg).Run()
}
