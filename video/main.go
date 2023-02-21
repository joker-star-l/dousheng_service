package main

import (
	"dousheng_service/video/config"
	"dousheng_service/video/infrastructure/gorm"
	"dousheng_service/video/infrastructure/kitex"
	"dousheng_service/video/infrastructure/minio"
	"dousheng_service/video/infrastructure/nacos"
	"dousheng_service/video/infrastructure/redis"
	"dousheng_service/video/infrastructure/snowflake"
	"dousheng_service/video/interfaces"
	"flag"
	util_hertz "github.com/joker-star-l/dousheng_common/util/hertz"
	"os"
)

func argParse() {
	flag.IntVar(&config.C.MachineId, "machineId", os.Getpid()%1000, "machineId, default is pid")
	flag.StringVar(&config.C.Env, "env", "dev", "env, default is dev")
	flag.IntVar(&config.C.HttpPort, "httpPort", 8082, "httpPort, default is 8081")
	flag.IntVar(&config.C.RpcPort, "rpcPort", 6062, "httpPort, default is 6061")
	flag.Parse()
}

func init() {
	argParse()
	nacos.Init()
	kitex.InitClient()
	gorm.Init()
	minio.Init()
	redis.Init()
	snowflake.Init()
}

func main() {
	h := util_hertz.InitServer(config.C.HttpPort)
	interfaces.InitRouter(h)
	h.Spin()
}
