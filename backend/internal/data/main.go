package data

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // 下划线表示初始化这个包的内容以便使用
)

// Test 函数名的首字母要大写才是以导出函数
func Test() string {
	return "hello internal world!"
}
