package impl

import (
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/model/entity"
	"backend/utility/injection"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/redis/go-redis/v9"
)

// 签到相关业务逻辑的具体实现

const (
	yearSignKeyFormat   = "user:checkins:daily:%d:%d"      // user:checkins:daily:12131321421312:2025
	monthRetroKeyFormat = "user:checkins:retro:%d:%d:%02d" // user:checkins:retro:12131321421312:2025:01

	defaultDailyPoints     = 1   // 每日签到积分
	defaultRetroCostPoints = 100 // 补签消耗积分

	maxRetroTimesPerMonth = 3 // 单月最多补签次数
)

type PointsTransactionType int

const (
	PointsTransactionTypeDaily       = 1 // 每日签到 1
	PointsTransactionTypeConsecutive = 2 // 连续签到 2
	PointsTransactionTypeRetro       = 3 // 补签 3
)

type ConsecutiveBonusType int32

const (
	// 连续签到奖励规则
	consecutiveBonus3  ConsecutiveBonusType = 1 // "连续签到3天奖励"
	consecutiveBonus7  ConsecutiveBonusType = 2 // "连续签到7天奖励"
	consecutiveBonus15 ConsecutiveBonusType = 3 // "连续签到15天奖励"
	consecutiveBonus30 ConsecutiveBonusType = 4 //"月度满签奖励"
)

var consecutiveBonusNames = map[ConsecutiveBonusType]string{
	consecutiveBonus3:  "连续签到3天奖励",
	consecutiveBonus7:  "连续签到7天奖励",
	consecutiveBonus15: "连续签到15天奖励",
	consecutiveBonus30: "月度满签奖励",
}

var PointsTransactionTypeMsgMap = map[PointsTransactionType]string{
	PointsTransactionTypeDaily:       "每日签到奖励",
	PointsTransactionTypeConsecutive: "连续签到奖励",
	PointsTransactionTypeRetro:       "补签%s消耗",
}

// consecutiveBonusRule 连续签到奖励规则
type consecutiveBonusRule struct {
	TriggerDays int                  // 触发连续签到奖励的天数
	Points      int64                // 连续签到奖励的积分
	BonusType   ConsecutiveBonusType // 连续签到奖励类型
}

var consecutiveBonusRules = []consecutiveBonusRule{
	{TriggerDays: 3, Points: 5, BonusType: consecutiveBonus3},
	{TriggerDays: 7, Points: 10, BonusType: consecutiveBonus7},
	{TriggerDays: 15, Points: 20, BonusType: consecutiveBonus15},
	{TriggerDays: 30, Points: 100, BonusType: consecutiveBonus30},
}

var (
	ErrInvalidRetroDate = errors.New("补签日期无效")
	ErrChecked          = errors.New("日期已签到")
	ErrRetroNotimes     = errors.New("本月补签次数已用完")

	ErrNoEnouthPoints = gerror.New("积分不足")
)

type Service struct {
	rc *redis.Client
}

func NewService() *Service {
	return &Service{
		rc: injection.MustInvoke[*redis.Client](), // 从注入器中获取 Redis 客户端实例
	}
}

// Daily 每日签到
func (s *Service) Daily(ctx context.Context, userId uint64) error {
	// 采用服务器时间进行每日签到，不依赖客户端传递的时间
	// 1. Redis 中使用 bitmap setbit 执行签到逻辑
	// 拿到当天是一年中的第几天，然后使用 setbit 记录这一天是否签到
	now := time.Now()
	year := now.Year()
	dayOfYearOffset := now.YearDay() - 1 // 因为 Redis bitmap 从 0 开始，所以要减一
	key := fmt.Sprintf(yearSignKeyFormat, userId, year)
	g.Log().Debugf(ctx, "key: %s dayOfYearOffset:%d", key, dayOfYearOffset)

	ret := s.rc.SetBit(ctx, key, int64(dayOfYearOffset), 1).Val()
	if ret == 1 {
		return errors.New("今日已签到")
	}

	// 2. 发放每日签到的积分
	err := AddPoints(ctx, &model.PointsTransactionInput{
		UserId: userId,
		Points: defaultDailyPoints,
		Desc:   PointsTransactionTypeMsgMap[PointsTransactionTypeDaily],
		Type:   int(PointsTransactionTypeDaily),
	})
	if err != nil {
		g.Log().Errorf(ctx, "事务处理失败: %v", err)
		return err
	}

	// 3. 发送连续签到的奖励积分
	return s.updateConsecutiveBonus(ctx, userId, year, int(now.Month()))
}

