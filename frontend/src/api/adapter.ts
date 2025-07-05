/**
 * API 适配层
 * @description 将后端API接口适配为前端组件需要的格式，提供缓存、降级处理和用户管理功能
 */

import {
  getCheckinCalendar,
  dailyCheckin,
  retroCheckin,
  getPointsStats,
  getPointsRecords,
  login as apiLogin,
  createUser,
  getUserProfile,
} from './index'
import { cacheManager, CacheKeys } from './utils/cache'
import { errorHandler } from './utils/error-handler'
import { tokenManager } from './utils/token-manager'
import { STORAGE_KEYS } from './constants'
import type { CheckInCalendarResponse, PointsInfo, CheckInResult, PointRecord } from '@/types'
import type { LoginRequest, RegisterRequest, User } from '@/types'
import type { CreateUserRequest } from './types'

// ========== 签到和积分相关适配器 ==========

/**
 * 获取积分信息（适配组件接口）
 * @description 获取用户积分信息，支持缓存机制
 * @returns 积分信息
 */
export async function getPointsInfo(): Promise<PointsInfo> {
  const cacheKey = CacheKeys.pointsInfo()

  // 尝试从缓存获取
  const cachedData = cacheManager.get<PointsInfo>(cacheKey)
  if (cachedData) {
    console.log('[Cache Hit] 返回缓存的积分数据')
    return cachedData
  }

  try {
    const stats = await getPointsStats()
    console.log('[API Success] 获取积分统计成功:', stats)

    // 构建基础积分信息
    const baseInfo: PointsInfo = {
      totalPoints: stats.total || 0,
      consecutiveDays: 0,
      retroAvailable: 3,
      checkedInToday: false,
      retroCheckedInDays: [],
    }

    // 更新缓存
    cacheManager.set(cacheKey, baseInfo)
    console.log('[Cache Set] 积分信息已缓存')

    return baseInfo
  } catch (error) {
    const standardError = errorHandler.handleUnknownError(error, {
      function: 'getPointsInfo',
      cacheKey,
    })

    console.error('[API Error] 获取积分信息失败:', standardError.message)

    // 返回默认值
    const defaultInfo: PointsInfo = {
      totalPoints: 0,
      consecutiveDays: 0,
      retroAvailable: 3,
      checkedInToday: false,
      retroCheckedInDays: [],
    }

    return defaultInfo
  }
}

/**
 * 获取签到日历详情（适配组件接口）
 * @description 获取指定年月的签到日历信息，支持缓存机制
 * @param year 年份
 * @param month 月份
 * @returns 签到日历详情
 */
export async function getCheckinCalendarDetail(
  year: number,
  month: number,
): Promise<CheckInCalendarResponse> {
  const cacheKey = CacheKeys.calendar(year, month)

  // 检查缓存
  const cachedData = cacheManager.get<CheckInCalendarResponse>(cacheKey)
  if (cachedData) {
    console.log('[Cache Hit] 返回缓存的日历数据')
    return cachedData
  }

  try {
    const yearMonth = `${year}-${String(month).padStart(2, '0')}`
    const response = await getCheckinCalendar(yearMonth)

    console.log('[API Success] 获取签到日历成功:', response)

    // 更新缓存
    cacheManager.set(cacheKey, response)
    console.log('[Cache Set] 日历数据已缓存')

    // 同时更新积分信息缓存中的相关字段
    const pointsInfoKey = CacheKeys.pointsInfo()
    const pointsInfo = cacheManager.get<PointsInfo>(pointsInfoKey)
    if (pointsInfo) {
      const updatedPointsInfo: PointsInfo = {
        ...pointsInfo,
        consecutiveDays: response.detail.consecutiveDays,
        retroAvailable: response.detail.remainRetroTimes,
        checkedInToday: response.detail.isCheckedInToday,
        retroCheckedInDays: response.detail.retroCheckedInDays,
      }
      cacheManager.set(pointsInfoKey, updatedPointsInfo)
      console.log('[Cache Update] 积分信息缓存已更新')
    }

    return response
  } catch (error) {
    const standardError = errorHandler.handleUnknownError(error, {
      function: 'getCheckinCalendarDetail',
      year,
      month,
      cacheKey,
    })

    console.error('[API Error] 获取签到日历详情失败:', standardError.message)

    // 返回空数据
    const emptyResponse: CheckInCalendarResponse = {
      year,
      month,
      detail: {
        checkedInDays: [],
        retroCheckedInDays: [],
        isCheckedInToday: false,
        remainRetroTimes: 3,
        consecutiveDays: 0,
      },
    }

    return emptyResponse
  }
}

/**
 * 清除缓存
 * @description 清除所有相关缓存，在签到、补签等操作后调用
 */
export function clearCache(): void {
  cacheManager.clear()
  console.log('[Cache Clear] 所有缓存已清除')
}

