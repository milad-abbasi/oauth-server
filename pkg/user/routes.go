package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(router *echo.Echo) {
	userRouter := router.Group("/user")

	userRouter.GET("/:id/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hi user")
	})
}
