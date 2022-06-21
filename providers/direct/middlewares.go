package direct

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pmoieni/auth/models"
)

const (
	authorizationHeader = "Authorization"
)

type ctxKey string // use a type of type string to make sure ctx key is coming from middleware
var userInfoCtxKey = ctxKey("user-email")

func CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		at := c.Request().Header.Get(authorizationHeader)
		userInfo, err := parseAccessTokenWithValidate(at)
		if err != nil {
			return &models.ErrorResponse{Status: http.StatusUnauthorized, Message: http.StatusText(http.StatusUnauthorized)}
		}

		privateClaims := userInfo.PrivateClaims()
		userEmail, ok := privateClaims["email"].(string)
		if !ok {
			return &models.ErrorResponse{Status: http.StatusUnauthorized, Message: http.StatusText(http.StatusUnauthorized)}
		}

		c.Set(string(userInfoCtxKey), userEmail)

		return next(c)
	}
}
