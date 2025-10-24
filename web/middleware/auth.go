package middleware

import (
	"go-tpl/infra/jwt"
	"go-tpl/logic/shared"
	"go-tpl/web/base"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	UserIDKey           = "user_id"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 中获取 token
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			base.FailWithError(c, shared.ErrNoToken)
			return
		}

		// 检查 Bearer 前缀
		if !strings.HasPrefix(authHeader, BearerPrefix) {
			base.FailWithError(c, shared.ErrNoToken)
			return
		}

		// 提取 token
		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		if tokenString == "" {
			base.FailWithError(c, shared.ErrNoToken)
			return
		}

		// 解析 token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			base.FailWithError(c, shared.ErrInvalidToken)
			return
		}

		// 将用户信息存入上下文
		c.Set(UserIDKey, claims.UserID)

		c.Next()
	}
}

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) (uint, bool) {
	value, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(uint)
	return userID, ok
}
