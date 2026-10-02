package util

import (
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"urlAPI/file"
)

const AppSettingsVersion = 2

const GenerationContextPromptKey = "generation_context"

const SummaryContextPromptKey = "summary_context"

type V2ProviderRow struct {
	Name         string
	APIKeyEnc    string
	TextModel    string
	SummaryModel string
	ImageModel   string
	ImageSize    string
	Endpoint     string
	Enabled      bool
}

type V2ServiceConfigRow struct {
	Service          string
	CacheMinutes     int
	FallbackImageURL string
	Settings         map[string]string
}

type V2PromptRow struct {
	Key      string
	Template string
}

type V2ConfigListItemRow struct {
	Scope     string
	Value     string
	SortOrder int
}

type V2AppSettingRow struct {
	Key   string
	Value string
}

type V2SettingsRows struct {
	Providers       []V2ProviderRow
	ServiceConfigs  []V2ServiceConfigRow
	Prompts         []V2PromptRow
	ConfigListItems []V2ConfigListItemRow
	AppSettings     []V2AppSettingRow
}

type NameValueRow struct {
	Name  string
	Value string
}

const SettingsTableName = "settings"

type ProviderSettings struct {
	OpenAI   ProviderConfig `json:"openai"`
	DeepSeek ProviderConfig `json:"deepseek"`
	Alibaba  ProviderConfig `json:"alibaba"`
	OtherAPI ProviderConfig `json:"otherapi"`
}

type ProviderConfig struct {
	APIKey       string `json:"api_key"`
	TextModel    string `json:"text_model"`
	SummaryModel string `json:"summary_model"`
	ImageModel   string `json:"image_model,omitempty"`
	ImageSize    string `json:"image_size,omitempty"`
	Endpoint     string `json:"endpoint"`
}

type FeatureSettings struct {
	TextEnabled   bool `json:"text_enabled"`
	ImageEnabled  bool `json:"image_enabled"`
	WebImgEnabled bool `json:"web_img_enabled"`
	RandomEnabled bool `json:"random_enabled"`
}

type TextSettings struct {
	GenerationAPI      string   `json:"generation_api"`
	SummaryAPI         string   `json:"summary_api"`
	CacheMinutes       int      `json:"cache_minutes"`
	FallbackImageURL   string   `json:"fallback_image_url"`
	EnabledPromptKeys  []string `json:"enabled_prompt_keys"`
	AcceptedPromptGlob []string `json:"accepted_prompt_glob"`
}

type ImageSettings struct {
	API                string   `json:"api"`
	CacheMinutes       int      `json:"cache_minutes"`
	FallbackImageURL   string   `json:"fallback_image_url"`
	AcceptedPromptGlob []string `json:"accepted_prompt_glob"`
}

type WebSettings struct {
	SummaryAPI       string   `json:"summary_api"`
	CacheMinutes     int      `json:"cache_minutes"`
	FallbackImageURL string   `json:"fallback_image_url"`
	RepoToken        string   `json:"repo_token"`
	YouTubeToken     string   `json:"youtube_token"`
	AllowedHosts     []string `json:"allowed_hosts"`
}

type RandomSettings struct {
	SourceRewriteFrom string `json:"source_rewrite_from"`
	FallbackImageURL  string `json:"fallback_image_url"`
	DefaultAPI        string `json:"default_api"`
}

type SecuritySettings struct {
	DashboardPasswordHash string   `json:"dashboard_password_hash"`
	DashboardAllowedIPs   []string `json:"dashboard_allowed_ips"`
	AllowedReferers       []string `json:"allowed_referers"`
}

type PromptSettings struct {
	GenerationContext string            `json:"generation_context"`
	SummaryContext    string            `json:"summary_context"`
	Templates         map[string]string `json:"templates"`
}

type TaskSettings struct {
	ExceptDomains []string `json:"except_domains"`
	ExceptInfos   []string `json:"except_infos"`
}

type AppSettings struct {
	Version   int              `json:"version"`
	Providers ProviderSettings `json:"providers"`
	Features  FeatureSettings  `json:"features"`
	Text      TextSettings     `json:"text"`
	Image     ImageSettings    `json:"image"`
	Web       WebSettings      `json:"web"`
	Random    RandomSettings   `json:"random"`
	Security  SecuritySettings `json:"security"`
	Prompts   PromptSettings   `json:"prompts"`
	Task      TaskSettings     `json:"task"`
}

