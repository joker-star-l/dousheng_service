package snowflake

import (
	"dousheng_service/user/infrastructure/config"
	"github.com/bwmarrin/snowflake"
)

var Snowflake *snowflake.Node

func Init() {
	Snowflake, _ = snowflake.NewNode(int64(config.C.MachineId))
}

func GenerateId() int64 {
	return Snowflake.Generate().Int64()
}