/**
 * 签到（适配组件接口）
 * @description 执行每日签到操作
 * @returns 签到结果
 */
export async function checkIn(): Promise<CheckInResult> {
  try {
    await dailyCheckin()
    console.log('[API Success] 签到成功')

    // 清除缓存，确保下次获取最新数据
    clearCache()

    return {
      success: true,
      points: 1, // 默认签到积分
      message: '签到成功！',
    }
  } catch (error) {
    const standardError = errorHandler.handleUnknownError(error, {
      function: 'checkIn',
    })

    console.error('[API Error] 签到失败:', standardError.message)

    return {
      success: false,
      points: 0,
      message: standardError.userMessage,
    }
  }
}

/**
 * 补签（适配组件接口）
 * @description 对指定日期进行补签操作
 * @param date 补签日期
 * @returns 补签结果
 */
export async function retroCheckIn(date: string): Promise<CheckInResult> {
  try {
    await retroCheckin(date)
    console.log('[API Success] 补签成功')

    // 清除缓存，确保下次获取最新数据
    clearCache()

    return {
      success: true,
      points: 5, // 默认补签积分
      message: '补签成功！',
    }
  } catch (error) {
    const standardError = errorHandler.handleUnknownError(error, {
      function: 'retroCheckIn',
      date,
    })

    console.error('[API Error] 补签失败:', standardError.message)

    return {
      success: false,
      points: 0,
      message: standardError.userMessage,
    }
  }
}

/**
 * 获取积分记录（适配组件接口）
 * @description 获取用户的积分变动记录
 * @returns 积分记录和分页信息
 */
export async function getPointRecords(): Promise<{ records: PointRecord[]; hasMore: boolean }> {
  try {
    const data = await getPointsRecords({ limit: 20, offset: 0 })
    console.log('[API Success] 获取积分记录成功')
    const records: PointRecord[] = data.list.map((record, index) => ({
      id: index + 1, // 临时ID
      type: record.transactionType as 'check-in' | 'retro-check-in' | 'reward' | 'register',
      points: record.pointsChange,
      description: record.description,
      date: record.transactionTime,
    }))

    return {
      records,
      hasMore: data.hasMore,
    }
  } catch (error) {
    const standardError = errorHandler.handleUnknownError(error, {
      function: 'getPointRecords',
    })

    console.error('[API Error] 获取积分记录失败:', standardError.message)

    return {
      records: [],
      hasMore: false,
    }
  }
}

// ========== 用户管理相关适配器 ==========

/**
 * 用户登录（适配store接口）
 * @description 用户登录并处理Token存储和用户信息获取
 * @param loginData 登录数据
 * @param remember 是否记住登录状态
 * @returns 登录结果，包含Token和用户信息
 */
export async function userLogin(loginData: LoginRequest, remember: boolean = false) {
  try {
    console.log('[User Login] 开始登录流程')
    const loginResult = await apiLogin(loginData)
    console.log('[User Login] API登录成功')

    // 存储新的token结构
    tokenManager.setTokens(loginResult.accessToken, loginResult.refreshToken, remember)

    // 获取用户信息
    try {
      const userProfile = await getUserProfile()
      const user: User = {
        id: 1, // 临时ID，实际应该从API返回
        username: userProfile.username,
        email: userProfile.email,
        avatar: userProfile.avatar,
        createdAt: new Date().toISOString(),
      }

      // 存储用户信息
      const storage = remember ? localStorage : sessionStorage
      storage.setItem(STORAGE_KEYS.USER_INFO, JSON.stringify(user))

      console.log('[User Login] 用户信息获取成功')

      return {
        accessToken: loginResult.accessToken,
        refreshToken: loginResult.refreshToken,
        user,
      }
    } catch (profileError) {
      console.warn('[User Login] 获取用户信息失败，使用基础信息', profileError)

      // 如果获取用户信息失败，返回基础用户信息
      const user: User = {
        id: 1,
        username: loginData.username,
        email: '',
        avatar: '',
        createdAt: new Date().toISOString(),
      }

      const storage = remember ? localStorage : sessionStorage
      storage.setItem(STORAGE_KEYS.USER_INFO, JSON.stringify(user))

      return {
        accessToken: loginResult.accessToken,
        refreshToken: loginResult.refreshToken,
        user,
      }
    }
  } catch (error) {
    const standardError = errorHandler.handleUnknownError(error, {
      function: 'userLogin',
      username: loginData.username,
      remember,
    })

    console.error('[User Login] 登录失败:', standardError.message)
    throw new Error(standardError.userMessage)
  }
}

/**
 * 用户注册（适配store接口）
 * @description 用户注册并自动登录
 * @param registerData 注册数据
 * @returns 注册并登录的结果
 */
