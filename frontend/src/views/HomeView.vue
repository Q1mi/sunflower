<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import CheckInCalendar from '@/components/CheckInCalendar.vue'
import PointsDetail from '@/components/PointsDetail.vue'
import UserProfile from '@/components/UserProfile.vue'
import { useUserStore } from '@/stores/user'
import { useCheckinStore } from '@/stores/checkin'

// 获取用户状态和签到状态
const userStore = useUserStore()
const checkinStore = useCheckinStore()

// 状态
const showPointsDetail = ref(false)

// 获取当前时间
const currentTime = ref('')

// 更新时间
const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('en-US', {
    hour: 'numeric',
    minute: '2-digit',
    hour12: true,
  })
}

// 显示积分明细
const handleShowPointsDetail = () => {
  showPointsDetail.value = true
}

// 返回主页
const handleBack = () => {
  showPointsDetail.value = false
}

// 初始化数据的函数
const initializeData = async () => {
  if (userStore.currentUser) {
    const now = new Date()
    // 先获取积分统计数据
    await checkinStore.fetchPointsInfo()
    // 再获取签到日历数据
    await checkinStore.fetchCalendarDetail(now.getFullYear(), now.getMonth() + 1)
  }
}

// 生命周期钩子
onMounted(async () => {
  updateTime()
  // 每分钟更新一次时间
  setInterval(updateTime, 60000)

  // 如果用户已登录，立即获取数据
  if (userStore.currentUser) {
    await initializeData()
  }
})

// 创建日历组件的引用
const calendarRef = ref<InstanceType<typeof CheckInCalendar> | null>(null)

// 监听用户初始化状态，确保在用户状态初始化完成后获取积分数据
watch(
  () => userStore.isInitialized,
  async (isInitialized) => {
    if (isInitialized && userStore.currentUser) {
      await initializeData()
    }
  },
  { immediate: true },
)

// 监听用户状态变化，当用户登录状态改变时刷新数据
watch(
  () => userStore.currentUser,
  async (newUser, oldUser) => {
    // 只有在用户状态真正发生变化时才刷新数据
    if (newUser && (!oldUser || newUser.id !== oldUser.id)) {
      // 用户登录状态变化时，重新获取数据
      await initializeData()
    } else if (!newUser && oldUser) {
      // 用户登出时，重置状态
      checkinStore.resetState()
    }
  },
)
</script>

<template>
  <div class="w-full h-screen bg-[#f9fafb] mx-auto flex flex-col">
    <!-- 主页 -->
    <div v-if="!showPointsDetail" class="flex-grow overflow-y-auto flex flex-col bg-[#f9fafb]">
      <header class="p-3 flex justify-between items-center">
        <div class="flex flex-col items-start">
          <h1 class="text-xl font-bold text-gray-800">向阳花签到</h1>
          <p class="text-sm text-gray-500">{{ currentTime }}</p>
          <p class="text-xs text-gray-600">欢迎，{{ userStore.currentUser?.username || '游客' }}</p>
        </div>

        <!-- 用户个人资料 -->
        <UserProfile />
      </header>

      <!-- 积分信息展示 -->
      <div class="bg-white p-2.5 rounded-lg shadow mx-3 mb-2">
        <div class="flex justify-between items-center">
          <div class="flex items-center">
            <span class="text-base font-bold text-gray-800">我的积分</span>
            <span class="ml-2 text-xl font-bold text-amber-500">{{
              checkinStore.pointsInfo.totalPoints
            }}</span>
          </div>
          <button
            class="text-[#FFDD94] hover:text-[#FA897B] text-sm font-medium"
            @click="handleShowPointsDetail"
          >
            查看详情 >
          </button>
        </div>
      </div>

      <!-- 日历签到组件 -->
      <CheckInCalendar ref="calendarRef" />

      <!-- 补签说明 -->
      <div class="p-2 bg-green-50 rounded-lg border border-green-200 shadow-sm mx-3 mb-3">
        <h3 class="text-xs font-bold text-green-600 mb-1">📅 补签说明</h3>
        <div class="text-xs text-gray-700 space-y-0.5">
          <p class="text-amber-600 font-medium">补签需消耗100积分</p>
          <p class="text-gray-500">在日历中点击未签到的日期可进行补签</p>
          <p class="text-gray-500">连续签到可获得额外奖励，详见日历下方奖励格子</p>
        </div>
      </div>
    </div>

    <!-- 积分明细 -->
    <PointsDetail v-else @back="handleBack" />
  </div>
</template>
