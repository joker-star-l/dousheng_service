package config

import "github.com/bwmarrin/snowflake"

var Snowflake *snowflake.Node

func init() {
	Snowflake, _ = snowflake.NewNode(int64(C.MachineId))
}

func GenerateId() int64 {
	return Snowflake.Generate().Int64()
}
