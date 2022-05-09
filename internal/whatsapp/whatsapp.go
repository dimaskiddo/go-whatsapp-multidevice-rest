package whatsapp

import (
	"bytes"
	"image/png"
	"io"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/image/webp"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/env"
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

// Login
// @Summary     Generate QR Code for WhatsApp Multi-Device Login
// @Description Get QR Code for WhatsApp Multi-Device Login
// @Tags        WhatsApp
// @Accept      */*
// @Produce     json
// @Produce     html
// @Param       output    query  string  false  "Change Output Format in HTML or JSON"  Enums(html, json)  default(html)
// @Success     200
// @Security    BearerAuth
// @Router      /api/v1/whatsapp/login [post]
func Login(c echo.Context) error {
	var err error
	jid := jwtPayload(c).JID

	var reqLogin typWhatsApp.RequestLogin
	reqLogin.Output = strings.TrimSpace(c.FormValue("output"))

	if len(reqLogin.Output) == 0 {
		reqLogin.Output = "html"
	}

	// Initialize WhatsApp Client
	pkgWhatsApp.WhatsAppInitClient(nil, jid)

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
        <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no" />
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

// Logout
// @Summary     Logout Device from WhatsApp Multi-Device
// @Description Make Device Logout from WhatsApp Multi-Device
// @Tags        WhatsApp
// @Accept      */*
// @Produce     json
// @Success     200
// @Security    BearerAuth
// @Router      /api/v1/whatsapp/logout [post]
func Logout(c echo.Context) error {
	var err error
	jid := jwtPayload(c).JID

	err = pkgWhatsApp.WhatsAppLogout(jid)
	if err != nil {
		return router.ResponseInternalError(c, err.Error())
	}

	return router.ResponseSuccess(c, "Successfully Logged Out")
}

// SendText
// @Summary     Logout Device from WhatsApp Multi-Device
// @Description Make Device Logout from WhatsApp Multi-Device
// @Tags        WhatsApp
// @Accept      */*
// @Produce     json
// @Param       msisdn    query  string  true  "Destination Phone Number"
// @Param       message   query  string  true  "Text Message Content"
// @Success     200
// @Security    BearerAuth
// @Router      /api/v1/whatsapp/send/text [post]
func SendText(c echo.Context) error {
	var err error
	jid := jwtPayload(c).JID

	var reqSendMessage typWhatsApp.RequestSendMessage
	reqSendMessage.RJID = strings.TrimSpace(c.FormValue("msisdn"))
	reqSendMessage.Message = strings.TrimSpace(c.FormValue("message"))

	if len(reqSendMessage.RJID) == 0 {
		return router.ResponseBadRequest(c, "Missing Form Value MSISDN")
	}

	var resSendMessage typWhatsApp.ResponseSendMessage
	resSendMessage.MsgID, err = pkgWhatsApp.WhatsAppSendText(jid, reqSendMessage.RJID, reqSendMessage.Message)
	if err != nil {
		return router.ResponseInternalError(c, err.Error())
	}

	return router.ResponseSuccessWithData(c, "Successfully Send Text Message", resSendMessage)
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

	var resSendMessage typWhatsApp.ResponseSendMessage
	resSendMessage.MsgID, err = pkgWhatsApp.WhatsAppSendLocation(jid, reqSendLocation.RJID, reqSendLocation.Latitude, reqSendLocation.Longitude)
	if err != nil {
		return router.ResponseInternalError(c, err.Error())
	}

	return router.ResponseSuccessWithData(c, "Successfully Send Location Message", resSendMessage)
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

	// Issue #7 Old Version Client Cannot Render WebP Format
	// If Media Type is "image" and MIME Type is "image/webp"
	// Then Convert it as PNG
	var fileBytes []byte

	isConvertMediaImageWebP, err := env.GetEnvBool("WHATSAPP_MEDIA_IMAGE_CONVERT_WEBP")
	if err != nil {
		isConvertMediaImageWebP = false
	}

	if mediaType == "image" && fileType == "image/webp" && isConvertMediaImageWebP {
		// Decode WebP Image
		fileWebP, err := webp.Decode(fileStream)
		if err != nil {
			return router.ResponseInternalError(c, "Error Decoding Image WebP Format")
		}

		// Encode to PNG Image
		filePNG := new(bytes.Buffer)
		err = png.Encode(filePNG, fileWebP)
		if err != nil {
			return router.ResponseInternalError(c, "Error Encoding Image PNG Format")
		}

		// Set File Stream Bytes and File Type
		// To New Encoded PNG Image and File Type to "image/png"
		fileBytes = filePNG.Bytes()
		fileType = "image/png"
	} else {
		// Convert File Stream in to Bytes
		// Since WhatsApp Proto for Media is only Accepting Bytes format
		fileBytes, err = convertFileToBytes(fileStream)
		if err != nil {
			return router.ResponseInternalError(c, err.Error())
		}
	}

	// Send Media Message Based on Media Type
	var resSendMessage typWhatsApp.ResponseSendMessage
	switch mediaType {
	case "document":
		resSendMessage.MsgID, err = pkgWhatsApp.WhatsAppSendDocument(jid, reqSendMessage.RJID, fileBytes, fileType, reqSendMessage.Message)

	case "image":
		resSendMessage.MsgID, err = pkgWhatsApp.WhatsAppSendImage(jid, reqSendMessage.RJID, fileBytes, fileType, reqSendMessage.Message)

	case "audio":
		resSendMessage.MsgID, err = pkgWhatsApp.WhatsAppSendAudio(jid, reqSendMessage.RJID, fileBytes, fileType)

	case "video":
		resSendMessage.MsgID, err = pkgWhatsApp.WhatsAppSendVideo(jid, reqSendMessage.RJID, fileBytes, fileType, reqSendMessage.Message)
	}

	// Return Internal Server Error
	// When Detected There are Some Errors While Sending The Media Message
	if err != nil {
		return router.ResponseInternalError(c, err.Error())
	}

	return router.ResponseSuccessWithData(c, "Successfully Send Media Message", resSendMessage)
}
