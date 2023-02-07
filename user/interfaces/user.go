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

func userRouter(h *server.Hertz) {
	r := h.Group("/douyin/user")
	{
		// 无需认证
		r.POST("/register", func(c context.Context, ctx *app.RequestContext) {
			user := &vo.LoginUser{}
			err := ctx.BindAndValidate(user)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse(err.Error()))
				return
			}
			tokenUser, err := service.Register(user)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse(err.Error()))
				return
			}
			token, _, _ := jwt.Middleware.TokenGenerator(tokenUser)
			ctx.JSON(consts.StatusOK, vo.UserLoginResponse{
				Response: common.SuccessResponse(),
				UserId:   tokenUser.Id,
				Token:    token,
			})
		})
		r.POST("/login", func(c context.Context, ctx *app.RequestContext) {
			user := &vo.LoginUser{}
			err := ctx.BindAndValidate(user)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse(err.Error()))
				return
			}
			tokenUser, err := service.Login(user)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse(err.Error()))
				return
			}
			token, _, _ := jwt.Middleware.TokenGenerator(tokenUser)
			ctx.JSON(consts.StatusOK, vo.UserLoginResponse{
				Response: common.SuccessResponse(),
				UserId:   tokenUser.Id,
				Token:    token,
			})
		})

		// 需要认证
		r.Use(jwt.Middleware.MiddlewareFunc())
		r.GET("/", func(c context.Context, ctx *app.RequestContext) {
			userId, err := strconv.ParseInt(ctx.Query("user_id"), 10, 0)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse("参数错误"))
				return
			}
			userInfo, err := service.UserInfo(userId)
			if err != nil {
				ctx.JSON(consts.StatusOK, common.ErrorResponse(err.Error()))
				return
			}
			ctx.JSON(consts.StatusOK, vo.UserInfoResponse{
				Response: common.SuccessResponse(),
				UserInfo: *userInfo,
			})
		})
	}
}
