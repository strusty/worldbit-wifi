package jwt

import (
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/strusty/worldbit-wifi/random"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var secret = ""

func SetSecret(s string) error {
	if s == "" {
		return errors.New("Secret cannot be empty")
	}
	if secret != "" {
		return errors.New("Secret is already set")
	}
	secret = s

	return nil
}

func SetRandomSecret() error {
	return SetSecret(random.String(256))
}

func Middleware(skipper middleware.Skipper) (echo.MiddlewareFunc, error) {
	if secret == "" {
		return nil, errors.New("Secret is not set")
	}

	return middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper:    skipper,
		SigningKey: []byte(secret),
	}), nil
}

func GenerateJWT(userID string) (string, error) {
	if secret == "" {
		return "", errors.New("Secret is not set")
	}
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["sub"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	return token.SignedString([]byte(secret))
}

func GetUserIDFromJWT(context echo.Context) (string, error) {
	if context.Get("user") == nil {
		return "", errors.New("Invalid token")
	}
	token := context.Get("user").(*jwt.Token)

	if token != nil {
		claims := token.Claims.(jwt.MapClaims)
		if claims["sub"] == nil {
			return "", errors.New("Invalid token")
		}

		userID := claims["sub"].(string)

		return userID, nil
	} else {
		return "", errors.New("Invalid token")
	}
}
