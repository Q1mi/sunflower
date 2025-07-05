<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { userRegister } from '@/api'
import type { RegisterRequest } from '@/types'

const router = useRouter()

// 表单数据
const formData = reactive<RegisterRequest>({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
})

// 状态
const loading = ref(false)
const errorMessage = ref('')

// 注册方法
const handleRegister = async () => {
  // 表单验证
  if (!formData.username || !formData.email || !formData.password) {
    errorMessage.value = '请填写所有必填字段'
    return
  }

  if (formData.password !== formData.confirmPassword) {
    errorMessage.value = '两次输入的密码不一致'
    return
  }

  if (formData.password.length < 6) {
    errorMessage.value = '密码长度不能少于6位'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    await userRegister(formData) // 注册成功后跳转到登录页
    router.push('/login')
  } catch (error: unknown) {
    console.error('注册失败', error)
    errorMessage.value = error instanceof Error ? error.message : '注册失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

// 跳转到登录页
const goToLogin = () => {
  router.push('/login')
}
</script>

<template>
  <div
    class="min-h-screen flex items-center justify-center bg-[#FFDD94]/10 py-12 px-4 sm:px-6 lg:px-8"
  >
    <div class="max-w-md w-full space-y-8 bg-white p-8 rounded-lg shadow-md">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">注册账号</h2>
        <p class="mt-2 text-center text-sm text-gray-600">加入向阳花签到，开始您的签到之旅</p>
      </div>

      <form class="mt-8 space-y-6" @submit.prevent="handleRegister">
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
            <label for="email" class="sr-only">邮箱</label>
            <input
              id="email"
              name="email"
              type="email"
              required
              v-model="formData.email"
              class="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-[#86E3CE] focus:border-[#86E3CE] focus:z-10 sm:text-sm"
              placeholder="邮箱"
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
              placeholder="密码（至少6位）"
            />
          </div>
          <div>
            <label for="confirmPassword" class="sr-only">确认密码</label>
            <input
              id="confirmPassword"
              name="confirmPassword"
              type="password"
              required
              v-model="formData.confirmPassword"
              class="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-[#86E3CE] focus:border-[#86E3CE] focus:z-10 sm:text-sm"
              placeholder="确认密码"
            />
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
            {{ loading ? '注册中...' : '注册' }}
          </button>
        </div>
      </form>

      <div class="text-center mt-4">
        <p class="text-sm text-gray-600">
          已有账号？
          <button @click="goToLogin" class="font-medium text-[#FA897B] hover:text-[#FFDD94]">
            立即登录
          </button>
        </p>
      </div>
    </div>
  </div>
</template>
