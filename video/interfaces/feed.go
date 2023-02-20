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
	"time"
)

func feedRouter(h *server.Hertz) {
	r := h.Group("/douyin/feed")
	{
		// 无需认证
		r.GET("/", func(c context.Context, ctx *app.RequestContext) {
			// 解析token
			userId := int64(0)
			claims, err := jwt.Middleware.GetClaimsFromJWT(c, ctx)
			if err == nil && claims != nil {
				userId, _ = strconv.ParseInt(claims[jwt.KeyData].(map[string]any)["id"].(string), 10, 0)
			}
			// 解析时间
			latestTime := time.Now()
			timestamp, err := strconv.ParseInt(ctx.Query("latest_time"), 10, 0)
			if err == nil {
				latestTime = time.Unix(timestamp/1000, 0)
			}
			// 处理
			videoInfoList, err := service.Feed(userId, latestTime)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse(err.Error()))
				return
			}
			response := vo.VideoInfoListResponse{
				Response:  common.SuccessResponse(),
				VideoList: videoInfoList,
			}
			if len(videoInfoList) > 0 {
				response.NextTime = videoInfoList[0].CreateTime
			}
			ctx.JSON(consts.StatusOK, response)
		})
	}
}
