package database

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
	"urlAPI/internal/model"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

/** @brief 全局 API Key 存储操作入口。 */
var APIKeyStore apiKeyStore

/** @brief API Key 存储操作集合。 */
type apiKeyStore struct{}

/** @brief 计算 API Key 的 SHA-256 哈希。 */
func hashKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}

/**
 * @brief 创建新的 API Key 记录。
 * @return string 仅创建时返回的明文 API Key。
 * @return error 创建失败时返回错误。
 */
func (s *apiKeyStore) Create(name, description, role string, quotaDay, quotaMonth int, allowedIPs []string, expiresAt time.Time) (string, error) {
	rawKey := "sk-" + uuid.New().String()
	hash := hashKey(rawKey)

	ipsJSON := ""
	if len(allowedIPs) > 0 {
		b, _ := json.Marshal(allowedIPs)
		ipsJSON = string(b)
	}

	key := model.APIKey{
		KeyHash:     hash,
		Name:        name,
		Description: description,
		Role:        role,
		QuotaDay:    quotaDay,
		QuotaMonth:  quotaMonth,
		AllowedIPs:  ipsJSON,
		ExpiresAt:   expiresAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := localDB.db.Create(&key).Error; err != nil {
		return "", errors.Wrap(err, "create api key")
	}
	return rawKey, nil
}

/**
 * @brief 校验 API Key 是否存在、启用且未过期。
 * @param key 明文 API Key。
 * @return *model.APIKey 命中的 API Key 记录。
 * @return error 校验失败时返回错误。
 */
func (s *apiKeyStore) Validate(key string) (*model.APIKey, error) {
	hash := hashKey(key)
	var apiKey model.APIKey
	if err := localDB.db.Where("key_hash = ?", hash).First(&apiKey).Error; err != nil {
		return nil, errors.New("invalid api key")
	}

	if !apiKey.Enabled {
		return nil, errors.New("api key disabled")
	}

	if !apiKey.ExpiresAt.IsZero() && time.Now().After(apiKey.ExpiresAt) {
		return nil, errors.New("api key expired")
	}

	return &apiKey, nil
}

/**
 * @brief 校验 API Key 并检查客户端 IP 白名单。
 * @param key 明文 API Key。
 * @param clientIP 客户端 IP。
 * @return *model.APIKey 命中的 API Key 记录。
 * @return error 校验失败时返回错误。
 */
func (s *apiKeyStore) ValidateWithIP(key string, clientIP string) (*model.APIKey, error) {
	apiKey, err := s.Validate(key)
	if err != nil {
		return nil, err
	}

	if apiKey.AllowedIPs != "" {
		var allowedIPs []string
		if err := json.Unmarshal([]byte(apiKey.AllowedIPs), &allowedIPs); err == nil && len(allowedIPs) > 0 {
			allowed := false
			for _, ip := range allowedIPs {
				if ip == "*" || ip == clientIP {
					allowed = true
					break
				}
			}
			if !allowed {
				return nil, errors.New("ip not allowed")
			}
		}
	}

	return apiKey, nil
}

/**
 * @brief 记录一次 API Key 调用日志。
 * @return error 写入失败时返回错误。
 */
func (s *apiKeyStore) RecordUsage(keyHash, endpoint, ip string, status, latency int) error {
	usage := model.APIKeyUsage{
		KeyHash:   keyHash,
		Endpoint:  endpoint,
		IP:        ip,
		Status:    status,
		Latency:   latency,
		CreatedAt: time.Now(),
	}
	return localDB.db.Create(&usage).Error
}

/**
 * @brief 增加 API Key 的日/月使用计数。
 * @param keyHash API Key 哈希。
 * @return error 更新失败时返回错误。
 */
func (s *apiKeyStore) IncrementUsage(keyHash string) error {
	now := time.Now()
	return localDB.db.Model(&model.APIKey{}).Where("key_hash = ?", keyHash).Updates(map[string]any{
		"usage_day":    gorm.Expr("CASE WHEN DATE(updated_at) = DATE(?) THEN usage_day + 1 ELSE 1 END", now),
		"usage_month":  gorm.Expr("CASE WHEN strftime('%Y-%m', updated_at) = strftime('%Y-%m', ?) THEN usage_month + 1 ELSE 1 END", now),
		"last_used_at": now,
		"updated_at":   now,
	}).Error
}

/**
 * @brief 检查 API Key 是否超过日/月配额。
 * @return error 超额或读取失败时返回错误。
 */
func (s *apiKeyStore) CheckQuota(keyHash string, quotaDay, quotaMonth int) error {
	if quotaDay <= 0 && quotaMonth <= 0 {
		return nil
	}

	var apiKey model.APIKey
	if err := localDB.db.Where("key_hash = ?", keyHash).First(&apiKey).Error; err != nil {
		return err
	}

	now := time.Now()
	if quotaDay > 0 {
		if now.Day() != apiKey.UpdatedAt.Day() || now.Month() != apiKey.UpdatedAt.Month() {
			apiKey.UsageDay = 0
		}
		if apiKey.UsageDay >= quotaDay {
			return errors.New("daily quota exceeded")
		}
	}

	if quotaMonth > 0 {
		if now.Month() != apiKey.UpdatedAt.Month() || now.Year() != apiKey.UpdatedAt.Year() {
			apiKey.UsageMonth = 0
		}
		if apiKey.UsageMonth >= quotaMonth {
			return errors.New("monthly quota exceeded")
		}
	}

	return nil
}

/**
 * @brief 按创建时间倒序列出 API Key。
 * @return []model.APIKey API Key 列表。
 * @return error 查询失败时返回错误。
 */
func (s *apiKeyStore) List() ([]model.APIKey, error) {
	var keys []model.APIKey
	err := localDB.db.Order("created_at DESC").Find(&keys).Error
	return keys, err
}

/**
 * @brief 删除指定 ID 的 API Key。
 * @param id API Key 主键。
 * @return error 删除失败时返回错误。
 */
func (s *apiKeyStore) Delete(id uint) error {
	return localDB.db.Delete(&model.APIKey{}, id).Error
}

/**
 * @brief 更新指定 ID 的 API Key 字段。
 * @param id API Key 主键。
 * @param updates 待更新字段映射。
 * @return error 更新失败时返回错误。
 */
func (s *apiKeyStore) Update(id uint, updates map[string]any) error {
	if value, ok := updates["expires_at"]; ok {
		switch typed := value.(type) {
		case string:
			if typed == "" {
				updates["expires_at"] = time.Time{}
			} else if parsed, err := time.Parse(time.RFC3339, typed); err == nil {
				updates["expires_at"] = parsed
			} else if parsed, err := time.Parse("2006-01-02T15:04", typed); err == nil {
				updates["expires_at"] = parsed
			}
		case nil:
			updates["expires_at"] = time.Time{}
		}
	}
	updates["updated_at"] = time.Now()
	return localDB.db.Model(&model.APIKey{}).Where("id = ?", id).Updates(updates).Error
}

/**
 * @brief 统计指定时间之后的 API Key 使用次数。
 * @return int64 使用记录数量。
 * @return error 查询失败时返回错误。
 */
func (s *apiKeyStore) GetUsageStats(keyHash string, since time.Time) (int64, error) {
	var count int64
	err := localDB.db.Model(&model.APIKeyUsage{}).Where("key_hash = ? AND created_at >= ?", keyHash, since).Count(&count).Error
	return count, err
}

/**
 * @brief 重置非当天更新的日使用量。
 * @return error 更新失败时返回错误。
 */
func (s *apiKeyStore) ResetDailyUsage() error {
	now := time.Now()
	return localDB.db.Model(&model.APIKey{}).Where("DATE(updated_at) != DATE(?)", now).Update("usage_day", 0).Error
}

/**
 * @brief 重置非当月更新的月使用量。
 * @return error 更新失败时返回错误。
 */
func (s *apiKeyStore) ResetMonthlyUsage() error {
	now := time.Now()
	return localDB.db.Model(&model.APIKey{}).Where("strftime('%Y-%m', updated_at) != strftime('%Y-%m', ?)", now).Update("usage_month", 0).Error
}
