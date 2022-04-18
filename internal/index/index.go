package index

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/auth"
	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/router"

	indexAuth "github.com/dimaskiddo/go-whatsapp-multidevice-rest/internal/index/auth"
	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/internal/index/model"
)

// Index
func Index(c echo.Context) error {
	return router.ResponseSuccess(c, "Go WhatsApp Multi-Device REST is running")
}

// Auth
func Auth(c echo.Context) error {
	var reqAuthBasicInfo model.ReqAuthBasicInfo
	var resAuthJWTData model.ResAuthJWTData

	// Parse Basic Auth Information from Rewrited Body Request
	// By Basic Auth Middleware
	_ = json.NewDecoder(c.Request().Body).Decode(&reqAuthBasicInfo)

	// Create JWT Claims
	jwtClaims := &indexAuth.AuthJWTClaims{
		indexAuth.AuthJWTClaimsPayload{
			MSISDN: reqAuthBasicInfo.Username,
		},
		jwt.StandardClaims{
			Issuer:    "go-whatsapp-multidevice-rest",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(auth.AuthJWTExpiredHour)).Unix(),
		},
	}

	// Create JWT Token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	// Generate Encoded JWT Token
	jwtTokenEncoded, err := jwtToken.SignedString([]byte(auth.AuthJWTSecret))
	if err != nil {
		return router.ResponseInternalError(c, "")
	}

	// Set Encoded JWT Token as Response Data
	resAuthJWTData.Token = jwtTokenEncoded

	// Return JWT Token in JSON Response
	return router.ResponseSuccessWithData(c, "", resAuthJWTData)
}
