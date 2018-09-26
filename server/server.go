package server

import (
	"git.sfxdx.ru/crystalline/wi-fi-backend/routing"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/auth"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/captcha"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/cloudtrax"
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
	cloudtraxService cloudtrax.Cloudtrax,
) Server {
	server := Server{
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
		cloudtraxService,
		twilioService,
	)
	paymentRouter.Register(server.Group("/crypto"))

	return server
}
