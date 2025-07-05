<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import type { LoginRequest } from '@/types'

const router = useRouter()
const userStore = useUserStore()

// 表单数据
const formData = reactive<LoginRequest>({
  username: '',
  password: '',
  remember: true,
})

// 状态
const loading = ref(false)
const errorMessage = ref('')

// 登录方法
const handleLogin = async () => {
  // 表单验证
  if (!formData.username || !formData.password) {
    errorMessage.value = '用户名和密码不能为空'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    await userStore.userLogin(formData) // 登录成功后跳转到首页
    router.push('/')
  } catch (error: unknown) {
    console.error('登录失败', error)
    errorMessage.value = error instanceof Error ? error.message : '登录失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

// 跳转到注册页
const goToRegister = () => {
  router.push('/register')
}
</script>

<template>
  <div
    class="min-h-screen flex items-center justify-center bg-[#FFDD94]/10 py-12 px-4 sm:px-6 lg:px-8"
  >
    <div class="max-w-md w-full space-y-8 bg-white p-8 rounded-lg shadow-md">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">登录向阳花签到</h2>
        <p class="mt-2 text-center text-sm text-gray-600">每日签到，收获积分和惊喜</p>
      </div>

      <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
        <div class="rounded-md shadow-sm space-y-4">
          <div>
            <label for="username" class="sr-only">用户名</label>
            <input
              id="username"
              name="username"
              type="text"
              required
              v-model="formData.username"
              class="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-[#86E3CE] focus:border-[#86E3CE] focus:z-10 sm:text-sm"
              placeholder="用户名"
            />
          </div>
          <div>
            <label for="password" class="sr-only">密码</label>
            <input
              id="password"
              name="password"
              type="password"
              required
              v-model="formData.password"
              class="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-[#86E3CE] focus:border-[#86E3CE] focus:z-10 sm:text-sm"
              placeholder="密码"
            />
          </div>
        </div>

        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <input
              id="remember"
              name="remember"
              type="checkbox"
              v-model="formData.remember"
              class="h-4 w-4 text-amber-600 focus:ring-amber-500 border-gray-300 rounded"
            />
            <label for="remember" class="ml-2 block text-sm text-gray-900"> 记住我 </label>
          </div>

          <div class="text-sm">
            <a href="#" class="font-medium text-amber-600 hover:text-amber-500"> 忘记密码？ </a>
          </div>
        </div>

        <div>
          <p v-if="errorMessage" class="text-red-500 text-sm mb-4">{{ errorMessage }}</p>
          <button
            type="submit"
            :disabled="loading"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-[#86E3CE] hover:bg-[#D0E6A5] focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-[#FFDD94]"
          >
            <span v-if="loading" class="absolute left-0 inset-y-0 flex items-center pl-3">
              <!-- 加载图标 -->
              <svg
                class="animate-spin h-5 w-5 text-white"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  class="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="4"
                ></circle>
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
            </span>
            {{ loading ? '登录中...' : '登录' }}
          </button>
        </div>
      </form>

      <div class="text-center mt-4">
        <p class="text-sm text-gray-600">
          还没有账号？
          <button @click="goToRegister" class="font-medium text-[#FA897B] hover:text-[#FFDD94]">
            立即注册
          </button>
        </p>
      </div>
    </div>
  </div>
</template>
