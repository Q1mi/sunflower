/**
 * Token 管理工具
 * @description 提供统一的Token管理功能，支持自动刷新和存储管理
 */

import { STORAGE_KEYS } from '../constants'

/** Token信息接口 */
export interface TokenInfo {
  /** 访问令牌 */
  accessToken: string
  /** 刷新令牌 */
  refreshToken: string
  /** 是否记住登录状态 */
  remember: boolean
  /** 令牌创建时间 */
  createTime: number
}

/** Token管理器类 */
class TokenManager {
  /** 正在刷新token的Promise，避免并发刷新 */
  private refreshPromise: Promise<string> | null = null

  /**
   * 获取访问令牌
   * @returns 访问令牌或null
   */
  getAccessToken(): string | null {
    return (
      localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN) ||
      sessionStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN)
    )
  }

  /**
   * 获取刷新令牌
   * @returns 刷新令牌或null
   */
  getRefreshToken(): string | null {
    return (
      localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN) ||
      sessionStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN)
    )
  }

  /**
   * 获取完整Token信息
   * @returns Token信息或null
   */
  getTokenInfo(): TokenInfo | null {
    const accessToken = this.getAccessToken()
    const refreshToken = this.getRefreshToken()

    if (!accessToken || !refreshToken) {
      return null
    }

    // 判断是否记住登录状态
    const remember = !!(
      localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN) ||
      localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN)
    )

    return {
      accessToken,
      refreshToken,
      remember,
      createTime: Date.now(),
    }
  }

  /**
   * 设置Token
   * @param accessToken 访问令牌
   * @param refreshToken 刷新令牌
   * @param remember 是否记住登录状态
   */
  setTokens(accessToken: string, refreshToken: string, remember: boolean = false): void {
    // 清除所有存储中的token
    this.clearTokens()

    const storage = remember ? localStorage : sessionStorage

    storage.setItem(STORAGE_KEYS.ACCESS_TOKEN, accessToken)
    storage.setItem(STORAGE_KEYS.REFRESH_TOKEN, refreshToken)

    console.log(`Token已保存到${remember ? 'localStorage' : 'sessionStorage'}`)
  }

  /**
   * 清除所有Token
   */
  clearTokens(): void {
    // 清除新版本token
    localStorage.removeItem(STORAGE_KEYS.ACCESS_TOKEN)
    localStorage.removeItem(STORAGE_KEYS.REFRESH_TOKEN)
    sessionStorage.removeItem(STORAGE_KEYS.ACCESS_TOKEN)
    sessionStorage.removeItem(STORAGE_KEYS.REFRESH_TOKEN)

    // 清除旧版本token（兼容性）
    localStorage.removeItem(STORAGE_KEYS.LEGACY_TOKEN)
    sessionStorage.removeItem(STORAGE_KEYS.LEGACY_TOKEN)
    localStorage.removeItem(STORAGE_KEYS.USER_INFO)
    sessionStorage.removeItem(STORAGE_KEYS.USER_INFO)

    console.log('所有Token已清除')
  }

  /**
   * 检查Token是否存在
   * @returns 是否存在有效Token
   */
  hasValidTokens(): boolean {
    return !!(this.getAccessToken() && this.getRefreshToken())
  }

  /**
   * 检查是否记住登录状态
   * @returns 是否记住登录状态
   */
  isRememberLogin(): boolean {
    return !!(
      localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN) ||
      localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN)
    )
  }

  /**
   * 设置刷新Promise（避免并发刷新）
   * @param promise 刷新Promise
   */
  setRefreshPromise(promise: Promise<string> | null): void {
    this.refreshPromise = promise
  }

  /**
   * 获取刷新Promise
   * @returns 刷新Promise
   */
  getRefreshPromise(): Promise<string> | null {
    return this.refreshPromise
  }

  /**
   * 获取Token统计信息
   * @returns Token统计信息
   */
  getTokenStats(): {
    hasAccessToken: boolean
    hasRefreshToken: boolean
    isRemember: boolean
    storageType: 'localStorage' | 'sessionStorage' | 'none'
  } {
    const hasAccessToken = !!this.getAccessToken()
    const hasRefreshToken = !!this.getRefreshToken()
    const isRemember = this.isRememberLogin()

    let storageType: 'localStorage' | 'sessionStorage' | 'none' = 'none'
    if (hasAccessToken || hasRefreshToken) {
      storageType = isRemember ? 'localStorage' : 'sessionStorage'
    }

    return {
      hasAccessToken,
      hasRefreshToken,
      isRemember,
      storageType,
    }
  }
}

/** 全局Token管理器实例 */
export const tokenManager = new TokenManager()

/** 导出类型 */
export type { TokenManager }
