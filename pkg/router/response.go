package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/log"
)

type ResSuccess struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResSuccessWithData struct {
	Status  bool        `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResError struct {
	Status bool   `json:"status"`
	Code   int    `json:"code"`
	Error  string `json:"error"`
}

func logSuccess(c echo.Context, code int, message string) {
	statusMessage := http.StatusText(code)

	if statusMessage == message || c.Request().RequestURI == BaseURL {
		log.Print(c).Info(fmt.Sprintf("%d %v", code, statusMessage))
	} else {
		log.Print(c).Info(fmt.Sprintf("%d %v", code, message))
	}
}

func logError(c echo.Context, code int, message string) {
	statusMessage := http.StatusText(code)

	if statusMessage == message {
		log.Print(c).Error(fmt.Sprintf("%d %v", code, statusMessage))
	} else {
		log.Print(c).Error(fmt.Sprintf("%d %v", code, message))
	}
}

func ResponseSuccess(c echo.Context, message string) error {
	var response ResSuccess

	response.Status = true
	response.Code = http.StatusOK

	if strings.TrimSpace(message) == "" {
		message = http.StatusText(response.Code)
	}
	response.Message = message

	logSuccess(c, response.Code, response.Message)
	return c.JSON(response.Code, response)
}

func ResponseSuccessWithData(c echo.Context, message string, data interface{}) error {
	var response ResSuccessWithData

	response.Status = true
	response.Code = http.StatusOK

	if strings.TrimSpace(message) == "" {
		message = http.StatusText(response.Code)
	}
	response.Message = message
	response.Data = data

	logSuccess(c, response.Code, response.Message)
	return c.JSON(response.Code, response)
}

func ResponseSuccessWithHTML(c echo.Context, html string) error {
	logSuccess(c, http.StatusOK, http.StatusText(http.StatusOK))
	return c.HTML(http.StatusOK, html)
}

func ResponseCreated(c echo.Context, message string) error {
	var response ResSuccess

	response.Status = true
	response.Code = http.StatusCreated

	if strings.TrimSpace(message) == "" {
		message = http.StatusText(response.Code)
	}
	response.Message = message

	logSuccess(c, response.Code, response.Message)
	return c.JSON(response.Code, response)
}

func ResponseNoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func ResponseNotFound(c echo.Context, message string) error {
	var response ResError

	response.Status = false
	response.Code = http.StatusNotFound

	if strings.TrimSpace(message) == "" {
		message = http.StatusText(response.Code)
	}
	response.Error = message

	logError(c, response.Code, response.Error)
	return c.JSON(response.Code, response)
}

func ResponseAuthenticate(c echo.Context) error {
	c.Response().Header().Set("WWW-Authenticate", `Basic realm="Authentication Required"`)
	return ResponseUnauthorized(c, "")
}

func ResponseUnauthorized(c echo.Context, message string) error {
	var response ResError

	response.Status = false
	response.Code = http.StatusUnauthorized

	if strings.TrimSpace(message) == "" {
		message = http.StatusText(response.Code)
	}
	response.Error = message

	logError(c, response.Code, response.Error)
	return c.JSON(response.Code, response)
}

func ResponseBadRequest(c echo.Context, message string) error {
	var response ResError

	response.Status = false
	response.Code = http.StatusBadRequest

	if strings.TrimSpace(message) == "" {
		message = http.StatusText(response.Code)
	}
	response.Error = message

	logError(c, response.Code, response.Error)
	return c.JSON(response.Code, response)
}

func ResponseInternalError(c echo.Context, message string) error {
	var response ResError

	response.Status = false
	response.Code = http.StatusInternalServerError

	if strings.TrimSpace(message) == "" {
		message = http.StatusText(response.Code)
	}
	response.Error = message

	logError(c, response.Code, response.Error)
	return c.JSON(response.Code, response)
}

func ResponseBadGateway(c echo.Context, message string) error {
	var response ResError

	response.Status = false
	response.Code = http.StatusBadGateway

	if strings.TrimSpace(message) == "" {
		message = http.StatusText(response.Code)
	}
	response.Error = message

	logError(c, response.Code, response.Error)
	return c.JSON(response.Code, response)
}
