<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { checkIn, retroCheckIn } from '@/api/adapter'
import { useCheckinStore } from '@/stores/checkin'
import type { CalendarDay } from '@/types'
import { Dialog, DialogPanel, DialogTitle, TransitionRoot, TransitionChild } from '@headlessui/vue'
import { toast } from 'vue3-toastify'
import 'vue3-toastify/dist/index.css'

// 使用签到状态管理
const checkinStore = useCheckinStore()

// 当前日期
const today = new Date()
const currentYear = today.getFullYear()
const currentMonth = today.getMonth()
const currentDay = today.getDate()

// 本地状态
const loading = ref(false)

// 日历数据
const year = ref(currentYear)
const month = ref(currentMonth)

// 计算属性
const daysInMonth = computed(() => {
  return new Date(year.value, month.value + 1, 0).getDate()
})

const firstDayOfMonth = computed(() => {
  return new Date(year.value, month.value, 1).getDay()
})

const monthName = computed(() => {
  return new Date(year.value, month.value, 1).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
  })
})

// 奖励日
// const rewardDays = [7, 15, 25, daysInMonth.value] // 移除未使用的变量

// 连续签到奖励配置
const consecutiveRewards = [
  { days: 3, points: 5, label: '连续3天', icon: '🔥' },
  { days: 7, points: 10, label: '连续7天', icon: '⭐' },
  { days: 15, points: 20, label: '连续15天', icon: '💎' },
  { days: 30, points: 100, label: '月满签', icon: '👑' },
]

// 计算奖励完成状态
const rewardStatus = computed(() => {
  // 使用日历详情中的连续天数，这个数据更准确且实时更新
  const consecutiveDays = checkinStore.calendarDetail.consecutiveDays
  return consecutiveRewards.map((reward) => ({
    ...reward,
    completed: consecutiveDays >= reward.days,
  }))
})

// 生成日历数据
const calendarDays = computed(() => {
  const days: CalendarDay[] = []

  // 添加月初前的空白天
  for (let i = 0; i < firstDayOfMonth.value; i++) {
    days.push({
      day: 0,
      isToday: false,
      isCheckedIn: false,
      isRetroCheckedIn: false,
      isPast: false,
      isRewardDay: false,
    })
  }

  // 添加当月的天数
  for (let day = 1; day <= daysInMonth.value; day++) {
    const isPast =
      year.value < currentYear ||
      (year.value === currentYear && month.value < currentMonth) ||
      (year.value === currentYear && month.value === currentMonth && day < currentDay)

    const isToday = year.value === currentYear && month.value === currentMonth && day === currentDay

    days.push({
      day,
      isToday,
      isCheckedIn: checkinStore.calendarDetail.checkedInDays.includes(day),
      isRetroCheckedIn: checkinStore.calendarDetail.retroCheckedInDays.includes(day),
      isPast,
      isRewardDay: false,
    })
  }

  return days
})

// 弹窗控制
const showRetroModal = ref(false)
const retroDay = ref(0)
const retroMessage = ref('')

const handleCheckIn = async () => {
  if (checkinStore.calendarDetail.isCheckedInToday) return

  loading.value = true
  try {
    const result = await checkIn()
    if (result.success) {
      // 重新获取积分信息和日历详情，确保数据同步
      await checkinStore.fetchPointsInfo()
      await checkinStore.fetchCalendarDetail(year.value, month.value + 1)

      toast.success(`🎉 签到成功！获得${result.points}积分`, {
        position: 'top-center',
        autoClose: 3000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
      })
    }
  } catch (error) {
    console.error('签到失败', error)
    toast.error('❌ 签到失败，请稍后重试', {
      position: 'top-center',
      autoClose: 3000,
      hideProgressBar: false,
      closeOnClick: true,
      pauseOnHover: true,
    })
  } finally {
    loading.value = false
  }
}

const handleRetroCheckIn = async (day: number) => {
  if (checkinStore.pointsInfo.totalPoints < 100) {
    retroMessage.value = `积分不足，无法补签！当前积分：${checkinStore.pointsInfo.totalPoints}，需要积分：100`
    showRetroModal.value = true
    return
  }
  if (checkinStore.calendarDetail.remainRetroTimes <= 0) {
    retroMessage.value = '本月补签次数已用完！'
    showRetroModal.value = true
    return
  }

  retroDay.value = day
  retroMessage.value = `确定消耗100积分和1次补签机会，补签 ${month.value + 1}月${day}日 吗？`
  showRetroModal.value = true
}

