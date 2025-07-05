// 定义接口返回的数据类型

// 导出用户相关类型
export * from './user'

// 积分信息接口返回类型
export interface PointsInfo {
  totalPoints: number // 总积分
  consecutiveDays: number // 当月连续签到天数
  retroAvailable: number // 剩余补签次数
  checkedInToday?: boolean // 今日是否已签到
  retroCheckedInDays?: number[] // 补签日期列表
}

// 签到接口返回类型
export interface CheckInResult {
  success: boolean // 签到是否成功
  points: number // 获得的积分
  message: string // 提示信息
}

// 补签接口返回类型
export interface RetroCheckInResult {
  success: boolean // 补签是否成功
  points: number // 获得的积分
  retroAvailable: number // 剩余补签次数
  message: string // 提示信息
}

// 积分记录项类型
export interface PointRecord {
  id: number // 记录ID
  type: string // 记录类型
  points: number // 积分变动
  date: string // 日期时间
  description: string // 描述
}

// 积分记录列表接口返回类型
export interface PointRecordList {
  records: PointRecord[] // 积分记录列表
  hasMore: boolean // 是否有更多记录
}

// 日历日期类型
export interface CalendarDay {
  day: number // 日期
  isToday: boolean // 是否是今天
  isCheckedIn: boolean // 是否已签到
  isRetroCheckedIn: boolean // 是否是补签
  isPast: boolean // 是否是过去的日期
  isRewardDay: boolean // 是否是奖励日
}

// 签到日历详情类型
export interface CheckInCalendarDetail {
  checkedInDays: number[] // 已签到的日期数组
  retroCheckedInDays: number[] // 补签的日期数组
  isCheckedInToday: boolean // 今日是否已签到
  remainRetroTimes: number // 剩余补签次数
  consecutiveDays: number // 连续签到天数
}

// 签到日历接口返回类型
export interface CheckInCalendarResponse {
  year: number // 年份
  month: number // 月份
  detail: CheckInCalendarDetail // 详情数据
}