func (providers ProviderSettings) ByName(name string) (ProviderConfig, bool) {
	switch name {
	case "openai":
		return providers.OpenAI, true
	case "deepseek":
		return providers.DeepSeek, true
	case "alibaba":
		return providers.Alibaba, true
	case "otherapi":
		return providers.OtherAPI, true
	default:
		return ProviderConfig{}, false
	}
}

func MigrateLegacySettings(legacy map[string][]string) AppSettings {
	settings := DefaultAppSettings()
	legacy = mergeDefaultLegacySettings(legacy)
	settings.Providers = ProviderSettings{
		OpenAI: ProviderConfig{
			APIKey:       legacyString(legacy, "openai", 0, settings.Providers.OpenAI.APIKey),
			TextModel:    legacyString(legacy, "openai", 1, settings.Providers.OpenAI.TextModel),
			SummaryModel: legacyString(legacy, "openai", 2, settings.Providers.OpenAI.SummaryModel),
			ImageModel:   legacyString(legacy, "openai", 3, settings.Providers.OpenAI.ImageModel),
			ImageSize:    legacyString(legacy, "openai", 4, settings.Providers.OpenAI.ImageSize),
			Endpoint:     legacyString(legacy, "openai", 5, settings.Providers.OpenAI.Endpoint),
		},
		DeepSeek: ProviderConfig{
			APIKey:       legacyString(legacy, "deepseek", 0, settings.Providers.DeepSeek.APIKey),
			TextModel:    legacyString(legacy, "deepseek", 1, settings.Providers.DeepSeek.TextModel),
			SummaryModel: legacyString(legacy, "deepseek", 2, settings.Providers.DeepSeek.SummaryModel),
			Endpoint:     legacyString(legacy, "deepseek", 3, settings.Providers.DeepSeek.Endpoint),
		},
		Alibaba: ProviderConfig{
			APIKey:       legacyString(legacy, "alibaba", 0, settings.Providers.Alibaba.APIKey),
			TextModel:    legacyString(legacy, "alibaba", 1, settings.Providers.Alibaba.TextModel),
			SummaryModel: legacyString(legacy, "alibaba", 2, settings.Providers.Alibaba.SummaryModel),
			ImageModel:   legacyString(legacy, "alibaba", 3, settings.Providers.Alibaba.ImageModel),
			ImageSize:    legacyString(legacy, "alibaba", 4, settings.Providers.Alibaba.ImageSize),
			Endpoint:     legacyString(legacy, "alibaba", 5, settings.Providers.Alibaba.Endpoint),
		},
		OtherAPI: ProviderConfig{
			APIKey:       legacyString(legacy, "otherapi", 0, settings.Providers.OtherAPI.APIKey),
			TextModel:    legacyString(legacy, "otherapi", 1, settings.Providers.OtherAPI.TextModel),
			SummaryModel: legacyString(legacy, "otherapi", 2, settings.Providers.OtherAPI.SummaryModel),
			Endpoint:     legacyString(legacy, "otherapi", 3, settings.Providers.OtherAPI.Endpoint),
		},
	}
	settings.Features = FeatureSettings{
		TextEnabled:   legacyBool(legacy, "txt", 0, settings.Features.TextEnabled),
		ImageEnabled:  legacyBool(legacy, "img", 0, settings.Features.ImageEnabled),
		WebImgEnabled: legacyBool(legacy, "web", 1, settings.Features.WebImgEnabled),
		RandomEnabled: legacyBool(legacy, "rand", 0, settings.Features.RandomEnabled),
	}
	settings.Text = TextSettings{
		GenerationAPI:      legacyString(legacy, "txt", 1, settings.Text.GenerationAPI),
		SummaryAPI:         legacyString(legacy, "txt", 2, settings.Text.SummaryAPI),
		CacheMinutes:       legacyInt(legacy, "txt", 3, settings.Text.CacheMinutes),
		FallbackImageURL:   legacyString(legacy, "txt", 4, settings.Text.FallbackImageURL),
		EnabledPromptKeys:  legacyList(legacy, "txtgenenabled"),
		AcceptedPromptGlob: legacyList(legacy, "txtacceptprompt"),
	}
	settings.Image = ImageSettings{
		API:                legacyString(legacy, "img", 1, settings.Image.API),
		CacheMinutes:       legacyInt(legacy, "img", 2, settings.Image.CacheMinutes),
		FallbackImageURL:   legacyString(legacy, "img", 3, settings.Image.FallbackImageURL),
		AcceptedPromptGlob: legacyList(legacy, "imgacceptprompt"),
	}
	settings.Web = WebSettings{
		SummaryAPI:       legacyString(legacy, "web", 2, settings.Web.SummaryAPI),
		CacheMinutes:     legacyInt(legacy, "web", 3, settings.Web.CacheMinutes),
		FallbackImageURL: legacyString(legacy, "web", 4, settings.Web.FallbackImageURL),
		RepoToken:        legacyString(legacy, "web", 5, settings.Web.RepoToken),
		YouTubeToken:     legacyString(legacy, "web", 6, settings.Web.YouTubeToken),
		AllowedHosts:     legacyList(legacy, "webimgallowed"),
	}
	settings.Random = RandomSettings{
		SourceRewriteFrom: legacyString(legacy, "rand", 1, settings.Random.SourceRewriteFrom),
		FallbackImageURL:  legacyString(legacy, "rand", 2, settings.Random.FallbackImageURL),
		DefaultAPI:        legacyString(legacy, "rand", 3, settings.Random.DefaultAPI),
	}
	settings.Security = SecuritySettings{
		DashboardPasswordHash: legacyString(legacy, "dash", 0, settings.Security.DashboardPasswordHash),
		DashboardAllowedIPs:   legacyList(legacy, "dashallowedip"),
		AllowedReferers:       legacyList(legacy, "allowedref"),
	}
	settings.Prompts = PromptSettings{
		GenerationContext: legacyString(legacy, "context", 0, settings.Prompts.GenerationContext),
		SummaryContext:    legacyString(legacy, "context", 1, settings.Prompts.SummaryContext),
		Templates: map[string]string{
			"laugh":    legacyString(legacy, "prompt", 0, settings.Prompts.Templates["laugh"]),
			"poem":     legacyString(legacy, "prompt", 1, settings.Prompts.Templates["poem"]),
			"sentence": legacyString(legacy, "prompt", 2, settings.Prompts.Templates["sentence"]),
		},
	}
	settings.Task = TaskSettings{
		ExceptDomains: legacyList(legacy, "taskexceptdomain"),
		ExceptInfos:   legacyList(legacy, "taskexceptinfo"),
	}
	return NormalizeSettings(settings)
}

