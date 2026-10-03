package op

import (
	"encoding/json"
	"net/url"
	"strings"
	"urlAPI/internal/auth"
	"urlAPI/internal/database"
	"urlAPI/util"

	"github.com/pkg/errors"
)

type providerSettingsDTO struct {
	APIKeySet    bool   `json:"api_key_set"`
	APIKey       string `json:"api_key,omitempty"`
	TextModel    string `json:"text_model"`
	SummaryModel string `json:"summary_model"`
	ImageModel   string `json:"image_model,omitempty"`
	ImageSize    string `json:"image_size,omitempty"`
	Endpoint     string `json:"endpoint"`
}

type textSettingsDTO struct {
	Enabled            bool     `json:"enabled"`
	GenerationAPI      string   `json:"generation_api"`
	SummaryAPI         string   `json:"summary_api"`
	CacheMinutes       int      `json:"cache_minutes"`
	FallbackImageURL   string   `json:"fallback_image_url"`
	EnabledPromptKeys  []string `json:"enabled_prompt_keys"`
	AcceptedPromptGlob []string `json:"accepted_prompt_glob"`
}

type imageSettingsDTO struct {
	Enabled            bool     `json:"enabled"`
	API                string   `json:"api"`
	CacheMinutes       int      `json:"cache_minutes"`
	FallbackImageURL   string   `json:"fallback_image_url"`
	AcceptedPromptGlob []string `json:"accepted_prompt_glob"`
}

type webSettingsDTO struct {
	Enabled          bool     `json:"enabled"`
	SummaryAPI       string   `json:"summary_api"`
	CacheMinutes     int      `json:"cache_minutes"`
	FallbackImageURL string   `json:"fallback_image_url"`
	RepoTokenSet     bool     `json:"repo_token_set"`
	RepoToken        string   `json:"repo_token,omitempty"`
	YouTubeTokenSet  bool     `json:"youtube_token_set"`
	YouTubeToken     string   `json:"youtube_token,omitempty"`
	AllowedHosts     []string `json:"allowed_hosts"`
}

type randomSettingsDTO struct {
	Enabled           bool   `json:"enabled"`
	SourceRewriteFrom string `json:"source_rewrite_from"`
	FallbackImageURL  string `json:"fallback_image_url"`
	DefaultAPI        string `json:"default_api"`
}

type dashboardSecurityDTO struct {
	PasswordHash        string   `json:"password_hash,omitempty"`
	DashboardAllowedIPs []string `json:"dashboard_allowed_ips"`
	AllowedReferers     []string `json:"allowed_referers"`
}

type promptSettingsDTO struct {
	GenerationContext string            `json:"generation_context"`
	SummaryContext    string            `json:"summary_context"`
	Templates         map[string]string `json:"templates"`
}

type taskBehaviorDTO struct {
	ExceptDomains []string `json:"except_domains"`
	ExceptInfos   []string `json:"except_infos"`
}

type textPromptSecurityDTO struct {
	AcceptedPromptGlob []string `json:"accepted_prompt_glob"`
}

type imagePromptSecurityDTO struct {
	AcceptedPromptGlob []string `json:"accepted_prompt_glob"`
}

func fetchSettings(info *Session) error {
	body, err := settingsBody(info.SettingPart, database.SettingsStore.Get())
	if err != nil {
		return errors.WithStack(err)
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		return errors.WithStack(err)
	}
	info.SettingBody = encoded
	return nil
}

func editSettings(info *Session) error {
	return database.UpdateAppSettings(func(settings *util.AppSettings) error {
		updated, err := applySettingsBody(info.SettingPart, *settings, info.SettingBody)
		if err != nil {
			return errors.WithStack(err)
		}
		if err := validateSettings(updated); err != nil {
			return errors.WithStack(err)
		}
		*settings = updated
		return nil
	})
}

