package config

import (
	util_net "github.com/joker-star-l/dousheng_common/util/net"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var C = Config{
	Local{
		env:      "dev",
		Ip:       util_net.LocalIp(),
		HttpName: "gateway",
		HttpPort: 7001,
		NacosClientParam: vo.NacosClientParam{
			ServerConfigs: []constant.ServerConfig{
				*constant.NewServerConfig("111.229.8.227", 8848),
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
	Remote{},
}

type Config struct {
	Local
	Remote
}

type Local struct {
	MachineId        int
	env              string
	Ip               string
	HttpName         string
	HttpPort         int
	NacosClientParam vo.NacosClientParam
	NacosConfigList  []vo.ConfigParam
}

type Remote struct {
	Gateway
}

type Gateway struct {
	// key: prefix, value: service
	Mapping []struct {
		Prefix  string `json:"prefix"`
		Service string `json:"service"`
	} `json:"mapping"`
}
