package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pmoieni/auth/providers/direct"
)

const (
	authorizationHeader    = "Authorization"
	refreshTokenCookieName = "BL_Direct_RT"
	refreshTokenCookiePath = "/api/v1/auth/direct"
)

type loginRes struct {
	IDToken     string `json:"id_token"`
	AccessToken string `json:"access_token"`
}

type refreshTokenRes struct {
	AccessToken string `json:"access_token"`
}

func Register(c echo.Context) (err error) {
	userInfo := direct.UserRegisterInfo{}
	if err = c.Bind(&userInfo); err != nil {
		return
	}

	err = userInfo.Register()
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, "success")
}

func Login(c echo.Context) (err error) {
	// get user credentials from request and bind it to UserLoginCreds type
	userCreds := direct.UserLoginCreds{}
	if err = c.Bind(&userCreds); err != nil {
		return
	}

	// Login the user with provided credentials
	tokens, err := userCreds.Login()
	if err != nil {
		return
	}

	rtCookie := http.Cookie{
		// Domain:   "example.com",
		Path:     refreshTokenCookiePath,
		Name:     refreshTokenCookieName,
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().UTC().Add(direct.RefreshTokenExpiry),
	}

	res := loginRes{
		IDToken:     tokens.IDToken,
		AccessToken: tokens.AccessToken,
	}

	c.SetCookie(&rtCookie)
	return c.JSON(http.StatusOK, res)
}

func Logout(c echo.Context) (err error) {
	// remove refresh token cookie
	cookie := &http.Cookie{
		// Domain:   "example.com",
		Path:     refreshTokenCookiePath,
		Name:     refreshTokenCookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
	}

	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, "success")
}

func RefreshToken(c echo.Context) (err error) {
	rtCookie, err := c.Cookie(refreshTokenCookieName)
	if err != nil {
		return
	}

	tokens, err := direct.RefreshToken(rtCookie.Value)
	if err != nil {
		return
	}

	newRTCookie := http.Cookie{
		// Domain:   "example.com",
		Path:     refreshTokenCookiePath,
		Name:     refreshTokenCookieName,
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().UTC().Add(direct.RefreshTokenExpiry),
	}

	res := refreshTokenRes{
		AccessToken: tokens.AccessToken,
	}

	c.SetCookie(&newRTCookie)
	return c.JSON(http.StatusOK, res)
}
