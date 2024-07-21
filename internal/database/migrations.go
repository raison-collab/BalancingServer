package database

import (
	"gorm.io/gorm"
	"log"
)

func migrateTasks(db *gorm.DB) error {
	log.Println("TaskModel")
	return db.AutoMigrate(TaskModel{})
}

func Migrate(db *gorm.DB) {
	log.Println("Запуск миграций...")

	if err := migrateTasks(db); err != nil {
		log.Fatalf("Ошибка миграции TaskModel: %v", err)
	}

	log.Println("Миграции завершены")
}
