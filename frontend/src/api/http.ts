/**
 * HTTP 客户端配置
 * @description 基于axios的HTTP客户端，提供统一的请求/响应处理、认证和错误处理
 */

import axios from 'axios'
import type { AxiosInstance, InternalAxiosRequestConfig, AxiosResponse, AxiosError } from 'axios'

import { API_CONFIG, API_ENDPOINTS, HTTP_STATUS, BUSINESS_ERROR_CODES } from './constants'
import { tokenManager } from './utils/token-manager'
import { errorHandler } from './utils/error-handler'

/** 扩展AxiosRequestConfig类型以支持_retry属性 */
interface ExtendedAxiosRequestConfig extends InternalAxiosRequestConfig {
  /** 是否为重试请求 */
  _retry?: boolean
}

/** 通用API响应类型 */
export interface HttpApiResponse<T = unknown> {
  /** 响应码 */
  code?: number
  /** 响应消息 */
  message?: string
  /** 响应数据 */
  data?: T
}

/** HTTP客户端类 */
class HttpClient {
  /** axios实例 */
  private instance: AxiosInstance

  constructor() {
    this.instance = this.createAxiosInstance()
    this.setupInterceptors()
  }

  /**
   * 创建axios实例
   * @returns axios实例
   */
  private createAxiosInstance(): AxiosInstance {
    return axios.create({
      baseURL: API_CONFIG.BASE_URL,
      timeout: API_CONFIG.TIMEOUT,
      headers: API_CONFIG.DEFAULT_HEADERS,
    })
  }

  /**
   * 设置拦截器
   */
  private setupInterceptors(): void {
    this.setupRequestInterceptor()
    this.setupResponseInterceptor()
  }

  /**
   * 设置请求拦截器
   */
  private setupRequestInterceptor(): void {
    this.instance.interceptors.request.use(
      (config: InternalAxiosRequestConfig) => {
        // 添加认证token
        const accessToken = tokenManager.getAccessToken()
        if (accessToken && config.headers) {
          config.headers.Authorization = `Bearer ${accessToken}`
        }

        // 请求日志
        console.log(
          `[HTTP Request] ${config.method?.toUpperCase()} ${config.url}`,
          config.data || config.params,
        )

        return config
      },
      (error) => {
        console.error('[HTTP Request Error]', error)
        return Promise.reject(error)
      },
    )
  }

  /**
   * 设置响应拦截器
   */
  private setupResponseInterceptor(): void {
    this.instance.interceptors.response.use(
      (response: AxiosResponse) => {
        // 响应日志
        console.log(`[HTTP Response] ${response.status} ${response.config.url}`, response.data)

        const apiResponse: HttpApiResponse = response.data

        // 检查业务响应状态
        if (apiResponse.code && apiResponse.code !== BUSINESS_ERROR_CODES.SUCCESS) {
          const error = errorHandler.handleBusinessError(
            apiResponse.code,
            apiResponse.message || 'Business error',
            {
              url: response.config.url,
              method: response.config.method,
              responseData: apiResponse,
            },
          )
          throw new Error(error.userMessage)
        }

        // 提取业务数据
        response.data = apiResponse.data ?? {}
        return response
      },
      async (error) => {
        return this.handleResponseError(error)
      },
    )
  }

  /**
   * 处理响应错误
   * @param error 错误对象
   * @returns Promise
   */
  private async handleResponseError(error: AxiosError): Promise<never> {
    const standardError = errorHandler.handleHttpError(error)

    // 处理401错误（token过期）
    if (standardError.code === HTTP_STATUS.UNAUTHORIZED) {
      const originalRequest = error.config as ExtendedAxiosRequestConfig

      // 如果是刷新token接口返回401，直接跳转登录页
      if (originalRequest?.url?.includes(API_ENDPOINTS.USER.REFRESH)) {
        this.redirectToLogin()
        return Promise.reject(new Error(standardError.userMessage))
      }

      // 如果有refreshToken且不是重试请求，尝试刷新token
      const refreshToken = tokenManager.getRefreshToken()
      if (refreshToken && originalRequest && !originalRequest._retry) {
        originalRequest._retry = true

        try {
          const newAccessToken = await this.refreshAccessToken()

          // 更新原请求的Authorization头
          if (originalRequest.headers) {
            originalRequest.headers.Authorization = `Bearer ${newAccessToken}`
          }

          // 重新发送原请求
          return this.instance(originalRequest)
        } catch (refreshError) {
          console.error('[Token Refresh Failed]', refreshError)
          this.redirectToLogin()
          return Promise.reject(new Error(standardError.userMessage))
        }
      } else {
        // 没有refreshToken或已经重试过，跳转到登录页
        this.redirectToLogin()
      }
    }

    return Promise.reject(new Error(standardError.userMessage))
  }

