package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pmoieni/auth/models"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	e := models.ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}

	if httpError, ok := err.(*models.ErrorResponse); ok {
		e.Status = httpError.Status
		e.Message = httpError.Message
	} else {
		c.Logger().Error(err)
	}

	c.JSON(e.Status, e.Message)
}