const confirmRetroCheckIn = async () => {
  loading.value = true
  try {
    // 构造日期字符串 YYYY-MM-DD
    const date = `${year.value}-${String(month.value + 1).padStart(2, '0')}-${String(retroDay.value).padStart(2, '0')}`
    const result = await retroCheckIn(date)

    if (result.success) {
      // 重新获取积分信息和日历详情，确保数据同步
      await checkinStore.fetchPointsInfo()
      await checkinStore.fetchCalendarDetail(year.value, month.value + 1)
      retroMessage.value = result.message
    } else {
      retroMessage.value = '补签失败，请稍后重试'
    }
  } catch (error) {
    console.error('补签失败', error)
    retroMessage.value = '补签失败，请稍后重试'
  } finally {
    loading.value = false
    showRetroModal.value = true
  }
}

const prevMonth = () => {
  if (month.value === 0) {
    year.value--
    month.value = 11
  } else {
    month.value--
  }
}

const nextMonth = () => {
  if (month.value === 11) {
    year.value++
    month.value = 0
  } else {
    month.value++
  }
}

// 监听年月变化，重新获取日历详情
watch(
  [year, month],
  () => {
    checkinStore.fetchCalendarDetail(year.value, month.value + 1)
  },
  { immediate: false },
)

// 生命周期钩子
onMounted(() => {
  // 如果当前月份的数据还没有加载，则加载
  if (checkinStore.calendarDetail.checkedInDays.length === 0) {
    checkinStore.fetchCalendarDetail(year.value, month.value + 1)
  }
})

// 暴露给父组件的方法
defineExpose({
  fetchCalendarDetail: () => checkinStore.fetchCalendarDetail(year.value, month.value + 1),
})
</script>

