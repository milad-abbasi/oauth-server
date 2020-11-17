package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/milad-abbasi/oauth-server/pkg/common/config"
	"github.com/milad-abbasi/oauth-server/pkg/common/middleware"
	"github.com/milad-abbasi/oauth-server/pkg/common/validator"
	"go.uber.org/zap"
)

type Controller struct {
	l  *zap.Logger
	r  *echo.Echo
	sv *validator.StructValidator
	s  *Service
}

func NewController(
	logger *zap.Logger,
	router *echo.Echo,
	validator *validator.StructValidator,
	service *Service,
) *Controller {
	return &Controller{
		l:  logger.Named("UserController"),
		r:  router,
		sv: validator,
		s:  service,
	}
}

func (con *Controller) RegisterRoutes() {
	jwtMiddleware := middleware.NewJwtMiddleware(config.MustGetEnv("TOKEN_SECRET"), &Identity{})
	router := con.r.Group("/user", jwtMiddleware.Guard)

	router.GET("/me", con.GetMe)
}

func (con *Controller) GetMe(c echo.Context) error {
	user := c.Get("user").(*Identity)
	u, err := con.s.UserRepo.FindOne(c.Request().Context(), &User{ID: user.ID})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Info{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	})
}
