/**
 * API 常量配置
 * @description 统一管理API相关的常量，便于维护和配置
 */

/** API基础配置 */
export const API_CONFIG = {
  /** 基础URL */
  BASE_URL: 'http://127.0.0.1:8000/api/v1',
  /** 请求超时时间 */
  TIMEOUT: 10000,
  /** 默认请求头 */
  DEFAULT_HEADERS: {
    'Content-Type': 'application/json',
  },
} as const

/** 缓存配置 */
export const CACHE_CONFIG = {
  /** 缓存持续时间（毫秒） */
  DURATION: 3000,
  /** 积分信息缓存键 */
  POINTS_INFO_KEY: 'points_info',
  /** 日历缓存键前缀 */
  CALENDAR_KEY_PREFIX: 'calendar',
} as const

/** 存储键名 */
export const STORAGE_KEYS = {
  /** 访问令牌 */
  ACCESS_TOKEN: 'accessToken',
  /** 刷新令牌 */
  REFRESH_TOKEN: 'refreshToken',
  /** 用户信息 */
  USER_INFO: 'user',
  /** 旧版本兼容 */
  LEGACY_TOKEN: 'token',
} as const

/** API端点 */
export const API_ENDPOINTS = {
  /** 用户相关 */
  USER: {
    CREATE: '/users',
    LOGIN: '/auth/login',
    REFRESH: '/auth/refresh',
    PROFILE: '/users/me',
  },
  /** 签到相关 */
  CHECKIN: {
    CALENDAR: '/checkins/calendar',
    DAILY: '/checkins',
    RETRO: '/checkins/retroactive',
  },
  /** 积分相关 */
  POINTS: {
    RECORDS: '/points/records',
    STATS: '/points/summary',
  },
} as const

/** HTTP状态码 */
export const HTTP_STATUS = {
  OK: 200,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  INTERNAL_SERVER_ERROR: 500,
} as const

/** 业务错误码 */
export const BUSINESS_ERROR_CODES = {
  SUCCESS: 200,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  VALIDATION_ERROR: 422,
  INTERNAL_ERROR: 500,
} as const

/** 错误消息 */
export const ERROR_MESSAGES = {
  NETWORK_ERROR: '网络连接失败，请检查网络设置',
  TOKEN_EXPIRED: '登录已过期，请重新登录',
  PERMISSION_DENIED: '权限不足',
  SERVER_ERROR: '服务器内部错误',
  VALIDATION_ERROR: '参数验证失败',
  UNKNOWN_ERROR: '未知错误',
} as const
