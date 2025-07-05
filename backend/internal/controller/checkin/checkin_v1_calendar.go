package checkin

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "backend/api/checkin/v1"
	"backend/internal/consts"
	"backend/internal/model"
)

func (c *ControllerV1) Calendar(ctx context.Context, req *v1.CalendarReq) (res *v1.CalendarRes, err error) {
	// 解析输入的参数
	t, err := time.Parse("2006-01", req.YearMonth)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter) // 参数错误
	}
	// 1. 从请求上下文中获取 userid
	userId, ok := ctx.Value(consts.CtxKeyUserID).(uint64)
	g.Log().Debugf(ctx, "从请求上下文中获取 userId: %d", userId)
	if !ok || userId == 0 {
		return nil, gerror.New("用户信息获取失败")
	}

	// 调用 service 层获取签到日历数据
	output, err := c.svc.MonthDetail(ctx, &model.MonthDetailInput{
		Year:   t.Year(),
		Month:  int(t.Month()),
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CalendarRes{
		Year:  t.Year(),
		Month: int(t.Month()),
		Detail: v1.DetailInfo{
			CheckedInDays:      output.CheckedInDays,
			RetroCheckedInDays: output.RetroCheckedInDays,
			IsCheckedInToday:   output.IsCheckedInToday,
			RemainRetroTimes:   output.RemainRetroTimes,
			ConsecutiveDays:    output.ConsecutiveDays,
		},
	}, nil
}
