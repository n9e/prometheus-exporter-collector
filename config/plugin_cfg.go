package config

import (
	"encoding/json"
	"fmt"
)

//{"target_urls":["http://192.168.2.30:9103/metrics?dns=192.168.1.7:3306","http://192.168.2.30:9103/metrics?dns=192.168.1.160:3306"],"service":"","endpoint":"", "step":60,"username":"","password":""}

type PluginCfg struct {
	TargetUrls []string `json:"target_urls"`
	Service    string   `json:"service"`
	Step       int      `json:"step"`
	Endpoint   string   `json:"endpoint"`
	Username   string   `json:"username"`
	Password   string   `json:"password"`
}

var Config *PluginCfg

func Get() *PluginCfg {
	return Config
}

func Parse(config string) error {
	err := json.Unmarshal([]byte(config), &Config)
	if err != nil {
		return fmt.Errorf("unmarshal config error:%v", err)
	}
	return nil
}
