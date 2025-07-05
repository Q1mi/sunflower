/**
 * 签到相关API接口
 * @description 提供签到日历、每日签到和补签功能的API接口
 */

import http from './http'
import { API_ENDPOINTS } from './constants'
import type { CheckinCalendarResponse, DailyCheckinResponse, RetroCheckinResponse } from './types'

/**
 * 获取签到日历详情
 * @description 获取指定年月的签到日历信息，包括已签到日期、补签日期等
 * @param yearMonth 年月字符串，格式: "2024-12"
 * @returns 签到日历详情
 * @example
 * ```typescript
 * const calendar = await getCheckinCalendar('2024-12')
 * console.log(calendar.detail.checkedInDays) // 已签到日期
 * ```
 */
export const getCheckinCalendar = async (yearMonth: string): Promise<CheckinCalendarResponse> => {
  const response = await http.get<CheckinCalendarResponse>(
    `${API_ENDPOINTS.CHECKIN.CALENDAR}?yearMonth=${yearMonth}`,
  )
  return response.data
}

/**
 * 每日签到
 * @description 执行当日签到操作，获取签到奖励
 * @returns 签到结果
 * @example
 * ```typescript
 * await dailyCheckin()
 * console.log('签到成功！')
 * ```
 */
export const dailyCheckin = async (): Promise<DailyCheckinResponse> => {
  const response = await http.post<DailyCheckinResponse>(API_ENDPOINTS.CHECKIN.DAILY, {})
  return response.data
}

/**
 * 补签操作
 * @description 对指定日期进行补签，消耗补签次数
 * @param date 补签日期，格式: "2024-12-15"
 * @returns 补签结果
 * @example
 * ```typescript
 * await retroCheckin('2024-12-15')
 * console.log('补签成功！')
 * ```
 */
export const retroCheckin = async (date: string): Promise<RetroCheckinResponse> => {
  const response = await http.post<RetroCheckinResponse>(API_ENDPOINTS.CHECKIN.RETRO, { date })
  return response.data
}
