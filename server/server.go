package server

import (
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/strusty/worldbit-wifi/jwt"
	"github.com/strusty/worldbit-wifi/routing"
	"github.com/strusty/worldbit-wifi/services/admins"
	"github.com/strusty/worldbit-wifi/services/auth"
	"github.com/strusty/worldbit-wifi/services/captcha"
	"github.com/strusty/worldbit-wifi/services/paypal"
	"github.com/strusty/worldbit-wifi/services/pricing_plans"
	"github.com/strusty/worldbit-wifi/services/radius"
	"github.com/strusty/worldbit-wifi/services/twilio"
	"github.com/strusty/worldbit-wifi/services/worldbit"
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
	paypalService paypal.PayPal,
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

	cryptoRouter := routing.NewCryptoRouter(
		worldbitService,
		radiusService,
		twilioService,
		pricingPlanService,
	)
	cryptoRouter.Register(server.Group("/crypto"))

	paypalRouter := routing.NewPayPalRouter(
		paypalService,
		radiusService,
		twilioService,
		pricingPlanService,
	)
	paypalRouter.Register(server.Group("/paypal"))

	pricingPlansRouter := routing.NewPricingPlansRouter(pricingPlanService)
	pricingPlansRouter.Register(server.Group("/plans"))

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
