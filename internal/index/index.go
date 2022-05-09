package index

import (
	"github.com/labstack/echo/v4"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/router"
)

// Index
// @Summary Show The Status of The Server
// @Description Get The Server Status
// @Tags Root
// @Accept */*
// @Produce json
// @Success 200
// @Router /api/v1/whatsapp [get]
func Index(c echo.Context) error {
	return router.ResponseSuccess(c, "Go WhatsApp Multi-Device REST is running")
}
