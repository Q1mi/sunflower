// 用户状态管理
import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  userLogin as apiUserLogin,
  userRegister as apiUserRegister,
  userLogout as apiUserLogout,
  getCurrentUser,
} from '@/api'
import type { User, LoginRequest, RegisterRequest } from '@/types'

export const useUserStore = defineStore('user', () => {
  // 状态
  const currentUser = ref<User | null>(null)
  const loading = ref(false)
  const error = ref('')
  const isInitialized = ref(false) // 添加初始化状态标记

  // 初始化用户状态
  async function initUserState() {
    // 如果已经初始化过，则跳过
    if (isInitialized.value) {
      return
    }

    loading.value = true
    error.value = ''
    try {
      const user = await getCurrentUser()
      currentUser.value = user
      isInitialized.value = true
    } catch (err: unknown) {
      console.error('获取用户信息失败', err)
      error.value = err instanceof Error ? err.message : '获取用户信息失败'
      currentUser.value = null
      isInitialized.value = true // 即使失败也标记为已初始化
    } finally {
      loading.value = false
    }
  }

  // 用户登录
  async function userLogin(loginData: LoginRequest) {
    loading.value = true
    error.value = ''
    try {
      const response = await apiUserLogin(loginData, loginData.remember || false)
      currentUser.value = response.user
      isInitialized.value = true
      return response
    } catch (err: unknown) {
      error.value = err instanceof Error ? err.message : '登录失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  // 用户注册
  async function userRegister(registerData: RegisterRequest) {
    loading.value = true
    error.value = ''
    try {
      const response = await apiUserRegister(registerData)
      currentUser.value = response.user
      isInitialized.value = true
      return response
    } catch (err: unknown) {
      error.value = err instanceof Error ? err.message : '注册失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  // 用户登出
  function userLogout() {
    currentUser.value = null
    isInitialized.value = false
    apiUserLogout()
  }

  // 判断用户是否已登录
  const isLoggedIn = () => !!currentUser.value

  return {
    currentUser,
    loading,
    error,
    isInitialized,
    initUserState,
    userLogin,
    userRegister,
    userLogout,
    isLoggedIn,
  }
})
