package whatsapp

import (
	"bytes"
	"io"
	"mime/multipart"
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

func convertFileToBytes(file multipart.File) ([]byte, error) {
	// Create Empty Buffer
	buffer := bytes.NewBuffer(nil)

	// Copy File Stream to Buffer
	_, err := io.Copy(buffer, file)
	if err != nil {
		return bytes.NewBuffer(nil).Bytes(), err
	}

	return buffer.Bytes(), nil
}

func sendMedia(c echo.Context, mediaType string) error {
	var err error
	jid := jwtPayload(c).JID

	var reqSendMessage typWhatsApp.RequestSendMessage
	reqSendMessage.RJID = strings.TrimSpace(c.FormValue("msisdn"))

	// Read Uploaded File Based on Send Media Type
	var fileStream multipart.File
	var fileHeader *multipart.FileHeader

	switch mediaType {
	case "document":
		fileStream, fileHeader, err = c.Request().FormFile("document")
		reqSendMessage.Message = fileHeader.Filename

	case "image":
		fileStream, fileHeader, err = c.Request().FormFile("image")
		reqSendMessage.Message = strings.TrimSpace(c.FormValue("caption"))

	case "audio":
		fileStream, fileHeader, err = c.Request().FormFile("audio")

	case "video":
		fileStream, fileHeader, err = c.Request().FormFile("video")
		reqSendMessage.Message = strings.TrimSpace(c.FormValue("caption"))
	}

	// Don't Forget to Close The File Stream
	defer fileStream.Close()

	// Get Uploaded File MIME Type
	fileType := fileHeader.Header.Get("Content-Type")

	// If There are Some Errors While Opeening The File Stream
	// Return Bad Request with Original Error Message
	if err != nil {
		return router.ResponseBadRequest(c, err.Error())
	}

	// Make Sure RJID is Filled
	if len(reqSendMessage.RJID) == 0 {
		return router.ResponseBadRequest(c, "Missing Form Value MSISDN")
	}

	// Convert File Stream in to Bytes
	// Since WhatsApp Proto for Media is only Accepting Bytes format
	fileBytes, err := convertFileToBytes(fileStream)
	if err != nil {
		return router.ResponseInternalError(c, err.Error())
	}

	// Send Media Message Based on Media Type
	switch mediaType {
	case "document":
		err = pkgWhatsApp.WhatsAppSendDocument(jid, reqSendMessage.RJID, fileBytes, fileType, reqSendMessage.Message)

	case "image":
		err = pkgWhatsApp.WhatsAppSendImage(jid, reqSendMessage.RJID, fileBytes, fileType, reqSendMessage.Message)

	case "audio":
		err = pkgWhatsApp.WhatsAppSendAudio(jid, reqSendMessage.RJID, fileBytes, fileType)

	case "video":
		err = pkgWhatsApp.WhatsAppSendVideo(jid, reqSendMessage.RJID, fileBytes, fileType, reqSendMessage.Message)
	}

	// Return Internal Server Error
	// When Detected There are Some Errors While Sending The Media Message
	if err != nil {
		return router.ResponseInternalError(c, err.Error())
	}

	return router.ResponseSuccess(c, "Successfully Send Media Message")
}

func Login(c echo.Context) error {
	var err error
	jid := jwtPayload(c).JID

	var reqLogin typWhatsApp.RequestLogin
	reqLogin.Output = strings.TrimSpace(c.FormValue("output"))

	if len(reqLogin.Output) == 0 {
		reqLogin.Output = "html"
	}

	// Initialize WhatsApp Client
	err = pkgWhatsApp.WhatsAppInitClient(nil, jid)
	if err != nil {
		return router.ResponseInternalError(c, err.Error())
	}

	// Get WhatsApp QR Code Image
	qrCodeImage, qrCodeTimeout, err := pkgWhatsApp.WhatsAppLogin(jid)
	if err != nil {
		return router.ResponseInternalError(c, err.Error())
	}

	// If Return is Not QR Code But Reconnected
	// Then Return OK With Reconnected Status
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
	var err error
	jid := jwtPayload(c).JID

	err = pkgWhatsApp.WhatsAppLogout(jid)
	if err != nil {
		return router.ResponseInternalError(c, err.Error())
	}

	return router.ResponseSuccess(c, "Successfully Logged Out")
}

func SendText(c echo.Context) error {
	var err error
	jid := jwtPayload(c).JID

	var reqSendMessage typWhatsApp.RequestSendMessage
	reqSendMessage.RJID = strings.TrimSpace(c.FormValue("msisdn"))
	reqSendMessage.Message = strings.TrimSpace(c.FormValue("message"))

	if len(reqSendMessage.RJID) == 0 {
		return router.ResponseBadRequest(c, "Missing Form Value MSISDN")
	}

	err = pkgWhatsApp.WhatsAppSendText(jid, reqSendMessage.RJID, reqSendMessage.Message)
	if err != nil {
		return router.ResponseInternalError(c, err.Error())
	}

	return router.ResponseSuccess(c, "Successfully Send Text Message")
}

func SendLocation(c echo.Context) error {
	var err error
	jid := jwtPayload(c).JID

	var reqSendLocation typWhatsApp.RequestSendLocation
	reqSendLocation.RJID = strings.TrimSpace(c.FormValue("msisdn"))

	reqSendLocation.Latitude, err = strconv.ParseFloat(strings.TrimSpace(c.FormValue("latitude")), 64)
	if err != nil {
		return router.ResponseInternalError(c, "Error While Decoding Latitude to Float64")
	}

	reqSendLocation.Longitude, err = strconv.ParseFloat(strings.TrimSpace(c.FormValue("longitude")), 64)
	if err != nil {
		return router.ResponseInternalError(c, "Error While Decoding Longitude to Float64")
	}

	if len(reqSendLocation.RJID) == 0 {
		return router.ResponseBadRequest(c, "Missing Form Value MSISDN")
	}

	err = pkgWhatsApp.WhatsAppSendLocation(jid, reqSendLocation.RJID, reqSendLocation.Latitude, reqSendLocation.Longitude)
	if err != nil {
		return router.ResponseInternalError(c, err.Error())
	}

	return router.ResponseSuccess(c, "Successfully Send Location Message")
}

func SendDocument(c echo.Context) error {
	return sendMedia(c, "document")
}

func SendImage(c echo.Context) error {
	return sendMedia(c, "image")
}

func SendAudio(c echo.Context) error {
	return sendMedia(c, "audio")
}

func SendVideo(c echo.Context) error {
	return sendMedia(c, "video")
}