// MonthDetail 签到详情
func (s *Service) MonthDetail(ctx context.Context, input *model.MonthDetailInput) (*model.MonthDetailOutput, error) {
	// 1. 从redis中分别取出签到bitmap和补签bitmap,分别得到签到日期和补签日期
	checkinBitmap, retroBitmap, err := s.getMonthBitmap(ctx, input.UserId, input.Year, input.Month)
	if err != nil {
		g.Log().Errorf(ctx, "获取年月bitmap失败: %v", err)
		return nil, err
	}
	g.Log().Debugf(ctx, "--> checkinBitmap: %031b retroBitmap:%031b", checkinBitmap, retroBitmap)
	monthDays := getMonthDays(input.Year, input.Month) // 当月天数
	checkinDays := parseBitmap2Days(checkinBitmap, monthDays)
	retroDays := parseBitmap2Days(retroBitmap, monthDays)

	// 2. 计算连续签到天数
	bitmap := checkinBitmap | retroBitmap
	maxConsecutive := calcMaxConsecutiveDays(bitmap, monthDays)
	// 3. 计算剩余补签次数
	remainRetroTimes := maxRetroTimesPerMonth - len(retroDays) // 用月度补签次数减去已补签天数
	// 4. 计算当天是否签到
	isCheckedToday, err := s.IsCheckedToday(ctx, input.UserId)
	if err != nil {
		g.Log().Errorf(ctx, "查询当天是否签到失败: %v", err)
		return nil, err
	}
	return &model.MonthDetailOutput{
		CheckedInDays:      checkinDays,
		RetroCheckedInDays: retroDays,
		ConsecutiveDays:    maxConsecutive,
		RemainRetroTimes:   remainRetroTimes,
		IsCheckedInToday:   isCheckedToday,
	}, nil
}

// Retro 根据输入的日期进行补签
func (s *Service) Retro(ctx context.Context, userId uint64, date time.Time) error {
	// 1. 判断补签日期是否有效
	if err := s.checkRetroDate(ctx, userId, date); err != nil {
		return err
	}
	// 2. 执行补签逻辑
	// 2.1 Redis里 补签的月度bitmap 中设置补签的标识
	retroKey := fmt.Sprintf(monthRetroKeyFormat, userId, date.Year(), date.Month())
	retroOffset := date.Day() - 1 // 索引是从0开始的，所以要减1
	err := s.rc.SetBit(ctx, retroKey, int64(retroOffset), 1).Err()
	if err != nil {
		g.Log().Errorf(ctx, "SetBit 设置补签状态失败: %v", err)
		return gerror.NewCode(gcode.CodeInternalError)
	}
	// 2.2 补签消耗积分、增加积分、增加积分记录
	// 正常应该把签到服务和积分服务分开，通过消息队列的方式实现事件驱动。
	// 签到服务负责签到/补签，发出消息；积分服务监听消息，处理积分的增加和扣减逻辑。
	if err := s.retroWithTransaction(ctx, userId, date); err != nil {
		// 如果数据库更新失败，则回滚 Redis 中的补签标识
		err := s.rc.SetBit(ctx, retroKey, int64(retroOffset), 0).Err()
		if err != nil {
			g.Log().Errorf(ctx, "SetBit 回滚补签状态失败: %v", err)
			return gerror.NewCode(gcode.CodeInternalError)
		}
	}
	// 3. 计算连续签到日期发放连续签到奖励
	return s.updateConsecutiveBonus(ctx, userId, date.Year(), int(date.Month()))
}

