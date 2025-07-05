/**
 * 用户账户相关API接口
 * @description 提供用户注册、登录、刷新令牌和获取用户信息的API接口
 */

import http from './http'
import { API_ENDPOINTS } from './constants'
import type {
  CreateUserRequest,
  CreateUserResponse,
  LoginRequest,
  LoginResponse,
  RefreshTokenRequest,
  RefreshTokenResponse,
  UserProfileResponse,
} from './types'

/**
 * 创建用户
 * @description 注册新用户账户
 * @param data 用户注册信息
 * @returns 创建结果
 * @example
 * ```typescript
 * const result = await createUser({
 *   username: 'testuser',
 *   password: 'password123',
 *   confirmPassword: 'password123'
 * })
 * ```
 */
export const createUser = async (data: CreateUserRequest): Promise<CreateUserResponse> => {
  const response = await http.post<CreateUserResponse>(API_ENDPOINTS.USER.CREATE, data)
  return response.data
}

/**
 * 用户登录
 * @description 用户登录获取访问令牌
 * @param data 登录凭据
 * @returns 登录结果，包含访问令牌和刷新令牌
 * @example
 * ```typescript
 * const result = await login({
 *   username: 'testuser',
 *   password: 'password123'
 * })
 * console.log(result.accessToken) // 访问令牌
 * ```
 */
export const login = async (data: LoginRequest): Promise<LoginResponse> => {
  const response = await http.post<LoginResponse>(API_ENDPOINTS.USER.LOGIN, data)
  return response.data
}

/**
 * 刷新访问令牌
 * @description 使用刷新令牌获取新的访问令牌
 * @param data 刷新令牌信息
 * @returns 新的令牌对
 * @example
 * ```typescript
 * const result = await refreshToken({
 *   refreshToken: 'your-refresh-token'
 * })
 * ```
 */
export const refreshToken = async (data: RefreshTokenRequest): Promise<RefreshTokenResponse> => {
  const response = await http.post<RefreshTokenResponse>(API_ENDPOINTS.USER.REFRESH, data)
  return response.data
}

/**
 * 获取用户个人信息
 * @description 获取当前登录用户的详细信息
 * @returns 用户信息
 * @example
 * ```typescript
 * const profile = await getUserProfile()
 * console.log(profile.username) // 用户名
 * ```
 */
export const getUserProfile = async (): Promise<UserProfileResponse> => {
  const response = await http.get<UserProfileResponse>(API_ENDPOINTS.USER.PROFILE)
  return response.data
}
