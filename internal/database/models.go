package database

import (
	"gorm.io/gorm"
)

// TaskModel модель задачи в базе данных
type TaskModel struct {
	gorm.Model        // Добавляет ID, CreatedAt, UpdatedAt, DeletedAt
	Bash       string `gorm:"type:text"`
	Ram        uint16
	Disk       uint
	CPU        uint
	Priority   uint8
	Status     bool
	PID        int
	Log        string `gorm:"type:text"`
}
