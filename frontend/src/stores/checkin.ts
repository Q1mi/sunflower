// 签到数据状态管理
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getPointsInfo, getCheckinCalendarDetail, clearCache } from '@/api/adapter'
import type { PointsInfo, CheckInCalendarDetail } from '@/types'

export const useCheckinStore = defineStore('checkin', () => {
  // 状态
  const pointsInfo = ref<PointsInfo>({
    totalPoints: 0,
    consecutiveDays: 0,
    retroAvailable: 0,
    checkedInToday: false,
    retroCheckedInDays: [],
  })

  const calendarDetail = ref<CheckInCalendarDetail>({
    checkedInDays: [],
    retroCheckedInDays: [],
    isCheckedInToday: false,
    remainRetroTimes: 0,
    consecutiveDays: 0,
  })

  const loading = ref(false)
  const error = ref('')

  // 获取积分信息
  async function fetchPointsInfo() {
    if (loading.value) return // 防止重复请求

    loading.value = true
    error.value = ''
    try {
      const data = await getPointsInfo()
      pointsInfo.value = data
      console.log('积分信息获取成功:', data)
    } catch (err: unknown) {
      console.error('获取积分信息失败', err)
      error.value = err instanceof Error ? err.message : '获取积分信息失败'
      // 不要抛出错误，让组件继续正常工作
    } finally {
      loading.value = false
    }
  }

  // 获取签到日历详情
  async function fetchCalendarDetail(year: number, month: number) {
    if (loading.value) return // 防止重复请求

    loading.value = true
    error.value = ''
    try {
      const response = await getCheckinCalendarDetail(year, month)
      calendarDetail.value = response.detail

      // 同步更新积分信息中的相关字段
      pointsInfo.value.consecutiveDays = response.detail.consecutiveDays
      pointsInfo.value.retroAvailable = response.detail.remainRetroTimes
      pointsInfo.value.checkedInToday = response.detail.isCheckedInToday
      pointsInfo.value.retroCheckedInDays = response.detail.retroCheckedInDays
    } catch (err: unknown) {
      console.error('获取签到日历详情失败', err)
      error.value = err instanceof Error ? err.message : '获取签到日历详情失败'
    } finally {
      loading.value = false
    }
  }

  // 刷新所有数据
  async function refreshData(year?: number, month?: number) {
    clearCache() // 清除缓存
    await fetchPointsInfo()
    if (year !== undefined && month !== undefined) {
      await fetchCalendarDetail(year, month)
    }
  }

  // 重置状态
  function resetState() {
    pointsInfo.value = {
      totalPoints: 0,
      consecutiveDays: 0,
      retroAvailable: 0,
      checkedInToday: false,
      retroCheckedInDays: [],
    }
    calendarDetail.value = {
      checkedInDays: [],
      retroCheckedInDays: [],
      isCheckedInToday: false,
      remainRetroTimes: 0,
      consecutiveDays: 0,
    }
    error.value = ''
  }

  return {
    pointsInfo,
    calendarDetail,
    loading,
    error,
    fetchPointsInfo,
    fetchCalendarDetail,
    refreshData,
    resetState,
  }
})
