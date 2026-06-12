package main

import "urlAPI/cmd"

/**
 * @file main.go
 * @brief urlAPI 程序入口。
 * @author 武汉大学开源软件与技术课程 2026
 * @copyright GPL-3.0
 */

/**
 * @brief 程序主入口。
 *
 * 启动命令行参数解析并分派具体子命令。
 */
func main() {
	cmd.Execute()
}
