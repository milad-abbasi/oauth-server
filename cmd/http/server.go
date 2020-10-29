package main

import (
	"github.com/brpaz/echozap"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	if err := godotenv.Load(".env"); err != nil {
		logger.Warn(err.Error())
	}

	router := echo.New()
	router.Debug = true
	router.Logger.SetLevel(log.DEBUG)

	router.Use(middleware.CORS())
	router.Use(middleware.Recover())
	router.Use(echozap.ZapLogger(logger))

	router.Logger.Fatal(router.Start(":1234"))
}
