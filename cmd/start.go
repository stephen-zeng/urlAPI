package cmd

import (
	"urlAPI/internal/bootstrap"
	"urlAPI/internal/server"
)

/**
 * @file start.go
 * @brief HTTP 服务启动命令。
 * @author 武汉大学开源软件与技术课程 2026
 * @copyright GPL-3.0
 */

/**
 * @brief 启动 HTTP 服务。
 * @return error 初始化依赖或启动服务失败时返回错误。
 */
func start() error {
	if err := bootstrap.Init(); err != nil {
		return err
	}
	defer bootstrap.Release()
	return server.Run(Port)
}
