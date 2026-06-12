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

/** @brief Gin 上下文中保存 API Key 信息的键名。 */
const APIKeyContextKey = "api_key"

/** @brief API Key 鉴权模式。 */
type AuthMode int

const (
	/** @brief 关闭 API Key 鉴权。 */
	AuthModeDisabled AuthMode = iota
	/** @brief API Key 可选，存在时校验。 */
	AuthModeOptional
	/** @brief API Key 必须存在且有效。 */
	AuthModeRequired
)

/** @brief API Key 鉴权中间件配置。 */
type AuthConfig struct {
	Mode         AuthMode
	RequireHMAC  bool
	AllowedRoles []string
}

/**
 * @brief 创建 API Key 鉴权中间件。
 * @param config 鉴权模式、HMAC 与角色限制配置。
 * @return gin.HandlerFunc Gin 中间件函数。
 */
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

/**
 * @brief 从请求头或查询参数中提取 API Key。
 * @param c Gin 请求上下文。
 * @return string 提取到的 API Key，未提供时为空字符串。
 */
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

/**
 * @brief 校验请求的 HMAC 签名。
 * @param c Gin 请求上下文。
 * @param apiKey 用于计算签名的 API Key。
 * @return bool 签名有效且时间戳未过期时返回 true。
 */
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

/**
 * @brief 将十进制字符串解析为 int64。
 * @param s 原始字符串。
 * @return int64 解析后的整数。
 * @return error 解析失败时返回错误。
 */
func parseInt64(s string) (int64, error) {
	var result int64
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

/**
 * @brief 从 Gin 上下文中获取已认证的 API Key 信息。
 * @param c Gin 请求上下文。
 * @return *database.APIKey API Key 记录。
 * @return bool 是否存在且类型正确。
 */
func GetAPIKey(c *gin.Context) (*database.APIKey, bool) {
	value, exists := c.Get(APIKeyContextKey)
	if !exists {
		return nil, false
	}
	key, ok := value.(*database.APIKey)
	return key, ok
}

/**
 * @brief 检查当前请求是否已通过 API Key 鉴权。
 * @param c Gin 请求上下文。
 * @return bool 上下文中存在 API Key 信息时返回 true。
 */
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get(APIKeyContextKey)
	return exists
}
