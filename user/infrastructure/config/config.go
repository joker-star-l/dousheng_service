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
		HttpName:  "user_service",
		HttpPort:  8081,
		RpcName:   "user_rpc_service",
		RpcPort:   6061,
		NacosClientParam: vo.NacosClientParam{
			ServerConfigs: []constant.ServerConfig{
				*constant.NewServerConfig("111.229.8.227", 8848),
			},
			ClientConfig: &constant.ClientConfig{
				CacheDir: "/tmp/dousheng_service/user/nacos/cache/",
				LogDir:   "/tmp/dousheng_service/user/nacos/cache/",
			},
		},
		NacosConfigList: []vo.ConfigParam{
			{DataId: "user_service.json", Group: "DEFAULT_GROUP"},
		},
	},
	UserService{},
}

type Config struct {
	common.Config
	UserService
}

type UserService struct {
	Dsn string `json:"dsn"`
}
