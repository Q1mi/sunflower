/**
 * 错误处理工具
 * @description 提供统一的错误处理机制，包括错误分类、日志记录和用户提示
 */

import { HTTP_STATUS, BUSINESS_ERROR_CODES, ERROR_MESSAGES } from '../constants'
import type { AxiosError } from 'axios'

/** 错误类型枚举 */
export enum ErrorType {
  /** 网络错误 */
  NETWORK = 'NETWORK',
  /** 业务错误 */
  BUSINESS = 'BUSINESS',
  /** 认证错误 */
  AUTH = 'AUTH',
  /** 权限错误 */
  PERMISSION = 'PERMISSION',
  /** 验证错误 */
  VALIDATION = 'VALIDATION',
  /** 服务器错误 */
  SERVER = 'SERVER',
  /** 未知错误 */
  UNKNOWN = 'UNKNOWN',
}

/** 错误级别枚举 */
export enum ErrorLevel {
  /** 信息 */
  INFO = 'INFO',
  /** 警告 */
  WARN = 'WARN',
  /** 错误 */
  ERROR = 'ERROR',
  /** 致命错误 */
  FATAL = 'FATAL',
}

/** 标准化错误接口 */
export interface StandardError {
  /** 错误类型 */
  type: ErrorType
  /** 错误级别 */
  level: ErrorLevel
  /** 错误码 */
  code: number | string
  /** 错误消息 */
  message: string
  /** 用户友好的错误消息 */
  userMessage: string
  /** 原始错误对象 */
  originalError?: unknown
  /** 错误发生时间 */
  timestamp: number
  /** 错误上下文信息 */
  context?: Record<string, unknown>
}

/** 错误处理器类 */
class ErrorHandler {
  /** 错误日志记录器 */
  private logError(error: StandardError): void {
    const logData = {
      type: error.type,
      level: error.level,
      code: error.code,
      message: error.message,
      timestamp: new Date(error.timestamp).toISOString(),
      context: error.context,
    }

    switch (error.level) {
      case ErrorLevel.INFO:
        console.info('[API Error]', logData)
        break
      case ErrorLevel.WARN:
        console.warn('[API Error]', logData)
        break
      case ErrorLevel.ERROR:
        console.error('[API Error]', logData)
        break
      case ErrorLevel.FATAL:
        console.error('[API Fatal Error]', logData)
        break
    }
  }

  /**
   * 处理HTTP错误
   * @param error Axios错误对象
   * @returns 标准化错误
   */
  handleHttpError(error: AxiosError): StandardError {
    const timestamp = Date.now()

    // 网络错误
    if (!error.response) {
      const standardError: StandardError = {
        type: ErrorType.NETWORK,
        level: ErrorLevel.ERROR,
        code: 'NETWORK_ERROR',
        message: error.message || 'Network request failed',
        userMessage: ERROR_MESSAGES.NETWORK_ERROR,
        originalError: error,
        timestamp,
        context: {
          url: error.config?.url,
          method: error.config?.method,
        },
      }
      this.logError(standardError)
      return standardError
    }

    const { status, data } = error.response
    const context = {
      url: error.config?.url,
      method: error.config?.method,
      status,
      responseData: data,
    }

    // 根据HTTP状态码分类错误
    switch (status) {
      case HTTP_STATUS.UNAUTHORIZED:
        return this.createStandardError({
          type: ErrorType.AUTH,
          level: ErrorLevel.WARN,
          code: status,
          message: (data as { message?: string })?.message || 'Unauthorized',
          userMessage: ERROR_MESSAGES.TOKEN_EXPIRED,
          originalError: error,
          timestamp,
          context,
        })

      case HTTP_STATUS.FORBIDDEN:
        return this.createStandardError({
          type: ErrorType.PERMISSION,
          level: ErrorLevel.WARN,
          code: status,
          message: (data as { message?: string })?.message || 'Forbidden',
          userMessage: ERROR_MESSAGES.PERMISSION_DENIED,
          originalError: error,
          timestamp,
          context,
        })

      case HTTP_STATUS.NOT_FOUND:
        return this.createStandardError({
          type: ErrorType.BUSINESS,
          level: ErrorLevel.WARN,
          code: status,
          message: (data as { message?: string })?.message || 'Not Found',
          userMessage: '请求的资源不存在',
          originalError: error,
          timestamp,
          context,
        })

      case HTTP_STATUS.INTERNAL_SERVER_ERROR:
        return this.createStandardError({
          type: ErrorType.SERVER,
          level: ErrorLevel.ERROR,
          code: status,
          message: (data as { message?: string })?.message || 'Internal Server Error',
          userMessage: ERROR_MESSAGES.SERVER_ERROR,
          originalError: error,
          timestamp,
          context,
        })

      default:
        return this.createStandardError({
          type: ErrorType.UNKNOWN,
          level: ErrorLevel.ERROR,
          code: status,
          message: (data as { message?: string })?.message || `HTTP ${status} Error`,
          userMessage: ERROR_MESSAGES.UNKNOWN_ERROR,
          originalError: error,
          timestamp,
          context,
        })
    }
  }

  /**
   * 处理业务错误
   * @param code 业务错误码
   * @param message 错误消息
   * @param context 错误上下文
   * @returns 标准化错误
   */
  handleBusinessError(
    code: number | string,
    message: string,
    context?: Record<string, unknown>,
  ): StandardError {
    const timestamp = Date.now()

    // 根据业务错误码分类
    let type = ErrorType.BUSINESS
    const level = ErrorLevel.WARN
    let userMessage = message

    if (code === BUSINESS_ERROR_CODES.VALIDATION_ERROR) {
      type = ErrorType.VALIDATION
      userMessage = ERROR_MESSAGES.VALIDATION_ERROR
    } else if (code === BUSINESS_ERROR_CODES.UNAUTHORIZED) {
      type = ErrorType.AUTH
      userMessage = ERROR_MESSAGES.TOKEN_EXPIRED
    } else if (code === BUSINESS_ERROR_CODES.FORBIDDEN) {
      type = ErrorType.PERMISSION
      userMessage = ERROR_MESSAGES.PERMISSION_DENIED
    }

    return this.createStandardError({
      type,
      level,
      code,
      message,
      userMessage,
      timestamp,
      context,
    })
  }

  /**
   * 创建标准化错误
   * @param errorData 错误数据
   * @returns 标准化错误
   */
  private createStandardError(errorData: StandardError): StandardError {
    this.logError(errorData)
    return errorData
  }

  /**
   * 处理未知错误
   * @param error 原始错误
   * @param context 错误上下文
   * @returns 标准化错误
   */
  handleUnknownError(error: unknown, context?: Record<string, unknown>): StandardError {
    const timestamp = Date.now()
    const message = error instanceof Error ? error.message : String(error)

    return this.createStandardError({
      type: ErrorType.UNKNOWN,
      level: ErrorLevel.ERROR,
      code: 'UNKNOWN_ERROR',
      message,
      userMessage: ERROR_MESSAGES.UNKNOWN_ERROR,
      originalError: error,
      timestamp,
      context,
    })
  }
}

/** 全局错误处理器实例 */
export const errorHandler = new ErrorHandler()

/** 错误处理装饰器 */
export function handleApiError<T extends (...args: unknown[]) => Promise<unknown>>(target: T): T {
  return (async (...args: Parameters<T>) => {
    try {
      return await target(...args)
    } catch (error) {
      const standardError = errorHandler.handleHttpError(error as AxiosError)
      throw new Error(standardError.userMessage)
    }
  }) as T
}

/** 导出类型 */
export type { ErrorHandler }
