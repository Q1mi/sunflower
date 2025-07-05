/**
 * API 统一导出文件
 * @description 基于后端 OpenAPI 文档的完整API接口导出
 * @version 1.0.0
 */

// ========== 类型定义导出 ==========
export * from './types'

// ========== 核心API接口导出 ==========

/** 用户账户相关接口 */
export { createUser, login, refreshToken, getUserProfile } from './account'

/** 签到相关接口 */
export { getCheckinCalendar, dailyCheckin, retroCheckin } from './checkin'

/** 积分相关接口 */
export { getPointsRecords, getPointsStats } from './points'

// ========== 工具类导出 ==========

/** HTTP客户端（用于自定义请求） */
export { default as http } from './http'

/** Token管理器 */
export { tokenManager } from './utils/token-manager'

/** 缓存管理器 */
export { cacheManager, CacheKeys } from './utils/cache'

/** 错误处理器 */
export { errorHandler } from './utils/error-handler'

// ========== 常量配置导出 ==========
export {
  API_CONFIG,
  API_ENDPOINTS,
  STORAGE_KEYS,
  HTTP_STATUS,
  BUSINESS_ERROR_CODES,
  ERROR_MESSAGES,
} from './constants'

// ========== 适配器接口导出 ==========

/** 前端组件适配接口 */
export {
  getPointsInfo,
  getCheckinCalendarDetail,
  checkIn,
  retroCheckIn,
  getPointRecords,
  clearCache,
} from './adapter'

/** 用户管理适配接口 */
export {
  userLogin,
  userRegister,
  getCurrentUser,
  userLogout,
  isLoggedIn,
  getUserTokenStats,
} from './adapter'
