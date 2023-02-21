package kitex

import (
	"dousheng_service/message/config"
	"dousheng_service/message/infrastructure/nacos"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/joker-star-l/dousheng_common/config/log"
	util_kitex "github.com/joker-star-l/dousheng_common/util/kitex"
	util_net "github.com/joker-star-l/dousheng_common/util/net"
	"github.com/joker-star-l/dousheng_idls/message/kitex_gen/api"
	"github.com/joker-star-l/dousheng_idls/message/kitex_gen/api/message"
	"github.com/kitex-contrib/registry-nacos/registry"
	"runtime"
)

func InitServer(handler api.Message) {
	go func() {
		options := []server.Option{
			server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.C.RpcName}),
			server.WithServiceAddr(util_net.IpFromStr(config.C.Ip, config.C.RpcPort)),
			server.WithRegistry(registry.NewNacosRegistry(nacos.Client)),
			server.WithMiddleware(util_kitex.Log),
			server.WithReusePort(true),
		}
		if runtime.GOOS != "windows" {
			options = append(options, server.WithMuxTransport())
		}
		svr := message.NewServer(handler, options...)
		err := svr.Run()
		if err != nil {
			log.Slog.Panicln(err.Error())
		}
	}()
}
