package captcha

type Captcha interface {
	CheckCaptcha(responseToken string) (bool, error)
}
