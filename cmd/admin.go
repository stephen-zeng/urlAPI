package cmd

import (
	"fmt"
	"log"
	"urlAPI/internal/bootstrap"
	"urlAPI/internal/database"
)

/**
 * @file admin.go
 * @brief 后台管理命令实现。
 * @author 武汉大学开源软件与技术课程 2026
 * @copyright GPL-3.0
 */

/**
 * @brief 执行后台管理命令。
 * @param args 管理子命令及其参数。
 * @return error 命令不存在或执行失败时返回错误。
 */
func admin(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing admin command")
	}
	if err := bootstrap.Init(); err != nil {
		return err
	}
	defer bootstrap.Release()
	switch args[0] {
	case "repwd":
		if err := resetPassword(); err != nil {
			return err
		}
		log.Println("Password has been reset to 123456, please change it ASAP.")
	case "clear":
		if err := database.ClearTask(); err != nil {
			return err
		}
		log.Println("Task Cleared")
	case "logout":
		if err := database.ClearSession(); err != nil {
			return err
		}
		log.Println("Session Restored")
	case "clear_ip_restriction":
		if err := clearIPRestrict(); err != nil {
			return err
		}
		log.Println("Cleared IP restriction")
	default:
		return fmt.Errorf("unknown admin command %q", args[0])
	}
	return nil
}

/**
 * @brief 将后台密码重置为默认值并清空现有会话。
 * @return error 保存设置或清空会话失败时返回错误。
 */
func resetPassword() error {
	settings := database.SettingsStore.Get()
	settings.Security.DashboardPasswordHash = "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92"
	if err := database.SaveAppSettings(settings); err != nil {
		return err
	}
	return database.ClearSession()
}

/**
 * @brief 清空后台 IP 白名单限制。
 * @return error 保存设置失败时返回错误。
 */
func clearIPRestrict() error {
	settings := database.SettingsStore.Get()
	settings.Security.DashboardAllowedIPs = []string{"*"}
	return database.SaveAppSettings(settings)
}
