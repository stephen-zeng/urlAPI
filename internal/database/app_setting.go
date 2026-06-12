package database

import (
	"encoding/base64"
	"encoding/json"
	"sync"
	"urlAPI/internal/model"
	"urlAPI/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

/** @brief 线程安全的应用设置缓存。 */
type appSettingsStore struct {
	mu       sync.RWMutex
	settings util.AppSettings
}

/** @brief 全局应用设置缓存入口。 */
var SettingsStore = appSettingsStore{}

/**
 * @brief 获取当前应用设置快照。
 * @return util.AppSettings 当前缓存中的应用设置。
 */
func (store *appSettingsStore) Get() util.AppSettings {
	store.mu.RLock()
	defer store.mu.RUnlock()
	return store.settings
}

/**
 * @brief 替换当前应用设置缓存。
 * @param settings 新的应用设置。
 */
func (store *appSettingsStore) Replace(settings util.AppSettings) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.settings = settings
}

/**
 * @brief 从数据库加载并持久化规范化后的应用设置。
 * @return error 加载或保存失败时返回错误。
 */
func initAppSettings() error {
	settings, err := loadAppSettings()
	if err != nil {
		return err
	}
	return SaveAppSettings(settings)
}

/**
 * @brief 保存完整应用设置并刷新缓存。
 * @param settings 待保存的应用设置。
 * @return error 持久化失败时返回错误。
 */
func SaveAppSettings(settings util.AppSettings) error {
	settings = util.NormalizeSettings(settings)
	if err := localDB.db.Transaction(func(tx *gorm.DB) error {
		rows := util.BuildV2SettingsRows(settings)
		if err := saveProviders(tx, rows.Providers); err != nil {
			return err
		}
		if err := saveServiceConfigs(tx, rows.ServiceConfigs); err != nil {
			return err
		}
		if err := savePrompts(tx, rows.Prompts); err != nil {
			return err
		}
		if err := saveConfigListItems(tx, rows.ConfigListItems); err != nil {
			return err
		}
		if err := saveScalarSettings(tx, rows.AppSettings); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return errors.WithStack(err)
	}
	SettingsStore.Replace(settings)
	return nil
}

/**
 * @brief 加载应用设置，必要时从旧版表迁移。
 * @return util.AppSettings 规范化后的应用设置。
 * @return error 加载失败时返回错误。
 */
func loadAppSettings() (util.AppSettings, error) {
	rows, err := readV2SettingsRows()
	if err != nil {
		return util.AppSettings{}, err
	}
	settings := util.BuildAppSettingsFromRows(rows, readNameValueRows)
	return util.NormalizeSettings(settings), nil
}

/**
 * @brief 读取 V2 设置相关表并组装为行集合。
 * @return util.V2SettingsRows V2 配置行集合。
 * @return error 查询或解析失败时返回错误。
 */
func readV2SettingsRows() (util.V2SettingsRows, error) {
	rows := util.V2SettingsRows{}
	if err := localDB.db.Find(&[]Provider{}).Error; err != nil {
		return rows, errors.WithStack(err)
	}
	var providers []Provider
	if err := localDB.db.Find(&providers).Error; err != nil {
		return rows, errors.WithStack(err)
	}
	for _, provider := range providers {
		rows.Providers = append(rows.Providers, util.V2ProviderRow{
			Name:             provider.Name,
			APIKeyEnc:        decodeSecret(provider.APIKeyEnc),
			TextModel:        provider.TextModel,
			SummaryModel:     provider.SummaryModel,
			ImageModel:       valueString(provider.ImageModel),
			ImageSize:        valueString(provider.ImageSize),
			EmbeddingModel:   provider.EmbeddingModel,
			Endpoint:         provider.Endpoint,
			APIType:          provider.APIType,
			Temperature:      provider.Temperature,
			MaxTokens:        provider.MaxTokens,
			TopP:             provider.TopP,
			PresencePenalty:  provider.PresencePenalty,
			FrequencyPenalty: provider.FrequencyPenalty,
			CustomHeaders:    provider.CustomHeaders,
			Enabled:          provider.Enabled,
		})
	}
	var services []ServiceConfig
	if err := localDB.db.Find(&services).Error; err != nil {
		return rows, errors.WithStack(err)
	}
	for _, service := range services {
		values := map[string]string{}
		if service.Settings != "" {
			if err := json.Unmarshal([]byte(service.Settings), &values); err != nil {
				return rows, errors.WithStack(err)
			}
		}
		rows.ServiceConfigs = append(rows.ServiceConfigs, util.V2ServiceConfigRow{
			Service:          service.Service,
			CacheMinutes:     service.CacheMinutes,
			FallbackImageURL: service.FallbackImageURL,
			Settings:         values,
		})
	}
	var prompts []Prompt
	if err := localDB.db.Find(&prompts).Error; err != nil {
		return rows, errors.WithStack(err)
	}
	for _, prompt := range prompts {
		rows.Prompts = append(rows.Prompts, util.V2PromptRow{Key: prompt.Key, Template: prompt.Template})
	}
	var items []ConfigListItem
	if err := localDB.db.Order("scope ASC, sort_order ASC, id ASC").Find(&items).Error; err != nil {
		return rows, errors.WithStack(err)
	}
	for _, item := range items {
		rows.ConfigListItems = append(rows.ConfigListItems, util.V2ConfigListItemRow{Scope: item.Scope, Value: item.Value, SortOrder: item.SortOrder})
	}
	var values []AppSetting
	if err := localDB.db.Find(&values).Error; err != nil {
		return rows, errors.WithStack(err)
	}
	for _, value := range values {
		rows.AppSettings = append(rows.AppSettings, util.V2AppSettingRow{Key: value.Key, Value: value.Value})
	}
	return rows, nil
}

func readNameValueRows(table string) []util.NameValueRow {
	var rows []struct {
		Name  string
		Value string
	}
	if err := localDB.db.Table(table).Find(&rows).Error; err != nil {
		return nil
	}
	ret := make([]util.NameValueRow, 0, len(rows))
	for _, row := range rows {
		ret = append(ret, util.NameValueRow{Name: row.Name, Value: row.Value})
	}
	return ret
}

func saveProviders(tx *gorm.DB, rows []util.V2ProviderRow) error {
	if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Provider{}).Error; err != nil {
		return err
	}
	for _, row := range rows {
		record := Provider{
			Name:             row.Name,
			APIKeyEnc:        encodeSecret(row.APIKeyEnc),
			TextModel:        row.TextModel,
			SummaryModel:     row.SummaryModel,
			ImageModel:       optionalString(row.ImageModel),
			ImageSize:        optionalString(row.ImageSize),
			EmbeddingModel:   row.EmbeddingModel,
			Endpoint:         row.Endpoint,
			APIType:          row.APIType,
			Temperature:      row.Temperature,
			MaxTokens:        row.MaxTokens,
			TopP:             row.TopP,
			PresencePenalty:  row.PresencePenalty,
			FrequencyPenalty: row.FrequencyPenalty,
			CustomHeaders:    row.CustomHeaders,
			Enabled:          row.Enabled,
		}
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
	}
	return nil
}

