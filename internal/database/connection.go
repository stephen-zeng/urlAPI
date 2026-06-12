package database

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

/**
 * @file connection.go
 * @brief 数据库连接、迁移和缓存初始化。
 * @author 武汉大学开源软件与技术课程 2026
 * @copyright GPL-3.0
 */

/**
 * @brief 初始化数据库连接及内存缓存。
 * @return error 初始化任一步骤失败时返回错误。
 */
func Init() error {
	figlet := figure.NewFigure("urlAPI", "", true)
	figlet.Print()
	if err := connect(); err != nil {
		return err
	}
	if err := migration(); err != nil {
		return errors.Wrap(err, "migration")
	}
	if err := initRepoMap(); err != nil {
		return err
	}
	if err := initSessionMap(); err != nil {
		return err
	}
	if err := initAppSettings(); err != nil {
		return errors.Wrap(err, "initAppSettings")
	}
	return nil
}

/**
 * @brief 执行数据库表结构迁移。
 * @return error 数据库迁移失败时返回错误。
 */
func migration() error {
	return localDB.db.AutoMigrate(
		&AppSetting{},
		&Provider{},
		&ServiceConfig{},
		&Prompt{},
		&ConfigListItem{},
		&Task{},
		&Session{},
		&Repo{},
		&APIKey{},
		&APIKeyUsage{},
	)
}

/**
 * @brief 建立 SQLite 数据库连接。
 * @return error 连接失败时返回错误。
 */
func connect() error {
	var err error
	os.Mkdir("assets", 0777)
	tmp, _ := sql.Open("sqlite3", dbPath)
	tmp.Close()
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return errors.Wrap(err, "gorm")
	}
	SetLocalDB(db)
	log.Println("Connected to database")
	return nil
}

/**
 * @brief 关闭数据库连接。
 */
func Disconnect() {
	sqlDB, _ := localDB.db.DB()
	defer sqlDB.Close()
	log.Println("Disconnected from database")
}

/**
 * @brief 从数据库初始化仓库缓存映射。
 * @return error 读取或反序列化失败时返回错误。
 */
func initRepoMap() error {
	var repos []Repo
	if err := localDB.db.Find(&repos).Error; err != nil {
		return errors.Wrap(err, "db find")
	}
	for _, repo := range repos {
		var repoList []string
		if err := json.Unmarshal([]byte(repo.Content), &repoList); err != nil {
			return errors.Wrap(err, "json")
		}
		RepoMap[repo.API+";"+repo.Info] = repoList
	}
	log.Println("Initialized RepoMap")
	return nil
}

/**
 * @brief 从数据库初始化会话缓存映射。
 * @return error 读取失败时返回错误。
 */
func initSessionMap() error {
	var sessions []Session
	if err := localDB.db.Find(&sessions).Error; err != nil {
		return errors.Wrap(err, "db")
	}
	for _, session := range sessions {
		SessionMap[session.Token] = session
	}
	log.Println("Initialized SessionMap")
	return nil
}

/**
 * @brief 清空任务表并重新创建结构。
 * @return error 删除或重建任务表失败时返回错误。
 */
func ClearTask() error {
	if localDB.db.Migrator().HasTable(&Task{}) {
		if err := localDB.db.Migrator().DropTable(&Task{}); err != nil {
			return errors.Wrap(err, "db")
		}
		if err := localDB.db.AutoMigrate(&Task{}); err != nil {
			return errors.Wrap(err, "db")
		}
	}
	return nil
}

/**
 * @brief 清空会话表并重置内存会话缓存。
 * @return error 删除或重建会话表失败时返回错误。
 */
func ClearSession() error {
	if localDB.db.Migrator().HasTable(&Session{}) {
		if err := localDB.db.Migrator().DropTable(&Session{}); err != nil {
			return errors.Wrap(err, "db")
		}
		if err := localDB.db.AutoMigrate(&Session{}); err != nil {
			return errors.Wrap(err, "db")
		}
	}
	SessionMap = make(map[string]Session)
	return nil
}