// retroWithTransaction 补签逻辑，使用事务保证原子性
func (s *Service) retroWithTransaction(ctx context.Context, userId uint64, date time.Time) error {
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1. 查询用户的当前积分，积分不够的不能补签
		var userPoint entity.UserPoints
		if err := tx.Model(dao.UserPoints.Table()).
			Where(dao.UserPoints.Columns().UserId, userId).
			Scan(&userPoint); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				g.Log().Errorf(ctx, "查询用户积分失败: %v", err)
				return err
			}
			userPoint = entity.UserPoints{
				UserId: userId,
			}
		}
		if userPoint.Points < defaultRetroCostPoints {
			return ErrNoEnouthPoints
		}
		// 2. 计算积分变化（补签消耗的、每日签到奖励）
		pointsChange := -defaultRetroCostPoints + defaultDailyPoints
		nowPoints := userPoint.Points + int64(pointsChange)
		nowTotalPoints := userPoint.PointsTotal + defaultDailyPoints // 只算得到的积分，不算消费的
		// 3. 积分记录中新增一条补签消耗100积分的记录
		retroCostRecord := entity.UserPointsTransactions{
			UserId:          userId,
			PointsChange:    -defaultRetroCostPoints,
			TransactionType: PointsTransactionTypeRetro,
			Description:     fmt.Sprintf(PointsTransactionTypeMsgMap[PointsTransactionTypeRetro], date.Format(time.DateOnly)),
			CurrentBalance:  userPoint.Points - defaultRetroCostPoints,
			CreatedAt:       gtime.NewFromTime(time.Now()),
			UpdatedAt:       gtime.NewFromTime(time.Now()),
		}
		if _, err := tx.Model(dao.UserPointsTransactions.Table()).Insert(&retroCostRecord); err != nil {
			g.Log().Errorf(ctx, "插入补签消耗的积分记录失败: %v", err)
			return err
		}
		// 4. 积分记录中新增每日签到固定得的奖励积分记录
		checkinBonusRecord := entity.UserPointsTransactions{
			UserId:          userId,
			PointsChange:    defaultDailyPoints,
			TransactionType: PointsTransactionTypeDaily,
			Description:     PointsTransactionTypeMsgMap[PointsTransactionTypeDaily],
			CurrentBalance:  nowPoints,
			CreatedAt:       gtime.NewFromTime(time.Now()),
			UpdatedAt:       gtime.NewFromTime(time.Now()),
		}
		if _, err := tx.Model(dao.UserPointsTransactions.Table()).Insert(&checkinBonusRecord); err != nil {
			g.Log().Errorf(ctx, "插入补签奖励积分记录失败: %v", err)
			return err
		}
		// 5. 更新用户积分
		userPoint.Points = nowPoints
		userPoint.PointsTotal = nowTotalPoints
		if _, err := tx.Model(dao.UserPoints.Table()).
			Where(dao.UserPoints.Columns().UserId, userId).
			Update(&userPoint); err != nil {
			g.Log().Errorf(ctx, "更新用户积分失败: %v", err)
			return err
		}
		return nil
	})

}

func (s *Service) checkRetroDate(ctx context.Context, userId uint64, date time.Time) error {
	// 补签日期不能是今天或者未来的日期
	now := time.Now()
	if date.Year() > now.Year() ||
		date.Month() != now.Month() ||
		(date.Year() == now.Year() && date.YearDay() >= now.YearDay()) {
		return ErrInvalidRetroDate
	}
	// 补签的日期不能是本月之前的日期
	// 补签的日期不能是已经签到的日期(签到或者补签过都算)
	checkinKey := fmt.Sprintf(yearSignKeyFormat, userId, date.Year())
	yearOffset := date.YearDay() - 1
	checked, err := s.rc.GetBit(ctx, checkinKey, int64(yearOffset)).Result()
	if err != nil {
		g.Log().Errorf(ctx, "GetBit 获取当天签到状态失败: %v", err)
		return err
	}
	if checked == 1 {
		return ErrInvalidRetroDate
	}
	retroKey := fmt.Sprintf(monthRetroKeyFormat, userId, date.Year(), date.Month())
	retroOffset := date.Day() - 1
	retroRet, err := s.rc.GetBit(ctx, retroKey, int64(retroOffset)).Result()
	if err != nil {
		g.Log().Errorf(ctx, "GetBit 获取当天补签状态失败: %v", err)
		return err
	}
	if retroRet == 1 {
		return ErrInvalidRetroDate
	}
	// 每个月补签不能超过三次
	retroCount, err := s.rc.BitCount(ctx, retroKey, nil).Result()
	if err != nil {
		g.Log().Errorf(ctx, "BitCount 获取补签次数失败: %v", err)
		return err
	}
	if retroCount >= maxRetroTimesPerMonth {
		return ErrRetroNotimes
	}
	return nil
}

func (s *Service) IsCheckedToday(ctx context.Context, userId uint64) (bool, error) {
	// 计算今天的年度索引，然后使用 getbit 判断这一天是否签到
	now := time.Now()
	year := now.Year()
	key := fmt.Sprintf(yearSignKeyFormat, userId, year)
	dayOffset := now.YearDay() - 1
	value, err := s.rc.GetBit(ctx, key, int64(dayOffset)).Result()
	if err != nil {
		g.Log().Errorf(ctx, "GetBit 获取当天签到状态失败: %v", err)
		return false, err
	}
	return value == 1, nil
}

