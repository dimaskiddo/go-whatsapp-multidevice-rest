package internal

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/auth"
	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/router"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/internal/index"
	indexAuth "github.com/dimaskiddo/go-whatsapp-multidevice-rest/internal/index/auth"
	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/internal/whatsapp"
)

func Routes(e *echo.Echo) {
	// Route for Index
	// ---------------------------------------------
	e.GET(router.BaseURL, index.Index)
	e.GET(router.BaseURL+"/auth", index.Auth, auth.BasicAuth())

	// Route for WhatsApp
	// ---------------------------------------------
	authJWTConfig := middleware.JWTConfig{
		Claims:     &indexAuth.AuthJWTClaims{},
		SigningKey: []byte(auth.AuthJWTSecret),
	}

	e.POST(router.BaseURL+"/login", whatsapp.Login, middleware.JWTWithConfig(authJWTConfig))
	e.POST(router.BaseURL+"/send/text", whatsapp.SendText, middleware.JWTWithConfig(authJWTConfig))
	e.POST(router.BaseURL+"/send/location", whatsapp.SendLocation, middleware.JWTWithConfig(authJWTConfig))
	e.POST(router.BaseURL+"/send/document", whatsapp.SendDocument, middleware.JWTWithConfig(authJWTConfig))
	e.POST(router.BaseURL+"/send/audio", whatsapp.SendAudio, middleware.JWTWithConfig(authJWTConfig))
	e.POST(router.BaseURL+"/send/image", whatsapp.SendImage, middleware.JWTWithConfig(authJWTConfig))
	e.POST(router.BaseURL+"/send/video", whatsapp.SendVideo, middleware.JWTWithConfig(authJWTConfig))
}
