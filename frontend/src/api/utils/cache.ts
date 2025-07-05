/**
 * 缓存管理工具
 * @description 提供统一的缓存管理功能，支持过期时间和自动清理
 */

import { CACHE_CONFIG } from '../constants'

/** 缓存项接口 */
interface CacheItem<T = unknown> {
  /** 缓存数据 */
  data: T
  /** 过期时间戳 */
  expireTime: number
  /** 创建时间戳 */
  createTime: number
}

/** 缓存管理器类 */
class CacheManager {
  /** 内存缓存存储 */
  private cache = new Map<string, CacheItem>()

  /**
   * 设置缓存
   * @param key 缓存键
   * @param data 缓存数据
   * @param duration 缓存持续时间（毫秒），默认使用配置值
   */
  set<T>(key: string, data: T, duration: number = CACHE_CONFIG.DURATION): void {
    const now = Date.now()
    const cacheItem: CacheItem<T> = {
      data,
      expireTime: now + duration,
      createTime: now,
    }
    this.cache.set(key, cacheItem)
  }

  /**
   * 获取缓存
   * @param key 缓存键
   * @returns 缓存数据，如果不存在或已过期则返回null
   */
  get<T>(key: string): T | null {
    const cacheItem = this.cache.get(key)
    if (!cacheItem) {
      return null
    }

    const now = Date.now()
    if (now > cacheItem.expireTime) {
      this.cache.delete(key)
      return null
    }

    return cacheItem.data as T
  }

  /**
   * 检查缓存是否存在且有效
   * @param key 缓存键
   * @returns 是否存在有效缓存
   */
  has(key: string): boolean {
    return this.get(key) !== null
  }

  /**
   * 删除指定缓存
   * @param key 缓存键
   */
  delete(key: string): void {
    this.cache.delete(key)
  }

  /**
   * 清空所有缓存
   */
  clear(): void {
    this.cache.clear()
  }

  /**
   * 清理过期缓存
   */
  cleanup(): void {
    const now = Date.now()
    for (const [key, item] of this.cache.entries()) {
      if (now > item.expireTime) {
        this.cache.delete(key)
      }
    }
  }

  /**
   * 获取缓存统计信息
   * @returns 缓存统计
   */
  getStats(): {
    total: number
    expired: number
    valid: number
  } {
    const now = Date.now()
    let expired = 0
    let valid = 0

    for (const item of this.cache.values()) {
      if (now > item.expireTime) {
        expired++
      } else {
        valid++
      }
    }

    return {
      total: this.cache.size,
      expired,
      valid,
    }
  }
}

/** 全局缓存管理器实例 */
export const cacheManager = new CacheManager()

/** 缓存键生成器 */
export const CacheKeys = {
  /** 积分信息缓存键 */
  pointsInfo: () => CACHE_CONFIG.POINTS_INFO_KEY,

  /** 日历缓存键 */
  calendar: (year: number, month: number) => `${CACHE_CONFIG.CALENDAR_KEY_PREFIX}_${year}_${month}`,

  /** 用户信息缓存键 */
  userProfile: (userId: string | number) => `user_profile_${userId}`,

  /** 积分记录缓存键 */
  pointsRecords: (limit: number, offset: number) => `points_records_${limit}_${offset}`,
} as const

// 定期清理过期缓存（每5分钟执行一次）
if (typeof window !== 'undefined') {
  setInterval(
    () => {
      cacheManager.cleanup()
    },
    5 * 60 * 1000,
  )
}