// parseBitmap2Days 根据当月的天数和bitmap，输出对应的签到/补签日期
func parseBitmap2Days(bitmap uint64, monthDays int) []int {
	days := make([]int, 0)
	for i := range monthDays {
		// 0010000000000000000000000000000
		// 1000000000000000000000000000000
		if (bitmap & (1 << (monthDays - 1 - i))) != 0 {
			days = append(days, i+1)
		}
	}
	return days
}

// updateConsecutiveBonus 更新连续签到奖励积分
func (s *Service) updateConsecutiveBonus(ctx context.Context, userId uint64, year, month int) error {
	// 1. 获取当前连续签到天数
	maxConsecutive, err := s.CalcMonthConsecutiveDays(ctx, userId, year, month)
	if err != nil {
		g.Log().Errorf(ctx, "计算连续签到天数失败: %v", err)
		return err
	}
	// 2. 计算连续签到奖励积分
	// 3. 更新用户积分汇总表和用户积分明细表
	// 如何避免重复发放连续签到奖励？ --> 使用 user_monthly_bonus_log 表记录用户指定月份已领取的奖励
	// 2.1 先查询用户已领取的奖励
	var bonusLogs []*entity.UserMonthlyBonusLog
	if err := dao.UserMonthlyBonusLog.Ctx(ctx).
		Where(dao.UserMonthlyBonusLog.Columns().UserId, userId).
		Where(dao.UserMonthlyBonusLog.Columns().YearMonth, fmt.Sprintf("%d%02d", year, month)). // 202505
		Scan(&bonusLogs); err != nil && !errors.Is(err, sql.ErrNoRows) {
		g.Log().Errorf(ctx, "查询用户已领取的奖励失败: %v", err)
		return err
	}
	// 把领取的奖励塞到map中，方便后续判断是否已经发放过奖励
	bonusLogsMap := make(map[ConsecutiveBonusType]bool)
	for _, v := range bonusLogs {
		bonusLogsMap[ConsecutiveBonusType(v.BonusType)] = true
	}
	// 遍历连续签到奖励配置，如果符合条件就发送奖励
	for _, rule := range consecutiveBonusRules {
		if maxConsecutive >= rule.TriggerDays && !bonusLogsMap[rule.BonusType] {
			// 发放连续签到奖励积分
			// 更新 user_points 表和 user_points_transactions 表
			if err := AddPoints(ctx, &model.PointsTransactionInput{
				UserId: userId,
				Points: rule.Points,
				Desc:   consecutiveBonusNames[rule.BonusType],
				Type:   int(PointsTransactionTypeConsecutive),
			}); err != nil {
				g.Log().Errorf(ctx, "发放连续签到奖励失败: %v", err)
				continue
			}
			// 记录到 user_monthly_bonus_log 表
			newLog := &entity.UserMonthlyBonusLog{
				UserId:      userId,
				YearMonth:   fmt.Sprintf("%d%02d", year, month),
				Description: consecutiveBonusNames[rule.BonusType],
				BonusType:   int(rule.BonusType),
				CreatedAt:   gtime.NewFromTime(time.Now()),
				UpdatedAt:   gtime.NewFromTime(time.Now()),
			}
			if _, err := dao.UserMonthlyBonusLog.Ctx(ctx).Insert(newLog); err != nil {
				// 积分已加，连续签到奖励已发，但月度奖励记录插入失败，需要手动处理
				g.Log().Errorf(ctx, "[NEED_HANDLE]插入用户月度奖励记录失败: %v", err)
				continue
			}
		}
	}

	return nil
}

// getMonthBitmap 获取当月 签到bitmap 和 补签的bitmap
func (s *Service) getMonthBitmap(ctx context.Context, userId uint64, year, month int) (uint64, uint64, error) {
	// 从用户年度签到记录中取出当月签到 bitmap
	key := fmt.Sprintf(yearSignKeyFormat, userId, year)
	firstOfMonthOffset := getFirstOfMonthOffset(year, month)
	monthDays := getMonthDays(year, month)
	bitWidthType := fmt.Sprintf("u%d", monthDays)
	values, err := s.rc.BitField(ctx, key, "GET", bitWidthType, firstOfMonthOffset).Result()
	if err != nil {
		g.Log().Errorf(ctx, "获取用户签到记录失败: %v", err)
		return 0, 0, err
	}
	if len(values) == 0 {
		values = []int64{0} // 如果没有查询到，则默认为0
	}
	checkinBitmap := uint64(values[0])
	g.Log().Debugf(ctx, "checkinBitmap: %0b", checkinBitmap)
	// 取出当月补签bitmap
	retroKey := fmt.Sprintf(monthRetroKeyFormat, userId, year, month)
	retroValues, err := s.rc.BitField(ctx, retroKey, "GET", bitWidthType, "#0").Result()
	if err != nil {
		g.Log().Errorf(ctx, "获取用户补签记录失败: %v", err)
		return 0, 0, err
	}
	if len(retroValues) == 0 {
		retroValues = []int64{0} // 没有查询到，则默认为0
	}
	retroBitmap := uint64(retroValues[0])
	return checkinBitmap, retroBitmap, nil
}