<template>
  <section class="flex-grow p-2">
    <!-- 签到统计信息 -->
    <div
      class="text-xs py-1 flex items-center justify-center border-b border-green-100 mb-2 bg-green-50 rounded-t-md"
    >
      <span class="text-gray-600">当月连签</span>
      <span class="font-bold text-green-500 mx-1">{{
        checkinStore.calendarDetail.consecutiveDays
      }}</span>
      <span class="mr-3 text-gray-600">天</span>
      <span class="text-gray-600">本月可补签</span>
      <span class="font-bold text-orange-500 mx-1">{{
        checkinStore.calendarDetail.remainRetroTimes
      }}</span>
      <span class="text-gray-600">次</span>
    </div>

    <div class="mb-2 flex justify-between items-center">
      <button class="text-gray-500 hover:text-[#86E3CE] px-2 py-1" @click="prevMonth">
        <span class="text-xl">&lt;</span>
      </button>
      <h2 class="text-base font-semibold text-gray-700">{{ monthName }}</h2>
      <button class="text-gray-500 hover:text-[#86E3CE] px-2 py-1" @click="nextMonth">
        <span class="text-xl">&gt;</span>
      </button>
    </div>

    <!-- 星期标题 -->
    <div class="grid grid-cols-7 gap-1 text-center mb-2">
      <div class="font-medium text-green-600 text-xs">日</div>
      <div class="font-medium text-gray-600 text-xs">一</div>
      <div class="font-medium text-gray-600 text-xs">二</div>
      <div class="font-medium text-gray-600 text-xs">三</div>
      <div class="font-medium text-gray-600 text-xs">四</div>
      <div class="font-medium text-gray-600 text-xs">五</div>
      <div class="font-medium text-green-600 text-xs">六</div>
    </div>

    <!-- 日历天数 -->
    <div class="grid grid-cols-7 gap-1">
      <template v-for="(day, index) in calendarDays" :key="index">
        <!-- 空白天 -->
        <div
          v-if="day.day === 0"
          class="aspect-square flex flex-col justify-center items-center text-sm bg-gray-50"
        ></div>

        <!-- 有效天数 -->
        <div
          v-else
          class="aspect-square flex flex-col justify-center items-center text-sm rounded-md cursor-pointer transition-all duration-150 ease-in-out"
          :class="{
            'bg-[#86E3CE]/20 text-[#86E3CE]':
              (day.isCheckedIn || day.isRetroCheckedIn) && !day.isToday,
            'bg-[#86E3CE] text-white': (day.isCheckedIn || day.isRetroCheckedIn) && day.isToday,
            'bg-[#FFDD94]/20 hover:bg-[#FFDD94]/30':
              day.isPast && !day.isCheckedIn && !day.isRetroCheckedIn,
            'border-2 border-[#86E3CE]': day.isToday,
            'bg-[#FFDD94]/30 text-amber-700':
              day.isToday && !day.isCheckedIn && !day.isRetroCheckedIn,
            'text-gray-400 bg-gray-100 cursor-default': !day.isPast && !day.isToday,
          }"
        >
          <!-- 日期数字 -->
          <span
            class="font-medium"
            :class="{ 'font-bold text-lg': day.isToday && !day.isCheckedIn }"
          >
            {{ day.day }}
            <span v-if="day.isRetroCheckedIn" class="text-xs text-orange-400 font-bold ml-1"
              >补</span
            >
          </span>

          <!-- 已签到标记 -->
          <span
            v-if="!day.isToday && (day.isCheckedIn || day.isRetroCheckedIn)"
            class="text-green-500 text-xs mt-0.5"
            >已签</span
          >

          <!-- 今日已签到文本 -->
          <template v-if="day.isToday && day.isCheckedIn">
            <span class="text-xs font-bold">已签</span>
          </template>

          <!-- 补签按钮 -->
          <button
            v-if="
              day.isPast &&
              !(day.isCheckedIn || day.isRetroCheckedIn) &&
              checkinStore.calendarDetail.remainRetroTimes > 0 &&
              checkinStore.pointsInfo.totalPoints >= 100
            "
            class="text-xs text-[#FA897B] mt-0.5"
            @click="handleRetroCheckIn(day.day)"
          >
            补签
          </button>
        </div>
      </template>
    </div>

    <!-- 连续签到奖励格子 -->
    <div
      class="mt-4 p-2 bg-gradient-to-r from-yellow-50 to-orange-50 rounded-lg border border-yellow-200"
    >
      <h3 class="text-xs font-bold text-yellow-700 mb-2 text-center">🎁 连续签到奖励</h3>
      <div class="grid grid-cols-4 gap-1.5">
        <div
          v-for="(reward, index) in rewardStatus"
          :key="index"
          class="relative flex flex-col items-center p-1.5 rounded-md border-2 transition-all duration-200"
          :class="{
            'bg-green-100 border-green-300': reward.completed,
            'bg-white border-gray-200 hover:border-yellow-300': !reward.completed,
          }"
        >
          <!-- 完成标记 -->
          <div
            v-if="reward.completed"
            class="absolute -top-0.5 -right-0.5 w-4 h-4 bg-green-500 rounded-full flex items-center justify-center"
          >
            <svg
              class="w-2.5 h-2.5 text-white"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="3"
                d="M5 13l4 4L19 7"
              />
            </svg>
          </div>

          <!-- 奖励图标 -->
          <div class="text-sm mb-0.5">{{ reward.icon }}</div>

          <!-- 奖励信息 -->
          <div class="text-center">
            <div class="text-xs font-medium text-gray-700 leading-tight">{{ reward.label }}</div>
            <div
              class="text-xs font-bold"
              :class="{
                'text-green-600': reward.completed,
                'text-yellow-600': !reward.completed,
              }"
            >
              +{{ reward.points }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>

  <footer class="p-3 border-t border-gray-200 bg-white">
    <button
      :class="{
        'w-full bg-[#86E3CE] hover:bg-[#FA897B] text-white font-bold py-2.5 px-4 rounded-lg shadow-md transition duration-150 ease-in-out':
          !checkinStore.calendarDetail.isCheckedInToday,
        'w-full bg-gray-400 cursor-not-allowed text-white font-bold py-2.5 px-4 rounded-lg shadow-md':
          checkinStore.calendarDetail.isCheckedInToday,
      }"
      :disabled="checkinStore.calendarDetail.isCheckedInToday"
      @click="handleCheckIn"
    >
      {{ checkinStore.calendarDetail.isCheckedInToday ? '今日已签到' : '今日签到' }}
    </button>
  </footer>

  <!-- 补签确认弹窗 -->
  <TransitionRoot appear :show="showRetroModal" as="template">
    <Dialog as="div" @close="showRetroModal = false" class="relative z-50">
      <TransitionChild
        as="template"
        enter="duration-300 ease-out"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="duration-200 ease-in"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <div class="fixed inset-0 bg-black/25 backdrop-blur-sm" />
      </TransitionChild>

      <div class="fixed inset-0 overflow-y-auto">
        <div class="flex min-h-full items-center justify-center p-4 text-center">
          <TransitionChild
            as="template"
            enter="duration-300 ease-out"
            enter-from="opacity-0 scale-95"
            enter-to="opacity-100 scale-100"
            leave="duration-200 ease-in"
            leave-from="opacity-100 scale-100"
            leave-to="opacity-0 scale-95"
          >
            <DialogPanel
              class="w-full max-w-sm transform overflow-hidden rounded-2xl bg-white p-6 text-left align-middle shadow-xl transition-all"
            >
              <div class="text-center py-2 pb-6">
                <!-- 动态图标区域 -->
                <div
                  class="mx-auto flex items-center justify-center h-18 w-18 rounded-full mb-6 shadow-lg"
                  :class="{
                    'bg-gradient-to-br from-green-50 to-emerald-100':
                      !retroMessage.includes('失败') &&
                      !retroMessage.includes('用完') &&
                      !retroMessage.includes('积分不足'),
                    'bg-gradient-to-br from-red-50 to-red-100':
                      retroMessage.includes('失败') ||
                      retroMessage.includes('用完') ||
                      retroMessage.includes('积分不足'),
                  }"
                >
                  <div
                    class="h-12 w-12 rounded-full flex items-center justify-center shadow-inner"
                    :class="{
                      'bg-gradient-to-br from-green-500 to-emerald-600':
                        !retroMessage.includes('失败') &&
                        !retroMessage.includes('用完') &&
                        !retroMessage.includes('积分不足'),
                      'bg-gradient-to-br from-red-500 to-red-600':
                        retroMessage.includes('失败') ||
                        retroMessage.includes('用完') ||
                        retroMessage.includes('积分不足'),
                    }"
                  >
                    <!-- 确认图标 -->
                    <svg
                      v-if="
                        !retroMessage.includes('成功') &&
                        !retroMessage.includes('失败') &&
                        !retroMessage.includes('用完') &&
                        !retroMessage.includes('积分不足')
                      "
                      class="h-6 w-6 text-white"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                    <!-- 成功图标 -->
                    <svg
                      v-else-if="retroMessage.includes('成功')"
                      class="h-6 w-6 text-white"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                    <!-- 失败/用完图标 -->
                    <svg
                      v-else
                      class="h-6 w-6 text-white"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L4.314 15.5c-.77.833.192 2.5 1.732 2.5z"
                      />
                    </svg>
                  </div>
                </div>

                <!-- 标题 -->
                <DialogTitle
                  as="h3"
                  class="text-xl font-semibold mb-3"
                  :class="{
                    'text-green-700':
                      !retroMessage.includes('失败') &&
                      !retroMessage.includes('用完') &&
                      !retroMessage.includes('积分不足'),
                    'text-red-700':
                      retroMessage.includes('失败') ||
                      retroMessage.includes('用完') ||
                      retroMessage.includes('积分不足'),
                  }"
                >
                  {{
                    retroMessage.includes('成功') ||
                    retroMessage.includes('失败') ||
                    retroMessage.includes('用完') ||
                    retroMessage.includes('积分不足')
                      ? '补签结果'
                      : '补签确认'
                  }}
                </DialogTitle>

                <!-- 积分奖励显示（仅成功时） -->
                <div v-if="retroMessage.includes('成功')" class="mb-4"></div>

                <p class="text-sm text-gray-600 mb-8 leading-relaxed px-2">{{ retroMessage }}</p>

                <!-- 结果按钮（成功/失败/用完） -->
                <div
                  v-if="
                    retroMessage.includes('成功') ||
                    retroMessage.includes('失败') ||
                    retroMessage.includes('用完') ||
                    retroMessage.includes('积分不足')
                  "
                  class="flex justify-center"
                >
                  <button
                    @click="showRetroModal = false"
                    class="w-full inline-flex justify-center items-center px-6 py-3 text-sm font-medium text-white rounded-xl shadow-lg focus:outline-none focus:ring-2 focus:ring-offset-2 transform hover:scale-[1.02] transition-all duration-200"
                    :class="{
                      'bg-gradient-to-r from-green-500 to-emerald-600 hover:from-green-600 hover:to-emerald-700 hover:shadow-xl focus:ring-green-500':
                        retroMessage.includes('成功'),
                      'bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 hover:shadow-xl focus:ring-red-500':
                        retroMessage.includes('失败') ||
                        retroMessage.includes('用完') ||
                        retroMessage.includes('积分不足'),
                    }"
                  >
                    <svg class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M5 13l4 4L19 7"
                      />
                    </svg>
                    知道了
                  </button>
                </div>

                <!-- 确认按钮组 -->
                <div v-else class="flex flex-col space-y-3">
                  <button
                    @click="confirmRetroCheckIn"
                    class="w-full inline-flex justify-center items-center px-6 py-3 text-sm font-medium text-white bg-gradient-to-r from-green-500 to-emerald-600 rounded-xl shadow-lg hover:from-green-600 hover:to-emerald-700 hover:shadow-xl focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 transform hover:scale-[1.02] transition-all duration-200"
                  >
                    <svg class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                    确定补签
                  </button>
                  <button
                    @click="showRetroModal = false"
                    class="w-full inline-flex justify-center items-center px-6 py-3 text-sm font-medium text-gray-700 bg-gray-50 rounded-xl border border-gray-200 hover:bg-gray-100 hover:border-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-300 focus:ring-offset-2 transition-all duration-200"
                  >
                    取消
                  </button>
                </div>
              </div>
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>
