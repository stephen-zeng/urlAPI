package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"urlAPI/internal/auth"
	"urlAPI/internal/model"
	"urlAPI/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type appSettingsStore struct {
	mu       sync.RWMutex
	settings util.AppSettings
}

var SettingsStore = appSettingsStore{}

func (store *appSettingsStore) Get() util.AppSettings {
	store.mu.RLock()
	defer store.mu.RUnlock()
	return store.settings
}

func (store *appSettingsStore) Replace(settings util.AppSettings) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.settings = settings
}

func initAppSettings() error {
	settings, err := loadAppSettings()
	if err != nil {
		return err
	}
	announce := ""
	switch hash := settings.Security.DashboardPasswordHash; {
	case hash == "":
		password, generated, err := initialAdminPassword()
		if err != nil {
			return err
		}
		if settings.Security.DashboardPasswordHash, err = auth.HashPassword(auth.ClientCredential(password)); err != nil {
			return err
		}
		if generated {
			announce = password
		} else {
			log.Printf("Dashboard password initialised from %s", auth.EnvAdminPassword)
		}
	case auth.IsLegacyHash(hash) && auth.IsDefaultCredential(hash):
		log.Println("WARNING: the dashboard still uses the default password 123456; change it now")
	}
	if err := SaveAppSettings(settings); err != nil {
		return err
	}
	if announce != "" {
		// Printed once, only when a new installation is initialised.
		fmt.Fprintf(os.Stderr, "\n=== urlAPI initial dashboard password: %s ===\n"+
			"It will not be shown again. Log in at /dash and change it, or run `urlAPI repwd` to generate a new one.\n\n", announce)
	}
	if !secrets.encrypting() && hasStoredSecrets(settings) {
		log.Printf("WARNING: %s is not set; provider API keys and tokens are stored unencrypted", auth.EnvSecretKey)
	}
	return nil
}

func initialAdminPassword() (password string, generated bool, err error) {
	if password = os.Getenv(auth.EnvAdminPassword); password != "" {
		return password, false, nil
	}
	password, err = auth.GeneratePassword()
	return password, true, err
}

func hasStoredSecrets(settings util.AppSettings) bool {
	p := settings.Providers
	return p.OpenAI.APIKey != "" || p.DeepSeek.APIKey != "" || p.Alibaba.APIKey != "" || p.OtherAPI.APIKey != "" ||
		settings.Web.RepoToken != "" || settings.Web.YouTubeToken != ""
}

// settingsWriteMu serialises read-modify-write updates of the settings.
var settingsWriteMu sync.Mutex

// UpdateAppSettings applies update to the current settings and saves the
// result, serialised against other updates.
func UpdateAppSettings(update func(*util.AppSettings) error) error {
	settingsWriteMu.Lock()
	defer settingsWriteMu.Unlock()
	settings := SettingsStore.Get()
	if err := update(&settings); err != nil {
		return err
	}
	return SaveAppSettings(settings)
}

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

func loadAppSettings() (util.AppSettings, error) {
	rows, err := readV2SettingsRows()
	if err != nil {
		return util.AppSettings{}, err
	}
	settings := util.BuildAppSettingsFromRows(rows, readNameValueRows)
	return util.NormalizeSettings(settings), nil
}

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
			Name:         provider.Name,
			APIKeyEnc:    secrets.open(provider.APIKeyEnc, providerSecretContext(provider.Name), decodeLegacyBase64),
			TextModel:    provider.TextModel,
			SummaryModel: provider.SummaryModel,
			ImageModel:   valueString(provider.ImageModel),
			ImageSize:    valueString(provider.ImageSize),
			Endpoint:     provider.Endpoint,
			Enabled:      provider.Enabled,
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
		if service.Service == "web" {
			for _, key := range []string{webRepoTokenKey, webYouTubeTokenKey} {
				values[key] = secrets.open(values[key], webSecretContext(key), identity)
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
		apiKey, err := secrets.seal(row.APIKeyEnc, providerSecretContext(row.Name), encodeLegacyBase64)
		if err != nil {
			return err
		}
		record := Provider{
			Name:         row.Name,
			APIKeyEnc:    apiKey,
			TextModel:    row.TextModel,
			SummaryModel: row.SummaryModel,
			ImageModel:   optionalString(row.ImageModel),
			ImageSize:    optionalString(row.ImageSize),
			Endpoint:     row.Endpoint,
			Enabled:      row.Enabled,
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
		values := row.Settings
		if row.Service == "web" {
			values = make(map[string]string, len(row.Settings))
			for key, value := range row.Settings {
				values[key] = value
			}
			for _, key := range []string{webRepoTokenKey, webYouTubeTokenKey} {
				sealed, err := secrets.seal(values[key], webSecretContext(key), identity)
				if err != nil {
					return err
				}
				values[key] = sealed
			}
		}
		payload, err := json.Marshal(values)
		if err != nil {
			return err
		}
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

func CreateAppSetting(setting *model.AppSetting) error {
	return errors.WithStack(localDB.db.Create(setting).Error)
}

func UpdateAppSetting(setting *model.AppSetting) error {
	return errors.WithStack(localDB.db.Save(setting).Error)
}

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

func DeleteAppSetting(setting *model.AppSetting) error {
	return errors.WithStack(localDB.db.Delete(setting).Error)
}
