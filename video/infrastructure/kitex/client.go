package kitex

import (
	"dousheng_service/video/config"
	"dousheng_service/video/infrastructure/nacos"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/transport"
	api "github.com/joker-star-l/dousheng_idls/user/kitex_gen/api/user"
	"github.com/kitex-contrib/registry-nacos/resolver"
	"log"
	"runtime"
)

var UserClient api.Client

func InitClient() {
	options := []client.Option{
		client.WithResolver(resolver.NewNacosResolver(nacos.Client)),
		client.WithTransportProtocol(transport.TTHeader),
	}
	if runtime.GOOS != "windows" {
		options = append(options, client.WithMuxConnection(1))
	}
	var err error
	UserClient, err = api.NewClient(config.C.UserRpcName, options...)
	if err != nil {
		log.Panicln(err)
	}
}
