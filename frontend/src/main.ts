import { createApp } from 'vue'
import { createPinia } from 'pinia'
import Vue3Toastify, { type ToastContainerOptions } from 'vue3-toastify'
import 'vue3-toastify/dist/index.css'

import App from './App.vue'
import router from './router'
import { useUserStore } from './stores/user'

// 先导入应用样式
import './assets/main.css'
// 后导入Font Awesome图标，确保图标样式优先级更高
import '@fortawesome/fontawesome-free/css/all.min.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(Vue3Toastify, {
  autoClose: 3000,
  position: 'top-center',
  hideProgressBar: false,
  closeOnClick: true,
  pauseOnHover: true,
  draggable: true,
  theme: 'light',
} as ToastContainerOptions)

// 初始化用户状态
const userStore = useUserStore(pinia)
userStore.initUserState()

app.mount('#app')
