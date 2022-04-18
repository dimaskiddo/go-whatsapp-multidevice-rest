package whatsapp

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/router"

	indexAuth "github.com/dimaskiddo/go-whatsapp-multidevice-rest/internal/index/auth"
)

func getJWTPayload(c echo.Context) indexAuth.AuthJWTClaimsPayload {
	jwtToken := c.Get("user").(*jwt.Token)
	jwtClaims := jwtToken.Claims.(*indexAuth.AuthJWTClaims)

	return jwtClaims.Data
}

func Login(c echo.Context) error {
	return router.ResponseSuccess(c, "")
}

func SendText(c echo.Context) error {
	return router.ResponseSuccess(c, "")
}

func SendLocation(c echo.Context) error {
	return router.ResponseSuccess(c, "")
}

func SendDocument(c echo.Context) error {
	return router.ResponseSuccess(c, "")
}

func SendImage(c echo.Context) error {
	return router.ResponseSuccess(c, "")
}

func SendAudio(c echo.Context) error {
	return router.ResponseSuccess(c, "")
}

func SendVideo(c echo.Context) error {
	return router.ResponseSuccess(c, "")
}
