package points

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "backend/api/points/v1"
	"backend/internal/consts"
)

func (c *ControllerV1) Summary(ctx context.Context, req *v1.SummaryReq) (res *v1.SummaryRes, err error) {
	// 1. 获取登录的用户id
	// 1. 从请求上下文中获取 userid
	userId, ok := ctx.Value(consts.CtxKeyUserID).(uint64)
	g.Log().Debugf(ctx, "从请求上下文中获取 userId: %d", userId)
	if !ok || userId == 0 {
		return nil, gerror.New("用户信息获取失败")
	}
	// 2. 调用 积分服务查询总积分
	total, err := c.svc.Summary(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &v1.SummaryRes{
		Total: total,
	}, nil
}
