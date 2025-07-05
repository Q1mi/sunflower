package user

import (
	"backend/internal/model"
	"backend/internal/model/entity"
	"context"
)

// 用户相关业务逻辑

// 把用户服务抽象成一个接口，列出来所有需要实现的方法
type UserInfoService interface {
	Create(ctx context.Context, input *model.CreateUserInput) (*model.CreateUserOutput, error)
	Login(ctx context.Context, input *model.LoginInput) (*model.LoginOutput, error)
	GetInfo(ctx context.Context, userId string) (*entity.Userinfo, error)
	RefreshToken(ctx context.Context, refreshToken string) (res *model.TokenOutput, err error)
}
