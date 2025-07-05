// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package checkin

import (
	"context"

	"backend/api/checkin/v1"
)

type ICheckinV1 interface {
	Daily(ctx context.Context, req *v1.DailyReq) (res *v1.DailyRes, err error)
	Calendar(ctx context.Context, req *v1.CalendarReq) (res *v1.CalendarRes, err error)
	Retro(ctx context.Context, req *v1.RetroReq) (res *v1.RetroRes, err error)
}
