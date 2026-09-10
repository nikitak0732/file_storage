package api

import (
	"log"
	"net/http"
	"test_astral/config"
	"test_astral/internal/middleware"
	"test_astral/internal/usecase/service"
	"test_astral/pkg/postgres"
	"time"
)

type Server struct {
	httpServer  *http.Server
	controller  *V1
	config      *config.Config
	connections Connections
}
type Connections struct {
	db *postgres.Postgres
}

func InitServer(srv *service.UseCase, pg *postgres.Postgres, cfg *config.Config, port string) *Server {
	mux := http.NewServeMux()

	r := NewRoutes(srv.File, srv.User, srv.Cache)

	handler := middleware.Chain(
		mux,                       // базовый хендлер
		middleware.WithRecovery(), // восстановление после паники
	)

	r.RegisterRoutes(mux)
	srvData := &Server{
		controller: r,
		config:     cfg,
		httpServer: &http.Server{
			Addr:         ":" + port,
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 45 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		connections: Connections{
			db: pg,
		},
	}
	log.Println("СЕРВАК ГОТОВ ГОТОВ")
	return srvData
}

func (s *Server) Start(cfg *config.Config) error {
	log.Println("Запусаем сервер...")
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Ошибка запуска сервера", err)
	}

	return nil
}
