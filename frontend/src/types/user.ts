// 用户相关类型定义

// 用户信息
export interface User {
  id: number
  username: string
  email: string
  avatar?: string
  createdAt: string
}

// 登录请求参数
export interface LoginRequest {
  username: string
  password: string
  remember?: boolean
}

// 注册请求参数
export interface RegisterRequest {
  username: string
  email: string
  password: string
  confirmPassword: string
}

// 登录/注册响应
export interface AuthResponse {
  user: User
  token: string
  message: string
}
