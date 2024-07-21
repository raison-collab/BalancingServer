package api

import (
	"BalancingServer/internal/service"
	"github.com/gin-gonic/gin"
)

// SetupRoutes настраивает маршруты API
func SetupRoutes(router *gin.Engine, taskService service.TaskService) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/tasks", createTaskHandler(taskService))
	}
}