func LegacySettingsMap(rows []NameValueRow) map[string][]string {
	legacy := make(map[string][]string, len(rows))
	for _, row := range rows {
		var values []string
		if err := json.Unmarshal([]byte(row.Value), &values); err != nil {
			continue
		}
		legacy[row.Name] = values
	}
	return legacy
}

func BuildV2SettingsRows(settings AppSettings) V2SettingsRows {
	settings = NormalizeSettings(settings)
	return V2SettingsRows{
		Providers: []V2ProviderRow{
			providerRow("openai", settings.Providers.OpenAI),
			providerRow("deepseek", settings.Providers.DeepSeek),
			providerRow("alibaba", settings.Providers.Alibaba),
			providerRow("otherapi", settings.Providers.OtherAPI),
		},
		ServiceConfigs: []V2ServiceConfigRow{
			{
				Service:          "text",
				CacheMinutes:     settings.Text.CacheMinutes,
				FallbackImageURL: settings.Text.FallbackImageURL,
				Settings: map[string]string{
					"generation_api": settings.Text.GenerationAPI,
					"summary_api":    settings.Text.SummaryAPI,
				},
			},
			{
				Service:          "image",
				CacheMinutes:     settings.Image.CacheMinutes,
				FallbackImageURL: settings.Image.FallbackImageURL,
				Settings: map[string]string{
					"api": settings.Image.API,
				},
			},
			{
				Service:          "web",
				CacheMinutes:     settings.Web.CacheMinutes,
				FallbackImageURL: settings.Web.FallbackImageURL,
				Settings: map[string]string{
					"summary_api":   settings.Web.SummaryAPI,
					"repo_token":    settings.Web.RepoToken,
					"youtube_token": settings.Web.YouTubeToken,
				},
			},
			{
				Service:          "random",
				FallbackImageURL: settings.Random.FallbackImageURL,
				Settings: map[string]string{
					"source_rewrite_from": settings.Random.SourceRewriteFrom,
					"default_api":         settings.Random.DefaultAPI,
				},
			},
		},
		Prompts:         buildV2PromptRows(settings.Prompts),
		ConfigListItems: buildV2ConfigListItemRows(settings),
		AppSettings: []V2AppSettingRow{
			{Key: "version", Value: strconv.Itoa(settings.Version)},
			{Key: "features.text_enabled", Value: boolString(settings.Features.TextEnabled)},
			{Key: "features.image_enabled", Value: boolString(settings.Features.ImageEnabled)},
			{Key: "features.web_img_enabled", Value: boolString(settings.Features.WebImgEnabled)},
			{Key: "features.random_enabled", Value: boolString(settings.Features.RandomEnabled)},
			{Key: "security.dashboard_password_hash", Value: settings.Security.DashboardPasswordHash},
		},
	}
}

