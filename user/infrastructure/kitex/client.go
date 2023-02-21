package kitex

import (
	"dousheng_service/user/config"
	"dousheng_service/user/infrastructure/nacos"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/transport"
	api "github.com/joker-star-l/dousheng_idls/message/kitex_gen/api/message"
	"github.com/kitex-contrib/registry-nacos/resolver"
	"log"
	"runtime"
)

var MessageClient api.Client

func InitClient() {
	options := []client.Option{
		client.WithResolver(resolver.NewNacosResolver(nacos.Client)),
		client.WithTransportProtocol(transport.TTHeader),
		client.WithFailureRetry(retry.NewFailurePolicy()),
	}
	if runtime.GOOS != "windows" {
		options = append(options, client.WithMuxConnection(1))
	}
	var err error
	MessageClient, err = api.NewClient(config.C.MessageRpcName, options...)
	if err != nil {
		log.Panicln(err)
	}
}
