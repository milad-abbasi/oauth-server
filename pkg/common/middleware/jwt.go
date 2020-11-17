package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/milad-abbasi/oauth-server/pkg/common/token"
	"gopkg.in/square/go-jose.v2/jwt"
)

type JwtMiddleware struct {
	Secret        string
	PrivateClaims interface{}
}

type DefaultPrivateClaims struct {
	Subject string `json:"sub"`
}

func NewJwtMiddleware(secret string, privateClaims interface{}) *JwtMiddleware {
	if privateClaims == nil {
		privateClaims = &DefaultPrivateClaims{}
	}

	return &JwtMiddleware{
		Secret:        secret,
		PrivateClaims: privateClaims,
	}
}

func (jm *JwtMiddleware) Guard(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authToken := c.Request().Header.Get("Authorization")
		if len(authToken) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "missing or malformed token")
		}

		var publicClaims jwt.Claims
		ok, err := token.ValidateSignedJwt(
			token.ExtractBearerJwt(authToken),
			&token.Expectation{Secret: jm.Secret, Time: time.Now()},
			&publicClaims,
			jm.PrivateClaims,
		)
		if !ok || err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing or malformed token")
		}

		c.Set("user", jm.PrivateClaims)

		return next(c)
	}
}
