package interfaces

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func feedRouter(h *server.Hertz) {
	r := h.Group("/douyin/feed")
	{
		r.GET("/", func(c context.Context, ctx *app.RequestContext) {
			ctx.JSON(consts.StatusOK, utils.H{
				"status_code": 0,
				"status_msg":  "success",
				"video_list":  []any{},
			})
		})
	}
}
