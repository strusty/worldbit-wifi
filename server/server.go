package server

import (
	"strings"

	"git.sfxdx.ru/crystalline/wi-fi-backend/jwt"
	"git.sfxdx.ru/crystalline/wi-fi-backend/routing"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/admins"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/auth"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/captcha"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/pricing_plans"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/radius"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/twilio"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/worldbit"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Server struct {
	*echo.Echo
}

func New(
	authService auth.Auth,
	captchaService captcha.Captcha,
	twilioService twilio.Twilio,
	worldbitService worldbit.Worldbit,
	radiusService radius.Radius,
	pricingPlanService pricing_plans.PricingPlans,
	adminService admins.Admins,
) (*Server, error) {
	server := &Server{
		Echo: echo.New(),
	}

	server.HTTPErrorHandler = func(err error, context echo.Context) {
		defer context.Logger().Error(err)
		if httpErr, ok := err.(*echo.HTTPError); ok {
			context.JSON(200, map[string]interface{}{
				"error": map[string]interface{}{
					"code":    httpErr.Code,
					"message": httpErr.Message,
				},
			})
		} else {
			context.JSON(200, map[string]interface{}{
				"error": map[string]interface{}{
					"code":    500,
					"message": err.Error(),
				},
			})
		}
	}

	// Middleware
	server.Pre(middleware.RemoveTrailingSlash())
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(middleware.CORS())

	authRouter := routing.NewAuthRouter(
		authService,
		captchaService,
		twilioService,
	)
	authRouter.Register(server.Group("/auth"))

	paymentRouter := routing.NewCryptoRouter(
		worldbitService,
		radiusService,
		twilioService,
		pricingPlanService,
	)
	paymentRouter.Register(server.Group("/crypto"))

	adminRouter := routing.NewAdminRouter(adminService, pricingPlanService)
	adminGroup := server.Group("/admin")
	jwtMiddleware, err := jwt.Middleware(func(c echo.Context) bool {
		if strings.HasPrefix(c.Path(), "/admin/login") {
			return true
		}
		return false
	})
	if err != nil {
		return nil, err
	}

	adminGroup.Use(jwtMiddleware)
	adminRouter.Register(adminGroup)

	return server, nil
}
