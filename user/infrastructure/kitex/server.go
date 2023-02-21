package kitex

import (
	"dousheng_service/user/application/service"
	"dousheng_service/user/config"
	"dousheng_service/user/infrastructure/nacos"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/joker-star-l/dousheng_common/config/log"
	util_kitex "github.com/joker-star-l/dousheng_common/util/kitex"
	util_net "github.com/joker-star-l/dousheng_common/util/net"
	"github.com/joker-star-l/dousheng_idls/user/kitex_gen/api/user"
	"github.com/kitex-contrib/registry-nacos/registry"
	"runtime"
)

func InitServer() {
	go func() {
		options := []server.Option{
			server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.C.RpcName}),
			server.WithServiceAddr(util_net.IpFromStr(config.C.Ip, config.C.RpcPort)),
			server.WithRegistry(registry.NewNacosRegistry(nacos.Client)),
			server.WithMiddleware(util_kitex.Log),
		}
		if runtime.GOOS != "windows" {
			options = append(options, server.WithMuxTransport())
		}
		svr := user.NewServer(new(service.UserImpl), options...)
		err := svr.Run()
		if err != nil {
			log.Slog.Panicln(err.Error())
		}
	}()
}
