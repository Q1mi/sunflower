package points

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "backend/api/points/v1"
	"backend/internal/consts"
	"backend/internal/model"
)

const defaultLimit = 10

func (c *ControllerV1) Records(ctx context.Context, req *v1.RecordsReq) (res *v1.RecordsRes, err error) {
	// 1. 从请求上下文中获取 userid
	userId, ok := ctx.Value(consts.CtxKeyUserID).(uint64)
	g.Log().Debugf(ctx, "从请求上下文中获取 userId: %d", userId)
	if !ok || userId == 0 {
		return nil, gerror.New("用户信息获取失败")
	}
	if req.Limit > 50 {
		req.Limit = defaultLimit
	}
	// 2. 调用积分服务查询积分记录
	output, err := c.svc.Records(ctx, &model.PointsRecordsInput{
		UserId: userId,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		return nil, err
	}
	// 格式化, output.List 转为 []v1.Record
	res = &v1.RecordsRes{
		HasMore: output.HasMore,
		Total:   output.Total,
	}
	list := make([]*v1.Record, 0, len(output.List))
	for _, item := range output.List {
		list = append(list, &v1.Record{
			PointsChange:    item.Points,
			TransactionType: item.TransactionType,
			Description:     item.Description,
			TransactionTime: item.Date,
		})
	}
	res.List = list
	return res, nil
}