func settingsBody(part string, settings util.AppSettings) (any, error) {
	switch part {
	case "provider.openai", "openai":
		return providerDTO(settings.Providers.OpenAI), nil
	case "provider.deepseek", "deepseek":
		return providerDTO(settings.Providers.DeepSeek), nil
	case "provider.alibaba", "alibaba":
		return providerDTO(settings.Providers.Alibaba), nil
	case "provider.otherapi", "otherapi":
		return providerDTO(settings.Providers.OtherAPI), nil
	case "feature.text", "txt":
		return textSettingsDTO{
			Enabled:            settings.Features.TextEnabled,
			GenerationAPI:      settings.Text.GenerationAPI,
			SummaryAPI:         settings.Text.SummaryAPI,
			CacheMinutes:       settings.Text.CacheMinutes,
			FallbackImageURL:   settings.Text.FallbackImageURL,
			EnabledPromptKeys:  settings.Text.EnabledPromptKeys,
			AcceptedPromptGlob: settings.Text.AcceptedPromptGlob,
		}, nil
	case "feature.image", "img":
		return imageSettingsDTO{
			Enabled:            settings.Features.ImageEnabled,
			API:                settings.Image.API,
			CacheMinutes:       settings.Image.CacheMinutes,
			FallbackImageURL:   settings.Image.FallbackImageURL,
			AcceptedPromptGlob: settings.Image.AcceptedPromptGlob,
		}, nil
	case "feature.web", "web":
		return webSettingsDTO{
			Enabled:          settings.Features.WebImgEnabled,
			SummaryAPI:       settings.Web.SummaryAPI,
			CacheMinutes:     settings.Web.CacheMinutes,
			FallbackImageURL: settings.Web.FallbackImageURL,
			RepoTokenSet:     settings.Web.RepoToken != "",
			YouTubeTokenSet:  settings.Web.YouTubeToken != "",
			AllowedHosts:     settings.Web.AllowedHosts,
		}, nil
	case "feature.random", "rand":
		return randomSettingsDTO{
			Enabled:           settings.Features.RandomEnabled,
			SourceRewriteFrom: settings.Random.SourceRewriteFrom,
			FallbackImageURL:  settings.Random.FallbackImageURL,
			DefaultAPI:        settings.Random.DefaultAPI,
		}, nil
	case "security.dashboard", "security":
		return dashboardSecurityDTO{
			DashboardAllowedIPs: settings.Security.DashboardAllowedIPs,
			AllowedReferers:     settings.Security.AllowedReferers,
		}, nil
	case "prompt", "contxt":
		return promptSettingsDTO{
			GenerationContext: settings.Prompts.GenerationContext,
			SummaryContext:    settings.Prompts.SummaryContext,
			Templates:         settings.Prompts.Templates,
		}, nil
	case "security.task_behavior", "taskBehavior":
		return taskBehaviorDTO{
			ExceptDomains: settings.Task.ExceptDomains,
			ExceptInfos:   settings.Task.ExceptInfos,
		}, nil
	case "security.text_prompt", "txtSecurity":
		return textPromptSecurityDTO{AcceptedPromptGlob: settings.Text.AcceptedPromptGlob}, nil
	case "security.image_prompt", "imgSecurity":
		return imagePromptSecurityDTO{AcceptedPromptGlob: settings.Image.AcceptedPromptGlob}, nil
	default:
		return nil, errors.New("Setting part not found")
	}
}

func applySettingsBody(part string, settings util.AppSettings, body json.RawMessage) (util.AppSettings, error) {
	switch part {
	case "provider.openai", "openai":
		var dto providerSettingsDTO
		if err := json.Unmarshal(body, &dto); err != nil {
			return settings, err
		}
		settings.Providers.OpenAI = applyProviderDTO(settings.Providers.OpenAI, dto)
	case "provider.deepseek", "deepseek":
		var dto providerSettingsDTO
		if err := json.Unmarshal(body, &dto); err != nil {
			return settings, err
		}
		settings.Providers.DeepSeek = applyProviderDTO(settings.Providers.DeepSeek, dto)
	case "provider.alibaba", "alibaba":
		var dto providerSettingsDTO
		if err := json.Unmarshal(body, &dto); err != nil {
			return settings, err
		}
		settings.Providers.Alibaba = applyProviderDTO(settings.Providers.Alibaba, dto)
	case "provider.otherapi", "otherapi":
		var dto providerSettingsDTO
		if err := json.Unmarshal(body, &dto); err != nil {
			return settings, err
		}
		settings.Providers.OtherAPI = applyProviderDTO(settings.Providers.OtherAPI, dto)
	case "feature.text", "txt":
		var dto textSettingsDTO
		if err := json.Unmarshal(body, &dto); err != nil {
			return settings, err
		}
		settings.Features.TextEnabled = dto.Enabled
		settings.Text.GenerationAPI = dto.GenerationAPI
		settings.Text.SummaryAPI = dto.SummaryAPI
		settings.Text.CacheMinutes = dto.CacheMinutes
		settings.Text.FallbackImageURL = dto.FallbackImageURL
		settings.Text.EnabledPromptKeys = dto.EnabledPromptKeys
		settings.Text.AcceptedPromptGlob = dto.AcceptedPromptGlob
	case "feature.image", "img":
		var dto imageSettingsDTO
		if err := json.Unmarshal(body, &dto); err != nil {
			return settings, err
		}
		settings.Features.ImageEnabled = dto.Enabled
		settings.Image.API = dto.API
		settings.Image.CacheMinutes = dto.CacheMinutes
		settings.Image.FallbackImageURL = dto.FallbackImageURL
		settings.Image.AcceptedPromptGlob = dto.AcceptedPromptGlob
	case "feature.web", "web":
		var dto webSettingsDTO
		if err := json.Unmarshal(body, &dto); err != nil {
			return settings, err
		}
		settings.Features.WebImgEnabled = dto.Enabled
		settings.Web.SummaryAPI = dto.SummaryAPI
		settings.Web.CacheMinutes = dto.CacheMinutes
		settings.Web.FallbackImageURL = dto.FallbackImageURL
		settings.Web.AllowedHosts = dto.AllowedHosts
		if dto.RepoToken != "" {
			settings.Web.RepoToken = dto.RepoToken
		}
		if dto.YouTubeToken != "" {
			settings.Web.YouTubeToken = dto.YouTubeToken
		}
	case "feature.random", "rand":
		var dto randomSettingsDTO
		if err := json.Unmarshal(body, &dto); err != nil {
			return settings, err
		}
		settings.Features.RandomEnabled = dto.Enabled
		settings.Random.SourceRewriteFrom = dto.SourceRewriteFrom
		settings.Random.FallbackImageURL = dto.FallbackImageURL
		settings.Random.DefaultAPI = dto.DefaultAPI
	case "security.dashboard", "security":
		var dto dashboardSecurityDTO
		if err := json.Unmarshal(body, &dto); err != nil {
			return settings, err
		}
		if dto.PasswordHash != "" {
			// password_hash is the client-side SHA-256 of the new password;
			// store only its Argon2id hash.
			hash, err := auth.HashPassword(dto.PasswordHash)
			if err != nil {
				return settings, err
			}
			settings.Security.DashboardPasswordHash = hash
		}
		settings.Security.DashboardAllowedIPs = dto.DashboardAllowedIPs
		settings.Security.AllowedReferers = dto.AllowedReferers
	case "prompt", "contxt":
		var dto promptSettingsDTO
		if err := json.Unmarshal(body, &dto); err != nil {
			return settings, err
		}
		settings.Prompts.GenerationContext = dto.GenerationContext
		settings.Prompts.SummaryContext = dto.SummaryContext
		settings.Prompts.Templates = dto.Templates
	case "security.task_behavior", "taskBehavior":
		var dto taskBehaviorDTO
		if err := json.Unmarshal(body, &dto); err != nil {
			return settings, err
		}
		settings.Task.ExceptDomains = dto.ExceptDomains
		settings.Task.ExceptInfos = dto.ExceptInfos
	case "security.text_prompt", "txtSecurity":
		var dto textPromptSecurityDTO
		if err := json.Unmarshal(body, &dto); err != nil {
			return settings, err
		}
		settings.Text.AcceptedPromptGlob = dto.AcceptedPromptGlob
	case "security.image_prompt", "imgSecurity":
		var dto imagePromptSecurityDTO
		if err := json.Unmarshal(body, &dto); err != nil {
			return settings, err
		}
		settings.Image.AcceptedPromptGlob = dto.AcceptedPromptGlob
	default:
		return settings, errors.New("Setting part not found")
	}
	return util.NormalizeSettings(settings), nil
}

