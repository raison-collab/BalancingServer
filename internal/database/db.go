package database

import (
	"BalancingServer/internal/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB подключается к базе данных
func ConnectDB(cfg config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDatabaseDSN()
	fmt.Println("DSN:", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	return db, nil
}
