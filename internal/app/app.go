package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"test_astral/config"
	"test_astral/internal/controller/api"
	repo "test_astral/internal/repo/pgsql"
	"test_astral/internal/usecase/service"
	"test_astral/pkg/postgres"
)

type UseCase struct {
	service *service.UseCase
	user    *service.UserService
	file    *service.FileService
}

func Run(cfg *config.Config) {

	pg, err := postgres.New(cfg)
	if err != nil {
		log.Fatal("Ошибка подключение к базе:", err)
	}
	log.Println("К базе подключились")
	uc := initUseCases(pg, cfg)
	log.Println("USECASE ГОТОВ")
	go func() {
		s := api.InitServer(uc.service, pg, cfg, "8080")
		if err := s.Start(cfg); err != nil {

		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	log.Println("ЖДЕМ СИГНАЛ ОСТАНОВКИ...")
	<-quit
	log.Println("Получен сигнал остановки...")
}

func initUseCases(pg *postgres.Postgres, cfg *config.Config) UseCase {
	fileRepo := repo.NewFileRepo(pg)
	userRepo := repo.NewUserRepo(pg)
	cache := service.NewCache(fileRepo)

	return UseCase{
		service: service.New(fileRepo, userRepo, cache, cfg),
	}
}