func saveServiceConfigs(tx *gorm.DB, rows []util.V2ServiceConfigRow) error {
	if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&ServiceConfig{}).Error; err != nil {
		return err
	}
	for _, row := range rows {
		payload, _ := json.Marshal(row.Settings)
		if err := tx.Create(&ServiceConfig{
			Service:          row.Service,
			CacheMinutes:     row.CacheMinutes,
			FallbackImageURL: row.FallbackImageURL,
			Settings:         string(payload),
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func savePrompts(tx *gorm.DB, rows []util.V2PromptRow) error {
	if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Prompt{}).Error; err != nil {
		return err
	}
	for _, row := range rows {
		if err := tx.Create(&Prompt{Key: row.Key, Template: row.Template}).Error; err != nil {
			return err
		}
	}
	return nil
}

func saveConfigListItems(tx *gorm.DB, rows []util.V2ConfigListItemRow) error {
	if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&ConfigListItem{}).Error; err != nil {
		return err
	}
	for _, row := range rows {
		if err := tx.Create(&ConfigListItem{Scope: row.Scope, Value: row.Value, SortOrder: row.SortOrder}).Error; err != nil {
			return err
		}
	}
	return nil
}

func saveScalarSettings(tx *gorm.DB, rows []util.V2AppSettingRow) error {
	if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&AppSetting{}).Error; err != nil {
		return err
	}
	for _, row := range rows {
		if err := tx.Create(&AppSetting{Key: row.Key, Value: row.Value}).Error; err != nil {
			return err
		}
	}
	return nil
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func valueString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func encodeSecret(value string) string {
	if value == "" {
		return ""
	}
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func decodeSecret(value string) string {
	if value == "" {
		return ""
	}
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return value
	}
	return string(decoded)
}

/** @brief 创建标量应用设置项。 */
func CreateAppSetting(setting *model.AppSetting) error {
	return errors.WithStack(localDB.db.Create(setting).Error)
}

/** @brief 更新标量应用设置项。 */
func UpdateAppSetting(setting *model.AppSetting) error {
	return errors.WithStack(localDB.db.Save(setting).Error)
}

/** @brief 读取标量应用设置项。 */
func ReadAppSetting(setting model.AppSetting) (*model.DBList, error) {
	var settings []model.AppSetting
	err := localDB.db.Where("key = ?", setting.Key).Find(&settings).Error
	if len(settings) == 0 {
		err = errors.WithStack(errors.New("AppSetting not found"))
	}
	ret := model.DBList{
		AppSettingList: settings,
	}
	return &ret, errors.WithStack(err)
}

/** @brief 删除标量应用设置项。 */
func DeleteAppSetting(setting *model.AppSetting) error {
	return errors.WithStack(localDB.db.Delete(setting).Error)
}
