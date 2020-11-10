package auth

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/milad-abbasi/oauth-server/pkg/user"
)

func RegisterRoutes(router *echo.Echo, sv *validator.Validate, userService *user.Service) {
	authRouter := router.Group("/auth")

	authRouter.POST("/register", func(c echo.Context) error {
		var dto RegisterDto
		if err := c.Bind(&dto); err != nil {
			return err
		}

		err := sv.Struct(&dto)

		switch err.(type) {
		case *validator.InvalidValidationError:
			return err
		case validator.ValidationErrors:
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		_, err = userService.Register(
			c.Request().Context(),
			&user.TinyUser{
				Name:     dto.Name,
				Email:    dto.Email,
				Password: dto.Password,
			},
		)
		if err != nil {
			if errors.Is(err, user.ErrUserExists) {
				return echo.NewHTTPError(http.StatusConflict, err)
			}

			return err
		}

		return c.NoContent(http.StatusNoContent)
	})

	authRouter.POST("/login", func(c echo.Context) error {
		return c.String(http.StatusOK, "login")
	})
}
