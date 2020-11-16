package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/milad-abbasi/oauth-server/pkg/common"
	"go.uber.org/zap"
)

type Controller struct {
	l *zap.Logger
	r *echo.Echo
	v *common.Validator
	s *Service
}

func NewController(
	logger *zap.Logger,
	router *echo.Echo,
	validator *common.Validator,
	service *Service,
) *Controller {
	return &Controller{
		l: logger.Named("UserController"),
		r: router,
		v: validator,
		s: service,
	}
}

func (con *Controller) RegisterRoutes() {
	router := con.r.Group("/user")
	router.GET("/me", con.GetMe)
}

func (con *Controller) GetMe(c echo.Context) error {
	return c.String(http.StatusOK, "hi user")
}
