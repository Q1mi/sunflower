package middleware

import (
	"backend/internal/consts"
	"backend/internal/model"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/golang-jwt/jwt/v5"
)

// Auth 认证中间件
func Auth(r *ghttp.Request) {
	// 获取请求上下文
	ctx := r.GetCtx()
	// 从请求中获取用户id(根据Access token获取用户id)
	// 从请求头中获取 jwt access token
	authorizationValue := r.GetHeader("Authorization")
	if len(authorizationValue) == 0 || !strings.HasPrefix(authorizationValue, "Bearer ") {
		r.Response.WriteStatusExit(http.StatusUnauthorized, "缺少token") // 401 未认证  403 未授权
	}
	accessToken := strings.TrimPrefix(authorizationValue, "Bearer ")
	// 解析token获取用户id
	var claim model.JWTClaims
	token, err := jwt.ParseWithClaims(accessToken, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(consts.JWTAccessTokenSecret), nil
	})
	if err != nil || !token.Valid {
		g.Log().Errorf(ctx, "token: %v, err:%+v", token, err)
		r.Response.WriteStatusExit(http.StatusUnauthorized, "无效的token") // 401 未认证  403 未授权
	}
	g.Log().Debugf(ctx, "claim: %v", claim)
	// 向请求的上下文中写入用户id

	r.SetCtxVar(consts.CtxKeyUserID, claim.UserId) // 使用自定义Key类型，防止被其他中间件覆盖
	// r.SetCtxVar("userId", "xxx")

	r.Middleware.Next() // 继续执行后续中间件
}
