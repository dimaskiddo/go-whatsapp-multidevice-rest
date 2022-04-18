package router

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func HttpErrorHandler(err error, c echo.Context) {
	report, _ := err.(*echo.HTTPError)

	response := &ResError{
		Status: false,
		Code:   report.Code,
		Error:  fmt.Sprintf("%v", report.Message),
	}

	logError(c, response.Code, response.Error)
	c.JSON(response.Code, response)
}
