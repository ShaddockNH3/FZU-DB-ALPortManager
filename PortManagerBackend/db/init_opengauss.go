//go:build opengauss

package db

import (
	"database/sql"
	"log"

	_ "gitee.com/opengauss/openGauss-connector-go-pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Init 用华为官方的 openGauss 驱动拿连接，再交给 gorm。
// 表已存在就不动结构，只补种子。
func Init() {
	sqlDB, err := sql.Open("opengauss", dsn())
	if err != nil {
		log.Fatalf("[DB] 打开连接失败: %v", err)
	}

	g, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		log.Fatalf("[DB] gorm 初始化失败: %v", err)
	}

	setup(g, true)
}
