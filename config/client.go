package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"lite-frp/tools/log"
)

type ClientCommonConf struct {
	ServerAddr string `ini:"server_addr"`
	ServerPort int    `ini:"server_port"`
	Type       string `ini:"type"`
	SK         string `ini:"sk"`
	LocalIP    string `ini:"local_ip"`
	LocalPort  int    `ini:"local_port"`
	RemotePort int    `ini:"remote_port"`
}

func GetDefaultClientConf() ClientCommonConf {
	return ClientCommonConf{
		ServerAddr: "0.0.0.0",
		ServerPort: 6000,
		Type:       "tcp",
		LocalIP:    "127.0.0.1",
		LocalPort:  0,
		RemotePort: 0,
	}
}

func ReadClientConfigFile(configPath string) (ClientCommonConf, error) {
	cfg, err := ini.Load(configPath)
	if err != nil {
		log.Error("%s file not find", configPath)
	}
	clientConf := GetDefaultClientConf()
	comm, err := cfg.GetSection("common")
	if err != nil {
		return ClientCommonConf{}, fmt.Errorf("invalid configuration file, not found [common] section")
	}
	err = comm.MapTo(&clientConf)
	if err != nil {
		return ClientCommonConf{}, err
	}
	local, err := cfg.GetSection("local")
	if err != nil {
		return ClientCommonConf{}, fmt.Errorf("invalid configuration file, not found [local] section")
	}
	err = local.MapTo(&clientConf)
	if err != nil {
		return ClientCommonConf{}, err
	}
	return clientConf, err
}
