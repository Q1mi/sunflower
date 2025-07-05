<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getPointRecords } from '@/api/adapter'
import type { PointRecord } from '@/types'

// 状态
const records = ref<PointRecord[]>([])
const loading = ref(false)
const hasMore = ref(false)

// 计算积分类型的图标和颜色
const getRecordStyle = (type: string) => {
  switch (type) {
    case 'check-in':
      return {
        icon: '✅',
        bgColor: 'bg-green-50',
        textColor: 'text-green-600',
        borderColor: 'border-green-200',
      }
    case 'retro-check-in':
      return {
        icon: '🔄',
        bgColor: 'bg-orange-50',
        textColor: 'text-orange-600',
        borderColor: 'border-orange-200',
      }
    case 'reward':
      return {
        icon: '🎁',
        bgColor: 'bg-yellow-50',
        textColor: 'text-yellow-600',
        borderColor: 'border-yellow-200',
      }
    case 'register':
      return {
        icon: '🎉',
        bgColor: 'bg-purple-50',
        textColor: 'text-purple-600',
        borderColor: 'border-purple-200',
      }
    default:
      return {
        icon: '💰',
        bgColor: 'bg-gray-50',
        textColor: 'text-gray-600',
        borderColor: 'border-gray-200',
      }
  }
}

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

// 方法
const fetchPointRecords = async () => {
  loading.value = true
  try {
    const result = await getPointRecords()
    records.value = result.records
    hasMore.value = result.hasMore
  } catch (error) {
    console.error('获取积分记录失败', error)
  } finally {
    loading.value = false
  }
}

// 生命周期钩子
onMounted(() => {
  fetchPointRecords()
})
</script>

<template>
  <div class="w-full h-screen bg-[#f9fafb] flex flex-col">
    <!-- 头部 -->
    <header class="bg-white shadow-sm border-b border-gray-100">
      <div class="flex items-center p-4">
        <button
          @click="$emit('back')"
          class="flex items-center justify-center w-8 h-8 rounded-full bg-[#86E3CE]/10 hover:bg-[#86E3CE]/20 text-[#86E3CE] hover:text-[#FA897B] transition-all duration-150 mr-3"
        >
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M15 19l-7-7 7-7"
            />
          </svg>
        </button>
        <h1 class="text-xl font-bold text-gray-800">积分明细</h1>
      </div>
    </header>

    <!-- 积分记录列表 -->
    <div class="flex-1 overflow-y-auto p-3">
      <div class="space-y-3">
        <!-- 加载状态 -->
        <div v-if="loading" class="flex items-center justify-center py-8">
          <div class="flex items-center space-x-2 text-[#86E3CE]">
            <div
              class="w-5 h-5 border-2 border-[#86E3CE] border-t-transparent rounded-full animate-spin"
            ></div>
            <span class="text-sm font-medium">加载中...</span>
          </div>
        </div>

        <!-- 积分记录 -->
        <template v-else-if="records.length > 0">
          <div
            v-for="record in records"
            :key="record.id"
            class="bg-white rounded-lg shadow-sm border border-gray-100 p-4 hover:shadow-md transition-shadow duration-150"
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-3">
                <!-- 图标 -->
                <div
                  :class="[
                    'w-10 h-10 rounded-full flex items-center justify-center border',
                    getRecordStyle(record.type).bgColor,
                    getRecordStyle(record.type).borderColor,
                  ]"
                >
                  <span class="text-lg">{{ getRecordStyle(record.type).icon }}</span>
                </div>

                <!-- 描述信息 -->
                <div>
                  <p class="font-medium text-gray-800 text-sm">{{ record.description }}</p>
                  <p class="text-xs text-gray-500 mt-1">{{ formatDate(record.date) }}</p>
                </div>
              </div>

              <!-- 积分变化 -->
              <div class="text-right">
                <span
                  class="text-lg font-bold"
                  :class="record.points > 0 ? 'text-[#86E3CE]' : 'text-[#FA897B]'"
                >
                  {{ record.points > 0 ? '+' : '' }}{{ record.points }}
                </span>
                <p class="text-xs text-gray-500 mt-1">积分</p>
              </div>
            </div>
          </div>
        </template>

        <!-- 空状态 -->
        <div
          v-else-if="!loading && records.length === 0"
          class="flex flex-col items-center justify-center py-12"
        >
          <div class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mb-4">
            <span class="text-2xl text-gray-400">📊</span>
          </div>
          <p class="text-gray-500 text-sm">暂无积分记录</p>
          <p class="text-gray-400 text-xs mt-1">完成签到后会显示积分变动记录</p>
        </div>

        <!-- 底部提示 -->
        <div v-if="!loading && !hasMore && records.length > 0" class="text-center py-6">
          <div class="inline-flex items-center space-x-2 text-gray-400">
            <div class="w-8 h-px bg-gray-300"></div>
            <span class="text-xs">没有更多记录了</span>
            <div class="w-8 h-px bg-gray-300"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