func providerRow(name string, provider ProviderConfig) V2ProviderRow {
	return V2ProviderRow{
		Name:         name,
		APIKeyEnc:    provider.APIKey,
		TextModel:    provider.TextModel,
		SummaryModel: provider.SummaryModel,
		ImageModel:   provider.ImageModel,
		ImageSize:    provider.ImageSize,
		Endpoint:     provider.Endpoint,
		Enabled:      true,
	}
}

func buildV2PromptRows(prompts PromptSettings) []V2PromptRow {
	rows := []V2PromptRow{
		{Key: GenerationContextPromptKey, Template: prompts.GenerationContext},
		{Key: SummaryContextPromptKey, Template: prompts.SummaryContext},
	}
	for key, template := range prompts.Templates {
		rows = append(rows, V2PromptRow{Key: key, Template: template})
	}
	return rows
}

func buildV2ConfigListItemRows(settings AppSettings) []V2ConfigListItemRow {
	rows := []V2ConfigListItemRow{}
	appendListItems := func(scope string, values []string) {
		for index, value := range values {
			rows = append(rows, V2ConfigListItemRow{Scope: scope, Value: value, SortOrder: index})
		}
	}
	appendListItems("text.enabled_prompt_keys", settings.Text.EnabledPromptKeys)
	appendListItems("text.accepted_prompt_glob", settings.Text.AcceptedPromptGlob)
	appendListItems("image.accepted_prompt_glob", settings.Image.AcceptedPromptGlob)
	appendListItems("web.allowed_hosts", settings.Web.AllowedHosts)
	appendListItems("security.dashboard_allowed_ips", settings.Security.DashboardAllowedIPs)
	appendListItems("security.allowed_referers", settings.Security.AllowedReferers)
	appendListItems("task.except_domains", settings.Task.ExceptDomains)
	appendListItems("task.except_infos", settings.Task.ExceptInfos)
	return rows
}

