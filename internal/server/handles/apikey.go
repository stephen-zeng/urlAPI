package handles

import (
	"net/http"
	"strconv"
	"time"
	"urlAPI/internal/database"

	"github.com/gin-gonic/gin"
)

/** @brief API Key 列表响应。 */
type APIKeyListResponse struct {
	Keys []APIKeyResponse `json:"keys"`
}

/** @brief API Key 元数据响应。 */
type APIKeyResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Role        string    `json:"role"`
	Enabled     bool      `json:"enabled"`
	QuotaDay    int       `json:"quota_day"`
	QuotaMonth  int       `json:"quota_month"`
	UsageDay    int       `json:"usage_day"`
	UsageMonth  int       `json:"usage_month"`
	LastUsedAt  time.Time `json:"last_used_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

/** @brief 创建 API Key 的请求参数。 */
type CreateAPIKeyRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Role        string   `json:"role"`
	QuotaDay    int      `json:"quota_day"`
	QuotaMonth  int      `json:"quota_month"`
	AllowedIPs  []string `json:"allowed_ips"`
	ExpiresAt   string   `json:"expires_at"`
}

/** @brief 创建 API Key 后返回的明文密钥响应。 */
type CreateAPIKeyResponse struct {
	APIKey string `json:"api_key"`
}

/** @brief 更新 API Key 的请求参数。 */
type UpdateAPIKeyRequest struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Role        string   `json:"role,omitempty"`
	Enabled     *bool    `json:"enabled,omitempty"`
	QuotaDay    *int     `json:"quota_day,omitempty"`
	QuotaMonth  *int     `json:"quota_month,omitempty"`
	AllowedIPs  []string `json:"allowed_ips,omitempty"`
}

/**
 * @brief 处理 API Key 列表查询请求。
 * @param c Gin 请求上下文。
 */
func ListAPIKeysHandler(c *gin.Context) {
	keys, err := database.APIKeyStore.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []APIKeyResponse
	for _, key := range keys {
		response = append(response, APIKeyResponse{
			ID:          key.ID,
			Name:        key.Name,
			Description: key.Description,
			Role:        key.Role,
			Enabled:     key.Enabled,
			QuotaDay:    key.QuotaDay,
			QuotaMonth:  key.QuotaMonth,
			UsageDay:    key.UsageDay,
			UsageMonth:  key.UsageMonth,
			LastUsedAt:  key.LastUsedAt,
			ExpiresAt:   key.ExpiresAt,
			CreatedAt:   key.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"keys": response})
}

/**
 * @brief 处理 API Key 创建请求。
 * @param c Gin 请求上下文。
 */
func CreateAPIKeyHandler(c *gin.Context) {
	var req CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	var expiresAt time.Time
	if req.ExpiresAt != "" {
		var err error
		expiresAt, err = time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expires_at format, use RFC3339"})
			return
		}
	}

	apiKey, err := database.APIKeyStore.Create(req.Name, req.Description, role, req.QuotaDay, req.QuotaMonth, req.AllowedIPs, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateAPIKeyResponse{APIKey: apiKey})
}

/**
 * @brief 处理 API Key 删除请求。
 * @param c Gin 请求上下文。
 */
func DeleteAPIKeyHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := database.APIKeyStore.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

/**
 * @brief 处理 API Key 更新请求。
 * @param c Gin 请求上下文。
 */
func UpdateAPIKeyHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	updates := make(map[string]any)
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}
	if req.QuotaDay != nil {
		updates["quota_day"] = *req.QuotaDay
	}
	if req.QuotaMonth != nil {
		updates["quota_month"] = *req.QuotaMonth
	}
	if req.AllowedIPs != nil {
		updates["allowed_ips"] = req.AllowedIPs
	}

	if err := database.APIKeyStore.Update(uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}
