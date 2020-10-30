package main

import (
	"fmt"

	"github.com/brpaz/echozap"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/milad-abbasi/oauth-server/pkg/config"
	userController "github.com/milad-abbasi/oauth-server/pkg/user/controller"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Error(err.Error())
		}
	}()

	if err := godotenv.Load(".env"); err != nil {
		logger.Warn(err.Error())
	}

	router := echo.New()
	router.Debug = true
	router.Logger.SetLevel(log.DEBUG)

	router.Use(middleware.CORS())
	router.Use(middleware.Recover())
	router.Use(echozap.ZapLogger(logger))
	userController.RegisterRoutes(router)

	router.Logger.Fatal(router.Start(fmt.Sprintf("0.0.0.0:%s", config.GetWithDefault("HTTP_PORT", "1234"))))
}
