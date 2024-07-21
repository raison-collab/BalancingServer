package scheduler

import (
	"BalancingServer/internal/config"
	"BalancingServer/internal/service"
	"log"
)

// RunScheduler запускает планировщик задач
func RunScheduler(taskService service.TaskService, cfg config.Config) {
	tasks, err := taskService.GetPendingTasks()
	if err != nil {
		log.Printf("Ошибка получения задач: %v", err)
		return
	}

	availableResources := cfg.ServerResources

	for _, task := range tasks {
		if task.CPU <= availableResources.CPU &&
			task.Ram <= availableResources.RAM &&
			task.Disk <= availableResources.Disk {

			if err := taskService.StartTask(task.ID); err != nil {
				log.Printf("Ошибка запуска задачи: %v", err)
				continue
			}

			availableResources.CPU -= task.CPU
			availableResources.RAM -= task.Ram
			availableResources.Disk -= task.Disk
		}
	}
}
