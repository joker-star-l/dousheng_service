package interfaces

import "github.com/cloudwego/hertz/pkg/app/server"

func InitRouter(h *server.Hertz) {
	feedRouter(h)
	publishRouter(h)
	commentRouter(h)
	favoriteRouter(h)
}
