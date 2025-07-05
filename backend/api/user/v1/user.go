package v1

import "github.com/gogf/gf/v2/frame/g"

// CreateReq 创建用户请求结构体
type CreateReq struct {
	g.Meta          `path:"/users" method:"post" tags:"用户模块" sm:"创建用户"`
	Username        string `p:"username" v:"required|length:3,12" dc:"用户名"`
	Email           string `p:"email" v:"required|email" dc:"邮箱"`
	Password        string `p:"password" v:"required|length:6,16" dc:"密码"`
	ConfirmPassword string `p:"confirmPassword" v:"required|same:Password#两次密码需一致" dc:"确认密码"`
}

// CreateRes 创建用户返回结构体
type CreateRes struct {
	g.Meta   `mime:"application/json"`
	UserId   uint64 `json:"userId" dc:"用户ID"`
	Username string `json:"username" dc:"用户名"`
}

type LoginReq struct {
	g.Meta   `path:"/auth/login" method:"post" tags:"用户模块" sm:"登录"`
	Username string `p:"username" v:"required|length:3,12" dc:"用户名"`
	Password string `p:"password" v:"required|length:6,16" dc:"密码"`
}

type LoginRes struct {
	AccessToken  string `json:"accessToken" dc:"访问令牌"`
	RefreshToken string `json:"refreshToken" dc:"刷新令牌"`
}

type MeReq struct {
	g.Meta `path:"/users/me" method:"get" tags:"用户模块" sm:"获取当前登录的用户信息"`
}

type MeRes struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type RefreshTokenReq struct {
	g.Meta       `path:"/auth/refresh" method:"post" tags:"用户模块" sm:"刷新令牌"`
	RefreshToken string `p:"refreshToken" v:"required" dc:"刷新令牌"`
}

type RefreshTokenRes struct {
	AccessToken  string `json:"accessToken" dc:"访问令牌"`
	RefreshToken string `json:"refreshToken" dc:"刷新令牌"`
}
