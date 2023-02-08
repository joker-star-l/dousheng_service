package gorm

import (
	"dousheng_service/user/infrastructure/config"
	"github.com/joker-star-l/dousheng_common/config/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(postgres.Open(config.C.Dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Slog.Errorf("db init error: %v", err)
		panic(any(err))
	}
}
