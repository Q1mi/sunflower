<script setup lang="ts">
import { computed, ref } from 'vue'
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'
import { Dialog, DialogPanel, DialogTitle, TransitionRoot, TransitionChild } from '@headlessui/vue'

const userStore = useUserStore()
const router = useRouter()

// 控制退出登录弹窗显示
const showLogoutModal = ref(false)

// 计算属性：用户头像
const userAvatar = computed(() => {
  return userStore.currentUser?.avatar || 'https://api.dicebear.com/7.x/avataaars/svg?seed=default'
})

// 计算属性：用户名
const username = computed(() => {
  return userStore.currentUser?.username || '游客'
})

// 显示退出登录确认弹窗
const handleLogout = () => {
  showLogoutModal.value = true
}

// 确认退出登录
const confirmLogout = () => {
  userStore.userLogout()
  router.push('/login')
  showLogoutModal.value = false
}

// 取消退出登录
const cancelLogout = () => {
  showLogoutModal.value = false
}
</script>

<template>
  <div class="user-profile flex items-center">
    <!-- 用户头像 -->
    <div class="avatar-container mr-3">
      <img
        :src="userAvatar"
        alt="用户头像"
        class="h-10 w-10 rounded-full border-2 border-white shadow-sm"
      />
    </div>

    <!-- 用户信息 -->
    <div class="user-info">
      <p class="font-medium text-gray-800">{{ username }}</p>
      <button
        @click="handleLogout"
        class="text-xs text-gray-500 hover:text-[#FA897B] transition-colors"
      >
        退出登录
      </button>
    </div>

    <!-- 退出登录确认弹窗 -->
    <TransitionRoot appear :show="showLogoutModal" as="template">
      <Dialog as="div" @close="cancelLogout" class="relative z-50">
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
                  <!-- 图标区域 -->
                  <div
                    class="mx-auto flex items-center justify-center h-18 w-18 rounded-full bg-gradient-to-br from-red-50 to-red-100 mb-6 shadow-lg"
                  >
                    <div
                      class="h-12 w-12 rounded-full bg-gradient-to-br from-red-500 to-red-600 flex items-center justify-center shadow-inner"
                    >
                      <svg
                        class="h-6 w-6 text-white"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                        />
                      </svg>
                    </div>
                  </div>

                  <!-- 标题和描述 -->
                  <DialogTitle as="h3" class="text-xl font-semibold text-red-700 mb-3">
                    退出登录
                  </DialogTitle>
                  <p class="text-sm text-gray-600 mb-8 leading-relaxed px-2">
                    确定要退出当前账户吗？<br />
                    <span class="text-gray-500">退出后需要重新登录才能使用完整功能</span>
                  </p>

                  <!-- 按钮组 -->
                  <div class="flex flex-col space-y-3">
                    <button
                      @click="confirmLogout"
                      class="w-full inline-flex justify-center items-center px-6 py-3 text-sm font-medium text-white bg-gradient-to-r from-red-500 to-red-600 rounded-xl shadow-lg hover:from-red-600 hover:to-red-700 hover:shadow-xl focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 transform hover:scale-[1.02] transition-all duration-200"
                    >
                      <svg
                        class="h-4 w-4 mr-2"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                        />
                      </svg>
                      确定退出
                    </button>
                    <button
                      @click="cancelLogout"
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
  </div>
</template>