// CalcMonthConsecutiveDays 计算本月连续签到天数
func (s *Service) CalcMonthConsecutiveDays(ctx context.Context, userId uint64, year, month int) (int, error) {
	monthDays := getMonthDays(year, month)
	checkinBitmap, retroBitmap, err := s.getMonthBitmap(ctx, userId, year, month)
	if err != nil {
		g.Log().Errorf(ctx, "获取用户签到记录失败: %v", err)
		return 0, err
	}
	// 逻辑或
	bitmap := checkinBitmap | retroBitmap
	return calcMaxConsecutiveDays(bitmap, monthDays), nil
}

// calcMaxConsecutiveDays 计算最大连续签到天数
func calcMaxConsecutiveDays(bitmap uint64, monthDays int) int {
	// 逐位判断，计算出连续签到天数
	maxCount := 0
	currCount := 0
	for i := range monthDays {
		// 从右向左逐位判断
		checked := (bitmap>>i)&1 == 1
		if checked {
			currCount++
		} else {
			if currCount > maxCount {
				maxCount = currCount
			}
			currCount = 0
		}
	}
	// 循环结束再最后比较一次
	if currCount > maxCount {
		maxCount = currCount
	}
	return maxCount
}

// getFirstOfMonthOffset 获取当月第一天在一年中的偏移量
func getFirstOfMonthOffset(year, month int) int {
	// 1. 获取当月第一天
	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	// 2. 计算偏移量
	return firstOfMonth.YearDay() - 1 // offset 从 0 开始
}

// getMonhDays 获取当月天数
func getMonthDays(year, month int) int {
	// 1. 获取当月第一天
	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	// 2. 获取当月最后一天
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	return lastOfMonth.Day() // 返回当月天数
}

// AddPoints 增加积分
func AddPoints(ctx context.Context, input *model.PointsTransactionInput) error {
	// 用户积分汇总表 user_points 增加积分
	// 2.1 先查询(新用户可能没有记录)
	var userPoint entity.UserPoints
	if err := dao.UserPoints.Ctx(ctx).
		Where(dao.UserPoints.Columns().UserId, input.UserId).
		Scan(&userPoint); err != nil && !errors.Is(err, sql.ErrNoRows) {
		g.Log().Errorf(ctx, "查询用户积分汇总表失败: %v", err)
		return err
	}
	// 如果查询不到，则插入一条记录
	if userPoint.Id == 0 {
		userPoint = entity.UserPoints{UserId: input.UserId} // 创建新对象
	}
	userPoint.Points = userPoint.Points + input.Points           // 增加每日签到积分
	userPoint.PointsTotal = userPoint.PointsTotal + input.Points // 累计积分

	// 2.2 事务更新 用户积分汇总表 和 用户积分明细表
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 用户积分明细表 user_points_transactions 增加记录
		newRecord := entity.UserPointsTransactions{
			UserId:          input.UserId,
			PointsChange:    input.Points,
			CurrentBalance:  userPoint.Points,
			TransactionType: input.Type,
			Description:     input.Desc,
			CreatedAt:       gtime.NewFromTime(time.Now()),
			UpdatedAt:       gtime.NewFromTime(time.Now()),
		}
		if _, err := tx.Model(&entity.UserPointsTransactions{}).Insert(&newRecord); err != nil {
			g.Log().Errorf(ctx, "插入用户积分明细表失败: %v", err)
			return err
		}
		if _, err := tx.Model(&entity.UserPoints{}).
			Where(dao.UserPoints.Columns().UserId, input.UserId).
			Save(&userPoint); err != nil {
			g.Log().Errorf(ctx, "更新用户积分汇总表失败: %v", err)
			return err
		}
		return nil
	})
	return err
}
