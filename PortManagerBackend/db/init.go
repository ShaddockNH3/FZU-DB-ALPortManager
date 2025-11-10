package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"FZU-DB-ALPortManager/model"
	"FZU-DB-ALPortManager/query"

	"FZU-DB-ALPortManager/pkg/constants"
)

func Init() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		constants.PostgresHost,
		constants.PostgresUser,
		constants.PostgresPassword,
		constants.PostgresDBName,
		constants.PostgresPort,
		constants.PostgresSSLMode,
		constants.PostgresTimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[DB] Failed to connect to database: %v", err)
	}

	query.SetDefault(db)

	err = db.AutoMigrate(&model.ShipInfo{})
	if err != nil {
		log.Fatalf("[DB] Failed to auto migrate: %v", err)
	}

	log.Println("[DB] Database connected and initialized successfully")
}
