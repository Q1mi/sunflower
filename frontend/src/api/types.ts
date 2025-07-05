/**
 * API 类型定义
 * @description 基于 OpenAPI 文档的完整类型定义，提供类型安全的API接口
 */

// ========== 通用响应类型 ==========

/** 通用API响应包装器 */
export interface ApiResponse<T = unknown> {
  /** 业务状态码 */
  code?: number
  /** 响应消息 */
  message?: string
  /** 响应数据 */
  data?: T
}

/** 分页查询参数 */
export interface PaginationParams {
  /** 每页数量，默认10 */
  limit?: number
  /** 偏移量，默认0 */
  offset?: number
}

/** 分页响应数据 */
export interface PaginationResponse<T> {
  /** 数据列表 */
  List: T[]
  /** 总数量 */
  total: number
  /** 是否有更多数据 */
  hasMore: boolean
}

// ========== 用户账户相关 ==========

/** 创建用户请求参数 */
export interface CreateUserRequest {
  /** 用户名，3-12字符 */
  username: string
  /** 密码，6-16字符 */
  password: string
  /** 确认密码 */
  confirmPassword: string
}

/** 创建用户响应 */
export type CreateUserResponse = Record<string, never>

/** 用户登录请求参数 */
export interface LoginRequest {
  /** 用户名，3-12字符 */
  username: string
  /** 密码，6-16字符 */
  password: string
}

/** 用户登录响应 */
export interface LoginResponse {
  /** 访问令牌 */
  accessToken: string
  /** 刷新令牌 */
  refreshToken: string
}

/** 刷新令牌请求参数 */
export interface RefreshTokenRequest {
  /** 刷新令牌 */
  refreshToken: string
}

/** 刷新令牌响应 */
export interface RefreshTokenResponse {
  /** 新的访问令牌 */
  accessToken: string
  /** 新的刷新令牌 */
  refreshToken: string
}

/** 用户信息响应 */
export interface UserProfileResponse {
  /** 用户名 */
  username: string
  /** 邮箱地址 */
  email: string
  /** 头像URL */
  avatar: string
}

// ========== 签到相关 ==========

/** 签到日历查询参数 */
export interface CheckinCalendarRequest {
  /** 年月格式，如 "2024-12" */
  yearMonth: string
}

/** 签到详细信息 */
export interface CheckinDetailInfo {
  /** 已签到日期数组 */
  checkedInDays: number[]
  /** 补签日期数组 */
  retroCheckedInDays: number[]
  /** 今日是否已签到 */
  isCheckedInToday: boolean
  /** 剩余补签次数 */
  remainRetroTimes: number
  /** 连续签到天数 */
  consecutiveDays: number
}

/** 签到日历响应 */
export interface CheckinCalendarResponse {
  /** 年份 */
  year: number
  /** 月份 */
  month: number
  /** 签到详细信息 */
  detail: CheckinDetailInfo
}

/** 每日签到请求参数 */
export type DailyCheckinRequest = Record<string, never>

/** 每日签到响应 */
export type DailyCheckinResponse = Record<string, never>

/** 补签请求参数 */
export interface RetroCheckinRequest {
  /** 补签日期，格式: "2024-12-15" */
  date: string
}

/** 补签响应 */
export type RetroCheckinResponse = Record<string, never>

// ========== 积分相关 ==========

/** 积分记录查询参数 */
export interface PointsRecordsRequest extends PaginationParams {
  /** 每页数量，默认10 */
  limit?: number
  /** 偏移量，默认0 */
  offset?: number
}

/** 积分记录项 */
export interface PointsRecord {
  /** 积分变动数量 */
  pointsChange: number
  /** 交易类型 */
  transactionType: string
  /** 描述信息 */
  description: string
  /** 交易时间 */
  transactionTime: string
}

/** 积分记录响应 */
export interface PointsRecordsResponse extends PaginationResponse<PointsRecord> {
  /** 积分记录列表 */
  List: PointsRecord[]
  /** 总数量 */
  total: number
  /** 是否有更多数据 */
  hasMore: boolean
}

/** 积分统计查询参数 */
export type PointsStatsRequest = Record<string, never>

/** 积分统计响应 */
export interface PointsStatsResponse {
  /** 总积分 */
  total: number
}

// ========== 前端组件使用的类型 ==========

/** 日历日期项 */
export interface CalendarDay {
  /** 日期 */
  day: number
  /** 是否为今天 */
  isToday: boolean
  /** 是否已签到 */
  isCheckedIn: boolean
  /** 是否已补签 */
  isRetroCheckedIn: boolean
  /** 是否为过去日期 */
  isPast: boolean
  /** 是否为奖励日 */
  isRewardDay: boolean
}

// ========== 兼容性类型别名 ==========

/** 签到日历详情（兼容性别名） */
export type CheckInCalendarDetail = CheckinDetailInfo

/** 签到日历响应（兼容性别名） */
export type CheckInCalendarResponse = CheckinCalendarResponse

// ========== 枚举类型 ==========

/** 交易类型枚举 */
export enum TransactionType {
  /** 签到 */
  CHECK_IN = 'check-in',
  /** 补签 */
  RETRO_CHECK_IN = 'retro-check-in',
  /** 奖励 */
  REWARD = 'reward',
  /** 注册 */
  REGISTER = 'register',
}

/** 用户状态枚举 */
export enum UserStatus {
  /** 活跃 */
  ACTIVE = 'active',
  /** 禁用 */
  DISABLED = 'disabled',
  /** 待验证 */
  PENDING = 'pending',
}

// ========== 工具类型 ==========

/** 可选字段类型 */
export type Optional<T, K extends keyof T> = Omit<T, K> & Partial<Pick<T, K>>

/** 必需字段类型 */
export type RequiredFields<T, K extends keyof T> = T & Required<Pick<T, K>>

/** API错误响应 */
export interface ApiErrorResponse {
  /** 错误码 */
  code: number
  /** 错误消息 */
  message: string
  /** 错误详情 */
  details?: Record<string, unknown>
  /** 时间戳 */
  timestamp: number
}
