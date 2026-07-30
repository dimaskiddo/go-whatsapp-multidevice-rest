package router

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HttpErrorHandler(err error, c echo.Context) {
	report, ok := err.(*echo.HTTPError)
	if !ok {
		// Handle non-HTTPError errors (e.g., panics recovered by recover middleware)
		c.JSON(http.StatusInternalServerError, &ResError{
			Status: false,
			Code:   http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	response := &ResError{
		Status: false,
		Code:   report.Code,
		Error:  fmt.Sprintf("%v", report.Message),
	}

	logError(c, response.Code, response.Error)
	c.JSON(response.Code, response)
}