func AppSettingsFromV2Rows(defaults AppSettings, rows V2SettingsRows) AppSettings {
	settings := defaults
	for _, provider := range rows.Providers {
		config := ProviderConfig{
			APIKey:       provider.APIKeyEnc,
			TextModel:    provider.TextModel,
			SummaryModel: provider.SummaryModel,
			ImageModel:   provider.ImageModel,
			ImageSize:    provider.ImageSize,
			Endpoint:     provider.Endpoint,
		}
		switch provider.Name {
		case "openai":
			settings.Providers.OpenAI = config
		case "deepseek":
			settings.Providers.DeepSeek = config
		case "alibaba":
			settings.Providers.Alibaba = config
		case "otherapi":
			settings.Providers.OtherAPI = config
		}
	}
	for _, service := range rows.ServiceConfigs {
		switch service.Service {
		case "text":
			settings.Text.CacheMinutes = service.CacheMinutes
			settings.Text.FallbackImageURL = service.FallbackImageURL
			settings.Text.GenerationAPI = service.Settings["generation_api"]
			settings.Text.SummaryAPI = service.Settings["summary_api"]
		case "image":
			settings.Image.CacheMinutes = service.CacheMinutes
			settings.Image.FallbackImageURL = service.FallbackImageURL
			settings.Image.API = service.Settings["api"]
		case "web":
			settings.Web.CacheMinutes = service.CacheMinutes
			settings.Web.FallbackImageURL = service.FallbackImageURL
			settings.Web.SummaryAPI = service.Settings["summary_api"]
			settings.Web.RepoToken = service.Settings["repo_token"]
			settings.Web.YouTubeToken = service.Settings["youtube_token"]
		case "random":
			settings.Random.FallbackImageURL = service.FallbackImageURL
			settings.Random.SourceRewriteFrom = service.Settings["source_rewrite_from"]
			settings.Random.DefaultAPI = service.Settings["default_api"]
		}
	}
	settings.Prompts.Templates = make(map[string]string)
	for _, prompt := range rows.Prompts {
		switch prompt.Key {
		case GenerationContextPromptKey:
			settings.Prompts.GenerationContext = prompt.Template
		case SummaryContextPromptKey:
			settings.Prompts.SummaryContext = prompt.Template
		default:
			settings.Prompts.Templates[prompt.Key] = prompt.Template
		}
	}
	listItems := make(map[string][]string)
	for _, item := range rows.ConfigListItems {
		listItems[item.Scope] = append(listItems[item.Scope], item.Value)
	}
	settings.Text.EnabledPromptKeys = listItems["text.enabled_prompt_keys"]
	settings.Text.AcceptedPromptGlob = listItems["text.accepted_prompt_glob"]
	settings.Image.AcceptedPromptGlob = listItems["image.accepted_prompt_glob"]
	settings.Web.AllowedHosts = listItems["web.allowed_hosts"]
	settings.Security.DashboardAllowedIPs = listItems["security.dashboard_allowed_ips"]
	settings.Security.AllowedReferers = listItems["security.allowed_referers"]
	settings.Task.ExceptDomains = listItems["task.except_domains"]
	settings.Task.ExceptInfos = listItems["task.except_infos"]
	for _, value := range rows.AppSettings {
		switch value.Key {
		case "version":
			settings.Version = parseInt(value.Value, settings.Version)
		case "features.text_enabled":
			settings.Features.TextEnabled = parseBool(value.Value, settings.Features.TextEnabled)
		case "features.image_enabled":
			settings.Features.ImageEnabled = parseBool(value.Value, settings.Features.ImageEnabled)
		case "features.web_img_enabled":
			settings.Features.WebImgEnabled = parseBool(value.Value, settings.Features.WebImgEnabled)
		case "features.random_enabled":
			settings.Features.RandomEnabled = parseBool(value.Value, settings.Features.RandomEnabled)
		case "security.dashboard_password_hash":
			settings.Security.DashboardPasswordHash = value.Value
		}
	}
	return NormalizeSettings(settings)
}

func ReadLegacySettings(rows []NameValueRow) map[string][]string {
	return LegacySettingsMap(rows)
}

func BuildAppSettingsFromRows(v2Rows V2SettingsRows, readRows func(table string) []NameValueRow) AppSettings {
	if HasV2Settings(v2Rows) {
		return AppSettingsFromV2Rows(DefaultAppSettings(), v2Rows)
	}
	return MigrateLegacySettings(ReadLegacySettings(readRows(SettingsTableName)))
}

func HasV2Settings(rows V2SettingsRows) bool {
	return len(rows.Providers) > 0 || len(rows.ServiceConfigs) > 0
}

