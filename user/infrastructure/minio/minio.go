package minio

import (
	"dousheng_service/user/config"
	"fmt"
	"github.com/joker-star-l/dousheng_common/config/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minio.Client

var DefaultAvatarAddress string
var DefaultBackgroundAddress string

func Init() {
	var err error
	Client, err = minio.New(config.C.Minio.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.C.Minio.AccessKeyID, config.C.Minio.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Slog.Panicln(err.Error())
	}

	DefaultAvatarAddress = "avatar/default.png"
	DefaultBackgroundAddress = "background/default.png"
}

func GetFullAddress(path string) string {
	return fmt.Sprintf("http://%s/%s/%s", config.C.Minio.EndPoint, config.C.Minio.Bucket, path)
}
