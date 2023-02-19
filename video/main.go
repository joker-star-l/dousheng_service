package main

import (
	"dousheng_service/video/infrastructure/config"
	"dousheng_service/video/infrastructure/nacos"
	"dousheng_service/video/interfaces"
	"flag"
	util_hertz "github.com/joker-star-l/dousheng_common/util/hertz"
	"os"
)

func argParse() {
	flag.IntVar(&config.C.MachineId, "machineId", os.Getpid(), "machineId, default is pid")
	flag.StringVar(&config.C.Env, "env", "dev", "env, default is dev")
	flag.IntVar(&config.C.HttpPort, "httpPort", 8082, "httpPort, default is 8081")
	flag.IntVar(&config.C.RpcPort, "rpcPort", 6062, "httpPort, default is 6061")
	flag.Parse()
}

func init() {
	argParse()
	nacos.Init()
}

func main() {
	h := util_hertz.InitServer(config.C.HttpPort)
	interfaces.InitRouter(h)
	h.Spin()
}
