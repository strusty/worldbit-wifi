package server

import (
	"git.sfxdx.ru/crystalline/wi-fi-backend/routing"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Server struct {
	*echo.Echo
}

func New() Server {
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

	authRouter := routing.NewAuthRouter()
	authRouter.Register(server.Group("/auth"))

	return server
}
