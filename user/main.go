package main

import (
	"dousheng_service/user/infrastructure/config"
	_ "dousheng_service/user/infrastructure/nacos"
	"dousheng_service/user/interfaces"
	"flag"
	util_hertz "github.com/joker-star-l/dousheng_common/util/hertz"
	"os"
)

func argParse() {
	flag.IntVar(&config.C.MachineId, "machineId", os.Getpid(), "machineId, default is pid")
	flag.StringVar(&config.C.Env, "env", "dev", "env, default is dev")
	flag.Parse()
}

func main() {
	argParse()
	h := util_hertz.InitServer(config.C.HttpPort)
	interfaces.InitRouter(h)
	h.Spin()
}
