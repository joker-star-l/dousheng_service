package main

import (
	"dousheng_service/user/infrastructure/config"
	"dousheng_service/user/infrastructure/gorm"
	"dousheng_service/user/infrastructure/kitex"
	"dousheng_service/user/infrastructure/nacos"
	"dousheng_service/user/infrastructure/snowflake"
	"dousheng_service/user/interfaces"
	"flag"
	"github.com/joker-star-l/dousheng_common/config/log"
	util_hertz "github.com/joker-star-l/dousheng_common/util/hertz"
	"os"
)

func argParse() {
	flag.IntVar(&config.C.MachineId, "machineId", os.Getpid(), "machineId, default is pid")
	flag.StringVar(&config.C.Env, "env", "dev", "env, default is dev")
	flag.Parse()
}

func init() {
	argParse()
	log.Slog.Infof("machineId: %d", config.C.MachineId)
	nacos.Init()
	kitex.InitServer()
	gorm.Init()
	snowflake.Init()
}

func main() {
	h := util_hertz.InitServer(config.C.HttpPort)
	interfaces.InitRouter(h)
	h.Spin()
}
