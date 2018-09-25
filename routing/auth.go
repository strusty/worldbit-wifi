package routing

import (
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/auth"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/captcha"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/twilio"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"net/http"
)

type AuthRouter struct {
	authService    auth.Auth
	captchaService captcha.Captcha
	twilioService  twilio.Twilio
}

func NewAuthRouter(
	authService auth.Auth,
	captchaService captcha.Captcha,
	twilioService twilio.Twilio,
) AuthRouter {
	return AuthRouter{
		authService:    authService,
		captchaService: captchaService,
		twilioService:  twilioService,
	}
}

func (router AuthRouter) Register(group *echo.Group) {
	group.POST("/sendCode", router.sendCode)
	group.POST("", router.authenticate)
}

func (router AuthRouter) sendCode(context echo.Context) error {
	request := new(SendCodeRequest)
	if err := context.Bind(request); err != nil {
		return err
	}

	code, err := router.authService.CreateCode(request.PhoneNumber)
	if err != nil {
		return err
	}

	if err := router.twilioService.SendConfirmationCode(request.PhoneNumber, code); err != nil {
		return err
	}

	return context.JSON(http.StatusOK, map[string]bool{
		"success": true,
	})
}

func (router AuthRouter) authenticate(context echo.Context) error {
	request := new(VerifyCodeRequest)
	if err := context.Bind(request); err != nil {
		return err
	}

	captchaVerified, err := router.captchaService.CheckCaptcha(request.Captcha)
	if err != nil {
		return err
	}

	if !captchaVerified {
		return errors.New("Captcha has not passed verification")
	}

	if err := router.authService.VerifyCode(request.ConfirmationCode); err != nil {
		return err
	}

	return context.JSON(http.StatusOK, map[string]bool{
		"success": true,
	})
}
