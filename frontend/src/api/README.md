# API 接口使用说明

本项目提供了完整的API接口封装，符合互联网大厂开发规范，支持用户管理、签到功能和积分系统。

## 🏗️ 架构设计

### 目录结构

```
src/api/
├── constants.ts         # 常量配置（API端点、错误码等）
├── types.ts            # 完整的TypeScript类型定义
├── http.ts             # HTTP客户端（axios封装）
├── account.ts          # 用户账户相关接口
├── checkin.ts          # 签到相关接口
├── points.ts           # 积分相关接口
├── adapter.ts          # 统一适配层（组件适配 + 用户管理）
├── index.ts            # 统一导出文件
├── utils/
│   ├── cache.ts        # 缓存管理工具
│   ├── token-manager.ts # Token管理工具
│   └── error-handler.ts # 错误处理工具
└── README.md           # 本文档
```

### 核心特性

- ✅ **类型安全**: 完整的TypeScript类型定义
- ✅ **统一错误处理**: 标准化的错误处理机制
- ✅ **智能缓存**: 自动缓存管理，提升性能
- ✅ **Token管理**: 自动刷新和存储管理
- ✅ **请求拦截**: 自动添加认证头和日志记录
- ✅ **响应拦截**: 统一数据格式和错误处理
- ✅ **统一适配器**: 将所有后端API适配为前端需要的格式

## 📚 使用指南

### 1. 基础API接口

#### 用户相关

```typescript
import { createUser, login, getUserProfile } from '@/api'

// 用户注册
const result = await createUser({
  username: 'testuser',
  password: 'password123',
  confirmPassword: 'password123',
})

// 用户登录
const loginResult = await login({
  username: 'testuser',
  password: 'password123',
})

// 获取用户信息
const profile = await getUserProfile()
```

#### 签到相关

```typescript
import { getCheckinCalendar, dailyCheckin, retroCheckin } from '@/api'

// 获取签到日历
const calendar = await getCheckinCalendar('2024-12')

// 每日签到
await dailyCheckin()

// 补签
await retroCheckin('2024-12-15')
```

#### 积分相关

```typescript
import { getPointsRecords, getPointsStats } from '@/api'

// 获取积分记录
const records = await getPointsRecords({ limit: 20, offset: 0 })

// 获取积分统计
const stats = await getPointsStats()
```

### 2. 适配器接口（推荐）

适配器接口提供了更友好的前端组件接口，包含缓存和错误处理：

```typescript
import { getPointsInfo, getCheckinCalendarDetail, checkIn, retroCheckIn } from '@/api'

// 获取积分信息（带缓存）
const pointsInfo = await getPointsInfo()

// 获取签到日历详情（带缓存）
const calendarDetail = await getCheckinCalendarDetail(2024, 12)

// 签到（自动清除缓存）
const result = await checkIn()

// 补签（自动清除缓存）
const retroResult = await retroCheckIn('2024-12-15')
```

### 3. 用户管理

```typescript
import { userLogin, userRegister, getCurrentUser, userLogout, isLoggedIn } from '@/api'

// 用户登录（自动处理Token存储）
const loginResult = await userLogin(
  { username: 'test', password: '123456' },
  true, // 记住登录状态
)

// 用户注册（自动登录）
const registerResult = await userRegister({
  username: 'newuser',
  password: 'password123',
  confirmPassword: 'password123',
})

// 获取当前用户
const currentUser = await getCurrentUser()

// 检查登录状态
const loggedIn = isLoggedIn()

// 用户登出
userLogout()
```

## 🛠️ 工具类使用

### 缓存管理

```typescript
import { cacheManager, CacheKeys } from '@/api'

// 设置缓存
cacheManager.set('key', data, 5000) // 5秒过期

// 获取缓存
const cachedData = cacheManager.get('key')

// 清除所有缓存
cacheManager.clear()

// 使用预定义的缓存键
const pointsKey = CacheKeys.pointsInfo()
const calendarKey = CacheKeys.calendar(2024, 12)
```

### Token管理

```typescript
import { tokenManager } from '@/api'

// 获取Token信息
const tokenInfo = tokenManager.getTokenInfo()

// 检查Token状态
const hasValidTokens = tokenManager.hasValidTokens()

// 获取Token统计
const stats = tokenManager.getTokenStats()
```

### 错误处理

```typescript
import { errorHandler } from '@/api'

try {
  // API调用
} catch (error) {
  const standardError = errorHandler.handleUnknownError(error, {
    function: 'myFunction',
    context: 'additional info',
  })
  console.error(standardError.userMessage)
}
```

## ⚙️ 配置说明

### API配置

```typescript
import { API_CONFIG } from '@/api'

// 基础配置
console.log(API_CONFIG.BASE_URL) // http://127.0.0.1:8000/api/v1
console.log(API_CONFIG.TIMEOUT) // 10000
```

### 端点配置

```typescript
import { API_ENDPOINTS } from '@/api'

// 用户相关端点
console.log(API_ENDPOINTS.USER.LOGIN) // /auth/login
console.log(API_ENDPOINTS.USER.PROFILE) // /users/me

// 签到相关端点
console.log(API_ENDPOINTS.CHECKIN.DAILY) // /checkins
```

## 🔧 自定义HTTP请求

如需发送自定义请求，可以直接使用HTTP客户端：

```typescript
import { http } from '@/api'

// 自定义GET请求
const response = await http.get('/custom/endpoint')

// 自定义POST请求
const result = await http.post('/custom/endpoint', { data: 'value' })
```

## 📝 类型定义

所有API接口都有完整的TypeScript类型定义：

```typescript
import type { LoginRequest, LoginResponse, CheckinCalendarResponse, PointsRecord } from '@/api'

// 使用类型定义
const loginData: LoginRequest = {
  username: 'test',
  password: '123456',
}
```

## 🚀 最佳实践

1. **优先使用适配器接口**: 适配器接口提供了更好的缓存和错误处理
2. **合理使用缓存**: 缓存会自动管理，无需手动清理
3. **统一错误处理**: 使用标准化的错误处理机制
4. **类型安全**: 始终使用TypeScript类型定义
5. **日志记录**: API调用会自动记录日志，便于调试

## 🔍 调试技巧

- 查看浏览器控制台的API请求日志
- 使用缓存统计信息监控缓存使用情况
- 检查Token状态排查认证问题
- 利用错误处理器的详细错误信息

## 📈 性能优化

- 自动缓存机制减少重复请求
- Token自动刷新避免频繁登录
- 请求去重防止并发问题
- 智能错误重试机制

## 🔄 架构优势

### 统一适配器设计

将原来的 `adapter.ts` 和 `user-adapter.ts` 合并为单一的 `adapter.ts` 文件，带来以下优势：

1. **减少文件数量**: 简化项目结构，便于维护
2. **统一管理**: 所有适配逻辑集中在一个文件中
3. **共享工具**: 缓存、错误处理等工具可以在不同适配器间共享
4. **一致性**: 统一的代码风格和处理模式
5. **易于扩展**: 新增适配器功能时只需在一个文件中添加
