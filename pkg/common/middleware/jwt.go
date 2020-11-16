package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/milad-abbasi/oauth-server/pkg/common"
	"github.com/milad-abbasi/oauth-server/pkg/token"
	"github.com/milad-abbasi/oauth-server/pkg/user"
	"gopkg.in/square/go-jose.v2/jwt"
)

func JWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authToken := c.Request().Header.Get("Authorization")
		if len(authToken) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "missing or malformed token")
		}

		var publicClaims jwt.Claims
		ok, err := token.ValidateSignedJwt(
			token.ExtractBearerJwt(authToken),
			&token.Expectation{Secret: common.MustGetEnv("TOKEN_SECRET"), Time: time.Now()},
			&publicClaims,
		)
		if !ok || err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing or malformed token")
		}

		c.Set("user", &user.Identity{ID: publicClaims.ID})

		return next(c)
	}
}
