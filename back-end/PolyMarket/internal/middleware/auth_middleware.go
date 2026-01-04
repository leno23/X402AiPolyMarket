package middleware

import (
	"context"
	"net/http"
	"strings"

	"X402AiPolyMarket/PolyMarket/internal/config"
	"X402AiPolyMarket/PolyMarket/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthMiddleware struct {
	config config.AuthConfig
}

func NewAuthMiddleware(c config.AuthConfig) *AuthMiddleware {
	return &AuthMiddleware{
		config: c,
	}
}

// ContextKey 用于在 context 中存储用户信息的 key
type ContextKey string

const (
	UserIDKey        ContextKey = "user_id"
	WalletAddressKey ContextKey = "wallet_address"
)

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求头获取 Token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.Unauthorized(w, "Missing authorization header")
			return
		}

		// 检查 Bearer 前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(w, "Invalid authorization header format")
			return
		}

		tokenString := parts[1]

		// 解析 Token
		claims, err := utils.ParseToken(tokenString, m.config.AccessSecret)
		if err != nil {
			logx.Errorf("Failed to parse token: %v", err)
			utils.Unauthorized(w, "Invalid or expired token")
			return
		}

		// 将用户信息存入 context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, WalletAddressKey, claims.WalletAddress)

		// 继续处理请求
		next(w, r.WithContext(ctx))
	}
}

// GetUserID 从 context 中获取用户 ID
func GetUserID(ctx context.Context) (uint64, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint64)
	return userID, ok
}

// GetWalletAddress 从 context 中获取钱包地址
func GetWalletAddress(ctx context.Context) (string, bool) {
	address, ok := ctx.Value(WalletAddressKey).(string)
	return address, ok
}

