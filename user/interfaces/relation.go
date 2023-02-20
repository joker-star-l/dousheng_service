package interfaces

import (
	"context"
	"dousheng_service/user/application/service"
	"dousheng_service/user/interfaces/vo"
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
			if from == to {
				ctx.JSON(consts.StatusOK, common.ErrorResponse("不能关注自己"))
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
		r.GET("/follow/list/", func(c context.Context, ctx *app.RequestContext) {
			tokenUser, _ := ctx.Get(jwt.KeyIdentity)
			userId, _ := strconv.ParseInt(tokenUser.(map[string]any)["id"].(string), 10, 0)
			queryId, err := strconv.ParseInt(ctx.Query("user_id"), 10, 0)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse("参数错误"))
				return
			}
			followList, err := service.GetFollowList(userId, queryId)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse(err.Error()))
				return
			}
			ctx.JSON(consts.StatusOK, vo.UserInfoListResponse{
				Response: common.SuccessResponse(),
				UserList: followList,
			})
		})
		r.GET("/follower/list/", func(c context.Context, ctx *app.RequestContext) {
			tokenUser, _ := ctx.Get(jwt.KeyIdentity)
			userId, _ := strconv.ParseInt(tokenUser.(map[string]any)["id"].(string), 10, 0)
			queryId, err := strconv.ParseInt(ctx.Query("user_id"), 10, 0)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse("参数错误"))
				return
			}
			followerList, err := service.GetFollowerList(userId, queryId)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse(err.Error()))
				return
			}
			ctx.JSON(consts.StatusOK, vo.UserInfoListResponse{
				Response: common.SuccessResponse(),
				UserList: followerList,
			})
		})
		r.GET("/friend/list/", func(c context.Context, ctx *app.RequestContext) {
			tokenUser, _ := ctx.Get(jwt.KeyIdentity)
			userId, _ := strconv.ParseInt(tokenUser.(map[string]any)["id"].(string), 10, 0)
			queryId, err := strconv.ParseInt(ctx.Query("user_id"), 10, 0)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse("参数错误"))
				return
			}
			if userId != queryId {
				ctx.JSON(consts.StatusOK, common.ErrorResponse("不能查看他人的好友"))
				return
			}
			friendList, err := service.GetFriendList(queryId)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse(err.Error()))
				return
			}
			ctx.JSON(consts.StatusOK, vo.FriendInfoListResponse{
				Response:   common.SuccessResponse(),
				FriendInfo: friendList,
			})
		})
	}
}
