package main

import (
	"dousheng_service/message/config"
	"dousheng_service/message/infrastructure/gorm"
	"dousheng_service/message/infrastructure/nacos"
	"dousheng_service/message/infrastructure/redis"
	"dousheng_service/message/infrastructure/snowflake"
	"dousheng_service/message/interfaces"
	"flag"
	"github.com/joker-star-l/dousheng_common/config/log"
	util_hertz "github.com/joker-star-l/dousheng_common/util/hertz"
	"os"
)

func argParse() {
	flag.IntVar(&config.C.MachineId, "machineId", os.Getpid()%1000, "machineId, default is pid")
	flag.StringVar(&config.C.Env, "env", "dev", "env, default is dev")
	flag.IntVar(&config.C.HttpPort, "httpPort", 8083, "httpPort, default is 8081")
	flag.IntVar(&config.C.RpcPort, "rpcPort", 6063, "httpPort, default is 6061")
	flag.Parse()
}

func init() {
	argParse()
	nacos.Init()
	gorm.Init()
	snowflake.Init()
	redis.Init()
}

func main() {
	log.Slog.Infoln(config.C.MachineId)
	h := util_hertz.InitServer(config.C.HttpPort)
	interfaces.InitRouter(h)
	h.Spin()
}
