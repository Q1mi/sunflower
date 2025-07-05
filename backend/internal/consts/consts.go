package consts

// CtxKey 上下文 key
type CtxKey string

const (
	JWTAccessTokenSecret    = "夏天夏天悄悄过去留下小秘密" // JWT 访问令牌密钥
	JWTRefreshTokenSecret   = "冬天冬天雪地滑落留下小秘密" // JWT 刷新令牌密钥
	JWTTokenExpireSeconds   = 3600            // JWT 令牌过期时间（秒）
	JWTRefreshExpireSeconds = 7 * 24 * 3600   // JWT 刷新令牌过期时间（秒）

	CtxKeyUserID CtxKey = "userId" // 用户ID上下文 key
)
