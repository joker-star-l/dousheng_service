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

func commentRouter(h *server.Hertz) {
	r := h.Group("/douyin/comment")
	{
		r.GET("/list/", func(c context.Context, ctx *app.RequestContext) {
			userId := jwt.ParseAndGetUserId(c, ctx)
			videoId, err := strconv.ParseInt(ctx.Query("video_id"), 10, 0)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse("参数错误"))
				return
			}
			commentList, err := service.GetCommentList(userId, videoId)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse(err.Error()))
				return
			}
			ctx.JSON(consts.StatusOK, vo.CommentListResponse{
				Response:    common.SuccessResponse(),
				CommentList: commentList,
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
			comment := &vo.Comment{}
			if actionType == "1" {
				commentText := ctx.Query("comment_text")
				if commentText == "" {
					ctx.JSON(consts.StatusOK, common.ErrorResponse("评论不能为空"))
					return
				}
				comment, err = service.Comment(userId, videoId, commentText)
			} else if actionType == "2" {
				commentId, err := strconv.ParseInt(ctx.Query("comment_id"), 10, 0)
				if err != nil {
					ctx.JSON(consts.StatusOK, common.ErrorResponse("参数错误"))
					return
				}
				err = service.DeleteComment(userId, commentId)
			}
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse(err.Error()))
				return
			}
			ctx.JSON(consts.StatusOK, vo.CommentResponse{
				Response: common.SuccessResponse(),
				Comment:  *comment,
			})
		})
	}
}
