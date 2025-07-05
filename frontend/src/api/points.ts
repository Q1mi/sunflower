/**
 * 积分相关API接口
 * @description 提供积分记录查询和积分统计功能的API接口
 */

import http from './http'
import { API_ENDPOINTS } from './constants'
import type { PointsRecordsRequest, PointsRecordsResponse, PointsStatsResponse } from './types'

/**
 * 获取积分记录
 * @description 分页获取用户的积分变动记录
 * @param params 查询参数，包含分页信息
 * @returns 积分记录列表
 * @example
 * ```typescript
 * const records = await getPointsRecords({ limit: 20, offset: 0 })
 * console.log(records.List) // 积分记录列表
 * console.log(records.hasMore) // 是否有更多数据
 * ```
 */
export const getPointsRecords = async (
  params?: PointsRecordsRequest,
): Promise<PointsRecordsResponse> => {
  const queryParams = new URLSearchParams()
  if (params?.limit) queryParams.append('limit', params.limit.toString())
  if (params?.offset) queryParams.append('offset', params.offset.toString())

  const url = queryParams.toString()
    ? `${API_ENDPOINTS.POINTS.RECORDS}?${queryParams.toString()}`
    : API_ENDPOINTS.POINTS.RECORDS

  const response = await http.get<PointsRecordsResponse>(url)
  return response.data
}

/**
 * 获取积分统计
 * @description 获取用户的积分统计信息，包括总积分等
 * @returns 积分统计信息
 * @example
 * ```typescript
 * const stats = await getPointsStats()
 * console.log(stats.total) // 总积分
 * ```
 */
export const getPointsStats = async (): Promise<PointsStatsResponse> => {
  const response = await http.get<PointsStatsResponse>(API_ENDPOINTS.POINTS.STATS)
  return response.data
}
