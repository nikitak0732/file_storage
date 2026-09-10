//go:build swagger

package main

import (
	_ "test_astral/internal/controller/DTO/constant"
	_ "test_astral/internal/controller/DTO/request"
	_ "test_astral/internal/controller/DTO/response"
	_ "test_astral/internal/controller/api"
)

// swag init -g swagger_main.go \
//   -o ./docs \
//   --parseInternal \
//   --parseDependency \
//   --parseDependencyLevel 4 \
//   --dir .,./internal/controller/api,./internal/controller/DTO/request,./internal/controller/DTO/response,./internal/controller/DTO/constant
// @title Backend API
// @version 1.0
// @description apishka!
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8081
// @BasePath /
// @schemes http

func main() {}
