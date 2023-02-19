package config

import (
	common "github.com/joker-star-l/dousheng_common/entity"
	util_net "github.com/joker-star-l/dousheng_common/util/net"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var C = Config{
	common.Config{
		MachineId: 2,
		Env:       "dev",
		Ip:        util_net.LocalIp(),
		HttpName:  "video_service",
		HttpPort:  8082,
		RpcName:   "video_rpc_service",
		RpcPort:   6062,
		NacosClientParam: vo.NacosClientParam{
			ServerConfigs: []constant.ServerConfig{
				*constant.NewServerConfig("111.229.8.227", 8848),
			},
			ClientConfig: &constant.ClientConfig{
				CacheDir: "/tmp/dousheng_service/video/nacos/cache/",
				LogDir:   "/tmp/dousheng_service/video/nacos/cache/",
			},
		},
		NacosConfigList: []vo.ConfigParam{
			{DataId: "video_service.json", Group: "DEFAULT_GROUP"},
		},
	},
	VideoService{},
}

type Config struct {
	common.Config
	VideoService
}

type VideoService struct {
}
