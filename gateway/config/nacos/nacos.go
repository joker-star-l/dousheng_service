package nacos

import (
	"dousheng_service/gateway/config"
	"github.com/joker-star-l/dousheng_common/config/log"
	util_json "github.com/joker-star-l/dousheng_common/util/json"
	util_nacos "github.com/joker-star-l/dousheng_common/util/nacos"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
)

var Client naming_client.INamingClient
var ConfigClient config_client.IConfigClient

func Init() {
	Client, ConfigClient = util_nacos.NewClient(config.C.NacosClientParam)
	util_nacos.RegisterService(Client, config.C.Ip, config.C.HttpPort, config.C.HttpName)
	if config.C.NacosConfigList != nil {
		for i := range config.C.NacosConfigList {
			util_nacos.GetAndListenJSONConfig(ConfigClient, &config.C, config.C.NacosConfigList[i])
		}
	}
	log.Slog.Infof("%s", util_json.Str(config.C.Gateway))
}
