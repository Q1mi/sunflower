package model

import "github.com/golang-jwt/jwt/v5"

// 定义 controller 层 与 service 层 之间交互的数据
// gf 框架推荐使用 input/output 结构体封装交互的数据

type CreateUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type CreateUserOutput struct {
	UserId   uint64 `json:"userId"`
	Username string `json:"username"`
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginOutput struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenOutput struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// JWTClaims 自定义声明结构体并内嵌 jwt.RegisteredClaims
type JWTClaims struct {
	UserId   uint64 `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
