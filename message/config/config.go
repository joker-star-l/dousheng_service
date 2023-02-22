package config

import (
	common "github.com/joker-star-l/dousheng_common/entity"
	util_net "github.com/joker-star-l/dousheng_common/util/net"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var C = Config{
	common.Config{
		MachineId: 1,
		Env:       "dev",
		Ip:        util_net.LocalIp(),
		HttpName:  "message_service",
		HttpPort:  8083,
		RpcName:   "message_rpc_service",
		RpcPort:   6063,
		NacosClientParam: vo.NacosClientParam{
			ServerConfigs: []constant.ServerConfig{
				*constant.NewServerConfig("127.0.0.1", 8848),
			},
			ClientConfig: &constant.ClientConfig{
				CacheDir: "/tmp/dousheng_service/message/nacos/cache/",
				LogDir:   "/tmp/dousheng_service/message/nacos/cache/",
			},
		},
		NacosConfigList: []vo.ConfigParam{
			{DataId: "message_service.json", Group: "DEFAULT_GROUP"},
		},
	},
	MessageService{},
}

type Config struct {
	common.Config
	MessageService
}

type MessageService struct {
	Dsn   string `json:"dsn"`
	Redis struct {
		Address  string `json:"address"`
		Password string `json:"password"`
	} `json:"redis"`
}
