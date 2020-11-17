package main

import (
	"context"
	"fmt"

	"github.com/brpaz/echozap"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/milad-abbasi/oauth-server/pkg/auth"
	"github.com/milad-abbasi/oauth-server/pkg/common/config"
	"github.com/milad-abbasi/oauth-server/pkg/common/validator"
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

	pgConfig, err := pgxpool.ParseConfig(config.MustGetEnv("POSTGRES_URI"))
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

	router := echo.New()
	router.Debug = true
	router.Logger.SetLevel(log.DEBUG)
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(
		middleware.CORS(),
		middleware.Recover(),
		echozap.ZapLogger(logger),
	)

	userRepo := userPgRepo.New(logger, pgPool)
	userService := user.NewService(logger, userRepo)
	authService := auth.NewService(logger, userService)

	structValidator := validator.NewStructValidator()
	auth.NewController(logger, router, structValidator, authService).RegisterRoutes()
	user.NewController(logger, router, structValidator, userService).RegisterRoutes()

	router.Logger.Fatal(router.Start(fmt.Sprintf("0.0.0.0:%s", config.GetEnvWithDefault("HTTP_PORT", "1234"))))
}
