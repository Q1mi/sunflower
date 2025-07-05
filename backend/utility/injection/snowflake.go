package injection

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/samber/do"
	"github.com/sony/sonyflake/v2"
)

func injectSnowflake(ctx context.Context, injector *do.Injector) {
	do.Provide(injector, func(i *do.Injector) (*sonyflake.Sonyflake, error) {
		// 完成 *sonyflake.Sonyflake 的初始化
		// 1. 读取配置文件中的起始时间
		type SnowflakeCfg struct {
			StartTime string `yaml:"startTime"`
		}
		var (
			err error
			cfg *SnowflakeCfg
		)
		err = g.Cfg().MustGet(ctx, "snowflake").Scan(&cfg)
		if err != nil {
			g.Log().Errorf(ctx, "manifest/config/config.yaml中必须指定snowflake的起始时间！: %v", err)
			return nil, err
		}
		if cfg == nil {
			return nil, fmt.Errorf("manifest/config/config.yaml中必须指定snowflake的起始时间！")
		}

		st, _ := time.Parse(time.DateOnly, cfg.StartTime)
		settings := sonyflake.Settings{
			StartTime: st,
		}
		sonyFlake, err := sonyflake.New(settings)
		if err != nil {
			g.Log().Errorf(ctx, "初始化sonyflake失败: %v", err)
			return nil, err
		}

		// 注册 sonyFlake 和一个关闭函数
		SetupShutdownHelper(injector, sonyFlake, func(svc *sonyflake.Sonyflake) error {
			return nil
		})
		return sonyFlake, nil
	})
}