  /**
   * 刷新访问令牌
   * @returns 新的访问令牌
   */
  private async refreshAccessToken(): Promise<string> {
    // 检查是否已有正在进行的刷新请求
    let refreshPromise = tokenManager.getRefreshPromise()
    if (refreshPromise) {
      return refreshPromise
    }

    const refreshToken = tokenManager.getRefreshToken()
    if (!refreshToken) {
      throw new Error('No refresh token available')
    }

    // 创建刷新请求
    refreshPromise = this.doRefreshToken(refreshToken)
    tokenManager.setRefreshPromise(refreshPromise)

    try {
      const newAccessToken = await refreshPromise
      return newAccessToken
    } finally {
      // 清除刷新Promise
      tokenManager.setRefreshPromise(null)
    }
  }

  /**
   * 执行token刷新
   * @param refreshToken 刷新令牌
   * @returns 新的访问令牌
   */
  private async doRefreshToken(refreshToken: string): Promise<string> {
    try {
      const response = await axios.post(`${API_CONFIG.BASE_URL}${API_ENDPOINTS.USER.REFRESH}`, {
        refreshToken,
      })

      const { accessToken: newAccessToken, refreshToken: newRefreshToken } = response.data.data

      // 更新token（保持原有的记住登录状态）
      const remember = tokenManager.isRememberLogin()
      tokenManager.setTokens(newAccessToken, newRefreshToken, remember)

      console.log('[Token Refreshed Successfully]')
      return newAccessToken
    } catch (error) {
      console.error('[Token Refresh Error]', error)
      throw error
    }
  }

  /**
   * 跳转到登录页
   */
  private redirectToLogin(): void {
    tokenManager.clearTokens()
    if (typeof window !== 'undefined') {
      window.location.href = '/login'
    }
  }

  /**
   * GET请求
   * @param url 请求URL
   * @param config 请求配置
   * @returns Promise
   */
  get<T = unknown>(url: string, config?: InternalAxiosRequestConfig): Promise<AxiosResponse<T>> {
    return this.instance.get(url, config)
  }

  /**
   * POST请求
   * @param url 请求URL
   * @param data 请求数据
   * @param config 请求配置
   * @returns Promise
   */
  post<T = unknown>(
    url: string,
    data?: unknown,
    config?: InternalAxiosRequestConfig,
  ): Promise<AxiosResponse<T>> {
    return this.instance.post(url, data, config)
  }

  /**
   * PUT请求
   * @param url 请求URL
   * @param data 请求数据
   * @param config 请求配置
   * @returns Promise
   */
  put<T = unknown>(
    url: string,
    data?: unknown,
    config?: InternalAxiosRequestConfig,
  ): Promise<AxiosResponse<T>> {
    return this.instance.put(url, data, config)
  }

  /**
   * DELETE请求
   * @param url 请求URL
   * @param config 请求配置
   * @returns Promise
   */
  delete<T = unknown>(url: string, config?: InternalAxiosRequestConfig): Promise<AxiosResponse<T>> {
    return this.instance.delete(url, config)
  }

  /**
   * PATCH请求
   * @param url 请求URL
   * @param data 请求数据
   * @param config 请求配置
   * @returns Promise
   */
  patch<T = unknown>(
    url: string,
    data?: unknown,
    config?: InternalAxiosRequestConfig,
  ): Promise<AxiosResponse<T>> {
    return this.instance.patch(url, data, config)
  }

  /**
   * 获取axios实例（用于特殊需求）
   * @returns axios实例
   */
  getInstance(): AxiosInstance {
    return this.instance
  }
}

/** 全局HTTP客户端实例 */
const http = new HttpClient()

export default http
export { tokenManager }
export type { HttpClient }
