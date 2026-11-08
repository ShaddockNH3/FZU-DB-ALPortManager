//go:build postgres

package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Init 用 gorm 官方驱动连库。建表与种子在 schema.go，和 openGauss 共用。
func Init() {
	g, err := gorm.Open(postgres.Open(dsn()), &gorm.Config{})
	if err != nil {
		log.Fatalf("[DB] 连接失败: %v", err)
	}

	setup(g, false)
}
