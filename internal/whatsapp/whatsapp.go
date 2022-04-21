package whatsapp

import (
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/router"
	pkgWhatsApp "github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/whatsapp"

	typAuth "github.com/dimaskiddo/go-whatsapp-multidevice-rest/internal/auth/types"
	typWhatsApp "github.com/dimaskiddo/go-whatsapp-multidevice-rest/internal/whatsapp/types"
)

func jwtPayload(c echo.Context) typAuth.AuthJWTClaimsPayload {
	jwtToken := c.Get("user").(*jwt.Token)
	jwtClaims := jwtToken.Claims.(*typAuth.AuthJWTClaims)

	return jwtClaims.Data
}

func Login(c echo.Context) error {
	jid := jwtPayload(c).JID

	var reqLogin typWhatsApp.RequestLogin
	reqLogin.Output = strings.TrimSpace(c.FormValue("output"))

	if len(reqLogin.Output) == 0 {
		reqLogin.Output = "html"
	}

	err := pkgWhatsApp.WhatAppConnect(jid)
	if err != nil {
		return router.ResponseInternalError(c, err.Error())
	}

	qrCodeImage, qrCodeTimeout, err := pkgWhatsApp.WhatsAppLogin(jid)
	if err != nil {
		return router.ResponseInternalError(c, err.Error())
	}

	if qrCodeImage == "WhatsApp Client is Reconnected" {
		return router.ResponseSuccess(c, qrCodeImage)
	}

	var resLogin typWhatsApp.ResponseLogin
	resLogin.QRCode = qrCodeImage
	resLogin.Timeout = qrCodeTimeout

	if reqLogin.Output == "html" {
		htmlContent := `
    <html>
      <head>
        <title>WhatsApp Multi-Device Login</title>
      </head>
      <body>
        <img src="` + resLogin.QRCode + `" />
        <p>
          <b>QR Code Scan</b>
          <br/>
          Timeout in ` + strconv.Itoa(resLogin.Timeout) + ` Second(s)
        </p>
      </body>
    </html>`

		return router.ResponseSuccessWithHTML(c, htmlContent)
	}

	return router.ResponseSuccessWithData(c, "Successfully Generated QR Code", resLogin)
}

func Logout(c echo.Context) error {
	jid := jwtPayload(c).JID

	err := pkgWhatsApp.WhatsAppLogout(jid)
	if err != nil {
		return router.ResponseInternalError(c, err.Error())
	}

	return router.ResponseSuccess(c, "Successfully Logged Out")
}

func SendText(c echo.Context) error {
	jid := jwtPayload(c).JID

	var reqSendMessage typWhatsApp.RequestSendMessage
	reqSendMessage.RJID = strings.TrimSpace(c.FormValue("msisdn"))
	reqSendMessage.Message = strings.TrimSpace(c.FormValue("message"))

	if len(reqSendMessage.RJID) == 0 {
		return router.ResponseBadRequest(c, "Missing Form Value MSISDN")
	}

	err := pkgWhatsApp.WhatsAppSendText(jid, reqSendMessage.RJID, reqSendMessage.Message)
	if err != nil {
		return router.ResponseInternalError(c, err.Error())
	}

	return router.ResponseSuccess(c, "Successfully Send Text Message")
}

/*
  TODO: Send Media
*/

/*
func SendLocation(c echo.Context) error {
	return router.ResponseSuccess(c, "Successfully Send Location Message")
}

func SendDocument(c echo.Context) error {
	return router.ResponseSuccess(c, "Successfully Send Document Message")
}

func SendImage(c echo.Context) error {
	return router.ResponseSuccess(c, "Successfully Send Image Message")
}

func SendAudio(c echo.Context) error {
	return router.ResponseSuccess(c, "Successfully Send Audio Message")
}

func SendVideo(c echo.Context) error {
	return router.ResponseSuccess(c, "Successfully Send Video Message")
}
*/
