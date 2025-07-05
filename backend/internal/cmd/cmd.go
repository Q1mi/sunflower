package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"backend/internal/controller/checkin"
	"backend/internal/controller/hello"
	"backend/internal/controller/points"
	"backend/internal/controller/user"
	"backend/internal/logic/middleware"
	"backend/utility/injection"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			// 服务注入
			injection.SetupDefaultInjector(ctx)
			defer injection.ShutdownDefaultInjector()

			s.Group("/api/v1", func(group *ghttp.RouterGroup) {
				// 注册通用响应中间件和CORS跨域中间件
				group.Middleware(ghttp.MiddlewareHandlerResponse, middleware.CORS)
				// 不需要登录也能访问的
				group.POST("/auth/login", user.NewV1(), "Login")          // 登录
				group.POST("/users", user.NewV1(), "Create")              // 创建用户
				group.POST("/auth/refresh", user.NewV1(), "RefreshToken") // 刷新token

				// 需要登录才能访问的
				group.Middleware(middleware.Auth)          // 用户认证中间件
				group.GET("/users/me", user.NewV1(), "Me") // 我的信息
				group.Bind(
					hello.NewV1(),
					checkin.NewV1(),
					points.NewV1(),
				)
			})
			s.Run()
			return nil
		},
	}
)