export async function userRegister(registerData: RegisterRequest) {
  try {
    console.log('[User Register] 开始注册流程')

    const createUserData: CreateUserRequest = {
      username: registerData.username,
      password: registerData.password,
      confirmPassword: registerData.confirmPassword,
    }

    await createUser(createUserData)
    console.log('[User Register] 注册成功，开始自动登录')

    // 注册成功后自动登录
    const loginData: LoginRequest = {
      username: registerData.username,
      password: registerData.password,
    }

    return await userLogin(loginData)
  } catch (error) {
    const standardError = errorHandler.handleUnknownError(error, {
      function: 'userRegister',
      username: registerData.username,
    })

    console.error('[User Register] 注册失败:', standardError.message)
    throw new Error(standardError.userMessage)
  }
}

/**
 * 获取当前用户信息（适配store接口）
 * @description 获取当前登录用户的完整信息
 * @returns 用户信息或null
 */
export async function getCurrentUser(): Promise<User | null> {
  try {
    const accessToken = tokenManager.getAccessToken()
    if (!accessToken) {
      console.log('[Get Current User] 未找到访问令牌')
      return null
    }

    // 先尝试从本地存储获取用户信息
    let storedUser = localStorage.getItem(STORAGE_KEYS.USER_INFO)
    if (!storedUser) {
      storedUser = sessionStorage.getItem(STORAGE_KEYS.USER_INFO)
    }

    if (storedUser) {
      try {
        const user = JSON.parse(storedUser) as User
        console.log('[Get Current User] 从本地存储获取用户信息成功')

        // 验证本地用户信息的有效性 - 尝试调用API验证token
        try {
          console.log('[Get Current User] 验证本地用户信息有效性')
          await getUserProfile()
          console.log('[Get Current User] 本地用户信息验证成功')
          return user
        } catch (validationError) {
          console.warn('[Get Current User] 本地用户信息验证失败，token可能已过期', validationError)

          // 清除无效的本地用户信息和token
          localStorage.removeItem(STORAGE_KEYS.USER_INFO)
          sessionStorage.removeItem(STORAGE_KEYS.USER_INFO)
          tokenManager.clearTokens()
          clearCache()

          console.log('[Get Current User] 已清除无效的本地用户信息，需要重新登录')
          return null
        }
      } catch (parseError) {
        console.warn('[Get Current User] 本地用户信息解析失败，从API获取', parseError)
        // 清除损坏的本地数据
        localStorage.removeItem(STORAGE_KEYS.USER_INFO)
        sessionStorage.removeItem(STORAGE_KEYS.USER_INFO)
      }
    }

    // 从API获取用户信息（首次获取或本地信息无效时）
    console.log('[Get Current User] 从API获取用户信息')
    const userProfile = await getUserProfile()
    const user: User = {
      id: 1, // 临时ID
      username: userProfile.username,
      email: userProfile.email,
      avatar: userProfile.avatar,
      createdAt: new Date().toISOString(),
    }

    // 更新本地存储（根据token存储位置决定）
    const storage = tokenManager.isRememberLogin() ? localStorage : sessionStorage
    storage.setItem(STORAGE_KEYS.USER_INFO, JSON.stringify(user))

    console.log('[Get Current User] 用户信息获取并缓存成功')
    return user
  } catch (error) {
    const standardError = errorHandler.handleUnknownError(error, {
      function: 'getCurrentUser',
    })

    console.error('[Get Current User] 获取用户信息失败:', standardError.message)

    // 清除无效的token和用户信息
    tokenManager.clearTokens()
    localStorage.removeItem(STORAGE_KEYS.USER_INFO)
    sessionStorage.removeItem(STORAGE_KEYS.USER_INFO)
    clearCache()

    console.log('[Get Current User] 已清除所有用户相关数据，需要重新登录')
    return null
  }
}

/**
 * 用户登出（适配store接口）
 * @description 清除用户登录状态和相关数据
 */
export function userLogout(): void {
  console.log('[User Logout] 开始登出流程')

  // 清除Token
  tokenManager.clearTokens()

  // 清除用户信息
  localStorage.removeItem(STORAGE_KEYS.USER_INFO)
  sessionStorage.removeItem(STORAGE_KEYS.USER_INFO)

  // 清除缓存
  clearCache()

  console.log('[User Logout] 登出完成')
}

/**
 * 检查用户登录状态
 * @description 检查用户是否已登录
 * @returns 是否已登录
 */
export function isLoggedIn(): boolean {
  const hasTokens = tokenManager.hasValidTokens()
  console.log('[Check Login Status] 登录状态:', hasTokens)
  return hasTokens
}

/**
 * 获取用户Token统计信息
 * @description 获取当前用户Token的详细状态信息
 * @returns Token统计信息
 */
export function getUserTokenStats() {
  const stats = tokenManager.getTokenStats()
  console.log('[User Token Stats]', stats)
  return stats
}
