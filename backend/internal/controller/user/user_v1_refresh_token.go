package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "backend/api/user/v1"
)

func (c *ControllerV1) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (res *v1.RefreshTokenRes, err error) {
	// 前置判断下 refreshToken 是否在黑名单中，如果在直接返回错误
	output, err := c.svc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		g.Log().Errorf(ctx, "刷新token失败: %v", err)
		return nil, gerror.New("刷新token失败")
	}
	// 让旧的 refreshToken 过期, 防止重放攻击
	// 可以利用 Redis Set 设置一个黑名单，把旧的 refreshToken 加入黑名单
	return &v1.RefreshTokenRes{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}, nil
}
