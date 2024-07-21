package api

import (
	"BalancingServer/internal/database"
	"BalancingServer/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreateTaskRequest представляет данные запроса для создания задачи
type CreateTaskRequest struct {
	Bash     string `json:"bash" binding:"required"`
	Ram      uint16 `json:"ram"`
	Disk     uint   `json:"disk"`
	CPU      uint   `json:"cpu"`
	Priority uint8  `json:"priority"`
}

// createTaskHandler создает новую задачу.
// @Summary Создание задачи
// @Description Добавляет новую задачу в очередь.
// @Tags Задачи
// @Accept json
// @Produce json
// @Param task body api.CreateTaskRequest true "Данные задачи"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/tasks [post]
func createTaskHandler(taskService service.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newTask database.TaskModel
		var createTaskRequest CreateTaskRequest

		if err := c.ShouldBindJSON(&createTaskRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
			return
		}

		newTask = database.TaskModel{
			Bash:     createTaskRequest.Bash,
			Ram:      createTaskRequest.Ram,
			Disk:     createTaskRequest.Disk,
			CPU:      createTaskRequest.CPU,
			Priority: createTaskRequest.Priority,
			Status:   false,
			PID:      0,
			Log:      "",
		}

		if err := taskService.CreateTask(&newTask); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании задачи"})
			return
		}

		response := map[string]interface{}{
			"id":   newTask.ID,
			"task": createTaskRequest,
		}

		c.JSON(http.StatusCreated, response)
	}
}

// getTaskStatusHandler возвращает статус задачи
func getTaskStatusHandler(taskService service.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID задачи"})
			return
		}

		status, err := taskService.GetTaskStatus(uint(taskID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении статуса задачи"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": status})
	}
}

// getTaskLogsHandler возвращает логи задачи
func getTaskLogsHandler(taskService service.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID задачи"})
			return
		}

		logs, err := taskService.GetTaskLogs(uint(taskID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении логов задачи"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"logs": logs})
	}
}
