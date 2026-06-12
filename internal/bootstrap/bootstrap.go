package bootstrap

import (
	"urlAPI/internal/database"
	"urlAPI/internal/op"
)

/**
 * @file bootstrap.go
 * @brief 服务基础依赖生命周期管理。
 * @author 武汉大学开源软件与技术课程 2026
 * @copyright GPL-3.0
 */

/**
 * @brief 初始化服务运行所需的基础依赖。
 * @return error 数据库或业务模块初始化失败时返回错误。
 */
func Init() error {
	if err := database.Init(); err != nil {
		return err
	}
	return op.Init()
}

/**
 * @brief 释放基础依赖占用的资源。
 */
func Release() {
	database.Disconnect()
}
