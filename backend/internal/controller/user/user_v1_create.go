package user

import (
	"context"

	v1 "backend/api/user/v1"
	"backend/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

// Create 用户注册流程
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (*v1.CreateRes, error) {
	// 参数校验（框架已经做完了）
	input := model.CreateUserInput{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}
	// 调用 service 层的逻辑
	output, err := c.svc.Create(ctx, &input)
	if err != nil {
		g.Log().Errorf(ctx, "Create user failed: %v", err)
		return nil, err
	}
	// 返回结果
	return &v1.CreateRes{
		UserId:   output.UserId,
		Username: output.Username,
	}, nil
}
