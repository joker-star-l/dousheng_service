package interfaces

import (
	"context"
	"dousheng_service/user/application/service"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/joker-star-l/dousheng_common/config/jwt"
	common "github.com/joker-star-l/dousheng_common/entity"
	"strconv"
)

func relationRouter(h *server.Hertz) {
	r := h.Group("/douyin/relation")
	// 需要认证
	r.Use(jwt.Middleware.MiddlewareFunc())
	{
		r.POST("/action/", func(c context.Context, ctx *app.RequestContext) {
			tokenUser, _ := ctx.Get(jwt.KeyIdentity)
			from, _ := strconv.ParseInt(tokenUser.(map[string]any)["id"].(string), 10, 0)
			to, err := strconv.ParseInt(ctx.Query("to_user_id"), 10, 0)
			actionType := ctx.Query("action_type")
			if err != nil || !(actionType == "1" || actionType == "2") {
				ctx.JSON(consts.StatusOK, common.ErrorResponse("参数错误"))
				return
			}
			if actionType == "1" {
				err = service.Follow(from, to)
			} else if actionType == "2" {
				err = service.CancelFollow(from, to)
			}
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse(err.Error()))
				return
			}
			ctx.JSON(consts.StatusOK, common.SuccessResponse())
		})
	}
}
