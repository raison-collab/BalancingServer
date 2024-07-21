package main

import (
	"BalancingServer/internal/api"
	"BalancingServer/internal/config"
	"BalancingServer/internal/database"
	"BalancingServer/internal/scheduler"
	"BalancingServer/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// арихитекутурв
//load-balancer/
//├── cmd/
//│   └── main.go                # Точка входа в приложение
//├── internal/
//│   ├── api/                   # Обработчики API
//│   │   ├── router.go          # Маршрутизация API
//│   │   ├── task_handler.go    # Обработчик запросов задач
//│   │   └── middleware/        # Middleware для API (аутентификация, логирование)
//│   ├── config/                # Загрузка конфигурации
//│   │   └── config.go
//│   ├── database/              # Взаимодействие с базой данных
//│   │   ├── models.go          # Модели GORM
//│   │   └── db.go              # Подключение и миграции
//│   ├── scheduler/             # Планировщик задач
//│   │   └── scheduler.go
//│   ├── executor/              # Исполнитель задач
//│   │   └── executor.go
//│   ├── taskqueue/             # Очередь задач (можно использовать in-memory или Redis)
//│   │   └── queue.go
//│   └── service/               # Бизнес-логика
//│       └── task_service.go    # Сервис для работы с задачами
//├── pkg/                      # Вспомогательные пакеты
//│   └── logger/                # Логирование
//└── go.mod                    # Файл модуля Go

func main() {
	loggerSetup()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	database.Migrate(db)

	taskService := service.NewTaskService(db)

	cronLocal := cron.New()
	_, err = cronLocal.AddFunc("@every 1m", func() {
		scheduler.RunScheduler(taskService, cfg)
	})
	if err != nil {
		log.Fatalf("Ошибка добавления задачи планировщика: %v", err)
	}
	cronLocal.Start()

	router := gin.Default()
	api.SetupRoutes(router, taskService)

	go func() {
		if err := router.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)); err != nil {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Остановка приложения...")
	cronLocal.Stop()

	log.Println("Приложение остановлено")
}

func loggerSetup() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Логер настроен")
}
