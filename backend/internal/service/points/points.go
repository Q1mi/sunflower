package points

import (
	"backend/internal/model"
	"context"
)

type Service interface {
	Summary(ctx context.Context, userId uint64) (int, error)
	Records(ctx context.Context, input *model.PointsRecordsInput) (*model.PointsRecordsOutput, error)
}
