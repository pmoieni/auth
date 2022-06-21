package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pmoieni/auth/providers/direct"
)

func ResetPassword(c echo.Context) (err error) {
	passwordResetInfo := direct.PasswordResetReq{}
	if err = c.Bind(&passwordResetInfo); err != nil {
		return
	}

	err = direct.ResetPassword(&passwordResetInfo)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "success")
}

func SendPasswordResetToken(c echo.Context) (err error) {
	passwordResetInfo := direct.PasswordResetTokenReq{}
	if err = c.Bind(&passwordResetInfo); err != nil {
		return
	}

	err = passwordResetInfo.SendPasswordResetToken()
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, "You will receive a code in your email if your email is registered.")
}
