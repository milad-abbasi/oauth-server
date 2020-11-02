package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(router *echo.Echo) {
	authRouter := router.Group("/auth")

	authRouter.POST("/register", func(c echo.Context) error {
		return c.String(http.StatusOK, "register")
	})

	authRouter.POST("/login", func(c echo.Context) error {
		return c.String(http.StatusOK, "login")
	})
}
