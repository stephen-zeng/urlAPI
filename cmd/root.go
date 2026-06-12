package cmd

import (
	"fmt"
	"os"
)

/**
 * @file root.go
 * @brief 命令行入口与参数分派。
 * @author 武汉大学开源软件与技术课程 2026
 * @copyright GPL-3.0
 */

/** @brief 服务默认监听端口。 */
var Port = "2233"

/**
 * @brief 执行命令行入口。
 *
 * 负责读取进程参数并在失败时输出错误后退出。
 */
func Execute() {
	if err := Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

/**
 * @brief 解析并执行命令行参数。
 * @param args 命令行参数列表，不包含程序名。
 * @return error 当参数非法或执行失败时返回错误。
 */
func Run(args []string) error {
	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "start", "server":
			return start()
		case "port":
			if index+1 >= len(args) {
				return fmt.Errorf("missing port value")
			}
			Port = args[index+1]
			index++
		case "admin":
			return admin(args[index+1:])
		case "repwd", "clear", "logout", "clear_ip_restriction":
			return admin(args[index:])
		default:
			return fmt.Errorf("unknown command %q", args[index])
		}
	}
	return start()
}
