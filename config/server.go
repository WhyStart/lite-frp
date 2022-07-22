package config

import (
	"fmt"
	"gopkg.in/ini.v1"
)

type ServerCommonConf struct {
	BindAddr string `ini:"bind_addr"`
	BindPort int    `ini:"bind_port"`
	Type     string `ini:"type"`
	Key      string `ini:"key"`
}

func GetDefaultServerConf() ServerCommonConf {
	return ServerCommonConf{
		BindAddr: "0.0.0.0",
		BindPort: 6000,
		Type:     "tcp",
		Key:      "",
	}
}

func ReadServerConfigFile(configPath string) (ServerCommonConf, error) {
	cfg, err := ini.Load(configPath)
	if err != nil {
		return ServerCommonConf{}, err
	}
	comm, err := cfg.GetSection("common")
	if err != nil {
		return ServerCommonConf{}, fmt.Errorf("invalid configuration file, not found [common] section")
	}
	clientConf := GetDefaultServerConf()
	err = comm.MapTo(&clientConf)
	if err != nil {
		return ServerCommonConf{}, err
	}
	return clientConf, err
}
