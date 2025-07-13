package impl

import (
	"backend/utility/injection"
	"context"
	"fmt"
	"time"

	_ "embed"

	"github.com/redis/go-redis/v9"
)

//go:embed remind.lua
var remindScript string

// CheckAndNotify 检查签到并发送通知
func CheckAndNotify(ctx context.Context, remindThreshold int) error {
	// 1. 获取所有符合条件的用户(避免遍历全量用户)
	// 可以通过扫 MySQL 用户表找到最近一天有登录过的
	// 或者可以在用户签到的时候记录一个 ZSet, userID:签到时间,如果用户量多需要拆分 Key
	userIDs := []uint64{25016147980058993}
	rc := injection.MustInvoke[*redis.Client]()
	// 2. 加载 lua script
	sha, err := rc.ScriptLoad(ctx, remindScript).Result()
	if err != nil {
		fmt.Printf("ScriptLoad err: %v\n", err)
		return err
	}
	// 3. 遍历判断每个用户
	for _, userID := range userIDs {
		key := fmt.Sprintf(yearSignKeyFormat, userID, time.Now().Year())

		// 计算当前日偏移量(当年第几天)
		dayOfYearOffset := time.Now().YearDay() - 1
		fmt.Printf("key: %s, dayOfYearOffset: %d\n", key, dayOfYearOffset)
		// 执行LUA脚本
		result, err := rc.EvalSha(ctx, sha, []string{key}, dayOfYearOffset, remindThreshold).Int()
		fmt.Printf("result: %d, err: %v\n", result, err)
		if err != nil {
			return err
		}

		if result == 1 {
			fmt.Printf("用户%d需要发送断签提醒\n", userID)
			// 发送到消息队列，执行后续的推送逻辑（APP Push 或 短信等）
		}
	}
	return nil
}
