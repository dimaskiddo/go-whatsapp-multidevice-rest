package auth

import (
	"encoding/base64"
	"io/ioutil"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/router"
)

// BasicAuth Function as Midleware for Basic Authorization
func BasicAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Parse HTTP Header Authorization
			authHeader := strings.SplitN(c.Request().Header.Get("Authorization"), " ", 2)

			// Check HTTP Header Authorization Section
			// Authorization Section Length Should Be 2
			// The First Authorization Section Should Be "Basic"
			if len(authHeader) != 2 || authHeader[0] != "Basic" {
				return router.ResponseAuthenticate(c)
			}

			// The Second Authorization Section Should Be The Credentials Payload
			// But We Should Decode it First From Base64 Encoding
			authPayload, err := base64.StdEncoding.DecodeString(authHeader[1])
			if err != nil {
				return router.ResponseInternalError(c, "")
			}

			// Split Decoded Authorization Payload Into Username and Password Credentials
			authCredentials := strings.SplitN(string(authPayload), ":", 2)

			// Check Credentials Section
			// It Should Have 2 Section, Username and Password
			if len(authCredentials) != 2 {
				return router.ResponseBadRequest(c, "")
			}

			// Validate Authentication Password
			if authCredentials[1] != AuthBasicPassword {
				return router.ResponseBadRequest(c, "Invalid Authentication")
			}

			// Make Credentials to JSON Format
			authInformation := `{"username": "` + authCredentials[0] + `"}`

			// Rewrite Body Content With Credentials in JSON Format
			c.Request().Header.Set("Content-Type", "application/json")
			c.Request().Body = ioutil.NopCloser(strings.NewReader(authInformation))

			// Call Next Handler Function With Current Request
			return next(c)
		}
	}
}
