package main

import (
	my_config "dousheng_service/gateway/config"
	my_nacos "dousheng_service/gateway/config/nacos"
	"flag"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/middlewares/client/sd"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/hertz-contrib/registry/nacos"
	"github.com/hertz-contrib/reverseproxy"
	"github.com/joker-star-l/dousheng_common/config/log"
	util_hertz "github.com/joker-star-l/dousheng_common/util/hertz"
	"os"
	"strings"
)

func argParse() {
	flag.IntVar(&my_config.C.MachineId, "machineId", os.Getpid(), "machineId, default is pid")
	flag.StringVar(&my_config.C.Env, "env", "dev", "env, default is dev")
	flag.IntVar(&my_config.C.HttpPort, "httpPort", 7001, "httpPort, default is 7001")
	flag.Parse()
}

func init() {
	argParse()
	my_nacos.Init()
}

func main() {
	// 创建客户端
	cli, err := client.NewClient()
	if err != nil {
		log.Slog.Panicf("hertz client init error: %v", err.Error())
	}
	resolver := nacos.NewNacosResolver(my_nacos.Client)
	cli.Use(sd.Discovery(resolver))

	// 创建代理
	proxy := &reverseproxy.ReverseProxy{}
	proxy.SetClient(cli)
	proxy.SetDirector(func(req *protocol.Request) {
		serviceName := "error"
		for _, m := range my_config.C.Mapping {
			if strings.HasPrefix(string(req.URI().Path()), m.Prefix) {
				serviceName = m.Service
				break
			}
		}
		uri := string(reverseproxy.JoinURLPath(req, "http://"+serviceName))
		log.Slog.Infof("request: %v", uri)
		req.SetRequestURI(uri)
		req.Header.SetHostBytes([]byte(serviceName))
		req.Options().Apply([]config.RequestOption{config.WithSD(true)})
	})

	// 创建服务端
	h := util_hertz.InitServer(my_config.C.HttpPort)
	h.Any("/*path", proxy.ServeHTTP)
	h.Spin()
}
