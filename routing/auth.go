package routing

import (
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type AuthRouter struct {
}

func NewAuthRouter() AuthRouter {
	return AuthRouter{}
}

func (router AuthRouter) Register(group *echo.Group) {
	group.POST("/sendCode", router.sendCode)
	group.POST("", router.authenticate)
}

func (router AuthRouter) sendCode(context echo.Context) error {
	return errors.New("not implemented")
}

func (router AuthRouter) authenticate(context echo.Context) error {
	return errors.New("not implemented")
}