func boolString(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func parseBool(value string, fallback bool) bool {
	switch value {
	case "true", "1":
		return true
	case "false", "0":
		return false
	default:
		return fallback
	}
}

func parseInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func DefaultAppSettings() AppSettings {
	f, err := file.Settings.Open("settings.json")
	if err != nil {
		return fallbackAppSettings()
	}
	defer f.Close()
	d, err := io.ReadAll(f)
	if err != nil {
		return fallbackAppSettings()
	}
	var settings AppSettings
	if err := json.Unmarshal(d, &settings); err != nil {
		return fallbackAppSettings()
	}
	return NormalizeSettings(settings)
}

func mergeDefaultLegacySettings(legacy map[string][]string) map[string][]string {
	if len(legacy) > 0 {
		return legacy
	}
	f, err := file.Settings.Open("setting.json")
	if err != nil {
		return legacy
	}
	defer f.Close()
	d, err := io.ReadAll(f)
	if err != nil {
		return legacy
	}
	var settingsInit struct {
		Names []string   `json:"names"`
		Edits [][]string `json:"edits"`
	}
	if err := json.Unmarshal(d, &settingsInit); err != nil {
		return legacy
	}
	ret := make(map[string][]string, len(settingsInit.Names))
	for index, name := range settingsInit.Names {
		if index >= len(settingsInit.Edits) {
			continue
		}
		ret[name] = settingsInit.Edits[index]
	}
	return ret
}

func fallbackAppSettings() AppSettings {
	return NormalizeSettings(AppSettings{
		Version: AppSettingsVersion,
		Providers: ProviderSettings{
			OpenAI: ProviderConfig{
				TextModel:    "gpt-4o",
				SummaryModel: "gpt-4o-mini",
				ImageModel:   "dall-e-3",
				ImageSize:    "1024x1024",
				Endpoint:     "https://api.openai.com/v1/chat/completions",
			},
			DeepSeek: ProviderConfig{
				TextModel:    "deepseek-chat",
				SummaryModel: "deepseek-chat",
				Endpoint:     "https://api.deepseek.com/chat/completions",
			},
			Alibaba: ProviderConfig{
				TextModel:    "deepseek-v3",
				SummaryModel: "qwen-turbo",
				ImageModel:   "wanx2.1-t2i-turbo",
				ImageSize:    "1024*768",
				Endpoint:     "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
			},
		},
		Text: TextSettings{
			GenerationAPI:    "alibaba",
			SummaryAPI:       "alibaba",
			CacheMinutes:     60,
			FallbackImageURL: "https://raw.githubusercontent.com/stephen-zeng/img/master/fallback.png",
		},
		Image: ImageSettings{
			API:              "alibaba",
			CacheMinutes:     60,
			FallbackImageURL: "https://raw.githubusercontent.com/stephen-zeng/urlAPI/img/master/fallback.png",
		},
		Web: WebSettings{
			SummaryAPI:       "alibaba",
			CacheMinutes:     10,
			FallbackImageURL: "https://raw.githubusercontent.com/stephen-zeng/urlAPI/img/master/fallback.png",
		},
		Random: RandomSettings{
			SourceRewriteFrom: "https://raw.githubusercontent.com",
			FallbackImageURL:  "https://raw.githubusercontent.com/stephen-zeng/urlAPI/master/fallback.png",
			DefaultAPI:        "github",
		},
		Security: SecuritySettings{
			DashboardPasswordHash: "",
			DashboardAllowedIPs:   []string{"*"},
			AllowedReferers:       []string{"localhost"},
		},
		Prompts: PromptSettings{
			GenerationContext: "你是一个助手，需要根据提示词写出对应的语句。语句中不要有空格，不要有换行。不要打招呼，直接给出答案。",
			SummaryContext:    "你是一个助手，需要根据提示词进行总结。总结中不要有空格，不要有换行。不要打招呼，直接给出答案。",
			Templates: map[string]string{
				"laugh":    "讲一个笑话，中文，不要换行，需要句中有标点符号",
				"poem":     "做几句诗歌，中文，不要换行，需要句中有标点符号",
				"sentence": "写几句心灵鸡汤，中文，不要换行，需要句中有标点符号",
			},
		},
		Task: TaskSettings{
			ExceptDomains: []string{"localhost"},
		},
	})
}

func NormalizeSettings(settings AppSettings) AppSettings {
	settings.Version = AppSettingsVersion
	settings.Text.EnabledPromptKeys = normalizeList(settings.Text.EnabledPromptKeys)
	settings.Text.AcceptedPromptGlob = normalizeList(settings.Text.AcceptedPromptGlob)
	settings.Image.AcceptedPromptGlob = normalizeList(settings.Image.AcceptedPromptGlob)
	settings.Web.AllowedHosts = normalizeList(settings.Web.AllowedHosts)
	settings.Security.DashboardAllowedIPs = normalizeList(settings.Security.DashboardAllowedIPs)
	settings.Security.AllowedReferers = normalizeList(settings.Security.AllowedReferers)
	settings.Task.ExceptDomains = normalizeList(settings.Task.ExceptDomains)
	settings.Task.ExceptInfos = normalizeList(settings.Task.ExceptInfos)
	if settings.Prompts.Templates == nil {
		settings.Prompts.Templates = make(map[string]string)
	}
	return settings
}

func legacyString(legacy map[string][]string, key string, index int, fallback string) string {
	values, ok := legacy[key]
	if !ok || index < 0 || index >= len(values) {
		return fallback
	}
	return values[index]
}

func legacyBool(legacy map[string][]string, key string, index int, fallback bool) bool {
	switch legacyString(legacy, key, index, "") {
	case "true":
		return true
	case "false":
		return false
	default:
		return fallback
	}
}

func legacyInt(legacy map[string][]string, key string, index int, fallback int) int {
	value, err := strconv.Atoi(legacyString(legacy, key, index, ""))
	if err != nil {
		return fallback
	}
	return value
}

func legacyList(legacy map[string][]string, key string) []string {
	return normalizeList(legacy[key])
}

func normalizeList(values []string) []string {
	ret := make([]string, 0, len(values))
	seen := make(map[string]bool)
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		ret = append(ret, value)
	}
	return ret
}
