package snowflake

import (
	"dousheng_service/video/infrastructure/config"
	"github.com/bwmarrin/snowflake"
	"github.com/joker-star-l/dousheng_common/config/log"
)

var Snowflake *snowflake.Node

func Init() {
	var err error
	Snowflake, err = snowflake.NewNode(int64(config.C.MachineId))
	if err != nil {
		log.Slog.Panic(err)
	}
}

func GenerateId() int64 {
	return Snowflake.Generate().Int64()
}
