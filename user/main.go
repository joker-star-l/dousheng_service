package main

import (
	"dousheng_service/user/application/service"
	"dousheng_service/user/config"
	"dousheng_service/user/infrastructure/gorm"
	"dousheng_service/user/infrastructure/kitex"
	"dousheng_service/user/infrastructure/minio"
	"dousheng_service/user/infrastructure/nacos"
	"dousheng_service/user/infrastructure/redis"
	"dousheng_service/user/infrastructure/snowflake"
	"dousheng_service/user/interfaces"
	"flag"
	util_hertz "github.com/joker-star-l/dousheng_common/util/hertz"
	"os"
)

func argParse() {
	flag.IntVar(&config.C.MachineId, "machineId", os.Getpid()%1000, "machineId, default is pid")
	flag.StringVar(&config.C.Env, "env", "dev", "env, default is dev")
	flag.IntVar(&config.C.HttpPort, "httpPort", 8081, "httpPort, default is 8081")
	flag.IntVar(&config.C.RpcPort, "rpcPort", 6061, "httpPort, default is 6061")
	flag.Parse()
}

func init() {
	argParse()
	nacos.Init()
	kitex.InitServer(new(service.UserImpl))
	kitex.InitClient()
	gorm.Init()
	snowflake.Init()
	minio.Init()
	redis.Init()
}

func main() {
	h := util_hertz.InitServer(config.C.HttpPort)
	interfaces.InitRouter(h)
	h.Spin()
}
