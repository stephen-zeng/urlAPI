package data

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // 下划线表示初始化这个包的内容以便使用
)

var (
	dbPath string = "assets/database.db"
	db     *sql.DB
	err    error
)

// Test 函数名的首字母要大写才是以导出函数
func Test() string {
	return "hello internal world!"
}

// Go中的单双引号不一样的

func init() {
	connect()
	dbInit()
}
func Data() {
	//for i := 0; i < 10; i++ {
	//	add(map[string]string{
	//		"time":   time.Now().Format(time.RFC3339), // RFC3339格式较为常见，方便解析
	//		"ip":     "127.0.0.1",
	//		"type":   "txt.generate.laugh",
	//		"status": "done",
	//		"target": strconv.Itoa(i),
	//	})
	//}
	fmt.Println(get(map[string]string{
		"ip": "127.0.0.1",
	}))
}
