import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { requiresAuth: true },
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
      meta: { guest: true },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/RegisterView.vue'),
      meta: { guest: true },
    },
  ],
})

// 全局前置守卫
router.beforeEach(async (to, from, next) => {
  // 检查路由是否需要认证
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth)
  const isGuestOnly = to.matched.some((record) => record.meta.guest)

  // 使用 userStore 获取用户状态
  const userStore = useUserStore()
  
  // 确保用户状态已初始化
  if (!userStore.isInitialized) {
    await userStore.initUserState()
  }

  // 使用 store 中的用户状态
  const currentUser = userStore.currentUser

  if (requiresAuth && !currentUser) {
    // 需要认证但用户未登录，重定向到登录页
    next({ name: 'login' })
  } else if (isGuestOnly && currentUser) {
    // 仅限游客但用户已登录，重定向到首页
    next({ name: 'home' })
  } else {
    // 其他情况正常通过
    next()
  }
})

export default router
