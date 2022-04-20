package index

import (
	"github.com/labstack/echo/v4"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/router"
)

// Index
func Index(c echo.Context) error {
	return router.ResponseSuccess(c, "Go WhatsApp Multi-Device REST is running")
}
