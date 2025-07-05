package impl

import (
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/model/entity"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// 积分服务的具体实现

type Service struct{}

func New() *Service {
	return &Service{}
}

// Summary 总积分查询
func (s *Service) Summary(ctx context.Context, userId uint64) (int, error) {
	var userPoint entity.UserPoints
	if err := dao.UserPoints.Ctx(ctx).
		Where(dao.UserPoints.Columns().UserId, userId).
		Scan(&userPoint); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil // 没有积分记录，返回0
		}
		g.Log().Errorf(ctx, "查询积分失败: %v", err)
		return 0, err
	}
	return int(userPoint.Points), nil
}

// Records 积分记录查询
func (s *Service) Records(ctx context.Context, input *model.PointsRecordsInput) (*model.PointsRecordsOutput, error) {
	// 1. 分页查询积分记录
	var (
		total   int
		records []entity.UserPointsTransactions
	)
	if err := dao.UserPointsTransactions.Ctx(ctx).
		Where(dao.UserPointsTransactions.Columns().UserId, input.UserId).
		OrderDesc(dao.UserPointsTransactions.Columns().CreatedAt). // 创建时间倒序
		Offset(input.Offset).
		Limit(input.Limit).
		ScanAndCount(&records, &total, false); err != nil {
		g.Log().Errorf(ctx, "查询积分记录失败: %v", err)
		return nil, err
	}
	// 2. 格式化输出,把数据库中记录格式化为需要的数据
	list := make([]*model.PointsRecordItem, 0, len(records))
	for _, v := range records {
		list = append(list, &model.PointsRecordItem{
			Points:          v.PointsChange,
			TransactionType: v.TransactionType,
			Description:     v.Description,
			Date:            v.CreatedAt.Time.Format(time.DateTime),
		})
	}
	return &model.PointsRecordsOutput{
		List:    list,
		HasMore: len(records) == input.Limit && total > (input.Offset+input.Limit),
		Total:   total,
	}, nil
}
