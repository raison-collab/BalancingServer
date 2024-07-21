package service

import (
	"BalancingServer/internal/database"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"os/exec"
	"time"
)

// TaskService интерфейс для работы с задачами
type TaskService interface {
	GetPendingTasks() ([]database.TaskModel, error)
	StartTask(taskID uint) error
	CreateTask(task *database.TaskModel) error
	GetTaskStatus(taskID uint) (bool, error)
	GetTaskLogs(taskID uint) (string, error)
}

// taskService реализация TaskService
type taskService struct {
	db *gorm.DB
}

// NewTaskService создает новый экземпляр TaskService
func NewTaskService(db *gorm.DB) TaskService {
	return &taskService{db: db}
}

// GetPendingTasks возвращает список задач, ожидающих выполнения
func (s *taskService) GetPendingTasks() ([]database.TaskModel, error) {
	var tasks []database.TaskModel
	if err := s.db.Where("status = ?", false).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// CreateTask создает новую задачу в базе данных
func (s *taskService) CreateTask(task *database.TaskModel) error {
	return s.db.Create(task).Error
}

// StartTask запускает задачу
func (s *taskService) StartTask(taskID uint) error {
	var task database.TaskModel
	if err := s.db.First(&task, taskID).Error; err != nil {
		return fmt.Errorf("задача не найдена: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // TODO: сделать timeout настраиваемым
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", task.Bash)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("ошибка запуска задачи: %w", err)
	}

	task.PID = cmd.Process.Pid
	if err := s.db.Save(&task).Error; err != nil {
		return fmt.Errorf("ошибка обновления PID задачи: %w", err)
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			if !errors.Is(err, context.DeadlineExceeded) {
				log.Printf("Ошибка выполнения задачи (превыщение таймаута) %d: %v", task.ID, err)
			}
		}

		task.Status = true
		if err := s.db.Save(&task).Error; err != nil {
			log.Printf("Ошибка обновления статуса задачи: %v", err)
		}
	}()

	return nil
}

// GetTaskStatus возвращает статус задачи
func (s *taskService) GetTaskStatus(taskID uint) (bool, error) {
	var task database.TaskModel
	if err := s.db.First(&task, taskID).Error; err != nil {
		return false, fmt.Errorf("задача не найдена: %w", err)
	}
	return task.Status, nil
}

// GetTaskLogs возвращает логи задачи
func (s *taskService) GetTaskLogs(taskID uint) (string, error) {
	var task database.TaskModel
	if err := s.db.First(&task, taskID).Error; err != nil {
		return "", fmt.Errorf("задача не найдена: %w", err)
	}
	return task.Log, nil
}
