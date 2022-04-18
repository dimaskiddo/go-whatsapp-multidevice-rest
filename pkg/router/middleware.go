package router

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func HttpRealIP() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if XForwardedFor := c.Request().Header.Get(http.CanonicalHeaderKey("X-Forwarded-For")); XForwardedFor != "" {
				dataIndex := strings.Index(XForwardedFor, ", ")
				if dataIndex == -1 {
					dataIndex = len(XForwardedFor)
				}

				c.Request().RemoteAddr = XForwardedFor[:dataIndex]
			} else if XRealIP := c.Request().Header.Get(http.CanonicalHeaderKey("X-Real-IP")); XRealIP != "" {
				c.Request().RemoteAddr = XRealIP
			}

			return next(c)
		}
	}
}
