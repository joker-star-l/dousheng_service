package interfaces

import (
	"context"
	"dousheng_service/video/application/service"
	"dousheng_service/video/interfaces/vo"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/joker-star-l/dousheng_common/config/jwt"
	common "github.com/joker-star-l/dousheng_common/entity"
	"strconv"
)

func publishRouter(h *server.Hertz) {
	r := h.Group("/douyin/publish")
	// 需要认证
	r.Use(jwt.Middleware.MiddlewareFunc())
	{
		r.GET("/list/", func(c context.Context, ctx *app.RequestContext) {
			tokenUser, _ := ctx.Get(jwt.KeyIdentity)
			userId, _ := strconv.ParseInt(tokenUser.(map[string]any)["id"].(string), 10, 0)
			queryId, err := strconv.ParseInt(ctx.Query("user_id"), 10, 0)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse("参数错误"))
				return
			}
			publishList, err := service.GetPublishList(userId, queryId)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse(err.Error()))
				return
			}
			ctx.JSON(consts.StatusOK, vo.VideoInfoListResponse{
				Response:  common.SuccessResponse(),
				VideoList: publishList,
			})
		})
		r.POST("/action/", func(c context.Context, ctx *app.RequestContext) {
			tokenUser, _ := ctx.Get(jwt.KeyIdentity)
			userId, _ := strconv.ParseInt(tokenUser.(map[string]any)["id"].(string), 10, 0)
			title := ctx.PostForm("title")
			file, err := ctx.FormFile("data")
			if title == "" || file == nil || err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse("参数错误"))
				return
			}
			err = service.Publish(userId, title, file)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse(err.Error()))
				return
			}
			ctx.JSON(consts.StatusOK, common.SuccessResponse())
		})
	}
}
