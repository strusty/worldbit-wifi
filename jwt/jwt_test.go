package jwt

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"
)

func generateContextWithAccessToken(token string) (context echo.Context) {
	e := echo.New()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/test/10",
		strings.NewReader(``),
	)

	request.Header.Set(echo.HeaderAuthorization, middleware.DefaultJWTConfig.AuthScheme+" "+token)

	return e.NewContext(request, recorder)
}

func TestGethyraJWT(t *testing.T) {
	_, err := GenerateJWT(1, false)
	assert.Error(t, err)

	_, err = Middleware(nil)
	assert.Error(t, err)

	assert.Error(t, SetSecret(""))

	assert.NoError(t, SetSecret("secret"))
	assert.Error(t, SetSecret("secret"))

	secret = ""

	assert.NoError(t, SetRandomSecret())
	assert.Error(t, SetRandomSecret())

	token, err := GenerateJWT(10, false)
	if assert.NoError(t, err) {
		middlewareFunc, err := Middleware(nil)
		if assert.NoError(t, err) {
			context := generateContextWithAccessToken(token)

			assert.NoError(t, middlewareFunc(func(context echo.Context) error {
				return nil
			})(context))

			userID, err := GetUserIDFromJWT(context)
			if assert.NoError(t, err) {
				assert.Equal(t, uint(10), userID)
			}

			_, err = GetUserIDFromJWT(generateContextWithAccessToken("adfdsf"))
			assert.Error(t, err)
		}
	}
}
