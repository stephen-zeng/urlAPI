package database

import (
	"urlAPI/internal/model"
)

var (
	dbPath = "assets/database.db"
	/** @brief 内置提示词名称到旧版配置索引的映射。 */
	PromptMap = map[string]int{
		"laugh":    0,
		"poem":     1,
		"sentence": 2,
	}
	/** @brief 仓库内容内存缓存映射。 */
	RepoMap = make(map[string][]string)
	/** @brief 后台会话内存缓存映射。 */
	SessionMap = make(map[string]model.Session)
)

/** @brief 仓库模型别名。 */
type Repo = model.Repo

/** @brief 会话模型别名。 */
type Session = model.Session

/** @brief 任务模型别名。 */
type Task = model.Task

/** @brief 应用设置模型别名。 */
type AppSetting = model.AppSetting

/** @brief 提供方配置模型别名。 */
type Provider = model.Provider

/** @brief 服务配置模型别名。 */
type ServiceConfig = model.ServiceConfig

/** @brief 提示词模型别名。 */
type Prompt = model.Prompt

/** @brief 配置列表项模型别名。 */
type ConfigListItem = model.ConfigListItem

/** @brief API Key 模型别名。 */
type APIKey = model.APIKey

/** @brief API Key 使用记录模型别名。 */
type APIKeyUsage = model.APIKeyUsage

/** @brief 数据库批量查询结果别名。 */
type DBList = model.DBList
