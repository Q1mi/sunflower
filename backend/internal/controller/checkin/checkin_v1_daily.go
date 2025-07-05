package checkin

import (
	v1 "backend/api/checkin/v1"
	"backend/internal/consts"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// Daily 每日签到接口实现
func (c *ControllerV1) Daily(ctx context.Context, req *v1.DailyReq) (*v1.DailyRes, error) {
	// 1. 从请求上下文中获取 userid
	userId, ok := ctx.Value(consts.CtxKeyUserID).(uint64)
	g.Log().Debugf(ctx, "从请求上下文中获取 userId: %d", userId)
	if !ok || userId == 0 {
		return nil, gerror.New("用户信息获取失败")
	}
	// 2. 调用服务层每日签到逻辑
	err := c.svc.Daily(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &v1.DailyRes{}, nil
}
