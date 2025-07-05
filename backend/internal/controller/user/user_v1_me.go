package user

import (
	"context"
	"strconv"

	v1 "backend/api/user/v1"
	"backend/internal/consts"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) Me(ctx context.Context, req *v1.MeReq) (res *v1.MeRes, err error) {
	// 从请求上下文中获取 userid
	userId, ok := ctx.Value(consts.CtxKeyUserID).(uint64)
	g.Log().Debugf(ctx, "从请求上下文中获取 userId: %d", userId)
	if !ok || userId == 0 {
		return nil, gerror.New("用户信息获取失败")
	}

	// 根据用户id获取用户信息
	userInfo, err := c.svc.GetInfo(ctx, strconv.FormatUint(userId, 10))
	if err != nil {
		return nil, gerror.New("获取用户信息失败")
	}
	// 返回用户信息
	return &v1.MeRes{
		Username: userInfo.Username,
		Avatar:   userInfo.Avatar,
		Email:    userInfo.Email,
	}, err
}
