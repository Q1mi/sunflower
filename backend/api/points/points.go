// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package points

import (
	"context"

	"backend/api/points/v1"
)

type IPointsV1 interface {
	Summary(ctx context.Context, req *v1.SummaryReq) (res *v1.SummaryRes, err error)
	Records(ctx context.Context, req *v1.RecordsReq) (res *v1.RecordsRes, err error)
}
