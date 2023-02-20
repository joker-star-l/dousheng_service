package interfaces

import (
	"context"
	"dousheng_service/message/application/service"
	"dousheng_service/message/interfaces/vo"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/joker-star-l/dousheng_common/config/jwt"
	common "github.com/joker-star-l/dousheng_common/entity"
	"strconv"
)

func messageRouter(h *server.Hertz) {
	r := h.Group("/douyin/message")
	{
		// 需要认证
		r.Use(jwt.Middleware.MiddlewareFunc())
		r.POST("/action/", func(c context.Context, ctx *app.RequestContext) {
			tokenUser, _ := ctx.Get(jwt.KeyIdentity)
			from, _ := strconv.ParseInt(tokenUser.(map[string]any)["id"].(string), 10, 0)
			to, err := strconv.ParseInt(ctx.Query("to_user_id"), 10, 0)
			actionType := ctx.Query("action_type")
			if err != nil || actionType != "1" {
				ctx.JSON(consts.StatusOK, common.ErrorResponse("参数错误"))
				return
			}
			content := ctx.Query("content")
			if content == "" {
				ctx.JSON(consts.StatusOK, common.ErrorResponse("不能发送空消息"))
				return
			}
			err = service.SendMessage(from, to, content)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse(err.Error()))
				return
			}
			ctx.JSON(consts.StatusOK, common.SuccessResponse())
		})
		r.GET("/chat/", func(c context.Context, ctx *app.RequestContext) {
			tokenUser, _ := ctx.Get(jwt.KeyIdentity)
			from, _ := strconv.ParseInt(tokenUser.(map[string]any)["id"].(string), 10, 0)
			to, err := strconv.ParseInt(ctx.Query("to_user_id"), 10, 0)
			time, err := strconv.ParseInt(ctx.Query("pre_msg_time"), 10, 0)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse("参数错误"))
				return
			}
			messageList, err := service.GetMessageList(from, to, time)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse(err.Error()))
				return
			}
			ctx.JSON(consts.StatusOK, vo.MessageListResponse{
				Response:    common.SuccessResponse(),
				MessageList: messageList,
			})
		})
	}
}
