package main

import (
	"context"
	"fmt"

	"github.com/brpaz/echozap"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/milad-abbasi/oauth-server/pkg/auth"
	"github.com/milad-abbasi/oauth-server/pkg/common"
	"github.com/milad-abbasi/oauth-server/pkg/user"
	userPgRepo "github.com/milad-abbasi/oauth-server/pkg/user/repository/postgres"
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

	pgConfig, err := pgxpool.ParseConfig(common.MustGetEnv("POSTGRES_URI"))
	if err != nil {
		logger.Fatal("Failed to parse postgres uri", zap.Error(err))
	}
	pgConfig.ConnConfig.Logger = zapadapter.NewLogger(logger)
	pgConfig.ConnConfig.LogLevel = pgx.LogLevelTrace

	pgPool, err := pgxpool.ConnectConfig(context.Background(), pgConfig)
	if err != nil {
		logger.Fatal("Failed to establish a connection to postgres", zap.Error(err))
	}
	defer pgPool.Close()

	userRepo := userPgRepo.New(logger, pgPool)
	userService := user.NewService(logger, userRepo)

	structValidator := validator.New()

	router := echo.New()
	router.Debug = true
	router.Logger.SetLevel(log.DEBUG)
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(
		middleware.CORS(),
		middleware.Recover(),
		echozap.ZapLogger(logger),
	)
	auth.RegisterRoutes(router, structValidator, userService)
	user.RegisterRoutes(router)

	router.Logger.Fatal(router.Start(fmt.Sprintf("0.0.0.0:%s", common.GetEnvWithDefault("HTTP_PORT", "1234"))))
}
