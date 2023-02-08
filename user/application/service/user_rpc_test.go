package service

import (
	"context"
	"dousheng_service/user/infrastructure/config"
	"dousheng_service/user/infrastructure/nacos"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/transport"
	"github.com/joker-star-l/dousheng_common/config/log"
	api "github.com/joker-star-l/dousheng_idls/user/kitex_gen/api/user"
	"github.com/kitex-contrib/registry-nacos/resolver"
	"runtime"
	"testing"
)

func TestUserImpl_UserInfo(t *testing.T) {
	nacos.Init()
	options := []client.Option{
		client.WithResolver(resolver.NewNacosResolver(nacos.Client)),
		client.WithTransportProtocol(transport.TTHeader),
	}
	if runtime.GOOS != "windows" {
		options = append(options, client.WithMuxConnection(1))
	}
	c, _ := api.NewClient(config.C.RpcName, options...)
	resp, _ := c.UserInfo(context.Background(), 1623010390776483840)
	log.Slog.Infoln(resp)
}
