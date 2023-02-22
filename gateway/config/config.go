package config

import (
	common "github.com/joker-star-l/dousheng_common/entity"
	util_net "github.com/joker-star-l/dousheng_common/util/net"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var C = Config{
	common.Config{
		Env:      "dev",
		Ip:       util_net.LocalIp(),
		HttpName: "gateway",
		HttpPort: 7001,
		NacosClientParam: vo.NacosClientParam{
			ServerConfigs: []constant.ServerConfig{
				*constant.NewServerConfig("127.0.0.1", 8848),
			},
			ClientConfig: &constant.ClientConfig{
				CacheDir: "/tmp/dousheng_service/gateway/nacos/cache/",
				LogDir:   "/tmp/dousheng_service/gateway/nacos/log/",
			},
		},
		NacosConfigList: []vo.ConfigParam{
			{DataId: "gateway.json", Group: "DEFAULT_GROUP"},
		},
	},
	Gateway{},
}

type Config struct {
	common.Config
	Gateway
}

type Gateway struct {
	// key: prefix, value: service
	Mapping []struct {
		Prefix  string `json:"prefix"`
		Service string `json:"service"`
	} `json:"mapping"`
}
