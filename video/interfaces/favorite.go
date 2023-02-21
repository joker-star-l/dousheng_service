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

func favoriteRouter(h *server.Hertz) {
	r := h.Group("/douyin/favorite")
	{
		r.GET("/list/", func(c context.Context, ctx *app.RequestContext) {
			userId := jwt.ParseAndGetUserId(c, ctx)
			queryId, err := strconv.ParseInt(ctx.Query("user_id"), 10, 0)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse("参数错误"))
				return
			}
			favoriteList, err := service.GetFavoriteList(userId, queryId)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse(err.Error()))
				return
			}
			ctx.JSON(consts.StatusOK, vo.VideoInfoListResponse{
				Response:  common.SuccessResponse(),
				VideoList: favoriteList,
			})
		})

		r.Use(jwt.Middleware.MiddlewareFunc())
		r.POST("/action/", func(c context.Context, ctx *app.RequestContext) {
			userId := jwt.GetUserId(ctx)
			videoId, err := strconv.ParseInt(ctx.Query("video_id"), 10, 0)
			actionType := ctx.Query("action_type")
			if err != nil || !(actionType == "1" || actionType == "2") {
				ctx.JSON(consts.StatusOK, common.ErrorResponse("参数错误"))
				return
			}
			if actionType == "1" {
				err = service.Favorite(userId, videoId)
			} else if actionType == "2" {
				err = service.CancelFavorite(userId, videoId)
			}
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse(err.Error()))
				return
			}
			ctx.JSON(consts.StatusOK, common.SuccessResponse())
		})
	}
}
