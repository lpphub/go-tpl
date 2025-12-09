package jwt

import (
	"go-tpl/infra"
	"go-tpl/infra/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWT(t *testing.T) {
	// 设置测试配置
	infra.Cfg = &config.Config{
		JWT: config.JWTConfig{
			Secret:            "test-secret-key",
			ExpireTime:        3600,      // 1 hour
			RefreshExpireTime: 7 * 86400, // 7 days
		},
	}

	t.Run("GenerateTokenPair", func(t *testing.T) {
		userID := uint(123)

		tokenPair, err := GenerateTokenPair(userID)
		require.NoError(t, err)
		assert.NotEmpty(t, tokenPair.AccessToken)
		assert.NotEmpty(t, tokenPair.RefreshToken)
		assert.NotEqual(t, tokenPair.AccessToken, tokenPair.RefreshToken)
	})

	t.Run("ParseAccessToken", func(t *testing.T) {
		userID := uint(456)

		tokenPair, err := GenerateTokenPair(userID)
		require.NoError(t, err)

		// 解析 access token
		claims, err := ParseToken(tokenPair.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, AccessTokenType, claims.Type)
	})

	t.Run("ParseRefreshToken", func(t *testing.T) {
		userID := uint(789)

		tokenPair, err := GenerateTokenPair(userID)
		require.NoError(t, err)

		// 解析 refresh token
		claims, err := ParseToken(tokenPair.RefreshToken)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, RefreshTokenType, claims.Type)
	})

	t.Run("RefreshToken", func(t *testing.T) {
		userID := uint(999)

		// 生成初始 token 对
		tokenPair, err := GenerateTokenPair(userID)
		require.NoError(t, err)

		// 使用 refresh token 生成新的 token 对
		newTokenPair, err := RefreshToken(tokenPair.RefreshToken)
		require.NoError(t, err)
		assert.NotEmpty(t, newTokenPair.AccessToken)
		assert.NotEmpty(t, newTokenPair.RefreshToken)

		// 由于 JWT 的 iat (issued at) 是基于当前时间的，如果生成得太快可能相同
		// 我们主要验证功能是否正常工作
		// 解析新 token 确保它们有效
		newAccessClaims, err := ParseToken(newTokenPair.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, userID, newAccessClaims.UserID)
		assert.Equal(t, AccessTokenType, newAccessClaims.Type)

		newRefreshClaims, err := ParseToken(newTokenPair.RefreshToken)
		require.NoError(t, err)
		assert.Equal(t, userID, newRefreshClaims.UserID)
		assert.Equal(t, RefreshTokenType, newRefreshClaims.Type)

		// 确保新旧 token 的用户 ID 一致
		oldAccessClaims, err := ParseToken(tokenPair.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, userID, oldAccessClaims.UserID)
	})

	t.Run("RefreshTokenWithAccessToken", func(t *testing.T) {
		userID := uint(111)

		tokenPair, err := GenerateTokenPair(userID)
		require.NoError(t, err)

		// 尝试使用 access token 刷新（应该失败）
		_, err = RefreshToken(tokenPair.AccessToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token type: expected refresh token")
	})

	t.Run("TokenExpiration", func(t *testing.T) {
		// 创建一个已过期的配置
		infra.Cfg.JWT.ExpireTime = 1
		infra.Cfg.JWT.RefreshExpireTime = 1

		userID := uint(222)

		tokenPair, err := GenerateTokenPair(userID)
		require.NoError(t, err)

		// 等待 token 过期
		time.Sleep(2 * time.Second)

		// 尝试解析过期的 token
		_, err = ParseToken(tokenPair.AccessToken)
		assert.Error(t, err)

		_, err = ParseToken(tokenPair.RefreshToken)
		assert.Error(t, err)
	})
}

func TestGenerateToken(t *testing.T) {
	// 设置测试配置
	infra.Cfg = &config.Config{
		JWT: config.JWTConfig{
			Secret:            "test-secret-key",
			ExpireTime:        3600,
			RefreshExpireTime: 7 * 86400,
		},
	}

	t.Run("BackwardCompatibility", func(t *testing.T) {
		userID := uint(333)

		// 测试旧的 GenerateToken 函数仍然工作
		token, err := GenerateToken(userID)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// 解析 token
		claims, err := ParseToken(token)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, AccessTokenType, claims.Type)
	})

	t.Run("DefaultValues", func(t *testing.T) {
		// 使用默认值
		infra.Cfg.JWT.ExpireTime = 0
		infra.Cfg.JWT.RefreshExpireTime = 0

		userID := uint(444)

		tokenPair, err := GenerateTokenPair(userID)
		require.NoError(t, err)

		// 验证 token 可以解析
		accessClaims, err := ParseToken(tokenPair.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, userID, accessClaims.UserID)

		refreshClaims, err := ParseToken(tokenPair.RefreshToken)
		require.NoError(t, err)
		assert.Equal(t, userID, refreshClaims.UserID)
	})
}
