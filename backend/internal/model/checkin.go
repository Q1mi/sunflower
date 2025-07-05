package model

// 签到相关参数

type MonthDetailInput struct {
	UserId uint64
	Year   int // 年份
	Month  int // 月份
}

type MonthDetailOutput struct {
	CheckedInDays      []int `json:"checkedInDays"`      // 签到的日期
	RetroCheckedInDays []int `json:"retroCheckedInDays"` // 补签的日期
	IsCheckedInToday   bool  `json:"isCheckedIn"`        // 当天是否签到
	RemainRetroTimes   int   `json:"remainRetroTimes"`   // 剩余补签次数
	ConsecutiveDays    int   `json:"consecutiveDays"`    // 连续签到天数
}