func providerDTO(provider util.ProviderConfig) providerSettingsDTO {
	return providerSettingsDTO{
		APIKeySet:    provider.APIKey != "",
		TextModel:    provider.TextModel,
		SummaryModel: provider.SummaryModel,
		ImageModel:   provider.ImageModel,
		ImageSize:    provider.ImageSize,
		Endpoint:     provider.Endpoint,
	}
}

func applyProviderDTO(provider util.ProviderConfig, dto providerSettingsDTO) util.ProviderConfig {
	if dto.APIKey != "" {
		provider.APIKey = dto.APIKey
	}
	provider.TextModel = dto.TextModel
	provider.SummaryModel = dto.SummaryModel
	provider.ImageModel = dto.ImageModel
	provider.ImageSize = dto.ImageSize
	provider.Endpoint = dto.Endpoint
	return provider
}

func validateSettings(settings util.AppSettings) error {
	if err := validateProviderAPI(settings.Text.GenerationAPI, true); err != nil {
		return err
	}
	if err := validateProviderAPI(settings.Text.SummaryAPI, true); err != nil {
		return err
	}
	if err := validateProviderAPI(settings.Image.API, false); err != nil {
		return err
	}
	if err := validateProviderAPI(settings.Web.SummaryAPI, true); err != nil {
		return err
	}
	if settings.Text.CacheMinutes < 0 || settings.Image.CacheMinutes < 0 || settings.Web.CacheMinutes < 0 {
		return errors.New("cache_minutes must be greater than or equal to 0")
	}
	for _, rawURL := range []string{settings.Text.FallbackImageURL, settings.Image.FallbackImageURL, settings.Web.FallbackImageURL, settings.Random.FallbackImageURL} {
		if err := validateOptionalURL(rawURL); err != nil {
			return err
		}
	}
	for _, endpoint := range []string{settings.Providers.OpenAI.Endpoint, settings.Providers.DeepSeek.Endpoint, settings.Providers.Alibaba.Endpoint, settings.Providers.OtherAPI.Endpoint} {
		if err := validateOptionalURL(endpoint); err != nil {
			return err
		}
	}
	return nil
}

func validateProviderAPI(api string, allowOther bool) error {
	switch api {
	case "openai", "deepseek", "alibaba":
		return nil
	case "otherapi":
		if allowOther {
			return nil
		}
	}
	return errors.New("invalid provider api")
}

func validateOptionalURL(rawURL string) error {
	if strings.TrimSpace(rawURL) == "" {
		return nil
	}
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return errors.New("invalid url")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("invalid url scheme")
	}
	return nil
}
