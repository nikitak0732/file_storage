package main

import (
	"log"
	"test_astral/config"
	"test_astral/internal/app"
)

// @title           Backend API
// @version         1.0
// @description     apishka
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath /
// @schemes   http

func main() {
	log.Println("ЗАПУСК ПРИЛОЖЕНИЯ")
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Ошибка считывания окружения: ", err)
	}
	log.Println("Конфиг получили")
	app.Run(cfg)
}
