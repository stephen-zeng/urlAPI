package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"urlAPI/internal/database"

	"github.com/gin-gonic/gin"
)

const APIKeyContextKey = "api_key"

// AuthMode 定义鉴权模式
type AuthMode int

const (
	AuthModeDisabled AuthMode = iota
	AuthModeOptional
	AuthModeRequired
)

// AuthConfig 鉴权配置
type AuthConfig struct {
	Mode       AuthMode
	RequireHMAC bool
	AllowedRoles []string
}

// APIKeyAuthMiddleware 创建 API Key 鉴权中间件
func APIKeyAuthMiddleware(config AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.Mode == AuthModeDisabled {
			c.Next()
			return
		}

		apiKey := extractAPIKey(c)
		if apiKey == "" {
			log.Printf("[Auth] Missing API key from %s | Path: %s", c.ClientIP(), c.Request.URL.Path)
			if config.Mode == AuthModeOptional {
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "Missing API key. Use Authorization: Bearer sk-xxx header",
					"type":    "authentication_error",
					"code":    "missing_api_key",
				},
			})
			c.Abort()
			return
		}

		// HMAC 签名验证
		if config.RequireHMAC {
			if !verifyHMAC(c, apiKey) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": gin.H{
						"message": "Invalid request signature",
						"type":    "authentication_error",
						"code":    "invalid_signature",
					},
				})
				c.Abort()
				return
			}
		}

		// 验证 API Key
		key, err := database.APIKeyStore.ValidateWithIP(apiKey, c.ClientIP())
		if err != nil {
			log.Printf("[Auth] Invalid API key from %s | Error: %v", c.ClientIP(), err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": err.Error(),
					"type":    "authentication_error",
					"code":    "invalid_api_key",
				},
			})
			c.Abort()
			return
		}
		log.Printf("[Auth] API key validated for %s | Key hash: %s", c.ClientIP(), key.KeyHash)

		// 角色检查
		if len(config.AllowedRoles) > 0 {
			allowed := false
			for _, role := range config.AllowedRoles {
				if key.Role == role {
					allowed = true
					break
				}
			}
			if !allowed {
				c.JSON(http.StatusForbidden, gin.H{
					"error": gin.H{
						"message": "Insufficient permissions",
						"type":    "authorization_error",
						"code":    "insufficient_permissions",
					},
				})
				c.Abort()
				return
			}
		}

		// 配额检查
		if err := database.APIKeyStore.CheckQuota(key.KeyHash, key.QuotaDay, key.QuotaMonth); err != nil {
			log.Printf("[Auth] Quota exceeded for key %s from %s | Error: %v", key.KeyHash, c.ClientIP(), err)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"message": err.Error(),
					"type":    "rate_limit_error",
					"code":    "quota_exceeded",
				},
			})
			c.Abort()
			return
		}

		// 将 key 信息存入上下文
		log.Printf("[Auth] Request authorized for %s | Key: %s | Role: %s", c.ClientIP(), key.KeyHash, key.Role)
		c.Set(APIKeyContextKey, key)
		c.Next()
	}
}

// extractAPIKey 从请求中提取 API Key
func extractAPIKey(c *gin.Context) string {
	// 1. Authorization: Bearer sk-xxx
	auth := c.GetHeader("Authorization")
	if auth != "" {
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			return parts[1]
		}
	}

	// 2. X-API-Key header
	if key := c.GetHeader("X-API-Key"); key != "" {
		return key
	}

	// 3. Query parameter ?api_key=xxx
	if key := c.Query("api_key"); key != "" {
		return key
	}

	return ""
}

// verifyHMAC 验证 HMAC 签名
func verifyHMAC(c *gin.Context, apiKey string) bool {
	timestamp := c.GetHeader("X-Timestamp")
	signature := c.GetHeader("X-Signature")

	if timestamp == "" || signature == "" {
		return false
	}

	// 检查时间戳是否在 5 分钟内
	ts, err := parseInt64(timestamp)
	if err != nil {
		return false
	}
	if time.Now().Unix()-ts > 300 {
		return false
	}

	// 重新计算签名
	method := c.Request.Method
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	body, _ := c.GetRawData()

	message := fmt.Sprintf("%s|%s|%s|%s|%s", method, path, query, timestamp, string(body))
	mac := hmac.New(sha256.New, []byte(apiKey))
	mac.Write([]byte(message))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSig))
}

func parseInt64(s string) (int64, error) {
	var result int64
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

// GetAPIKey 从上下文中获取 API Key 信息
func GetAPIKey(c *gin.Context) (*database.APIKey, bool) {
	value, exists := c.Get(APIKeyContextKey)
	if !exists {
		return nil, false
	}
	key, ok := value.(*database.APIKey)
	return key, ok
}

// IsAuthenticated 检查请求是否已通过 API Key 鉴权
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get(APIKeyContextKey)
	return exists
}
