package v1

import "github.com/gogf/gf/v2/frame/g"

// 定义签到相关的 API

type DailyReq struct {
	g.Meta `path:"/checkins" method:"POST" tags:"签到" summary:"每日签到"`
}

type DailyRes struct{}

type CalendarReq struct {
	g.Meta    `path:"/checkins/calendar" method:"GET" tags:"签到" summary:"获取签到日历"`
	YearMonth string `p:"yearMonth" v:"required#请选择年月,例如2025-02"`
}

type CalendarRes struct {
	Year   int        `json:"year"`
	Month  int        `json:"month"`
	Detail DetailInfo `json:"detail"`
}

type DetailInfo struct {
	CheckedInDays      []int `json:"checkedInDays"`      // 签到的日期
	RetroCheckedInDays []int `json:"retroCheckedInDays"` // 补签的日期
	IsCheckedInToday   bool  `json:"isCheckedInToday"`   // 当天是否签到
	RemainRetroTimes   int   `json:"remainRetroTimes"`   // 剩余补签次数
	ConsecutiveDays    int   `json:"consecutiveDays"`    // 连续签到天数
}

// RetroReq 补签请求结构体
type RetroReq struct {
	g.Meta `path:"/checkins/retroactive" method:"POST" tags:"签到" summary:"补签"`
	Date   string `p:"date" v:"required#请选择补签日期"`
}

type RetroRes struct{}
